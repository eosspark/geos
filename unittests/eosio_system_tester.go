package unittests

import (
	. "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/log"
	"io/ioutil"
)

type EosioSystemTester struct {
	BaseTester
	abiSer      abi_serializer.AbiSerializer
	tokenAbiSer abi_serializer.AbiSerializer
}

func newEosioSystemTester(pushGenesis bool, readMode DBReadMode) *EosioSystemTester {
	e := &EosioSystemTester{}
	e.DefaultExpirationDelta = 6
	e.DefaultBilledCpuTimeUs = 2000
	e.AbiSerializerMaxTime = 1000 * 1000
	e.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	e.LastProducedBlock = make(map[common.AccountName]common.BlockIdType)

	e.init(pushGenesis, readMode)
	return e
}

func initEosioSystemTester() *EosioSystemTester {
	e := newEosioSystemTester(true, SPECULATIVE)

	e.ProduceBlocks(2, false)
	e.CreateAccounts([]common.AccountName{eosioToken, eosioRam, eosioRamFee, eosioStake,
		eosioBpay, eosioVpay, eosioSaving}, false, true)
	e.ProduceBlocks(100, false)

	//eosio.token
	wasmName := "test_contracts/eosio.token.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	e.SetCode(eosioToken, code, nil)
	abiName := "test_contracts/eosio.token.abi"
	abi, _ := ioutil.ReadFile(abiName)
	e.SetAbi(eosioToken, abi, nil)
	accnt := entity.AccountObject{Name: eosioToken}
	e.Control.DB.Find("byName", accnt, &accnt)
	abiDef := abi_serializer.AbiDef{}
	if !abi_serializer.ToABI(accnt.Abi, &abiDef) {
		log.Error("eosio_system_tester::initEosioSystemTester failed with ToAbi")
	}
	//TODO
	//e.tokenAbiSer.SetAbi(&abiDef,&e.AbiSerializerMaxTime)
	e.CreateCurrency(eosioToken, common.DefaultConfig.SystemAccountName, CoreFromString("10000000000.0000"))
	e.Issue(common.DefaultConfig.SystemAccountName, CoreFromString("1000000000.0000"), common.DefaultConfig.SystemAccountName)
	currencyBalance := e.GetBalance(eosio)
	expectedBalance := CoreFromString("1000000000.0000")
	if currencyBalance != expectedBalance {
		log.Error("error, initEosioSystemTester failed")
	}

	//eosio.system
	wasmName = "test_contracts/eosio.system.wasm"
	code, _ = ioutil.ReadFile(wasmName)
	e.SetCode(eosio, code, nil)
	abiName = "test_contracts/eosio.system.abi"
	abi, _ = ioutil.ReadFile(abiName)
	e.SetAbi(eosio, abi, nil)

	accnt = entity.AccountObject{Name: eosio}
	e.Control.DB.Find("byName", accnt, &accnt)
	abiDef = abi_serializer.AbiDef{}
	if !abi_serializer.ToABI(accnt.Abi, &abiDef) {
		log.Error("eosio_system_tester::initEosioSystemTester failed with ToAbi")
	}
	//TODO
	//e.abiSer.SetAbi(&abiDef,&e.AbiSerializerMaxTime)

	e.ProduceBlocks(1, false)

	e.CreateAccountWithResources(
		alice,
		common.DefaultConfig.SystemAccountName,
		CoreFromString("1.0000"),
		false,
		CoreFromString("10.0000"),
		CoreFromString("10.0000"),
	)
	e.CreateAccountWithResources(
		bob,
		common.DefaultConfig.SystemAccountName,
		CoreFromString("0.4500"),
		false,
		CoreFromString("10.0000"),
		CoreFromString("10.0000"),
	)
	e.CreateAccountWithResources(
		carol,
		common.DefaultConfig.SystemAccountName,
		CoreFromString("1.0000"),
		false,
		CoreFromString("10.0000"),
		CoreFromString("10.0000"),
	)
	if CoreFromString("1000000000.0000").Amount != e.GetBalance(eosio).Amount+e.GetBalance(eosioRamFee).Amount+
		e.GetBalance(eosioStake).Amount+e.GetBalance(eosioRam).Amount {
		log.Error("error")
	}
	return e
}

