// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasmgo_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/eosspark/container/sets/treeset"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/math"
	"github.com/eosspark/eos-go/crypto"
	abi "github.com/eosspark/eos-go/crypto/abi_serializer"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	arithmetic "github.com/eosspark/eos-go/common/arithmetic_types"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/wasmgo"
	"github.com/stretchr/testify/assert"
)

const crypto_api_exception int = 0
const DUMMY_ACTION_DEFAULT_A = 0x45
const DUMMY_ACTION_DEFAULT_B = 0xab11cd1244556677
const DUMMY_ACTION_DEFAULT_C = 0x7451ae12

type dummy_action struct {
	A byte
	B uint64
	C int32
}

func (d *dummy_action) getName() common.AccountName {
	return common.AccountName(common.N("dummy_action"))
}

func (d *dummy_action) getAccount() common.AccountName {
	return common.AccountName(common.N("testapi"))
}

func TestAction(t *testing.T) {
	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.CreateAccounts([]common.AccountName{common.N("acc1")}, false, true)
		b.CreateAccounts([]common.AccountName{common.N("acc2")}, false, true)
		b.CreateAccounts([]common.AccountName{common.N("acc3")}, false, true)
		b.CreateAccounts([]common.AccountName{common.N("acc4")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(1, false)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_action", "assert_true")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})

		retException := callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_action", "assert_false")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		da := dummy_action{DUMMY_ACTION_DEFAULT_A, DUMMY_ACTION_DEFAULT_B, DUMMY_ACTION_DEFAULT_C}
		load, _ := rlp.EncodeToBytes(&da)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_action", "read_action_normal")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})

		load = bytes.Repeat([]byte{byte(0x01)}, 1<<16)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_action", "read_action_to_0")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})
		load = bytes.Repeat([]byte{byte(0x01)}, 1<<16+1)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_action", "read_action_to_0")}, load, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.OverlappingMemoryError{}.Code(), "access violation")
		assert.Equal(t, retException, true)

		load = bytes.Repeat([]byte{byte(0x01)}, 1)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_action", "read_action_to_64k")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})
		load = bytes.Repeat([]byte{byte(0x01)}, 3)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_action", "read_action_to_64k")}, load, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.OverlappingMemoryError{}.Code(), "access violation")
		assert.Equal(t, retException, true)

		b.ProduceBlocks(1, false)
		// test require_notice
		testRequireNotice := func(b *BaseTester, data []byte, scope []common.AccountName) {
			trx := NewTransaction()

			pl := []types.PermissionLevel{{scope[0], common.PermissionName(common.N("active"))}}
			a := testApiAction{wasmTestAction("test_action", "require_notice")}
			act := newAction(pl, &a)
			act.Data = data
			trx.Transaction.Actions = append(trx.Transaction.Actions, act)

			b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)
			privKey := b.getPrivateKey(common.AccountName(common.N("inita")), "active")
			chainId := b.Control.GetChainId()
			trx.Sign(&privKey, &chainId)

			ret := b.PushTransaction(trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
			assert.Equal(t, ret.Receipt.Status, types.TransactionStatusExecuted)

		}
		retException = checkException(testRequireNotice, b, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.UnsatisfiedAuthorization{}.Code(), "transaction declares authority")
		assert.Equal(t, retException, true)

		// test require_auth
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_action", "require_auth")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.MissingAuthException{}.Code(), "missing authority of")
		assert.Equal(t, retException, true)

		a3only := []types.PermissionLevel{{common.AccountName(common.N("acc3")), common.PermissionName(common.N("active"))}}
		load, _ = rlp.EncodeToBytes(&a3only)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_action", "require_auth")}, load, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.MissingAuthException{}.Code(), "missing authority of")
		assert.Equal(t, retException, true)

		a4only := []types.PermissionLevel{{common.AccountName(common.N("acc4")), common.PermissionName(common.N("active"))}}
		load, _ = rlp.EncodeToBytes(&a4only)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_action", "require_auth")}, load, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.MissingAuthException{}.Code(), "missing authority of")
		assert.Equal(t, retException, true)

		//a3a4 := []types.PermissionLevel{{common.AccountName(common.N("acc3")), common.PermissionName(common.N("active"))}, {common.AccountName(common.N("acc4")), common.PermissionName(common.N("active"))}}
		////a3a4Scope := []common.AccountName{common.AccountName(common.N("acc3")), common.AccountName(common.N("acc4"))}
		//{
		//
		//	trx := NewTransaction()
		//	a := testApiAction{wasmTestAction("test_action", "require_notice")}
		//	act := newAction(a3a4, &a)
		//
		//	data, _ := rlp.EncodeToBytes(&a3a4)
		//	act.Data = data
		//
		//	act.Authorization = append(act.Authorization, types.PermissionLevel{common.AccountName(common.N("testapi")), common.PermissionName(common.N("active"))})
		//	trx.Transaction.Actions = append(trx.Transaction.Actions, act)
		//
		//	b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)
		//
		//	chainId := b.Control.GetChainId()
		//	privKey := b.getPrivateKey(common.AccountName(common.N("testapi")), "active")
		//	trx.Sign(&privKey, &chainId)
		//	privKey = b.getPrivateKey(common.AccountName(common.N("acc3")), "active")
		//	trx.Sign(&privKey, &chainId)
		//	privKey = b.getPrivateKey(common.AccountName(common.N("acc4")), "active")
		//	trx.Sign(&privKey, &chainId)
		//
		//	ret := b.PushTransaction(trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		//	assert.Equal(t, ret.Receipt.Status, types.TransactionStatusExecuted)
		//}

		now := b.Control.HeadBlockTime().AddUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
		n := now.TimeSinceEpoch().Count()
		load, _ = rlp.EncodeToBytes(&n)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_action", "test_current_time")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})
		b.ProduceBlocks(1, false)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_action", "test_current_time")}, load, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), "tmp == current_time()")
		assert.Equal(t, retException, true)

		account := common.AccountName(common.N("testapi"))
		load, _ = rlp.EncodeToBytes(&account)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_action", "test_current_receiver")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_transaction", "send_action_sender")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})
		b.ProduceBlocks(1, false)

		now = b.Control.HeadBlockTime().AddUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
		n = now.TimeSinceEpoch().Count()
		load, _ = rlp.EncodeToBytes(&n)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_action", "test_publication_time")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_action", "test_abort")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.AbortCalled{}.Code(), "abort() called")
		assert.Equal(t, retException, true)

		da = dummy_action{DUMMY_ACTION_DEFAULT_A, DUMMY_ACTION_DEFAULT_B, DUMMY_ACTION_DEFAULT_C}
		load, _ = rlp.EncodeToBytes(&da)
		callTestF2(t, b, &da, load, []common.AccountName{common.AccountName(common.N("testapi"))})

		b.close()

	})

}

func checkException(f func(b *BaseTester, data []byte, scope []common.AccountName), b *BaseTester, data []byte, scope []common.AccountName, errCode exception.ExcTypes, errMsg string) (ret bool) {

	returning := false
	try.Try(func() {
		f(b, data, scope)
	}).Catch(func(e exception.Exception) {
		if e.Code() == errCode || inString(e.What(), errMsg) {
			returning = true
		}
	}).End()

	if returning {
		return returning
	}

	return false

}

func TestRequireNoticeTests(t *testing.T) {
	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.CreateAccounts([]common.AccountName{common.N("acc5")}, false, true)
		b.ProduceBlocks(1, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.SetCode(common.AccountName(common.N("acc5")), code, nil)
		b.ProduceBlocks(1, false)

		trx := NewTransaction()
		a := testApiAction{wasmTestAction("test_action", "require_notice_tests")}
		pl := []types.PermissionLevel{{common.AccountName(common.N("testapi")), common.PermissionName(common.N("active"))}}

		act := newAction(pl, &a)
		trx.Transaction.Actions = append(trx.Transaction.Actions, act)
		b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

		privKey := b.getPrivateKey(common.AccountName(common.N("testapi")), "active")
		chainId := b.Control.GetChainId()
		trx.Sign(&privKey, &chainId)
		ret := b.PushTransaction(trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
		assert.Equal(t, ret.Receipt.Status, types.TransactionStatusExecuted)

		b.close()

	})

}

func TestRamBillingInNotifyTests(t *testing.T) {
	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.CreateAccounts([]common.AccountName{common.N("testapi2")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(1, false)
		b.SetCode(common.AccountName(common.N("testapi2")), code, nil)
		b.ProduceBlocks(1, false)

		data := arithmetic.Int128{uint64(common.N("testapi")), uint64(common.N("testapi2"))}
		load, _ := rlp.EncodeToBytes(&data)
		retException := callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_action", "test_ram_billing_in_notify")}, load, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.SubjectiveBlockProductionException{}.Code(), "Cannot charge RAM to other accounts during notify.")
		assert.Equal(t, retException, true)

		data = arithmetic.Int128{0, uint64(common.N("testapi2"))}
		load, _ = rlp.EncodeToBytes(&data)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_action", "test_ram_billing_in_notify")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})

		data = arithmetic.Int128{uint64(common.N("testapi2")), uint64(common.N("testapi2"))}
		load, _ = rlp.EncodeToBytes(&data)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_action", "test_ram_billing_in_notify")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})

		b.close()

	})

}

type cfAction struct {
	Payload uint32
	Cfd_idx uint32
}

func (n *cfAction) getAccount() common.AccountName {
	return common.AccountName(common.N("testapi"))
}

func (n *cfAction) getName() common.AccountName {
	return common.AccountName(common.N("cf_action"))
}

type actionInterface interface {
	getAccount() common.AccountName
	getName() common.AccountName
}

func newSignedTransaction(control *chain.Controller) *types.SignedTransaction {
	trxHeader := types.TransactionHeader{
		Expiration: common.NewTimePointSecTp(control.PendingBlockTime()).AddSec(60), //common.MaxTimePointSec(),
		// RefBlockNum:      4,
		// RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 100000,
		MaxCpuUsageMS:    200,
		DelaySec:         0,
	}

	trx := types.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*types.Action{},
		Actions:               []*types.Action{},
		TransactionExtensions: []*types.Extension{},
		//RecoveryCache:         make(map[ecc.Signature]types.CachedPubKey),
	}

	headBlockId := control.HeadBlockId()
	trx.SetReferenceBlock(&headBlockId)
	signedTrx := types.NewSignedTransaction(&trx, []ecc.Signature{}, []common.HexBytes{})

	return signedTrx
}

func pushSignedTransaction(control *chain.Controller, trx *types.SignedTransaction) *types.TransactionTrace {
	metaTrx := types.NewTransactionMetadataBySignedTrx(trx, common.CompressionNone)
	return control.PushTransaction(metaTrx, common.TimePoint(common.MaxMicroseconds()), 0)
}

