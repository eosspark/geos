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

	Initialize(options *cli.Context)
	StartUp()
	ShutDown()

	GetName() PluginTypeName
	setName(name PluginTypeName)
	GetState() State
	setState(state State)

	bind(plugin Plugin)
}

type PluginTypeName = string

type State int

const (
	Registered  = State(iota + 1) ///< the plugin is constructed but doesn't do anything
	Initialized                   ///< the plugin has initialized any state required but is idle
	Started                       ///< the plugin is actively running
	Stopped                       ///< the plugin is no longer running
)

type AbstractPlugin struct {
	self  Plugin //must be pointer
	name  PluginTypeName
	state State
}

func (a *AbstractPlugin) ShutDown() {
	if a.state == Started {
		a.state = Stopped
		a.self.PluginShutdown()
	}
}

func (a *AbstractPlugin) Initialize(options *cli.Context) {
	if a.state == Registered {
		a.state = Initialized
		a.self.PluginInitialize(options)
		App().PluginInitialized(a.self)
	}
	Assert(a.state == Initialized, "plugin initialize failed")
}

func (a *AbstractPlugin) StartUp() {
	if a.state == Initialized {
		a.state = Started
		//static_cast<Impl*>(this)->plugin_requires([&](auto& plug){ plug.initialize(options); });
		//static_cast<Impl*>(this)->plugin_initialize(options);
		a.self.PluginStartup()
		App().PluginStarted(a.self)
	}
	Assert(a.state == Started, "plugin startup failed")
}

func (a *AbstractPlugin) GetName() PluginTypeName {
	return a.name
}

func (a *AbstractPlugin) setName(name PluginTypeName) {
	a.name = name
}

func (a *AbstractPlugin) GetState() State {
	return a.state
}

func (a *AbstractPlugin) setState(state State) {
	a.state = state
}

func (a *AbstractPlugin) bind(plugin Plugin) {
	a.self = plugin
}
