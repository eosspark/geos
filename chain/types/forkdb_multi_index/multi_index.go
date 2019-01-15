package forkdb_multi_index

import (
	"bytes"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type BlockStatePtr = *types.BlockState

//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/"
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container/..."
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/hashed_index/..."
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/ordered_index/..."

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container" MultiIndex(ByBlockId,ByBlockIdNode,BlockStatePtr)
func (m *MultiIndex) GetByBlockId() *ByBlockId         { return m.super }
func (m *MultiIndex) GetByPrev() *ByPrev               { return m.super.super }
func (m *MultiIndex) GetByBlockNum() *ByBlockNum       { return m.super.super.super }
func (m *MultiIndex) GetByLibBlockNum() *ByLibBlockNum { return m.super.super.super.super }

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/multiindex/hashed_index" ByBlockId(MultiIndex,MultiIndexNode,ByPrev,ByPrevNode,BlockStatePtr,common.BlockIdType,ByBlockIdFunc)
var ByBlockIdFunc = func(n BlockStatePtr) common.BlockIdType { return n.BlockId }

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/multiindex/ordered_index" ByPrev(MultiIndex,MultiIndexNode,ByBlockNum,ByBlockNumNode,BlockStatePtr,common.BlockIdType,ByPrevFunc,ByPrevCompare,true)
var ByPrevFunc = func(n BlockStatePtr) common.BlockIdType { return n.Header.Previous }
var ByPrevCompare = func(a, b common.BlockIdType) int { return bytes.Compare(a.Bytes(), b.Bytes()) }

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/multiindex/ordered_index" ByBlockNum(MultiIndex,MultiIndexNode,ByLibBlockNum,ByLibBlockNumNode,BlockStatePtr,ByBlockNumComposite,ByBlockNumFunc,ByBlockNumCompare,true)
type ByBlockNumComposite struct {
	BlockNum       *uint32
	InCurrentChain *bool
}

var ByBlockNumFunc = func(n BlockStatePtr) ByBlockNumComposite { return ByBlockNumComposite{&n.BlockNum, &n.InCurrentChain} }
var ByBlockNumCompare = func(aBlock, bBlock ByBlockNumComposite) int {
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

//go:generate gotemplate -outfmt "gen_%v" "github.com/eosspark/eos-go/common/container/multiindex/ordered_index" ByLibBlockNum(MultiIndex,MultiIndexNode,MultiIndexBase,MultiIndexBaseNode,BlockStatePtr,ByLibBlockNumComposite,ByLibBlockNumFunc,ByLibBlockNumCompare,true)
//go:generate go build
type ByLibBlockNumComposite struct {
	DposIrreversibleBlocknum *uint32
	BftIrreversibleBlocknum  *uint32
	BlockNum                 *uint32
}

var ByLibBlockNumFunc = func(n BlockStatePtr) ByLibBlockNumComposite { return ByLibBlockNumComposite{&n.DposIrreversibleBlocknum, &n.BftIrreversibleBlocknum, &n.BlockNum} }
var ByLibBlockNumCompare = func(aBlock, bBlock ByLibBlockNumComposite) int {
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

