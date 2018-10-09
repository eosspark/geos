package p2p

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"

	"time"
)

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
)

/*
   Index by id
   Index by is_known, block_num, validated_time, this is the order we will broadcast
   to peer.
   Index by is_noticed, validated_time
*/
type transactionState struct {
	id              common.TransactionIdType
	isKnownByPeer   bool //true if we sent or received this trx to this peer or received notice from peer
	isNoticedToPeer bool //have we sent peer notice we know it (true if we receive from this peer)
	blockNum        uint32
	expires         common.TimePointSec
	requestedTime   common.TimePoint
}

type updateTxnExpiry struct {
	newExpiry common.TimePointSec
}

func (u *updateTxnExpiry) updateTxnExpiry(e common.TimePointSec) {
	u.newExpiry = e
}

func (u *updateTxnExpiry) operator(ts *transactionState) { //TODO operator()??
	ts.expires = u.newExpiry
}

/*typedef multi_index_container<
transaction_state,
indexed_by<
ordered_unique< tag<by_id>, member<transaction_state, transaction_id_type, &transaction_state::id > >,
ordered_non_unique< tag< by_expiry >, member< transaction_state,fc::time_point_sec,&transaction_state::expires >>,
ordered_non_unique<
tag<by_block_num>,
member< transaction_state,
uint32_t,
&transaction_state::block_num > >
>

> transaction_state_index;
*/

type peerBlockState struct {
	id          common.BlockIdType
	blockNum    uint32
	isKnown     bool
	isNoticed   bool
	requestTime common.TimePoint
}

/*
struct update_request_time {
void operator() (struct transaction_state &ts) {
ts.requested_time = time_point::now();
}
void operator () (struct eosio::peer_block_state &bs) {
bs.requested_time = time_point::now();
}
} set_request_time;
*/
type updateRequestTime struct { //TODO func()

}

func (u *updateRequestTime) operator1(ts *transactionState) { //TODO operator1
	ts.requestedTime = common.Now()
}
func (u *updateRequestTime) operator2(bs *peerBlockState) { //TODO operator2
	bs.requestTime = common.Now()
}

var setRequestTime updateRequestTime

/*
typedef multi_index_container<
eosio::peer_block_state,
indexed_by<
ordered_unique< tag<by_id>, member<eosio::peer_block_state, block_id_type, &eosio::peer_block_state::id > >,
ordered_unique< tag<by_block_num>, member<eosio::peer_block_state, uint32_t, &eosio::peer_block_state::block_num > >
>
> peer_block_state_index;
*/

type updateKnownByPeer struct {
}

func (u *updateKnownByPeer) operator1(bs *peerBlockState) { //TODO operator1
	bs.isKnown = true
}
func (u *updateKnownByPeer) operator2(ts *transactionState) { //TODO operator2
	ts.isKnownByPeer = true
}

/*
struct update_known_by_peer {
void operator() (eosio::peer_block_state& bs) {
bs.is_known = true;
}
void operator() (transaction_state& ts) {
ts.is_known_by_peer = true;
}
} set_is_known;
*/

var setIsKnown updateKnownByPeer

type updateBlockNum struct {
	newBnum uint32
}

func (u *updateBlockNum) updateBlockNum(bnum uint32) {
	u.newBnum = bnum
}

func (u *updateBlockNum) operator1(nts *nodeTransactionState) { //TODO operator1
	if nts.request != 0 {
		nts.trueBlock = u.newBnum
	} else {
		nts.blockNum = u.newBnum
	}
}

func (u *updateBlockNum) operator2(ts *transactionState) { //TODO operator2
	ts.blockNum = u.newBnum
}

func (u *updateBlockNum) operator3(pbs *peerBlockState) { //TODO operator3
	pbs.blockNum = u.newBnum
}

/*
struct update_block_num {
uint32_t new_bnum;
update_block_num(uint32_t bnum) : new_bnum(bnum) {}
void operator() (node_transaction_state& nts) {
if (nts.requests ) {
nts.true_block = new_bnum;
}
else {
nts.block_num = new_bnum;
}
}
void operator() (transaction_state& ts) {
ts.block_num = new_bnum;
}
void operator() (peer_block_state& pbs) {
pbs.block_num = new_bnum;
}
};
*/

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

