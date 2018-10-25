package include

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/eosspark/eos-go/plugins/appbase/app"
)

/** these notifications get called from the plugin when their state changes so that
 * the application can call shutdown in the reverse order.
 */
//Register(abstract_plugin) error

type Plugin interface {
	SetProgramOptions()
	PluginInitialize()
	PluginStartUp()
	PluginShutDown()


	GetName() string
	GetState() State
	Initialize(options *cli.App)
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

func (a *AbstractPlugin) Initialize(options *cli.App) {


	if a.State == Registered {
		a.PluginInitialize()

		app.App.PluginInitialized(a.Plugin)
	}
}

func (a *AbstractPlugin) GetName() string {
	return a.Name
}

func (a *AbstractPlugin) GetState() State {
	return a.State
}


