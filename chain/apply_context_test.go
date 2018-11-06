package chain

import (
	"github.com/eosspark/eos-go/common"
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
