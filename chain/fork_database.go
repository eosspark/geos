package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
)

type ForkDatabase struct {
	DB database.DataBase
	//Index   *ForkMultiIndexType `json:"index"`
	Head *types.BlockState `json:"head"`
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
	forkDB, err := newForkDatabase(stateDir, common.DefaultConfig.ForkDBName, true)
	if err != nil {
		log.Error("GetForkDbInstance is error ,detail:%s", err.Error())
	}
	return forkDB
}

func newForkDatabase(path string, fileName string, rw bool) (*ForkDatabase, error) {

	db, err := database.NewDataBase(path + "/" + fileName)
	if err != nil {
		log.Error("newForkDatabase is error:%s", err.Error())
		return nil, err
	}

	return &ForkDatabase{DB: db}, err
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

	err := f.DB.Insert(s)
	if err != nil {
		log.Error("ForkDB SetHead is error:%s", err.Error())
	}

}

func (f *ForkDatabase) AddBlockState(blockState *types.BlockState) *types.BlockState {

	result := types.BlockState{}
	err := f.DB.Insert(blockState)
	//TODO try catch
	//try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "ForkDB AddBlockState Is Error:", err)

	multiIndex, err := f.DB.GetIndex("byLibBlockNum", &result)
	//try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "ForkDB AddBlockState Is Error:", err)

	err = multiIndex.BeginData(&result)
	//try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "ForkDB AddBlockState Is Error:", err)

	lib := f.Head.DposIrreversibleBlocknum
	oldest, err := f.DB.GetIndex("byBlockNum", &result)
	//try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "ForkDB AddBlockState Is Error:", err)
	err = oldest.BeginData(&result)
	try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "ForkDB AddBlockState Is Error:%s", err)

	if result.BlockNum < lib {
		f.Prune(&result)
	}
	return blockState
}
func (f *ForkDatabase) AddSignedBlockState(signedBlcok *types.SignedBlock, trust bool) *types.BlockState {
	try.EosAssert(signedBlcok != nil, &exception.ForkDatabaseException{}, "attempt to add null block")
	try.EosAssert(f.Head != nil, &exception.ForkDbBlockNotFound{}, "no head block set")

	blockId := signedBlcok.BlockID()
	blockState := types.BlockState{}

	index, err := f.DB.GetIndex("byBlockId", &blockState) //find all data multiIndex
	if err != nil {
		log.Error("AddSignedBlockState is error,detail:%s", err.Error())
	}
	blockState.BlockId = blockId
	bs := types.BlockState{}
	index.Find(blockState, &bs)
	try.EosAssert(&bs == nil, &exception.ForkDatabaseException{}, "we already know about this block")

	previous := types.BlockState{}
	b := types.BlockState{}
	b.BlockId = signedBlcok.Previous
	erro := index.Find(b, &previous)
	if erro != nil {
		log.Error("AddSignedBlockState is error,detail:%v", err.Error())
	}
	try.EosAssert(&previous != nil, &exception.UnlinkableBlockException{}, "unlinkable block:%d,%d", signedBlcok.BlockID(), signedBlcok.Previous)

	result := types.NewBlockState3(&previous.BlockHeaderState, signedBlcok, trust)

	return f.AddBlockState(result)
}

func (f *ForkDatabase) Add(c *types.HeaderConfirmation) {
	header := f.GetBlock(&c.BlockId)

	try.EosAssert(header != nil, &exception.ForkDbBlockNotFound{}, "unable to find block id:%d ", c.BlockId)
	header.AddConfirmation(c)

	if header.BftIrreversibleBlocknum < header.BlockNum && len(header.Confirmations) >= ((len(header.ActiveSchedule.Producers)*2)/3+1) {
		f.SetBftIrreversible(c.BlockId)
	}
}
func (f *ForkDatabase) Header() *types.BlockState { return f.Head }

/*type BranchType struct {
	branch []BlockState
}*/

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
	removeQueue := []common.BlockIdType{*id}

	for i := 0; i < len(removeQueue); i++ {
		p := types.BlockState{}
		err := f.DB.Find("ID", p, &p)
		try.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDB Remove Find Block Is Error:%s", err)
		f.DB.Remove(&p)
		previdx, err := f.DB.GetIndex("byPrev", &types.BlockState{})
		param := &types.BlockState{}
		param.BlockId = removeQueue[i]
		previtr, err := previdx.LowerBound(param)
		try.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDB Remove LowerBound Is Error:%s", err)
		pre := types.BlockState{}
		err = previtr.Data(pre)
		try.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDB Remove previtr.Data Is Error:%s", err)
		for previtr != previdx.End() && pre.Header.Previous == removeQueue[i] {
			removeQueue = append(removeQueue, pre.BlockId)
		}
	}
	mi, err := f.DB.GetIndex("byLibBlockNum", &types.BlockState{})
	try.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDbBlockNotFound byLibBlockNum Is Error:%s", err)
	err = mi.BeginData(f.Head)
	try.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDbBlockNotFound BeginData Is Error:%s", err)
}

