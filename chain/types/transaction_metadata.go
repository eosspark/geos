package types

import (
	"github.com/eosspark/container/sets/treeset"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
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
	PublicKey treeset.Set
}

func NewTransactionMetadata(ptrx *PackedTransaction) *TransactionMetadata {
	hashed := crypto.Hash256(ptrx)
	signedTransaction := ptrx.GetSignedTransaction()
	return &TransactionMetadata{
		ID:        signedTransaction.ID(),
		SignedID:  common.TransactionIdType(*hashed),
		Trx:       signedTransaction,
		PackedTrx: ptrx,
	}
}

func NewTransactionMetadataBySignedTrx(t *SignedTransaction, c CompressionType) *TransactionMetadata {
	hashed := crypto.Hash256(t)
	packedTrx := NewPackedTransactionBySignedTrx(t, c)
	return &TransactionMetadata{
		ID:        t.ID(),
		SignedID:  common.TransactionIdType(*hashed),
		Trx:       t,
		PackedTrx: packedTrx,
		Implicit:  false,
	}
}

func (t *TransactionMetadata) RecoverKeys(chainID *common.ChainIdType) *treeset.Set {
	////if( !signing_keys || signing_keys->first != chain_id ) TODO !signing_keys ？？  ->&tm.SigningKeys ==nil
	//if t.SigningKeys.ID != *chainID { // Unlikely for more than one chain_id to be used in one nodeos instance
	//	t.SigningKeys = Pair{
	//		ID:        *chainID,
	//		PublicKey: t.Trx.GetSignatureKeys(chainID, false, true),
	//	}
	//}
	//return &t.SigningKeys.PublicKey

	if common.Empty(t.SigningKeys) || t.SigningKeys.ID != *chainID { // Unlikely for more than one chain_id to be used in one nodeos instance
		t.SigningKeys = Pair{
			ID:        *chainID,
			PublicKey: t.Trx.GetSignatureKeys(chainID, false, true),
		}
	}
	return &t.SigningKeys.PublicKey

}

func (t *TransactionMetadata) TotalActions() uint32 {
	return uint32(len(t.Trx.ContextFreeActions) + len(t.Trx.Actions))
}
