package unittests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/chain/types/generated_containers"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"

	"github.com/stretchr/testify/assert"
)

var charlie = common.N("charlie")
var path string = "/tmp/data/"

type WhitelistBlacklistTester struct {
	tempDir           string
	chain             *BaseTester
	ActorWhitelist    generated.AccountNameSet
	ActorBlacklist    generated.AccountNameSet
	ContractWhitelist generated.AccountNameSet
	ContractBlacklist generated.AccountNameSet
	ActionBlacklist   generated.NamePairSet //pair<account_name, action_name>
	ResourceGreylist  generated.AccountNameSet
	LastProducedBlock map[common.Name]common.BlockIdType
}

type TransferArgs struct {
	From     common.Name
	To       common.Name
	Quantity common.Asset
	Memo     string
}

func NewWhitelistBlacklistTester() *WhitelistBlacklistTester {
	w := &WhitelistBlacklistTester{}

	w.ActorWhitelist = *generated.NewAccountNameSet()
	w.ActorBlacklist = *generated.NewAccountNameSet()
	w.ContractWhitelist = *generated.NewAccountNameSet()
	w.ContractBlacklist = *generated.NewAccountNameSet()
	w.ActionBlacklist = *generated.NewNamePairSet()

	w.ResourceGreylist = *generated.NewAccountNameSet()

	w.LastProducedBlock = make(map[common.AccountName]common.BlockIdType)

	return w
}
func newBaseTesterblack(pushGenesis bool, cfg *chain.Config) *BaseTester {
	t := &BaseTester{}
	t.DefaultExpirationDelta = 6
	t.DefaultBilledCpuTimeUs = 2000
	t.AbiSerializerMaxTime = 1000 * 1000
	t.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	t.LastProducedBlock = make(map[common.AccountName]common.BlockIdType)

	t.Control = chain.NewController(cfg)
	t.Control.Startup()
	if pushGenesis {
		t.pushGenesisBlock()
	}
	return t
}
func (w *WhitelistBlacklistTester) initConfig(bootstrap bool, exist bool) {
	cfg := w.GetDefaultChainConfiguration(path)
	if !exist {
		cfg.BlocksDir = path + cfg.BlocksDir
		cfg.StateDir = path + cfg.StateDir
	} else {
		cfg.BlocksDir = path + "node2/" + cfg.BlocksDir
		cfg.StateDir = path + "node2/" + cfg.StateDir
	}

	cfg.StateSize = 1024 * 1024 * 8
	cfg.StateGuardSize = 0
	cfg.ReversibleCacheSize = 1024 * 1024 * 8
	cfg.ReversibleGuardSize = 0
	cfg.ContractsConsole = true

	cfg.Genesis = types.NewGenesisState()
	cfg.Genesis.InitialTimestamp, _ = common.FromIsoString("2020-01-01T00:00:00.000")
	cfg.Genesis.InitialKey = BaseTester{}.getPublicKey(eosio, "active")

	cfg.ActorWhitelist = w.ActorWhitelist
	cfg.ActorBlacklist = w.ActorBlacklist
	cfg.ContractWhitelist = w.ContractWhitelist
	cfg.ContractBlacklist = w.ContractBlacklist
	cfg.ActionBlacklist = w.ActionBlacklist

	cfg.ResourceGreylist = w.ResourceGreylist

	if !bootstrap {
		w.chain = newBaseTesterblack(false, cfg)
		w.chain.SetLastProducedBlockMap(w.LastProducedBlock)
		return
	}
	w.chain = newBaseTesterblack(false, cfg)
	w.chain.SetLastProducedBlockMap(w.LastProducedBlock)
	accounts := []common.AccountName{eosioToken, alice, bob, charlie}
	w.chain.CreateAccounts(accounts, false, true)
	wasmName := "test_contracts/eosio.token.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	w.chain.SetCode(eosioToken, code, nil)
	abiName := "test_contracts/eosio.token.abi"
	abi, _ := ioutil.ReadFile(abiName)
	w.chain.SetAbi(eosioToken, abi, nil)

	act := common.ActionName(common.N("create"))
	createData := common.Variants{
		"issuer":         eosioToken,
		"maximum_supply": "1000000.00 TOK",
		"can_freeze":     0,
		"can_recall":     0,
		"can_whitelist":  0,
	}

	w.chain.PushAction2(&eosioToken, &act, eosioToken, &createData, w.chain.DefaultExpirationDelta, 0)

	actName := common.ActionName(common.N("issue"))
	issueData := common.Variants{
		"to":       eosioToken,
		"quantity": "1000000.00 TOK",
		"memo":     "issue",
	}
	w.chain.PushAction2(&eosioToken, &actName, eosioToken, &issueData, w.chain.DefaultExpirationDelta, 0)
	w.chain.ProduceBlocks(1, false)
}

