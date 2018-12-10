package chain_plugin

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/abi_serializer"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/appbase/app"
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

func (ro *ReadOnly) SetShortenAbiErrors(f bool) {
	ro.shortenAbiErrors = f
}

func (ro *ReadOnly) GetInfo() *InfoResp {
	rm := ro.db.GetMutableResourceLimitsManager()
	return &InfoResp{
		ServerVersion:            strconv.Itoa(int(app.App().GetVersion())), //eosio::utilities::common::itoh(static_cast<uint32_t>(app().version())),
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

type GetBlockParams struct {
	BlockNumOrID string `json:"block_num_or_id"`
}

func (ro *ReadOnly) GetBlock(params GetBlockParams) *BlockResp {
	block := &types.SignedBlock{}
	EosAssert(len(params.BlockNumOrID) != 0 && len(params.BlockNumOrID) <= 64, &exception.BlockIdTypeException{},
		"Invalid Block number or ID,must be greater than 0 and less than 64 characters ")

	Try(func() {
		//blockID := common.BlockIdType(*crypto.NewSha256String(params)) //TODO panic??
		//block = ro.db.FetchBlockById(blockID)
		//if common.Empty(block) {
		blockNum, _ := strconv.Atoi(params.BlockNumOrID)
		block = ro.db.FetchBlockByNumber(uint32(blockNum)) // TODO Uint64
		//}
	}).EosRethrowExceptions(&exception.BlockIdTypeException{}, "Invalid block ID: %s", params).End()

	EosAssert(!common.Empty(block), &exception.UnknownBlockException{}, "Could not find block: %s", params)

	refBlockPrefix := uint32(block.BlockID().Hash[1])
	return &BlockResp{
		SignedBlock:    *block,
		ID:             block.BlockID(),
		BlockNum:       block.BlockNumber(),
		RefBlockPrefix: refBlockPrefix,
	}
}

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

	var abi abi_serializer.AbiDef
	if abi_serializer.ToABI(accnt.Abi, &abi) {
		result.ABI = abi
	}

	return result
}

// read_only::get_required_keys_result read_only::get_required_keys( const get_required_keys_params& params )const {
//    transaction pretty_input;
//    auto resolver = make_resolver(this, abi_serializer_max_time);
//    try {
//       abi_serializer::from_variant(params.transaction, pretty_input, resolver, abi_serializer_max_time);
//    } EOS_RETHROW_EXCEPTIONS(chain::transaction_type_exception, "Invalid transaction")

//    auto required_keys_set = db.get_authorization_manager().get_required_keys( pretty_input, params.available_keys, fc::seconds( pretty_input.delay_sec ));
//    get_required_keys_result result;
//    result.required_keys = required_keys_set;
//    return result;
// }

//struct get_required_keys_params {
//fc::variant transaction;
//flat_set<public_key_type> available_keys;
//};
//struct get_required_keys_result {
//flat_set<public_key_type> required_keys;
//};
//
type GetRequiredKeysParams struct {
	Transaction   map[string]interface{} `json:"transaction"`
	AvailableKeys []ecc.PublicKey        `json:"available_keys"`
}
type GetRequiredKeysResult struct {
	RequiredKeys []ecc.PublicKey `json:"required_keys"`
}

func (ro *ReadOnly) GetRequiredKeys(params *GetRequiredKeysParams) GetRequiredKeysResult {
	trx := types.Transaction{}
	re, err := json.Marshal(params.Transaction)
	fmt.Println(re, err)
	err = json.Unmarshal(re, &trx)
	fmt.Printf("trx:    ************** %#v,%s", trx, err)

	return GetRequiredKeysResult{
		RequiredKeys: []ecc.PublicKey{ecc.MustNewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")}}
}

//rekey = {"available_keys":[],"transaction":{"expiration":19991,"ref_block_num":90,"ref_block_prefix":888,"max_net_usage_words":0,"max_cpu_usage_ms":0,"delay_sec":899,"context_free_actions":"hello","actions":null,"transaction_extensions":null,"signatures":[],"context_free_data":[]}}
