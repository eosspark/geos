package console_plugin

import (
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/libraries/asio"
	"github.com/eosspark/eos-go/plugins/console_plugin/console"
)

type ConsolePluginImpl struct {
	enable  bool
	jsPath  string
	datadir string
	preload []string
	exec    string
	console *console.Console

	stop    chan struct{} // Channel to wait for termination notifications
	baseUrl string
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
	config := console.Config{
		DataDir: impl.datadir,
		DocRoot: impl.jsPath,
		Client:  console.NewClient(impl.baseUrl),
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
