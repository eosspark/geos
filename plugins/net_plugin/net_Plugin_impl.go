package net_plugin

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"strings"
	"time"

	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	. "github.com/eosspark/eos-go/libraries/asio"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/eosspark/eos-go/plugins/net_plugin/multi_index"
	"github.com/eosspark/eos-go/plugins/net_plugin/multi_index/node_transaction"
	"github.com/eosspark/eos-go/plugins/producer_plugin"
)

var netLog log.Logger
var FcLog = log.GetLoggerMap()["net_plugin"]

type possibleConnections byte

const (
	// default value initializers
	defSendBufferSizeMb        = 4
	defSendBufferSize          = 1024 * 1024 * defSendBufferSizeMb
	defMaxClients              = 25 // 0 for unlimited clients
	defMaxNodesPerHost         = 1
	defConnRetryWait           = 30
	defTxnExpireWait           = time.Duration(3 * time.Second)
	defRespExpectedWait        = time.Duration(5 * time.Second)
	defSyncFetchSpan           = 100
	defMaxJustSend      uint32 = 1500 // roughly 1 "mtu"
	largeMsgNotify      bool   = false

	messageHeaderSize = 4

	/*	   For a while, network version was a 16 bit value equal to the second set of 16 bits
		   of the current build's git commit id. We are now replacing that with an integer protocol
		   identifier. Based on historical analysis of all git commit identifiers, the larges gap
		   between ajacent commit id values is shown below.
		   these numbers were found with the following commands on the master branch:

		   git log | grep "^commit" | awk '{print substr($2,5,4)}' | sort -u > sorted.txt
		   rm -f gap.txt; prev=0; for a in $(cat sorted.txt); do echo $prev $((0x$a - 0x$prev)) $a >> gap.txt; prev=$a; done; sort -k2 -n gap.txt | tail

		   DO NOT EDIT net_version_base OR net_version_range!
	*/
	netVersionBase  uint16 = 0x04b5
	netVersionRange uint16 = 106

	//If there is a change to network protocol or behavior, increment net version to identify
	//the need for compatibility hooks
	protoBase         uint16 = 0
	protoExplicitSync uint16 = 1

	netVersion uint16 = protoExplicitSync

	nonePossible      possibleConnections = 0
	producersPossible possibleConnections = 1 << 0
	specifiedPossible possibleConnections = 1 << 1
	anyPossible       possibleConnections = 1 << 2
)

type netPluginIMpl struct {
	syncMaster *syncManager
	dispatcher *dispatchManager

	localTxns *node_transaction.NodeTransactionIndex

	Listener           net.Listener
	p2PAddress         string
	resolver           *ReactiveSocket
	maxClientCount     uint32
	maxNodesPerHost    uint32
	numClients         uint32
	suppliedPeers      []string
	AllowedPeers       []ecc.PublicKey                  //< peer keys allowed to connect
	privateKeys        map[ecc.PublicKey]ecc.PrivateKey //< overlapping with producer keys, also authenticating non-producing nodes
	allowedConnections possibleConnections

	connectorCheck   *DeadlineTimer
	transactionCheck *DeadlineTimer
	keepAliceTimer   *DeadlineTimer

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

	connections []*Connection

	ChainPlugin *chain_plugin.ChainPlugin
	context     context.Context

	Self *NetPlugin
}

func NewNetPluginIMpl(io *IoContext) *netPluginIMpl {
	impl := &netPluginIMpl{
		maxClientCount:             0,
		maxNodesPerHost:            1,
		numClients:                 0,
		allowedConnections:         nonePossible,
		keepaliveInterval:          32 * time.Second,
		peerAuthenticationInterval: 1 * time.Second,
		maxCleanupTimeMs:           0,
		networkVersionMatch:        false,
		txnExpPeriod:               defTxnExpireWait,
		privateKeys:                make(map[ecc.PublicKey]ecc.PrivateKey),
		localTxns:                  node_transaction.NewNodeTransactionIndex(),
		connections:                make([]*Connection, 0),
		resolver:                   NewReactiveSocket(io),
		context:                    context.Background(),
		suppliedPeers:              make([]string, 0),
	}

	impl.syncMaster = NewSyncManager(impl, 100)
	impl.dispatcher = NewDispatchManager(impl)

	netLog = log.New("net")
	netLog.SetHandler(log.TerminalHandler)
	//impl.log.SetHandler(DiscardHandler())

	return impl
}

