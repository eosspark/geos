package net_plugin

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type multiIndexNet struct {
	indexs map[string]*indexNet
}

type indexNet struct {
	target     string
	uniqueness bool
	less       bool
	Value      common.Bucket
}

type iteratorNet struct {
	keySet common.Bucket
}

func newNodeMultinetIndex() *multiIndexNet {
	mi := multiIndexNet{}
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
	mi.indexs = make(map[string]*indexNet)
	index := &indexNet{target: "byId", uniqueness: true}
	index2 := &indexNet{target: "byBlockNum", uniqueness: true}
	mi.indexs["byId"] = index
	mi.indexs["byBlockNum"] = index2
	return &mi
}

func (m *multiIndexNet) GetIndex(tag string) *indexNet {
	if index, ok := m.indexs[tag]; ok {
		return index
	}
	return nil
}

func (m *multiIndexNet) insertNode(n *nodeTransactionState) {
	idx := &indexNet{}
	idx.target = "byId"
	idx.uniqueness = true
	if m.indexs[idx.target].Value.Len() > 0 {
		m.indexs[idx.target].Value.Insert(n)
	} else {
		bt := common.Bucket{}
		bt.Insert(n)
		bt.Compare = CompareById
		idx.Value = bt
		m.indexs["byId"] = idx
	}

	expiryIdx := &indexNet{}
	expiryIdx.target = "byExpiry"
	expiryIdx.uniqueness = false
	if m.indexs[idx.target].Value.Len() > 0 {
		m.indexs[idx.target].Value.Insert(n)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareByExpiry
		bt.Insert(n)
		idx.Value = bt
		m.indexs["byExpiry"] = expiryIdx
	}

	numIdx := &indexNet{}
	numIdx.target = "byBlockNum"
	numIdx.uniqueness = false

	if m.indexs[numIdx.target].Value.Len() > 0 {
		m.indexs[numIdx.target].Value.Insert(n)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareByBlockNum
		bt.Insert(n)
		numIdx.Value = bt
		m.indexs["byBlockNum"] = numIdx
	}
}

func (m *multiIndexNet) insertTrx(trx *transactionState) {
	trxIdx := &indexNet{}
	trxIdx.target = "byId"
	trxIdx.uniqueness = true

	if m.indexs[trxIdx.target].Value.Len() > 0 {
		m.indexs[trxIdx.target].Value.Insert(trx)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareById
		bt.Insert(trx)
		trxIdx.Value = bt
		m.indexs["byId"] = trxIdx
	}

	expiryIdx := &indexNet{}
	expiryIdx.target = "byExpiry"
	expiryIdx.uniqueness = false

	if m.indexs[expiryIdx.target].Value.Len() > 0 {
		m.indexs[expiryIdx.target].Value.Insert(trx)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareByExpiry
		bt.Insert(trx)
		expiryIdx.Value = bt
		m.indexs["byExpiry"] = expiryIdx
	}

	numIdx := &indexNet{}
	numIdx.target = "byBlockNum"
	numIdx.uniqueness = false

	if m.indexs[numIdx.target].Value.Len() > 0 {
		m.indexs[numIdx.target].Value.Insert(trx)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareByBlockNum
		bt.Insert(trx)
		numIdx.Value = bt
		m.indexs["byBlockNum"] = numIdx
	}
}

func (m *multiIndexNet) insertPeer(pbs *peerBlockState) {
	peerIdx := &indexNet{}
	peerIdx.target = "byId"
	peerIdx.uniqueness = true

	if m.indexs[peerIdx.target].Value.Len() > 0 {
		m.indexs[peerIdx.target].Value.Insert(pbs)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareById
		bt.Insert(pbs)
		peerIdx.Value = bt
		m.indexs["byId"] = peerIdx
	}

	numIdx := &indexNet{}
	numIdx.target = "byBlockNum"
	numIdx.uniqueness = true

	if m.indexs[numIdx.target].Value.Len() > 0 {
		m.indexs[numIdx.target].Value.Insert(pbs)
	} else {
		bt := common.Bucket{}
		bt.Compare = CompareByBlockNum
		bt.Insert(pbs)
		numIdx.Value = bt
		m.indexs["byBlockNum"] = numIdx
	}
}

func (idx *indexNet) FindTrxById(id common.BlockIdType) interface{} {
	trx := transactionState{}
	trx.id = id
	bt := idx.Value
	exist, sub := bt.Find(&trx)
	if exist {
		return bt.Data[sub].(*transactionState)
	}
	return nil
}

func (idx *indexNet) FindPeerById(id common.BlockIdType) interface{} {
	peer := peerBlockState{}
	peer.id = id
	bt := idx.Value
	exist, sub := bt.Find(&peer)
	if exist {
		return bt.Data[sub].(*peerBlockState)
	}
	return nil
}

func (idx *indexNet) FindNodeById(id common.BlockIdType) interface{} {
	node := nodeTransactionState{}
	node.id = id
	bt := idx.Value
	exist, sub := bt.Find(&node)
	if exist {
		return bt.Data[sub].(*nodeTransactionState)
	}
	return nil
}

func (idx *indexNet) FindPeerByBlockNum(blockNum uint32) interface{} {
	peer := peerBlockState{}
	peer.blockNum = blockNum
	exist, sub := idx.Value.Find(&peer)
	if exist {
		return idx.Value.Data[sub].(*peerBlockState)
	}
	return nil
}

func (idx *indexNet) Begin() *iteratorNet {
	itr := iteratorNet{}
	if idx.Value.Len() > 0 {
		itr.keySet.Data[0] = idx.Value.Data[0]
	}
	return &itr
}

