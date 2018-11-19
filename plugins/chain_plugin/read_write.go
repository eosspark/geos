package chain_plugin

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
)

type ReadWrite struct {
	db                   *chain.Controller
	abiSerializerMaxTime common.Microseconds
}

func NewReadWrite(db *chain.Controller, abiSerializerMaxTime common.Microseconds) *ReadWrite {
	rw := new(ReadWrite)
	rw.db = db
	rw.abiSerializerMaxTime = abiSerializerMaxTime
	return rw
}

func (ReadWrite) PushTransaction() {

}