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

	"github.com/eosspark/eos-go/cvm/exec"
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

	wasm := exec.NewWasmInterface()

	applyContext := &exec.ApplyContext{
		Receiver: exec.AccountName(exec.N("walker")),
		Contract: exec.AccountName(exec.N("walker")),
		Action:   exec.ActionName(exec.N("hi")),
	}

	wasm.Apply("00000000", code, applyContext)

}

// hello.wasm is from eos hello contract
## go run hello.wasm
