package chain_plugin

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/abi_serializer"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"strconv"
	"github.com/eosspark/eos-go/plugins/appbase/app"
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

func (ro *ReadOnly) GetInfo() *InfoResp {
	rm := ro.db.GetMutableResourceLimitsManager()
	return &InfoResp{
		//ServerVersion:            "0f6695cb", //eosio::utilities::common::itoh(static_cast<uint32_t>(app().version())),
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

func (ro *ReadOnly) GetBlock(params string) *BlockResp {
	block := &types.SignedBlock{}
	EosAssert(len(params) != 0 && len(params) <= 64, &exception.BlockIdTypeException{},
		"Invalid Block number or ID,must be greater than 0 and less than 64 characters ")

	Try(func() {
		blockID := common.BlockIdType(*crypto.NewSha256String(params)) //TODO panic??
		block = ro.db.FetchBlockById(blockID)
		if common.Empty(block) {
			blockNum, _ := strconv.Atoi(params)
			block = ro.db.FetchBlockByNumber(uint32(blockNum)) // TODO Uint64
		}
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
