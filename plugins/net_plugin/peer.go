package net_plugin

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"io"
	"net"
	"reflect"
	"runtime"
	"time"
)

const (
	tsBufferSize int = 32
)

type Peer struct {
	//blkState peerBlockStateIndex
	//trxState transactionStateIndex
	//peerRequested optional<syncState>//this peer is requesting info from us
	//socket socketPtr

	connection           net.Conn
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
	responseExpected   time.Timer
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

}

type PeerStatus struct {
	Peer          string
	Connecting    bool
	Syncing       bool
	LastHandshake HandshakeMessage
}

func NewPeer(conn net.Conn, reader io.Reader) *Peer {
	return &Peer{
		connection:         conn,
		reader:             reader,
		peerAddr:           conn.RemoteAddr().String(),
		lastHandshakeSent:  &HandshakeMessage{},
		lastHandshakeRecv:  &HandshakeMessage{},
		sentHandshakeCount: 0,
	}
}

func (p *Peer) getStatus() *PeerStatus {
	return &PeerStatus{
		Peer:          p.peerAddr,
		Connecting:    p.connecting,
		Syncing:       p.syncing,
		LastHandshake: *p.lastHandshakeRecv,
	}
}

func (p *Peer) connected() bool {

	return false
}

func (p *Peer) current() bool {

	return false
}

func (p *Peer) reset() {

}

func (p *Peer) close() {

}

func (p *Peer) sendHandshake(impl *netPluginIMpl) {
	handshakePopulate(impl, p.lastHandshakeSent)
	p.sentHandshakeCount += 1
	p.lastHandshakeSent.Generation = p.sentHandshakeCount
	//fc_dlog(logger, "Sending handshake generation ${g} to ${ep}",
	//	("g",last_handshake_sent.generation)("ep", peer_name()));

	fmt.Printf("Sending handshake generation %d to %s\n", p.lastHandshakeSent.Generation, p.peerAddr)
	p.write(p.lastHandshakeSent)
}

func handshakePopulate(impl *netPluginIMpl, hello *HandshakeMessage) {
	hello.NetworkVersion = netVersionBase + netVersion
	hello.ChainID = impl.chainID
	hello.NodeID = impl.nodeID
	hello.Key = *impl.getAuthenticationKey()
	hello.Time = common.Now()
	hello.Token = crypto.Hash256(hello.Time)
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
		try.Try(func() {
			//hello.LastIrreversibleBlockID = cc.get_block_id_for_num(hello.last_irreversible_block_num)
			//hello.LastIrreversibleBlockID =
		}).Catch(func(ex exception.UnknownBlockException) {
			//ilog("caught unkown_block");
			fmt.Println("caught unkown_block")
			//hello.LastIrreversibleBlockNum =0
		}).End()
	}
	if hello.HeadNum > 0 {
		try.Try(func() {
			//hello.id = cc.get_block_id_for_num( hello.head_num )
			//hello.id =
		}).Catch(func(ex exception.UnknownBlockException) {
			hello.HeadNum = 0
		}).End()
	}

}

