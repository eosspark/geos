package net_plugin

import (
	"errors"
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

func newNodeTransactionIndex() *multiIndexNet {
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

func newTransactionStateIndex() *multiIndexNet {
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

func newPeerBlockStatueIndex() *multiIndexNet {
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

func (m *multiIndexNet) insertNodeTrx(n *nodeTransactionState) {
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

func (m *multiIndexNet) insertPeerBlock(pbs *peerBlockState) {
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

func (idx *indexNet) findLocalTrxById(id common.TransactionIdType) *nodeTransactionState {
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

func (m *multiIndexNet) eraseRegion(begin int, end int, tag string) {
	idx := m.getIndex(tag)
	if idx.value.Len() > 0 {
		tmp := idx.value.Data[begin:end]
		if len(tmp) > 0 {
			for _, t := range tmp {
				for _, v := range m.indexs {
					bt := v.value
					ext, _ := bt.Find(t)
					if ext {
						v.value.Eraser(t)
					}
				}
			}
		}
	}
}

func (idx *indexNet) upperBound(eo common.ElementObject) *iteratorNet {
	itr := iteratorNet{}
	itr.idx = idx
	obj, sub := idx.value.UpperBound(eo)
	if sub >= 0 {
		itr.value = obj
		itr.currentSub = sub
	}
	return &itr
}

func (idx *indexNet) lowerBound(eo common.ElementObject) *iteratorNet {
	itr := iteratorNet{}
	itr.idx = idx
	obj, sub := idx.value.LowerBound(eo)
	if sub >= 0 {
		itr.value = obj
		itr.currentSub = sub
	}
	return &itr
}

func (m *multiIndexNet) modify(old common.ElementObject, isKey bool, updata func(in common.ElementObject)) {
	if isKey {
		m.erase(old)
		updata(old)
		m.insert(old)
	} else {
		for _, v := range m.indexs {
			bt := v.value
			ext, sub := bt.Find(old)
			if ext {
				updata(old)
				v.value.Data[sub] = old
			}
		}
	}
}

func (m *multiIndexNet) insert(eo common.ElementObject) /*(bool,error)*/ {
	switch m.objectName {
	case "node":
		m.insertNodeTrx(eo.(*nodeTransactionState))
	case "trx":
		m.insertTrx(eo.(*transactionState))
	case "peer":
		m.insertPeerBlock(eo.(*peerBlockState))
	}
}

func getInstance(objTag string) (*multiIndexNet, error) {
	var m *multiIndexNet
	switch objTag {
	case "node":
		m = newNodeTransactionIndex()
	case "peer":
		m = newPeerBlockStatueIndex()
	case "trx":
		m = newTransactionStateIndex()
	}
	if m == nil {
		return nil, errors.New("multiIndexNet getInstance is error,objTag must be [node、peer、trx]")
	}
	return m, nil
}

func (itr *iteratorNet) next() bool {
	itr.currentSub++
	if itr.currentSub < itr.idx.value.Len() {
		itr.value = itr.idx.value.Data[itr.currentSub].(*types.BlockState)
		return true
	} else {
		return false
	}
}

func (idx *multiIndexNet) clear() bool {
	idx.indexs = nil
	return true
}
