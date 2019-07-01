package net_plugin

import (
	"encoding/binary"
	"fmt"
	"net"
	"reflect"
	"runtime"

	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/libraries/asio"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	. "github.com/eosspark/eos-go/plugins/net_plugin/multi_index"
	"github.com/eosspark/eos-go/plugins/net_plugin/multi_index/peer_block_state"
	"github.com/eosspark/eos-go/plugins/net_plugin/multi_index/transaction_state"
)

type PeerStatus struct {
	Peer          string           `json:"peer"`
	Connecting    bool             `json:"connecting"`
	Syncing       bool             `json:"syncing"`
	LastHandshake HandshakeMessage `json:"last_handshake"`
}

type queuedWrite struct {
	buff     []byte
	callback func(err error, n int)
}

//Index by start_block_num
type syncState struct {
	startBlock uint32
	endBlock   uint32
	last       uint32           //last sent or received
	startTime  common.TimePoint //time request made or received
}

func newSyncState(start, end, lastActed uint32) *syncState {
	return &syncState{
		startBlock: start,
		endBlock:   end,
		last:       lastActed,
		startTime:  common.Now(),
	}
}

type Connection struct {
	blkState      *peer_block_state.PeerBlockStateIndex
	trxState      *transaction_state.TransactionStateIndex
	peerRequested *syncState //this peer is requesting info from us

	socket               *asio.ReactiveSocket
	conn                 net.Conn
	pendingMessageBuffer [1024 * 1024]byte
	nodeID               common.NodeIdType
	lastHandshakeRecv    *HandshakeMessage
	lastHandshakeSent    *HandshakeMessage
	sentHandshakeCount   uint16
	connecting           bool
	syncing              bool
	protocolVersion      uint16
	peerAddr             string
	responseExpected     *asio.DeadlineTimer
	pendingFetch         *RequestMessage
	noRetry              GoAwayReason
	forkHead             common.BlockIdType
	forkHeadNum          uint32
	lastReq              *RequestMessage

	//outstandingReadBytes int //optional
	bufTemp []byte

	// Members set from network data
	org    common.TimePoint //originate timestamp
	rec    common.TimePoint //receive timestamp
	dst    common.TimePoint //destination timestamp
	xmt    common.TimePoint // transmit timestamp
	offset float64          //peer offset

	writeQueue []queuedWrite
	outQueue   []queuedWrite

	impl *netPluginIMpl
}

func NewConnectionByEndPoint(endpoint string, impl *netPluginIMpl) *Connection {
	conn := &Connection{
		blkState:      peer_block_state.NewPeerBlockStateIndex(),
		trxState:      transaction_state.NewTransactionStateIndex(),
		peerRequested: &syncState{},
		//conn:,
		socket:             asio.NewReactiveSocket(App().GetIoService()),
		nodeID:             common.NodeIdType(crypto.NewSha256Nil()),
		lastHandshakeRecv:  &HandshakeMessage{},
		lastHandshakeSent:  &HandshakeMessage{},
		sentHandshakeCount: 0,
		connecting:         false,
		syncing:            false,
		protocolVersion:    0,
		peerAddr:           endpoint,
		//responseExpected:
		pendingFetch: &RequestMessage{},
		noRetry:      noReason,
		forkHead:     common.BlockIdNil(),
		forkHeadNum:  0,
		lastReq:      &RequestMessage{},
		impl:         impl,
	}
	netLog.Warn("accepted network connection")
	conn.initialize()

	return conn
}

func NewConnectionByConn(socket *asio.ReactiveSocket, c net.Conn, impl *netPluginIMpl) *Connection {
	conn := &Connection{
		blkState:           peer_block_state.NewPeerBlockStateIndex(),
		trxState:           transaction_state.NewTransactionStateIndex(),
		peerRequested:      &syncState{},
		conn:               c,
		socket:             socket,
		nodeID:             common.NodeIdType(crypto.NewSha256Nil()),
		lastHandshakeRecv:  &HandshakeMessage{},
		lastHandshakeSent:  &HandshakeMessage{},
		sentHandshakeCount: 0,
		connecting:         true,
		syncing:            false,
		protocolVersion:    0,
		peerAddr:           c.RemoteAddr().String(), //,
		//responseExpected:,
		pendingFetch: &RequestMessage{},
		noRetry:      noReason,
		forkHead:     common.BlockIdNil(),
		forkHeadNum:  0,
		lastReq:      &RequestMessage{},
		impl:         impl,
	}
	netLog.Warn("accepted network connection")
	conn.initialize()

	return conn
}

