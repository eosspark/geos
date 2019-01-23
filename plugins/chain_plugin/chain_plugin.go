package chain_plugin

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	. "github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/plugins/chain_interface"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"strings"
)

const ChainPlug = PluginTypeName("ChainPlugin")

var chainPlugin Plugin = App().RegisterPlugin(ChainPlug, NewChainPlugin())

type ChainPlugin struct {
	AbstractPlugin
	my *ChainPluginImpl
}

func NewChainPlugin() *ChainPlugin {
	plugin := &ChainPlugin{}
	plugin.my = NewChainPluginImpl()
	return plugin
}

func (c *ChainPlugin) SetProgramOptions(options *[]cli.Flag) {
	//cfg
	*options = append(*options,
		cli.StringFlag{
			Name:  "blocks-dir",
			Usage: "the location of the blocks directory (absolute path or relative to application data dir)",
			Value: App().DataDir() + "/blocks",
		},
		cli.StringSliceFlag{
			Name:  "checkpoint",
			Usage: "Pairs of [BLOCK_NUM,BLOCK_ID] that should be enforced as checkpoints.",
		},

		cli.StringFlag{
			Name:  "wasm-runtime",
			Usage: "Override default WASM runtime.",
		},
		cli.UintFlag{
			Name:  "abi-serializer-max-time-ms",
			Usage: "Override default maximum ABI serialization time allowed in ms",
		},
		//TODO UNUSED
		//cli.Uint64Flag{
		//	Name:  "chain-state-db-size-mb",
		//	Usage: "Maximum size (in MiB) of the chain state database",
		//},
		//cli.Uint64Flag{
		//	Name:  "chain-state-db-guard-size-mb",
		//	Usage: "Safely shut down node when free space remaining in the chain state database drops below this size (in MiB).",
		//},
		//cli.Uint64Flag{
		//	Name:  "reversible-blocks-db-size-mb",
		//	Usage: "Maximum size (in MiB) of the reversible blocks database",
		//},
		//cli.Uint64Flag{
		//	Name:  "reversible-blocks-db-guard-size-mb",
		//	Usage: "Safely shut down node when free space remaining in the reverseible blocks database drops below this size (in MiB).",
		//},
		cli.BoolFlag{
			Name:  "contracts-console",
			Usage: "print contract's output to console",
		},
		cli.StringSliceFlag{
			Name:  "actor-whitelist",
			Usage: "Account added to actor whitelist (may specify multiple times)",
		},
		cli.StringSliceFlag{
			Name:  "actor-blacklist",
			Usage: "Account added to actor blacklist (may specify multiple times)",
		},
		cli.StringSliceFlag{
			Name:  "contract-whitelist",
			Usage: "Contract account added to contract whitelist (may specify multiple times)",
		},
		cli.StringSliceFlag{
			Name:  "contract-blacklist",
			Usage: "Contract account added to contract blacklist (may specify multiple times)",
		},
		cli.StringSliceFlag{
			Name:  "action-blacklist",
			Usage: "Action (in the form code::action) added to action blacklist (may specify multiple times)",
		},
		cli.StringSliceFlag{
			Name:  "Key-blacklist",
			Usage: "Public key added to blacklist of keys that should not be included in authorities (may specify multiple times)",
		},
		cli.StringFlag{
			Name: "read-mode",
			Usage: "Database read mode (\"speculative\", \"head\", or \"read-only\").\n" + // or \"irreversible\").\n"
				"In \"speculative\" mode database contains changes done up to the head block plus changes made by transactions not yet included to the blockchain.\n" +
				"In \"head\" mode database contains changes done up to the current head block.\n" +
				"In \"read-only\" mode database contains incoming block changes but no speculative transaction processing.\n",
		},
		cli.StringFlag{
			Name: "validation-mode",
			Usage: "Chain validation mode (\"full\" or \"light\").\n" +
				"In \"full\" mode all incoming blocks will be fully validated.\n" +
				"In \"light\" mode all incoming blocks headers will be fully validated; transactions in those validated blocks will be trusted \n",
		},
		cli.BoolFlag{
			Name:  "disable-ram-billing-notify-checks",
			Usage: "Disable the check which subjectively fails a transaction if a contract bills more RAM to another account within the context of a notification handler (i.e. when the receiver is not the code of the action).",
		},
		cli.StringSliceFlag{
			Name:  "trusted-producer",
			Usage: "Indicate a producer whose blocks headers signed by it will be fully validated, but transactions in those validated blocks will be trusted.",
		},
	)

	//cli
	*options = append(*options,
		cli.StringFlag{
			Name:  "genesis-json",
			Usage: "File to read Genesis State from",
		},
		cli.StringFlag{
			Name:  "genesis-timestamp",
			Usage: "override the initial timestamp in the Genesis State file",
		},
		cli.BoolFlag{
			Name:  "print-genesis-json",
			Usage: "extract genesis_state from blocks.log as JSON, print to console, and exit",
		},
		cli.StringFlag{
			Name:  "extract-genesis-json",
			Usage: "extract genesis_state from blocks.log as JSON, write into specified file, and exit",
		},
		cli.BoolFlag{
			Name:  "fix-reversible-blocks",
			Usage: "recovers reversible block database if that database is in a bad state",
		},
		cli.BoolFlag{
			Name:  "force-all-checks",
			Usage: "do not skip any checks that can be skipped while replaying irreversible blocks",
		},
		cli.BoolFlag{
			Name:  "disable-replay-opts",
			Usage: "disable optimizations that specifically target replay",
		},
		cli.BoolFlag{
			Name:  "replay-blockchain",
			Usage: "clear chain state database and replay all blocks",
		},
		cli.BoolFlag{
			Name:  "hard-replay-blockchain",
			Usage: "clear chain state database, recover as many blocks as possible from the block log, and then replay those blocks",
		},
		cli.BoolFlag{
			Name:  "delete-all-blocks",
			Usage: "clear chain state database and block log",
		},
		cli.UintFlag{
			Name:  "truncate-at-block",
			Usage: "stop hard replay / block log recovery at this block number (if set to non-zero number)",
		},
		cli.StringFlag{
			Name:  "import-reversible-blocks",
			Usage: "replace reversible block database with blocks imported from specified file and then exit",
		},
		cli.StringFlag{
			Name:  "export-reversible-blocks",
			Usage: "export reversible block database in portable format into specified file and then exit",
		},
		cli.StringFlag{
			Name:  "snapshot",
			Usage: "File to read Snapshot State from",
		},
	)
}

