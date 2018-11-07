package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	arithmetic "github.com/eosspark/eos-go/common/arithmetic_types"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	"github.com/stretchr/testify/assert"
	"math"
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

		// trxHeader := types.TransactionHeader{
		// 	Expiration:       common.MaxTimePointSec(),
		// 	RefBlockNum:      4,
		// 	RefBlockPrefix:   3832731038,
		// 	MaxNetUsageWords: 0,
		// 	MaxCpuUsageMS:    0,
		// 	DelaySec:         0,
		// }

		// trx := types.Transaction{
		// 	TransactionHeader:     trxHeader,
		// 	ContextFreeActions:    []*types.Action{},
		// 	Actions:               []*types.Action{&act},
		// 	TransactionExtensions: []*types.Extension{},
		// }
		// signedTrx := types.NewSignedTransaction(&trx, []ecc.Signature{}, []common.HexBytes{})
		// privateKey, _ := ecc.NewRandomPrivateKey()
		// chainIdType := common.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
		// signedTrx.Sign(privateKey, &chainIdType)
		// trxContext := NewTransactionContext(control, signedTrx, trx.ID(), common.Now())

		a := newApplyContext(control, &act)

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

	})

}

func TestDbSecondaryKeyIdx64(t *testing.T) {

	t.Run("", func(t *testing.T) {

		control := GetControllerInstance()
		blockTimeStamp := common.NewBlockTimeStamp(common.Now())
		control.StartBlock(blockTimeStamp, 0)

		account1 := "hello"
		createNewAccount(control, account1)

		buffer, _ := rlp.EncodeToBytes("0123456")
		act := types.Action{
			Account: common.AccountName(common.N(account1)),
			Name:    common.ActionName(common.N("hi")),
			Data:    buffer,
			Authorization: []types.PermissionLevel{
				//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
				{Actor: common.AccountName(common.N("eosio")), Permission: common.PermissionName(common.N("active"))},
			},
		}

		a := newApplyContext(control, &act)

		for i := 0; i < 10; i++ {
			primaryKey := int64(i + 1)
			secondaryKey := uint64(i + 1)

			a.DbStoreI64(int64(common.N("xiaoyu")), int64(common.N("accounts")), int64(common.N("eosio")), primaryKey, []byte{byte(i)})
			a.Idx64Store(int64(common.N("xiaoyu")), int64(common.N("accounts")), int64(common.N("eosio")), primaryKey, &secondaryKey)
		}

		var secondaryKey uint64
		primaryKey := uint64(5)
		itrFind := a.Idx64FindPrimary(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)
		assert.Equal(t, secondaryKey, uint64(5))

		primaryKey = uint64(1)
		itrFind = a.Idx64FindPrimary(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)

		secondaryKey = uint64(5)
		itrFind = a.Idx64Lowerbound(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)
		a.Idx64Remove(itrFind)

		secondaryKey = uint64(5)
		itrFind = a.Idx64Lowerbound(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)
		assert.Equal(t, secondaryKey, uint64(6))

		itrFind = a.Idx64Upperbound(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)
		secondaryKey = uint64(700)
		a.Idx64Update(itrFind, int64(common.N("eosio")), &secondaryKey)
		secondaryKey = uint64(7)
		primaryKey = uint64(7)
		itrFind = a.Idx64FindPrimary(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)
		assert.Equal(t, secondaryKey, uint64(700))

		secondaryKey = uint64(6)
		itrFind = a.Idx64FindSecondary(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)
		assert.Equal(t, primaryKey, uint64(6))

		itrFind = a.Idx64Next(itrFind, &primaryKey)
		assert.Equal(t, primaryKey, uint64(8))
		secondaryKey = uint64(9)
		itrFind = a.Idx64FindSecondary(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)
		itrFind = a.Idx64Next(itrFind, &primaryKey)
		itrFind = a.Idx64Next(itrFind, &primaryKey)
		assert.Equal(t, primaryKey, uint64(7))

		for i := 0; i < 5; i++ {
			itrFind = a.Idx64Previous(itrFind, &primaryKey)
		}
		assert.Equal(t, primaryKey, uint64(4))
		for i := 0; i < 4; i++ {
			itrFind = a.Idx64Previous(itrFind, &primaryKey)
		}
		assert.Equal(t, primaryKey, uint64(1))

		control.Close()

	})
}

