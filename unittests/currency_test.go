package unittests

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/pkg/testutil/assert"
	. "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/abi_serializer"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"io/ioutil"
	"testing"
)

var eosioToken = common.AccountName(common.N("eosio.token"))
var DEFAULT_EXPIRATION_DELTA uint32 = 6
var DEFAULT_BILLED_CPU_TIME_US uint32 = 2000

type CurrencyTester struct {
	abiSer           abi_serializer.AbiDef
	eosioToken       string
	validatingTester *ValidatingTester
}

/*func initBaseTester() *BaseTester {
	bt := newBaseTester(true, SPECULATIVE)
	return bt
}*/

func NewCurrencyTester() *CurrencyTester {
	ct := &CurrencyTester{}
	ct.eosioToken = eosioToken.String()
	bt := newValidatingTester(true, SPECULATIVE)
	//ct.abiSer = abi_serializer.NewABI()
	bt.CreateDefaultAccount(common.N(ct.eosioToken))
	ct.validatingTester = bt
	//bt.SetCode2(common.N("eosio"),eosioTokenWast)
	wasmName := "test_contracts/eosio.token.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	bt.SetCode(eosioToken, code, nil)
	abiName := "test_contracts/eosio.token.abi"
	abi, err := ioutil.ReadFile(abiName)
	if err != nil {
		log.Error("pushGenesisBlock is err : %v", err)
	}
	bt.SetAbi(common.AccountName(eosioToken), abi, nil)
	accountName := common.AccountName(common.N(ct.eosioToken))

	createData := common.Variants{
		"issuer":         eosioToken,
		"maximum_supply": "1000000000.0000 EOS",
		"can_freeze":     0,
		"can_recall":     0,
		"can_whitelist":  0,
	}
	cn := common.ActionName(common.N("create"))
	result := ct.PushAction(&accountName, &cn, &createData)
	log.Info("NewCurrencyTester push action issue:%v", result.BlockNum)
	data := common.Variants{
		"to":       eosioToken,
		"quantity": "1000000.0000 EOS",
		"memo":     "test",
	}
	in := common.ActionName(common.N("issue"))
	result = ct.PushAction(&accountName, &in, &data)

	log.Info("NewCurrencyTester push action issue::%v", result.BlockNum)
	signedBlock := bt.DefaultProduceBlock()
	log.Info("NewCurrencyTester produceBlock result:%v", signedBlock.Producer.String())
	ct.validatingTester = bt
	return ct
}

func (c *CurrencyTester) PushAction(signer *common.AccountName, name *common.ActionName, data *common.Variants) *types.TransactionTrace {
	action := types.Action{eosioToken, *name, []types.PermissionLevel{{*signer, common.DefaultConfig.ActiveName}}, nil}
	acnt := c.validatingTester.Control.GetAccount(eosioToken)
	a := acnt.GetAbi()
	buf, _ := json.Marshal(data)
	action.Data, _ = a.EncodeAction(*name, buf)
	trx := types.NewSignedTransactionNil()
	trx.Actions = append(trx.Actions, &action)

	c.validatingTester.SetTransactionHeaders(&trx.Transaction, c.validatingTester.DefaultBilledCpuTimeUs, 0)
	key := c.validatingTester.getPrivateKey(*signer, "active")
	chainID := c.validatingTester.Control.GetChainId()
	trx.Sign(&key, &chainID)
	return c.validatingTester.PushTransaction(trx, common.MaxTimePoint(), DEFAULT_BILLED_CPU_TIME_US)
}

func (ct *CurrencyTester) GetBalance(account *common.AccountName) *common.Asset {
	symbol := common.Symbol{Precision: 4, Symbol: "EOS"}
	actionName := common.N(ct.eosioToken)
	asset := ct.validatingTester.GetCurrencyBalance(&actionName, &symbol, account)
	return &asset
}

func (c *CurrencyTester) Transfer(from *common.AccountName, to *common.AccountName, quantity string, memo string) *types.TransactionTrace {
	q := common.N("transfer")
	data := common.Variants{
		"from":     from,
		"to":       to,
		"quantity": quantity,
		"memo":     memo}
	trace := c.PushAction(from, &q, &data)
	c.validatingTester.DefaultProduceBlock()
	return trace
}

