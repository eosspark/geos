package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/exception"
	"math"
	"github.com/eosspark/eos-go/exception/try"
)

type BaseTester struct {
	ActionResult           string
	DefaultExpirationDelta uint32
	DefaultBilledCpuTimeUs uint32
	AbiSerializerMaxTime   common.Microseconds
	//TempDir                tempDirectory
	Control                 *Controller
	BlockSigningPrivateKeys map[ecc.PublicKey]ecc.PrivateKey
	Cfg                     Config
	ChainTransactions       map[common.IdType]types.TransactionReceipt
	LastProducedBlock       map[common.AccountName]common.IdType
}

func (t BaseTester) init(pushGenesis bool, mode DBReadMode) {
	t.DefaultExpirationDelta = 6
	t.DefaultBilledCpuTimeUs = 2000
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
	//	//headTime := t.Control.HeadBlockTime()
	//	//nextTime := headTime + skipTime
	//	//
	//	//if !t.Control.PendingBlockState() /*|| t.Control.PendingBlockState().Header.Timestamp != nextTime*/ {
	//	//
	//	//}
}

func (t BaseTester) SetTransactionHeaders(trx *types.Transaction, expiration uint32, delaySec uint32) {
	trx.Expiration = common.TimePointSec((common.Microseconds(t.Control.HeadBlockTime()) + common.Seconds(int64(expiration))).ToSeconds())
	headBlockId := t.Control.HeadBlockId()
	trx.SetReferenceBlock(&headBlockId)

	trx.MaxNetUsageWords = 0
	trx.MaxCpuUsageMS = 0
	trx.DelaySec = delaySec
}

func (t BaseTester) CreateAccounts(names []common.AccountName, multiSig bool, includeCode bool) []*types.TransactionTrace{
	traces := make ([]*types.TransactionTrace, len(names))
	for i, n := range names {
		traces[i] = t.createAccount(n, common.DefaultConfig.SystemAccountName, multiSig, includeCode)
	}
	return traces
}

func (t BaseTester) createAccount(name common.AccountName, creator common.AccountName, multiSig bool, includeCode bool) *types.TransactionTrace{
	trx := types.SignedTransaction{}
	t.SetTransactionHeaders(&trx.Transaction, 6,0) //TODO: test
	ownerAuth := types.Authority{}
	if multiSig {
		ownerAuth = types.Authority{Threshold:2, Keys:[]types.KeyWeight{{t.getPublicKey(name, "owner"), 1}}}
	} else {
		ownerAuth = types.NewAuthority(t.getPublicKey(name, "owner"), 0)
	}
	activeAuth := types.NewAuthority(t.getPublicKey(name, "active"), 0)

	sortPermissions := func(auth *types.Authority){

	}
	if includeCode {//TODO
		exception.EosAssert(ownerAuth.Threshold <= math.MaxUint16, nil,"threshold is too high")
		exception.EosAssert(uint64(activeAuth.Threshold) <= uint64(math.MaxUint64), nil,"threshold is too high")
		//ownerAuth.Accounts
		sortPermissions(&ownerAuth, )
		//ownerAuth.Accounts
		sortPermissions(&ownerAuth)
	}
	//trx.Actions
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	pk := t.getPrivateKey(creator, "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&pk, &chainId)
	return t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) PushTransaction(trx *types.SignedTransaction, deadline common.TimePoint, billedCpuTimeUs uint32) (trace *types.TransactionTrace) {
	try.Try(func(){//TODO
		//if !t.Control.PendingBlockState() {
		//}
		//c := 0
		//size, _ := rlp.EncodeSize(trx)
		//if size > 1000{
		//	c = 1
		//}
	}).Catch(func(){

	}).End()
	return trace
}

func (t BaseTester) getPrivateKey(keyName common.Name, role string) ecc.PrivateKey {
	//TODO: wait for testing
	priKey, _ := ecc.NewPrivateKey(crypto.Hash256(keyName.String()+role).String())
	return *priKey
}

func (t BaseTester) getPublicKey(keyName common.Name, role string) ecc.PublicKey {
	priKey := t.getPrivateKey(keyName, role)
	return priKey.PublicKey()
}