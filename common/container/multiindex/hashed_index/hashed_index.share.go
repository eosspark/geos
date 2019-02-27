//+build shared

package hashed_index

import (
	"github.com/eosspark/eos-go/common/allocator"
	"github.com/eosspark/eos-go/common/container"
	"github.com/eosspark/eos-go/common/container/multiindex"
	. "github.com/eosspark/eos-go/common/offsetptr"
	"unsafe"
)

// template type HashedIndex(FinalIndex,FinalNode,SuperIndex,SuperNode,Value,Key,KeyFunc,Hasher,Allocator)
type Key = int
type Value = int

var Allocator allocator.MemoryManager = nil
var KeyFunc = func(Value) Key { return 0 }
var Hasher = func(h Key) uintptr { return container.Hash(h) }

//TODO nonunique_hashtable const Multiply = false

type HashedIndex struct {
	super Pointer `*SuperIndex` // index on the HashedUniqueIndex, IndexBase is the last super index
	final Pointer `*FinalIndex` // index under the HashedUniqueIndex, MultiIndex is the final index
	b     Pointer `*Buckets`
	l     uintptr
}

func (h *HashedIndex) init(final *FinalIndex) {
	h.final.Set(unsafe.Pointer(final))
	h.super.Set(unsafe.Pointer(NewSuperIndex()))
	h.b.Set(unsafe.Pointer(NewBuckets(0)))
	//h.b.NewBuckets(0)
	h.l = 0
	(*SuperIndex)(h.super.Get()).init(final)
}

func (h *HashedIndex) free() {
	(*SuperIndex)(h.super.Get()).free()
	Allocator.DeAllocate(unsafe.Pointer(h))
}

/*generic class*/
type SuperIndex struct {
	init    func(*FinalIndex)
	free    func()
	clear   func()
	insert  func(Value, *FinalNode) (*SuperNode, bool)
	erase   func(*SuperNode) bool
	erase_  func(multiindex.IteratorType)
	modify  func(*SuperNode) (*SuperNode, bool)
	modify_ func(multiindex.IteratorType, func(*Value)) bool
}

func NewSuperIndex() *SuperIndex {
	return (*SuperIndex)(Allocator.Allocate(unsafe.Sizeof(SuperIndex{})))
}

/*generic class*/
type FinalIndex struct {
	insert func(Value) (*FinalNode, bool)
	erase  func(*FinalNode)
	modify func(func(*Value), *FinalNode) (*FinalNode, bool)
}

type Buckets struct {
	array Pointer `Pointer<*HashedIndexNode>`
	len   uintptr
}

const _SizeofBuckets = unsafe.Sizeof(Buckets{})
const _SizeofBucket = unsafe.Sizeof(&HashedIndexNode{})

func NewBuckets(size uintptr) *Buckets {
	b := (*Buckets)(Allocator.Allocate(_SizeofBuckets))
	if size == 0 {
		size++
	}
	b.array.Set(Allocator.Allocate(_SizeofBucket * size))
	//b.array = (**HashedIndexNode)(Allocator.Allocate(_SizeofBucket * size))
	allocator.Memset(b.array.Get(), 0, _SizeofBucket*size)
	//allocator.Memset(unsafe.Pointer(b.array), 0, _SizeofBucket*size)
	b.len = size
	return b
}

func (b *Buckets) Put(i uintptr, e *HashedIndexNode) {
	(*Pointer)(unsafe.Pointer(uintptr(b.array.Get()) + _SizeofBucket*i)).Set(unsafe.Pointer(e))
	//*(**HashedIndexNode)(unsafe.Pointer(uintptr(b.array.Get()) + _SizeofBucket*i)) = e
	//*(**HashedIndexNode)(unsafe.Pointer(uintptr(unsafe.Pointer(b.array)) + _SizeofBucket*i)) = e
}

func (b *Buckets) At(i uintptr) *HashedIndexNode {
	p := (*Pointer)(unsafe.Pointer(uintptr(b.array.Get()) + _SizeofBucket*i))
	if p.IsSelf() {
		p.Set(nil)
	}

	return (*HashedIndexNode)(p.Get())
	//return *(**HashedIndexNode)(unsafe.Pointer(uintptr(b.array.Get()) + _SizeofBucket*i))
	//return *(**HashedIndexNode)(unsafe.Pointer(uintptr(unsafe.Pointer(b.array)) + _SizeofBucket*i))
}

type HashedIndexNode struct {
	super  Pointer `*SuperNode` // index-node on the HashedUniqueIndexNode, IndexBaseNode is the last super node
	final  Pointer `*FinalNode` // index-node under the HashedUniqueIndexNode, MultiIndexNode is the final index
	bucket uintptr              // buckets position
	key    Key                  // k of hashtable
	next   Pointer `*HashedIndexNode`
}

func NewHashedIndexNode(bucket uintptr, key Key) *HashedIndexNode {
	n := (*HashedIndexNode)(Allocator.Allocate(unsafe.Sizeof(HashedIndexNode{})))
	n.key = key
	n.bucket = bucket
	n.super.Set(nil)
	n.final.Set(nil)
	n.next.Set(nil)
	//n.next = nil
	return n
}

