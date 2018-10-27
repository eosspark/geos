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

func TestIteratorCache(t *testing.T) {

	t.Run("", func(t *testing.T) {

		i := NewIteratorCache()

		var firstTable *entity.TableIdObject
		var itrTable int = 0
		var k int = 0
		for k = 0; k < 8; k++ {
			table := entity.TableIdObject{
				ID:    common.IdType(k),
				Code:  common.AccountName(common.N("eosio.token")),
				Scope: common.ScopeName(common.N("xiaoyu")),
				Table: common.TableName(common.N("accounts")),
				Payer: common.AccountName(k),
			}

			iterator := i.cacheTable(&table)

			if k == 4 {
				itrTable = iterator
				firstTable = &table
			}
		}

		var itrKeyvalue int = 0
		var j int
		for j = 0; j < 31; j++ {
			value, _ := rlp.EncodeToBytes(common.N("walker"))
			keyvalue := entity.KeyValueObject{
				ID:         common.IdType(k + j),
				TId:        firstTable.ID,
				PrimaryKey: uint64(j),
				Payer:      common.AccountName(j),
				Value:      value,
			}

			itr := i.add(&keyvalue)
			if j == 10 {
				itrKeyvalue = itr
			}
		}

		obj := (i.get(itrKeyvalue)).(*entity.KeyValueObject)
		objTable := i.getTable(obj.TId)
		assert.Equal(t, firstTable.Payer, objTable.Payer)

		var name uint64
		rlp.DecodeBytes(obj.Value, &name)
		assert.Equal(t, common.S(name), "walker")

		itr := i.getEndIteratorByTableID(firstTable.ID)
		assert.Equal(t, itrTable, itr)

		objTable = i.findTablebyEndIterator(itr)
		assert.Equal(t, firstTable.Payer, objTable.Payer)

		i.remove(itrKeyvalue)
		//		obj = (i.get(itrKeyvalue)).(*entity.KeyValueObject)

		value, _ := rlp.EncodeToBytes(common.N("take"))
		keyvalue := entity.KeyValueObject{
			ID:         common.IdType(k + j),
			TId:        firstTable.ID,
			PrimaryKey: uint64(j),
			Payer:      common.AccountName(j),
			Value:      value,
		}
		itrKeyvalue = i.add(&keyvalue)
		obj = (i.get(itrKeyvalue)).(*entity.KeyValueObject)
		assert.Equal(t, keyvalue.ID, obj.ID)

	})

}

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
				types.PermissionLevel{Actor: common.AccountName(common.N("eosio")), Permission: common.PermissionName(common.N("active"))},
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

		var itrStore int
		for i := 0; i < 10; i++ {
			itrStore = a.DbStoreI64(int64(common.N("xiaoyu")), int64(common.N("accounts")), int64(common.N("eosio")), int64(i), []byte{byte(i)})
			if i == 1 {
				itr = itrStore
			}

		}

		a.DbFindI64(int64(common.N("eosio")), int64(common.N("xiaoyu")), int64(common.N("accounts")), 5)

		//itrFind := a.DbFindI64(int64(common.N("eosio")), int64(common.N("xiaoyu")), int64(common.N("accounts")), 5)
		//var primary uint64
		//itr = a.DbPreviousI64(itrFind, &primary)
		//itr = a.DbNextI64(itrFind, &primary)

		//assert.Equal(t, itrFind, itr)

		control.Close()
		control.Clean()

	})

}
