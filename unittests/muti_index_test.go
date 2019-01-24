package unittests

import (
	"github.com/stretchr/testify/assert"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"io/ioutil"
	"testing"
)

func TestMultiIndexLoad(t *testing.T) {
	t.Run("", func(t *testing.T) {

		b := newBaseTester(true, chain.SPECULATIVE)
		b.ProduceBlocks(2, false)

		account := common.N("multitest")
		b.CreateAccounts([]common.AccountName{account}, false, true)
		b.ProduceBlocks(2, false)

		wasm, _ := ioutil.ReadFile("test_contracts/multi_index_test.wasm")
		abi, _ := ioutil.ReadFile("test_contracts/multi_index_test.abi")
		b.SetCode(account, wasm, nil)
		b.SetAbi(account, abi, nil)
		b.ProduceBlocks(1, false)

		var trxId1 common.BlockIdType
		{
			trx := types.SignedTransaction{}
			actData := common.Variants{"what": 0}
			act := b.GetAction(account,
				common.N("trigger"),
				[]common.PermissionLevel{{account, common.DefaultConfig.ActiveName}},
				&actData)

			trx.Actions = append(trx.Actions, act)
			b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

			privKey := b.getPrivateKey(account, "active")
			chainId := b.Control.GetChainId()
			trx.Sign(&privKey, &chainId)
			b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
			trxId1 = trx.ID()
		}

		var trxId2 common.BlockIdType
		{
			trx := types.SignedTransaction{}
			actData := common.Variants{"what": 1}
			act := b.GetAction(account,
				common.N("trigger"),
				[]common.PermissionLevel{{account, common.DefaultConfig.ActiveName}},
				&actData)

			trx.Actions = append(trx.Actions, act)
			b.SetTransactionHeaders(&trx.Transaction, b.DefaultExpirationDelta, 0)

			privKey := b.getPrivateKey(account, "active")
			chainId := b.Control.GetChainId()
			trx.Sign(&privKey, &chainId)
			b.PushTransaction(&trx, common.MaxTimePoint(), b.DefaultBilledCpuTimeUs)
			trxId2 = trx.ID()
		}
		b.ProduceBlocks(1, false)

		assert.Equal(t, b.ChainHasTransaction(&trxId1), true)
		assert.Equal(t, b.ChainHasTransaction(&trxId2), true)

		b.close()

	})
}
