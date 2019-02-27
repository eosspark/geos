//+build !shared

package multi_index_container

import (
	"github.com/eosspark/eos-go/common/container"
	"github.com/eosspark/eos-go/common/container/multiindex"
)

// template type MultiIndex(SuperIndex,SuperNode,Value)
type Value = int

type MultiIndex struct {
	super *SuperIndex
	count int
}

func New() *MultiIndex {
	m := &MultiIndex{}
	m.super = &SuperIndex{}
	m.super.init(m)
	return m
}

/*generic class*/
type SuperIndex struct {
	init    func(*MultiIndex)
	clear   func()
	insert  func(Value, *MultiIndexNode) (*SuperNode, bool)
	erase   func(*SuperNode)
	erase_  func(multiindex.IteratorType)
	modify  func(*SuperNode) (*SuperNode, bool)
	modify_ func(multiindex.IteratorType, func(*Value)) bool
}

type MultiIndexNode struct {
	super *SuperNode
}

/*generic class*/
type SuperNode struct {
	final *MultiIndexNode

	GetSuperNode func() (n interface{})
	GetFinalNode func() (n interface{})

	value func() *Value
}

//method for MultiIndex
func (m *MultiIndex) GetSuperIndex() interface{} { return m.super }
func (m *MultiIndex) GetFinalIndex() interface{} { return nil }

func (m *MultiIndex) GetIndex() interface{} {
	return nil
}

func (m *MultiIndex) Size() int {
	return m.count
}

func (m *MultiIndex) Clear() {
	m.super.clear()
	m.count = 0
}

func (m *MultiIndex) Insert(v Value) bool {
	_, res := m.insert(v)
	return res
}

func (m *MultiIndex) insert(v Value) (*MultiIndexNode, bool) {
	fn := &MultiIndexNode{}
	n, res := m.super.insert(v, fn)
	if res {
		fn.super = n
		m.count++
		return fn, true
	}
	return nil, false
}

func (m *MultiIndex) Erase(iter multiindex.IteratorType) {
	m.super.erase_(iter)
}

func (m *MultiIndex) erase(n *MultiIndexNode) {
	m.super.erase(n.super)
	m.count--
}

func (m *MultiIndex) Modify(iter multiindex.IteratorType, mod func(*Value)) bool {
	return m.super.modify_(iter, mod)
}

func (m *MultiIndex) modify(mod func(*Value), n *MultiIndexNode) (*MultiIndexNode, bool) {
	defer func() {
		if e := recover(); e != nil {
			container.Logger.Error("#multi modify failed: %v", e)
			m.erase(n)
			m.count--
			panic(e)
		}
	}()
	mod(n.value())
	if sn, res := m.super.modify(n.super); !res {
		m.count--
		return nil, false
	} else {
		n.super = sn
		return n, true
	}
}

func (n *MultiIndexNode) GetSuperNode() interface{} { return n.super }
func (n *MultiIndexNode) GetFinalNode() interface{} { return nil }

func (n *MultiIndexNode) value() *Value {
	return n.super.value()
}

/// IndexBase
type MultiIndexBase struct {
	final *MultiIndex
}

type MultiIndexBaseNode struct {
	final *MultiIndexNode
	pv    *Value
}

func (i *MultiIndexBase) init(final *MultiIndex) {
	i.final = final
}

func (i *MultiIndexBase) clear() {}

func (i *MultiIndexBase) GetSuperIndex() interface{} { return nil }

func (i *MultiIndexBase) GetFinalIndex() interface{} { return i.final }

func (i *MultiIndexBase) insert(v Value, fn *MultiIndexNode) (*MultiIndexBaseNode, bool) {
	return &MultiIndexBaseNode{fn, &v}, true
}

func (i *MultiIndexBase) erase(n *MultiIndexBaseNode) {
	n.pv = nil
}

func (i *MultiIndexBase) erase_(iter multiindex.IteratorType) {
	container.Logger.Warn("erase iterator doesn't match all index")
}

func (i *MultiIndexBase) modify(n *MultiIndexBaseNode) (*MultiIndexBaseNode, bool) {
	return n, true
}

func (i *MultiIndexBase) modify_(iter multiindex.IteratorType, mod func(*Value)) bool {
	container.Logger.Warn("modify iterator doesn't match all index")
	return false
}

func (n *MultiIndexBaseNode) value() *Value {
	return n.pv
}
