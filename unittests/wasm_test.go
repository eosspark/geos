package unittests

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

type assertdef struct {
	Condition int8
	Message   string
}

func (d *assertdef) getAccount() common.AccountName {
	return common.N("asserter")
}

func (d *assertdef) getName() common.AccountName {
	return common.N("procassert")
}

type provereset struct{}

func (d *provereset) getAccount() common.AccountName {
	return common.N("asserter")
}

func (d *provereset) getName() common.AccountName {
	return common.N("provereset")
}

type actionInterface interface {
	getAccount() common.AccountName
	getName() common.AccountName
}

func newAction(permissionLevel []types.PermissionLevel, a actionInterface) *types.Action {

	payload, _ := rlp.EncodeToBytes(a)
	act := types.Action{
		Account:       common.AccountName(a.getAccount()),
		Name:          common.AccountName(a.getName()),
		Data:          payload,
		Authorization: permissionLevel,
	}
	return &act
}

func TestBasic(t *testing.T) {
	name := "test_contracts/asserter.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		asserter := common.N("asserter")
		procassert := common.N("procassert")

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{asserter}, false, true)
		b.ProduceBlocks(1, false)
		b.SetCode(asserter, code, nil)
		b.ProduceBlocks(1, false)

		var noAssertID common.TransactionIdType
		{
			trx := types.SignedTransaction{}
			pl := []types.PermissionLevel{{asserter, common.DefaultConfig.ActiveName}}
			action := assertdef{1, "Should Not Assert!"}
			act := newAction(pl, &action)
			trx.Actions = append(trx.Actions, act)
			b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

			privKey := b.getPrivateKey(asserter, "active")
			chainId := b.Control.GetChainId()
			trx.Sign(&privKey, &chainId)

			result := b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
			assert.Equal(t, result.Receipt.Status, types.TransactionStatusExecuted)
			assert.Equal(t, len(result.ActionTraces), 1)
			assert.Equal(t, result.ActionTraces[0].Receipt.Receiver, asserter)
			assert.Equal(t, result.ActionTraces[0].Act.Account, asserter)
			assert.Equal(t, result.ActionTraces[0].Act.Name, procassert)
			assert.Equal(t, len(result.ActionTraces[0].Act.Authorization), 1)
			assert.Equal(t, result.ActionTraces[0].Act.Authorization[0].Actor, asserter)
			assert.Equal(t, result.ActionTraces[0].Act.Authorization[0].Permission, common.DefaultConfig.ActiveName)

			noAssertID = trx.ID()
		}
		b.ProduceBlocks(1, false)
		assert.Equal(t, b.ChainHasTransaction(&noAssertID), true)
		receipt := b.GetTransactionReceipt(&noAssertID)
		assert.Equal(t, receipt.Status, types.TransactionStatusExecuted)

		var yesAssertID common.TransactionIdType
		{
			trx := types.SignedTransaction{}

			pl := []types.PermissionLevel{{asserter, common.DefaultConfig.ActiveName}}
			action := assertdef{0, "Should Assert!"}
			act := newAction(pl, &action)
			trx.Actions = append(trx.Actions, act)
			b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)
			privKey := b.getPrivateKey(asserter, "active")
			chainId := b.Control.GetChainId()
			trx.Sign(&privKey, &chainId)
			yesAssertID = trx.ID()

			returning := false
			try.Try(func() {
				b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
			}).Catch(func(e exception.Exception) {
				errCode := exception.EosioAssertCodeException{}.Code()
				if e.Code() == errCode {
					returning = true
				}
			}).End()
			assert.Equal(t, returning, true)
		}

		b.ProduceBlocks(1, false)
		hasTx := b.ChainHasTransaction(&yesAssertID)
		assert.Equal(t, hasTx, false)

		b.close()
	})
}

