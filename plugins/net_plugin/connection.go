package net_plugin

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	. "github.com/eosspark/eos-go/plugins/net_plugin/multi_index"
	"github.com/eosspark/eos-go/plugins/net_plugin/multi_index/peer_block_state"
	"github.com/eosspark/eos-go/plugins/net_plugin/multi_index/transaction_state"
	"io"
	"net"
	"reflect"
	"runtime"
	"time"
)

const (
	tsBufferSize int = 32
)

type Connection struct {
	blkState      *peer_block_state.PeerBlockStateIndex
	trxState      *transaction_state.TransactionStateIndex
	peerRequested *syncState //this peer is requesting info from us
	//socket socketPtr

	socket               net.Conn
	reader               io.Reader
	pendingMessageBuffer [1024 * 1024]byte
	//outstandingReadBytes optional<int>
	blkBuffer []byte

	nodeID             common.NodeIdType
	lastHandshakeRecv  *HandshakeMessage
	lastHandshakeSent  *HandshakeMessage
	sentHandshakeCount uint16 //int16
	connecting         bool
	syncing            bool
	protocolVersion    uint16
	peerAddr           string
	responseExpected   *asio.DeadlineTimer
	//pendingFetch optional<request_message>
	pendingFetch *RequestMessage

	noRetry     GoAwayReason
	forkHead    common.BlockIdType
	forkHeadNum uint32
	//lastReq optional<request_message>
	lastReq *RequestMessage

	// Members set from network data
	org common.TimePoint //originate timestamp
	rec common.TimePoint //receive timestamp
	dst common.TimePoint //destination timestamp
	xmt common.TimePoint // transmit timestamp

	// Computed data
	offset float64 //peer offset double TODO

	ts [tsBufferSize]byte //working buffer for making human readable timestamps

	impl *netPluginIMpl
}

type PeerStatus struct {
	Peer          string
	Connecting    bool
	Syncing       bool
	LastHandshake HandshakeMessage
}

func NewConnectionByEndPoint(endpoint string, impl *netPluginIMpl) *Connection {
	conn := &Connection{
		blkState:           peer_block_state.NewPeerBlockStateIndex(),
		trxState:           transaction_state.NewTransactionStateIndex(),
		peerRequested:      new(syncState),
		lastHandshakeSent:  &HandshakeMessage{},
		lastHandshakeRecv:  &HandshakeMessage{},
		sentHandshakeCount: 0,
		connecting:         false,
		syncing:            false,
		protocolVersion:    0,
		noRetry:            noReason,
		forkHeadNum:        0,
		impl:               impl,
		//forkHead:
		//nodeID:
		//lastReq:
		peerAddr: endpoint,
	}
	netLog.Warn("accepted network connection")
	conn.initialize()

	return conn
}

func NewConnectionByConn(c net.Conn, impl *netPluginIMpl) *Connection {
	conn := &Connection{
		blkState:           peer_block_state.NewPeerBlockStateIndex(),
		trxState:           transaction_state.NewTransactionStateIndex(),
		peerRequested:      new(syncState),
		socket:             c,
		lastHandshakeSent:  &HandshakeMessage{},
		lastHandshakeRecv:  &HandshakeMessage{},
		sentHandshakeCount: 0,
		connecting:         true,
		syncing:            false,
		protocolVersion:    0,
		noRetry:            noReason,
		forkHeadNum:        0,
		impl:               impl,
		//forkHead:
		//nodeID:
		//lastReq:
		//peerAddr:           c.RemoteAddr().String(),

	}
	netLog.Warn("accepted network connection")
	conn.initialize()

	return conn
}

func (c *Connection) initialize() {
	rnd := c.nodeID.Bytes()
	rnd[0] = 0
	c.responseExpected = asio.NewDeadlineTimer(App().GetIoService())
}

func NewPeer(impl *netPluginIMpl, conn net.Conn, reader io.Reader) *Connection {
	return &Connection{
		blkState:      peer_block_state.NewPeerBlockStateIndex(),
		trxState:      transaction_state.NewTransactionStateIndex(),
		peerRequested: new(syncState),

		socket:             conn,
		reader:             reader,
		peerAddr:           conn.RemoteAddr().String(),
		lastHandshakeSent:  &HandshakeMessage{},
		lastHandshakeRecv:  &HandshakeMessage{},
		sentHandshakeCount: 0,
		impl:               impl,
	}
}

func (c *Connection) getStatus() *PeerStatus {
	return &PeerStatus{
		Peer:          c.peerAddr,
		Connecting:    c.connecting,
		Syncing:       c.syncing,
		LastHandshake: *c.lastHandshakeRecv,
	}
}

