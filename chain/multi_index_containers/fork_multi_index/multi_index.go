package fork_multi_index

import (
	"bytes"
	"container/list"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

//go:generate go install "github.com/eosspark/container/templates/..."
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/container/templates/treemap" byPrevIndex(common.BlockIdType,IndexKey,byPrevCompare)
var byPrevCompare = func(a, b interface{}) int {
	return bytes.Compare(a.(common.BlockIdType).Bytes(), b.(common.BlockIdType).Bytes())
}

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/container/templates/treemap" byBlockNumIndex(ByBlockNumComposite,IndexKey,byBlockNumCompare)
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
//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/container/templates/treemap" byLibBlockNumIndex(ByLibBlockNumComposite,IndexKey,byLibBlockNumCompare)
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

type MultiIndex struct {
	base          IndexBase
	ByBlockId     map[common.BlockIdType]IndexKey
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

func New() *MultiIndex {
	m := MultiIndex{}
	m.base = list.New()
	m.ByBlockId = make(map[common.BlockIdType]IndexKey)
	m.ByPrev = *newMultiByPrevIndex()
	m.ByBlockNum = *newMultiByBlockNumIndex()
	m.ByLibBlockNum = *newMultiByLibBlockNumIndex()
	return &m
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
}

func (m *MultiIndex) Modify(itr IndexKey, modifier func(b *types.BlockState)) bool {
	node := itr.Value.(*Node)
	delete(m.ByBlockId, node.hashByBlockId)
	node.iteratorByPrev.Delete()
	node.iteratorByBlockNum.Delete()
	node.iteratorByLibBlockNum.Delete()

	modifier(node.value)

	return m.insert(node.value, itr)
}

func (m *MultiIndex) Insert(n *types.BlockState) bool {
	itr := m.base.PushBack(&Node{value: n})
	return m.insert(n, itr)
}

func (m *MultiIndex) insert(n *types.BlockState, itr IndexKey) bool {
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
	m.ByBlockId = make(map[common.BlockIdType]IndexKey)
	m.ByPrev.Clear()
	m.ByBlockNum.Clear()
	m.ByLibBlockNum.Clear()
}
