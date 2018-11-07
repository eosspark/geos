package main

import (
	_ "github.com/eosspark/eos-go/plugins/appbase/plugin/net_plugin"
	_ "github.com/eosspark/eos-go/plugins/appbase/plugin/producer_plugin"
	_ "github.com/eosspark/eos-go/plugins/appbase/plugin/http_plugin"
	_ "github.com/eosspark/eos-go/plugins/appbase/plugin/chain_plugin"
	. "github.com/eosspark/eos-go/plugins/appbase/app/include"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/exception/try"
	"os"
	"os/signal"
	. "github.com/eosspark/eos-go/exception"
)

const (
	OTHER_FAIL              = -2
	INITIALIZEFAIL          = -1
	SUCCESS                 = 0
	BAD_ALLOC               = 1
	DATABASE_DIRTY          = 2
	FIXED_REVERSIBLE        = 3
	EXTRACTED_GENESIS       = 4
	NODE_MANAGEMENT_SUCCESS = 5
)


var basicPlugin  = []string{"ProducerPlugin","ChainPlugin","NetPlugin","HttpPlugin"}


//var pro producer_plugin.Producer_plugin
//var net net_plugin.Net_plugin

func main() {
	defer try.HandleReturn()
	try.Try(func() {
			App.SetVersion(Version)
			App.SetDefaultDataDir()
			App.SetDefaultConfigDir()
			if !App.Initialize(basicPlugin) {
				try.Return()
			}

			App.My.Options.Run(os.Args)
			App.StartUp()
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)
			select {
				case <-sigChan:
					App.ShutDown()
			}
		}).Catch(func(e ExtractGenesisStateException) {

		}).Catch(func(e FixedReversibleDbException) {

		}).End()
}
