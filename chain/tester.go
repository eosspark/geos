package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"math"
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
	ChainTransactions       map[common.BlockIdType]types.TransactionReceipt
	LastProducedBlock       map[common.AccountName]common.BlockIdType
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
	t.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	//t.Control.AcceptedBlock.Connect() // TODO: Control.signal
}

func (t BaseTester) close() {
	t.Control.Close()
	t.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
}

func (t BaseTester) IsSameChain(other *BaseTester) bool {
	return t.Control.HeadBlockId() == other.Control.HeadBlockId()
}

func (t BaseTester) PushBlock(b *types.SignedBlock) *types.SignedBlock {
	t.Control.AbortBlock()
	t.Control.PushBlock(b, types.Complete)
	return &types.SignedBlock{}
}

func (t BaseTester) pushGenesisBlock() {
	//t.setCode()
	//t.setAbi()
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

	if common.Empty(t.Control.PendingBlockState()) || t.Control.PendingBlockState().Header.Timestamp != types.BlockTimeStamp(nextTime) {
		t.startBlock(nextTime)
	}
	Hbs := t.Control.HeadBlockState()
	producer := Hbs.GetScheduledProducer(types.BlockTimeStamp(nextTime))
	privKey := ecc.PrivateKey{}
	privateKey, ok := t.BlockSigningPrivateKeys[producer.BlockSigningKey]
	if !ok {
		privKey = t.getPrivateKey(producer.ProducerName, "active")
	} else {
		privKey = privateKey
	}

	if !skipPendingTrxs {
		unappliedTrxs := t.Control.GetUnappliedTransactions()
		for _, trx := range unappliedTrxs {
			trace := t.Control.pushTransaction(trx, common.MaxTimePoint(), 0, false)
			if !common.Empty(trace.Except) {
				try.EosThrow(trace.Except, "tester produceBlock is error:%#v", trace.Except)
			}
		}

		scheduledTrxs := t.Control.GetScheduledTransactions()
		for len(scheduledTrxs) > 0 {
			for _, trx := range scheduledTrxs {
				trace := t.Control.pushScheduledTransactionById(&trx, common.MaxTimePoint(), 0, false)
				if !common.Empty(trace.Except) {
					try.EosThrow(trace.Except, "tester produceBlock is error:%#v", trace.Except)
				}
			}
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
	t.startBlock(nextTime + common.TimePoint(common.Seconds(common.DefaultConfig.BlockIntervalUs)))
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
		traces[i] = t.createAccount(n, common.DefaultConfig.SystemAccountName, multiSig, includeCode)
	}
	return traces
}

func (t BaseTester) createAccount(name common.AccountName, creator common.AccountName, multiSig bool, includeCode bool) *types.TransactionTrace {
	trx := types.SignedTransaction{}
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0) //TODO: test
	ownerAuth := types.Authority{}
	if multiSig {
		ownerAuth = types.Authority{Threshold: 2, Keys: []types.KeyWeight{{t.getPublicKey(name, "owner"), 1}}}
	} else {
		ownerAuth = types.NewAuthority(t.getPublicKey(name, "owner"), 0)
	}
	activeAuth := types.NewAuthority(t.getPublicKey(name, "active"), 0)

	sortPermissions := func(auth *types.Authority) {

	}
	if includeCode { //TODO
		try.EosAssert(ownerAuth.Threshold <= math.MaxUint16, nil, "threshold is too high")
		try.EosAssert(uint64(activeAuth.Threshold) <= uint64(math.MaxUint64), nil, "threshold is too high")
		//ownerAuth.Accounts
		sortPermissions(&ownerAuth)
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
	_, r := false, (*types.TransactionTrace)(nil)
	try.Try(func() {
		if t.Control.PendingBlockState() != nil {
			t.startBlock(t.Control.HeadBlockTime() + common.TimePoint(common.Seconds(common.DefaultConfig.BlockIntervalUs)))
		}
		mtrx := types.TransactionMetadata{}
		mtrx.Trx = trx
		trace = t.Control.pushTransaction(&mtrx, deadline, billedCpuTimeUs, true)
		if trace.ExceptPtr != nil {
			try.EosThrow(trace.ExceptPtr, "tester PushTransaction is error :%#v", trace.ExceptPtr)
		}
		if trace.Except != nil {
			try.EosThrow(trace.Except, "tester PushTransaction is error :%#v", trace.Except)
		}
		r = trace
		return
	}).FcCaptureAndRethrow().End()
	return r
}

func (t BaseTester) PushAction(act *types.Action, authorizer common.AccountName) {
	trx := types.SignedTransaction{}
	if !common.Empty(authorizer) {
		act.Authorization = []types.PermissionLevel{{authorizer, common.DefaultConfig.ActiveName}}
	}
	trx.Actions = append(trx.Actions, act)
	t.SetTransactionHeaders(&trx.Transaction, 0, 0) //TODO
	if common.Empty(authorizer) {
		chainId := t.Control.GetChainId()
		privateKey := t.getPrivateKey(authorizer, "active")
		trx.Sign(&privateKey, &chainId)
	}
	try.Try(func() {
		t.PushTransaction(&trx, 0, 0) //TODO
	}).Catch(func(ex exception.Exception) {
		log.Error("tester PushAction is error: %#v", ex.Message())
	}).End()
	t.ProduceBlock(common.Microseconds(common.DefaultConfig.BlockIntervalMs), 0)
	/*BOOST_REQUIRE_EQUAL(true, chain_has_transaction(trx.id()))
	success()*/
	return
}

type VariantsObject []map[string]interface{}

func (t BaseTester) PushAction2(code *common.AccountName, acttype *common.AccountName,
	actor common.AccountName, data *VariantsObject, expiration uint32, delaySec uint32) *types.TransactionTrace {
	auths := make([]types.PermissionLevel, 0)
	auths = append(auths, types.PermissionLevel{actor, common.DefaultConfig.ActiveName})
	return t.PushAction4(code, acttype, &auths, data, expiration, delaySec)
}

func (t BaseTester) PushAction3(code *common.AccountName, acttype *common.AccountName,
	actors *[]common.AccountName, data *VariantsObject, expiration uint32, delaySec uint32) *types.TransactionTrace {
	auths := make([]types.PermissionLevel, 0)
	for _, actor := range auths {
		auths = append(auths, actor)
	}
	return t.PushAction4(code, acttype, &auths, data, expiration, delaySec)
}

func (t BaseTester) PushAction4(code *common.AccountName, acttype *common.AccountName,
	actors *[]types.PermissionLevel, data *VariantsObject, expiration uint32, delaySec uint32) *types.TransactionTrace {
	try.Try(func() {
		trx := types.SignedTransaction{}
		action := t.GetAction(*code, *acttype, *actors, data)
		trx.Actions = append(trx.Actions, action)
	})
	return nil
}
func (t BaseTester) GetAction(code common.AccountName, actType common.AccountName,
	auths []types.PermissionLevel, data *VariantsObject) *types.Action {
	/*acnt := t.Control.GetAccount(code)
	abi := acnt.GetAbi()
	abis := types.AbiSerializer{}
	actionTypeName := abis.getActionType(actType)*/

	return nil
}

func (t BaseTester) getPrivateKey(keyName common.Name, role string) ecc.PrivateKey {
	//TODO: wait for testing
	priKey, _ := ecc.NewPrivateKey(crypto.Hash256(keyName.String() + role).String())
	return *priKey
}

func (t BaseTester) getPublicKey(keyName common.Name, role string) ecc.PublicKey {
	priKey := t.getPrivateKey(keyName, role)
	return priKey.PublicKey()
}

func (t BaseTester) ProduceBlock(skipTime common.Microseconds, skipFlag uint32) *types.SignedBlock {
	return t.produceBlock(skipTime, false, skipFlag)
}

func (t BaseTester) ProduceEmptyBlock(skipTime common.Microseconds, skipFlag uint32) *types.SignedBlock {
	t.Control.AbortBlock()
	return t.produceBlock(skipTime, true, skipFlag)
}

func (t BaseTester) LinkAuthority(account common.AccountName, code common.AccountName, req common.PermissionName, rtype common.ActionName) {
	trx := types.SignedTransaction{}
	link := linkAuth{Account: account, Code: code, Type: rtype, Requirement: req}
	data, _ := rlp.EncodeToBytes(link)
	act := types.Action{Account: link.getName(), Name: link.getName(), Authorization: []types.PermissionLevel{{account, common.DefaultConfig.ActiveName}}, Data: data}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	privKey := t.getPrivateKey(account, "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&privKey, &chainId)
	t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) UnlinkAuthority(account common.AccountName, code common.AccountName, rtype common.ActionName) {
	trx := types.SignedTransaction{}
	unlink := unlinkAuth{Account: account, Code: code, Type: rtype}
	data, _ := rlp.EncodeToBytes(unlink)
	act := types.Action{Account: unlink.getName(), Name: unlink.getName(), Authorization: []types.PermissionLevel{{account, common.DefaultConfig.ActiveName}}, Data: data}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta, 0)
	privKey := t.getPrivateKey(account, "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&privKey, &chainId)
	t.PushTransaction(&trx, common.MaxTimePoint(), t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) SetAuthority(account common.AccountName, perm common.PermissionName, auth types.Authority, parent common.PermissionName, auths *[]types.PermissionLevel, keys *[]ecc.PrivateKey) {
	//TODO: try
	trx := types.SignedTransaction{}
	update := updateAuth{Account: account, Permission: perm, Parent: parent, Auth: auth}
	data, _ := rlp.EncodeToBytes(update)
	act := types.Action{Account: update.getName(), Name: update.getName(), Authorization: *auths, Data: data}
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
	//TODO: try
	trx := types.SignedTransaction{}
	delete := deleteAuth{Account: account, Permission: perm}
	data, _ := rlp.EncodeToBytes(delete)
	act := types.Action{Account: delete.getName(), Name: delete.getName(), Authorization: *auths, Data: data}
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
	setCode := setCode{Account: account, VmType: 0, VmVersion: 0, Code: wasm}
	data, _ := rlp.EncodeToBytes(setCode)
	act := types.Action{Account: setCode.getName(), Name: setCode.getName(), Authorization: []types.PermissionLevel{{account, common.DefaultConfig.ActiveName}}, Data: data}
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

func (t BaseTester) SetAbi(account common.AccountName, abiJson *byte, signer *ecc.PrivateKey) {
	// abi := fc::json::from_string(abi_json).template as<abi_def>()
	abi := types.AbiDef{}
	trx := types.SignedTransaction{}
	abiBytes, _ := rlp.EncodeToBytes(abi)
	setAbi := setAbi{Account: account, Abi: abiBytes}
	data, _ := rlp.EncodeToBytes(setAbi)
	act := types.Action{Account: setAbi.getName(), Name: setAbi.getName(), Authorization: []types.PermissionLevel{{account, common.DefaultConfig.ActiveName}}, Data: data}
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
	val, _ := t.ChainTransactions[*txId]
	return &val
}
