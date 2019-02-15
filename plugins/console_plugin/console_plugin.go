package console_plugin

import (
	"github.com/eosspark/eos-go/common"
	. "github.com/eosspark/eos-go/exception/try"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"github.com/eosspark/eos-go/plugins/console_plugin/console"
	"github.com/urfave/cli"
	"strings"
)

const ConsolePlug = PluginTypeName("ConsolePlugin")

var consolePlugin Plugin = App().RegisterPlugin(ConsolePlug, NewConsolePlugin(App().GetIoService()))

type ConsolePlugin struct {
	AbstractPlugin
	my *ConsolePluginImpl
}

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
		cli.StringFlag{
			Name:  "attach",
			Usage: "Start an interactive JavaScript environment (connect to node)",
			Value: common.HttpEndPoint,
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
		},

		//// RPC settings
		//cli.BoolFlag{
		//	Name:  "rpc",
		//	Usage: "Enable the HTTP-RPC server",
		//},
		//cli.StringFlag{
		//	Name:  "rpcaddr",
		//	Usage: "HTTP-RPC server listening interface",
		//	//Value: node.DefaultHTTPHost,
		//},
		//cli.IntFlag{
		//	Name:  "rpcport",
		//	Usage: "HTTP-RPC server listening port",
		//	//Value: node.DefaultHTTPPort,
		//},
		//cli.StringFlag{
		//	Name:  "rpccorsdomain",
		//	Usage: "Comma separated list of domains from which to accept cross origin requests (browser enforced)",
		//	Value: "",
		//},
		//cli.StringFlag{
		//	Name:  "rpcvhosts",
		//	Usage: "Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard.",
		//	//Value: strings.Join(node.DefaultConfig.HTTPVirtualHosts, ","),
		//},
		//cli.StringFlag{
		//	Name:  "rpcapi",
		//	Usage: "API's offered over the HTTP-RPC interface",
		//	Value: "",
		//},
		//cli.BoolFlag{
		//	Name:  "ipcdisable",
		//	Usage: "Disable the IPC-RPC server",
		//},
		//DirectoryFlag{//IPCPathFlag =
		//	Name:  "ipcpath",
		//	Usage: "Filename for IPC socket/pipe within the datadir (explicit paths escape it)",
		//},
		//cli.BoolFlag{
		//	Name:  "ws",
		//	Usage: "Enable the WS-RPC server",
		//},
		//cli.StringFlag{
		//	Name:  "wsaddr",
		//	Usage: "WS-RPC server listening interface",
		//	//Value: node.DefaultWSHost,
		//},
		//cli.IntFlag{
		//	Name:  "wsport",
		//	Usage: "WS-RPC server listening port",
		//	//Value: node.DefaultWSPort,
		//},
		//cli.StringFlag{
		//	Name:  "wsapi",
		//	Usage: "API's offered over the WS-RPC interface",
		//	Value: "",
		//},
		//cli.StringFlag{
		//	Name:  "wsorigins",
		//	Usage: "Origins from which to accept websockets requests",
		//	Value: "",
		//}
	)
}

func (cp *ConsolePlugin) PluginInitialize(c *cli.Context) {
	Try(func() {
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
		cp.my.baseUrl = c.String("attach")
		console.BaseUrl = cp.my.baseUrl

		cp.my.enable = c.Bool("console")
		if cp.my.enable {
			err := cp.my.localConsole()
			if err != nil {
				FcThrow("Failed to start the JavaScript console : %s", err)
			}
		}
	}).FcLogAndRethrow().End()
}

func (cp *ConsolePlugin) PluginStartup() {
	if cp.my.enable {
		// Otherwise print the welcome screen and enter interactive mode
		cp.my.console.Welcome()

		cp.my.console.Interactive()
	}
}

func (cp *ConsolePlugin) PluginShutdown() {
	cp.my.console.Stop(false)

}
