package net_api_plugin

import (
	"encoding/json"

	"github.com/eosspark/eos-go/common"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/http_plugin"
	"github.com/eosspark/eos-go/plugins/net_plugin"

	"github.com/urfave/cli"
)

const NetApiPlug = PluginTypeName("NetApiPlugin")

var netApiPlugin = App().RegisterPlugin(NetApiPlug, NewNetApiPlugin())

type NetApiPlugin struct {
	AbstractPlugin
	log log.Logger
}

type NetAPiPluginImpl struct {
}

func NewNetApiPlugin() *NetApiPlugin {
	plugin := &NetApiPlugin{}
	plugin.log = log.New("netApiPlugin")
	plugin.log.SetHandler(log.TerminalHandler)
	return plugin
}

func (n *NetApiPlugin) SetProgramOptions(options *[]cli.Flag) {}

func (n *NetApiPlugin) PluginInitialize(options *cli.Context) {
	Try(func() {
		httpPlugin := App().GetPlugin(http_plugin.HttpPlug).(*http_plugin.HttpPlugin)
		if !httpPlugin.IsOnLoopBack() {
			if !httpPlugin.IsSecure() {
				n.log.Warn("\n" +
					"**********SECURITY WARNING**********\n" +
					"*                                  *\n" +
					"* --         Net API            -- *\n" +
					"* - EXPOSED to the LOCAL NETWORK - *\n" +
					"* - USE ONLY ON SECURE NETWORKS! - *\n" +
					"*                                  *\n" +
					"************************************\n")
			}
		}
	}).FcLogAndRethrow().End()
}

func (n *NetApiPlugin) PluginStartup() {
	n.log.Info("starting net_api_plugin")

	netMgr := App().GetPlugin(net_plugin.NetPlug).(*net_plugin.NetPlugin)
	httpPlugin := App().GetPlugin(http_plugin.HttpPlug).(*http_plugin.HttpPlugin)

	httpPlugin.AddHandler(common.NetConnect, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			//127.0.0.1:9111 ->[34 49 50 55 46 48 46 48 46 49 58 57 49 49 49 34],delete 34...34 ""
			result := netMgr.Connect(string(body[1 : len(body)-1]))
			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "net", "connect", string(body), cb)
		}).End()
	})

	httpPlugin.AddHandler(common.NetDisconnect, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			result := netMgr.Disconnect(string(body[1 : len(body)-1]))
			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "net", "disconnect", string(body), cb)
		}).End()
	})

	httpPlugin.AddHandler(common.NetStatus, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			result := netMgr.Status(string(body[1 : len(body)-1]))
			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "net", "status", string(body), cb)
		}).End()
	})

	httpPlugin.AddHandler(common.NetConnections, func(source string, body []byte, cb http_plugin.UrlResponseCallback) {
		Try(func() {
			result := netMgr.Connections()
			if byte, err := json.Marshal(result); err == nil {
				cb(200, byte)
			} else {
				Throw(err)
			}

		}).Catch(func(e interface{}) {
			http_plugin.HandleException(e, "net", "connections", string(body), cb)
		}).End()
	})
}

func (n *NetApiPlugin) PluginShutdown() {}