func (impl *netPluginIMpl) startListenLoop() {
	socket := NewReactiveSocket(App().GetIoService())

	socket.AsyncAccept(impl.Listener, func(conn net.Conn, err error) {
		if conn == nil {
			netLog.Error("Error connect, nil")
			impl.startListenLoop()
			return
		}

		if err == nil {
			visitors := uint32(0)
			fromAddr := uint32(0)
			pAddr := conn.RemoteAddr().String()
			netLog.Info("accept connection: %s,visitor: %d, fromAddr: %d", pAddr, visitors, fromAddr)

			for _, c := range impl.connections {
				if c.conn != nil {
					if len(c.peerAddr) == 0 {
						visitors++
						if pAddr == c.conn.RemoteAddr().String() {
							fromAddr++
						}
					}
				}
			}
			if impl.numClients != visitors {
				FcLog.Info("checking max client, visitors = %d, num clients %d", visitors, impl.numClients)
				impl.numClients = visitors
			}
			if fromAddr < impl.maxNodesPerHost && (impl.maxClientCount == 0 || impl.numClients < impl.maxClientCount) {
				impl.numClients++
				c := NewConnectionByConn(socket, conn, impl)
				impl.connections = append(impl.connections, c)
				impl.startReadMessage(socket, c)

			} else {
				if fromAddr >= impl.maxNodesPerHost {
					netLog.Error("Number of connections (%d) from %s exceeds limit", fromAddr+1, pAddr)
				} else {
					netLog.Error("Error max_client_count %d exceeded", impl.maxClientCount)
				}
				conn.Close()
			}
		} else {
			netLog.Error("Error accepting connection: %s", err.Error())
		}

		impl.startListenLoop()
	})
}

func (impl *netPluginIMpl) startReadMessage(socket *ReactiveSocket, conn *Connection) {
	returning := false
	pendingMessageBuffer := make([]byte, 0)
	Try(func() {
		if conn.conn == nil {
			returning = true
			return
		}

		buf := make([]byte, 4096)
		socket.AsyncRead(conn.conn, buf, func(n int, err error) {
			if conn == nil {
				returning = true
				return
			}

			Try(func() {
				if err == nil {
					if len(conn.bufTemp) > 0 {
						pendingMessageBuffer = append(conn.bufTemp, buf[:n]...)
					} else {
						pendingMessageBuffer = buf[:n]
					}
					conn.bufTemp = nil
					n = len(pendingMessageBuffer)

					for i := 0; i < n; {
						bytesInBuf := n - i
						if bytesInBuf < messageHeaderSize {
							conn.bufTemp = append(conn.bufTemp, pendingMessageBuffer[i:n]...)
							break
						} else {
							messageLength := int(binary.LittleEndian.Uint32(pendingMessageBuffer[i : i+4]))
							if messageLength > defSendBufferSize*2 || messageLength == 0 {
								netLog.Error("incoming message length unexpected %d, from %s", messageLength, conn.conn.RemoteAddr())
								impl.close(conn)
								returning = true
								return
							}
							totalMessageBytes := messageLength + messageHeaderSize
							if bytesInBuf >= totalMessageBytes {
								if !conn.processNextMessage(pendingMessageBuffer[i+4 : i+4+totalMessageBytes]) {
									returning = true
									return
								}
								i = i + totalMessageBytes
							} else {
								conn.bufTemp = append(conn.bufTemp, pendingMessageBuffer[i:n]...)
								break
							}
						}
					}
					impl.startReadMessage(socket, conn)

				} else {
					pName := conn.peerAddr
					if err == io.EOF {
						netLog.Error("Error reading message from %s: %s", pName, err)
					} else {
						netLog.Info("Peer %s closed connection", pName)
					}
					impl.close(conn)
				}

			}).Catch(func(ex interface{}) {
				pName := "no connection name"
				if conn != nil {
					pName = conn.PeerName()
				}
				netLog.Error("Exception in handling read data from %s %s", pName, ex)
				impl.close(conn)

			}).End()

			if returning {
				return
			}
		})
	}).Catch(func(e interface{}) {
		pName := "no connection name"
		if conn != nil {
			pName = conn.PeerName()
		}
		netLog.Error("Undefined exception handling reading %s", pName)
		impl.close(conn)

	}).End()

}

