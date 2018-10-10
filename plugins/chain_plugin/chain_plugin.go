package chain_plugin

import "github.com/eosspark/eos-go/common"

var IsActive bool = false
var chain *ChainPlugin

type ChainPlugin struct {
	AbiSerializerMaxTimeMs common.Microseconds
}

func GetInstance() *ChainPlugin {
	if !IsActive {
		chain = newChainPlugin()
		IsActive = true
	}
	chain.Init()
	return chain
}

func newChainPlugin() *ChainPlugin {
	chain = &ChainPlugin{}
	return chain
}

func (chain *ChainPlugin) GetAbiSerializerMaxTime() common.Microseconds {
	return chain.AbiSerializerMaxTimeMs
}

func (chain *ChainPlugin) Init() {
	chain.AbiSerializerMaxTimeMs = 1000 //TODO tmp value
}
