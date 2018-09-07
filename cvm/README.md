cvm
=====

fork from https://github.com/go-interpreter/wagon

**NOTE:** `cvm` requires `Go >= 1.9.x`.

## examples

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/eosgo/common"
	"github.com/eosgo/control"
	"github.com/eosgo/cvm/exec"
)

func main() {

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		flag.PrintDefaults()
		os.Exit(1)
	}

	name := flag.Arg(0)

	code, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}

	apply_context := &control.Apply_context{Receiver: common.AccountName(exec.N("walker")), Code: common.AccountName(exec.N("walker")), Action: common.ActionName(exec.N("hi"))}
	
	wasm := exec.NewWasm_interface()
	wasm.Apply("00000000", code, apply_context)

}

go run hello.wasm // from eos hello contract