func TestProveMemReset(t *testing.T) {
	name := "test_contracts/asserter.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		asserter := common.N("asserter")

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{asserter}, false, true)
		b.ProduceBlocks(1, false)
		b.SetCode(asserter, code, nil)
		b.ProduceBlocks(1, false)

		for i := 0; i < 5; i++ {
			trx := types.SignedTransaction{}

			pl := []types.PermissionLevel{{asserter, common.DefaultConfig.ActiveName}}
			action := provereset{}
			act := newAction(pl, &action)
			trx.Actions = append(trx.Actions, act)
			b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)
			privKey := b.getPrivateKey(asserter, "active")
			chainId := b.Control.GetChainId()
			trx.Sign(&privKey, &chainId)

			b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
			b.ProduceBlocks(1, false)

			trxId := trx.ID()
			assert.Equal(t, b.ChainHasTransaction(&trxId), true)
			receipt := b.GetTransactionReceipt(&trxId)
			assert.Equal(t, receipt.Status, types.TransactionStatusExecuted)
		}

		b.close()
	})
}

func TestAbiFromVariant(t *testing.T) {
	wasm := "test_contracts/asserter.wasm"
	abi := "test_contracts/asserter.abi"
	t.Run(filepath.Base(wasm), func(t *testing.T) {
		code, err := ioutil.ReadFile(wasm)
		if err != nil {
			t.Fatal(err)
		}

		abiCode, _ := ioutil.ReadFile(abi)
		asserter := common.N("asserter")

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{asserter}, false, true)
		b.ProduceBlocks(1, false)
		b.SetCode(asserter, code, nil)
		b.SetAbi(asserter, abiCode, nil)
		b.ProduceBlocks(1, false)

		trx := types.SignedTransaction{}

		//prettyTrx := common.Variants{
		//	"actions": common.Variants{
		//		"actions": "asserter",
		//		"name":    "procassert",
		//		"authorization": common.Variants{
		//			"actor":      "asserter",
		//			"permission": "active"}}}

		//abi_serializer::from_variant(pretty_trx, trx, resolver, abi_serializer_max_time);
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)
		privKey := b.getPrivateKey(asserter, "active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey, &chainId)
		b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		b.ProduceBlocks(1, false)
		trxId := trx.ID()
		assert.Equal(t, b.ChainHasTransaction(&trxId), true)
		receipt := b.GetTransactionReceipt(&trxId)
		assert.Equal(t, receipt.Status, types.TransactionStatusExecuted)

		b.close()
	})
}

func TestSoftfloat32(t *testing.T) {
	wasm := "test_contracts/f32_test.wasm"
	t.Run(filepath.Base(wasm), func(t *testing.T) {
		code, err := ioutil.ReadFile(wasm)
		if err != nil {
			t.Fatal(err)
		}

		f32_tests := common.N("f32_tests")

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{f32_tests}, false, true)
		b.ProduceBlocks(1, false)
		b.SetCode(f32_tests, code, nil)
		b.ProduceBlocks(10, false)

		trx := types.SignedTransaction{}
		act := types.Action{
			Account:       f32_tests,
			Name:          common.N(""),
			Authorization: []types.PermissionLevel{{f32_tests, common.DefaultConfig.ActiveName}}}
		trx.Actions = append(trx.Actions, &act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

		privKey := b.getPrivateKey(f32_tests, "active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey, &chainId)
		b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		b.ProduceBlocks(1, false)

		//trxId := trx.ID()
		//assert.Equal(t, b.ChainHasTransaction(&trxId), true)

		b.close()
	})
}

func TestErrorfloat32(t *testing.T) {
	wasm := "test_contracts/f32_error.wasm"
	t.Run(filepath.Base(wasm), func(t *testing.T) {
		code, err := ioutil.ReadFile(wasm)
		if err != nil {
			t.Fatal(err)
		}

		f32_tests := common.N("f32_tests")

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{f32_tests}, false, true)
		b.ProduceBlocks(1, false)
		b.SetCode(f32_tests, code, nil)
		b.ProduceBlocks(10, false)

		trx := types.SignedTransaction{}
		act := types.Action{
			Account:       f32_tests,
			Name:          common.N(""),
			Authorization: []types.PermissionLevel{{f32_tests, common.DefaultConfig.ActiveName}}}
		trx.Actions = append(trx.Actions, &act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

		privKey := b.getPrivateKey(f32_tests, "active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey, &chainId)
		b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		b.ProduceBlocks(1, false)

		//trxId := trx.ID()
		//assert.Equal(t, b.ChainHasTransaction(&trxId), true)

		b.close()
	})
}
