package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
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
	genHeader.ProducerToLastProduced = *types.NewAccountNameUint32Map()
	genHeader.ProducerToLastImpliedIrb = *types.NewAccountNameUint32Map()
	genHeader.BlockSigningKey = initPubKey
	genHeader.Header.ProducerSignature = *ecc.NewSigNil()
	blockState := types.NewBlockState(genHeader)
	blockState.SignedBlock = new(types.SignedBlock)
	blockState.SignedBlock.SignedBlockHeader = genHeader.Header
	blockState.Header.ProducerSignature = *ecc.NewSigNil()
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

	log.Info("%v", mi.Indexs["byBlockId"].Value.Size())

	assert.Equal(t, 11, mi.Indexs["byBlockId"].Value.Size())
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

	itr := idxFork.Value.LowerBound(tm)
	log.Info("tm:%#v", tm)
	log.Info("current sub:%#v", itr.Value().(*types.BlockState).BlockNum)
	assert.Equal(t, tm.BlockNum, itr.Value().(*types.BlockState).BlockNum)

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

	itr := idxFork.Value.UpperBound(bs)
	log.Info("result:%#v,%v", itr.Value().(*types.BlockState), tm.BlockNum)
	assert.Equal(t, uint32(2), itr.Value().(*types.BlockState).BlockNum)
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

	assert.Equal(t, uint32(2), itr.Value().(*types.BlockState).BlockNum)
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
	assert.Equal(t, uint32(1), itr.Value().(*types.BlockState).BlockNum)
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
	val := idxFork.Value.Iterator()
	val.Next()
	assert.Equal(t, val.Value().(*types.BlockState), obj)
}

func Test_LowerBound_NotFound(t *testing.T) {
	mi, bs := initMulti()
	b := types.BlockState{}
	b.BlockNum = 100
	numIdx := mi.GetIndex("byBlockNum")
	itr := numIdx.Value.LowerBound(&b)
	try.Try(func() {
		val := itr.Value()
		if itr != nil || val.(*types.BlockState).BlockNum != bs.BlockNum || val.(*types.BlockState).InCurrentChain != true {
		}
		fmt.Println("exec there is error")
	}).Catch(func(e error) {
		assert.Errorf(t, e, "invalid memory address or nil pointer dereference")
	})
}
