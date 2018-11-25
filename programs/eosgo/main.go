package main

import (
	. "github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	. "github.com/eosspark/eos-go/plugins/appbase/app/include"
	"os"
	"strings"

	//plugins
	//_ "github.com/eosspark/eos-go/plugins/chain_plugin"
	_ "github.com/eosspark/eos-go/plugins/console_plugin"
	_ "github.com/eosspark/eos-go/plugins/producer_plugin"
)

const (
	OTHER_FAIL              = -2
	INITIALIZE_FAIL         = -1
	SUCCESS                 = 0
	BAD_ALLOC               = 1
	DATABASE_DIRTY          = 2
	FIXED_REVERSIBLE        = 3
	EXTRACTED_GENESIS       = 4
	NODE_MANAGEMENT_SUCCESS = 5
)

var basicPlugin = []PluginName{ProducerPlug, ConsolePlug}

func main() {

	try.Try(func() {
		App().SetVersion(Version)
		App().SetDefaultDataDir()
		App().SetDefaultConfigDir()
		if !App().Initialize(basicPlugin) {
			os.Exit(INITIALIZE_FAIL)
		}
		App().StartUp()
		App().Exec()

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

		//}).Catch(func(e try.RuntimeError) {
		//	if strings.Contains(e.Message, "database dirty flag set") {
		//		log.Error("database dirty flag set (likely due to unclean shutdown): replay required")
		//		os.Exit(DATABASE_DIRTY)
		//
		//	} else if strings.Contains(e.Message, "database metadata dirty flag set") {
		//		log.Error("database metadata dirty flag set (likely due to unclean shutdown): replay required")
		//		os.Exit(DATABASE_DIRTY)
		//
		//	} else {
		//		log.Error("%s", e.Message)
		//	}
		//	os.Exit(OTHER_FAIL)

	}).Catch(func(e error) {
		log.Error("%s", e.Error())

	}).Catch(func(interface{}) {
		log.Error("unknown exception")
		os.Exit(OTHER_FAIL)

	}).End()

	os.Exit(SUCCESS)
}