func (p *Peer) PeerName() string {
	return p.peerAddr
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
func (p *Peer) syncWait() {

}

func (p *Peer) cancelWait() {

}

func (p *Peer) fetchWait() {

}

func (p *Peer) cancelSync(reason GoAwayReason) {
	//fc_dlog(logger,"cancel sync reason = ${m}, write queue size ${o} peer ${p}",
	//	("m",reason_str(reason)) ("o", write_queue.size())("p", peer_name()));

	fmt.Println("cancel sync reason = %s, write queue size %d peer %s\n", ReasonToString[reason], p.write)

	p.cancelWait()
	//p.flushQueues()
	switch reason {
	case validation, fatalOther:
		p.noRetry = reason
		p.write(&GoAwayMessage{Reason: reason})
	default:
		//fc_dlog(logger, "sending empty request but not calling sync wait on ${p}", ("p",peer_name()))
		fmt.Println("sending empty request but not calling sync wait on %s\n", p.peerAddr)
		p.write(&SyncRequestMessage{0, 0})
	}

}

//void txn_send_pending(const vector<transaction_id_type> &ids);
//void txn_send(const vector<transaction_id_type> &txn_lis);
//
//void blk_send_branch();
//void blk_send(const vector<block_id_type> &txn_lis);
//void stop_send();
//
//void enqueue( const net_message &msg, bool trigger_send = true );
//void cancel_sync(go_away_reason);
//void flush_queues();
//bool enqueue_sync_block();

func (p *Peer) requestSyncBlocks(start, end uint32) {
	syncRequest := SyncRequestMessage{
		StartBlock: start,
		EndBlock:   end,
	}
	p.write(&syncRequest)
	p.syncWait()
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
		fmt.Printf("Handshake message validation: last irreversible block %d is greater than head block %d\n",
			msg.LastIrreversibleBlockNum, msg.HeadNum)
		valid = false
	}
	if len(msg.P2PAddress) == 0 {
		fmt.Println("Handshake message validation: p2p_address is null string")
		valid = false
	}
	if len(msg.OS) == 0 {
		fmt.Println("Handshake message validation: os field is null string")
		valid = false
	}
	if (common.CompareString(msg.Signature, ecc.NewSigNil()) != 0 || common.CompareString(msg.Token, crypto.NewSha256Nil()) != 0) &&
		common.CompareString(msg.Token, crypto.Hash256(msg.Time)) != 0 {
		fmt.Println("Handshake message validation: token field invalid")
		valid = false
	}

	return valid
}

//sendTime populate and queue time_message immediately using incoming time_message
func (p *Peer) sendTime(msg *TimeMessage) {
	xpkt := &TimeMessage{
		Org: msg.Org,
		Rec: msg.Dst,
		Xmt: common.Now(),
	}
	p.write(xpkt)
}

func (p *Peer) sendTimeTicker() {
	xpkt := &TimeMessage{
		Org: p.rec,
		Rec: p.dst,
		Xmt: common.Now(),
	}
	p.org = xpkt.Xmt
	p.write(xpkt)
}

func (p *Peer) addPeerBlock(entry *peerBlockState) bool { //TODO
	//auto bptr = blk_state.get<by_id>().find(entry.id);
	//bool added = (bptr == blk_state.end());
	//if (added){
	//	blk_state.insert(entry);
	//}
	//else {
	//blk_state.modify(bptr,set_is_known);
	//if (entry.block_num == 0) {
	//blk_state.modify(bptr,update_block_num(entry.block_num));
	//}
	//else {
	//blk_state.modify(bptr,set_request_time);
	//}
	//}
	//return added;
	return false
}

func (p *Peer) read(impl *netPluginIMpl) {
	defer func() {
		p.connection.Close()
		impl.loopWG.Done()
	}()

	impl.loopWG.Add(1)
	fmt.Println("start read message!")

	for {
		p2pMessage, err := ReadP2PMessageData(p.reader)
		if err != nil {
			fmt.Println("Error reading from p2p client:", err)
			continue
		}
		time.Sleep(100 * time.Millisecond) //TODO for testing
		//data, err := json.Marshal(p2pMessage)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//fmt.Println(p.peerAddr, ": Receive P2PMessag ", string(data))

		switch msg := p2pMessage.(type) {
		case *HandshakeMessage:
			impl.handleHandshakeMsg(p, msg)
		case *ChainSizeMessage:
			impl.handleChainSizeMsg(p, msg)
		case *GoAwayMessage:
			impl.handleGoawayMsg(p, msg)
			fmt.Printf("GO AWAY Reason[%d] \n", msg.Reason)
		case *TimeMessage:
			impl.handleTimeMsg(p, msg)
		case *NoticeMessage:
			impl.handleNoticeMsg(p, msg)
		case *RequestMessage:
			impl.handleRequestMsg(p, msg)
		case *SyncRequestMessage:
			impl.handleSyncRequestMsg(p, msg)
		case *SignedBlockMessage:
			impl.handleSignedBlock(p, msg)
		case *PackedTransactionMessage:
			impl.handlePackTransaction(p, msg)
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

func (p *Peer) write(message P2PMessage) {
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

	p.connection.Write(sendBuf)

	//fmt.Println(p.peerAddr, ": Message bytes", sendBuf)
	data, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p.peerAddr, ": send Message json:", string(data))

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