func (t *WhitelistBlacklistTester) GetDefaultChainConfiguration(path string) *chain.Config {
	cfg := &chain.Config{}
	//tempDirSuffix := "_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	//cfg.BlocksDir = path + tempDirSuffix
	//cfg.StateDir = path + tempDirSuffix
	cfg.StateSize = 1024 * 1024 * 8
	cfg.StateGuardSize = 0
	cfg.ReversibleCacheSize = 1024 * 1024 * 8
	cfg.ReversibleGuardSize = 0
	cfg.ContractsConsole = true

	cfg.Genesis = types.NewGenesisState()
	cfg.Genesis.InitialTimestamp, _ = common.FromIsoString("2020-01-01T00:00:00.000")
	cfg.Genesis.InitialKey = BaseTester{}.getPublicKey(eosio, "active")

	cfg.ActorWhitelist = *generated.NewAccountNameSet()
	cfg.ActorBlacklist = *generated.NewAccountNameSet()
	cfg.ContractWhitelist = *generated.NewAccountNameSet()
	cfg.ContractBlacklist = *generated.NewAccountNameSet()
	cfg.ActionBlacklist = *generated.NewNamePairSet()
	cfg.KeyBlacklist = *generated.NewPublicKeySet()
	cfg.ResourceGreylist = *generated.NewAccountNameSet()
	cfg.TrustedProducers = *generated.NewAccountNameSet()

	return cfg
}

func (t *WhitelistBlacklistTester) Shutdown() {
	//assert.Equal(!common.Empty(t.chain),&Exception{},"chain is not up")
	lastProducedBlock := t.chain.GetLastProducedBlockMap()
	log.Info("%v", lastProducedBlock)
	t.chain.close()
}

func (w *WhitelistBlacklistTester) Transfer(from common.AccountName, to common.AccountName, quantity string) *types.TransactionTrace {

	act := types.Action{}
	act.Account = eosioToken
	act.Name = common.N("transfer")
	act.Authorization = []common.PermissionLevel{{Actor: from, Permission: common.DefaultConfig.ActiveName}}
	data := common.Variants{
		"from":     from,
		"to":       to,
		"quantity": quantity,
		"memo":     "transfer",
	}
	trace := w.chain.PushAction2(&eosioToken, &act.Name, from, &data, w.chain.DefaultExpirationDelta, 0)
	w.chain.ProduceBlocks(1, false)
	return trace
}

