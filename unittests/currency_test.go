package unittests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"

	"github.com/stretchr/testify/assert"
)

type CurrencyTester struct {
	abiSer           abi_serializer.AbiDef
	eosioToken       string
	validatingTester *ValidatingTester
}

func NewCurrencyTester() *CurrencyTester {
	ct := &CurrencyTester{}
	ct.eosioToken = eosioToken.String()
	bt := newValidatingTester(true, SPECULATIVE)
	bt.CreateDefaultAccount(common.N(ct.eosioToken))
	ct.validatingTester = bt
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
	action := types.Action{Account: eosioToken, Name: *name, Authorization: []common.PermissionLevel{{Actor: *signer, Permission: common.DefaultConfig.ActiveName}}, Data: nil}
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
	return c.validatingTester.PushTransaction(trx, common.MaxTimePoint(), c.validatingTester.DefaultBilledCpuTimeUs)
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

func (c *CurrencyTester) Issue(to *common.AccountName, quantity string, memo string) *types.TransactionTrace {
	data := common.Variants{
		"to":       to,
		"quantity": quantity,
		"memo":     memo}
	actionName := common.N("issue")
	trace := c.PushAction(&eosioToken, &actionName, &data)
	c.validatingTester.ProduceBlocks(1, false)
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
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
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
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		assert.Equal(t, *ct.GetBalance(&alice), expected)
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
		ct.validatingTester.DefaultProduceBlock() //
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
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

		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
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
		ct.validatingTester.DefaultProduceBlock() //
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
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
		zero := "0.0000 EOS"
		val := "100.0000 EOS"
		data := common.Variants{
			"from":     eosioToken,
			"to":       alice,
			"quantity": val,
			"memo":     "all in! Alice"}
		trace := ct.PushAction(&eosioToken, &actionName, &data)
		ct.validatingTester.DefaultProduceBlock()
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		z := common.Asset{}.FromString(&zero)

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
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace2.ID))
		assert.Equal(t, *ct.GetBalance(&bob), expected)
		assert.Equal(t, *ct.GetBalance(&alice), z)
		ct.validatingTester.close()
	}).FcLogAndRethrow().End()
}

func TestSymbol(t *testing.T) {
	{
		dollar := common.Symbol{Precision: 2, Symbol: "DLLR"}
		sy := "2,DLLR"
		dollar2 := common.Symbol{}.FromString(&sy)
		assert.Equal(t, dollar2, dollar)
		assert.Equal(t, dollar.Decimals(), uint8(2))
		assert.Equal(t, dollar.Name(), "DLLR")
		assert.Equal(t, dollar.Valid(), true)
	}
	{
		def := CORE_SYMBOL
		assert.Equal(t, def.Decimals(), uint8(4))
		assert.Equal(t, def.Name(), CORE_SYMBOL_NAME)
	}
	{
		returning := false
		try.Try(func() {
			sy := ""
			common.Symbol{}.FromString(&sy)
		}).Catch(func(e exception.SymbolTypeException) {
			returning = true
		}).End()
		if returning {
			assert.Equal(t, true, returning)
		}
	}
	{
		returning := false
		try.Try(func() {
			sy := "RND"
			common.Symbol{}.FromString(&sy)
		}).Catch(func(e exception.SymbolTypeException) {
			returning = true
		}).End()
		assert.Equal(t, true, returning)
	}
	{
		returning := false
		try.Try(func() {
			sy := "6,EoS"
			common.Symbol{}.FromString(&sy)
		}).Catch(func(e exception.SymbolTypeException) {

			returning = true
		}).End()
		assert.Equal(t, true, returning)
	}
	{
		str := "10 CUR"
		asset := common.Asset{}.FromString(&str)
		assert.Equal(t, asset.Amount, int64(10))
		assert.Equal(t, asset.Decimals(), uint8(0))
		assert.Equal(t, asset.Symbol.Symbol, "CUR")
	}
	{
		returning := false
		try.Try(func() {
			str := "10CUR"
			common.Asset{}.FromString(&str)
		}).Catch(func(e exception.AssetTypeException) {
			returning = true
		}).End()
		assert.Equal(t, true, returning)
	}
	{
		returning := false
		try.Try(func() {
			str := "10. CUR"
			common.Asset{}.FromString(&str)
		}).Catch(func(e exception.AssetTypeException) {
			returning = true
		}).End()
		assert.Equal(t, true, returning)
	}
	{
		returning := false
		try.Try(func() {
			str := "10"
			common.Asset{}.FromString(&str)
		}).Catch(func(e exception.AssetTypeException) {
			returning = true
		}).End()
		assert.Equal(t, true, returning)
	}
	{
		str := "-001000000.00010 CUR"
		asset := common.Asset{}.FromString(&str)
		assert.Equal(t, asset.Amount, int64(-100000000010))
		assert.Equal(t, asset.Decimals(), uint8(5))
		assert.Equal(t, asset.Symbol.Symbol, "CUR")
		assert.Equal(t, asset.String(), "-1000000.00010 CUR")
	}
	{
		str := "-000000000.00100 CUR"
		asset := common.Asset{}.FromString(&str)
		assert.Equal(t, asset.Amount, int64(-100))
		assert.Equal(t, asset.Decimals(), uint8(5))
		assert.Equal(t, asset.Symbol.Symbol, "CUR")
		assert.Equal(t, asset.String(), "-0.00100 CUR")
	}

	{
		str := "-0.0001 PPP"
		asset := common.Asset{}.FromString(&str)
		assert.Equal(t, asset.Amount, int64(-1))
		assert.Equal(t, asset.Decimals(), uint8(4))
		assert.Equal(t, asset.Symbol.Symbol, "PPP")
		assert.Equal(t, asset.String(), "-0.0001 PPP")
	}
}

