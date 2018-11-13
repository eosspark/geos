package chain_plugin

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
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




