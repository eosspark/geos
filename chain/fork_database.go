package chain

import (
	"github.com/eosspark/eos-go/chain/multi_index_containers/fork_multi_index"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/appbase/app/include"
	"os"
)

type ForkDatabase struct {
	Index fork_multi_index.MultiIndex
	Head  *types.BlockState `json:"head"`
	//MultiIndexFork *MultiIndexFork
	ForkDbPath string
	fileStream *os.File

	Irreversible include.Signal
}

func GetForkDbInstance(stateDir string) *ForkDatabase {
	forkDB := &ForkDatabase{Index: *fork_multi_index.New()}
	//forkDB, err := newForkDatabase(stateDir, common.DefaultConfig.ForkDbName, true)
	//try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "%s", err)
	return forkDB
}

func newForkDatabase(stateDir string, fileName string, rw bool) (*ForkDatabase, error) {
	fk := &ForkDatabase{Index: *fork_multi_index.New()}
	//_, err := os.Stat(stateDir)
	//if err != nil {
	//	os.MkdirAll(stateDir, os.ModePerm)
	//}
	//fk.ForkDbPath = stateDir + "/" + fileName
	//fStream, err := os.OpenFile(fk.ForkDbPath, os.O_RDWR, os.ModePerm)
	//if fStream == nil {
	//	fStream, err = os.Create(fk.ForkDbPath)
	//	try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "%s", err)
	//}
	//
	//b, err := ioutil.ReadAll(fStream)
	//try.EosAssert(err == nil, &exception.ForkDatabaseException{}, "%s", err)
	//
	//if len(b) > 0 {
	//	mif := &MultiIndexFork{}
	//	mif.Indexs = make(map[string]*IndexFork)
	//	index := &IndexFork{Target: "byBlockId", Uniqueness: true}
	//	index2 := &IndexFork{Target: "byPrev", Uniqueness: false}
	//	index3 := &IndexFork{Target: "byBlockNum", Uniqueness: false}
	//	index4 := &IndexFork{Target: "byLibBlockNum", Uniqueness: false}
	//	mif.Indexs["byBlockId"] = index
	//	mif.Indexs["byPrev"] = index2
	//	mif.Indexs["byBlockNum"] = index3
	//	mif.Indexs["byLibBlockNum"] = index4
	//	bt := treeset.NewMultiWith(types.BlockIdTypes, types.CompareBlockId)
	//	index.Value = bt
	//	mif.Indexs[index.Target] = index
	//	bt2 := treeset.NewMultiWith(types.BlockIdTypes, types.ComparePrev)
	//	index2.Value = bt2
	//	mif.Indexs[index2.Target] = index2
	//	bt3 := treeset.NewMultiWith(types.BlockNumType, types.CompareBlockNum)
	//	index3.Value = bt3
	//	mif.Indexs[index3.Target] = index3
	//	bt4 := treeset.NewMultiWith(types.BlockNumType, types.CompareLibNum)
	//	index4.Value = bt4
	//	mif.Indexs[index4.Target] = index4
	//	fmt.Println("forkDatabase:", len(mif.Indexs), mif.Indexs)
	//	fmt.Println("value:  ", reflect.TypeOf(mif.Indexs["byBlockId"].Value).String())
	//
	//	fmt.Println("78:  ", index.Value.ValueType.String())
	//
	//	err := rlp.DecodeBytes(b, mif)
	//	fmt.Println("forkDatabase error:", err)
	//	fmt.Println("forkDatabase2:", len(mif.Indexs), mif.Indexs)
	//	fk.MultiIndexFork = mif
	//
	//}
	//
	//fk.fileStream = fStream
	//err = os.Remove(fk.ForkDbPath)
	return fk, nil
}

func (f *ForkDatabase) SetHead(s *types.BlockState) {
	inserted := f.Index.Insert(s)
	EosAssert(s.BlockId == s.Header.BlockID(), &ForkDatabaseException{},
		"block state id:%d, is different from block state header id:%d", s.BlockId, s.Header.BlockID())

	EosAssert(inserted, &ForkDatabaseException{}, "unable to insert block state, duplicate state detected")
	if f.Head == nil {
		f.Head = s
	} else if f.Head.BlockNum < s.BlockNum {
		f.Head = s
	}
}

