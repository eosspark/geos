package chain

import (
	"bytes"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/chain/types/forkdb_multi_index"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/plugins/appbase/app/include"
	"io/ioutil"
	"os"
	"strconv"
)

type ForkDatabase struct {
	Index      *forkdb_multi_index.MultiIndex
	Head       *types.BlockState `json:"head"`
	dataDir    string
	fileStream *os.File

	Irreversible include.Signal
}

func NewForkDatabase(dataDir string) *ForkDatabase {
	f := &ForkDatabase{Index: forkdb_multi_index.NewMultiIndex()}
	f.dataDir = dataDir

	if !common.FileExist(dataDir) {
		os.MkdirAll(dataDir, os.ModePerm)
	}

	forkDbDat := f.dataDir + common.DefaultConfig.ForkDbName
	if common.FileExist(forkDbDat) {
		content, err := ioutil.ReadFile(forkDbDat)
		Throw(err)

		decode := rlp.NewDecoder(content)
		var size uint
		decode.Decode(&size)

		for i := uint(0); i < size; i++ {
			s := types.BlockState{}
			decode.Decode(&s)
			f.SetHead(&s)
		}

		headId := common.BlockIdType{}
		decode.Decode(&headId)

		f.Head = f.GetBlock(&headId)
		err = os.Remove(forkDbDat)
		Throw(err)
	}

	return f
}

func (f *ForkDatabase) Close() {
	if f.Index.Size() == 0 {
		return
	}

	//TODO pack BlockState to forkDbDat
	/*forkDbDat := f.dataDir + common.DefaultConfig.ForkDbName
	file, err := os.OpenFile(forkDbDat, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModeAppend)
	Throw(err)

	out := rlp.NewEncoder(file)

	numBlockInForkDB := uint(f.Index.Size())
	out.Encode(numBlockInForkDB)

	f.Index.ByBlockNum.Each(func(key fork_multi_index.ByBlockNumComposite, value fork_multi_index.IndexKey) {
		s := f.Index.Value(value)
		out.Encode(s)
	})

	if f.Head != nil {
		out.Encode(f.Head.BlockId)
	} else {
		out.Encode(crypto.NewSha256Nil())
	}*/

	// we don't normally indicate the head block as irreversible
	// we cannot normally prune the lib if it is the head block because
	// the next block needs to build off of the head block. We are exiting
	// now so we can prune this block as irreversible before exiting.
	lib := f.Head.DposIrreversibleBlocknum
	oldest := f.Index.GetByBlockNum().Begin().Value()
	if oldest.BlockNum <= lib {
		f.Prune(oldest)
	}

	f.Index.Clear()
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
	inserted := f.Index.Insert(b)
	EosAssert(inserted, &ForkDatabaseException{}, "duplicate block added?")

	f.Head = f.Index.GetByLibBlockNum().Begin().Value()

	lib := f.Head.DposIrreversibleBlocknum
	oldest := f.Index.GetByBlockNum().Begin().Value()

	if oldest.BlockNum < lib {
		f.Prune(oldest)
	}

	return b
}

func (f *ForkDatabase) AddSignedBlock(b *types.SignedBlock, trust bool) *types.BlockState {
	EosAssert(b != nil, &ForkDatabaseException{}, "attempt to add null block")
	EosAssert(f.Head != nil, &ForkDbBlockNotFound{}, "no head block set")

	byIdIdx := f.Index.GetByBlockId()
	_, existing := byIdIdx.Find(b.BlockID())
	EosAssert(!existing, &ForkDatabaseException{}, "we already know about this block")

	prior, hasPrior := byIdIdx.Find(b.Previous)
	EosAssert(hasPrior, &ForkDatabaseException{}, "unlinkable block:%s,%s", b.BlockID().String(), b.Previous.String())

	bhs := prior.Value().BlockHeaderState
	result := types.NewBlockState3(&bhs, b, trust)
	EosAssert(result != nil, &ForkDatabaseException{}, "fail to add new block state")
	return f.AddBlockState(result)
}

func (f *ForkDatabase) AddConfirmation(c *types.HeaderConfirmation) {
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

/**
 *  Given two head blocks, return two branches of the fork graph that
 *  end with a common ancestor (same prior block)
 */
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
		itr, existing := f.Index.GetByBlockId().Find(removeQueue[i])
		if existing {
			f.Index.Erase(itr)
		}

		prevIdx := f.Index.GetByPrev()
		prevItr := prevIdx.LowerBound(removeQueue[i])
		for !prevItr.IsEnd() && prevItr.Value().Header.Previous == removeQueue[i] {
			removeQueue = append(removeQueue, prevItr.Value().BlockId)
			prevItr.Next()
		}

	}

	f.Head = f.Index.GetByLibBlockNum().Begin().Value()
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

	byIdIdx := f.Index.GetByBlockId()
	itr, existing := byIdIdx.Find(h.BlockId)
	EosAssert(existing, &ForkDbBlockNotFound{}, "could not find block in fork database")

	f.Index.Modify(itr, func(bsp *forkdb_multi_index.BlockStatePtr) {
		(*bsp).InCurrentChain = inCurrentChain
	})
}

