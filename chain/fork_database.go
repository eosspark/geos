package chain

import (
	"fmt"
	"github.com/eosspark/container/sets/treeset"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"io/ioutil"
	"os"
	"reflect"
)

type ForkDatabase struct {
	Head           *types.BlockState `json:"head"`
	MultiIndexFork *MultiIndexFork
	ForkDbPath     string
	fileStream     *os.File
}

func GetForkDbInstance(stateDir string) *ForkDatabase {
	forkDB := &ForkDatabase{MultiIndexFork: newMultiIndexFork()}
	//forkDB, err := newForkDatabase(stateDir, common.DefaultConfig.ForkDbName, true)
	//try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "%s", err)
	return forkDB
}

func newForkDatabase(stateDir string, fileName string, rw bool) (*ForkDatabase, error) {
	fk := &ForkDatabase{MultiIndexFork: newMultiIndexFork()}
	_, err := os.Stat(stateDir)
	if err != nil {
		os.MkdirAll(stateDir, os.ModePerm)
	}
	fk.ForkDbPath = stateDir + "/" + fileName
	fStream, err := os.OpenFile(fk.ForkDbPath, os.O_RDWR, os.ModePerm)
	if fStream == nil {
		fStream, err = os.Create(fk.ForkDbPath)
		try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "%s", err)
	}

	b, err := ioutil.ReadAll(fStream)
	try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "%s", err)

	if len(b) > 0 {
		mif := &MultiIndexFork{}
		mif.Indexs = make(map[string]*IndexFork)
		index := &IndexFork{Target: "byBlockId", Uniqueness: true}
		index2 := &IndexFork{Target: "byPrev", Uniqueness: false}
		index3 := &IndexFork{Target: "byBlockNum", Uniqueness: false}
		index4 := &IndexFork{Target: "byLibBlockNum", Uniqueness: false}
		mif.Indexs["byBlockId"] = index
		mif.Indexs["byPrev"] = index2
		mif.Indexs["byBlockNum"] = index3
		mif.Indexs["byLibBlockNum"] = index4
		bt := treeset.NewMultiWith(types.BlockIdTypes, types.CompareBlockId)
		index.Value = bt
		mif.Indexs[index.Target] = index
		bt2 := treeset.NewMultiWith(types.BlockIdTypes, types.ComparePrev)
		index2.Value = bt2
		mif.Indexs[index2.Target] = index2
		bt3 := treeset.NewMultiWith(types.BlockNumType, types.CompareBlockNum)
		index3.Value = bt3
		mif.Indexs[index3.Target] = index3
		bt4 := treeset.NewMultiWith(types.BlockNumType, types.CompareLibNum)
		index4.Value = bt4
		mif.Indexs[index4.Target] = index4
		fmt.Println("forkDatabase:", len(mif.Indexs), mif.Indexs)
		fmt.Println("value:  ", reflect.TypeOf(mif.Indexs["byBlockId"].Value).String())

		fmt.Println("78:  ", index.Value.ValueType.String())

		err := rlp.DecodeBytes(b, mif)
		fmt.Println("forkDatabase error:", err)
		fmt.Println("forkDatabase2:", len(mif.Indexs), mif.Indexs)
		fk.MultiIndexFork = mif

	}

	fk.fileStream = fStream
	err = os.Remove(fk.ForkDbPath)
	return fk, err
}

func (f *ForkDatabase) SetHead(s *types.BlockState) {

	try.EosAssert(s.BlockId == s.Header.BlockID(), &exception.ForkDatabaseException{},
		"block state id:%d, is different from block state header id:%d", s.BlockId, s.Header.BlockID())

	try.EosAssert(s.BlockNum == s.Header.BlockNumber(), &exception.ForkDatabaseException{}, "unable to insert block state, duplicate state detected")
	if f.Head == nil {
		f.Head = s
	} else if f.Head.BlockNum < s.BlockNum {
		f.Head = s
	}

	ok := f.MultiIndexFork.Insert(s)
	if !ok {
		log.Error("forkDatabase SetHead insert is error:%#v", s)
	}

}