func (f *ForkDatabase) AddBlockState(b *types.BlockState) *types.BlockState {
	EosAssert(b != nil, &ForkDatabaseException{}, "attempt to add null block state")
	EosAssert(f.Head != nil, &ForkDbBlockNotFound{}, "no head block set")

	inserted := f.Index.Insert(b)
	EosAssert(inserted, &ForkDatabaseException{}, "duplicate block added?")

	libItr := f.Index.ByLibBlockNum.Begin()
	f.Head = f.Index.Value(libItr.Value())
	lib := f.Head.DposIrreversibleBlocknum

	nitr := f.Index.ByBlockNum.Begin()
	oldest := f.Index.Value(nitr.Value())

	if oldest.BlockNum < lib {
		f.Prune(oldest)
	}

	return b
}

func (f *ForkDatabase) AddSignedBlock(b *types.SignedBlock, trust bool) *types.BlockState {
	EosAssert(b != nil, &ForkDatabaseException{}, "attempt to add null block")
	EosAssert(f.Head != nil, &ForkDbBlockNotFound{}, "no head block set")

	byIdIdx := f.Index.ByBlockId
	_, existing := byIdIdx[b.BlockID()]
	EosAssert(!existing, &ForkDatabaseException{}, "we already know about this block")

	prior := f.Index.Value(byIdIdx[b.Previous])
	EosAssert(prior != nil, &ForkDatabaseException{}, "unlinkable block:%s,%s", b.BlockID().String(), b.Previous.String())

	result := types.NewBlockState3(&prior.BlockHeaderState, b, trust)
	EosAssert(result != nil, &ForkDatabaseException{}, "fail to add new block state")
	return f.AddBlockState(result)
}

func (f *ForkDatabase) Add(c *types.HeaderConfirmation) {
	header := f.GetBlock(&c.BlockId)

	EosAssert(header != nil, &ForkDbBlockNotFound{}, "unable to find block id:%#v ", c.BlockId)
	header.AddConfirmation(c)

	if header.BftIrreversibleBlocknum < header.BlockNum && len(header.Confirmations) >= ((len(header.ActiveSchedule.Producers)*2)/3+1) {
		f.SetBftIrreversible(c.BlockId)
	}
}
func (f *ForkDatabase) Header() *types.BlockState { return f.Head }

type FetchBranch struct {
	first  []*types.BlockState
	second []*types.BlockState
}

func (f *ForkDatabase) FetchBranchFrom(first *common.BlockIdType, second *common.BlockIdType) FetchBranch {
	result := FetchBranch{}
	firstBranch := f.GetBlock(first)
	secondBranch := f.GetBlock(second)

	for firstBranch.BlockNum > secondBranch.BlockNum {
		result.first = append(result.first, firstBranch)
		firstBranch = f.GetBlock(&firstBranch.Header.Previous)
		EosAssert(firstBranch != nil, &ForkDbBlockNotFound{}, "block %s does not exist", firstBranch.Header.Previous)
	}

	for secondBranch.BlockNum > firstBranch.BlockNum {
		result.second = append(result.second, secondBranch)
		secondBranch = f.GetBlock(&secondBranch.Header.Previous)
		EosAssert(secondBranch != nil, &ForkDbBlockNotFound{}, "block %s does not exist", secondBranch.Header.Previous)
	}

	for firstBranch.Header.Previous != secondBranch.Header.Previous {
		result.first = append(result.first, firstBranch)
		result.second = append(result.second, secondBranch)
		firstBranch = f.GetBlock(&firstBranch.Header.Previous)
		secondBranch = f.GetBlock(&secondBranch.Header.Previous)
		EosAssert(firstBranch != nil && secondBranch != nil, &ForkDbBlockNotFound{},
			"either block %s or %s does not exist",
			firstBranch.Header.Previous, secondBranch.Header.Previous)
	}

	if firstBranch != nil && secondBranch != nil {
		result.first = append(result.first, firstBranch)
		result.second = append(result.second, secondBranch)
	}

	return result
}

