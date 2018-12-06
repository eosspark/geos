// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasmgo_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/eosspark/container/sets/treeset"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/math"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/wasmgo"
	"github.com/stretchr/testify/assert"

	arithmetic "github.com/eosspark/eos-go/common/arithmetic_types"
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

		control := startBlock()
		createNewAccount2(control, "testapi", "eosio")
		createNewAccount2(control, "acc1", "eosio")
		createNewAccount2(control, "acc2", "eosio")
		createNewAccount2(control, "acc3", "eosio")
		createNewAccount2(control, "acc4", "eosio")

		SetCode(control, "testapi", code)

		permissions := []types.PermissionLevel{types.PermissionLevel{common.AccountName(common.N("testapi")), common.PermissionName(common.N("active"))}}
		privateKeys := []*ecc.PrivateKey{getPrivateKey("inita", "active")}
		RunCheckException(control, "test_action", "require_notice", []byte{}, "testapi", permissions, privateKeys, exception.UnsatisfiedAuthorization{}.Code(), exception.UnsatisfiedAuthorization{}.What())

		permissions = []types.PermissionLevel{
			types.PermissionLevel{common.AccountName(common.N("testapi")), common.PermissionName(common.N("active"))},
			types.PermissionLevel{common.AccountName(common.N("acc3")), common.PermissionName(common.N("active"))},
			types.PermissionLevel{common.AccountName(common.N("acc4")), common.PermissionName(common.N("active"))},
		}
		privateKeys = []*ecc.PrivateKey{
			getPrivateKey("testapi", "active"),
			getPrivateKey("acc3", "active"),
			getPrivateKey("acc4", "active"),
		}
		ret := pushAction2(control, "test_action", "require_auth", []byte{}, "testapi", permissions, privateKeys)
		assert.Equal(t, ret.Receipt.Status, types.TransactionStatusExecuted)

		// now := control.HeadBlockTime().AddUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
		// n := now.TimeSinceEpoch().Count()
		// //fmt.Println(now)
		// b, _ := rlp.EncodeToBytes(&n)
		// callTestFunction2(control, "test_action", "test_current_time", b, "testapi")

		account := common.AccountName(common.N("testapi"))
		b, _ := rlp.EncodeToBytes(&account)
		callTestFunction2(control, "test_action", "test_current_receiver", b, "testapi")
		callTestFunction2(control, "test_transaction", "send_action_sender", b, "testapi")

	})

}