func (f *ForkDatabase) AddBlockState(b *types.BlockState) *types.BlockState {

	ok := f.MultiIndexFork.Insert(b)
	if !ok {
		log.Error("ForkDatabase AddBlockState insert is error:%#v", b)
	}
	try.EosAssert(ok, &exception.ForkDatabaseException{}, "ForkDB AddBlockState Is Error:", ok)
	idx := f.MultiIndexFork.GetIndex("byLibBlockNum")
	if idx != nil {
		result, err := idx.Begin()
		if err != nil {
			log.Error("ForkDatabase AddBlockState MultiIndexFork Begin is error:%#v", err.Error())
		}
		f.Head = result
		lib := f.Head.DposIrreversibleBlocknum
		mItr := f.MultiIndexFork.GetIndex("byLibBlockNum").Value.Iterator()
		if mItr.Next() {
			itrVal := mItr.Value()
			oldest := itrVal.(*types.BlockState)
			if oldest.BlockNum < lib {
				f.Prune(oldest)
			}
		}
	} else {
		try.EosAssert(idx != nil, &exception.ForkDatabaseException{}, "ForkDatabase AddBlockState MultiIndexFork Begin is not found!")
	}

	return b
}
func (f *ForkDatabase) AddSignedBlockState(signedBlock *types.SignedBlock, trust bool) *types.BlockState {
	try.EosAssert(signedBlock != nil, &exception.ForkDatabaseException{}, "attempt to add null block")
	try.EosAssert(f.Head != nil, &exception.ForkDbBlockNotFound{}, "no head block set")

	blockId := signedBlock.BlockID()
	//blockState := types.BlockState{}

	block := f.MultiIndexFork.find(blockId) //find all data multiIndex
	try.EosAssert(block == nil, &exception.ForkDatabaseException{}, "we already know about this block")
	prior := f.MultiIndexFork.find(signedBlock.Previous)
	try.EosAssert(prior == nil, &exception.ForkDatabaseException{}, "unlinkable block:%#v,%#v", blockId, signedBlock.Previous)

	//previous := types.BlockState{}
	b := types.BlockState{}
	b.BlockId = signedBlock.Previous

	result := types.NewBlockState3(&prior.BlockHeaderState, signedBlock, trust)
	try.EosAssert(result != nil, &exception.ForkDatabaseException{}, "fail to add new block state")
	return f.AddBlockState(result)
}

func (f *ForkDatabase) Add(c *types.HeaderConfirmation) {
	header := f.GetBlock(&c.BlockId)

	try.EosAssert(header != nil, &exception.ForkDbBlockNotFound{}, "unable to find block id:%#v ", c.BlockId)
	header.AddConfirmation(c)

	if header.BftIrreversibleBlocknum < header.BlockNum && len(header.Confirmations) >= ((len(header.ActiveSchedule.Producers)*2)/3+1) {
		f.SetBftIrreversible(c.BlockId)
	}
}
func (f *ForkDatabase) Header() *types.BlockState { return f.Head }

type FetchBranch struct {
	first  []types.BlockState
	second []types.BlockState
}

func (f *ForkDatabase) FetchBranchFrom(first *common.BlockIdType, second *common.BlockIdType) FetchBranch {
	result := FetchBranch{}
	firstBranch := f.GetBlock(first)

	secondBranch := f.GetBlock(second)

	for firstBranch.BlockNum > secondBranch.BlockNum {
		result.first = append(result.first, *firstBranch)
		firstBranch = f.GetBlock(&firstBranch.Header.Previous)
		try.EosAssert(firstBranch != nil, &exception.ForkDbBlockNotFound{}, "block %d does not exist", firstBranch.Header.Previous)
	}

	for firstBranch.BlockNum < secondBranch.BlockNum {
		result.second = append(result.second, *secondBranch)
		secondBranch = f.GetBlock(&firstBranch.Header.Previous)
		try.EosAssert(secondBranch != nil, &exception.ForkDbBlockNotFound{}, "block %d does not exist", secondBranch.Header.Previous)
	}

	for firstBranch.Header.Previous != secondBranch.Header.Previous {
		result.first = append(result.first, *firstBranch)
		result.second = append(result.second, *secondBranch)
		firstBranch = f.GetBlock(&firstBranch.Header.Previous)
		secondBranch = f.GetBlock(&secondBranch.Header.Previous)
		try.EosAssert(firstBranch != nil && secondBranch != nil, &exception.ForkDbBlockNotFound{},
			"either block %d or %d does not exist",
			firstBranch.Header.Previous, secondBranch.Header.Previous)
	}

	if firstBranch != nil && secondBranch != nil {
		result.first = append(result.first, *firstBranch)
		result.second = append(result.second, *secondBranch)
	}

	return result
}

func (f *ForkDatabase) Remove(id *common.BlockIdType) {
	b := types.BlockState{}
	b.BlockId = *id
	f.MultiIndexFork.erase(&b)
	result, err := f.MultiIndexFork.GetIndex("byLibBlockNum").Begin()
	if err != nil {
		log.Error("ForkDataBase Remove index Begin is error:%s", err)
	}
	f.Head = result

}

