package unittests

import (
	"bytes"
	"encoding/json"
	"github.com/eosspark/container/sets/treeset"
	. "github.com/eosspark/eos-go/chain"
	abi "github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/chain_interface"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

var CORE_SYMBOL = common.Symbol{Precision: 4, Symbol: "SYS"}
var CORE_SYMBOL_NAME = "SYS"
var EPSINON = float64(0.001)
var eosio = common.N("eosio")
var eosioToken = common.N("eosio.token")
var eosioRam = common.N("eosio.ram")
var eosioRamFee = common.N("eosio.ramfee")
var eosioStake = common.N("eosio.stake")
var eosioBpay = common.N("eosio.bpay")
var eosioVpay = common.N("eosio.vpay")
var eosioSaving = common.N("eosio.saving")
var eosioName = common.N("eosio.names")
var eosioMsig = common.N("eosio.msig")
var eosioSudo = common.N("eosio.sudo")
var alice = common.N("alice1111111")
var bob = common.N("bob111111111")
var carol = common.N("carol1111111")
var test1 = common.N("testram11111")
var test2 = common.N("testram22222")

type ActionResult = string

func CatchThrowException(t *testing.T, expectException interface{}, function func()) {
	returning := false
	try.Try(func() {
		function()
	}).Catch(func(e exception.Exception) {
		returning = reflect.TypeOf(expectException) == reflect.TypeOf(e)
	}).End()
	assert.True(t, returning)
}

func CatchThrowMsg(t *testing.T, expectMsg string, function func()) {
	returning := false
	try.Try(func() {
		function()
	}).Catch(func(e exception.Exception) {
		returning = strings.Contains(e.DetailMessage(), expectMsg)
	}).End()
	assert.True(t, returning)
}

func CatchThrowExceptionAndMsg(t *testing.T, expectException interface{}, expectMsg string, function func()) {
	returning := false
	try.Try(func() {
		function()
	}).Catch(func(e exception.Exception) {
		returning = reflect.TypeOf(expectException) == reflect.TypeOf(e) && strings.Contains(e.DetailMessage(), expectMsg)
	}).End()
	assert.True(t, returning)
}

type BaseTester struct {
	ActionResult           string
	DefaultExpirationDelta uint32
	DefaultBilledCpuTimeUs uint32
	AbiSerializerMaxTime   common.Microseconds
	//TempDir                tempDirectory
	Control                 *Controller
	BlockSigningPrivateKeys map[string]ecc.PrivateKey //map[ecc.PublicKey]ecc.PrivateKey
	Cfg                     Config
	ChainTransactions       map[common.BlockIdType]types.TransactionReceipt
	LastProducedBlock       map[common.AccountName]common.BlockIdType
}

func newBaseTester(pushGenesis bool, readMode DBReadMode) *BaseTester {
	t := &BaseTester{}
	t.DefaultExpirationDelta = 6
	t.DefaultBilledCpuTimeUs = 2000
	t.AbiSerializerMaxTime = 1000 * 1000
	t.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	t.LastProducedBlock = make(map[common.AccountName]common.BlockIdType)

	t.init(pushGenesis, readMode)
	return t
}

//for forked_test
func newBaseTesterSecNode(pushGenesis bool, readMode DBReadMode) *BaseTester {
	t := &BaseTester{}
	t.DefaultExpirationDelta = 6
	t.DefaultBilledCpuTimeUs = 2000
	t.AbiSerializerMaxTime = 1000 * 1000
	t.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	t.LastProducedBlock = make(map[common.AccountName]common.BlockIdType)

	cfg := newConfig(readMode)

	t.Control = NewController(cfg)
	if pushGenesis {
		t.pushGenesisBlock()
	}
	return t
}

func (t *BaseTester) init(pushGenesis bool, readMode DBReadMode) {
	t.Cfg = *newConfig(readMode)

	t.open()
	if pushGenesis {
		t.pushGenesisBlock()
	}
}

func newConfig(readMode DBReadMode) *Config {
	cfg := &Config{}
	tempDirSuffix := "_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	cfg.BlocksDir = common.DefaultConfig.DefaultBlocksDirName + tempDirSuffix
	cfg.StateDir = common.DefaultConfig.DefaultStateDirName + tempDirSuffix
	cfg.ReversibleDir = common.DefaultConfig.DefaultReversibleBlocksDirName + tempDirSuffix
	cfg.StateSize = 1024 * 1024 * 8
	cfg.StateGuardSize = 0
	cfg.ReversibleCacheSize = 1024 * 1024 * 8
	cfg.ReversibleGuardSize = 0
	cfg.ContractsConsole = true
	cfg.ReadMode = readMode

	cfg.Genesis = types.NewGenesisState()
	cfg.Genesis.InitialTimestamp, _ = common.FromIsoString("2020-01-01T00:00:00.000")
	cfg.Genesis.InitialKey = BaseTester{}.getPublicKey(eosio, "active")

	cfg.ActorWhitelist = *treeset.NewWith(common.TypeName, common.CompareName)
	cfg.ActorBlacklist = *treeset.NewWith(common.TypeName, common.CompareName)
	cfg.ContractWhitelist = *treeset.NewWith(common.TypeName, common.CompareName)
	cfg.ContractBlacklist = *treeset.NewWith(common.TypeName, common.CompareName)
	cfg.ActionBlacklist = *treeset.NewWith(common.TypePair, common.ComparePair)
	cfg.KeyBlacklist = *treeset.NewWith(ecc.TypePubKey, ecc.ComparePubKey)
	cfg.ResourceGreylist = *treeset.NewWith(common.TypeName, common.CompareName)
	cfg.TrustedProducers = *treeset.NewWith(common.TypeName, common.CompareName)

	//cfg.VmType = common.DefaultConfig.DefaultWasmRuntime // TODO

	return cfg
}

func (t *BaseTester) open() {
	t.Control = NewController(&t.Cfg)
	//TODO
	//t.Control.AddIndices()
	//t.Control.startUp()
	t.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	ab := chain_interface.AcceptedBlockCaller{Caller: t.acceptedBlock}
	t.Control.AcceptedBlock.Connect(&ab)
}

func (t *BaseTester) acceptedBlock(b *types.BlockState) {
	try.EosAssert(b.SignedBlock != nil, &exception.BlockLogNotFound{}, "tester acceptedBlock is not found")
	for _, receipt := range b.SignedBlock.Transactions {
		if !common.Empty(receipt.Trx.PackedTransaction) {
			t.ChainTransactions[receipt.Trx.PackedTransaction.ID()] = receipt
		} else {
			id := receipt.Trx.TransactionID
			t.ChainTransactions[id] = receipt
		}
	}
}

func (t *BaseTester) close() {
	t.Control.Close()
	t.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
}

func (t BaseTester) IsSameChain(other *BaseTester) bool {
	return t.Control.HeadBlockId() == other.Control.HeadBlockId()
}

func (t BaseTester) PushBlock(b *types.SignedBlock) *types.SignedBlock {
	t.Control.AbortBlock()
	t.Control.PushBlock(b, types.Complete)
	id := b.BlockID()
	if v, ok := t.LastProducedBlock[b.Producer]; !ok || types.NumFromID(&id) > types.NumFromID(&v) {
		t.LastProducedBlock[b.Producer] = b.BlockID()
	}
	return b
}

func (t BaseTester) PushBlock2(b *types.SignedBlock, status types.BlockStatus) *types.SignedBlock {
	t.Control.AbortBlock()
	t.Control.PushBlock(b, status)
	id := b.BlockID()
	if v, ok := t.LastProducedBlock[b.Producer]; !ok || types.NumFromID(&id) > types.NumFromID(&v) {
		t.LastProducedBlock[b.Producer] = b.BlockID()
	}
	return b
}

func (t BaseTester) pushGenesisBlock() {
	wasmName := "test_contracts/eosio.bios.wasm"
	code, err := ioutil.ReadFile(wasmName)
	if err != nil {
		log.Error("pushGenesisBlock is err : %v", err)
	}
	t.SetCode(eosio, code, nil)
	abiName := "test_contracts/eosio.bios.abi"
	abi, err := ioutil.ReadFile(abiName)
	if err != nil {
		log.Error("pushGenesisBlock is err : %v", err)
	}
	t.SetAbi(eosio, abi, nil)
}

func (t BaseTester) ProduceBlocks(n uint32, empty bool) {
	if empty {
		for i := 0; uint32(i) < n; i++ {
			t.ProduceEmptyBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
		}
	} else {
		for i := 0; uint32(i) < n; i++ {
			t.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
		}
	}
}

func (t BaseTester) produceBlock(skipTime common.Microseconds, skipPendingTrxs bool, skipFlag uint32) *types.SignedBlock {
	headTime := t.Control.HeadBlockTime()
	nextTime := headTime + common.TimePoint(skipTime)
	if common.Empty(t.Control.PendingBlockState()) || t.Control.PendingBlockState().Header.Timestamp != types.NewBlockTimeStamp(nextTime) {
		t.startBlock(nextTime)
	}
	Hbs := t.Control.HeadBlockState()
	producer := Hbs.GetScheduledProducer(types.NewBlockTimeStamp(nextTime))
	privKey := ecc.PrivateKey{}
	privateKey, ok := t.BlockSigningPrivateKeys[producer.BlockSigningKey.String()]
	if !ok {
		privKey = t.getPrivateKey(producer.ProducerName, "active")
	} else {
		privKey = privateKey
	}

	if !skipPendingTrxs {
		unappliedTrxs := t.Control.GetUnappliedTransactions()
		for _, trx := range unappliedTrxs {
			trace := t.Control.PushTransaction(trx, common.MaxTimePoint(), 0)
			if trace.Except != nil {
				try.Throw(trace.Except)
			}
		}

		scheduledTrxs := t.Control.GetScheduledTransactions()
		for len(scheduledTrxs) > 0 {
			for _, trx := range scheduledTrxs {
				trace := t.Control.PushScheduledTransaction(&trx, common.MaxTimePoint(), 0)
				if trace.Except != nil {
					//try.Throw(trace.Except)
				}
			}
			scheduledTrxs = t.Control.GetScheduledTransactions()
		}
	}
	t.Control.FinalizeBlock()
	t.Control.SignBlock(func(d common.DigestType) ecc.Signature {
		sign, err := privKey.Sign(d.Bytes())
		if err != nil {
			log.Error(err.Error())
		}
		return sign
	})
	t.Control.CommitBlock(true)
	b := t.Control.HeadBlockState()
	t.LastProducedBlock[t.Control.HeadBlockState().Header.Producer] = b.BlockId
	t.startBlock(nextTime + common.TimePoint(common.TimePoint(common.DefaultConfig.BlockIntervalUs)))
	return t.Control.HeadBlockState().SignedBlock
}

func (t BaseTester) startBlock(blockTime common.TimePoint) {
	headBlockNumber := t.Control.HeadBlockNum()
	producer := t.Control.HeadBlockState().GetScheduledProducer(types.NewBlockTimeStamp(blockTime))
	lastProducedBlockNum := t.Control.LastIrreversibleBlockNum()
	itr := t.LastProducedBlock[producer.ProducerName]
	if !common.Empty(itr) {
		if t.Control.LastIrreversibleBlockNum() > types.NumFromID(&itr) {
			lastProducedBlockNum = t.Control.LastIrreversibleBlockNum()
		} else {
			lastProducedBlockNum = types.NumFromID(&itr)
		}
	}
	t.Control.AbortBlock()
	t.Control.StartBlock(types.NewBlockTimeStamp(blockTime), uint16(headBlockNumber-lastProducedBlockNum))
}

func (t BaseTester) SetTransactionHeaders(trx *types.Transaction, expiration uint32, delaySec uint32) {
	trx.Expiration = common.TimePointSec((common.Microseconds(t.Control.HeadBlockTime()) + common.Seconds(int64(expiration))).ToSeconds())
	headBlockId := t.Control.HeadBlockId()
	trx.SetReferenceBlock(&headBlockId)

	trx.MaxNetUsageWords = 0
	trx.MaxCpuUsageMS = 0
	trx.DelaySec = delaySec
}

func (t BaseTester) CreateAccounts(names []common.AccountName, multiSig bool, includeCode bool) []*types.TransactionTrace {
	traces := make([]*types.TransactionTrace, len(names))
	for i, n := range names {
		traces[i] = t.CreateAccount(n, eosio, multiSig, includeCode)
	}
	return traces
}

func (t BaseTester) CreateAccount(name common.AccountName, creator common.AccountName, multiSig bool, includeCode bool) *types.TransactionTrace {
	trx := types.SignedTransaction{}
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	ownerAuth := types.Authority{}
	if multiSig {
		ownerAuth = types.Authority{
			Threshold: 2,
			Keys:      []types.KeyWeight{{Key: t.getPublicKey(name, "owner"), Weight: 1}},
			Accounts:  []types.PermissionLevelWeight{{Permission: types.PermissionLevel{Actor: creator, Permission: common.DefaultConfig.ActiveName}, Weight: 1}},
		}
	} else {
		ownerAuth = types.NewAuthority(t.getPublicKey(name, "owner"), 0)
	}
	activeAuth := types.NewAuthority(t.getPublicKey(name, "active"), 0)

	sortPermissions := func(auth *types.Authority) {

	}
	if includeCode {
		try.EosAssert(ownerAuth.Threshold <= math.MaxUint16, nil, "threshold is too high")
		try.EosAssert(uint64(activeAuth.Threshold) <= uint64(math.MaxUint64), nil, "threshold is too high")
		ownerAuth.Accounts = append(ownerAuth.Accounts, types.PermissionLevelWeight{
			Permission: types.PermissionLevel{Actor: name, Permission: common.DefaultConfig.EosioCodeName},
			Weight:     types.WeightType(ownerAuth.Threshold),
		})
		sortPermissions(&ownerAuth)
		activeAuth.Accounts = append(activeAuth.Accounts, types.PermissionLevelWeight{
			Permission: types.PermissionLevel{Actor: name, Permission: common.DefaultConfig.EosioCodeName},
			Weight:     types.WeightType(activeAuth.Threshold),
		})
		sortPermissions(&activeAuth)
	}
	new := NewAccount{
		Creator: creator,
		Name:    name,
		Owner:   ownerAuth,
		Active:  activeAuth,
	}
	data, _ := rlp.EncodeToBytes(new)
	act := &types.Action{
		Account:       new.GetAccount(),
		Name:          new.GetName(),
		Authorization: []types.PermissionLevel{{creator, common.DefaultConfig.ActiveName}},
		Data:          data,
	}
	trx.Actions = append(trx.Actions, act)

	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	pk := t.getPrivateKey(creator, "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&pk, &chainId)
	return t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) PushTransaction(trx *types.SignedTransaction, deadline common.TimePoint, billedCpuTimeUs uint32) (trace *types.TransactionTrace) {
	_, r := false, (*types.TransactionTrace)(nil)
	try.Try(func() {
		if t.Control.PendingBlockState() == nil {
			t.startBlock(t.Control.HeadBlockTime().AddUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs)))
		}
		c := types.CompressionNone
		size, _ := rlp.EncodeSize(trx)
		if size > 1000 {
			c = types.CompressionZlib
		}
		mtrx := types.NewTransactionMetadataBySignedTrx(trx, c)
		trace = t.Control.PushTransaction(mtrx, deadline, billedCpuTimeUs)
		if trace.ExceptPtr != nil {
			try.Throw(trace.ExceptPtr)
			//try.EosThrow(trace.ExceptPtr, "tester PushTransaction is error :%#v", trace.ExceptPtr.DetailMessage())
		}
		if trace.Except != nil {
			try.Throw(trace.ExceptPtr)
			//try.EosThrow(trace.Except, "tester PushTransaction is error :%#v", trace.Except.DetailMessage())
		}
		r = trace
		return
		//}).FcCaptureAndRethrow().End()
	}).FcRethrowExceptions(log.LvlWarn, "transaction_header: %#v, billed_cpu_time_us: %d", trx.Transaction.TransactionHeader, billedCpuTimeUs)
	return r
}

