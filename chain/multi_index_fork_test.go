package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initMulti() (*MultiIndexFork, *types.BlockState) {
	initPriKey, _ := ecc.NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
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
	genHeader.Header.Timestamp = types.BlockTimeStamp(1162425600) //slot of 2018-6-2 00:00:00:000
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

func TestMultiIndexFork_Insert_Repeat(t *testing.T) {
	mi, bs := initMulti()

	for i := 0; i < 10; i++ {
		t := 1162425602 + 200
		tmp := types.BlockTimeStamp(t)
		bhs := bs.GenerateNext(tmp)
		bhs.BlockId = bhs.Header.BlockID()
		blockState := types.NewBlockState(bhs)
		mi.Insert(blockState)
	}

	log.Info("%v", mi.Indexs["byBlockId"].Value.Len())

	assert.Equal(t, 11, mi.Indexs["byBlockId"].Value.Len())
	result := mi.find(bs.BlockId)
	log.Info("%v", result)
	assert.Equal(t, bs, result)

}

func TestMultiIndexFork_FindById(t *testing.T) {
	mi, bs := initMulti()
	result := mi.find(bs.BlockId)
	assert.Equal(t, bs, result)
}

func TestIndexFork_LowerBound_byBlockNum(t *testing.T) {
	mi, bs := initMulti()
	var tm *types.BlockState
	for i := 0; i < 10; i++ {
		t := 1162425602 + 200
		tmp := types.BlockTimeStamp(t)
		bhs := bs.GenerateNext(tmp)
		bhs.BlockId = bhs.Header.BlockID()
		blockState := types.NewBlockState(bhs)
		blockState.InCurrentChain = true
		tm = blockState
		mi.Insert(blockState)
	}

	idxFork := mi.Indexs["byBlockNum"]

	val, sub := idxFork.Value.LowerBound(tm)
	log.Info("tm:%#v", tm)
	log.Info("current sub:%#v", sub)
	assert.Equal(t, 1, sub)
	assert.Equal(t, tm, val)
}

func TestIndexFork_UpperBound_byBlockNum(t *testing.T) {
	mi, bs := initMulti()
	var tm *types.BlockState
	for i := 0; i < 10; i++ {
		t := 1162425602 + 200
		tmp := types.BlockTimeStamp(t)
		bhs := bs.GenerateNext(tmp)
		bhs.BlockId = bhs.Header.BlockID()
		blockState := types.NewBlockState(bhs)
		blockState.InCurrentChain = true
		//fmt.Println("insert:", blockState.BlockNum)
		if i == 9 {
			tm = blockState
		}

		mi.Insert(blockState)
	}

	idxFork := mi.Indexs["byBlockNum"]

	val, sub := idxFork.Value.UpperBound(bs)
	log.Info("result:%#v,%v", sub, tm.BlockNum)
	assert.Equal(t, uint32(1), val.BlockNum)
}

func TestMultiIndexFork_LowerBound_lib(t *testing.T) {
	mi, bs := initMulti()
	var tm *types.BlockState
	for i := 0; i < 10; i++ {
		t := 1162425602 + 200
		tmp := types.BlockTimeStamp(t)
		bhs := bs.GenerateNext(tmp)
		bhs.BlockId = bhs.Header.BlockID()
		blockState := types.NewBlockState(bhs)
		blockState.InCurrentChain = true
		//fmt.Println("insert:", blockState.BlockNum)
		if i == 9 {
			tm = blockState
		}

		mi.Insert(blockState)
	}

	idxFork := mi.Indexs["byLibBlockNum"]
	itr := idxFork.lowerBound(tm)
	//fmt.Println(itr.idx.value.Data[itr.currentSub])
	assert.Equal(t, 0, itr.CurrentSub)
	assert.Equal(t, uint32(2), itr.Value.BlockNum)
}

func TestMultiIndexFork_UpperBound_lib(t *testing.T) {
	mi, bs := initMulti()
	var tm *types.BlockState
	for i := 0; i < 10; i++ {
		t := 1162425602 + 200
		tmp := types.BlockTimeStamp(t)
		bhs := bs.GenerateNext(tmp)
		bhs.BlockId = bhs.Header.BlockID()
		blockState := types.NewBlockState(bhs)
		blockState.InCurrentChain = true
		//fmt.Println("insert:", blockState.BlockNum)
		if i == 9 {
			tm = blockState
		}

		mi.Insert(blockState)
	}

	idxFork := mi.Indexs["byLibBlockNum"]
	itr := idxFork.upperBound(tm)
	//fmt.Println(itr.idx.value.Data[itr.currentSub])
	assert.Equal(t, 9, itr.CurrentSub)
	assert.Equal(t, uint32(2), itr.Value.BlockNum)
}

func TestIndexFork_Begin(t *testing.T) {
	mi, bs := initMulti()
	for i := 0; i < 10; i++ {
		t := 1162425602 + 200
		tmp := types.BlockTimeStamp(t)
		bhs := bs.GenerateNext(tmp)
		bhs.BlockId = bhs.Header.BlockID()
		blockState := types.NewBlockState(bhs)
		blockState.InCurrentChain = true
		//fmt.Println("insert:", blockState.BlockNum)

		mi.Insert(blockState)
	}

	idxFork := mi.Indexs["byLibBlockNum"]

	obj, _ := idxFork.Begin()

	assert.Equal(t, idxFork.Value.Data[0], obj)
}

func Test_LowerBound_NotFound(t *testing.T) {
	mi, bs := initMulti()
	b := types.BlockState{}
	b.BlockNum = bs.BlockNum
	numIdx := mi.GetIndex("byBlockNum")
	bs.BlockNum = 100
	obj, _ := numIdx.Value.LowerBound(&b)
	try.Try(func() {
		//obj := val.(*types.BlockState)
		if obj != nil || obj.BlockNum != bs.BlockNum || obj.InCurrentChain != true {
			//return &types.BlockState{}
		}
		fmt.Println(obj)
	}).Catch(func(ex exception.Exception) { //TODO catch exception code
		assert.Equal(t, 3100002, int(ex.Code()))
	}).End()
}