/// remove all of the invalid forks built of this id including this id
func (f *ForkDatabase) Remove(id *common.BlockIdType) {
	removeQueue := []common.BlockIdType{*id}
	for i := 0; i < len(removeQueue); i++ {
		itr, existing := f.Index.ByBlockId[removeQueue[i]]
		if existing {
			f.Index.Erase(itr)
		}

		prevIdx := f.Index.ByPrev
		prevItr := prevIdx.LowerBound(removeQueue[i])
		for prevItr.HasNext() {
			if bsp := f.Index.Value(prevItr.Value()); bsp.Header.Previous == removeQueue[i] {
				removeQueue = append(removeQueue, bsp.BlockId)
				prevItr.Next()
				continue
			}
			break
		}

	}

	libItr := f.Index.ByLibBlockNum.Begin()
	f.Head = f.Index.Value(libItr.Value())
}

func (f *ForkDatabase) SetValidity(h *types.BlockState, valid bool) {
	if !valid {
		f.Remove(&h.BlockId)
	} else {
		/// remove older than irreversible and mark block as valid
		h.Validated = true
	}
}
func (f *ForkDatabase) MarkInCurrentChain(h *types.BlockState, inCurrentChain bool) {
	if h.InCurrentChain == inCurrentChain {
		return
	}

	byIdIdx := f.Index.ByBlockId
	itr, existing := byIdIdx[h.BlockId]
	EosAssert(existing, &ForkDbBlockNotFound{}, "could not find block in fork database")

	f.Index.Modify(itr, func(bsp *types.BlockState) {
		bsp.InCurrentChain = inCurrentChain
	})
}

func (f *ForkDatabase) Prune(h *types.BlockState) {
	num := h.BlockNum

	byBn := f.Index.ByBlockNum
	bni := byBn.Begin()
	for bni.HasNext() {
		if bsp := f.Index.Value(bni.Value()); bsp.BlockNum < num {
			f.Prune(bsp)
			bni = byBn.Begin()
			continue
		}
		break
	}

	itr, existing := f.Index.ByBlockId[h.BlockId]
	if existing {
		//TODO
		f.Irreversible.Emit(f.Index.Value(itr))
		f.Index.Erase(itr)
	}

	numIdx := f.Index.ByBlockNum
	nitr := numIdx.LowerBound(fork_multi_index.ByBlockNumComposite{BlockNum: &num})

	for nitr.HasNext() {
		if itrToRemove := f.Index.Value(nitr.Value()); itrToRemove.BlockNum == num {
			nitr.Next()
			id := itrToRemove.BlockId
			f.Remove(&id)
			continue
		}
		break
		//		itrToRemove := f.Index.Value(nitr.Value())
		//		if itrToRemove != nil && itrToRemove.BlockNum == num{
		//			nitr.Next()
		//			id := itrToRemove.BlockId
		//			f.Remove(&id)
		//			continue
		//		} else {
		//			if itrToRemove != nil {
		//fmt.Println("IndexKey:", itrToRemove.BlockNum)
		//			} else {
		////fmt.Println("IndexKey is nil index:", nitr.Value())
		//			}
		//			break
		//		}
	}
}

func (f *ForkDatabase) GetBlock(id *common.BlockIdType) *types.BlockState {
	b, existing := f.Index.Find(*id)
	if existing {
		return b
	}
	return nil
}

func (f *ForkDatabase) GetBlockInCurrentChainByNum(n uint32) *types.BlockState {
	numIdx := f.Index.ByBlockNum
	nitr := numIdx.LowerBound(fork_multi_index.ByBlockNumComposite{BlockNum: &n})

	if nitr.IsEnd() {
		return nil
	}

	if bsp := f.Index.Value(nitr.Value()); bsp.BlockNum != n || bsp.InCurrentChain != true {
		return nil
	} else {
		return bsp
	}
}

func (f *ForkDatabase) SetBftIrreversible(id common.BlockIdType) {
	idx := f.Index.ByBlockId
	itr := idx[id]
	blockNum := f.Index.Value(itr).BlockNum
	f.Index.Modify(itr, func(bsp *types.BlockState) {
		bsp.BftIrreversibleBlocknum = bsp.BlockNum
	})

	update := func(in []common.BlockIdType) []common.BlockIdType {
		updated := []common.BlockIdType{}
		for _, i := range in {
			pidx := f.Index.ByPrev
			pitr := pidx.LowerBound(i)
			epitr := pidx.UpperBound(i)
			for pitr != epitr {
				f.Index.Modify(pitr.Value(), func(bsp *types.BlockState) {
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
	for len(queue) > 0 {
		update(queue)
	}
}

func (f *ForkDatabase) Close() {

	log.Info("ForkDatabase Close tmp code")
	f.Index.Clear()
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
