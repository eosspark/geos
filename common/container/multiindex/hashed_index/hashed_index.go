//+build !shared

package hashed_index

import (
	"github.com/eosspark/eos-go/common/container"
	"github.com/eosspark/eos-go/common/container/multiindex"
)

// template type HashedUniqueIndex(FinalIndex,FinalNode,SuperIndex,SuperNode,Value,Hash,KeyFunc)
type Hash = int
type Value = int

var KeyFunc = func(Value) Hash { return 0 }

type HashedUniqueIndex struct {
	super *SuperIndex                     // index on the HashedUniqueIndex, IndexBase is the last super index
	final *FinalIndex                     // index under the HashedUniqueIndex, MultiIndex is the final index
	inner map[Hash]*HashedUniqueIndexNode // use hashmap to safe HashedUniqueIndex's k/v(HashedUniqueIndexNode)
}

func (i *HashedUniqueIndex) init(final *FinalIndex) {
	i.final = final
	i.inner = map[Hash]*HashedUniqueIndexNode{}
	i.super = &SuperIndex{}
	i.super.init(final)
}

/*generic class*/
type SuperIndex struct {
	init    func(*FinalIndex)
	clear   func()
	insert  func(Value, *FinalNode) (*SuperNode, bool)
	erase   func(*SuperNode) bool
	erase_  func(multiindex.IteratorType)
	modify  func(*SuperNode) (*SuperNode, bool)
	modify_ func(multiindex.IteratorType, func(*Value)) bool
}

/*generic class*/
type FinalIndex struct {
	insert func(Value) (*FinalNode, bool)
	erase  func(*FinalNode)
	modify func(func(*Value), *FinalNode) (*FinalNode, bool)
}

type HashedUniqueIndexNode struct {
	super *SuperNode // index-node on the HashedUniqueIndexNode, IndexBaseNode is the last super node
	final *FinalNode // index-node under the HashedUniqueIndexNode, MultiIndexNode is the final index
	hash  Hash       // k of hashmap
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

func (i *HashedUniqueIndex) GetSuperIndex() interface{} { return i.super }
func (i *HashedUniqueIndex) GetFinalIndex() interface{} { return i.final }

func (n *HashedUniqueIndexNode) GetSuperNode() interface{} { return n.super }
func (n *HashedUniqueIndexNode) GetFinalNode() interface{} { return n.final }

func (n *HashedUniqueIndexNode) value() *Value {
	return n.super.value()
}

func (i *HashedUniqueIndex) Size() int {
	return len(i.inner)
}

func (i *HashedUniqueIndex) Empty() bool {
	return len(i.inner) == 0
}

func (i *HashedUniqueIndex) clear() {
	i.inner = map[Hash]*HashedUniqueIndexNode{}
	i.super.clear()
}

func (i *HashedUniqueIndex) Insert(v Value) (Iterator, bool) {
	fn, res := i.final.insert(v)
	if res {
		return i.makeIterator(fn), true
	}
	return i.End(), false
}

func (i *HashedUniqueIndex) insert(v Value, fn *FinalNode) (*HashedUniqueIndexNode, bool) {
	hash := KeyFunc(v)
	node := HashedUniqueIndexNode{hash: hash}
	if _, ok := i.inner[hash]; ok {
		container.Logger.Warn("#hash index insert failed")
		return nil, false
	}
	i.inner[hash] = &node
	sn, res := i.super.insert(v, fn)
	if res {
		node.final = fn
		node.super = sn
		return &node, true
	}
	delete(i.inner, hash)
	return nil, false
}

func (i *HashedUniqueIndex) Find(k Hash) (Iterator, bool) {
	node, res := i.inner[k]
	if res {
		return Iterator{i, node, between}, true
	}
	return i.End(), false
}

func (i *HashedUniqueIndex) Each(f func(key Hash, obj Value)) {
	for k, v := range i.inner {
		f(k, *v.value())
	}
}

func (i *HashedUniqueIndex) Erase(iter Iterator) {
	i.final.erase(iter.node.final)
}

func (i *HashedUniqueIndex) erase(n *HashedUniqueIndexNode) {
	delete(i.inner, n.hash)
	i.super.erase(n.super)
}

func (i *HashedUniqueIndex) erase_(iter multiindex.IteratorType) {
	if itr, ok := iter.(Iterator); ok {
		i.Erase(itr)
	} else {
		i.super.erase_(iter)
	}
}

func (i *HashedUniqueIndex) Modify(iter Iterator, mod func(*Value)) bool {
	if _, b := i.final.modify(mod, iter.node.final); b {
		return true
	}
	return false
}

func (i *HashedUniqueIndex) modify(n *HashedUniqueIndexNode) (*HashedUniqueIndexNode, bool) {
	delete(i.inner, n.hash)

	hash := KeyFunc(*n.value())
	if _, exist := i.inner[hash]; exist {
		container.Logger.Warn("#hash index modify failed")
		i.super.erase(n.super)
		return nil, false
	}

	i.inner[hash] = n

	if sn, res := i.super.modify(n.super); !res {
		delete(i.inner, hash)
		return nil, false
	} else {
		n.super = sn
	}

	return n, true
}

func (i *HashedUniqueIndex) modify_(iter multiindex.IteratorType, mod func(*Value)) bool {
	if itr, ok := iter.(Iterator); ok {
		return i.Modify(itr, mod)
	} else {
		return i.super.modify_(iter, mod)
	}
}

func (i *HashedUniqueIndex) Values() []Value {
	vs := make([]Value, 0, i.Size())
	i.Each(func(key Hash, obj Value) {
		vs = append(vs, obj)
	})
	return vs
}

type Iterator struct {
	index    *HashedUniqueIndex
	node     *HashedUniqueIndexNode
	position pos
}

type pos byte

const (
	//begin   = 0
	between = 1
	end     = 2
)

func (i *HashedUniqueIndex) makeIterator(fn *FinalNode) Iterator {
	node := fn.GetSuperNode()
	for {
		if node == nil {
			panic("Wrong index node type!")

		} else if n, ok := node.(*HashedUniqueIndexNode); ok {
			return Iterator{i, n, between}
		} else {
			node = node.(multiindex.NodeType).GetSuperNode()
		}
	}
}

func (i *HashedUniqueIndex) End() Iterator {
	return Iterator{i, nil, end}
}

func (iter Iterator) Value() (v Value) {
	if iter.position == between {
		return *iter.node.value()
	}
	return
}

func (iter Iterator) HasNext() bool {
	container.Logger.Warn("hashed index iterator is unmoveable")
	return false
}

func (iter Iterator) IsEnd() bool {
	return iter.position == end
}
