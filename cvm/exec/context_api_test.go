// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec_test

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eosspark/eos-go/cvm/exec"
	"github.com/eosspark/eos-go/rlp"
	"github.com/stretchr/testify/assert"
)

func TestContextApis(t *testing.T) {
	fnames, err := filepath.Glob(filepath.Join("testdata_context", "*.wasm"))
	if err != nil {
		t.Fatal(err)
	}
	for _, fname := range fnames {
		name := fname
		t.Run(filepath.Base(name), func(t *testing.T) {
			code, err := ioutil.ReadFile(name)
			if err != nil {
				t.Fatal(err)
			}

			_, fileName := filepath.Split(name)
			if strings.Compare(fileName, "hello.wasm") == 0 {
				fmt.Println(fileName)
				wasm := exec.NewWasmInterface()
				applyContext := &chain.ApplyContext{
					Receiver: common.AccountName(exec.N("hello")),
					Act: types.Action{
						Account: common.AccountName(exec.N("hello")),
						Name:    common.ActionName(exec.N("hi")),
						Data:    []byte{0x00, 0x00, 0x00, 0x00, 0x5c, 0x05, 0xa3, 0xe1}, //'{"walker"}'
					},
				}

				codeVersion := rlp.NewSha256Byte([]byte(code)).String()
				wasm.Apply(codeVersion, code, applyContext)

				//print "hello,walker"
				//fmt.Println(applyContext.PendingConsoleOutput)
				if strings.Compare(applyContext.PendingConsoleOutput, "Hello, walker") != 0 {
					t.Fatalf("error excute hello.wasm")
				}
			}

		})
	}
}

func TestContextAtion(t *testing.T) {

	name := "testdata_context/action.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(name)
		wasm := exec.NewWasmInterface()
		applyContext := &chain.ApplyContext{
			Receiver: common.AccountName(exec.N("ctx.action")),
			Act: types.Action{
				Account: common.AccountName(exec.N("ctx.action")),
				Name:    common.ActionName(exec.N("test")),
				Data:    []byte{0x00, 0x00, 0x00, 0x00, 0x5c, 0x05, 0xa3, 0xe1}, //'{"walker"}'
			},
		}

		codeVersion := rlp.NewSha256Byte([]byte(code)).String()
		wasm.Apply(codeVersion, code, applyContext)

		//print "hello,walker"
		fmt.Println(applyContext.PendingConsoleOutput)
		if strings.Compare(applyContext.PendingConsoleOutput, "receiver:ctx.action,code:ctx.action,action:test,hello, walker") != 0 {
			t.Fatalf("error excute action.wasm")
		}

	})

}

func TestContextConsole(t *testing.T) {

	name := "testdata_context/console.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(name)
		wasm := exec.NewWasmInterface()
		applyContext := &chain.ApplyContext{
			Receiver: common.AccountName(exec.N("ctx.console")),
			Act: types.Action{
				Account: common.AccountName(exec.N("ctx.console")),
				Name:    common.ActionName(exec.N("test")),
				Data:    []byte{0x00, 0x00, 0x00, 0x00, 0x5c, 0x05, 0xa3, 0xe1}, //'{"walker"}'
			},
		}

		codeVersion := rlp.NewSha256Byte([]byte(code)).String()
		wasm.Apply(codeVersion, code, applyContext)

		//print "hello,walker"
		fmt.Println(applyContext.PendingConsoleOutput)
		if strings.Compare(applyContext.PendingConsoleOutput, "hello,mic,-3,20,3.14E+38,3.14E+300,walker,0x313233343536") != 0 {
			t.Fatalf("error excute console.wasm")
		}

	})

}

func TestContextMemory(t *testing.T) {

	name := "testdata_context/memory.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(name)
		wasm := exec.NewWasmInterface()
		applyContext := &chain.ApplyContext{
			Receiver: common.AccountName(exec.N("ctx.memory")),
			Act: types.Action{
				Account: common.AccountName(exec.N("ctx.memory")),
				Name:    common.ActionName(exec.N("test")),
				Data:    []byte{0x00, 0x00, 0x00, 0x00, 0x5c, 0x05, 0xa3, 0xe1}, //'{"walker"}'
			},
		}

		codeVersion := rlp.NewSha256Byte([]byte(code)).String()
		wasm.Apply(codeVersion, code, applyContext)

		//print "hello,walker"
		fmt.Println(applyContext.PendingConsoleOutput)
		if strings.Compare(applyContext.PendingConsoleOutput, "cccccccccccccchecksum256 ok") != 0 {
			t.Fatalf("error excute memory.wasm")
		}

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
		wasm := exec.NewWasmInterface()
		param, _ := rlp.EncodeToBytes(exec.N("walker"))
		applyContext := &chain.ApplyContext{
			Receiver: common.AccountName(exec.N("ctx.auth")),
			Act: types.Action{
				Account: common.AccountName(exec.N("ctx.auth")),
				Name:    common.ActionName(exec.N("test")),
				Data:    param,
				Authorization: []types.PermissionLevel{{
					Actor:      common.AccountName(exec.N("walker")),
					Permission: common.PermissionName(exec.N("active")),
				}},
			},
			UsedAuthorizations: make([]bool, 1),
		}

		codeVersion := rlp.NewSha256Byte([]byte(code)).String()
		wasm.Apply(codeVersion, code, applyContext)

		//fmt.Println(applyContext.PendingConsoleOutput)
		//if strings.Compare(applyContext.PendingConsoleOutput, "walker has authorization,walker is account") != 0 {
		//	t.Fatalf("error excute memory.wasm")
		//}

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

		//fromEOSIOC := "SIG_K1_K2WBNtiTY8o4mqFSz7HPnjkiT9JhUYGFa81RrzaXr3aWRF1F8qwVfutJXroqiL35ZiHTcvn8gPWGYJDwnKZTCcbAGJy4i9"
		//assert.Equal(t, fromEOSIOC, sig.String())

		//load, _ := rlp.EncodeToBytes(digest)

		load := digest
		//fmt.Println(string(digest))

		p, _ := rlp.EncodeToBytes(pubKey)
		load = append(load, p...)
		//fmt.Println(string(p))

		s, _ := rlp.EncodeToBytes(sig)
		load = append(load, s...)
		//fmt.Println(string(s))

		fmt.Println("d:", hex.EncodeToString(digest), " s:", hex.EncodeToString(s), " p:", hex.EncodeToString(p))
		fmt.Println("load:", hex.EncodeToString(load))
		callTestFunction(code, "test_crypto", "test_recover_key_assert_true", load)

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

func callTestFunction(code []byte, cls string, method string, payload []byte) {

	wasm := exec.NewWasmInterface()
	action := wasmTestAction(cls, method)
	applyContext := &chain.ApplyContext{
		Receiver: common.AccountName(exec.N("test.api")),
		Act: types.Action{
			Account: common.AccountName(exec.N("test.api")),
			Name:    common.ActionName(action),
			Data:    payload,
		},
	}

	fmt.Println(action)
	codeVersion := rlp.NewSha256Byte([]byte(code)).String()
	wasm.Apply(codeVersion, code, applyContext)

}

func callTestFunctionCheckException(cls string, method string, payload []byte) {

}

func sigDigest(chainID, payload []byte) []byte {
	h := sha256.New()
	_, _ = h.Write(chainID)
	_, _ = h.Write(payload)
	return h.Sum(nil)
}
