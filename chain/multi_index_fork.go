package chain

import (
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type multiIndexFork struct {
	indexs map[string]*indexFork
}

type indexFork struct {
	target     string
	uniqueness bool
	less       bool
	value      common.Bucket
}

type iteratorFork struct {
	currentSub int
	value      *types.BlockState
	idx        *indexFork
}

func newMultiIndexFork() *multiIndexFork {
	mi := multiIndexFork{}
	mi.indexs = make(map[string]*indexFork)
	index := &indexFork{target: "byBlockId", uniqueness: true}
	index2 := &indexFork{target: "byPrev", uniqueness: false}
	index3 := &indexFork{target: "byBlockNum", uniqueness: false}
	index4 := &indexFork{target: "byLibBlockNum", uniqueness: false}
	mi.indexs["byBlockId"] = index
	mi.indexs["byPrev"] = index2
	mi.indexs["byBlockNum"] = index3
	mi.indexs["byLibBlockNum"] = index4

	return &mi
}

func (m *multiIndexFork) Insert(b *types.BlockState) bool {

	index := &indexFork{}
	index.target = "byBlockId"
	index.uniqueness = true
	if m.indexs[index.target].value.Len() > 0 {
		m.indexs[index.target].value.Insert(b)
	} else {
		bt := common.Bucket{}
		bt.Compare = types.CompareBlockId
		bt.Insert(b)
		index.value = bt
		m.indexs["byBlockId"] = index
	}

	index2 := &indexFork{}
	index2.target = "byPrev"
	index2.uniqueness = false
	index2.less = true

	if m.indexs[index2.target].value.Len() > 0 {
		m.indexs[index2.target].value.Insert(b)
	} else {
		bt := common.Bucket{}
		bt.Compare = types.ComparePrev
		bt.Insert(b)
		index2.value = bt
		m.indexs[index2.target] = index2
	}

	index3 := &indexFork{}
	index3.target = "byBlockNum"
	index3.uniqueness = false
	index3.less = true

	if m.indexs[index3.target].value.Len() > 0 {
		m.indexs[index3.target].value.Insert(b)
	} else {
		bt := common.Bucket{}
		bt.Compare = types.CompareBlockNum
		bt.Insert(b)
		index3.value = bt
		m.indexs[index3.target] = index3
	}

	index4 := &indexFork{}
	index4.target = "byLibBlockNum"
	index4.uniqueness = false
	index4.less = false
	if m.indexs[index4.target].value.Len() > 0 {
		m.indexs[index4.target].value.Insert(b)
	} else {
		bt := common.Bucket{}
		bt.Compare = types.CompareLibNum
		bt.Insert(b)
		index4.value = bt
		m.indexs[index4.target] = index4
	}

	return true
}

func (m *multiIndexFork) GetIndex(tag string) *indexFork {
	if index, ok := m.indexs[tag]; ok {
		return index
	}
	return nil
}

func (idx *indexFork) Begin() (*types.BlockState, error) { //syscall.Mmap()

	if idx.value.Len() > 0 {
		return idx.value.Data[0].(*types.BlockState), nil
	} else {
		return nil, errors.New("MultiIndexFork Begin : iterator is nil")
	}
}

func (idx *indexFork) upperBound(b *types.BlockState) *iteratorFork {
	itr := iteratorFork{}
	itr.idx = idx
	if idx.value.Len() > 0 {
		ext := idx.searchSub(b)
		if idx.less {
			for i := ext; i < idx.value.Len(); i++ {
				if idx.value.Compare(b, idx.value.Data[i]) < 0 {
					itr.value = idx.value.Data[i-1].(*types.BlockState)
					itr.currentSub = i - 1
					break
				} else if i == idx.value.Len()-1 && idx.value.Compare(b, idx.value.Data[i]) == 0 {
					itr.value = idx.value.Data[i].(*types.BlockState)
					itr.currentSub = i
				}
			}
		} else {
			for i := ext; i > 0; i-- {
				if idx.value.Compare(b, idx.value.Data[i]) < 0 {
					itr.value = idx.value.Data[i+1].(*types.BlockState)
					itr.currentSub = i + 1
					break
				} else if i == 0 && idx.value.Compare(b, idx.value.Data[i]) == 0 {
					itr.value = idx.value.Data[i].(*types.BlockState)
					itr.currentSub = i
				}
			}
		}
		return &itr
	}
	return nil
}

func (idx *indexFork) searchSub(b *types.BlockState) int {
	length := idx.value.Len()
	i, j := 0, length-1
	for i < j {
		h := int(uint(i+j) >> 1)
		if i <= h && h < j {
			ext := idx.value.Compare(idx.value.Data[h], b)
			if idx.less {
				if ext < 0 {
					i = h + 1
				} else {
					j = h
				}
			} else {
				if ext > 0 {
					i = h + 1
				} else {
					j = h
				}
			}
		}
	}
	return i
}

func (idx *indexFork) lowerBound(b *types.BlockState) *iteratorFork {
	itr := iteratorFork{}
	itr.idx = idx
	first := 0
	if idx.value.Len() > 0 {
		ext := idx.searchSub(b)
		first = ext
		if idx.less {
			fmt.Println("less search")
			for i := first; i > 0; i-- {
				if idx.value.Compare(idx.value.Data[i], b) == -1 {
					itr.value = idx.value.Data[i+1].(*types.BlockState)
					itr.currentSub = i + 1
					break
				} else if i == 0 && idx.value.Compare(idx.value.Data[i], b) == 0 {
					itr.value = idx.value.Data[i].(*types.BlockState)
					itr.currentSub = i
					break
				}
			}
		} else {
			fmt.Println("greater search")
			for i := first; i < idx.value.Len(); i++ {
				if i >= 1 && idx.value.Compare(idx.value.Data[i], b) == -1 {
					itr.value = idx.value.Data[i-1].(*types.BlockState)
					itr.currentSub = i - 1
					break
				} else if i == 0 && idx.value.Compare(idx.value.Data[i], b) == 0 {
					itr.value = idx.value.Data[i].(*types.BlockState)
					itr.currentSub = i
					break
				}
			}
		}
		return &itr
	}
	return nil
}

func (itr *iteratorFork) next() {
	itr.currentSub++
	itr.value = itr.idx.value.Data[itr.currentSub].(*types.BlockState)
}

func (m *multiIndexFork) Find(id common.BlockIdType) *types.BlockState {
	b := types.BlockState{}
	b.BlockId = id
	idx := m.indexs["byBlockId"]
	bucket := idx.value
	exist, sub := bucket.Find(&b)
	if exist {
		return idx.value.Data[sub].(*types.BlockState)
	} else {
		return nil
	}
}

func (m *multiIndexFork) FindByPrev(prev common.BlockIdType) *types.BlockState {
	b := types.BlockState{}
	b.Header.Previous = prev
	idx := m.indexs["byBlockId"]
	bucket := idx.value
	exist, sub := bucket.Find(&b)
	if exist {
		return idx.value.Data[sub].(*types.BlockState)
	} else {
		return nil
	}
}

func (m *multiIndexFork) erase(b *types.BlockState) bool {
	if len(m.indexs) > 0 {
		for _, v := range m.indexs {
			bt := v.value
			ext, _ := bt.Find(b)
			if ext {
				v.value.Eraser(b)
			}
		}
	}
	return true
}

func (m *multiIndexFork) modify(b *types.BlockState) {
	m.erase(b)
	m.Insert(b)
}

/*type iteratorFork interface {

	Next() bool

	Prev() bool

	Key() []byte

	value() []byte
}*/
