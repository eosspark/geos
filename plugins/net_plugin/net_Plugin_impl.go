package net_plugin

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/eosspark/eos-go/plugins/net_plugin/multi_index"
	"github.com/eosspark/eos-go/plugins/net_plugin/multi_index/node_transaction"
	"github.com/eosspark/eos-go/plugins/producer_plugin"
	"net"
	"time"
)

var netLog log.Logger
var fcLog log.Logger
var peerLog log.Logger

type possibleConnections byte

const (
	nonePossible      possibleConnections = 0
	producersPossible possibleConnections = 1 << 0
	specifiedPossible possibleConnections = 1 << 1
	anyPossible       possibleConnections = 1 << 2
)

type netPluginIMpl struct {
	Listener       net.Listener
	ListenEndpoint string

	p2PAddress         string
	maxClientCount     uint32
	maxNodesPerHost    uint32
	numClients         uint32
	suppliedPeers      []string
	AllowedPeers       []ecc.PublicKey                  //< peer keys allowed to connect
	privateKeys        map[ecc.PublicKey]ecc.PrivateKey //< overlapping with producer keys, also authenticating non-producing nodes
	allowedConnections possibleConnections
	done               bool
	connectorCheck     *asio.DeadlineTimer
	transactionCheck   *asio.DeadlineTimer
	keepAliceTimer     *asio.DeadlineTimer

	connectorPeriod            time.Duration
	txnExpPeriod               time.Duration
	respExpectedPeriod         time.Duration
	keepaliveInterval          time.Duration //32*time.Sencond
	peerAuthenticationInterval time.Duration //< Peer clock may be no more than 1 second skewed from our clock, including network latency.

	maxCleanupTimeMs    int
	networkVersionMatch bool
	chainID             common.ChainIdType
	nodeID              common.NodeIdType
	userAgentName       string

	useSocketReadWatermark bool

	localTxns *node_transaction.NodeTransactionIndex

	connections []*Connection

	syncMaster *syncManager
	dispatcher *dispatchManager

	ChainPlugin *chain_plugin.ChainPlugin

	Self *NetPlugin
}

func NewNetPluginIMpl() *netPluginIMpl {
	impl := &netPluginIMpl{
		maxClientCount:             0,
		maxNodesPerHost:            1,
		numClients:                 0,
		allowedConnections:         nonePossible,
		done:                       false,
		keepaliveInterval:          32 * time.Second,
		peerAuthenticationInterval: 1 * time.Second,
		maxCleanupTimeMs:           0,
		networkVersionMatch:        false,
		txnExpPeriod:               defTxnExpireWait,
		useSocketReadWatermark:     false,
		privateKeys:                make(map[ecc.PublicKey]ecc.PrivateKey),
		localTxns:                  node_transaction.NewNodeTransactionIndex(),
		connections:                make([]*Connection, 0),
	}

	impl.syncMaster = NewSyncManager(impl, 250)
	impl.dispatcher = NewDispatchManager(impl)

	netLog = log.New("net")
	netLog.SetHandler(log.TerminalHandler)
	//impl.log.SetHandler(log.DiscardHandler())

	fcLog = log.New("fc")
	fcLog.SetHandler(log.TerminalHandler)

	peerLog = log.New("peer")
	peerLog.SetHandler(log.TerminalHandler)

	return impl
}

func (impl *netPluginIMpl) startListenLoop() {
	socket := asio.NewReactiveSocket(App().GetIoService())

	socket.AsyncAccept(impl.Listener, func(conn net.Conn, err error) {
		if conn == nil {
			log.Error("Error connect, nil")
			impl.startListenLoop()
			return
		}

		if err == nil {
			visitors := uint32(0)
			fromAddr := uint32(0)
			paddr := conn.RemoteAddr().String()
			log.Info("accept connection: %s,visitor: %d, fromAddr: %d", paddr, visitors, fromAddr)

			for _, c := range impl.connections {
				if c.socket != nil { //
					if len(c.peerAddr) == 0 {
						visitors++
						if paddr == c.socket.RemoteAddr().String() {
							fromAddr++
						}
					}
				}
			}
			if impl.numClients != visitors {
				netLog.Info("checking max client, visitors = %d, num clients %d", visitors, impl.numClients)
				impl.numClients = visitors
			}
			if fromAddr < impl.maxNodesPerHost && (impl.maxClientCount == 0 || impl.numClients < impl.maxClientCount) {
				impl.numClients++
				c := NewConnectionByConn(conn, impl)
				impl.connections = append(impl.connections, c)
				//start_session( c );TODO
			} else {
				if fromAddr >= impl.maxNodesPerHost {
					log.Error("Number of connections (%d) from %s exceeds limit", fromAddr+1, paddr)
				} else {
					log.Error("Error max_client_count %d exceeded", impl.maxClientCount)
				}
				conn.Close()
			}
		} else {
			log.Error("Error accepting connection: %s", err.Error())

			//// For the listed error codes below, recall start_listen_loop()
			//switch (ec.value()) {
			//case ECONNABORTED:
			//case EMFILE:
			//case ENFILE:
			//case ENOBUFS:
			//case ENOMEM:
			//case EPROTO:
			//	break;
			//default:
			//	return;
			//}
		}
		impl.startListenLoop()
	})
}