func (impl *netPluginIMpl) connect(c *Connection) {
	if c.noRetry != noReason {
		FcLog.Debug("Skipping connect due to go_away reason %s", ReasonStr[c.noRetry])
		return
	}

	colon := strings.IndexAny(c.peerAddr, ":")
	if colon == -1 || colon == 0 || colon == len(c.peerAddr)-1 {
		netLog.Error("Invalid peer address. must be \"host:port\" : %s", c.peerAddr)
		i, itr := impl.findConnection(c.peerAddr)
		if i != -1 {
			itr.reset()
			impl.close(itr)
			impl.eraseConnection(itr)
		}
		return
	}

	host := c.peerAddr[:colon]
	port := c.peerAddr[colon+1 : len(c.peerAddr)]
	netLog.Info("host:post =%s,%s", host, port)

	impl.resolver.AsyncResolve(impl.context, host, port, func(address string, err error) {
		if c == nil {
			return
		}
		if err == nil {
			impl.connect2(c, address)
		} else {
			netLog.Error("Unable to resolve %s:%s,%s", host, port, err)
		}
	})
}

func (impl *netPluginIMpl) connect2(c *Connection, endPoint string) {
	if c.noRetry != GoAwayReason(noReason) {
		rsn := ReasonStr[c.noRetry]
		netLog.Warn("no retry is %s", rsn)
		return
	}

	c.connecting = true
	c.socket.AsyncConnect("tcp", endPoint, func(conn net.Conn, err error) {
		if c == nil {
			return
		}
		if err == nil {
			c.conn = conn
			impl.startReadMessage(c.socket, c)
			c.sendHandshake()

		} else {
			netLog.Error("connection failed to %s:%s", c.PeerName(), err.Error())
			c.connecting = false
			impl.close(c)
		}
	})

}

func (impl *netPluginIMpl) close(c *Connection) {
	if len(c.peerAddr) == 0 {
		if impl.numClients == 0 {
			FcLog.Warn("num_clients already at 0")
		} else {
			impl.numClients--
		}
	}
	impl.eraseConnection(c)
	c.close()
}

func (impl *netPluginIMpl) countOpenSockets() int {
	return len(impl.connections)
}

func (impl *netPluginIMpl) sendAll(msg NetMessage, verify func(c *Connection) bool) {
	for _, c := range impl.connections {
		if c.current() && verify(c) {
			c.enqueue(msg, true)
		}
	}
}

func (impl *netPluginIMpl) AcceptedBlockHeader(block *types.BlockState) {
	FcLog.Debug("signaled,id = %s", block.BlockId)
}

func (impl *netPluginIMpl) AcceptedBlock(block *types.BlockState) {
	FcLog.Debug("signaled,id = %s", block.BlockId)
	//impl.dispatcher.bcastBlock(impl, block.SignedBlock)
}

func (impl *netPluginIMpl) IrreversibleBlock(block *types.BlockState) {
	FcLog.Debug("signaled,id = %s", block.BlockId)
}

