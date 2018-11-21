package main

import (
	. "github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app/include"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	_ "github.com/eosspark/eos-go/plugins/appbase/plugin/chain_plugin"
	_ "github.com/eosspark/eos-go/plugins/appbase/plugin/http_plugin"
	_ "github.com/eosspark/eos-go/plugins/appbase/plugin/net_plugin"
	_ "github.com/eosspark/eos-go/plugins/appbase/plugin/producer_plugin"
	"os"
	"os/signal"
	"strings"
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

var basicPlugin = []string{"ProducerPlugin", "ChainPlugin", "NetPlugin", "HttpPlugin"}

//var pro producer_plugin.Producer_plugin
//var net net_plugin.Net_plugin

func main() {

	try.Try(func() {
		App.SetVersion(Version)
		App.SetDefaultDataDir()
		App.SetDefaultConfigDir()
		if !App.Initialize(basicPlugin) {
			os.Exit(INITIALIZEFAIL)
		}
		App.StartUp()
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		select {
		case <-sigChan:
			App.ShutDown()
		}
	}).Catch(func(e ExtractGenesisStateException) {
		os.Exit(EXTRACTED_GENESIS)
	}).Catch(func(e FixedReversibleDbException) {
		os.Exit(FIXED_REVERSIBLE)
	}).Catch(func(e NodeManagementSuccess) {
		os.Exit(NODE_MANAGEMENT_SUCCESS)
	}).Catch(func(e Exception) {
		if e.Code() == StdExceptionCode {
			if strings.Contains(e.Message(), "database dirty flag set") {
				log.Error("database dirty flag set (likely due to unclean shutdown): replay required")
				os.Exit(DATABASE_DIRTY)
			} else if strings.Contains(e.Message(), "database metadata dirty flag set") {
				log.Error("database metadata dirty flag set (likely due to unclean shutdown): replay required")
				os.Exit(DATABASE_DIRTY)
			}
		}
		log.Error(e.Message())
		os.Exit(OTHER_FAIL)
	}).Catch(func(e try.RuntimeError) {
		if strings.Contains(e.Message,"database dirty flag set") {
			log.Error("database dirty flag set (likely due to unclean shutdown): replay required")
			os.Exit(DATABASE_DIRTY)
		}else if strings.Contains(e.Message,"database metadata dirty flag set") {
			log.Error("database metadata dirty flag set (likely due to unclean shutdown): replay required")
			os.Exit(DATABASE_DIRTY)
		}else {
			log.Error("%s",e.Message)
		}
		os.Exit(OTHER_FAIL)
	}).Catch(func(e Exception) {
		log.Error("%s", e.Message())
	}).Catch(func(... interface{}) {
		log.Error("unknown exception")
		os.Exit(OTHER_FAIL)
	}).End()

	os.Exit(SUCCESS)
}