func (impl *netPluginIMpl) connect(peer *Connection) {

}

func (impl *netPluginIMpl) close(c *Connection) {
	if len(c.peerAddr) == 0 {
		if impl.numClients == 0 {
			netLog.Warn("num_clients already at 0") //FC
		} else {
			impl.numClients--
		}
	}
	c.close()
}

func (impl *netPluginIMpl) countOpenSockets() int {
	return len(impl.connections)
}

func (impl *netPluginIMpl) sendAll(msg P2PMessage, verify func(c *Connection) bool) {
	for _, c := range impl.connections {
		if c.current() && verify(c) {
			c.write(msg)
		}
	}
}

func (impl *netPluginIMpl) AcceptedBlockHeader(block *types.BlockState) {
	fcLog.Debug("signaled,id = %s", block.BlockId)
}

func (impl *netPluginIMpl) AcceptedBlock(block *types.BlockState) {
	fcLog.Debug("signaled,id = %s", block.BlockId)
	impl.dispatcher.bcastBlock(impl, block.SignedBlock)
}

func (impl *netPluginIMpl) IrreversibleBlock(block *types.BlockState) {
	fcLog.Debug("signaled,id = %s", block.BlockId)
}

func (impl *netPluginIMpl) AcceptedTransaction(md *types.TransactionMetadata) {
	fcLog.Debug("signaled,id = %s", md.ID)
	impl.dispatcher.bcastTransaction(md.PackedTrx)
}

func (impl *netPluginIMpl) AppliedTransaction(txn *types.TransactionTrace) {
	fcLog.Debug("signaled,id = %s", txn.ID)
}

func (impl *netPluginIMpl) AcceptedConfirmation(head types.HeaderConfirmation) {
	fcLog.Debug("signaled,id = %s", head.BlockId)
}

func (impl *netPluginIMpl) TransactionAck(results common.Pair) {
	packedTrx := results.Second.(types.PackedTransaction) //TODO  std::pair<fc::exception_ptr, packed_transaction_ptr>&
	id := packedTrx.ID()
	if results.First != nil {
		fcLog.Info("signaled NACK, trx-id = %s :%s", id, results.First)
		impl.dispatcher.rejectedTransaction(&id)
	} else {
		fcLog.Info("signaled ACK,trx-id = %s", id)
		impl.dispatcher.bcastTransaction(&packedTrx)
	}
}

func (impl *netPluginIMpl) startMonitors() {
	impl.connectorCheck = asio.NewDeadlineTimer(App().GetIoService())
	impl.transactionCheck = asio.NewDeadlineTimer(App().GetIoService())

	impl.startConnTimer(impl.connectorPeriod, nil)
	impl.startTxnTimer()
}

func (impl *netPluginIMpl) startConnTimer(du time.Duration, fromConnection *Connection) {
	impl.connectorCheck.ExpiresFromNow(impl.connectorPeriod)
	impl.connectorCheck.AsyncWait(func(err error) {
		if err != nil {
			log.Error("Error from connection check monitor: %s", err.Error())
			impl.startConnTimer(impl.connectorPeriod, nil)
		} else {
			impl.connectionMonitor(fromConnection)
		}
	})
}

