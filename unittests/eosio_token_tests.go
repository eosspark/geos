package unittests

import (
	"bytes"
	"github.com/docker/docker/pkg/testutil/assert"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/abi_serializer"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"io/ioutil"
	"strings"
	"testing"
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

	e.init(pushGenesis, readMode)

	return e
}

func initEosioTokenTester() *EosioSystemTester {
	e := newEosioSystemTester(true, chain.SPECULATIVE)
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
	abiName := "../wasmgo/testdata_context/eosio.token.abi"
	abi, _ := ioutil.ReadFile(abiName)
	e.SetAbi(common.N("eosio.token"), abi, nil)
	accnt := entity.AccountObject{Name: common.N("eosio.token")}
	e.Control.DB.Find("byName", accnt, &accnt)
	abiDef := abi_serializer.AbiDef{}
	if !abi_serializer.ToABI(accnt.Abi, &abiDef) {
		log.Error("eosio_token_tester::initEosioTokenTester failed with ToAbi")
	}

	//abiSer.SetAbi(&abiDef, e.AbiSerializerMaxTime)

	return e
}

func (e EosioTokenTester) pushAction(singer common.AccountName, name common.ActionName, data *VariantsObject) string {
	act := e.GetAction(common.N("eosio"), name, []types.PermissionLevel{}, data)
	return e.PushAction(act, singer)
}

func (e EosioTokenTester) getStats(symbolName string) *VariantsObject {

	//symb := common.Symbol{}
	//symbolCode := symb.FromString(symbolName).ToSymbolCode()
	//bytes := e.GetRowByAccount(common.N("eosio.token"), symbolCode, common.N("stat"), symbolCode)
	return nil

}

func (e EosioTokenTester) getAccount(acc common.AccountName, symbolName string) *VariantsObject {

	//symb := common.Symbol{}
	//symbolCode := symb.FromString(symbolName).ToSymbolCode()
	//bytes := e.GetRowByAccount(common.N("eosio.token"), acc, common.N("accounts"), symbolCode)
	return nil

}

func (e EosioTokenTester) create(issuer common.AccountName, maximum_supply common.Asset) string {
	return e.pushAction(
		issuer,
		common.N("create"),
		&VariantsObject{"issuer": issuer, "maximum_supply": maximum_supply})
}

func (e EosioTokenTester) issue(issuer common.AccountName, to common.AccountName, quantity common.Asset, memo string) string {
	return e.pushAction(
		issuer,
		common.N("issue"),
		&VariantsObject{"to": to, "quantity": quantity, "memo": memo})
}

func (e EosioTokenTester) transfer(from common.AccountName, to common.AccountName, quantity common.Asset, memo string) string {
	return e.pushAction(
		from,
		common.N("transfer"),
		&VariantsObject{"from": from, "to": to, "quantity": quantity, "memo": memo})
}

func equal(v1 *VariantsObject, v2 *VariantsObject) bool {
	b1, _ := rlp.EncodeToBytes(v1)
	b2, _ := rlp.EncodeToBytes(v2)

	if bytes.Compare(b1, b2) == 0 {
		return true
	}

	return false

}

func inString(s1, s2 string) bool {
	if strings.Index(s1, s2) <= 0 {
		return false
	}

	return true
}

func TestCreate(t *testing.T) {
	eosioToken := newEosioTokenTester(true, chain.SPECULATIVE)
	symbol := "1000.000 TKN"
	eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))
	stats := eosioToken.getStats("3,TKN")
	obj := VariantsObject{
		"supply":     "0.000 TKN",
		"max_supply": "1000.000 TKN",
		"issuer":     "alice"}
	ret := equal(stats, &obj)
	assert.Equal(t, ret, true)
	eosioToken.ProduceBlocks(1, false)
}

func TestCreateNegativeMaxSupply(t *testing.T) {
	eosioToken := newEosioTokenTester(true, chain.SPECULATIVE)

	returning := false
	try.Try(func() {
		eosioToken.create(common.N("alice"), common.Asset{}.FromString("-1000.000 TKN"))
	}).Catch(func(e exception.Exception) {
		if inString(exception.GetDetailMessage(e), "max-supply must be positive") {
			returning = true
		}
	}).End()
	assert.Equal(t, returning, true)
}

func TestSymbolAlreadyExists(t *testing.T) {
	eosioToken := newEosioTokenTester(true, chain.SPECULATIVE)

	symbol := "100 TKN"
	eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))
	stats := eosioToken.getStats("0,TKN")
	obj := VariantsObject{
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
}

