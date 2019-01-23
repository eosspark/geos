package mock_block

import "fmt"

type MockBlock struct {
	Id        int
	Perv      int
	Num       int
	Dpos      int
	Bft       int
	InCurrent bool
}

func (m *MockBlock) String() string {
	if m.InCurrent {
		return fmt.Sprintf("%d%d%d%d%dT", m.Id, m.Perv, m.Num, m.Dpos, m.Bft)
	}
	return fmt.Sprintf("%d%d%d%d%dF", m.Id, m.Perv, m.Num, m.Dpos, m.Bft)
}

type ValueType = *MockBlock

//go:generate go install "github.com/eosspark/eos-go/common/container/"
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/"
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container/..."
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/hashed_index/..."
//go:generate go install "github.com/eosspark/eos-go/common/container/multiindex/ordered_index/..."

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/multi_index_container" TestIndex(ById,ByIdNode,ValueType)

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/hashed_index" ById(TestIndex,TestIndexNode,ByPrev,ByPrevNode,ValueType,int,ByIdHashFunc)
var ByIdHashFunc = func(n ValueType) int { return n.Id }

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/ordered_index" ByPrev(TestIndex,TestIndexNode,ByNum,ByNumNode,ValueType,int,ByPrevKeyFunc,ByPrevCompare,true)
var ByPrevKeyFunc = func(n ValueType) int { return n.Perv }
var ByPrevCompare = func(av, bv int) int {
	switch {
	case av > bv:
		return 1
	case av < bv:
		return -1
	default:
		return 0
	}
}

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/ordered_index" ByNum(TestIndex,TestIndexNode,ByLibNum,ByLibNumNode,ValueType,ByNumComposite,ByNumKeyFunc,ByNumCompare,true)
type ByNumComposite struct {
	Num       *int
	InCurrent *bool
}

var ByNumKeyFunc = func(n ValueType) ByNumComposite { return ByNumComposite{&n.Num, &n.InCurrent} }
var ByNumCompare = func(aKey, bKey ByNumComposite) int {
	if aKey.Num != nil && bKey.Num != nil {
		if r := ByPrevCompare(*aKey.Num, *bKey.Num); r != 0 {
			return r
		}
	}
	if aKey.InCurrent != nil && bKey.InCurrent != nil {
		if *aKey.InCurrent && !*bKey.InCurrent {
			return -1
		} else if !*aKey.InCurrent && *bKey.InCurrent {
			return 1
		}
	}
	return 0
}

//go:generate gotemplate "github.com/eosspark/eos-go/common/container/multiindex/ordered_index" ByLibNum(TestIndex,TestIndexNode,TestIndexBase,TestIndexBaseNode,ValueType,ByLibNumComposite,ByLibNumKeyFunc,ByLibNumCompare,true)
type ByLibNumComposite struct {
	Dpos *int
	Bft  *int
	Num  *int
}

var ByLibNumKeyFunc = func(n ValueType) ByLibNumComposite { return ByLibNumComposite{&n.Dpos, &n.Bft, &n.Num} }
var ByLibNumCompare = func(aKey, bKey ByLibNumComposite) int {
	if aKey.Dpos != nil && bKey.Dpos != nil {
		if r := ByPrevCompare(*aKey.Dpos, *bKey.Dpos); r != 0 {
			return (-1) * r
		}
	}
	if aKey.Bft != nil && bKey.Bft != nil {
		if r := ByPrevCompare(*aKey.Bft, *bKey.Bft); r != 0 {
			return (-1) * r
		}
	}
	if aKey.Num != nil && bKey.Num != nil {
		if r := ByPrevCompare(*aKey.Num, *bKey.Num); r != 0 {
			return (-1) * r
		}
	}

	return 0
}

//go:generate go build

func (m *TestIndex) GetById() *ById         { return m.super }
func (m *TestIndex) GetByPrev() *ByPrev     { return m.super.super }
func (m *TestIndex) GetByNum() *ByNum       { return m.super.super.super }
func (m *TestIndex) GetByLibNum() *ByLibNum { return m.super.super.super.super }
