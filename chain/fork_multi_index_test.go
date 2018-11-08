package chain

import (
	"encoding/binary"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initMulti() (*ForkMultiIndex, *types.BlockState) {
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

	blockState := types.NewBlockState(*genHeader)
	blockState.SignedBlock = new(types.SignedBlock)
	blockState.SignedBlock.SignedBlockHeader = genHeader.Header
	blockState.InCurrentChain = true
	mi := NewForkMultiIndex()
	mi.Insert(blockState)

	return mi, blockState
}

func GenerateNext() (*ForkMultiIndex, *types.BlockState) {
	mi, bs := initMulti()

	//t100 := common.BlockTimeStamp(1162425600 + 100)
	t2 := common.BlockTimeStamp(1162425602)

	bhs := bs.GenerateNext(t2)
	bhs.BlockId = bhs.Header.BlockID()
	blockState := types.NewBlockState(*bhs)
	mi.Insert(blockState)
	log.Info("%v", mi.Indexs["byBlockId"].KeyValue.Len())
	return mi, bs
}

func TestNewForkContainer_Insert(t *testing.T) {
	mi, bs := initMulti()
	idx := mi.GetIndex("byBlockId")
	b := idx.Begin()
	idxEle, _ := b.KeySet.FindData(b.Key)
	val := idxEle.(*IndexElement).value

	assert.Equal(t, bs, val.(*types.BlockState))

}

func TestNewForkContainer_Begin(t *testing.T) {
	mi, bs := initMulti()
	idx := mi.GetIndex("byBlockId")
	b := idx.Begin()
	assert.Equal(t, bs, b.KeySet.Data[0].(*IndexElement).value)
	prevIdx := mi.GetIndex("byPrev")
	pit := prevIdx.Begin()
	idxEle, _ := pit.KeySet.FindData(pit.Key)
	uniqueKey := idxEle.(*IndexElement).value
	idxEle, _ = idx.KeyValue.FindData(uniqueKey.([]byte))
	obj := idxEle.(*IndexElement).value.(*types.BlockState)
	log.Info("begin obj:%#v", obj)
	assert.Equal(t, bs, obj)
}

func TestForkContainer_LowerBound(t *testing.T) {
	mi, bs := GenerateNext()
	idx := mi.GetIndex("byBlockNum")
	bn := make([]byte, 8)
	binary.BigEndian.PutUint64(bn, uint64(bs.BlockNum))
	s := []byte("byBlockNum_")
	s = append(s, bn...)
	itr := idx.LowerBound(s)
	idxEle := itr.KeySet.Data[0]
	assert.Equal(t, append([]byte("byBlockId_"), bs.BlockId.Bytes()...), idxEle.(*IndexElement).value.([]byte))
	log.Info("second key:%#v", idxEle.(*IndexElement).value.([]byte))

}

func TestForkContainer_UpperBound(t *testing.T) {
	mi, bs := GenerateNext()

	idx := mi.GetIndex("byBlockNum")
	bn := make([]byte, 8)
	binary.BigEndian.PutUint64(bn, uint64(bs.BlockNum))
	s := []byte("byBlockNum_")
	s = append(s, bn...)
	itr := idx.UpperBound(s)
	idxEle := itr.KeySet.Data[0]
	fmt.Println(idxEle.(*IndexElement).value)
	assert.Equal(t, append([]byte("byBlockId_"), bs.BlockId.Bytes()...), idxEle.(*IndexElement).value.([]byte))
}
