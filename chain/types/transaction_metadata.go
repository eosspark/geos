package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
)

type TransactionMetadata struct {
	ID          common.TransactionIdType `json:"id"`
	SignedID    common.TransactionIdType `json:"signed_id"`
	Trx         *SignedTransaction       `json:"trx"`
	PackedTrx   *PackedTransaction       `json:"packed_trx"`
	SigningKeys Pair                     `json:"signing_keys"`
	Accepted    bool                     `json:"accepted"`
	Implicit    bool                     `json:"implicit"`
	Scheduled   bool                     `json:"scheduled"`
}

type Pair struct {
	ID        common.ChainIdType
	PublicKey []ecc.PublicKey
}

func NewTransactionMetadata(ptrx *PackedTransaction) *TransactionMetadata {
	hashed := crypto.Hash256(ptrx)
	signedTransaction := ptrx.GetSignedTransaction()
	return &TransactionMetadata{
		ID:        signedTransaction.ID(),
		SignedID:  common.TransactionIdType(hashed),
		Trx:       signedTransaction,
		PackedTrx: ptrx,
	}
}

func NewTransactionMetadataBySignedTrx(t *SignedTransaction, c common.CompressionType) *TransactionMetadata {
	hashed := crypto.Hash256(t)
	packedTrx := NewPackedTransactionBySignedTrx(t, c)
	return &TransactionMetadata{
		ID:        t.ID(),
		SignedID:  common.TransactionIdType(hashed),
		Trx:       t,
		PackedTrx: packedTrx,
	}
}

func (tm *TransactionMetadata) RecoverKeys(chainID common.ChainIdType) []ecc.PublicKey {
	//if( !signing_keys || signing_keys->first != chain_id ) TODO !signing_keys ？？  ->&tm.SigningKeys ==nil
	if tm.SigningKeys.ID != chainID { // Unlikely for more than one chain_id to be used in one nodeos instance
		tm.SigningKeys = Pair{
			ID:        chainID,
			PublicKey: tm.Trx.GetSignatureKeys(chainID, false, true),
		}
	}
	return tm.SigningKeys.PublicKey

}

func (tm *TransactionMetadata) TotalActions() uint32 {
	return uint32(len(tm.Trx.ContextFreeActions) + len(tm.Trx.Actions))
}
