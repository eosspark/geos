package types

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/log"
)

var isFdActive bool = false

type ForkDatabase struct {
	DB database.DataBase
	//Index   *ForkMultiIndexType `json:"index"`
	Head *BlockState `json:"head"`
	//DataDir string
}

//var IrreversibleBlock chan BlockState = make(chan BlockState)

/*type ForkMultiIndexType struct {
	ByBlockID  common.BlockIdType `storm:"unique" json:"id"`
	ByPrev     common.BlockIdType `storm:"index"  json:"prev"`
	ByBlockNum common.Tuple       `storm:"index"  json:"block_num"`
	BlockState *BlockState         `storm:"inline"`
}*/

func GetForkDbInstance(stateDir string) *ForkDatabase {
	forkDB := &ForkDatabase{}
	if !isFdActive {
		forkd, err := newForkDatabase(stateDir, common.DefaultConfig.ForkDBName, true)
		if err != nil {
			log.Error("GetForkDbInstance is error ,detail:", err)
		}
		//forkd.DB.Insert("test")
		forkDB = forkd
		//
		isFdActive = true
	}
	return forkDB
}

func newForkDatabase(path string, fileName string, rw bool) (*ForkDatabase, error) {

	db, err := database.NewDataBase(path + "/" + fileName)
	if err != nil {
		log.Error("newForkDatabase is error:", err)
		return nil, err
	}

	return &ForkDatabase{DB: db}, err
}

func (f *ForkDatabase) SetHead(s *BlockState) {

	exception.EosAssert(s.ID == s.Header.BlockID(), &exception.ForkDatabaseException{},
		"block state id:d%, is different from block state header id:d%", s.ID, s.Header.BlockID())

	exception.EosAssert(s.BlockNum == s.Header.BlockNumber(), &exception.ForkDatabaseException{}, "unable to insert block state, duplicate state detected")
	if f.Head == nil {
		f.Head = s
	} else if f.Head.BlockNum < s.BlockNum {
		f.Head = s
	}

	err := f.DB.Insert(s)
	if err != nil {
		fmt.Println("ForkDB SetHead is error:", err)
	}

}

func (f *ForkDatabase) AddBlockState(blockState *BlockState) *BlockState {

	result := BlockState{}
	err := f.DB.Insert(blockState)
	//TODO try catch
	//exception.EosAssert(err == nil, &exception.ForkDatabaseException{}, "ForkDB AddBlockState Is Error:", err)

	multiIndex, err := f.DB.GetIndex("byLibBlockNum", &result)
	//exception.EosAssert(err == nil, &exception.ForkDatabaseException{}, "ForkDB AddBlockState Is Error:", err)

	err = multiIndex.Begin(&result)
	//exception.EosAssert(err == nil, &exception.ForkDatabaseException{}, "ForkDB AddBlockState Is Error:", err)

	lib := f.Head.DposIrreversibleBlocknum
	oldest, err := f.DB.GetIndex("byBlockNum", &result)
	//exception.EosAssert(err == nil, &exception.ForkDatabaseException{}, "ForkDB AddBlockState Is Error:", err)
	err = oldest.Begin(&result)
	exception.EosAssert(err == nil, &exception.ForkDatabaseException{}, "ForkDB AddBlockState Is Error:", err)

	if result.BlockNum < lib {
		f.Prune(&result)
	}
	return blockState
}
func (f *ForkDatabase) AddSignedBlockState(signedBlcok *SignedBlock, trust bool) *BlockState {
	exception.EosAssert(signedBlcok != nil, &exception.ForkDatabaseException{}, "attempt to add null block")
	exception.EosAssert(f.Head != nil, &exception.ForkDbBlockNotFound{}, "no head block set")

	blockId := signedBlcok.BlockID()
	blockState := BlockState{}

	index, err := f.DB.GetIndex("byBlockId", &blockState) //find all data multiIndex
	if err != nil {
		log.Error("AddSignedBlockState is error,detail:", err)
	}
	blockState.ID = blockId
	bs := BlockState{}
	index.Find(blockState, &bs)
	exception.EosAssert(&bs == nil, &exception.ForkDatabaseException{}, "we already know about this block")

	previous := BlockState{}
	b := BlockState{}
	b.ID = signedBlcok.Previous
	erro := index.Find(b, &previous)
	if erro != nil {
		log.Error("AddSignedBlockState is error,detail:", err)
	}
	exception.EosAssert(&previous != nil, &exception.UnlinkableBlockException{}, "unlinkable block", signedBlcok.BlockID(), signedBlcok.Previous)

	result := NewBlockState3(&previous.BlockHeaderState, signedBlcok, trust)

	return f.AddBlockState(result)
}

