package net_plugin

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"github.com/eosspark/eos-go/plugins/chain_interface"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/urfave/cli"
	"net"
	"time"
)

const NetPlug = PluginTypeName("NetPlugin")

var netPlugin Plugin = App().RegisterPlugin(NetPlug, NewNetPlugin(App().GetIoService()))

type NetPlugin struct {
	AbstractPlugin
	my *netPluginIMpl
}

func NewNetPlugin(io *asio.IoContext) *NetPlugin {
	plugin := &NetPlugin{}

	plugin.my = NewNetPluginIMpl()
	plugin.my.Self = plugin

	return plugin
}

func (n *NetPlugin) SetProgramOptions(options *[]cli.Flag) {
	*options = append(*options,
		cli.StringFlag{
			Name:  "p2p-listen-endpoint",
			Usage: "The actual host:port used to listen for incoming p2p connections.",
			Value: "0.0.0.0:9800",
		},
		cli.StringFlag{
			Name:  "p2p-server-address",
			Usage: "An externally accessible host:port for identifying this node. Defaults to p2p-listen-endpoint.",
			Value: "",
		},
		cli.StringSliceFlag{
			Name:  "p2p-peer-address",
			Usage: "The public endpoint of a peer node to connect to. Use multiple p2p-peer-address options as needed to compose a network.",
		},
		cli.IntFlag{
			Name:  "p2p-max-nodes-per-host",
			Usage: "Maximum number of client nodes from any single IP address",
			Value: defMaxNodesPerHost,
		},
		cli.StringFlag{
			Name:  "agent-name",
			Usage: "The name supplied to identify this node amongst the peers.",
			Value: "EOS Test Agent",
		},
		cli.StringSliceFlag{
			Name: "allowed-connection",
			Usage: "Can be 'any' or 'producers' or 'specified' or 'none'. If 'specified', " + "peer-key must be specified at least once. " +
				"If only 'producers', peer-key is not required. 'producers' and 'specified' may be combined.",
			Value: &cli.StringSlice{"any"},
		},
		cli.StringSliceFlag{
			Name:  "peer_key",
			Usage: "Optional public key of peer allowed to connect.  May be used multiple times.",
		},
		cli.StringSliceFlag{
			Name:  "peer-private-key",
			Usage: "Tuple of [PublicKey, WIF private key] (may specify multiple times)",
		},
		cli.IntFlag{
			Name:  "max-clients",
			Usage: "Maximum number of clients from which connections are accepted, use 0 for no limit",
			Value: defMaxClients,
		},
		cli.IntFlag{
			Name:  "connection-cleanup-period",
			Usage: "number of seconds to wait before cleaning up dead connections",
			Value: defConnRetryWait,
		},
		cli.IntFlag{
			Name:  "max-cleanup-time-msec",
			Usage: "max connection cleanup time per cleanup call in millisec",
			Value: 10,
		},
		cli.BoolFlag{
			Name:  "network-version-match",
			Usage: "True to require exact match of peer network version.",
		},
		cli.UintFlag{
			Name:  "sync-fetch-span",
			Usage: "number of blocks to retrieve in a chunk from any individual peer during synchronization",
			Value: defSyncFetchSpan,
		},
		cli.UintFlag{
			Name:  "max-implicit-request",
			Usage: "maximum sizes of transaction or block messages that are sent without first sending a notice",
			Value: uint(defMaxJustSend),
		},
		cli.BoolFlag{ //false
			Name:  "use-socket-read-watermark",
			Usage: "Enable expirimental socket read watermark optimization",
		},
		cli.StringFlag{
			Name: "peer-log-format",
			Usage: "The string used to format peers when logging messages about them.  Variables are escaped with ${<variable name>}.\n" +
				"Available Variables:\n" +
				"   _name  \tself-reported name\n\n" +
				"   _id    \tself-reported ID (64 hex characters)\n\n" +
				"   _sid   \tfirst 8 characters of _peer.id\n\n" +
				"   _ip    \tremote IP address of peer\n\n" +
				"   _port  \tremote port number of peer\n\n" +
				"   _lip   \tlocal IP address connected to peer\n\n" +
				"   _lport \tlocal port number connected to peer\n\n",
			Value: "[\"${_name}\" ${_ip}:${_port}]",
		},
	)
}
func (n *NetPlugin) PluginInitialize(c *cli.Context) {
	Try(func() {
		netLog.Info("Initialize net plugin")

		n.my.networkVersionMatch = c.Bool("network-version-match")
		n.my.connectorPeriod = time.Duration(c.Int("connection-cleanup-period")) * time.Second
		n.my.maxCleanupTimeMs = c.Int("max-cleanup-time-msec")
		n.my.txnExpPeriod = defTxnExpireWait
		n.my.respExpectedPeriod = defRespExpectedWait
		n.my.dispatcher.justSendItMax = uint32(c.Int("max-implicit-request"))
		n.my.maxClientCount = uint32(c.Int("max-clients"))
		n.my.maxNodesPerHost = uint32(c.Int("p2p-max-nodes-per-host"))
		n.my.numClients = 0
		n.my.useSocketReadWatermark = c.Bool("use-socket-read-watermark")
		n.my.ListenEndpoint = c.String("p2p-listen-endpoint")
		n.my.p2PAddress = c.String("p2p-listen-endpoint")
		n.my.suppliedPeers = c.StringSlice("p2p-peer-address")
		n.my.userAgentName = c.String("agent-name")

		allowedRemotes := c.StringSlice("allowed-connection")
		for _, allowedRemote := range allowedRemotes {
			switch allowedRemote {
			case "any":
				n.my.allowedConnections |= anyPossible
			case "producers":
				n.my.allowedConnections |= producersPossible
			case "specified":
				n.my.allowedConnections |= specifiedPossible
			case "none":
				n.my.allowedConnections |= nonePossible
			}
		}

		if n.my.allowedConnections&specifiedPossible != 0 {
			EosAssert(c.IsSet("peer-key"), &exception.PluginConfigException{}, "At least one peer-key must accompany 'allowed-connection=specified'")
		}

		if c.IsSet("peer_key") {
			keyStrings := c.StringSlice("peer-key")
			for _, keyString := range keyStrings {
				pubKey, err := ecc.NewPublicKey(keyString)
				if err != nil {
					Throw(err)
				}
				n.my.AllowedPeers = append(n.my.AllowedPeers, pubKey)
			}
		}
		if c.IsSet("peer-private-key") {
			keyIdToWifPairStrings := c.StringSlice("peer-private-key")
			var keyIdToWifPair []string
			for _, keyIdToWifPairString := range keyIdToWifPairStrings {
				json.Unmarshal([]byte(keyIdToWifPairString), &keyIdToWifPair)
				pubKey, err := ecc.NewPublicKey(keyIdToWifPair[0])
				if err != nil {
					Throw(err)
				}
				prikey, err := ecc.NewPrivateKey(keyIdToWifPair[1])
				if err != nil {
					Throw(err)
				}
				if prikey.PublicKey() != pubKey {
					Throw(fmt.Errorf("the privateKey and PublicKey are not pairs"))
				}
				n.my.privateKeys[pubKey] = *prikey
			}
		}

		n.my.ChainPlugin = App().FindPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin)
		EosAssert(n.my.ChainPlugin != nil, &exception.MissingChainPluginException{}, "")
		n.my.chainID = n.my.ChainPlugin.GetChainId()

		nodeID := make([]byte, 32)
		rand.Read(nodeID)
		nodeIdHash := *crypto.NewSha256Byte(nodeID)
		n.my.nodeID = common.NodeIdType(nodeIdHash)
		netLog.Info("my node_id is %s", n.my.nodeID)
		n.my.connections = make([]*Connection, 0)

		n.my.keepAliceTimer = asio.NewDeadlineTimer(App().GetIoService())
		n.my.ticker()

	}).FcLogAndRethrow().End()
}

