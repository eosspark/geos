package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/abi_serializer"
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
	wasmName := "../wasmgo/testdata_context/eosio.token.wasm"
	code, err := ioutil.ReadFile(wasmName)
	if err != nil {
		log.Error("pushGenesisBlock is err : %v", err)
	}
	e.SetCode(common.N("eosio.token"), code, nil)
	abiName := "../wasmgo/testdata_context/eosio.token.abi"
	abi, err := ioutil.ReadFile(abiName)
	if err != nil {
		log.Error("pushGenesisBlock is err : %v", err)
	}
	e.SetAbi(common.N("eosio.token"), abi, nil)
	accnt := entity.AccountObject{Name: common.N("eosio.token")}
	e.Control.DB.Find("byName", accnt, &accnt)
	abiDef := abi_serializer.AbiDef{}
	if !abi_serializer.ToABI(accnt.Abi,&abiDef) {
		log.Error("eosio_system_tester::initEosioSystemTester failed with ToAbi")
	}
	e.tokenAbiSer.SetAbi(&abiDef,&e.AbiSerializerMaxTime)

	e.CreateCurrency(common.N("eosio.token"),common.DefaultConfig.SystemAccountName,CoreFromString("10000000000.0000"))
	e.Issue(common.DefaultConfig.SystemAccountName,CoreFromString("10000000000.0000"),common.DefaultConfig.SystemAccountName)

	e.close()
}

func (e EosioSystemTester) GetBalance(act common.AccountName) common.Asset {
	data := e.GetRowByAccount(uint64(common.N("eosio.token")), uint64(act), uint64(common.N("accounts")),&common.DefaultConfig.SystemAccountName)
	if len(data) == 0 {
		return common.Asset{Amount:0,Symbol:CORE_SYMBOL}
	} else {
		return common.Asset{Amount:10,Symbol:CORE_SYMBOL}
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
		contract,
		&act,
		e.DefaultExpirationDelta,
		0,
	)
}