func TestContextFreeAction(t *testing.T) {
	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		control := startBlock()
		createNewAccount2(control, "testapi", "eosio")
		createNewAccount2(control, "dummy", "eosio")

		SetCode(control, "testapi", code)
		trx := newSignedTransaction(control)

		// need at least one normal action
		ret := pushSignedTransaction(control, trx)
		assert.Equal(t, ret.Except.Code(), exception.TxNoAuths{}.Code())

		cfa := cfAction{100, 0}
		act := newAction([]types.PermissionLevel{}, &cfa)
		trx.Transaction.ContextFreeActions = append(trx.Transaction.ContextFreeActions, act)
		var raw uint32 = 100
		data, _ := rlp.EncodeToBytes(raw)
		trx.ContextFreeData = append(trx.ContextFreeData, data)
		raw = 200
		data, _ = rlp.EncodeToBytes(raw)
		trx.ContextFreeData = append(trx.ContextFreeData, data)
		// signing a transaction with only context_free_actions should not be allowed
		ret = pushSignedTransaction(control, trx)
		assert.Equal(t, ret.Except.Code(), exception.TxNoAuths{}.Code())

		da := dummy_action{DUMMY_ACTION_DEFAULT_A, DUMMY_ACTION_DEFAULT_B, DUMMY_ACTION_DEFAULT_C}
		permissions := []types.PermissionLevel{
			types.PermissionLevel{common.AccountName(common.N("testapi")), common.PermissionName(common.N("active"))},
		}
		act = newAction(permissions, &da)
		trx.Transaction.Actions = append(trx.Transaction.Actions, act)

		privateKeys := []*ecc.PrivateKey{
			getPrivateKey("testapi", "active"),
		}
		chainIdType := control.GetChainId()
		for _, privateKey := range privateKeys {
			trx.Sign(privateKey, &chainIdType)
		}
		// add a normal action along with cfa
		ret = pushSignedTransaction(control, trx)

		da = dummy_action{DUMMY_ACTION_DEFAULT_A, 200, DUMMY_ACTION_DEFAULT_C}
		act = newAction(permissions, &da)
		trx.Transaction.Actions = []*types.Action{}
		trx.Transaction.Actions = append(trx.Transaction.Actions, act)

		trx.Signatures = []ecc.Signature{}
		for _, privateKey := range privateKeys {
			trx.Sign(privateKey, &chainIdType)
		}
		// attempt to access context free api in non context free action
		ret = pushSignedTransaction(control, trx)
		assert.Equal(t, ret.Except.Code(), exception.UnaccessibleApi{}.Code())

		act = newAction(permissions, &da)
		trx = newSignedTransaction(control)
		trx.Transaction.ContextFreeActions = append(trx.Transaction.ContextFreeActions, act)
		raw = 100
		data, _ = rlp.EncodeToBytes(raw)
		trx.ContextFreeData = append(trx.ContextFreeData, data)
		raw = 200
		data, _ = rlp.EncodeToBytes(raw)
		trx.ContextFreeData = append(trx.ContextFreeData, data)
		trx.Transaction.Actions = append(trx.Transaction.Actions, act)
		for i := 200; i <= 211; i++ {
			trx.Transaction.ContextFreeActions = []*types.Action{}
			trx.ContextFreeData = []common.HexBytes{}
			cfa.Payload = uint32(i)
			cfa.Cfd_idx = 1
			cfa_act := newAction([]types.PermissionLevel{}, &cfa)

			trx.Transaction.ContextFreeActions = append(trx.Transaction.ContextFreeActions, cfa_act)
			trx.Signatures = []ecc.Signature{}
			for _, privateKey := range privateKeys {
				trx.Sign(privateKey, &chainIdType)
			}

			// attempt to access non context free api
			ret := pushSignedTransaction(control, trx)
			assert.Equal(t, ret.Except.Code(), exception.UnaccessibleApi{}.Code())
		}

		ret = callTestFunction2(control, "test_transaction", "send_cf_action", []byte{}, "testapi")
		assert.Equal(t, len(ret.ActionTraces), 1)
		assert.Equal(t, len(ret.ActionTraces[0].InlineTraces), 1)
		assert.Equal(t, ret.ActionTraces[0].InlineTraces[0].Receipt.Receiver, common.AccountName(common.N("dummy")))
		assert.Equal(t, ret.ActionTraces[0].InlineTraces[0].Act.Account, common.AccountName(common.N("dummy")))
		assert.Equal(t, ret.ActionTraces[0].InlineTraces[0].Act.Name, common.AccountName(common.N("event1")))
		assert.Equal(t, len(ret.ActionTraces[0].InlineTraces[0].Act.Authorization), 0)

		retException := callTestFunctionException2(control, "test_transaction", "send_cf_action_fail", []byte{}, "testapi", exception.EosioAssertMessageException{}.Code(), "context free actions cannot have authorizations")
		assert.Equal(t, retException, true)

		trx1 := newSignedTransaction(control)
		cfa = cfAction{100, 0}
		act = newAction([]types.PermissionLevel{}, &cfa)
		raw = 100
		data, _ = rlp.EncodeToBytes(raw)
		trx1.ContextFreeData = append(trx1.ContextFreeData, data)
		trx1.Transaction.ContextFreeActions = append(trx1.Transaction.ContextFreeActions, act)

		trx2 := newSignedTransaction(control)
		raw = 200
		data, _ = rlp.EncodeToBytes(raw)
		trx2.ContextFreeData = append(trx2.ContextFreeData, data)
		trx2.Transaction.ContextFreeActions = append(trx2.Transaction.ContextFreeActions, act)
		chainIdType = control.GetChainId()
		privKey := getPrivateKey("dummy", "active")
		assert.Equal(t, trx1.Sign(privKey, &chainIdType).String() == trx2.Sign(privKey, &chainIdType).String(), false)

		stopBlock(control)

	})

}

type newAccount struct {
	Creator common.AccountName
	Name    common.AccountName
	Owner   types.Authority
	Active  types.Authority
}

func (n *newAccount) getAccount() common.AccountName {
	return common.AccountName(common.DefaultConfig.SystemAccountName)
}

func (n *newAccount) getName() common.AccountName {
	return common.AccountName(common.N("newaccount"))
}

type testApiAction struct {
	actionName uint64
}

func (a *testApiAction) getAccount() common.AccountName {
	return common.AccountName(common.N("testapi"))
}

func (a *testApiAction) getName() common.AccountName {
	return common.AccountName(a.actionName)
}

func TestStatefulApi(t *testing.T) {
	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		control := startBlock()
		createNewAccount2(control, "testapi", "eosio")
		SetCode(control, "testapi", code)

		creator := common.AccountName(common.N("eosio"))
		name := "testapi2"
		c := newAccount{
			Creator: creator,
			Name:    common.AccountName(common.N(name)),
			Owner: types.Authority{
				Threshold: 1,
				Keys:      []types.KeyWeight{{Key: *getPublicKey(name, "owner"), Weight: 1}},
			},
			Active: types.Authority{
				Threshold: 1,
				Keys:      []types.KeyWeight{{Key: *getPublicKey(name, "active"), Weight: 1}},
			},
		}
		permissions := []types.PermissionLevel{
			types.PermissionLevel{creator, common.PermissionName(common.N("active"))},
		}

		act := newAction(permissions, &c)
		trx := newSignedTransaction(control)
		trx.Transaction.Actions = append(trx.Transaction.Actions, act)

		da := testApiAction{wasmTestAction("test_transaction", "stateful_api")}
		act = newAction([]types.PermissionLevel{}, &da)
		trx.Transaction.ContextFreeActions = append(trx.Transaction.ContextFreeActions, act)
		privKey := getPrivateKey("eosio", "active")
		chainIdType := control.GetChainId()
		trx.Sign(privKey, &chainIdType)

		ret := pushSignedTransaction(control, trx)
		assert.Equal(t, ret.Except.Code(), exception.UnaccessibleApi{}.Code())

		stopBlock(control)

	})

}

// func TestStatefulApi(t *testing.T) {
// 	name := "testdata_context/test_api.wasm"
// 	t.Run(filepath.Base(name), func(t *testing.T) {
// 		code, err := ioutil.ReadFile(name)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		control := startBlock()

//         ret := callTestFunction2(control, "test_checktime", "checktime_pass", []byte{}, "testapi")

// 		stopBlock(control)
// 	})

// }

func TestTransaction(t *testing.T) {
	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(100, false)
		b.SetCode(common.N("testapi"), code, nil)
		b.ProduceBlocks(1, false)

		{
			trx := NewTransaction()
			a := testApiAction{wasmTestAction("test_transaction", "require_auth")}
			pl := []types.PermissionLevel{}
			act := newAction(pl, &a)
			trx.Transaction.Actions = append(trx.Transaction.Actions, act)
			b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

			returning := false
			try.Try(func() {
				b.PushTransaction(trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
			}).Catch(func(e exception.Exception) {
				if inString(exception.GetDetailMessage(e), "transaction must have at least one authorization") {
					returning = true
				}
			}).End()
			assert.Equal(t, returning, true)
		}

		callTestF2(t, b, &testApiAction{wasmTestAction("test_transaction", "send_action")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_transaction", "send_action_empty")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})

		retException := callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_transaction", "send_action_large")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.InlineActionTooBig{}.Code(), "inline action too big")
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_transaction", "send_action_inline_fail")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), "test_action::assert_false")
		assert.Equal(t, retException, true)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_transaction", "send_transaction")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_transaction", "send_transaction_empty")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.TxNoAuths{}.Code(), "transaction must have at least one authorization")
		assert.Equal(t, retException, true)

		// {
		// 	produce_blocks(10);
		// 	transaction_trace_ptr trace;
		// 	auto c = control->applied_transaction.connect([&]( const transaction_trace_ptr& t) { if (t && t->receipt && t->receipt->status != transaction_receipt::executed) { trace = t; } } );

		// 	// test error handling on deferred transaction failure
		// 	CALL_TEST_FUNCTION(*this, "test_transaction", "send_transaction_trigger_error_handler", {});

		// 	BOOST_REQUIRE(trace);
		// 	BOOST_CHECK_EQUAL(trace->receipt->status, transaction_receipt::soft_fail);
		// 	c.disconnect();
		// }
		ret := callTestF2(t, b, &testApiAction{wasmTestAction("test_transaction", "test_read_transaction")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		assert.Equal(t, ret.ID.String(), ret.ActionTraces[0].Console)

		bn := b.Control.HeadBlockNum()
		load, _ := rlp.EncodeToBytes(&bn)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_transaction", "test_tapos_block_num")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})

		hh := b.Control.HeadBlockId().Hash[1]
		load, _ = rlp.EncodeToBytes(&hh)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_transaction", "test_tapos_block_prefix")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_transaction", "send_action_recurse")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.TransactionException{}.Code(), "max inline action depth per transaction reached")
		assert.Equal(t, retException, true)

		b.close()

	})

}

