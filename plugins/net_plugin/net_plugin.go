package net_plugin

import (
	"bufio"
	"fmt"
	"net"

	"encoding/hex"
	"encoding/json"
	"github.com/eosspark/eos-go/exception"
	"time"

	"crypto/rand"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/log"
	"gopkg.in/urfave/cli.v1"
	"github.com/eosspark/eos-go/exception/try"
)

const (
	p2pChainIDString string = "cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"
)

type NetPlugin struct {
	my *netPluginIMpl
}

var netlog log.Logger

func init() {
	netlog = log.New("net_plugin")
	netlog.SetHandler(log.TerminalHandler)
}

func SetProgramOptions() {

}

func NewNetPlugin() *NetPlugin {
	np := new(NetPlugin)
	my := NewNetPluginIMpl()
	np.my = my

	return np
}

func (n *NetPlugin) NetPluginInitialize(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "p2p-listen-endpoint",
			Usage: "The actual host:port used to listen for incoming p2p connections.",
			Value: "0.0.0.0:9876",
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
			Value: &cli.StringSlice{"any"}, //TODO
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
		cli.BoolFlag{ //false
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
	}

	app.Action = func(c *cli.Context) {
		fmt.Println("Initialize net plugin")

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
		n.my.p2PAddress = c.String("p2p-server-address")
		n.my.suppliedPeers = c.StringSlice("p2p-peer-address")
		n.my.userAgentName = c.String("agent-name")

		allowecRemotes := c.StringSlice("allowed-connection")
		for _, allowedRemote := range allowecRemotes {
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
			try.EosAssert(c.IsSet("peer-key"), &exception.PluginConfigException{}, "At least one peer-key must accompany 'allowed-connection=specified'")
		}

		if c.IsSet("peer_key") {
			keyStrings := c.StringSlice("peer-key")
			for _, keyString := range keyStrings {
				pubKey, err := ecc.NewPublicKey(keyString)
				if err != nil {
					panic(err)
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
					panic(err)
				}
				prikey, err := ecc.NewPrivateKey(keyIdToWifPair[1])
				if err != nil {
					panic(err)
				}
				if prikey.PublicKey() != pubKey {
					panic(fmt.Errorf("the privateKey and PublicKey are not pairs"))
				}
				n.my.privateKeys[pubKey] = *prikey
			}
		}

		//	my->chain_plug = app().find_plugin<chain_plugin>();
		//	EOS_ASSERT( my->chain_plug, chain::missing_chain_plugin_exception, ""  );
		//	my->chain_id = app().get_plugin<chain_plugin>().get_chain_id();
		//fc::rand_pseudo_bytes( my->node_id.data(), my->node_id.data_size());
		//	ilog( "my node_id is ${id}", ("id", my->node_id));

		//n.my.keepAliceTimer.Reset(0)
		//n.my.tiker()

		cID, _ := hex.DecodeString(p2pChainIDString)
		cIdHash := *crypto.NewSha256Byte(cID)
		n.my.chainID = common.ChainIdType(cIdHash)

		nodeID := make([]byte, 32)
		rand.Read(nodeID)
		nodeIdHash := *crypto.NewSha256Byte(nodeID)
		n.my.nodeID = common.NodeIdType(nodeIdHash)

		fmt.Println("chainID: ", n.my.chainID)
		fmt.Println("nodeID: ", n.my.nodeID)

		n.my.peers = make(map[string]*Peer, 25)

		//np.my.keepAliceTimer.Reset(0)

	}

	n.my.loopWG.Add(1)
	go n.my.ticker()

}

func (np *NetPlugin) PluginStartup() {

	//ilog("starting listener, max clients is ${mc}",("mc",my->max_client_count));
	fmt.Printf("starting listener, max clients is %d\n", np.my.maxClientCount)

	np.my.loopWG.Add(3)
	go np.my.startListenLoop()
	go np.my.startConnTimer()
	go np.my.startTxnTimer()

	//chain::controller&cc = my->chain_plug->chain();
	//	{
	//		cc.accepted_block_header.connect( boost::bind(&net_plugin_impl::accepted_block_header, my.get(), _1));
	//		cc.accepted_block.connect(  boost::bind(&net_plugin_impl::accepted_block, my.get(), _1));
	//		cc.irreversible_block.connect( boost::bind(&net_plugin_impl::irreversible_block, my.get(), _1));
	//		cc.accepted_transaction.connect( boost::bind(&net_plugin_impl::accepted_transaction, my.get(), _1));
	//		cc.applied_transaction.connect( boost::bind(&net_plugin_impl::applied_transaction, my.get(), _1));
	//		cc.accepted_confirmation.connect( boost::bind(&net_plugin_impl::accepted_confirmation, my.get(), _1));
	//	}
	//	my->incoming_transaction_ack_subscription = app().get_channel<channels::transaction_ack>().subscribe(boost::bind(&net_plugin_impl::transaction_ack, my.get(), _1));
	//	if( cc.get_read_mode() == chain::db_read_mode::READ_ONLY ) {
	//	my->max_nodes_per_host = 0;
	//	ilog( "node in read-only mode setting max_nodes_per_host to 0 to prevent connections" );
	//	}

	for _, seedNode := range np.my.suppliedPeers {
		re := np.connect(seedNode)
		if re != "added connection" {
			fmt.Println(re)
		}
	}
	fmt.Println("******************** go *******************")

	np.my.loopWG.Wait()

}

func (np *NetPlugin) PluginShutDown() {
	//ilog( "shutdown.." )
	fmt.Println("shutdown...")
	np.my.done = true
	//if( my->acceptor ) {
	//	ilog( "close acceptor" );
	//	my->acceptor->close();
	//
	//	ilog( "close ${s} connections",( "s",my->connections.size()) );
	//	auto cons = my->connections;
	//	for( auto con : cons ) {
	//	my->close( con);
	//	}

	//ilog( "exit shutdown" )
	fmt.Println("exit shutdown")
}

//func (np *NetPlugin) numPeers() int {
//	return np.my.countOpenSockets()
//}

//connect used to trigger a new connetion RPC API
func (np *NetPlugin) connect(host string) string {
	_, ok := np.my.peers[host]
	if ok {
		return "already connected"
	}

	con, err := net.Dial("tcp", host)
	if err != nil {
		return err.Error()
	}

	np.my.peers[host] = NewPeer(con, bufio.NewReader(con))
	////fc_dlog(logger,"adding new connection to the list")
	fmt.Println("connecting to: ", con.RemoteAddr(), "adding new peer to the list")
	np.my.loopWG.Add(1)
	go np.my.peers[host].read(np.my)

	return "added connection"

}

func (np *NetPlugin) disconnect(host string) string {
	for name, peer := range np.my.peers {
		if name == host {
			peer.connection.Close()
			delete(np.my.peers, host)
			return "connection removed"
		}
	}
	return "no known connection for host"
}

func (np *NetPlugin) status(host string) PeerStatus {
	con, ok := np.my.peers[host]
	if ok {
		return *con.getStatus()
	}
	return PeerStatus{}
}

func (np *NetPlugin) connections() []PeerStatus {
	result := make([]PeerStatus, len(np.my.peers))
	for _, c := range np.my.peers {
		result = append(result, *c.getStatus())
	}
	return result
}
