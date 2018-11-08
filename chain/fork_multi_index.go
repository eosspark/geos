package chain

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/log"
	"strings"
)

type IndexElement struct {
	Key   []byte
	value interface{}
}

func (idx IndexElement) GetKey() []byte {
	return idx.Key
}

type ForkMultiIndex struct {
	Indexs map[string]*Index
}

type Index struct {
	Target     string
	Uniqueness bool
	Less       bool
	KeyValue   common.FlatSet
}

type FkIterator struct {
	Key    []byte
	KeySet common.FlatSet
}

func NewForkMultiIndex() *ForkMultiIndex {
	mi := ForkMultiIndex{}
	mi.Indexs = make(map[string]*Index)
	index := &Index{Target: "byBlockId", Uniqueness: true}
	index2 := &Index{Target: "byPrev", Uniqueness: false}
	index3 := &Index{Target: "byBlockNum", Uniqueness: false}
	index4 := &Index{Target: "byLibBlockNum", Uniqueness: false}
	mi.Indexs["byBlockId"] = index
	mi.Indexs["byPrev"] = index2
	mi.Indexs["byBlockNum"] = index3
	mi.Indexs["byLibBlockNum"] = index4

	return &mi
}

func (m *ForkMultiIndex) Insert(b *types.BlockState) bool {

	index := &Index{}
	index.Target = "byBlockId"
	index.Uniqueness = true
	mainKey := computeMainKey(b.BlockId.Bytes())
	indexElement := IndexElement{mainKey, &b}
	prevKey := computePrevKey(b.BlockHeaderState.Header.Previous.Bytes())
	uniqueKey2 := append(prevKey, '_')
	uniqueKey2 = append(uniqueKey2, mainKey...)
	indexElement2 := IndexElement{[]byte(uniqueKey2), &mainKey}

	blockNumKey := computeBlockNumKey(b.BlockNum, b.InCurrentChain)
	uniqueKey3 := append(blockNumKey, '_')
	uniqueKey3 = append(uniqueKey3, mainKey...)
	indexElement3 := IndexElement{[]byte(uniqueKey3), &mainKey}
	libBlockNumKey := computeLibNumKey(b.BlockHeaderState.DposIrreversibleBlocknum, b.BlockHeaderState.BftIrreversibleBlocknum, b.BlockNum)
	uniqueKey4 := append(libBlockNumKey, '_')
	uniqueKey4 = append(uniqueKey3, mainKey...)
	indexElement4 := IndexElement{[]byte(uniqueKey4), &mainKey}
	if m.Indexs[index.Target].KeyValue.Len() > 0 {
		m.Indexs[index.Target].KeyValue.Insert(&indexElement)
	} else {
		fs := common.FlatSet{}
		fs.Insert(&indexElement)
		index.KeyValue = fs
		m.Indexs["byBlockId"] = index
	}
	index2 := &Index{}
	index2.Target = "byPrev"
	index2.Uniqueness = false
	index2.Less = true

	if m.Indexs[index2.Target].KeyValue.Len() > 0 {
		m.Indexs[index2.Target].KeyValue.Insert(&indexElement2)
	} else {
		fs2 := common.FlatSet{}
		fs2.Insert(&indexElement2)
		index2.KeyValue = fs2
		m.Indexs[index2.Target] = index2
	}

	index3 := &Index{}
	index3.Target = "byBlockNum"
	index3.Uniqueness = false
	index3.Less = true
	if m.Indexs[index3.Target].KeyValue.Len() > 0 {
		m.Indexs[index3.Target].KeyValue.Insert(&indexElement3)
	} else {
		fs3 := common.FlatSet{}
		fs3.Insert(&indexElement3)
		index3.KeyValue = fs3
		m.Indexs[index3.Target] = index3
	}

	index4 := &Index{}
	index4.Target = "byLibBlockNum"
	index4.Uniqueness = false
	index4.Less = false
	if m.Indexs[index4.Target].KeyValue.Len() > 0 {
		m.Indexs[index4.Target].KeyValue.Insert(&indexElement4)
	} else {
		fs4 := common.FlatSet{}
		fs4.Insert(&indexElement4)
		index4.KeyValue = fs4
		m.Indexs[index4.Target] = index4
	}

	return true
}