func TestChain(t *testing.T) {
	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)

		producers := []common.AccountName{
			common.N("inita"),
			common.N("initb"),
			common.N("initc"),
			common.N("initd"),
			common.N("inite"),
			common.N("initf"),
			common.N("initg"),
			common.N("inith"),
			common.N("initi"),
			common.N("initj"),
			common.N("initk"),
			common.N("initl"),
			common.N("initm"),
			common.N("initn"),
			common.N("inito"),
			common.N("initp"),
			common.N("initq"),
			common.N("initr"),
			common.N("inits"),
			common.N("initt"),
			common.N("initu"),
		}

		b.CreateAccounts(producers, false, true)
		b.SetProducers(producers)

		b.SetCode(common.N("testapi"), code, nil)
		b.ProduceBlocks(10, false)

		ps := b.Control.ActiveProducers().Producers
		prods := make([]common.AccountName, len(ps))

		for i := 0; i < len(prods); i++ {
			prods[i] = ps[i].ProducerName
			fmt.Println("prod", i, " ", common.S(uint64(prods[i])))
		}

		load, _ := rlp.EncodeToBytes(&prods)
		ret := callTestF2(t, b, &testApiAction{wasmTestAction("test_chain", "test_activeprods")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})
		fmt.Println(ret.ActionTraces[0].Console)

		b.close()

	})

}

func TestPrint(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(10, false)

		ret := callTestF2(t, b, &testApiAction{wasmTestAction("test_print", "test_prints")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		retCnsl := ret.ActionTraces[0].Console
		assert.Equal(t, retCnsl, "abcefg")

		ret = callTestF2(t, b, &testApiAction{wasmTestAction("test_print", "test_prints_l")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		retCnsl = ret.ActionTraces[0].Console
		assert.Equal(t, retCnsl, "abatest")

		ret = callTestF2(t, b, &testApiAction{wasmTestAction("test_print", "test_printi")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		retCnsl = ret.ActionTraces[0].Console
		assert.Equal(t, retCnsl[0:1], string(strconv.FormatInt(0, 10)))
		assert.Equal(t, retCnsl[1:7], string(strconv.FormatInt(556644, 10)))
		assert.Equal(t, retCnsl[7:9], string(strconv.FormatInt(-1, 10)))

		ret = callTestF2(t, b, &testApiAction{wasmTestAction("test_print", "test_printui")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		retCnsl = ret.ActionTraces[0].Console
		assert.Equal(t, retCnsl[0:1], string(strconv.FormatInt(0, 10)))
		assert.Equal(t, retCnsl[1:7], string(strconv.FormatInt(556644, 10)))
		v := -1
		assert.Equal(t, retCnsl[7:len(retCnsl)], string(strconv.FormatUint(uint64(v), 10))) //-1 / 1844674407370955161

		ret = callTestF2(t, b, &testApiAction{wasmTestAction("test_print", "test_printn")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		retCnsl = ret.ActionTraces[0].Console
		assert.Equal(t, retCnsl[0:5], "abcde")
		assert.Equal(t, retCnsl[5:10], "ab.de")
		assert.Equal(t, retCnsl[10:16], "1q1q1q")
		assert.Equal(t, retCnsl[16:27], "abcdefghijk")
		assert.Equal(t, retCnsl[27:39], "abcdefghijkl")
		assert.Equal(t, retCnsl[39:52], "abcdefghijkl1")
		assert.Equal(t, retCnsl[52:65], "abcdefghijkl1")
		assert.Equal(t, retCnsl[65:78], "abcdefghijkl1")

		ret = callTestF2(t, b, &testApiAction{wasmTestAction("test_print", "test_printi128")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		retCnsl = ret.ActionTraces[0].Console
		s := strings.Split(retCnsl, "\n")
		assert.Equal(t, s[0], "1")
		assert.Equal(t, s[1], "0")
		assert.Equal(t, s[2], "-170141183460469231731687303715884105728")
		assert.Equal(t, s[3], "-87654323456")

		ret = callTestF2(t, b, &testApiAction{wasmTestAction("test_print", "test_printui128")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		retCnsl = ret.ActionTraces[0].Console
		s = strings.Split(retCnsl, "\n")
		assert.Equal(t, s[0], "340282366920938463463374607431768211455")
		assert.Equal(t, s[1], "0")
		assert.Equal(t, s[2], "87654323456")

		ret = callTestF2(t, b, &testApiAction{wasmTestAction("test_print", "test_printsf")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		retCnsl = ret.ActionTraces[0].Console
		s = strings.Split(retCnsl, "\n")
		assert.Equal(t, s[0], "5.000000e-01")
		assert.Equal(t, s[1], "-3.750000e+00")
		assert.Equal(t, s[2], "6.666667e-07")

		ret = callTestF2(t, b, &testApiAction{wasmTestAction("test_print", "test_printdf")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		retCnsl = ret.ActionTraces[0].Console
		s = strings.Split(retCnsl, "\n")
		assert.Equal(t, s[0], "5.000000000000000e-01")
		assert.Equal(t, s[1], "-3.750000000000000e+00")
		assert.Equal(t, s[2], "6.666666666666666e-07")

		//ret = callTestF2(t, b, &testApiAction{wasmTestAction("test_print", "test_printqf")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		//retCnsl = ret.ActionTraces[0].Console
		//s = strings.Split(retCnsl, "\n")
		//assert.Equal(t, s[0], "5.000000000000000000e-01")
		//assert.Equal(t, s[1], "-3.750000000000000000e+00")
		//assert.Equal(t, s[2], "6.666666666666666667e-07")

		b.close()

	})

}

//func TestEosioSystem(t *testing.T) {
//
//	name := "testdata_context/eosio.system.wasm"
//	t.Run(filepath.Base(name), func(t *testing.T) {
//		code, err := ioutil.ReadFile(name)
//		if err != nil {
//			t.Fatal(err)
//		}
//
//		b := newBaseTester(true, chain.SPECULATIVE)
//		b.ProduceBlocks(2, false)
//		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
//		b.ProduceBlocks(10, false)
//		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
//		b.ProduceBlocks(10, false)
//
//		ret := callTestF2(t, b, &testApiAction{wasmTestAction("", "set_abi")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
//		retCnsl := ret.ActionTraces[0].Console
//		assert.Equal(t, retCnsl, "abcefg")
//
//		b.close()
//
//	})
//
//}

func TestTypes(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(10, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(10, false)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_types", "types_size")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_types", "char_to_symbol")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_types", "string_to_name")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_types", "name_class")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})

		b.close()

	})

}

func TestMemory(t *testing.T) {

	name := "testdata_context/test_api_mem.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(10, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(10, false)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_memory", "test_memory_allocs")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		b.ProduceBlocks(10, false)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_memory", "test_memory_hunk")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		b.ProduceBlocks(10, false)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_memory", "test_memory_hunks")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		b.ProduceBlocks(10, false)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_memory", "test_memset_memcpy")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		b.ProduceBlocks(10, false)
		retException := callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_memory", "test_memcpy_overlap_start")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.OverlappingMemoryError{}.Code(), exception.OverlappingMemoryError{}.What())
		assert.Equal(t, retException, true)
		b.ProduceBlocks(10, false)

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_memory", "test_memcpy_overlap_end")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.OverlappingMemoryError{}.Code(), exception.OverlappingMemoryError{}.What())
		assert.Equal(t, retException, true)
		b.ProduceBlocks(10, false)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_memory", "test_memcmp")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		b.ProduceBlocks(10, false)

		testMemoryOob := func(t *testing.T, f string) {
			returning := false
			try.Try(func() {
				callTestF2(t, b, &testApiAction{wasmTestAction("test_memory", f)}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
			}).Catch(func(e exception.Exception) {
				returning = true
			}).End()
			assert.Equal(t, returning, true)

		}

		testMemoryOob(t, "test_outofbound_0")
		testMemoryOob(t, "test_outofbound_1")
		testMemoryOob(t, "test_outofbound_2")
		testMemoryOob(t, "test_outofbound_3")
		testMemoryOob(t, "test_outofbound_4")
		testMemoryOob(t, "test_outofbound_5")
		testMemoryOob(t, "test_outofbound_6")
		testMemoryOob(t, "test_outofbound_7")
		testMemoryOob(t, "test_outofbound_8")
		testMemoryOob(t, "test_outofbound_9")
		testMemoryOob(t, "test_outofbound_10")
		testMemoryOob(t, "test_outofbound_11")
		testMemoryOob(t, "test_outofbound_12")
		testMemoryOob(t, "test_outofbound_13")

		b.close()
	})

}

//extended_memory_test_initial_memory

func TestExtendedMemoryTestInitial(t *testing.T) {

	name := "testdata_context/test_api_mem.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(10, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(10, false)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_extended_memory", "test_initial_buffer")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		b.close()
	})
}

func TestExtendedMemoryTestPage(t *testing.T) {

	name := "testdata_context/test_api_mem.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(10, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(10, false)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_extended_memory", "test_page_memory")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		b.close()
	})
}

func TestExtendedMemoryTestPageExceed(t *testing.T) {

	name := "testdata_context/test_api_mem.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(10, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(10, false)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_extended_memory", "test_page_memory_exceeded")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		b.close()
	})
}

func TestExtendedMemoryTestPageNegativeBytes(t *testing.T) {

	name := "testdata_context/test_api_mem.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(10, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(10, false)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_extended_memory", "test_page_memory_negative_bytes")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		b.close()
	})
}

type check_auth struct {
	Account    common.AccountName
	Permission common.PermissionName
	Pubkeys    []ecc.PublicKey
}

func TestPermission(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(1, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(1, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(1, false)

		getResultInt64 := func() uint64 {

			tab := entity.TableIdObject{
				Code:  common.AccountName(common.N("testapi")),
				Scope: common.ScopeName(common.N("testapi")),
				Table: common.TableName(common.N("testapi")),
			}
			err := b.Control.DB.Find("byCodeScopeTable", tab, &tab)
			try.EosAssert(err == nil, &exception.AssertException{}, "Table id not found")

			obj := entity.KeyValueObject{TId: tab.ID}
			idx, _ := b.Control.DB.GetIndex("byScopePrimary", &obj)
			itr, _ := idx.LowerBound(&obj)
			objLowerbound := entity.KeyValueObject{}
			itr.Data(&objLowerbound)
			try.EosAssert(!idx.CompareEnd(itr) && objLowerbound.TId == tab.ID, &exception.AssertException{}, "lower_bound failed")
			try.EosAssert(len(objLowerbound.Value) > 0, &exception.AssertException{}, "unexpected result size")

			var ret uint64
			rlp.DecodeBytes(objLowerbound.Value, &ret)
			return ret
		}

		checkAuth := check_auth{
			common.AccountName(common.N("testapi")),
			common.PermissionName(common.N("active")),
			[]ecc.PublicKey{b.getPublicKey(common.N("testapi"), "active")},
		}
		load, _ := rlp.EncodeToBytes(checkAuth)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_permission", "check_authorization")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})
		assert.Equal(t, uint64(1), getResultInt64())

		publicKey, _ := ecc.NewPublicKey("EOS7GfRtyDWWgxV88a5TRaYY59XmHptyfjsFmHHfioGNJtPjpSmGX")
		checkAuth = check_auth{
			common.AccountName(common.N("testapi")),
			common.PermissionName(common.N("active")),
			[]ecc.PublicKey{publicKey},
		}
		load, _ = rlp.EncodeToBytes(checkAuth)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_permission", "check_authorization")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})
		assert.Equal(t, uint64(0), getResultInt64())

		checkAuth = check_auth{
			common.AccountName(common.N("testapi")),
			common.PermissionName(common.N("active")),
			[]ecc.PublicKey{b.getPublicKey(common.N("testapi"), "active"), publicKey},
		}
		load, _ = rlp.EncodeToBytes(checkAuth)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_permission", "check_authorization")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})
		assert.Equal(t, uint64(0), getResultInt64())

		checkAuth = check_auth{
			common.AccountName(common.N("noname")),
			common.PermissionName(common.N("active")),
			[]ecc.PublicKey{b.getPublicKey(common.N("testapi"), "active")},
		}
		load, _ = rlp.EncodeToBytes(checkAuth)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_permission", "check_authorization")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})
		assert.Equal(t, uint64(0), getResultInt64())

		checkAuth = check_auth{
			common.AccountName(common.N("testapi")),
			common.PermissionName(common.N("active")),
			[]ecc.PublicKey{},
		}
		load, _ = rlp.EncodeToBytes(checkAuth)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_permission", "check_authorization")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})
		assert.Equal(t, uint64(0), getResultInt64())

		checkAuth = check_auth{
			common.AccountName(common.N("testapi")),
			common.PermissionName(common.N("noname")),
			[]ecc.PublicKey{b.getPublicKey(common.N("testapi"), "active")},
		}
		load, _ = rlp.EncodeToBytes(checkAuth)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_permission", "check_authorization")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})
		assert.Equal(t, uint64(0), getResultInt64())

		b.close()

	})
}

