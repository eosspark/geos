package ordered_index

import (
	"fmt"
	"github.com/eosspark/eos-go/libraries/container"
	"github.com/eosspark/eos-go/libraries/multiindex"
)

// template type OrderedIndex(FinalIndex,FinalNode,SuperIndex,SuperNode,Value,Key,KeyFunc,Comparator,Multiply)
type Value = int
type Key = int

var KeyFunc = func(v Value) Key { return v }
var Comparator = func(a, b Key) int { return 0 }

const Multiply = false

// OrderedIndex holds elements of the red-black tree
type OrderedIndex struct {
	super *SuperIndex // index on the OrderedIndex, IndexBase is the last super index
	final *FinalIndex // index under the OrderedIndex, MultiIndex is the final index

	Root *OrderedIndexNode
	size int
}

func (tree *OrderedIndex) init(final *FinalIndex) {
	tree.final = final
	tree.super = &SuperIndex{}
	tree.super.init(final)
}

func (tree *OrderedIndex) clear() {
	tree.Clear()
	tree.super.clear()
}

/*generic class*/
type SuperIndex struct {
	init    func(*FinalIndex)
	clear   func()
	insert  func(Value, *FinalNode) (*SuperNode, bool)
	erase   func(*SuperNode) bool
	erase_  func(multiindex.IteratorType)
	modify  func(*SuperNode) (*SuperNode, bool)
	modify_ func(multiindex.IteratorType, func(*Value)) bool
}

/*generic class*/
type FinalIndex struct {
	insert func(Value) (*FinalNode, bool)
	erase  func(*FinalNode)
	modify func(func(*Value), *FinalNode) (*FinalNode, bool)
}

// OrderedIndexNode is a single element within the tree
type OrderedIndexNode struct {
	Key    Key
	super  *SuperNode
	final  *FinalNode
	color  color
	Left   *OrderedIndexNode
	Right  *OrderedIndexNode
	Parent *OrderedIndexNode
}

/*generic class*/
type SuperNode struct {
	value func() *Value
}

/*generic class*/
type FinalNode struct {
	GetSuperNode func() interface{}
	GetFinalNode func() interface{}
}

func (node *OrderedIndexNode) value() *Value {
	return node.super.value()
}

type color bool

const (
	black, red color = true, false
)

func (tree *OrderedIndex) Insert(v Value) (Iterator, bool) {
	fn, res := tree.final.insert(v)
	if res {
		return tree.makeIterator(fn), true
	}
	return tree.End(), false
}

func (tree *OrderedIndex) insert(v Value, fn *FinalNode) (*OrderedIndexNode, bool) {
	key := KeyFunc(v)

	node, res := tree.put(key)
	if !res {
		container.Logger.Warn("#ordered index insert failed")
		return nil, false
	}
	sn, res := tree.super.insert(v, fn)
	if res {
		node.super = sn
		node.final = fn
		return node, true
	}
	tree.remove(node)
	return nil, false
}

func (tree *OrderedIndex) Erase(iter Iterator) (itr Iterator) {
	itr = iter
	itr.Next()
	tree.final.erase(iter.node.final)
	return
}

func (tree *OrderedIndex) Erases(first, last Iterator) {
	for first != last {
		first = tree.Erase(first)
	}
}

func (tree *OrderedIndex) erase(n *OrderedIndexNode) {
	tree.remove(n)
	tree.super.erase(n.super)
	n.super = nil
	n.final = nil
}

func (tree *OrderedIndex) erase_(iter multiindex.IteratorType) {
	if itr, ok := iter.(Iterator); ok {
		tree.Erase(itr)
	} else {
		tree.super.erase_(iter)
	}
}

func (tree *OrderedIndex) Modify(iter Iterator, mod func(*Value)) bool {
	if _, b := tree.final.modify(mod, iter.node.final); b {
		return true
	}
	return false
}

func (tree *OrderedIndex) modify(n *OrderedIndexNode) (*OrderedIndexNode, bool) {
	n.Key = KeyFunc(*n.value())

	if !tree.inPlace(n) {
		tree.remove(n)
		node, res := tree.put(n.Key)
		if !res {
			container.Logger.Warn("#ordered index modify failed")
			tree.super.erase(n.super)
			return nil, false
		}

		//n.Left = node.Left
		//if n.Left != nil {
		//	n.Left.Parent = n
		//}
		//n.Right = node.Right
		//if n.Right != nil {
		//	n.Right.Parent = n
		//}
		//n.Parent = node.Parent
		//if n.Parent != nil {
		//	if n.Parent.Left == node {
		//		n.Parent.Left = n
		//	} else {
		//		n.Parent.Right = n
		//	}
		//} else {
		//	tree.Root = n
		//}
		node.super = n.super
		node.final = n.final
		n = node
	}

	if sn, res := tree.super.modify(n.super); !res {
		tree.remove(n)
		return nil, false
	} else {
		n.super = sn
	}

	return n, true
}

