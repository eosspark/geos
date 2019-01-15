package unittests

import (
	. "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/chain_interface"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

type EosioSudoTester struct {
	BaseTester
	abiSer abi_serializer.AbiSerializer
}

func newEosioSudoTester(pushGenesis bool, readMode DBReadMode) *EosioSudoTester {
	e := &EosioSudoTester{}
	e.BaseTester = *newBaseTester(pushGenesis, readMode)
	return e
}

func initEosioSudoTester() *EosioSudoTester {
	e := newEosioSudoTester(true, SPECULATIVE)
	e.CreateAccounts([]common.AccountName{eosioMsig, producer1, producer2, producer3, producer4, producer5, alice, bob, carol}, false, true)
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

	trx := types.SignedTransaction{}
	auth := types.Authority{Threshold: 1, Accounts: []types.PermissionLevelWeight{{Permission: types.PermissionLevel{Actor: eosio, Permission: common.DefaultConfig.ActiveName}, Weight: 1}}}
	newAccount := NewAccount{
		Creator: eosio,
		Name:    eosioSudo,
		Owner:   auth,
		Active:  auth,
	}
	actionData, _ := rlp.EncodeToBytes(newAccount)
	action := types.Action{
		Account:       newAccount.GetAccount(),
		Name:          newAccount.GetName(),
		Authorization: []types.PermissionLevel{{eosio, common.DefaultConfig.ActiveName}},
		Data:          actionData,
	}
	trx.Actions = append(trx.Actions, &action)
	e.SetTransactionHeaders(&trx.Transaction, e.DefaultExpirationDelta, 0)
	chainId := e.Control.GetChainId()
	pk := e.getPrivateKey(eosio, "active")
	trx.Sign(&pk, &chainId)
	e.PushTransaction(&trx, common.MaxTimePoint(), e.DefaultBilledCpuTimeUs)

	data = common.Variants{
		"account": eosioSudo,
		"is_priv": 1,
	}
	acttype = common.N("setpriv")
	e.PushAction2(&eosio, &acttype, eosio, &data, e.DefaultExpirationDelta, 0)

	systemPrivateKey := e.getPrivateKey(eosio, "active")
	wasmName = "test_contracts/eosio.sudo.wasm"
	code, _ = ioutil.ReadFile(wasmName)
	e.SetCode(eosioSudo, code, &systemPrivateKey)
	abiName = "test_contracts/eosio.sudo.abi"
	abi, _ = ioutil.ReadFile(abiName)
	e.SetAbi(eosioSudo, abi, &systemPrivateKey)

	e.ProduceBlocks(1, false)

	authority := types.Authority{
		Threshold: 1,
		Keys:      []types.KeyWeight{{Key: e.getPublicKey(eosio, "active"), Weight: 1}},
		Accounts: []types.PermissionLevelWeight{
			{Permission: types.PermissionLevel{Actor: common.DefaultConfig.ProducersAccountName, Permission: common.DefaultConfig.ActiveName}, Weight: 1},
		},
	}
	e.SetAuthority(
		eosio,
		common.DefaultConfig.ActiveName,
		authority,
		common.DefaultConfig.OwnerName,
		&[]types.PermissionLevel{{eosio, common.DefaultConfig.OwnerName}},
		&[]ecc.PrivateKey{e.getPrivateKey(eosio, "active")},
	)

	e.SetProducers(&[]common.AccountName{producer1, producer2, producer3, producer4, producer5})

	e.ProduceBlocks(1, false)

	accnt := entity.AccountObject{Name: eosioSudo}
	e.Control.DB.Find("byName", accnt, &accnt)
	abiDef := abi_serializer.AbiDef{}
	if !abi_serializer.ToABI(accnt.Abi, &abiDef) {
		log.Error("eosio_system_tester::InitializeMultisig failed with ToAbi")
	}
	e.abiSer.SetAbi(&abiDef, e.AbiSerializerMaxTime)

	for e.Control.PendingBlockState().Header.Producer == eosio {
		e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	}
	return e
}

func (e EosioSudoTester) Propose(proposer common.AccountName, proposalName common.Name, trx types.Transaction, request []types.PermissionLevel) *types.TransactionTrace {
	propose := common.Variants{
		"proposer":      proposer,
		"proposal_name": proposalName,
		"trx":           trx,
		"requested":     request,
	}
	acttype := common.N("propose")
	return e.PushAction2(&eosioMsig, &acttype, proposer, &propose, e.DefaultExpirationDelta, 0)
}

func (e EosioSudoTester) Approve(approver common.AccountName, proposer common.AccountName, proposalName common.Name, level types.PermissionLevel) *types.TransactionTrace {
	approve := common.Variants{
		"proposer":      proposer,
		"proposal_name": proposalName,
		"level":         level,
	}
	acttype := common.N("approve")
	return e.PushAction2(&eosioMsig, &acttype, approver, &approve, e.DefaultExpirationDelta, 0)
}

func (e EosioSudoTester) UnApprove(unapprover common.AccountName, proposer common.AccountName, proposalName common.Name, level types.PermissionLevel) *types.TransactionTrace {
	unapprove := common.Variants{
		"proposer":      proposer,
		"proposal_name": proposalName,
		"level":         level,
	}
	acttype := common.N("unapprove")
	return e.PushAction2(&eosioMsig, &acttype, unapprover, &unapprove, e.DefaultExpirationDelta, 0)
}

func (e EosioSudoTester) Exec(executer common.AccountName, proposer common.AccountName, proposalName common.Name) *types.TransactionTrace {
	exec := common.Variants{
		"proposer":      proposer,
		"proposal_name": proposalName,
		"executer":      executer,
	}
	acttype := common.N("exec")
	return e.PushAction2(&eosioMsig, &acttype, executer, &exec, e.DefaultExpirationDelta, 0)
}

func (e EosioSudoTester) SudoExec(executer common.AccountName, trx *types.Transaction, expiration uint32) types.Transaction {
	auth := []types.PermissionLevel{{executer, common.DefaultConfig.ActiveName}, {eosioSudo, common.DefaultConfig.ActiveName}}
	//exec := common.Variants{
	//	"executer": executer,
	//	"trx":      *trx,
	//}
	type Exec struct {
		Executer common.AccountName
		Trx types.Transaction
	}
	exec := Exec{executer, *trx}
	data, _ := rlp.EncodeToBytes(exec)
	actObj := types.Action{
		Account:       eosioSudo,
		Name:          common.N("exec"),
		Authorization: auth,
		Data:          data,
	}
	trx2 := types.Transaction{}
	trx2.Actions = append(trx2.Actions, &actObj)
	e.SetTransactionHeaders(&trx2, expiration, 0)
	return trx2
}

func (e EosioSudoTester) ReqAuth(from common.AccountName, auths []types.PermissionLevel, expiration uint32) types.Transaction {
	var auth []types.PermissionLevel
	for _, level := range auths {
		auth = append(auth, level)
	}
	//reqauth := common.Variants{"from": from}
	type Reqauth struct {
		From common.AccountName
	}
	reqauth := Reqauth{From: from}
	data, _ := rlp.EncodeToBytes(reqauth)
	act := types.Action{
		Account:       eosio,
		Name:          common.N("reqauth"),
		Authorization: auths,
		Data:          data,
	}
	var trx types.Transaction
	trx.Actions = append(trx.Actions, &act)
	e.SetTransactionHeaders(&trx, expiration, 0)
	return trx
}

func TestSudoExecDirect(t *testing.T) {
	e := initEosioSudoTester()
	trx := e.ReqAuth(bob, []types.PermissionLevel{{bob,common.DefaultConfig.ActiveName}}, e.DefaultExpirationDelta)
	sudoTrx := types.SignedTransaction{Transaction: e.SudoExec(alice, &trx, e.DefaultExpirationDelta)}
	chainId := e.Control.GetChainId()
	pk := e.getPrivateKey(alice, "active")
	sudoTrx.Sign(&pk, &chainId)
	prodVector := []common.AccountName{producer1, producer2, producer3, producer4}
	trace := types.TransactionTrace{}
	appliedTrxCaller := chain_interface.AppliedTransactionCaller{Caller: func(t *types.TransactionTrace) {if t.Scheduled {trace = *t}}}
	e.Control.AppliedTransaction.Connect(&appliedTrxCaller)
	for _, actor := range prodVector {
		pk = e.getPrivateKey(actor, "active")
		sudoTrx.Sign(&pk, &chainId)
	}
	e.PushTransaction(&sudoTrx, common.MaxTimePoint(), e.DefaultBilledCpuTimeUs)
	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	assert.Equal(t, int(1), len(trace.ActionTraces))
	assert.Equal(t, eosio, trace.ActionTraces[0].Act.Account)
	assert.Equal(t, common.N("reqauth"), trace.ActionTraces[0].Act.Name)
	assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
}

func TestSudoWithMsig(t *testing.T) {
	e := initEosioSudoTester()
	trx := e.ReqAuth(bob, []types.PermissionLevel{{bob,common.DefaultConfig.ActiveName}}, e.DefaultExpirationDelta)
	sudoTrx := e.SudoExec(alice, &trx, e.DefaultExpirationDelta)
	active := common.DefaultConfig.ActiveName
	e.Propose(carol, common.N("first"), sudoTrx, []types.PermissionLevel{
		{alice, active}, {producer1, active},
		{producer2, active}, {producer3, active},
		{producer4, active}, {producer5, active},
	})
	e.Approve(alice, carol, common.N("first"), types.PermissionLevel{Actor:alice, Permission:active})

	// More than 2/3 of block producers approve
	e.Approve(producer1, carol, common.N("first"), types.PermissionLevel{Actor:producer1, Permission:active})
	e.Approve(producer2, carol, common.N("first"), types.PermissionLevel{Actor:producer2, Permission:active})
	e.Approve(producer3, carol, common.N("first"), types.PermissionLevel{Actor:producer3, Permission:active})
	e.Approve(producer4, carol, common.N("first"), types.PermissionLevel{Actor:producer4, Permission:active})

	var traces []types.TransactionTrace
	appliedTrxCaller := chain_interface.AppliedTransactionCaller{Caller: func(t *types.TransactionTrace) {if t.Scheduled {traces = append(traces, *t)}}}
	e.Control.AppliedTransaction.Connect(&appliedTrxCaller)

	// Now the proposal should be ready to execute
	e.Exec(alice, carol, common.N("first"))
	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	assert.Equal(t, int(2), len(traces))

	assert.Equal(t, int(1), len(traces[0].ActionTraces))
	assert.Equal(t, eosioSudo, traces[0].ActionTraces[0].Act.Account)
	assert.Equal(t, common.N("exec"), traces[0].ActionTraces[0].Act.Name)
	assert.Equal(t, types.TransactionStatusExecuted, traces[0].Receipt.Status)

	assert.Equal(t, int(1), len(traces[1].ActionTraces))
	assert.Equal(t, eosio, traces[1].ActionTraces[0].Act.Account)
	assert.Equal(t, common.N("reqauth"), traces[1].ActionTraces[0].Act.Name)
	assert.Equal(t, types.TransactionStatusExecuted, traces[1].Receipt.Status)
}

func TestSudoWithMsigUnapprove(t *testing.T) {
	e := initEosioSudoTester()
	trx := e.ReqAuth(bob, []types.PermissionLevel{{bob,common.DefaultConfig.ActiveName}}, e.DefaultExpirationDelta)
	sudoTrx := e.SudoExec(alice, &trx, e.DefaultExpirationDelta)
	active := common.DefaultConfig.ActiveName
	e.Propose(carol, common.N("first"), sudoTrx, []types.PermissionLevel{
		{alice, active}, {producer1, active},
		{producer2, active}, {producer3, active},
		{producer4, active}, {producer5, active},
	})
	e.Approve(alice, carol, common.N("first"), types.PermissionLevel{Actor:alice, Permission:active})

	// 3 of the 4 needed producers approve
	e.Approve(producer1, carol, common.N("first"), types.PermissionLevel{Actor:producer1, Permission:active})
	e.Approve(producer2, carol, common.N("first"), types.PermissionLevel{Actor:producer2, Permission:active})
	e.Approve(producer3, carol, common.N("first"), types.PermissionLevel{Actor:producer3, Permission:active})
	e.Approve(producer4, carol, common.N("first"), types.PermissionLevel{Actor:producer4, Permission:active})

	// first producer takes back approval
	e.UnApprove(producer1, carol, common.N("first"), types.PermissionLevel{Actor:producer1, Permission:active})

	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	// The proposal should not have sufficient approvals to pass the authorization checks of eosio.sudo::exec.
	exec := func() {e.Exec(alice, carol, common.N("first"))}
	CatchThrowExceptionAndMsg(t, &exception.EosioAssertMessageException{}, "transaction authorization failed", exec)
}

func TestSudoWithMsigProducersChange(t *testing.T) {
	e := initEosioSudoTester()
	e.CreateAccounts([]common.AccountName{common.N("newprod1")}, false, true)
	trx := e.ReqAuth(bob, []types.PermissionLevel{{bob,common.DefaultConfig.ActiveName}}, e.DefaultExpirationDelta)
	sudoTrx := e.SudoExec(alice, &trx, 36000)
	active := common.DefaultConfig.ActiveName
	e.Propose(carol, common.N("first"), sudoTrx, []types.PermissionLevel{
		{alice, active}, {producer1, active},
		{producer2, active}, {producer3, active},
		{producer4, active}, {producer5, active},
	})
	e.Approve(alice, carol, common.N("first"), types.PermissionLevel{Actor:alice, Permission:active})

	// 2 of the 4 needed producers approve
	e.Approve(producer1, carol, common.N("first"), types.PermissionLevel{Actor:producer1, Permission:active})
	e.Approve(producer2, carol, common.N("first"), types.PermissionLevel{Actor:producer2, Permission:active})

	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	e.SetProducers(&[]common.AccountName{producer1, producer2, producer3, producer4, producer5, common.N("newprod1")})

	for len(e.Control.PendingBlockState().ActiveSchedule.Producers) != int(6) {
		e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	}

	e.Approve(producer3, carol, common.N("first"), types.PermissionLevel{Actor:producer3, Permission:active})
	e.Approve(producer4, carol, common.N("first"), types.PermissionLevel{Actor:producer4, Permission:active})
	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	exec := func() {e.Exec(alice, carol, common.N("first"))}
	CatchThrowExceptionAndMsg(t, &exception.EosioAssertMessageException{}, "transaction authorization failed", exec)

	approve := func() {e.Approve(common.N("newprod1"), carol, common.N("first"), types.PermissionLevel{Actor:common.N("newprod1"), Permission:active})}
	CatchThrowExceptionAndMsg(t, &exception.EosioAssertMessageException{}, "approval is not on the list of requested approvals", approve)

	// But prod5 still can provide the fifth approval necessary to satisfy the 2/3+1 threshold of the new producer set
	e.Approve(producer5, carol, common.N("first"), types.PermissionLevel{Actor:producer5, Permission:active})

	var traces []types.TransactionTrace
	appliedTrxCaller := chain_interface.AppliedTransactionCaller{Caller: func(t *types.TransactionTrace) {if t.Scheduled {traces = append(traces, *t)}}}
	e.Control.AppliedTransaction.Connect(&appliedTrxCaller)

	// Now the proposal should be ready to execute
	e.Exec(alice, carol, common.N("first"))
	e.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)

	assert.Equal(t, int(2), len(traces))

	assert.Equal(t, int(1), len(traces[0].ActionTraces))
	assert.Equal(t, eosioSudo, traces[0].ActionTraces[0].Act.Account)
	assert.Equal(t, common.N("exec"), traces[0].ActionTraces[0].Act.Name)
	assert.Equal(t, types.TransactionStatusExecuted, traces[0].Receipt.Status)

	assert.Equal(t, int(1), len(traces[1].ActionTraces))
	assert.Equal(t, eosio, traces[1].ActionTraces[0].Act.Account)
	assert.Equal(t, common.N("reqauth"), traces[1].ActionTraces[0].Act.Name)
	assert.Equal(t, types.TransactionStatusExecuted, traces[1].Receipt.Status)
}