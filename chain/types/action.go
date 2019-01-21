package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/rlp"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
)

// Action
type Action struct {
	Account       common.AccountName `json:"account"`
	Name          common.ActionName  `json:"name"`
	Authorization []PermissionLevel  `json:"authorization,omitempty"`
	Data          common.HexBytes    `json:"data"`
}

func (a Action) DataAs(t interface{}) {
	err := rlp.DecodeBytes(a.Data, t)
	if err != nil {
		EosThrow(&ParseErrorException{}, "action data parse error: %s", err.Error())
	}
}

type ContractTypesInterface interface {
	GetAccount() common.AccountName
	GetName() common.ActionName
}

type ActionReceipt struct {
	Receiver       common.AccountName            `json:"receiver"`
	ActDigest      crypto.Sha256                 `json:"act_digest"`
	GlobalSequence uint64                        `json:"global_sequence"`
	RecvSequence   uint64                        `json:"recv_sequence"`
	AuthSequence   map[common.AccountName]uint64 `json:"auth_sequence"`
	CodeSequence   common.Vuint32                `json:"code_sequence" eos:"vuint32"`
	AbiSequence    common.Vuint32                `json:"abi_sequence" eos:"vuint32"`
}

func (a *ActionReceipt) Digest() crypto.Sha256 {
	return *crypto.Hash256(a)
}
