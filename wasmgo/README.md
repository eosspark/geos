wasmgo
=====

fork from https://github.com/go-interpreter/wagon

**NOTE:** `cvm` requires `Go >= 1.9.x`.

## examples

```
package main

import (
	"fmt"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/wasmgo"
	"io/ioutil"
	"log"
)

func main() {

	name := "hello.wasm"
	code, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}

	wasmgo := wasmgo.NewWasmGo()
	param, _ := rlp.EncodeToBytes(common.N("walker")) //[]byte{0x00, 0x00, 0x00, 0x00, 0x5c, 0x05, 0xa3, 0xe1}
	applyContext := &chain.ApplyContext{
		Receiver: common.AccountName(common.N("hello")),
		Act: &types.Action{
			Account: common.AccountName(common.N("hello")),
			Name:    common.ActionName(common.N("hi")),
			Data:    param,
		},
	}

	codeVersion := crypto.NewSha256Byte([]byte(code))
	wasmgo.Apply(codeVersion, code, applyContext)

	//print "hello, walker"
	fmt.Println(applyContext.PendingConsoleOutput)

}
```

## go run hello.wasm
// hello.wasm from eosio hello contract