type nodeTransactionState struct {
	id            common.TransactionIdType
	expires       common.TimePointSec //time after which this may be purged.Expires increased while the txn is "in flight" to another peer
	packedTxn     types.PackedTransaction
	serializedTxn common.HexBytes // the received raw bundle
	blockNum      uint32          // block transaction was included in
	trueBlock     uint32          // used to reset block_uum when request is 0
	request       uint16          // the number of "in flight" requests for this txn
}

type updateInFlight struct {
	incr int32
}

func newUpdateInFlight(delta int32) *updateInFlight {
	return &updateInFlight{
		incr: delta,
	}
}
func (u *updateInFlight) operator(nts *nodeTransactionState) {
	exp := nts.expires.SecSinceEpoch()
	nts.expires = common.NewTimePointSecTp(common.TimePoint(exp + uint32(u.incr*60)))
	if nts.request == 0 {
		nts.trueBlock = nts.blockNum
		nts.blockNum = 0
	}
	nts.request += uint16(u.incr) //TODO int32 -> uint16
	if nts.request == 0 {
		nts.blockNum = nts.trueBlock
	}
}

//incrInFlight := newUpdateInFlight(1)
//decrInFlight := newUpdateInFlight(-1)

/*
struct update_in_flight {
int32_t incr;
update_in_flight (int32_t delta) : incr (delta) {}
void operator() (node_transaction_state& nts) {
int32_t exp = nts.expires.sec_since_epoch();
nts.expires = fc::time_point_sec (exp + incr * 60);
if( nts.requests == 0 ) {
nts.true_block = nts.block_num;
nts.block_num = 0;
}
nts.requests += incr;
if( nts.requests == 0 ) {
nts.block_num = nts.true_block;
}
}
} incr_in_flight(1), decr_in_flight(-1);
*/

/*
struct by_expiry;
struct by_block_num;

typedef multi_index_container<
node_transaction_state,
indexed_by<
ordered_unique<
tag< by_id >,
member < node_transaction_state,
transaction_id_type,
&node_transaction_state::id > >,
ordered_non_unique<
tag< by_expiry >,
member< node_transaction_state,
fc::time_point_sec,
&node_transaction_state::expires >
>,
ordered_non_unique<
tag<by_block_num>,
member< node_transaction_state,
uint32_t,
&node_transaction_state::block_num > >
>
>
node_transaction_index;
*/

//func (N *NetPluginIMpl)findConnection(host string) *connection_ptr{
//
//}

type stages byte

const (
	libCatchup = stages(iota)
	headCatchup
	inSync
)

type syncManager struct {
	syncKnownLibNum      uint32
	syncLastRequestedNum uint32
	syncNextExpectedNum  uint32
	syncReqSpan          uint32
	//source ConnectionPtr
	state   stages
	_blocks common.BlockIdType //<deque<block_id_type>>
	//chainPlugin *chainPlugin

}

func NewsyncManager(span uint32) *syncManager {
	//chainPlugin :=
	return &syncManager{
		syncKnownLibNum:      0,
		syncLastRequestedNum: 0,
		syncNextExpectedNum:  1,
		syncReqSpan:          span,
		//source:
		state: inSync,
	}
}

func stageStr(s stages) string {
	switch s {
	case libCatchup:
		return "lib catchup"
	case headCatchup:
		return "head catchup"
	case inSync:
		return "in sync"
	default:
		return "unkown"
	}
}

func (s *syncManager) setStage(newstate stages) {
	if s.state == newstate {
		return
	}
	fmt.Printf("old state %s becoming %s \n", stageStr(s.state), stageStr(newstate))
	s.state = newstate
}

func (s *syncManager) syncRequired() bool {
	fmt.Printf("last req = %d,last recv = %d known = %d our head %d\n", s.syncLastRequestedNum, s.syncNextExpectedNum, s.syncKnownLibNum, 100) //chain_plug->chain( ).head_block_num( )
	return s.syncLastRequestedNum < s.syncKnownLibNum || 0 < s.syncLastRequestedNum                                                            //100  ---->  chain_plug->chain( ).head_block_num( )
}

