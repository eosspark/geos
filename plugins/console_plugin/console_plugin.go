package console_plugin

import (
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/appbase/asio"

	"github.com/eosspark/eos-go/common"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/urfave/cli"
	"strings"
)

type ConsolePlugin struct {
	AbstractPlugin
	//ConfirmedBlock Signal //TODO signal ConfirmedBlock
	my *ConsolePluginImpl
}

func init() {
	plug := NewConsolePlugin(App().GetIoService())
	plug.Plugin = plug
	plug.Name = PluginName(ConsolePlug)
	plug.State = State(Registered)
	App().RegisterPlugin(plug)
}

//TODO: io from appbase
func NewConsolePlugin(io *asio.IoContext) *ConsolePlugin {
	p := new(ConsolePlugin)

	p.my = NewConsolePluginImpl(io)
	p.my.Self = p

	return p
}

func (cp *ConsolePlugin) SetProgramOptions(options *[]cli.Flag) {
	*options = append(*options,
		cli.BoolFlag{
			Name:  "console",
			Usage: "Start an interactive JavaScript environment",
		},
		cli.StringFlag{ // ATM the url is left to the user and deployment to
			Name:  "jspath",
			Usage: "JavaScript root path for `loadScript`",
			Value: ".",
		},
		cli.StringFlag{
			Name:  "console-datadir",
			Usage: "Data directory for the databases and keystore",
			Value: "./console_history",
		},
		cli.StringFlag{
			Name:  "exec",
			Usage: "Execute JavaScript statement",
		},
		cli.StringFlag{
			Name:  "preload",
			Usage: "Comma separated list of JavaScript files to preload into the console",
		})
}

func (cp *ConsolePlugin) PluginInitialize(c *cli.Context) {
	//Try(func() {
	cp.my.datadir = c.String("console-datadir")
	cp.my.exec = c.String("exec")

	cp.my.jsPath = c.String("jspath")
	preloadJSString := c.String("preload")
	// MakeConsolePreloads retrieves the absolute paths for the console JavaScript
	// scripts to preload before starting.
	if preloadJSString != "" {
		for _, file := range strings.Split(preloadJSString, ",") {
			cp.my.preload = append(cp.my.preload, common.AbsolutePath(cp.my.jsPath, strings.TrimSpace(file)))
		}
	} else {
		cp.my.preload = nil
	}

	cp.my.enable = c.Bool("console")
	if cp.my.enable {
		err := cp.my.localConsole()
		if err != nil {
			FcThrow("Failed to start the JavaScript console : %s", err)
		}
	}

	//}).FcLogAndRethrow().End()
}

func (cp *ConsolePlugin) PluginStartup() {
	if cp.my.enable {
		// Otherwise print the welcome screen and enter interactive mode
		cp.my.console.Welcome()

		cp.my.console.Interactive()
		cp.my.log.Info("interactive start")
	}
}

func (cp *ConsolePlugin) PluginShutdown() {
	cp.my.console.Stop(false) //TODO
}