func (c *Connection) connected() bool {

	return false
}

func (c *Connection) current() bool { //TODO

	return true
}

func (c *Connection) reset() {
	c.peerRequested = nil //TODO
	c.blkState = nil
	c.trxState = nil
}

func (c *Connection) close() { //TODO

}

func (c *Connection) sendHandshake(impl *netPluginIMpl) {
	handshakePopulate(impl, c.lastHandshakeSent)
	c.sentHandshakeCount += 1
	c.lastHandshakeSent.Generation = c.sentHandshakeCount
	//fc_dlog(logger, "Sending handshake generation ${g} to ${ep}",
	//	("g",last_handshake_sent.generation)("ep", peer_name()));

	fmt.Printf("Sending handshake generation %d to %s\n", c.lastHandshakeSent.Generation, c.peerAddr)
	c.write(c.lastHandshakeSent)
}

func handshakePopulate(impl *netPluginIMpl, hello *HandshakeMessage) {
	hello.NetworkVersion = netVersionBase + netVersion
	hello.ChainID = impl.chainID
	hello.NodeID = impl.nodeID
	hello.Key = *impl.getAuthenticationKey()
	hello.Time = common.Now()
	hello.Token = *crypto.Hash256(hello.Time)
	hello.Signature = *impl.signCompact(&hello.Key, &hello.Token)

	// If we couldn't sign, don't send a token.
	if common.Empty(hello.Signature) {
		hello.Token = *crypto.NewSha256Nil()
	}

	hello.P2PAddress = impl.p2PAddress + "-" + hello.NodeID.String()[:7]

	switch runtime.GOOS {
	case "darwin":
		hello.OS = "osx"
	case "linux":
		hello.OS = "linux"
	case "windows":
		hello.OS = "win32"
	default:
		hello.OS = "other"
	}
	hello.Agent = impl.userAgentName

	//controller& cc = my_impl->chain_plug->chain();
	hello.HeadID = common.BlockIdNil()
	//hello.head_num = cc.fork_db_head_block_num();
	//hello.last_irreversible_block_num = cc.last_irreversible_block_num();

	//hello.HeadNum =  cc.fork_db_head_block_num();
	//hello.LastIrreversibleBlockID = hello.LastIrreversibleBlockID

	if hello.LastIrreversibleBlockNum > 0 {
		Try(func() {
			//hello.LastIrreversibleBlockID = cc.get_block_id_for_num(hello.last_irreversible_block_num)
			//hello.LastIrreversibleBlockID =
		}).Catch(func(ex *exception.UnknownBlockException) {
			//ilog("caught unkown_block");
			fmt.Println("caught unkown_block")
			//hello.LastIrreversibleBlockNum =0
		}).End()
	}
	if hello.HeadNum > 0 {
		Try(func() {
			//hello.id = cc.get_block_id_for_num( hello.head_num )
			//hello.id =
		}).Catch(func(ex *exception.UnknownBlockException) {
			hello.HeadNum = 0
		}).End()
	}

}

func (c *Connection) PeerName() string {
	return c.peerAddr
}

//void connection::sync_wait( ) {
//response_expected->expires_from_now( my_impl->resp_expected_period);
//connection_wptr c(shared_from_this());
//response_expected->async_wait( [c]( boost::system::error_code ec){
//connection_ptr conn = c.lock();
//if (!conn) {
//// connection was destroyed before this lambda was delivered
//return;
//}
//
//conn->sync_timeout(ec);
//} );
//}
func (c *Connection) syncWait() {

}

func (c *Connection) cancelWait() {

}

func (c *Connection) fetchWait() {

}

func (c *Connection) cancelSync(reason GoAwayReason) {
	//fc_dlog(logger,"cancel sync reason = ${m}, write queue size ${o} peer ${p}",
	//	("m",reason_str(reason)) ("o", write_queue.size())("p", peer_name()));

	netLog.Debug("cancel sync reason = %s, write queue size %d peer %s", ReasonToString[reason], c.write)

	c.cancelWait()
	//c.flushQueues()
	switch reason {
	case validation, fatalOther:
		c.noRetry = reason
		c.write(&GoAwayMessage{Reason: reason})
	default:
		//fc_dlog(logger, "sending empty request but not calling sync wait on ${p}", ("p",peer_name()))
		netLog.Debug("sending empty request but not calling sync wait on %s", c.peerAddr)
		c.write(&SyncRequestMessage{0, 0})
	}

}

