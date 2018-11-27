package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"io/ioutil"
	"os"
)

type ForkDatabase struct {
	Head           *types.BlockState `json:"head"`
	MultiIndexFork *MultiIndexFork
	ForkDbPath     string
	fileStream     *os.File
}

func GetForkDbInstance(stateDir string) *ForkDatabase {

	forkDB, err := newForkDatabase(stateDir, common.DefaultConfig.ForkDBName, true)
	try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "%s", err)
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
		rlp.DecodeBytes(b, mif)
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
		oldest, err := f.MultiIndexFork.GetIndex("byLibBlockNum").Begin()
		if oldest.BlockNum < lib {
			f.Prune(oldest)
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
	bni, err := idx.Begin()
	if err != nil {
		log.Error("ForkDatabase Prune MultiIndexFork Begin is error:%#v", err.Error())
	}
	for bni != nil && bni.BlockNum < num {
		f.Prune(bni)
		bni, _ = idx.Begin()
	}
	p := f.MultiIndexFork.find(h.BlockId)
	if p != nil {
		//irreversible(*itr) TODO channel
		f.MultiIndexFork.erase(p)
	}
	numidx := f.MultiIndexFork.GetIndex("byBlockNum")
	obj, sub := numidx.Value.LowerBound(h)
	if sub >= 0 {
		//obj := val.(*types.BlockState)
		for obj.BlockNum == num {

			f.Remove(&obj.BlockId)
		}
	}

}

func (f *ForkDatabase) GetBlock(id *common.BlockIdType) *types.BlockState {

	b := f.MultiIndexFork.find(*id)
	if b != nil {
		return b
	}
	return &types.BlockState{}
}

func (f *ForkDatabase) GetBlockInCurrentChainByNum(n uint32) *types.BlockState {
	b := types.BlockState{}
	b.BlockNum = n
	numIdx := f.MultiIndexFork.GetIndex("byBlockNum")
	obj, _ := numIdx.Value.LowerBound(&b)

	if obj != nil || obj.BlockNum != n || obj.InCurrentChain != true {
		return &types.BlockState{}
	}
	return obj
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
			for pitr != epitr {
				if pitr.Value.BftIrreversibleBlocknum < blockNum {
					pitr.Value.BftIrreversibleBlocknum = blockNum
					updated = append(updated, pitr.Value.BlockId)
				}
				f.MultiIndexFork.modify(pitr.Value)
				pitr.next()
			}
		}
		return updated
	}
	queue := []common.BlockIdType{id}
	update(queue)
}

func (f *ForkDatabase) Close() {

	bts, err := rlp.EncodeToBytes(f.MultiIndexFork)
	fout, err := os.Create(f.ForkDbPath)
	f.fileStream = fout
	_, err = f.fileStream.Write(bts)
	err = f.fileStream.Close()
	try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "%s", err)
	if err != nil {
		log.Error("ForkDatabase Close is error:%s", err)

	}
}