func TestRequireRecipient(t *testing.T) {
	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		control := startBlock()
		createNewAccount2(control, "testapi", "eosio")
		createNewAccount2(control, "testapi2", "eosio")
		createNewAccount2(control, "acc5", "eosio")

		SetCode(control, "testapi", code)
		SetCode(control, "acc5", code)

		//permissions := []types.PermissionLevel{
		//	types.PermissionLevel{common.AccountName(common.N("testapi")), common.PermissionName(common.N("active"))},
		//}
		//privateKeys := []*ecc.PrivateKey{
		//	getPrivateKey("testapi", "active"),
		//}
		//ret := pushAction2(control, "test_action", "require_notice_tests", []byte{}, "testapi", permissions, privateKeys)
		//assert.Equal(t, ret.Receipt.Status, types.TransactionStatusExecuted)

		SetCode(control, "testapi2", code)
		data := arithmetic.Int128{uint64(common.N("testapi")), uint64(common.N("testapi2"))}
		b, _ := rlp.EncodeToBytes(&data)
		ret := callTestFunctionException2(control, "test_action", "test_ram_billing_in_notify", b, "testapi", exception.SubjectiveBlockProductionException{}.Code(), "Cannot charge RAM to other accounts during notify.")
		assert.Equal(t, true, ret)

		data = arithmetic.Int128{0, uint64(common.N("testapi2"))}
		b, _ = rlp.EncodeToBytes(&data)
		callTestFunction2(control, "test_action", "test_ram_billing_in_notify", b, "testapi")

		data = arithmetic.Int128{uint64(common.N("testapi2")), uint64(common.N("testapi2"))}
		b, _ = rlp.EncodeToBytes(&data)
		callTestFunction2(control, "test_action", "test_ram_billing_in_notify", b, "testapi")

		stopBlock(control)

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

func TestContextAction(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		control := startBlock()
		createNewAccount(control, "testapi")
		createNewAccount(control, "acc1")
		createNewAccount(control, "acc2")
		createNewAccount(control, "acc3")
		createNewAccount(control, "acc4")

		SetCode(control, "testapi", code)

		dummy13 := dummy_action{DUMMY_ACTION_DEFAULT_A, DUMMY_ACTION_DEFAULT_B, DUMMY_ACTION_DEFAULT_C}

		callTestFunction(control, code, "test_action", "assert_true", []byte{}, "testapi")
		callTestFunctionCheckException(control, code, "test_action", "assert_false", []byte{}, "testapi", exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())

		b, _ := rlp.EncodeToBytes(&dummy13)
		callTestFunction(control, code, "test_action", "read_action_normal", b, "testapi")

		//rawBytes := []byte{(1 << 16)}
		b = bytes.Repeat([]byte{byte(0x01)}, 1<<16)
		callTestFunction(control, code, "test_action", "read_action_to_0", b, "testapi")

		b = bytes.Repeat([]byte{byte(0x01)}, 1<<16+1)
		callTestFunctionCheckException(control, code, "test_action", "read_action_to_0", b, "testapi", exception.OverlappingMemoryError{}.Code(), exception.OverlappingMemoryError{}.What())

		b = bytes.Repeat([]byte{byte(0x01)}, 1)
		callTestFunction(control, code, "test_action", "read_action_to_64k", b, "testapi")

		b = bytes.Repeat([]byte{byte(0x01)}, 3)
		callTestFunctionCheckException(control, code, "test_action", "read_action_to_64k", b, "testapi", exception.OverlappingMemoryError{}.Code(), exception.OverlappingMemoryError{}.What())

		ret := pushAction(control, code, "test_action", "require_notice", b, "testapi")
		assert.Equal(t, ret, "assertion failure with message: Should've failed")

		callTestFunctionCheckException(control, code, "test_action", "require_auth", []byte{}, "testapi", exception.MissingAuthException{}.Code(), exception.MissingAuthException{}.What())
		a3only := []types.PermissionLevel{{common.AccountName(common.N("acc3")), common.PermissionName(common.N("active"))}}
		b, _ = rlp.EncodeToBytes(a3only)
		callTestFunctionCheckException(control, code, "test_action", "require_auth", b, "testapi", exception.MissingAuthException{}.Code(), exception.MissingAuthException{}.What())

		a4only := []types.PermissionLevel{{common.AccountName(common.N("acc4")), common.PermissionName(common.N("active"))}}
		b, _ = rlp.EncodeToBytes(a4only)
		callTestFunctionCheckException(control, code, "test_action", "require_auth", b, "testapi", exception.MissingAuthException{}.Code(), exception.MissingAuthException{}.What())

		stopBlock(control)

	})

}