func (impl *netPluginIMpl) connectionMonitor(fromConnection *Connection) {
	maxTime := common.Now()
	maxTime = maxTime.AddUs(common.Milliseconds(int64(impl.maxCleanupTimeMs)))

	var i int
	var it *Connection
	if fromConnection != nil {
		i, it = impl.findConnection(fromConnection.peerAddr)
	} else {
		i, it = 0, impl.connections[0]
	}

	for ; i < len(impl.connections); i++ {
		if common.Now().Sub(maxTime) >= 0 {
			impl.startConnTimer(time.Millisecond, it)
			return
		}
		if it.socket == nil && !it.connecting {
			if len(it.peerAddr) > 0 {
				impl.connect(it)
			} else {
				impl.eraseConnection(it)
				continue
			}
		}
	}
	impl.startConnTimer(impl.connectorPeriod, nil)
}

func (impl *netPluginIMpl) startTxnTimer() {
	impl.transactionCheck.ExpiresFromNow(impl.txnExpPeriod)
	impl.transactionCheck.AsyncWait(func(err error) {
		if err != nil {
			log.Error("Error from transaction check monitor: %s", err.Error())
			impl.startTxnTimer()
		} else {
			impl.expireTxns()
		}
	})

}

func (impl *netPluginIMpl) expireTxns() {
	impl.startTxnTimer()

	old := impl.localTxns.GetByExpiry()
	exUp := old.UpperBound(common.NewTimePointSecTp(common.Now()))
	exLo := old.LowerBound(common.TimePointSec(0))
	old.Erases(exLo, exUp)

	stale := impl.localTxns.GetByBlockNum()
	cc := impl.ChainPlugin.Chain()
	bn := cc.LastIrreversibleBlockNum()
	stale.Erases(stale.LowerBound(1), stale.UpperBound(bn))

	for _, c := range impl.connections {
		staleTxn := c.trxState.GetByBlockNum()
		staleTxn.Erases(staleTxn.LowerBound(1), staleTxn.UpperBound(bn))

		staleTxnE := c.trxState.GetByExpiry()
		staleTxnE.Erases(staleTxnE.LowerBound(common.NewTimePointSecTp(0)), staleTxnE.UpperBound(common.NewTimePointSecTp(common.Now())))

		staleBlk := c.blkState.GetByBlockNum()
		staleBlk.Erases(staleBlk.LowerBound(1), staleBlk.UpperBound(bn))
	}
}

//ticker Peer heartbeat
func (impl *netPluginIMpl) ticker() {
	impl.keepAliceTimer.ExpiresFromNow(impl.keepaliveInterval)
	impl.keepAliceTimer.AsyncWait(func(err error) {
		impl.ticker()
		if err != nil {
			log.Warn("Peer keep live ticked sooner than expected: %s", err)
		}
		for _, peer := range impl.connections {
			peer.sendTimeTicker()
		}
	})
}

//authenticatePeer determine if a peer is allowed to connect.
//Checks current connection mode and key authentication.
//return False if the peer should not connect, True otherwise.
func (impl *netPluginIMpl) authenticatePeer(msg *HandshakeMessage) bool {
	var allowedIt, privateIt, foundProducerKey bool

	if impl.allowedConnections == nonePossible {
		return false
	}
	if impl.allowedConnections == anyPossible {
		return true
	}
	if impl.allowedConnections&(producersPossible|specifiedPossible) != 0 {
		for _, pubkey := range impl.AllowedPeers {
			if pubkey == msg.Key {
				allowedIt = true
			}
		}
		_, privateIt = impl.privateKeys[msg.Key]

		pp := App().FindPlugin(producer_plugin.ProducerPlug).(*producer_plugin.ProducerPlugin)
		if pp != nil {
			foundProducerKey = pp.IsProducerKey(msg.Key)
		}
		if !allowedIt && !privateIt && !foundProducerKey {
			netLog.Error("Peer %s sent a handshake with an unauthorized key: %s", msg.P2PAddress, msg.Key)
			return false
		}
	}

	msgTime := msg.Time
	t := common.Now()
	if time.Duration(uint64(t-msgTime))*time.Microsecond > impl.peerAuthenticationInterval {
		netLog.Error("Peer %s sent a handshake with a timestamp skewed by more than 1 second", msg.P2PAddress)
		return false
	}

	if msg.Signature.String() != ecc.NewSigNil().String() && msg.Token.Equals(*crypto.NewSha256Nil()) {
		hash := crypto.Hash256(msg.Time)
		if !hash.Equals(msg.Token) {
			netLog.Error("Peer %s sent a handshake with an invalid token.", msg.P2PAddress)
			return false
		}

		peerKey, err := msg.Signature.PublicKey(msg.Token.Bytes())
		if err != nil {
			netLog.Error("Peer %s sent a handshake with an unrecoverable key.", msg.P2PAddress)
			return false
		}
		if (impl.allowedConnections&(producersPossible|specifiedPossible)) != 0 && peerKey.String() != msg.Key.String() {
			netLog.Error("Peer %s sent a handshake with an unauthenticated key.", msg.P2PAddress)
			return false
		}
	} else if impl.allowedConnections&(producersPossible|specifiedPossible) != 0 {
		netLog.Debug("Peer sent a handshake with blank signature and token,but this node accepts only authenticate connections.")
		return false
	}

	return true
}

