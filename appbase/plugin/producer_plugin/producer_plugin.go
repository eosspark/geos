package producer_plugin

import (
	. "github.com/eosspark/eos-go/appbase/app/include"
	"fmt"

	"github.com/eosspark/eos-go/appbase/app"
)

//var Producer app.Plugin = app.Instance.RegisterPlugin(app.Plugin{1,"producer_plugin"}) //使用后便会被加载

type ProducerPlugin struct {
	AbstractPlugin
}

var producerPlugin *ProducerPlugin

func init() {
	producerPlugin = new(ProducerPlugin)
	producerPlugin.Name = "ProducerPlugin"
	producerPlugin.State = Registered
	app.App.RegisterPlugin(producerPlugin)
}

func (producerPlugin *ProducerPlugin) SetProgramOptions() {
	fmt.Println("ProducerPlugin SetProgramOptions")

}

func (producerPlugin *ProducerPlugin) PluginInitialize() {
	producerPlugin.State = Initialized
}

var loop = true

func (producerPlugin *ProducerPlugin) PluginStartUp() {
	go func() {
		for loop {
			fmt.Println("ProducerPlugin PluginStartUp")
		}
	}()
}

func (producerPlugin *ProducerPlugin) PluginShutDown() {
	loop = false
	fmt.Println("ProducerPlugin PluginShutDown")
}
func (producerPlugin *ProducerPlugin) GetName() string {
	return producerPlugin.Name
}

func (producerPlugin *ProducerPlugin) GetState() State {
	return producerPlugin.State
}
