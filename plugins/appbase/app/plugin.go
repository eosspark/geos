package app

import (
	"github.com/urfave/cli"

	. "github.com/eosspark/eos-go/exception"
)

/** these notifications get called from the plugin when their state changes so that
 * the application can call shutdown in the reverse order.
 */
//Register(abstract_plugin) error

type Plugin interface {
	SetProgramOptions(options *cli.App)
	PluginInitialize()
	PluginStartUp()
	PluginShutDown()

	GetName() string
	GetState() State
	Initialize(options *cli.App)
	StartUp()
}

type State int

const (
	Registered  = State(iota + 1) ///< the plugin is constructed but doesn't do anything
	Initialized                   ///< the plugin has initialized any state required but is idle
	Started                       ///< the plugin is actively running
	Stopped                       ///< the plugin is no longer running
)

type AbstractPlugin struct {
	Plugin
	Name  string
	State State
}



func (a *AbstractPlugin) ShutDown() {

}



func (a *AbstractPlugin) Initialize(options *cli.App) {
	if a.State == Registered {
		a.State = Initialized
		a.PluginInitialize()
		App.PluginInitialized(a)
	}
}

func (a *AbstractPlugin) StartUp() {
	if a.State == Initialized {
		a.State = Started
		//为了确保每个plugin依赖的其他plugin也保证initialize
		//static_cast<Impl*>(this)->plugin_requires([&](auto& plug){ plug.initialize(options); });
		//static_cast<Impl*>(this)->plugin_initialize(options);
		a.PluginStartUp()
		EosAssert(false, &ExtractGenesisStateException{}, "error")
		App.PluginStarted(a)
	}
}


func (a *AbstractPlugin) GetName() string {
	return a.Name
}

func (a *AbstractPlugin) GetState() State {
	return a.State
}