func (n *HashedIndexNode) free() {
	if n != nil {
		Allocator.DeAllocate(unsafe.Pointer(n))
	}
}

/*generic class*/
type SuperNode struct {
	value func() *Value
}

/*generic class*/
type FinalNode struct {
	GetSuperNode func() interface{}
	GetFinalNode func() interface{}
}

func (h *HashedIndex) GetSuperIndex() interface{} { return h.super }
func (h *HashedIndex) GetFinalIndex() interface{} { return h.final }

func (n *HashedIndexNode) GetSuperNode() interface{} { return n.super }
func (n *HashedIndexNode) GetFinalNode() interface{} { return n.final }

func (n *HashedIndexNode) value() *Value {
	return (*SuperNode)(n.super.Get()).value()
}

func (n *HashedIndexNode) Get(key Key) *HashedIndexNode {
	for node := n; node != nil; node = (*HashedIndexNode)(node.next.Get()) {
		//for node := n; node != nil; node = node.next {
		if node.key == key {
			return node
		}
	}
	return nil
}

func (h *HashedIndex) Size() int {
	return int(h.l)
}

func (h *HashedIndex) Empty() bool {
	return h.l == 0
}

func (h *HashedIndex) clear() {
	//h.inner = map[Hash]*HashedIndexNode{}
	(*SuperIndex)(h.super.Get()).clear()
}

func (h *HashedIndex) Insert(v Value) (Iterator, bool) {
	fn, res := (*FinalIndex)(h.final.Get()).insert(v)
	if res {
		return h.makeIterator(fn), true
	}
	return h.End(), false
}

func (h *HashedIndex) insert(v Value, fn *FinalNode) (*HashedIndexNode, bool) {
	key := KeyFunc(v)
	b := (*Buckets)(h.b.Get())
	bucket := Hasher(key) % b.len
	//bucket := Hasher(key) % h.b.len
	node := NewHashedIndexNode(bucket, key)

	if !h.link(bucket, key, node) {
		container.Logger.Warn("#hash index insert failed")
		node.free()
		return nil, false
	}

	if h.l > b.len {
		//if h.l > h.b.len {
		h.resize()
	}

	sn, res := (*SuperIndex)(h.super.Get()).insert(v, fn)
	if res {
		node.final.Set(unsafe.Pointer(fn))
		node.super.Set(unsafe.Pointer(sn))
		return node, true
	}

	//rollback for failed insert of SuperIndex
	h.unlink(node) //should never failed
	return nil, false
}

func (h *HashedIndex) Find(k Key) Iterator {
	b := (*Buckets)(h.b.Get())
	return Iterator{h, b.At(Hasher(k) % b.len).Get(k)}
	//return Iterator{h, h.b.At(Hasher(k) % h.b.len).Get(k)}
}

func (h *HashedIndex) Erase(iter Iterator) {
	(*FinalIndex)(h.final.Get()).erase((*FinalNode)(iter.node.final.Get()))
}

func (h *HashedIndex) erase(n *HashedIndexNode) {
	h.unlink(n) //should never failed
	(*SuperIndex)(h.super.Get()).erase((*SuperNode)(n.super.Get()))
	n.free()
}

func (h *HashedIndex) erase_(iter multiindex.IteratorType) {
	if itr, ok := iter.(Iterator); ok {
		h.Erase(itr)
	} else {
		(*SuperIndex)(h.super.Get()).erase_(iter)
	}
}

func (h *HashedIndex) Modify(iter Iterator, mod func(*Value)) bool {
	if _, b := (*FinalIndex)(h.final.Get()).modify(mod, (*FinalNode)(iter.node.final.Get())); b {
		return true
	}
	return false
}

func (h *HashedIndex) modify(n *HashedIndexNode) (*HashedIndexNode, bool) {
	key := KeyFunc(*n.value())
	bucket := Hasher(key) % (*Buckets)(h.b.Get()).len
	//bucket := Hasher(key) % h.b.len

	if !h.inPlace(bucket, key, n) {
		h.unlink(n)
		if !h.link(bucket, key, n) {
			container.Logger.Warn("#hash index modify failed")
			(*SuperIndex)(h.super.Get()).erase((*SuperNode)(n.super.Get()))
			n.free()
			return nil, false
		}
	}

	if sn, res := (*SuperIndex)(h.super.Get()).modify((*SuperNode)(n.super.Get())); !res {
		h.unlink(n)
		n.free()
		return nil, false
	} else {
		n.super.Set(unsafe.Pointer(sn))
	}

	return n, true
}

func (h *HashedIndex) modify_(iter multiindex.IteratorType, mod func(*Value)) bool {
	if itr, ok := iter.(Iterator); ok {
		return h.Modify(itr, mod)
	} else {
		return (*SuperIndex)(h.super.Get()).modify_(iter, mod)
	}
}

