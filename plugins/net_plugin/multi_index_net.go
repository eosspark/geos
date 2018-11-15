package net_plugin

import (
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type multiIndexNet struct {
	indexs     map[string]*indexNet
	objectName string //peer、trx、node
}

type indexNet struct {
	target     string
	uniqueness bool
	less       bool
	value      common.Bucket
}

type iteratorNet struct {
	currentSub int
	idx        *indexNet
	value      common.ElementObject
}

func newNodeMultinetIndex() *multiIndexNet {
	mi := multiIndexNet{}
	mi.objectName = "node"
	mi.indexs = make(map[string]*indexNet)
	index := &indexNet{target: "byId", uniqueness: true}
	index2 := &indexNet{target: "byExpiry", uniqueness: false}
	index3 := &indexNet{target: "byBlockNum", uniqueness: false}
	mi.indexs["byId"] = index
	mi.indexs["byExpiry"] = index2
	mi.indexs["byBlockNum"] = index3
	return &mi
}

func newTrxMultinetIndex() *multiIndexNet {
	mi := multiIndexNet{}
	mi.objectName = "trx"
	mi.indexs = make(map[string]*indexNet)
	index := &indexNet{target: "byId", uniqueness: true}
	index2 := &indexNet{target: "byExpiry", uniqueness: false}
	index3 := &indexNet{target: "byBlockNum", uniqueness: false}
	mi.indexs["byId"] = index
	mi.indexs["byExpiry"] = index2
	mi.indexs["byBlockNum"] = index3
	return &mi
}

func newPeerMultinetIndex() *multiIndexNet {
	mi := multiIndexNet{}
	mi.objectName = "peer"
	mi.indexs = make(map[string]*indexNet)
	index := &indexNet{target: "byId", uniqueness: true}
	index2 := &indexNet{target: "byBlockNum", uniqueness: true}
	mi.indexs["byId"] = index
	mi.indexs["byBlockNum"] = index2
	return &mi
}

func (m *multiIndexNet) getIndex(tag string) *indexNet {
	if index, ok := m.indexs[tag]; ok {
		return index
	}
	return nil
}

func (m *multiIndexNet) insertNode(n *nodeTransactionState) {
	idx := &indexNet{}
	idx.target = "byId"
	idx.uniqueness = true
	if m.indexs[idx.target].value.Len() > 0 {
		m.indexs[idx.target].value.Insert(n)
	} else {
		bt := common.Bucket{}
		bt.Insert(n)
		bt.Compare = CompareById
		idx.value = bt
		m.indexs["byId"] = idx
	}

	expiryIdx := &indexNet{}
	expiryIdx.target = "byExpiry"
	expiryIdx.uniqueness = false
	if m.indexs[idx.target].value.Len() > 0 {
		m.indexs[idx.target].value.Insert(n)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareByExpiry
		bt.Insert(n)
		idx.value = bt
		m.indexs["byExpiry"] = expiryIdx
	}

	numIdx := &indexNet{}
	numIdx.target = "byBlockNum"
	numIdx.uniqueness = false

	if m.indexs[numIdx.target].value.Len() > 0 {
		m.indexs[numIdx.target].value.Insert(n)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareByBlockNum
		bt.Insert(n)
		numIdx.value = bt
		m.indexs["byBlockNum"] = numIdx
	}
}

func (m *multiIndexNet) insertTrx(trx *transactionState) {
	trxIdx := &indexNet{}
	trxIdx.target = "byId"
	trxIdx.uniqueness = true

	if m.indexs[trxIdx.target].value.Len() > 0 {
		m.indexs[trxIdx.target].value.Insert(trx)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareById
		bt.Insert(trx)
		trxIdx.value = bt
		m.indexs["byId"] = trxIdx
	}

	expiryIdx := &indexNet{}
	expiryIdx.target = "byExpiry"
	expiryIdx.uniqueness = false

	if m.indexs[expiryIdx.target].value.Len() > 0 {
		m.indexs[expiryIdx.target].value.Insert(trx)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareByExpiry
		bt.Insert(trx)
		expiryIdx.value = bt
		m.indexs["byExpiry"] = expiryIdx
	}

	numIdx := &indexNet{}
	numIdx.target = "byBlockNum"
	numIdx.uniqueness = false

	if m.indexs[numIdx.target].value.Len() > 0 {
		m.indexs[numIdx.target].value.Insert(trx)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareByBlockNum
		bt.Insert(trx)
		numIdx.value = bt
		m.indexs["byBlockNum"] = numIdx
	}
}

func (m *multiIndexNet) insertPeer(pbs *peerBlockState) {
	peerIdx := &indexNet{}
	peerIdx.target = "byId"
	peerIdx.uniqueness = true

	if m.indexs[peerIdx.target].value.Len() > 0 {
		m.indexs[peerIdx.target].value.Insert(pbs)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareById
		bt.Insert(pbs)
		peerIdx.value = bt
		m.indexs["byId"] = peerIdx
	}

	numIdx := &indexNet{}
	numIdx.target = "byBlockNum"
	numIdx.uniqueness = true

	if m.indexs[numIdx.target].value.Len() > 0 {
		m.indexs[numIdx.target].value.Insert(pbs)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareByBlockNum
		bt.Insert(pbs)
		numIdx.value = bt
		m.indexs["byBlockNum"] = numIdx
	}
}

func (idx *indexNet) findTrxById(id common.BlockIdType) *transactionState {
	trx := transactionState{}
	trx.id = id
	bt := idx.value
	exist, sub := bt.Find(&trx)
	if exist {
		return bt.Data[sub].(*transactionState)
	}
	return nil
}

func (idx *indexNet) findPeerById(id common.BlockIdType) *peerBlockState {
	peer := peerBlockState{}
	peer.id = id
	bt := idx.value
	exist, sub := bt.Find(&peer)
	if exist {
		return bt.Data[sub].(*peerBlockState)
	}
	return nil
}

func (idx *indexNet) findNodeById(id common.BlockIdType) *nodeTransactionState {
	node := nodeTransactionState{}
	node.id = id
	bt := idx.value
	exist, sub := bt.Find(&node)
	if exist {
		return bt.Data[sub].(*nodeTransactionState)
	}
	return nil
}

func (idx *indexNet) findPeerByBlockNum(blockNum uint32) *peerBlockState {
	peer := peerBlockState{}
	peer.blockNum = blockNum
	exist, sub := idx.value.Find(&peer)
	if exist {
		return idx.value.Data[sub].(*peerBlockState)
	}
	return nil
}

func (idx *indexNet) begin() *iteratorNet {
	itr := iteratorNet{}
	if idx.value.Len() > 0 {
		itr.value = idx.value.Data[0]
		itr.currentSub = 0
	}
	return &itr
}

func (m *multiIndexNet) erase(i common.ElementObject) {
	if len(m.indexs) > 0 {
		for _, v := range m.indexs {
			bt := v.value
			ext, _ := bt.Find(i)
			if ext {
				v.value.Eraser(i)
			}
		}
	}
}

func (idx *indexNet) upperBound(eo common.ElementObject) *iteratorNet {
	itr := iteratorNet{}
	itr.idx = idx
	if idx.value.Len() > 0 {
		ext := idx.searchSub(eo)
		if idx.less {
			for i := ext; i < idx.value.Len(); i++ {
				if idx.value.Compare(idx.value.Data[i], eo) > 0 {
					itr.value = idx.value.Data[i-1].(*types.BlockState)
					itr.currentSub = i - 1
					break
				} else if i == idx.value.Len()-1 && idx.value.Compare(eo, idx.value.Data[i]) == 0 {
					itr.value = idx.value.Data[i].(*types.BlockState)
					itr.currentSub = i
				}
			}
		}
		return &itr
	}
	return nil
}

func (idx *indexNet) searchSub(eo common.ElementObject) int {
	length := idx.value.Len()
	i, j := 0, length-1
	for i < j {
		h := int(uint(i+j) >> 1)
		if i <= h && h < j {
			ext := idx.value.Compare(idx.value.Data[h], eo)
			if ext < 0 {
				i = h + 1
			} else {
				j = h
			}
		}
	}
	return i
}

func (idx *indexNet) lowerBound(eo common.ElementObject) *iteratorNet {
	itr := iteratorNet{}
	itr.idx = idx
	first := 0
	if idx.value.Len() > 0 {
		ext := idx.searchSub(eo)
		first = ext
		if idx.less {
			fmt.Println("less search")
			for i := first; i > 0; i-- {
				if idx.value.Compare(idx.value.Data[i], eo) == -1 {
					itr.value = idx.value.Data[i+1].(*types.BlockState)
					itr.currentSub = i + 1
					break
				} else if i == 0 && idx.value.Compare(idx.value.Data[i], eo) == 0 {
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

func (m *multiIndexNet) modify(eo common.ElementObject) {
	m.erase(eo)

	m.insert(eo)
}

func (m *multiIndexNet) insert(eo common.ElementObject) /*(bool,error)*/ {
	switch m.objectName {
	case "node":
		m.insertNode(eo.(*nodeTransactionState))
	case "trx":
		m.insertTrx(eo.(*transactionState))
	case "peer":
		m.insertPeer(eo.(*peerBlockState))
	}
}

func getInstance(objTag string) (*multiIndexNet, error) {
	var m *multiIndexNet
	switch objTag {
	case "node":
		m = newNodeMultinetIndex()
	case "peer":
		m = newPeerMultinetIndex()
	case "trx":
		m = newTrxMultinetIndex()
	}
	if m == nil {
		return nil, errors.New("multiIndexNet getInstance is error,objTag must be [node、peer、trx]")
	}
	return m, nil
}

func (idx *multiIndexNet) clear() bool {
	idx.indexs = nil
	return true
}
