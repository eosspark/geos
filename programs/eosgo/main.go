package main

import (
	"os"
	"strings"

	. "github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/eosspark/eos-go/plugins/http_plugin"
	"github.com/eosspark/eos-go/plugins/net_plugin"
	"github.com/eosspark/eos-go/plugins/producer_plugin"
	"github.com/eosspark/eos-go/plugins/wallet_plugin"

	_ "github.com/eosspark/eos-go/plugins/chain_api_plugin"
	_ "github.com/eosspark/eos-go/plugins/console_plugin"
	_ "github.com/eosspark/eos-go/plugins/net_api_plugin"
	_ "github.com/eosspark/eos-go/plugins/wallet_api_plugin"
	_ "github.com/eosspark/eos-go/plugins/wallet_plugin"
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

const EOSGO_VERSION = 0x0

//go run main.go -e -p eosio --private-key [\"EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV\",\"5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3\"]
//--plugin ChainApiPlugin --plugin WalletPlugin --plugin WalletApiPlugin --plugin ChainApiPlugin --plugin NetApiPlugin --plugin ProducerPlugin --max-transaction-age=999999999
func main() {
	try.Try(func() {
		App().SetVersion(EOSGO_VERSION)
		App().SetDefaultDataDir()
		App().SetDefaultConfigDir()
		if !App().Initialize([]PluginTypeName{
			chain_plugin.ChainPlug,
			producer_plugin.ProducerPlug,
			http_plugin.HttpPlug,
			net_plugin.NetPlug,
			wallet_plugin.WalletPlug,
		}) {
			os.Exit(INITIALIZE_FAIL)
		}
		App().StartUp()
		App().Exec()

	}).Catch(func(e *ExtractGenesisStateException) {
		os.Exit(EXTRACTED_GENESIS)

	}).Catch(func(e *FixedReversibleDbException) {
		os.Exit(FIXED_REVERSIBLE)

	}).Catch(func(e *NodeManagementSuccess) {
		os.Exit(NODE_MANAGEMENT_SUCCESS)

	}).Catch(func(e Exception) {
		if e.Code() == StdExceptionCode {
			if strings.Contains(e.DetailMessage(), "database dirty flag set") {
				log.Error("database dirty flag set (likely due to unclean shutdown): replay required")
				os.Exit(DATABASE_DIRTY)
			} else if strings.Contains(e.DetailMessage(), "database metadata dirty flag set") {
				log.Error("database metadata dirty flag set (likely due to unclean shutdown): replay required")
				os.Exit(DATABASE_DIRTY)
			}
		}
		log.Error(e.DetailMessage())
		os.Exit(OTHER_FAIL)

	}).Catch(func(interface{}) {
		log.Error("unknown exception")
		os.Exit(OTHER_FAIL)

	}).End()

	os.Exit(SUCCESS)
}
