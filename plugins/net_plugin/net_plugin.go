package net_plugin

import (
	"bufio"
	//"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"net"
	"time"
)

const (
	p2pListenEndPoint string = "127.0.0.1:8000" //TODO for testing
	p2pNodeIDString   string = "f1259a544acbe6fbaa3d13965e1b767991c9d444e3bead117ce01a0d5c96e1ef"
	p2pChainIDString  string = "cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"
)

var p2pAddress = []string{
	//"127.0.0.1:9876",
	//"127.0.0.1:7777",
	//"127.0.0.1:7778",
	//"127.0.0.1:7779",
}

var (
	p2pListenAddress = flag.String("p2p-listen-endpoint", "0.0.0.0:9876",
		"The actual host:port used to listen for incoming p2p connections.")
	p2pServerAddress = flag.String("p2p-server-address", "",
		"An externally accessible host:port for identifying this node. Defaults to p2p-listen-endpoint.")
	p2pPeerAddress = flag.String("p2p-peer-address", "",
		"The public endpoint of a peer node to connect to. Use multiple p2p-peer-address options as needed to compose a network.")
	p2pMaxNodesPerHost = flag.Int("p2p-max-nodes-per-host", defMaxNodesPerHost,
		"Maximum number of client nodes from any single IP address")
	agentName = flag.String("agent-name", "\"EOS Test Agent\"",
		"The name supplied to identify this node amongst the peers.")
	allowedConnection = flag.String("allowed-connection", "any",
		"Can be 'any' or 'producers' or 'specified' or 'none'. If 'specified', "+"peer-key must be specified at least once. "+
			"If only 'producers', peer-key is not required. 'producers' and 'specified' may be combined.") //TODO vector<string>??
	peerKey                 = flag.String("peer_key", "", "Optional public key of peer allowed to connect.  May be used multiple times.") //TODO  multitoken
	peerPrivateKey          = flag.String("peer_private_key", "", "Tuple of [PublicKey, WIF private key] (may specify multiple times)")   //TODO  multitoken
	maxClients              = flag.Int("max-clients", defMaxClients, "Maximum number of clients from which connections are accepted, use 0 for no limit")
	connectionCleanupPeriod = flag.Int("connection-cleanup-period", defConnRetryWait, "number of seconds to wait before cleaning up dead connections")
	maxCleanupTimeMsec      = flag.Int("max-cleanup-time-msec", 10, "max connection cleanup time per cleanup call in millisec")
	networkVersionMatch     = flag.Bool("network-version-match", false, "True to require exact match of peer network version.")
	syncFetchSpan           = flag.Uint("sync-fetch-span", defSyncFetchSpan, "number of blocks to retrieve in a chunk from any individual peer during synchronization")
	maxImplicitRequest      = flag.Uint("max-implicit-request", uint(defMaxJustSend), "maximum sizes of transaction or block messages that are sent without first sending a notice")
	useSocketReadWatermark  = flag.Bool("use-socket-read-watermark", false, "Enable expirimental socket read watermark optimization")
	peerLogFormat           = flag.String("peer-log-format", "[\"${_name}\" ${_ip}:${_port}]",
		"The string used to format peers when logging messages about them.  Variables are escaped with ${<variable name>}.\n"+
			"Available Variables:\n"+
			"   _name  \tself-reported name\n\n"+
			"   _id    \tself-reported ID (64 hex characters)\n\n"+
			"   _sid   \tfirst 8 characters of _peer.id\n\n"+
			"   _ip    \tremote IP address of peer\n\n"+
			"   _port  \tremote port number of peer\n\n"+
			"   _lip   \tlocal IP address connected to peer\n\n"+
			"   _lport \tlocal port number connected to peer\n\n")
)

type NetPlugin struct {
	my *netPluginIMpl
}

func SetProgramOptions() {

}