func (f *ForkDatabase) Add(c *HeaderConfirmation) {
	header := f.GetBlock(&c.BlockId)

	exception.EosAssert(header != nil, &exception.ForkDbBlockNotFound{}, "unable to find block id ", c.BlockId)
	header.AddConfirmation(c)

	if header.BftIrreversibleBlocknum < header.BlockNum && len(header.Confirmations) >= ((len(header.ActiveSchedule.Producers)*2)/3+1) {
		f.SetBftIrreversible(c.BlockId)
	}
}
func (f *ForkDatabase) Header() *BlockState { return f.Head }

/*type BranchType struct {
	branch []BlockState
}*/

type FetchBranch struct {
	first  []BlockState
	second []BlockState
}

func (f *ForkDatabase) FetchBranchFrom(first *common.BlockIdType, second *common.BlockIdType) FetchBranch {
	result := FetchBranch{}
	firstBranch := f.GetBlock(first)

	secondBranch := f.GetBlock(second)

	for firstBranch.BlockNum > secondBranch.BlockNum {
		result.first = append(result.first, *firstBranch)
		firstBranch = f.GetBlock(&firstBranch.Header.Previous)
		exception.EosAssert(firstBranch != nil, &exception.ForkDbBlockNotFound{}, "block d% does not exist", firstBranch.Header.Previous)
	}

	for firstBranch.BlockNum < secondBranch.BlockNum {
		result.second = append(result.second, *secondBranch)
		secondBranch = f.GetBlock(&firstBranch.Header.Previous)
		exception.EosAssert(secondBranch != nil, &exception.ForkDbBlockNotFound{}, "block d% does not exist", secondBranch.Header.Previous)
	}

	for firstBranch.Header.Previous != secondBranch.Header.Previous {
		result.first = append(result.first, *firstBranch)
		result.second = append(result.second, *secondBranch)
		firstBranch = f.GetBlock(&firstBranch.Header.Previous)
		secondBranch = f.GetBlock(&secondBranch.Header.Previous)
		exception.EosAssert(firstBranch != nil && secondBranch != nil, &exception.ForkDbBlockNotFound{},
			"either block d% or d% does not exist",
			firstBranch.Header.Previous, secondBranch.Header.Previous)
	}

	if firstBranch != nil && secondBranch != nil {
		result.first = append(result.first, *firstBranch)
		result.second = append(result.second, *secondBranch)
	}

	return result
}

func (f *ForkDatabase) Remove(id *common.BlockIdType) {
	removeQueue := []common.BlockIdType{*id}

	for i := 0; i < len(removeQueue); i++ {
		p := BlockState{}
		err := f.DB.Find("ID", p, &p)
		exception.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDB Remove Find Block Is Error:", err)
		f.DB.Remove(&p)
		previdx, err := f.DB.GetIndex("byPrev", &BlockState{})
		param := &BlockState{}
		param.ID = removeQueue[i]
		previtr, err := previdx.LowerBound(param)
		exception.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDB Remove LowerBound Is Error:", err)
		pre := BlockState{}
		err = previtr.Data(pre)
		exception.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDB Remove previtr.Data Is Error:", err)
		for previtr != previdx.End() && pre.Header.Previous == removeQueue[i] {
			removeQueue = append(removeQueue, pre.ID)
		}
	}
	mi, err := f.DB.GetIndex("byLibBlockNum", &BlockState{})
	exception.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDB Remove f.DB.GetIndex Is Error:", err)
	err = mi.Begin(f.Head)
	exception.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDB Remove mi.Begin(f.Head) Is Error:", err)
}

func (f *ForkDatabase) GetBlockInCurrentChainByNum(n uint32) *BlockState {
	numidx, err := f.DB.GetIndex("byBlockNum", &BlockState{})
	exception.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDB Remove mi.Begin(f.Head) Is Error:", err)
	param := BlockState{}
	param.BlockNum = n
	nitr, err := numidx.LowerBound(param)
	exception.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDB Remove LowerBound Is Error:", err)
	result := BlockState{}
	nitr.Data(result)
	return &result
}

