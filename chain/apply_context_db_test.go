package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDbPrimaryKey(t *testing.T) {

	t.Run("", func(t *testing.T) {

		control := GetControllerInstance()
		blockTimeStamp := common.NewBlockTimeStamp(common.Now())
		control.StartBlock(blockTimeStamp, 0)

		buffer, _ := rlp.EncodeToBytes("0123456")
		act := types.Action{
			Account: common.AccountName(common.N("eosio")),
			Name:    common.ActionName(common.N("hello")),
			Data:    buffer,
			Authorization: []types.PermissionLevel{
				//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
				{Actor: common.AccountName(common.N("eosio")), Permission: common.PermissionName(common.N("active"))},
			},
		}

		trxHeader := types.TransactionHeader{
			Expiration:       common.MaxTimePointSec(),
			RefBlockNum:      4,
			RefBlockPrefix:   3832731038,
			MaxNetUsageWords: 0,
			MaxCpuUsageMS:    0,
			DelaySec:         0,
		}

		trx := types.Transaction{
			TransactionHeader:     trxHeader,
			ContextFreeActions:    []*types.Action{},
			Actions:               []*types.Action{&act},
			TransactionExtensions: []*types.Extension{},
		}
		signedTrx := types.NewSignedTransaction(&trx, []ecc.Signature{}, []common.HexBytes{})
		privateKey, _ := ecc.NewRandomPrivateKey()
		chainIdType := common.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
		signedTrx.Sign(privateKey, &chainIdType)
		trxContext := NewTransactionContext(control, signedTrx, trx.ID(), common.Now())

		a := NewApplyContext(control, trxContext, &act, 0)

		//DbStoreI64
		itr := a.DbStoreI64(int64(common.N("xiaoyu")), int64(common.N("accounts")), int64(common.N("eosio")), 1, buffer)

		tab := a.FindTable(int64(common.N("eosio")), int64(common.N("xiaoyu")), int64(common.N("accounts")))
		fmt.Println(common.S(uint64(tab.Code)))

		obj := (a.KeyvalCache.get(itr)).(*entity.KeyValueObject)
		assert.Equal(t, []byte(obj.Value), buffer)

		//DbUpdateI64
		buffer, _ = rlp.EncodeToBytes("0123456789")
		a.DbUpdateI64(itr, int64(common.N("eosio")), buffer)

		//DbGetI64
		buf := make([]byte, 11)
		bufferSize := a.DbGetI64(itr, buf, 11)

		var ret string
		rlp.DecodeBytes(buf, &ret)
		assert.Equal(t, ret, "0123456789")
		assert.Equal(t, bufferSize, 11)

		//DbRemoveI64
		a.DbRemoveI64(itr)

		for i := 0; i < 10; i++ {
			a.DbStoreI64(int64(common.N("xiaoyu")), int64(common.N("accounts")), int64(common.N("eosio")), int64(i), []byte{byte(i)})

		}

		itrFind := a.DbFindI64(int64(common.N("eosio")), int64(common.N("xiaoyu")), int64(common.N("accounts")), 5)

		var primary uint64
		itr = a.DbNextI64(itrFind, &primary)
		assert.Equal(t, primary, uint64(6))

		itr = a.DbPreviousI64(itrFind, &primary)
		assert.Equal(t, primary, uint64(4))

		itrLowerbound := a.DbLowerboundI64(int64(common.N("eosio")), int64(common.N("xiaoyu")), int64(common.N("accounts")), 5)
		bufferBound := []byte{byte(20)}
		a.DbGetI64(itrLowerbound, bufferBound, len(bufferBound))
		assert.Equal(t, bufferBound, []byte{byte(5)})

		itrUpperbound := a.DbUpperboundI64(int64(common.N("eosio")), int64(common.N("xiaoyu")), int64(common.N("accounts")), 6)
		a.DbGetI64(itrUpperbound, bufferBound, len(bufferBound))
		assert.Equal(t, bufferBound, []byte{byte(7)})

		control.Close()
		control.Clean()

	})

}

func TestDbSecondaryKeyIdx64(t *testing.T) {

	t.Run("", func(t *testing.T) {

	})
}
