package chain_api_plugin

import (
	"encoding/json"
	"github.com/eosspark/eos-go/chain"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/eosspark/eos-go/plugins/http_plugin"
	"github.com/urfave/cli"
)

const ChainAPiPlug = PluginTypeName("ChainAPiPlugin")

var templatePlugin = App().RegisterPlugin(ChainAPiPlug, NewChainAPiPlugin())

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
	ROApi := App().GetPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin).GetReadOnlyApi()
	//RWApi := App().GetPlugin(chain_plugin.ChainPlug).(*chain_plugin.ChainPlugin).GetReadWriteApi()

	httpPlugin := App().GetPlugin(http_plugin.HttpPlug).(*http_plugin.HttpPlugin)

	ROApi.SetShortenAbiErrors(httpPlugin.VerboseErrors())

	httpPlugin.AddHandler("/v1/chain/get_info", func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			result := ROApi.GetInfo()
			byte, err := json.Marshal(result)
			if err != nil {
				EosThrow(&EofException{}, "marshal get_info result: %s", err.Error())
			}
			cb(200, byte)
		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "chain", "get_info", string(body), cb)
		})
	})

	httpPlugin.AddHandler("/v1/chain/get_block", func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}
			var param chain_plugin.GetBlockParams
			err := json.Unmarshal(body, &param)
			if err != nil {
				EosThrow(&EofException{}, "marshal get_block params: %s", err.Error())
			}
			result := ROApi.GetBlock(param)
			byte, err := json.Marshal(result)
			if err != nil {
				EosThrow(&EofException{}, "marshal get_info result: %s", err.Error())
			}
			cb(200, byte)
		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "chain", "get_block", string(body), cb)
		})
	})

}

func (c *ChainApiPlugin) PluginShutdown() {
}

// {
// 	std::string("/v1/""chain""/""get_info"), [this, ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_info(fc::json::from_string(body).as < chain_apis::read_only::get_info_params > ());
// 			cb(200 l, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_info", body, cb);
// 		}
// 	}
// }

// {
// 	std::string("/v1/""chain""/""get_block"),
// 	  [this, ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_block(fc::json::from_string(body).as < chain_apis::read_only::get_block_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_block", body, cb);
// 		}
// 	}
// },
// apis returns the collection of RPC descriptors this node offers.

//fc::static_variant < fc::exception_ptr, chain_apis::read_write::push_transaction_results > & result
//func next(result  Variants){
//	err,ok:= result["exception"]
//	if ok {
//		Try(func() {
//
//		}).Catch(func(interface{}) {
//		//http_plugin::handle_exception("chain", "push_transaction", body, cb);
//		})
//	}else{
//
//		fmt.Fprintf(ctx, "%s", reJson)
//
//		//fmt.Fprintf(ctx, "%s", re)
//		ctx.SetContentType("text/plain; charset=utf8")
//		// Set arbitrary headers
//		ctx.Response.Header.Set("X-My-Header", "my-header-value")
//		// Set cookies
//		var c fasthttp.Cookie
//		c.SetKey("cookie-name")
//		c.SetValue("cookie-value")
//		ctx.Response.Header.SetCookie(&c)
//	}
//
//}

//[cb, body](const fc::static_variant < fc::exception_ptr, chain_apis::read_write::push_transaction_results > & result) {
//if (result.contains < fc::exception_ptr > ()) {
//try {
//result.get < fc::exception_ptr > () - > dynamic_rethrow_exception();
//} catch (...) {
//http_plugin::handle_exception("chain", "push_transaction", body, cb);
//}
//} else {
//cb(202, result.visit(async_result_visitor()));
//}
//}
