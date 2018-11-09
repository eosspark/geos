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

type indexElement struct {
	key   []byte
	value interface{}
}

func (idx indexElement) GetKey() []byte {
	return idx.key
}

type forkMultiIndex struct {
	Indexs map[string]*fkIndex
}

type fkIndex struct {
	target     string
	uniqueness bool
	less       bool
	keyValue   common.FlatSet
}

type FkIterator struct {
	Key    []byte
	KeySet common.FlatSet
}

func NewForkMultiIndex() *forkMultiIndex {
	mi := forkMultiIndex{}
	mi.Indexs = make(map[string]*fkIndex)
	index := &fkIndex{target: "byBlockId", uniqueness: true}
	index2 := &fkIndex{target: "byPrev", uniqueness: false}
	index3 := &fkIndex{target: "byBlockNum", uniqueness: false}
	index4 := &fkIndex{target: "byLibBlockNum", uniqueness: false}
	mi.Indexs["byBlockId"] = index
	mi.Indexs["byPrev"] = index2
	mi.Indexs["byBlockNum"] = index3
	mi.Indexs["byLibBlockNum"] = index4

	return &mi
}

func (m *forkMultiIndex) Insert(b *types.BlockState) bool {

	index := &fkIndex{}
	index.target = "byBlockId"
	index.uniqueness = true
	mainKey := computeMainKey(b.BlockId.Bytes())
	idxEle := indexElement{mainKey, b}
	prevKey := computePrevKey(b.BlockHeaderState.Header.Previous.Bytes())
	uniqueKey2 := append(prevKey, '_')
	uniqueKey2 = append(uniqueKey2, mainKey...)
	indexElement2 := indexElement{[]byte(uniqueKey2), &mainKey}

	blockNumKey := computeBlockNumKey(b.BlockNum, b.InCurrentChain)
	uniqueKey3 := append(blockNumKey, '_')
	uniqueKey3 = append(uniqueKey3, mainKey...)
	indexElement3 := indexElement{[]byte(uniqueKey3), &mainKey}
	libBlockNumKey := computeLibNumKey(b.BlockHeaderState.DposIrreversibleBlocknum, b.BlockHeaderState.BftIrreversibleBlocknum, b.BlockNum)
	uniqueKey4 := append(libBlockNumKey, '_')
	uniqueKey4 = append(uniqueKey3, mainKey...)
	indexElement4 := indexElement{[]byte(uniqueKey4), &mainKey}
	if m.Indexs[index.target].keyValue.Len() > 0 {
		m.Indexs[index.target].keyValue.Insert(&idxEle)
	} else {
		fs := common.FlatSet{}
		fs.Insert(&idxEle)
		index.keyValue = fs
		m.Indexs["byBlockId"] = index
	}
	index2 := &fkIndex{}
	index2.target = "byPrev"
	index2.uniqueness = false
	index2.less = true

	if m.Indexs[index2.target].keyValue.Len() > 0 {
		m.Indexs[index2.target].keyValue.Insert(&indexElement2)
	} else {
		fs2 := common.FlatSet{}
		fs2.Insert(&indexElement2)
		index2.keyValue = fs2
		m.Indexs[index2.target] = index2
	}

	index3 := &fkIndex{}
	index3.target = "byBlockNum"
	index3.uniqueness = false
	index3.less = true
	if m.Indexs[index3.target].keyValue.Len() > 0 {
		m.Indexs[index3.target].keyValue.Insert(&indexElement3)
	} else {
		fs3 := common.FlatSet{}
		fs3.Insert(&indexElement3)
		index3.keyValue = fs3
		m.Indexs[index3.target] = index3
	}

	index4 := &fkIndex{}
	index4.target = "byLibBlockNum"
	index4.uniqueness = false
	index4.less = false
	if m.Indexs[index4.target].keyValue.Len() > 0 {
		m.Indexs[index4.target].keyValue.Insert(&indexElement4)
	} else {
		fs4 := common.FlatSet{}
		fs4.Insert(&indexElement4)
		index4.keyValue = fs4
		m.Indexs[index4.target] = index4
	}

	return true
}

func (m *forkMultiIndex) GetIndex(tag string) *fkIndex {
	if index, ok := m.Indexs[tag]; ok {
		return index
	}
	return nil
}

func (idx *fkIndex) Begin() *FkIterator { //syscall.Mmap()
	itr := FkIterator{}
	if idx.keyValue.Len() > 0 {
		idxEle := idx.keyValue.Data[0]
		itr.Key = idxEle.GetKey()
		itr.KeySet = idx.keyValue
	}
	return &itr
}

/*func (idx *Index) Next() *FkIterator{

}*/

func (idx *fkIndex) UpperBound(b []byte) *FkIterator {
	itr := FkIterator{}
	if idx.keyValue.Len() > 0 {
		for _, idxEle := range idx.keyValue.Data {
			tagKey := idxEle.(*indexElement).GetKey()
			if bytes.Compare(tagKey, b) == 1 {
				itr.Key = tagKey
				break
			}
		}
		return idx.LowerBound(itr.Key)
	}
	return nil
}

