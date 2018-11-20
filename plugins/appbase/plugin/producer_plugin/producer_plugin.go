package producer_plugin

import (
	"fmt"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/urfave/cli"
)

//var Producer app.Plugin = app.Instance.RegisterPlugin(app.Plugin{1,"producer_plugin"}) //使用后便会被加载





type ProducerPlugin struct {
	AbstractPlugin
}

var producerPlugin *ProducerPlugin

func init() {
	producerPlugin = new(ProducerPlugin)
	producerPlugin.Plugin = producerPlugin
	producerPlugin.Name = "ProducerPlugin"
	producerPlugin.State = Registered
	App.RegisterPlugin(producerPlugin)
}

func (producerPlugin *ProducerPlugin) SetProgramOptions(options *cli.App) {
	fmt.Println("ProducerPlugin SetProgramOptions")
	fmt.Println(options.Name)

}

func (producerPlugin *ProducerPlugin) PluginInitialize() {
	fmt.Println("producerPlugin PluginInitialize")
}

//var loop = true

func (producerPlugin *ProducerPlugin) PluginStartUp() {
	fmt.Println("producerPlugin PluginStartUp")
	//go func() {
	//	for loop {
	//		fmt.Println("ProducerPlugin PluginStartUp")
	//	}
	//}()
}

func (producerPlugin *ProducerPlugin) PluginShutDown() {
	//loop = false
	fmt.Println("ProducerPlugin PluginShutDown")
}

