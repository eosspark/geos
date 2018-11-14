package chain_plugin

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/exception"
)

//var IsActive bool = false
//var chainPlugin *ChainPlugin

type ChainPlugin struct {
	my *ChainPluginImpl
}

func NewChainPlugin() *ChainPlugin {
	c := NewChainPlugin()
	c.my = NewChainPluginImpl()
	return c
}

func (c *ChainPlugin) SetProgramOptions(options *cli.App) {

}

func (c *ChainPlugin) PluginInitialize(options *cli.App) {

}

func (c *ChainPlugin) PluginStartup(options *cli.App) {

}

func (c *ChainPlugin) PluginShutdown(options *cli.App) {

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