func (c *Connection) txnSendPending(ids []*common.TransactionIdType) {
	for tx := c.impl.localTxns.GetById().Begin(); !tx.IsEnd(); tx.Next() {
		if len(tx.Value().SerializedTxn) > 0 && tx.Value().BlockNum == 0 {
			found := false
			for _, known := range ids {
				if known.Equals(tx.Value().ID) {
					found = true
					break
				}
			}
			if !found {
				//my_impl->local_txns.modify(tx,incr_in_flight);
				c.impl.localTxns.Modify(tx, func(state *NodeTransactionState) {
					exp := state.Expires.SecSinceEpoch()
					state.Expires = common.TimePointSec(exp + 1*60)
					if state.Requests == 0 {
						state.TrueBlock = state.BlockNum
						state.BlockNum = 0
					}
					state.Requests = state.Requests + 1
					if state.Requests == 0 {
						state.BlockNum = state.TrueBlock
					}
				})

				//TODO queueWrite()
				//queue_write(std::make_shared<vector<char>>(tx->serialized_txn),
				//	true,
				//	[tx_id=tx->id](boost::system::error_code ec, std::size_t ) {
				//auto& local_txns = my_impl->local_txns;
				//auto tx = local_txns.get<by_id>().find(tx_id);
				//if (tx != local_txns.end()) {
				//local_txns.modify(tx, decr_in_flight);
				//} else {
				//fc_wlog(logger, "Local pending TX erased before queued_write called callback");
				//}
				//});

			}

		}
	}

}

func (c *Connection) txnSend(ids []*common.TransactionIdType) {
	for _, t := range ids {
		tx := c.impl.localTxns.GetById().Find(*t)
		if !tx.IsEnd() && len(tx.Value().SerializedTxn) > 0 {
			c.impl.localTxns.Modify(tx, func(state *NodeTransactionState) {
				exp := state.Expires.SecSinceEpoch()
				state.Expires = common.TimePointSec(exp + 1*60)
				if state.Requests == 0 {
					state.TrueBlock = state.BlockNum
					state.BlockNum = 0
				}
				state.Requests = state.Requests + 1
				if state.Requests == 0 {
					state.BlockNum = state.TrueBlock
				}
			})
			//queue_write(std::make_shared<vector<char>>(tx->serialized_txn),
			//	true,
			//	[t](boost::system::error_code ec, std::size_t ) {
			//	auto& local_txns = my_impl->local_txns;
			//	auto tx = local_txns.get<by_id>().find(t);
			//	if (tx != local_txns.end()) {
			//		local_txns.modify(tx, decr_in_flight);
			//	} else {
			//		fc_wlog(logger, "Local TX erased before queued_write called callback");
			//	}
			//});
		}
	}
}

func (c *Connection) blkSendBranch() {
	cc := c.impl.ChainPlugin.Chain()
	headNum := cc.ForkDbHeadBlockNum()

	var note NoticeMessage

	note.KnownBlocks.Mode = normal
	note.KnownBlocks.Pending = 0
	fcLog.Debug("head_num = %d", headNum)
	if headNum == 0 {
		c.write(&note) //todo
		return
	}

	var (
		headID        common.BlockIdType
		libID         common.BlockIdType
		remoteHeadID  common.BlockIdType
		remoteHeadNum uint32
	)
	returning := false
	Try(func() {
		if c.lastHandshakeRecv.Generation >= 1 {
			remoteHeadID = c.lastHandshakeRecv.HeadID
			remoteHeadNum = types.NumFromID(&remoteHeadID)
			fcLog.Debug("maybe truncating branch at = %d : %s", remoteHeadNum, remoteHeadID)
		}
		// base our branch off of the last handshake we sent the peer instead of our current
		// LIB which could have moved forward in time as packets were in flight.
		if c.lastHandshakeSent.Generation >= 1 {
			libID = c.lastHandshakeSent.LastIrreversibleBlockID
		} else {
			libID = cc.LastIrreversibleBlockId()
		}
		headID = cc.ForkDbHeadBlockId()
	}).Catch(func(ex *exception.AssertException) {
		log.Error("unable to retrieve block info: %s for %s", ex.What(), c.peerAddr)
		c.write(&note)
		returning = true
	}).Catch(func(ex exception.Exception) {

	}).Catch(func(interface{}) {

	}).End()
	if returning {
		return
	}

	var (
		bStack []*types.SignedBlock
		nullID common.BlockIdType
	)
	breaking := false
	for bid := headID; bid != nullID && bid != libID; {
		Try(func() {
			if remoteHeadID.Equals(bid) {
				breaking = true
				return
			}

			b := cc.FetchBlockById(bid)
			if b != nil {
				bid = b.Previous
				bStack = append(bStack, b)
			} else {
				breaking = true
				return
			}
		}).Catch(func(interface{}) {
			breaking = true
			return
		}).End()

		if breaking {
			break
		}
	}

	count := len(bStack)
	if len(bStack) > 0 {
		if bStack[count-1].Previous.Equals(libID) || bStack[count-1].Previous.Equals(remoteHeadID) {
			for i := len(bStack); i > 0; i-- {
				c.write(&SignedBlockMessage{*bStack[i-1]})
			}
		}
		fcLog.Info("Sent %d blocks on my fork", count)
	} else {
		fcLog.Info("Nothing to send on fork request")
	}

	c.syncing = false
}