func (f *ForkDatabase) GetBlockInCurrentChainByNum(n uint32) *types.BlockState {
	numidx, err := f.DB.GetIndex("byBlockNum", &types.BlockState{SignedBlock: &types.SignedBlock{}})
	try.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDbBlockNotFound byBlockNum Is Error:%s", err)
	param := types.BlockState{}
	param.BlockNum = n
	nitr, err := numidx.LowerBound(param)
	try.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "ForkDbBlockNotFound LowerBound Is Error:%s", err)
	result := types.BlockState{}
	nitr.Data(result)
	return &result
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
	byIdIdx, err := f.DB.GetIndex("byBlockId", &types.BlockState{})
	try.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "could not find block in fork database")
	param := types.BlockState{}
	param.ID = h.ID
	result := types.BlockState{}
	err = byIdIdx.Find(param, &result)
	try.EosAssert(&result == nil, &exception.ForkDbBlockNotFound{}, "could not find block in fork database")
	f.DB.Modify(result, func(b *types.BlockState) {
		b.InCurrentChain = inCurrentChain
	})
}

func (f *ForkDatabase) Prune(h *types.BlockState) {

	num := h.BlockNum
	param := types.BlockState{}
	mIndex, err := f.DB.GetIndex("byBlockNum", &param)
	err = mIndex.BeginData(&param)
	bItr := mIndex.IteratorTo(&param)
	for !mIndex.CompareEnd(bItr) && param.BlockNum < num {
		f.Prune(&param)
		err = mIndex.BeginData(&param)
		bItr = mIndex.IteratorTo(&param)
		f.DB.Remove(bItr)
	}
	p := types.BlockState{}
	p.ID = h.ID
	result := types.BlockState{}
	err = f.DB.Find("ID", p, &result)
	try.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "Prune could not find block in fork database")
	if !common.Empty(result) {
		//irreversible(*itr) TODO channel
		//my->index.erase(itr)
		//f.DB.Remove(result)
	}

	in := types.BlockState{}
	in.BlockNum = num
	numIdx, err := f.DB.GetIndex("byBlockNun", &in)
	nitr, err := numIdx.LowerBound(&in)

	err = nitr.Data(&in)
	for !numIdx.CompareEnd(nitr) && in.BlockNum == num {

		itrToRemove := nitr
		nitr.Next()
		err = itrToRemove.Data(&in)
		id := in.BlockId
		f.Remove(&id)
	}
}

func (f *ForkDatabase) GetBlock(id *common.BlockIdType) *types.BlockState {

	blockState := types.BlockState{}
	blockState.BlockId = *id
	multiIndex, err := f.DB.GetIndex("ID", &blockState)
	if err != nil {
		log.Error("ForkDb GetBlock Is Error:%s", err.Error())
		return &types.BlockState{}
	}
	err = multiIndex.BeginData(&blockState)
	if err != nil {
		log.Error("ForkDB GetBlock MultiIndex.Begin Is Error :%s", err.Error())
	}
	return &blockState
}

func (f *ForkDatabase) SetBftIrreversible(id common.BlockIdType) {
	param := types.BlockState{}
	param.BlockId = id
	result := types.BlockState{}
	idx, err := f.DB.GetIndex("byID", param)
	try.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "could not find block in fork database")
	err = idx.Find(param, &result)
	try.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "could not find idx in fork database")
	blockNum := result.BlockNum
	f.DB.Modify(result, func(b *types.BlockState) {
		b.BftIrreversibleBlocknum = b.BlockNum
	})
	update := func(in []common.BlockIdType) []common.BlockIdType {
		updated := []common.BlockIdType{}
		for _, i := range in {
			pidx, err := f.DB.GetIndex("byPrev", types.BlockState{})
			//try.EosAssert(err==nil,&exception.ForkDbBlockNotFound{},"SetBftIrreversible could not find idx in fork database")
			in := types.BlockState{}
			in.BlockId = i
			pitr, err := pidx.LowerBound(in)
			epitr, err := pidx.UpperBound(in)
			try.EosAssert(err == nil, &exception.ForkDbBlockNotFound{}, "SetBftIrreversible could not find idx in fork database")
			for pitr != epitr {
				f.DB.Modify(pitr, func(bsp *types.BlockState) {
					if bsp.BftIrreversibleBlocknum < blockNum {
						bsp.BftIrreversibleBlocknum = blockNum
						updated = append(updated, bsp.BlockId)
					}
				})
				pitr.Next()
			}
		}
		return updated
	}
	queue := []common.BlockIdType{id}
	update(queue)
}

func (f *ForkDatabase) Close() {
	//isFdActive = false
	f.DB.Close()
}
