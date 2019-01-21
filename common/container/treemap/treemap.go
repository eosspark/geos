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
package treemap

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/container/utils"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/container"
	rbt "github.com/eosspark/eos-go/common/container/tree"
	"github.com/eosspark/eos-go/crypto/rlp"
	"strings"
)

// template type Map(K,V,Compare,Multi)
type K int
type V int

var Compare = utils.IntComparator
var Multi = false

func assertMapImplementation() {
	var _ container.Map = (*Map)(nil)
}

// Map holds the elements in a red-black Tree
type Map struct {
	*rbt.Tree
}

// NewWith instantiates a Tree map with the custom comparator.
func New() *Map {
	return &Map{Tree: rbt.NewWith(Compare, Multi)}
}

func CopyFrom(tm *Map) *Map {
	return &Map{Tree: rbt.CopyFrom(tm.Tree)}
}

// Put inserts key-value pair into the map.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map) Put(key K, value V) {
	m.Tree.Put(key, value)
}

func (m *Map) Insert(key K, value V) Iterator {
	return Iterator{m.Tree.Insert(key, value)}
}

// Get searches the element in the map by key and returns its value or nil if key is not found in Tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map) Get(key K) Iterator {
	return Iterator{m.Tree.Get(key)}
}

// Remove removes the element from the map by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (m *Map) Remove(key K) {
	m.Tree.Remove(key)
}

// Keys returns all keys in-order
func (m *Map) Keys() []K {
	keys := make([]K, m.Tree.Size())
	it := m.Tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key().(K)
	}
	return keys
}

// Values returns all values in-order based on the key.
func (m *Map) Values() []V {
	values := make([]V, m.Tree.Size())
	it := m.Tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value().(V)
	}
	return values
}

// Each calls the given function once for each element, passing that element's key and value.
func (m *Map) Each(f func(key K, value V)) {
	Iterator := m.Iterator()
	for Iterator.Next() {
		f(Iterator.Key(), Iterator.Value())
	}
}

// Find passes each element of the container to the given function and returns
// the first (key,value) for which the function is true or nil,nil otherwise if no element
// matches the criteria.
func (m *Map) Find(f func(key K, value V) bool) (k K, v V) {
	Iterator := m.Iterator()
	for Iterator.Next() {
		if f(Iterator.Key(), Iterator.Value()) {
			return Iterator.Key(), Iterator.Value()
		}
	}
	return
}

// String returns a string representation of container
func (m Map) String() string {
	str := "TreeMap\nmap["
	it := m.Iterator()
	for it.Next() {
		str += fmt.Sprintf("%v:%v ", it.Key(), it.Value())
	}
	return strings.TrimRight(str, " ") + "]"

}

// Iterator holding the Iterator's state
type Iterator struct {
	rbt.Iterator
}

// Iterator returns a stateful Iterator whose elements are key/value pairs.
func (m *Map) Iterator() Iterator {
	return Iterator{Iterator: m.Tree.Iterator()}
}

// Begin returns First Iterator whose position points to the first element
// Return End Iterator when the map is empty
func (m *Map) Begin() Iterator {
	return Iterator{m.Tree.Begin()}
}

// End returns End Iterator
func (m *Map) End() Iterator {
	return Iterator{m.Tree.End()}
}

// Value returns the current element's value.
// Does not modify the state of the Iterator.
func (iterator Iterator) Value() V {
	return iterator.Iterator.Value().(V)
}

// Key returns the current element's key.
// Does not modify the state of the Iterator.
func (iterator Iterator) Key() K {
	return iterator.Iterator.Key().(K)
}

func (m *Map) LowerBound(key K) Iterator {
	return Iterator{m.Tree.LowerBound(key)}
}

func (m *Map) UpperBound(key K) Iterator {
	return Iterator{m.Tree.UpperBound(key)}

}

// ToJSON outputs the JSON representation of the map.
type pair struct {
	Key K `json:"key"`
	Val V `json:"val"`
}

func (m Map) MarshalJSON() ([]byte, error) {
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

func (m Map) Pack() (re []byte, err error) {
	re = append(re, common.WriteUVarInt(m.Size())...)
	m.Each(func(key K, value V) {
		rekey, _ := rlp.EncodeToBytes(key)
		re = append(re, rekey...)
		reVal, _ := rlp.EncodeToBytes(value)
		re = append(re, reVal...)
	})
	return re, nil
}

func (m *Map) Unpack(in []byte) (int, error) {
	m.Tree = rbt.NewWith(Compare, Multi)

	decoder := rlp.NewDecoder(in)
	l, err := decoder.ReadUvarint64()
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(l); i++ {
		k, v := new(K), new(V)
		decoder.Decode(k)
		decoder.Decode(v)
		m.Put(*k, *v)
	}
	return decoder.GetPos(), nil
}
