package entity

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/rlp"
)

type AccountObject struct {
	ID             types.IdType       `storm:"id,increment" json:"id"`
	Name           common.AccountName `storm:"unique" json :"name"`
	VmType         uint8              //c++ default value 0
	VmVersion      uint8              //c++ default value 0
	Privileged     bool               //c++ default value false
	LastCodeUpdate common.TimePoint
	CodeVersion    rlp.Sha256
	CreationDate   common.BlockTimeStamp
	Code           common.HexBytes
	Abi            common.HexBytes
}

type AccountSequenceObject struct {
	ID           types.IdType       `storm:"id,increment" json:"id"`
	Name         common.AccountName `storm:"unique" json:name`
	RecvSequence uint64             //default value 0
	authSequence uint64
	CodeSequence uint64
	AbiSequence  uint64
}

func (self *AccountObject) SetAbi(a types.AbiDef) {
	d, _ := rlp.EncodeToBytes(a)
	self.Abi = d
}

func (self *AccountObject) GetAbi() types.AbiDef {
	abiDef := types.AbiDef{}
	if len(self.Abi) != 0 {
		fmt.Println("abi_not_found_exception ,No ABI set on account", self.Name)
	}
	err := rlp.DecodeBytes(self.Abi, abiDef)
	if err != nil {
		fmt.Println("account_object GetAbi DecodeBytes is error:", err.Error())
	}
	return abiDef
}
