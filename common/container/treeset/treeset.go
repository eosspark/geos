// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package treeset implements a Tree backed by a red-black Tree.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Set_%28abstract_data_type%29
package treeset

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/container/utils"
	"github.com/eosspark/eos-go/common/container"
	rbt "github.com/eosspark/eos-go/common/container/tree"
	"strings"
)

// template type Set(V,Compare)
type V = int

var Compare = utils.IntComparator

func assertSetImplementation() {
	var _ container.Set = (*Set)(nil)
}

// Set holds elements in a red-black Tree
type Set struct {
	*rbt.Tree
}

var itemExists = struct{}{}

// NewWith instantiates a new empty set with the custom comparator.

func New(Value ...V) *Set {
	set := &Set{Tree: rbt.NewWith(Compare, false)}
	set.Add(Value...)
	return set
}

func CopyFrom(ts *Set) *Set {
	return &Set{Tree: rbt.CopyFrom(ts.Tree)}
}

type MultiSet = Set

func NewMulti(items ...V) *MultiSet {
	set := &Set{Tree: rbt.NewWith(Compare, true)}
	set.Add(items...)
	return set
}

func SetIntersection(a *Set, b *Set, callback func(elem V)) {
	aIterator := a.Iterator()
	bIterator := b.Iterator()

	if !aIterator.First() || !bIterator.First() {
		return
	}

	for aHasNext, bHasNext := true, true; aHasNext && bHasNext; {
		comp := Compare(aIterator.Value(), bIterator.Value())
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

// Add adds the item one to the set.Returns false and the interface if it already exists
func (set *Set) AddItem(item V) (bool, V) {
	itr := set.Tree.Insert(item, itemExists)
	if itr.IsEnd() {
		return false, item
	}
	return true, itr.Value().(V)
}

// Add adds the items (one or more) to the set.
func (set *Set) Add(items ...V) {
	for _, item := range items {
		set.Tree.Put(item, itemExists)
	}
}

// Remove removes the items (one or more) from the set.
func (set *Set) Remove(items ...V) {
	for _, item := range items {
		set.Tree.Remove(item)
	}

}

// Values returns all items in the set.
func (set *Set) Values() []V {
	keys := make([]V, set.Size())
	it := set.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Value()
	}
	return keys
}

// Contains checks weather items (one or more) are present in the set.
// All items have to be present in the set for the method to return true.
// Returns true if no arguments are passed at all, i.e. set is always superset of empty set.
func (set *Set) Contains(items ...V) bool {
	for _, item := range items {
		if iter := set.Get(item); iter.IsEnd() {
			return false
		}
	}
	return true
}

// String returns a string representation of container
func (set *Set) String() string {
	str := "TreeSet\n"
	items := make([]string, 0)
	for _, v := range set.Tree.Keys() {
		items = append(items, fmt.Sprintf("%v", v))
	}
	str += strings.Join(items, ", ")
	return str
}

// Iterator returns a stateful iterator whose values can be fetched by an index.
type Iterator struct {
	rbt.Iterator
}

// Iterator holding the iterator's state
func (set *Set) Iterator() Iterator {
	return Iterator{Iterator: set.Tree.Iterator()}
}

// Begin returns First Iterator whose position points to the first element
// Return End Iterator when the map is empty
func (set *Set) Begin() Iterator {
	return Iterator{set.Tree.Begin()}
}

// End returns End Iterator
func (set *Set) End() Iterator {
	return Iterator{set.Tree.End()}
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator *Iterator) Value() V {
	return iterator.Iterator.Key().(V)
}

// Each calls the given function once for each element, passing that element's index and value.
func (set *Set) Each(f func(value V)) {
	iterator := set.Iterator()
	for iterator.Next() {
		f(iterator.Value())
	}
}

// Find passes each element of the container to the given function and returns
// the first (index,value) for which the function is true or -1,nil otherwise
// if no element matches the criteria.
func (set *Set) Find(f func(value V) bool) (v V) {
	iterator := set.Iterator()
	for iterator.Next() {
		if f(iterator.Value()) {
			return iterator.Value()
		}
	}
	return
}

func (set *Set) LowerBound(item V) Iterator {
	return Iterator{set.Tree.LowerBound(item)}
}

func (set *Set) UpperBound(item V) Iterator {
	return Iterator{set.Tree.UpperBound(item)}
}

// ToJSON outputs the JSON representation of the set.
func (set Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.Values())
}

// FromJSON populates the set from the input JSON representation.
func (set *Set) UnmarshalJSON(data []byte) error {
	elements := make([]V, 0)
	err := json.Unmarshal(data, &elements)
	if err == nil {
		set.Clear()
		set.Add(elements...)
	}
	return err
}
