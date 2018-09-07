package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
)

type GenesisState struct {
	EosioRootKey     string        `json:"eosio_root_key"`
	InitialTimestamp common.Tstamp `json:"initial_timestamp"`
	InitialKey       ecc.PublicKey `json:"initial_key"`
}

func (gs *GenesisState) initial() {

}