func TestProxy(t *testing.T) {
	try.Try(func() {
		ct := NewCurrencyTester()
		ct.validatingTester.ProduceBlocks(2, false)
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
			act.Authorization = []common.PermissionLevel{{Actor: common.N("alice"), Permission: common.DefaultConfig.ActiveName}}
			data := common.Variants{
				"owner": alice,
				"delay": 10,
			}

			//trace := ct.validatingTester.PushAction(&alice, &act.Name, &data)
			trace := ct.validatingTester.PushAction2(&proxy, &act.Name, alice, &data, ct.validatingTester.DefaultExpirationDelta, 0)
			ct.validatingTester.ProduceBlocks(1, false)
			assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
			ct.validatingTester.ProduceBlocks(1, false)
		}

		{
			act1 := types.Action{}
			act1.Account = eosioToken
			act1.Name = common.N("transfer")
			act1.Authorization = []common.PermissionLevel{{Actor: eosioToken, Permission: common.DefaultConfig.ActiveName}}
			data1 := common.Variants{
				"from":     eosioToken,
				"to":       proxy,
				"quantity": "5.0000 EOS",
				"memo":     "fund Proxy",
			}

			//trace1 := ct.PushAction(&eosioToken, &act1.Name, &data1)
			trace1 := ct.validatingTester.PushAction2(&eosioToken, &act1.Name, eosioToken, &data1, ct.validatingTester.DefaultExpirationDelta, 0)
			tt := ct.validatingTester.Control.HeadBlockTime().TimeSinceEpoch().Count()
			expectedDelivery := tt + common.Seconds(10).Count()
			s := "5.0000 EOS"
			expected := common.Asset{}.FromString(&s)
			s1 := "0.0000 EOS"
			expected1 := common.Asset{}.FromString(&s1)
			for ct.validatingTester.Control.HeadBlockTime().TimeSinceEpoch().Count() < expectedDelivery {
				ct.validatingTester.ProduceBlocks(1, false)
				assert.Equal(t, *ct.GetBalance(&proxy), expected)
				assert.Equal(t, *ct.GetBalance(&alice), expected1)
			}

			ct.validatingTester.ProduceBlocks(1, false)
			assert.Equal(t, *ct.GetBalance(&proxy), expected1)
			assert.Equal(t, *ct.GetBalance(&alice), expected)
			assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace1.ID))
			ct.validatingTester.close()
		}

	}).FcLogAndRethrow().End()
}

