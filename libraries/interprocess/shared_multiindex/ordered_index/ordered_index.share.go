package ordered_index

import (
	"fmt"
	"github.com/eosspark/eos-go/libraries/container"
	"github.com/eosspark/eos-go/libraries/interprocess/allocator"
	. "github.com/eosspark/eos-go/libraries/interprocess/offsetptr"
	"github.com/eosspark/eos-go/libraries/multiindex"
	"unsafe"
)

// template type OrderedIndex(FinalIndex,FinalNode,SuperIndex,SuperNode,Value,Key,KeyFunc,Comparator,Multiply,Allocator)
type Value = int
type Key = int

var KeyFunc = func(v Value) Key { return v }
var Comparator = func(a, b Key) int { return 0 }
var Allocator allocator.MemoryManager = nil

const Multiply = false

// OrderedIndex holds elements of the red-black tree
type OrderedIndex struct {
	super Pointer `*SuperIndex` // index on the OrderedIndex, IndexBase is the last super index
	final Pointer `*FinalIndex` // index under the OrderedIndex, MultiIndex is the final index

	Root Pointer `*OrderedIndexNode`
	size int
}

func (tree *OrderedIndex) init(final *FinalIndex) {
	tree.Root.Set(nil)
	tree.size = 0
	tree.final.Set(unsafe.Pointer(final))
	//tree.final = final
	tree.super.Set(unsafe.Pointer(NewSuperIndex()))
	//tree.super = NewSuperIndex()
	(*SuperIndex)(tree.super.Get()).init(final)
	//tree.super.init(final)
}

func (tree *OrderedIndex) free() {
	(*SuperIndex)(tree.super.Get()).free()
}

func (tree *OrderedIndex) clear() {
	tree.Clear()
	(*SuperIndex)(tree.super.Get()).clear()
}

/*generic class*/
type SuperIndex struct {
	init    func(*FinalIndex)
	free    func()
	clear   func()
	insert  func(Value, *FinalNode) (*SuperNode, bool)
	erase   func(*SuperNode) bool
	erase_  func(multiindex.IteratorType)
	modify  func(*SuperNode) (*SuperNode, bool)
	modify_ func(multiindex.IteratorType, func(*Value)) bool
}

const _SizeofSuperIndex = unsafe.Sizeof(SuperIndex{})

func NewSuperIndex() *SuperIndex {
	if Allocator == nil {
		return &SuperIndex{}
	}
	return (*SuperIndex)(Allocator.Allocate(_SizeofSuperIndex))
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
	super  Pointer `*SuperNode`
	final  Pointer `*FinalNode`
	color  color
	Left   Pointer `*OrderedIndexNode`
	Right  Pointer `*OrderedIndexNode`
	Parent Pointer `*OrderedIndexNode`
}

type np_ = *OrderedIndexNode

const _SizeofOrderedIndexNode = unsafe.Sizeof(OrderedIndexNode{})

