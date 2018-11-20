package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto"
	"math"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
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

func (t BaseTester) IsSameChain(other *BaseTester) bool {
	return t.Control.HeadBlockId() == other.Control.HeadBlockId()
}

func (t BaseTester) PushBlock(b *types.SignedBlock) *types.SignedBlock {
	t.Control.AbortBlock()
	//t.control.PushBlock(b)
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
	//head := t.Control.HeadBlockState()
	headTime := t.Control.HeadBlockTime()
	nextTime := headTime + common.TimePoint(skipTime)

	if common.Empty(t.Control.PendingBlockState()) || t.Control.PendingBlockState().Header.Timestamp != types.BlockTimeStamp(nextTime) {
		t.startBlock(nextTime)
	}
	Hbs := t.Control.HeadBlockState()
	producer := Hbs.GetScheduledProducer(types.BlockTimeStamp(nextTime))
	privKey := ecc.PrivateKey{}
	privateKey, ok  := t.BlockSigningPrivateKeys[producer.BlockSigningKey]
	if !ok {
		privKey = t.getPrivateKey(producer.ProducerName, "active")
	} else {
		privKey = privateKey
	}

	if !skipPendingTrxs {
		//unappliedTrxs := t.Control.GetUnAppliedTransactions()
		var unappliedTrxs []*types.TransactionMetadata
		for _, trx := range unappliedTrxs {
			trace := t.Control.pushTransaction(trx, common.MaxTimePoint(), 0, false)
			if !common.Empty(trace.Except) {
				//trace.Except.DynamicRethrowException
			}
		}

		//scheduledTrxs := t.Control.GetScheduledTransactions()
		//for len(scheduledTrxs) > 0 {
		//	for _, trx := range scheduledTrxs {
		//		trace := t.Control.pushScheduledTransaction(&trx, common.MaxTimePoint(), 0, false)
		//		if !common.Empty(trace.Except) {
		//			//trace.Except.DynamicRethrowException
		//		}
		//	}
		//}
	}

	t.Control.FinalizeBlock()
	t.Control.SignBlock(func(d common.DigestType) ecc.Signature {
		sign, err := privKey.Sign(d.Bytes())
		if err != nil{
			log.Error(err.Error())
		}
		return sign
	})

	return t.Control.HeadBlockState().SignedBlock
}

