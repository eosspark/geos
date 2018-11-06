package chain_plugin

import  (
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/urfave/cli"
	"fmt"
)

type ChainPlugin struct {
	AbstractPlugin
}

var chainPlugin *ChainPlugin

func init() {
	chainPlugin = new(ChainPlugin)
	chainPlugin.Plugin = chainPlugin
	chainPlugin.Name = "ChainPlugin"
	chainPlugin.State = Registered
	App.RegisterPlugin(chainPlugin)
}


func (chainPlugin *ChainPlugin) SetProgramOptions(options *cli.App) {
	fmt.Println("ChainPlugin SetProgramOptions")
	fmt.Println(options.Name)
}
func (chainPlugin *ChainPlugin) PluginInitialize() {
	fmt.Println("chainPlugin PluginInitialize")
}
func (chainPlugin *ChainPlugin) PluginStartUp() {
	fmt.Println("chainPlugin PluginStartUp")
}
func (chainPlugin *ChainPlugin)PluginShutDown() {

}
