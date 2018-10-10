package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
)

type TransactionMetadata struct {
	ID          common.TransactionIdType
	SignedID    common.TransactionIdType
	Trx         SignedTransaction
	PackedTrx   PackedTransaction
	SigningKeys pair
	Accepted    bool //default value false
	Implicit    bool //default value false
	Scheduled   bool //default value false
}

type pair struct {
	id        common.ChainIdType
	publicKey []ecc.PublicKey
}

func TransactionMetadata1(ptrx PackedTransaction) *TransactionMetadata {
	tm := TransactionMetadata{}
	tm.Trx = *ptrx.GetSignedTransaction()
	tm.PackedTrx = ptrx
	tm.ID = ptrx.ID()

	return &tm
}

func (tm *TransactionMetadata) RecoverKeys(chainId common.ChainIdType) []ecc.PublicKey {
	if /*unsafe.Sizeof(tm.SigningKeys) || */ tm.SigningKeys.id != chainId {
		tm.SigningKeys.id = chainId
		//tm.SigningKeys.publicKey = tm.Trx.GetSignatureKeys(chainId)
	}
	return nil
}
