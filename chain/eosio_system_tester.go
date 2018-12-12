package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/abi_serializer"
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
	e.AbiSerializerMaxTime = 1000*1000
	e.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	e.LastProducedBlock = make(map[common.AccountName]common.BlockIdType)

	e.init(pushGenesis, readMode)
	return e
}

func initEosioSystemTester() {
	e := newEosioSystemTester(true, SPECULATIVE)

	e.ProduceBlocks(2, false)
	e.CreateAccounts([]common.AccountName{common.N("eosio.token"), common.N("eosio.ram"), common.N("eosio.ramfee"), common.N("eosio.stake"),
		common.N("eosio.bpay"), common.N("eosio.vpay"), common.N("eosio.saving")}, false, true)
	e.ProduceBlocks(100, false)

	//eosio.token
	wasmName := "../wasmgo/testdata_context/eosio.token.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	e.SetCode(common.N("eosio.token"), code, nil)
	abiName := "../wasmgo/testdata_context/eosio.token.abi"
	abi, _ := ioutil.ReadFile(abiName)
	e.SetAbi(common.N("eosio.token"), abi, nil)
	accnt := entity.AccountObject{Name: common.N("eosio.token")}
	e.Control.DB.Find("byName", accnt, &accnt)
	abiDef := abi_serializer.AbiDef{}
	if !abi_serializer.ToABI(accnt.Abi,&abiDef) {
		log.Error("eosio_system_tester::initEosioSystemTester failed with ToAbi")
	}
	//TODO
	//e.tokenAbiSer.SetAbi(&abiDef,&e.AbiSerializerMaxTime)
	e.CreateCurrency(common.N("eosio.token"),common.DefaultConfig.SystemAccountName,CoreFromString("10000000000.0000"))
	e.Issue(common.DefaultConfig.SystemAccountName,CoreFromString("10000000000.0000"),common.DefaultConfig.SystemAccountName)
	currencyBalance := e.GetBalance(common.N("eosio"))
	expectedBalance := CoreFromString("10000000000.0000")
	if currencyBalance != expectedBalance {
		log.Error("error, initEosioSystemTester failed")
	}

	//eosio.system
	wasmName = "../wasmgo/testdata_context/eosio.system.wasm"
	code, _ = ioutil.ReadFile(wasmName)
	e.SetCode(common.N("eosio"), code, nil)
	abiName = "../wasmgo/testdata_context/eosio.system.abi"
	abi, _ = ioutil.ReadFile(abiName)
	e.SetAbi(common.N("eosio"), abi, nil)

	accnt = entity.AccountObject{Name: common.N("eosio")}
	e.Control.DB.Find("byName", accnt, &accnt)
	abiDef = abi_serializer.AbiDef{}
	if !abi_serializer.ToABI(accnt.Abi,&abiDef) {
		log.Error("eosio_system_tester::initEosioSystemTester failed with ToAbi")
	}
	//TODO
	//e.abiSer.SetAbi(&abiDef,&e.AbiSerializerMaxTime)

	e.ProduceBlocks(1,false)

	e.CreateAccountWithResources(
		common.N("alice1111111"),
		common.DefaultConfig.SystemAccountName,
		CoreFromString("1.0000"),
		false,
		CoreFromString("10.0000"),
		CoreFromString("10.0000"),
	)
	e.CreateAccountWithResources(
		common.N("bob111111111"),
		common.DefaultConfig.SystemAccountName,
		CoreFromString("0.4500"),
		false,
		CoreFromString("10.0000"),
		CoreFromString("10.0000"),
	)
	e.CreateAccountWithResources(
		common.N("carol1111111"),
		common.DefaultConfig.SystemAccountName,
		CoreFromString("1.0000"),
		false,
		CoreFromString("10.0000"),
		CoreFromString("10.0000"),
	)
	if CoreFromString("10000000000.0000").Amount != e.GetBalance(common.N("eosio")).Amount + e.GetBalance(common.N("eosio.ramfee")).Amount +
		e.GetBalance(common.N("eosio.stake")).Amount + e.GetBalance(common.N("eosio.ram")).Amount {
				log.Error("error")
	}
	e.close()
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

	new := newAccount{
		Creator: creator,
		Name:    name,
		Owner:   ownerAuth,
		Active:  activeAuth,
	}
	data, _ := rlp.EncodeToBytes(new)
	act := &types.Action{
		Account:       new.getAccount(),
		Name:          new.getName(),
		Authorization: []types.PermissionLevel{{creator, common.DefaultConfig.ActiveName}},
		Data:          data,
	}
	trx.Actions = append(trx.Actions, act)

	buyRamData := VariantsObject{
		"payer":    creator,
		"receiver": name,
		"quant":    ramFunds,
	}
	buyRam := e.GetAction(
		common.N("eosio"),
		common.N("buyram"),
		[]types.PermissionLevel{{creator,common.DefaultConfig.ActiveName}},
		&buyRamData,
	)
	trx.Actions = append(trx.Actions, buyRam)

	delegateData := VariantsObject{
		"from":               creator,
		"receiver":           name,
		"stake_net_quantity": net,
		"stake_cpu_quantity": cpu,
		"transfer":           0,
	}
	delegate := e.GetAction(
		common.N("eosio"),
		common.N("delegatebw"),
		[]types.PermissionLevel{{creator,common.DefaultConfig.ActiveName}},
		&delegateData,
	)
	trx.Actions = append(trx.Actions, delegate)


	e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)
	pk := e.getPrivateKey(creator, "active")
	chainId := e.Control.GetChainId()
	trx.Sign(&pk, &chainId)
	return e.PushTransaction(&trx, common.MaxTimePoint(), e.DefaultBilledCpuTimeUs)
}

func (e EosioSystemTester) GetBalance(act common.AccountName) common.Asset {
	a := common.AccountName(5462355)
	data := e.GetRowByAccount(uint64(common.N("eosio.token")), uint64(act), uint64(common.N("accounts")),&a)
	if len(data) == 0 {
		return common.Asset{Amount:0,Symbol:CORE_SYMBOL}
	} else {
		asset := common.Asset{}
		rlp.DecodeBytes(data, &asset)
		return asset
	}
}

func (e EosioSystemTester) CreateCurrency(contract common.Name, manager common.Name, maxSupply common.Asset) {
	act := VariantsObject{
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
	act := VariantsObject{
		"to":       to,
		"quantity": amount,
		"memo":     "",
	}
	acttype := common.N("issue")
	contract := common.N("eosio.token")
	e.PushAction2(
		&contract,
		&acttype,
		manager,
		&act,
		e.DefaultExpirationDelta,
		0,
	)
}