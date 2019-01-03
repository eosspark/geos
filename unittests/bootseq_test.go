package unittests

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

type genesisAccount struct {
	Aname          common.AccountName
	InitialBlanace uint64
}

var testGenesis = []genesisAccount{
	{common.N("b1"), 1000000000000},
	{common.N("whale4"), 400000000000},
	{common.N("whale3"), 300000000000},
	{common.N("whale2"), 200000000000},
	{common.N("proda"), 10000000000},
	{common.N("prodb"), 10000000000},
	{common.N("prodc"), 10000000000},
	{common.N("prodd"), 10000000000},
	{common.N("prode"), 10000000000},
	{common.N("prodf"), 10000000000},
	{common.N("prodg"), 10000000000},
	{common.N("prodh"), 10000000000},
	{common.N("prodi"), 10000000000},
	{common.N("prodj"), 10000000000},
	{common.N("prodk"), 10000000000},
	{common.N("prodl"), 10000000000},
	{common.N("prodm"), 10000000000},
	{common.N("prodn"), 10000000000},
	{common.N("prodo"), 10000000000},
	{common.N("prodp"), 10000000000},
	{common.N("prodq"), 10000000000},
	{common.N("prodr"), 10000000000},
	{common.N("prods"), 10000000000},
	{common.N("prodt"), 10000000000},
	{common.N("produ"), 10000000000},
	{common.N("runnerup1"), 10000000000},
	{common.N("runnerup2"), 10000000000},
	{common.N("runnerup3"), 10000000000},
	{common.N("minow1"), 1000000},
	{common.N("minow2"), 10000},
	{common.N("minow3"), 10000},
	{common.N("masses"), 8000000000000},
}

type BootSeqTester struct {
	BaseTester
	abiSer abi_serializer.AbiSerializer
}

func newBootSeqTester(pushGenesis bool, readMode chain.DBReadMode) *BootSeqTester {
	e := &BootSeqTester{}
	e.DefaultExpirationDelta = 6
	e.DefaultBilledCpuTimeUs = 2000
	e.AbiSerializerMaxTime = 1000 * 1000
	e.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	e.LastProducedBlock = make(map[common.AccountName]common.BlockIdType)

	e.init(pushGenesis, readMode)

	return e
}

// func initBootSeqTester() *BootSeqTester {
// 	e := newEosioTokenTester(true, chain.SPECULATIVE)
// 	e.ProduceBlocks(2, false)
// 	e.CreateAccounts([]common.AccountName{
// 		common.N("alice"),
// 		common.N("bob"),
// 		common.N("carol"),
// 		common.N("eosio.token")}, false, true)
// 	e.ProduceBlocks(2, false)

// 	//eosio.token
// 	wasmName := "test_contracts/eosio.token.wasm"
// 	code, _ := ioutil.ReadFile(wasmName)
// 	e.SetCode(common.N("eosio.token"), code, nil)
// 	abiName := "test_contracts/eosio.token.abi"
// 	abi, _ := ioutil.ReadFile(abiName)
// 	e.SetAbi(common.N("eosio.token"), abi, nil)
// 	accnt := entity.AccountObject{Name: common.N("eosio.token")}
// 	e.Control.DB.Find("byName", accnt, &accnt)
// 	abiDef := abi_serializer.AbiDef{}
// 	if !abi_serializer.ToABI(accnt.Abi, &abiDef) {
// 		log.Error("eosio_token_tester::initEosioTokenTester failed with ToAbi")
// 	}

// 	e.abiSer.SetAbi(&abiDef, e.AbiSerializerMaxTime)

// 	return e
// }

func (e *BootSeqTester) getGlobalState() common.Variants {

	data := e.GetRowByAccount(uint64(common.DefaultConfig.SystemAccountName), uint64(common.DefaultConfig.SystemAccountName), uint64(common.N("global")), uint64(common.N("global")))

	return e.abiSer.BinaryToVariant("eosio_global_state", data, e.AbiSerializerMaxTime, false)
}

func (e *BootSeqTester) buyram(payer common.AccountName, receiver common.AccountName, ram common.Asset) *types.TransactionTrace {
	actType := common.N("buyram")
	r := e.PushAction2(&common.DefaultConfig.SystemAccountName, &actType, payer,
		&common.Variants{
			"payer":    payer,
			"receiver": receiver,
			"quant":    ram},
		e.DefaultExpirationDelta,
		0)

	e.ProduceBlocks(1, false)
	return r
}

func (e *BootSeqTester) delegateBandwidth(from common.AccountName, receiver common.AccountName, net common.Asset, cpu common.Asset, transfer uint8) *types.TransactionTrace {
	actType := common.N("delegatebw")
	r := e.PushAction2(&common.DefaultConfig.SystemAccountName, &actType, from,
		&common.Variants{
			"from":               from,
			"receiver":           receiver,
			"stake_net_quantity": net,
			"stake_cpu_quantity": cpu,
			"transfer":           transfer},
		e.DefaultExpirationDelta,
		0)
	e.ProduceBlocks(1, false)
	return r
}