func (f *ForkDatabase) SetValidity(h *types.BlockState, valid bool) {
	if !valid {
		f.Remove(&h.BlockId)
	} else {
		h.Validated = true
	}
}
func (f *ForkDatabase) MarkInCurrentChain(h *types.BlockState, inCurrentChain bool) {
	if h.InCurrentChain == inCurrentChain {
		return
	}
	result := f.MultiIndexFork.find(h.BlockId)
	try.EosAssert(result == nil, &exception.ForkDbBlockNotFound{}, "could not find block in fork database")
	result.InCurrentChain = inCurrentChain
	f.MultiIndexFork.modify(result)
}

func (f *ForkDatabase) Prune(h *types.BlockState) {
	num := h.BlockNum
	idx := f.MultiIndexFork.GetIndex("byBlockNum")
	//bni, err := idx.Begin()
	byBn := idx.Value.Iterator()

	for !byBn.Last() && byBn.Next() && byBn.Value().(*types.BlockState).BlockNum < num {
		f.Prune(byBn.Value().(*types.BlockState))
		byBn.Begin()
	}
	obj := f.MultiIndexFork.find(h.BlockId)
	if obj != nil {
		//irreversible(*itr) TODO channel
		f.MultiIndexFork.erase(obj)
	}
	numidx := f.MultiIndexFork.GetIndex("byBlockNum")
	mItr := numidx.Value.LowerBound(h)
	for !mItr.Last() && mItr.Value().(*types.BlockState).BlockNum == num {
		obj := mItr.Value().(*types.BlockState)
		mItr.Next()
		f.Remove(&obj.BlockId)
	}
}

func (f *ForkDatabase) GetBlock(id *common.BlockIdType) *types.BlockState {

	b := f.MultiIndexFork.find(*id)
	if b != nil {
		return b
	}
	return nil
}

func (f *ForkDatabase) GetBlockInCurrentChainByNum(n uint32) *types.BlockState {
	b := types.BlockState{}
	b.BlockNum = n
	numIdx := f.MultiIndexFork.GetIndex("byBlockNum")
	mItr := numIdx.Value.LowerBound(&b)

	if mItr == nil || mItr.Value().(*types.BlockState).BlockNum != n || mItr.Value().(*types.BlockState).InCurrentChain != true {
		return nil
	}
	return mItr.Value().(*types.BlockState)
}

func (f *ForkDatabase) SetBftIrreversible(id common.BlockIdType) {
	param := types.BlockState{}
	param.BlockId = id
	b := f.MultiIndexFork.find(id)
	blockNum := b.BlockNum
	b.BftIrreversibleBlocknum = b.BlockNum
	f.MultiIndexFork.modify(b)

	update := func(in []common.BlockIdType) []common.BlockIdType {
		updated := []common.BlockIdType{}
		for _, i := range in {
			pidx := f.MultiIndexFork.GetIndex("byPrev")
			//try.EosAssert(err==nil,&exception.ForkDbBlockNotFound{},"SetBftIrreversible could not find idx in fork database")
			b := types.BlockState{}
			b.BlockId = i
			pitr := pidx.lowerBound(&b)
			epitr := pidx.upperBound(&b)
			try.EosAssert(pitr == nil, &exception.ForkDbBlockNotFound{}, "SetBftIrreversible could not find idx in fork database")
			//for !pitr.Equal(*epitr) {
			for *pitr != *epitr {
				if pitr.Value().(*types.BlockState).BftIrreversibleBlocknum < blockNum {
					pitr.Value().(*types.BlockState).BftIrreversibleBlocknum = blockNum
					updated = append(updated, pitr.Value().(*types.BlockState).BlockId)
				}
				f.MultiIndexFork.modify(pitr.Value().(*types.BlockState))
				pitr.Next()
			}
		}
		return updated
	}
	queue := []common.BlockIdType{id}
	update(queue)
}

func (f *ForkDatabase) Close() {

	log.Info("ForkDatabase Close tmp code")
	/*bts, err := rlp.EncodeToBytes(f.MultiIndexFork)
	fmt.Println("EncodeToBytes :  ",err)
	fout, err := os.Create(f.ForkDbPath)
	f.fileStream = fout
	_, err = f.fileStream.Write(bts)
	err = f.fileStream.Close()
	try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "%s", err)
	if err != nil {
		log.Error("ForkDatabase Close is error:%s", err)
	}*/
}
