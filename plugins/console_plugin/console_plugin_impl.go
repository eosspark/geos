package console_plugin

import (
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"github.com/eosspark/eos-go/plugins/console_plugin/console"
	"github.com/eosspark/eos-go/plugins/http_plugin/rpc"
	"net"
)

type ConsolePluginImpl struct {
	enable  bool
	jsPath  string
	datadir string
	preload []string
	exec    string
	console *console.Console

	rpcAPIs       []rpc.API   // List of APIs currently provided by the node
	inprocHandler *rpc.Server // In-process RPC request handler to process the API requests

	ipcEndpoint string       // IPC endpoint to listen at (empty = IPC disabled)
	ipcListener net.Listener // IPC RPC listener socket to serve API requests
	ipcHandler  *rpc.Server  // IPC RPC request handler to process the API requests

	httpEndpoint  string       // HTTP endpoint (interface + port) to listen at (empty = HTTP disabled)
	httpWhitelist []string     // HTTP RPC modules to allow through this endpoint
	httpListener  net.Listener // HTTP RPC listener socket to server API requests
	httpHandler   *rpc.Server  // HTTP RPC request handler to process the API requests

	config *Config

	wsEndpoint string       // Websocket endpoint (interface + port) to listen at (empty = websocket disabled)
	wsListener net.Listener // Websocket RPC listener socket to server API requests
	wsHandler  *rpc.Server  // Websocket RPC request handler to process the API requests

	stop chan struct{} // Channel to wait for termination notifications

	Self *ConsolePlugin
	log  log.Logger
}

func NewConsolePluginImpl(io *asio.IoContext) *ConsolePluginImpl {
	impl := new(ConsolePluginImpl)
	impl.log = log.New("console")
	impl.log.SetHandler(log.TerminalHandler)

	impl.config = &DefaultConfig

	//impl.ipcEndpoint=      impl.config.IPCEndpoint()
	impl.httpEndpoint = impl.config.HTTPEndpoint()
	//impl.wsEndpoint=   impl.config.WSEndpoint()

	return impl
}

func (impl *ConsolePluginImpl) localConsole() error {

	client, err := rpc.Dial("http://127.0.0.1:8888")
	if err != nil {
		impl.log.Error("dial client is wrong: %s", err)
	}

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