//func (s *syncManager) isActive(c net.Conn) bool {
//	if s.state == headCatchup && c {
//		fhset := c.forkHead != common.BlockIdType()
//		// return c.forkHead != common.BlockIdType() && c.forkHeadNum < chainPlugin
//	}
//	return s.state != inSync
//
//}

// bool sync_manager::is_active(connection_ptr c) {
//    if (state == head_catchup && c) {
//       bool fhset = c->fork_head != block_id_type();
//       fc_dlog(logger, "fork_head_num = ${fn} fork_head set = ${s}",
//               ("fn", c->fork_head_num)("s", fhset));
//          return c->fork_head != block_id_type() && c->fork_head_num < chain_plug->chain().head_block_num();
//    }
//    return state != in_sync;
// }

//func (s *syncManager) resetLibNum(c net.Conn) {
//	if s.state == inSync {
//		s.source.reset()
//	}
//	if c.Current() {
//
//	}
//}

// void sync_manager::reset_lib_num(connection_ptr c) {
//    if(state == in_sync) {
//       source.reset();
//    }
//    if( c->current() ) {
//       if( c->last_handshake_recv.last_irreversible_block_num > sync_known_lib_num) {
//          sync_known_lib_num =c->last_handshake_recv.last_irreversible_block_num;
//       }
//    } else if( c == source ) {
//       sync_last_requested_num = 0;
//       request_next_chunk();
//    }
// }
//func (s *syncManager) requestNextChunk(conn ConnectionPtr) {
//
//}

//func (s *syncManager) sendHandshakes() {
//
//	// for _,ci := range
//
//	// for( auto &ci : my_impl->connections) {
//	//    if( ci->current()) {
//	//       ci->send_handshake();
//	//    }
//	// }
//}

//func (s *syncManager) recvHandshake(c net.Conn, msg *HandshakeMessage) {
//
//}

// void sync_manager::recv_handshake (connection_ptr c, const handshake_message &msg) {
//    controller& cc = chain_plug->chain();
//    uint32_t lib_num = cc.last_irreversible_block_num( );
//    uint32_t peer_lib = msg.last_irreversible_block_num;
//    reset_lib_num(c);
//    c->syncing = false;

//    //--------------------------------
//    // sync need checks; (lib == last irreversible block)
//    //
//    // 0. my head block id == peer head id means we are all caugnt up block wise
//    // 1. my head block num < peer lib - start sync locally
//    // 2. my lib > peer head num - send an last_irr_catch_up notice if not the first generation
//    //
//    // 3  my head block num <= peer head block num - update sync state and send a catchup request
//    // 4  my head block num > peer block num ssend a notice catchup if this is not the first generation
//    //
//    //-----------------------------

//    uint32_t head = cc.head_block_num( );
//    block_id_type head_id = cc.head_block_id();
//    if (head_id == msg.head_id) {
//       fc_dlog(logger, "sync check state 0");
//       // notify peer of our pending transactions
//       notice_message note;
//       note.known_blocks.mode = none;
//       note.known_trx.mode = catch_up;
//       note.known_trx.pending = my_impl->local_txns.size();

//       // transaction_id_type id ;
//       // block_id_type bid ;
//       // note.known_trx.ids.push_back(id);
//       // note.known_trx.ids.push_back(id);
//       // note.known_blocks.pending = 1000;// walker none;
//       // note.known_blocks.ids.push_back(bid);// walker none;
//       // note.known_blocks.ids.push_back(bid);// walker none;
//       c->enqueue( note );
//       return;
//    }
//    if (head < peer_lib) {
//       fc_dlog(logger, "sync check state 1");
//       wlog("sync check state 1");
//       // wait for receipt of a notice message before initiating sync
//       if (c->protocol_version < proto_explicit_sync) {
//          start_sync( c, peer_lib);
//       }
//       return;
//    }
//    if (lib_num > msg.head_num ) {
//       fc_dlog(logger, "sync check state 2");
//       if (msg.generation > 1 || c->protocol_version > proto_base) {
//          notice_message note;
//          note.known_trx.pending = lib_num;
//          note.known_trx.mode = last_irr_catch_up;
//          note.known_blocks.mode = last_irr_catch_up;
//          note.known_blocks.pending = head;
//          c->enqueue( note );
//       }
//       c->syncing = true;
//       return;
//    }

