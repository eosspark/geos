package app

import (
	"fmt"
	. "github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	. "github.com/eosspark/eos-go/plugins/appbase/app/include"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	. "github.com/eosspark/eos-go/plugins/chain_interface"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

type ApplicationImpl struct {
	version uint64
	Options *cli.App

	DateDir     asio.Path
	ConfigDir   asio.Path
	LoggingConf asio.Path
}

//var App_global *app

type Application struct {
	plugins            map[PluginTypeName]Plugin //< all registered plugins
	initializedPlugins []Plugin                  //< stored in the order they were started running
	runningPlugins     []Plugin                  //<  stored in the order they were started running

	//methods
	channels map[ChannelsType]*Channel
	methods  map[MethodsType]*Method
	iosv     *asio.IoContext
	my       *ApplicationImpl
}

func newApplication() *Application {
	iosv := asio.NewIoContext()
	appImpl := &ApplicationImpl{
		Version,
		cli.NewApp(),
		"data-dir",
		"config-dir",
		"logging.json"}

	app := &Application{
		plugins:            make(map[PluginTypeName]Plugin),
		initializedPlugins: make([]Plugin, 0),
		runningPlugins:     make([]Plugin, 0),
		channels:           make(map[ChannelsType]*Channel),
		methods:            make(map[MethodsType]*Method),
		iosv:               iosv,
		my:                 appImpl,
	}

	return app
}

var app *Application

func App() *Application {
	if app == nil {
		app = newApplication()
	}
	return app
}

func (app *Application) RegisterPlugin(typename PluginTypeName, plugin Plugin) Plugin {
	if p, existing := app.plugins[typename]; existing {
		return p
	}

	plugin.setName(typename)
	plugin.setState(Registered)
	plugin.bind(plugin)

	app.plugins[typename] = plugin

	return plugin
}

func (app *Application) setProgramOptions() {
	for _, v := range app.plugins {
		v.SetProgramOptions(&app.my.Options.Flags)
	}
	app.my.Options.Flags = append(app.my.Options.Flags,
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
	)
	cli.HelpFlag = cli.BoolFlag{
		Name:  "help, h",
		Usage: "Print this help message and exit.",
	}
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, v",
		Usage: "Print version information.",
	}

}

func (app *Application) Initialize(basicPlugin []PluginTypeName) bool {
	var AP []Plugin
	for i := 0; i < len(basicPlugin); i++ {
		if p := app.FindPlugin(basicPlugin[i]); p != nil {
			AP = append(AP, p)
		}
	}
	return app.InitializeImpl(AP)
}

func (app *Application) InitializeImpl(p []Plugin) bool {
	returning, r := false, false
	try.Try(func() {
		app.setProgramOptions()

		app.my.Options.Action = func(c *cli.Context) error {
			for i := 0; i < len(p); i++ {
				if p[i].GetState() == Registered {
					p[i].Initialize(c)
				}
			}

			//help、version  will be deal with urfave.cli
			if c.String("data-dir") != "" {
				app.my.DateDir = homeDir() + c.String("data-dir")
			}
			if c.String("config-dir") != "" {
				app.my.ConfigDir = homeDir() + c.String("config-dir")
			}

			return nil
		}

		app.my.Options.Run(os.Args)

	}).Catch(func(e Exception) {
		fmt.Println(e)
		returning, r = true, false
		return
	}).End()

	if returning {
		return r
	}
	//need to add function--promise the plugins relative should be initialized

	return true
}

func (app *Application) GetChannel(channelType ChannelsType) *Channel {
	if v, ok := app.channels[channelType]; ok {
		return v
	} else {
		channel := NewChannel(app.iosv)
		app.channels[channelType] = channel
		return channel
	}
}

func (app *Application) GetMethod(methodsType MethodsType) *Method {
	if v, ok := app.methods[methodsType]; ok {
		return v
	} else {
		method := NewMethod()
		app.methods[methodsType] = method
		return method
	}
}

func (app *Application) GetIoService() *asio.IoContext {
	return app.iosv
}

func (app *Application) PluginInitialized(p Plugin) {
	app.initializedPlugins = append(app.initializedPlugins, p)
}

func (app *Application) PluginStarted(p Plugin) {
	app.runningPlugins = append(app.runningPlugins, p)
}

func (app *Application) StartUp() {
	for i := range app.initializedPlugins {
		app.initializedPlugins[i].StartUp()
	}
}

func (app *Application) ShutDown() {
	for _, v := range app.plugins {
		v.ShutDown()
	}
	app.runningPlugins = app.runningPlugins[:0]
	app.initializedPlugins = app.initializedPlugins[:0]

	for k := range app.plugins {
		delete(app.plugins, k)
	}
	app.iosv = nil
	app = nil
}

func (app *Application) Quit() {
	app.iosv.Stop()
}

func (app *Application) Exec() {
	sigint := asio.NewSignalSet(app.iosv, syscall.SIGINT)
	sigint.AsyncWait(func(err error) {
		app.Quit()
		sigint.Cancel()
	})
	sigterm := asio.NewSignalSet(app.iosv, syscall.SIGTERM)
	sigterm.AsyncWait(func(err error) {
		app.Quit()
		sigterm.Cancel()
	})
	sigpipe := asio.NewSignalSet(app.iosv, syscall.SIGPIPE)
	sigpipe.AsyncWait(func(err error) {
		app.Quit()
		sigpipe.Cancel()
	})

	app.iosv.Run()

	app.ShutDown()
}

func (app *Application) FindPlugin(name PluginTypeName) (plugin Plugin) {
	if v, ok := app.plugins[name]; ok {
		return v
	}
	return nil
}

func (app *Application) GetPlugin(name PluginTypeName) (plugin Plugin) {
	p := app.FindPlugin(name)
	if p == nil {
		fmt.Println("unable to find plugin") //need to fix
	}
	return p
}

func (app *Application) GetVersion() uint64 {
	return app.my.version
}

func (app *Application) SetVersion(version uint64) {
	app.my.version = version
}

func (app *Application) SetDefaultConfigDir() {
	app.my.ConfigDir = DefaultConfigDir()
}

func (app *Application) SetDefaultDataDir() {
	app.my.DateDir = DefaultDataDir()
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

func (app *Application) DataDir() asio.Path {
	return app.my.DateDir
}
