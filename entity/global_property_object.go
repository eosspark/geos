package entity

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type GlobalPropertyObject struct {
	ID                       common.IdType `multiIndex:"id,increment"`
	ProposedScheduleBlockNum uint32
	ProposedSchedule         types.SharedProducerScheduleType
	Configuration            types.ChainConfig
}

type DynamicGlobalPropertyObject struct {
	ID                   common.IdType `multiIndex:"id,increment"` //c++ chainbase.hpp id_type
	GlobalActionSequence uint64
}

func (g GlobalPropertyObject) IsEmpty() bool {
	return g.ID == 0 && g.ProposedScheduleBlockNum == 0 &&
		g.ProposedSchedule.IsEmpty() && g.Configuration.IsEmpty()
}

/*
type Config struct {
	ActorWhitelist          common.FlatSet //common.AccountName
	ActorBlacklist          common.FlatSet //common.AccountName
	ContractWhitelist       common.FlatSet //common.AccountName
	ContractBlacklist       common.FlatSet //common.AccountName]struct{}
	ActionBlacklist         common.FlatSet //common.Pair //see actionBlacklist
	KeyBlacklist            common.FlatSet
	blocksDir               string
	stateDir                string
	stateSize               uint64
	stateGuardSize          uint64
	reversibleCacheSize     uint64
	reversibleGuardSize     uint64
	readOnly                bool
	forceAllChecks          bool
	disableReplayOpts       bool
	disableReplay           bool
	contractsConsole        bool
	allowRamBillingInNotify bool
	genesis                 types.GenesisState
	vmType                  wasmgo.WasmGo
	readMode                DBReadMode
	blockValidationMode     ValidationMode
	resourceGreylist        map[common.AccountName]struct{}
	trustedProducers        map[common.AccountName]struct{}
}*/