func TestContextPrint(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		control := startBlock()
		createNewAccount(control, "testapi")

		result := callTestFunction(control, code, "test_print", "test_prints", []byte{}, "testapi")
		assert.Equal(t, result, "abcefg")

		result = callTestFunction(control, code, "test_print", "test_prints_l", []byte{}, "testapi")
		assert.Equal(t, result, "abatest")

		result = callTestFunction(control, code, "test_print", "test_printi", []byte{}, "testapi")
		assert.Equal(t, result[0:1], string(strconv.FormatInt(0, 10)))
		assert.Equal(t, result[1:7], string(strconv.FormatInt(556644, 10)))
		assert.Equal(t, result[7:9], string(strconv.FormatInt(-1, 10)))

		result = callTestFunction(control, code, "test_print", "test_printui", []byte{}, "testapi")
		assert.Equal(t, result[0:1], string(strconv.FormatInt(0, 10)))
		assert.Equal(t, result[1:7], string(strconv.FormatInt(556644, 10)))

		v := -1
		assert.Equal(t, result[7:len(result)], string(strconv.FormatUint(uint64(v), 10))) //-1 / 1844674407370955161

		result = callTestFunction(control, code, "test_print", "test_printn", []byte{}, "testapi")
		assert.Equal(t, result[0:5], "abcde")
		assert.Equal(t, result[5:10], "ab.de")
		assert.Equal(t, result[10:16], "1q1q1q")
		assert.Equal(t, result[16:27], "abcdefghijk")
		assert.Equal(t, result[27:39], "abcdefghijkl")
		assert.Equal(t, result[39:52], "abcdefghijkl1")
		assert.Equal(t, result[52:65], "abcdefghijkl1")
		assert.Equal(t, result[65:78], "abcdefghijkl1")

		result = callTestFunction(control, code, "test_print", "test_printi128", []byte{}, "testapi")
		s := strings.Split(result, "\n")
		assert.Equal(t, s[0], "1")
		assert.Equal(t, s[1], "0")
		assert.Equal(t, s[2], "-170141183460469231731687303715884105728")
		assert.Equal(t, s[3], "-87654323456")

		result = callTestFunction(control, code, "test_print", "test_printui128", []byte{}, "testapi")
		s = strings.Split(result, "\n")
		assert.Equal(t, s[0], "340282366920938463463374607431768211455")
		assert.Equal(t, s[1], "0")
		assert.Equal(t, s[2], "87654323456")

		result = callTestFunction(control, code, "test_print", "test_printsf", []byte{}, "testapi")
		r := strings.Split(result, "\n")
		assert.Equal(t, r[0], "5.000000e-01")
		assert.Equal(t, r[1], "-3.750000e+00")
		assert.Equal(t, r[2], "6.666667e-07")

		result = callTestFunction(control, code, "test_print", "test_printdf", []byte{}, "testapi")
		r = strings.Split(result, "\n")
		assert.Equal(t, r[0], "5.000000000000000e-01")
		assert.Equal(t, r[1], "-3.750000000000000e+00")
		assert.Equal(t, r[2], "6.666666666666666e-07")

		//result = callTestFunction(control, code, "test_print", "test_printqf", []byte{}, "testapi")
		//r = strings.Split(result, "\n")
		//assert.Equal(t, r[0], "5.000000000000000000e-01")
		//assert.Equal(t, r[1], "-3.750000000000000000e+00")
		//assert.Equal(t, r[2], "6.666666666666666667e-07")

		stopBlock(control)

	})

}

func TestContextTypes(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		control := startBlock()
		createNewAccount(control, "testapi")

		callTestFunction(control, code, "test_types", "types_size", []byte{}, "testapi")
		callTestFunction(control, code, "test_types", "char_to_symbol", []byte{}, "testapi")
		callTestFunction(control, code, "test_types", "string_to_name", []byte{}, "testapi")
		callTestFunction(control, code, "test_types", "name_class", []byte{}, "testapi")

		stopBlock(control)

	})

}

func TestContextMemory(t *testing.T) {

	name := "testdata_context/test_api_mem.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		control := startBlock()
		createNewAccount(control, "testapi")

		callTestFunction(control, code, "test_memory", "test_memory_allocs", []byte{}, "testapi")
		callTestFunction(control, code, "test_memory", "test_memory_hunk", []byte{}, "testapi")
		callTestFunction(control, code, "test_memory", "test_memory_hunks", []byte{}, "testapi")
		//callTestFunction(control, code, "test_memory", "test_memory_hunks_disjoint", []byte{}, "testapi")
		callTestFunction(control, code, "test_memory", "test_memset_memcpy", []byte{}, "testapi")

		callTestFunctionCheckException(control, code, "test_memory", "test_memcpy_overlap_start", []byte{}, "testapi", exception.OverlappingMemoryError{}.Code(), exception.OverlappingMemoryError{}.What())
		callTestFunctionCheckException(control, code, "test_memory", "test_memcpy_overlap_end", []byte{}, "testapi", exception.OverlappingMemoryError{}.Code(), exception.OverlappingMemoryError{}.What())

		callTestFunction(control, code, "test_memory", "test_memcmp", []byte{}, "testapi")

		//callTestFunction(control, code, "test_memory", "test_outofbound_0", []byte{}, "testapi")
		// callTestFunction(control, code, "test_memory", "test_outofbound_1", []byte{}, "testapi")
		// callTestFunction(control, code, "test_memory", "test_outofbound_2", []byte{}, "testapi")
		// callTestFunction(control, code, "test_memory", "test_outofbound_3", []byte{}, "testapi")
		// callTestFunction(control, code, "test_memory", "test_outofbound_4", []byte{}, "testapi")
		// callTestFunction(control, code, "test_memory", "test_outofbound_5", []byte{}, "testapi")
		// callTestFunction(control, code, "test_memory", "test_outofbound_6", []byte{}, "testapi")
		// callTestFunction(control, code, "test_memory", "test_outofbound_7", []byte{}, "testapi")
		// callTestFunction(control, code, "test_memory", "test_outofbound_8", []byte{}, "testapi")
		// callTestFunction(control, code, "test_memory", "test_outofbound_9", []byte{}, "testapi")
		// callTestFunction(control, code, "test_memory", "test_outofbound_10", []byte{}, "testapi")
		// callTestFunction(control, code, "test_memory", "test_outofbound_11", []byte{}, "testapi")
		// callTestFunction(control, code, "test_memory", "test_outofbound_12", []byte{}, "testapi")
		// callTestFunction(control, code, "test_memory", "test_outofbound_13", []byte{}, "testapi")

		callTestFunction(control, code, "test_extended_memory", "test_initial_buffer", []byte{}, "testapi")
		callTestFunction(control, code, "test_extended_memory", "test_page_memory", []byte{}, "testapi")
		callTestFunction(control, code, "test_extended_memory", "test_page_memory_exceeded", []byte{}, "testapi")
		callTestFunction(control, code, "test_extended_memory", "test_page_memory_negative_bytes", []byte{}, "testapi")

		stopBlock(control)
	})

}

