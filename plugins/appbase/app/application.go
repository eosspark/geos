package app

import (
	. "github.com/eosspark/eos-go/plugins/appbase/app/include"
	"github.com/urfave/cli"
	. "github.com/eosspark/eos-go/exception"
	"fmt"
	"github.com/eosspark/eos-go/exception/try"
	"os"
	"path/filepath"
	"runtime"
)

// 完成初步架构设计
type applicationImpl struct {
	Version   uint64
	Options   *cli.App
	ConfigDir string
	DateDir   string
}

//var App_global *app

type application struct {
	My                 *applicationImpl
	Plugins            map[string]Plugin //< all registered plugins
	initializedPlugins []Plugin         //< stored in the order they were started running
	runningPlugins     []Plugin          //<  stored in the order they were started running
}


var appImpl = &applicationImpl{Version, cli.NewApp(), "", ""}

var App = &application{appImpl, make(map[string]Plugin), make([]Plugin, 0), make([]Plugin, 0)}

func (app *application) RegisterPlugin(plugin Plugin) Plugin {
	if p, existing := app.Plugins[plugin.GetName()]; existing {
		return p
	}
	app.Plugins[plugin.GetName()] = plugin
	return plugin
}

func setProgramOptions() {
	for _,v := range App.Plugins {
		v.SetProgramOptions(App.My.Options)
	}

	App.My.Options.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 8000,
			Usage: "listening port",
		},

		cli.StringFlag{
			Name:  "print-default-config",
			Usage: "Print default configuration template",
		},
		cli.StringFlag{
			Name:  "data-dir,d",
			Usage: "Directory containing program runtime data",
		},
		cli.StringFlag{
			Name:  "config-dir",
			Usage: "Directory containing configuration files such as config.ini",
		},
		cli.StringFlag{
			Name:  "config,c",
			Usage: "Configuration file name relative to config-dir",
		},
		cli.StringFlag{
			Name:  "logconf",
			Usage: "Logging configuration file name/path for library users",
		},
	}
	cli.HelpFlag = cli.BoolFlag{
		Name:  "help, h",
		Usage: "Print this help message and exit.",
	}
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, v",
		Usage: "Print version information.",
	}

}

func (app *application) Initialize(basicPlugin []string) bool {
	var AP []Plugin
	for i := 0; i < len(basicPlugin); i++ {
		if p := FindPlugin(basicPlugin[i]); p != nil {
			AP = append(AP, p)
		}
	}
	return App.InitializeImpl(AP)
}

func (app *application) InitializeImpl(p []Plugin) (r bool) {
	setProgramOptions()

	app.My.Options.Action = func(c *cli.Context) error {

		//help、version  will be deal with urfave.cli
		if c.String("data-dir") != "" {
			app.My.DateDir = homeDir() + c.String("data-dir")
		}
		if c.String("config-dir") != "" {
			app.My.ConfigDir = homeDir() + c.String("config-dir")
		}


		return nil
	}

	defer try.HandleReturn()
	try.Try(func() {
		for i := 0; i < len(p); i++ {
			if p[i].GetState() == Registered {
				p[i].Initialize(app.My.Options)
			}
		}
	}).Catch(func(e Exception) {
		fmt.Println(e)
		r = false
		try.Return()
	}).End()
	//need to add function--promise the plugins relative should be initialized

	return true
}

func (app *application) PluginInitialized(p Plugin) {
	app.initializedPlugins = append(app.initializedPlugins, p)
}

func (app *application) PluginStarted(p Plugin) {
	app.runningPlugins = append(app.runningPlugins,p)
}

func (app *application) StartUp() {
	for i,_ := range app.initializedPlugins {
		app.initializedPlugins[i].StartUp()
	}
}

func (app *application) ShutDown() {
	for _, v := range app.Plugins {
		v.PluginShutDown()
	}
	app.runningPlugins = app.runningPlugins[:0]
	app.initializedPlugins = app.initializedPlugins[:0]

	for k, v := range app.Plugins {
		v.PluginShutDown()
		delete(app.Plugins, k)
	}

}

func FindPlugin(name string) (plugin Plugin) {
	if v, ok := App.Plugins[name]; ok {
		return v
	}
	return nil
}

func (app *application) SetVersion(Version uint64) {
	App.My.Version = Version
}

func (app *application) FindPlugin(name string) (a *AbstractPlugin) {

	return nil
}

func GetVersion() uint64 {
	return App.My.Version
}

func (app *application) SetDefaultConfigDir() {
	App.My.ConfigDir = DefaultConfigDir()
}

func (app *application) SetDefaultDataDir() {
	App.My.DateDir = DefaultDataDir()
}

func DefaultConfigDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Application Support", "eosgo", "nodes", "config")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "eosgo", "nodes", "config")
		} else {
			return filepath.Join(home, ".clef")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Application Support", "eosgo", "nodes", "data")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "eosgo", "nodes", "data")
		} else {
			return filepath.Join(home, ".clef")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	//if usr, err := user.Current(); err == nil {
	//	return usr.HomeDir
	//}
	return ""
}