func TestActorWhitelist(t *testing.T) {
	os.RemoveAll(path)
	w := NewWhitelistBlacklistTester()
	w.initConfig(true, false)

	w.ActorWhitelist.Add(common.DefaultConfig.SystemAccountName)
	w.ActorWhitelist.Add(eosioToken)
	w.ActorWhitelist.Add(alice)
	w.chain.Control.SetActorWhiteList(&w.ActorWhitelist)
	w.Transfer(eosioToken, alice, "1000.00 TOK")
	w.Transfer(alice, bob, "100.00 TOK")

	var ex string
	{
		try.Try(func() {
			w.Transfer(bob, alice, "1.00 TOK")
		}).Catch(func(e exception.ActorWhitelistException) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "authorizing actor(s) in transaction are not on the actor whitelist"))
	}

	{
		trx := types.SignedTransaction{}
		action := types.Action{}
		action.Account = eosioToken
		action.Name = common.N("transfer")
		action.Authorization = []common.PermissionLevel{{Actor: alice, Permission: common.DefaultConfig.ActiveName}, {Actor: bob, Permission: common.DefaultConfig.ActiveName}}
		data := common.Variants{
			"from":     alice,
			"to":       bob,
			"quantity": "10.00 TOK",
			"memo":     "transfer",
		}
		acnt := w.chain.Control.GetAccount(eosioToken)
		a := acnt.GetAbi()
		buf, _ := json.Marshal(data)
		action.Data, _ = a.EncodeAction(action.Account, buf)
		trx.Actions = append(trx.Actions, &action)
		w.chain.SetTransactionHeaders(&trx.Transaction, w.chain.DefaultBilledCpuTimeUs, 0)
		prikey := w.chain.getPrivateKey(alice, "active")
		chainID := w.chain.Control.GetChainId()
		trx.Sign(&prikey, &chainID)
		prikey2 := w.chain.getPrivateKey(bob, "active")
		chainID2 := w.chain.Control.GetChainId()
		trx.Sign(&prikey2, &chainID2)
		try.Try(func() {
			w.chain.PushTransaction(&trx, common.MaxTimePoint(), w.chain.DefaultBilledCpuTimeUs)
		}).Catch(func(e exception.ActorWhitelistException) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "authorizing actor(s) in transaction are not on the actor whitelist"))

		w.chain.ProduceBlocks(1, false)
	}

}

func TestActorBlacklist(t *testing.T) {
	os.RemoveAll(path)
	w := NewWhitelistBlacklistTester()
	w.initConfig(true, false)
	w.ActorBlacklist.Add(bob)
	w.chain.Control.SetActorBlackList(&w.ActorBlacklist)

	w.Transfer(eosioToken, alice, "1000.00 TOK")
	w.Transfer(alice, bob, "100.00 TOK")

	var ex string
	{
		try.Try(func() {
			w.Transfer(bob, alice, "1.00 TOK")
		}).Catch(func(e exception.ActorBlacklistException) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "authorizing actor(s) in transaction are on the actor blacklist: bob111111111"))
	}

	{
		trx := types.SignedTransaction{}
		action := types.Action{}
		action.Account = eosioToken
		action.Name = common.N("transfer")
		action.Authorization = []common.PermissionLevel{{Actor: alice, Permission: common.DefaultConfig.ActiveName}, {Actor: bob, Permission: common.DefaultConfig.ActiveName}}
		data := common.Variants{
			"from":     alice,
			"to":       bob,
			"quantity": "10.00 TOK",
			"memo":     "transfer",
		}
		acnt := w.chain.Control.GetAccount(eosioToken)
		a := acnt.GetAbi()
		buf, _ := json.Marshal(data)
		action.Data, _ = a.EncodeAction(action.Account, buf)
		trx.Actions = append(trx.Actions, &action)
		w.chain.SetTransactionHeaders(&trx.Transaction, w.chain.DefaultBilledCpuTimeUs, 0)
		prikey := w.chain.getPrivateKey(alice, "active")
		chainID := w.chain.Control.GetChainId()
		trx.Sign(&prikey, &chainID)
		prikey2 := w.chain.getPrivateKey(bob, "active")
		chainID2 := w.chain.Control.GetChainId()
		trx.Sign(&prikey2, &chainID2)
		try.Try(func() {
			w.chain.PushTransaction(&trx, common.MaxTimePoint(), w.chain.DefaultBilledCpuTimeUs)
		}).Catch(func(e exception.ActorBlacklistException) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "authorizing actor(s) in transaction are on the actor blacklist: bob111111111"))

		w.chain.ProduceBlocks(1, false)
	}
}

