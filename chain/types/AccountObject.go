package types

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/rlp"
)

type AccountObject struct {
	ID             IdType             `storm:"id,increment" json:"ById"`
	Name           common.AccountName `storm:"unique" json :"name"`
	VmType         uint8              //c++ default value 0
	VmVersion      uint8              //c++ default value 0
	Privileged     bool               //c++ default value false
	LastCodeUpdate common.TimePoint
	CodeVersion    rlp.Sha256
	CreationDate   common.BlockTimeStamp
	Code           string
	Abi            string
}

type AccountSequenceObject struct {
	ID           IdType             `storm:"id,increment" json:"ById"`
	Name         common.AccountName `storm:"unique" json:ByName`
	RecvSequence uint64             //default value 0
	authSequence uint64
	CodeSequence uint64
	AbiSequence  uint64
}

func (self *AccountObject) SetAbi(a AbiDef) {
	d, _ := rlp.EncodeToBytes(a)
	fmt.Println(d)
}

func (self *AccountObject) GetAbi() *AbiDef {
	if len(self.Abi) != 0 {
		//self.Abi
	}
	return nil
}