func TestDbSecondaryKeyIdx64_2(t *testing.T) {

	t.Run("", func(t *testing.T) {

		control := GetControllerInstance()
		blockTimeStamp := common.NewBlockTimeStamp(common.Now())
		control.StartBlock(blockTimeStamp, 0)

		account1 := "hello"
		createNewAccount(control, account1)

		buffer, _ := rlp.EncodeToBytes("0123456")
		act := types.Action{
			Account: common.AccountName(common.N(account1)),
			Name:    common.ActionName(common.N("hi")),
			Data:    buffer,
			Authorization: []types.PermissionLevel{
				//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
				{Actor: common.AccountName(common.N("eosio")), Permission: common.PermissionName(common.N("active"))},
			},
		}

		a := newApplyContext(control, &act)

		for i := 0; i < 10; i++ {
			primaryKey := int64(i + 1)
			secondaryKey := uint64(i + 1)

			a.DbStoreI64(int64(common.N("xiaoyu")), int64(common.N("accounts")), int64(common.N("eosio")), primaryKey, []byte{byte(i)})
			a.Idx64Store(int64(common.N("xiaoyu")), int64(common.N("accounts")), int64(common.N("eosio")), primaryKey, &secondaryKey)
		}

		var secondaryKey uint64
		primaryKey := uint64(5)
		itrFind := a.idx64.lowerboundPrimary(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &primaryKey)
		a.idx64.get(itrFind, &secondaryKey, &primaryKey)
		assert.Equal(t, secondaryKey, uint64(5))

		primaryKey = uint64(5)
		itrFind = a.idx64.upperboundPrimary(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &primaryKey)
		a.idx64.get(itrFind, &secondaryKey, &primaryKey)
		assert.Equal(t, secondaryKey, uint64(6))

		for i := 0; i < 4; i++ {
			itrFind = a.idx64.nextPrimary(itrFind, &primaryKey)
		}
		assert.Equal(t, secondaryKey, uint64(10))

		for i := 0; i < 10; i++ {
			itrFind = a.idx64.previousPrimary(itrFind, &primaryKey)
		}
		assert.Equal(t, secondaryKey, uint64(1))

		control.Close()

	})
}

func TestDbSecondaryKeyDouble(t *testing.T) {

	t.Run("", func(t *testing.T) {

		control := GetControllerInstance()
		blockTimeStamp := common.NewBlockTimeStamp(common.Now())
		control.StartBlock(blockTimeStamp, 0)

		account1 := "hello"
		createNewAccount(control, account1)

		buffer, _ := rlp.EncodeToBytes("0123456")
		act := types.Action{
			Account: common.AccountName(common.N(account1)),
			Name:    common.ActionName(common.N("hi")),
			Data:    buffer,
			Authorization: []types.PermissionLevel{
				//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
				{Actor: common.AccountName(common.N("eosio")), Permission: common.PermissionName(common.N("active"))},
			},
		}

		a := newApplyContext(control, &act)

		for i := 0; i < 10; i++ {
			primaryKey := int64(i + 1)
			secondaryKey := arithmetic.Float64(math.Float64bits(float64(i+1) + 1.5))

			a.DbStoreI64(int64(common.N("xiaoyu")), int64(common.N("accounts")), int64(common.N("eosio")), primaryKey, []byte{byte(i)})
			a.IdxDoubleStore(int64(common.N("xiaoyu")), int64(common.N("accounts")), int64(common.N("eosio")), primaryKey, &secondaryKey)
		}

		var secondaryKey arithmetic.Float64
		primaryKey := uint64(5)
		itrFind := a.IdxDoubleFindPrimary(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)
		assert.Equal(t, secondaryKey, arithmetic.Float64(math.Float64bits(float64(5)+1.5)))

		primaryKey = uint64(1)
		itrFind = a.IdxDoubleFindPrimary(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)
		fmt.Println(math.Float64frombits(uint64(secondaryKey)))

		secondaryKey = arithmetic.Float64(math.Float64bits(float64(5+1) + 1.5))
		itrFind = a.IdxDoubleLowerbound(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)

		a.IdxDoubleRemove(itrFind)

		secondaryKey = arithmetic.Float64(math.Float64bits(float64(5+1) + 1.5))
		itrFind = a.IdxDoubleUpperbound(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)
		secondaryKey = arithmetic.Float64(math.Float64bits(float64(700+1) + 1.5))
		a.IdxDoubleUpdate(itrFind, int64(common.N("eosio")), &secondaryKey)
		secondaryKey = arithmetic.Float64(math.Float64bits(float64(5+1) + 1.5))
		primaryKey = uint64(7)
		itrFind = a.IdxDoubleFindPrimary(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)
		fmt.Println(math.Float64frombits(uint64(secondaryKey)))
		assert.Equal(t, secondaryKey, arithmetic.Float64(math.Float64bits(float64(700+1)+1.5)))

		secondaryKey = arithmetic.Float64(math.Float64bits(float64(4+1) + 1.5))
		itrFind = a.IdxDoubleFindSecondary(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &secondaryKey, &primaryKey)
		assert.Equal(t, primaryKey, uint64(5))

		for i := 0; i < 4; i++ {
			itrFind = a.IdxDoublePrevious(itrFind, &primaryKey)
		}
		assert.Equal(t, primaryKey, uint64(1))

		for i := 0; i < 9; i++ {
			itrFind = a.IdxDoubleNext(itrFind, &primaryKey)
		}
		assert.Equal(t, primaryKey, uint64(7))
		control.Close()

	})
}