func TestInputQuantity(t *testing.T) {
	ct := NewCurrencyTester()
	ct.validatingTester.ProduceBlocks(2, false)
	alice := common.N("alice")
	carl := common.N("carl")
	s := "100.0000 EOS"
	expected := common.Asset{}.FromString(&s)
	ct.validatingTester.CreateAccounts([]common.AccountName{common.N("alice"), common.N("bob"), common.N("carl")}, false, true)
	{
		trace := ct.Transfer(&eosioToken, &alice, "100.0000 EOS", "test transfer to alice 100.0000 EOS")
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		assert.Equal(t, *ct.GetBalance(&alice), expected)
		assert.Equal(t, *ct.GetBalance(&alice), expected)
	}

	{
		returning := false
		try.Try(func() {
			ct.Transfer(&alice, &carl, "20.50 USD", "throw")
		}).Catch(func(e exception.EosioAssertMessageException) {
			returning = true
		}).End()
		if returning {
			assert.Equal(t, true, returning)
		}
	}

	{
		s := "125.0256 EOS"
		expected := common.Asset{}.FromString(&s)
		trace := ct.Issue(&alice, "25.0256 EOS", "Issue")
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
		assert.Equal(t, *ct.GetBalance(&alice), expected)
	}
	ct.validatingTester.close()
}

