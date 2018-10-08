package p2p

//import (
//	"fmt"
//	"github.com/eosspark/eos-go/chain/types"
//	"github.com/eosspark/eos-go/common"
//	"github.com/eosspark/eos-go/crypto/ecc"
//	"github.com/eosspark/eos-go/crypto/rlp"
//	"net"
//	"time"
//)
//
//type possibleConnections byte
//
//const (
//	nonePossible      possibleConnections = 0
//	producersPossible possibleConnections = 1 << 0
//	specifiedPossible possibleConnections = 1 << 1
//	anyPossible       possibleConnections = 1 << 2
//)
//
//type netPluginIMpl struct {
//	//conn               *net.Conn
//	//ListenEndpoint     net.Addr
//
//	p2PAddress         string
//	maxClientCount     uint32
//	maxNodesPerHost    uint32
//	numClients         uint32
//	suppliedPeers      []string
//	AllowedPeers       []ecc.PublicKey                  //< peer keys allowed to connect
//	privateKeys        map[ecc.PublicKey]ecc.PrivateKey //< overlapping with producer keys, also authenticating non-producing nodes
//	allowedConnections possibleConnections
//	done               bool
//	connectorCheck     time.Timer
//	transactionCheck   time.Timer
//	keepAliceTimer     time.Timer
//
//	connectorPeriod            time.Duration
//	txnExpPeriod               time.Duration
//	respExpectedPeriod         time.Duration
//	keepaliveINterval          time.Duration //32*time.Sencond
//	peerAuthenticationInterval time.Duration //< Peer clock may be no more than 1 second skewed from our clock, including network latency.
//
//	maxCleanupTimeMs    int
//	networkVersionMatch bool
//	chainID             common.ChainIdType
//	nodeID              common.NodeIdType
//	userAgentName       string
//	startedSessions     int
//
//	//LocalTxns           NodeTransactionIndex
//
//	//Connections connection_ptr
//	//SyncMaster SyncManager
//	//Dispatcher DispatchManager
//	//ChainPlugin *ChainPlugin
//	//Resolver tcp::resolver
//
//	//incomingTransactionAckSubscription chan //incomingTransactionAckSubscription  channel_type::handle
//}
//
//func NewNetPluginIMpl() *netPluginIMpl {
//	return &netPluginIMpl{
//		maxClientCount:             0,
//		maxNodesPerHost:            1,
//		numClients:                 0,
//		allowedConnections:         nonePossible,
//		done:                       false,
//		keepaliveINterval:          32 * time.Second,
//		peerAuthenticationInterval: 1 * time.Second,
//		maxCleanupTimeMs:           0,
//		networkVersionMatch:        false,
//		startedSessions:            0,
//		txnExpPeriod:               defTxnExpireWait,
//	}
//}
//
//func (impl *netPluginIMpl) findConnection(host string) {
//
//}
//
//func (impl *netPluginIMpl) connect(c connectionPtr) {
//
//}
//
//func (impl *netPluginIMpl) connect(c connectionPtr, endpointItr net.Conn) {
//
//}
//
//func (impl *netPluginIMpl) startSession(c connectionPtr) bool {
//
//}
//
//func (impl *netPluginIMpl) startListenLoop() {
//
//}
//
//func (impl *netPluginIMpl) startReadMessage(c connectionPtr) {
//
//}
//
//func (impl *netPluginIMpl) close(c connectionPtr) {
//
//}
//
//func (impl *netPluginIMpl) countOpenSockets() int {
//
//	return 0
//}
//
//func (impl *netPluginIMpl) sendAll(msg *P2PMessage, verify func()) {
//
//}
//
//func (impl *netPluginIMpl) AcceptedBlockHeader(block types.BlockState) {
//	//fc_dlog(logger,"signaled, id = ${id}",("id", block->id))
//	fmt.Printf("signed,id =%v", block.ID)
//}
//
//func (impl *netPluginIMpl) AcceptedBlock(block types.BlockState) {
//	//fc_dlog(logger,"signaled, id = ${id}",("id", block.ID))
//	fmt.Printf("signaled,id = %v\n", block.ID)
//	//dispatcher.bcast_block(*block.SignedBlock)
//	//wlog("广播signed block：${block}",("block",*block.SignedBlock))
//}
//
//func (impl *netPluginIMpl) IrreversibleBlock(block types.BlockState) {
//	//fc_dlog(logger,"signaled, id = ${id}",("id", block.ID))
//}
//
//func (impl *netPluginIMpl) AcceptedTransaction(md *TransactionMetadataPtr) {
//	//fc_dlog(logger,"signaled, id = ${id}",("id", md.id))
//	//      dispatcher.bcast_transaction(md.packed_trx)
//}
//
//func (impl *netPluginIMpl) AppliedTransaction(txn *TransactionTracePtr) {
//	//fc_dlog(logger,"signaled, id = ${id}",("id", txn.id))
//}
//
//func (impl *netPluginIMpl) AcceptedConfirmation(head *types.HeaderConfirmation) {
//	//fc_dlog(logger,"signaled, id = ${id}",("id", head.BlockId))
//}
//
//func (impl *netPluginIMpl) TransactionAck(results common.Tuple) {
//	id := results[1].(common.TransactionIdType)
//	if results[0] != nil {
//		//fc_ilog(logger,"signaled NACK, trx-id = ${id} : ${why}",("id", id)("why", results.first->to_detail_string()));
//		//dispatche.ejected_transaction(id)
//		fmt.Println(id)
//	} else {
//		//fc_ilog(logger,"signaled ACK, trx-id = ${id}",("id", id));
//		//dispatcher->bcast_transaction(*results[1])
//		//elog("广播transactoin: ${sig}",("sig",*results.second));
//	}
//}
//
//func (impl *netPluginIMpl) startConnTimer(du time.Duration, fromConnection Connection) {
//
//}
//
////void net_plugin_impl::start_conn_timer(boost::asio::steady_timer::duration du, std::weak_ptr<connection> from_connection) {
////connector_check->expires_from_now( du);
////connector_check->async_wait( [this, from_connection](boost::system::error_code ec) {
////if( !ec) {
////connection_monitor(from_connection);
////}
////else {
////elog( "Error from connection check monitor: ${m}",( "m", ec.message()));
////start_conn_timer( connector_period, std::weak_ptr<connection>());
////}
////});
////}
//
//func (impl *netPluginIMpl) startTxnTimer() {
//
//	for {
//		select {
//		case <-time.After(txnExpPeriod):
//
//		}
//	}
//}
//
////void net_plugin_impl::start_txn_timer() {
////transaction_check->expires_from_now( txn_exp_period);
////transaction_check->async_wait( [this](boost::system::error_code ec) {
////if( !ec) {
////expire_txns( );
////}
////else {
////elog( "Error from transaction check monitor: ${m}",( "m", ec.message()));
////start_txn_timer( );
////}
////});
////}
//
//func startMonitors() {
//
//}
//
////void net_plugin_impl::start_monitors() {
////connector_check.reset(new boost::asio::steady_timer( app().get_io_service()));
////transaction_check.reset(new boost::asio::steady_timer( app().get_io_service()));
////start_conn_timer(connector_period, std::weak_ptr<connection>());
////start_txn_timer();
////}
//
//func expireTxns() {
//
//}
//
//func connectionMonitor(fromConnection *Connection) {
//
//}
//
//// tiker Peer heartbeat
//func (impl *netPluginIMpl) tiker() {
//
//}
//
////void net_plugin_impl::ticker() {
////keepalive_timer->expires_from_now (keepalive_interval);
////keepalive_timer->async_wait ([this](boost::system::error_code ec) {
////ticker ();
////if (ec) {
////wlog ("Peer keepalive ticked sooner than expected: ${m}", ("m", ec.message()));
////}
////for (auto &c : connections ) {
////if (c->socket->is_open()) {
////c->send_time();
////}
////}
////});
////}
//
//// authenticatePeer determine if a peer is allowed to connect.
//// Checks current connection mode and key authentication.
//// return False if the peer should not connect, True otherwise.
//func (n *netPluginIMpl) authenticatePeer(msg *HandshakeMessage) bool {
//	var allowedIt, privateIt ,foundProducerKey bool
//
//	if n.allowedConnections == nonePossible {
//		return false
//	}
//	if n.allowedConnections == anyPossible {
//		return true
//	}
//	if n.allowedConnections&(producersPossible|specifiedPossible) != 0 {
//		for _, pubkey := range n.AllowedPeers {
//			if pubkey == msg.Key {
//				allowedIt = true
//			}
//		}
//		_,privateIt = n.privateKeys[msg.Key]
//
//		//producer_plugin* pp = app().find_plugin<producer_plugin>();
//		//if(pp != nullptr)
//		//found_producer_key = pp->is_producer_key(msg.key);
//
//		if allowedIt && privateIt && !foundProducerKey {
//			//elog( "Peer ${peer} sent a handshake with an unauthorized key: ${key}.",
//			//	("peer", msg.p2p_address)("key", msg.key));
//
//			return false
//		}
//	}
//	msgTime := msg.Time
//	t :=common.Now()
//	if time.Duration(uint64((t -msgTime))*time.Microsecond) > n.peerAuthenticationInterval {
//
//	}
//
//
//}
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
////
////
////namespace sc = std::chrono;
////sc::system_clock::duration msg_time(msg.time);
////auto time = sc::system_clock::now().time_since_epoch();
////if(time - msg_time > peer_authentication_interval) {
////elog( "Peer ${peer} sent a handshake with a timestamp skewed by more than ${time}.",
////("peer", msg.p2p_address)("time", "1 second")); // TODO Add to_variant for std::chrono::system_clock::duration
////return false;
////}
////
////if(msg.sig != chain::signature_type() && msg.token != sha256()) {
////sha256 hash = fc::sha256::hash(msg.time);
////if(hash != msg.token) {
////elog( "Peer ${peer} sent a handshake with an invalid token.",
////("peer", msg.p2p_address));
////return false;
////}
////chain::public_key_type peer_key;
////try {
////peer_key = crypto::public_key(msg.sig, msg.token, true);
////}
////catch (fc::exception& /*e*/) {
////elog( "Peer ${peer} sent a handshake with an unrecoverable key.",
////("peer", msg.p2p_address));
////return false;
////}
////if((allowed_connections & (Producers | Specified)) && peer_key != msg.key) {
////elog( "Peer ${peer} sent a handshake with an unauthenticated key.",
////("peer", msg.p2p_address));
////return false;
////}
////}
////else if(allowed_connections & (Producers | Specified)) {
////dlog( "Peer sent a handshake with blank signature and token, but this node accepts only authenticated connections.");
////return false;
////}
////return true;
////}
//
//
//// getAuthenticationKey retrieve public key used to authenticate with peers.
//// Finds a key to use for authentication.  If this node is a producer, use
//// the front of the producer key map.  If the node is not a producer but has
//// a configured private key, use it.  If the node is neither a producer nor has
//// a private key, returns an empty key.
//// On a node with multiple private keys configured, the key with the first
//// numerically smaller byte will always be used.
//func (n *netPluginIMpl) getAuthenticationKey() *ecc.PublicKey {
//	if len(n.privateKeys) > 0 {
//		for pubkey, _ := range n.privateKeys { //TODO easier  ？？？
//			return &pubkey
//		}
//		/*producer_plugin* pp = app().find_plugin<producer_plugin>(); //TODO EOSIO not used
//		if(pp != nullptr && pp->get_state() == abstract_plugin::started)
//		   return pp->first_producer_public_key();*/
//		return &ecc.PublicKey{}
//	}
//
//	return &ecc.PublicKey{}
//}
//
//// signCompact returns a signature of the digest using the corresponding private key of the signer.
//// If there are no configured private keys, returns an empty signature.
//func (n *netPluginIMpl) signCompact(signer *ecc.PublicKey, digest *rlp.Sha256) *ecc.Signature {
//	privateKeyPtr, ok := n.privateKeys[*signer]
//	if ok {
//		signature, err := privateKeyPtr.Sign(digest.Bytes())
//		if err != nil {
//			panic(err)
//		}
//		return &signature
//	} else { //TODO producer_plugin
//		//producerPlugin
//		//
//		//return pp.signCompact(signer,digest)
//
//		//producer_plugin* pp = app().find_plugin<producer_plugin>();
//		//if(pp != nullptr && pp->get_state() == abstract_plugin::started)
//		//return pp->sign_compact(signer, digest);
//	}
//	return &ecc.Signature{}
//}
//
//func (n *netPluginIMpl) toProtocolVersion(v uint16) uint16 {
//	if v > netVersionBase {
//		v -= netVersionBase
//		if v <= netVersionRange {
//			return v
//		}
//	}
//	return 0
//}
