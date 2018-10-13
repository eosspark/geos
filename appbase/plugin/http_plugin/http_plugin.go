package http_plugin

import (
	. "github.com/eosspark/eos-go/appbase/app/include"
	"github.com/eosspark/eos-go/appbase/app"
)


type HttpPlugin struct {
	AbstractPlugin
}

func init () {
	var httpPlugin HttpPlugin
	httpPlugin.State = Registered
	httpPlugin.Name = "HttpPlugin"
	app.App.RegisterPlugin(&httpPlugin)

}

func (httpPlugin *HttpPlugin) SetProgramOptions() {

}

func (httpPlugin *HttpPlugin) PluginInitialize() {

}
func (httpPlugin *HttpPlugin) PluginStartUp() {

}
func (httpPlugin *HttpPlugin) PluginShutDown() {

}

func (httpPlugin *HttpPlugin) GetName() string {
	return httpPlugin.Name
}

func (httpPlugin *HttpPlugin) GetState() State {
	return httpPlugin.State
}