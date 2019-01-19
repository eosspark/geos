package int

import "github.com/eosspark/eos-go/common/container"

//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container/..."
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/hashed_index/..."
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/ordered_index/..."

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container" TestIndex(ById,ByIdNode,int)

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/hashed_index" ById(TestIndex,TestIndexNode,ByNum,ByNumNode,int,int,ByIdHashFunc)
var ByIdHashFunc = func(n int) int { return n }
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/ordered_index" ByNum(TestIndex,TestIndexNode,ByPrev,ByPrevNode,int,int,ByNumKeyFunc,ByNumCompare,true)
var ByNumKeyFunc = func(n int) int { return n + 1 }
var ByNumCompare = container.IntComparator
//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/hashed_index" ByPrev(TestIndex,TestIndexNode,TestIndexBase,TestIndexBaseNode,int,int,ByPrevHashFunc)
var ByPrevHashFunc = func(n int) int { return n + 2 }
//go:generate go build

func (m *TestIndex) GetById() *ById     { return m.super }
func (m *TestIndex) GetByNum() *ByNum   { return m.super.super }
func (m *TestIndex) GetByPrev() *ByPrev { return m.super.super.super }
