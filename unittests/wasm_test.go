package unittests

import (
	"fmt"
	"github.com/eosspark/eos-go/chain"
	//"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/wasmgo"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type assertdef struct {
	Condition int8
	Message   string
}

func (d *assertdef) GetAccount() common.AccountName {
	return common.N("asserter")
}

func (d *assertdef) GetName() common.AccountName {
	return common.N("procassert")
}

type provereset struct{}

func (d *provereset) GetAccount() common.AccountName {
	return common.N("asserter")
}

func (d *provereset) GetName() common.AccountName {
	return common.N("provereset")
}

type actionInterface interface {
	GetAccount() common.AccountName
	GetName() common.AccountName
}

func newAction(permissionLevel []types.PermissionLevel, a actionInterface) *types.Action {

	payload, _ := rlp.EncodeToBytes(a)
	act := types.Action{
		Account:       common.AccountName(a.GetAccount()),
		Name:          common.AccountName(a.GetName()),
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
		// prettyTrx := common.Variants{
		// 	"actions": common.Variants{
		// 		"account": "asserter",
		// 		"name":    "procassert",
		// 		"authorization": common.Variants{
		// 			"actor":      "asserter",
		// 			"permission": "active"},
		// 		"data":common.Variants{
		// 			"condition":1,
		// 			"message":"Should Not Assert"}}}
		//abi_serializer::from_variant(pretty_trx, trx, resolver, abi_serializer_max_time);

		actData := common.Variants{
			"message": common.Variants{
				"condition": 1,
				"message":   "Should Not Assert"}}
		act := b.GetAction(
			asserter, common.N("procassert"),
			[]types.PermissionLevel{{asserter, common.DefaultConfig.ActiveName}},
			&actData)
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

		tester := common.N("f32.tests")

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{tester}, false, true)
		b.ProduceBlocks(1, false)
		b.SetCode(tester, code, nil)
		b.ProduceBlocks(10, false)

		trx := types.SignedTransaction{}
		act := types.Action{
			Account:       tester,
			Name:          common.N(""),
			Authorization: []types.PermissionLevel{{tester, common.DefaultConfig.ActiveName}}}
		trx.Actions = append(trx.Actions, &act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

		privKey := b.getPrivateKey(tester, "active")
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

		f32_tests := common.N("f32.tests")

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

func TestFloat64(t *testing.T) {
	wasm := "test_contracts/f64_test.wasm"
	t.Run(filepath.Base(wasm), func(t *testing.T) {
		code, err := ioutil.ReadFile(wasm)
		if err != nil {
			t.Fatal(err)
		}

		f64_tests := common.N("f_tests")

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{f64_tests}, false, true)
		b.ProduceBlocks(1, false)
		b.SetCode(f64_tests, code, nil)
		b.ProduceBlocks(10, false)

		trx := types.SignedTransaction{}
		act := types.Action{
			Account:       f64_tests,
			Name:          common.N(""),
			Authorization: []types.PermissionLevel{{f64_tests, common.DefaultConfig.ActiveName}}}
		trx.Actions = append(trx.Actions, &act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

		privKey := b.getPrivateKey(f64_tests, "active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey, &chainId)
		b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		b.ProduceBlocks(1, false)

		//trxId := trx.ID()
		//assert.Equal(t, b.ChainHasTransaction(&trxId), true)
		b.close()
	})
}

func TestFloat64Bitwise(t *testing.T) {
	wasm := "test_contracts/f64_test_bitwise.wasm"
	t.Run(filepath.Base(wasm), func(t *testing.T) {
		code, err := ioutil.ReadFile(wasm)
		if err != nil {
			t.Fatal(err)
		}

		f64_tests := common.N("f_tests")

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{f64_tests}, false, true)
		b.ProduceBlocks(1, false)
		b.SetCode(f64_tests, code, nil)
		b.ProduceBlocks(10, false)

		trx := types.SignedTransaction{}
		act := types.Action{
			Account:       f64_tests,
			Name:          common.N(""),
			Authorization: []types.PermissionLevel{{f64_tests, common.DefaultConfig.ActiveName}}}
		trx.Actions = append(trx.Actions, &act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

		privKey := b.getPrivateKey(f64_tests, "active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey, &chainId)
		b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		b.ProduceBlocks(1, false)

		//trxId := trx.ID()
		//assert.Equal(t, b.ChainHasTransaction(&trxId), true)
		b.close()
	})
}

func wast2wasm(wast []uint8) []uint8 {
	wastTmp := "wast_tmp.wast"
	wasmTmp := "wast_tmp.wasm"
	os.Remove(wastTmp)
	os.Remove(wasmTmp)
	ioutil.WriteFile(wastTmp, wast, os.ModePerm)
	cmd := exec.Command("./test_contracts/wat2wasm", wastTmp)
	cmd.Run()
	code, _ := ioutil.ReadFile(wasmTmp)
	os.Remove(wastTmp)
	os.Remove(wasmTmp)
	return code
}

func TestF32F64overflow(t *testing.T) {
	t.Run("", func(t *testing.T) {

		f_tests := common.N("f_tests")
		b := newBaseTester(true, chain.SPECULATIVE)

		var count uint64 = 0
		check := func(wastTemplate string, op string, param string) bool {
			count += 16
			tester := common.AccountName(uint64(f_tests) + count)
			b.CreateAccounts([]common.AccountName{tester}, false, true)
			b.ProduceBlocks(1, false)

			wast := fmt.Sprintf(wastTemplate, op, param)
			wasm := wast2wasm([]byte(wast))
			b.SetCode(tester, wasm, nil)
			b.ProduceBlocks(10, false)

			trx := types.SignedTransaction{}
			act := types.Action{
				Account:       tester,
				Name:          common.N(""),
				Authorization: []types.PermissionLevel{{tester, common.DefaultConfig.ActiveName}}}
			trx.Actions = append(trx.Actions, &act)
			b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

			privKey := b.getPrivateKey(tester, "active")
			chainId := b.Control.GetChainId()
			trx.Sign(&privKey, &chainId)

			returning := false
			try.Try(func() {
				b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
				b.ProduceBlocks(1, false)
				trxId := trx.ID()
				assert.Equal(t, b.ChainHasTransaction(&trxId), true)
				returning = true
			}).Catch(func(e exception.Exception) {
			}).End()
			return returning
		}

		//// float32 => int32
		// 2^31
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_s_f32", "f32.const 2147483648"), false)
		// the maximum value below 2^31 representable in IEEE float32
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_s_f32", "f32.const 2147483520"), true)
		// -2^31
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_s_f32", "f32.const -2147483648"), true)
		// the maximum value below -2^31 in IEEE float32
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_s_f32", "f32.const -2147483904"), false)

		//// float32 => uint32
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_u_f32", "f32.const 0"), true)
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_u_f32", "f32.const -1"), false)
		// max value below 2^32 in IEEE float32
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_u_f32", "f32.const 4294967040"), true)
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_u_f32", "f32.const 4294967296"), false)

		//// double => int32
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_s_f64", "f64.const 2147483648"), false)
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_s_f64", "f64.const 2147483647"), true)
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_s_f64", "f64.const -2147483648"), true)
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_s_f64", "f64.const -2147483649"), false)

		//// double => uint32
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_u_f64", "f64.const 0"), true)
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_u_f64", "f64.const -1"), false)
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_u_f64", "f64.const 4294967295"), true)
		assert.Equal(t, check(i32_overflow_wast, "i32_trunc_u_f64", "f64.const 4294967296"), false)

		//// float32 => int64
		// 2^63
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_s_f32", "f32.const 9223372036854775808"), false)
		// the maximum value below 2^63 representable in IEEE float32
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_s_f32", "f32.const 922337148709896192"), true)
		// -2^63
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_s_f32", "f32.const -9223372036854775808"), true)
		// the maximum value below -2^63 in IEEE float32
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_s_f32", "f32.const -9223373136366403584"), false)

		//// float32 => uint64
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_u_f32", "f32.const -1"), false)
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_u_f32", "f32.const 0"), true)
		// max value below 2^64 in IEEE float32
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_u_f32", "f32.const 18446742974197923840"), true)
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_u_f32", "f32.const 18446744073709551616"), false)

		//// double => int64
		// 2^63
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_s_f64", "f32.const 9223372036854775808"), false)
		// the maximum value below 2^63 representable in IEEE float64
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_s_f64", "f32.const 9223372036854774784"), true)
		// -2^63
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_s_f64", "f64.const -9223372036854775808"), true)
		// the maximum value below -2^63 in IEEE float64
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_s_f64", "f64.const -9223372036854777856"), false)

		//// double => uint64
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_u_f64", "f64.const -1"), false)
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_u_f64", "f64.const 0"), true)
		// max value below 2^64 in IEEE float64
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_u_f64", "f64.const 18446744073709549568"), true)
		assert.Equal(t, check(i64_overflow_wast, "i64_trunc_u_f64", "f64.const 18446744073709551616"), false)

		b.close()

	})
}

func TestMisaligned(t *testing.T) {
	t.Run("", func(t *testing.T) {
		aligncheck := common.N("aligncheck")
		b := newBaseTester(true, chain.SPECULATIVE)
		b.CreateAccounts([]common.AccountName{aligncheck}, false, true)
		b.ProduceBlocks(1, false)

		checkAligned := func(wast string) {

			wasm := wast2wasm([]byte(wast))
			b.SetCode(aligncheck, wasm, nil)
			b.ProduceBlocks(10, false)

			trx := types.SignedTransaction{}
			act := types.Action{
				Account:       aligncheck,
				Name:          common.N(""),
				Authorization: []types.PermissionLevel{{aligncheck, common.DefaultConfig.ActiveName}}}
			trx.Actions = append(trx.Actions, &act)
			b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

			privKey := b.getPrivateKey(aligncheck, "active")
			chainId := b.Control.GetChainId()
			trx.Sign(&privKey, &chainId)
			b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
			b.ProduceBlocks(1, false)
			trxId := trx.ID()
			assert.Equal(t, b.ChainHasTransaction(&trxId), true)
		}

		checkAligned(aligned_ref_wast)
		checkAligned(misaligned_ref_wast)
		checkAligned(aligned_const_ref_wast)
		checkAligned(misaligned_const_ref_wast)

		b.close()

	})
}

func TestWeightedCpuLimit(t *testing.T) {
	t.Run("", func(t *testing.T) {

		b := newBaseTester(true, chain.SPECULATIVE)
		mgr := b.Control.GetMutableResourceLimitsManager()

		f_tests := common.N("f_tests")
		acc2 := common.N("acc2")
		b.CreateAccounts([]common.AccountName{f_tests}, false, true)
		b.CreateAccounts([]common.AccountName{acc2}, false, true)

		//pass := false

		code := `(module
		(import "env" "require_auth" (func $require_auth (param i64)))
		(import "env" "eosio_assert" (func $eosio_assert (param i32 i32)))
		(table 0 anyfunc)
		(memory $0 1)
		(export "apply" (func $apply))
		(func $i64_trunc_u_f64 (param $0 f64) (result i64) (i64.trunc_u/f64 (get_local $0)))
		(func $test (param $0 i64))
		(func $apply (param $0 i64)(param $1 i64)(param $2 i64)`

		for i := 0; i < 1024; i++ {
			code += "(call $test (call $i64_trunc_u_f64 (f64.const 1)))\n"
		}

		code += "))"
		b.ProduceBlocks(1, false)

		wasm := wast2wasm([]byte(code))
		b.SetCode(common.N("f_tests"), wasm, nil)
		b.ProduceBlocks(10, false)

		mgr.SetAccountLimits(f_tests, -1, -1, 1)
		var count int = 0

		for count < 4 {

			trx := types.SignedTransaction{}

			for i := 0; i < 2; i++ {

				actionName := common.ActionName(uint64(common.N("")) + uint64(i*16))
				act := types.Action{
					Account:       f_tests,
					Name:          actionName,
					Authorization: []types.PermissionLevel{{f_tests, common.DefaultConfig.ActiveName}}}
				trx.Actions = append(trx.Actions, &act)
			}
			b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

			privKey := b.getPrivateKey(f_tests, "active")
			chainId := b.Control.GetChainId()
			trx.Sign(&privKey, &chainId)

			returning := false
			try.Try(func() {
				b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
				b.ProduceBlocks(1, false)
				//trxId := trx.ID()
				//assert.Equal(t, b.ChainHasTransaction(&trxId), true)
				count++
			}).Catch(func(e exception.Exception) {
				//if (e.Code() == exception.LeewayDeadlineException{}.Code()) { //catch by check time
				assert.Equal(t, count, 3)
				returning = true
				//}
			}).End()

			if returning {
				break
			}

			//BOOST_REQUIRE_EQUAL(true, validate());
			if count == 2 {
				mgr.SetAccountLimits(acc2, -1, -1, 100000000)
			}
		}

		assert.Equal(t, count, 3)

		b.close()

	})
}

func TestCheckEntryBehavior(t *testing.T) {
	t.Run("", func(t *testing.T) {

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)

		entrycheck := common.N("entrycheck")
		b.CreateAccounts([]common.AccountName{entrycheck}, false, true)
		b.ProduceBlocks(1, false)

		wasm := wast2wasm([]byte(entry_wast))
		b.SetCode(entrycheck, wasm, nil)
		b.ProduceBlocks(10, false)

		trx := types.SignedTransaction{}
		act := types.Action{
			Account:       entrycheck,
			Name:          common.N(""),
			Authorization: []types.PermissionLevel{{entrycheck, common.DefaultConfig.ActiveName}}}
		trx.Actions = append(trx.Actions, &act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

		privKey := b.getPrivateKey(entrycheck, "active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey, &chainId)
		b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		b.ProduceBlocks(1, false)

		//trxId := trx.ID()
		//assert.Equal(t, b.ChainHasTransaction(&trxId), true)

		// receipt := b.GetTransactionReceipt(&trxId)
		// assert.Equal(t, receipt.Status, types.TransactionStatusExecuted)

		b.close()

	})
}

func TestCheckEntryBehavior2(t *testing.T) {
	t.Run("", func(t *testing.T) {

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)

		entrycheck := common.N("entrycheck")
		b.CreateAccounts([]common.AccountName{entrycheck}, false, true)
		b.ProduceBlocks(1, false)

		wasm := wast2wasm([]byte(entry_wast_2))
		b.SetCode(entrycheck, wasm, nil)
		b.ProduceBlocks(10, false)

		trx := types.SignedTransaction{}
		act := types.Action{
			Account:       entrycheck,
			Name:          common.N(""),
			Authorization: []types.PermissionLevel{{entrycheck, common.DefaultConfig.ActiveName}}}
		trx.Actions = append(trx.Actions, &act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

		privKey := b.getPrivateKey(entrycheck, "active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey, &chainId)
		b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		b.ProduceBlocks(1, false)

		//trxId := trx.ID()
		//assert.Equal(t, b.ChainHasTransaction(&trxId), true)

		// receipt := b.GetTransactionReceipt(&trxId)
		// assert.Equal(t, receipt.Status, types.TransactionStatusExecuted)

		b.close()

	})
}

func TestSimpleNoMemoryCheck(t *testing.T) {
	t.Run("", func(t *testing.T) {

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)

		nomem := common.N("nomem")
		b.CreateAccounts([]common.AccountName{nomem}, false, true)
		b.ProduceBlocks(1, false)

		wasm := wast2wasm([]byte(simple_no_memory_wast))
		b.SetCode(nomem, wasm, nil)
		b.ProduceBlocks(1, false)

		trx := types.SignedTransaction{}
		act := types.Action{
			Account:       nomem,
			Name:          common.N(""),
			Authorization: []types.PermissionLevel{{nomem, common.DefaultConfig.ActiveName}}}
		trx.Actions = append(trx.Actions, &act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

		privKey := b.getPrivateKey(nomem, "active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey, &chainId)
		// b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		// b.ProduceBlocks(1, false)

		returning := false
		try.Try(func() {
			b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		}).Catch(func(e exception.Exception) {
			if (e.Code() == exception.WasmExecutionError{}.Code()) { //catch by check time
				returning = true
			}
		}).End()

		assert.Equal(t, returning, true)

		b.close()

	})
}

func TestCheckGlobalReset(t *testing.T) {
	t.Run("", func(t *testing.T) {

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)

		globalreset := common.N("globalreset")
		b.CreateAccounts([]common.AccountName{globalreset}, false, true)
		b.ProduceBlocks(1, false)

		wasm := wast2wasm([]byte(mutable_global_wast))
		b.SetCode(globalreset, wasm, nil)
		b.ProduceBlocks(1, false)

		trx := types.SignedTransaction{}
		act := types.Action{
			Account:       globalreset,
			Name:          common.N(""),
			Authorization: []types.PermissionLevel{{globalreset, common.DefaultConfig.ActiveName}}}
		trx.Actions = append(trx.Actions, &act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

		privKey := b.getPrivateKey(globalreset, "active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey, &chainId)

		returning := false
		try.Try(func() {
			b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		}).Catch(func(e exception.Exception) {
			if (e.Code() == exception.WasmExecutionError{}.Code()) { //catch by check time
				returning = true
			}
		}).End()

		assert.Equal(t, returning, true)
		b.close()

	})
}

func TestStlTest(t *testing.T) {
	t.Run("", func(t *testing.T) {
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)

		stltest := common.N("stltest")
		alice := common.N("alice")
		bob := common.N("bob")

		b.CreateAccounts([]common.AccountName{stltest, alice, bob}, false, true)
		b.ProduceBlocks(1, false)

		wasm := "test_contracts/stltest.wasm"
		abi := "test_contracts/stltest.abi"
		code, _ := ioutil.ReadFile(wasm)
		abiCode, _ := ioutil.ReadFile(abi)

		b.SetCode(stltest, code, nil)
		b.SetAbi(stltest, abiCode, nil)
		b.ProduceBlocks(1, false)

		trx := types.SignedTransaction{}
		actData := common.Variants{
			"message": common.Variants{
				"from":    "bob",
				"to":      "alice",
				"message": "Hi Alice!"}}
		act := b.GetAction(stltest,
			common.N("message"),
			[]types.PermissionLevel{{stltest, common.DefaultConfig.ActiveName}},
			&actData)

		trx.Actions = append(trx.Actions, act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

		privKey := b.getPrivateKey(stltest, "active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey, &chainId)

		b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		trxId := trx.ID()
		assert.Equal(t, b.ChainHasTransaction(&trxId), true)

		b.close()
	})
}

func TestBigMemory(t *testing.T) {
	t.Run("", func(t *testing.T) {
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)

		bigmem := common.N("bigmem")

		b.CreateAccounts([]common.AccountName{bigmem}, false, true)
		b.ProduceBlocks(1, false)

		wast := fmt.Sprintf(biggest_memory_wast, wasmgo.MaximumLinearMemory/(64*1024))
		wasm := wast2wasm([]byte(wast))
		b.SetCode(bigmem, wasm, nil)
		b.ProduceBlocks(1, false)

		trx := types.SignedTransaction{}
		act := types.Action{bigmem, common.N(""), []types.PermissionLevel{{bigmem, common.DefaultConfig.ActiveName}}, nil}

		trx.Actions = append(trx.Actions, &act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

		privKey := b.getPrivateKey(bigmem, "active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey, &chainId)
		b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		b.ProduceBlocks(1, false)

		// trxId := trx.ID()
		// assert.Equal(t, b.ChainHasTransaction(&trxId), true)

		wast = fmt.Sprintf(too_big_memory_wast, wasmgo.MaximumLinearMemory/(64*1024)+1)
		wasm = wast2wasm([]byte(wast))

		returning := false
		try.Try(func() {
			b.SetCode(bigmem, wasm, nil)
		}).Catch(func(e exception.Exception) {
			if (e.Code() == exception.WasmExecutionError{}.Code()) {
				returning = true
			}
		}).End()
		assert.Equal(t, returning, true)

		b.close()
	})
}

func TestTableInit(t *testing.T) {
	t.Run("", func(t *testing.T) {
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)

		tableinit := common.N("tableinit")

		b.CreateAccounts([]common.AccountName{tableinit}, false, true)
		b.ProduceBlocks(1, false)

		wasm := wast2wasm([]byte(valid_sparse_table))
		b.SetCode(tableinit, wasm, nil)
		b.ProduceBlocks(1, false)

		wasm = wast2wasm([]byte(too_big_table))
		returning := false
		try.Try(func() {
			b.SetCode(tableinit, wasm, nil)
		}).Catch(func(e exception.Exception) {
			if (e.Code() == exception.WasmExecutionError{}.Code()) {
				returning = true
			}
		}).End()
		assert.Equal(t, returning, true)

		b.close()
	})
}

func TestMemoryInitBorder(t *testing.T) {
	t.Run("", func(t *testing.T) {
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)

		memoryborder := common.N("memoryborder")

		b.CreateAccounts([]common.AccountName{memoryborder}, false, true)
		b.ProduceBlocks(1, false)

		wasm := wast2wasm([]byte(memory_init_borderline))
		b.SetCode(memoryborder, wasm, nil)
		b.ProduceBlocks(1, false)

		wasm = wast2wasm([]byte(memory_init_toolong))
		returning := false
		try.Try(func() {
			b.SetCode(memoryborder, wasm, nil)
		}).Catch(func(e exception.Exception) {
			if (e.Code() == exception.WasmExecutionError{}.Code()) {
				returning = true
			}
		}).End()
		assert.Equal(t, returning, true)

		wasm = wast2wasm([]byte(memory_init_negative))
		returning = false
		try.Try(func() {
			b.SetCode(memoryborder, wasm, nil)
		}).Catch(func(e exception.Exception) {
			if (e.Code() == exception.WasmExecutionError{}.Code()) {
				returning = true
			}
		}).End()
		assert.Equal(t, returning, true)

		b.close()
	})
}

func TestImports(t *testing.T) {
	t.Run("", func(t *testing.T) {
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		imports := common.N("imports")

		b.CreateAccounts([]common.AccountName{imports}, false, true)
		b.ProduceBlocks(1, false)

		wasm := wast2wasm([]byte(memory_table_import))
		returning := false
		try.Try(func() {
			b.SetCode(imports, wasm, nil)
		}).Catch(func(e exception.Exception) {
			returning = true
		}).End()
		assert.Equal(t, returning, true)

		b.close()
	})
}

func TestNestedLimit(t *testing.T) {
	t.Run("", func(t *testing.T) {
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		nested := common.N("nested")

		b.CreateAccounts([]common.AccountName{nested}, false, true)
		b.ProduceBlocks(1, false)

		nested2 := func(command string) {
			wast := `(module (export "apply" (func $apply)) (func $apply (param $0 i64) (param $1 i64) (param $2 i64)`
			for i := 0; i < 1023; i++ {
				wast += fmt.Sprintf(command, i)
			}
			for i := 0; i < 1023; i++ {
				wast += ")"
			}
			wast += "))"
			wasm := wast2wasm([]byte(wast))
			b.SetCode(nested, wasm, nil)

		}

		nestedException := func(command string) bool {
			wast := `(module (export "apply" (func $apply)) (func $apply (param $0 i64) (param $1 i64) (param $2 i64)`
			for i := 0; i < 1024; i++ {
				wast += fmt.Sprintf(command, i)
			}
			for i := 0; i < 1024; i++ {
				wast += ")"
			}
			wast += "))"
			wasm := wast2wasm([]byte(wast))
			returning := false
			try.Try(func() {
				b.SetCode(nested, wasm, nil)
			}).Catch(func(e exception.Exception) {
				if (e.Code() == exception.WasmExecutionError{}.Code()) {
					returning = true
				}
			}).End()
			return returning
			//assert.Equal(t, returning, true)
		}

		//nested loops
		nested2("(loop (drop (i32.const %d ))")
		ret := nestedException("(loop (drop (i32.const %d ))")
		assert.Equal(t, ret, true)

		//nested blocks
		nested2("(block (drop (i32.const %d ))")
		ret = nestedException("(block (drop (i32.const %d ))")
		assert.Equal(t, ret, true)

		//nested ifs
		nested2("if (i32.wrap/i64 (get_local $0)) (then (drop (i32.const %d ))")
		ret = nestedException("if (i32.wrap/i64 (get_local $0)) (then (drop (i32.const %d ))")
		assert.Equal(t, ret, true)

		// mixed nested
		{
			wast := `(module (export "apply" (func $apply)) (func $apply (param $0 i64) (param $1 i64) (param $2 i64)`
			for i := 0; i < 223; i++ {
				wast += fmt.Sprintf("if (i32.wrap/i64 (get_local $0)) (then (drop (i32.const %d ))", i)
			}
			for i := 0; i < 400; i++ {
				wast += fmt.Sprintf("(block (drop (i32.const %d ))", i)
			}
			for i := 0; i < 400; i++ {
				wast += fmt.Sprintf("(loop (drop (i32.const %d ))", i)
			}
			for i := 0; i < 800; i++ {
				wast += ")"
			}
			for i := 0; i < 223; i++ {
				wast += "))"
			}
			wast += "))"
			wasm := wast2wasm([]byte(wast))
			b.SetCode(nested, wasm, nil)
		}

		{
			wast := `(module (export "apply" (func $apply)) (func $apply (param $0 i64) (param $1 i64) (param $2 i64)`
			for i := 0; i < 224; i++ {
				wast += fmt.Sprintf("if (i32.wrap/i64 (get_local $0)) (then (drop (i32.const %d ))", i)
			}
			for i := 0; i < 400; i++ {
				wast += fmt.Sprintf("(block (drop (i32.const %d ))", i)
			}
			for i := 0; i < 400; i++ {
				wast += fmt.Sprintf("(loop (drop (i32.const %d ))", i)
			}
			for i := 0; i < 800; i++ {
				wast += ")"
			}
			for i := 0; i < 224; i++ {
				wast += "))"
			}
			wast += "))"
			wasm := wast2wasm([]byte(wast))

			returning := false
			try.Try(func() {
				b.SetCode(nested, wasm, nil)
			}).Catch(func(e exception.Exception) {
				if (e.Code() == exception.WasmExecutionError{}.Code()) {
					returning = true
				}
			}).End()
			assert.Equal(t, returning, true)
		}

		b.close()
	})
}

func TestLotsoGlobals(t *testing.T) {
	t.Run("", func(t *testing.T) {
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		nested := common.N("nested")

		b.CreateAccounts([]common.AccountName{nested}, false, true)
		b.ProduceBlocks(1, false)

		wast := `(module (export "apply" (func $apply)) (func $apply (param $0 i64) (param $1 i64) (param $2 i64)`
		for i := 0; i < 85; i++ {
			wast += fmt.Sprintf("(global $g%d (mut i32) (i32.const 0))", i)
			wast += fmt.Sprintf("(global $g%d (mut i64) (i64.const 0))", i+100)
		}
		//that gives us 1020 bytes of mutable globals
		//add a few immutable ones for good measure
		for i := 0; i < 10; i++ {
			wast += fmt.Sprintf("(global $g%d i32 (i32.const 0))", i+100)
		}

		wasm := wast2wasm([]byte(wast + ")"))
		b.SetCode(nested, wasm, nil)

		//1024 should pass
		wasm = wast2wasm([]byte(wast + "(global $z (mut i32) (i32.const -12)))"))
		b.SetCode(nested, wasm, nil)
		//1028 should fail
		wasm = wast2wasm([]byte(wast + "(global $z (mut i64) (i64.const -12)))"))
		returning := false
		try.Try(func() {
			b.SetCode(nested, wasm, nil)
		}).Catch(func(e exception.Exception) {
			if (e.Code() == exception.WasmExecutionError{}.Code()) {
				returning = true
			}
		}).End()
		assert.Equal(t, returning, true)

		b.close()
	})
}

func TestOffsetCheck(t *testing.T) {
	t.Run("", func(t *testing.T) {
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		account := common.N("offsets")

		b.CreateAccounts([]common.AccountName{account}, false, true)
		b.ProduceBlocks(1, false)

		loadops := []string{
			"i32.load", "i64.load", "f32.load", "f64.load", "i32.load8_s", "i32.load8_u",
			"i32.load16_s", "i32.load16_u", "i64.load8_s", "i64.load8_u", "i64.load16_s",
			"i64.load16_u", "i64.load32_s", "i64.load32_u"}

		storeops := [][]string{
			{"i32.store", "i32"},
			{"i64.store", "i64"},
			{"f32.store", "f32"},
			{"f64.store", "f64"},
			{"i32.store8", "i32"},
			{"i32.store16", "i32"},
			{"i64.store8", "i64"},
			{"i64.store16", "i64"},
			{"i64.store32", "i64"}}

		for _, s := range loadops {
			wast := fmt.Sprintf("(module (memory $0 %d ) (func $apply (param $0 i64) (param $1 i64) (param $2 i64)", wasmgo.MaximumLinearMemory/(64*1024))
			wast += fmt.Sprintf("(drop (%s offset=%d (i32.const 0)))", s, wasmgo.MaximumLinearMemory-2)
			wast += `) (export "apply" (func $apply)) )`

			wasm := wast2wasm([]byte(wast))
			b.SetCode(account, wasm, nil)
			b.ProduceBlocks(1, false)
		}

		for _, s := range storeops {
			wast := fmt.Sprintf("(module (memory $0 %d ) (func $apply (param $0 i64) (param $1 i64) (param $2 i64)", wasmgo.MaximumLinearMemory/(64*1024))
			wast += fmt.Sprintf("(%s offset=%d (i32.const 0)( %s.const 0))", s[0], wasmgo.MaximumLinearMemory-2, s[1])
			wast += `) (export "apply" (func $apply)) )`

			wasm := wast2wasm([]byte(wast))
			b.SetCode(account, wasm, nil)
			b.ProduceBlocks(1, false)
		}

		for _, s := range loadops {
			wast := fmt.Sprintf("(module (memory $0 %d ) (func $apply (param $0 i64) (param $1 i64) (param $2 i64)", wasmgo.MaximumLinearMemory/(64*1024))
			wast += fmt.Sprintf("(drop (%s offset=%d (i32.const 0)))", s, wasmgo.MaximumLinearMemory+4)
			wast += `) (export "apply" (func $apply)) )`

			wasm := wast2wasm([]byte(wast))
			returning := false
			try.Try(func() {
				b.SetCode(account, wasm, nil)
			}).Catch(func(e exception.Exception) {
				if (e.Code() == exception.WasmExecutionError{}.Code()) {
					returning = true
				}
			}).End()
			assert.Equal(t, returning, true)
			b.ProduceBlocks(1, false)
		}

		for _, s := range storeops {
			wast := fmt.Sprintf("(module (memory $0 %d ) (func $apply (param $0 i64) (param $1 i64) (param $2 i64)", wasmgo.MaximumLinearMemory/(64*1024))
			wast += fmt.Sprintf("(%s offset=%d (i32.const 0)( %s.const 0))", s[0], wasmgo.MaximumLinearMemory+4, s[1])
			wast += `) (export "apply" (func $apply)) )`

			wasm := wast2wasm([]byte(wast))
			returning := false
			try.Try(func() {
				b.SetCode(account, wasm, nil)
			}).Catch(func(e exception.Exception) {
				if (e.Code() == exception.WasmExecutionError{}.Code()) {
					returning = true
				}
			}).End()
			assert.Equal(t, returning, true)
			b.ProduceBlocks(1, false)
		}

		b.close()
	})
}

func TestNoop(t *testing.T) {
	t.Run("", func(t *testing.T) {
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)

		noop := common.N("noop")
		alice := common.N("alice")

		b.CreateAccounts([]common.AccountName{noop, alice}, false, true)
		b.ProduceBlocks(1, false)

		wasm := "test_contracts/noop.wasm"
		abi := "test_contracts/noop.abi"
		code, _ := ioutil.ReadFile(wasm)
		abiCode, _ := ioutil.ReadFile(abi)

		b.SetCode(noop, code, nil)
		b.SetAbi(noop, abiCode, nil)

		{
			b.ProduceBlocks(5, false)
			trx := types.SignedTransaction{}
			actData := common.Variants{
				"anyaction": common.Variants{
					"from": "noop",
					"type": "some type",
					"data": "some data goes here"}}
			act := b.GetAction(noop,
				common.N("anyaction"),
				[]types.PermissionLevel{{noop, common.DefaultConfig.ActiveName}},
				&actData)

			trx.Actions = append(trx.Actions, act)
			b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

			privKey := b.getPrivateKey(noop, "active")
			chainId := b.Control.GetChainId()
			trx.Sign(&privKey, &chainId)

			b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
			trxId := trx.ID()
			assert.Equal(t, b.ChainHasTransaction(&trxId), true)
		}

		{
			b.ProduceBlocks(5, false)
			trx := types.SignedTransaction{}
			actData := common.Variants{
				"anyaction": common.Variants{
					"from": "alice",
					"type": "some type",
					"data": "some data goes here"}}
			act := b.GetAction(noop,
				common.N("anyaction"),
				[]types.PermissionLevel{{alice, common.DefaultConfig.ActiveName}},
				&actData)

			trx.Actions = append(trx.Actions, act)
			b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

			privKey := b.getPrivateKey(alice, "active")
			chainId := b.Control.GetChainId()
			trx.Sign(&privKey, &chainId)

			b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
			trxId := trx.ID()
			assert.Equal(t, b.ChainHasTransaction(&trxId), true)
		}

		b.close()
	})
}

func TestEosioAbi(t *testing.T) {
	t.Run("", func(t *testing.T) {
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)

		//accnt := b.Control.GetAccount(common.DefaultConfig.SystemAccountName)
		//abi := accnt.GetAbi()
		//abiSerializer := abi_serializer.NewAbiSerializer(abi, b.AbiSerializerMaxTime)

		trx := types.SignedTransaction{}
		alice := common.N("alice")

		ownerAuth := types.NewAuthority(b.getPublicKey(alice, "owner"), uint32(b.AbiSerializerMaxTime))

		pl := []types.PermissionLevel{{common.DefaultConfig.SystemAccountName, common.PermissionName(common.N("active"))}}
		a := chain.NewAccount{
			common.DefaultConfig.SystemAccountName,
			alice,
			ownerAuth,
			types.NewAuthority(b.getPublicKey(alice, "active"), uint32(b.AbiSerializerMaxTime))}

		act := newAction(pl, &a)
		trx.Actions = append(trx.Actions, act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

		privKey := b.getPrivateKey(common.DefaultConfig.SystemAccountName, "active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey, &chainId)

		//result := b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)

		// fc::variant pretty_output;
		// // verify to_variant works on eos native contract type: newaccount
		// // see abi_serializer::to_abi()
		// abi_serializer::to_variant(*result, pretty_output, get_resolver(), abi_serializer_max_time);
		// BOOST_TEST(fc::json::to_string(pretty_output).find("newaccount") != std::string::npos);

		b.close()
	})
}