func (e EosioSystemTester) CreateAccountWithResources(name common.AccountName, creator common.AccountName,
	ramFunds common.Asset, multiSig bool, net common.Asset, cpu common.Asset) *types.TransactionTrace {
	trx := types.SignedTransaction{}
	e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)
	ownerAuth := types.Authority{}
	if multiSig {
		ownerAuth = types.Authority{
			Threshold: 2,
			Keys:      []types.KeyWeight{{Key: e.getPublicKey(name, "owner"), Weight: 1}},
			Accounts:  []types.PermissionLevelWeight{{Permission: types.PermissionLevel{Actor: creator, Permission: common.DefaultConfig.ActiveName}, Weight: 1}},
		}
	} else {
		ownerAuth = types.NewAuthority(e.getPublicKey(name, "owner"), 0)
	}
	activeAuth := types.NewAuthority(e.getPublicKey(name, "active"), 0)

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

	buyRamData := common.Variants{
		"payer":    creator,
		"receiver": name,
		"quant":    ramFunds,
	}
	buyRam := e.GetAction(
		eosio,
		common.N("buyram"),
		[]types.PermissionLevel{{creator, common.DefaultConfig.ActiveName}},
		&buyRamData,
	)
	trx.Actions = append(trx.Actions, buyRam)

	delegateData := common.Variants{
		"from":               creator,
		"receiver":           name,
		"stake_net_quantity": net,
		"stake_cpu_quantity": cpu,
		"transfer":           0,
	}
	delegate := e.GetAction(
		eosio,
		common.N("delegatebw"),
		[]types.PermissionLevel{{creator, common.DefaultConfig.ActiveName}},
		&delegateData,
	)
	trx.Actions = append(trx.Actions, delegate)

	e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)
	pk := e.getPrivateKey(creator, "active")
	chainId := e.Control.GetChainId()
	trx.Sign(&pk, &chainId)
	return e.PushTransaction(&trx, common.MaxTimePoint(), e.DefaultBilledCpuTimeUs)
}

func (e EosioSystemTester) CreateAccountWithResources2(name common.AccountName, creator common.AccountName, ramBytes uint32) *types.TransactionTrace {
	trx := types.SignedTransaction{}
	e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)

	ownerAuth := types.NewAuthority(e.getPublicKey(name, "owner"), 0)
	activeAuth := types.NewAuthority(e.getPublicKey(name, "active"), 0)

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

	buyRamBytesData := common.Variants{
		"payer":    creator,
		"receiver": name,
		"bytes":    ramBytes,
	}
	buyRam := e.GetAction(
		eosio,
		common.N("buyrambytes"),
		[]types.PermissionLevel{{creator, common.DefaultConfig.ActiveName}},
		&buyRamBytesData,
	)
	trx.Actions = append(trx.Actions, buyRam)

	delegateData := common.Variants{
		"from":               creator,
		"receiver":           name,
		"stake_net_quantity": CoreFromString("10.0000"),
		"stake_cpu_quantity": CoreFromString("10.0000"),
		"transfer":           0,
	}
	delegate := e.GetAction(
		eosio,
		common.N("delegatebw"),
		[]types.PermissionLevel{{creator, common.DefaultConfig.ActiveName}},
		&delegateData,
	)
	trx.Actions = append(trx.Actions, delegate)

	e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)
	pk := e.getPrivateKey(creator, "active")
	chainId := e.Control.GetChainId()
	trx.Sign(&pk, &chainId)
	return e.PushTransaction(&trx, common.MaxTimePoint(), e.DefaultBilledCpuTimeUs)
}

func (e EosioSystemTester) BuyRam(payer common.AccountName, receiver common.AccountName, eosin common.Asset) ActionResult {
	buyRam := common.Variants{
		"payer":    payer,
		"receiver": receiver,
		"quant":    eosin,
	}
	act := common.N("buyram")
	return e.EsPushAction(&payer, &act, &buyRam, true)
}

