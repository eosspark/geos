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

func assertSetImplementation() {
	var _ sets.Set = (*Set)(nil)
}

// Set holds elements in a red-black tree
type Set struct {
	ValueType reflect.Type
	tree      *rbt.Tree
}

var itemExists = struct{}{}

// NewWith instantiates a new empty set with the custom comparator.
func NewWith(valueType reflect.Type, comparator utils.Comparator, values ...interface{}) *Set {
	set := &Set{ValueType: valueType, tree: rbt.NewWith(comparator)}
	if len(values) > 0 {
		set.Add(values...)
	}
	return set
}

// NewWithIntComparator instantiates a new empty set with the IntComparator, i.e. keys are of type int.
func NewWithIntComparator(values ...interface{}) *Set {
	set := &Set{ValueType: utils.TypeInt, tree: rbt.NewWithIntComparator()}
	if len(values) > 0 {
		set.Add(values...)
	}
	return set
}

// NewWithStringComparator instantiates a new empty set with the StringComparator, i.e. keys are of type string.
func NewWithStringComparator(values ...interface{}) *Set {
	set := &Set{ValueType: utils.TypeString, tree: rbt.NewWithStringComparator()}
	if len(values) > 0 {
		set.Add(values...)
	}
	return set
}

func CopyFrom(ts *Set) *Set {
	return &Set{ValueType: ts.ValueType, tree: rbt.CopyFrom(ts.tree)}
}

func SetIntersection(a *Set, b *Set, callback func(elem interface{})) {
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

func (set *Set) GetComparator() utils.Comparator {
	return set.tree.Comparator
}

// Add adds the item one to the set.Returns false and the interface if it already exists
func (set *Set) AddItem(item interface{}) (bool, interface{}) {
	opt, k, _ := set.tree.PutItem(item, itemExists)
	return opt, k
}

// Add adds the items (one or more) to the set.
func (set *Set) Add(items ...interface{}) {
	for _, item := range items {
		set.tree.Put(item, itemExists)
	}
}

// Remove removes the items (one or more) from the set.
func (set *Set) Remove(items ...interface{}) {
	for _, item := range items {
		set.tree.Remove(item)
	}
}

// Contains checks weather items (one or more) are present in the set.
// All items have to be present in the set for the method to return true.
// Returns true if no arguments are passed at all, i.e. set is always superset of empty set.
func (set *Set) Contains(items ...interface{}) bool {
	for _, item := range items {
		if _, contains := set.tree.Get(item); !contains {
			return false
		}
	}
	return true
}

// Empty returns true if set does not contain any elements.
func (set *Set) Empty() bool {
	return set.tree.Size() == 0
}

// Size returns number of elements within the set.
func (set *Set) Size() int {
	return set.tree.Size()
}

// Clear clears all values in the set.
func (set *Set) Clear() {
	set.tree.Clear()
}

// Values returns all items in the set.
func (set *Set) Values() []interface{} {
	return set.tree.Keys()
}

// String returns a string representation of container
func (set *Set) String() string {
	str := "TreeSet\n"
	items := []string{}
	for _, v := range set.tree.Keys() {
		items = append(items, fmt.Sprintf("%v", v))
	}
	str += strings.Join(items, ", ")
	return str
}
