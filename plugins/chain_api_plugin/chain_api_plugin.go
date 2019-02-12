package chain_api_plugin

import (
	"encoding/json"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/eosspark/eos-go/plugins/http_plugin"
	"github.com/urfave/cli"
)

const ChainAPiPlug = PluginTypeName("ChainAPiPlugin")

var chainAPiPlugin = App().RegisterPlugin(ChainAPiPlug, NewChainAPiPlugin())

type ChainApiPlugin struct {
	AbstractPlugin
	my  *ChainApiPluginImpl
	log log.Logger
}

type ChainApiPluginImpl struct {
	db *chain.Controller
}

func NewChainAPiPlugin() *ChainApiPlugin {
	plugin := &ChainApiPlugin{}
	plugin.my = &ChainApiPluginImpl{}
	plugin.log = log.New("ChainAPiPlugin")
	plugin.log.SetHandler(log.TerminalHandler)
	return plugin
}

func (c *ChainApiPlugin) SetProgramOptions(options *[]cli.Flag) {
}

func (c *ChainApiPlugin) PluginInitialize(options *cli.Context) {
}

func (c *ChainApiPlugin) PluginStartup() {
	c.log.Info("starting chain_api_plugin")
	c.my.db = App().GetPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin).Chain()

	httpPlugin := App().GetPlugin(http_plugin.HttpPlug).(*http_plugin.HttpPlugin)

	ROApi := App().GetPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin).GetReadOnlyApi()

	//TODO read_only api
	ROApi.SetShortenAbiErrors(httpPlugin.VerboseErrors())

	httpPlugin.AddHandler(common.GetInfoFunc, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			result := ROApi.GetInfo()

			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "chain", "get_info", string(body), cb)
		}).End()
	})

	httpPlugin.AddHandler(common.GetBlockFunc, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte("{}")
			}

			var param chain_plugin.GetBlockParams
			if err := json.Unmarshal(body, &param); err != nil {
				EosThrow(&EofException{}, "marshal get_block params: %s", err.Error())
			}

			result := ROApi.GetBlock(param)

			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "chain", "get_block", string(body), cb)
		}).End()
	})

	httpPlugin.AddHandler(common.GetBlockHeaderStateFunc, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte("{}")
			}

			var param chain_plugin.GetBlockHeaderStateParams
			if err := json.Unmarshal(body, &param); err != nil {
				EosThrow(&EofException{}, "marshal get_block_header_state params: %s", err.Error())
			}

			result := ROApi.GetBlockHeaderState(param)

			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "chain", "get_block_header_state", string(body), cb)
		}).End()
	})

	httpPlugin.AddHandler(common.GetAccountFunc, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte("{}")
			}

			var param chain_plugin.GetAccountParams
			if err := json.Unmarshal(body, &param); err != nil {
				EosThrow(&EofException{}, "marshal get_account params: %s", err.Error())
			}

			result := ROApi.GetAccount(param)

			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "chain", "get_account", string(body), cb)
		}).End()
	})

	httpPlugin.AddHandler(common.GetAbiFunc, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte("{}")
			}

			var param chain_plugin.GetAbiParams
			if err := json.Unmarshal(body, &param); err != nil {
				EosThrow(&EofException{}, "marshal get_abi params: %s", err.Error())
			}

			result := ROApi.GetAbi(param)

			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "chain", "get_abi", string(body), cb)
		}).End()
	})

	httpPlugin.AddHandler(common.GetCodeFunc, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte("{}")
			}

			var param chain_plugin.GetCodeParams
			if err := json.Unmarshal(body, &param); err != nil {
				EosThrow(&EofException{}, "marshal get_code params: %s", err.Error())
			}

			result := ROApi.GetCode(param)

			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "chain", "get_code", string(body), cb)
		}).End()
	})

	httpPlugin.AddHandler(common.GetCurrencyBalanceFunc, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte("{}")
			}

			var param chain_plugin.GetCurrencyBalanceParams
			if err := json.Unmarshal(body, &param); err != nil {
				EosThrow(&EofException{}, "marshal get_currency_balance params: %s", err.Error())
			}

			result := ROApi.GetCurrencyBalance(param)

			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "chain", "get_currency_balance", string(body), cb)
		}).End()
	})

	httpPlugin.AddHandler(common.GetCurrencyStatsFunc, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte("{}")
			}

			var param chain_plugin.GetCurrencyStatsParams
			if err := json.Unmarshal(body, &param); err != nil {
				EosThrow(&EofException{}, "marshal get_currency_stats params: %s", err.Error())
			}

			result := ROApi.GetCurrencyStats(param)

			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "chain", "get_currency_stats", string(body), cb)
		}).End()
	})

	httpPlugin.AddHandler(common.GetRequiredKeys, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte("{}")
			}

			var param chain_plugin.GetRequiredKeysParams
			if err := json.Unmarshal(body, &param); err != nil {
				EosThrow(&EofException{}, "marshal get_required_keys params: %s", err.Error())
			}

			result := ROApi.GetRequiredKeys(param)

			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "chain", "get_required_keys", string(body), cb)
		}).End()
	})

	httpPlugin.AddHandler(common.GetTableFunc, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte("{}")
			}

			var param chain_plugin.GetTableRowsParams
			if err := json.Unmarshal(body, &param); err != nil {
				EosThrow(&EofException{}, "marshal get_table params: %s", err.Error())
			}

			result := ROApi.GetTableRows(param)

			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "chain", "get_table", string(body), cb)
		}).End()
	})

	//TODO read_write api
	RWApi := App().GetPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin).GetReadWriteApi()

	httpPlugin.AddHandler(common.PushTxnFunc, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		if len(body) == 0 {
			body = []byte("{}")
		}

		RWApi.Validate()

		var param chain_plugin.PushTransactionParams
		if err := json.Unmarshal(body, &param); err != nil {
			EosThrow(&EofException{}, "marshal push_transaction params: %s", err.Error())
		}

		RWApi.PushTransaction(param, func(result common.StaticVariant) {
			if exception, ok := result.(Exception); ok {
				http_plugin.HandleException(exception, "chain", "push_transaction", string(body), cb)
			} else {
				if byte, err := json.Marshal(result); err == nil {
					cb(200, byte)
				} else {
					http_plugin.HandleException(err, "chain", "push_transaction", string(body), cb)
				}
			}
		})
	})

}

func (c *ChainApiPlugin) PluginShutdown() {
}