func (e EosioSystemTester) BuyRamBytes(payer common.AccountName, receiver common.AccountName, numBytes uint32) ActionResult {
	buyRamBytes := common.Variants{
		"payer":    payer,
		"receiver": receiver,
		"bytes":    numBytes,
	}
	act := common.N("buyrambytes")
	return e.EsPushAction(&payer, &act, &buyRamBytes, true)
}

func (e EosioSystemTester) SellRam(account common.AccountName, numBytes uint64) ActionResult{
	sellRam := common.Variants{
		"account": account,
		"bytes":   numBytes,
	}
	act := common.N("sellram")
	return e.EsPushAction(&account, &act, &sellRam, true)
}

func (e EosioSystemTester) EsPushAction(signer *common.AccountName, name *common.ActionName, data *common.Variants, auth bool) ActionResult {
	var authorizer common.AccountName
	if auth == true {
		authorizer = *signer
	} else {
		if *signer == bob {
			authorizer = alice
		} else {
			authorizer = bob
		}
	}
	act := e.GetAction(eosio, *name, []types.PermissionLevel{}, data)
	return e.PushAction(act, common.AccountName(authorizer))
}

func (e EosioSystemTester) Stake(from common.AccountName, to common.AccountName, net common.Asset, cpu common.Asset) ActionResult {
	stake := common.Variants{
		"from":               from,
		"receiver":           to,
		"stake_net_quantity": net,
		"stake_cpu_quantity": cpu,
		"transfer":           0,
	}
	act := common.N("delegatebw")
	return e.EsPushAction(&from, &act, &stake, true)
}

func (e EosioSystemTester) GetBalance(act common.AccountName) common.Asset {
	PrimaryKey := uint64(CORE_SYMBOL.ToSymbolCode())
	data := e.GetRowByAccount(uint64(eosioToken), uint64(act), uint64(common.N("accounts")), PrimaryKey)
	if len(data) == 0 {
		return common.Asset{Amount: 0, Symbol: CORE_SYMBOL}
	} else {
		asset := common.Asset{}
		rlp.DecodeBytes(data, &asset)
		return asset
	}
}

func (e EosioSystemTester) GetTotalStake(act uint64) common.Variants {
	type UserResources struct {
		Owner     common.AccountName
		NetWeight common.Asset
		CpuWeight common.Asset
		RamBytes  int64
	}
	data := e.GetRowByAccount(uint64(eosio), uint64(act), uint64(common.N("userres")), act)
	if len(data) == 0 {
		return common.Variants{}
	} else {
		res := UserResources{}
		rlp.DecodeBytes(data, &res)
		return common.Variants{
			"owner":      res.Owner,
			"net_weight": res.NetWeight,
			"cpu_weight": res.CpuWeight,
			"ram_bytes":  res.RamBytes,
		}
	}
}

func (e EosioSystemTester) CreateCurrency(contract common.Name, manager common.Name, maxSupply common.Asset) {
	act := common.Variants{
		"issuer":         manager,
		"maximum_supply": maxSupply,
	}
	acttype := common.N("create")
	e.PushAction2(
		&contract,
		&acttype,
		contract,
		&act,
		e.DefaultExpirationDelta,
		0,
	)
}

func (e EosioSystemTester) Issue(to common.Name, amount common.Asset, manager common.Name) {
	act := common.Variants{
		"to":       to,
		"quantity": amount,
		"memo":     "issue",
	}
	acttype := common.N("issue")
	contract := eosioToken
	e.PushAction2(
		&contract,
		&acttype,
		manager,
		&act,
		e.DefaultExpirationDelta,
		0,
	)
}

func (e EosioSystemTester) Transfer(from common.Name, to common.Name, amount common.Asset, manager common.Name){
	act := common.Variants{
		"from":     from,
		"to":       to,
		"quantity": amount,
		"memo":     "transfer",
	}
	acttype := common.N("transfer")
	contract := eosioToken
	e.PushAction2(
		&contract,
		&acttype,
		manager,
		&act,
		e.DefaultExpirationDelta,
		0,
	)
}