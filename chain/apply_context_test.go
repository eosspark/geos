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
		var k int = 0
		for k = 0; k < 8; k++ {
			table := entity.TableIdObject{
				ID:    common.IdType(k),
				Code:  common.AccountName(common.N("eosio.token")),
				Scope: common.ScopeName(common.N("xiaoyu")),
				Table: common.TableName(common.N("accounts")),
				Payer: common.AccountName(k),
			}

			i.cacheTable(&table)

			if k == 0 {
				firstTable = &table
			}
		}

		var itrKeyvalue int = 0
		for j := 0; j < 32; j++ {
			value, _ := rlp.EncodeToBytes(common.N("walker"))
			keyvalue := entity.KeyValueObject{
				ID:         common.IdType(k + j),
				TId:        firstTable.ID,
				PrimaryKey: uint64(j),
				Payer:      common.AccountName(j),
				Value:      value,
			}

			if j == 31 {
				itrKeyvalue = i.add(&keyvalue)
			}
		}

		obj := (i.get(itrKeyvalue)).(*entity.KeyValueObject)
		objTable := i.getTable(obj.TId)

		assert.Equal(t, firstTable.Payer, objTable.Payer)

		var name uint64
		rlp.DecodeBytes(obj.Value, &name)
		assert.Equal(t, common.S(name), "walker")

		itr := i.getEndIteratorByTableID(firstTable.ID)
		objTable = i.findTablebyEndIterator(itr)

		assert.Equal(t, firstTable.Payer, objTable.Payer)

	})

}
