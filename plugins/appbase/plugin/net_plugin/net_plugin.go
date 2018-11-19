package net_plugin

import (
	"fmt"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/urfave/cli"
)




type NetPlugin struct {
	AbstractPlugin
}

var netPlugin *NetPlugin

func init() {
	netPlugin = new(NetPlugin)
	netPlugin.Plugin = netPlugin
	netPlugin.Name = "NetPlugin"
	netPlugin.State = Registered
	App.RegisterPlugin(netPlugin)
}


func (netPlugin *NetPlugin) Exec() bool {
	fmt.Println("执行net_plugin的exec内容")
	return true
}

func (netPlugin *NetPlugin) SetProgramOptions(options *cli.App) {
	fmt.Println("NetPlugin SetProgramOptions")
	fmt.Println(options.Name)
}

func (netPlugin *NetPlugin) PluginInitialize() {
	fmt.Println("netPlugin PluginInitialize")
}

//var loop = true

func (netPlugin *NetPlugin) PluginStartUp() {

	fmt.Println("netPlugin PluginStartUp")

	//go func() {
	//	for loop {
	//		fmt.Println("NetPlugin PluginStartUp")
	//	}
	//}()
}

func (netPlugin *NetPlugin) PluginShutDown() {
	//loop = false
	fmt.Println("NetPlugin PluginShutDown")
}

