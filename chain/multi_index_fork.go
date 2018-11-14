package chain

import (
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
	KeySet common.Bucket
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

func (idx *indexFork) Begin() *iteratorFork { //syscall.Mmap()
	itr := iteratorFork{}
	if idx.value.Len() > 0 {
		itr.KeySet.Data[0] = idx.value.Data[0]
	}
	return &itr
}

/*func (idx *Index) Next() *iteratorFork{

}*/

func (idx *indexFork) UpperBound(b *types.BlockState) *iteratorFork {
	itr := iteratorFork{}
	var tagObj *types.BlockState
	if idx.value.Len() > 0 {
		for _, idxEle := range idx.value.Data {
			tagObj = idxEle.(*types.BlockState)
			if idx.value.Compare(idxEle.(*types.BlockState), b) == 1 {
				itr.KeySet.Insert(tagObj)
				break
			}
		}
		return idx.LowerBound(tagObj)
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
			if ext < 0 {
				i = h + 1
			} else {
				j = h
			}
		}
	}
	return i
}

func (idx *indexFork) LowerBound(b *types.BlockState) *iteratorFork {
	itr := iteratorFork{}
	first := 0
	if idx.value.Len() > 0 {
		//start
		ext := idx.searchSub(b)
		first = ext
		for i := first; i < idx.value.Len(); i++ {
			if idx.value.Compare(idx.value.Data[i], b) > 0 || (i == idx.value.Len()-1 && idx.value.Compare(idx.value.Data[i], b) == 0) {
				if i == idx.value.Len() {
					itr.KeySet.Data = idx.value.Data[first:idx.value.Len()]
				} else {
					itr.KeySet.Data = idx.value.Data[first : i+1]
				}
				break
			}
		}
		return &itr
	}
	return nil
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

/*func (m *multiIndexFork) FindByBlockNum(prev uint32,inCurrentChain bool) *types.BlockState {
	key := computePrevKey(prev.Bytes())
	idx := m.indexs["byLibBlockNum"]
	idxEle, _ := idx.keyValue.FindData(key)
	mainKey := idxEle.(*indexElement).value
	idxe := m.GetIndex("byBlockId")
	objEle, _ := idxe.keyValue.FindData(mainKey.([]byte))
	bs := objEle.(*indexElement).value.(*types.BlockState)
	return bs
}*/

/*func (m *multiIndexFork) FindByLibBlockNum(prev uint32,inCurrentChain bool) *types.BlockState {
	key := computePrevKey(prev.Bytes())
	idx := m.indexs["byLibBlockNum"]
	idxEle, _ := idx.keyValue.FindData(key)
	mainKey := idxEle.(*indexElement).value
	idxe := m.GetIndex("byBlockId")
	objEle, _ := idxe.keyValue.FindData(mainKey.([]byte))
	bs := objEle.(*indexElement).value.(*types.BlockState)
	return bs
}*/

/*func (idx *indexFork) eraseMainKey(id *common.BlockIdType) bool {
	keyArray := make([][]byte, 4)
	key := computeMainKey(id.BigEndianBytes())
	keyArray = append(keyArray, key)
	ele, sub := idx.keyValue.FindData(key)
	if sub >= 0 {
		block := ele.(*types.BlockState)
		prevKey := computePrevKey(block.SignedBlock.Previous.BigEndianBytes())
		keyArray = append(keyArray, prevKey)
		blockNumKey := computeBlockNumKey(block.BlockNum, block.InCurrentChain)
		keyArray = append(keyArray, blockNumKey)
		libNumKey := computeLibNumKey(block.DposIrreversibleBlocknum, block.BftIrreversibleBlocknum, block.BlockNum)
		keyArray = append(keyArray, libNumKey)
		for _, k := range keyArray {
			boo := idx.keyValue.Remove(indexElement{key: key})
			if !boo {
				log.Error("fork_contanier eraseMainKey is error:%#v", k)
			}
		}

	}
	return false
}*/

func (m *multiIndexFork) erase(b *types.BlockState) bool {
	if len(m.indexs) > 0 {
		for _, v := range m.indexs {
			bt := v.value
			ext, _ := bt.Find(b)
			if ext {
				v.value.Easer(b)
			}
		}
	}
	return true
}

func (m *multiIndexFork) modify(b *types.BlockState) {
	m.erase(b)
	m.Insert(b)
}

func computeLibNumKey(dposIrreversibleBlocknum uint32, bftIrreversibleBlocknum uint32, blockNum uint32) []byte {
	str := fmt.Sprintf("%d_%d_%d", dposIrreversibleBlocknum, bftIrreversibleBlocknum, blockNum)
	val := []byte(str)
	return append([]byte("byLibBlockNum_"), val...)
}

/*type iteratorFork interface {
	Next() bool

	Prev() bool

	Key() []byte

	value() []byte
}*/