func TestCrypto(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(10, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(10, false)

		trx := NewTransaction()
		//pl := []types.PermissionLevel{
		//	types.PermissionLevel{common.AccountName(common.N("testapi")), common.PermissionName(common.N("active"))},
		//}
		//act := newAction(pl, &testApiAction{wasmTestAction("test_crypto", "test_recover_key")})

		// payload, err := hex.DecodeString("88e4b25a00006c08ac5b595b000000000000")
		// trx.ContextFreeData = []common.HexBytes{payload}

		privKey := b.getPrivateKey(common.N("testapi"), "active")
		chainId := b.Control.GetChainId()
		signatures := trx.Sign(&privKey, &chainId)

		b.ProduceBlocks(1, false)

		digest := trx.Transaction.SigDigest(&chainId, []common.HexBytes{})
		load, _ := rlp.EncodeToBytes(digest)
		//load := digest
		//fmt.Println("digest:", hex.EncodeToString(load))
		//fmt.Println("digest:", load)

		pk := b.getPublicKey(common.N("testapi"), "active")
		p, _ := rlp.EncodeToBytes(pk)
		load = append(load, p...)

		//fmt.Println("publickey:", hex.EncodeToString(p))
		//fmt.Println("publickey:", p)

		sig, _ := rlp.EncodeToBytes(signatures)
		load = append(load, sig...)
		//fmt.Println("sig:", hex.EncodeToString(sig))
		//fmt.Println("sig:", sig)

		//fmt.Println("load:", hex.EncodeToString(load))

		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "test_recover_key")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "test_recover_key_assert_true")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})

		load[len(load)-1] = 0
		retException := callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_crypto", "test_recover_key_assert_false")}, load, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.CryptoApiException{}.Code(), exception.CryptoApiException{}.What())
		assert.Equal(t, retException, true)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "test_sha1")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "test_sha256")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "test_sha512")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "test_ripemd160")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "sha1_no_data")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "sha256_no_data")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "sha512_no_data")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "ripemd160_no_data")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "assert_sha256_true")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "assert_sha1_true")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "assert_sha512_true")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_crypto", "assert_ripemd160_true")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_crypto", "assert_sha256_false")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.CryptoApiException{}.Code(), exception.CryptoApiException{}.What())
		assert.Equal(t, retException, true)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_crypto", "assert_sha1_false")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.CryptoApiException{}.Code(), exception.CryptoApiException{}.What())
		assert.Equal(t, retException, true)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_crypto", "assert_sha512_false")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.CryptoApiException{}.Code(), exception.CryptoApiException{}.What())
		assert.Equal(t, retException, true)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_crypto", "assert_ripemd160_false")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.CryptoApiException{}.Code(), exception.CryptoApiException{}.What())
		assert.Equal(t, retException, true)

		b.close()

	})
}

func TestFixedPoint(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(10, false)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_fixedpoint", "create_instances")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_fixedpoint", "test_addition")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_fixedpoint", "test_subtraction")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_fixedpoint", "test_multiplication")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_fixedpoint", "test_division")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})

		retException := callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_fixedpoint", "test_division_by_0")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		b.close()

	})
}

type testPermissionLastUsedAction struct {
	Account      common.AccountName
	Permission   common.PermissionName
	LastUsedTime common.TimePoint
}

func TestAccountCreationTime(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(1, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(1, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(1, false)
		b.CreateAccounts([]common.AccountName{common.N("alice")}, false, true)
		aliceCreationTime := b.Control.PendingBlockTime()
		b.ProduceBlocks(10, false)

		usedAction := testPermissionLastUsedAction{
			common.N("alice"),
			common.N("active"),
			aliceCreationTime,
		}
		load, _ := rlp.EncodeToBytes(usedAction)
		callTestF2(t, b, &testApiAction{wasmTestAction("test_permission", "test_account_creation_time")}, load, []common.AccountName{common.AccountName(common.N("testapi"))})

		b.close()

	})
}

func callTestException(control *chain.Controller, cls string, method string, payload []byte, authorizer string, billedCpuTimeUs uint32, max_cpu_usage_ms int64, errCode exception.ExcTypes, errMsg string) bool {

	//wasm := wasmgo.NewWasmGo()
	action := wasmTestAction(cls, method)
	fmt.Println(cls, method, action)

	act := types.Action{
		Account:       common.AccountName(common.N(authorizer)),
		Name:          common.ActionName(action),
		Data:          payload,
		Authorization: []types.PermissionLevel{types.PermissionLevel{Actor: common.AccountName(common.N(authorizer)), Permission: common.PermissionName(common.N("active"))}},
	}

	privateKeys := []*ecc.PrivateKey{getPrivateKey(authorizer, "active")}
	trx := newTransaction(control, &act, privateKeys)
	ret := pushTransactionForCallTest(control, trx, billedCpuTimeUs, max_cpu_usage_ms)

	return ret.Except.Code() == errCode

}

func pushTransactionForCallTest(control *chain.Controller, trx *types.TransactionMetadata, billedCpuTimeUs uint32, max_cpu_usage_ms int64) *types.TransactionTrace {
	return control.PushTransaction(trx, common.Now()+common.TimePoint(common.Milliseconds(max_cpu_usage_ms)), billedCpuTimeUs)
}

func TestChecktimePass(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(1, false)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_pass")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})

		b.close()

	})
}

func TestChecktimeFail(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(1, false)

		//var x, cpu, net int64
		//b.Control.GetMutableResourceLimitsManager().GetAccountLimits(common.AccountName(common.N("testapi")), &x, &cpu, &net)
		//fmt.Println("ram:", x, " cpu:", cpu, " net:", net)

		ret := callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 200, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		assert.Equal(t, ret, true)

		ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 0, 200, exception.TxCpuUsageExceeded{}.Code(), exception.TxCpuUsageExceeded{}.What())
		assert.Equal(t, ret, true)
		//
		//ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 0, 200, exception.BlockCpuUsageExceeded{}.Code(), exception.BlockCpuUsageExceeded{}.What())
		//assert.Equal(t, ret, true)
		//
		//ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_sha1_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		//assert.Equal(t, ret, true)
		//
		//ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_assert_sha1_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		//assert.Equal(t, ret, true)
		//
		//ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_sha256_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		//assert.Equal(t, ret, true)
		//
		//ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_assert_sha256_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		//assert.Equal(t, ret, true)
		//
		//ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_assert_sha512_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		//assert.Equal(t, ret, true)
		//
		//ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_ripemd160_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		//assert.Equal(t, ret, true)
		//
		//ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_sha1_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		//assert.Equal(t, ret, true)
		//
		//ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_assert_ripemd160_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		//assert.Equal(t, ret, true)

		b.close()

	})
}

func TestDatastream(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(10, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(10, false)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_datastream", "test_basic")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})

		b.close()

	})

}

func TestCompilerBuiltin(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(1, false)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_compiler_builtins", "test_multi3")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_compiler_builtins", "test_divti3")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		ret := callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_compiler_builtins", "test_divti3_by_0")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, exception.ArithmeticException{}.Code(), "divide by zero")
		assert.Equal(t, ret, true)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_compiler_builtins", "test_udivti3")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		ret = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_compiler_builtins", "test_udivti3_by_0")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, exception.ArithmeticException{}.Code(), "divide by zero")
		assert.Equal(t, ret, true)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_compiler_builtins", "test_modti3")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		ret = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_compiler_builtins", "test_modti3_by_0")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, exception.ArithmeticException{}.Code(), "divide by zero")
		assert.Equal(t, ret, true)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_compiler_builtins", "test_lshlti3")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_compiler_builtins", "test_lshrti3")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_compiler_builtins", "test_ashlti3")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_compiler_builtins", "test_ashrti3")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})

		b.close()

	})

}

type invalidAccessAction struct {
	Code  uint64
	Val   uint64
	Index uint32
	Store bool
}

