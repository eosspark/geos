package net_plugin

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

type netMultiIndex struct {
	indexs map[string]*netIndex
}

type netIndexElement struct {
	key   []byte
	value interface{}
}

func (idx netIndexElement) GetKey() []byte {
	return idx.key
}

type netIndex struct {
	target     string
	uniqueness bool
	less       bool
	keyValue   common.FlatSet
}

type netIterator struct {
	key    []byte
	keySet common.FlatSet
}

func newNodeMultinetIndex() *netMultiIndex {
	mi := netMultiIndex{}
	mi.indexs = make(map[string]*netIndex)
	index := &netIndex{target: "byId", uniqueness: true}
	index2 := &netIndex{target: "byExpiry", uniqueness: false}
	index3 := &netIndex{target: "byBlockNum", uniqueness: false}
	mi.indexs["byId"] = index
	mi.indexs["byExpiry"] = index2
	mi.indexs["byBlockNum"] = index3
	return &mi
}

func newTrxMultinetIndex() *netMultiIndex {
	mi := netMultiIndex{}
	mi.indexs = make(map[string]*netIndex)
	index := &netIndex{target: "byId", uniqueness: true}
	index2 := &netIndex{target: "byExpiry", uniqueness: false}
	index3 := &netIndex{target: "byBlockNum", uniqueness: false}
	mi.indexs["byId"] = index
	mi.indexs["byExpiry"] = index2
	mi.indexs["byBlockNum"] = index3
	return &mi
}

func newPeerMultinetIndex() *netMultiIndex {
	mi := netMultiIndex{}
	mi.indexs = make(map[string]*netIndex)
	index := &netIndex{target: "byId", uniqueness: true}
	index2 := &netIndex{target: "byBlockNum", uniqueness: true}
	mi.indexs["byId"] = index
	mi.indexs["byBlockNum"] = index2
	return &mi
}

func (m *netMultiIndex) GetIndex(tag string) *netIndex {
	if index, ok := m.indexs[tag]; ok {
		return index
	}
	return nil
}

func (m *netMultiIndex) insertNode(n *nodeTransactionState) {
	idx := &netIndex{}
	idx.target = "byId"
	idx.uniqueness = true
	idKey := computeIdKey(n.id.Bytes())
	idIdxElement := netIndexElement{idKey, &n}
	if m.indexs[idx.target].keyValue.Len() > 0 {
		m.indexs[idx.target].keyValue.Insert(&idIdxElement)
	} else {
		fs := common.FlatSet{}
		fs.Insert(&idIdxElement)
		idx.keyValue = fs
		m.indexs["byId"] = idx
	}

	expiryIdx := &netIndex{}
	expiryIdx.target = "byExpiry"
	expiryIdx.uniqueness = false
	exKey := computeExpiryKey(n.id.Bytes())
	exIdxElement := netIndexElement{exKey, &idKey}
	if m.indexs[idx.target].keyValue.Len() > 0 {
		m.indexs[idx.target].keyValue.Insert(&exIdxElement)
	} else {
		fs := common.FlatSet{}
		fs.Insert(&exIdxElement)
		idx.keyValue = fs
		m.indexs["byExpiry"] = expiryIdx
	}

	numIdx := &netIndex{}
	numIdx.target = "byBlockNum"
	numIdx.uniqueness = false
	numKey := computeBlockNumKey(n.blockNum)
	numIdxElement := netIndexElement{numKey, &idKey}
	if m.indexs[numIdx.target].keyValue.Len() > 0 {
		m.indexs[numIdx.target].keyValue.Insert(&numIdxElement)
	} else {
		fs := common.FlatSet{}
		fs.Insert(&numIdxElement)
		numIdx.keyValue = fs
		m.indexs["byBlockNum"] = numIdx
	}
}

func (m *netMultiIndex) insertTrx(trx *transactionState) {
	trxIdx := &netIndex{}
	trxIdx.target = "byId"
	trxIdx.uniqueness = true
	idKey := computeIdKey(trx.id.Bytes())
	idIdxElement := netIndexElement{idKey, &trx}
	if m.indexs[trxIdx.target].keyValue.Len() > 0 {
		m.indexs[trxIdx.target].keyValue.Insert(&idIdxElement)
	} else {
		fs := common.FlatSet{}
		fs.Insert(&idIdxElement)
		trxIdx.keyValue = fs
		m.indexs["byId"] = trxIdx
	}

	expiryIdx := &netIndex{}
	expiryIdx.target = "byExpiry"
	expiryIdx.uniqueness = false
	exKey := computeExpiryKey(trx.id.Bytes())
	exIdxElement := netIndexElement{exKey, &idKey}
	if m.indexs[expiryIdx.target].keyValue.Len() > 0 {
		m.indexs[expiryIdx.target].keyValue.Insert(&exIdxElement)
	} else {
		fs := common.FlatSet{}
		fs.Insert(&exIdxElement)
		expiryIdx.keyValue = fs
		m.indexs["byExpiry"] = expiryIdx
	}

	numIdx := &netIndex{}
	numIdx.target = "byBlockNum"
	numIdx.uniqueness = false
	numKey := computeBlockNumKey(trx.blockNum)
	numIdxElement := netIndexElement{numKey, &idKey}
	if m.indexs[numIdx.target].keyValue.Len() > 0 {
		m.indexs[numIdx.target].keyValue.Insert(&numIdxElement)
	} else {
		fs := common.FlatSet{}
		fs.Insert(&numIdxElement)
		numIdx.keyValue = fs
		m.indexs["byBlockNum"] = numIdx
	}
}

