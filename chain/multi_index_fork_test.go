package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initMulti() (*multiIndexFork, *types.BlockState) {
	initPriKey, _ := ecc.NewPrivateKey("5KYZdUEo39z3FPrtuX2QbbwGnNP5zTd7yyr2SC1j299sBCnWjss")
	initPubKey := initPriKey.PublicKey()
	eosio := common.AccountName(common.N("eosio"))
	eos := common.AccountName(common.N("eos"))
	tester := common.AccountName(common.N("tester"))

	initSchedule := types.ProducerScheduleType{0, []types.ProducerKey{
		{eosio, initPubKey},
		{eos, initPubKey},
		{tester, initPubKey},
	}}

	genHeader := new(types.BlockHeaderState)
	genHeader.ActiveSchedule = initSchedule
	genHeader.PendingSchedule = initSchedule
	genHeader.Header.Timestamp = common.BlockTimeStamp(1162425600) //slot of 2018-6-2 00:00:00:000
	genHeader.BlockId = genHeader.Header.BlockID()
	genHeader.BlockNum = genHeader.Header.BlockNumber()

	genHeader.BlockSigningKey = initPubKey

	blockState := types.NewBlockState(genHeader)
	blockState.SignedBlock = new(types.SignedBlock)
	blockState.SignedBlock.SignedBlockHeader = genHeader.Header
	blockState.InCurrentChain = true
	mi := newMultiIndexFork()
	mi.Insert(blockState)

	return mi, blockState
}

func TestMultiIndexFork_Insert(t *testing.T) {
	mi, bs := initMulti()

	//fmt.Println("byBlockId", mi.indexs["byBlockId"].Value.Compare == nil )

	for i := 0; i < 10; i++ {
		t := 1162425602 + 200
		tmp := common.BlockTimeStamp(t)
		bhs := bs.GenerateNext(tmp)
		bhs.BlockId = bhs.Header.BlockID()
		blockState := types.NewBlockState(bhs)
		mi.Insert(blockState)
	}

	log.Info("%v", mi.indexs["byBlockId"].Value.Len())

	assert.Equal(t, 11, mi.indexs["byBlockId"].Value.Len())
	result := mi.Find(bs.BlockId)
	assert.Equal(t, bs, result)

}

func TestMultiIndexFork_Find(t *testing.T) {
	mi, bs := initMulti()
	result := mi.Find(bs.BlockId)
	assert.Equal(t, bs, result)
}

func TestIndexFork_LowerBound(t *testing.T) {
	mi, bs := initMulti()
	var tm *types.BlockState
	for i := 0; i < 10; i++ {
		t := 1162425602 + 200
		tmp := common.BlockTimeStamp(t)
		bhs := bs.GenerateNext(tmp)
		bhs.BlockId = bhs.Header.BlockID()
		blockState := types.NewBlockState(bhs)
		blockState.InCurrentChain = true
		fmt.Println("insert:", blockState.BlockNum)
		tm = blockState
		mi.Insert(blockState)
	}

	idxFork := mi.indexs["byBlockNum"]
	for i := 0; i < idxFork.Value.Len(); i++ {
		fmt.Println(idxFork.Value.Data[i].(*types.BlockState).BlockNum)
	}

	itr := idxFork.LowerBound(tm)
	fmt.Println("result:", itr.KeySet.Len(), bs.BlockNum)
	assert.Equal(t, 10, itr.KeySet.Len())
}

func TestIndexFork_UpperBound(t *testing.T) {
	mi, bs := initMulti()
	//var tm *types.BlockState
	for i := 0; i < 10; i++ {
		t := 1162425602 + 200
		tmp := common.BlockTimeStamp(t)
		bhs := bs.GenerateNext(tmp)
		bhs.BlockId = bhs.Header.BlockID()
		blockState := types.NewBlockState(bhs)
		blockState.InCurrentChain = true
		fmt.Println("insert:", blockState.BlockNum)
		//	tm = blockState
		mi.Insert(blockState)
	}

	idxFork := mi.indexs["byBlockNum"]
	for i := 0; i < idxFork.Value.Len(); i++ {
		fmt.Println(idxFork.Value.Data[i].(*types.BlockState).BlockNum)
	}

	itr := idxFork.UpperBound(bs)
	fmt.Println("result:", itr.KeySet.Len(), bs.BlockNum)
	assert.Equal(t, 10, itr.KeySet.Len())
}

func Search(n int, f func(int) bool) int {

	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1)
		if !f(h) {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
