// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exec_test

import (
	"fmt"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/eosspark/eos-go/cvm/exec"
	"github.com/eosspark/eos-go/rlp"
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
			t.Fatalf("error excute hello.wasm")
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
		if strings.Compare(applyContext.PendingConsoleOutput, "hello,mic,-3,20,3.14E+38,3.14E+300,walker,0x313233343536") != 0 {
			t.Fatalf("error excute hello.wasm")
		}

	})

}
