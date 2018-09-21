package types

import (
	"github.com/eosspark/eos-go/common"
)

type TransactionMetadata struct {
	ID          common.TransactionIDType
	SignedID    common.TransactionIDType
	Trx         SignedTransaction
	PackedTrx   PackedTransaction
	SigningKeys pair
	Accepted    bool //default value false
	Implicit    bool //default value false
	Scheduled   bool //default value false
}

type pair struct {
	id        common.ChainIDType
	publicKey []common.PublicKeyType
}

func NewTransactionMetadata(ptrx PackedTransaction) *TransactionMetadata {
	//TODO
	return new(TransactionMetadata)
}

func (tm *TransactionMetadata) RecoverKeys(chainId common.ChainIDType) []common.PublicKeyType {
	if /*unsafe.Sizeof(tm.SigningKeys) || */ tm.SigningKeys.id != chainId {
		tm.SigningKeys.id = chainId
		//tm.SigningKeys.publicKey = tm.Trx.GetSignatureKeys(chainId)
	}
	return nil
}