func ClearDirectoryContents(p string) {
	if fileStat, err := os.Stat(p); err != nil || !fileStat.IsDir() {
		return
	}

	filepath.Walk(p, func(subPath string, info os.FileInfo, err error) error {
		return os.RemoveAll(subPath)
	})
}

func (c *ChainPlugin) PluginInitialize(options *cli.Context) {
	log.Info("initializing chain plugin")

	Try(func() {
		types.NewGenesisState()
	}).Catch(func(e Exception) {
		log.Error("EOSIO_ROOT_KEY ('%s') is invalid. Recompile with a valid public key.", types.EosioRootKey)
		Throw(e)
	}).End()

	c.my.ChainConfig = chain.NewConfig()

	for _, actor := range options.StringSlice("actor-whitelist") {
		c.my.ChainConfig.ActorWhitelist.Add(N(actor))
	}
	for _, actor := range options.StringSlice("actor-blacklist") {
		c.my.ChainConfig.ActorBlacklist.Add(N(actor))
	}
	for _, constract := range options.StringSlice("contract-whitelist") {
		c.my.ChainConfig.ContractWhitelist.Add(N(constract))
	}
	for _, constract := range options.StringSlice("contract-blacklist") {
		c.my.ChainConfig.ContractBlacklist.Add(N(constract))
	}

	for _, producer := range options.StringSlice("trusted-producer") {
		c.my.ChainConfig.TrustedProducers.Add(N(producer))
	}

	for _, action := range options.StringSlice("action-blacklist") {
		pos := strings.Index(action, "::")
		EosAssert(pos != -1, &PluginConfigException{}, "Invalid entry in action-blacklist: '%s'", action)
		code := N(action[0:pos])
		act := N(action[pos+2:])
		c.my.ChainConfig.ActionBlacklist.Add(MakePair(code, act))
	}

	for _, keyStr := range options.StringSlice("key-blacklist") {
		key, err := ecc.NewPublicKey(keyStr)
		if err != nil {
			Throw(err)
		}
		c.my.ChainConfig.KeyBlacklist.Add(key)
	}

	c.my.BlockDir = options.String("blocks-dir")

	for _, cp := range options.StringSlice("checkpoint") {
		//TODO handle checkpoint
		log.Debug(cp)
	}

	//TODO wasm-runtime: just use wasmgo?

	if ms := options.Uint("abi-serializer-max-time-ms"); ms > 0 {
		c.my.AbiSerializerMaxTimeMs = Microseconds(ms * 1000)
	} else {
		c.my.AbiSerializerMaxTimeMs = Microseconds(DefaultConfig.DefaultAbiSerializerMaxTimeMs)
	}

	c.my.ChainConfig.BlocksDir = c.my.BlockDir
	c.my.ChainConfig.StateDir = App().DataDir() + "/" + DefaultConfig.DefaultStateDirName
	c.my.ChainConfig.ReadOnly = c.my.Readonly

	//TODO UNUSED
	//if mb := options.Uint64("chain-state-db-size-mb"); mb > 0 {
	//	c.my.ChainConfig.StateSize = mb * 1024 * 1024
	//}
	//if mb := options.Uint64("chain-state-db-guard-size-mb"); mb > 0 {
	//	c.my.ChainConfig.StateGuardSize = mb * 1024 * 1024
	//}
	//if mb := options.Uint64("reversible-blocks-db-size-mb"); mb > 0 {
	//	c.my.ChainConfig.ReversibleCacheSize = mb * 1024 * 1024
	//}
	//if mb := options.Uint64("reversible-blocks-db-guard-size-mb"); mb > 0 {
	//	c.my.ChainConfig.ReversibleGuardSize = mb * 1024 * 1024
	//}

	//TODO handle wasm-runtime
	c.my.ChainConfig.ForceAllChecks = options.Bool("force-all-checks")
	c.my.ChainConfig.DisableReplayOpts = options.Bool("disable-replay-opts")
	c.my.ChainConfig.ContractsConsole = options.Bool("contracts-console")
	c.my.ChainConfig.AllowRamBillingInNotify = options.Bool("disable-ram-billing-notify-checks")

	if options.String("extract-genesis-json") != "" || options.Bool("print-genesis-json") {
		gs := types.NewGenesisState()

		if _, err := os.Stat(c.my.BlockDir + "/block.log"); !(err != nil && os.IsNotExist(err)) {
			*gs = chain.ExtractGenesisState(c.my.BlockDir)
		} else {
			log.Warn("No blocks.log found at '%s'. Using default genesis state.", c.my.BlockDir+"/blocks.log")
		}

		if options.Bool("print-genesis-json") {
			gsJson, err := json.Marshal(gs)
			if err != nil {
				log.Error("genesis_state to json error: %s", err.Error())
			} else {
				log.Info("Genesis JSON:\n%s", string(gsJson))
			}
		}

		if json := options.String("extract-genesis-json"); json != "" {
			//TODO: save json
		}

		EosThrow(&ExtractGenesisStateException{}, "extracted genesis state from blocks.log")
	}

	if blocks := options.String("export-reversible-blocks"); blocks != "" {
		//TODO: export-reversible-blocks
		EosThrow(&NodeManagementSuccess{}, "exported reversible blocks")
	}

	if options.Bool("delete-all-blocks") {
		log.Info("Deleting state database and blocks")
		if options.Uint("truncate-at-block") > 0 {
			log.Warn("The --truncate-at-block option does not make sense when deleting all blocks.")
		}
		ClearDirectoryContents(c.my.ChainConfig.StateDir)
		os.RemoveAll(c.my.BlockDir)

	} else if options.Bool("hard-replay-blockchain") {
		log.Info("Hard replay requested: deleting state database")
		ClearDirectoryContents(c.my.ChainConfig.StateDir)
		//TODO: hard-replay-blockchain

	} else if options.Bool("replay-blockchain") {
		log.Info("Replay requested: deleting state database")
		if options.Uint("truncate-at-block") > 0 {
			log.Warn("The --truncate-at-block option does not work for a regular replay of the blockchain.")
		}
		ClearDirectoryContents(c.my.ChainConfig.StateDir)
		if options.Bool("fix-reversible-blocks") {
			if !c.RecoverReversibleBlocks(fmt.Sprintf("%s/%s", c.my.ChainConfig.BlocksDir,
				DefaultConfig.DefaultReversibleBlocksDirName),
				uint32(c.my.ChainConfig.ReversibleCacheSize), "", 0) {
				log.Info("Reversible blocks database was not corrupted.")
			}
		}
	} else if options.Bool("fix-reversible-blocks") {
		if !c.RecoverReversibleBlocks(fmt.Sprintf("%s/%s", c.my.ChainConfig.BlocksDir,
			DefaultConfig.DefaultReversibleBlocksDirName),
			uint32(c.my.ChainConfig.ReversibleCacheSize), "", uint32(options.Uint("truncate-at-block"))) {
			log.Info("Reversible blocks database verified to not be corrupted. Now exiting...")
		} else {
			log.Info("Exiting after fixing reversible blocks database...")
		}
		EosThrow(&FixedReversibleDbException{}, "fixed corrupted reversible blocks database")

	} else if options.Uint("truncate-at-block") > 0 {
		log.Warn("The --truncate-at-block option can only be used with --fix-reversible-blocks without a replay or with --hard-replay-blockchain.")

	} else if reversibleBlocksFile := options.String("import-reversible-blocks"); reversibleBlocksFile != "" {
		log.Info("Importing reversible blocks from '%s'", reversibleBlocksFile)
		path := fmt.Sprintf("%s/%s", c.my.ChainConfig.BlocksDir, DefaultConfig.DefaultReversibleBlocksDirName)
		os.RemoveAll(path)

		c.ImportReversibleBlocks(path, uint32(c.my.ChainConfig.ReversibleCacheSize), reversibleBlocksFile)

		EosThrow(&NodeManagementSuccess{}, "imported reversible blocks")
	}

	if options.String("import-reversible-blocks") != "" {
		log.Warn("The --import-reversible-blocks option should be used by itself.")
	}

	if options.String("snapshot") != "" {
		//TODO: snapshot
	} else {
		if genesisFile := options.String("genesis-json"); genesisFile != "" {
			//TODO: genesis-json

		} else if genesisTimestamp := options.String("genesis-timestamp"); genesisTimestamp != "" {
			//TODO: genesis-timestamp

		} else { //TODO: else if( fc::is_regular_file( my->blocks_dir / "blocks.log" )) {}
			log.Warn("Starting up fresh blockchain with default genesis state.")
		}
	}

	if readMode := options.String("read-mode"); readMode != "" {
		if dbReadMode, ok := chain.DBReadModeFromString(readMode); ok {
			c.my.ChainConfig.ReadMode = dbReadMode
		} else {
			log.Error("The read-mode option %s format failed", readMode)
		}
		EosAssert(c.my.ChainConfig.ReadMode != chain.IRREVERSIBLE, &PluginConfigException{}, "irreversible mode not currently supported.")
	} else {
		c.my.ChainConfig.ReadMode = chain.SPECULATIVE
	}

	if validationMode := options.String("validation-mode"); validationMode != "" {
		if blockValidationMode, ok := chain.ValidationModeFromString(validationMode); ok {
			c.my.ChainConfig.BlockValidationMode = blockValidationMode
		}
	} else {
		c.my.ChainConfig.BlockValidationMode = chain.FULL
	}

	c.my.Chain = chain.NewController(c.my.ChainConfig) //TODO
	//c.my.Chain = chain.GetControllerInstance()
	c.my.ChainId = c.my.Chain.GetChainId()

	// set up method providers
	//TODO

	// relay signals to channels
	//TODO
}

