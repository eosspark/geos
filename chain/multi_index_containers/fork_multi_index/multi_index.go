package fork_multi_index

import (
	"bytes"
	"container/list"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type MultiIndex struct {
	base          IndexBase
	ByBlockId     byBlockIdIndex
	ByPrev        byPrevIndex
	ByBlockNum    byBlockNumIndex
	ByLibBlockNum byLibBlockNumIndex
}

type IndexBase = *list.List
type IndexKey = *list.Element

type Node struct {
	value                 *types.BlockState
	hashByBlockId         common.BlockIdType
	iteratorByPrev        iteratorByPrevIndex
	iteratorByBlockNum    iteratorByBlockNumIndex
	iteratorByLibBlockNum iteratorByLibBlockNumIndex
}
type byBlockIdIndex = map[common.BlockIdType]IndexKey

//go:generate go install "github.com/eosspark/eos-go/common/container/..."
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/treemap" byPrevIndex(common.BlockIdType,IndexKey,byPrevCompare,true)
var byPrevCompare = func(a, b interface{}) int {
	return bytes.Compare(a.(common.BlockIdType).Bytes(), b.(common.BlockIdType).Bytes())
}

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/treemap" byBlockNumIndex(ByBlockNumComposite,IndexKey,byBlockNumCompare,true)
type ByBlockNumComposite struct {
	BlockNum       *uint32
	InCurrentChain *bool
}

var byBlockNumCompare = func(a, b interface{}) int {
	aBlock, bBlock := a.(ByBlockNumComposite), b.(ByBlockNumComposite)

	if aBlock.BlockNum != nil && bBlock.BlockNum != nil {
		if *aBlock.BlockNum < *bBlock.BlockNum {
			return -1
		} else if *aBlock.BlockNum > *bBlock.BlockNum {
			return 1
		}
	}

	if aBlock.InCurrentChain != nil && bBlock.InCurrentChain != nil {
		if *aBlock.InCurrentChain && !*bBlock.InCurrentChain {
			return -1
		} else if !*aBlock.InCurrentChain && *bBlock.InCurrentChain {
			return 1
		}
	}

	return 0
}

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/treemap" byLibBlockNumIndex(ByLibBlockNumComposite,IndexKey,byLibBlockNumCompare,true)
type ByLibBlockNumComposite struct {
	DposIrreversibleBlocknum *uint32
	BftIrreversibleBlocknum  *uint32
	BlockNum                 *uint32
}

var byLibBlockNumCompare = func(a, b interface{}) int {
	aBlock, bBlock := a.(ByLibBlockNumComposite), b.(ByLibBlockNumComposite)
	if aBlock.DposIrreversibleBlocknum != nil && bBlock.DposIrreversibleBlocknum != nil {
		if *aBlock.DposIrreversibleBlocknum > *bBlock.DposIrreversibleBlocknum {
			return -1
		} else if *aBlock.DposIrreversibleBlocknum < *bBlock.DposIrreversibleBlocknum {
			return 1
		}
	}

	if aBlock.BftIrreversibleBlocknum != nil && bBlock.BftIrreversibleBlocknum != nil {
		if *aBlock.BftIrreversibleBlocknum > *bBlock.BftIrreversibleBlocknum {
			return -1
		} else if *aBlock.BftIrreversibleBlocknum < *bBlock.BftIrreversibleBlocknum {
			return 1
		}
	}

	if aBlock.BlockNum != nil && bBlock.BlockNum != nil {
		if *aBlock.BlockNum > *bBlock.BlockNum {
			return -1
		} else if *aBlock.BlockNum < *bBlock.BlockNum {
			return 1
		}
	}

	return 0
}

func New() *MultiIndex {
	return &MultiIndex{
		base:          list.New(),
		ByBlockId:     byBlockIdIndex{},
		ByPrev:        *newByPrevIndex(),
		ByBlockNum:    *newByBlockNumIndex(),
		ByLibBlockNum: *newByLibBlockNumIndex(),
	}
}

func (m *MultiIndex) Value(k IndexKey) *types.BlockState {
	return k.Value.(*Node).value
}

func (m *MultiIndex) Erase(itr IndexKey) {
	node := itr.Value.(*Node)
	delete(m.ByBlockId, node.hashByBlockId)
	node.iteratorByPrev.Delete()
	node.iteratorByBlockNum.Delete()
	node.iteratorByLibBlockNum.Delete()
	m.base.Remove(itr)
}

func (m *MultiIndex) Modify(itr IndexKey, modifier func(b *types.BlockState)) bool {
	node := itr.Value.(*Node)
	modifier(node.value)
	v := node.value

	if v.BlockId != node.hashByBlockId {
		delete(m.ByBlockId, node.hashByBlockId)
		if _, ok := m.ByBlockId[v.BlockId]; ok {
			m.base.Remove(itr)
			return false
		}
		node.hashByBlockId = v.BlockId
	}

	node.iteratorByPrev = node.iteratorByPrev.Modify(v.Header.Previous, itr)
	if node.iteratorByPrev.IsEnd() {
		delete(m.ByBlockId, node.hashByBlockId)
		m.base.Remove(itr)
		return false
	}

	node.iteratorByBlockNum = node.iteratorByBlockNum.Modify(ByBlockNumComposite{&v.BlockNum, &v.InCurrentChain}, itr)
	if node.iteratorByBlockNum.IsEnd() {
		delete(m.ByBlockId, node.hashByBlockId)
		node.iteratorByPrev.Delete()
		m.base.Remove(itr)
		return false
	}

	node.iteratorByLibBlockNum = node.iteratorByLibBlockNum.Modify(ByLibBlockNumComposite{
		&v.DposIrreversibleBlocknum,
		&v.BftIrreversibleBlocknum,
		&v.BlockNum}, itr)
	if node.iteratorByLibBlockNum.IsEnd() {
		delete(m.ByBlockId, node.hashByBlockId)
		node.iteratorByPrev.Delete()
		node.iteratorByBlockNum.Delete()
		m.base.Remove(itr)
		return false
	}

	return true
}

func (m *MultiIndex) Insert(n *types.BlockState) bool {
	itr := m.base.PushBack(&Node{value: n})
	node := itr.Value.(*Node)

	if _, ok := m.ByBlockId[n.BlockId]; ok {
		m.base.Remove(itr)
		return false
	}
	m.ByBlockId[n.BlockId] = itr
	node.hashByBlockId = n.BlockId

	node.iteratorByPrev = m.ByPrev.Insert(n.Header.Previous, itr)
	if node.iteratorByPrev.IsEnd() {
		delete(m.ByBlockId, n.BlockId)
		m.base.Remove(itr)
		return false
	}

	node.iteratorByBlockNum = m.ByBlockNum.Insert(ByBlockNumComposite{&n.BlockNum, &n.InCurrentChain}, itr)
	if node.iteratorByBlockNum.IsEnd() {
		node.iteratorByPrev.Delete()
		delete(m.ByBlockId, n.BlockId)
		m.base.Remove(itr)
		return false
	}

	node.iteratorByLibBlockNum = m.ByLibBlockNum.Insert(ByLibBlockNumComposite{&n.DposIrreversibleBlocknum,
		&n.BftIrreversibleBlocknum, &n.BlockNum}, itr)
	if node.iteratorByLibBlockNum.IsEnd() {
		node.iteratorByBlockNum.Delete()
		node.iteratorByPrev.Delete()
		delete(m.ByBlockId, n.BlockId)
		m.base.Remove(itr)
		return false
	}

	return true
}

func (m *MultiIndex) Find(id common.BlockIdType) (*types.BlockState, bool) {
	if itr, found := m.ByBlockId[id]; found {
		return m.Value(itr), found
	}

	return nil, false
}

func (m *MultiIndex) Size() int {
	return m.base.Len()
}

func (m *MultiIndex) Clear() {
	m.base.Init()
	m.ByBlockId = byBlockIdIndex{}
	m.ByPrev.Clear()
	m.ByBlockNum.Clear()
	m.ByLibBlockNum.Clear()
}