func (impl *netPluginIMpl) AcceptedTransaction(md *types.TransactionMetadata) {
	FcLog.Debug("signaled,id = %s", md.ID)
	//netLog.Error("comming an transaction:%s %s ",md.ID,md.Trx)
	//impl.dispatcher.bcastTransaction(md.PackedTrx)
}

func (impl *netPluginIMpl) AppliedTransaction(txn *types.TransactionTrace) {
	FcLog.Debug("signaled,id = %s", txn.ID)
}

func (impl *netPluginIMpl) AcceptedConfirmation(head *types.HeaderConfirmation) {
	FcLog.Debug("signaled,id = %s", head.BlockId)
}

func (impl *netPluginIMpl) TransactionAck(results common.Pair) {
	packedTrx, _ := results.Second.(*types.PackedTransaction)

	id := packedTrx.ID()
	if results.First != nil {
		FcLog.Info("signaled NACK, trx-id = %s :%s", id, results.First)
		impl.dispatcher.rejectedTransaction(id)
	} else {
		FcLog.Info("signaled ACK,trx-id = %s", id)
		impl.dispatcher.bcastTransaction(packedTrx)
	}
}

func (impl *netPluginIMpl) startMonitors() {
	impl.connectorCheck = NewDeadlineTimer(App().GetIoService())
	impl.transactionCheck = NewDeadlineTimer(App().GetIoService())

	impl.startConnTimer(impl.connectorPeriod, nil)
	impl.startTxnTimer()
}

func (impl *netPluginIMpl) startConnTimer(du time.Duration, fromConnection *Connection) {
	impl.connectorCheck.ExpiresFromNow(impl.connectorPeriod)
	impl.connectorCheck.AsyncWait(func(err error) {
		if err != nil {
			netLog.Error("Error from connection check monitor: %s", err.Error())
			impl.startConnTimer(impl.connectorPeriod, nil)
		} else {
			impl.connectionMonitor(fromConnection)
		}
	})
}