func NewOrderedIndexNode(key Key, color color) (n *OrderedIndexNode) {
	n = np_(Allocator.Allocate(_SizeofOrderedIndexNode))
	n.Key = key
	n.color = color
	n.super.Set(nil)
	n.final.Set(nil)
	n.Left.Set(nil)
	n.Right.Set(nil)
	n.Parent.Set(nil)

	return n
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

func (node *OrderedIndexNode) free() {
	if node != nil {
		Allocator.DeAllocate(unsafe.Pointer(node))
	}
	// else free by golang gc
}

func (node *OrderedIndexNode) value() *Value {
	return (*SuperNode)(node.super.Get()).value()
	//return node.super.value()
}

type color bool

const (
	black, red color = true, false
)

func (tree *OrderedIndex) Insert(v Value) (Iterator, bool) {
	fn, res := (*FinalIndex)(tree.final.Get()).insert(v)
	//fn, res := tree.final.insert(v)
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
	sn, res := (*SuperIndex)(tree.super.Get()).insert(v, fn)
	//sn, res := tree.super.insert(v, fn)
	if res {
		node.super.Set(unsafe.Pointer(sn))
		//node.super = sn
		node.final.Set(unsafe.Pointer(fn))
		//node.final = fn
		return node, true
	}
	tree.remove(node)
	return nil, false
}

func (tree *OrderedIndex) Erase(iter Iterator) (itr Iterator) {
	itr = iter
	itr.Next()
	(*FinalIndex)(tree.final.Get()).erase((*FinalNode)(iter.node.final.Get()))
	//tree.final.erase(iter.node.final)
	return
}

func (tree *OrderedIndex) Erases(first, last Iterator) {
	for first != last {
		first = tree.Erase(first)
	}
}

func (tree *OrderedIndex) erase(n *OrderedIndexNode) {
	tree.remove(n)
	(*SuperIndex)(tree.super.Get()).erase((*SuperNode)(n.super.Get()))
	//tree.super.erase(n.super)
	n.free()
	//n.super = nil
	//n.final = nil
}

func (tree *OrderedIndex) erase_(iter multiindex.IteratorType) {
	if itr, ok := iter.(Iterator); ok {
		tree.Erase(itr)
	} else {
		(*SuperIndex)(tree.super.Get()).erase_(iter)
		//tree.super.erase_(iter)
	}
}

func (tree *OrderedIndex) Modify(iter Iterator, mod func(*Value)) bool {
	if _, b := (*FinalIndex)(tree.final.Get()).modify(mod, (*FinalNode)(iter.node.final.Get())); b {
		//if _, b := tree.final.modify(mod, iter.node.final); b {
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
			(*SuperIndex)(tree.super.Get()).erase((*SuperNode)(n.super.Get()))
			//tree.super.erase(n.super)
			return nil, false
		}

		node.super.Forward(&n.super)
		node.final.Forward(&n.final)

		n.free()
		n = node
	}

	if sn, res := (*SuperIndex)(tree.super.Get()).modify((*SuperNode)(n.super.Get())); !res {
		//if sn, res := tree.super.modify(n.super); !res {
		tree.remove(n)
		return nil, false
	} else {
		n.super.Set(unsafe.Pointer(sn))
		//n.super = sn
	}

	return n, true
}

