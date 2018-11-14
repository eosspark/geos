package main

import (
	"fmt"
	"os"
	"gopkg.in/urfave/cli.v1"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"github.com/eosspark/eos-go/plugins/producer_plugin"
	"log"
	"syscall"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/plugins/producer_plugin/testing"
)


func main() {

	/*
	go run main.go -e -p eosio -p yuanc --private-key '["EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM","5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss"]' --private-key '["EOS5jeUuKEZ8s8LLoxz4rNysYdHWboup8KtkyJzZYQzcVKFGek9Zu","5Ja3h2wJNUnNcoj39jDMHGigsazvbGHAeLYEHM5uTwtfUoRDoYP"]'
	 */
	fmt.Println(os.Args)

	app := cli.NewApp()
	iosv := asio.NewIoContext()

	chainTester :=testing.NewChainTester(0, common.AccountName(common.N("eosio")),common.AccountName(common.N("yuanc")))
	testing.Control = chainTester.Control

	producerPlugin := producer_plugin.NewProducerPlugin(iosv)

	producerPlugin.SetProgramOptions(&app.Flags)
	app.Action = func(c *cli.Context) {
		producerPlugin.PluginInitialize(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	producerPlugin.PluginStartup()

	sigint := asio.NewSignalSet(iosv, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIPE)
	sigint.AsyncWait(func(err error) {
		iosv.Stop()
		sigint.Cancel()
	})

	iosv.Run()

	producerPlugin.PluginShutdown()
}