func TestDeferredFailure(t *testing.T) {
	ct := NewCurrencyTester()
	ct.validatingTester.ProduceBlocks(2, false)
	alice := common.N("alice")
	proxy := common.N("proxy")
	bob := common.N("bob")
	ct.validatingTester.CreateAccounts([]common.AccountName{alice, bob, proxy}, false, true)
	ct.validatingTester.DefaultProduceBlock()
	wasmName := "test_contracts/proxy.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	ct.validatingTester.SetCode(proxy, code, nil)
	ct.validatingTester.SetCode(bob, code, nil)
	ct.validatingTester.ProduceBlocks(1, false)
	abiName := "test_contracts/proxy.abi"
	abi, err := ioutil.ReadFile(abiName)
	if err != nil {
		log.Error("pushGenesisBlock is err : %v", err)
	}
	ct.validatingTester.SetAbi(proxy, abi, nil)
	ct.validatingTester.SetAbi(bob, abi, nil)
	// set up proxy owner bob
	{
		act := types.Action{}
		act.Account = proxy
		act.Name = common.N("setowner")
		act.Authorization = []common.PermissionLevel{{Actor: common.N("bob"), Permission: common.DefaultConfig.ActiveName}}
		data := common.Variants{
			"owner": bob,
			"delay": 10,
		}
		//trace := ct.PushAction(&proxy, &act.Name,&data)
		trace := ct.validatingTester.PushAction2(&proxy, &act.Name, bob, &data, ct.validatingTester.DefaultExpirationDelta, 0)
		ct.validatingTester.ProduceBlocks(1, false)
		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
	}

	gto := entity.GeneratedTransactionObject{}

	index, _ := ct.validatingTester.Control.DataBase().GetIndex("byTrxId", &gto)
	count := 0
	itr := index.Begin()
	if !index.CompareEnd(itr) {
		count = 1
	}

	assert.Equal(t, 0, count)

	actionName := common.N("transfer")
	data := common.Variants{
		"from":     eosioToken,
		"to":       proxy,
		"quantity": "5.0000 EOS",
		"memo":     "fund Proxy",
	}
	ct.validatingTester.PushAction2(&eosioToken, &actionName, eosioToken, &data, ct.validatingTester.DefaultExpirationDelta, 0)
	expectedDelivery := ct.validatingTester.Control.PendingBlockTime().TimeSinceEpoch().Count() + common.Seconds(10).Count()
	var deferredId common.TransactionIdType
	if !index.CompareEnd(index.Begin()) {
		count = 1
	}
	assert.Equal(t, 1, count)
	assert.Equal(t, false, ct.validatingTester.ChainHasTransaction(&deferredId))
	s := "5.0000 EOS"
	s2 := "0.0000 EOS"
	for ct.validatingTester.Control.PendingBlockTime().TimeSinceEpoch().Count() < expectedDelivery {
		ct.validatingTester.ProduceBlocks(1, false)
		assert.Equal(t, *ct.GetBalance(&proxy), common.Asset{}.FromString(&s))
		assert.Equal(t, *ct.GetBalance(&bob), common.Asset{}.FromString(&s2))
		assert.Equal(t, 1, count)
		assert.Equal(t, false, ct.validatingTester.ChainHasTransaction(&deferredId))
	}
	itr = index.Begin()
	err = itr.Data(&gto)
	if err != nil {
		fmt.Println("TestDeferredFailure:find gto is error")
	}
	deferredId = gto.TrxId
	//next deferred trx
	expectedRedelivery := ct.validatingTester.Control.PendingBlockTime().TimeSinceEpoch().Count() + common.Seconds(10).Count()
	ct.validatingTester.ProduceBlocks(1, false)

	assert.Equal(t, 1, count)

	assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&deferredId))
	assert.Equal(t, ct.validatingTester.GetTransactionReceipt(&deferredId).Status, types.TransactionStatusSoftFail)
	//first deferred trx end

	// set up alice owner
	count = 0
	if !index.CompareEnd(itr) {
		count = 1
	}
	index.Begin().Data(&gto)
	deferredId2 := gto.TrxId
	// set up alice owner
	{
		act := types.Action{}
		act.Account = bob
		act.Name = common.N("setowner")
		act.Authorization = []common.PermissionLevel{{Actor: common.N("alice"), Permission: common.DefaultConfig.ActiveName}}
		data := common.Variants{
			"owner": alice,
			"delay": 0,
		}
		//trace := ct.PushAction(&proxy, &act.Name,&data)
		trace := ct.validatingTester.PushAction2(&bob, &act.Name, alice, &data, ct.validatingTester.DefaultExpirationDelta, 0)
		ct.validatingTester.ProduceBlocks(1, false)

		assert.Equal(t, true, ct.validatingTester.ChainHasTransaction(&trace.ID))
	}

	for ct.validatingTester.Control.PendingBlockTime().TimeSinceEpoch().Count() < expectedRedelivery {
		ct.validatingTester.ProduceBlocks(1, false)
		assert.Equal(t, *ct.GetBalance(&proxy), common.Asset{}.FromString(&s))
		assert.Equal(t, *ct.GetBalance(&alice), common.Asset{}.FromString(&s2))
		assert.Equal(t, *ct.GetBalance(&bob), common.Asset{}.FromString(&s2))
		assert.Equal(t, 1, count)
		assert.Equal(t, false, ct.validatingTester.ChainHasTransaction(&deferredId2))
	}
	assert.Equal(t, 1, count)
	// Second deferred transaction should be retired in this block and should succeed,
	// which should move tokens from the proxy contract to the bob contract, thereby trigger the bob contract to
	// schedule a third deferred transaction with no delay.
	// That third deferred transaction (which moves tokens from the bob contract to account alice) should be executed immediately
	// after in the same block (note that this is the current deferred transaction scheduling policy in tester and it may change).
	//fmt.Println("1",ct.validatingTester.Control.PendingBlockTime())
	//ct.validatingTester.ProduceBlocks(1, false)

	ct.validatingTester.ProduceBlocks(1, false)

	index3, _ := ct.validatingTester.Control.DataBase().GetIndex("byTrxId", &gto)
	if index3.CompareEnd(index3.Begin()) {
		count = 0
	}
	assert.Equal(t, 0, count)
	assert.Equal(t, ct.validatingTester.GetTransactionReceipt(&deferredId2).Status, types.TransactionStatusExecuted)

	assert.Equal(t, *ct.GetBalance(&proxy), common.Asset{}.FromString(&s2))
	assert.Equal(t, *ct.GetBalance(&alice), common.Asset{}.FromString(&s))
	assert.Equal(t, *ct.GetBalance(&bob), common.Asset{}.FromString(&s2))

	ct.validatingTester.close()
}
