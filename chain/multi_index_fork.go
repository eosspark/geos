package chain

import (
	"errors"
	"github.com/eosspark/container/sets/treeset"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type MultiIndexFork struct {
	Indexs map[string]*IndexFork
}

type IndexFork struct {
	Target     string
	Uniqueness bool
	Less       bool
	Value      *treeset.MultiSet
}

func newMultiIndexFork() *MultiIndexFork {
	mi := MultiIndexFork{}
	mi.Indexs = make(map[string]*IndexFork)
	index := &IndexFork{Target: "byBlockId", Uniqueness: true}
	index2 := &IndexFork{Target: "byPrev", Uniqueness: false}
	index3 := &IndexFork{Target: "byBlockNum", Uniqueness: false}
	index4 := &IndexFork{Target: "byLibBlockNum", Uniqueness: false}
	mi.Indexs["byBlockId"] = index
	mi.Indexs["byPrev"] = index2
	mi.Indexs["byBlockNum"] = index3
	mi.Indexs["byLibBlockNum"] = index4

	return &mi
}

func (m *MultiIndexFork) Insert(b *types.BlockState) bool {

	index := &IndexFork{}
	index.Target = "byBlockId"
	index.Uniqueness = true
	val := m.Indexs[index.Target].Value
	if val != nil && val.Size() >= 0 {
		val.Add(b)
	} else {
		bt := treeset.NewMultiWith(types.BlockIdTypes, types.CompareBlockId)
		bt.Add(b)
		index.Value = bt
		m.Indexs["byBlockId"] = index
	}

	index2 := &IndexFork{}
	index2.Target = "byPrev"
	index2.Uniqueness = false
	index2.Less = true

	if m.Indexs[index2.Target].Value != nil && m.Indexs[index2.Target].Value.Size() > 0 {
		m.Indexs[index2.Target].Value.Add(b)
	} else {
		bt := treeset.NewMultiWith(types.BlockIdTypes, types.ComparePrev)
		//bt.Compare = types.ComparePrev
		bt.Add(b)
		index2.Value = bt
		m.Indexs[index2.Target] = index2
	}

	index3 := &IndexFork{}
	index3.Target = "byBlockNum"
	index3.Uniqueness = false
	index3.Less = true

	if m.Indexs[index3.Target].Value != nil && m.Indexs[index3.Target].Value.Size() > 0 {
		m.Indexs[index3.Target].Value.Add(b)
	} else {
		bt := treeset.NewMultiWith(types.BlockNumType, types.CompareBlockNum)
		//bt.Compare = types.CompareBlockNum
		bt.Add(b)
		index3.Value = bt
		m.Indexs[index3.Target] = index3
	}

	index4 := &IndexFork{}
	index4.Target = "byLibBlockNum"
	index4.Uniqueness = false
	index4.Less = false
	if m.Indexs[index4.Target].Value != nil && m.Indexs[index4.Target].Value.Size() > 0 {
		m.Indexs[index4.Target].Value.Add(b)
	} else {
		bt := treeset.NewMultiWith(types.BlockNumType, types.CompareLibNum)
		//bt.Compare = types.CompareLibNum
		bt.Add(b)
		index4.Value = bt
		m.Indexs[index4.Target] = index4
	}

	return true
}

func (m *MultiIndexFork) GetIndex(tag string) *IndexFork {
	if index, ok := m.Indexs[tag]; ok {
		return index
	}
	return nil
}

func (idx *IndexFork) Begin() (*types.BlockState, error) {
	itr := idx.Value.Iterator()
	itr.Begin()
	if itr.Next() {
		return itr.Value().(*types.BlockState), nil
	}
	return nil, errors.New("MultiIndexFork Begin : iterator is nil")
}

func (idx *IndexFork) Iterator() treeset.MultiSetIterator {
	return idx.Value.Iterator()
}

func (idx *IndexFork) upperBound(b *types.BlockState) *treeset.MultiSetIterator {

	return idx.Value.UpperBound(b)
}

func (idx *IndexFork) lowerBound(b *types.BlockState) *treeset.MultiSetIterator {

	return idx.Value.LowerBound(b)
}

func (m *MultiIndexFork) find(id common.BlockIdType) *types.BlockState {
	b := types.BlockState{}
	b.BlockId = id
	idx := m.Indexs["byBlockId"]
	multiSet := idx.Value
	mItr, exist := multiSet.Get(&b)
	if exist {
		return mItr.Value().(*types.BlockState)
	} else {
		return nil
	}
}

func (m *MultiIndexFork) FindByPrev(prev common.BlockIdType) *types.BlockState {
	b := types.BlockState{}
	b.Header.Previous = prev
	idx := m.Indexs["byBlockId"]
	multiSet := idx.Value
	itr, exist := multiSet.Get(&b)
	if exist {
		return itr.Value().(*types.BlockState)
	} else {
		return nil
	}
}

func (m *MultiIndexFork) erase(b *types.BlockState) {
	if len(m.Indexs) > 0 {
		for _, v := range m.Indexs {
			v.Value.Remove(b)
		}
	}
}

func (m *MultiIndexFork) modify(b *types.BlockState) {
	m.erase(b)
	m.Insert(b)
}
