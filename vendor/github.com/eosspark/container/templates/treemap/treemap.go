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
	"encoding/json"
	"fmt"
	"github.com/eosspark/container/templates"
	rbt "github.com/eosspark/container/trees/redblacktree"
	"github.com/eosspark/container/utils"
	"strings"
)

// template type Map(K,V,Compare)
type K int
type V int

var Compare = utils.IntComparator

func assertMapImplementation() {
	var _ templates.Map = (*Map)(nil)
}

// Map holds the elements in a red-black tree
type Map struct {
	isMulti bool
	tree    *rbt.Tree
}

// NewWith instantiates a tree map with the custom comparator.
func New() *Map {
	return &Map{tree: rbt.NewWith(Compare)}
}

func CopyFrom(tm *Map) *Map {
	return &Map{tree: rbt.CopyFrom(tm.tree)}
}

func (m *Map) GetComparator() utils.Comparator {
	return m.tree.Comparator
}

// Put inserts key-value pair into the map.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map) Put(key K, value V) {
	if m.isMulti {
		m.tree.MultiPut(key, value)
	} else {
		m.tree.Put(key, value)
	}
}

// Get searches the element in the map by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map) Get(key K) (value V, found bool) {
	if v, ok := m.tree.Get(key); ok {
		return v.(V), ok
	}
	return
}

// Remove removes the element from the map by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map) Remove(key K) {
	if m.isMulti {
		m.tree.MultiRemove(key)
	} else {
		m.tree.Remove(key)
	}
}

// Empty returns true if map does not contain any elements
func (m *Map) Empty() bool {
	return m.tree.Empty()
}

// Size returns number of elements in the map.
func (m *Map) Size() int {
	return m.tree.Size()
}

// Keys returns all keys in-order
func (m *Map) Keys() []K {
	keys := make([]K, m.tree.Size())
	it := m.tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key().(K)
	}
	return keys
}

// Values returns all values in-order based on the key.
func (m *Map) Values() []V {
	values := make([]V, m.tree.Size())
	it := m.tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value().(V)
	}
	return values
}

// Clear removes all elements from the map.
func (m *Map) Clear() {
	m.tree.Clear()
}

// Min returns the minimum key and its value from the tree map.
// Returns nil, nil if map is empty.
func (m *Map) Min() (key K, value V) {
	if node := m.tree.Left(); node != nil {
		return node.Key.(K), node.Value.(V)
	}
	return
}

// Max returns the maximum key and its value from the tree map.
// Returns nil, nil if map is empty.
func (m *Map) Max() (key K, value V) {
	if node := m.tree.Right(); node != nil {
		return node.Key.(K), node.Value.(V)
	}
	return
}

// Each calls the given function once for each element, passing that element's key and value.
func (m *Map) Each(f func(key K, value V)) {
	iterator := m.Iterator()
	for iterator.Next() {
		f(iterator.Key(), iterator.Value())
	}
}

// Map invokes the given function once for each element and returns a container
// containing the values returned by the given function as key/value pairs.
func (m *Map) Map(f func(key1 K, value1 V) (K, V)) *Map {
	newMap := &Map{tree: rbt.NewWith(m.tree.Comparator)}
	iterator := m.Iterator()
	for iterator.Next() {
		key2, value2 := f(iterator.Key(), iterator.Value())
		newMap.Put(key2, value2)
	}
	return newMap
}

// Select returns a new container containing all elements for which the given function returns a true value.
func (m *Map) Select(f func(key K, value V) bool) *Map {
	newMap := &Map{tree: rbt.NewWith(m.tree.Comparator)}
	iterator := m.Iterator()
	for iterator.Next() {
		if f(iterator.Key(), iterator.Value()) {
			newMap.Put(iterator.Key(), iterator.Value())
		}
	}
	return newMap
}

