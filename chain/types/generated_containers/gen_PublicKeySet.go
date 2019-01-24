// Code generated by gotemplate. DO NOT EDIT.

// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package treeset implements a Tree backed by a red-black Tree.
//
// Structure is not thread safe.
//
// Reference: http://en.wikipedia.org/wiki/Set_%28abstract_data_type%29
package generated

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/container"
	rbt "github.com/eosspark/eos-go/common/container/tree"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
)

// template type Set(V,Compare,Multi)

func assertPublicKeySetImplementation() {
	var _ container.Set = (*PublicKeySet)(nil)
}

// Set holds elements in a red-black Tree
type PublicKeySet struct {
	*rbt.Tree
}

var itemExistsPublicKeySet = struct{}{}

// NewWith instantiates a new empty set with the custom comparator.

func NewPublicKeySet(Value ...ecc.PublicKey) *PublicKeySet {
	set := &PublicKeySet{Tree: rbt.NewWith(ecc.ComparePubKey, false)}
	set.Add(Value...)
	return set
}

func CopyFromPublicKeySet(ts *PublicKeySet) *PublicKeySet {
	return &PublicKeySet{Tree: rbt.CopyFrom(ts.Tree)}
}

func PublicKeySetIntersection(a *PublicKeySet, b *PublicKeySet, callback func(elem ecc.PublicKey)) {
	aIterator := a.Iterator()
	bIterator := b.Iterator()

	if !aIterator.First() || !bIterator.First() {
		return
	}

	for aHasNext, bHasNext := true, true; aHasNext && bHasNext; {
		comp := ecc.ComparePubKey(aIterator.Value(), bIterator.Value())
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
func (set *PublicKeySet) AddItem(item ecc.PublicKey) (bool, ecc.PublicKey) {
	itr := set.Tree.Insert(item, itemExistsPublicKeySet)
	if itr.IsEnd() {
		return false, item
	}
	return true, itr.Key().(ecc.PublicKey)
}

// Add adds the items (one or more) to the set.
func (set *PublicKeySet) Add(items ...ecc.PublicKey) {
	for _, item := range items {
		set.Tree.Put(item, itemExistsPublicKeySet)
	}
}

// Remove removes the items (one or more) from the set.
func (set *PublicKeySet) Remove(items ...ecc.PublicKey) {
	for _, item := range items {
		set.Tree.Remove(item)
	}

}

// Values returns all items in the set.
func (set *PublicKeySet) Values() []ecc.PublicKey {
	keys := make([]ecc.PublicKey, set.Size())
	it := set.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Value()
	}
	return keys
}

// Contains checks weather items (one or more) are present in the set.
// All items have to be present in the set for the method to return true.
// Returns true if no arguments are passed at all, i.e. set is always superset of empty set.
func (set *PublicKeySet) Contains(items ...ecc.PublicKey) bool {
	for _, item := range items {
		if iter := set.Get(item); iter.IsEnd() {
			return false
		}
	}
	return true
}

// String returns a string representation of container
func (set *PublicKeySet) String() string {
	str := "TreeSet\n"
	items := make([]string, 0)
	for _, v := range set.Tree.Keys() {
		items = append(items, fmt.Sprintf("%v", v))
	}
	str += strings.Join(items, ", ")
	return str
}

// Iterator returns a stateful iterator whose values can be fetched by an index.
type IteratorPublicKeySet struct {
	rbt.Iterator
}

// Iterator holding the iterator's state
func (set *PublicKeySet) Iterator() IteratorPublicKeySet {
	return IteratorPublicKeySet{Iterator: set.Tree.Iterator()}
}

// Begin returns First Iterator whose position points to the first element
// Return End Iterator when the map is empty
func (set *PublicKeySet) Begin() IteratorPublicKeySet {
	return IteratorPublicKeySet{set.Tree.Begin()}
}

// End returns End Iterator
func (set *PublicKeySet) End() IteratorPublicKeySet {
	return IteratorPublicKeySet{set.Tree.End()}
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator IteratorPublicKeySet) Value() ecc.PublicKey {
	return iterator.Iterator.Key().(ecc.PublicKey)
}

// Each calls the given function once for each element, passing that element's index and value.
func (set *PublicKeySet) Each(f func(value ecc.PublicKey)) {
	iterator := set.Iterator()
	for iterator.Next() {
		f(iterator.Value())
	}
}

// Find passes each element of the container to the given function and returns
// the first (index,value) for which the function is true or -1,nil otherwise
// if no element matches the criteria.
func (set *PublicKeySet) Find(f func(value ecc.PublicKey) bool) (v ecc.PublicKey) {
	iterator := set.Iterator()
	for iterator.Next() {
		if f(iterator.Value()) {
			return iterator.Value()
		}
	}
	return
}

func (set *PublicKeySet) LowerBound(item ecc.PublicKey) IteratorPublicKeySet {
	return IteratorPublicKeySet{set.Tree.LowerBound(item)}
}

func (set *PublicKeySet) UpperBound(item ecc.PublicKey) IteratorPublicKeySet {
	return IteratorPublicKeySet{set.Tree.UpperBound(item)}
}

// ToJSON outputs the JSON representation of the set.
func (set PublicKeySet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.Values())
}

// FromJSON populates the set from the input JSON representation.
func (set *PublicKeySet) UnmarshalJSON(data []byte) error {
	elements := make([]ecc.PublicKey, 0)
	err := json.Unmarshal(data, &elements)
	if err == nil {
		set.Tree = rbt.NewWith(ecc.ComparePubKey, false)
		set.Add(elements...)
	}
	return err
}

func (set PublicKeySet) Pack() (re []byte, err error) {
	re = append(re, common.WriteUVarInt(set.Size())...)
	set.Each(func(value ecc.PublicKey) {
		reVal, _ := rlp.EncodeToBytes(value)
		re = append(re, reVal...)
	})
	return re, nil
}

func (set *PublicKeySet) Unpack(in []byte) (int, error) {
	set.Tree = rbt.NewWith(ecc.ComparePubKey, false)

	decoder := rlp.NewDecoder(in)
	l, err := decoder.ReadUvarint64()
	if err != nil {
		return 0, err
	}

	for i := 0; i < int(l); i++ {
		v := new(ecc.PublicKey)
		decoder.Decode(v)
		set.Add(*v)
	}
	return decoder.GetPos(), nil
}
