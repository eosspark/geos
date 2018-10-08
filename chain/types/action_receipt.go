package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
)

type ActionReceipt struct {
	Receiver       common.AccountName            `json:"receiver"`
	ActDigest      crypto.Sha256                 `json:"act_digest"`
	GlobalSequence uint64                        `json:"global_sequence"`
	RecvSequence   uint64                        `json:"recv_sequence"`
	AuthSequence   map[common.AccountName]uint64 `json:"auth_sequence"`
	CodeSequence   uint32                        `json:"code_sequence" eos:vuint32` //TODO
	AbiSequence    uint32                        `json:"abi_sequence" eos:vuint32`	//TODO
}

func (self *ActionReceipt) Digest() crypto.Sha256 {
	return crypto.Hash256(self)
}
