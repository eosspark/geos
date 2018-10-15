package net_plugin

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"io"
	"net"
	"reflect"
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
	lastHandshajeSent  *HandshakeMessage
	sentHandshakeCount int16
	connecting         bool
	syncing            bool
	protocolVersion    uint16
	peerAddr           string
	responseExpected   time.Timer
	//pendingFetch optional<request_message>

	noRetry     GoAwayReason
	forkHead    common.BlockIdType
	forkHeadNum uint32
	//lastReq optional<request_message>

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

//func NewPeer() *Peer {
//	return &Peer{}
//}
func NewPeer(conn net.Conn, reader io.Reader) *Peer {
	return &Peer{
		connection: conn,
		reader:     reader,
		peerAddr:   conn.RemoteAddr().String(),
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
func (p *Peer) sendHandshake() {

}

func (p *Peer) PeerName() string {

	return ""
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
//void request_sync_blocks (uint32_t start, uint32_t end);
//
//void cancel_wait();
//void sync_wait();
//void fetch_wait();
//void sync_timeout(boost::system::error_code ec);
//void fetch_timeout(boost::system::error_code ec);
//
//void queue_write(std::shared_ptr<vector<char>> buff,
//bool trigger_send,
//std::function<void(boost::system::error_code, std::size_t)> callback);
//void do_queue_write();
//
///** \brief Process the next message from the pending message buffer
// *
// * Process the next message from the pending_message_buffer.
// * message_length is the already determined length of the data
// * part of the message and impl in the net plugin implementation
// * that will handle the message.
// * Returns true is successful. Returns false if an error was
// * encountered unpacking or processing the message.
// */
//bool process_next_message(net_plugin_impl& impl, uint32_t message_length);
//
//bool add_peer_block(const peer_block_state &pbs);

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

func (p *Peer) handleHandshakeMsg(np *netPluginIMpl, msg *HandshakeMessage) {
	fmt.Println("receive a handshake message")
	if !isValid(msg) {
		fmt.Println("bad handshake message")
		goAwayMsg := &GoAwayMessage{
			Reason: fatalOther,
			NodeID: *crypto.NewSha256Nil(),
		}
		p.write(goAwayMsg)
		return
	}

	//controller& cc = chain_plug->chain();
	//uint32_t lib_num = cc.last_irreversible_block_num( );
	//uint32_t peer_lib = msg.last_irreversible_block_num;

	//libNum := uint32(100) //TODO
	//peerLib := msg.LastIrreversibleBlockNum

	if msg.Generation == 1 {
		if crypto.Sha256(msg.NodeID).Compare(crypto.Sha256(p.nodeID)) {
			//elog( "Self connection detected. Closing connection")
			fmt.Println("Self connection detected. Closing connection")
			goAwayMsg := &GoAwayMessage{
				Reason: fatalOther,
				NodeID: *crypto.NewSha256Nil(),
			}
			p.write(goAwayMsg)
		}

		//TODO check for duplicate!!
		//if( c->peer_addr.empty() || c->last_handshake_recv.node_id == fc::sha256()) {
		//fc_dlog(logger, "checking for duplicate" )
		//}

		if msg.ChainID.String() != p2pChainIDString {
			//elog( "Peer on a different chain. Closing connection")
			fmt.Println("Peer on different chain. Closing connection")
			goAwayMsg := &GoAwayMessage{
				Reason: wrongChain,
				NodeID: *crypto.NewSha256Nil(),
			}
			p.write(goAwayMsg)
			return
		}

		p.protocolVersion = toProtocolVersion(msg.NetworkVersion)
		if p.protocolVersion != netVersion {
			if np.networkVersionMatch {
				//elog("Peer network version does not match expected ${nv} but got ${mnv}",
				//	("nv", net_version)("mnv", c->protocol_version))
				fmt.Printf("Peer network version does not match expected %d but got %d", netVersion, p.protocolVersion)
				goAwayMsg := &GoAwayMessage{
					Reason: wrongVersion,
					NodeID: *crypto.NewSha256Nil(),
				}
				p.write(goAwayMsg)
				return
			} else {
				//ilog("Local network version: ${nv} Remote version: ${mnv}",
				//	("nv", net_version)("mnv", c->protocol_version))
				fmt.Printf("local network version: %d Remote version: %d", netVersion, p.protocolVersion)
			}
		}

		if p.nodeID.String() != msg.NodeID.String() {
			p.nodeID = msg.NodeID
		}

		//fmt.Println("authrnticatePeer")//TODO check for authenticatePeer!!!
		//if !np.authenticatePeer(msg){
		//	//elog("Peer not authenticated.  Closing connection.")
		//	fmt.Println("Peer not authenticated. Closing connection")
		//	goAwayMsg := &GoAwayMessage{
		//		Reason: authentication,
		//		NodeID: *crypto.NewSha256Nil(),
		//	}
		//	p.write(goAwayMsg)
		//	return
		//}

		//onFork := false//TODO check for onFork!!!
		////fc_dlog(logger, "lib_num = ${ln} peer_lib = ${pl}",("ln",lib_num)("pl",peer_lib));
		//if peerLib <= libNum && peerLib > 0 {
		//	//peerLibID := cc.getBlockIdForNum(peerLib)
		//	//onFork = msg.LastIrreversibleBlockID != peerLibID
		//	onFork = true
		//
		//	if onFork {
		//		//elog( "Peer chain is forked");
		//		fmt.Println("Peer chain is forked")
		//		goAwayMsg := &GoAwayMessage{
		//			Reason: forked,
		//			NodeID: *crypto.NewSha256Nil(),
		//		}
		//		p.write(goAwayMsg)
		//		return
		//	}
		//}

		if p.sentHandshakeCount == 0 {
			p.sendHandshake()
		}
	}

	p.lastHandshakeRecv = msg
	//c->_logger_variant.reset();

	//sync_master->recv_handshake(p,msg)

}

func (p *Peer) handleChainSizeMsg(msg *ChainSizeMessage) {
	//peer_ilog(c, "received chain_size_message")
	fmt.Println("receives chain_size_message")
}

func (p *Peer) handleGoawayMsg(msg *GoAwayMessage) {
	rsn := ReasonToString[msg.Reason]
	fmt.Printf("receive go_away_message reason = %s\n", rsn)
	p.noRetry = msg.Reason
	//if msg.Reason == duplicate {
	//	p.nodeID = msg.NodeID
	//}
	//c.flushQueues()
	//close(c)

}

// handleTimeMsg process time_message
// Calculate offset, delay and dispersion.  Note carefully the
// implied processing.  The first-order difference is done
// directly in 64-bit arithmetic, then the result is converted
// to floating double.  All further processing is in
// floating-double arithmetic with rounding done by the hardware.
// This is necessary in order to avoid overflow and preserve precision.
func (p *Peer) handleTimeMsg(msg *TimeMessage) {
	fmt.Println("receive time_message")
	/* We've already lost however many microseconds it took to dispatch
	 * the message, but it can't be helped.
	 */
	msg.Dst = common.Now()

	// If the transmit timestamp is zero, the peer is horribly broken.
	if msg.Xmt == 0 {
		return /* invalid timestamp */
	}
	if msg.Xmt == p.xmt {
		return /* duplicate packet */
	}

	p.xmt = msg.Xmt
	p.rec = msg.Rec
	p.dst = msg.Dst
	if msg.Org == 0 {
		p.sendTime(msg)
		return // We don't have enough data to perform the calculation yet.
	}

	//p.offset = float64((p.rec-p.org)+(msg.Xmt-p.dst)) / 2
	//fmt.Println(p.offset)

	//NsecPerUsec := float64(1000)
	//fmt.Printf("Clock offset is %v ns  %v us\n", p.offset, p.offset/NsecPerUsec)
	//if(logger.is_enabled(fc::log_level::all))
	//logger.log(FC_LOG_MESSAGE(all, "Clock offset is ${o}ns (${us}us)", ("o", c->offset)("us", c->offset/NsecPerUsec)));
	p.org = 0
	p.rec = 0

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

func (p *Peer) handleNoticeMsg(msg *NoticeMessage) {
	// peer tells us about one or more blocks or txns. When done syncing, forward on
	// notices of previously unknown blocks or txns,
	//
	fmt.Println("received notice_message")
	p.connecting = false
	req := RequestMessage{}
	sendReq := false
	if msg.KnownTrx.Mode != none {
		fmt.Printf("this is a %s notice with %d blocks \n",
			modeTostring[msg.KnownTrx.Mode], msg.KnownTrx.Pending)
	}
	switch msg.KnownTrx.Mode {
	case none:
	case lastIrrCatchUp:
		//c.lastHandshakeRecv.HeadNum = &msg.KnownTrx.Pending
		req.ReqTrx.Mode = none
	case catchUp:
		if msg.KnownTrx.Pending > 0 {
			//plan to get all except what we already know about
			req.ReqTrx.Mode = catchUp
			sendReq = true
			//knownSum := local_txns.size()
			//if( known_sum ) {
			//	for( const auto& t : local_txns.get<by_id>( ) ) {
			//	req.req_trx.ids.push_back( t.id )
			//	}
			//}
		}
	case normal:
		//dispatcher.recvNotice(c,msg,false)
	}

	if msg.KnownBlocks.Mode != none {
		fmt.Printf("this is a %s notice with  %d blocks\n",
			modeTostring[msg.KnownBlocks.Mode], msg.KnownBlocks.Pending)
	}
	switch msg.KnownBlocks.Mode {
	case none:
		if msg.KnownTrx.Mode != normal {
			return
		}
	case lastIrrCatchUp:
	case catchUp:
		//syncMaster.recvNotice(c,msg)
	case normal:
		//dispatcher.recvNotice(c,msg,false)
	default:
		fmt.Printf("bad notice_message : invalid knwon_blocks.mode %d\n", msg.KnownBlocks.Mode)
	}
	//fc_dlog(logger, "send req = ${sr}", ("sr",send_req));

	if sendReq {
		//c.enqueue(req)
	}

}

func (p *Peer) handleRequestMsg(msg *RequestMessage) {
	switch msg.ReqBlocks.Mode {
	case catchUp:
		fmt.Println("received request_message:catch_up")
		//c.blkSendBranch()
	case normal:
		fmt.Println("receive request_message:normal")
		//c.blkSend(msg.ReqBlocks.IDs)
	default:

	}

	switch msg.ReqTrx.Mode {
	case catchUp:
		//c.txnSendPending(msg.ReqTrx.IDs)
	case normal:
		//c.txnSend(msg.ReqTrx.IDs)
	case none:
		if msg.ReqBlocks.Mode == none {
			//c.stopSend()
		}
	default:

	}
}

func (p *Peer) handleSyncRequestMsg(msg *SyncRequestMessage) {
	if msg.EndBlock == 0 {
		//c.peerRequested.reset()
		//c.flushQueues()
	} else {
		//c.peerRequested = syncState(msg.StartBlock,msg.EndBlock,msg.StartBlock-1)
		//c.enqueueSyncBlock()
	}

}

func (p *Peer) handleSignedBlock(msg *SignedBlockMessage) {
	fmt.Println("receive signed_block message")
	//cc := chain_plug->chain();
	// blkID := msg.ID();
	//blkNum := msg.BlockNum();
	//fmt.Printf("canceling wait on %s\n",c.perrName())
	//c.cancel_wait();
	fmt.Printf("signed Block : %v\n", msg)
}

func (p *Peer) handlePackTransaction(msg *PackedTransactionMessage) {
	fmt.Println("receive packed transaction")
	tid := msg.ID()
	fmt.Println(tid)
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

func (p *Peer) read(np *netPluginIMpl) {
	fmt.Println("start read message!")

	for {
		p2pMessage, err := ReadP2PMessageData(p.reader)
		if err != nil {
			fmt.Println("Error reading from p2p client:", err)
			continue
		}

		data, err := json.Marshal(p2pMessage)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(p.peerAddr, ": Receive P2PMessag ", string(data))

		switch msg := p2pMessage.(type) {
		case *HandshakeMessage:
			p.handleHandshakeMsg(np, msg)
		case *ChainSizeMessage:
			p.handleChainSizeMsg(msg)
		case *GoAwayMessage:
			p.handleGoawayMsg(msg)
			fmt.Printf("GO AWAY Reason[%d] \n", msg.Reason)
		case *TimeMessage:
			p.handleTimeMsg(msg)
		case *NoticeMessage:
			p.handleNoticeMsg(msg)
		case *RequestMessage:
			p.handleRequestMsg(msg)
		case *SyncRequestMessage:
			p.handleSyncRequestMsg(msg)
		case *SignedBlockMessage:
			p.handleSignedBlock(msg)
		case *PackedTransactionMessage:
			p.handlePackTransaction(msg)
		default:
			fmt.Println("unsupport p2pmessage type")
		}

	}
}

//data := make([]byte, 100)
//n, err := p.connection.Read(data)
//if err!=nil{
//	fmt.Println(err)
//}
////fmt.Println(string(data[:n]))
//fmt.Println(data[:n])
func ReadP2PMessageData(r io.Reader) (p2pMessage P2PMessage, err error) {
	data := make([]byte, 0)
	lengthBytes := make([]byte, 4, 4)
	_, err = io.ReadFull(r, lengthBytes)
	if err != nil {
		return
	}
	data = append(data, lengthBytes...)
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
	data = append(data, payloadBytes...)
	fmt.Println("receive data:  ", data)

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

	//fmt.Println(p.peerAddr, ": 已发送Message", sendBuf)
	data, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p.peerAddr, ": 已发送Message", string(data))

	return
}

func toProtocolVersion(v uint16) uint16 {
	if v > netVersionBase {
		v -= netVersionBase
		if v <= netVersionRange {
			return v
		}
	}
	return 0
}
