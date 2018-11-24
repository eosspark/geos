// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasmgo_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
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

func (d *dummy_action) get_name() uint64 {
	return common.N("dummy_action")
}

func (d *dummy_action) get_account() uint64 {
	return common.N("testapi")
}

func TestContextAction(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		control := startBlock()
		createNewAccount(control, "testapi")

		dummy13 := dummy_action{DUMMY_ACTION_DEFAULT_A, DUMMY_ACTION_DEFAULT_B, DUMMY_ACTION_DEFAULT_C}

		callTestFunction(control, code, "test_action", "assert_true", []byte{}, "testapi")
		callTestFunction(control, code, "test_action", "assert_false", []byte{}, "testapi")

		b, _ := rlp.EncodeToBytes(&dummy13)
		callTestFunction(control, code, "test_action", "read_action_normal", b, "testapi")

		//rawBytes := []byte{(1 << 16)}
		b = bytes.Repeat([]byte{byte(0x01)}, 1<<16)
		callTestFunction(control, code, "test_action", "read_action_to_0", b, "testapi")
		b = bytes.Repeat([]byte{byte(0x01)}, 1<<16+1)
		callTestFunction(control, code, "test_action", "read_action_to_0", b, "testapi")

		b = bytes.Repeat([]byte{byte(0x01)}, 1)
		callTestFunction(control, code, "test_action", "read_action_to_64k", b, "testapi")
		b = bytes.Repeat([]byte{byte(0x01)}, 3)
		callTestFunction(control, code, "test_action", "read_action_to_64k", b, "testapi")

		callTestFunction(control, code, "test_action", "require_auth", []byte{}, "testapi")

		a3only := []types.PermissionLevel{{common.AccountName(common.N("acc3")), common.PermissionName(common.N("active"))}}
		b, _ = rlp.EncodeToBytes(a3only)
		callTestFunction(control, code, "test_action", "require_auth", b, "testapi")

		a4only := []types.PermissionLevel{{common.AccountName(common.N("acc4")), common.PermissionName(common.N("active"))}}
		b, _ = rlp.EncodeToBytes(a4only)
		callTestFunction(control, code, "test_action", "require_auth", b, "testapi")

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

		trace := callTestFunction(control, code, "test_print", "test_prints", []byte{}, "testapi")
		result := trace
		assert.Equal(t, result, "abcefg")

		trace = callTestFunction(control, code, "test_print", "test_prints_l", []byte{}, "testapi")
		result = trace
		assert.Equal(t, result, "abatest")

		trace = callTestFunction(control, code, "test_print", "test_printi", []byte{}, "testapi")
		result = trace
		assert.Equal(t, result[0:1], string(strconv.FormatInt(0, 10)))
		assert.Equal(t, result[1:7], string(strconv.FormatInt(556644, 10)))
		assert.Equal(t, result[7:9], string(strconv.FormatInt(-1, 10)))

		trace = callTestFunction(control, code, "test_print", "test_printui", []byte{}, "testapi")
		result = trace
		assert.Equal(t, result[0:1], string(strconv.FormatInt(0, 10)))
		assert.Equal(t, result[1:7], string(strconv.FormatInt(556644, 10)))

		v := -1
		assert.Equal(t, result[7:len(result)], string(strconv.FormatUint(uint64(v), 10))) //-1 / 1844674407370955161

		trace = callTestFunction(control, code, "test_print", "test_printn", []byte{}, "testapi")
		result = trace
		assert.Equal(t, result[0:5], "abcde")
		assert.Equal(t, result[5:10], "ab.de")
		assert.Equal(t, result[10:16], "1q1q1q")
		assert.Equal(t, result[16:27], "abcdefghijk")
		assert.Equal(t, result[27:39], "abcdefghijkl")
		assert.Equal(t, result[39:52], "abcdefghijkl1")
		assert.Equal(t, result[52:65], "abcdefghijkl1")
		assert.Equal(t, result[65:78], "abcdefghijkl1")

		trace = callTestFunction(control, code, "test_print", "test_printi128", []byte{}, "testapi")
		result = trace

		s := strings.Split(result, "\n")
		assert.Equal(t, s[0], "1")
		assert.Equal(t, s[1], "0")
		assert.Equal(t, s[2], "-170141183460469231731687303715884105728")
		assert.Equal(t, s[3], "-87654323456")

		trace = callTestFunction(control, code, "test_print", "test_printui128", []byte{}, "testapi")
		result = trace
		s = strings.Split(result, "\n")
		assert.Equal(t, s[0], "340282366920938463463374607431768211455")
		assert.Equal(t, s[1], "0")
		assert.Equal(t, s[2], "87654323456")

		trace = callTestFunction(control, code, "test_print", "test_printsf", []byte{}, "testapi")
		result = trace
		r := strings.Split(result, "\n")
		assert.Equal(t, r[0], "5.000000e-01")
		assert.Equal(t, r[1], "-3.750000e+00")
		assert.Equal(t, r[2], "6.666667e-07")

		trace = callTestFunction(control, code, "test_print", "test_printdf", []byte{}, "testapi")
		result = trace
		r = strings.Split(result, "\n")
		assert.Equal(t, r[0], "5.000000000000000e-01")
		assert.Equal(t, r[1], "-3.750000000000000e+00")
		assert.Equal(t, r[2], "6.666666666666666e-07")

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

func TestContextChecktime(t *testing.T) {

	name := "testdata_context/test_api.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		control := startBlock()
		createNewAccount(control, "testapi")

		callTestFunction(control, code, "test_checktime", "checktime_pass", []byte{}, "testapi")
		//callTestFunction(control, code, "test_checktime", "checktime_failure", []byte{}, "testapi")
		callTestFunction(control, code, "test_checktime", "checktime_sha1_failure", []byte{}, "testapi")
		callTestFunction(control, code, "test_checktime", "checktime_assert_sha1_failure", []byte{}, "testapi")
		callTestFunction(control, code, "test_checktime", "checktime_sha256_failure", []byte{}, "testapi")
		callTestFunction(control, code, "test_checktime", "checktime_assert_sha256_failure", []byte{}, "testapi")
		callTestFunction(control, code, "test_checktime", "checktime_sha512_failure", []byte{}, "testapi")
		callTestFunction(control, code, "test_checktime", "checktime_assert_sha512_failure", []byte{}, "testapi")
		callTestFunction(control, code, "test_checktime", "checktime_ripemd160_failure", []byte{}, "testapi")
		callTestFunction(control, code, "test_checktime", "checktime_assert_ripemd160_failure", []byte{}, "testapi")

		stopBlock(control)

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

		action1 := invalidAccessAction{common.N("testapi"), 10, 0, true}
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
	wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	privKey, _ := ecc.NewPrivateKey(wif)
	pubKey := privKey.PublicKey()

	creator := chain.NewAccount{
		Creator: common.AccountName(common.N("eosio")),
		Name:    common.AccountName(common.N(name)),
		Owner: types.Authority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{Key: pubKey, Weight: 1}},
		},
		Active: types.Authority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{Key: pubKey, Weight: 1}},
		},
	}

	buffer, _ := rlp.EncodeToBytes(&creator)

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

	//create new account
	chain.ApplyEosioNewaccount(a)
}

func pushAction(control *chain.Controller, code []byte, cls string, method string, payload []byte, authorizer string) (ret string) {

	wasm := wasmgo.NewWasmGo()
	action := wasmTestAction(cls, method)

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
		ret = e.Message()
		try.Return()
	}).End()

	return ""
}

func startBlock() *chain.Controller {
	control := chain.GetControllerInstance()
	blockTimeStamp := types.NewBlockTimeStamp(common.Now())
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

	//fmt.Println(cls, method, action)
	codeVersion := crypto.NewSha256Byte([]byte(code))

	defer try.HandleReturn()
	try.Try(func() {
		wasm.Apply(codeVersion, code, applyContext)
	}).Catch(func(e exception.Exception) {
		ret = e.Message()
		try.Return()
	}).End()

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