func TestContextAuth(t *testing.T) {

	name := "testdata_context/auth.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(name)
		wasm := wasmgo.NewWasmGo()
		param, _ := rlp.EncodeToBytes(common.N("walker"))
		applyContext := &chain.ApplyContext{
			Receiver: common.AccountName(common.N("ctx.auth")),
			Act: &types.Action{
				Account: common.AccountName(common.N("ctx.auth")),
				Name:    common.ActionName(common.N("test")),
				Data:    param,
				Authorization: []types.PermissionLevel{{
					Actor:      common.AccountName(common.N("walker")),
					Permission: common.PermissionName(common.N("active")),
				}},
			},
			UsedAuthorizations: make([]bool, 1),
		}

		codeVersion := crypto.NewSha256Byte([]byte(code))
		wasm.Apply(codeVersion, code, applyContext)

		result := fmt.Sprintf("%v", applyContext.PendingConsoleOutput)
		assert.Equal(t, result, "walker has authorization,walker is account")

	})

}

func TestContextCrypto(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(name)

		wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
		privKey, err := ecc.NewPrivateKey(wif)

		chainID, err := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
		payload, err := hex.DecodeString("88e4b25a00006c08ac5b595b000000000000")
		digest := sigDigest(chainID, payload)
		sig, err := privKey.Sign(digest)
		pubKey, err := sig.PublicKey(digest)

		load := digest

		p, _ := rlp.EncodeToBytes(pubKey)
		load = append(load, p...)

		s, _ := rlp.EncodeToBytes(sig)
		load = append(load, s...)

		fmt.Println("load:", hex.EncodeToString(load))

		control := startBlock()
		createNewAccount(control, "testapi")

		callTestFunction(control, code, "test_crypto", "test_recover_key", load, "testapi")
		callTestFunction(control, code, "test_crypto", "test_recover_key_assert_true", load, "testapi")
		callTestFunction(control, code, "test_crypto", "test_sha1", []byte{}, "testapi")
		callTestFunction(control, code, "test_crypto", "test_sha256", []byte{}, "testapi")
		callTestFunction(control, code, "test_crypto", "test_sha512", []byte{}, "testapi")
		callTestFunction(control, code, "test_crypto", "test_ripemd160", []byte{}, "testapi")
		callTestFunction(control, code, "test_crypto", "sha1_no_data", []byte{}, "testapi")
		callTestFunction(control, code, "test_crypto", "sha256_no_data", []byte{}, "testapi")
		callTestFunction(control, code, "test_crypto", "sha512_no_data", []byte{}, "testapi")
		callTestFunction(control, code, "test_crypto", "ripemd160_no_data", []byte{}, "testapi")
		callTestFunction(control, code, "test_crypto", "assert_sha256_true", []byte{}, "testapi")
		callTestFunction(control, code, "test_crypto", "assert_sha1_true", []byte{}, "testapi")
		callTestFunction(control, code, "test_crypto", "assert_sha512_true", []byte{}, "testapi")
		callTestFunction(control, code, "test_crypto", "assert_ripemd160_true", []byte{}, "testapi")

		callTestFunctionCheckException(control, code, "test_crypto", "assert_sha256_false", []byte{}, "testapi", exception.CryptoApiException{}.Code(), exception.CryptoApiException{}.What())
		callTestFunctionCheckException(control, code, "test_crypto", "assert_sha1_false", []byte{}, "testapi", exception.CryptoApiException{}.Code(), exception.CryptoApiException{}.What())
		callTestFunctionCheckException(control, code, "test_crypto", "assert_sha512_false", []byte{}, "testapi", exception.CryptoApiException{}.Code(), exception.CryptoApiException{}.What())
		callTestFunctionCheckException(control, code, "test_crypto", "assert_ripemd160_false", []byte{}, "testapi", exception.CryptoApiException{}.Code(), exception.CryptoApiException{}.What())

		stopBlock(control)

	})
}

