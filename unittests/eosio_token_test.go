package unittests

import (
	//"bytes"
	//"encoding/hex"
	"github.com/docker/docker/pkg/testutil/assert"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/abi_serializer"
	//"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"io/ioutil"
	"strings"
	"testing"
	//"unsafe"
)

type EosioTokenTester struct {
	BaseTester
	abiSer abi_serializer.AbiSerializer
}

func newEosioTokenTester(pushGenesis bool, readMode chain.DBReadMode) *EosioTokenTester {
	e := &EosioTokenTester{}
	e.DefaultExpirationDelta = 6
	e.DefaultBilledCpuTimeUs = 2000
	e.AbiSerializerMaxTime = 1000 * 1000
	e.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	e.LastProducedBlock = make(map[common.AccountName]common.BlockIdType)

	e.init(pushGenesis, readMode)

	return e
}

func initEosioTokenTester() *EosioTokenTester {
	e := newEosioTokenTester(true, chain.SPECULATIVE)
	e.ProduceBlocks(2, false)
	e.CreateAccounts([]common.AccountName{
		common.N("alice"),
		common.N("bob"),
		common.N("carol"),
		common.N("eosio.token")}, false, true)
	e.ProduceBlocks(2, false)

	//eosio.token
	wasmName := "test_contracts/eosio.token.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	e.SetCode(common.N("eosio.token"), code, nil)
	abiName := "test_contracts/eosio.token.abi"
	abi, _ := ioutil.ReadFile(abiName)
	e.SetAbi(common.N("eosio.token"), abi, nil)
	accnt := entity.AccountObject{Name: common.N("eosio.token")}
	e.Control.DB.Find("byName", accnt, &accnt)
	abiDef := abi_serializer.AbiDef{}
	if !abi_serializer.ToABI(accnt.Abi, &abiDef) {
		log.Error("eosio_token_tester::initEosioTokenTester failed with ToAbi")
	}

	e.abiSer.SetAbi(&abiDef, e.AbiSerializerMaxTime)

	return e
}

func (e *EosioTokenTester) pushAction(signer common.AccountName, name common.ActionName, data *common.Variants) string {

	act := e.GetAction(common.N("eosio.token"), name, []types.PermissionLevel{}, data)

	log.Info("action:%v", act)
	return e.PushAction(act, signer)
}

func (e *EosioTokenTester) getStats(symbolName string) *common.Variants {

	symb := common.Symbol{}
	symbol := symb.FromString(&symbolName)
	//symbolCode := 5131092//symbol.ToSymbolCode()
	symbolCode := symbol.ToSymbolCode()
	bytes := e.GetRowByAccount(uint64(common.N("eosio.token")), uint64(symbolCode), uint64(common.N("stat")), uint64(symbolCode))

	v := e.abiSer.BinaryToVariant("currency_stats", bytes, e.AbiSerializerMaxTime, false)
	return &v

}

func (e *EosioTokenTester) getAccount(acc common.AccountName, symbolName string) *common.Variants {

	symb := common.Symbol{}
	symbol := symb.FromString(&symbolName)
	symbolCode := symbol.ToSymbolCode()
	bytes := e.GetRowByAccount(uint64(common.N("eosio.token")), uint64(acc), uint64(common.N("accounts")), uint64(symbolCode))

	v := e.abiSer.BinaryToVariant("account", bytes, e.AbiSerializerMaxTime, false)
	return &v

}

func (e *EosioTokenTester) create(issuer common.AccountName, maximum_supply common.Asset) string {
	return e.pushAction(
		common.N("eosio.token"),
		common.N("create"),
		&common.Variants{"issuer": issuer, "maximum_supply": maximum_supply})
}

func (e *EosioTokenTester) issue(issuer common.AccountName, to common.AccountName, quantity common.Asset, memo string) string {
	return e.pushAction(
		issuer,
		common.N("issue"),
		&common.Variants{"to": to, "quantity": quantity, "memo": memo})
}

func (e *EosioTokenTester) transfer(from common.AccountName, to common.AccountName, quantity common.Asset, memo string) string {
	return e.pushAction(
		from,
		common.N("transfer"),
		&common.Variants{"from": from, "to": to, "quantity": quantity, "memo": memo})
}

