package producer_api_plugin

import (
	"encoding/json"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types/generated_containers"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/http_plugin"
	"github.com/eosspark/eos-go/plugins/producer_plugin"
	"github.com/urfave/cli"
)

type ProducerApiPlugin struct {
	AbstractPlugin
	my  *ProducerApiPluginImpl
	log log.Logger
}

type ProducerApiPluginImpl struct {
	db *chain.Controller
}

func NewProducerApiPlugin() *ProducerApiPlugin {
	plugin := &ProducerApiPlugin{}
	plugin.my = &ProducerApiPluginImpl{}
	plugin.log = log.New("ProducerApiPlugin")
	plugin.log.SetHandler(log.TerminalHandler)
	return plugin
}

const ProducerApiPlug = PluginTypeName("ProducerApiPlugin")

var producerApiPlugin = App().RegisterPlugin(ProducerApiPlug, NewProducerApiPlugin())

func (c *ProducerApiPlugin) SetProgramOptions(options *[]cli.Flag) {
}

func (c *ProducerApiPlugin) PluginInitialize(options *cli.Context) {
	Try(func() {
		httpPlugin := App().GetPlugin(http_plugin.HttpPlug).(*http_plugin.HttpPlugin)
		if !httpPlugin.IsOnLoopBack() {
			if !httpPlugin.IsSecure() {
				c.log.Warn("\n" +
					"**********SECURITY WARNING**********\n" +
					"*                                  *\n" +
					"* --         PRODUCER API       -- *\n" +
					"* - EXPOSED to the LOCAL NETWORK - *\n" +
					"* - USE ONLY ON SECURE NETWORKS! - *\n" +
					"*                                  *\n" +
					"************************************\n")
			}
		}
	}).FcLogAndRethrow().End()
}

func (c *ProducerApiPlugin) PluginStartup() {
	c.log.Info("starting producer_api_plugin")
	httpPlugin := App().GetPlugin(http_plugin.HttpPlug).(*http_plugin.HttpPlugin)
	proApi := App().GetPlugin(producer_plugin.ProducerPlug).(*producer_plugin.ProducerPlugin)

	httpPlugin.AddHandler(common.ProducerSetWhitelistBlacklist, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {

		Try(func() {
			params := producer_plugin.WhitelistAndBlacklistParam{}
			w := producer_plugin.WhitelistAndBlacklist{}
			if err := json.Unmarshal(body, &params); err != nil {
				EosThrow(&exception.EofException{}, "marshal SetBlacklistWhitelist params: %s", err.Error())
			}
			if len(params.ContractBlacklist) > 0 {
				w.ContractBlacklist = generated.NewAccountNameSet()
				w.ContractBlacklist.Add(params.ContractBlacklist...)
			}
			if len(params.ContractWhitelist) > 0 {
				w.ContractWhitelist = generated.NewAccountNameSet()
				w.ContractWhitelist.Add(params.ContractWhitelist...)
			}
			if len(params.ActorBlacklist) > 0 {
				w.ActorBlacklist = generated.NewAccountNameSet()
				w.ActorBlacklist.Add(params.ActorBlacklist...)
			}
			if len(params.ActorWhitelist) > 0 {
				w.ActorWhitelist = generated.NewAccountNameSet()
				w.ActorWhitelist.Add(params.ActorWhitelist...)
			}
			if len(params.ActionBlacklist) > 0 {
				w.ActionBlacklist = generated.NewNamePairSet()
				for _, v := range params.ActionBlacklist {
					n := common.NamePair{v.First, v.Second}
					w.ActionBlacklist.Add(n)
				}
			}
			if len(params.KeyBlacklist) > 0 {
				w.KeyBlacklist = generated.NewPublicKeySet()
				for _, keyStr := range params.KeyBlacklist {
					key, _ := ecc.NewPublicKey(keyStr)
					w.KeyBlacklist.Add(key)
				}
			}
			proApi.SetWhitelistBlacklist(w)
			if byte, err := json.Marshal("set is ok"); err == nil {
				cb(200, byte)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "producer", "set_whitelist_blacklist", string(body), cb)
		}).End()
	})

	httpPlugin.AddHandler(common.ProducerGetWhitelistBlacklist, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			data := proApi.GetWhitelistBlacklist()
			result, err := json.Marshal(data)
			if err != nil {
				log.Error("producer_plugin ProducerGetWhitelistBlacklist is error:", err)
			}
			cb(200, result)
		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "producer", "get_whitelist_blacklist", string(body), cb)
		}).End()
	})
}

func (c *ProducerApiPlugin) PluginShutdown() {
}