func TestContextFixedPoint(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		control := startBlock()
		createNewAccount(control, "testapi")

		callTestFunction(control, code, "test_fixedpoint", "create_instances", []byte{}, "testapi")
		callTestFunction(control, code, "test_fixedpoint", "test_addition", []byte{}, "testapi")
		callTestFunction(control, code, "test_fixedpoint", "test_subtraction", []byte{}, "testapi")
		callTestFunction(control, code, "test_fixedpoint", "test_multiplication", []byte{}, "testapi")
		callTestFunction(control, code, "test_fixedpoint", "test_division", []byte{}, "testapi")
		callTestFunctionCheckException(control, code, "test_fixedpoint", "test_division_by_0", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())

		stopBlock(control)

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

		ret := callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		assert.Equal(t, ret, true)

		ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 0, 200, exception.TxCpuUsageExceeded{}.Code(), exception.TxCpuUsageExceeded{}.What())
		assert.Equal(t, ret, true)

		ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 0, 200, exception.BlockCpuUsageExceeded{}.Code(), exception.BlockCpuUsageExceeded{}.What())
		assert.Equal(t, ret, true)

		ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_sha1_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		assert.Equal(t, ret, true)

		ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_assert_sha1_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		assert.Equal(t, ret, true)

		ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_sha256_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		assert.Equal(t, ret, true)

		ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_assert_sha256_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		assert.Equal(t, ret, true)

		ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_assert_sha512_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		assert.Equal(t, ret, true)

		ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_ripemd160_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		assert.Equal(t, ret, true)

		ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_sha1_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		assert.Equal(t, ret, true)

		ret = callTestExceptionF2(t, b, &testApiAction{wasmTestAction("test_checktime", "checktime_assert_ripemd160_failure")}, []byte{}, []common.AccountName{common.AccountName(common.N("testapi"))}, 5000, 10, exception.DeadlineException{}.Code(), exception.DeadlineException{}.What())
		assert.Equal(t, ret, true)

		b.close()

	})
}

func TestContextDatastream(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		control := startBlock()
		createNewAccount(control, "testapi")

		callTestFunction(control, code, "test_datastream", "test_basic", []byte{}, "testapi")
		stopBlock(control)

	})
}

func TestContextCompilerBuiltin(t *testing.T) {

	name := "testdata_context/compiler_builtin.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		control := startBlock()
		createNewAccount(control, "testapi")

		callTestFunction(control, code, "test_compiler_builtins", "test_ashrti3", []byte{}, "testapi")
		callTestFunction(control, code, "test_compiler_builtins", "test_ashlti3", []byte{}, "testapi")
		callTestFunction(control, code, "test_compiler_builtins", "test_lshrti3", []byte{}, "testapi")
		callTestFunction(control, code, "test_compiler_builtins", "test_lshlti3", []byte{}, "testapi")

		callTestFunction(control, code, "test_compiler_builtins", "test_umodti3", []byte{}, "testapi")
		callTestFunctionCheckException(control, code, "test_compiler_builtins", "test_umodti3_by_0", []byte{}, "testapi",
			exception.ArithmeticException{}.Code(), exception.ArithmeticException{}.What())

		callTestFunction(control, code, "test_compiler_builtins", "test_modti3", []byte{}, "testapi")
		callTestFunctionCheckException(control, code, "test_compiler_builtins", "test_modti3_by_0", []byte{}, "testapi",
			exception.ArithmeticException{}.Code(), exception.ArithmeticException{}.What())

		callTestFunction(control, code, "test_compiler_builtins", "test_udivti3", []byte{}, "testapi")
		callTestFunctionCheckException(control, code, "test_compiler_builtins", "test_udivti3_by_0", []byte{}, "testapi",
			exception.ArithmeticException{}.Code(), exception.ArithmeticException{}.What())

		callTestFunction(control, code, "test_compiler_builtins", "test_divti3", []byte{}, "testapi")
		callTestFunctionCheckException(control, code, "test_compiler_builtins", "test_divti3_by_0", []byte{}, "testapi",
			exception.ArithmeticException{}.Code(), exception.ArithmeticException{}.What())

		callTestFunction(control, code, "test_compiler_builtins", "test_multi3", []byte{}, "testapi")

		stopBlock(control)
	})
}