func (c *Connection) initialize() {
	rnd := &c.nodeID
	rnd.Hash[0] = 0
	c.responseExpected = asio.NewDeadlineTimer(App().GetIoService())
}

func (c *Connection) getStatus() *PeerStatus {
	return &PeerStatus{
		Peer:          c.peerAddr,
		Connecting:    c.connecting,
		Syncing:       c.syncing,
		LastHandshake: *c.lastHandshakeRecv,
	}
}

func (c *Connection) connected() bool { //TODO
	return c.socket != nil && !c.connecting
}

func (c *Connection) current() bool {
	return c.connected() && !c.syncing
}

func (c *Connection) reset() {
	c.peerRequested = &syncState{}
	c.blkState = peer_block_state.NewPeerBlockStateIndex()
	c.trxState = transaction_state.NewTransactionStateIndex()
}

func (c *Connection) close() {
	if c.socket != nil {
		c.socket = nil
	} else {
		netLog.Warn("no socket to close")
	}
	c.flushQueues()
	c.connecting = false
	c.syncing = false
	if !common.Empty(c.lastReq) {
		c.impl.dispatcher.retryFetch(c)
	}
	c.reset()
	c.sentHandshakeCount = 0
	c.lastHandshakeRecv = &HandshakeMessage{}
	c.lastHandshakeSent = &HandshakeMessage{}
	c.impl.syncMaster.resetLibNum(c)
	FcLog.Debug("cancel wait on %s", c.PeerName())
	c.cancelWait()
	c.bufTemp = nil
}

func (c *Connection) sendHandshake() {
	c.handshakePopulate(c.impl, c.lastHandshakeSent)
	c.sentHandshakeCount += 1
	c.lastHandshakeSent.Generation = c.sentHandshakeCount

	FcLog.Info("Sending handshake generation %d to %s", c.lastHandshakeSent.Generation, c.peerAddr)
	c.enqueue(c.lastHandshakeSent, true)
}

func (c *Connection) handshakePopulate(impl *netPluginIMpl, hello *HandshakeMessage) {
	hello.NetworkVersion = netVersionBase + netVersion
	hello.ChainID = impl.chainID
	hello.NodeID = impl.nodeID
	hello.Key = *impl.getAuthenticationKey()
	hello.Time = common.Now()
	hello.Token = *crypto.Hash256(hello.Time)
	hello.Signature = *impl.signCompact(&hello.Key, &hello.Token)

	// If we couldn't sign, don't send a token.
	if common.Empty(hello.Signature) {
		hello.Token = crypto.NewSha256Nil()
	}

	hello.P2PAddress = impl.p2PAddress + " - " + hello.NodeID.String()[:7]

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

	cc := impl.ChainPlugin.Chain()
	hello.HeadID = common.BlockIdNil()
	hello.LastIrreversibleBlockID = common.BlockIdNil()
	hello.HeadNum = cc.ForkDbHeadBlockNum()
	hello.LastIrreversibleBlockNum = cc.LastIrreversibleBlockNum()

	if hello.LastIrreversibleBlockNum > 0 {
		Try(func() {
			hello.LastIrreversibleBlockID = cc.GetBlockIdForNum(hello.LastIrreversibleBlockNum)
		}).Catch(func(ex *exception.UnknownBlockException) {
			netLog.Info("caught unknown_block")
			hello.LastIrreversibleBlockNum = 0
		}).End()
	}
	if hello.HeadNum > 0 {
		Try(func() {
			hello.HeadID = cc.GetBlockIdForNum(hello.HeadNum)
		}).Catch(func(ex *exception.UnknownBlockException) {
			hello.HeadNum = 0
		}).End()
	}

}