//getAuthenticationKey retrieve public key used to authenticate with peers.
//Finds a key to use for authentication.  If this node is a producer, use
//the front of the producer key map.  If the node is not a producer but has
//a configured private key, use it.  If the node is neither a producer nor has
//a private key, returns an empty key.
//On a node with multiple private keys configured, the key with the first
//numerically smaller byte will always be used.
func (impl *netPluginIMpl) getAuthenticationKey() *ecc.PublicKey {
	if len(impl.privateKeys) > 0 {
		for pubKey := range impl.privateKeys { //TODO easier  ？？？
			return &pubKey
		}
		/*producer_plugin* pp = app().find_plugin<producer_plugin>();
		if(pp != nullptr && pp->get_state() == abstract_plugin::started)
		   return pp->first_producer_public_key();*/
		return &ecc.PublicKey{}
	}

	return &ecc.PublicKey{}
}

//signCompact returns a signature of the digest using the corresponding private key of the signer.
//If there are no configured private keys, returns an empty signature.
func (impl *netPluginIMpl) signCompact(signer *ecc.PublicKey, digest *crypto.Sha256) *ecc.Signature {
	privateKeyPtr, ok := impl.privateKeys[*signer]
	if ok {
		signature, err := privateKeyPtr.Sign(digest.Bytes())
		if err != nil {
			netLog.Error("signCompact is error: %s", err.Error())
			return ecc.NewSigNil()
		}
		return &signature
	} else {
		pp := App().FindPlugin("ProducerPlugin").(*producer_plugin.ProducerPlugin)
		if pp != nil && pp.GetState() == Started {
			return pp.SignCompact(signer, *digest)
		}
	}
	return ecc.NewSigNil()
}

func (impl *netPluginIMpl) handleChainSizeMsg(c *Connection, msg *ChainSizeMessage) {
	netLog.Info("%s : receives chain_size_message", c.peerAddr)
}