func (n *NetPlugin) PluginStartup() {
	netLog.Info("starting listener, max clients is %d", n.my.maxClientCount)

	var err error
	n.my.Listener, err = net.Listen("tcp", n.my.ListenEndpoint)
	if err != nil {
		netLog.Error("Error getting remote endpoint:", err)
	}
	netLog.Info("Listening on: %s", n.my.ListenEndpoint)

	n.my.startListenLoop()
	n.my.startMonitors()

	cc := n.my.ChainPlugin.Chain()
	{
		cc.AcceptedBlockHeader.Connect(&chain_interface.AcceptedBlockHeaderCaller{Caller: n.my.AcceptedBlockHeader})
		cc.AcceptedBlock.Connect(&chain_interface.AcceptedBlockCaller{Caller: n.my.AcceptedBlock})
		cc.IrreversibleBlock.Connect(&chain_interface.IrreversibleBlockCaller{Caller: n.my.IrreversibleBlock})
		cc.AcceptedTransaction.Connect(&chain_interface.AcceptedTransactionCaller{Caller: n.my.AcceptedTransaction})
		cc.AppliedTransaction.Connect(&chain_interface.AppliedTransactionCaller{Caller: n.my.AppliedTransaction})
		cc.AcceptedConfirmation.Connect(&chain_interface.AcceptedConfirmationCaller{Caller: n.my.AcceptedConfirmation})
	}
	App().GetChannel(chain_interface.TransactionAck).Subscribe(&chain_interface.TransactionAckCaller{Caller: n.my.TransactionAck})

	if cc.GetReadMode() == chain.READONLY {
		n.my.maxNodesPerHost = 0
		netLog.Info("node in read-only mode setting max_nodes_per_host to 0 to prevent connections")
	}

	for _, seedNode := range n.my.suppliedPeers {
		re := n.Connect(seedNode)
		if re != "added connection" {
			netLog.Error(re)
		}
	}
}

func (n *NetPlugin) PluginShutdown() {
	Try(func() {
		netLog.Info("shutdown...")
		if n.my.Listener != nil {
			netLog.Info("close acceptor")
			n.my.Listener.Close()
		}
		netLog.Info("close %d connections", len(n.my.connections))
		peers := n.my.connections
		for _, p := range peers {
			n.my.close(p)
		}
		netLog.Info("exit shutdown")
	}).FcCaptureAndRethrow().End()
}

//Connect used to trigger a new connection RPC API
func (n *NetPlugin) Connect(host string) string {
	i, _ := n.my.findConnection(host)
	if i >= 0 {
		return "already connected"
	}

	c := NewConnectionByEndPoint(host, n.my)
	netLog.Info("adding new peer to the list") //FC
	n.my.connections = append(n.my.connections, c)
	netLog.Info("calling active connector") //FC
	n.my.connect(c)
	return "added connection"
}

func (n *NetPlugin) Disconnect(host string) string {
	for _, con := range n.my.connections {
		if con.peerAddr == host {
			con.reset()
			n.my.close(con)
			n.my.eraseConnection(con)
			return "connection removed"
		}
	}
	return "no known connection for host"
}

func (n *NetPlugin) Status(host string) PeerStatus {
	i, con := n.my.findConnection(host)
	if i >= 0 {
		return *con.getStatus()
	}
	return PeerStatus{}
}

func (n *NetPlugin) Connections() []PeerStatus {
	result := make([]PeerStatus, len(n.my.connections))
	for _, c := range n.my.connections {
		result = append(result, *c.getStatus())
	}
	return result
}
