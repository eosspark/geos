package p2p

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"github.com/eosspark/eos-go/rlp"
	"time"
)

const (
	tsBufferSize int = 32
)

type Connection struct {
	//blkState peerBlockStateIndex
	//trxState transactionStateIndex
	//peerRequested optional<syncState>//this peer is requesting info from us
	//socket socketPtr
	pendingMessageBuffer [1024 * 1024]byte
	//outstandingReadBytes optional<int>
	blkBuffer []byte

	nodeID             rlp.Sha256
	lastHandshakeRecv  HandshakeMessage
	lastHandshajeSent  HandshakeMessage
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

type connectionStatus struct {
	peer          string
	connecting    bool
	syncing       bool
	lastHandshake HandshakeMessage
}

func NewConection() *Connection {
	return &Connection{}
}

func (c *Connection) getStatus() *connectionStatus {
	return &connectionStatus{
		peer:          c.peerAddr,
		connecting:    c.connecting,
		syncing:       c.syncing,
		lastHandshake: c.lastHandshakeRecv,
	}
}

func (c *Connection) connected() bool {

	return false
}

func (c *Connection) current() bool {

	return false
}

func (c *Connection) reset() {

}

func (c *Connection) close() {

}
func (c *Connection) sendHandshake() {

}

func (c *Connection) PeerName() string {

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
	if (common.CompareString(msg.Signature, ecc.NewSigNil()) != 0 || common.CompareString(msg.Token, rlp.NewSha256Nil()) != 0) &&
		common.CompareString(msg.Token, rlp.Hash256(msg.Time)) != 0 {
		fmt.Println("Handshake message validation: token field invalid")
		valid = false
	}
	return valid
}

func (c *Connection) handleHandshakeMsg(msg *HandshakeMessage) {
	if !isValid(msg) {
		fmt.Println("bad handshake message")
		//c->enqueue( GoAwayMessage{GoAwayMessage ,rlp.NewSha256Nil()});
		return
	}

}

func (c *Connection) handleChainSizeMsg(msg *ChainSizeMessage) {
	//peer_ilog(c, "received chain_size_message")
	fmt.Println("receives chain_size_message")
}

func (c *Connection) handleGoawayMsg(msg *GoAwayMessage) {
	rsn := ReasonToString[msg.Reason]
	fmt.Printf("receive go_away_message reason = %s\n", rsn)
	c.noRetry = msg.Reason
	if msg.Reason == duplicate {
		c.nodeID = msg.NodeID
	}
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
func (c *Connection) handleTimeMsg(msg *TimeMessage) {
	fmt.Println("receive time_message")
	/* We've already lost however many microseconds it took to dispatch
	 * the message, but it can't be helped.
	 */
	msg.Dst = common.Now()

	// If the transmit timestamp is zero, the peer is horribly broken.
	if msg.Xmt == 0 {
		return /* invalid timestamp */
	}
	if msg.Xmt == c.xmt {
		return /* duplicate packet */
	}

	c.xmt = msg.Xmt
	c.rec = msg.Rec
	c.dst = msg.Dst
	if msg.Org == 0 {
		c.sendTime(msg)
		return // We don't have enough data to perform the calculation yet.
	}

	c.offset = float64((c.rec-c.org)+(msg.Xmt-c.dst)) / 2
	fmt.Println(c.offset)

	NsecPerUsec := float64(1000)
	fmt.Printf("Clock offset is %v ns  %v us", c.offset, c.offset/NsecPerUsec)
	//if(logger.is_enabled(fc::log_level::all))
	//logger.log(FC_LOG_MESSAGE(all, "Clock offset is ${o}ns (${us}us)", ("o", c->offset)("us", c->offset/NsecPerUsec)));
	c.org = 0
	c.rec = 0

}

//sendTime populate and queue time_message immediately using incoming time_message
func (c *Connection) sendTime(msg *TimeMessage) {
	xpkt := &TimeMessage{
		Org: msg.Org,
		Rec: msg.Dst,
		Xmt: common.Now(),
	}
	enqueue(xpkt)
}

func (c *Connection) handleNoticeMsg(msg *NoticeMessage) {
	// peer tells us about one or more blocks or txns. When done syncing, forward on
	// notices of previously unknown blocks or txns,
	//
	fmt.Println("received notice_message")
	c.connecting = false
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

func (c *Connection) handleRequestMsg(msg *RequestMessage) {
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

func (c *Connection) handleSyncRequestMsg(msg *SyncRequestMessage) {
	if msg.EndBlock == 0 {
		//c.peerRequested.reset()
		//c.flushQueues()
	} else {
		//c.peerRequested = syncState(msg.StartBlock,msg.EndBlock,msg.StartBlock-1)
		//c.enqueueSyncBlock()
	}

}

func (c *Connection) handleSignedBlock(msg *SignedBlockMessage) {
	fmt.Println("receive signed_block message")
	//cc := chain_plug->chain();
	// blkID := msg.ID();
	//blkNum := msg.BlockNum();
	//fmt.Printf("canceling wait on %s\n",c.perrName())
	//c.cancel_wait();
	fmt.Printf("signed Block : %v\n", msg)
}

func (c *Connection) handlePackTransaction(msg *PackedTransactionMessage) {
	fmt.Println("receive packed transaction")
	tid := msg.ID()
	fmt.Println(tid)
}

func enqueue(t interface{}) { //TODO

}
