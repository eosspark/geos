package fork_multi_index

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"testing"
)

func generateBlockState() *types.BlockState {
	initPriKey, _ := ecc.NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	initPubKey := initPriKey.PublicKey()
	eosio := common.AccountName(common.N("eosio"))
	eos := common.AccountName(common.N("eos"))
	tester := common.AccountName(common.N("tester"))

	initSchedule := types.ProducerScheduleType{Producers: []types.ProducerKey{
		{ProducerName: eosio, BlockSigningKey: initPubKey},
		{ProducerName: eos, BlockSigningKey: initPubKey},
		{ProducerName: tester, BlockSigningKey: initPubKey},
	}}

	genHeader := new(types.BlockHeaderState)
	genHeader.Header = *types.NewSignedBlockHeader()
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

	return blockState
}

func generateBlock(prev *types.BlockState) *types.BlockState {
	bs := types.NewBlockState2(&prev.BlockHeaderState, prev.SignedBlock.Timestamp.Next())
	bs.BlockId = bs.Header.BlockID()
	bs.DposIrreversibleBlocknum = prev.DposIrreversibleBlocknum + 1

	return bs
}

func printBlock(blk *types.BlockState) {
	fmt.Printf("id:%v, previous:%v, num:%d, inCurrentChain:%t, Dpos:%d, Bft:%d \n",
		blk.BlockId, blk.Header.Previous, blk.BlockNum, blk.InCurrentChain,
		blk.DposIrreversibleBlocknum, blk.BftIrreversibleBlocknum)
}

func TestNew(t *testing.T) {
	midx := New()
	bs1 := generateBlockState()
	bs2 := generateBlock(bs1)
	bs3 := generateBlock(bs2)
	bs4 := generateBlock(bs3)

	midx.Insert(bs1)
	midx.Insert(bs2)
	midx.Insert(bs3)
	midx.Insert(bs4)

	byIdIdx := midx.ByBlockId
	for id := range byIdIdx {
		printBlock(midx.Value(byIdIdx[id]))
	}

	fmt.Println()

	byPrevIdx := midx.ByPrev
	for itr := byPrevIdx.Begin(); itr.HasNext(); itr.Next() {
		printBlock(midx.Value(itr.Value()))
	}

}