type invalidAccessAction struct {
	Code  uint64
	Val   uint64
	Index uint32
	Store bool
}

func TestContextDB(t *testing.T) {

	name := "testdata_context/test_api_db.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		control := startBlock()
		createNewAccount(control, "testapi")
		createNewAccount(control, "testapi2")

		callTestFunction(control, code, "test_db", "primary_i64_general", []byte{}, "testapi")
		callTestFunction(control, code, "test_db", "primary_i64_lowerbound", []byte{}, "testapi")
		callTestFunction(control, code, "test_db", "primary_i64_upperbound", []byte{}, "testapi")
		callTestFunction(control, code, "test_db", "idx64_general", []byte{}, "testapi")
		callTestFunction(control, code, "test_db", "idx64_lowerbound", []byte{}, "testapi")
		callTestFunction(control, code, "test_db", "idx64_upperbound", []byte{}, "testapi")

		action1 := invalidAccessAction{uint64(common.N("testapi")), 10, 0, true}
		actionData1, _ := rlp.EncodeToBytes(&action1)
		ret := pushAction(control, code, "test_db", "test_invalid_access", actionData1, "testapi")
		assert.Equal(t, ret, "")

		action2 := invalidAccessAction{action1.Code, 20, 0, true}
		actionData2, _ := rlp.EncodeToBytes(&action2)
		ret = pushAction(control, code, "test_db", "test_invalid_access", actionData2, "testapi2")
		assert.Equal(t, ret, "db access violation")

		action1.Store = false
		actionData3, _ := rlp.EncodeToBytes(&action1)
		ret = pushAction(control, code, "test_db", "test_invalid_access", actionData3, "testapi")
		assert.Equal(t, ret, "")

		action1.Store = true
		action1.Index = 1
		actionData1, _ = rlp.EncodeToBytes(&action1)
		ret = pushAction(control, code, "test_db", "test_invalid_access", actionData1, "testapi")
		assert.Equal(t, ret, "")

		action2.Index = 1
		actionData2, _ = rlp.EncodeToBytes(&action2)
		ret = pushAction(control, code, "test_db", "test_invalid_access", actionData2, "testapi2")
		assert.Equal(t, ret, "db access violation")

		action1.Store = false
		actionData3, _ = rlp.EncodeToBytes(&action1)
		ret = pushAction(control, code, "test_db", "test_invalid_access", actionData3, "testapi")
		assert.Equal(t, ret, "")

		retException := callTestFunctionCheckException(control, code, "test_db", "idx_double_nan_create_fail", []byte{}, "testapi",
			exception.TableAccessViolation{}.Code(), exception.TableAccessViolation{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_db", "idx_double_nan_modify_fail", []byte{}, "testapi",
			exception.TableAccessViolation{}.Code(), exception.TableAccessViolation{}.What())
		assert.Equal(t, retException, true)

		var loopupType uint32 = 0
		l, _ := rlp.EncodeToBytes(&loopupType)
		retException = callTestFunctionCheckException(control, code, "test_db", "idx_double_nan_lookup_fail", l, "testapi",
			exception.TableAccessViolation{}.Code(), exception.TableAccessViolation{}.What())
		assert.Equal(t, retException, true)

		loopupType = 1
		l, _ = rlp.EncodeToBytes(&loopupType)
		callTestFunctionCheckException(control, code, "test_db", "idx_double_nan_lookup_fail", l, "testapi",
			exception.TableAccessViolation{}.Code(), exception.TableAccessViolation{}.What())
		assert.Equal(t, retException, true)

		loopupType = 2
		l, _ = rlp.EncodeToBytes(&loopupType)
		retException = callTestFunctionCheckException(control, code, "test_db", "idx_double_nan_lookup_fail", l, "testapi",
			exception.TableAccessViolation{}.Code(), exception.TableAccessViolation{}.What())
		assert.Equal(t, retException, true)

		//fmt.Println(ret)

		stopBlock(control)

	})
}

