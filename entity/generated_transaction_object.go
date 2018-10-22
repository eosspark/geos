package entity

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	arithmetic "github.com/eosspark/eos-go/common/arithmetic_types"
	"github.com/eosspark/eos-go/crypto/rlp"
)

type GeneratedTransactionObject struct {
	Id         common.IdType            `multiIndex:"id,increment,byExpiration,byDelay"`
	TrxId      common.TransactionIdType `multiIndex:"byTrxId,orderedUnique"`
	Sender     common.AccountName       `multiIndex:"bySenderId,orderedUnique"`
	SenderId   arithmetic.Uint128       `multiIndex:"bySenderId,orderedUnique"`
	Payer      common.AccountName
	DelayUntil common.TimePoint         `multiIndex:"byDelay,orderedUnique"`
	Expiration common.TimePoint         `multiIndex:"byExpiration,orderedUnique"`
	Published  common.TimePoint
	PackedTrx  common.HexBytes //c++ shared_string
}

func (g *GeneratedTransactionObject) Set(trx *types.Transaction) uint32 {
	g.PackedTrx, _ = rlp.EncodeToBytes(trx)
	return uint32(len(g.PackedTrx))
}

type GeneratedTransaction struct {
	TrxId      common.TransactionIdType
	Sender     common.AccountName
	SenderId   arithmetic.Uint128
	Payer      common.AccountName
	DelayUntil common.TimePoint
	Expiration common.TimePoint
	Published  common.TimePoint
	PackedTrx  []byte
}

func GeneratedTransactions(gto *GeneratedTransactionObject) *GeneratedTransaction {
	gt := GeneratedTransaction{}
	gt.TrxId = gto.TrxId
	gt.Sender = gto.Sender
	gt.SenderId = gto.SenderId
	gt.Payer = gto.Payer
	gt.DelayUntil = gto.DelayUntil
	gt.Expiration = gto.Expiration
	gt.Published = gto.Published
	gt.PackedTrx[0] = gto.PackedTrx[0]
	gt.PackedTrx[1] = gto.PackedTrx[len(gto.PackedTrx)]
	return &gt
}