func (impl *netPluginIMpl) connectionMonitor(fromConnection *Connection) {
	maxTime := common.Now().AddUs(common.Milliseconds(int64(impl.maxCleanupTimeMs)))

	var i int
	var it *Connection
	if fromConnection != nil {
		i, it = impl.findConnection(fromConnection.peerAddr)
	} else {
		if len(impl.connections) > 0 {
			i, it = 0, impl.connections[0]
		} else {
			impl.startConnTimer(impl.connectorPeriod, nil)
			return
		}
	}

	for ; i < len(impl.connections); i++ {
		if common.Now().Sub(maxTime) >= 0 {
			impl.startConnTimer(time.Millisecond, it)
			return
		}
		if it.conn == nil && !it.connecting {
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
			netLog.Error("Error from transaction check monitor: %s", err.Error())
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
			netLog.Warn("Peer keep live ticked sooner than expected: %s", err)
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
		for _, pubKey := range impl.AllowedPeers {
			if pubKey == msg.Key {
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

	if msg.Signature.String() != ecc.NewSigNil().String() && msg.Token.Equals(crypto.NewSha256Nil()) {
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
		FcLog.Debug("Peer sent a handshake with blank signature and token,but this node accepts only authenticate connections.")
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
		for pubKey := range impl.privateKeys {
			return &pubKey
		}
	}
	return &ecc.PublicKey{}
}

//signCompact returns a signature of the digest using the corresponding private key of the signer.
//If there are no configured private keys, returns an empty signature.
func (impl *netPluginIMpl) signCompact(signer *ecc.PublicKey, digest *crypto.Sha256) *ecc.Signature {
	priKey, ok := impl.privateKeys[*signer]
	if ok {
		signature, err := priKey.Sign(digest.Bytes())
		if err != nil {
			netLog.Error("signCompact is error: %s", err.Error())
			return ecc.NewSigNil()
		}
		return &signature
	}
	pp := App().FindPlugin("ProducerPlugin").(*producer_plugin.ProducerPlugin)
	if pp != nil && pp.GetState() == Started {
		return pp.SignCompact(signer, *digest)
	}
	return ecc.NewSigNil()
}

func (impl *netPluginIMpl) handleChainSize(c *Connection, msg *ChainSizeMessage) {
	FcLog.Info("%s : receives chain_size_message %s", c.peerAddr, msg.String())
}

func (impl *netPluginIMpl) handleHandshake(c *Connection, msg *HandshakeMessage) {
	FcLog.Info("%s : receives handshake_message %s", c.peerAddr, msg.String())
	if !isValid(msg) {
		FcLog.Error("%s : bad handshake message", c.peerAddr)
		goAwayMsg := &GoAwayMessage{
			Reason: fatalOther,
			NodeID: crypto.NewSha256Nil(),
		}
		c.enqueue(goAwayMsg, true)
		return
	}

	cc := impl.ChainPlugin.Chain()
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
				NodeID: crypto.NewSha256Nil(),
			}
			c.enqueue(goAwayMsg, true)
		}

		if len(c.peerAddr) == 0 || c.lastHandshakeRecv.NodeID.Equals(crypto.NewSha256Nil()) {
			FcLog.Debug("checking for duplicate")
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
					FcLog.Debug("sending go_away duplicate to %s", msg.P2PAddress)
					goAwayMsg := &GoAwayMessage{
						Reason: duplicate,
						NodeID: c.nodeID,
					}
					c.enqueue(goAwayMsg, true)
					c.noRetry = duplicate
					return
				}
			}
		} else {
			FcLog.Debug("skipping dulicate check, addr ==%s, id == %s", c.peerAddr, c.lastHandshakeRecv.NodeID)
		}

		if !msg.ChainID.Equals(impl.chainID) {
			netLog.Error("Peer on different chain. Closing connection")
			goAwayMsg := &GoAwayMessage{
				Reason: wrongChain,
				NodeID: crypto.NewSha256Nil(),
			}
			c.enqueue(goAwayMsg, true)
			return
		}

		c.protocolVersion = impl.toProtocolVersion(msg.NetworkVersion)
		if c.protocolVersion != netVersion {
			if impl.networkVersionMatch {
				netLog.Error("Peer network version does not match expected %d but got %d", netVersion, c.protocolVersion)
				goAwayMsg := &GoAwayMessage{
					Reason: wrongVersion,
					NodeID: crypto.NewSha256Nil(),
				}
				c.enqueue(goAwayMsg, true)
				return
			} else {
				FcLog.Info("local network version: %d Remote version: %d", netVersion, c.protocolVersion)
			}
		}

		if !c.nodeID.Equals(msg.NodeID) {
			c.nodeID = msg.NodeID
		}

		if !impl.authenticatePeer(msg) {
			netLog.Error("Peer not authenticated. Closing connection")
			goAwayMsg := &GoAwayMessage{
				Reason: authentication,
				NodeID: crypto.NewSha256Nil(),
			}
			c.enqueue(goAwayMsg, true)
			return
		}

		onFork := false
		FcLog.Debug("lib_num = %d peer_lib = %d", libNum, peerLib)
		if peerLib <= libNum && peerLib > 0 {
			Try(func() {
				peerLibID := cc.GetBlockIdForNum(peerLib)
				onFork = !msg.LastIrreversibleBlockID.Equals(peerLibID)
			}).Catch(func(ex UnknownBlockException) {
				FcLog.Warn("peer last irreversible block %d is unknown", peerLib)
				onFork = true
			}).Catch(func(e interface{}) {
				netLog.Warn("caught an exception getting block id for %d", peerLib)
				onFork = true
			}).End()

			if onFork {
				netLog.Error("Peer chain is forked")
				goAwayMsg := &GoAwayMessage{
					Reason: forked,
					NodeID: crypto.NewSha256Nil(),
				}
				c.enqueue(goAwayMsg, true)
				return
			}
		}

		if c.sentHandshakeCount == 0 {
			c.sendHandshake()
		}
	}

	c.lastHandshakeRecv = msg
	impl.syncMaster.recvHandshake(c, msg)
}