func (m *netMultiIndex) insertPeer(pbs *peerBlockState) {
	peerIdx := &netIndex{}
	peerIdx.target = "byId"
	peerIdx.uniqueness = true
	idKey := computeIdKey(pbs.id.Bytes())
	idIdxElement := netIndexElement{idKey, &pbs}
	if m.indexs[peerIdx.target].keyValue.Len() > 0 {
		m.indexs[peerIdx.target].keyValue.Insert(&idIdxElement)
	} else {
		fs := common.FlatSet{}
		fs.Insert(&idIdxElement)
		peerIdx.keyValue = fs
		m.indexs["byId"] = peerIdx
	}

	numIdx := &netIndex{}
	numIdx.target = "byBlockNum"
	numIdx.uniqueness = true
	numKey := computeBlockNumKey(pbs.blockNum)
	numIdxElement := netIndexElement{numKey, &idKey}
	if m.indexs[numIdx.target].keyValue.Len() > 0 {
		m.indexs[numIdx.target].keyValue.Insert(&numIdxElement)
	} else {
		fs := common.FlatSet{}
		fs.Insert(&numIdxElement)
		numIdx.keyValue = fs
		m.indexs["byBlockNum"] = numIdx
	}
}

func (m *netMultiIndex) FindById(id common.BlockIdType) interface{} {
	mk := computeIdKey(id.Bytes())
	idx := m.indexs["byId"]
	fs := idx.keyValue
	idxElements, _ := fs.FindData(mk)
	return idxElements.(*netIndexElement).value
}

func (m *netMultiIndex) FindPeerByBlockNum(blockNum uint32) interface{} {
	numKey := computeBlockNumKey(blockNum)
	numIdx := m.indexs["byBlockNum"]
	numFs := numIdx.keyValue
	numIdxElements, _ := numFs.FindData(numKey)
	idValue := numIdxElements.(*netIndexElement).value
	idx := m.indexs["byId"]
	fs := idx.keyValue
	idxElements, _ := fs.FindData(idValue.([]byte))
	return idxElements.(*netIndexElement).value
}

/*func (m *netMultiIndex) FindTrxByBlockNum(blockNum uint32) *netIndex{
	numKey := computeBlockNumKey(blockNum)
	numIdx := m.indexs["byBlockNum"]
	numFs := numIdx.keyValue
	numIdxElements,_:=numFs.FindData(numKey)
	idValue := numIdxElements.(*netIndexElement).value

	return nil
}*/

func (idx *netIndex) Begin() *netIterator {
	itr := netIterator{}
	if idx.keyValue.Len() > 0 {
		idxEle := idx.keyValue.Data[0]
		itr.key = idxEle.GetKey()
		itr.keySet = idx.keyValue
	}
	return &itr
}

func (idx *netIndex) UpperBound(b []byte) *netIterator {
	itr := netIterator{}
	if idx.keyValue.Len() > 0 {
		for _, idxEle := range idx.keyValue.Data {
			tagKey := idxEle.(*netIndexElement).GetKey()
			if bytes.Compare(tagKey, b) == 1 {
				itr.key = tagKey
				break
			}
		}
		return idx.LowerBound(itr.key)
	}
	return nil
}

func (idx *netIndex) LowerBound(b []byte) *netIterator {
	itr := netIterator{}
	first := 0
	length := idx.keyValue.Len()
	if length > 0 {
		//start
		i, j := 0, length-1
		for i < j {
			h := int(uint(i+j) >> 1)
			if i <= h && h < j {
				ext := strings.Index(string(idx.keyValue.Data[h].GetKey()), string(b))
				if ext >= 0 {
					first = h
					break
				} else {
					i = h + 1
				}
			}
		}
		if first == -1 {
			return nil
		}
		idxEle := idx.keyValue.Data[first]
		itr.key = idxEle.(*netIndexElement).GetKey()
		//end
		si, sj := 0, length-1
		for i < j {
			h := int(si + sj>>1)
			if si <= h && h < sj {
				et := strings.Index(string(idx.keyValue.Data[h].GetKey()), string(b))
				if et < 0 {
					itr.keySet.Data = idx.keyValue.Data[first:h]
					break
				} else {
					si = h + 1
				}
			}
		}
		if len(itr.key) > 0 && itr.keySet.Len() == 0 {
			itr.keySet.Data = idx.keyValue.Data[first:]
		}
		return &itr
	}
	return nil
}

