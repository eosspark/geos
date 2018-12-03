package http_plugin

import (
	"os"
	"testing"

	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"gopkg.in/urfave/cli.v1"
)

func makeArguments(values ...string) {
	options := append([]string(values), "--") // use "--" to divide arguments

	osArgs := make([]string, len(os.Args)+len(options))
	copy(osArgs[:1], os.Args[:1])
	copy(osArgs[1:len(options)+1], options)
	copy(osArgs[len(options)+1:], os.Args[1:])

	os.Args = osArgs
}

func TestHttpPlugin_PluginInitialize(t *testing.T) {
	makeArguments("--http-server-address", "127.0.0.1:8989")

	io := asio.NewIoContext()
	app := cli.NewApp()
	app.Name = "nodeos"
	app.Version = "0.1.0beta"

	httpPlugin := NewHttpPlugin(io)
	httpPlugin.PluginInitialize(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("net_plugin start !!")
	httpPlugin.PluginStartup()

}