func (impl *netPluginIMpl) handleHandshakeMsg(c *Connection, msg *HandshakeMessage) {
	netLog.Info("%s : receives handshake_message", c.peerAddr)
	if !isValid(msg) {
		netLog.Error("%s : bad handshake message", c.peerAddr)
		goAwayMsg := &GoAwayMessage{
			Reason: fatalOther,
			NodeID: *crypto.NewSha256Nil(),
		}
		c.write(goAwayMsg)
		return
	}

	//netLog.Info("%s : receive a handshake message %v", c.peerAddr, msg)

	cc := App().FindPlugin("ChainPlugin").(*chain_plugin.ChainPlugin).Chain()
	libNum := cc.LastIrreversibleBlockNum()
	peerLib := msg.LastIrreversibleBlockNum
	if c.connecting {
		c.connecting = false
	}

	if msg.Generation == 1 {
		if msg.NodeID.Equals(c.nodeID) {
			netLog.Error("Self connection detected. Closing connection")
			goAwayMsg := &GoAwayMessage{
				Reason: fatalOther,
				NodeID: *crypto.NewSha256Nil(),
			}
			c.write(goAwayMsg)
		}

		if len(c.peerAddr) == 0 || c.lastHandshakeRecv.NodeID.Equals(*crypto.NewSha256Nil()) {
			netLog.Info("checking for duplicate")
			for _, check := range impl.connections {
				if check == c {
					continue
				}
				if check.connected() && check.PeerName() == msg.P2PAddress {
					// It's possible that both peers could arrive here at relatively the same time, so
					// we need to avoid the case where they would both tell a different connection to go away.
					// Using the sum of the initial handshake times of the two connections, we will
					// arbitrarily (but consistently between the two peers) keep one of them.
					if msg.Time+c.lastHandshakeSent.Time <= check.lastHandshakeSent.Time+check.lastHandshakeRecv.Time {
						continue
					}
					netLog.Debug("sending go_away duplicate to %s", msg.P2PAddress)
					goAwayMsg := &GoAwayMessage{
						Reason: duplicate,
						NodeID: c.nodeID,
					}
					//c.enqueue(goAwayMsg)
					c.write(goAwayMsg)
					c.noRetry = duplicate
					return
				}
			}
		} else {
			netLog.Debug("skipping dulicate check, addr ==%s, id == %s", c.peerAddr, c.lastHandshakeRecv.NodeID)
		}

		if !msg.ChainID.Equals(impl.chainID) {
			netLog.Error("Peer on different chain. Closing connection")
			goAwayMsg := &GoAwayMessage{
				Reason: wrongChain,
				NodeID: *crypto.NewSha256Nil(),
			}
			c.write(goAwayMsg)
			return
		}

		c.protocolVersion = toProtocolVersion(msg.NetworkVersion)
		if c.protocolVersion != netVersion {
			if impl.networkVersionMatch {
				netLog.Error("Peer network version does not match expected %d but got %d", netVersion, c.protocolVersion)
				goAwayMsg := &GoAwayMessage{
					Reason: wrongVersion,
					NodeID: *crypto.NewSha256Nil(),
				}
				c.write(goAwayMsg)
				return
			} else {
				netLog.Info("local network version: %d Remote version: %d", netVersion, c.protocolVersion)
			}
		}

		if !c.nodeID.Equals(msg.NodeID) {
			c.nodeID = msg.NodeID
		}

		if !impl.authenticatePeer(msg) {
			netLog.Error("Peer not authenticated. Closing connection")
			goAwayMsg := &GoAwayMessage{
				Reason: authentication,
				NodeID: *crypto.NewSha256Nil(),
			}
			c.write(goAwayMsg)
			return
		}

		onFork := false
		netLog.Debug("lib_num = %d peer_lib = %d")
		if peerLib <= libNum && peerLib > 0 {
			Try(func() {
				peerLibID := cc.GetBlockIdForNum(peerLib)
				onFork = !msg.LastIrreversibleBlockID.Equals(peerLibID)
			}).Catch(func(ex UnknownBlockException) {
				netLog.Warn("peer last irreversible block %d is unknown", peerLib)
				onFork = true
			}).Catch(func(e interface{}) {
				netLog.Warn("caught an exception getting block id for %d", peerLib)
				onFork = true
			}).End()

			if onFork {
				netLog.Error("Peer chain is forked")
				goAwayMsg := &GoAwayMessage{
					Reason: forked,
					NodeID: *crypto.NewSha256Nil(),
				}
				c.write(goAwayMsg)
				return
			}
		}

		if c.sentHandshakeCount == 0 {
			c.sendHandshake(impl)
		}
	}

	c.lastHandshakeRecv = msg
	impl.syncMaster.recvHandshake(c, msg)

}

func (impl *netPluginIMpl) handleGoawayMsg(c *Connection, msg *GoAwayMessage) {
	rsn := ReasonToString[msg.Reason]
	netLog.Info("%s : receive a go_away_message", c.peerAddr)
	netLog.Info("receive go_away_message reason = %s", rsn)
	c.noRetry = msg.Reason
	if msg.Reason == duplicate {
		c.nodeID = msg.NodeID
	}
	//p.flushQueues()
	c.close()
}

//handleTimeMsg process time_message
//Calculate offset, delay and dispersion.  Note carefully the
//implied processing.  The first-order difference is done
//directly in 64-bit arithmetic, then the result is converted
//to floating double.  All further processing is in
//floating-double arithmetic with rounding done by the hardware.
//This is necessary in order to avoid overflow and preserve precision.
func (impl *netPluginIMpl) handleTimeMsg(c *Connection, msg *TimeMessage) {
	netLog.Info("receive time_message")
	netLog.Info("%s: receive a time message %v", c.peerAddr, msg)

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
	//fmt.Println(c.dst)
	if msg.Org == 0 {
		c.sendTime(msg)
		return // We don't have enough data to perform the calculation yet.
	}

	//p.offset = float64((p.rec-p.org)+(msg.Xmt-p.dst)) / 2
	//NsecPerUsec := float64(1000)
	//netLog.Info("Clock offset is %v ns  %v us", p.offset, p.offset/NsecPerUsec)

	c.org = 0
	c.rec = 0
}