func (t BaseTester) startBlock(blockTime common.TimePoint) {

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
	t.SetTransactionHeaders(&trx.Transaction, t.DefaultExpirationDelta,0) //TODO: test
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
		try.EosAssert(ownerAuth.Threshold <= math.MaxUint16, nil,"threshold is too high")
		try.EosAssert(uint64(activeAuth.Threshold) <= uint64(math.MaxUint64), nil,"threshold is too high")
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

//func (t BaseTester) GetAction(code common.AccountName, actType common.AccountName, auths []types.PermissionLevel) types.Action {
//	acnt := t.Control.GetAccount(code)
//	abi := acnt.GetAbi()
//	abis := types.AbiSerializer{}
//	actionTypeName := abis.getActionType(actType)
//}

func (t BaseTester) getPrivateKey(keyName common.Name, role string) ecc.PrivateKey {
	//TODO: wait for testing
	priKey, _ := ecc.NewPrivateKey(crypto.Hash256(keyName.String()+role).String())
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

func (t BaseTester) LinkAuthority(account common.AccountName, code common.AccountName, req common.PermissionName, rtype common.ActionName){
	trx := types.SignedTransaction{}
	link := linkAuth{Account:account, Code:code, Type:rtype, Requirement:req}
	data, _ := rlp.EncodeToBytes(link)
	act := types.Action{Account:link.getName(),Name:link.getName(),Authorization:[]types.PermissionLevel{{account,common.DefaultConfig.ActiveName}}, Data:data}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction,t.DefaultExpirationDelta,0)
	privKey := t.getPrivateKey(account, "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&privKey, &chainId)
	t.PushTransaction(&trx,common.MaxTimePoint(),t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) UnlinkAuthority(account common.AccountName, code common.AccountName, rtype common.ActionName){
	trx := types.SignedTransaction{}
	unlink := unlinkAuth{Account:account, Code:code, Type:rtype}
	data, _ := rlp.EncodeToBytes(unlink)
	act := types.Action{Account:unlink.getName(),Name:unlink.getName(),Authorization:[]types.PermissionLevel{{account,common.DefaultConfig.ActiveName}}, Data:data}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction,t.DefaultExpirationDelta,0)
	privKey := t.getPrivateKey(account, "active")
	chainId := t.Control.GetChainId()
	trx.Sign(&privKey, &chainId)
	t.PushTransaction(&trx,common.MaxTimePoint(),t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) SetAuthority(account common.AccountName, perm common.PermissionName, auth types.Authority, parent common.PermissionName, auths *[]types.PermissionLevel, keys *[]ecc.PrivateKey){
	//TODO: try
	trx := types.SignedTransaction{}
	update := updateAuth{Account:account, Permission:perm, Parent:parent, Auth:auth}
	data, _ := rlp.EncodeToBytes(update)
	act := types.Action{Account:update.getName(),Name:update.getName(),Authorization:*auths, Data:data}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction,t.DefaultExpirationDelta,0)
	chainId := t.Control.GetChainId()
	for _, key := range *keys{
		trx.Sign(&key,&chainId)
	}
	t.PushTransaction(&trx,common.MaxTimePoint(),t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) SetAuthority2(account common.AccountName, perm common.PermissionName, auth types.Authority, parent common.PermissionName){
	permL := types.PermissionLevel{Actor:account,Permission:common.DefaultConfig.OwnerName}
	privKey := t.getPrivateKey(account,"owner")
	t.SetAuthority(account,perm,auth,parent,&[]types.PermissionLevel{permL},&[]ecc.PrivateKey{privKey})
}

func (t BaseTester) DeleteAuthority(account common.AccountName, perm common.PermissionName, auths *[]types.PermissionLevel, keys *[]ecc.PrivateKey){
	//TODO: try
	trx := types.SignedTransaction{}
	delete := deleteAuth{Account:account, Permission:perm}
	data, _ := rlp.EncodeToBytes(delete)
	act := types.Action{Account:delete.getName(),Name:delete.getName(),Authorization:*auths, Data:data}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction,t.DefaultExpirationDelta,0)
	chainId := t.Control.GetChainId()
	for _, key := range *keys{
		trx.Sign(&key,&chainId)
	}
	t.PushTransaction(&trx,common.MaxTimePoint(),t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) DeleteAuthority2(account common.AccountName, perm common.PermissionName){
	permL := types.PermissionLevel{Actor:account,Permission:common.DefaultConfig.OwnerName}
	privKey := t.getPrivateKey(account,"owner")
	t.DeleteAuthority(account,perm,&[]types.PermissionLevel{permL},&[]ecc.PrivateKey{privKey})
}

func (t BaseTester) SetCode(account common.AccountName, wasm []uint8, signer *ecc.PrivateKey){
	trx := types.SignedTransaction{}
	setCode := setCode{Account:account, VmType:0, VmVersion:0, Code:wasm}
	data, _ := rlp.EncodeToBytes(setCode)
	act := types.Action{Account:setCode.getName(),Name:setCode.getName(),Authorization:[]types.PermissionLevel{{account,common.DefaultConfig.ActiveName}}, Data:data}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction,t.DefaultExpirationDelta,0)
	chainId := t.Control.GetChainId()
	if signer != nil {
		trx.Sign(signer,&chainId)
	} else {
		privKey := t.getPrivateKey(account,"active")
		trx.Sign(&privKey,&chainId)
	}
	t.PushTransaction(&trx,common.MaxTimePoint(),t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) SetCode2(account common.AccountName, wast *byte, signer *ecc.PrivateKey){
	//t.SetCode(account, wastToWasm(wast), signer)
}

func (t BaseTester) SetAbi(account common.AccountName, abiJson *byte, signer *ecc.PrivateKey){
	// abi := fc::json::from_string(abi_json).template as<abi_def>()
	abi := types.AbiDef{}
	trx := types.SignedTransaction{}
	abiBytes, _ := rlp.EncodeToBytes(abi)
	setAbi := setAbi{Account:account,Abi:abiBytes}
	data, _ := rlp.EncodeToBytes(setAbi)
	act := types.Action{Account:setAbi.getName(),Name:setAbi.getName(),Authorization:[]types.PermissionLevel{{account,common.DefaultConfig.ActiveName}}, Data:data}
	trx.Actions = append(trx.Actions, &act)
	t.SetTransactionHeaders(&trx.Transaction,t.DefaultExpirationDelta,0)
	chainId := t.Control.GetChainId()
	if signer != nil {
		trx.Sign(signer,&chainId)
	} else {
		privKey := t.getPrivateKey(account,"active")
		trx.Sign(&privKey,&chainId)
	}
	t.PushTransaction(&trx,common.MaxTimePoint(),t.DefaultBilledCpuTimeUs)
}

func (t BaseTester) ChainHasTransaction(txId *common.IdType) bool{
	_, ok := t.ChainTransactions[*txId]
	return ok
}

func (t BaseTester) GetTransactionReceipt(txId *common.IdType) *types.TransactionReceipt{
	val, _ := t.ChainTransactions[*txId]
	return &val
}

func (t BaseTester) GetCurrencyBalance(code *common.AccountName, assetSymbol *common.Symbol, account *common.AccountName) common.Asset{
	db := t.Control.DB
	table := entity.TableIdObject{Code:*code, Scope:*account, Table:common.TableName(*account)}
	err := db.Find("byCodeScopeTable", table, &table)
	result := int64(0)
	if err != nil{
		log.Error("GetCurrencyBalance is error: %s", err)
	} else {
		//TODO
		//obj := entity.KeyValueObject{ID:table.ID,assetSymbol.ToSymbolCode().value}
		obj := entity.KeyValueObject{}
		err := db.Find("byScopePrimary", obj, &obj)
		if err != nil {
			log.Error("GetCurrencyBalance is error: %s", err)
		} else {
			//TODO
		}
	}
	return common.Asset{Amount:result, Symbol:*assetSymbol}
}


