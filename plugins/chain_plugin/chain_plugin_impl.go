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
	ChainId     *common.ChainIdType `eos:"optional"`

	//fc::optional<vm_type>            wasm_runtime;
	AbiSerializerMaxTimeMs common.Microseconds
	//fc::optional<bfs::path>          snapshot_path;

	// retained references to channels for easy publication

	// retained references to methods for easy calling

	// method provider handles

	// scoped connections for chain controller
}

func NewChainPluginImpl() *ChainPluginImpl {
	return new(ChainPluginImpl)
}