func (c *Connection) blkSend(ids []*common.BlockIdType) {
	cc := c.impl.ChainPlugin.Chain()
	var count int
	breaking := false
	for _, blkID := range ids {
		count++
		Try(func() {
			b := cc.FetchBlockById(*blkID)
			if b != nil {
				netLog.Debug("found block for id ar num %d", b.BlockNumber())
				//enqueue(net_message(*b))//TODO
				c.write(&SignedBlockMessage{*b})
			} else {
				netLog.Info("fetch block by id returned null, id %s on block %d of %d for %s",
					blkID, count, len(ids), c.peerAddr)
				breaking = true
			}
		}).Catch(func(ex *exception.AssertException) {
			netLog.Error("caught assert on fetch_block_by_id, %s, id %s on block %d of %d for %s",
				ex.What(), blkID, count, len(ids), c.peerAddr)
			breaking = true
		}).Catch(func(interface{}) {
			netLog.Error("caught others exception fetching block id %s on block %d of %d for %s",
				blkID, count, len(ids), c.peerAddr)
			breaking = true
		}).End()

		if breaking {
			break
		}
	}
}

func (c *Connection) stopSend() {
	c.syncing = false
}

//void enqueue( const net_message &msg, bool trigger_send = true );
//void cancel_sync(go_away_reason);
//void flush_queues();

func (c *Connection) enqueueSyncBlock() bool {
	cc := App().FindPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin).Chain()
	if c.peerRequested == nil {
		return false
	}

	c.peerRequested.last++
	num := c.peerRequested.last
	triggerSend := num == c.peerRequested.startBlock
	if num == c.peerRequested.endBlock {
		c.peerRequested = nil
	}

	result := false
	Try(func() {
		sb := cc.FetchBlockByNumber(num)
		if sb != nil {
			//enqueue( *sb, trigger_send);//todo
			fmt.Println(triggerSend)
			result = true
		}
	}).Catch(func(e interface{}) {
		log.Warn("write loop exception")
	}).End()

	return result
}

func (c *Connection) requestSyncBlocks(start, end uint32) {
	syncRequest := SyncRequestMessage{
		StartBlock: start,
		EndBlock:   end,
	}
	c.write(&syncRequest)
	c.syncWait()
}

//void sync_timeout(boost::system::error_code ec);
//void fetch_timeout(boost::system::error_code ec);
//
//void queue_write(std::shared_ptr<vector<char>> buff,
//bool trigger_send,
//std::function<void(boost::system::error_code, std::size_t)> callback);
//void do_queue_write();

func isValid(msg *HandshakeMessage) bool {
	// Do some basic validation of an incoming handshake_message, so things
	// that really aren't handshake messages can be quickly discarded without
	// affecting state.
	valid := true
	if msg.LastIrreversibleBlockNum > msg.HeadNum {
		netLog.Warn("Handshake message validation: last irreversible block %d is greater than head block %d",
			msg.LastIrreversibleBlockNum, msg.HeadNum)
		valid = false
	}
	if len(msg.P2PAddress) == 0 {
		netLog.Warn("Handshake message validation: p2p_address is null string")
		valid = false
	}
	if len(msg.OS) == 0 {
		netLog.Warn("Handshake message validation: os field is null string")
		valid = false
	}
	if (common.CompareString(msg.Signature, ecc.NewSigNil()) != 0 || msg.Token.Equals(*crypto.NewSha256Nil())) &&
		msg.Token.Equals(*crypto.Hash256(msg.Time)) {
		netLog.Warn("Handshake message validation: token field invalid")
		valid = false
	}

	return valid
}

//sendTime populate and queue time_message immediately using incoming time_message
func (c *Connection) sendTime(msg *TimeMessage) {
	xpkt := &TimeMessage{
		Org: msg.Org,
		Rec: msg.Dst,
		Xmt: common.Now(),
	}
	c.write(xpkt)
}