func (m *ForkMultiIndex) GetIndex(tag string) *Index {
	if index, ok := m.Indexs[tag]; ok {
		return index
	}
	return nil
}

func (idx *Index) Begin() *FkIterator { //syscall.Mmap()
	itr := FkIterator{}
	if idx.KeyValue.Len() > 0 {
		idxEle := idx.KeyValue.Data[0]
		itr.Key = idxEle.GetKey()
		itr.KeySet = idx.KeyValue
	}
	return &itr
}

/*func (idx *Index) Next() *FkIterator{

}*/

func (idx *Index) UpperBound(b []byte) *FkIterator {
	itr := FkIterator{}
	if idx.KeyValue.Len() > 0 {
		for _, idxEle := range idx.KeyValue.Data {
			tagKey := idxEle.(*IndexElement).GetKey()
			if bytes.Compare(tagKey, b) == 1 {
				itr.Key = tagKey
				break
			}
		}
		return idx.LowerBound(itr.Key)
	}
	return nil
}

func (idx *Index) searchSub(key []byte) int {
	length := idx.KeyValue.Len()
	i, j := 0, length-1
	for i < j {
		h := int(uint(i+j) >> 1)
		if i <= h && h < j {
			ext := strings.Index(string(idx.KeyValue.Data[h].GetKey()), string(key))
			if ext >= 0 {
				return h
			} else {
				i = h + 1
			}
			/*if bytes.Compare(idx.KeyValue.Data[h].GetKey(), key) == -1 {
				i = h + 1
			} else if bytes.Compare(idx.KeyValue.Data[h].GetKey(), key) == 0 {
				return h
			} else {
				j = h
			}*/
		}
	}
	return i
}

func (idx *Index) LowerBound(b []byte) *FkIterator {
	itr := FkIterator{}
	first := 0
	if idx.KeyValue.Len() > 0 {
		//start
		ext := idx.searchSub(b)
		idxEle := idx.KeyValue.Data[ext]
		itr.Key = idxEle.(*IndexElement).GetKey()
		first = ext
		/*for s,idxEle := range idx.KeyValue.Data{
			tagKey:=idxEle.(*IndexElement).GetKey()
			ext:=strings.Index(string(tagKey),string(b))
			if ext>=0{
				itr.Key = tagKey
				first =s
				break
			}
		}*/
		//end
		/*for s,idxEle := range idx.KeyValue.Data[first:]{
			tagKey:=idxEle.(*IndexElement).GetKey()
			ext:=strings.Index(string(tagKey),string(b))
			if ext<0{
				itr.KeySet.Data = idx.KeyValue.Data[first:s]
				break
			}
		}*/
		i, j := 0, idx.KeyValue.Len()-1
		for i < j {
			h := int(i + j>>1)
			if i <= h && h < j {
				et := strings.Index(string(idx.KeyValue.Data[h].GetKey()), string(b))
				if et < 0 {
					itr.KeySet.Data = idx.KeyValue.Data[first:h]
					break
				} else {
					i = h + 1
				}
			}
		}
		if len(itr.Key) > 0 && itr.KeySet.Len() == 0 {
			itr.KeySet.Data = idx.KeyValue.Data[first:]
		}
		return &itr
	}
	return nil
}

func (m *ForkMultiIndex) FindById(id common.BlockIdType) *types.BlockState {
	mk := computeMainKey(id.Bytes())
	idx := m.Indexs["byBlockId"]
	fs := idx.KeyValue
	idxElements, _ := fs.FindData(mk)
	return idxElements.(*IndexElement).value.(*types.BlockState)
}

func (m *ForkMultiIndex) FindByPrev(prev common.BlockIdType) *types.BlockState {
	key := computePrevKey(prev.Bytes())
	idx := m.Indexs["byPrev"]
	idxEle, _ := idx.KeyValue.FindData(key)
	mainKey := idxEle.(*IndexElement).value
	idxe := m.GetIndex("byBlockId")
	objEle, _ := idxe.KeyValue.FindData(mainKey.([]byte))
	bs := objEle.(*IndexElement).value.(*types.BlockState)
	return bs
}

