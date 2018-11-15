package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
)

type ForkDatabase struct {
	Head           *types.BlockState `json:"head"`
	multiIndexFork *multiIndexFork
}

//var IrreversibleBlock chan BlockState = make(chan BlockState)

/*type ForkMultiIndexType struct {
	ByBlockID  common.BlockIdType `storm:"unique" json:"id"`
	ByPrev     common.BlockIdType `storm:"index"  json:"prev"`
	ByBlockNum common.Tuple       `storm:"index"  json:"block_num"`
	BlockState *BlockState         `storm:"inline"`
}*/

func GetForkDbInstance(stateDir string) *ForkDatabase {

	forkDB, err := newForkDatabase(stateDir, common.DefaultConfig.ForkDBName, true)
	if err != nil {
		log.Error("GetForkDbInstance is error ,detail:%s", err.Error())
	}
	return forkDB
}

func newForkDatabase(path string, fileName string, rw bool) (*ForkDatabase, error) {
	return &ForkDatabase{multiIndexFork: newMultiIndexFork()}, nil
}

func (f *ForkDatabase) SetHead(s *types.BlockState) {

	try.EosAssert(s.BlockId == s.Header.BlockID(), &exception.ForkDatabaseException{},
		"block state id:%d, is different from block state header id:%d", s.ID, s.Header.BlockID())

	try.EosAssert(s.BlockNum == s.Header.BlockNumber(), &exception.ForkDatabaseException{}, "unable to insert block state, duplicate state detected")
	if f.Head == nil {
		f.Head = s
	} else if f.Head.BlockNum < s.BlockNum {
		f.Head = s
	}

	ok := f.multiIndexFork.Insert(s)
	if !ok {
		log.Error("forkDatabase SetHead insert is error:%#v", s)
	}

}

func (f *ForkDatabase) AddBlockState(b *types.BlockState) *types.BlockState {

	ok := f.multiIndexFork.Insert(b)
	if !ok {
		log.Error("ForkDatabase AddBlockState insert is error:%#v", b)
	}
	try.EosAssert(ok, &exception.ForkDatabaseException{}, "ForkDB AddBlockState Is Error:", ok)

	result, err := f.multiIndexFork.GetIndex("byLiBlockNum").Begin()
	if err != nil {
		log.Error("ForkDatabase AddBlockState multiIndexFork Begin is error:%#v", err.Error())
	}
	f.Head = result
	lib := f.Head.DposIrreversibleBlocknum
	oldest, err := f.multiIndexFork.GetIndex("byLibBlockNum").Begin()
	if oldest.BlockNum < lib {
		f.Prune(oldest)
	}
	return b
}
func (f *ForkDatabase) AddSignedBlockState(signedBlock *types.SignedBlock, trust bool) *types.BlockState {
	try.EosAssert(signedBlock != nil, &exception.ForkDatabaseException{}, "attempt to add null block")
	try.EosAssert(f.Head != nil, &exception.ForkDbBlockNotFound{}, "no head block set")

	blockId := signedBlock.BlockID()
	//blockState := types.BlockState{}

	block := f.multiIndexFork.find(blockId) //find all data multiIndex
	try.EosAssert(block == nil, &exception.ForkDatabaseException{}, "we already know about this block")
	prior := f.multiIndexFork.find(signedBlock.Previous)
	try.EosAssert(prior == nil, &exception.ForkDatabaseException{}, "unlinkable block:%#v,%#v", blockId, signedBlock.Previous)

	//previous := types.BlockState{}
	b := types.BlockState{}
	b.BlockId = signedBlock.Previous

	result := types.NewBlockState3(&prior.BlockHeaderState, signedBlock, trust)

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
	f.multiIndexFork.erase(&b)
	result, err := f.multiIndexFork.GetIndex("byLibBlockNum").Begin()
	if err != nil {
		log.Error("ForkDataBase Remove index Begin is error:%#v", err.Error())
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
	result := f.multiIndexFork.find(h.BlockId)
	try.EosAssert(result == nil, &exception.ForkDbBlockNotFound{}, "could not find block in fork database")
	result.InCurrentChain = inCurrentChain
	f.multiIndexFork.modify(result)
}

func (f *ForkDatabase) Prune(h *types.BlockState) {

	num := h.BlockNum
	idx := f.multiIndexFork.GetIndex("byBlockNum")
	bni, err := idx.Begin()
	if err != nil {
		log.Error("ForkDatabase Prune multiIndexFork Begin is error:%#v", err.Error())
	}
	for bni != nil && bni.BlockNum < num {
		f.Prune(bni)
		bni, _ = idx.Begin()
	}
	p := f.multiIndexFork.find(h.BlockId)
	if p != nil {
		//irreversible(*itr) TODO channel
		f.multiIndexFork.erase(p)
	}
	numidx := f.multiIndexFork.GetIndex("byBlockNum")
	val, sub := numidx.value.LowerBound(h)
	if sub >= 0 {
		obj := val.(*types.BlockState)
		for obj.BlockNum == num {

			f.Remove(&obj.BlockId)
		}
	}

}

func (f *ForkDatabase) GetBlock(id *common.BlockIdType) *types.BlockState {

	b := f.multiIndexFork.find(*id)
	if b != nil {
		return b
	}
	return &types.BlockState{}
}

func (f *ForkDatabase) GetBlockInCurrentChainByNum(n uint32) *types.BlockState {
	b := types.BlockState{}
	b.BlockNum = n
	numIdx := f.multiIndexFork.GetIndex("byBlockNum")
	val, _ := numIdx.value.LowerBound(&b)
	obj := val.(*types.BlockState)
	if val != nil || obj.BlockNum != n || obj.InCurrentChain != true {
		return &types.BlockState{}
	}
	return obj
}

func (f *ForkDatabase) SetBftIrreversible(id common.BlockIdType) {
	param := types.BlockState{}
	param.BlockId = id
	b := f.multiIndexFork.find(id)
	blockNum := b.BlockNum
	b.BftIrreversibleBlocknum = b.BlockNum
	f.multiIndexFork.modify(b)

	update := func(in []common.BlockIdType) []common.BlockIdType {
		updated := []common.BlockIdType{}
		for _, i := range in {
			pidx := f.multiIndexFork.GetIndex("byPrev")
			//try.EosAssert(err==nil,&exception.ForkDbBlockNotFound{},"SetBftIrreversible could not find idx in fork database")
			b := types.BlockState{}
			b.BlockId = i
			pitr := pidx.lowerBound(&b)
			epitr := pidx.upperBound(&b)
			//try.EosAssert(pitr == nil, &exception.ForkDbBlockNotFound{}, "SetBftIrreversible could not find idx in fork database")
			for pitr != epitr {
				if pitr.value.BftIrreversibleBlocknum < blockNum {
					pitr.value.BftIrreversibleBlocknum = blockNum
					updated = append(updated, pitr.value.BlockId)
				}
				f.multiIndexFork.modify(pitr.value)
				pitr.next()
			}
		}
		return updated
	}
	queue := []common.BlockIdType{id}
	update(queue)
}

func (f *ForkDatabase) Close() {
	//isFdActive = false
	//f.DB.Close()
}
