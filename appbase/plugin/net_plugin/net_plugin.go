package net_plugin

import (
	"appbase/app"
	. "appbase/app/include"
	"fmt"
)

type NetPlugin struct {
	AbstractPlugin
}

var netPlugin *NetPlugin

func init() {
	netPlugin = new(NetPlugin)
	netPlugin.Name = "NetPlugin"
	netPlugin.State = Registered
	app.App.RegisterPlugin(netPlugin)
}

func (netPlugin *NetPlugin) GetName() string {
	return netPlugin.Name
}

func (netPlugin *NetPlugin) GetState() State {
	return netPlugin.State
}

func (netPlugin *NetPlugin) Exec() bool {
	fmt.Println("执行net_plugin的exec内容")
	return true
}

func (netPlugin *NetPlugin) SetProgramOptions() {
	fmt.Println("NetPlugin SetProgramOptions")

}

func (netPlugin *NetPlugin) PluginInitialize() {
	fmt.Println("NetPlugin完成自己初始化操作")
	netPlugin.State = Initialized

}

var loop = true

func (netPlugin *NetPlugin) PluginStartUp() {

	go func() {
		for loop {
			fmt.Println("NetPlugin PluginStartUp")
		}
	}()
}

func (netPlugin *NetPlugin) PluginShutDown() {
	loop = false
	fmt.Println("NetPlugin PluginShutDown")
}

//var Net = NewNet_Plugin()

//func NewNet_Plugin() *NetPlugin {
//	net := new(NetPlugin)
//	net.Plugin.Initialize = net.Initialize
//	net.Plugin.StartUp = net.StartUp
//	net.Plugin.Exec = net.Exec
//	return net
//}
//
//func (net_plugin *NetPlugin) Initialize () bool {
//	fmt.Println("执行net_plugin的Initialize内容")
//	return true
//}
//
