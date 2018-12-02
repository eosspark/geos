package chain_plugin

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
	"github.com/urfave/cli"
)

const ChainPlug = PluginTypeName("ChainPlugin")

var chainPlugin Plugin = App().RegisterPlugin(ChainPlug, NewChainPlugin(App().GetIoService()))

type ChainPlugin struct {
	AbstractPlugin
	my *ChainPluginImpl
}

func NewChainPlugin(io *asio.IoContext) *ChainPlugin {
	plugin := &ChainPlugin{}

	plugin.my = NewChainPluginImpl()
	plugin.my.Self = plugin
	return plugin
}

func (c *ChainPlugin) SetProgramOptions(options *[]cli.Flag) {

}

func (c *ChainPlugin) PluginInitialize(options *cli.Context) {
	log.Info("initializing chain plugin")

	c.my.ChainConfig = &chain.Config{}
}

func (c *ChainPlugin) PluginStartup() {
	log.Info("chain plugin startup")
}

func (c *ChainPlugin) PluginShutdown() {
	log.Info("chain plugin shutdown")
}

func (c *ChainPlugin) GetReadOnlyApi() *ReadOnly {
	return NewReadOnly(c.Chain(), c.GetAbiSerializerMaxTime())
}

func (c *ChainPlugin) GetReadWriteApi() *ReadWrite {
	return NewReadWrite(c.Chain(), c.GetAbiSerializerMaxTime())
}

func (c *ChainPlugin) Chain() *chain.Controller {
	return c.my.Chain
}

func (c *ChainPlugin) GetChainId() common.ChainIdType {
	return *c.my.ChainId
}

func (c *ChainPlugin) GetAbiSerializerMaxTime() common.Microseconds {
	return c.my.AbiSerializerMaxTimeMs
}

func (c *ChainPlugin) HandleGuardException(e exception.GuardExceptions) {
	//TODO
	//log_guard_exception(e);
	//
	//// quit the app
	//app().quit();
}

//func (chain *ChainPlugin) GetAbiSerializerMaxTime() common.Microseconds {
//	return chain.AbiSerializerMaxTimeMs
//}
//
//func (chain *ChainPlugin) Init() {
//	chain.AbiSerializerMaxTimeMs = 1000 //TODO tmp value
//}