func TestContextMultiIndex(t *testing.T) {

	name := "testdata_context/test_api_multi_index.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		control := startBlock()
		createNewAccount(control, "testapi")
		createNewAccount(control, "testapi2")

		callTestFunction(control, code, "test_multi_index", "idx64_general", []byte{}, "testapi")
		callTestFunction(control, code, "test_multi_index", "idx64_store_only", []byte{}, "testapi")
		callTestFunction(control, code, "test_multi_index", "idx64_check_without_storing", []byte{}, "testapi")

		retException := callTestFunctionCheckException(control, code, "test_multi_index", "idx64_pk_iterator_exceed_end", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_sk_iterator_exceed_end", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_pk_iterator_exceed_begin", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_sk_iterator_exceed_begin", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_pass_pk_ref_to_other_table", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_pass_sk_ref_to_other_table", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_pass_pk_end_itr_to_iterator_to", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_pass_pk_end_itr_to_modify", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_pass_pk_end_itr_to_erase", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_pass_sk_end_itr_to_iterator_to", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_pass_sk_end_itr_to_modify", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_pass_sk_end_itr_to_erase", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_modify_primary_key", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		//assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_run_out_of_avl_pk", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_require_find_fail", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_require_find_fail_with_msg", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_require_find_sk_fail", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		retException = callTestFunctionCheckException(control, code, "test_multi_index", "idx64_require_find_sk_fail_with_msg", []byte{}, "testapi",
			exception.EosioAssertMessageException{}.Code(), exception.EosioAssertMessageException{}.What())
		assert.Equal(t, retException, true)

		callTestFunction(control, code, "test_multi_index", "idx64_sk_cache_pk_lookup", []byte{}, "testapi")
		callTestFunction(control, code, "test_multi_index", "idx64_pk_cache_sk_lookup", []byte{}, "testapi")

		stopBlock(control)

	})
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
	errCode exception.ExcTypes, errMsg string) (ret bool) {

	defer try.HandleReturn()
	try.Try(func() {
		pushAction2(control, cls, method, payload, authorizer, permissionLevel, privateKeys)
	}).Catch(func(e exception.Exception) {
		if e.Code() == errCode {
			fmt.Println(errMsg)
			ret = true
			try.Return()
		}
	}).End()

	ret = false
	return
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

	defer try.HandleReturn()
	try.Try(func() {
		wasm.Apply(codeVersion, code, applyContext)
	}).Catch(func(e exception.Exception) {
		ret = exception.GetDetailMessage(e)
		try.Return()
	}).End()

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

func callTestFunctionCheckException(control *chain.Controller, code []byte, cls string, method string, payload []byte, authorizer string, errCode exception.ExcTypes, errMsg string) (ret bool) {

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
	defer try.HandleReturn()
	try.Try(func() {
		wasm.Apply(codeVersion, code, applyContext)
	}).Catch(func(e exception.Exception) {
		if e.Code() == errCode {
			fmt.Println(errMsg)
			ret = true
			try.Return()
		}
	}).End()

	ret = false
	return

}

func sigDigest(chainID, payload []byte) []byte {
	h := sha256.New()
	_, _ = h.Write(chainID)
	_, _ = h.Write(payload)
	return h.Sum(nil)
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

	// action := wasmTestAction(cls, method)
	// fmt.Println(cls, method, action)
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

	//trx.get_signature_keys(test.control->get_chain_id() );
	//ret := t.PushTransaction(trx, common.Now()+common.TimePoint(common.Milliseconds(max_cpu_usage_ms)), billedCpuTimeUs)
	//assert.Equal(test, ret.Except.Code() == errCode, true)

	//ret := false
	defer try.HandleReturn()
	try.Try(func() {
		t.PushTransaction(trx, common.Now()+common.TimePoint(common.Milliseconds(max_cpu_usage_ms)), billedCpuTimeUs)
	}).Catch(func(e exception.Exception) {
		if e.Code() == errCode {
			fmt.Println(e.Message())
			ret = true
			try.Return()
		}
	}).End()

	t.ProduceBlocks(1, false)

	ret = false
	return
}

func callTestF2(test *testing.T, t *BaseTester, a actionInterface, data []byte, scope []common.AccountName) *types.TransactionTrace {

	trx := NewTransaction()

	// action := wasmTestAction(cls, method)
	// fmt.Println(cls, method, action)
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

	//trx.get_signature_keys(test.control->get_chain_id() );

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
	t.Control = chain.NewController(t.Cfg)

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
	cfg.ContractsConsole = false
	cfg.ReadMode = readMode

	cfg.Genesis = types.NewGenesisState()
	cfg.Genesis.InitialTimestamp, _ = common.FromIsoString("2020-01-01T00:00:00.000")
	cfg.Genesis.InitialKey = BaseTester{}.getPublicKey(common.DefaultConfig.SystemAccountName, "active")

	cfg.ActorWhitelist = *treeset.NewWith(common.CompareName)
	cfg.ActorBlacklist = *treeset.NewWith(common.CompareName)
	cfg.ContractWhitelist = *treeset.NewWith(common.CompareName)
	cfg.ContractBlacklist = *treeset.NewWith(common.CompareName)
	cfg.ActionBlacklist = *treeset.NewWith(common.ComparePair)
	cfg.KeyBlacklist = *treeset.NewWith(ecc.ComparePubKey)
	cfg.ResourceGreylist = *treeset.NewWith(common.CompareName)
	cfg.TrustedProducers = *treeset.NewWith(common.CompareName)

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
	//abiName := "testdata_context/eosio.bios.abi"
	//abi, _ := ioutil.ReadFile(abiName)
	////if err != nil {
	////	log.Error("pushGenesisBlock is err : %v", err)
	////}
	//t.SetAbi(common.DefaultConfig.SystemAccountName, abi, nil)
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
	producer := Hbs.GetScheduledProducer(types.BlockTimeStamp(nextTime))
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
			try.EosThrow(trace.ExceptPtr, "tester PushTransaction is error :%#v", trace.ExceptPtr.Message())
		}
		if !common.Empty(trace.Except) {
			try.EosThrow(trace.Except, "tester PushTransaction is error :%#v", trace.Except.Message())
		}
		r = trace
		return
	}).FcCaptureAndRethrow().End()
	return r
}