func TestBootstrap(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		s := "1000000.0000 EOS"
		asset := common.Asset{}
		expected := asset.FromString(&s)
		actionName := common.N(ct.eosioToken)
		accountName := common.N(ct.eosioToken)

		actual := ct.validatingTester.GetCurrencyBalance(&actionName, &expected.Symbol, &accountName)
		assert.Equal(t, expected, actual)
		ct.validatingTester.close()
	}).FcLogAndRethrow().End()
}

func TestCurrencyTransfer(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		alice := common.N("alice")
		ct.validatingTester.CreateAccounts([]common.AccountName{common.N("alice")}, false, true)
		accountName := common.AccountName(eosioToken)
		actionName := common.N("transfer")
		data := common.Variants{
			"from":     eosioToken,
			"to":       "alice",
			"quantity": "100.0000 EOS",
			"memo":     "fund Alice"}
		trace := ct.PushAction(&accountName, &actionName, &data)
		ct.validatingTester.DefaultProduceBlock()

		s := "100.0000 EOS"
		expected := common.Asset{}.FromString(&s)
		fmt.Println("tmp code :%s,%v", trace.ID, ct.GetBalance(&alice))
		assert.Equal(t, *ct.GetBalance(&alice), expected)
		ct.validatingTester.close()
	}).FcLogAndRethrow().End()
}

func TestDuplicateTransfer(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		ct.validatingTester.CreateAccounts([]common.AccountName{common.N("alice")}, false, true)
		accountName := common.AccountName(common.N(ct.eosioToken))
		actionName := common.ActionName(common.N("transfer"))
		asset := common.Asset{}
		s := "100.0000 EOS"
		expected := asset.FromString(&s)
		alice := common.N("alice")
		data := common.Variants{
			"from":     eosioToken,
			"to":       alice,
			"quantity": "100.0000 EOS",
			"memo":     "fund Alice"}
		trace := ct.PushAction(&accountName, &actionName, &data)

		try.Try(func() {
			ct.PushAction(&accountName, &actionName, &data)
		}).Catch(func(e error) {
			assert.Error(t, e, "Duplicate transaction")
		})

		ct.validatingTester.DefaultProduceBlock()
		log.Info("tmp code:%d,%v", trace.BlockNum, expected)
		//assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		//assert.Equal(t, ct.GetBalance(&alice), expected)
		ct.validatingTester.close()
	}).FcLogAndRethrow().End()
}

func TestAddTransfer(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		ct.validatingTester.CreateAccounts([]common.AccountName{common.N("alice")}, false, true)
		alice := common.N("alice")
		actionName := common.ActionName(common.N("transfer"))
		asset := common.Asset{}
		s := "100.0000 EOS"
		expected := asset.FromString(&s)
		data := common.Variants{
			"from":     eosioToken,
			"to":       "alice",
			"quantity": s,
			"memo":     "fund Alice"}
		trace := ct.PushAction(&eosioToken, &actionName, &data)
		log.Info("tmp code:%d,%v", trace.BlockNum, expected)
		ct.validatingTester.DefaultProduceBlock() //
		//assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		assert.Equal(t, *ct.GetBalance(&alice), expected)

		asset2 := common.Asset{}
		st := "110.0000 EOS"

		exp := asset2.FromString(&st)
		transferData := common.Variants{
			"from":     eosioToken,
			"to":       alice,
			"quantity": "10.0000 EOS",
			"memo":     "fund Alice"}

		try.Try(func() {
			ct.PushAction(&eosioToken, &actionName, &transferData)
		}).Catch(func(e exception.TxDuplicate) {
			log.Error(e.String())
		})

		ct.validatingTester.DefaultProduceBlock()

		//assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		assert.Equal(t, *ct.GetBalance(&alice), exp)
		ct.validatingTester.close()
	}).FcLogAndRethrow().End()
}

