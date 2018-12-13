// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package treeset implements a tree backed by a red-black tree.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Set_%28abstract_data_type%29
package treeset

import (
	"fmt"
	"github.com/eosspark/container/sets"
	rbt "github.com/eosspark/container/trees/redblacktree"
	"github.com/eosspark/container/utils"
	"strings"
	"reflect"
)

func assertMultiSetImplementation() {
	var _ sets.Set = (*MultiSet)(nil)
}

// Set holds elements in a red-black tree
type MultiSet struct {
	ValueType reflect.Type
	tree      *rbt.Tree
}

//var itemExists = struct{}{}

// NewWith instantiates a new empty set with the custom comparator.
func NewMultiWith(valueType reflect.Type, comparator utils.Comparator, values ...interface{}) *MultiSet {
	set := &MultiSet{ValueType: valueType, tree: rbt.NewWith(comparator)}
	if len(values) > 0 {
		set.Add(values...)
	}
	return set
}

// NewWithIntComparator instantiates a new empty set with the IntComparator, i.e. keys are of type int.
func NewMultiWithIntComparator(values ...interface{}) *MultiSet {
	set := &MultiSet{ValueType: utils.TypeInt, tree: rbt.NewWithIntComparator()}
	if len(values) > 0 {
		set.Add(values...)
	}
	return set
}

// NewWithStringComparator instantiates a new empty set with the StringComparator, i.e. keys are of type string.
func NewMultiWithStringComparator(values ...interface{}) *MultiSet {
	set := &MultiSet{ValueType: utils.TypeString, tree: rbt.NewWithStringComparator()}
	if len(values) > 0 {
		set.Add(values...)
	}
	return set
}

func CopyFromMulti(ts *MultiSet) *MultiSet {
	return &MultiSet{ValueType: ts.ValueType, tree: rbt.CopyFrom(ts.tree)}
}

func MultiSetIntersection(a *MultiSet, b *MultiSet, callback func(elem interface{})) {
	aIterator := a.Iterator()
	bIterator := b.Iterator()

	if !aIterator.First() || !bIterator.First() {
		return
	}

	comparator := a.GetComparator()

	for aHasNext, bHasNext := true, true; aHasNext && bHasNext; {
		comp := comparator(aIterator.Value(), bIterator.Value())
		switch {
		case comp > 0:
			bHasNext = bIterator.Next()
		case comp < 0:
			aHasNext = aIterator.Next()
		default:
			callback(aIterator.Value())
			aHasNext = aIterator.Next()
			bHasNext = bIterator.Next()
		}
	}
}

func (set *MultiSet) GetComparator() utils.Comparator {
	return set.tree.Comparator
}

// Add adds the items (one or more) to the set.
func (set *MultiSet) Add(items ...interface{}) {
	for _, item := range items {
		set.tree.MultiPut(item, itemExists)
	}
}

// Remove removes the items (one or more) from the set.
func (set *MultiSet) Remove(items ...interface{}) {
	for _, item := range items {
		set.tree.MultiRemove(item)
	}
}

// Contains checks weather items (one or more) are present in the set.
// All items have to be present in the set for the method to return true.
// Returns true if no arguments are passed at all, i.e. set is always superset of empty set.
func (set *MultiSet) Contains(items ...interface{}) bool {
	for _, item := range items {
		if _, contains := set.tree.Get(item); !contains {
			return false
		}
	}
	return true
}

func (set *MultiSet) Get(key interface{}) (MultiSetIterator, bool) {
	iterator, found := set.tree.MultiGet(key)
	return MultiSetIterator{iterator: iterator, tree: set.tree}, found
}

// Empty returns true if set does not contain any elements.
func (set *MultiSet) Empty() bool {
	return set.tree.Size() == 0
}

// Size returns number of elements within the set.
func (set *MultiSet) Size() int {
	return set.tree.Size()
}

// Clear clears all values in the set.
func (set *MultiSet) Clear() {
	set.tree.Clear()
}

// Values returns all items in the set.
func (set *MultiSet) Values() []interface{} {
	return set.tree.Keys()
}

// String returns a string representation of container
func (set *MultiSet) String() string {
	str := "MultiTreeSet\n"
	items := []string{}
	for _, v := range set.tree.Keys() {
		items = append(items, fmt.Sprintf("%v", v))
	}
	str += strings.Join(items, ", ")
	return str
}

func (set *MultiSet) UpperBound(in interface{}) *MultiSetIterator{
	if set.Size()>0{
		mItr:=set.Iterator()
		for mItr.Next(){
			comp:=set.GetComparator()(in,mItr.Value())
			if comp == -1{
				return &mItr
			}
		}
	}
	return nil
}

func (set *MultiSet) LowerBound(in interface{}) *MultiSetIterator{
	if set.Size()>0{
		mItr:=set.Iterator()
		for mItr.Next(){
			comp:=set.GetComparator()(mItr.Value(),in)
			if comp==0{
				return &mItr
			}
		}
	}
	return nil
}