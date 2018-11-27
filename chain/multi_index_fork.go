package chain

import (
	"errors"
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
	Value      BucketFork
}

type IteratorFork struct {
	CurrentSub int
	Value      *types.BlockState
	Idx        *IndexFork
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
	if m.Indexs[index.Target].Value.Len() > 0 {
		m.Indexs[index.Target].Value.Insert(b)
	} else {
		bt := BucketFork{}
		bt.Compare = types.CompareBlockId
		bt.Insert(b)
		index.Value = bt
		m.Indexs["byBlockId"] = index
	}

	index2 := &IndexFork{}
	index2.Target = "byPrev"
	index2.Uniqueness = false
	index2.Less = true

	if m.Indexs[index2.Target].Value.Len() > 0 {
		m.Indexs[index2.Target].Value.Insert(b)
	} else {
		bt := BucketFork{}
		bt.Compare = types.ComparePrev
		bt.Insert(b)
		index2.Value = bt
		m.Indexs[index2.Target] = index2
	}

	index3 := &IndexFork{}
	index3.Target = "byBlockNum"
	index3.Uniqueness = false
	index3.Less = true

	if m.Indexs[index3.Target].Value.Len() > 0 {
		m.Indexs[index3.Target].Value.Insert(b)
	} else {
		bt := BucketFork{}
		bt.Compare = types.CompareBlockNum
		bt.Insert(b)
		index3.Value = bt
		m.Indexs[index3.Target] = index3
	}

	index4 := &IndexFork{}
	index4.Target = "byLibBlockNum"
	index4.Uniqueness = false
	index4.Less = false
	if m.Indexs[index4.Target].Value.Len() > 0 {
		m.Indexs[index4.Target].Value.Insert(b)
	} else {
		bt := BucketFork{}
		bt.Compare = types.CompareLibNum
		bt.Insert(b)
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

func (idx *IndexFork) Begin() (*types.BlockState, error) { //syscall.Mmap()

	if idx.Value.Len() > 0 {
		return idx.Value.Data[0], nil
	} else {
		return nil, errors.New("MultiIndexFork Begin : iterator is nil")
	}
}

func (idx *IndexFork) upperBound(b *types.BlockState) *IteratorFork {
	itr := IteratorFork{}
	itr.Idx = idx
	obj, sub := idx.Value.UpperBound(b)
	if sub >= 0 {
		itr.Value = obj
		itr.CurrentSub = sub
	}
	return &itr
}

func (idx *IndexFork) lowerBound(b *types.BlockState) *IteratorFork {
	itr := IteratorFork{}
	itr.Idx = idx
	obj, sub := idx.Value.LowerBound(b)
	if sub >= 0 {
		itr.Value = obj
		itr.CurrentSub = sub
	}
	return &itr
}

func (itr *IteratorFork) next() bool {
	itr.CurrentSub++
	if itr.CurrentSub < itr.Idx.Value.Len() {
		itr.Value = itr.Idx.Value.Data[itr.CurrentSub]
		return true
	} else {
		return false
	}
}

func (m *MultiIndexFork) find(id common.BlockIdType) *types.BlockState {
	b := types.BlockState{}
	b.BlockId = id
	idx := m.Indexs["byBlockId"]
	bucket := idx.Value
	exist, sub := bucket.Find(&b)
	if exist {
		return idx.Value.Data[sub]
	} else {
		return nil
	}
}

func (m *MultiIndexFork) FindByPrev(prev common.BlockIdType) *types.BlockState {
	b := types.BlockState{}
	b.Header.Previous = prev
	idx := m.Indexs["byBlockId"]
	bucket := idx.Value
	exist, sub := bucket.Find(&b)
	if exist {
		return idx.Value.Data[sub]
	} else {
		return nil
	}
}

func (m *MultiIndexFork) erase(b *types.BlockState) bool {
	if len(m.Indexs) > 0 {
		for _, v := range m.Indexs {
			bt := v.Value
			ext, _ := bt.Find(b)
			if ext {
				v.Value.Eraser(b)
			}
		}
	}
	return true
}

func (m *MultiIndexFork) modify(b *types.BlockState) {
	m.erase(b)
	m.Insert(b)
}

/*type IteratorFork interface {

	Next() bool

	Prev() bool

	Key() []byte

	Value() []byte
}*/