func TestContractWhitelist(t *testing.T) {
	os.RemoveAll(path)
	w := NewWhitelistBlacklistTester()
	w.initConfig(true, false)
	w.ContractWhitelist.Add(common.DefaultConfig.SystemAccountName)
	w.ContractWhitelist.Add(eosioToken)
	w.ContractWhitelist.Add(bob)
	w.chain.Control.SetContractWhiteList(&w.ContractWhitelist)

	w.Transfer(eosioToken, alice, "1000.00 TOK")
	w.Transfer(alice, eosioToken, "1.00 TOK")
	w.Transfer(alice, bob, "1.00 TOK")
	w.Transfer(alice, charlie, "100.00 TOK")
	w.Transfer(charlie, alice, "1.00 TOK")
	w.chain.ProduceBlocks(1, false)

	{
		wasmName := "test_contracts/eosio.token.wasm"
		code, _ := ioutil.ReadFile(wasmName)
		w.chain.SetCode(bob, code, nil)
		abiName := "test_contracts/eosio.token.abi"
		abi, err := ioutil.ReadFile(abiName)
		if err != nil {
			log.Error("eosio.token.abi is err : %v", err)
		}
		w.chain.SetAbi(bob, abi, nil)
		w.chain.ProduceBlocks(1, false)

		w.chain.SetCode(charlie, code, nil)
		w.chain.SetAbi(charlie, code, nil)
		w.chain.ProduceBlocks(1, false)
	}
	var ex string
	{
		try.Try(func() {
			w.Transfer(alice, charlie, "1.00 TOK")
		}).Catch(func(e exception.ContractWhitelistException) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "account charlie is not on the contract whitelist"))

		w.chain.ProduceBlocks(1, false)
	}

	{
		act := common.ActionName(common.N("create"))
		createData := common.Variants{
			"issuer":         bob,
			"maximum_supply": "1000000.00 CUR",
			"can_freeze":     0,
			"can_recall":     0,
			"can_whitelist":  0,
		}

		w.chain.PushAction2(&bob, &act, bob, &createData, w.chain.DefaultExpirationDelta, 0)
	}
	{
		try.Try(func() {
			charlie := charlie
			act := common.ActionName(common.N("create"))
			data := common.Variants{
				"issuer":         "chalie",
				"maximum_supply": "1000000.00 CUR",
				"can_freeze":     0,
				"can_recall":     0,
				"can_whitelist":  0,
			}
			w.chain.PushAction2(&charlie, &act, charlie, &data, w.chain.DefaultExpirationDelta, 0)
		}).Catch(func(e exception.ContractWhitelistException) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "account charlie is not on the contract whitelist"))
	}
	w.chain.ProduceBlocks(1, false)
}

func TestContractBlacklist(t *testing.T) {
	os.RemoveAll(path)
	w := NewWhitelistBlacklistTester()
	w.initConfig(true, false)
	w.ContractBlacklist.Add(charlie)
	w.chain.Control.SetContractBlackList(&w.ContractBlacklist)

	w.Transfer(eosioToken, alice, "1000.00 TOK")
	w.Transfer(alice, eosioToken, "1.00 TOK")
	w.Transfer(alice, bob, "1.00 TOK")
	w.Transfer(alice, charlie, "100.00 TOK")
	w.Transfer(charlie, alice, "1.00 TOK")
	w.chain.ProduceBlocks(1, false)

	{
		wasmName := "test_contracts/eosio.token.wasm"
		code, _ := ioutil.ReadFile(wasmName)
		w.chain.SetCode(bob, code, nil)
		abiName := "test_contracts/eosio.token.abi"
		abi, err := ioutil.ReadFile(abiName)
		if err != nil {
			log.Error("eosio.token.abi is err : %v", err)
		}
		w.chain.SetAbi(bob, abi, nil)
		w.chain.ProduceBlocks(1, false)

		w.chain.SetCode(charlie, code, nil)
		w.chain.SetAbi(charlie, code, nil)
		w.chain.ProduceBlocks(1, false)
	}
	var ex string
	w.Transfer(alice, bob, "1.00 TOK")

	{
		try.Try(func() {
			w.Transfer(alice, charlie, "1.00 TOK")
		}).Catch(func(e exception.ContractBlacklistException) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "account charlie is on the contract blacklist"))
	}

	{
		act := common.ActionName(common.N("create"))
		createData := common.Variants{
			"issuer":         bob,
			"maximum_supply": "1000000.00 CUR",
		}

		w.chain.PushAction2(&bob, &act, bob, &createData, w.chain.DefaultExpirationDelta, 0)
	}

	{
		try.Try(func() {
			act := common.ActionName(common.N("create"))
			createData := common.Variants{
				"issuer":         charlie,
				"maximum_supply": "1000000.00 CUR",
			}
			w.chain.PushAction2(&charlie, &act, charlie, &createData, w.chain.DefaultExpirationDelta, 0)
		}).Catch(func(e exception.ContractBlacklistException) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "account charlie is on the contract blacklist"))
	}
	w.chain.ProduceBlocks(1, false)
}

