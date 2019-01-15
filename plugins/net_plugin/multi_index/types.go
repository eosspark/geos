package multi_index

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

//    typedef multi_index_container<
//       node_transaction_state,
//       indexed_by<
//          ordered_unique<
//             tag< by_id >,
//             member < node_transaction_state,
//                      transaction_id_type,
//                      &node_transaction_state::id > >,
//          ordered_non_unique<
//             tag< by_expiry >,
//             member< node_transaction_state,
//                     fc::time_point_sec,
//                     &node_transaction_state::expires >
//             >,
//          ordered_non_unique<
//             tag<by_block_num>,
//             member< node_transaction_state,
//                     uint32_t,
//                     &node_transaction_state::block_num > >
//          >
//       >
//    node_transaction_index;

type NodeTransactionState struct {
	ID            common.TransactionIdType
	Expires       common.TimePointSec
	PackedTxn     types.PackedTransaction
	SerializedTxn []byte
	BlockNum      uint32
	TrueBlock     uint32
	Requests      uint16
}

// typedef multi_index_container<
//    transaction_state,
//    indexed_by<
//       ordered_unique< tag<by_id>, member<transaction_state, transaction_id_type, &transaction_state::id > >,
//       ordered_non_unique< tag< by_expiry >, member< transaction_state,fc::time_point_sec,&transaction_state::expires >>,
//       ordered_non_unique<
//          tag<by_block_num>,
//          member< transaction_state,
//                  uint32_t,
//                  &transaction_state::block_num > >
//       >

//    > transaction_state_index;

type TransactionState struct {
	ID              common.TransactionIdType
	IsKnownByPeer   bool
	IsNoticedToPeer bool
	BlockNum        uint32
	Expires         common.TimePointSec
	RequestedTime   common.TimePoint
}

// typedef multi_index_container<
//    eosio::peer_block_state,
//    indexed_by<
//       ordered_unique< tag<by_id>, member<eosio::peer_block_state, block_id_type, &eosio::peer_block_state::id > >,
//       ordered_unique< tag<by_block_num>, member<eosio::peer_block_state, uint32_t, &eosio::peer_block_state::block_num > >
//       >
//    > peer_block_state_index;

type PeerBlockState struct {
	ID            common.BlockIdType
	BlockNum      uint32
	IsKnown       bool
	IsNoticed     bool
	RequestedTime common.TimePoint
}