//    if (head <= msg.head_num ) {
//       fc_dlog(logger, "sync check state 3");
//       verify_catchup (c, msg.head_num, msg.head_id);
//       return;
//    }
//    else {
//       fc_dlog(logger, "sync check state 4");
//       if (msg.generation > 1 ||  c->protocol_version > proto_base) {
//          notice_message note;
//          note.known_trx.mode = none;
//          note.known_blocks.mode = catch_up;
//          note.known_blocks.pending = head;
//          note.known_blocks.ids.push_back(head_id);
//          c->enqueue( note );
//       }
//       c->syncing = true;
//       return;
//    }
//    elog ("sync check failed to resolve status");
// }

//func (s *syncManager) startSync(c net.Conn, target uint32) {
//	if target > s.syncKnownLibNum {
//		s.syncKnownLibNum = target
//	}
//	if !s.syncRequired() {
//		bnum := 100 //chain_plug->chain().last_irreversible_block_num()
//		hnum := 100 //chain_plug->chain().head_block_num()
//		fmt.Printf("we are already caught up, my irr = %d,head =%d,target = %d\n", bnum, hnum, target)
//		return
//	}
//	if s.state == inSync {
//		s.setStage(libCatchup)
//		s.syncNextExpectedNum = 99 + 1 //chain_plug->chain().last_irreversible_block_num() + 1
//	}
//	fmt.Printf("Catching up with chain, our last req is %d, theirs is %d peer %s\n", s.syncLastRequestedNum, target, "walker") //walker  c->peer_name()
//	// s.requestNextChunk(c)
//}

// void sync_manager::start_sync( connection_ptr c, uint32_t target) {

//    if (!sync_required()) {
//       uint32_t bnum = chain_plug->chain().last_irreversible_block_num();
//       uint32_t hnum = chain_plug->chain().head_block_num();
//       fc_dlog( logger, "We are already caught up, my irr = ${b}, head = ${h}, target = ${t}",
//                ("b",bnum)("h",hnum)("t",target));
//       return;
//    }

//    if (state == in_sync) {
//       set_state(lib_catchup);
//       sync_next_expected_num = chain_plug->chain().last_irreversible_block_num() + 1;
//    }

//    fc_ilog(logger, "Catching up with chain, our last req is ${cc}, theirs is ${t} peer ${p}",
//            ( "cc",sync_last_requested_num)("t",target)("p",c->peer_name()));

//    wlog("Catching up with chain, our last req is ${cc}, theirs is ${t} peer ${p}",
//            ( "cc",sync_last_requested_num)("t",target)("p",c->peer_name()));
//    request_next_chunk(c);
// }

//func (s *syncManager) reassignFetch(c net.Conn, reason GoAwayReason) {
//	fmt.Printf("reassign_fetch, our last req is %d, next expected is %d peer %s\n", s.syncLastRequestedNum, s.syncNextExpectedNum, "walker") //walker c->peer_name()
//	if c == source {
//		c.cancelSync(reason)
//		s.syncLastRequestedNum = 0
//		s.requestNextChunk()
//	}
//
//}
//
//func (s *syncManager) verifyCatchup(c net.Conn, num uint32, id common.BlockIdType) {
//
//}

