package wallet_api_plugin

import (
	"encoding/json"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/http_plugin"
	"github.com/eosspark/eos-go/plugins/wallet_plugin"
	"github.com/urfave/cli"
)

const WalletApiPlug = PluginTypeName("WalletApiPlugin")

var templatePlugin = App().RegisterPlugin(WalletApiPlug, NewWalletApiPlugin())

type WalletApiPlugin struct {
	AbstractPlugin
	log log.Logger
}

type WalletApiPluginImpl struct {
}

func NewWalletApiPlugin() *WalletApiPlugin {
	plugin := &WalletApiPlugin{}
	plugin.log = log.New("WalletApiPlugin")
	plugin.log.SetHandler(log.TerminalHandler)
	return plugin
}

func (w *WalletApiPlugin) SetProgramOptions(options *[]cli.Flag) {

}

func (w *WalletApiPlugin) PluginInitialize(options *cli.Context) {
	Try(func() {
		httpPlugin := App().GetPlugin(http_plugin.HttpPlug).(*http_plugin.HttpPlugin)
		if !httpPlugin.IsOnLoopBack() {
			if !httpPlugin.IsSecure() {
				w.log.Error("\n" +
					"********!!!SECURITY ERROR!!!********\n" +
					"*                                  *\n" +
					"* --       Wallet API           -- *\n" +
					"* - EXPOSED to the LOCAL NETWORK - *\n" +
					"* -  HTTP RPC is NOT encrypted   - *\n" +
					"* - Password and/or Private Keys - *\n" +
					"* - are at HIGH risk of exposure - *\n" +
					"*                                  *\n" +
					"************************************\n")
			} else {
				w.log.Warn("\n" +
					"**********SECURITY WARNING**********\n" +
					"*                                  *\n" +
					"* --       Wallet API           -- *\n" +
					"* - EXPOSED to the LOCAL NETWORK - *\n" +
					"* - Password and/or Private Keys - *\n" +
					"* -   are at risk of exposure    - *\n" +
					"*                                  *\n" +
					"************************************\n")
			}
		}
	}).FcLogAndRethrow().End()
}

