package main

import (
	"fmt"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"github.com/eosspark/eos-go/plugins/http_plugin"
	"github.com/urfave/cli"
	"os"
)

func MakeArguments(values ...string) {
	options := append([]string(values), "--") // use "--" to divide arguments

	osArgs := make([]string, len(os.Args)+len(options))
	copy(osArgs[:1], os.Args[:1])
	copy(osArgs[1:len(options)+1], options)
	copy(osArgs[len(options)+1:], os.Args[1:])

	os.Args = osArgs
}

type httpPluginTester struct {
	*http_plugin.HttpPlugin
	io  *asio.IoContext
	app *cli.App
}

func main() {
	ppt := new(httpPluginTester)
	ppt.io = asio.NewIoContext()
	ppt.app = cli.NewApp()

	app := cli.NewApp()
	app.Name = "nodeos"
	app.Version = "0.1.0beta"

	ppt.HttpPlugin = http_plugin.NewHttpPlugin(ppt.io)
	ppt.SetProgramOptions(&ppt.app.Flags)
	MakeArguments("--http-server-address 127.0.0.1:8891")
	ppt.app.Action = func(c *cli.Context) {
		ppt.PluginInitialize(c)
	}
	err := ppt.app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}

	ppt.PluginStartup()

	for {

	}
}
