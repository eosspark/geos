package console_plugin

import (
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"github.com/eosspark/eos-go/plugins/console_plugin/console"
	"github.com/eosspark/eos-go/plugins/console_plugin/console/js/eosapi"
	"github.com/eosspark/eos-go/plugins/console_plugin/console/rpc"
)

type ConsolePluginImpl struct {
	enable  bool
	jsPath  string
	datadir string
	preload []string
	exec    string
	console *console.Console
	Self    *ConsolePlugin
	log     log.Logger
}

func NewConsolePluginImpl(io *asio.IoContext) *ConsolePluginImpl {
	impl := new(ConsolePluginImpl)
	impl.log = log.New("console")
	impl.log.SetHandler(log.TerminalHandler)
	return impl
}

func (impl *ConsolePluginImpl) localConsole() error {
	// Register all the APIs exposed by the services
	handler := rpc.NewServer()
	apis := apis()
	for _, api := range apis {
		if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
			log.Error(err.Error())
			panic(err)
		}
		log.Debug("InProc registered :  namespace =%s", api.Namespace)
	}

	client := rpc.DialInProc(handler)

	//config := console.Config{
	//	DataDir: "console_history",//mpl.datadir,
	//	DocRoot: "test",//impl.jsPath,
	//	Client:  client,
	//	Preload: nil,//impl.preload,
	//}
	config := console.Config{
		DataDir: impl.datadir,
		DocRoot: impl.jsPath,
		Client:  client,
		Preload: impl.preload,
	}

	console, err := console.New(config)
	if err != nil {
		log.Error("Failed to start the JavaScript console: %v", err)
	}

	//// If only a short execution was requested, evaluate and return
	if impl.exec != "" {
		console.Evaluate(impl.exec)
		return nil
	}
	impl.console = console

	return nil

}

// apis returns the collection of RPC descriptors this node offers.
func apis() []rpc.API {
	return []rpc.API{
		{
			Namespace: "api",
			Version:   "1.0",
			Service:   eosapi.NewEosApi(),
		},
		//{
		//	Namespace: "produce",
		//	Version:   "1.0",
		//	Service:   eosapi.NewEosApi(),
		//},
		//{
		//	Namespace: "net",
		//	Version:   "1.0",
		//	Service:   net_plugin.NewNetPlugin(),
		//},
		//{
		//	Namespace: "wallet",
		//	Version:   "1.0",
		//	Service:   wallet_plugin.NewWalletPlugin(),
		//},
	}
}
