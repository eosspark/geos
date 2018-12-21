// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package treeset

import (
	"github.com/eosspark/container/containers"
	rbt "github.com/eosspark/container/trees/redblacktree"
)

func assertMultiSetIteratorImplementation() {
	var _ containers.ReverseIteratorWithIndex = (*Iterator)(nil)
}

// Iterator returns a stateful iterator whose values can be fetched by an index.
type MultiSetIterator struct {
	//index    int
	iterator rbt.Iterator
	tree     *rbt.Tree
}

// Iterator holding the iterator's state
func (set *MultiSet) Iterator() MultiSetIterator {
	//return MultiSetIterator{index: -1, iterator: set.tree.Iterator(), tree: set.tree}
	return MultiSetIterator{iterator: set.tree.Iterator(), tree: set.tree}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's index and value can be retrieved by Index() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (iterator *MultiSetIterator) Next() bool {
	//if iterator.index < iterator.tree.Size() {
	//	iterator.index++
	//}
	return iterator.iterator.Next()
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *MultiSetIterator) Prev() bool {
	//if iterator.index >= 0 {
	//	iterator.index--
	//}
	return iterator.iterator.Prev()
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator *MultiSetIterator) Value() interface{} {
	return iterator.iterator.Key()
}

// Index returns the current element's index.
// Does not modify the state of the iterator.
//func (iterator *MultiSetIterator) Index() int {
//	return iterator.index
//}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (iterator *MultiSetIterator) Begin() {
	//iterator.index = -1
	iterator.iterator.Begin()
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (iterator *MultiSetIterator) End() {
	//iterator.index = iterator.tree.Size()
	iterator.iterator.End()
}

// First moves the iterator to the first element and returns true if there was a first element in the container.
// If First() returns true, then first element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *MultiSetIterator) First() bool {
	iterator.Begin()
	return iterator.Next()
}

// Last moves the iterator to the last element and returns true if there was a last element in the container.
// If Last() returns true, then last element's index and value can be retrieved by Index() and Value().
// Modifies the state of the iterator.
func (iterator *MultiSetIterator) Last() bool {
	iterator.End()
	return iterator.Prev()
}