func (f *ForkDatabase) Prune(h *types.BlockState) {
	num := h.BlockNum

	byBn := f.Index.GetByBlockNum()
	bni := byBn.Begin()
	for !bni.IsEnd() && bni.Value().BlockNum < num {
		f.Prune(bni.Value())
		bni = byBn.Begin()
	}

	itr, existing := f.Index.GetByBlockId().Find(h.BlockId)
	if existing {
		f.Irreversible.Emit(itr.Value())
		f.Index.Erase(itr)
	}

	numIdx := f.Index.GetByBlockNum()
	nitr := numIdx.LowerBound(forkdb_multi_index.ByBlockNumComposite{BlockNum: &num})

	for !nitr.IsEnd() && nitr.Value().BlockNum == num {
		itrToRemove := nitr
		nitr.Next()
		id := itrToRemove.Value().BlockId
		f.Remove(&id)
	}
}

func (f *ForkDatabase) GetBlock(id *common.BlockIdType) *types.BlockState {
	b, existing := f.Index.GetByBlockId().Find(*id)
	if existing {
		return b.Value()
	}
	return nil
}

func (f *ForkDatabase) GetBlockInCurrentChainByNum(n uint32) *types.BlockState {
	numIdx := f.Index.GetByBlockNum()
	nitr := numIdx.LowerBound(forkdb_multi_index.ByBlockNumComposite{BlockNum: &n})

	if nitr.IsEnd() || nitr.Value().BlockNum != n || nitr.Value().InCurrentChain != true {
		return nil
	}
	return nitr.Value()
}

/**
 *  This method will set this block as being BFT irreversible and will update
 *  all blocks which build off of it to have the same bft_irb if their existing
 *  bft irb is less than this block num.
 *
 *  This will require a search over all forks
 */
func (f *ForkDatabase) SetBftIrreversible(id common.BlockIdType) {
	idx := f.Index.GetByBlockId()
	itr, _ := idx.Find(id)
	blockNum := itr.Value().BlockNum
	idx.Modify(itr, func(bsp *forkdb_multi_index.BlockStatePtr) {
		(*bsp).BftIrreversibleBlocknum = (*bsp).BlockNum
	})

	/** to prevent stack-overflow, we perform a bredth-first traversal of the
	 * fork database. At each stage we iterate over the leafs from the prior stage
	 * and find all nodes that link their previous. If we update the bft lib then we
	 * add it to a queue for the next layer.  This lambda takes one layer and returns
	 * all block ids that need to be iterated over for next layer.
	 */
	update := func(in []common.BlockIdType) []common.BlockIdType {
		updated := []common.BlockIdType{}
		for _, i := range in {
			pidx := f.Index.GetByPrev()
			pitr := pidx.LowerBound(i)
			epitr := pidx.UpperBound(i)
			for pitr != epitr {
				pidx.Modify(pitr, func(bsp *forkdb_multi_index.BlockStatePtr) {
					if (*bsp).BftIrreversibleBlocknum < blockNum {
						(*bsp).BftIrreversibleBlocknum = blockNum
						updated = append(updated, (*bsp).BlockId)
					}
				})
				pitr.Next()
			}
		}
		return updated
	}

	queue := []common.BlockIdType{id}
	for len(queue) > 0 {
		queue = update(queue)
	}
}

func (f *ForkDatabase) ToString() string {
	var buffer bytes.Buffer
	buffer.WriteString("forkdb current head status:")
	buffer.WriteString("\n[")
	buffer.WriteString("DposIrreversibleBlocknum:")
	buffer.WriteString(strconv.Itoa(int(f.Head.DposIrreversibleBlocknum)))
	buffer.WriteString("\n")
	buffer.WriteString("BftIrreversibleBlocknum:")
	buffer.WriteString(strconv.Itoa(int(f.Head.BftIrreversibleBlocknum)))
	buffer.WriteString("\n")
	buffer.WriteString("Blocknum:")
	buffer.WriteString(strconv.Itoa(int(f.Head.BlockNum)))
	buffer.WriteString("\n")
	buffer.WriteString("InCurrentChain:")
	if f.Head.InCurrentChain {
		buffer.WriteString("true")
	} else {
		buffer.WriteString("false")
	}

	buffer.WriteString("]")
	return buffer.String()
}