func NewNetPlugin() *NetPlugin {
	fmt.Println("Initialize net plugin")
	impl := NewNetPluginIMpl()

	impl.networkVersionMatch = *networkVersionMatch
	//impl.syncMaster.reset()
	//impl.dispatcher.reset()

	impl.connectorPeriod = time.Duration(*connectionCleanupPeriod) * time.Second //*connectionCleanupPeriod*time.Second //TODO
	impl.maxCleanupTimeMs = *maxCleanupTimeMsec
	impl.txnExpPeriod = defTxnExpireWait
	impl.respExpectedPeriod = defRespExpectedWait
	//impl.dispatcher.justSendItMax = maxImplicitRequest
	impl.maxClientCount = uint32(*maxClients)
	impl.maxNodesPerHost = uint32(*p2pMaxNodesPerHost)
	impl.numClients = 0
	impl.useSocketReadWatermark = *useSocketReadWatermark
	//impl.resolver =

	//impl.ListenEndpoint = *p2pListenAddress
	impl.ListenEndpoint = p2pListenEndPoint
	impl.p2PAddress = *p2pServerAddress
	impl.suppliedPeers = []string{*p2pPeerAddress}

	//if( options.count( "p2p-peer-address" )) {
	//	my->supplied_peers = options.at( "p2p-peer-address" ).as<vector<string> >();
	//}
	impl.userAgentName = *agentName

	//allowecRemotes := *allowedConnection
	//for _,allowedRemote := range allowecRemotes{
	//	switch allowecRemote{
	//	case "any":
	//		impl.allowedConnections |= anyPossible
	//	case "producers":
	//		impl.allowedConnections |= producersPossible
	//	case "specified":
	//		impl.allowedConnections |= specifiedPossible
	//	case "none":
	//		impl.allowedConnections |= nonePossible
	//	}
	//}

	if impl.allowedConnections&specifiedPossible != 0 {
		//EOS_ASSERT( options.count( "peer-key" ),
		//	plugin_config_exception,
		//	"At least one peer-key must accompany 'allowed-connection=specified'" );
	}

	//keyStrings := *peerKey
	//for _,keyString := range keyStrings{
	//	pubKey,err := ecc.NewPublicKey(keyString)
	//	if err !=nil{
	//		panic(err)
	//	}
	//	impl.AllowedPeers = append(impl.AllowedPeers,pubKey)
	//}

	//keyIdToWifPairStrings := *peerPrivateKey
	//for _,keyIdToWifPairString := range keyIdToWifPairStrings{
	//	keyIdToWifPair := dejsonify<std::pair<chain::public_key_type, std::string>>(
	//		key_id_to_wif_pair_string )
	//
	//	impl.privateKeys[keyIdToWifPair.first] =  keyIdToWifPair.second
	//}

	//	my->chain_plug = app().find_plugin<chain_plugin>();
	//	EOS_ASSERT( my->chain_plug, chain::missing_chain_plugin_exception, ""  );
	//	my->chain_id = app().get_plugin<chain_plugin>().get_chain_id();
	//fc::rand_pseudo_bytes( my->node_id.data(), my->node_id.data_size());
	//	ilog( "my node_id is ${id}", ("id", my->node_id));

	//impl.keepAliceTimer.Reset(0)
	//impl.tiker()

	cID, _ := hex.DecodeString(p2pChainIDString)
	cIdHash := *crypto.NewSha256Byte(cID)
	impl.chainID = common.ChainIdType(cIdHash)

	//nodeID := make([]byte, 32)
	//rand.Read(nodeID)
	//nodeIdHash := *crypto.NewSha256Byte(nodeID)
	//impl.nodeID = common.NodeIdType(nodeIdHash)

	nodeID, _ := hex.DecodeString(p2pNodeIDString)
	nodeIdHash := *crypto.NewSha256Byte(nodeID)
	impl.nodeID = common.NodeIdType(nodeIdHash)
	fmt.Println("chainID: ", impl.chainID)
	fmt.Println("nodeID: ", impl.nodeID)

	impl.peers = make(map[string]*Peer, 25)
	np := new(NetPlugin)
	np.my = impl
	//np.my.keepAliceTimer.Reset(0)
	go np.my.ticker()

	return np
}

func (np *NetPlugin) PluginStartup() {

	//ilog("starting listener, max clients is ${mc}",("mc",my->max_client_count));
	fmt.Printf("starting listener, max clients is %d\n", np.my.maxClientCount)

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

	np.my.suppliedPeers = p2pAddress //TODO for testing

	for _, seedNode := range np.my.suppliedPeers {
		np.connect(seedNode)
		fmt.Println(len(np.my.suppliedPeers))

	}
	fmt.Println("******************** go *******************")

	for {

	}

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