func (t BaseTester) PushAction(act *types.Action, authorizer common.AccountName) ActionResult {
	trx := types.SignedTransaction{}
	if !common.Empty(authorizer) {
		act.Authorization = []types.PermissionLevel{{authorizer, common.DefaultConfig.ActiveName}}
	}
	trx.Actions = append(trx.Actions, act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	if !common.Empty(authorizer) {
		chainId := t.Control.GetChainId()
		privateKey := t.getPrivateKey(authorizer, "active")
		trx.Sign(&privateKey, &chainId)
	}
	try.Try(func() {
		t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
	}).Catch(func(ex exception.Exception) {
		//log.Error("tester PushAction is error: %v", ex.DetailMessage())
		try.Throw(ex)
	}).End()
	t.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	//BOOST_REQUIRE_EQUAL(true, chain_has_transaction(trx.id()))
	return t.Success()
}

func (t BaseTester) PushAction2(code *common.AccountName, acttype *common.AccountName,
	actor common.AccountName, data *common.Variants, expiration uint32, delaySec uint32) *types.TransactionTrace {
	auths := make([]types.PermissionLevel, 0)
	auths = append(auths, types.PermissionLevel{Actor: actor, Permission: common.DefaultConfig.ActiveName})
	return t.PushAction4(code, acttype, &auths, data, expiration, delaySec)
}

func (t BaseTester) PushAction3(code *common.AccountName, acttype *common.AccountName,
	actors []*common.AccountName, data *common.Variants, expiration uint32, delaySec uint32) *types.TransactionTrace {
	auths := make([]types.PermissionLevel, 0)
	for _, actor := range actors{
		auths = append(auths, types.PermissionLevel{Actor:*actor, Permission:common.DefaultConfig.ActiveName})
	}
	return t.PushAction4(code, acttype, &auths, data, expiration, delaySec)
}

func (t BaseTester) PushAction4(code *common.AccountName, acttype *common.AccountName,
	auths *[]types.PermissionLevel, data *common.Variants, expiration uint32, delaySec uint32) *types.TransactionTrace {
	trx := types.SignedTransaction{}
	try.Try(func() {
		action := t.GetAction(*code, *acttype, *auths, data)
		trx.Actions = append(trx.Actions, action)
	})
	t.SetTransactionHeaders(&trx.Transaction, expiration, delaySec)
	chainId := t.Control.GetChainId()
	key := ecc.PrivateKey{}
	for _, auth := range *auths {
		key = t.getPrivateKey(auth.Actor, auth.Permission.String())
		trx.Sign(&key, &chainId)
	}
	return t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) GetResolver() func(name common.AccountName) *abi.AbiSerializer {
	return func(name common.AccountName) *abi.AbiSerializer {
		var r *abi.AbiSerializer
		try.Try(func() {
			accObj := entity.AccountObject{Name: name}
			t.Control.DB.Find("byName", accObj, &accObj)
			var abid abi.AbiDef
			if abi.ToABI(accObj.Abi, &abid) {
				r = abi.NewAbiSerializer(&abid, t.AbiSerializerMaxTime)
			}
		}).FcRethrowExceptions(log.LvlError, "Failed to find or parse ABI for %s", name)
		return r
	}
}

func (t BaseTester) GetAction(code common.AccountName, actType common.AccountName,
	auths []types.PermissionLevel, data *common.Variants) *types.Action {

	acnt := t.Control.GetAccount(code)
	a := acnt.GetAbi()
	action := types.Action{code, actType, auths, nil}
	buf, _ := json.Marshal(data)
	action.Data, _ = a.EncodeAction(actType, buf)
	return &action

	//action := types.Action{code, actType, auths, nil}
	//try.Try(func() {
	//	acnt := t.Control.GetAccount(code)
	//	a := acnt.GetAbi()
	//	abis := abi.NewAbiSerializer(a, t.AbiSerializerMaxTime)
	//	actionTypeName := abis.GetActionType(actType)
	//	try.FcAssert(reflect.TypeOf(actionTypeName).Kind() == reflect.String, "unknown action type %s", actType)
	//
	//	action.Data = abis.VariantToBinary(actionTypeName, data, t.AbiSerializerMaxTime)
	//
	//}).FcCaptureAndRethrow().End()
	//return &action
}

func (t BaseTester) getPrivateKey(keyName common.Name, role string) ecc.PrivateKey {
	pk := &ecc.PrivateKey{}
	if keyName == eosio {
		pk, _ = ecc.NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	} else {
		rawPrivKey := crypto.Hash256String(keyName.String() + role).Bytes()
		g := bytes.NewReader(rawPrivKey)
		pk, _ = ecc.NewDeterministicPrivateKey(g)
	}
	return *pk
}

func (t BaseTester) getPublicKey(keyName common.Name, role string) ecc.PublicKey {
	priKey := t.getPrivateKey(keyName, role)
	return priKey.PublicKey()
}

func (t BaseTester) ProduceBlocksUntileEndOfRound() {
	var blocksPerRound uint64
	for {
		blocksPerRound = uint64(len(t.Control.ActiveProducers().Producers) * common.DefaultConfig.ProducerRepetitions)
		t.ProduceBlocks(1, false)

		if uint64(t.Control.HeadBlockNum())%blocksPerRound == blocksPerRound-1 {
			break
		}
	}
}

func (t BaseTester) ProduceBlocksForNrounds(numOfRounds int) {

	for i := 0; i < numOfRounds; i++ {
		t.ProduceBlocksUntileEndOfRound()
	}
}

func (t BaseTester) ProduceMinNumOfBlocksToSpendTimeWoInactiveProd(targetElapsedTime common.Microseconds) {

	var elapsedTime common.Microseconds

	for elapsedTime < targetElapsedTime {
		for i := 0; i < len(t.Control.HeadBlockState().ActiveSchedule.Producers); i++ {
			timeToSkip := common.Microseconds(int64(common.DefaultConfig.ProducerRepetitions) * common.DefaultConfig.BlockIntervalMs)
			t.produceBlock(timeToSkip, false, 0)

			elapsedTime += timeToSkip
		}

		timeToSkip := common.Seconds(23 * 60 * 60)
		t.produceBlock(timeToSkip, false, 0)

		elapsedTime += timeToSkip
	}
}

func (t BaseTester) ProduceBlock(skipTime common.Microseconds, skipFlag uint32) *types.SignedBlock {
	return t.produceBlock(skipTime, false, skipFlag)
}

func (t BaseTester) ProduceBlockNoValidation(skipTime common.Microseconds, skipFlag uint32) *types.SignedBlock {
	return t.produceBlock(skipTime, false, skipFlag|2)
}

func (t BaseTester) ProduceEmptyBlock(skipTime common.Microseconds, skipFlag uint32) *types.SignedBlock {
	t.Control.AbortBlock()
	return t.produceBlock(skipTime, true, skipFlag)
}

func (t BaseTester) PushReqAuth(from common.AccountName, auths *[]types.PermissionLevel, keys *[]ecc.PrivateKey) *types.TransactionTrace {
	trx := types.SignedTransaction{}
	type params struct {
		From common.AccountName
	}
	ps := params{From: from}
	data, _ := rlp.EncodeToBytes(ps)
	act := types.Action{
		Account:       eosio,
		Name:          common.ActionName(common.N("reqauth")),
		Authorization: *auths,
		Data:          data,
	}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	chainId := t.Control.GetChainId()
	for _, iter := range *keys {
		trx.Sign(&iter, &chainId)
	}
	return t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) PushReqAuth2(from common.AccountName, role string, multiSig bool) *types.TransactionTrace {
	if multiSig {
		auths := []types.PermissionLevel{{Actor: from, Permission: common.DefaultConfig.OwnerName}}
		keys := []ecc.PrivateKey{t.getPrivateKey(from, role), t.getPrivateKey(eosio, "active")}
		return t.PushReqAuth(from, &auths, &keys)
	} else {
		auths := []types.PermissionLevel{{Actor: from, Permission: common.DefaultConfig.OwnerName}}
		keys := []ecc.PrivateKey{t.getPrivateKey(from, role)}
		return t.PushReqAuth(from, &auths, &keys)
	}
}

func (t BaseTester) PushDummy(from common.AccountName, v *string, billedCpuTimeUs uint32) *types.TransactionTrace {
	//TODO
	trx := types.SignedTransaction{}
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	privKey := t.getPrivateKey(from, "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&privKey, &chainId)
	return t.PushTransaction(&trx, common.MaxTimePoint(), billedCpuTimeUs)
}

func (t BaseTester) Transfer(from common.AccountName, to common.AccountName, amount common.Asset, memo string, currency common.AccountName) *types.TransactionTrace {
	trx := types.SignedTransaction{}
	data := common.Variants{
		"from":     from,
		"to":       to,
		"quantity": amount,
		"memo":     memo,
	}
	acttype := common.N("transfer")
	act := t.GetAction(currency, acttype, []types.PermissionLevel{{from, common.N("active")}}, &data)
	trx.Actions = append(trx.Actions, act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	privKey := t.getPrivateKey(from, "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&privKey, &chainId)
	return t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) Transfer2(from common.AccountName, to common.AccountName, amount string, memo string, currency common.AccountName) *types.TransactionTrace {
	return t.Transfer(from, to, common.Asset{}.FromString(&amount), memo, currency)
}

func (t BaseTester) Issue(to common.AccountName, amount string, currency common.AccountName) *types.TransactionTrace {
	trx := types.SignedTransaction{}
	data := common.Variants{
		"to":       to,
		"quantity": amount,
	}
	acttype := common.N("issue")
	act := t.GetAction(currency, acttype, []types.PermissionLevel{{currency, common.N("active")}}, &data)
	trx.Actions = append(trx.Actions, act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	privKey := t.getPrivateKey(currency, "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&privKey, &chainId)
	return t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) LinkAuthority(account common.AccountName, code common.AccountName, req common.PermissionName, rtype common.ActionName) {
	trx := types.SignedTransaction{}
	link := LinkAuth{Account: account, Code: code, Type: rtype, Requirement: req}
	data, _ := rlp.EncodeToBytes(link)
	act := types.Action{
		Account:       link.GetAccount(),
		Name:          link.GetName(),
		Authorization: []types.PermissionLevel{{account, common.DefaultConfig.ActiveName}},
		Data:          data,
	}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	privKey := t.getPrivateKey(account, "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&privKey, &chainId)
	t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) UnlinkAuthority(account common.AccountName, code common.AccountName, rtype common.ActionName) {
	trx := types.SignedTransaction{}
	unlink := UnLinkAuth{Account: account, Code: code, Type: rtype}
	data, _ := rlp.EncodeToBytes(unlink)
	act := types.Action{
		Account:       unlink.GetAccount(),
		Name:          unlink.GetName(),
		Authorization: []types.PermissionLevel{{account, common.DefaultConfig.ActiveName}},
		Data:          data,
	}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	privKey := t.getPrivateKey(account, "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&privKey, &chainId)
	t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) SetAuthority(account common.AccountName, perm common.PermissionName, auth types.Authority, parent common.PermissionName, auths *[]types.PermissionLevel, keys *[]ecc.PrivateKey) {
	trx := types.SignedTransaction{}
	update := UpdateAuth{Account: account, Permission: perm, Parent: parent, Auth: auth}
	data, _ := rlp.EncodeToBytes(update)
	act := types.Action{
		Account:       update.GetAccount(),
		Name:          update.GetName(),
		Authorization: *auths,
		Data:          data,
	}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	chainId := t.Control.GetChainId()
	for _, key := range *keys {
		trx.Sign(&key, &chainId)
	}
	t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) SetAuthority2(account common.AccountName, perm common.PermissionName, auth types.Authority, parent common.PermissionName) {
	permL := types.PermissionLevel{Actor: account, Permission: common.DefaultConfig.OwnerName}
	privKey := t.getPrivateKey(account, "owner")
	t.SetAuthority(account, perm, auth, parent, &[]types.PermissionLevel{permL}, &[]ecc.PrivateKey{privKey})
}

func (t BaseTester) DeleteAuthority(account common.AccountName, perm common.PermissionName, auths *[]types.PermissionLevel, keys *[]ecc.PrivateKey) {
	trx := types.SignedTransaction{}
	delete := DeleteAuth{Account: account, Permission: perm}
	data, _ := rlp.EncodeToBytes(delete)
	act := types.Action{
		Account:       delete.GetAccount(),
		Name:          delete.GetName(),
		Authorization: *auths,
		Data:          data,
	}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	chainId := t.Control.GetChainId()
	for _, key := range *keys {
		trx.Sign(&key, &chainId)
	}
	t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) DeleteAuthority2(account common.AccountName, perm common.PermissionName) {
	permL := types.PermissionLevel{Actor: account, Permission: common.DefaultConfig.OwnerName}
	privKey := t.getPrivateKey(account, "owner")
	t.DeleteAuthority(account, perm, &[]types.PermissionLevel{permL}, &[]ecc.PrivateKey{privKey})
}

func (t BaseTester) SetCode(account common.AccountName, wasm []uint8, signer *ecc.PrivateKey) {
	trx := types.SignedTransaction{}
	setCode := SetCode{Account: account, VmType: 0, VmVersion: 0, Code: wasm}
	data, _ := rlp.EncodeToBytes(setCode)
	act := types.Action{
		Account:       setCode.GetAccount(),
		Name:          setCode.GetName(),
		Authorization: []types.PermissionLevel{{account, common.DefaultConfig.ActiveName}},
		Data:          data,
	}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	chainId := t.Control.GetChainId()
	if signer != nil {
		trx.Sign(signer, &chainId)
	} else {
		privKey := t.getPrivateKey(account, "active")
		trx.Sign(&privKey, &chainId)
	}
	t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) SetCode2(account common.AccountName, wast *byte, signer *ecc.PrivateKey) {
	//t.SetCode(account, wastToWasm(wast), signer)
}

func (t BaseTester) SetAbi(account common.AccountName, abiJson []byte, signer *ecc.PrivateKey) {
	abiEt := abi.AbiDef{}
	err := json.Unmarshal(abiJson, &abiEt)
	if err != nil {
		log.Error("unmarshal abiJson is wrong :%s", err)
	}
	trx := types.SignedTransaction{}
	abiBytes, _ := rlp.EncodeToBytes(abiEt)
	setAbi := SetAbi{Account: account, Abi: abiBytes}
	data, _ := rlp.EncodeToBytes(setAbi)
	act := types.Action{
		Account:       setAbi.GetAccount(),
		Name:          setAbi.GetName(),
		Authorization: []types.PermissionLevel{{account, common.DefaultConfig.ActiveName}},
		Data:          data,
	}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	chainId := t.Control.GetChainId()
	if signer != nil {
		trx.Sign(signer, &chainId)
	} else {
		privKey := t.getPrivateKey(account, "active")
		trx.Sign(&privKey, &chainId)
	}
	t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) ChainHasTransaction(txId *common.BlockIdType) bool {
	_, ok := t.ChainTransactions[*txId]
	return ok
}

func (t BaseTester) GetTransactionReceipt(txId *common.BlockIdType) *types.TransactionReceipt {
	if val, ok := t.ChainTransactions[*txId]; ok {
		return &val
	}
	return nil
}

func (t BaseTester) GetCurrencyBalance(code *common.AccountName, assetSymbol *common.Symbol, account *common.AccountName) common.Asset {
	db := t.Control.DB
	table := entity.TableIdObject{Code: *code, Scope: *account, Table: common.N("accounts")}
	err := db.Find("byCodeScopeTable", table, &table)
	result := int64(0)
	if err != nil {
		log.Error("GetCurrencyBalance is error: %s", err)
	} else {
		obj := entity.KeyValueObject{TId: table.ID, PrimaryKey: uint64(assetSymbol.ToSymbolCode())}
		err := db.Find("byScopePrimary", obj, &obj)
		if err != nil {
			log.Error("GetCurrencyBalance is error: %s", err)
		} else {
			rlp.DecodeBytes(obj.Value, &result)
		}
	}
	return common.Asset{Amount: result, Symbol: *assetSymbol}
}

func (t BaseTester) GetRowByAccount(code uint64, scope uint64, table uint64, act uint64) []byte {
	var data []byte
	db := t.Control.DB
	tId := entity.TableIdObject{Code: common.AccountName(code), Scope: common.ScopeName(scope), Table: common.TableName(table)}
	err := db.Find("byCodeScopeTable", tId, &tId)
	if err != nil {
		return data
		//log.Error("GetRowByAccount is error: %s", err)
	}
	idx, _ := db.GetIndex("byScopePrimary", entity.KeyValueObject{})
	obj := entity.KeyValueObject{TId: tId.ID, PrimaryKey: act}
	itr, _ := idx.LowerBound(&obj)
	if idx.CompareEnd(itr) {
		return data
	}

	objLowerBound := entity.KeyValueObject{}
	itr.Data(&objLowerBound)
	if objLowerBound.TId != tId.ID || objLowerBound.PrimaryKey != act {
		return data
	}

	data = make([]byte, len([]byte(objLowerBound.Value)))
	copy(data, []byte(objLowerBound.Value))
	return data
}

func (t BaseTester) Uint64ToUint8Vector(x uint64) []uint8 {
	//TODO
	v := make([]byte, 8)
	for i := uint8(0); i < 8; i++ {
		v[i] = uint8(x >> (8 * i))
	}
	return v
}

func (t BaseTester) StringToUint8Vector(s string) []uint8 {
	return []byte(s)
}

func (t BaseTester) ToUint64(x common.Variant) uint64 {
	var bytes []uint8
	common.FromVariant(x, &bytes)
	var re uint64
	try.FcAssert(len(bytes) == 8)
	for i := uint8(0); i < 8; i++ {
		re += uint64(bytes[i]) << (8 * i)
	}

	return re
}
func (t BaseTester) ToString(x common.Variant) string {
	var bytes []uint8
	common.FromVariant(x, &bytes)
	return string(bytes)
}

func (t BaseTester) SyncWith(other *BaseTester) {
	if t.Control.HeadBlockId() == other.Control.HeadBlockId() {
		return
	}

	if t.Control.HeadBlockNum() == other.Control.HeadBlockNum() {
		other.SyncWith(&t)
		return
	}

	syncDbs := func(a *BaseTester, b *BaseTester) {
		for i := uint32(1); i <= a.Control.HeadBlockNum(); i++ {
			block := a.Control.FetchBlockByNumber(i)
			if !common.Empty(block) {
				b.Control.AbortBlock()
				b.Control.PushBlock(block, types.Complete)
			}
		}
	}
	syncDbs(&t, other)
	syncDbs(other, &t)
}

func (t BaseTester) PushGenesisBlock() {
	t.SetCode2(eosio, nil, nil)
	t.SetAbi(eosio, nil, nil)
}

func (t BaseTester) GetProducerKeys(producerNames *[]common.AccountName) []types.ProducerKey {
	var schedule []types.ProducerKey
	for _, producerName := range *producerNames {
		pk := types.ProducerKey{ProducerName: common.AccountName(producerName), BlockSigningKey: t.getPublicKey(common.AccountName(producerName), "active")}
		schedule = append(schedule, pk)
	}
	return schedule
}

func (t BaseTester) SetProducers(producerNames *[]common.AccountName) *types.TransactionTrace {
	schedule := t.GetProducerKeys(producerNames)
	actName := common.N("setprods")
	return t.PushAction2(
		&eosio,
		&actName,
		eosio,
		&common.Variants{"schedule": schedule},
		t.DefaultExpirationDelta,
		0,
	)
}

func (t BaseTester) FindTable(code common.Name, scope common.Name, table common.Name) *entity.TableIdObject {
	tId := entity.TableIdObject{Code: code, Scope: scope, Table: table}
	err := t.Control.DB.Find("byCodeScopeTable", tId, &tId)
	if err != nil {
		log.Error("FindTable is error: %s", err)
	}
	return &tId
}

func (t BaseTester) Success() ActionResult {
	return "success"
}

type ValidatingTester struct {
	BaseTester
	ValidatingControl                 *Controller
	VCfg                              Config
	NumBlocksToProducerBeforeShutdown uint32
	SkipValidate                      bool
}

func newValidatingTester(pushGenesis bool, readMode DBReadMode) *ValidatingTester {
	vt := &ValidatingTester{}
	vt.DefaultExpirationDelta = 6
	vt.DefaultBilledCpuTimeUs = 2000
	vt.AbiSerializerMaxTime = 1000 * 1000
	vt.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	vt.LastProducedBlock = make(map[common.AccountName]common.BlockIdType)
	vt.VCfg = *newConfig(readMode)

	vt.ValidatingControl = NewController(&vt.VCfg)
	//TODO
	//vt.ValidatingControl.AddIndices()
	//vt.ValidatingControl.startUp()

	vt.init(pushGenesis, readMode)
	return vt
}

func (vt ValidatingTester) DefaultProduceBlock() *types.SignedBlock {
	skipTime := common.DefaultConfig.BlockIntervalMs
	skipFlag := uint32(0)
	sb := vt.produceBlock(common.Milliseconds(skipTime), false, skipFlag|2)
	vt.ValidatingControl.PushBlock(sb, types.Complete)
	return sb
}

func NewValidatingTesterTrustedProducers(trustedProducers *treeset.Set) *ValidatingTester {
	vt := &ValidatingTester{}
	vt.DefaultExpirationDelta = 6
	vt.DefaultBilledCpuTimeUs = 2000
	vt.AbiSerializerMaxTime = 1000 * 1000
	vt.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	vt.LastProducedBlock = make(map[common.AccountName]common.BlockIdType)
	vt.VCfg = *newConfig(SPECULATIVE)
	vt.VCfg.TrustedProducers = *trustedProducers
	vt.ValidatingControl = NewController(&vt.VCfg)

	vt.init(true, SPECULATIVE)

	return vt
}

func (vt ValidatingTester) PushAction(act *types.Action, authorizer common.AccountName) ActionResult {
	trx := types.SignedTransaction{}
	if !common.Empty(authorizer) {
		act.Authorization = []types.PermissionLevel{{authorizer, common.DefaultConfig.ActiveName}}
	}
	trx.Actions = append(trx.Actions, act)
	vt.SetTransactionHeaders(&trx.Transaction, vt.DefaultExpirationDelta, 0)
	if !common.Empty(authorizer) {
		chainId := vt.Control.GetChainId()
		privateKey := vt.getPrivateKey(authorizer, "active")
		trx.Sign(&privateKey, &chainId)
	}
	try.Try(func() {
		vt.PushTransaction(&trx, common.MaxTimePoint(), vt.DefaultBilledCpuTimeUs)
	}).Catch(func(ex exception.Exception) {
		//log.Error("tester PushAction is error: %v", ex.DetailMessage())
		try.Throw(ex)
	}).End()
	vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	//BOOST_REQUIRE_EQUAL(true, chain_has_transaction(trx.id()))
	return vt.Success()
}

func (vt ValidatingTester) ProduceBlocks(n uint32, empty bool) {
	if empty {
		for i := 0; uint32(i) < n; i++ {
			vt.ProduceEmptyBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
		}
	} else {
		for i := 0; uint32(i) < n; i++ {
			vt.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
		}
	}
}

func (vt ValidatingTester) ProduceBlock(skipTime common.Microseconds, skipFlag uint32) *types.SignedBlock {
	sb := vt.produceBlock(skipTime, false, skipFlag|2)
	vt.ValidatingControl.PushBlock(sb, types.Complete)
	return sb
}

func (vt ValidatingTester) ProduceEmptyBlock(skipTime common.Microseconds, skipFlag uint32) *types.SignedBlock {
	sb := vt.produceBlock(skipTime, true, skipFlag|2)
	vt.ValidatingControl.PushBlock(sb, types.Complete)
	return sb
}

func (vt ValidatingTester) Validate() bool {
	return true
}

func (vt ValidatingTester) ValidatePushBlock(sb *types.SignedBlock) {
	vt.ValidatingControl.PushBlock(sb, types.Complete)
}

func (vt *ValidatingTester) close() {
	vt.Control.Close()
	vt.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
}

func CoreFromString(s string) common.Asset {
	str := s + " " + CORE_SYMBOL_NAME
	return common.Asset{}.FromString(&str)
}

func (vt *ValidatingTester) CreateDefaultAccount(name common.AccountName) *types.TransactionTrace {
	includeCode := true
	creator := eosio
	trx := types.SignedTransaction{}
	vt.SetTransactionHeaders(&trx.Transaction, vt.DefaultExpirationDelta, 0)

	ownerAuth := types.NewAuthority(vt.getPublicKey(name, "owner"), 0)

	activeAuth := types.NewAuthority(vt.getPublicKey(name, "active"), 0)

	sortPermissions := func(auth *types.Authority) {}
	if includeCode {
		try.EosAssert(ownerAuth.Threshold <= math.MaxUint16, nil, "threshold is too high")
		try.EosAssert(uint64(activeAuth.Threshold) <= uint64(math.MaxUint64), nil, "threshold is too high")
		ownerAuth.Accounts = append(ownerAuth.Accounts, types.PermissionLevelWeight{
			Permission: types.PermissionLevel{Actor: name, Permission: common.DefaultConfig.EosioCodeName},
			Weight:     types.WeightType(ownerAuth.Threshold),
		})
		sortPermissions(&ownerAuth)
		activeAuth.Accounts = append(activeAuth.Accounts, types.PermissionLevelWeight{
			Permission: types.PermissionLevel{Actor: name, Permission: common.DefaultConfig.EosioCodeName},
			Weight:     types.WeightType(activeAuth.Threshold),
		})
		sortPermissions(&activeAuth)
	}
	new := NewAccount{
		Creator: creator,
		Name:    name,
		Owner:   ownerAuth,
		Active:  activeAuth,
	}
	data, _ := rlp.EncodeToBytes(new)
	act := &types.Action{
		Account:       new.GetAccount(),
		Name:          new.GetName(),
		Authorization: []types.PermissionLevel{{creator, common.DefaultConfig.ActiveName}},
		Data:          data,
	}
	trx.Actions = append(trx.Actions, act)

	vt.SetTransactionHeaders(&trx.Transaction, vt.DefaultExpirationDelta, 0)
	pk := vt.getPrivateKey(creator, "active")
	chainId := vt.Control.GetChainId()
	trx.Sign(&pk, &chainId)
	return vt.PushTransaction(&trx, common.MaxTimePoint(), vt.DefaultBilledCpuTimeUs)
}
