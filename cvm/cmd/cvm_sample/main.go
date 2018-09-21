package main

import (
	"flag"
	"github.com/eosspark/eos-go/chain/types"
	"io/ioutil"
	"log"
	"os"

	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
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
	applyContext := &chain.ApplyContext{
		Receiver: common.AccountName(exec.N("walker")),
		Act: types.Action{
			Account: common.AccountName(exec.N("walker")),
			Name:    common.ActionName(exec.N("hi")),
			Data:    []byte{0x00, 0x00, 0x00, 0x00, 0x5c, 0x05, 0xa3, 0xe1}, //'{"walker"}'
		},
	}

	// applyContext := &exec.ApplyContext{
	// 	Receiver: exec.AccountName(exec.N("walker")),
	// 	Contract: exec.AccountName(exec.N("walker")),
	// 	Action:   exec.ActionName(exec.N("hi")),
	// }

	wasm.Apply("00000000", code, applyContext)

}
