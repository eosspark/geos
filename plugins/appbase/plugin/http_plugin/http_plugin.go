package http_plugin

import (
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/urfave/cli"
	"fmt"
)


type HttpPlugin struct {
	AbstractPlugin
}
var httpPlugin *HttpPlugin

func init () {
	httpPlugin = new(HttpPlugin)
	httpPlugin.Plugin = httpPlugin
	httpPlugin.State = Registered
	httpPlugin.Name = "HttpPlugin"
	App.RegisterPlugin(httpPlugin)
}

func (httpPlugin *HttpPlugin) SetProgramOptions(options *cli.App) {
	fmt.Println("httpPlugin SetProgramOptions")
	fmt.Println(options.Name)
}

func (httpPlugin *HttpPlugin) PluginInitialize() {
	fmt.Println("httpPlugin PluginInitialize")
}
func (httpPlugin *HttpPlugin) PluginStartUp() {
	fmt.Println("httpPlugin PluginStartUp")
}
func (httpPlugin *HttpPlugin) PluginShutDown() {

}

