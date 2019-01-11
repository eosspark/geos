package unittests

import (
	. "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/log"
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

	e.SetProducers(&[]common.AccountName{eosioMsig, producer1, producer2, producer3, producer4, producer5})

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
	exec := common.Variants{
		"executer": executer,
		"trx":      trx,
	}
	data, _ := rlp.EncodeToBytes(exec)
	actObj := types.Action{
		Account:       eosioSudo,
		Name:          common.N("exec"),
		Authorization: auth,
		Data:          data,
	}
	trx2 := types.Transaction{}
	trx2.Actions = append(trx2.Actions, &actObj)
	e.SetTransactionHeaders(&trx2, e.DefaultExpirationDelta, 0)
	return trx2
}

func (e EosioSudoTester) ReqAuth(from common.AccountName, auths []types.PermissionLevel, expiration uint32) types.Transaction {
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
	e.SetTransactionHeaders(&trx, expiration, 0)
	return trx
}

func TestSudoExecDirect(t *testing.T) {
	e := initEosioSudoTester()
	trx := e.ReqAuth(bob, []types.PermissionLevel{{bob,common.DefaultConfig.ActiveName}}, e.DefaultExpirationDelta)
	sudoTrx := types.SignedTransaction{Transaction: trx}
	chainId := e.Control.GetChainId()
	pk := e.getPrivateKey(alice, "active")
	sudoTrx.Sign(&pk, &chainId)
	prodVector := []common.AccountName{producer1, producer2, producer3, producer4}
	for _, actor := range prodVector {
		pk = e.getPrivateKey(actor, "active")
		sudoTrx.Sign(&pk, &chainId)
	}
	trace := e.PushTransaction(&sudoTrx, common.MaxTimePoint(), e.DefaultBilledCpuTimeUs)
	assert.Equal(t, int(1), len(trace.ActionTraces))
	assert.Equal(t, eosio, trace.ActionTraces[0].Act.Account)
	assert.Equal(t, common.N("reqauth"), trace.ActionTraces[0].Act.Name)
	assert.Equal(t, types.TransactionStatusExecuted, trace.Receipt.Status)
}