package unittests

import (
	. "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/log"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

type EosioMsigTester struct {
	BaseTester
	abiSer abi_serializer.AbiSerializer
}

func newEosioMsigTester(pushGenesis bool, readMode DBReadMode) *EosioMsigTester {
	e := &EosioMsigTester{}
	e.BaseTester = *newBaseTester(pushGenesis, readMode)
	return e
}

func initEosioMsigTester() *EosioMsigTester {
	e := newEosioMsigTester(true, SPECULATIVE)
	e.CreateAccounts([]common.AccountName{eosioMsig, eosioStake, eosioRam, eosioRamFee, alice, bob, carol}, false, true)
	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	data := common.Variants{
		"account": eosioMsig,
		"is_priv": 1,
	}
	acttype := common.N("setpriv")
	e.PushAction2(&eosio, &acttype, eosio, &data, e.DefaultExpirationDelta, 0)
	wasmName := "test_contracts/eosio.msig.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	e.SetCode(eosioMsig, code, nil)
	abiName := "test_contracts/eosio.msig.abi"
	abi, _ := ioutil.ReadFile(abiName)
	e.SetAbi(eosioMsig, abi, nil)
	e.ProduceBlocks(1, false)
	accnt := entity.AccountObject{Name: eosioMsig}
	e.Control.DB.Find("byName", accnt, &accnt)
	abiDef := abi_serializer.AbiDef{}
	if !abi_serializer.ToABI(accnt.Abi, &abiDef) {
		log.Error("eosio_system_tester::InitializeMultisig failed with ToAbi")
	}
	e.abiSer.SetAbi(&abiDef, e.AbiSerializerMaxTime)
	return e
}

func (e EosioMsigTester) CreateAccountWithResources(name common.AccountName, creator common.AccountName,
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

	newAccount := NewAccount{
		Creator: creator,
		Name:    name,
		Owner:   ownerAuth,
		Active:  activeAuth,
	}
	data, _ := rlp.EncodeToBytes(newAccount)
	act := &types.Action{
		Account:       newAccount.GetAccount(),
		Name:          newAccount.GetName(),
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

func (e EosioMsigTester) CreateCurrency(contract common.Name, manager common.Name, maxSupply common.Asset) {
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

func (e EosioMsigTester) Issue(to common.Name, amount common.Asset, manager common.Name) {
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

func (e EosioMsigTester) Transfer(from common.Name, to common.Name, amount common.Asset, manager common.Name) {
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

func (e EosioMsigTester) GetBalance(act common.AccountName) common.Asset {
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

func (e EosioMsigTester) PushAction(signer common.AccountName, name common.ActionName, data common.Variants, auth bool) *types.TransactionTrace {
	var accounts []*common.AccountName
	if auth {
		accounts = append(accounts, &signer)
	}
	trace := e.PushAction3(&eosioMsig, &name, accounts, &data, e.DefaultExpirationDelta, 0)
	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	return trace
}

func (e EosioMsigTester) ReqAuth(from common.AccountName, auths []types.PermissionLevel, maxSerialization common.Microseconds) types.Transaction {
	var auth []types.PermissionLevel
	for _, level := range auths {
		auth = append(auth, level)
	}
	reqauth := common.Variants{"from": from}
	data, _ := rlp.EncodeToBytes(reqauth)
	act := types.Action{
		Account:       eosio,
		Name:          common.N("requauth"),
		Authorization: auths,
		Data:          data,
	}
	var trx types.Transaction
	trx.Actions = append(trx.Actions, &act)
	e.SetTransactionHeaders(&trx, 0, 0)
	expiration, _ := common.FromIsoString("2020-01-01T00:00:30.000")
	trx.Expiration = common.NewTimePointSecTp(expiration)
	trx.RefBlockNum = 2
	trx.RefBlockPrefix = 3
	return trx
}

func (e EosioMsigTester) Propose(proposer common.AccountName, proposalName common.Name, trx types.Transaction, request []types.PermissionLevel) *types.TransactionTrace {
	propose := common.Variants{
		"proposer":      proposer,
		"proposal_name": proposalName,
		"trx":           trx,
		"requested":     request,
	}
	acttype := common.N("propose")
	return e.PushAction(proposer, acttype, propose, true)
}

func (e EosioMsigTester) Approve(approver common.AccountName, proposer common.AccountName, proposalName common.Name, level types.PermissionLevel) *types.TransactionTrace {
	approve := common.Variants{
		"proposer":      proposer,
		"proposal_name": proposalName,
		"level":         level,
	}
	acttype := common.N("approve")
	return e.PushAction(approver, acttype, approve, true)
}

func (e EosioMsigTester) Exec(executer common.AccountName, proposer common.AccountName, proposalName common.Name) *types.TransactionTrace {
	exec := common.Variants{
		"proposer":      proposer,
		"proposal_name": proposalName,
		"executer":      executer,
	}
	acttype := common.N("exec")
	return e.PushAction(executer, acttype, exec, true)
}

func TestProposeApproveExecute(t *testing.T) {
	e := initEosioMsigTester()
	trx := e.ReqAuth(alice, []types.PermissionLevel{{Actor: alice, Permission: common.DefaultConfig.ActiveName}}, e.AbiSerializerMaxTime)
	trace := e.Propose(alice, common.N("first"), trx, []types.PermissionLevel{{alice, common.DefaultConfig.ActiveName}})
	assert.True(t, e.ChainHasTransaction(&trace.ID))

	//fail to execute before approval
	exec := func() {e.Exec(alice, alice, common.N("first"))}
	CatchThrowExceptionAndMsg(t, &exception.EosioAssertMessageException{}, "transaction authorization failed", exec)

	//approve and execute
	trace = e.Approve(alice, alice, common.N("first"), types.PermissionLevel{Actor:alice, Permission:common.DefaultConfig.ActiveName})
	assert.True(t, e.ChainHasTransaction(&trace.ID))
	trace = e.Exec(alice, alice, common.N("first"))
	assert.True(t, e.ChainHasTransaction(&trace.ID))

}
