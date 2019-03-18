package int

import (
	"github.com/eosspark/eos-go/common/container"
	"github.com/eosspark/eos-go/common/container/allocator/callocator"
)

var alloc = callocator.Instance

//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container/..."
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/hashed_index/..."
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/ordered_index/..."

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container" TestIndex(ById,ByIdNode,int,alloc)

// go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/hashed_index" ById(TestIndex,TestIndexNode,ByNum,ByNumNode,int,int,ByIdHashFunc,alloc)
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/hashed_index" ById(TestIndex,TestIndexNode,TestIndexBase,TestIndexBaseNode,int,int,ByIdKeyFunc,ByIdHashFunc,alloc)
var ByIdKeyFunc = func(n int) int { return n }
var ByIdHashFunc = func(n int) uintptr { return container.Hash(n) }

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/ordered_index" ByNum(TestIndex,TestIndexNode,ByPrev,ByPrevNode,int,int,ByNumKeyFunc,ByNumCompare,true,alloc)
var ByNumKeyFunc = func(n int) int { return n + 1 }
var ByNumCompare = container.IntComparator

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/hashed_index" ByPrev(TestIndex,TestIndexNode,TestIndexBase,TestIndexBaseNode,int,int,ByPrevKeyFunc,ByPrevHashFunc,alloc)
var ByPrevKeyFunc = func(n int) int { return n + 2 }
var ByPrevHashFunc = func(n int) uintptr { return container.Hash(n + 2) }

//go:generate go build

func (m *TestIndex) GetById() *ById { return (*ById)(m.super.Get()) }

//func (m *TestIndex) GetByNum() *ByNum   { return (*ByNum)(m.GetById().super.Get()) }
//func (m *TestIndex) GetByPrev() *ByPrev { return (*ByPrev)(m.GetByNum().super.Get()) }