func (f *ForkDatabase) SetValidity(h *BlockState, valid bool) {
	if !valid {
		f.Remove(&h.ID)
	} else {
		h.Validated = true
	}
}
func (f *ForkDatabase) MarkInCurrentChain(h *BlockState, inCurrentChain bool) {
	if h.InCurrentChain == inCurrentChain {
		return
	}
	byIdIdx, err := f.DB.GetIndex("byBlockId", &BlockState{})
	exception.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "could not find block in fork database")
	param := BlockState{}
	param.ID = h.ID
	result := BlockState{}
	err = byIdIdx.Find(param, &result)
	exception.EosAssert(&result == nil, &exception.ForkDbBlockNotFound{}, "could not find block in fork database")
	f.DB.Modify(result, func(b *BlockState) {
		b.InCurrentChain = inCurrentChain
	})
}

func (f *ForkDatabase) Prune(h *BlockState) {

	num := h.BlockNum
	param := BlockState{}
	mIndex, err := f.DB.GetIndex("byBlockNum", &param)
	err = mIndex.Begin(&param)
	bItr := mIndex.IteratorTo(&param)
	for !mIndex.CompareEnd(bItr) && param.BlockNum < num {
		f.Prune(&param)
		err = mIndex.Begin(&param)
		bItr = mIndex.IteratorTo(&param)
		f.DB.Remove(bItr)
	}
	p := BlockState{}
	p.ID = h.ID
	result := BlockState{}
	err = f.DB.Find("ID", p, &result)
	exception.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "Prune could not find block in fork database")
	if !common.Empty(result) {
		//irreversible(*itr) TODO channel
		//my->index.erase(itr)
		//f.DB.Remove(result)
	}

	in := BlockState{}
	in.BlockNum = num
	numIdx, err := f.DB.GetIndex("byBlockNun", &in)
	nitr, err := numIdx.LowerBound(&in)

	err = nitr.Data(&in)
	for !numIdx.CompareEnd(nitr) && in.BlockNum == num {

		itrToRemove := nitr
		nitr.Next()
		err = itrToRemove.Data(&in)
		id := in.ID
		f.Remove(&id)
	}
}

func (f *ForkDatabase) GetBlock(id *common.BlockIdType) *BlockState {

	blockState := BlockState{}
	blockState.ID = *id
	multiIndex, err := f.DB.GetIndex("ID", &blockState)
	if err != nil {
		fmt.Println("ForkDb GetBlock Is Error:", err)
		return &BlockState{}
	}
	err = multiIndex.Begin(&blockState)
	if err != nil {
		fmt.Println("ForkDB GetBlock MultiIndex.Begin Is Error :", err)
	}
	return &blockState
}

func (f *ForkDatabase) SetBftIrreversible(id common.BlockIdType) {
	param := BlockState{}
	param.ID = id
	result := BlockState{}
	idx, err := f.DB.GetIndex("byID", param)
	exception.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "could not find block in fork database")
	err = idx.Find(param, &result)
	exception.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "could not find idx in fork database")
	blockNum := result.BlockNum
	f.DB.Modify(result, func(b *BlockState) {
		b.BftIrreversibleBlocknum = b.BlockNum
	})
	update := func(in []common.BlockIdType) []common.BlockIdType {
		updated := []common.BlockIdType{}
		for _, i := range in {
			pidx, err := f.DB.GetIndex("byPrev", BlockState{})
			//exception.EosAssert(err==nil,&exception.ForkDbBlockNotFound{},"SetBftIrreversible could not find idx in fork database")
			in := BlockState{}
			in.ID = i
			pitr, err := pidx.LowerBound(in)
			epitr, err := pidx.UpperBound(in)
			exception.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "SetBftIrreversible could not find idx in fork database")
			for pitr != epitr {
				f.DB.Modify(pitr, func(bsp *BlockState) {
					if bsp.BftIrreversibleBlocknum < blockNum {
						bsp.BftIrreversibleBlocknum = blockNum
						updated = append(updated, bsp.ID)
					}

				})
				pitr.Next()
			}
		}

		fmt.Println(updated)
		return updated
	}

	fmt.Println(blockNum)

	queue := []common.BlockIdType{id}
	update(queue)
}