func TestDB(t *testing.T) {

	name := "testdata_context/test_api_db.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.CreateAccounts([]common.AccountName{common.N("testapi2")}, false, true)
		b.ProduceBlocks(10, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.SetCode(common.AccountName(common.N("testapi2")), code, nil)
		b.ProduceBlocks(1, false)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_db", "primary_i64_general")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_db", "primary_i64_lowerbound")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_db", "primary_i64_upperbound")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_db", "idx64_general")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_db", "idx64_lowerbound")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_db", "idx64_upperbound")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})

		actionInvalidAccess1 := invalidAccessAction{uint64(common.N("testapi")), 10, 0, true}
		payload, _ := rlp.EncodeToBytes(&actionInvalidAccess1)
		act := types.Action{
			Account:       common.AccountName(common.N("testapi")),
			Name:          common.AccountName(wasmTestAction("test_db", "test_invalid_access")),
			Data:          payload,
			Authorization: []types.PermissionLevel{{common.AccountName(common.N("testapi")), common.PermissionName(common.N("active"))}},
		}
		ret := b.PushAction(t, &act, common.AccountName(common.N("testapi")))
		assert.Equal(t, ret, "")

		actionInvalidAccess2 := invalidAccessAction{actionInvalidAccess1.Code, 20, 0, true}
		payload, _ = rlp.EncodeToBytes(&actionInvalidAccess2)
		act = types.Action{
			Account:       common.AccountName(common.N("testapi2")),
			Name:          common.AccountName(wasmTestAction("test_db", "test_invalid_access")),
			Data:          payload,
			Authorization: []types.PermissionLevel{{common.AccountName(common.N("testapi2")), common.PermissionName(common.N("active"))}},
		}
		ret = b.PushAction(t, &act, common.AccountName(common.N("testapi2")))
		assert.Equal(t, inString(ret, "db access violation"), true)

		actionInvalidAccess1.Store = false
		payload, _ = rlp.EncodeToBytes(&actionInvalidAccess1)
		act = types.Action{
			Account:       common.AccountName(common.N("testapi")),
			Name:          common.AccountName(wasmTestAction("test_db", "test_invalid_access")),
			Data:          payload,
			Authorization: []types.PermissionLevel{{common.AccountName(common.N("testapi")), common.PermissionName(common.N("active"))}},
		}
		ret = b.PushAction(t, &act, common.AccountName(common.N("testapi")))
		assert.Equal(t, ret, "")

		actionInvalidAccess1.Store = true
		actionInvalidAccess1.Index = 1
		payload, _ = rlp.EncodeToBytes(&actionInvalidAccess1)
		act = types.Action{
			Account:       common.AccountName(common.N("testapi")),
			Name:          common.AccountName(wasmTestAction("test_db", "test_invalid_access")),
			Data:          payload,
			Authorization: []types.PermissionLevel{{common.AccountName(common.N("testapi")), common.PermissionName(common.N("active"))}},
		}
		ret = b.PushAction(t, &act, common.AccountName(common.N("testapi")))
		assert.Equal(t, ret, "")

		actionInvalidAccess2.Index = 1
		payload, _ = rlp.EncodeToBytes(&actionInvalidAccess2)
		act = types.Action{
			Account:       common.AccountName(common.N("testapi2")),
			Name:          common.AccountName(wasmTestAction("test_db", "test_invalid_access")),
			Data:          payload,
			Authorization: []types.PermissionLevel{{common.AccountName(common.N("testapi2")), common.PermissionName(common.N("active"))}},
		}
		ret = b.PushAction(t, &act, common.AccountName(common.N("testapi2")))
		assert.Equal(t, inString(ret, "db access violation"), true)

		actionInvalidAccess1.Store = true
		payload, _ = rlp.EncodeToBytes(&actionInvalidAccess1)
		act = types.Action{
			Account:       common.AccountName(common.N("testapi")),
			Name:          common.AccountName(wasmTestAction("test_db", "test_invalid_access")),
			Data:          payload,
			Authorization: []types.PermissionLevel{{common.AccountName(common.N("testapi")), common.PermissionName(common.N("active"))}},
		}
		ret = b.PushAction(t, &act, common.AccountName(common.N("testapi")))
		assert.Equal(t, ret, "")

		retException := callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_db", "idx_double_nan_create_fail")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.TransactionException{}.Code(), exception.TransactionException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_db", "idx_double_nan_modify_fail")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.TransactionException{}.Code(), exception.TransactionException{}.What())
		assert.Equal(t, retException, true)

		var loopupType uint32 = 0
		l, _ := rlp.EncodeToBytes(&loopupType)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_db", "idx_double_nan_lookup_fail")}, l, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.TransactionException{}.Code(), exception.TransactionException{}.What())
		assert.Equal(t, retException, true)

		loopupType = 1
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_db", "idx_double_nan_lookup_fail")}, l, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.TransactionException{}.Code(), exception.TransactionException{}.What())
		assert.Equal(t, retException, true)

		loopupType = 2
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_db", "idx_double_nan_lookup_fail")}, l, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.TransactionException{}.Code(), exception.TransactionException{}.What())
		assert.Equal(t, retException, true)

		b.close()

	})
}

func TestMultiIndex(t *testing.T) {

	name := "testdata_context/test_api_multi_index.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(1, false)
		b.CreateAccounts([]common.AccountName{common.N("testapi")}, false, true)
		b.ProduceBlocks(1, false)
		b.SetCode(common.AccountName(common.N("testapi")), code, nil)
		b.ProduceBlocks(1, false)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_general")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_store_only")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_check_without_storing")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx_double_general")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})

		retException := callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_pk_iterator_exceed_end")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_sk_iterator_exceed_end")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_pk_iterator_exceed_begin")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_sk_iterator_exceed_begin")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_pass_pk_ref_to_other_table")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_pass_sk_ref_to_other_table")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_pass_pk_end_itr_to_iterator_to")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_pass_pk_end_itr_to_modify")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_pass_pk_end_itr_to_erase")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_pass_sk_end_itr_to_iterator_to")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_pass_sk_end_itr_to_modify")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_pass_sk_end_itr_to_erase")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_modify_primary_key")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_require_find_fail_with_msg")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_require_find_sk_fail")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)
		retException = callTestFunctionCheckExceptionF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_require_find_sk_fail_with_msg")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))},
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		callTestF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_sk_cache_pk_lookup")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})
		callTestF2(t, b, &testApiAction{wasmTestAction("test_multi_index", "idx64_pk_cache_sk_lookup")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))})

		b.close()

	})
}

func inString(s1, s2 string) bool {
	if strings.Index(s1, s2) <= 0 {
		return false
	}

	return true
}

func (t BaseTester) PushAction(test *testing.T, act *types.Action, authorizer common.AccountName) string {

	trx := NewTransaction()

	if !common.Empty(authorizer) {
		act.Authorization = append(act.Authorization, types.PermissionLevel{authorizer, common.PermissionName(common.N("active"))})
	}

	trx.Transaction.Actions = append(trx.Transaction.Actions, act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)

	if !common.Empty(authorizer) {
		privKey := t.getPrivateKey(authorizer, "active")
		chainId := t.Control.GetChainId()
		trx.Sign(&privKey, &chainId)
	}

	//defer try.HandleReturn()
	returning, ret := false, ""
	try.Try(func() {
		t.PushTransaction(trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
	}).Catch(func(e exception.Exception) {
		returning, ret = true, exception.GetDetailMessage(e)
		return
		//try.Return()
	}).End()

	if returning {
		return ret
	}

	t.ProduceBlocks(1, false)
	return ""
}

func DJBH(str string) uint32 {
	var hash uint32 = 5381
	bytes := []byte(str)

	for i := 0; i < len(bytes); i++ {
		hash = 33*hash ^ uint32(bytes[i])
	}
	return hash
}

func wasmTestAction(cls string, method string) uint64 {
	fmt.Println(cls, ".", method)
	return uint64(DJBH(cls))<<32 | uint64(DJBH(method))
}

func newApplyContext(control *chain.Controller, action *types.Action) *chain.ApplyContext {

	//pack a singedTrx
	trxHeader := types.TransactionHeader{
		Expiration:       common.MaxTimePointSec(),
		RefBlockNum:      4,
		RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 0,
		MaxCpuUsageMS:    0,
		DelaySec:         0,
	}

	trx := types.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*types.Action{},
		Actions:               []*types.Action{action},
		TransactionExtensions: []*types.Extension{},
	}
	signedTrx := types.NewSignedTransaction(&trx, []ecc.Signature{}, []common.HexBytes{})
	privateKey, _ := ecc.NewRandomPrivateKey()
	chainIdType := common.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
	signedTrx.Sign(privateKey, &chainIdType)
	trxContext := chain.NewTransactionContext(control, signedTrx, trx.ID(), common.Now())

	//pack a applycontext from control, trxContext and act
	a := chain.NewApplyContext(control, trxContext, action, 0)
	return a
}

func createNewAccount(control *chain.Controller, name string) {

	//action for create a new account
	// wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	// privKey, _ := ecc.NewPrivateKey(wif)
	// pubKey := privKey.PublicKey()

	creator := "eosio"

	c := chain.NewAccount{
		Creator: common.AccountName(common.N(creator)),
		Name:    common.AccountName(common.N(name)),
		Owner: types.Authority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{Key: *getPublicKey(name, "owner"), Weight: 1}},
		},
		Active: types.Authority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{Key: *getPublicKey(name, "active"), Weight: 1}},
		},
	}

	buffer, _ := rlp.EncodeToBytes(&c)

	act := types.Action{
		Account: common.AccountName(common.N("eosio")),
		Name:    common.ActionName(common.N("newaccount")),
		Data:    buffer,
		Authorization: []types.PermissionLevel{
			//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
			{Actor: common.AccountName(common.N("eosio")), Permission: common.PermissionName(common.N("active"))},
		},
	}

	a := newApplyContext(control, &act)
	//trx := newTransaction(control, &action, []*ecc.PrivateKey{getPrivateKey(creator, "active")})
	//pushTransaction(control, trx)

	//create new account
	chain.ApplyEosioNewaccount(a)
}

func createNewAccount2(control *chain.Controller, name string, creator string) {

	//action for create a new account
	// wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	// privKey, _ := ecc.NewPrivateKey(wif)
	// pubKey := privKey.PublicKey()

	c := chain.NewAccount{
		Creator: common.AccountName(common.N(creator)),
		Name:    common.AccountName(common.N(name)),
		Owner: types.Authority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{Key: *getPublicKey(name, "owner"), Weight: 1}},
		},
		Active: types.Authority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{Key: *getPublicKey(name, "active"), Weight: 1}},
		},
	}

	buffer, _ := rlp.EncodeToBytes(&c)

	action := types.Action{
		Account: common.AccountName(common.N("eosio")),
		Name:    common.ActionName(common.N("newaccount")),
		Data:    buffer,
		Authorization: []types.PermissionLevel{
			//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
			{Actor: common.AccountName(common.N("eosio")), Permission: common.PermissionName(common.N("active"))},
		},
	}

	//a := newApplyContext(control, &act)
	trx := newTransaction(control, &action, []*ecc.PrivateKey{getPrivateKey(creator, "active")})
	pushTransaction(control, trx)

	//create new account
	//chain.ApplyEosioNewaccount(a)
}

func getPrivateKey(account string, permission string) *ecc.PrivateKey {

	var privKey *ecc.PrivateKey
	if account == "eosio" {
		wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
		privKey, _ = ecc.NewPrivateKey(wif)
	} else {
		a := crypto.Hash256(account + "@" + permission)
		g := bytes.NewReader(a.Bytes())
		privKey, _ = ecc.NewDeterministicPrivateKey(g)
	}

	return privKey
}

func getPublicKey(account string, permission string) *ecc.PublicKey {
	PublicKey := getPrivateKey(account, permission).PublicKey()
	return &PublicKey
}

func SetCode(control *chain.Controller, account string, code []byte) {

	setCode := chain.SetCode{
		Account:   common.AccountName(common.N(account)),
		VmType:    0,
		VmVersion: 0,
		Code:      code,
	}
	buffer, _ := rlp.EncodeToBytes(&setCode)
	action := types.Action{
		Account: common.AccountName(common.N("eosio")),
		Name:    common.ActionName(common.N("setcode")),
		Data:    buffer,
		Authorization: []types.PermissionLevel{
			{Actor: common.AccountName(common.N(account)), Permission: common.PermissionName(common.N("active"))},
		},
	}

	// wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	// privateKey, _ := ecc.NewPrivateKey(wif)

	trx := newTransaction(control, &action, []*ecc.PrivateKey{getPrivateKey(account, "active")})
	pushTransaction(control, trx)
}

func pushTransaction(control *chain.Controller, trx *types.TransactionMetadata) *types.TransactionTrace {
	return control.PushTransaction(trx, common.TimePoint(common.MaxMicroseconds()), 0)
}