func (idx *fkIndex) searchSub(key []byte) int {
	length := idx.keyValue.Len()
	i, j := 0, length-1
	for i < j {
		h := int(uint(i+j) >> 1)
		if i <= h && h < j {
			ext := strings.Index(string(idx.keyValue.Data[h].GetKey()), string(key))
			if ext >= 0 {
				return h
			} else {
				i = h + 1
			}
			/*if bytes.Compare(idx.keyValue.Data[h].GetKey(), key) == -1 {
				i = h + 1
			} else if bytes.Compare(idx.keyValue.Data[h].GetKey(), key) == 0 {
				return h
			} else {
				j = h
			}*/
		}
	}
	return i
}

func (idx *fkIndex) LowerBound(b []byte) *FkIterator {
	itr := FkIterator{}
	first := 0
	if idx.keyValue.Len() > 0 {
		//start
		ext := idx.searchSub(b)
		idxEle := idx.keyValue.Data[ext]
		itr.Key = idxEle.(*indexElement).GetKey()
		first = ext
		/*for s,idxEle := range idx.keyValue.Data{
			tagKey:=idxEle.(*indexElement).GetKey()
			ext:=strings.Index(string(tagKey),string(b))
			if ext>=0{
				itr.Key = tagKey
				first =s
				break
			}
		}*/
		//end
		/*for s,idxEle := range idx.keyValue.Data[first:]{
			tagKey:=idxEle.(*indexElement).GetKey()
			ext:=strings.Index(string(tagKey),string(b))
			if ext<0{
				itr.KeySet.Data = idx.keyValue.Data[first:s]
				break
			}
		}*/
		i, j := 0, idx.keyValue.Len()-1
		for i < j {
			h := int(i + j>>1)
			if i <= h && h < j {
				et := strings.Index(string(idx.keyValue.Data[h].GetKey()), string(b))
				if et < 0 {
					itr.KeySet.Data = idx.keyValue.Data[first:h]
					break
				} else {
					i = h + 1
				}
			}
		}
		if len(itr.Key) > 0 && itr.KeySet.Len() == 0 {
			itr.KeySet.Data = idx.keyValue.Data[first:]
		}
		return &itr
	}
	return nil
}

func (m *forkMultiIndex) FindById(id common.BlockIdType) *types.BlockState {
	mk := computeMainKey(id.Bytes())
	idx := m.Indexs["byBlockId"]
	fs := idx.keyValue
	idxElements, _ := fs.FindData(mk)
	return idxElements.(*indexElement).value.(*types.BlockState)
}

func (m *forkMultiIndex) FindByPrev(prev common.BlockIdType) *types.BlockState {
	key := computePrevKey(prev.Bytes())
	idx := m.Indexs["byPrev"]
	idxEle, _ := idx.keyValue.FindData(key)
	mainKey := idxEle.(*indexElement).value
	idxe := m.GetIndex("byBlockId")
	objEle, _ := idxe.keyValue.FindData(mainKey.([]byte))
	bs := objEle.(*indexElement).value.(*types.BlockState)
	return bs
}

/*func (m *forkMultiIndex) FindByBlockNum(prev uint32,inCurrentChain bool) *types.BlockState {
	key := computePrevKey(prev.Bytes())
	idx := m.Indexs["byLibBlockNum"]
	idxEle, _ := idx.keyValue.FindData(key)
	mainKey := idxEle.(*indexElement).value
	idxe := m.GetIndex("byBlockId")
	objEle, _ := idxe.keyValue.FindData(mainKey.([]byte))
	bs := objEle.(*indexElement).value.(*types.BlockState)
	return bs
}*/

/*func (m *forkMultiIndex) FindByLibBlockNum(prev uint32,inCurrentChain bool) *types.BlockState {
	key := computePrevKey(prev.Bytes())
	idx := m.Indexs["byLibBlockNum"]
	idxEle, _ := idx.keyValue.FindData(key)
	mainKey := idxEle.(*indexElement).value
	idxe := m.GetIndex("byBlockId")
	objEle, _ := idxe.keyValue.FindData(mainKey.([]byte))
	bs := objEle.(*indexElement).value.(*types.BlockState)
	return bs
}*/

func (idx *fkIndex) eraseMainKey(id *common.BlockIdType) bool {
	keyArray := make([][]byte, 4)
	key := computeMainKey(id.Bytes())
	keyArray = append(keyArray, key)
	ele, sub := idx.keyValue.FindData(key)
	if sub >= 0 {
		block := ele.(*types.BlockState)
		prevKey := computePrevKey(block.SignedBlock.Previous.Bytes())
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
}

func (idx *fkIndex) EraseKey(key []byte) bool {
	idPtr, _ := idx.keyValue.FindData(key)
	m := idPtr.(*indexElement).value
	id := crypto.NewSha256Byte(m.([]byte))
	return idx.eraseMainKey(id)
}

func computeMainKey(val []byte) []byte {
	return append([]byte("byBlockId_"), val...)
}

func (m *forkMultiIndex) modify() {

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
	idxEle := indexElement{[]byte(uniqueKey), []byte(uniqueKey)}
	keySet.Insert(idxEle)
	idx.ManyKeyMap[string(first)] = keySet
}*/

/*type fkIterator interface {
	Next() bool

	Prev() bool

	Key() []byte

	Value() []byte
}*/