func (e *BootSeqTester) createCurrency(contract common.AccountName, manager common.AccountName, maxsupply common.Asset, signer *ecc.PrivateKey) *types.TransactionTrace {
	actType := common.N("create")
	r := e.PushAction2(&contract, &actType, contract,
		&common.Variants{
			"issuer":         manager,
			"maximum_supply": maxsupply},
		e.DefaultExpirationDelta,
		0)
	return r
}

func (e *BootSeqTester) issue(contract common.AccountName, manager common.AccountName, to common.AccountName, amount common.Asset) *types.TransactionTrace {
	actType := common.N("issue")
	r := e.PushAction2(&contract, &actType, manager,
		&common.Variants{
			"to":       to,
			"quantity": amount,
			"memo":     ""},
		e.DefaultExpirationDelta,
		0)
	e.ProduceBlocks(1, false)
	return r
}

func (e *BootSeqTester) claimRewards(owner common.AccountName) *types.TransactionTrace {
	actType := common.N("claimrewards")
	r := e.PushAction2(&common.DefaultConfig.SystemAccountName, &actType, owner,
		&common.Variants{
			"owner": owner},
		e.DefaultExpirationDelta,
		0)
	e.ProduceBlocks(1, false)
	return r
}

func (e *BootSeqTester) setPrivileged(account common.AccountName) *types.TransactionTrace {
	actType := common.N("setpriv")
	r := e.PushAction2(&common.DefaultConfig.SystemAccountName, &actType, common.DefaultConfig.SystemAccountName,
		&common.Variants{
			"account": account,
			"is_priv": 1},
		e.DefaultExpirationDelta,
		0)
	e.ProduceBlocks(1, false)
	return r
}

func (e *BootSeqTester) registerProducer(producer common.AccountName) *types.TransactionTrace {
	actType := common.N("regproducer")
	r := e.PushAction2(&common.DefaultConfig.SystemAccountName, &actType, producer,
		&common.Variants{
			"producer":     producer,
			"is_priv":      1,
			"producer_key": e.getPublicKey(producer, "active"),
			"url":          "",
			"location":     0},
		e.DefaultExpirationDelta,
		0)
	e.ProduceBlocks(1, false)
	return r
}

func (e *BootSeqTester) undelegateBandwidth(from common.AccountName, receiver common.AccountName, net common.Asset, cpu common.Asset) *types.TransactionTrace {
	actType := common.N("undelegatebw")
	r := e.PushAction2(&common.DefaultConfig.SystemAccountName, &actType, from,
		&common.Variants{
			"from":                 from,
			"receiver":             receiver,
			"unstake_net_quantity": net,
			"unstake_cpu_quantity": cpu},
		e.DefaultExpirationDelta,
		0)
	e.ProduceBlocks(1, false)
	return r
}

func (e *BootSeqTester) getBalance(act common.AccountName) common.Asset {
	account := common.N("eosio.token")
	return e.GetCurrencyBalance(&account, &CORE_SYMBOL, &act)
}

func (e *BootSeqTester) setCodeAbi(account common.AccountName, wasm []byte, abi []byte, signer *ecc.PrivateKey) {

	e.SetCode(account, wasm, signer)
	e.SetAbi(account, abi, signer)

	if account == common.DefaultConfig.SystemAccountName {

		accnt := e.Control.GetAccount(account)
		abiDef := abi_serializer.AbiDef{}
		if !abi_serializer.ToABI(accnt.Abi, &abiDef) {
			//log.Error("eosio_token_tester::initEosioTokenTester failed with ToAbi")
		}
		e.abiSer.SetAbi(&abiDef, e.AbiSerializerMaxTime)
	}
	e.ProduceBlocks(1, false)
}

func TestBootSeq(t *testing.T) {

	b := newBootSeqTester(true, chain.SPECULATIVE)

	b.CreateAccounts([]common.AccountName{
		common.N("eosio.msig"),
		common.N("eosio.token"),
		common.N("eosio.ram"),
		common.N("eosio.ramfee"),
		common.N("eosio.stake"),
		common.N("eosio.vpay"),
		common.N("eosio.bpay"),
		common.N("eosio.saving")},
		false, true)

	wasm, _ := ioutil.ReadFile("test_contracts/eosio_msig.wasm")
	abi, _ := ioutil.ReadFile("test_contracts/eosio_msig.abi")
	b.setCodeAbi(common.N("eosio.msig"), wasm, abi, nil)

	wasm, _ = ioutil.ReadFile("test_contracts/eosio_token.wasm")
	abi, _ = ioutil.ReadFile("test_contracts/eosio_token.abi")
	b.setCodeAbi(common.N("eosio.token"), wasm, abi, nil)

	b.setPrivileged(common.N("eosio.msig"))
	b.setPrivileged(common.N("eosio.token"))

	eosioMsigAcc := b.Control.GetAccount(common.N("eosio.msig"))
	assert.Equal(t, eosioMsigAcc.Privileged, true)

	eosioTokenAcc := b.Control.GetAccount(common.N("eosio.token"))
	assert.Equal(t, eosioTokenAcc.Privileged, true)

}