/*func (m *ForkMultiIndex) FindByBlockNum(prev uint32,inCurrentChain bool) *types.BlockState {
	key := computePrevKey(prev.Bytes())
	idx := m.Indexs["byLibBlockNum"]
	idxEle, _ := idx.KeyValue.FindData(key)
	mainKey := idxEle.(*IndexElement).value
	idxe := m.GetIndex("byBlockId")
	objEle, _ := idxe.KeyValue.FindData(mainKey.([]byte))
	bs := objEle.(*IndexElement).value.(*types.BlockState)
	return bs
}*/

/*func (m *ForkMultiIndex) FindByLibBlockNum(prev uint32,inCurrentChain bool) *types.BlockState {
	key := computePrevKey(prev.Bytes())
	idx := m.Indexs["byLibBlockNum"]
	idxEle, _ := idx.KeyValue.FindData(key)
	mainKey := idxEle.(*IndexElement).value
	idxe := m.GetIndex("byBlockId")
	objEle, _ := idxe.KeyValue.FindData(mainKey.([]byte))
	bs := objEle.(*IndexElement).value.(*types.BlockState)
	return bs
}*/

func (idx *Index) eraseMainKey(id *common.BlockIdType) bool {
	keyArray := make([][]byte, 4)
	key := computeMainKey(id.Bytes())
	keyArray = append(keyArray, key)
	ele, sub := idx.KeyValue.FindData(key)
	if sub >= 0 {
		block := ele.(*types.BlockState)
		prevKey := computePrevKey(block.SignedBlock.Previous.Bytes())
		keyArray = append(keyArray, prevKey)
		blockNumKey := computeBlockNumKey(block.BlockNum, block.InCurrentChain)
		keyArray = append(keyArray, blockNumKey)
		libNumKey := computeLibNumKey(block.DposIrreversibleBlocknum, block.BftIrreversibleBlocknum, block.BlockNum)
		keyArray = append(keyArray, libNumKey)
		for _, k := range keyArray {
			boo := idx.KeyValue.Remove(k)
			if !boo {
				log.Error("fork_contanier eraseMainKey is error:%#v", k)
			}
		}

	}
	return false
}

func (idx *Index) EraseKey(key []byte) bool {
	idPtr, _ := idx.KeyValue.FindData(key)
	m := idPtr.(*IndexElement).value
	id := crypto.NewSha256Byte(m.([]byte))
	return idx.eraseMainKey(id)
}

func computeMainKey(val []byte) []byte {
	return append([]byte("byBlockId_"), val...)
}

func (m *ForkMultiIndex) modify() {

}

func computePrevKey(val []byte) []byte {
	return append([]byte("byPrev_"), val...)
}

func computeBlockNumKey(blockNum uint32, inCurrent bool) []byte {
	bn := make([]byte, 8)
	if inCurrent {
		binary.BigEndian.PutUint64(bn, uint64(blockNum))
	} else {
		binary.BigEndian.PutUint64(bn, uint64(^blockNum+1))
	}
	return append([]byte("byBlockNum_"), bn...)
}

func computeLibNumKey(dposIrreversibleBlocknum uint32, bftIrreversibleBlocknum uint32, blockNum uint32) []byte {
	str := fmt.Sprintf("%d_%d_%d", dposIrreversibleBlocknum, bftIrreversibleBlocknum, blockNum)
	val := []byte(str)
	return append([]byte("byLibBlockNum_"), val...)
}

/*func (idx *Index) computeManyKey(first []byte, mainKey []byte) {
	uniqueKey := string(first) + "_" + string(mainKey)
	if len(idx.ManyKeyMap) == 0{
		idx.ManyKeyMap = make(map[string]common.FlatSet)
	}
	keySet := idx.ManyKeyMap[string(first)]
	idxEle := IndexElement{[]byte(uniqueKey), []byte(uniqueKey)}
	keySet.Insert(idxEle)
	idx.ManyKeyMap[string(first)] = keySet
}*/

/*type fkIterator interface {
	Next() bool

	Prev() bool

	Key() []byte

	Value() []byte
}*/
