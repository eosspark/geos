package main

import (
	"fmt"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/cvm/exec"
	"io/ioutil"
	"log"
)

func main() {

	name := "hello.wasm"
	code, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}

	wasm := exec.NewWasmInterface()
	param, _ := rlp.EncodeToBytes(exec.N("walker"))
	applyContext := &chain.ApplyContext{
		Receiver: common.AccountName(exec.N("hello")),
		Act: types.Action{
			Account: common.AccountName(exec.N("hello")),
			Name:    common.ActionName(exec.N("hi")),
			//Data:    []byte{0x00, 0x00, 0x00, 0x00, 0x5c, 0x05, 0xa3, 0xe1}, //'{"walker"}'
			Data: param,
		},
	}

	codeVersion := rlp.NewSha256Byte([]byte(code)).String()
	wasm.Apply(codeVersion, code, applyContext)

	//print "hello, walker"
	fmt.Println(applyContext.PendingConsoleOutput)

}