func TestActionBlacklist(t *testing.T) {
	os.RemoveAll(path)
	w := NewWhitelistBlacklistTester()
	w.initConfig(true, false)
	w.ContractWhitelist.Add(common.DefaultConfig.SystemAccountName)
	w.ContractWhitelist.Add(eosioToken)
	w.ContractWhitelist.Add(bob)
	w.ContractWhitelist.Add(charlie)

	abl := common.NamePair{First: charlie, Second: common.N("create")}
	w.ActionBlacklist.Add(abl)

	w.chain.Control.SetContractWhiteList(&w.ContractWhitelist)
	w.chain.Control.SetActionBlackList(&w.ActionBlacklist)

	w.Transfer(eosioToken, alice, "1000.00 TOK")
	w.chain.ProduceBlocks(1, false)

	{
		wasmName := "test_contracts/eosio.token.wasm"
		code, _ := ioutil.ReadFile(wasmName)
		w.chain.SetCode(bob, code, nil)
		abiName := "test_contracts/eosio.token.abi"
		abi, err := ioutil.ReadFile(abiName)
		if err != nil {
			log.Error("eosio.token.abi is err : %v", err)
		}
		w.chain.SetAbi(bob, abi, nil)
		w.chain.ProduceBlocks(1, false)

		w.chain.SetCode(charlie, code, nil)
		w.chain.SetAbi(charlie, code, nil)
		w.chain.ProduceBlocks(1, false)

		w.Transfer(alice, bob, "1.00 TOK")
		w.Transfer(alice, charlie, "1.00 TOK")

		act := common.ActionName(common.N("create"))
		createData := common.Variants{
			"issuer":         bob,
			"maximum_supply": "1000000.00 CUR",
		}
		w.chain.PushAction2(&bob, &act, bob, &createData, w.chain.DefaultExpirationDelta, 0)
	}

	{
		var ex string
		try.Try(func() {
			act := common.ActionName(common.N("create"))
			createData := common.Variants{
				"issuer":         charlie,
				"maximum_supply": "1000000.00 CUR",
			}
			w.chain.PushAction2(&charlie, &act, charlie, &createData, w.chain.DefaultExpirationDelta, 0)
		}).Catch(func(e exception.ActionBlacklistException) {
			ex = e.DetailMessage()
		}).End()
		assert.True(t, inString(ex, "action 'charlie::create' is on the action blacklist"))
	}
	w.chain.ProduceBlocks(1, false)
}

func TestBlacklistEosio(t *testing.T) {
	os.RemoveAll(path)
	w := NewWhitelistBlacklistTester()
	w.initConfig(true, false)
	w.chain.ProduceBlocks(1, false)

	wasmName := "test_contracts/eosio.token.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	w.chain.SetCode(common.DefaultConfig.SystemAccountName, code, nil)
	abiName := "test_contracts/eosio.token.abi"
	abi, err := ioutil.ReadFile(abiName)
	if err != nil {
		log.Error("eosio.token.abi is err : %v", err)
	}
	w.chain.SetAbi(common.DefaultConfig.SystemAccountName, abi, nil)
	w.chain.ProduceBlocks(1, false)
	w.Shutdown()

	w.ContractBlacklist.Add(common.DefaultConfig.SystemAccountName)
	w.initConfig(false, false)
	w2 := NewWhitelistBlacklistTester()
	w2.initConfig(true, true)
	w2.chain.ProduceBlocks(1, false)
	for w2.chain.Control.HeadBlockNum() < w.chain.Control.HeadBlockNum() {
		b := w.chain.Control.FetchBlockByNumber(w2.chain.Control.HeadBlockNum() + 1)
		w2.chain.PushBlock(b)
	}

	w.chain.ProduceBlocks(2, false)

	for w2.chain.Control.HeadBlockNum() < w.chain.Control.HeadBlockNum() {
		b := w.chain.Control.FetchBlockByNumber(w2.chain.Control.HeadBlockNum() + 1)
		w2.chain.PushBlock(b)
	}
}

