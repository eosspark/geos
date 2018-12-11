package chain

import (
	"fmt"
	"github.com/eosspark/container/maps/treemap"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"testing"
	"time"
	"github.com/eosspark/container/utils"
)

func initForkDatabase() (*MultiIndexFork, *types.BlockState) {
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
	genHeader.ProducerToLastProduced = *treemap.NewWith(common.TypeName, utils.TypeUInt32, common.CompareName)
	genHeader.ProducerToLastImpliedIrb = *treemap.NewWith(common.TypeName, utils.TypeUInt32, common.CompareName)
	genHeader.BlockSigningKey = initPubKey
	genHeader.Header.ProducerSignature = *ecc.NewSigNil()
	blockState := types.NewBlockState(genHeader)
	blockState.SignedBlock = new(types.SignedBlock)
	blockState.SignedBlock.SignedBlockHeader = genHeader.Header
	blockState.Header.ProducerSignature = *ecc.NewSigNil()
	blockState.InCurrentChain = true

	mi := newMultiIndexFork()

	mi.Insert(blockState)
	fork := GetForkDbInstance("/tmp/data")
	fork.AddBlockState(blockState)
	//fmt.Println("%#v",b.BlockNum)
	//fork.Close()
	return mi, blockState
}

func TestForkDatabase_Close(t *testing.T) {
	start := time.Now()
	fmt.Println(start)
	mi, bs := initForkDatabase()
	for i := 0; i < 10; i++ {
		t := 1162425602 + 200
		tmp := types.BlockTimeStamp(t)
		bhs := bs.GenerateNext(tmp)
		bhs.BlockId = bhs.Header.BlockID()
		blockState := types.NewBlockState(bhs)
		mi.Insert(blockState)
	}
	fork := GetForkDbInstance("/tmp/data")
	fork.MultiIndexFork = mi
	fork.Close()
	end := time.Now()
	fmt.Println(end.Sub(start))
	fork2 := GetForkDbInstance("/tmp/data")
	fmt.Printf("%#v\n\n\n", fork.MultiIndexFork.Indexs["byLibBlockNum"].Value.Data[0])
	fmt.Printf("%#v", fork2.MultiIndexFork.Indexs["byLibBlockNum"].Value.Data[0])
	//assert.Equal(t, fork.MultiIndexFork, fork2.MultiIndexFork)
}

func TestGetForkDbInstance(t *testing.T) {
	//start:=time.Now()
	//fork:=GetForkDbInstance("/tmp/data")
	//if len(fork.MultiIndexFork.Indexs)>0{
	//	for _,v:=range fork.MultiIndexFork.Indexs{
	//		if v.Value.Len()>0{
	//			for _,d:=range v.Value.Data{
	//				fmt.Println(v.Value.Len())
	//				fmt.Print("%#v", d)
	//			}
	//		}
	//	}
	//}
	//end :=time.Now()
	//fmt.Println(end.Sub(start))
}