func (impl *netPluginIMpl) handleGoaway(c *Connection, msg *GoAwayMessage) {
	rsn := ReasonStr[msg.Reason]
	FcLog.Info("%s : receive go_away_message reason = %s", c.peerAddr, rsn)
	c.noRetry = msg.Reason
	if msg.Reason == duplicate {
		c.nodeID = msg.NodeID
	}
	c.flushQueues()
	impl.close(c)
}

//handleTime processes time_message
//Calculate offset, delay and dispersion.  Note carefully the
//implied processing.  The first-order difference is done
//directly in 64-bit arithmetic, then the result is converted
//to floating double.  All further processing is in
//floating-double arithmetic with rounding done by the hardware.
//This is necessary in order to avoid overflow and preserve precision.
func (impl *netPluginIMpl) handleTime(c *Connection, msg *TimeMessage) {
	FcLog.Info("%s :received time_message %s", c.peerAddr, msg.String())
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

	//c.offset = float64((c.rec-c.org)+(msg.Xmt-c.dst)) / 2
	//NsecPerUsec := float64(1000)
	//FcLog.Info("Clock offset is %v ns  %v us", c.offset, c.offset/NsecPerUsec)

	c.org = 0
	c.rec = 0
}

func (impl *netPluginIMpl) handleNotice(c *Connection, msg *NoticeMessage) {
	// peer tells us about one or more blocks or txns. When done syncing, forward on
	// notices of previously unknown blocks or txns,
	FcLog.Info("%s : received notice_message %s", c.peerAddr, msg.String())

	c.connecting = false
	req := RequestMessage{}
	sendReq := false

	if msg.KnownTrx.Mode != none {
		FcLog.Debug("this is a %s notice with %d transactions", modeStr[msg.KnownTrx.Mode], msg.KnownTrx.Pending)
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
				ltx := impl.localTxns.GetById()
				for _, t := range ltx.Values() {
					req.ReqTrx.IDs = append(req.ReqTrx.IDs, t.ID)
				}
			}
		}
	case normal:
		impl.dispatcher.recvNotice(c, msg, false)
	}

	if msg.KnownBlocks.Mode != none {
		FcLog.Debug("this is a %s notice with  %d blocks", modeStr[msg.KnownBlocks.Mode], msg.KnownBlocks.Pending)
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
		FcLog.Error("bad notice_message : invalid known.mode %s", modeStr[msg.KnownBlocks.Mode])
	}

	FcLog.Debug("send req = %t", sendReq)
	if sendReq {
		c.enqueue(&req, true)
	}
}