func (c *ChainPlugin) PluginStartup() {
	//log.Info("Blockchain started; head block is #%d, genesis timestamp is %s",
	//	c.my.Chain.HeadBlockNum(), c.my.ChainConfig.Genesis.InitialTimestamp)
	//my->chain->head_block_num(), my->chain_config->genesis.initial_timestamp
	Try(func() {
		c.my.Chain.Startup()
	}).Catch(func(e *DatabaseGuardException){
		c.logGuardException(e)
		Throw(e)
	})

	if !c.my.Readonly {
		log.Info("starting chain in read/write mode");
	}

	log.Info("Blockchain started; head block is #%d, genesis timestamp is %s", c.my.Chain.HeadBlockNum(), c.my.ChainConfig.Genesis.InitialTimestamp)
}

func (c *ChainPlugin) PluginShutdown() {
	c.my.Chain.Close()
	log.Info("chain plugin shutdown")
}

func (c *ChainPlugin) GetReadOnlyApi() *ReadOnly {
	return NewReadOnly(c.Chain(), c.GetAbiSerializerMaxTime())
}

func (c *ChainPlugin) GetReadWriteApi() *ReadWrite {
	return NewReadWrite(c.Chain(), c.GetAbiSerializerMaxTime())
}

func (c *ChainPlugin) AcceptBlock(block *types.SignedBlock) {
	c.my.IncomingBlockSyncMethod.CallMethods(block)
}