func (impl *netPluginIMpl) handleNoticeMsg(c *Connection, msg *NoticeMessage) {
	// peer tells us about one or more blocks or txns. When done syncing, forward on
	// notices of previously unknown blocks or txns,
	netLog.Info("%s : receive notice_message", c.peerAddr)
	netLog.Info("%s : received notice_message %v", c.peerAddr, msg)

	c.connecting = false
	req := RequestMessage{}
	sendReq := false

	if msg.KnownTrx.Mode != none {
		netLog.Debug("this is a %s notice with %d transactions",
			modeTostring[msg.KnownTrx.Mode], msg.KnownTrx.Pending)
	}

	switch msg.KnownTrx.Mode {
	case none:
	case lastIrrCatchUp:
		c.lastHandshakeRecv.HeadNum = msg.KnownTrx.Pending
		req.ReqTrx.Mode = none
	case catchUp:
		if msg.KnownTrx.Pending > 0 {
			//plan to get all except what we already know about
			req.ReqTrx.Mode = catchUp
			sendReq = true
			knownSum := impl.localTxns.Size()
			if knownSum > 0 {
				//for( const auto& t : local_txns.get<by_id>( ) ) {//TODO
				//	req.req_trx.ids.push_back( t.id );
				//}

				//ltx :=impl.localTxns.GetById()
				//req.ReqTrx.IDs =append(req.ReqTrx.IDs,ltx.Begin())
			}
		}
	case normal:
		impl.dispatcher.recvNotice(c, msg, false)
	}

	if msg.KnownBlocks.Mode != none {
		netLog.Debug("this is a %s notice with  %d blocks",
			modeTostring[msg.KnownBlocks.Mode], msg.KnownBlocks.Pending)
	}
	switch msg.KnownBlocks.Mode {
	case none:
		if msg.KnownTrx.Mode != normal {
			return
		}
	case lastIrrCatchUp, catchUp:
		impl.syncMaster.recvNotice(c, msg)
	case normal:
		impl.dispatcher.recvNotice(c, msg, false)
	default:
		netLog.Error("bad notice_message : invalid known.mode %d", msg.KnownBlocks.Mode)
	}
	netLog.Debug("send req = %t", sendReq)
	if sendReq {
		c.write(&req)
	}
}

func (impl *netPluginIMpl) handleRequestMsg(c *Connection, msg *RequestMessage) {
	netLog.Info("%s: received request_message %v", c.peerAddr, msg)

	switch msg.ReqBlocks.Mode {
	case catchUp:
		peerLog.Info("%s : received request_message:catch_up", c.peerAddr)
		c.blkSendBranch()
	case normal:
		peerLog.Info("%s : receive request_message:normal", c.peerAddr)
		c.blkSend(msg.ReqBlocks.IDs)
	default:
	}

	switch msg.ReqTrx.Mode {
	case catchUp:
		c.txnSendPending(msg.ReqTrx.IDs)
	case normal:
		c.txnSend(msg.ReqTrx.IDs)
	case none:
		if msg.ReqBlocks.Mode == none {
			c.stopSend()
		}
	default:
	}
}

func (impl *netPluginIMpl) handleSyncRequestMsg(c *Connection, msg *SyncRequestMessage) {
	netLog.Info("%s : received sync_request_message %v", c.peerAddr, msg)
	if msg.EndBlock == 0 {
		c.peerRequested = nil //TODO
		//c.peerRequested.reset()
		//c.flushQueues()
	} else {
		c.peerRequested = newSyncState(msg.StartBlock, msg.EndBlock, msg.StartBlock-1)
		c.enqueueSyncBlock()
	}
}

func (impl *netPluginIMpl) handlePackTransaction(c *Connection, msg *PackedTransactionMessage) {
	fcLog.Debug("got a packed transaction ,cancel wait")
	peerLog.Info(" %s receive packed transaction", c.peerAddr)

	cc := impl.ChainPlugin.Chain()
	if cc.GetReadMode() == chain.READONLY {
		fcLog.Debug("got a txn in read-only mode - dropping")
		return
	}

	if impl.syncMaster.isActive(c) {
		fcLog.Debug("got a txn during sync - dropping")
		return
	}
	tid := msg.ID()
	c.cancelWait()
	if !impl.localTxns.GetById().Find(tid).IsEnd() {
		fcLog.Debug("got a duplicate transaction - dropping")
		return
	}

	impl.dispatcher.recvTransaction(c, &tid)

	impl.ChainPlugin.AcceptTransaction(&msg.PackedTransaction, func(result common.StaticVariant) {
		if exception, ok := result.(Exception); ok {
			peerLog.Debug("bad packed_transaction : %s", exception.DetailMessage())
		} else {
			trace, _ := result.(types.TransactionTrace)
			if trace.Except == nil {
				fcLog.Debug("chain accepted transaction")
				impl.dispatcher.bcastTransaction(&msg.PackedTransaction)
				return
			}
			peerLog.Error("bad packed_transaction : %s", trace.Except.DetailMessage())
		}
		impl.dispatcher.rejectedTransaction(&tid)
	})
}

