// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package redblacktree implements a red-black tree.
//
// Used by TreeSet and TreeMap.
//
// Structure is not thread safe.
//
// References: http://en.wikipedia.org/wiki/Red%E2%80%93black_tree
package redblacktree

import (
	"fmt"
	"github.com/eosspark/eos-go/libraries/container"
)

type color bool

const (
	black, red color = true, false
)

// Tree holds elements of the red-black tree
type Tree struct {
	isMulti    bool
	Root       *Node
	size       int
	Comparator func(a, b interface{}) int
}

// Node is a single element within the tree
type Node struct {
	Key    interface{}
	Value  interface{}
	color  color
	Left   *Node
	Right  *Node
	Parent *Node
}

// NewWith instantiates a red-black tree with the custom comparator.
func NewWith(comparator func(a, b interface{}) int, isMulti bool) *Tree {
	return &Tree{Comparator: comparator, isMulti: isMulti}
}

// NewWithIntComparator instantiates a red-black tree with the IntComparator, i.e. keys are of type int.
func NewWithIntComparator(isMulti bool) *Tree {
	return &Tree{Comparator: func(a, b interface{}) int {
		return container.IntComparator(a.(int), b.(int))
	}, isMulti: isMulti}
}

// NewWithStringComparator instantiates a red-black tree with the StringComparator, i.e. keys are of type string.
func NewWithStringComparator(isMulti bool) *Tree {
	return &Tree{Comparator: func(a, b interface{}) int {
		return container.StringComparator(a.(string), b.(string))
	}, isMulti: isMulti}
}

func CopyFrom(rbt *Tree) *Tree {
	t := &Tree{Comparator: rbt.Comparator, isMulti: rbt.isMulti}
	t.size = rbt.size
	t.Root = copyNode(rbt.Root)
	return t
}

func copyNode(nd *Node) *Node {
	if nd == nil {
		return nil
	}

	n := &Node{Key: nd.Key, Value: nd.Value, color: nd.color}

	if nd.Left != nil {
		n.Left = copyNode(nd.Left)
	}
	if nd.Right != nil {
		n.Right = copyNode(nd.Right)
	}

	if n.Left != nil {
		n.Left.Parent = n
	}
	if n.Right != nil {
		n.Right.Parent = n
	}
	if nd.Parent != nil {
		n.Parent = nd.Parent
	}

	return n
}

func (tree *Tree) New(comparator func(a, b interface{}) int) {
	tree.Comparator = comparator
}

func (tree *Tree) CopyFrom(rbt *Tree) {
	tree.Comparator = rbt.Comparator
	tree.size = rbt.size
	tree.Root = copyNode(rbt.Root)
}

// Put inserts node into the tree.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Put(key interface{}, value interface{}) {
	tree.put(key, value)
}

func (tree *Tree) put(key interface{}, value interface{}) (*Node, bool) {
	var insertedNode *Node
	if tree.Root == nil {
		// Assert key is of comparator's type for initial tree
		tree.Comparator(key, key)
		tree.Root = &Node{Key: key, Value: value, color: red}
		insertedNode = tree.Root
	} else {
		node := tree.Root
		loop := true
		if tree.isMulti {
			for loop {
				compare := tree.Comparator(key, node.Key)
				switch {
				case compare < 0:
					if node.Left == nil {
						node.Left = &Node{Key: key, Value: value, color: red}
						insertedNode = node.Left
						loop = false
					} else {
						node = node.Left
					}
				case compare >= 0:
					if node.Right == nil {
						node.Right = &Node{Key: key, Value: value, color: red}
						insertedNode = node.Right
						loop = false
					} else {
						node = node.Right
					}
				}
			}
		} else {
			for loop {
				compare := tree.Comparator(key, node.Key)
				switch {
				case compare == 0:
					node.Key = key
					node.Value = value
					return node, false
				case compare < 0:
					if node.Left == nil {
						node.Left = &Node{Key: key, Value: value, color: red}
						insertedNode = node.Left
						loop = false
					} else {
						node = node.Left
					}
				case compare > 0:
					if node.Right == nil {
						node.Right = &Node{Key: key, Value: value, color: red}
						insertedNode = node.Right
						loop = false
					} else {
						node = node.Right
					}
				}
			}
		}
		insertedNode.Parent = node
	}
	tree.insertCase1(insertedNode)
	tree.size++

	return insertedNode, true
}

