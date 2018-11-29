package eosapi

import (
	"context"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
)

type EOSApi struct {
}

func NewEosApi() *EOSApi {
	return &EOSApi{}
}

// read_only::get_info_results read_only::get_info(const read_only::get_info_params&) const {
//    const auto& rm = db.get_resource_limits_manager();
//    return {
//       eosio::utilities::common::itoh(static_cast<uint32_t>(app().version())),
//       db.get_chain_id(),
//       db.fork_db_head_block_num(),
//       db.last_irreversible_block_num(),
//       db.last_irreversible_block_id(),
//       db.fork_db_head_block_id(),
//       db.fork_db_head_block_time(),
//       db.fork_db_head_block_producer(),
//       rm.get_virtual_block_cpu_limit(),
//       rm.get_virtual_block_net_limit(),
//       rm.get_block_cpu_limit(),
//       rm.get_block_net_limit(),
//       //std::bitset<64>(db.get_dynamic_global_properties().recent_slots_filled).to_string(),
//       //__builtin_popcountll(db.get_dynamic_global_properties().recent_slots_filled) / 64.0,
//       app().version_string(),
//    };
// }
func (api *EOSApi) GetInfo(ctx context.Context) *InfoResp {
	return &InfoResp{
		ServerVersion:            "0f6695cb",
		ChainID:                  common.BlockIdNil(),
		HeadBlockNum:             17673,
		LastIrreversibleBlockNum: 17672,
		LastIrreversibleBlockID:  common.BlockIdNil(),
		HeadBlockID:              common.BlockIdNil(),
		HeadBlockTime:            common.Now(),
		HeadBlockProducer:        common.AccountName(common.N("eosio")),
		VirtualBlockCPULimit:     200000000,
		VirtualBlockNetLimit:     1048576000,
		BlockCPULimit:            199900,
		BlockNetLimit:            1048576,
		ServerVersionString:      "TODO walker",
	}
}

type Keys struct {
	Pri ecc.PrivateKey `json:"Private Key"`
	Pub ecc.PublicKey  `json:"Public Key"`
}

func (api *EOSApi) CreateKey() *Keys {
	prikey, _ := ecc.NewRandomPrivateKey()
	return &Keys{Pri: *prikey, Pub: prikey.PublicKey()}

}
