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

// {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_block_header_state"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_block_header_state(fc::json::from_string(body).as < chain_apis::read_only::get_block_header_state_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_block_header_state", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_account"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_account(fc::json::from_string(body).as < chain_apis::read_only::get_account_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_account", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_code"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_code(fc::json::from_string(body).as < chain_apis::read_only::get_code_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_code", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_code_hash"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_code_hash(fc::json::from_string(body).as < chain_apis::read_only::get_code_hash_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_code_hash", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_abi"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_abi(fc::json::from_string(body).as < chain_apis::read_only::get_abi_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_abi", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_raw_code_and_abi"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_raw_code_and_abi(fc::json::from_string(body).as < chain_apis::read_only::get_raw_code_and_abi_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_raw_code_and_abi", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_raw_abi"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_raw_abi(fc::json::from_string(body).as < chain_apis::read_only::get_raw_abi_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_raw_abi", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_table_rows"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_table_rows(fc::json::from_string(body).as < chain_apis::read_only::get_table_rows_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_table_rows", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_table_by_scope"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_table_by_scope(fc::json::from_string(body).as < chain_apis::read_only::get_table_by_scope_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_table_by_scope", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_currency_balance"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_currency_balance(fc::json::from_string(body).as < chain_apis::read_only::get_currency_balance_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_currency_balance", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_currency_stats"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_currency_stats(fc::json::from_string(body).as < chain_apis::read_only::get_currency_stats_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_currency_stats", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_producers"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_producers(fc::json::from_string(body).as < chain_apis::read_only::get_producers_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_producers", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_producer_schedule"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_producer_schedule(fc::json::from_string(body).as < chain_apis::read_only::get_producer_schedule_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_producer_schedule", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_scheduled_transactions"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_scheduled_transactions(fc::json::from_string(body).as < chain_apis::read_only::get_scheduled_transactions_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_scheduled_transactions", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"abi_json_to_bin"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.abi_json_to_bin(fc::json::from_string(body).as < chain_apis::read_only::abi_json_to_bin_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "abi_json_to_bin", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"abi_bin_to_json"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.abi_bin_to_json(fc::json::from_string(body).as < chain_apis::read_only::abi_bin_to_json_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "abi_bin_to_json", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_required_keys"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_required_keys(fc::json::from_string(body).as < chain_apis::read_only::get_required_keys_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_required_keys", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"get_transaction_id"),
// 	[ro_api](string, string body, url_response_callback cb) mutable {
// 		ro_api.validate();
// 		try {
// 			if (body.empty()) body = "{}";
// 			auto result = ro_api.get_transaction_id(fc::json::from_string(body).as < chain_apis::read_only::get_transaction_id_params > ());
// 			cb(200, fc::json::to_string(result));
// 		} catch (...) {
// 			http_plugin::handle_exception("chain", "get_transaction_id", body, cb);
// 		}
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"push_block"),
// 	[rw_api](string, string body, url_response_callback cb) mutable {
// 		if (body.empty()) body = "{}";
// 		rw_api.validate();
// 		rw_api.push_block(fc::json::from_string(body).as < chain_apis::read_write::push_block_params > (), [cb, body](const fc::static_variant < fc::exception_ptr, chain_apis::read_write::push_block_results > & result) {
// 			if (result.contains < fc::exception_ptr > ()) {
// 				try {
// 					result.get < fc::exception_ptr > () - > dynamic_rethrow_exception();
// 				} catch (...) {
// 					http_plugin::handle_exception("chain", "push_block", body, cb);
// 				}
// 			} else {
// 				cb(202, result.visit(async_result_visitor()));
// 			}
// 		});
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"push_transaction"),
// 	[rw_api](string, string body, url_response_callback cb) mutable {
// 		if (body.empty()) body = "{}";
// 		rw_api.validate();
// 		rw_api.push_transaction(fc::json::from_string(body).as < chain_apis::read_write::push_transaction_params > (), [cb, body](const fc::static_variant < fc::exception_ptr, chain_apis::read_write::push_transaction_results > & result) {
// 			if (result.contains < fc::exception_ptr > ()) {
// 				try {
// 					result.get < fc::exception_ptr > () - > dynamic_rethrow_exception();
// 				} catch (...) {
// 					http_plugin::handle_exception("chain", "push_transaction", body, cb);
// 				}
// 			} else {
// 				cb(202, result.visit(async_result_visitor()));
// 			}
// 		});
// 	}
// }, {
// 	std::string("/v1/"
// 		"chain"
// 		"/"
// 		"push_transactions"),
// 	[rw_api](string, string body, url_response_callback cb) mutable {
// 		if (body.empty()) body = "{}";
// 		rw_api.validate();
// 		rw_api.push_transactions(fc::json::from_string(body).as < chain_apis::read_write::push_transactions_params > (), [cb, body](const fc::static_variant < fc::exception_ptr, chain_apis::read_write::push_transactions_results > & result) {
// 			if (result.contains < fc::exception_ptr > ()) {
// 				try {
// 					result.get < fc::exception_ptr > () - > dynamic_rethrow_exception();
// 				} catch (...) {
// 					http_plugin::handle_exception("chain", "push_transactions", body, cb);
// 				}
// 			} else {
// 				cb(202, result.visit(async_result_visitor()));
// 			}
// 		});
// 	}
// }