func (impl *netPluginIMpl) handleRequest(c *Connection, msg *RequestMessage) {
	FcLog.Info("%s: received request_message %s", c.peerAddr, msg.String())

	switch msg.ReqBlocks.Mode {
	case catchUp:
		FcLog.Info("%s : received request_message:catch_up", c.peerAddr)
		c.blkSendBranch()
	case normal:
		FcLog.Info("%s : receive request_message:normal", c.peerAddr)
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

func (impl *netPluginIMpl) handleSyncRequest(c *Connection, msg *SyncRequestMessage) {
	FcLog.Info("%s : received sync_request_message %v", c.peerAddr, msg.String())

	if msg.EndBlock == 0 {
		c.peerRequested = &syncState{}
		c.flushQueues()
	} else {
		c.peerRequested = newSyncState(msg.StartBlock, msg.EndBlock, msg.StartBlock-1)
		c.enqueueSyncBlock()
	}
}

func (impl *netPluginIMpl) handlePackTransaction(c *Connection, msg *PackedTransactionMessage) {
	FcLog.Info(" %s received packed transaction %s", c.peerAddr, msg.String())

	cc := impl.ChainPlugin.Chain()
	if cc.GetReadMode() == chain.READONLY {
		FcLog.Debug("got a txn in read-only mode - dropping")
		return
	}

	if impl.syncMaster.isActive(c) {
		FcLog.Debug("got a txn during sync - dropping")
		return
	}
	tid := msg.ID()
	//c.cancelWait()
	if !impl.localTxns.GetById().Find(tid).IsEnd() {
		FcLog.Debug("got a duplicate transaction - dropping")
		return
	}

	impl.dispatcher.recvTransaction(c, tid)

	impl.ChainPlugin.AcceptTransaction(&msg.PackedTransaction, func(result common.StaticVariant) {
		if exception, ok := result.(Exception); ok {
			FcLog.Debug("bad packed_transaction : %s", exception.DetailMessage())
		} else {
			trace, _ := result.(types.TransactionTrace)
			if trace.Except == nil {
				FcLog.Debug("chain accepted transaction")
				impl.dispatcher.bcastTransaction(&msg.PackedTransaction)
				return
			}
			FcLog.Error("bad packed_transaction : %s", trace.Except.DetailMessage())
		}
		impl.dispatcher.rejectedTransaction(tid)
	})
}

func (impl *netPluginIMpl) handleSignedBlock(c *Connection, msg *SignedBlockMessage) {
	FcLog.Info("%s : receive signed_block message %d, %v", c.peerAddr, msg.BlockNumber(), msg.String())

	cc := impl.ChainPlugin.Chain()
	blkID := msg.BlockID()
	blkNum := msg.BlockNumber()
	//FcLog.Debug("canceling wait on %s", c.peerAddr)
	//c.cancelWait()

	returning := false
	Try(func() {
		if cc.FetchBlockById(blkID) != nil {
			impl.syncMaster.recvBlock(c, blkID, blkNum)
			returning = true
		}
	}).Catch(func(e interface{}) {
		netLog.Error("Caught an unknown exception trying to recall blockID")
	}).End()

	if returning {
		return
	}

	impl.dispatcher.recvBlock(c, blkID, blkNum)

	age := common.Now().Sub(msg.Timestamp.ToTimePoint())
	FcLog.Info("received signed_block : %d block age in secs = %d", blkNum, age.ToSeconds())

	reason := fatalOther
	Try(func() {
		impl.ChainPlugin.AcceptBlock(&msg.SignedBlock)
		reason = noReason
	}).Catch(func(ex UnlinkableBlockException) {
		FcLog.Error("bad signed_block : %s", ex.DetailMessage())
		reason = unlinkable
	}).Catch(func(ex BlockValidateException) {
		FcLog.Error("bad signed_block : %s", ex.DetailMessage())
		reason = validation
	}).Catch(func(ex AssertException) {
		FcLog.Error("bad signed_block : %s", ex.DetailMessage())
		netLog.Error("unable to accept block on assert exception %s from %s", ex.DetailMessage(), c.peerAddr)
	}).Catch(func(ex FcException) {
		FcLog.Error("bad signed_block : %s", ex.DetailMessage())
		netLog.Error("accept_block threw a non-assert exception %s from %s", ex.DetailMessage(), c.peerAddr)
		reason = noReason
	}).Catch(func(ex interface{}) {
		FcLog.Error("bad signed_block : unknown exception")
		netLog.Error("handle sync block caught something else from %s", c.peerAddr)
	}).End()

	if reason == noReason {
		var id common.TransactionIdType
		for _, recpt := range msg.Transactions {
			if recpt.Trx.TransactionID.Equals(crypto.NewSha256Nil()) {
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

func (impl *netPluginIMpl) toProtocolVersion(v uint16) uint16 {
	if v >= netVersionBase {
		v -= netVersionBase
		if v <= netVersionRange {
			return v
		}
	}
	return 0
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
