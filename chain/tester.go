package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
)

type BaseTester struct {
	ActionResult           string
	DefaultExpirationDelta uint32
	AbiSerializerMaxTime   common.Microseconds
	//TempDir                tempDirectory
	Control                 *Controller
	BlockSigningPrivateKeys map[ecc.PublicKey]ecc.PrivateKey
	Cfg                     Config
	ChainTransactions       map[common.IdType]types.TransactionReceipt
	LastProducedBlock       map[common.AccountName]common.IdType
}

func (t BaseTester) init(pushGenesis bool, mode DBReadMode) {
	t.Cfg.blocksDir = common.DefaultConfig.DefaultBlocksDirName
	t.Cfg.stateDir = common.DefaultConfig.DefaultStateDirName
	t.Cfg.stateSize = 1024 * 1024 * 8
	t.Cfg.stateGuardSize = 0
	t.Cfg.reversibleCacheSize = 1024 * 1024 * 8
	t.Cfg.reversibleGuardSize = 0
	t.Cfg.contractsConsole = true
	t.Cfg.readMode = mode

	t.Cfg.genesis.InitialTimestamp, _ = common.FromIsoString("2020-01-01T00:00:00.000")
	//t.Cfg.genesis.InitialKey = t.GetPublicKeys(common.DefaultConfig.SystemAccountName, "active")

	t.open()
	if pushGenesis {
		t.pushGenesisBlock()
	}
}

func (t BaseTester) initCfg(config Config) {
	t.Cfg = config
	t.open()
}

func (t BaseTester) open() {
	t.Control.Config = t.Cfg
	//t.Control.startUp() //TODO
	t.ChainTransactions = make(map[common.IdType]types.TransactionReceipt)
	//t.Control.AcceptedBlock.Connect() // TODO: Control.signal
}

func (t BaseTester) close() {
	t.Control.Close()
	t.ChainTransactions = make(map[common.IdType]types.TransactionReceipt)
}

func (t BaseTester) pushGenesisBlock() {
	//t.setCode()
	//t.setAbi()
}

func (t BaseTester) PushBlock(b *types.SignedBlock) *types.SignedBlock {
	t.Control.AbortBlock()
	//t.control.PushBlock(b)
	return &types.SignedBlock{}
}

func (t BaseTester) ProduceBlock(skipTime common.TimePoint, skipPendingTrxs bool, skipFlag uint32) {
	//head := t.Control.HeadBlockState()
	//headTime := t.Control.HeadBlockTime()
	//nextTime := headTime + skipTime
	//
	//if !t.Control.PendingBlockState() /*|| t.Control.PendingBlockState().Header.Timestamp != nextTime*/ {
	//
	//}
}