func TestDeferredBlacklistFailure(t *testing.T) {
	os.RemoveAll(path)
	w := NewWhitelistBlacklistTester()
	w.initConfig(true, false)
	w.chain.ProduceBlocks(1, false)
	wasmName := "test_contracts/deferred_test.wasm"
	code, _ := ioutil.ReadFile(wasmName)

	abiName := "test_contracts/deferred_test.abi"
	abi, err := ioutil.ReadFile(abiName)
	if err != nil {
		log.Error("deferred_test.abi is err : %v", err)
	}
	w.chain.SetCode(bob, code, nil)
	w.chain.SetAbi(bob, abi, nil)

	w.chain.SetCode(charlie, code, nil)
	w.chain.SetAbi(charlie, abi, nil)
	w.chain.ProduceBlocks(1, false)

	act := common.ActionName(common.N("defercall"))
	data := common.Variants{
		"payer":     alice,
		"sender_id": 0,
		"contract":  charlie,
		"payload":   10,
	}
	w.chain.PushAction2(&bob, &act, alice, &data, w.chain.DefaultExpirationDelta, 0)
	w.chain.ProduceBlocks(2, false)
	w.Shutdown()

	w.ContractBlacklist.Add(charlie)
	w.initConfig(false, false)

	w2 := NewWhitelistBlacklistTester()
	w2.initConfig(false, true)
	for w2.chain.Control.HeadBlockNum() < w.chain.Control.HeadBlockNum() {
		b := w.chain.Control.FetchBlockByNumber(w2.chain.Control.HeadBlockNum() + 1)
		w2.chain.PushBlock(b)
	}
	data2 := common.Variants{
		"payer":     alice,
		"sender_id": 1,
		"contract":  charlie,
		"payload":   10,
	}
	w.chain.PushAction2(&bob, &act, alice, &data2, w.chain.DefaultExpirationDelta, 0)

	var ex string
	try.Try(func() {
		w.chain.ProduceBlocks(1, false)
	}).Catch(func(e exception.Exception) {
		ex = e.DetailMessage()
	})
	assert.True(t, inString(ex, "account charlie is on the contract blacklist"))

	w.chain.ProduceBlocks(2, true)

	for w2.chain.Control.HeadBlockNum() < w.chain.Control.HeadBlockNum() {
		b := w.chain.Control.FetchBlockByNumber(w2.chain.Control.HeadBlockNum() + 1)
		w2.chain.PushBlock(b)
	}
}

func TestBlacklistOnerror(t *testing.T) {
	os.RemoveAll(path)
	w := NewWhitelistBlacklistTester()
	w.initConfig(true, false)
	w.chain.ProduceBlocks(1, false)
	wasmName := "test_contracts/deferred_test.wasm"
	code, _ := ioutil.ReadFile(wasmName)

	abiName := "test_contracts/deferred_test.abi"
	abi, err := ioutil.ReadFile(abiName)
	if err != nil {
		log.Error("deferred_test.abi is err : %v", err)
	}
	w.chain.SetCode(bob, code, nil)
	w.chain.SetAbi(bob, abi, nil)

	w.chain.SetCode(charlie, code, nil)
	w.chain.SetAbi(charlie, abi, nil)
	w.chain.ProduceBlocks(1, false)

	act := common.ActionName(common.N("defercall"))
	data := common.Variants{
		"payer":     alice,
		"sender_id": 0,
		"contract":  charlie,
		"payload":   13,
	}
	w.chain.PushAction2(&bob, &act, alice, &data, w.chain.DefaultExpirationDelta, 0)

	w.chain.ProduceBlocks(1, false)

	w.Shutdown()
	abl := common.NamePair{First: common.DefaultConfig.SystemAccountName, Second: common.N("onerror")}
	w.ActionBlacklist.Add(abl)
	w.initConfig(false, false)

	w.chain.PushAction2(&bob, &act, alice, &data, w.chain.DefaultExpirationDelta, 0)

	var ex string
	try.Try(func() {
		w.chain.ProduceBlocks(1, false)
	}).Catch(func(e exception.Exception) {
		ex = e.DetailMessage()
	})
	assert.True(t, inString(ex, "action 'eosio::onerror' is on the action blacklist"))
}
