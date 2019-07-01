package chain_plugin

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/libraries/asio"
	"github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/appbase/app/include"
	. "github.com/eosspark/eos-go/plugins/chain_interface"
)

type ChainPluginImpl struct {
	BlockDir asio.Path
	Readonly bool
	//flat_map<uint32_t,block_id_type> loaded_checkpoints;

	ForkDB      *chain.ForkDatabase
	BlockLogger *chain.BlockLog

	ChainConfig *chain.Config
	Chain       *chain.Controller
	ChainId     common.ChainIdType

	//fc::optional<vm_type>            wasm_runtime;
	AbiSerializerMaxTimeMs common.Microseconds
	//fc::optional<bfs::path>          snapshot_path;

	// retained references to channels for easy publication
	PreAcceptedBlockChannel     include.Channel
	AcceptedBlockHeaderChannel  include.Channel
	AcceptedBlockChannel        include.Channel
	IrreversibleBlockChannel    include.Channel
	AcceptedTransactionChannel  include.Channel
	AppliedTransactionChannel   include.Channel
	AcceptedConfirmationChannel include.Channel
	IncomingBlockChannel        include.Channel

	// retained references to methods for easy calling
	IncomingBlockSyncMethod        *include.Method
	IncomingTransactionAsyncMethod *include.Method

	// method provider handles
	GetBlockByNumberProvider               *include.Method
	GetBlockByIdProvider                   *include.Method
	GetHeadBlockIdProvider                 *include.Method
	GetLastIrreversibleBlockNumberProvider *include.Method

	// scoped connections for chain controller
}

func NewChainPluginImpl() *ChainPluginImpl {
	return &ChainPluginImpl{
		IncomingBlockSyncMethod:        app.App().GetMethod(BlockSync),
		IncomingTransactionAsyncMethod: app.App().GetMethod(TransactionAsync),
	}
}
