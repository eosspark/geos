// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package treemap implements a map backed by red-black tree.
//
// Elements are ordered by key in the map.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Associative_array
package treemap

import (
	"fmt"
	"github.com/eosspark/container/maps"
	rbt "github.com/eosspark/container/trees/redblacktree"
	"github.com/eosspark/container/utils"
	"strings"
	"reflect"
)

func assertMultiMapImplementation() {
	var _ maps.Map = (*MultiMap)(nil)
}

// Map holds the elements in a red-black tree
type MultiMap struct {
	KeyType   reflect.Type
	ValueType reflect.Type
	tree      *rbt.Tree
}

// NewWith instantiates a tree map with the custom comparator.
func NewMultiWith(keyType reflect.Type, valueType reflect.Type, comparator utils.Comparator) *MultiMap {
	return &MultiMap{KeyType: keyType, ValueType: valueType, tree: rbt.NewWith(comparator)}
}

// NewWithIntComparator instantiates a tree map with the IntComparator, i.e. keys are of type int.
func NewMultiWithIntComparator(valueType reflect.Type) *MultiMap {
	return &MultiMap{KeyType: utils.TypeInt, ValueType: valueType, tree: rbt.NewWithIntComparator()}
}

// NewWithStringComparator instantiates a tree map with the StringComparator, i.e. keys are of type string.
func NewMultiWithStringComparator(valueType reflect.Type) *MultiMap {
	return &MultiMap{KeyType: utils.TypeString, ValueType: valueType, tree: rbt.NewWithStringComparator()}
}

func CopyFromMulti(tm *MultiMap) *MultiMap {
	return &MultiMap{KeyType: tm.KeyType, ValueType: tm.ValueType, tree: rbt.CopyFrom(tm.tree)}
}

func (m *MultiMap) GetComparator() utils.Comparator {
	return m.tree.Comparator
}

// Put inserts key-value pair into the map.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *MultiMap) Put(key interface{}, value interface{}) {
	m.tree.MultiPut(key, value)
}

// Contains checks weather items (one or more) are present in the set.
// All items have to be present in the set for the method to return true.
// Returns true if no arguments are passed at all, i.e. set is always superset of empty set.
func (m *MultiMap) Get(key interface{}) (value interface{}, found bool) {
	return m.tree.Get(key)
}

// Get searches the element in the map by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *MultiMap) Gets(key interface{}) (MultiMapIterator, MultiMapIterator) {
	lower, upper := m.tree.MultiGet(key)
	return MultiMapIterator{iterator: lower}, MultiMapIterator{iterator: upper}
}

// Remove removes the element from the map by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *MultiMap) Remove(key interface{}) {
	m.tree.MultiRemove(key)
}

// Empty returns true if map does not contain any elements
func (m *MultiMap) Empty() bool {
	return m.tree.Empty()
}

// Size returns number of elements in the map.
func (m *MultiMap) Size() int {
	return m.tree.Size()
}

// Keys returns all keys in-order
func (m *MultiMap) Keys() []interface{} {
	return m.tree.Keys()
}

// Values returns all values in-order based on the key.
func (m *MultiMap) Values() []interface{} {
	return m.tree.Values()
}

// Clear removes all elements from the map.
func (m *MultiMap) Clear() {
	m.tree.Clear()
}

// Min returns the minimum key and its value from the tree map.
// Returns nil, nil if map is empty.
func (m *MultiMap) Min() (key interface{}, value interface{}) {
	if node := m.tree.Left(); node != nil {
		return node.Key, node.Value
	}
	return nil, nil
}

// Max returns the maximum key and its value from the tree map.
// Returns nil, nil if map is empty.
func (m *MultiMap) Max() (key interface{}, value interface{}) {
	if node := m.tree.Right(); node != nil {
		return node.Key, node.Value
	}
	return nil, nil
}

// Floor finds the floor key-value pair for the input key.
// In case that no floor is found, then both returned values will be nil.
// It's generally enough to check the first value (key) for nil, which determines if floor was found.
//
// Floor key is defined as the largest key that is smaller than or equal to the given key.
// A floor key may not be found, either because the map is empty, or because
// all keys in the map are larger than the given key.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *MultiMap) Floor(key interface{}) (foundKey interface{}, foundValue interface{}) {
	node, found := m.tree.Floor(key)
	if found {
		return node.Key, node.Value
	}
	return nil, nil
}

// Ceiling finds the ceiling key-value pair for the input key.
// In case that no ceiling is found, then both returned values will be nil.
// It's generally enough to check the first value (key) for nil, which determines if ceiling was found.
//
// Ceiling key is defined as the smallest key that is larger than or equal to the given key.
// A ceiling key may not be found, either because the map is empty, or because
// all keys in the map are smaller than the given key.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *MultiMap) Ceiling(key interface{}) (foundKey interface{}, foundValue interface{}) {
	node, found := m.tree.Ceiling(key)
	if found {
		return node.Key, node.Value
	}
	return nil, nil
}

// String returns a string representation of container
func (m *MultiMap) String() string {
	str := "TreeMap\nmap["
	it := m.Iterator()
	for it.Next() {
		str += fmt.Sprintf("%v:%v ", it.Key(), it.Value())
	}
	return strings.TrimRight(str, " ") + "]"

}