func (w *WalletApiPlugin) PluginStartup() {
	w.log.Info("starting wallet_api_plugin")

	walletMgr := App().GetPlugin(wallet_plugin.WalletPlug).(*wallet_plugin.WalletPlugin).GetWalletManager()
	h := App().GetPlugin(http_plugin.HttpPlug).(*http_plugin.HttpPlugin)

	h.AddHandler("/v1/wallet/set_timeout", func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}
			var param int64
			err := json.Unmarshal(body, &param)
			if err != nil {
				EosThrow(&EofException{}, "unmarshal set_time params: %s", err.Error())
			}

			walletMgr.SetTimeOut(param)

			byte, _ := json.Marshal(walletApiPluginEmpty{})
			cb(200, byte)
		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "set_timeout", string(body), cb)
		}).End()
	})

	h.AddHandler("/v1/wallet/sign_digest", func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}

			type signDigest struct {
				Digest common.DigestType
				Key    ecc.PublicKey
			}
			var param signDigest
			err := json.Unmarshal(body, &param)
			if err != nil {
				EosThrow(&EofException{}, "unmarshal sign_digest params: %s", err.Error())
			}

			result := walletMgr.SignDigest(param.Digest, param.Key)

			byte, _ := json.Marshal(result)
			cb(201, byte)
		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "sign_digest", string(body), cb)
		}).End()
	})

	//{
	//	std::string("/v1/""wallet""/""sign_transaction"),
	//	[ & wallet_mgr](string, string body, url_response_callback cb) mutable {
	//		try {
	//			if (body.empty()) body = "{}";
	//			const auto & vs = fc::json::json::from_string(body).as < fc::variants > ();
	//			auto result = wallet_mgr.sign_transaction(vs.at(0).as < chain::signed_transaction > (), vs.at(1).as < flat_set < public_key_type > > (), vs.at(2).as < chain::chain_id_type > ());
	//			cb(201, fc::json::to_string(result));
	//		} catch (...) {
	//			http_plugin::handle_exception("wallet", "sign_transaction", body, cb);
	//		}
	//	}
	//},
	h.AddHandler(common.WalletSignTrx, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}
			var param wallet_plugin.SignTrxParams
			err := json.Unmarshal(body, &param)
			if err != nil {
				EosThrow(&EofException{}, "unmarshal sign_transaction params: %s", err.Error())
			}

			result := walletMgr.SignTransaction(param.Txn, param.Keys, param.ChainID)

			byte, _ := json.Marshal(result)
			cb(201, byte)
		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "sign_transaction", string(body), cb)
		}).End()
	})

	//{
	//	std::string("/v1/""wallet""/""create"),
	//	[ & wallet_mgr](string, string body, url_response_callback cb) mutable {
	//		try {
	//			if (body.empty()) body = "{}";
	//			auto result = wallet_mgr.create(fc::json::from_string(body).as < std::string > ());
	//			cb(201, fc::json::to_string(result));
	//		} catch (...) {
	//			http_plugin::handle_exception("wallet", "create", body, cb);
	//		}
	//	}
	//},

	h.AddHandler(common.WalletCreate, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}
			var param string
			err := json.Unmarshal(body, &param)
			if err != nil {
				EosThrow(&EofException{}, "unmarshal create params: %s", err.Error())
			}

			result, err := walletMgr.Create(param)

			if err != nil {
				EosThrow(&EofException{}, "create wallet error: %s", err.Error())
			} else {
				byte, _ := json.Marshal(result)
				cb(201, byte)
			}
		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "create", string(body), cb)
		}).End()
	})

	//{
	//	std::string("/v1/""wallet""/""open"),
	//	[ & wallet_mgr](string, string body, url_response_callback cb) mutable {
	//		try {
	//			if (body.empty()) body = "{}";
	//			wallet_mgr.open(fc::json::from_string(body).as < std::string > ());
	//			eosio::detail::wallet_api_plugin_empty result;
	//			cb(200, fc::json::to_string(result));
	//		} catch (...) {
	//			http_plugin::handle_exception("wallet", "open", body, cb);
	//		}
	//	}
	//},

	h.AddHandler(common.WalletOpen, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}
			var param string
			err := json.Unmarshal(body, &param)
			if err != nil {
				EosThrow(&EofException{}, "unmarshal open params: %s", err.Error())
			}
			walletMgr.Open(param)

			byte, _ := json.Marshal(walletApiPluginEmpty{})
			cb(200, byte)

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "open", string(body), cb)
		}).End()
	})

	//{
	//	std::string("/v1/""wallet""/""lock_all"),
	//	[ & wallet_mgr](string, string body, url_response_callback cb) mutable {
	//		try {
	//			if (body.empty()) body = "{}";
	//			wallet_mgr.lock_all();
	//			eosio::detail::wallet_api_plugin_empty result;
	//			cb(200, fc::json::to_string(result));
	//		} catch (...) {
	//			http_plugin::handle_exception("wallet", "lock_all", body, cb);
	//		}
	//	}
	//},

	h.AddHandler(common.WalletLockAll, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			walletMgr.LockAllwallets()
			byte, _ := json.Marshal(walletApiPluginEmpty{})
			cb(200, byte)
		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "lock_all", string(body), cb)
		}).End()
	})

	//{
	//	std::string("/v1/""wallet""/""lock"),
	//	[ & wallet_mgr](string, string body, url_response_callback cb) mutable {
	//		try {
	//			if (body.empty()) body = "{}";
	//			wallet_mgr.lock(fc::json::from_string(body).as < std::string > ());
	//			eosio::detail::wallet_api_plugin_empty result;
	//			cb(200, fc::json::to_string(result));
	//		} catch (...) {
	//			http_plugin::handle_exception("wallet", "lock", body, cb);
	//		}
	//	}
	//},
	h.AddHandler(common.WalletLock, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}
			var param string
			err := json.Unmarshal(body, &param)
			if err != nil {
				EosThrow(&EofException{}, "unmarshal lock params: %s", err.Error())
			}
			walletMgr.Lock(param)

			byte, _ := json.Marshal(walletApiPluginEmpty{})
			cb(200, byte)

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "lock", string(body), cb)
		}).End()
	})

	//{
	//	std::string("/v1/""wallet""/""unlock"),
	//	[ & wallet_mgr](string, string body, url_response_callback cb) mutable {
	//		try {
	//			if (body.empty()) body = "{}";
	//			const auto & vs = fc::json::json::from_string(body).as < fc::variants > ();
	//			wallet_mgr.unlock(vs.at(0).as < std::string > (), vs.at(1).as < std::string > ());
	//			eosio::detail::wallet_api_plugin_empty result;
	//			cb(200, fc::json::to_string(result));
	//		} catch (...) {
	//			http_plugin::handle_exception("wallet", "unlock", body, cb);
	//		}
	//	}
	//},

	h.AddHandler(common.WalletUnlock, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}
			var param wallet_plugin.UnlockParams
			err := json.Unmarshal(body, &param)
			if err != nil {
				EosThrow(&EofException{}, "unmarshal unlock params: %s", err.Error())
			}
			err = walletMgr.Unlock(param.Name, param.Password)
			if err != nil {
				EosThrow(&EofException{}, "unlock wallet error: %s", err.Error())
			} else {
				byte, _ := json.Marshal(walletApiPluginEmpty{})
				cb(200, byte)
			}
		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "unlock", string(body), cb)
		}).End()
	})

	//{
	//	std::string("/v1/""wallet""/""import_key"),
	//	[ & wallet_mgr](string, string body, url_response_callback cb) mutable {
	//		try {
	//			if (body.empty()) body = "{}";
	//			const auto & vs = fc::json::json::from_string(body).as < fc::variants > ();
	//			wallet_mgr.import_key(vs.at(0).as < std::string > (), vs.at(1).as < std::string > ());
	//			eosio::detail::wallet_api_plugin_empty result;
	//			cb(201, fc::json::to_string(result));
	//		} catch (...) {
	//			http_plugin::handle_exception("wallet", "import_key", body, cb);
	//		}
	//	}
	//},
	h.AddHandler(common.WalletImportKey, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}
			var param wallet_plugin.ImportKeyParams
			err := json.Unmarshal(body, &param)
			if err != nil {
				EosThrow(&EofException{}, "unmarshal import_key params: %s", err.Error())
			}

			walletMgr.ImportKey(param.Name, param.Key)

			byte, _ := json.Marshal(walletApiPluginEmpty{})
			cb(200, byte)

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "import_key", string(body), cb)
		}).End()
	})
	//{
	//	std::string("/v1/""wallet""/""remove_key"),
	//	[ & wallet_mgr](string, string body, url_response_callback cb) mutable {
	//		try {
	//			if (body.empty()) body = "{}";
	//			const auto & vs = fc::json::json::from_string(body).as < fc::variants > ();
	//			wallet_mgr.remove_key(vs.at(0).as < std::string > (), vs.at(1).as < std::string > (), vs.at(2).as < std::string > ());
	//			eosio::detail::wallet_api_plugin_empty result;
	//			cb(201, fc::json::to_string(result));
	//		} catch (...) {
	//			http_plugin::handle_exception("wallet", "remove_key", body, cb);
	//		}
	//	}
	//},

	h.AddHandler(common.WalletRemoveKey, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}
			var param wallet_plugin.RemoveKeyParams
			err := json.Unmarshal(body, &param)
			if err != nil {
				EosThrow(&EofException{}, "unmarshal remove_key params: %s", err.Error())
			}
			walletMgr.RemoveKey(param.Name, param.Password, param.Key)
			byte, _ := json.Marshal(walletApiPluginEmpty{})
			cb(200, byte)

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "remove_key", string(body), cb)
		}).End()
	})

	//{
	//	std::string("/v1/""wallet""/""create_key"),
	//	[ & wallet_mgr](string, string body, url_response_callback cb) mutable {
	//		try {
	//			if (body.empty()) body = "{}";
	//			const auto & vs = fc::json::json::from_string(body).as < fc::variants > ();
	//			auto result = wallet_mgr.create_key(vs.at(0).as < std::string > (), vs.at(1).as < std::string > ());
	//			cb(201, fc::json::to_string(result));
	//		} catch (...) {
	//			http_plugin::handle_exception("wallet", "create_key", body, cb);
	//		}
	//	}
	//},
	h.AddHandler(common.WalletCreateKey, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}
			var param wallet_plugin.CreateKeyParams
			err := json.Unmarshal(body, &param)
			if err != nil {
				EosThrow(&EofException{}, "unmarshal create_key params: %s", err.Error())
			}
			result := walletMgr.CreateKey(param.Name, param.KeyType)
			byte, _ := json.Marshal(result)
			cb(200, byte)

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "create_key", string(body), cb)
		}).End()
	})

	//{
	//	std::string("/v1/""wallet""/""list_wallets"),
	//	[ & wallet_mgr](string, string body, url_response_callback cb) mutable {
	//		try {
	//			if (body.empty()) body = "{}";
	//			auto result = wallet_mgr.list_wallets();
	//			cb(200, fc::json::to_string(result));
	//		} catch (...) {
	//			http_plugin::handle_exception("wallet", "list_wallets", body, cb);
	//		}
	//	}
	//},
	h.AddHandler(common.WalletList, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}

			result := walletMgr.ListWallets()
			byte, _ := json.Marshal(result)
			cb(200, byte)

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "list_wallets", string(body), cb)
		}).End()
	})
	//{
	//	std::string("/v1/""wallet""/""list_keys"),
	//	[ & wallet_mgr](string, string body, url_response_callback cb) mutable{
	//	try{
	//		if (body.empty()) body = "{}";
	//		const auto & vs = fc::json::json::from_string(body).as < fc::variants > ();
	//		auto result = wallet_mgr.list_keys(vs.at(0).as < std::string > (), vs.at(1).as < std::string > ());
	//		cb(200, fc::json::to_string(result));
	//	} catch (...){
	//	http_plugin::handle_exception("wallet", "list_keys", body, cb);
	//}
	//}
	h.AddHandler(common.WalletListKeys, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}
			var param wallet_plugin.ListKeysParams
			err := json.Unmarshal(body, &param)
			if err != nil {
				EosThrow(&EofException{}, "unmarshal list_keys params: %s", err.Error())
			}
			result := walletMgr.ListKeys(param.Name, param.Password)

			byte, err := json.Marshal(result)
			if err != nil {
				EosThrow(&EofException{}, "marshal list_keys result: %s", err.Error())
			}
			cb(200, byte)
		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "list_keys", string(body), cb)
		}).End()
	})
	//{
	//	std::string("/v1/""wallet""/""get_public_keys"),
	//	[ & wallet_mgr](string, string body, url_response_callback cb) mutable {
	//		try {
	//			if (body.empty()) body = "{}";
	//			auto result = wallet_mgr.get_public_keys();
	//			cb(200, fc::json::to_string(result));
	//		} catch (...) {
	//			http_plugin::handle_exception("wallet", "get_public_keys", body, cb);
	//		}
	//	}
	//}
	h.AddHandler(common.WalletPublicKeys, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			if len(body) == 0 {
				body = []byte{123, 125} //"{}"
			}
			result := walletMgr.GetPublicKeys()

			byte, err := json.Marshal(result)
			if err != nil {
				EosThrow(&EofException{}, "marshal get_public_keys result: %s", err.Error())
			}
			cb(200, byte)

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "wallet", "get_public_keys", string(body), cb)
		}).End()
	})

}

func (w *WalletApiPlugin) PluginShutdown() {
}

type walletApiPluginEmpty struct {
}