// void sync_manager::verify_catchup(connection_ptr c, uint32_t num, block_id_type id) {
//    request_message req;
//    req.req_blocks.mode = catch_up;
//    for (auto cc : my_impl->connections) {
//       if (cc->fork_head == id ||
//           cc->fork_head_num > num)
//          req.req_blocks.mode = none;
//       break;
//    }
//    if( req.req_blocks.mode == catch_up ) {//所有conn中最长的链
//       c->fork_head = id;
//       c->fork_head_num = num;
//       ilog ("got a catch_up notice while in ${s}, fork head num = ${fhn} target LIB = ${lib} next_expected = ${ne}", ("s",stage_str(state))("fhn",num)("lib",sync_known_lib_num)("ne", sync_next_expected_num));
//       if (state == lib_catchup)
//          return;
//       set_state(head_catchup);
//    }
//    else {
//       c->fork_head = block_id_type();
//       c->fork_head_num = 0;
//    }
//    req.req_trx.mode = none;
//    c->enqueue( req );
// }

//func (s *syncManager) recvNotice(c net.Conn, msg *NoticeMessage) {
//
//}

// void sync_manager::recv_notice (connection_ptr c, const notice_message &msg) {
//    fc_ilog (logger, "sync_manager got ${m} block notice",("m",modes_str(msg.known_blocks.mode)));
//    if (msg.known_blocks.mode == catch_up) {
//       if (msg.known_blocks.ids.size() == 0) {
//          elog ("got a catch up with ids size = 0");
//       }
//       else {
//          verify_catchup(c,  msg.known_blocks.pending, msg.known_blocks.ids.back());
//       }
//    }
//    else {
//       c->last_handshake_recv.last_irreversible_block_num = msg.known_trx.pending;
//       reset_lib_num (c);
//       start_sync(c, msg.known_blocks.pending);
//    }
// }

//func (s *syncManager) rejectedBlock(c net.conn, blkNum uint32) {
//	if s.state != inSync {
//		fmt.Printf("block %d not accepted from %s", blkNum, "walker") //walker c->peer_name()
//		s.syncLastRequestedNum = 0
//		s.source.reset()
//		my_impl.close(c)
//		s.setStage(inSync)
//		s.snedHandshakes()
//	}
//}

// void sync_manager::rejected_block (connection_ptr c, uint32_t blk_num) {
//   if (state != in_sync ) {
//      fc_ilog (logger, "block ${bn} not accepted from ${p}",("bn",blk_num)("p",c->peer_name()));
//      sync_last_requested_num = 0;
//      source.reset();
//      my_impl->close(c);
//      set_state(in_sync);
//      send_handshakes();
//   }
// }
//func (s *syncManager) recvBlock(c net.Conn, blkID *common.BlockIdType, blkNum uint32) {
//
//}

// void sync_manager::recv_block (connection_ptr c, const block_id_type &blk_id, uint32_t blk_num) {
//    fc_dlog(logger," got block ${bn} from ${p}",("bn",blk_num)("p",c->peer_name()));
//    if (state == lib_catchup) {
//       if (blk_num != sync_next_expected_num) {
//          fc_ilog (logger, "expected block ${ne} but got ${bn}",("ne",sync_next_expected_num)("bn",blk_num));
//          my_impl->close(c);
//          return;
//       }
//       sync_next_expected_num = blk_num + 1;
//    }
//    if (state == head_catchup) {
//       fc_dlog (logger, "sync_manager in head_catchup state");
//       set_state(in_sync);
//       source.reset();

//       block_id_type null_id;
//       for (auto cp : my_impl->connections) {
//          if (cp->fork_head == null_id) {
//             continue;
//          }
//          if (cp->fork_head == blk_id || cp->fork_head_num < blk_num) {
//             c->fork_head = null_id;
//             c->fork_head_num = 0;
//          }
//          else {
//             set_state(head_catchup);
//          }
//       }
//    }
//    else if (state == lib_catchup) {
//       if( blk_num == sync_known_lib_num ) {
//          fc_dlog( logger, "All caught up with last known last irreversible block resending handshake");

//          wlog("All caught up with last known last irreversible block resending handshake");
//          set_state(in_sync);
//          send_handshakes();
//       }
//       else if (blk_num == sync_last_requested_num) {
//          // source->request_sync_blocks(start, end);
//          // sync_last_requested_num = end;
//          request_next_chunk();
//       }
//       else {
//          fc_dlog(logger,"calling sync_wait on connection ${p}",("p",c->peer_name()));
//          c->sync_wait();
//       }
//    }
// }