func newTransaction(control *chain.Controller, action *types.Action, privateKeys []*ecc.PrivateKey) *types.TransactionMetadata {

	trxHeader := types.TransactionHeader{
		Expiration: common.NewTimePointSecTp(control.PendingBlockTime()).AddSec(6),
		// RefBlockNum:      4,
		// RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 0,
		MaxCpuUsageMS:    0,
		DelaySec:         0,
	}

	trx := types.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*types.Action{},
		Actions:               []*types.Action{action},
		TransactionExtensions: []*types.Extension{},
		//RecoveryCache:         make(map[ecc.Signature]types.CachedPubKey),
	}
	headBlockId := control.HeadBlockId()
	trx.SetReferenceBlock(&headBlockId)

	signedTrx := types.NewSignedTransaction(&trx, []ecc.Signature{}, []common.HexBytes{})
	//privateKey, _ := ecc.NewRandomPrivateKey()
	//chainIdType := common.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
	chainIdType := control.GetChainId()
	for _, privateKey := range privateKeys {
		signedTrx.Sign(privateKey, &chainIdType)
	}

	metaTrx := types.NewTransactionMetadataBySignedTrx(signedTrx, common.CompressionNone)

	return metaTrx
}

func RunCheckException(control *chain.Controller, cls string, method string, payload []byte, authorizer string, permissionLevel []types.PermissionLevel, privateKeys []*ecc.PrivateKey,
	errCode exception.ExcTypes, errMsg string) bool {

	//defer try.HandleReturn()
	returning := false
	try.Try(func() {
		pushAction2(control, cls, method, payload, authorizer, permissionLevel, privateKeys)
	}).Catch(func(e exception.Exception) {
		if e.Code() == errCode {
			fmt.Println(errMsg)
			//ret = true
			returning = true
			//try.Return()
		}
	}).End()

	if returning {
		return returning
	}
	//ret = false
	return false
}

func pushAction2(control *chain.Controller, cls string, method string, payload []byte, authorizer string, permissionLevel []types.PermissionLevel, privateKeys []*ecc.PrivateKey) *types.TransactionTrace {

	//wasm := wasmgo.NewWasmGo()
	action := wasmTestAction(cls, method)
	fmt.Println(cls, method, action)

	//fmt.Println(cls, method, action)
	act := types.Action{
		Account: common.AccountName(common.N(authorizer)),
		Name:    common.ActionName(action),
		Data:    payload,
		//Authorization: []types.PermissionLevel{types.PermissionLevel{Actor: common.AccountName(common.N(authorizer)), Permission: common.PermissionName(common.N("active"))}},
		Authorization: permissionLevel,
	}

	//applyContext := newApplyContext(control, &act)
	//codeVersion := crypto.NewSha256Byte([]byte(code))

	trx := newTransaction(control, &act, privateKeys)
	return pushTransaction(control, trx)
}

func callTestFunction2(control *chain.Controller, cls string, method string, payload []byte, authorizer string) *types.TransactionTrace {

	//wasm := wasmgo.NewWasmGo()
	action := wasmTestAction(cls, method)
	fmt.Println(cls, method, action)

	act := types.Action{
		Account:       common.AccountName(common.N(authorizer)),
		Name:          common.ActionName(action),
		Data:          payload,
		Authorization: []types.PermissionLevel{types.PermissionLevel{Actor: common.AccountName(common.N(authorizer)), Permission: common.PermissionName(common.N("active"))}},
	}

	privateKeys := []*ecc.PrivateKey{getPrivateKey(authorizer, "active")}
	trx := newTransaction(control, &act, privateKeys)
	return pushTransaction(control, trx)

}

func callTestFunctionException2(control *chain.Controller, cls string, method string, payload []byte, authorizer string, errCode exception.ExcTypes, errMsg string) bool {

	action := wasmTestAction(cls, method)
	fmt.Println(cls, method, action)

	act := types.Action{
		Account:       common.AccountName(common.N(authorizer)),
		Name:          common.ActionName(action),
		Data:          payload,
		Authorization: []types.PermissionLevel{types.PermissionLevel{Actor: common.AccountName(common.N(authorizer)), Permission: common.PermissionName(common.N("active"))}},
	}

	privateKeys := []*ecc.PrivateKey{getPrivateKey(authorizer, "active")}
	trx := newTransaction(control, &act, privateKeys)

	//pushTransaction(control, trx)

	//defer try.HandleReturn()
	//try.Try(func() {
	//	pushTransaction(control, trx)
	//}).Catch(func(e exception.Exception) {
	//	if e.Code() == errCode {
	//		fmt.Println(errMsg)
	//		ret = true
	//		try.Return()
	//	}
	//}).End()

	ret := pushTransaction(control, trx)
	return ret.Except.Code() == errCode

}

func pushAction(control *chain.Controller, code []byte, cls string, method string, payload []byte, authorizer string) (ret string) {

	wasm := wasmgo.NewWasmGo()
	action := wasmTestAction(cls, method)
	fmt.Println(cls, method, action)

	//fmt.Println(cls, method, action)
	//createNewAccount(control, authorizer)
	act := types.Action{
		Account:       common.AccountName(common.N(authorizer)),
		Name:          common.ActionName(action),
		Data:          payload,
		Authorization: []types.PermissionLevel{types.PermissionLevel{Actor: common.AccountName(common.N(authorizer)), Permission: common.PermissionName(common.N("active"))}},
	}

	applyContext := newApplyContext(control, &act)
	codeVersion := crypto.NewSha256Byte([]byte(code))

	//defer try.HandleReturn()
	returning, ret := false, ""
	try.Try(func() {
		wasm.Apply(codeVersion, code, applyContext)
	}).Catch(func(e exception.Exception) {
		ret = exception.GetDetailMessage(e)
		//try.Return()
		returning = true
	}).End()

	if returning {
		return ret
	}

	return ""
}

func startBlock() *chain.Controller {
	control := chain.GetControllerInstance()
	//blockTimeStamp := types.NewBlockTimeStamp(common.Now())

	blockTimeStamp := types.NewBlockTimeStamp(control.HeadBlockTime() + common.TimePoint(common.Milliseconds(common.DefaultConfig.BlockIntervalMs)))
	control.StartBlock(blockTimeStamp, 0)
	return control
}

func stopBlock(c *chain.Controller) {
	c.Close()
}

func callTestFunction(control *chain.Controller, code []byte, cls string, method string, payload []byte, authorizer string) (ret string) {

	wasm := wasmgo.NewWasmGo()
	action := wasmTestAction(cls, method)

	act := types.Action{
		Account:       common.AccountName(common.N(authorizer)),
		Name:          common.ActionName(action),
		Data:          payload,
		Authorization: []types.PermissionLevel{types.PermissionLevel{Actor: common.AccountName(common.N(authorizer)), Permission: common.PermissionName(common.N("active"))}},
	}

	applyContext := newApplyContext(control, &act)

	fmt.Println(cls, method, action)
	codeVersion := crypto.NewSha256Byte([]byte(code))

	//defer try.HandleReturn()
	//try.Try(func() {
	wasm.Apply(codeVersion, code, applyContext)
	//}).Catch(func(e exception.Exception) {
	//	ret = e.Message()
	//	try.Return()
	//}).End()

	return applyContext.PendingConsoleOutput

}

func callTestFunctionCheckException(control *chain.Controller, code []byte, cls string, method string, payload []byte, authorizer string, errCode exception.ExcTypes, errMsg string) bool {

	wasm := wasmgo.NewWasmGo()
	action := wasmTestAction(cls, method)

	// control := chain.GetControllerInstance()
	// blockTimeStamp := types.NewBlockTimeStamp(common.Now())
	// control.StartBlock(blockTimeStamp, 0)

	act := types.Action{
		Account:       common.AccountName(common.N(authorizer)),
		Name:          common.ActionName(action),
		Data:          payload,
		Authorization: []types.PermissionLevel{types.PermissionLevel{Actor: common.AccountName(common.N(authorizer)), Permission: common.PermissionName(common.N("active"))}},
	}

	applyContext := newApplyContext(control, &act)
	fmt.Println(cls, method, action)
	codeVersion := crypto.NewSha256Byte([]byte(code))

	//ret := false
	//defer try.HandleReturn()
	returning := false
	try.Try(func() {
		wasm.Apply(codeVersion, code, applyContext)
	}).Catch(func(e exception.Exception) {
		if e.Code() == errCode {
			fmt.Println(errMsg)
			//ret = true
			//try.Return()
			returning = true
		}
	}).End()

	if returning {
		return returning
	}

	//ret = false
	return false

}

