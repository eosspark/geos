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
	"github.com/eosspark/eos-go/plugins/console_plugin"
	"github.com/eosspark/eos-go/plugins/producer_plugin"
	"github.com/eosspark/eos-go/plugins/template_plugin"
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

//go run main.go -e -p eosio --private-key [\"EOS859gxfnXyUriMgUeThh1fWv3oqcpLFyHa3TfFYC4PK2HqhToVM\",\"5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss\"] --console
func main() {

	try.Try(func() {

		App().SetVersion(Version)
		App().SetDefaultDataDir()
		App().SetDefaultConfigDir()
		if !App().Initialize([]PluginTypeName{
			producer_plugin.ProducerPlug,
			console_plugin.ConsolePlug,
			template_plugin.TemplatePlug,
		}) {
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

	}).Catch(func(e try.RuntimeError) {
		if strings.Contains(e.Message, "database dirty flag set") {
			log.Error("database dirty flag set (likely due to unclean shutdown): replay required")
			os.Exit(DATABASE_DIRTY)

		} else if strings.Contains(e.Message, "database metadata dirty flag set") {
			log.Error("database metadata dirty flag set (likely due to unclean shutdown): replay required")
			os.Exit(DATABASE_DIRTY)

		} else {
			log.Error("%s", e.Message)
		}
		os.Exit(OTHER_FAIL)

	}).Catch(func(e error) {
		log.Error("%s", e.Error())

	}).Catch(func(interface{}) {
		log.Error("unknown exception")
		os.Exit(OTHER_FAIL)

	}).End()

	os.Exit(SUCCESS)
}
