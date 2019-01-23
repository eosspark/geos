package unittests

import (
	"github.com/eosspark/container/sets/treeset"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	. "github.com/eosspark/eos-go/exception"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockWithInvalidTx(t *testing.T) {
	main := newBaseTester(true, chain.SPECULATIVE)
	defer main.close()
	var err error

	// First we create a valid block with valid transaction
	main.CreateAccount(common.N("newacc"), common.DefaultConfig.SystemAccountName, false, true)
	b := main.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	// Make a copy of the valid block and corrupt the transaction
	copyB := b
	signedTx := copyB.Transactions[len(copyB.Transactions)-1].Trx.PackedTransaction.GetSignedTransaction()
	act := signedTx.Actions[len(signedTx.Actions)-1]
	actData := chain.NewAccount{}
	act.DataAs(&actData)
	// Make the transaction invalid by having the new account name the same as the creator name
	actData.Name = actData.Creator
	act.Data, err = rlp.EncodeToBytes(actData)
	assert.NoError(t, err)

	// Re-sign the transaction
	signedTx.Signatures = make([]ecc.Signature, 0)
	priKey, chainId := main.getPrivateKey(common.DefaultConfig.SystemAccountName, "active"), main.Control.GetChainId()
	signedTx.Sign(&priKey, &chainId)
	// Replace the valid transaction with the invalid transaction
	invalidPackedTx := types.NewPackedTransactionBySignedTrx(signedTx, types.CompressionNone)
	copyB.Transactions[len(copyB.Transactions)-1].Trx.PackedTransaction = invalidPackedTx

	// Re-sign the block'
	headerBmroot := crypto.Hash256(common.MakePair(copyB.Digest(), main.Control.HeadBlockState().BlockrootMerkle.GetRoot()))
	sigDigest := crypto.Hash256(common.MakePair(headerBmroot, main.Control.HeadBlockState().PendingScheduleHash))
	copyB.ProducerSignature, err = priKey.Sign(sigDigest.Bytes())
	assert.NoError(t, err)

	// Push block with invalid transaction to other chain
	validator := newBaseTester(true, chain.SPECULATIVE)
	validator.Control.AbortBlock()

	CheckThrowException(t, &AccountNameExistsException{}, func() { validator.Control.PushBlock(copyB, types.Complete) })
}

type blockPair struct {
	first  *types.SignedBlock
	second *types.SignedBlock
}

func CorruptTrxInBlock(main *ValidatingTester, actName common.AccountName) (blockPair, error) {
	var err error

	// First we create a valid block with valid transaction
	main.CreateAccount(actName, common.DefaultConfig.SystemAccountName, false, true)
	b := main.ProduceBlockNoValidation(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	//return blockPair{b, b}, nil
	// Make a copy of the valid block and corrupt the transaction
	copyB := b
	signedTx := copyB.Transactions[len(copyB.Transactions)-1].Trx.PackedTransaction.GetSignedTransaction()
	// Corrupt one signature
	priKey, chainId := main.getPrivateKey(actName, "active"), main.Control.GetChainId()
	signedTx.Signatures = make([]ecc.Signature, 0)
	signedTx.Sign(&priKey, &chainId)

	// Replace the valid transaction with the invalid transaction
	invalidPackedTx := types.NewPackedTransactionBySignedTrx(signedTx, types.CompressionNone)
	copyB.Transactions[len(copyB.Transactions)-1].Trx.PackedTransaction = invalidPackedTx

	// Re-calculate the transaction merkle
	trxs := copyB.Transactions
	trxDigests := make([]common.DigestType, 0, len(trxs))
	for _, a := range trxs {
		trxDigests = append(trxDigests, a.Digest())
	}
	copyB.TransactionMRoot = types.Merkle(trxDigests)

	// Re-sign the block
	headerBmroot := crypto.Hash256(common.MakePair(copyB.Digest(), main.Control.HeadBlockState().BlockrootMerkle.GetRoot()))
	sigDigest := crypto.Hash256(common.MakePair(headerBmroot, main.Control.HeadBlockState().PendingScheduleHash))
	priKey = main.getPrivateKey(b.Producer, "active")
	if copyB.ProducerSignature, err = priKey.Sign(sigDigest.Bytes()); err != nil {
		return blockPair{}, err
	}

	return blockPair{b, copyB}, nil
}

func TestTrustedProducer(t *testing.T) {
	trustedProducers := treeset.NewWith(common.TypeName, common.CompareName, common.N("defproducera"), common.N("defproducerc"))
	main := NewValidatingTesterTrustedProducers(trustedProducers)
	defer main.close()
	// only using validating_tester to keep the 2 chains in sync, not to validate that the validating_node matches the main node,
	// since it won't be
	main.SkipValidate = true

	// First we create a valid block with valid transaction
	producers := []common.AccountName{common.N("defproducera"),
		common.N("defproducerb"), common.N("defproducerc"), common.N("defproducerd")}

	for _, prod := range producers {
		main.CreateAccount(prod, common.DefaultConfig.SystemAccountName, false, true)
	}
	b := main.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	trace := main.SetProducers(&producers)
	assert.Equal(t, Exception(nil), trace.Except)

	for b.Producer != common.N("defproducera") {
		b = main.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	}

	blocks, err := CorruptTrxInBlock(main, common.N("tstproducera"))
	assert.NoError(t, err)

	main.ValidatePushBlock(blocks.second)
}

func TestTrustedProducerVerify2nd(t *testing.T) {
	trustedProducers := treeset.NewWith(common.TypeName, common.CompareName, common.N("defproducera"), common.N("defproducerc"))
	main := NewValidatingTesterTrustedProducers(trustedProducers)
	defer main.close()
	// only using validating_tester to keep the 2 chains in sync, not to validate that the validating_node matches the main node,
	// since it won't be
	main.SkipValidate = true

	// First we create a valid block with valid transaction
	producers := []common.AccountName{common.N("defproducera"),
		common.N("defproducerb"), common.N("defproducerc"), common.N("defproducerd")}

	for _, prod := range producers {
		main.CreateAccount(prod, common.DefaultConfig.SystemAccountName, false, true)
	}
	b := main.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	trace := main.SetProducers(&producers)
	assert.Equal(t, Exception(nil), trace.Except)

	for b.Producer != common.N("defproducerc") {
		b = main.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	}

	blocks, err := CorruptTrxInBlock(main, common.N("tstproducera"))
	assert.NoError(t, err)

	main.ValidatePushBlock(blocks.second)
}

func TestUntrustedProducer(t *testing.T) {
	trustedProducers := treeset.NewWith(common.TypeName, common.CompareName, common.N("defproducera"), common.N("defproducerc"))
	main := NewValidatingTesterTrustedProducers(trustedProducers)
	defer main.close()
	// only using validating_tester to keep the 2 chains in sync, not to validate that the validating_node matches the main node,
	// since it won't be
	main.SkipValidate = true

	// First we create a valid block with valid transaction
	producers := []common.AccountName{common.N("defproducera"),
		common.N("defproducerb"), common.N("defproducerc"), common.N("defproducerd")}

	for _, prod := range producers {
		main.CreateAccount(prod, common.DefaultConfig.SystemAccountName, false, true)
	}
	b := main.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	trace := main.SetProducers(&producers)
	assert.Equal(t, Exception(nil), trace.Except)

	for b.Producer != common.N("defproducerb") {
		b = main.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	}

	blocks, err := CorruptTrxInBlock(main, common.N("tstproducera"))
	assert.NoError(t, err)

	CheckThrowException(t, &UnsatisfiedAuthorization{}, func() { main.ValidatePushBlock(blocks.second) })
}