func (c *ChainPlugin) AcceptTransaction(trx *types.PackedTransaction, next chain_interface.NextFunction) {
	c.my.IncomingTransactionAsyncMethod.CallMethods(trx, false, next)
}

func (c *ChainPlugin) RecoverReversibleBlocks(dbDir string, cacheSize uint32, newDbDir string, truncateAtBlock uint32) bool {
	//TODO: recover_reversible_blocks
	return true
}

func (c *ChainPlugin) ImportReversibleBlocks(reversibleDir string, cacheSize uint32, reversibleBlocksFile string) bool {
	//TODO: import_reversible_blocks
	return true
}

func (c *ChainPlugin) ExportReversibleBlocks(reversibleDir string, reversibleBlocksFile string) bool {
	//TODO: export_reversible_blocks
	return true
}

func (c *ChainPlugin) Chain() *chain.Controller {
	return c.my.Chain
}

func (c *ChainPlugin) GetChainId() ChainIdType {
	return c.my.ChainId
}

func (c *ChainPlugin) GetAbiSerializerMaxTime() Microseconds {
	return c.my.AbiSerializerMaxTimeMs
}

func (c *ChainPlugin) logGuardException(e GuardExceptions) {
	if e.Code() == (DatabaseGuardException{}).Code() {
		log.Error("Database has reached an unsafe level of usage, shutting down to avoid corrupting the database.  " +
			"Please increase the value set for \"chain-state-db-size-mb\" and restart the process!")
	} else if e.Code() == (ReversibleGuardException{}).Code() {
		log.Error("Reversible block database has reached an unsafe level of usage, shutting down to avoid corrupting the database.  " +
			"Please increase the value set for \"reversible-blocks-db-size-mb\" and restart the process!")
	}

	log.Debug("Details: %s", e.DetailMessage())
}

func (c *ChainPlugin) HandleGuardException(e GuardExceptions) {
	c.logGuardException(e)

	// quit the app
	App().Quit()
}

func (c *ChainPlugin) HandleDbExhaustion() {
	log.Info("database memory exhausted: increase chain-state-db-size-mb and/or reversible-blocks-db-size-mb")
	os.Exit(1)
}