func (tree *OrderedIndex) modify_(iter multiindex.IteratorType, mod func(*Value)) bool {
	if itr, ok := iter.(Iterator); ok {
		return tree.Modify(itr, mod)
	} else {
		return (*SuperIndex)(tree.super.Get()).modify_(iter, mod)
		//return tree.super.modify_(iter, mod)
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
	node := np_(tree.Root.Get())

	if node == nil {
		return result
	}

	for {
		if Comparator(key, node.Key) > 0 {
			if !node.Right.IsNil() {
				//if node.Right != nil {
				node = np_(node.Right.Get())
				//node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			//result.node = node
			result.position = between
			if !node.Left.IsNil() {
				//if node.Left != nil {
				node = np_(node.Left.Get())
				//node = node.Left
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
	node := np_(tree.Root.Get())

	if node == nil {
		return result
	}

	for {
		if Comparator(key, node.Key) >= 0 {
			if !node.Right.IsNil() {
				//if node.Right != nil {
				node = np_(node.Right.Get())
				//node = node.Right
			} else {
				return result
			}
		} else {
			result.node = node
			//result.node = node
			result.position = between
			if !node.Left.IsNil() {
				//if node.Left != nil {
				node = np_(node.Left.Get())
				//node = node.Left
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
	if tree.Root.IsNil() {
		//if tree.Root == nil {
		// Assert key is of comparator's type for initial tree
		Comparator(key, key)
		tree.Root.Set(unsafe.Pointer(NewOrderedIndexNode(key, red)))
		//tree.Root = &OrderedIndexNode{Key: key, color: red}
		insertedNode = np_(tree.Root.Get())
	} else {
		node := np_(tree.Root.Get())
		loop := true
		if Multiply {
			for loop {
				compare := Comparator(key, node.Key)
				//compare := Comparator(key, node.Key)
				switch {
				case compare < 0:
					if node.Left.IsNil() {
						//if node.Left == nil {
						node.Left.Set(unsafe.Pointer(NewOrderedIndexNode(key, red)))
						//node.Left = NewOrderedIndexNode(key, red)
						insertedNode = np_(node.Left.Get())
						//insertedNode = node.Left
						loop = false
					} else {
						node = np_(node.Left.Get())
						//node = node.Left
					}
				case compare >= 0:
					if node.Right.IsNil() {
						//if node.Right == nil {
						node.Right.Set(unsafe.Pointer(NewOrderedIndexNode(key, red)))
						//node.Right = NewOrderedIndexNode(key, red)
						insertedNode = np_(node.Right.Get())
						//insertedNode = node.Right
						loop = false
					} else {
						node = np_(node.Right.Get())
						//node = node.Right
					}
				}
			}
		} else {
			for loop {
				compare := Comparator(key, node.Key)
				//compare := Comparator(key, node.Key)
				switch {
				case compare == 0:
					node.Key = key
					//node.Key = key
					return node, false
					//return node, false
				case compare < 0:
					if node.Left.IsNil() {
						//if node.Left == nil {
						node.Left.Set(unsafe.Pointer(NewOrderedIndexNode(key, red)))
						//node.Left = NewOrderedIndexNode(key, red)
						insertedNode = np_(node.Left.Get())
						//insertedNode = node.Left
						loop = false
					} else {
						node = np_(node.Left.Get())
						//node = node.Left
					}
				case compare > 0:
					if node.Right.IsNil() {
						//if node.Right == nil {
						node.Right.Set(unsafe.Pointer(NewOrderedIndexNode(key, red)))
						//node.Right = NewOrderedIndexNode(key, red)
						insertedNode = np_(node.Right.Get())
						//insertedNode = node.Right
						loop = false
					} else {
						node = np_(node.Right.Get())
						//node = node.Right
					}
				}
			}
		}
		insertedNode.Parent.Set(unsafe.Pointer(node))
		//insertedNode.Parent = node
	}
	tree.insertCase1(insertedNode)
	tree.size++

	return insertedNode, true
}

func (tree *OrderedIndex) swapNode(node *OrderedIndexNode, pred *OrderedIndexNode) {
	if node == pred {
		return
	}

	tmp := OrderedIndexNode{color: pred.color}
	tmp.Left.Forward(&pred.Left)
	tmp.Right.Forward(&pred.Right)
	tmp.Parent.Forward(&pred.Parent)
	//tmp := OrderedIndexNode{color: pred.color, Left: pred.Left, Right: pred.Right, Parent: pred.Parent}

	pred.color = node.color
	node.color = tmp.color

	pred.Right.Forward(&node.Right)
	//pred.Right = node.Right
	if !pred.Right.IsNil() {
		//if pred.Right != nil {
		np_(pred.Right.Get()).Parent.Set(unsafe.Pointer(pred))
		//pred.Right.Parent = pred
	}
	node.Right.Forward(&tmp.Right)
	//node.Right = tmp.Right
	if !node.Right.IsNil() {
		//if node.Right != nil {
		np_(pred.Right.Get()).Parent.Set(unsafe.Pointer(node))
		//node.Right.Parent = node
	}

	if np_(pred.Parent.Get()) == node {
		//if pred.Parent == node {
		pred.Left.Set(unsafe.Pointer(node))
		//pred.Left = node
		node.Left.Forward(&tmp.Left)
		//node.Left = tmp.Left
		if !node.Left.IsNil() {
			//if node.Left != nil {
			np_(node.Left.Get()).Parent.Set(unsafe.Pointer(node))
			//node.Left.Parent = node
		}

		pred.Parent.Forward(&node.Parent)
		//pred.Parent = node.Parent
		if !pred.Parent.IsNil() {
			//if pred.Parent != nil {
			if np_(np_(pred.Parent.Get()).Left.Get()) == node {
				//if pred.Parent.Left == node {
				np_(pred.Parent.Get()).Left.Set(unsafe.Pointer(pred))
				//pred.Parent.Left = pred
			} else {
				np_(pred.Parent.Get()).Right.Set(unsafe.Pointer(pred))
				//pred.Parent.Right = pred
			}
		} else {
			tree.Root.Set(unsafe.Pointer(pred))
			//tree.Root = pred
		}
		node.Parent.Set(unsafe.Pointer(pred))
		//node.Parent = pred

	} else {
		pred.Left.Forward(&node.Left)
		//pred.Left = node.Left
		if !pred.Left.IsNil() {
			//if pred.Left != nil {
			np_(pred.Left.Get()).Parent.Set(unsafe.Pointer(pred))
			//pred.Left.Parent = pred
		}
		node.Left.Forward(&tmp.Left)
		//node.Left = tmp.Left
		if !node.Left.IsNil() {
			//if node.Left != nil {
			np_(node.Left.Get()).Parent.Set(unsafe.Pointer(node))
			//node.Left.Parent = node
		}

		pred.Parent.Forward(&node.Parent)
		//pred.Parent = node.Parent
		if !pred.Parent.IsNil() {
			if np_(np_(pred.Parent.Get()).Left.Get()) == node {
				//if pred.Parent.Left == node {
				np_(pred.Parent.Get()).Left.Set(unsafe.Pointer(pred))
				//pred.Parent.Left = pred
			} else {
				np_(pred.Parent.Get()).Right.Set(unsafe.Pointer(pred))
				//pred.Parent.Right = pred
			}
		} else {
			tree.Root.Set(unsafe.Pointer(pred))
			//tree.Root = pred
		}

		node.Parent.Forward(&tmp.Parent)
		//node.Parent = tmp.Parent
		if !node.Parent.IsNil() {
			//if node.Parent != nil {
			if np_(np_(node.Parent.Get()).Left.Get()) == pred {
				//if node.Parent.Left == pred {
				np_(node.Parent.Get()).Left.Set(unsafe.Pointer(node))
				//node.Parent.Left = node
			} else {
				np_(node.Parent.Get()).Right.Set(unsafe.Pointer(node))
				//node.Parent.Right = node
			}
		} else {
			tree.Root.Set(unsafe.Pointer(node))
			//tree.Root = node
		}
	}
}

func (tree *OrderedIndex) remove(node *OrderedIndexNode) {
	var child *OrderedIndexNode
	if node == nil {
		return
	}
	if !node.Left.IsNil() && !node.Right.IsNil() {
		//if node.Left != nil && node.Right != nil {
		pred := np_(node.Left.Get()).maximumNode()
		//pred := node.Left.maximumNode()
		tree.swapNode(node, pred)
	}
	if node.Left.IsNil() || node.Right.IsNil() {
		//if node.Left == nil || node.Right == nil {
		if node.Right.IsNil() {
			//if node.Right == nil {
			child = np_(node.Left.Get())
			//child = node.Left
		} else {
			child = np_(node.Right.Get())
			//child = node.Right
		}
		if node.color == black {
			node.color = nodeColor(child)
			tree.deleteCase1(node)
		}
		tree.replaceNode(node, child)
		if node.Parent.IsNil() && child != nil {
			//if node.Parent == nil && child != nil {
			child.color = black
		}
	}
	tree.size--
}

func (tree *OrderedIndex) lookup(key Key) *OrderedIndexNode {
	node := np_(tree.Root.Get())
	//node := tree.Root
	for node != nil {
		compare := Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = np_(node.Left.Get())
			//node = node.Left
		case compare > 0:
			node = np_(node.Right.Get())
			//node = node.Right
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
	current := np_(tree.Root.Get())
	//current := tree.Root
	for current != nil {
		parent = current
		current = np_(current.Left.Get())
		//current = current.Left
	}
	return parent
}

// Right returns the right-most (max) node or nil if tree is empty.
func (tree *OrderedIndex) Right() *OrderedIndexNode {
	var parent *OrderedIndexNode
	current := np_(tree.Root.Get())
	//current := tree.Root
	for current != nil {
		parent = current
		current = np_(current.Right.Get())
	}
	return parent
}

// Clear removes all nodes from the tree.
func (tree *OrderedIndex) Clear() {
	if Allocator != nil {
		//TODO DeAllocator
	}
	tree.Root.Set(nil)
	//tree.Root = nil
	tree.size = 0
}

// String returns a string representation of container
func (tree *OrderedIndex) String() string {
	str := "OrderedIndex\n"
	if !tree.Empty() {
		output(np_(tree.Root.Get()), "", true, &str)
		//output(tree.Root, "", true, &str)
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
	if !node.Right.IsNil() {
		//if node.Right != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "│   "
		} else {
			newPrefix += "    "
		}
		output(np_(node.Right.Get()), newPrefix, false, str)
		//output(node.Right, newPrefix, false, str)
	}
	*str += prefix
	if isTail {
		*str += "└── "
	} else {
		*str += "┌── "
	}
	*str += node.String() + "\n"
	if !node.Left.IsNil() {
		//if node.Left != nil {
		newPrefix := prefix
		if isTail {
			newPrefix += "    "
		} else {
			newPrefix += "│   "
		}
		output(np_(node.Left.Get()), newPrefix, true, str)
		//output(node.Left, newPrefix, true, str)
	}
}

func (node *OrderedIndexNode) grandparent() *OrderedIndexNode {
	if node != nil && !node.Parent.IsNil() {
		//if node != nil && node.Parent != nil {
		return np_(np_(node.Parent.Get()).Parent.Get())
		//return node.Parent.Parent
	}
	return nil
}

func (node *OrderedIndexNode) uncle() *OrderedIndexNode {
	if node == nil || node.Parent.IsNil() || np_(node.Parent.Get()).Parent.IsNil() {
		//if node == nil || node.Parent == nil || node.Parent.Parent == nil {
		return nil
	}
	return np_(node.Parent.Get()).sibling()
	//return node.Parent.sibling()
}

func (node *OrderedIndexNode) sibling() *OrderedIndexNode {
	if node == nil || node.Parent.IsNil() {
		//if node == nil || node.Parent == nil {
		return nil
	}
	if node == np_(np_(node.Parent.Get()).Left.Get()) {
		//if node == node.Parent.Left {
		return np_(np_(node.Parent.Get()).Right.Get())
		//return node.Parent.Get().Right
	}
	return np_(np_(node.Parent.Get()).Left.Get())
	//return node.Parent.Left
}

func (node *OrderedIndexNode) isLeaf() bool {
	if node == nil {
		return true
	}
	if node.Right.IsNil() && node.Left.IsNil() {
		//if node.Right == nil && node.Left == nil {
		return true
	}
	return false
}

func (tree *OrderedIndex) rotateLeft(node *OrderedIndexNode) {
	right := np_(node.Right.Get())
	tree.replaceNode(node, right)
	node.Right.Forward(&right.Left)
	if !right.Left.IsNil() {
		//if right.Left != nil {
		np_(right.Left.Get()).Parent.Set(unsafe.Pointer(node))
		//right.Left.Parent = node
	}
	right.Left.Set(unsafe.Pointer(node))
	//right.Left = node
	node.Parent.Set(unsafe.Pointer(right))
	//node.Parent = right
}

func (tree *OrderedIndex) rotateRight(node *OrderedIndexNode) {
	left := np_(node.Left.Get())
	//left := node.Left
	tree.replaceNode(node, left)
	node.Left.Forward(&left.Right)
	if !left.Right.IsNil() {
		//if left.Right != nil {
		np_(left.Right.Get()).Parent.Set(unsafe.Pointer(node))
		//left.Right.Parent = node
	}
	left.Right.Set(unsafe.Pointer(node))
	//left.Right = node
	node.Parent.Set(unsafe.Pointer(left))
}

func (tree *OrderedIndex) replaceNode(old *OrderedIndexNode, new *OrderedIndexNode) {
	if old.Parent.IsNil() {
		//if old.Parent == nil {
		tree.Root.Set(unsafe.Pointer(new))
		//tree.Root = new
	} else {
		if old == np_(np_(old.Parent.Get()).Left.Get()) {
			//if old == old.Parent.Left {
			np_(old.Parent.Get()).Left.Set(unsafe.Pointer(new))
			//old.Parent.Left = new
		} else {
			np_(old.Parent.Get()).Right.Set(unsafe.Pointer(new))
			//old.Parent.Right = new
		}
	}
	if new != nil {
		new.Parent.Forward(&old.Parent)
	}
}

func (tree *OrderedIndex) insertCase1(node *OrderedIndexNode) {
	if node.Parent.IsNil() {
		//if node.Parent == nil {
		node.color = black
	} else {
		tree.insertCase2(node)
	}
}

func (tree *OrderedIndex) insertCase2(node *OrderedIndexNode) {
	if nodeColor(np_(node.Parent.Get())) == black {
		//if nodeColor(node.Parent) == black {
		return
	}
	tree.insertCase3(node)
}

func (tree *OrderedIndex) insertCase3(node *OrderedIndexNode) {
	uncle := node.uncle()
	if nodeColor(uncle) == red {
		np_(node.Parent.Get()).color = black
		//node.Parent.color = black
		uncle.color = black
		node.grandparent().color = red
		tree.insertCase1(node.grandparent())
	} else {
		tree.insertCase4(node)
	}
}

func (tree *OrderedIndex) insertCase4(node *OrderedIndexNode) {
	grandparent := node.grandparent()
	if node == np_(np_(node.Parent.Get()).Right.Get()) && node.Parent.Get() == grandparent.Left.Get() {
		//if node == node.Parent.Right && node.Parent == grandparent.Left {
		tree.rotateLeft(np_(node.Parent.Get()))
		//tree.rotateLeft(node.Parent)
		node = np_(node.Left.Get())
		//node = node.Left
	} else if node == np_(np_(node.Parent.Get()).Left.Get()) && node.Parent.Get() == grandparent.Right.Get() {
		//} else if node == node.Parent.Left && node.Parent == grandparent.Right {
		tree.rotateRight(np_(node.Parent.Get()))
		//tree.rotateRight(node.Parent)
		node = np_(node.Right.Get())
		//node = node.Right
	}
	tree.insertCase5(node)
}

func (tree *OrderedIndex) insertCase5(node *OrderedIndexNode) {
	np_(node.Parent.Get()).color = black
	//node.Parent.color = black
	grandparent := node.grandparent()
	grandparent.color = red
	if node == np_(np_(node.Parent.Get()).Left.Get()) && node.Parent.Get() == grandparent.Left.Get() {
		//if node == node.Parent.Left && node.Parent == grandparent.Left {
		tree.rotateRight(grandparent)
	} else if node == np_(np_(node.Parent.Get()).Right.Get()) && node.Parent.Get() == grandparent.Right.Get() {
		//} else if node == node.Parent.Right && node.Parent == grandparent.Right {
		tree.rotateLeft(grandparent)
	}
}

func (node *OrderedIndexNode) maximumNode() *OrderedIndexNode {
	if node == nil {
		return nil
	}
	for !node.Right.IsNil() {
		//for node.Right != nil {
		node = np_(node.Right.Get())
		//node = node.Right
	}
	return node
}

func (tree *OrderedIndex) deleteCase1(node *OrderedIndexNode) {
	if node.Parent.IsNil() {
		//if node.Parent == nil {
		return
	}
	tree.deleteCase2(node)
}

func (tree *OrderedIndex) deleteCase2(node *OrderedIndexNode) {
	sibling := node.sibling()
	if nodeColor(sibling) == red {
		np_(node.Parent.Get()).color = red
		//node.Parent.color = red
		sibling.color = black
		if node == np_(np_(node.Parent.Get()).Left.Get()) {
			//if node == node.Parent.Left {
			tree.rotateLeft(np_(node.Parent.Get()))
			//tree.rotateLeft(node.Parent)
		} else {
			tree.rotateRight(np_(node.Parent.Get()))
			//tree.rotateRight(node.Parent)
		}
	}
	tree.deleteCase3(node)
}

func (tree *OrderedIndex) deleteCase3(node *OrderedIndexNode) {
	sibling := node.sibling()
	if nodeColor(np_(node.Parent.Get())) == black &&
	//if nodeColor(node.Parent) == black &&
		nodeColor(sibling) == black &&
	//nodeColor(sibling) == black &&
		nodeColor(np_(sibling.Left.Get())) == black &&
	//nodeColor(sibling.Left) == black &&
		nodeColor(np_(sibling.Right.Get())) == black {
		//nodeColor(sibling.Right) == black {
		sibling.color = red
		tree.deleteCase1(np_(node.Parent.Get()))
		//tree.deleteCase1(node.Parent)
	} else {
		tree.deleteCase4(node)
	}
}

func (tree *OrderedIndex) deleteCase4(node *OrderedIndexNode) {
	sibling := node.sibling()
	if nodeColor(np_(node.Parent.Get())) == red &&
	//if nodeColor(node.Parent) == red &&
		nodeColor(sibling) == black &&
		nodeColor(np_(sibling.Left.Get())) == black &&
	//nodeColor(sibling.Left) == black &&
		nodeColor(np_(sibling.Right.Get())) == black {
		//nodeColor(sibling.Right) == black {
		sibling.color = red
		np_(node.Parent.Get()).color = black
		//node.Parent.color = black
	} else {
		tree.deleteCase5(node)
	}
}

func (tree *OrderedIndex) deleteCase5(node *OrderedIndexNode) {
	sibling := node.sibling()
	if node == np_(np_(node.Parent.Get()).Left.Get()) &&
	//if node == node.Parent.Left &&
		nodeColor(sibling) == black &&
		nodeColor(np_(sibling.Left.Get())) == red &&
	//nodeColor(sibling.Left) == red &&
		nodeColor(np_(sibling.Right.Get())) == black {
		//nodeColor(sibling.Right) == black {
		sibling.color = red
		np_(sibling.Left.Get()).color = black
		//sibling.Left.color = black
		tree.rotateRight(sibling)
	} else if node == np_(np_(node.Parent.Get()).Right.Get()) &&
	//} else if node == node.Parent.Right &&
		nodeColor(sibling) == black &&
		nodeColor(np_(sibling.Right.Get())) == red &&
	//nodeColor(sibling.Right) == red &&
		nodeColor(np_(sibling.Left.Get())) == black {
		//nodeColor(sibling.Left) == black {
		sibling.color = red
		np_(sibling.Right.Get()).color = black
		//sibling.Right.color = black
		tree.rotateLeft(sibling)
	}
	tree.deleteCase6(node)
}

func (tree *OrderedIndex) deleteCase6(node *OrderedIndexNode) {
	sibling := node.sibling()
	sibling.color = nodeColor(np_(node.Parent.Get()))
	//sibling.color = nodeColor(node.Parent)
	np_(node.Parent.Get()).color = black
	//node.Parent.color = black
	if node == np_(np_(node.Parent.Get()).Left.Get()) && nodeColor(np_(sibling.Right.Get())) == red {
		//if node == node.Parent.Left && nodeColor(sibling.Right) == red {
		np_(sibling.Right.Get()).color = black
		//sibling.Right.color = black
		tree.rotateLeft(np_(node.Parent.Get()))
		//tree.rotateLeft(node.Parent)
	} else if nodeColor(np_(sibling.Left.Get())) == red {
		//} else if nodeColor(sibling.Left) == red {
		np_(sibling.Left.Get()).color = black
		//sibling.Left.color = black
		tree.rotateRight(np_(node.Parent.Get()))
		//tree.rotateRight(node.Parent)
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

		} else if n, ok := node.(np_); ok {
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
	if !iterator.node.Right.IsNil() {
		//if iterator.node.Right != nil {
		iterator.node = np_(iterator.node.Right.Get())
		//iterator.node = iterator.node.Right
		for !iterator.node.Left.IsNil() {
			//for iterator.node.Left != nil {
			iterator.node = np_(iterator.node.Left.Get())
			//iterator.node = iterator.node.Left
		}
		goto between
	}
	if !iterator.node.Parent.IsNil() {
		//if iterator.node.Parent != nil {
		node := iterator.node
		for !iterator.node.Parent.IsNil() {
			//for iterator.node.Parent != nil {
			iterator.node = np_(iterator.node.Parent.Get())
			//iterator.node = iterator.node.Parent
			if node == np_(iterator.node.Left.Get()) {
				//if node == iterator.node.Left {
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
	if !iterator.node.Left.IsNil() {
		//if iterator.node.Left != nil {
		iterator.node = np_(iterator.node.Left.Get())
		//iterator.node = iterator.node.Left
		for !iterator.node.Right.IsNil() {
			//for iterator.node.Right != nil {
			iterator.node = np_(iterator.node.Right.Get())
			//iterator.node = iterator.node.Right
		}
		goto between
	}
	if !iterator.node.Parent.IsNil() {
		//if iterator.node.Parent != nil {
		node := iterator.node
		for !iterator.node.Parent.IsNil() {
			//for iterator.node.Parent != nil {
			iterator.node = np_(iterator.node.Parent.Get())
			//iterator.node = iterator.node.Parent
			if node == np_(iterator.node.Right.Get()) {
				//if node == iterator.node.Right {
				goto between
			}
			node = iterator.node
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