func equal(v1 *common.Variants, v2 *common.Variants) bool {

	for k, v := range *v1 {

		if value, ok := (*v2)[k]; !ok || value != v {
			return false
		}
	}

	return true

}

func inString(s1, s2 string) bool {
	if strings.Index(s1, s2) <= 0 {
		return false
	}

	return true
}

func TestCreate(t *testing.T) {
	eosioToken := initEosioTokenTester()
	symbol := "1000.000 TKN"
	eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))
	stats := eosioToken.getStats("3,TKN")
	obj := common.Variants{
		"supply":     "0.000 TKN",
		"max_supply": "1000.000 TKN",
		"issuer":     "alice"}
	ret := equal(stats, &obj)
	assert.Equal(t, ret, true)
	eosioToken.ProduceBlocks(1, false)
	eosioToken.close()
}

func TestCreateNegativeMaxSupply(t *testing.T) {
	eosioToken := initEosioTokenTester()

	returning := false
	try.Try(func() {
		symbol := "-1000.000 TKN"
		eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))
	}).Catch(func(e exception.Exception) {
		if inString(exception.GetDetailMessage(e), "max-supply must be positive") {
			returning = true
		}
	}).End()
	assert.Equal(t, returning, true)
	eosioToken.close()
}

func TestSymbolAlreadyExists(t *testing.T) {
	eosioToken := initEosioTokenTester()

	symbol := "100 TKN"
	eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))
	stats := eosioToken.getStats("0,TKN")
	obj := common.Variants{
		"supply":     "0 TKN",
		"max_supply": "100 TKN",
		"issuer":     "alice"}
	ret := equal(stats, &obj)
	assert.Equal(t, ret, true)
	eosioToken.ProduceBlocks(1, false)

	returning := false
	try.Try(func() {
		symbol = "100 TKN"
		eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))
	}).Catch(func(e exception.Exception) {
		if inString(exception.GetDetailMessage(e), "token with symbol already exists") {
			returning = true
		}
	}).End()
	assert.Equal(t, returning, true)
	eosioToken.close()
}

func TestCreateMaxSupply(t *testing.T) {
	eosioToken := initEosioTokenTester()
	symbol := "4611686018427387903 TKN"
	eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))
	stats := eosioToken.getStats("0,TKN")
	obj := common.Variants{
		"supply":     "0 TKN",
		"max_supply": "4611686018427387903 TKN",
		"issuer":     "alice"}
	ret := equal(stats, &obj)
	assert.Equal(t, ret, true)
	eosioToken.ProduceBlocks(1, false)

	max := common.Asset{10, common.Symbol{0, "NKT"}}
	max.Amount = 4611686018427387904
	returning := false
	try.Try(func() {
		eosioToken.create(common.N("alice"), max)
	}).Catch(func(e exception.Exception) {
		if inString(exception.GetDetailMessage(e), "magnitude of asset amount must be less than 2^62") {
			returning = true
		}
	}).End()
	assert.Equal(t, returning, true)
	eosioToken.close()
}

func TestCreateMaxDecimals(t *testing.T) {
	eosioToken := initEosioTokenTester()
	symbol := "1.000000000000000000 TKN"
	eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))
	stats := eosioToken.getStats("18,TKN")
	obj := common.Variants{
		"supply":     "0.000000000000000000 TKN",
		"max_supply": "1.000000000000000000 TKN",
		"issuer":     "alice"}
	ret := equal(stats, &obj)
	assert.Equal(t, ret, true)
	eosioToken.ProduceBlocks(1, false)

	//max := common.Asset{10, common.Symbol{0, "NKT"}}

	//bytesDest := (*byte)(unsafe.Pointer((&max.Amount))
	//bytesSource := hex.DecodeString("8ac7230489e80000")
	//
	//returning := false
	//try.Try(func() {
	//	eosioToken.create(common.N("alice"), max)
	//}).Catch(func(e exception.Exception) {
	//	if inString(exception.GetDetailMessage(e), "magnitude of asset amount must be less than 2^62") {
	//		returning = true
	//	}
	//}).End()
	//assert.Equal(t, returning, true)
	eosioToken.close()
}