func (c *Connection) PeerName() string {
	if len(c.lastHandshakeRecv.P2PAddress) != 0 {
		return c.lastHandshakeRecv.P2PAddress
	}

	if len(c.peerAddr) != 0 {
		return c.peerAddr
	}
	return "connecting client"
}

func (c *Connection) cancelWait() {
	if c.responseExpected != nil {
		c.responseExpected.Cancel()
	}
}

func (c *Connection) syncWait() {

	//c.responseExpected.ExpiresFromNow(c.impl.respExpectedPeriod)
	//c.responseExpected.AsyncWait(func(err error) {
	//	if c == nil {
	//		// connection was destroyed before this lambda was delivered
	//		return
	//	}
	//	c.syncTimeout(err)
	//})
}

func (c *Connection) fetchWait() {
	//c.responseExpected.ExpiresFromNow(c.impl.respExpectedPeriod)
	//c.responseExpected.AsyncWait(func(err error) {
	//	if c == nil {
	//		// connection was destroyed before this lambda was delivered
	//		return
	//	}
	//	c.fetchTimeout(err)
	//})
}

func (c *Connection) syncTimeout(err error) { //TODO not same as C++
	if err == nil {
		c.impl.syncMaster.reassignFetch(c, benignOther)
	} else {
		netLog.Error("setting timer for sync request fot error %s", err)
	}
}

func (c *Connection) fetchTimeout(err error) { //TODO not same as C++
	if err == nil {
		if c.pendingFetch != nil && c.pendingFetch.ReqTrx.empty() || c.pendingFetch.ReqBlocks.empty() {
			c.impl.dispatcher.retryFetch(c)
		}
	} else {
		netLog.Error("setting time for fetch request got error %s", err)
	}

}

func (c *Connection) cancelSync(reason GoAwayReason) {
	FcLog.Debug("cancel sync reason = %s, write queue size %d peer %s", ReasonStr[reason], len(c.writeQueue), c.peerAddr)

	c.cancelWait()
	c.flushQueues()

	switch reason {
	case validation, fatalOther:
		c.noRetry = reason
		c.enqueue(&GoAwayMessage{Reason: reason}, true)
	default:
		FcLog.Debug("sending empty request but not calling sync wait on %s", c.peerAddr)
		c.enqueue(&SyncRequestMessage{0, 0}, true)
	}
}

func (c *Connection) txnSendPending(ids []common.TransactionIdType) {
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

				c.queueWrite(tx.Value().SerializedTxn, true, func(err error, n int) {
					localTxns := c.impl.localTxns
					tx := localTxns.GetById().Find(tx.Value().ID)
					if !tx.IsEnd() {
						localTxns.Modify(tx, func(state *NodeTransactionState) {
							exp := state.Expires.SecSinceEpoch()
							state.Expires = common.TimePointSec(exp - 1*60)
							if state.Requests == 0 {
								state.TrueBlock = state.BlockNum
								state.BlockNum = 0
							}
							state.Requests = state.Requests - 1
							if state.Requests == 0 {
								state.BlockNum = state.TrueBlock
							}
						})
					} else {
						FcLog.Warn("Local pending TX erased before queued_write called callback")
					}
				})
			}
		}
	}
}

func (c *Connection) txnSend(ids []common.TransactionIdType) {
	for _, t := range ids {
		tx := c.impl.localTxns.GetById().Find(t)
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

			c.queueWrite(tx.Value().SerializedTxn, true, func(err error, n int) {
				localTxns := c.impl.localTxns
				tx := localTxns.GetById().Find(t)
				if !tx.IsEnd() {
					localTxns.Modify(tx, func(state *NodeTransactionState) {
						exp := state.Expires.SecSinceEpoch()
						state.Expires = common.TimePointSec(exp - 1*60)
						if state.Requests == 0 {
							state.TrueBlock = state.BlockNum
							state.BlockNum = 0
						}
						state.Requests = state.Requests - 1
						if state.Requests == 0 {
							state.BlockNum = state.TrueBlock
						}
					})
				} else {
					FcLog.Warn("Local TX erased before queued_write called callback")
				}
			})
		}
	}
}

