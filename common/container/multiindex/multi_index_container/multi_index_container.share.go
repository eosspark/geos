//+build shared

package multi_index_container

import (
	"github.com/eosspark/eos-go/common/allocator"
	"github.com/eosspark/eos-go/common/container"
	"github.com/eosspark/eos-go/common/container/multiindex"
	. "github.com/eosspark/eos-go/common/offsetptr"
	"unsafe"
)

// template type MultiIndex(SuperIndex,SuperNode,Value,Allocator)
type Value = int

var Allocator allocator.MemoryManager = nil

type MultiIndex struct {
	super Pointer `*SuperIndex`
	count int
}

const _SizeofMultiIndex = unsafe.Sizeof(MultiIndex{})

func New() (m *MultiIndex) {
	m = (*MultiIndex)(Allocator.Allocate(_SizeofMultiIndex))
	m.super.Set(Allocator.Allocate(_SizeofSuperIndex))
	m.count = 0

	(*SuperIndex)(m.super.Get()).init(m)
	return m
}

func (m *MultiIndex) Free() {
	(*SuperIndex)(m.super.Get()).free()
	Allocator.DeAllocate(unsafe.Pointer(m))
}

/*generic class*/
type SuperIndex struct {
	init    func(*MultiIndex)
	free    func()
	clear   func()
	insert  func(Value, *MultiIndexNode) (*SuperNode, bool)
	erase   func(*SuperNode)
	erase_  func(multiindex.IteratorType)
	modify  func(*SuperNode) (*SuperNode, bool)
	modify_ func(multiindex.IteratorType, func(*Value)) bool
}

const _SizeofSuperIndex = unsafe.Sizeof(SuperIndex{})

type MultiIndexNode struct {
	super Pointer `*SuperNode`
}

const _SizeofMultiIndexNode = unsafe.Sizeof(MultiIndexNode{})

func NewMultiIndexNode() *MultiIndexNode {
	n := (*MultiIndexNode)(Allocator.Allocate(_SizeofMultiIndexNode))
	n.super.Set(nil)
	return n
}

func (n *MultiIndexNode) free() {
	if n != nil {
		Allocator.DeAllocate(unsafe.Pointer(n))
	}
}

/*generic class*/
type SuperNode struct {
	final Pointer `*MultiIndexNode`

	GetSuperNode func() (n interface{})
	GetFinalNode func() (n interface{})
	value        func() *Value
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
	(*SuperIndex)(m.super.Get()).clear()
	m.count = 0
}

func (m *MultiIndex) Insert(v Value) bool {
	_, res := m.insert(v)
	return res
}

func (m *MultiIndex) insert(v Value) (*MultiIndexNode, bool) {
	fn := NewMultiIndexNode()
	n, res := (*SuperIndex)(m.super.Get()).insert(v, fn)
	if res {
		fn.super.Set(unsafe.Pointer(n))
		m.count++
		return fn, true
	}

	fn.free()
	return nil, false
}

func (m *MultiIndex) Erase(iter multiindex.IteratorType) {
	(*SuperIndex)(m.super.Get()).erase_(iter)
}

func (m *MultiIndex) erase(n *MultiIndexNode) {
	m.count-- // only sub count when MultiIndexNode erase itself
	(*SuperIndex)(m.super.Get()).erase((*SuperNode)(n.super.Get()))
	n.free() // free memory finally
}

func (m *MultiIndex) Modify(iter multiindex.IteratorType, mod func(*Value)) bool {
	return (*SuperIndex)(m.super.Get()).modify_(iter, mod)
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
	if sn, res := (*SuperIndex)(m.super.Get()).modify((*SuperNode)(n.super.Get())); !res {
		//delete for failure
		m.count--
		n.free()
		return nil, false
	} else {
		n.super.Set(unsafe.Pointer(sn))
		return n, true
	}
}

func (n *MultiIndexNode) GetSuperNode() interface{} { return (*SuperNode)(n.super.Get()) }
func (n *MultiIndexNode) GetFinalNode() interface{} { return nil }

func (n *MultiIndexNode) value() *Value {
	return (*SuperNode)(n.super.Get()).value()
}

/// IndexBase
type MultiIndexBase struct {
	final Pointer `*MultiIndex`
}

func (i *MultiIndexBase) init(final *MultiIndex) {
	i.final.Set(unsafe.Pointer(final))
}

func (i *MultiIndexBase) free() {
	if i != nil {
		Allocator.DeAllocate(unsafe.Pointer(i))
	}
}

func (i *MultiIndexBase) clear() {}

func (i *MultiIndexBase) GetSuperIndex() interface{} { return nil }

func (i *MultiIndexBase) GetFinalIndex() interface{} { return i.final }

func (i *MultiIndexBase) insert(v Value, fn *MultiIndexNode) (*MultiIndexBaseNode, bool) {
	return NewMultiIndexBaseNode(fn, v), true
}

func (i *MultiIndexBase) erase(n *MultiIndexBaseNode) {
	n.free()
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

type MultiIndexBaseNode struct {
	final Pointer `*MultiIndexNode`
	pv    Pointer `*Value`
}

const _SizeofMultiIndexBaseNode = unsafe.Sizeof(MultiIndexBaseNode{})

func NewMultiIndexBaseNode(final *MultiIndexNode, pv Value) (mn *MultiIndexBaseNode) {
	mn = (*MultiIndexBaseNode)(Allocator.Allocate(_SizeofMultiIndexBaseNode))
	pvAlloc := Allocator.Allocate(unsafe.Sizeof(pv))
	mn.pv.Set(pvAlloc)
	*(*Value)(pvAlloc) = pv
	mn.final.Set(unsafe.Pointer(final))
	return
}

func (n *MultiIndexBaseNode) free() {
	if n != nil {
		Allocator.DeAllocate(n.pv.Get())
		Allocator.DeAllocate(unsafe.Pointer(n))
	}
}

func (n *MultiIndexBaseNode) value() *Value {
	return (*Value)(n.pv.Get())
}
