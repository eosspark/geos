package chain_plugin

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"strconv"
)

type ReadOnly struct {
	db                   *chain.Controller
	abiSerializerMaxTime common.Microseconds
	shortenAbiErrors     bool
}

func NewReadOnly(db *chain.Controller, abiSerializerMaxTime common.Microseconds) *ReadOnly {
	ro := new(ReadOnly)
	ro.db = db
	ro.abiSerializerMaxTime = abiSerializerMaxTime
	return ro
}

//read_only::get_info_results read_only::get_info(const read_only::get_info_params&) const {
//   const auto& rm = db.get_resource_limits_manager();
//   return {
//      eosio::utilities::common::itoh(static_cast<uint32_t>(app().version())),
//      db.get_chain_id(),
//      db.fork_db_head_block_num(),
//      db.last_irreversible_block_num(),
//      db.last_irreversible_block_id(),
//      db.fork_db_head_block_id(),
//      db.fork_db_head_block_time(),
//      db.fork_db_head_block_producer(),
//      rm.get_virtual_block_cpu_limit(),
//      rm.get_virtual_block_net_limit(),
//      rm.get_block_cpu_limit(),
//      rm.get_block_net_limit(),
//      //std::bitset<64>(db.get_dynamic_global_properties().recent_slots_filled).to_string(),
//      //__builtin_popcountll(db.get_dynamic_global_properties().recent_slots_filled) / 64.0,
//      app().version_string(),
//   };
//}

func (ro *ReadOnly) GetInfo() *InfoResp {
	rm := ro.db.GetMutableResourceLimitsManager()
	return &InfoResp{
		ServerVersion:            "0f6695cb", //eosio::utilities::common::itoh(static_cast<uint32_t>(app().version())),
		ChainID:                  ro.db.GetChainId(),
		HeadBlockNum:             ro.db.ForkDbHeadBlockNum(),
		LastIrreversibleBlockNum: ro.db.LastIrreversibleBlockNum(),
		LastIrreversibleBlockID:  ro.db.LastIrreversibleBlockId(),
		HeadBlockID:              ro.db.ForkDbHeadBlockId(),
		HeadBlockTime:            ro.db.ForkDbHeadBlockTime(),
		HeadBlockProducer:        ro.db.ForkDbHeadBlockProducer(),
		VirtualBlockCPULimit:     rm.GetVirtualBlockCpuLimit(),
		VirtualBlockNetLimit:     rm.GetVirtualBlockNetLimit(),
		BlockCPULimit:            rm.GetBlockCpuLimit(),
		BlockNetLimit:            rm.GetBlockNetLimit(),
		ServerVersionString:      "TODO walker", //app().version_string(),
	}
}

// fc::variant read_only::get_block(const read_only::get_block_params& params) const {
//    signed_block_ptr block;
//    EOS_ASSERT(!params.block_num_or_id.empty() && params.block_num_or_id.size() <= 64, chain::block_id_type_exception, "Invalid Block number or ID, must be greater than 0 and less than 64 characters" );
//    try {
//       block = db.fetch_block_by_id(fc::variant(params.block_num_or_id).as<block_id_type>());
//       if (!block) {
//          block = db.fetch_block_by_number(fc::to_uint64(params.block_num_or_id));
//       }

//    } EOS_RETHROW_EXCEPTIONS(chain::block_id_type_exception, "Invalid block ID: ${block_num_or_id}", ("block_num_or_id", params.block_num_or_id))

//    EOS_ASSERT( block, unknown_block_exception, "Could not find block: ${block}", ("block", params.block_num_or_id));

//    fc::variant pretty_output;
//    abi_serializer::to_variant(*block, pretty_output, make_resolver(this, abi_serializer_max_time), abi_serializer_max_time);

//    uint32_t ref_block_prefix = block->id()._hash[1];

//    return fc::mutable_variant_object(pretty_output.get_object())
//            ("id", block->id())
//            ("block_num",block->block_num())
//            ("ref_block_prefix", ref_block_prefix);
// }
func (ro *ReadOnly) GetBlock(params string) *BlockResp {
	block := &types.SignedBlock{}
	EosAssert(len(params) != 0 && len(params) <= 64, &exception.BlockIdTypeException{},
		"Invalid Block number or ID,must be greater than 0 and less than 64 characters")

	Try(func() {
		blockID := common.BlockIdType(*crypto.NewSha256String(params)) //TODO panic??
		block = ro.db.FetchBlockById(blockID)
		if common.Empty(block) {
			blockNum, _ := strconv.Atoi(params)
			block = ro.db.FetchBlockByNumber(uint32(blockNum)) // TODO Uint64
		}
	}).EosRethrowExceptions(&exception.BlockIdTypeException{}, "Invalid block ID: %s", params)

	EosAssert(!common.Empty(block), &exception.UnknownBlockException{}, "Could not find block: %s", params)

	refBlockPrefix := uint32(block.BlockID().Hash[1])
	return &BlockResp{
		SignedBlock:    *block,
		ID:             block.BlockID(),
		BlockNum:       block.BlockNumber(),
		RefBlockPrefix: refBlockPrefix,
	}
}

//read_only::get_abi_results read_only::get_abi( const get_abi_params& params )const {
//   get_abi_results result;
//   result.account_name = params.account_name;
//   const auto& d = db.db();
//   const auto& accnt  = d.get<account_object,by_name>( params.account_name );
//
//   abi_def abi;
//   if( abi_serializer::to_abi(accnt.abi, abi) ) {
//      result.abi = std::move(abi);
//   }
//
//   return result;
//}

func (ro *ReadOnly) GetAbi(name common.AccountName) GetABIResp {
	result := GetABIResp{}
	result.AccountName = name
	d := ro.db.DataBase()
	accountObject := entity.AccountObject{}
	mutilx, _ := d.GetIndex("byName", &accountObject)
	accnt := entity.AccountObject{}
	accnt.Name = name
	err := mutilx.Find(&accnt, &accnt)
	if err != nil {
		log.Error("not find account_object %s", name)
	}

	//var abi  types.AbiDef //TODO need serializer::to_abi
	//   abi_def abi;
	//   if( abi_serializer::to_abi(accnt.abi, abi) ) {
	//      result.abi = std::move(abi);
	//   }
	return result
}
