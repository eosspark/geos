package chain_plugin

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/urfave/cli"
)

const ChainPlug = PluginTypeName("ChainPlugin")

var chainPlugin Plugin = App().RegisterPlugin(ChainPlug, NewChainPlugin())

type ChainPlugin struct {
	AbstractPlugin
	My *ChainPluginImpl
}

func NewChainPlugin() *ChainPlugin {
	plugin := &ChainPlugin{}

	plugin.My = NewChainPluginImpl()
	plugin.My.Self = plugin
	return plugin
}

func (c *ChainPlugin) SetProgramOptions(options *[]cli.Flag) {
	*options = append(*options,
		cli.StringFlag{
			Name: "blocks-dir",
			Usage: "the location of the blocks directory (absolute path or relative to application data dir)",
		},
		//checkpoint
		cli.StringFlag{
			Name:"wasm-runtime",
			Usage:"Override default WASM runtime.",
		},
		//abi-serializer-max-time-ms
		cli.Uint64Flag{
			Name:"chain-state-db-size-mb",
			Usage:"Maximum size (in MiB) of the chain state database",
		},
		cli.Uint64Flag{
			Name:"chain-state-db-guard-size-mb",
			Usage:"Safely shut down node when free space remaining in the chain state database drops below this size (in MiB).",
		},
		cli.Uint64Flag{
			Name:"reversible-blocks-db-size-mb",
			Usage:"Maximum size (in MiB) of the reversible blocks database",
		},
		cli.Uint64Flag{
			Name:  "reversible-blocks-db-guard-size-mb",
			Usage: "Safely shut down node when free space remaining in the reverseible blocks database drops below this size (in MiB).",
		},
		cli.BoolFlag{
			Name:"contracts-console",
			Usage:"print contract's output to console",
		},



	)
}

func (c *ChainPlugin) PluginInitialize(options *cli.Context) {
	//TODO: option initialize
	c.My.Chain = chain.GetControllerInstance()
	c.My.ChainId = c.My.Chain.GetChainId()
}

func (c *ChainPlugin) PluginStartup() {
	//log.Info("Blockchain started; head block is #%d, genesis timestamp is %s",
	//	c.My.Chain.HeadBlockNum(), c.My.ChainConfig.Genesis.InitialTimestamp)
		//my->chain->head_block_num(), my->chain_config->genesis.initial_timestamp
}

func (c *ChainPlugin) PluginShutdown() {
	c.My.Chain.Close()
	log.Info("chain plugin shutdown")
}

func (c *ChainPlugin) GetReadOnlyApi() *ReadOnly {
	return NewReadOnly(c.Chain(), c.GetAbiSerializerMaxTime())
}

func (c *ChainPlugin) GetReadWriteApi() *ReadWrite {
	return NewReadWrite(c.Chain(), c.GetAbiSerializerMaxTime())
}

func (c *ChainPlugin) Chain() *chain.Controller {
	return c.My.Chain
}

func (c *ChainPlugin) GetChainId() common.ChainIdType {
	return c.My.ChainId
}

func (c *ChainPlugin) GetAbiSerializerMaxTime() common.Microseconds {
	return c.My.AbiSerializerMaxTimeMs
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
