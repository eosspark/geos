package chain_plugin

import  (
	. "github.com/eosspark/eos-go/appbase/app/include"
	"github.com/eosspark/eos-go/appbase/app"
)

type ChainPlugin struct {
	AbstractPlugin
}

func init() {
	var chainPlugin = new(ChainPlugin)
	chainPlugin.Name = "ChainPlugin"
	chainPlugin.State = Registered
	app.App.RegisterPlugin(chainPlugin)
}

func (chainPlugin *ChainPlugin)SetProgramOptions() {

}
func (chainPlugin *ChainPlugin) PluginInitialize() {

}
func (chainPlugin *ChainPlugin) PluginStartUp() {

}
func (chainPlugin *ChainPlugin)PluginShutDown() {

}

func (chainPlugin *ChainPlugin) GetName() string {
	return chainPlugin.Name
}
func (chainPlugin *ChainPlugin) GetState() State {
	return chainPlugin.State

}