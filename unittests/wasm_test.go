package unittests

import (
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

type assertdef struct {
	Condition int8
	Message   string
}

func (d *assertdef) getAccount() common.AccountName {
	return common.N("asserter")
}

func (d *assertdef) getName() common.AccountName {
	return common.N("procassert")
}

type provereset struct{}

func (d *provereset) getAccount() common.AccountName {
	return common.N("asserter")
}

func (d *provereset) getName() common.AccountName {
	return common.N("procassert")
}

type actionInterface interface {
	getAccount() common.AccountName
	getName() common.AccountName
}

func newAction(permissionLevel []types.PermissionLevel, a actionInterface) *types.Action {

	payload, _ := rlp.EncodeToBytes(a)
	act := types.Action{
		Account:       common.AccountName(a.getAccount()),
		Name:          common.AccountName(a.getName()),
		Data:          payload,
		Authorization: permissionLevel,
	}
	return &act
}

func TestBasic(t *testing.T) {
	name := "test_contracts/asserter.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		asserter := common.N("asserter")
		procassert := common.N("procassert")

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)
		b.CreateAccounts([]common.AccountName{asserter}, false, true)
		b.ProduceBlocks(1, false)
		b.SetCode(asserter, code, nil)
		b.ProduceBlocks(1, false)
		{
			trx := types.SignedTransaction{}
			pl := []types.PermissionLevel{{asserter, common.DefaultConfig.ActiveName}}
			action := assertdef{1, "Should Not Assert!"}
			act := newAction(pl, &action)
			trx.Actions = append(trx.Actions, act)
			b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

			privKey := b.getPrivateKey(asserter, "active")
			chainId := b.Control.GetChainId()
			trx.Sign(&privKey, &chainId)

			result := b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
			assert.Equal(t, result.Receipt.Status, types.TransactionStatusExecuted)
			assert.Equal(t, len(result.ActionTraces), 1)
			assert.Equal(t, result.ActionTraces[0].Receipt.Receiver, asserter)
			assert.Equal(t, result.ActionTraces[0].Act.Account, asserter)
			assert.Equal(t, result.ActionTraces[0].Act.Name, procassert)
			assert.Equal(t, len(result.ActionTraces[0].Act.Authorization), 1)
			assert.Equal(t, result.ActionTraces[0].Act.Authorization[0].Actor, asserter)
			assert.Equal(t, result.ActionTraces[0].Act.Authorization[0].Permission, common.DefaultConfig.ActiveName)
		}

		b.close()

	})

}
