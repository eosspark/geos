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
	value      Bucket
}

type iteratorNet struct {
	currentSub int
	idx        *indexNet
	value      ElementObject
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
		bt := Bucket{}
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
		bt := Bucket{}
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
		bt := Bucket{}
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
		bt := Bucket{}
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
		bt := Bucket{}
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
		bt := Bucket{}
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
		bt := Bucket{}
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
		bt := Bucket{}
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

func (m *multiIndexNet) erase(i ElementObject) {
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

func (idx *indexNet) upperBound(eo ElementObject) *iteratorNet {
	itr := iteratorNet{}
	itr.idx = idx
	obj, sub := idx.value.UpperBound(eo)
	if sub >= 0 {
		itr.value = obj
		itr.currentSub = sub
	}
	return &itr
}

func (idx *indexNet) lowerBound(eo ElementObject) *iteratorNet {
	itr := iteratorNet{}
	itr.idx = idx
	obj, sub := idx.value.LowerBound(eo)
	if sub >= 0 {
		itr.value = obj
		itr.currentSub = sub
	}
	return &itr
}

func (m *multiIndexNet) modify(old ElementObject, isKey bool, updata func(in ElementObject)) {
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

func (m *multiIndexNet) insert(eo ElementObject) /*(bool,error)*/ {
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


type ElementObject interface {
	ElementObject() //default implements interface only taget
}

type Bucket struct {
	Data    []ElementObject
	Compare func(first ElementObject, second ElementObject) int
}

func (b *Bucket) Len() int {
	return len(b.Data)
}

func (b *Bucket) GetData(i int) (ElementObject, error) {
	if len(b.Data)-1 >= i {
		return b.Data[i], nil
	}
	return nil, errors.New("not found data")
}

func (b *Bucket) Clear() {
	if len(b.Data) > 0 {
		b.Data = nil
	}
}

func (b *Bucket) Find(element ElementObject) (bool, int) {
	r := b.searchSub(element)
	if r >= 0 && b.Compare(element, b.Data[r]) == 0 {
		return true, r
	}
	return false, -1
}

func (b *Bucket) Eraser(element ElementObject) bool {
	result := false
	if b.Len() == 0 {
		return result
	}
	exist, sub := b.Find(element)
	if exist /* && f.Len()>=1*/ {
		b.Data = append(b.Data[:sub], b.Data[sub+1:]...)
		result = true
	}
	return result
}

func (b *Bucket) searchSub(obj ElementObject) int {
	length := b.Len()
	if length == 0 {
		return -1
	}
	i, j := 0, length-1
	for i < j {
		h := int(uint(i+j) >> 1)
		if i <= h && h < j {
			if b.Compare(b.Data[h], obj) == -1 {
				i = h + 1
			} else if b.Compare(b.Data[h], obj) == 0 {
				return h
			} else {
				j = h
			}
		}
	}
	return i
}

func (b *Bucket) Insert(obj ElementObject) (*ElementObject, error) {
	if b.Compare == nil {
		return nil, errors.New("Bucket Compare is nil")
	}
	var result ElementObject
	length := b.Len()
	target := b.Data
	if length == 0 {
		b.Data = append(b.Data, obj)
		result = b.Data[0]
	} else {
		r := b.searchSub(obj)
		start := b.Compare(target[0], obj)
		end := b.Compare(obj, target[length-1])
		if (start == -1 || start == 0) && (end == -1 || end == 0) {
			//Insert middle
			elemnts := []ElementObject{}
			first := target[:r]
			second := target[r:length]
			elemnts = append(elemnts, first...)
			elemnts = append(elemnts, obj)
			elemnts = append(elemnts, second...)
			b.Data = elemnts
			result = elemnts[r]
		} else {
			//insert target before
			if b.Compare(obj, target[0]) == -1 {
				elemnts := []ElementObject{}
				elemnts = append(elemnts, obj)
				elemnts = append(elemnts, target...)
				b.Data = elemnts
				result = elemnts[0]
			} else if b.Compare(obj, target[length-1]) == 1 { //target append
				target = append(target, obj)
				result = target[length]
				b.Data = target
			}
		}
	}
	return &result, nil
}

func (b *Bucket) LowerBound(eo ElementObject) (ElementObject, int) {
	first := 0
	if b.Len() > 0 {
		ext := b.searchSub(eo)
		first = ext
		for i := first; i >= 0; i-- {
			if b.Compare(b.Data[i], eo) == -1 {
				value := b.Data[i+1]
				currentSub := i + 1
				return value, currentSub
			} else if i == 0 && b.Compare(b.Data[i], eo) == 0 {
				value := b.Data[i]
				currentSub := i
				return value, currentSub
			}
		}
	}
	return nil, -1
}

func (b *Bucket) UpperBound(eo ElementObject) (ElementObject, int) {
	if b.Len() > 0 {
		ext := b.searchSub(eo)
		for i := ext; i < b.Len(); i++ {
			if b.Compare(b.Data[i], eo) > 0 {
				value := b.Data[i-1]
				currentSub := i - 1
				return value, currentSub
			} else if i == b.Len()-1 && b.Compare(eo, b.Data[i]) == 0 {
				value := b.Data[i]
				currentSub := i
				return value, currentSub
			}
		}
	}
	return nil, -1
}

