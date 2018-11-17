package entity

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
)

type AccountObject struct {
	ID             common.IdType      `multiIndex:"id,increment"`
	Name           common.AccountName `multiIndex:"byName,orderedUnique"`
	VmType         uint8              //c++ default value 0
	VmVersion      uint8              //c++ default value 0
	Privileged     bool               //c++ default value false
	LastCodeUpdate common.TimePoint
	CodeVersion    crypto.Sha256
	CreationDate   types.BlockTimeStamp
	Code           common.HexBytes
	Abi            common.HexBytes
}

type AccountSequenceObject struct {
	ID           common.IdType      `multiIndex:"id,increment"`
	Name         common.AccountName `multiIndex:"byName,orderedUnique"`
	RecvSequence uint64             //default value 0
	AuthSequence uint64
	CodeSequence uint64
	AbiSequence  uint64
}

func (a *AccountObject) SetAbi(ad types.AbiDef) {
	d, _ := rlp.EncodeToBytes(ad)
	a.Abi = d
}

func (a *AccountObject) GetAbi() types.AbiDef {
	abiDef := types.AbiDef{}
	if len(a.Abi) != 0 {
		try.EosAssert(len(a.Abi) != 0, &exception.AbiNotFoundException{}, "No ABI set on account :", a.Name)
	}
	err := rlp.DecodeBytes(a.Abi, abiDef)
	if err != nil {
		fmt.Println("account_object GetAbi DecodeBytes is error:", err.Error())
	}
	return abiDef
}
