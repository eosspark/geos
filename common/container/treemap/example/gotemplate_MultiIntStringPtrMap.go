// Code generated by gotemplate. DO NOT EDIT.

// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package treemap implements a map backed by red-black Tree.
//
// Elements are ordered by key in the map.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Associative_array
package example

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/container"
	rbt "github.com/eosspark/eos-go/common/container/tree"
	"github.com/eosspark/eos-go/crypto/rlp"
)

// template type Map(K,V,Compare,Multi)

func assertMultiIntStringPtrMapImplementation() {
	var _ container.Map = (*MultiIntStringPtrMap)(nil)
}

// Map holds the elements in a red-black Tree
type MultiIntStringPtrMap struct {
	*rbt.Tree
}

// NewWith instantiates a Tree map with the custom comparator.
func NewMultiIntStringPtrMap() *MultiIntStringPtrMap {
	return &MultiIntStringPtrMap{Tree: rbt.NewWith(IntComparator, true)}
}

func CopyFromMultiIntStringPtrMap(tm *MultiIntStringPtrMap) *MultiIntStringPtrMap {
	return &MultiIntStringPtrMap{Tree: rbt.CopyFrom(tm.Tree)}
}

// Put inserts key-value pair into the map.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *MultiIntStringPtrMap) Put(key int, value *string) {
	m.Tree.Put(key, value)
}

func (m *MultiIntStringPtrMap) Insert(key int, value *string) IteratorMultiIntStringPtrMap {
	return IteratorMultiIntStringPtrMap{m.Tree.Insert(key, value)}
}

// Get searches the element in the map by key and returns its value or nil if key is not found in Tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *MultiIntStringPtrMap) Get(key int) IteratorMultiIntStringPtrMap {
	return IteratorMultiIntStringPtrMap{m.Tree.Get(key)}
}

// Remove removes the element from the map by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *MultiIntStringPtrMap) Remove(key int) {
	m.Tree.Remove(key)
}

// Keys returns all keys in-order
func (m *MultiIntStringPtrMap) Keys() []int {
	keys := make([]int, m.Tree.Size())
	it := m.Tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key().(int)
	}
	return keys
}

// Values returns all values in-order based on the key.
func (m *MultiIntStringPtrMap) Values() []*string {
	values := make([]*string, m.Tree.Size())
	it := m.Tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value().(*string)
	}
	return values
}

// Each calls the given function once for each element, passing that element's key and value.
func (m *MultiIntStringPtrMap) Each(f func(key int, value *string)) {
	Iterator := m.Iterator()
	for Iterator.Next() {
		f(Iterator.Key(), Iterator.Value())
	}
}

// Find passes each element of the container to the given function and returns
// the first (key,value) for which the function is true or nil,nil otherwise if no element
// matches the criteria.
func (m *MultiIntStringPtrMap) Find(f func(key int, value *string) bool) (k int, v *string) {
	Iterator := m.Iterator()
	for Iterator.Next() {
		if f(Iterator.Key(), Iterator.Value()) {
			return Iterator.Key(), Iterator.Value()
		}
	}
	return
}

// String returns a string representation of container
func (m MultiIntStringPtrMap) String() string {
	str := "TreeMap\nmap["
	it := m.Iterator()
	for it.Next() {
		str += fmt.Sprintf("%v:%v ", it.Key(), it.Value())
	}
	return strings.TrimRight(str, " ") + "]"

}

// Iterator holding the Iterator's state
type IteratorMultiIntStringPtrMap struct {
	rbt.Iterator
}

// Iterator returns a stateful Iterator whose elements are key/value pairs.
func (m *MultiIntStringPtrMap) Iterator() IteratorMultiIntStringPtrMap {
	return IteratorMultiIntStringPtrMap{Iterator: m.Tree.Iterator()}
}

// Begin returns First Iterator whose position points to the first element
// Return End Iterator when the map is empty
func (m *MultiIntStringPtrMap) Begin() IteratorMultiIntStringPtrMap {
	return IteratorMultiIntStringPtrMap{m.Tree.Begin()}
}

// End returns End Iterator
func (m *MultiIntStringPtrMap) End() IteratorMultiIntStringPtrMap {
	return IteratorMultiIntStringPtrMap{m.Tree.End()}
}

// Value returns the current element's value.
// Does not modify the state of the Iterator.
func (iterator IteratorMultiIntStringPtrMap) Value() *string {
	return iterator.Iterator.Value().(*string)
}

// Key returns the current element's key.
// Does not modify the state of the Iterator.
func (iterator IteratorMultiIntStringPtrMap) Key() int {
	return iterator.Iterator.Key().(int)
}

func (m *MultiIntStringPtrMap) LowerBound(key int) IteratorMultiIntStringPtrMap {
	return IteratorMultiIntStringPtrMap{m.Tree.LowerBound(key)}
}

func (m *MultiIntStringPtrMap) UpperBound(key int) IteratorMultiIntStringPtrMap {
	return IteratorMultiIntStringPtrMap{m.Tree.UpperBound(key)}

}

// ToJSON outputs the JSON representation of the map.
type pairMultiIntStringPtrMap struct {
	Key int     `json:"key"`
	Val *string `json:"val"`
}

func (m MultiIntStringPtrMap) MarshalJSON() ([]byte, error) {
	elements := make([]pairMultiIntStringPtrMap, 0, m.Size())
	it := m.Iterator()
	for it.Next() {
		elements = append(elements, pairMultiIntStringPtrMap{it.Key(), it.Value()})
	}
	return json.Marshal(&elements)
}

// FromJSON populates the map from the input JSON representation.
func (m *MultiIntStringPtrMap) UnmarshalJSON(data []byte) error {
	elements := make([]pairMultiIntStringPtrMap, 0)
	err := json.Unmarshal(data, &elements)
	if err == nil {
		m.Clear()
		for _, pair := range elements {
			m.Put(pair.Key, pair.Val)
		}
	}
	return err
}

func (m MultiIntStringPtrMap) Pack() (re []byte, err error) {
	re = append(re, common.WriteUVarInt(m.Size())...)
	m.Each(func(key int, value *string) {
		rekey, _ := rlp.EncodeToBytes(key)
		re = append(re, rekey...)
		reVal, _ := rlp.EncodeToBytes(value)
		re = append(re, reVal...)
	})
	return re, nil
}

func (m *MultiIntStringPtrMap) Unpack(in []byte) (int, error) {
	m.Tree = rbt.NewWith(IntComparator, true)

	decoder := rlp.NewDecoder(in)
	l, err := decoder.ReadUvarint64()
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(l); i++ {
		k, v := new(int), new(*string)
		decoder.Decode(k)
		decoder.Decode(v)
		m.Put(*k, *v)
	}
	return decoder.GetPos(), nil
}