func (tree *Tree) Insert(key interface{}, value interface{}) Iterator {
	node, new := tree.put(key, value)
	if new {
		return Iterator{tree: tree, node: node, position: between}
	}
	return tree.End()
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Get(key interface{}) Iterator {
	if tree.isMulti {
		lower := tree.LowerBound(key)
		if !lower.IsEnd() && tree.Comparator(key, lower.Key()) == 0 {
			return lower
		}
		return tree.End()
	} else {
		if node := tree.lookup(key); node != nil {
			return Iterator{tree, node, between}
		}
		return tree.End()
	}
}

// LowerBound returns an iterator pointing to the first element that is not less than the given key.
// Complexity: O(log N).
func (tree *Tree) LowerBound(key interface{}) Iterator {
	result := tree.End()
	node := tree.Root

	if node == nil {
		return result
	}

	for {
		if tree.Comparator(key, node.Key) > 0 {
			if node.Right != nil {
				node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			result.position = between
			if node.Left != nil {
				node = node.Left
			} else {
				return result
			}
		}
	}
}

// UpperBound returns an iterator pointing to the first element that is greater than the given key.
// Complexity: O(log N).
func (tree *Tree) UpperBound(key interface{}) Iterator {
	result := tree.End()
	node := tree.Root

	if node == nil {
		return result
	}

	for {
		if tree.Comparator(key, node.Key) >= 0 {
			if node.Right != nil {
				node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			result.position = between
			if node.Left != nil {
				node = node.Left
			} else {
				return result
			}
		}
	}
}

// Remove remove the node from the tree by key.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Remove(key interface{}) {
	if tree.isMulti {
		for lower := tree.LowerBound(key); lower.position != end; {
			if tree.Comparator(lower.Key(), key) == 0 {
				node := lower.node
				lower.Next()
				tree.remove(node)
			} else {
				break
			}
		}
	} else {
		node := tree.lookup(key)
		tree.remove(node)
	}
}

func (tree *Tree) RemoveOne(itr Iterator) {
	tree.remove(itr.node)
}

func (tree *Tree) swapNode(node *Node, pred *Node) {
	if node == pred {
		return
	}

	tmp := Node{pred.Key, pred.Value, pred.color,
		pred.Left, pred.Right, pred.Parent}

	pred.color = node.color
	node.color = tmp.color

	pred.Right = node.Right
	if pred.Right != nil {
		pred.Right.Parent = pred
	}
	node.Right = tmp.Right
	if node.Right != nil {
		node.Right.Parent = node
	}

	if pred.Parent == node {
		pred.Left = node
		node.Left = tmp.Left
		if node.Left != nil {
			node.Left.Parent = node
		}

		pred.Parent = node.Parent
		if pred.Parent != nil {
			if pred.Parent.Left == node {
				pred.Parent.Left = pred
			} else {
				pred.Parent.Right = pred
			}
		} else {
			tree.Root = pred
		}
		node.Parent = pred

	} else {
		pred.Left = node.Left
		if pred.Left != nil {
			pred.Left.Parent = pred
		}
		node.Left = tmp.Left
		if node.Left != nil {
			node.Left.Parent = node
		}

		pred.Parent = node.Parent
		if pred.Parent != nil {
			if pred.Parent.Left == node {
				pred.Parent.Left = pred
			} else {
				pred.Parent.Right = pred
			}
		} else {
			tree.Root = pred
		}

		node.Parent = tmp.Parent
		if node.Parent != nil {
			if node.Parent.Left == pred {
				node.Parent.Left = node
			} else {
				node.Parent.Right = node
			}
		} else {
			tree.Root = node
		}
	}
}

func (tree *Tree) remove(node *Node) {
	var child *Node
	if node == nil {
		return
	}
	if node.Left != nil && node.Right != nil {
		pred := node.Left.maximumNode()
		tree.swapNode(node, pred)
	}
	if node.Left == nil || node.Right == nil {
		if node.Right == nil {
			child = node.Left
		} else {
			child = node.Right
		}
		if node.color == black {
			node.color = nodeColor(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.Parent == nil && child != nil {
			child.color = black
		}
	}
	tree.size--
}

//func (tree *Tree) remove(node *Node) {
//	var child *Node
//	if node == nil {
//		return
//	}
//	if node.Left != nil && node.Right != nil {
//		pred := node.Left.maximumNode()
//		node.Key = pred.Key
//		node.Value = pred.Value
//		node = pred
//	}
//	if node.Left == nil || node.Right == nil {
//		if node.Right == nil {
//			child = node.Left
//		} else {
//			child = node.Right
//		}
//		if node.color == black {
//			node.color = nodeColor(child)
//			tree.deleteCase1(node)
//		}
//		tree.replaceNode(node, child)
//		if node.Parent == nil && child != nil {
//			child.color = black
//		}
//	}
//	tree.size--
//}

// Empty returns true if tree does not contain any nodes
func (tree *Tree) Empty() bool {
	return tree.size == 0
}

// Size returns number of nodes in the tree.
func (tree *Tree) Size() int {
	return tree.size
}

// Keys returns all keys in-order
func (tree *Tree) Keys() []interface{} {
	keys := make([]interface{}, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

// Values returns all values in-order based on the key.
func (tree *Tree) Values() []interface{} {
	values := make([]interface{}, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}
	return values
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *Tree) Left() *Node {
	var parent *Node
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Left
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *Tree) Right() *Node {
	var parent *Node
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Right
	}
	return parent
}

// Floor Finds floor node of the input key, return the floor node or nil if no floor is found.
// Second return parameter is true if floor was found, otherwise false.
//
// Floor node is defined as the largest node that is smaller than or equal to the given node.
// A floor node may not be found, either because the tree is empty, or because
// all nodes in the tree are larger than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Floor(key interface{}) (floor *Node, found bool) {
	found = false
	node := tree.Root
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			node = node.Left
		case compare > 0:
			floor, found = node, true
			node = node.Right
		}
	}
	if found {
		return floor, true
	}
	return nil, false
}

// Ceiling finds ceiling node of the input key, return the ceiling node or nil if no ceiling is found.
// Second return parameter is true if ceiling was found, otherwise false.
//
// Ceiling node is defined as the smallest node that is larger than or equal to the given node.
// A ceiling node may not be found, either because the tree is empty, or because
// all nodes in the tree are smaller than the given node.
//
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *Tree) Ceiling(key interface{}) (ceiling *Node, found bool) {
	found = false
	node := tree.Root
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node, true
		case compare < 0:
			ceiling, found = node, true
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	if found {
		return ceiling, true
	}
	return nil, false
}

// Clear removes all nodes from the tree.
func (tree *Tree) Clear() {
	tree.Root = nil
	tree.size = 0
}

// String returns a string representation of container
func (tree *Tree) String() string {
	str := "RedBlackTree\n"
	if !tree.Empty() {
		output(tree.Root, "", true, &str)
	}
	return str
}

func (node *Node) String() string {
	if !node.color {
		return fmt.Sprintf("(%v:%v,%v)", node.Key, node.Value, "red")
	}
	return fmt.Sprintf("(%v,%v)", node.Key, node.Value)
}

func output(node *Node, prefix string, isTail bool, str *string) {
	if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(node.Right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if node.Left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(node.Left, newPrefix, true, str)
	}
}

func (tree *Tree) lookup(key interface{}) *Node {
	node := tree.Root
	for node != nil {
		compare := tree.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return nil
}

func (node *Node) grandparent() *Node {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

func (node *Node) uncle() *Node {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return node.Parent.sibling()
}

func (node *Node) sibling() *Node {
	if node == nil || node.Parent == nil {
		return nil
	}
	if node == node.Parent.Left {
		return node.Parent.Right
	}
	return node.Parent.Left
}

func (node *Node) isLeaf() bool {
	if node == nil {
		return true
	}
	if node.Right == nil && node.Left == nil {
		return true
	}
	return false
}

func (tree *Tree) rotateLeft(node *Node) {
	right := node.Right
	tree.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (tree *Tree) rotateRight(node *Node) {
	left := node.Left
	tree.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Right = node
	node.Parent = left
}

func (tree *Tree) replaceNode(old *Node, new *Node) {
	if old.Parent == nil {
		tree.Root = new
	} else {
		if old == old.Parent.Left {
			old.Parent.Left = new
		} else {
			old.Parent.Right = new
		}
	}
	if new != nil {
		new.Parent = old.Parent
	}
}

func (tree *Tree) insertCase1(node *Node) {
	if node.Parent == nil {
		node.color = black
	} else {
		tree.insertCase2(node)
	}
}

func (tree *Tree) insertCase2(node *Node) {
	if nodeColor(node.Parent) == black {
		return
	}
	tree.insertCase3(node)
}

func (tree *Tree) insertCase3(node *Node) {
	uncle := node.uncle()
	if nodeColor(uncle) == red {
		node.Parent.color = black
		uncle.color = black
		node.grandparent().color = red
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *Tree) insertCase4(node *Node) {
	grandparent := node.grandparent()
	if node == node.Parent.Right && node.Parent == grandparent.Left {
		tree.rotateLeft(node.Parent)
		node = node.Left
	} else if node == node.Parent.Left && node.Parent == grandparent.Right {
		tree.rotateRight(node.Parent)
		node = node.Right
	}
	tree.insertCase5(node)
}

func (tree *Tree) insertCase5(node *Node) {
	node.Parent.color = black
	grandparent := node.grandparent()
	grandparent.color = red
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *Node) maximumNode() *Node {
	if node == nil {
		return nil
	}
	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (tree *Tree) deleteCase1(node *Node) {
	if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *Tree) deleteCase2(node *Node) {
	sibling := node.sibling()
	if nodeColor(sibling) == red {
		node.Parent.color = red
		sibling.color = black
		if node == node.Parent.Left {
			tree.rotateLeft(node.Parent)
		} else {
			tree.rotateRight(node.Parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *Tree) deleteCase3(node *Node) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == black &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *Tree) deleteCase4(node *Node) {
	sibling := node.sibling()
	if nodeColor(node.Parent) == red &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == black &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		node.Parent.color = black
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *Tree) deleteCase5(node *Node) {
	sibling := node.sibling()
	if node == node.Parent.Left &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Left) == red &&
		nodeColor(sibling.Right) == black {
		sibling.color = red
		sibling.Left.color = black
		tree.rotateRight(sibling)
	} else if node == node.Parent.Right &&
		nodeColor(sibling) == black &&
		nodeColor(sibling.Right) == red &&
		nodeColor(sibling.Left) == black {
		sibling.color = red
		sibling.Right.color = black
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *Tree) deleteCase6(node *Node) {
	sibling := node.sibling()
	sibling.color = nodeColor(node.Parent)
	node.Parent.color = black
	if node == node.Parent.Left && nodeColor(sibling.Right) == red {
		sibling.Right.color = black
		tree.rotateLeft(node.Parent)
	} else if nodeColor(sibling.Left) == red {
		sibling.Left.color = black
		tree.rotateRight(node.Parent)
	}
}

func nodeColor(node *Node) color {
	if node == nil {
		return black
	}
	return node.color
}
