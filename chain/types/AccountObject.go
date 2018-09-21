package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/rlp"
)

type AccountObject struct {
	ID IdType
	Name common.AccountName
	VmType	uint8	//c++ default value 0
	VmVersion uint8	//c++ default value 0
	Privileged	bool //c++ default value false
	LastCodeUpdate	common.TimePoint
	CodeVersion rlp.Sha256
	CreationDate	common.BlockTimeStamp
	Code	string
	ABI		string
}

/*func (ao *AccountObject) SetAbi(Abi){

}*/