func (tree *OrderedIndex) modify_(iter multiindex.IteratorType, mod func(*Value)) bool {
	if itr, ok := iter.(Iterator); ok {
		return tree.Modify(itr, mod)
	} else {
		return tree.super.modify_(iter, mod)
	}
}

// Get searches the node in the tree by key and returns its value or nil if key is not found in tree.
// Second return parameter is true if key was found, otherwise false.
// Key should adhere to the comparator's type assertion, otherwise method panics.
func (tree *OrderedIndex) Find(key Key) Iterator {
	if Multiply {
		lower := tree.LowerBound(key)
		if !lower.IsEnd() && Comparator(key, lower.Key()) == 0 {
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
func (tree *OrderedIndex) LowerBound(key Key) Iterator {
	result := tree.End()
	node := tree.Root

	if node == nil {
		return result
	}

	for {
		if Comparator(key, node.Key) > 0 {
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
func (tree *OrderedIndex) UpperBound(key Key) Iterator {
	result := tree.End()
	node := tree.Root

	if node == nil {
		return result
	}

	for {
		if Comparator(key, node.Key) >= 0 {
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
func (tree *OrderedIndex) Remove(key Key) {
	if Multiply {
		for lower := tree.LowerBound(key); lower.position != end; {
			if Comparator(lower.Key(), key) == 0 {
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

func (tree *OrderedIndex) put(key Key) (*OrderedIndexNode, bool) {
	var insertedNode *OrderedIndexNode
	if tree.Root == nil {
		// Assert key is of comparator's type for initial tree
		Comparator(key, key)
		tree.Root = &OrderedIndexNode{Key: key, color: red}
		insertedNode = tree.Root
	} else {
		node := tree.Root
		loop := true
		if Multiply {
			for loop {
				compare := Comparator(key, node.Key)
				switch {
				case compare < 0:
					if node.Left == nil {
						node.Left = &OrderedIndexNode{Key: key, color: red}
						insertedNode = node.Left
						loop = false
					} else {
						node = node.Left
					}
				case compare >= 0:
					if node.Right == nil {
						node.Right = &OrderedIndexNode{Key: key, color: red}
						insertedNode = node.Right
						loop = false
					} else {
						node = node.Right
					}
				}
			}
		} else {
			for loop {
				compare := Comparator(key, node.Key)
				switch {
				case compare == 0:
					node.Key = key
					return node, false
				case compare < 0:
					if node.Left == nil {
						node.Left = &OrderedIndexNode{Key: key, color: red}
						insertedNode = node.Left
						loop = false
					} else {
						node = node.Left
					}
				case compare > 0:
					if node.Right == nil {
						node.Right = &OrderedIndexNode{Key: key, color: red}
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

func (tree *OrderedIndex) swapNode(node *OrderedIndexNode, pred *OrderedIndexNode) {
	if node == pred {
		return
	}

	tmp := OrderedIndexNode{color: pred.color, Left: pred.Left, Right: pred.Right, Parent: pred.Parent}

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

func (tree *OrderedIndex) remove(node *OrderedIndexNode) {
	var child *OrderedIndexNode
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

func (tree *OrderedIndex) lookup(key Key) *OrderedIndexNode {
	node := tree.Root
	for node != nil {
		compare := Comparator(key, node.Key)
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

// Empty returns true if tree does not contain any nodes
func (tree *OrderedIndex) Empty() bool {
	return tree.size == 0
}

// Size returns number of nodes in the tree.
func (tree *OrderedIndex) Size() int {
	return tree.size
}

// Keys returns all keys in-order
func (tree *OrderedIndex) Keys() []Key {
	keys := make([]Key, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		keys[i] = it.Key()
	}
	return keys
}

// Values returns all values in-order based on the key.
func (tree *OrderedIndex) Values() []Value {
	values := make([]Value, tree.size)
	it := tree.Iterator()
	for i := 0; it.Next(); i++ {
		values[i] = it.Value()
	}
	return values
}

// Left returns the left-most (min) node or nil if tree is empty.
func (tree *OrderedIndex) Left() *OrderedIndexNode {
	var parent *OrderedIndexNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Left
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *OrderedIndex) Right() *OrderedIndexNode {
	var parent *OrderedIndexNode
	current := tree.Root
	for current != nil {
		parent = current
		current = current.Right
	}
	return parent
}

// Clear removes all nodes from the tree.
func (tree *OrderedIndex) Clear() {
	tree.Root = nil
	tree.size = 0
}

// String returns a string representation of container
func (tree *OrderedIndex) String() string {
	str := "OrderedIndex\n"
	if !tree.Empty() {
		output(tree.Root, "", true, &str)
	}
	return str
}

func (node *OrderedIndexNode) String() string {
	if !node.color {
		return fmt.Sprintf("(%v,%v)", node.Key, "red")
	}
	return fmt.Sprintf("(%v)", node.Key)
}

func output(node *OrderedIndexNode, prefix string, isTail bool, str *string) {
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

func (node *OrderedIndexNode) grandparent() *OrderedIndexNode {
	if node != nil && node.Parent != nil {
		return node.Parent.Parent
	}
	return nil
}

func (node *OrderedIndexNode) uncle() *OrderedIndexNode {
	if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return node.Parent.sibling()
}

func (node *OrderedIndexNode) sibling() *OrderedIndexNode {
	if node == nil || node.Parent == nil {
		return nil
	}
	if node == node.Parent.Left {
		return node.Parent.Right
	}
	return node.Parent.Left
}

func (node *OrderedIndexNode) isLeaf() bool {
	if node == nil {
		return true
	}
	if node.Right == nil && node.Left == nil {
		return true
	}
	return false
}

func (tree *OrderedIndex) rotateLeft(node *OrderedIndexNode) {
	right := node.Right
	tree.replaceNode(node, right)
	node.Right = right.Left
	if right.Left != nil {
		right.Left.Parent = node
	}
	right.Left = node
	node.Parent = right
}

func (tree *OrderedIndex) rotateRight(node *OrderedIndexNode) {
	left := node.Left
	tree.replaceNode(node, left)
	node.Left = left.Right
	if left.Right != nil {
		left.Right.Parent = node
	}
	left.Right = node
	node.Parent = left
}

func (tree *OrderedIndex) replaceNode(old *OrderedIndexNode, new *OrderedIndexNode) {
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

func (tree *OrderedIndex) insertCase1(node *OrderedIndexNode) {
	if node.Parent == nil {
		node.color = black
	} else {
		tree.insertCase2(node)
	}
}

func (tree *OrderedIndex) insertCase2(node *OrderedIndexNode) {
	if nodeColor(node.Parent) == black {
		return
	}
	tree.insertCase3(node)
}

func (tree *OrderedIndex) insertCase3(node *OrderedIndexNode) {
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

func (tree *OrderedIndex) insertCase4(node *OrderedIndexNode) {
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

func (tree *OrderedIndex) insertCase5(node *OrderedIndexNode) {
	node.Parent.color = black
	grandparent := node.grandparent()
	grandparent.color = red
	if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *OrderedIndexNode) maximumNode() *OrderedIndexNode {
	if node == nil {
		return nil
	}
	for node.Right != nil {
		node = node.Right
	}
	return node
}

func (tree *OrderedIndex) deleteCase1(node *OrderedIndexNode) {
	if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *OrderedIndex) deleteCase2(node *OrderedIndexNode) {
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

func (tree *OrderedIndex) deleteCase3(node *OrderedIndexNode) {
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

func (tree *OrderedIndex) deleteCase4(node *OrderedIndexNode) {
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

func (tree *OrderedIndex) deleteCase5(node *OrderedIndexNode) {
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

func (tree *OrderedIndex) deleteCase6(node *OrderedIndexNode) {
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

func nodeColor(node *OrderedIndexNode) color {
	if node == nil {
		return black
	}
	return node.color
}

//////////////iterator////////////////

func (tree *OrderedIndex) makeIterator(fn *FinalNode) Iterator {
	node := fn.GetSuperNode()
	for {
		if node == nil {
			panic("Wrong index node type!")

		} else if n, ok := node.(*OrderedIndexNode); ok {
			return Iterator{tree: tree, node: n, position: between}
		} else {
			node = node.(multiindex.NodeType).GetSuperNode()
		}
	}
}

// Iterator holding the iterator's state
type Iterator struct {
	tree     *OrderedIndex
	node     *OrderedIndexNode
	position position
}

type position byte

const (
	begin, between, end position = 0, 1, 2
)

// Iterator returns a stateful iterator whose elements are key/value pairs.
func (tree *OrderedIndex) Iterator() Iterator {
	return Iterator{tree: tree, node: nil, position: begin}
}

func (tree *OrderedIndex) Begin() Iterator {
	itr := Iterator{tree: tree, node: nil, position: begin}
	itr.Next()
	return itr
}

func (tree *OrderedIndex) End() Iterator {
	return Iterator{tree: tree, node: nil, position: end}
}

// Next moves the iterator to the next element and returns true if there was a next element in the container.
// If Next() returns true, then next element's key and value can be retrieved by Key() and Value().
// If Next() was called for the first time, then it will point the iterator to the first element if it exists.
// Modifies the state of the iterator.
func (iterator *Iterator) Next() bool {
	if iterator.position == end {
		goto end
	}
	if iterator.position == begin {
		left := iterator.tree.Left()
		if left == nil {
			goto end
		}
		iterator.node = left
		goto between
	}
	if iterator.node.Right != nil {
		iterator.node = iterator.node.Right
		for iterator.node.Left != nil {
			iterator.node = iterator.node.Left
		}
		goto between
	}
	if iterator.node.Parent != nil {
		node := iterator.node
		for iterator.node.Parent != nil {
			iterator.node = iterator.node.Parent
			if node == iterator.node.Left {
				goto between
			}
			node = iterator.node
		}
	}

end:
	iterator.node = nil
	iterator.position = end
	return false

between:
	iterator.position = between
	return true
}

// Prev moves the iterator to the previous element and returns true if there was a previous element in the container.
// If Prev() returns true, then previous element's key and value can be retrieved by Key() and Value().
// Modifies the state of the iterator.
func (iterator *Iterator) Prev() bool {
	if iterator.position == begin {
		goto begin
	}
	if iterator.position == end {
		right := iterator.tree.Right()
		if right == nil {
			goto begin
		}
		iterator.node = right
		goto between
	}
	if iterator.node.Left != nil {
		iterator.node = iterator.node.Left
		for iterator.node.Right != nil {
			iterator.node = iterator.node.Right
		}
		goto between
	}
	if iterator.node.Parent != nil {
		node := iterator.node
		for iterator.node.Parent != nil {
			iterator.node = iterator.node.Parent
			if node == iterator.node.Right {
				goto between
			}
			node = iterator.node
			//if iterator.tree.Comparator(node.Key, iterator.node.Key) >= 0 {
			//	goto between
			//}
		}
	}

begin:
	iterator.node = nil
	iterator.position = begin
	return false

between:
	iterator.position = between
	return true
}

func (iterator Iterator) HasNext() bool {
	return iterator.position != end
}

func (iterator *Iterator) HasPrev() bool {
	return iterator.position != begin
}

// Value returns the current element's value.
// Does not modify the state of the iterator.
func (iterator Iterator) Value() Value {
	return *iterator.node.value()
}

// Key returns the current element's key.
// Does not modify the state of the iterator.
func (iterator Iterator) Key() Key {
	return iterator.node.Key
}

// Begin resets the iterator to its initial state (one-before-first)
// Call Next() to fetch the first element if any.
func (iterator *Iterator) Begin() {
	iterator.node = nil
	iterator.position = begin
}

func (iterator Iterator) IsBegin() bool {
	return iterator.position == begin
}

// End moves the iterator past the last element (one-past-the-end).
// Call Prev() to fetch the last element if any.
func (iterator *Iterator) End() {
	iterator.node = nil
	iterator.position = end
}

func (iterator Iterator) IsEnd() bool {
	return iterator.position == end
}

// Delete remove the node which pointed by the iterator
// Modifies the state of the iterator.
func (iterator *Iterator) Delete() {
	node := iterator.node
	//iterator.Prev()
	iterator.tree.remove(node)
}

func (tree *OrderedIndex) inPlace(n *OrderedIndexNode) bool {
	prev := Iterator{tree, n, between}
	next := Iterator{tree, n, between}
	prev.Prev()
	next.Next()

	var (
		prevResult int
		nextResult int
	)

	if prev.IsBegin() {
		prevResult = 1
	} else {
		prevResult = Comparator(n.Key, prev.Key())
	}

	if next.IsEnd() {
		nextResult = -1
	} else {
		nextResult = Comparator(n.Key, next.Key())
	}

	return (Multiply && prevResult >= 0 && nextResult <= 0) ||
		(!Multiply && prevResult > 0 && nextResult < 0)
}