func (h *HashedIndex) Values() []Value {
	vs := make([]Value, 0, h.Size())
	for it := h.Begin(); it.HasNext(); it.Next() {
		vs = append(vs, it.Value())
	}
	return vs
}

func (h *HashedIndex) resize() {
	b := (*Buckets)(h.b.Get())
	n := b.len * 2
	//n := h.b.len * 2
	tmp := NewBuckets(n)
	for bucket := uintptr(0); bucket < b.len; bucket++ {
		//for bucket := uintptr(0); bucket < h.b.len; bucket++ {
		first := b.At(bucket)
		//first := h.b.At(bucket)
		for first != nil {
			newBucket := Hasher(first.key) % n
			b.Put(bucket, (*HashedIndexNode)(first.next.Get()))
			//h.b.Put(bucket, first.next)

			first.bucket = newBucket
			first.next.Set(unsafe.Pointer(tmp.At(newBucket)))
			//first.next = tmp.At(newBucket)

			tmp.Put(newBucket, first)
			first = b.At(bucket)
			//first = h.b.At(bucket)
		}
	}

	Allocator.DeAllocate(unsafe.Pointer(b))
	//Allocator.DeAllocate(unsafe.Pointer(h.b))
	h.b.Set(unsafe.Pointer(tmp))
	//h.b = tmp
}

func (h *HashedIndex) inPlace(buc uintptr, k Key, x *HashedIndexNode) bool {
	b := (*Buckets)(h.b.Get())
	found := false
	for y := b.At(buc); y != nil; y = (*HashedIndexNode)(y.next.Get()) {
		//for y := h.b.At(buc); y != nil; y = y.next {
		if x == y {
			found = true
		} else if k == y.key {
			found = false
		}
	}
	return found
}

func (h *HashedIndex) unlink(n *HashedIndexNode) {
	b := (*Buckets)(h.b.Get())
	var node, prev *HashedIndexNode

	found := false
	for node = b.At(n.bucket); node != nil; node = (*HashedIndexNode)(node.next.Get()) {
		//for node = h.b.At(n.bucket); node != nil; node = node.next {
		if node == n {
			break
		}
		prev = node
	}

	if !found {
		container.Logger.Warn("unlink node not found")
	}

	if prev != nil {
		prev.next.Forward(&node.next)
	} else {
		b.Put(n.bucket, (*HashedIndexNode)(node.next.Get()))
		//h.b.Put(n.bucket, node.next)
	}

	h.l--
}

func (h *HashedIndex) link(buc uintptr, k Key, n *HashedIndexNode) bool {
	b := (*Buckets)(h.b.Get())
	var prev *HashedIndexNode
	for node := b.At(buc); node != nil; node = (*HashedIndexNode)(node.next.Get()) {
		//for node := h.b.At(buc); node != nil; node = node.next {
		if node.key == k {
			return false
		}
		prev = node
	}

	if prev != nil {
		prev.next.Set(unsafe.Pointer(n))
		//prev.next = n
	} else {
		b.Put(buc, n)
		//h.b.Put(buc, n)
	}

	n.key = k
	n.next.Set(nil)
	//n.next = nil
	h.l++

	return true
}

type Iterator struct {
	index *HashedIndex
	node  *HashedIndexNode
}

func (h *HashedIndex) Begin() Iterator {
	if h.l == 0 {
		return h.End()
	}

	b := (*Buckets)(h.b.Get())
	for i := uintptr(0); i < b.len; i++ {
		//for i := uintptr(0); i < h.b.len; i++ {
		if e := b.At(i); e != nil {
			//if e := h.b.At(i); e != nil {
			return Iterator{h, e}
		}
	}

	panic(container.ErrFatalAddress)
}

func (h *HashedIndex) makeIterator(fn *FinalNode) Iterator {
	node := fn.GetSuperNode()
	for {
		if node == nil {
			panic("Wrong index node type!")

		} else if n, ok := node.(*HashedIndexNode); ok {
			return Iterator{h, n}
		} else {
			node = node.(multiindex.NodeType).GetSuperNode()
		}
	}
}

func (h *HashedIndex) End() Iterator {
	return Iterator{h, nil}
}

func (iter Iterator) Value() (v Value) {
	return *iter.node.value()
}

func (iter *Iterator) Next() bool {
	b := (*Buckets)(iter.index.b.Get())

	if !iter.node.next.IsNil() {
		//if iter.node.next != nil {
		iter.node = (*HashedIndexNode)(iter.node.next.Get())
		//iter.node = iter.node.next
		return true
	}

	for bucket := (Hasher(iter.node.key) % b.len) + 1; bucket < b.len; bucket++ {
		//for bucket := (Hasher(iter.node.key) % iter.index.b.len) + 1; bucket < iter.index.b.len; bucket++ {
		if entry := b.At(bucket); entry != nil {
			//if entry := iter.index.b.At(bucket); entry != nil {
			iter.node = entry
			return true
		}
	}

	iter.node = nil
	return false
}

func (iter Iterator) HasNext() bool {
	return iter.node != nil
}

func (iter Iterator) IsEnd() bool {
	return iter.node == nil
}