func (impl *netPluginIMpl) handleSignedBlock(c *Connection, msg *SignedBlockMessage) {
	netLog.Info("%s : receive signed_block message %v", c.peerAddr, msg)

	cc := impl.ChainPlugin.Chain()
	blkID := msg.BlockID()
	blkNum := msg.BlockNumber()
	fcLog.Debug("canceling wait on %s", c.peerAddr)
	c.cancelWait()

	Try(func() {
		if cc.FetchBlockById(blkID) != nil {
			impl.syncMaster.recvBlock(c, blkID, blkNum)
			return
		}
	}).Catch(func(e interface{}) {
		// should this even be caught?
		log.Error("Caught an unknown exception trying to recall blockID")
	}).End()

	impl.dispatcher.recvBlock(c, &blkID, blkNum)
	age := common.Now().Sub(msg.Timestamp.ToTimePoint())
	peerLog.Info("received signed_block : %d block age in secs = %d", blkNum, age.ToSeconds())

	reason := GoAwayReason(fatalOther)
	Try(func() {
		sbp := msg.SignedBlock
		impl.ChainPlugin.AcceptBlock(&sbp)
		reason = noReason
	}).Catch(func(ex UnlinkableBlockException) {
		peerLog.Error("bad signed_block : %s", ex.DetailMessage())
		reason = unlinkable
	}).Catch(func(ex BlockValidateException) {
		peerLog.Error("bad signed_block : %s", ex.DetailMessage())
		reason = validation
	}).Catch(func(ex AssertException) {
		peerLog.Error("bad signed_block : %s", ex.DetailMessage())
		log.Error("unable to accept block on assert exception %s from %s", ex.DetailMessage(), c.peerAddr)
	}).Catch(func(ex FcException) {
		peerLog.Error("bad signed_block : %s", ex.DetailMessage())
		log.Error("accept_block threw a non-assert exception %s from %s", ex.DetailMessage(), c.peerAddr)
		reason = noReason
	}).Catch(func(ex interface{}) {
		peerLog.Error("bad signed_block : unknown exception")
		log.Error("handle sync block caught something else from %s", c.peerAddr)
	}).End()

	if reason == noReason {
		var id common.TransactionIdType
		for _, recpt := range msg.Transactions {
			if recpt.Trx.TransactionID.Equals(*crypto.NewSha256Nil()) {
				id = recpt.Trx.TransactionID
			} else {
				id = recpt.Trx.PackedTransaction.ID()
			}

			ltx := impl.localTxns.GetById().Find(id)
			if !ltx.IsEnd() {
				impl.localTxns.Modify(ltx, func(state *multi_index.NodeTransactionState) {
					if state.Requests > 0 {
						(*state).TrueBlock = blkNum
					}
					(*state).BlockNum = blkNum
				})
			}

			ctx := c.trxState.GetById().Find(id)
			if !ctx.IsEnd() {
				c.trxState.Modify(ctx, func(state *multi_index.TransactionState) {
					(*state).BlockNum = blkNum
				})
			}
		}

		impl.syncMaster.recvBlock(c, blkID, blkNum)
	} else {
		impl.syncMaster.rejectedBlock(c, blkNum)
	}

}

func (impl *netPluginIMpl) findConnection(host string) (int, *Connection) {
	for i, c := range impl.connections {
		if c.peerAddr == host {
			return i, c
		}
	}
	return -1, &Connection{}
}

func (impl *netPluginIMpl) eraseConnection(it *Connection) {
	for i, c := range impl.connections {
		if c.peerAddr == it.peerAddr {
			impl.connections = append(impl.connections[:i], impl.connections[i+1:]...)
		}
	}
}