// Any passes each element of the container to the given function and
// returns true if the function ever returns true for any element.
func (m *Map) Any(f func(key K, value V) bool) bool {
	iterator := m.Iterator()
	for iterator.Next() {
		if f(iterator.Key(), iterator.Value()) {
			return true
		}
	}
	return false
}

// All passes each element of the container to the given function and
// returns true if the function returns true for all elements.
func (m *Map) All(f func(key K, value V) bool) bool {
	iterator := m.Iterator()
	for iterator.Next() {
		if !f(iterator.Key(), iterator.Value()) {
			return false
		}
	}
	return true
}

// Find passes each element of the container to the given function and returns
// the first (key,value) for which the function is true or nil,nil otherwise if no element
// matches the criteria.
func (m *Map) Find(f func(key K, value V) bool) (k K, v V) {
	iterator := m.Iterator()
	for iterator.Next() {
		if f(iterator.Key(), iterator.Value()) {
			return iterator.Key(), iterator.Value()
		}
	}
	return
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
func (m *Map) Floor(key K) (foundKey K, foundValue V) {
	node, found := m.tree.Floor(key)
	if found {
		return node.Key.(K), node.Value.(V)
	}
	return
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
func (m *Map) Ceiling(key K) (foundKey K, foundValue V) {
	node, found := m.tree.Ceiling(key)
	if found {
		return node.Key.(K), node.Value.(V)
	}
	return
}

// String returns a string representation of container
func (m *Map) String() string {
	str := "TreeMap\nmap["
	it := m.Iterator()
	for it.Next() {
		str += fmt.Sprintf("%v:%v ", it.Key(), it.Value())
	}
	return strings.TrimRight(str, " ") + "]"

}

// Iterator holding the iterator's state
type Iterator struct {
	iterator rbt.Iterator
}

// Iterator returns a stateful iterator whose elements are key/value pairs.
func (m *Map) Iterator() Iterator {
	return Iterator{iterator: m.tree.Iterator()}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (iterator *Iterator) Next() bool {
	return iterator.iterator.Next()
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator) Prev() bool {
	return iterator.iterator.Prev()
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator *Iterator) Value() V {
	return iterator.iterator.Value().(V)
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
func (iterator *Iterator) Key() K {
	return iterator.iterator.Key().(K)
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (iterator *Iterator) Begin() {
	iterator.iterator.Begin()
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (iterator *Iterator) End() {
	iterator.iterator.End()
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, then first element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator
func (iterator *Iterator) First() bool {
	return iterator.iterator.First()
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator) Last() bool {
	return iterator.iterator.Last()
}

// ToJSON outputs the JSON representation of the map.
type pair struct {
	Key K
	Val V
}

func (m *Map) MarshalJSON() ([]byte, error) {
	elements := make([]pair, 0, m.Size())
	it := m.Iterator()
	for it.Next() {
		elements = append(elements, pair{it.Key(), it.Value()})
	}
	return json.Marshal(&elements)
}

// FromJSON populates the map from the input JSON representation.
func (m *Map) UnmarshalJSON(data []byte) error {
	elements := make([]pair, 0)
	err := json.Unmarshal(data, &elements)
	if err == nil {
		m.Clear()
		for _, pair := range elements {
			m.Put(pair.Key, pair.Val)
		}
	}
	return err
}

type MultiMap struct {
	Map
}

func NewMulti() *MultiMap {
	return &MultiMap{Map{tree: rbt.NewWith(Compare), isMulti: true}}
}

func (m *MultiMap) Get(key K) (front, end Iterator) {
	lower, upper := m.tree.MultiGet(key)
	return Iterator{lower}, Iterator{upper}
}

func (m *MultiMap) LowerBound(key K) *Iterator {
	if itr := m.tree.LowerBound(key); itr != m.tree.End() {
		return &Iterator{itr}
	}
	return nil
}

func (m *MultiMap) UpperBound(key K) *Iterator {
	if itr := m.tree.UpperBound(key); itr != m.tree.End() {
		return &Iterator{itr}
	}
	return nil
}
