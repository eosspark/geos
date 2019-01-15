package producer_multi_index

import (
	"container/list"
	"github.com/eosspark/eos-go/common"
)

type TransactionIdWithExpiry struct {
	TrxId  common.TransactionIdType
	Expiry common.TimePoint
}

type MultiIndex struct {
	base     IndexBase
	ById     byIdIndex
	ByExpiry byExpireIndex
}

type IndexBase = *list.List
type IndexKey = *list.Element

type Node struct {
	value            TransactionIdWithExpiry
	hashById         common.TransactionIdType
	iteratorByExpire iteratorByExpireIndex
}

type byIdIndex = map[common.TransactionIdType]IndexKey

//go:generate go install "github.com/eosspark/eos-go/common/container/..."
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/treemap" byExpireIndex(common.TimePoint,IndexKey,byExpireCompare,true)
var byExpireCompare = func(a, b interface{}) int {
	atime, btime := a.(common.TimePoint), b.(common.TimePoint)
	switch {
	case atime > btime:
		return 1
	case atime < btime:
		return -1
	default:
		return 0
	}
}

func New() *MultiIndex {
	return &MultiIndex{
		base:     list.New(),
		ById:     byIdIndex{},
		ByExpiry: *newByExpireIndex(),
	}
}

func (m *MultiIndex) Size() int {
	return m.base.Len()
}

func (m *MultiIndex) Value(k IndexKey) TransactionIdWithExpiry {
	return k.Value.(*Node).value
}

func (m *MultiIndex) Insert(n *TransactionIdWithExpiry) bool {
	itr := m.base.PushBack(&Node{value: *n})
	node := itr.Value.(*Node)

	if _, ok := m.ById[n.TrxId]; ok {
		m.base.Remove(itr)
		return false
	}
	m.ById[n.TrxId] = itr
	node.hashById = n.TrxId

	node.iteratorByExpire = m.ByExpiry.Insert(n.Expiry, itr)
	if node.iteratorByExpire.IsEnd() {
		delete(m.ById, n.TrxId)
		m.base.Remove(itr)
		return false
	}

	return true
}

func (m *MultiIndex) Erase(itr IndexKey) {
	node := itr.Value.(*Node)
	delete(m.ById, node.hashById)
	node.iteratorByExpire.Delete()
	m.base.Remove(itr)
}