func TestOverspend(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		ct.validatingTester.CreateAccounts([]common.AccountName{common.N("alice"), common.N("bob")}, false, true)
		accountName := common.AccountName(common.N(ct.eosioToken))
		actionName := common.ActionName(common.N("transfer"))
		alice := common.AccountName(common.N("alice"))
		asset := common.Asset{}
		s := "100.0000 EOS"
		expected := asset.FromString(&s)
		data := common.Variants{
			"from":     ct.eosioToken,
			"to":       "alice",
			"quantity": s,
			"memo":     "fund Alice"}
		trace := ct.PushAction(&accountName, &actionName, &data)
		log.Info("tmp code:tmp code:%d,%v", trace.BlockNum)
		ct.validatingTester.DefaultProduceBlock() //
		//assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		assert.Equal(t, *ct.GetBalance(&alice), expected)

		s2 := "101.0000 EOS"
		//expected2 := asset2.FromString(&s2)
		data2 := common.Variants{
			"from":     "alice",
			"to":       "bob",
			"quantity": s2,
			"memo":     "fund Alice"}
		returning := false
		try.Try(func() {
			ct.PushAction(&alice, &actionName, &data2)
			bob := common.AccountName(common.N("bob"))
			ct.validatingTester.DefaultProduceBlock() //
			tt := "0.0000 EOS"
			assert.Equal(t, *ct.GetBalance(&alice), expected)
			assert.Equal(t, *ct.GetBalance(&bob), asset.FromString(&tt))
		}).Catch(func(e exception.EosioAssertMessageException) {
			if inString(e.DetailMessage(), "overdrawn balance") {
				returning = true
			}
		}).End()
		assert.Equal(t, true, returning)
		ct.validatingTester.close()
	}).FcLogAndRethrow().End()
}

func TestFullspend(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		ct.validatingTester.CreateAccounts([]common.AccountName{common.N("alice"), common.N("bob")}, false, true)
		actionName := common.ActionName(common.N("transfer"))
		alice := common.AccountName(common.N("alice"))
		//zero := "0.0000 EOS"
		val := "100.0000 EOS"
		data := common.Variants{
			"from":     eosioToken,
			"to":       alice,
			"quantity": val,
			"memo":     "all in! Alice"}
		trace := ct.PushAction(&eosioToken, &actionName, &data)
		ct.validatingTester.DefaultProduceBlock()
		//assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		log.Info("trace id:%v", trace)
		//z :=common.Asset{}.FromString(&zero)

		assert.Equal(t, *ct.GetBalance(&alice), common.Asset{}.FromString(&val))
		bob := common.AccountName(common.N("bob"))
		data2 := common.Variants{
			"from":     alice,
			"to":       bob,
			"quantity": "100.0000 EOS",
			"memo":     "all in! Alice"}
		trace2 := ct.PushAction(&alice, &actionName, &data2)
		log.Info("trace2 id:%d", trace2.ID)
		ct.validatingTester.DefaultProduceBlock()

		s := "100.0000 EOS"
		expected := common.Asset{}.FromString(&s)
		assert.Equal(t, *ct.GetBalance(&bob), expected)
		//assert.Equal(t, *ct.GetBalance(&alice), z)
		ct.validatingTester.close()
	}).FcLogAndRethrow().End()
}

func TestSymbol(t *testing.T) {
	{

	}
}

func TestProxy(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		ct.validatingTester.ProduceBlocks(2, true)
		alice := common.N("alice")
		proxy := common.N("proxy")
		ct.validatingTester.CreateAccounts([]common.AccountName{alice, proxy}, false, true)
		ct.validatingTester.DefaultProduceBlock()
		wasmName := "test_contracts/proxy.wasm"
		code, _ := ioutil.ReadFile(wasmName)
		ct.validatingTester.SetCode(proxy, code, nil)
		abiName := "test_contracts/proxy.abi"
		abi, _ := ioutil.ReadFile(abiName)

		ct.validatingTester.SetAbi(proxy, abi, nil)

		{
			act := types.Action{}
			act.Account = proxy
			act.Name = common.N("setowner")
			act.Authorization = []types.PermissionLevel{{common.N("alice"), common.DefaultConfig.ActiveName}}
			data := common.Variants{
				"owner": alice,
				"delay": 10,
			}

			trace := ct.PushAction(&alice, &act.Name, &data)
			fmt.Println(trace.ID)
			ct.validatingTester.DefaultProduceBlock()
			//assert.Equal(t, true, trace.ID)
			ct.validatingTester.close()
		}
	}).FcLogAndRethrow().End()
}