func TestDbSecondaryKeyIdxDouble_2(t *testing.T) {

	t.Run("", func(t *testing.T) {

		control := GetControllerInstance()
		blockTimeStamp := common.NewBlockTimeStamp(common.Now())
		control.StartBlock(blockTimeStamp, 0)

		account1 := "hello"
		createNewAccount(control, account1)

		buffer, _ := rlp.EncodeToBytes("0123456")
		act := types.Action{
			Account: common.AccountName(common.N(account1)),
			Name:    common.ActionName(common.N("hi")),
			Data:    buffer,
			Authorization: []types.PermissionLevel{
				//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
				{Actor: common.AccountName(common.N("eosio")), Permission: common.PermissionName(common.N("active"))},
			},
		}

		a := newApplyContext(control, &act)

		for i := 0; i < 10; i++ {
			primaryKey := int64(i + 1)
			secondaryKey := arithmetic.Float64(math.Float64bits(float64(i+1) + 1.5))

			a.DbStoreI64(int64(common.N("xiaoyu")), int64(common.N("accounts")), int64(common.N("eosio")), primaryKey, []byte{byte(i)})
			a.IdxDoubleStore(int64(common.N("xiaoyu")), int64(common.N("accounts")), int64(common.N("eosio")), primaryKey, &secondaryKey)
		}

		var secondaryKey arithmetic.Float64
		primaryKey := uint64(5)
		itrFind := a.idxDouble.lowerboundPrimary(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &primaryKey)
		a.idxDouble.get(itrFind, &secondaryKey, &primaryKey)
		assert.Equal(t, secondaryKey, arithmetic.Float64(math.Float64bits(float64(4+1)+1.5)))

		primaryKey = uint64(5)
		itrFind = a.idxDouble.upperboundPrimary(int64(common.N(account1)), int64(common.N("xiaoyu")), int64(common.N("accounts")), &primaryKey)
		a.idxDouble.get(itrFind, &secondaryKey, &primaryKey)
		assert.Equal(t, secondaryKey, arithmetic.Float64(math.Float64bits(float64(5+1)+1.5)))

		for i := 0; i < 4; i++ {
			itrFind = a.idxDouble.nextPrimary(itrFind, &primaryKey)
		}
		assert.Equal(t, primaryKey, uint64(10))

		for i := 0; i < 10; i++ {
			itrFind = a.idxDouble.previousPrimary(itrFind, &primaryKey)
		}
		assert.Equal(t, primaryKey, uint64(1))

		control.Close()

	})
}