func (c *Connection) sendTimeTicker() {
	xpkt := &TimeMessage{
		Org: c.rec,
		Rec: c.dst,
		Xmt: common.Now(),
	}
	c.org = xpkt.Xmt
	c.write(xpkt)
}

func (c *Connection) addPeerBlock(entry *PeerBlockState) bool {
	bPtr := c.blkState.GetById().Find(entry.ID)
	added := bPtr.IsEnd()
	if added {
		c.blkState.Insert(*entry)
	} else {
		//blk_state.modify(bptr,set_is_known);
		c.blkState.Modify(bPtr, func(state *PeerBlockState) {
			(*state).IsKnown = true
		})
		if entry.BlockNum == 0 {
			c.blkState.Modify(bPtr, func(state *PeerBlockState) {
				(*state).BlockNum = entry.BlockNum
			})
		} else {
			c.blkState.Modify(bPtr, func(state *PeerBlockState) {
				(*state).RequestedTime = common.Now()
			})
		}
	}
	return added
}

func (c *Connection) read(impl *netPluginIMpl) {
	defer func() {
		c.socket.Close()

	}()

	netLog.Info("start read message!")

	for {
		p2pMessage, err := ReadP2PMessageData(c.reader)
		if err != nil {
			fmt.Println("Error reading from p2p client:", err)
			continue
		}
		time.Sleep(100 * time.Millisecond) //TODO for testing
		//data, err := json.Marshal(p2pMessage)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//fmt.Println(c.peerAddr, ": Receive P2PMessag ", string(data))

		switch msg := p2pMessage.(type) {
		case *HandshakeMessage:
			impl.handleHandshakeMsg(c, msg)
		case *ChainSizeMessage:
			impl.handleChainSizeMsg(c, msg)
		case *GoAwayMessage:
			impl.handleGoawayMsg(c, msg)
			fmt.Printf("GO AWAY Reason[%d] \n", msg.Reason)
		case *TimeMessage:
			impl.handleTimeMsg(c, msg)
		case *NoticeMessage:
			impl.handleNoticeMsg(c, msg)
		case *RequestMessage:
			impl.handleRequestMsg(c, msg)
		case *SyncRequestMessage:
			impl.handleSyncRequestMsg(c, msg)
		case *SignedBlockMessage:
			impl.handleSignedBlock(c, msg)
		case *PackedTransactionMessage:
			impl.handlePackTransaction(c, msg)
		default:
			fmt.Println("unsupport p2pmessage type")
		}

	}
}

func ReadP2PMessageData(r io.Reader) (p2pMessage P2PMessage, err error) {
	//data := make([]byte, 0)
	lengthBytes := make([]byte, 4, 4)
	_, err = io.ReadFull(r, lengthBytes)
	if err != nil {
		return
	}
	//data = append(data, lengthBytes...)
	size := binary.LittleEndian.Uint32(lengthBytes)
	payloadBytes := make([]byte, size, size)
	count, err := io.ReadFull(r, payloadBytes)
	if count != int(size) {
		err = fmt.Errorf("readfull not full read [%d] expected[%d]", count, size)
		return
	}
	if err != nil {
		fmt.Println("readfull ,error:", err)
		return
	}
	//data = append(data, payloadBytes...)
	//fmt.Printf("receive data:  %#v\n", data)

	messagetype := P2PMessageType(payloadBytes[0])
	attr, ok := messagetype.reflectTypes()
	if !ok {
		fmt.Errorf("decode, unknown p2p message type [%d]", messagetype)
		return
	}
	msg := reflect.New(attr.ReflectType)
	err = rlp.DecodeBytes(payloadBytes[1:], msg.Interface())
	if err != nil {
		return nil, err
	}

	p2pMessage = msg.Interface().(P2PMessage)

	return
}

func (c *Connection) write(message P2PMessage) {
	payload, err := rlp.EncodeToBytes(message)
	if err != nil {
		err = fmt.Errorf("p2p message, %s", err)
		return
	}
	messageLen := uint32(len(payload) + 1)
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, messageLen)
	sendBuf := append(buf, byte(message.GetType()))
	sendBuf = append(sendBuf, payload...)

	c.socket.Write(sendBuf)

	data, _ := json.Marshal(message)
	netLog.Info("%s: send Message json: %s", c.peerAddr, string(data))

	return
}

func toProtocolVersion(v uint16) uint16 {
	if v >= netVersionBase {
		v -= netVersionBase
		if v <= netVersionRange {
			return v
		}
	}
	return 0
}
