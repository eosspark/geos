package chain_plugin

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/plugins/appbase/asio"
)

type ChainPluginImpl struct {
	BlockDir asio.Path
	Readonly bool
	//flat_map<uint32_t,block_id_type> loaded_checkpoints;

	//fc::optional<fork_database>      fork_db;
	//fc::optional<block_log>          block_logger;

	ChainConfig *chain.Config       `eos:"optional"`
	Chain       *chain.Controller   `eos:"optional"`
	ChainId     common.ChainIdType `eos:"optional"`

	//fc::optional<vm_type>            wasm_runtime;
	AbiSerializerMaxTimeMs common.Microseconds
	//fc::optional<bfs::path>          snapshot_path;

	// retained references to channels for easy publication

	// retained references to methods for easy calling

	// method provider handles

	// scoped connections for chain controller
	Self *ChainPlugin
}

func NewChainPluginImpl() *ChainPluginImpl {
	return new(ChainPluginImpl)
}

func (c *ChainPluginImpl) GetInfo() *InfoResp {
	return &InfoResp{
		ServerVersion:            "0f6695cb",
		ChainID:                  common.BlockIdNil(),
		HeadBlockNum:             17673,
		LastIrreversibleBlockNum: 17672,
		LastIrreversibleBlockID:  common.BlockIdNil(),
		HeadBlockID:              common.BlockIdNil(),
		HeadBlockTime:            common.Now(),
		HeadBlockProducer:        common.AccountName(common.N("eosio")),
		VirtualBlockCPULimit:     200000000,
		VirtualBlockNetLimit:     1048576000,
		BlockCPULimit:            199900,
		BlockNetLimit:            1048576,
		ServerVersionString:      "TODO walker",
	}

}