func (m *netMultiIndex) eraseTrxById(id *common.BlockIdType) bool {
	keyArray := make([][]byte, 4)
	key := computeIdKey(id.Bytes())
	keyArray = append(keyArray, key)
	if len(m.indexs) > 0 {

		for _, idx := range m.indexs {
			ele, sub := idx.keyValue.FindData(key)
			if sub >= 0 {
				block := ele.(*types.BlockState)
				prevKey := computeExpiryKey(block.SignedBlock.Previous.Bytes())
				keyArray = append(keyArray, prevKey)
				blockNumKey := computeBlockNumKey(block.BlockNum)
				keyArray = append(keyArray, blockNumKey)
				for _, k := range keyArray {
					boo := idx.keyValue.Remove(netIndexElement{key: k})
					if !boo {
						log.Error("netMultiIndex eraseTrx is error:%#v", k)
					}
				}
			}
		}
	}
	return true
}

func (m *netMultiIndex) eraseTrx(startKey []byte, endKey []byte) {
	idx := m.indexs["byId"]
	_, start := idx.keyValue.FindData(startKey)
	_, end := idx.keyValue.FindData(endKey)
	if start > end || start < 0 || end < 0 {
		return
	}
	tmpIds := idx.keyValue.Data[start:end]
	for _, idxEle := range tmpIds {
		val := idxEle.(*netIndexElement).value
		id := crypto.NewSha256Byte(val.([]byte))
		m.eraseTrxById(id)
	}
}

/*func (idx *netIndex) findTrxByID(id common.TransactionIdType) *transactionState{
	idKey :=computeIdKey(id.Bytes())
	nie,_ := idx.keyValue.FindData(idKey)
	return nie.(*netIndexElement).value.(*transactionState)
}*/

func (m *netMultiIndex) updatePeer(pbs peerBlockState) {
	idx := m.indexs["byId"]
	idKey := computeIdKey(pbs.id.Bytes())
	fmt.Println("updatePeer before:%#v", idx.keyValue.Data)
	idxEle, t := idx.keyValue.FindData(idKey)
	param := netIndexElement{idxEle.GetKey(), &pbs}
	idx.keyValue.Data[t] = param
	fmt.Println("updatePeer result:%#v", idx.keyValue.Data)
}

func (m *netMultiIndex) updateTrx(trx transactionState) {
	idx := m.indexs["byId"]
	idKey := computeIdKey(trx.id.Bytes())
	fmt.Println("updateTrx before:%#v", idx.keyValue.Data)
	idxEle, t := idx.keyValue.FindData(idKey)
	param := netIndexElement{idxEle.GetKey(), &trx}
	idx.keyValue.Data[t] = param
	fmt.Println("updateTrx result:%#v", idx.keyValue.Data)
}

func (m *netMultiIndex) updateNodeTrx(trx nodeTransactionState) {
	idx := m.indexs["byId"]
	idKey := computeIdKey(trx.id.Bytes())
	fmt.Println("updateNodeTrx before:%#v", idx.keyValue.Data)
	idxEle, t := idx.keyValue.FindData(idKey)
	param := netIndexElement{idxEle.GetKey(), &trx}
	idx.keyValue.Data[t] = param
	fmt.Println("updateNodeTrx result:%#v", idx.keyValue.Data)
}

//modify key recompute

func (m *netMultiIndex) modifyTrx(trx *transactionState) {
	suc := m.eraseTrxById(&trx.id)
	if suc {
		m.insertTrx(trx)
	}
}

func (m *netMultiIndex) modifyNodeTrx(nodeTrx *nodeTransactionState) {
	suc := m.eraseTrxById(&nodeTrx.id)
	if suc {
		m.insertNode(nodeTrx)
	}
}

func (m *netMultiIndex) modifyPeer(pbs *peerBlockState) {
	suc := m.eraseTrxById(&pbs.id)
	if suc {
		m.insertPeer(pbs)
	}
}

func (idx *netMultiIndex) clear() bool {
	idx.indexs = nil
	return true
}

func computeIdKey(val []byte) []byte {
	return append([]byte("byId_"), val...)
}

func computeExpiryKey(val []byte) []byte {
	return append([]byte("byExpiry_"), val...)
}

func computeBlockNumKey(blockNum uint32) []byte {
	bn := make([]byte, 8)
	binary.BigEndian.PutUint64(bn, uint64(blockNum))
	return append([]byte("byBlockNum_"), bn...)
}