func (c *Connection) blkSendBranch() {
	cc := c.impl.ChainPlugin.Chain()
	headNum := cc.ForkDbHeadBlockNum()

	var note NoticeMessage

	note.KnownBlocks.Mode = normal
	note.KnownBlocks.Pending = 0
	FcLog.Debug("head_num = %d", headNum)
	if headNum == 0 {
		c.enqueue(&note, true)
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
			FcLog.Debug("maybe truncating branch at = %d : %s", remoteHeadNum, remoteHeadID)
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
		netLog.Error("unable to retrieve block info: %s for %s", ex.What(), c.peerAddr)
		c.enqueue(&note, true)
		returning = true
	}).Catch(func(ex exception.Exception) {}).Catch(func(interface{}) {}).End()

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
				c.enqueue(&SignedBlockMessage{*bStack[i-1]}, true)
			}
		}
		FcLog.Info("Sent %d blocks on my fork", count)
	} else {
		FcLog.Info("Nothing to send on fork request")
	}

	c.syncing = false
}

func (c *Connection) blkSend(ids []common.BlockIdType) {
	cc := c.impl.ChainPlugin.Chain()
	var count int
	breaking := false
	for _, blkID := range ids {
		count++
		Try(func() {
			b := cc.FetchBlockById(blkID)
			if b != nil {
				FcLog.Debug("found block for id ar num %d", b.BlockNumber())
				c.enqueue(&SignedBlockMessage{*b}, true)
			} else {
				FcLog.Info("fetch block by id returned null, id %s on block %d of %d for %s",
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

func (c *Connection) flushQueues() {
	c.writeQueue = nil
}

func (c *Connection) enqueueSyncBlock() bool {
	cc := App().FindPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin).Chain()
	//if common.Empty(c.peerRequested){
	//	return false
	//}

	if c.peerRequested.startTime == 0 { //TODO check nil
		return false
	}

	c.peerRequested.last++
	num := c.peerRequested.last
	triggerSend := num == c.peerRequested.startBlock
	if num == c.peerRequested.endBlock {
		c.peerRequested = &syncState{}
	}

	result := false
	Try(func() {
		sb := cc.FetchBlockByNumber(num)
		if sb != nil {
			c.enqueue(&SignedBlockMessage{*sb}, triggerSend)
			result = true
		}
	}).Catch(func(e interface{}) {
		FcLog.Warn("write loop exception")
	}).End()

	return result
}

func (c *Connection) requestSyncBlocks(start, end uint32) {
	syncRequest := SyncRequestMessage{
		StartBlock: start,
		EndBlock:   end,
	}
	c.enqueue(&syncRequest, true)
	//c.syncWait()
}

func isValid(msg *HandshakeMessage) bool {
	// Do some basic validation of an incoming handshake_message, so things
	// that really aren't handshake messages can be quickly discarded without
	// affecting state.
	valid := true
	if msg.LastIrreversibleBlockNum > msg.HeadNum {
		FcLog.Warn("Handshake message validation: last irreversible block %d is greater than head block %d",
			msg.LastIrreversibleBlockNum, msg.HeadNum)
		valid = false
	}
	if len(msg.P2PAddress) == 0 {
		FcLog.Warn("Handshake message validation: p2p_address is null string")
		valid = false
	}
	if len(msg.OS) == 0 {
		FcLog.Warn("Handshake message validation: os field is null string")
		valid = false
	}
	if (common.CompareString(msg.Signature, ecc.NewSigNil()) != 0 || msg.Token.Equals(crypto.NewSha256Nil())) &&
		msg.Token.Equals(*crypto.Hash256(msg.Time)) {
		FcLog.Warn("Handshake message validation: token field invalid")
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
	c.enqueue(xpkt, true)
}

func (c *Connection) sendTimeTicker() {
	xpkt := &TimeMessage{
		Org: c.rec,
		Rec: c.dst,
		Xmt: common.Now(),
	}
	c.org = xpkt.Xmt
	c.enqueue(xpkt, true)
}

func (c *Connection) addPeerBlock(entry *PeerBlockState) bool {
	bPtr := c.blkState.GetById().Find(entry.ID)
	added := bPtr.IsEnd()
	if added {
		c.blkState.Insert(*entry)
	} else {
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

func (c *Connection) processNextMessage(payloadBytes []byte) bool {
	result := true
	Try(func() {
		messageType := NetMessageType(payloadBytes[0])
		attr, ok := messageType.reflectTypes()
		if !ok {
			Throw(fmt.Errorf("processNextMessage, unknown p2p message type %d", messageType))
		}
		msg := reflect.New(attr.ReflectType)
		err := rlp.DecodeBytes(payloadBytes[1:], msg.Interface())
		if err != nil {
			Throw(err)
		}

		netMsg := msg.Interface().(NetMessage)
		switch msg := netMsg.(type) {
		case *HandshakeMessage:
			c.impl.handleHandshake(c, msg)
		case *ChainSizeMessage:
			c.impl.handleChainSize(c, msg)
		case *GoAwayMessage:
			c.impl.handleGoaway(c, msg)
		case *TimeMessage:
			c.impl.handleTime(c, msg)
		case *NoticeMessage:
			c.impl.handleNotice(c, msg)
		case *RequestMessage:
			c.impl.handleRequest(c, msg)
		case *SyncRequestMessage:
			c.impl.handleSyncRequest(c, msg)
		case *SignedBlockMessage:
			c.impl.handleSignedBlock(c, msg)
		case *PackedTransactionMessage:
			c.impl.handlePackTransaction(c, msg)
		default:
			Throw(fmt.Errorf("unsuppoted p2p message type %d", messageType))
		}

	}).Catch(func(e exception.FcException) {
		netLog.Error("read message is error:%s", e.DetailMessage())
		c.impl.close(c)
		result = false
	}).End()
	return result
}

func (c *Connection) enqueue(m NetMessage, triggerSend bool) {
	closeAfterSend := noReason
	if m.GetType() == GoAwayMessageType {
		closeAfterSend = m.(*GoAwayMessage).Reason
	}

	payload, _ := rlp.EncodeToBytes(m)
	messageLen := uint32(len(payload) + 1)
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, messageLen)
	sendBuf := append(buf, byte(m.GetType()))
	sendBuf = append(sendBuf, payload...)

	FcLog.Debug("send message :%d,%s", m.GetType(), m.String())

	c.queueWrite(sendBuf, triggerSend, func(err error, n int) {
		if c != nil {
			if closeAfterSend != noReason {
				netLog.Error("sent a go away message: %s, closing connection to %s", ReasonStr[closeAfterSend], c.PeerName())
				c.impl.close(c)
				return
			}
		} else {
			FcLog.Warn("connection expired before enqueued net_message called callback!")
		}
	})
}

func (c *Connection) queueWrite(buf []byte, triggerSend bool, callback func(err error, n int)) {
	c.writeQueue = append(c.writeQueue, queuedWrite{buf, callback})
	if len(c.outQueue) == 0 && triggerSend {
		c.doQueueWrite()
	}
}

func (c *Connection) doQueueWrite() {
	if len(c.writeQueue) == 0 || len(c.outQueue) > 0 {
		return
	}

	if c.socket == nil {
		FcLog.Error("socket not open to %s", c.peerAddr)
		c.impl.close(c)
		return
	}

	bufs := make([]byte, 0)
	for i := 0; i < len(c.writeQueue); i++ {
		m := c.writeQueue[i]
		bufs = append(bufs, m.buff...)
		c.outQueue = append(c.outQueue, m)
	}
	c.writeQueue = nil

	c.socket.AsyncWrite(c.conn, bufs, func(n int, err error) {
		Try(func() {
			for _, m := range c.outQueue {
				m.callback(err, n)
			}
			if err != nil {
				pName := "no connection name"
				if c != nil {
					pName = c.PeerName()
				}
				netLog.Error("error sending to peer %s,error is %s", pName, err.Error())
				FcLog.Info("connection closure detected o write to %s", pName)
				c.impl.close(c)
				return
			}
			c.outQueue = nil

			c.enqueueSyncBlock()
			c.doQueueWrite()

		}).Catch(func(e interface{}) {
			conn := c
			var pName string
			if conn != nil {
				pName = c.peerAddr
			} else {
				pName = "no connection name"
			}

			netLog.Error("Exception in do_queue_write to %s: %s", pName, e)
		}).End()
	})
}
