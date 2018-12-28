package fork_multi_index

import (
	"bytes"
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
	increaseKey   IndexKey
	base          IndexBase
	ByBlockId     map[common.BlockIdType]IndexKey
	ByPrev        byPrevIndex
	ByBlockNum    byBlockNumIndex
	ByLibBlockNum byLibBlockNumIndex
}

type IndexBase = map[IndexKey]*Node

type IndexKey = uint32

type Node struct {
	value                 *types.BlockState
	hashByBlockId         common.BlockIdType
	iteratorByPrev        iteratorByPrevIndex
	iteratorByBlockNum    iteratorByBlockNumIndex
	iteratorByLibBlockNum iteratorByLibBlockNumIndex
}

func New() *MultiIndex {
	m := MultiIndex{}
	m.base = make(IndexBase)
	m.ByBlockId = make(map[common.BlockIdType]IndexKey)
	m.ByPrev = *newMultiByPrevIndex()
	m.ByBlockNum = *newMultiByBlockNumIndex()
	m.ByLibBlockNum = *newMultiByLibBlockNumIndex()
	return &m
}

func (m *MultiIndex) Value(k IndexKey) *types.BlockState {
	if node, existing := m.base[k]; existing {
		return node.value
	}
	return nil
}

func (m *MultiIndex) Erase(k IndexKey) bool {
	if node, existing := m.base[k]; existing {
		delete(m.ByBlockId, node.hashByBlockId)
		node.iteratorByPrev.Delete()
		node.iteratorByBlockNum.Delete()
		node.iteratorByLibBlockNum.Delete()
		delete(m.base, k)
if len(m.base) != m.ByBlockNum.Size() { println("Erase Failed")}
		return true
	}
	return false
}

func (m *MultiIndex) Modify(k IndexKey, modify func(b *types.BlockState)) bool {
	if node, existing := m.base[k]; existing {
		node.iteratorByPrev.Delete()
		node.iteratorByBlockNum.Delete()
		node.iteratorByLibBlockNum.Delete()
		bsp := node.value
		modify(bsp)
		node.iteratorByPrev = m.ByPrev.Insert(bsp.Header.Previous, k)
		node.iteratorByBlockNum = m.ByBlockNum.Insert(ByBlockNumComposite{&bsp.BlockNum, &bsp.InCurrentChain}, k)
		node.iteratorByLibBlockNum = m.ByLibBlockNum.Insert(ByLibBlockNumComposite{
			&bsp.DposIrreversibleBlocknum, &bsp.BftIrreversibleBlocknum, &bsp.BlockNum}, k)
	}
	return false
}

func (m *MultiIndex) Insert(n *types.BlockState) bool {
	m.increaseKey ++
	return m.insert(n, m.increaseKey)
}

func (m *MultiIndex) insert(n *types.BlockState, key IndexKey) bool {
	if _, ok := m.ByBlockId[n.BlockId]; ok {
		return false
	}
	m.ByBlockId[n.BlockId] = m.increaseKey
	iteratorByPrev := m.ByPrev.Insert(n.Header.Previous, m.increaseKey)
	iteratorByBlockNum := m.ByBlockNum.Insert(ByBlockNumComposite{&n.BlockNum, &n.InCurrentChain}, m.increaseKey)
	iteratorByLibBlockNum := m.ByLibBlockNum.Insert(ByLibBlockNumComposite{&n.DposIrreversibleBlocknum,
		&n.BftIrreversibleBlocknum, &n.BlockNum}, m.increaseKey)

	m.base[m.increaseKey] = &Node{n, n.BlockId, iteratorByPrev, iteratorByBlockNum, iteratorByLibBlockNum}
	m.increaseKey++
	return true
}

func (m *MultiIndex) Find(id common.BlockIdType) (*types.BlockState, bool) {
	n, found := m.ByBlockId[id]
	return m.base[n].value, found
}

func (m *MultiIndex) Size() int {
	return len(m.ByBlockId)
}

func (m *MultiIndex) Clear() {
	m.base = make(IndexBase)
	m.ByBlockId = make(map[common.BlockIdType]IndexKey)
	m.ByPrev.Clear()
	m.ByBlockNum.Clear()
	m.ByLibBlockNum.Clear()
	m.increaseKey = 0
}