func TestIssue(t *testing.T) {
	eosioToken := initEosioTokenTester()
	symbol := "1000.000 TKN"
	eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))

	quantity := "500.000 TKN"
	eosioToken.issue(common.N("alice"), common.N("alice"), common.Asset{}.FromString(&quantity), "hola")

	stats := eosioToken.getStats("3,TKN")
	obj := common.Variants{
		"supply":     "500.000 TKN",
		"max_supply": "1000.000 TKN",
		"issuer":     "alice"}
	ret := equal(stats, &obj)
	assert.Equal(t, ret, true)

	aliceBalance := eosioToken.getAccount(common.N("alice"), "3,TKN")
	obj = common.Variants{"balance": "500.000 TKN"}
	ret = equal(aliceBalance, &obj)
	assert.Equal(t, ret, true)

	returning := false
	try.Try(func() {
		quantity = "500.001 TKN"
		eosioToken.issue(common.N("alice"), common.N("alice"), common.Asset{}.FromString(&quantity), "hola")
	}).Catch(func(e exception.Exception) {
		if inString(exception.GetDetailMessage(e), "quantity exceeds available supply") {
			returning = true
		}
	}).End()
	assert.Equal(t, returning, true)

	returning = false
	try.Try(func() {
		quantity = "-1.000 TKN"
		eosioToken.issue(common.N("alice"), common.N("alice"), common.Asset{}.FromString(&quantity), "hola")
	}).Catch(func(e exception.Exception) {
		if inString(exception.GetDetailMessage(e), "must issue positive quantity") {
			returning = true
		}
	}).End()
	assert.Equal(t, returning, true)

	quantity = "1.000 TKN"
	eosioToken.issue(common.N("alice"), common.N("alice"), common.Asset{}.FromString(&quantity), "hola")
	eosioToken.close()

}

func TestTransfer(t *testing.T) {
	eosioToken := initEosioTokenTester()
	symbol := "1000 CERO"
	eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))
	eosioToken.ProduceBlocks(1, false)

	quantity := "1000 CERO"
	eosioToken.issue(common.N("alice"), common.N("alice"), common.Asset{}.FromString(&quantity), "hola")

	stats := eosioToken.getStats("0,CERO")
	obj := common.Variants{
		"supply":     "1000 CERO",
		"max_supply": "1000 CERO",
		"issuer":     "alice"}
	ret := equal(stats, &obj)
	assert.Equal(t, ret, true)

	aliceBalance := eosioToken.getAccount(common.N("alice"), "0,CERO")
	obj = common.Variants{"balance": "1000 CERO"}
	ret = equal(aliceBalance, &obj)
	assert.Equal(t, ret, true)

	quantity = "300 CERO"
	eosioToken.transfer(common.N("alice"), common.N("bob"), common.Asset{}.FromString(&quantity), "hola")

	aliceBalance = eosioToken.getAccount(common.N("alice"), "0,CERO")
	obj = common.Variants{
		"balance":   "700 CERO",
		"frozen":    0,
		"whitelist": 1,
	}
	ret = equal(aliceBalance, &obj)
	assert.Equal(t, ret, true)

	bobBalance := eosioToken.getAccount(common.N("bob"), "0,CERO")
	obj = common.Variants{
		"balance":   "300 CERO",
		"frozen":    0,
		"whitelist": 1,
	}
	ret = equal(bobBalance, &obj)
	assert.Equal(t, ret, true)

	returning := false
	try.Try(func() {
		quantity = "701 CERO"
		eosioToken.transfer(common.N("alice"), common.N("bob"), common.Asset{}.FromString(&quantity), "hola")
	}).Catch(func(e exception.Exception) {
		if inString(exception.GetDetailMessage(e), "overdrawn balance") {
			returning = true
		}
	}).End()
	assert.Equal(t, returning, true)

	returning = false
	try.Try(func() {
		quantity = "-1000 CERO"
		eosioToken.transfer(common.N("alice"), common.N("bob"), common.Asset{}.FromString(&quantity), "hola")
	}).Catch(func(e exception.Exception) {
		if inString(exception.GetDetailMessage(e), "must transfer positive quantity") {
			returning = true
		}
	}).End()
	assert.Equal(t, returning, true)
	eosioToken.close()
}