func sigDigest(chainID, payload []byte) []byte {
	h := sha256.New()
	_, _ = h.Write(chainID)
	_, _ = h.Write(payload)
	return h.Sum(nil)
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

func NewTransaction() *types.SignedTransaction {
	trx := types.Transaction{
		TransactionHeader:     types.TransactionHeader{},
		ContextFreeActions:    []*types.Action{},
		Actions:               []*types.Action{},
		TransactionExtensions: []*types.Extension{},
	}
	signedTrx := types.NewSignedTransaction(&trx, []ecc.Signature{}, []common.HexBytes{})
	return signedTrx
}

func callTestExceptionF2(test *testing.T, t *BaseTester, a actionInterface, data []byte, scope []common.AccountName, billedCpuTimeUs uint32, max_cpu_usage_ms int64, errCode exception.ExcTypes, errMsg string) (ret bool) {

	trx := NewTransaction()

	pl := []types.PermissionLevel{{scope[0], common.PermissionName(common.N("active"))}}
	if len(scope) > 1 {
		for i, account := range scope {
			if i == 0 {
				continue
			}
			pl = append(pl, types.PermissionLevel{account, common.PermissionName(common.N("active"))})
		}
	}

	act := newAction(pl, a)
	act.Data = data
	act.Authorization = pl
	trx.Transaction.Actions = append(trx.Transaction.Actions, act)

	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)

	privKey := t.getPrivateKey(scope[0], "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&privKey, &chainId)

	//defer try.HandleReturn()
	returning := false
	try.Try(func() {
		t.PushTransaction(trx, common.Now()+common.TimePoint(common.Milliseconds(max_cpu_usage_ms)), billedCpuTimeUs)
	}).Catch(func(e exception.Exception) {
		if e.Code() == errCode || inString(e.What(), errMsg) {
			fmt.Println(e.String())
			returning = true
		}

	}).End()

	if returning {
		return returning
	}

	t.ProduceBlocks(1, false)

	//ret = false
	return false
}

func callTestFunctionCheckExceptionF2(test *testing.T, t *BaseTester, a actionInterface, data []byte, scope []common.AccountName, errCode exception.ExcTypes, errMsg string) (ret bool) {

	trx := NewTransaction()

	pl := []types.PermissionLevel{{scope[0], common.PermissionName(common.N("active"))}}
	if len(scope) > 1 {
		for i, account := range scope {
			if i == 0 {
				continue
			}
			pl = append(pl, types.PermissionLevel{account, common.PermissionName(common.N("active"))})
		}
	}

	act := newAction(pl, a)
	act.Data = data
	act.Authorization = pl
	trx.Transaction.Actions = append(trx.Transaction.Actions, act)

	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)

	privKey := t.getPrivateKey(scope[0], "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&privKey, &chainId)

	//defer try.HandleReturn()
	returning := false
	try.Try(func() {
		t.PushTransaction(trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
	}).Catch(func(e exception.Exception) {
		//fmt.Println(exception.GetDetailMessage(e))
		if e.Code() == errCode || inString(exception.GetDetailMessage(e), errMsg) {
			//ret = true
			//try.Return()

			returning = true
		}
	}).End()

	if returning {
		return returning
	}

	t.ProduceBlocks(1, false)
	return false
}

func callTestF2(test *testing.T, t *BaseTester, a actionInterface, data []byte, scope []common.AccountName) *types.TransactionTrace {

	trx := NewTransaction()

	pl := []types.PermissionLevel{{scope[0], common.PermissionName(common.N("active"))}}
	if len(scope) > 1 {
		for i, account := range scope {
			if i == 0 {
				continue
			}
			pl = append(pl, types.PermissionLevel{account, common.PermissionName(common.N("active"))})
		}
	}

	act := newAction(pl, a)
	act.Data = data
	act.Authorization = pl
	trx.Transaction.Actions = append(trx.Transaction.Actions, act)

	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)

	privKey := t.getPrivateKey(scope[0], "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&privKey, &chainId)

	ret := t.PushTransaction(trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
	assert.Equal(test, ret.Receipt.Status, types.TransactionStatusExecuted)

	t.ProduceBlocks(1, false)

	return ret
}

type BaseTester struct {
	ActionResult           string
	DefaultExpirationDelta uint32
	DefaultBilledCpuTimeUs uint32
	AbiSerializerMaxTime   common.Microseconds
	//TempDir                tempDirectory
	Control                 *chain.Controller
	BlockSigningPrivateKeys map[string]ecc.PrivateKey //map[ecc.PublicKey]ecc.PrivateKey
	Cfg                     chain.Config
	ChainTransactions       map[common.BlockIdType]types.TransactionReceipt
	LastProducedBlock       map[common.AccountName]common.BlockIdType
}

func newBaseTester(pushGenesis bool, readMode chain.DBReadMode) *BaseTester {
	t := &BaseTester{}
	t.DefaultExpirationDelta = 6
	t.DefaultBilledCpuTimeUs = 2000
	t.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	t.LastProducedBlock = make(map[common.AccountName]common.BlockIdType)

	t.init(pushGenesis, readMode)
	return t
}

func (t *BaseTester) init(pushGenesis bool, readMode chain.DBReadMode) {
	t.Cfg = *newConfig(readMode)
	t.Control = chain.NewController(&t.Cfg)

	t.open()

	if pushGenesis {
		t.pushGenesisBlock()
	}
}

func newConfig(readMode chain.DBReadMode) *chain.Config {
	cfg := &chain.Config{}
	cfg.BlocksDir = common.DefaultConfig.DefaultBlocksDirName
	cfg.StateDir = common.DefaultConfig.DefaultStateDirName
	cfg.ReversibleDir = common.DefaultConfig.DefaultReversibleBlocksDirName
	cfg.StateSize = 1024 * 1024 * 8
	cfg.StateGuardSize = 0
	cfg.ReversibleCacheSize = 1024 * 1024 * 8
	cfg.ReversibleGuardSize = 0
	//cfg.ContractsConsole = true
	cfg.ReadMode = readMode

	cfg.Genesis = types.NewGenesisState()
	cfg.Genesis.InitialTimestamp, _ = common.FromIsoString("2020-01-01T00:00:00.000")
	cfg.Genesis.InitialKey = BaseTester{}.getPublicKey(common.DefaultConfig.SystemAccountName, "active")

	cfg.ActorWhitelist = *treeset.NewWith(common.TypeName, common.CompareName)
	cfg.ActorBlacklist = *treeset.NewWith(common.TypeName, common.CompareName)
	cfg.ContractWhitelist = *treeset.NewWith(common.TypeName, common.CompareName)
	cfg.ContractBlacklist = *treeset.NewWith(common.TypeName, common.CompareName)
	cfg.ActionBlacklist = *treeset.NewWith(common.TypePair, common.ComparePair)
	cfg.KeyBlacklist = *treeset.NewWith(ecc.TypePubKey, ecc.ComparePubKey)
	cfg.ResourceGreylist = *treeset.NewWith(common.TypeName, common.CompareName)
	cfg.TrustedProducers = *treeset.NewWith(common.TypeName, common.CompareName)

	//cfg.VmType = common.DefaultConfig.DefaultWasmRuntime // TODO

	return cfg
}

func (t *BaseTester) open() {
	//t.Control.Config = t.Cfg
	//t.Control.startUp() //TODO
	//t.Control.StartBlock()
	t.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	//t.Control.AcceptedBlock.Connect() // TODO: Control.signal
}

func (t *BaseTester) close() {
	t.Control.Close()
	t.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
}

func (t BaseTester) PushBlock(b *types.SignedBlock) *types.SignedBlock {
	t.Control.AbortBlock()
	t.Control.PushBlock(b, types.Complete)
	return &types.SignedBlock{}
}

func (t BaseTester) pushGenesisBlock() {
	wasmName := "testdata_context/eosio.bios.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	//if err != nil {
	//	log.Error("pushGenesisBlock is err : %v", err)
	//}
	t.SetCode(common.DefaultConfig.SystemAccountName, code, nil)
	abiName := "testdata_context/eosio.bios.abi"
	abi, _ := ioutil.ReadFile(abiName)
	////if err != nil {
	////	log.Error("pushGenesisBlock is err : %v", err)
	////}
	t.SetAbi(common.DefaultConfig.SystemAccountName, abi, nil)
}

func (t BaseTester) ProduceBlocks(n uint32, empty bool) {
	if empty {
		for i := 0; uint32(i) < n; i++ {
			t.ProduceEmptyBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
		}
	} else {
		for i := 0; uint32(i) < n; i++ {
			t.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
		}
	}
}

func (t BaseTester) produceBlock(skipTime common.Microseconds, skipPendingTrxs bool, skipFlag uint32) *types.SignedBlock {
	headTime := t.Control.HeadBlockTime()
	nextTime := headTime + common.TimePoint(skipTime)
	if common.Empty(t.Control.PendingBlockState()) || t.Control.PendingBlockState().Header.Timestamp != types.NewBlockTimeStamp(nextTime) {
		t.startBlock(nextTime)
	}
	Hbs := t.Control.HeadBlockState()
	producer := Hbs.GetScheduledProducer(types.NewBlockTimeStamp(nextTime))
	//producer := Hbs.GetScheduledProducer(types.BlockTimeStamp(nextTime))
	privKey := ecc.PrivateKey{}
	privateKey, ok := t.BlockSigningPrivateKeys[producer.BlockSigningKey.String()]
	if !ok {
		privKey = t.getPrivateKey(producer.ProducerName, "active")
	} else {
		privKey = privateKey
	}

	if !skipPendingTrxs {
		unappliedTrxs := t.Control.GetUnappliedTransactions()
		for _, trx := range unappliedTrxs {
			trace := t.Control.PushTransaction(trx, common.MaxTimePoint(), 0)
			if !common.Empty(trace.Except) {
				try.EosThrow(trace.Except, "tester produceBlock is error:%#v", trace.Except)
			}
		}

		// scheduledTrxs := t.Control.GetScheduledTransactions()
		// for len(scheduledTrxs) > 0 {
		// 	for _, trx := range scheduledTrxs {
		// 		trace := t.Control.PushScheduledTransaction(&trx, common.MaxTimePoint(), 0)
		// 		if !common.Empty(trace.Except) {
		// 			try.EosThrow(trace.Except, "tester produceBlock is error:%#v", trace.Except)
		// 		}
		// 	}
		// }
	}

	t.Control.FinalizeBlock()
	t.Control.SignBlock(func(d common.DigestType) ecc.Signature {
		sign, _ := privKey.Sign(d.Bytes())
		//if err != nil {
		//	log.Error(err.Error())
		//}
		return sign
	})
	t.Control.CommitBlock(true)
	b := t.Control.HeadBlockState()
	t.LastProducedBlock[t.Control.HeadBlockState().Header.Producer] = b.BlockId
	t.startBlock(nextTime + common.TimePoint(common.TimePoint(common.DefaultConfig.BlockIntervalUs)))
	return t.Control.HeadBlockState().SignedBlock
}

func (t BaseTester) startBlock(blockTime common.TimePoint) {
	headBlockNumber := t.Control.HeadBlockNum()
	producer := t.Control.HeadBlockState().GetScheduledProducer(types.NewBlockTimeStamp(blockTime))
	lastProducedBlockNum := t.Control.LastIrreversibleBlockNum()
	itr := t.LastProducedBlock[producer.ProducerName]
	if !common.Empty(itr) {
		if t.Control.LastIrreversibleBlockNum() > types.NumFromID(&itr) {
			lastProducedBlockNum = t.Control.LastIrreversibleBlockNum()
		} else {
			lastProducedBlockNum = types.NumFromID(&itr)
		}
	}
	t.Control.AbortBlock()
	t.Control.StartBlock(types.NewBlockTimeStamp(blockTime), uint16(headBlockNumber-lastProducedBlockNum))
}

func (t BaseTester) SetTransactionHeaders(trx *types.Transaction, expiration uint32, delaySec uint32) {
	trx.Expiration = common.TimePointSec((common.Microseconds(t.Control.HeadBlockTime()) + common.Seconds(int64(expiration))).ToSeconds())
	headBlockId := t.Control.HeadBlockId()
	trx.SetReferenceBlock(&headBlockId)

	trx.MaxNetUsageWords = 0
	trx.MaxCpuUsageMS = 0
	trx.DelaySec = delaySec
}

func (t BaseTester) CreateAccounts(names []common.AccountName, multiSig bool, includeCode bool) []*types.TransactionTrace {
	traces := make([]*types.TransactionTrace, len(names))
	for i, n := range names {
		traces[i] = t.CreateAccount(n, common.DefaultConfig.SystemAccountName, multiSig, includeCode)
	}
	return traces
}

func (t BaseTester) CreateAccount(name common.AccountName, creator common.AccountName, multiSig bool, includeCode bool) *types.TransactionTrace {
	trx := types.SignedTransaction{}
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0) //TODO: test
	ownerAuth := types.Authority{}
	if multiSig {
		ownerAuth = types.Authority{
			Threshold: 2,
			Keys:      []types.KeyWeight{{Key: t.getPublicKey(name, "owner"), Weight: 1}},
			Accounts:  []types.PermissionLevelWeight{{Permission: types.PermissionLevel{Actor: creator, Permission: common.DefaultConfig.ActiveName}, Weight: 1}},
		}
	} else {
		ownerAuth = types.NewAuthority(t.getPublicKey(name, "owner"), 0)
	}
	activeAuth := types.NewAuthority(t.getPublicKey(name, "active"), 0)

	sortPermissions := func(auth *types.Authority) {

	}
	if includeCode {
		try.EosAssert(ownerAuth.Threshold <= math.MaxUint16, nil, "threshold is too high")
		try.EosAssert(uint64(activeAuth.Threshold) <= uint64(math.MaxUint64), nil, "threshold is too high")
		ownerAuth.Accounts = append(ownerAuth.Accounts, types.PermissionLevelWeight{
			Permission: types.PermissionLevel{Actor: name, Permission: common.DefaultConfig.EosioCodeName},
			Weight:     types.WeightType(ownerAuth.Threshold),
		})
		sortPermissions(&ownerAuth)
		activeAuth.Accounts = append(activeAuth.Accounts, types.PermissionLevelWeight{
			Permission: types.PermissionLevel{Actor: name, Permission: common.DefaultConfig.EosioCodeName},
			Weight:     types.WeightType(activeAuth.Threshold),
		})
		sortPermissions(&activeAuth)
	}
	new := newAccount{
		Creator: creator,
		Name:    name,
		Owner:   ownerAuth,
		Active:  activeAuth,
	}
	data, _ := rlp.EncodeToBytes(new)
	act := &types.Action{
		Account:       new.getAccount(),
		Name:          new.getName(),
		Authorization: []types.PermissionLevel{{creator, common.DefaultConfig.ActiveName}},
		Data:          data,
	}
	trx.Actions = append(trx.Actions, act)

	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	pk := t.getPrivateKey(creator, "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&pk, &chainId)
	return t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) PushTransaction(trx *types.SignedTransaction, deadline common.TimePoint, billedCpuTimeUs uint32) (trace *types.TransactionTrace) {
	_, r := false, (*types.TransactionTrace)(nil)
	try.Try(func() {
		if t.Control.PendingBlockState() == nil {
			t.startBlock(t.Control.HeadBlockTime().AddUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs)))
		}
		c := common.CompressionNone
		size, _ := rlp.EncodeSize(trx)
		if size > 1000 {
			c = common.CompressionZlib
		}
		mtrx := types.NewTransactionMetadataBySignedTrx(trx, c)
		trace = t.Control.PushTransaction(mtrx, deadline, billedCpuTimeUs)
		if trace.ExceptPtr != nil {
			try.EosThrow(trace.ExceptPtr, "tester PushTransaction is error :%#v", trace.ExceptPtr.String())
		}
		if !common.Empty(trace.Except) {
			try.EosThrow(trace.Except, "tester PushTransaction is error :%#v", trace.Except.String())
		}
		r = trace
		return
	}).FcCaptureAndRethrow().End()
	return r
}

// func (t BaseTester) PushAction(act *types.Action, authorizer common.AccountName) {
// 	trx := types.SignedTransaction{}
// 	if !common.Empty(authorizer) {
// 		act.Authorization = []types.PermissionLevel{{authorizer, common.DefaultConfig.ActiveName}}
// 	}
// 	trx.Actions = append(trx.Actions, act)
// 	t.SetTransactionHeaders(&trx.Transaction, 0, 0) //TODO
// 	if common.Empty(authorizer) {
// 		chainId := t.Control.GetChainId()
// 		privateKey := t.getPrivateKey(authorizer, "active")
// 		trx.Sign(&privateKey, &chainId)
// 	}
// 	try.Try(func() {
// 		t.PushTransaction(&trx, 0, 0) //TODO
// 	}).Catch(func(ex exception.Exception) {
// 		//log.Error("tester PushAction is error: %#v", ex.Message())
// 	}).End()
// 	t.ProduceBlock(common.Microseconds(common.DefaultConfig.BlockIntervalMs), 0)
// 	/*BOOST_REQUIRE_EQUAL(true, chain_has_transaction(trx.id()))
// 	success()*/
// 	return
// }

func (t BaseTester) getPrivateKey(keyName common.Name, role string) ecc.PrivateKey {
	pk := &ecc.PrivateKey{}
	if keyName == common.DefaultConfig.SystemAccountName {
		pk, _ = ecc.NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	} else {
		rawPrivKey := crypto.Hash256(keyName.String() + role).Bytes()
		g := bytes.NewReader(rawPrivKey)
		pk, _ = ecc.NewDeterministicPrivateKey(g)
	}
	return *pk
}

func (t BaseTester) getPublicKey(keyName common.Name, role string) ecc.PublicKey {
	priKey := t.getPrivateKey(keyName, role)
	return priKey.PublicKey()
}

func (t BaseTester) ProduceBlock(skipTime common.Microseconds, skipFlag uint32) *types.SignedBlock {
	return t.produceBlock(skipTime, false, skipFlag)
}

func (t BaseTester) ProduceEmptyBlock(skipTime common.Microseconds, skipFlag uint32) *types.SignedBlock {
	t.Control.AbortBlock()
	return t.produceBlock(skipTime, true, skipFlag)
}

func (t BaseTester) PushDummy(from common.AccountName, v *string, billedCpuTimeUs uint32) *types.TransactionTrace {
	//TODO
	trx := types.SignedTransaction{}
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	privKey := t.getPrivateKey(from, "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&privKey, &chainId)
	return t.PushTransaction(&trx, common.MaxTimePoint(), billedCpuTimeUs)
}

func (t BaseTester) SetCode(account common.AccountName, wasm []byte, signer *ecc.PrivateKey) {
	trx := types.SignedTransaction{}
	setCode := chain.SetCode{Account: account, VmType: 0, VmVersion: 0, Code: wasm}
	data, _ := rlp.EncodeToBytes(setCode)
	act := types.Action{
		//Account:       setCode.getAccount(),
		//Name:          setCode.getName(),
		Account:       common.AccountName(common.N("eosio")),
		Name:          common.ActionName(common.N("setcode")),
		Authorization: []types.PermissionLevel{{account, common.DefaultConfig.ActiveName}},
		Data:          data,
	}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	chainId := t.Control.GetChainId()
	if signer != nil {
		trx.Sign(signer, &chainId)
	} else {
		privKey := t.getPrivateKey(account, "active")
		trx.Sign(&privKey, &chainId)
	}
	t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

type setAbi struct {
	Account common.AccountName
	Abi     []byte
}

func (s setAbi) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s setAbi) getName() common.ActionName {
	return common.ActionName(common.N("setabi"))
}
func (t BaseTester) SetAbi(account common.AccountName, abiJson []byte, signer *ecc.PrivateKey) {
	abiEt := abi.AbiDef{}
	json.Unmarshal(abiJson, &abiEt)
	//if err != nil {
	//	log.Error("unmarshal abiJson is wrong :%s", err)
	//}
	trx := types.SignedTransaction{}
	abiBytes, _ := rlp.EncodeToBytes(abiEt)
	setAbi := setAbi{Account: account, Abi: abiBytes}
	data, _ := rlp.EncodeToBytes(setAbi)
	act := types.Action{
		Account:       setAbi.getAccount(),
		Name:          setAbi.getName(),
		Authorization: []types.PermissionLevel{{account, common.DefaultConfig.ActiveName}},
		Data:          data,
	}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	chainId := t.Control.GetChainId()
	if signer != nil {
		trx.Sign(signer, &chainId)
	} else {
		privKey := t.getPrivateKey(account, "active")
		trx.Sign(&privKey, &chainId)
	}
	t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

type VariantsObject map[string]interface{}

func (t BaseTester) SetProducers(accounts []common.AccountName) *types.TransactionTrace {

	schedule := t.GetProducerKeys(&accounts)

	return t.PushAction2(common.DefaultConfig.SystemAccountName,
		common.N("setprods"),
		common.DefaultConfig.SystemAccountName,
		&VariantsObject{"schedule": schedule},
		t.DefaultExpirationDelta,
		0)

}

func (t BaseTester) GetProducerKeys(producerNames *[]common.AccountName) []types.ProducerKey {
	var schedule []types.ProducerKey
	for _, producerName := range *producerNames {
		pk := types.ProducerKey{ProducerName: common.AccountName(producerName), BlockSigningKey: t.getPublicKey(common.AccountName(producerName), "active")}
		schedule = append(schedule, pk)
	}
	return schedule
}

func (t BaseTester) PushAction2(code common.AccountName, acttype common.AccountName,
	actor common.AccountName, data *VariantsObject, expiration uint32, delaySec uint32) *types.TransactionTrace {
	auths := make([]types.PermissionLevel, 0)
	auths = append(auths, types.PermissionLevel{Actor: actor, Permission: common.DefaultConfig.ActiveName})
	return t.PushAction4(code, acttype, &auths, data, expiration, delaySec)
}

func (t BaseTester) PushAction4(code common.AccountName, acttype common.AccountName,
	auths *[]types.PermissionLevel, data *VariantsObject, expiration uint32, delaySec uint32) *types.TransactionTrace {
	trx := types.SignedTransaction{}
	try.Try(func() {
		action := t.GetAction(code, acttype, *auths, data)
		trx.Actions = append(trx.Actions, action)
	})
	t.SetTransactionHeaders(&trx.Transaction, expiration, delaySec)
	chainId := t.Control.GetChainId()
	key := ecc.PrivateKey{}
	for _, auth := range *auths {
		key = t.getPrivateKey(auth.Actor, auth.Permission.String())
		trx.Sign(&key, &chainId)
	}
	return t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) GetAction(code common.AccountName, actType common.AccountName,
	auths []types.PermissionLevel, data *VariantsObject) *types.Action {
	acnt := t.Control.GetAccount(code)
	a := acnt.GetAbi()
	action := types.Action{code, actType, auths, nil}
	//actionTypeName := a.ActionForName(actType).Type
	buf, _ := json.Marshal(data)
	//if err != nil {
	//	log.Error("tester GetAction Marshal is error:%s", err)
	//}
	//action.Data, _ = a.EncodeAction(common.N(actionTypeName), buf) //TODO
	action.Data, _ = a.EncodeAction(actType, buf)
	//if err != nil {
	//	log.Error("tester GetAction EncodeAction is error:%s", err)
	//}
	//log.Error("action:%s", action)
	return &action
}

//func (t BaseTester) SetAbi(account common.AccountName, abiJson []byte, signer *ecc.PrivateKey) {
//	abiEt := abi.AbiDef{}
//	err := json.Unmarshal(abiJson, &abiEt)
//	//if err != nil {
//	//	log.Error("unmarshal abiJson is wrong :%s", err)
//	//}
//	trx := types.SignedTransaction{}
//	abiBytes, _ := rlp.EncodeToBytes(abiEt)
//	setAbi := setAbi{Account: account, Abi: abiBytes}
//	data, _ := rlp.EncodeToBytes(setAbi)
//	act := types.Action{
//		// Account:       account,
//		// Name:          setAbi.getName(),
//		Account: common.AccountName(common.N("eosio")),
//		Name:    common.ActionName(common.N("setabi")),
//		Authorization: []types.PermissionLevel{{account, common.DefaultConfig.ActiveName}},
//		Data:          data,
//	}
//	trx.Actions = append(trx.Actions, &act)
//	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
//	chainId := t.Control.GetChainId()
//	if signer != nil {
//		trx.Sign(signer, &chainId)
//	} else {
//		privKey := t.getPrivateKey(account, "active")
//		trx.Sign(&privKey, &chainId)
//	}
//	t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
//}