func TestCreateMaxSupply(t *testing.T) {
	eosioToken := newEosioTokenTester(true, chain.SPECULATIVE)
	symbol := "4611686018427387903 TKN"
	eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))
	stats := eosioToken.getStats("0,TKN")
	obj := VariantsObject{
		"supply":     "0 TKN",
		"max_supply": "4611686018427387903 TKN",
		"issuer":     "alice"}
	ret := equal(stats, &obj)
	assert.Equal(t, ret, true)
	eosioToken.ProduceBlocks(1, false)

	// max := common.Asset{10, common.Symbol{0, "NKT"}}
	// max.Amount = 4611686018427387904
	// returning := false
	// try.Try(func() {
	// 	eosioToken.create(common.N("alice"), max)
	// }).Catch(func(e exception.Exception) {
	// 	if inString(exception.GetDetailMessage(e), "magnitude of asset amount must be less than 2^62") {
	// 		returning = true
	// 	}
	// }).End()
	// assert.Equal(t, returning, true)
}

func TestCreateMaxDecimals(t *testing.T) {
	eosioToken := newEosioTokenTester(true, chain.SPECULATIVE)
	symbol := "1.000000000000000000 TKN"
	eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))
	stats := eosioToken.getStats("18,TKN")
	obj := VariantsObject{
		"supply":     "0.000000000000000000 TKN",
		"max_supply": "1.000000000000000000 TKN",
		"issuer":     "alice"}
	ret := equal(stats, &obj)
	assert.Equal(t, ret, true)
	eosioToken.ProduceBlocks(1, false)

	// max := common.Asset{10, common.Symbol{0, "NKT"}}
	// max.Amount = 0x8ac7230489e80000L
	// returning := false
	// try.Try(func() {
	// 	eosioToken.create(common.N("alice"), max)
	// }).Catch(func(e exception.Exception) {
	// 	if inString(exception.GetDetailMessage(e), "magnitude of asset amount must be less than 2^62") {
	// 		returning = true
	// 	}
	// }).End()
	// assert.Equal(t, returning, true)
}

func TestIssue(t *testing.T) {
	eosioToken := newEosioTokenTester(true, chain.SPECULATIVE)
	symbol := "1000.000 TKN"
	eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))

	quantity := "500.000 TKN"
	eosioToken.issue(common.N("alice"), common.N("alice"), common.Asset{}.FromString(&quantity), "hola")

	stats := eosioToken.getStats("3,TKN")
	obj := VariantsObject{
		"supply":     "500.000 TKN",
		"max_supply": "1000.000 TKN",
		"issuer":     "alice"}
	ret := equal(stats, &obj)
	assert.Equal(t, ret, true)

	aliceBalance := eosioToken.getAccount(common.N("alice"), "3,TKN")
	obj = VariantsObject{"balance": "500.000 TKN"}
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

}

func TestTransfer(t *testing.T) {
	eosioToken := newEosioTokenTester(true, chain.SPECULATIVE)
	symbol := "1000 CERO"
	eosioToken.create(common.N("alice"), common.Asset{}.FromString(&symbol))
	eosioToken.ProduceBlocks(1, false)

	quantity := "1000 CERO"
	eosioToken.issue(common.N("alice"), common.N("alice"), common.Asset{}.FromString(&quantity), "hola")

	stats := eosioToken.getStats("0,CERO")
	obj := VariantsObject{
		"supply":     "1000 CERO",
		"max_supply": "1000 CERO",
		"issuer":     "alice"}
	ret := equal(stats, &obj)
	assert.Equal(t, ret, true)

	aliceBalance := eosioToken.getAccount(common.N("alice"), "0,CERO")
	obj = VariantsObject{"balance": "1000 CERO"}
	ret = equal(aliceBalance, &obj)
	assert.Equal(t, ret, true)

	quantity = "300 CERO"
	eosioToken.issue(common.N("alice"), common.N("bob"), common.Asset{}.FromString(&quantity), "hola")

	aliceBalance = eosioToken.getAccount(common.N("alice"), "0,CERO")
	obj = VariantsObject{
		"balance":   "700 CERO",
		"frozen":    0,
		"whitelist": 1,
	}
	ret = equal(aliceBalance, &obj)
	assert.Equal(t, ret, true)

	bobBalance := eosioToken.getAccount(common.N("bob"), "0,CERO")
	obj = VariantsObject{
		"balance":   "300 CERO",
		"frozen":    0,
		"whitelist": 1,
	}
	ret = equal(bobBalance, &obj)
	assert.Equal(t, ret, true)

	returning := false
	try.Try(func() {
		quantity = "701 CERO"
		eosioToken.issue(common.N("alice"), common.N("bob"), common.Asset{}.FromString(&quantity), "hola")
	}).Catch(func(e exception.Exception) {
		if inString(exception.GetDetailMessage(e), "overdrawn balance") {
			returning = true
		}
	}).End()
	assert.Equal(t, returning, true)

	returning = false
	try.Try(func() {
		quantity = "-1000 CERO"
		eosioToken.issue(common.N("alice"), common.N("bob"), common.Asset{}.FromString(&quantity), "hola")
	}).Catch(func(e exception.Exception) {
		if inString(exception.GetDetailMessage(e), "must transfer positive quantity") {
			returning = true
		}
	}).End()
	assert.Equal(t, returning, true)

}