func (t BaseTester) PushAction(act *types.Action, authorizer common.AccountName) {
	trx := types.SignedTransaction{}
	if !common.Empty(authorizer) {
		act.Authorization = []types.PermissionLevel{{authorizer, common.DefaultConfig.ActiveName}}
	}
	trx.Actions = append(trx.Actions, act)
	t.SetTransactionHeaders(&trx.Transaction, 0, 0) //TODO
	if common.Empty(authorizer) {
		chainId := t.Control.GetChainId()
		privateKey := t.getPrivateKey(authorizer, "active")
		trx.Sign(&privateKey, &chainId)
	}
	try.Try(func() {
		t.PushTransaction(&trx, 0, 0) //TODO
	}).Catch(func(ex exception.Exception) {
		//log.Error("tester PushAction is error: %#v", ex.Message())
	}).End()
	t.ProduceBlock(common.Microseconds(common.DefaultConfig.BlockIntervalMs), 0)
	/*BOOST_REQUIRE_EQUAL(true, chain_has_transaction(trx.id()))
	success()*/
	return
}

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

// func (t BaseTester) SetAbi(account common.AccountName, abiJson []byte, signer *ecc.PrivateKey) {
// 	abiEt := abi.AbiDef{}
// 	err := json.Unmarshal(abiJson, &abiEt)
// 	if err != nil {
// 		log.Error("unmarshal abiJson is wrong :%s", err)
// 	}
// 	trx := types.SignedTransaction{}
// 	abiBytes, _ := rlp.EncodeToBytes(abiEt)
// 	setAbi := setAbi{Account: account, Abi: abiBytes}
// 	data, _ := rlp.EncodeToBytes(setAbi)
// 	act := types.Action{
// 		// Account:       account,
// 		// Name:          setAbi.getName(),
// 		Account: common.AccountName(common.N("eosio")),
// 		Name:    common.ActionName(common.N("setabi")),
// 		Authorization: []types.PermissionLevel{{account, common.DefaultConfig.ActiveName}},
// 		Data:          data,
// 	}
// 	trx.Actions = append(trx.Actions, &act)
// 	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
// 	chainId := t.Control.GetChainId()
// 	if signer != nil {
// 		trx.Sign(signer, &chainId)
// 	} else {
// 		privKey := t.getPrivateKey(account, "active")
// 		trx.Sign(&privKey, &chainId)
// 	}
// 	t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
// }
