package app

import (
	"github.com/urfave/cli"

		. "github.com/eosspark/eos-go/exception/try"
)

/** these notifications get called from the plugin when their state changes so that
 * the application can call shutdown in the reverse order.
 */
//Register(abstract_plugin) error

type Plugin interface {
	SetProgramOptions(options *[]cli.Flag)
	PluginInitialize(*cli.Context)
	PluginStartup()
	PluginShutdown()

	GetName() PluginName
	GetState() State
	Initialize(options *cli.Context)
	StartUp()
}

type PluginName string

const (
	ProducerPlug = PluginName("ProducerPlugin")
	ChainPlug    = PluginName("ChainPlugin")
	NetPlug      = PluginName("NetPlugin")
	HttpPlug     = PluginName("HttpPlugin")
)

type State int

const (
	Registered  = State(iota + 1) ///< the plugin is constructed but doesn't do anything
	Initialized                   ///< the plugin has initialized any state required but is idle
	Started                       ///< the plugin is actively running
	Stopped                       ///< the plugin is no longer running
)

type AbstractPlugin struct {
	Plugin
	Name  PluginName
	State State
}

func (a *AbstractPlugin) ShutDown() {

}

func (a *AbstractPlugin) Initialize(options *cli.Context) {
	if a.State == Registered {
		a.State = Initialized
		a.PluginInitialize(options)
		App().PluginInitialized(a.Plugin)
	}
}

func (a *AbstractPlugin) StartUp() {
	if a.State == Initialized {
		a.State = Started
		//为了确保每个plugin依赖的其他plugin也保证initialize
		//static_cast<Impl*>(this)->plugin_requires([&](auto& plug){ plug.initialize(options); });
		//static_cast<Impl*>(this)->plugin_initialize(options);
		a.PluginStartup()
		App().PluginStarted(a.Plugin)
	}

	Assert(a.State == State(Started), "plugin startup failed")
}

func (a *AbstractPlugin) GetName() PluginName {
	return a.Name
}

func (a *AbstractPlugin) GetState() State {
	return a.State
}