func (m *multiIndexNet) eraseNode(i *nodeTransactionState) {
	if len(m.indexs) > 0 {
		for _, v := range m.indexs {
			bt := v.Value
			ext, _ := bt.Find(i)
			if ext {
				v.Value.Easer(i)
			}
		}
	}
}

func (m *multiIndexNet) erasePeer(i *peerBlockState) {
	if len(m.indexs) > 0 {
		for _, v := range m.indexs {
			bt := v.Value
			ext, _ := bt.Find(i)
			if ext {
				v.Value.Easer(i)
			}
		}
	}
}

func (m *multiIndexNet) eraseTrx(i *transactionState) {
	if len(m.indexs) > 0 {
		for _, v := range m.indexs {
			bt := v.Value
			ext, _ := bt.Find(i)
			if ext {
				v.Value.Easer(i)
			}
		}
	}
}

func (idx *indexNet) trxUpperBound(trx *transactionState) *iteratorNet {
	itr := iteratorNet{}
	var tagObj *transactionState
	if idx.Value.Len() > 0 {
		for _, idxEle := range idx.Value.Data {
			tagObj = idxEle.(*transactionState)
			if idx.Value.Compare(idxEle.(*transactionState), trx) == 1 {
				itr.keySet.Insert(tagObj)
				break
			}
		}
		return idx.trxLowerBound(tagObj)
	}
	return nil
}

func (idx *indexNet) searchSub(b *types.BlockState) int {
	length := idx.Value.Len()
	i, j := 0, length-1
	for i < j {
		h := int(uint(i+j) >> 1)
		if i <= h && h < j {
			ext := idx.Value.Compare(idx.Value.Data[h], b)
			if ext < 0 {
				i = h + 1
			} else {
				j = h
			}
		}
	}
	return i
}

func (idx *indexNet) trxLowerBound(b *transactionState) *iteratorNet {
	itr := iteratorNet{}
	first := 0
	if idx.Value.Len() > 0 {
		//start
		length := idx.Value.Len()
		i, j := 0, length-1
		for i < j {
			h := int(uint(i+j) >> 1)
			if i <= h && h < j {
				ext := idx.Value.Compare(idx.Value.Data[h], b)
				if ext < 0 {
					i = h + 1
				} else {
					j = h
				}
			}
		}
		first = i
		for i := first; i < idx.Value.Len(); i++ {
			if idx.Value.Compare(idx.Value.Data[i], b) > 0 || (i == idx.Value.Len()-1 && idx.Value.Compare(idx.Value.Data[i], b) == 0) {
				if i == idx.Value.Len() {
					itr.keySet.Data = idx.Value.Data[first:idx.Value.Len()]
				} else {
					itr.keySet.Data = idx.Value.Data[first : i+1]
				}
				break
			}
		}
		return &itr
	}
	return nil
}

/*func (idx *indexNet) findTrxByID(id common.TransactionIdType) *transactionState{
	idKey :=computeIdKey(id.Bytes())
	nie,_ := idx.Value.FindData(idKey)
	return nie.(*netIndexElement).value.(*transactionState)
}*/

/*func (m *multiIndexNet) updatePeer(pbs peerBlockState) {
	idx := m.indexs["byId"]
	idKey := computeIdKey(pbs.id.BigEndianBytes())
	idxEle, t := idx.Value.Find(idKey)
	param := netIndexElement{idxEle.GetKey(), &pbs}
	idx.Value.Data[t] = param
	fmt.Println("updatePeer result:%#v", idx.Value.Data)
}

func (m *multiIndexNet) updateTrx(trx transactionState) {
	idx := m.indexs["byId"]
	idKey := computeIdKey(trx.id.BigEndianBytes())
	fmt.Println("updateTrx before:%#v", idx.Value.Data)
	idxEle, t := idx.Value.FindData(idKey)
	param := netIndexElement{idxEle.GetKey(), &trx}
	idx.Value.Data[t] = param
	fmt.Println("updateTrx result:%#v", idx.Value.Data)
}

func (m *multiIndexNet) updateNodeTrx(trx nodeTransactionState) {
	idx := m.indexs["byId"]
	idKey := computeIdKey(trx.id.BigEndianBytes())
	fmt.Println("updateNodeTrx before:%#v", idx.Value.Data)
	idxEle, t := idx.Value.FindData(idKey)
	param := netIndexElement{idxEle.GetKey(), &trx}
	idx.Value.Data[t] = param
	fmt.Println("updateNodeTrx result:%#v", idx.Value.Data)
}
*/
//modify key recompute

func (m *multiIndexNet) modifyTrx(trx *transactionState) {
	m.eraseTrx(trx)

	m.insertTrx(trx)
}

func (m *multiIndexNet) modifyNode(node *nodeTransactionState) {
	m.eraseNode(node)

	m.insertNode(node)

}

func (m *multiIndexNet) modifyPeer(pbs *peerBlockState) {
	m.erasePeer(pbs)
	m.insertPeer(pbs)
}

func (idx *multiIndexNet) clear() bool {
	idx.indexs = nil
	return true
}

/*func computeIdKey(val []byte) []byte {
	return append([]byte("byId_"), val...)
}

func computeExpiryKey(val []byte) []byte {
	return append([]byte("byExpiry_"), val...)
}

func computeBlockNumKey(blockNum uint32) []byte {
	bn := make([]byte, 8)
	binary.BigEndian.PutUint64(bn, uint64(blockNum))
	return append([]byte("byBlockNum_"), bn...)
}*/
