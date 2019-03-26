// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redblacktree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestRedBlackTreeMultiPut(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite

	if actualValue := tree.Size(); actualValue != 8 {
		t.Errorf("Got %v expected %v", actualValue, 8)
	}
	if actualValue, expectedValue := fmt.Sprintf("%d%d%d%d%d%d%d%d", tree.Keys()...), "11234567"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s%s%s%s%s%s%s%s", tree.Values()...), "xabcdefg"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	tests1 := [][]interface{}{
		{1, "x", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, "e", true},
		{6, "f", true},
		{7, "g", true},
		{8, nil, false},
	}

	for _, test := range tests1 {
		// retrievals
		actualValue := tree.Get(test[0])
		if actualValue.HasNext() != test[2] || (actualValue.HasNext() && actualValue.Value() != test[1]) {
			t.Errorf("Got %v expected %v", actualValue.Value(), test[1])
		}
	}
}

func TestRedBlackTreeMultiRemove(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite

	tree.Remove(5)
	tree.Remove(6)
	tree.Remove(7)
	tree.Remove(8)
	tree.Remove(5)

	if actualValue, expectedValue := fmt.Sprintf("%d%d%d%d%d", tree.Keys()...), "11234"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s%s%s%s%s", tree.Values()...), "xabcd"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue := tree.Size(); actualValue != 5 {
		t.Errorf("Got %v expected %v", actualValue, 5)
	}

	tests2 := [][]interface{}{
		{1, "x", true},
		{2, "b", true},
		{3, "c", true},
		{4, "d", true},
		{5, nil, false},
		{6, nil, false},
		{7, nil, false},
		{8, nil, false},
	}

	for _, test := range tests2 {
		actualValue := tree.Get(test[0])
		if actualValue.HasNext() != test[2] || (actualValue.HasNext() && actualValue.Value() != test[1]) {
			t.Errorf("Got %v expected %v", actualValue.Value(), test[1])
		}
	}

	tree.Remove(1)
	tree.Remove(4)
	tree.Remove(2)
	tree.Remove(3)
	tree.Remove(2)
	tree.Remove(2)

	if actualValue, expectedValue := fmt.Sprintf("%s", tree.Keys()), "[]"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s", tree.Values()), "[]"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if empty, size := tree.Empty(), tree.Size(); empty != true || size != -0 {
		t.Errorf("Got %v expected %v", empty, true)
	}

}

func TestRedBlackMultiTreeLeftAndRight(t *testing.T) {
	tree := NewWithIntComparator(true)

	if actualValue := tree.Left(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := tree.Right(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	tree.Put(1, "a")
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x") // overwrite
	tree.Put(2, "b")

	if actualValue, expectedValue := fmt.Sprintf("%d", tree.Left().Key), "1"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s", tree.Left().Value), "a"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}

	if actualValue, expectedValue := fmt.Sprintf("%d", tree.Right().Key), "7"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s", tree.Right().Value), "g"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestRedBlackMultiTreeCeilingAndFloor(t *testing.T) {
	tree := NewWithIntComparator(true)

	if node, found := tree.Floor(0); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
	if node, found := tree.Ceiling(0); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(1, "x")
	tree.Put(2, "b")

	if node, found := tree.Floor(4); node.Key != 4 || !found {
		t.Errorf("Got %v expected %v", node.Key, 4)
	}
	if node, found := tree.Floor(0); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	if node, found := tree.Ceiling(4); node.Key != 4 || !found {
		t.Errorf("Got %v expected %v", node.Key, 4)
	}
	if node, found := tree.Ceiling(8); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
}

func TestRedBlackMultiTreeIteratorNextOnEmpty(t *testing.T) {
	tree := NewWithIntComparator(true)
	it := tree.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty tree")
	}
}

func TestRedBlackMultiTreeIteratorPrevOnEmpty(t *testing.T) {
	tree := NewWithIntComparator(true)
	it := tree.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty tree")
	}
}

func TestRedBlackMultiTreeIterator1Next(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite
	// │   ┌── 7
	// └── 6
	//     │   ┌── 5
	//     └── 4
	//         │   ┌── 3
	//         └── 2
	//             └── 1
	it := tree.Iterator()
	count := 0
	for it.Next() {
		count++
		key := it.Key()
		switch key {
		case count:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestRedBlackMultiTreeIterator1Prev(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(5, "e")
	tree.Put(6, "f")
	tree.Put(7, "g")
	tree.Put(3, "c")
	tree.Put(4, "d")
	tree.Put(2, "b")
	tree.Put(1, "a") //overwrite
	// │   ┌── 7
	// └── 6
	//     │   ┌── 5
	//     └── 4
	//         │   ┌── 3
	//         └── 2
	//             └── 1
	it := tree.Iterator()
	for it.Next() {
	}
	countDown := tree.size
	for it.Prev() {
		key := it.Key()
		switch key {
		case countDown:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
		countDown--
	}
	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestRedBlackMultiTreeIterator2Next(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it := tree.Iterator()
	count := 0
	for it.Next() {
		count++
		key := it.Key()
		switch key {
		case count:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestRedBlackMultiTreeIterator2Prev(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it := tree.Iterator()
	for it.Next() {
	}
	countDown := tree.size
	for it.Prev() {
		key := it.Key()
		switch key {
		case countDown:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
		countDown--
	}
	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestRedBlackMultiTreeIterator3Next(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(1, "a")
	it := tree.Iterator()
	count := 0
	for it.Next() {
		count++
		key := it.Key()
		switch key {
		case count:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestRedBlackMultiTreeIterator3Prev(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(1, "a")
	it := tree.Iterator()
	for it.Next() {
	}
	countDown := tree.size
	for it.Prev() {
		key := it.Key()
		switch key {
		case countDown:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := key, countDown; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
		countDown--
	}
	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestRedBlackMultiTreeIterator4Next(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(13, 5)
	tree.Put(8, 3)
	tree.Put(17, 7)
	tree.Put(1, 1)
	tree.Put(11, 4)
	tree.Put(15, 6)
	tree.Put(25, 9)
	tree.Put(6, 2)
	tree.Put(22, 8)
	tree.Put(27, 10)
	// │           ┌── 27
	// │       ┌── 25
	// │       │   └── 22
	// │   ┌── 17
	// │   │   └── 15
	// └── 13
	//     │   ┌── 11
	//     └── 8
	//         │   ┌── 6
	//         └── 1
	it := tree.Iterator()
	count := 0
	for it.Next() {
		count++
		value := it.Value()
		switch value {
		case count:
			if actualValue, expectedValue := value, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := value, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
	}
	if actualValue, expectedValue := count, tree.Size(); actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestRedBlackMultiTreeIterator4Prev(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(13, 5)
	tree.Put(8, 3)
	tree.Put(17, 7)
	tree.Put(1, 1)
	tree.Put(11, 4)
	tree.Put(15, 6)
	tree.Put(25, 9)
	tree.Put(6, 2)
	tree.Put(22, 8)
	tree.Put(27, 10)
	// │           ┌── 27
	// │       ┌── 25
	// │       │   └── 22
	// │   ┌── 17
	// │   │   └── 15
	// └── 13
	//     │   ┌── 11
	//     └── 8
	//         │   ┌── 6
	//         └── 1
	it := tree.Iterator()
	count := tree.Size()
	for it.Next() {
	}
	for it.Prev() {
		value := it.Value()
		switch value {
		case count:
			if actualValue, expectedValue := value, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			if actualValue, expectedValue := value, count; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		}
		count--
	}
	if actualValue, expectedValue := count, 0; actualValue != expectedValue {
		t.Errorf("Size different. Got %v expected %v", actualValue, expectedValue)
	}
}

func TestRedBlackMultiTreeMultiRemove(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(1, 11)
	tree.Put(1, 12)
	tree.Put(1, 13)
	tree.Put(2, 21)
	tree.Put(2, 22)
	tree.Put(3, 31)
	tree.Put(3, 32)
	tree.Put(5, 5)
	tree.Put(4, 41)
	tree.Put(4, 42)

	tree.Remove(1)
	assert.Equal(t, []int{21, 22, 31, 32, 41, 42, 5}, tree.Values())
	tree.Remove(3)
	assert.Equal(t, []int{21, 22, 41, 42, 5}, tree.Values())
	tree.Remove(5)
	assert.Equal(t, []int{21, 22, 41, 42}, tree.Values())
}

func TestRedBlackTreeLowerUpperBound(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(1, 1)
	tree.Put(2, 2)
	tree.Put(3, 31)
	tree.Put(3, 32)
	tree.Put(3, 33)
	tree.Put(3, 34)
	tree.Put(4, 41)
	tree.Put(4, 42)
	tree.Put(5, 5)
	tree.Put(6, 61)
	tree.Put(6, 62)
	tree.Put(8, 8)

	assert := func(actual, expect int) {
		if actual != expect {
			t.Fatalf("got %d, but expected %d", actual, expect)
		}
	}

	lower := tree.LowerBound(3)
	upper := tree.UpperBound(3)

	assert(31, lower.node.Value.(int))
	assert(41, upper.node.Value.(int))

	lower = tree.LowerBound(6)
	upper = tree.UpperBound(6)

	assert(61, lower.node.Value.(int))
	//if upper != tree.End() {
	//	t.Fatal("upper expected error")
	//}

	lower = tree.LowerBound(1)
	upper = tree.UpperBound(1)

	assert(1, lower.node.Value.(int))
	assert(2, upper.node.Value.(int))

	lower = tree.LowerBound(0)
	upper = tree.UpperBound(0)

	fmt.Println(lower.Value())
	fmt.Println(upper.Value())

	lower = tree.LowerBound(7)
	upper = tree.UpperBound(7)

	fmt.Println(lower.Value())
	fmt.Println(upper.Value())

	lower = tree.LowerBound(8)
	upper = tree.UpperBound(8)

	fmt.Println(lower.Value())
	fmt.Println(upper.Value())
}

func TestRedBlackMultiTreeIteratorDelete(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(1, 11)
	tree.Put(1, 12)
	tree.Put(1, 13)
	tree.Put(2, 21)
	tree.Put(2, 22)
	tree.Put(3, 31)
	tree.Put(3, 32)
	tree.Put(5, 5)
	tree.Put(4, 41)
	tree.Put(4, 42)

	expects := []int{11, 12, 13, 21, 22, 31, 32, 41, 42, 5}
	index := 1

	for itr := tree.Begin(); itr.HasNext(); itr.Next() {
		tree.RemoveOne(itr)
		assert.Equal(t, expects[index:], tree.Values())
		index++
	}
}

func TestRedBlackMultiTreeIteratorBegin(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it := tree.Iterator()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	it.Begin()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	for it.Next() {
	}

	it.Begin()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	it.Next()
	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
	}
}

func TestRedBlackMultiTreeIteratorEnd(t *testing.T) {
	tree := NewWithIntComparator(true)
	it := tree.Iterator()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	it.End()
	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it.End()
	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	it.Prev()
	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
	}
}

func TestRedBlackMultiTreeIteratorFirst(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it := tree.Iterator()
	if actualValue, expectedValue := it.First(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
	}
}

func TestRedBlackMultiTreeIteratorLast(t *testing.T) {
	tree := NewWithIntComparator(true)
	tree.Put(3, "c")
	tree.Put(1, "a")
	tree.Put(2, "b")
	it := tree.Iterator()
	if actualValue, expectedValue := it.Last(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
	}
}

func benchmarkMultiGet(b *testing.B, tree *Tree, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Get(n)
		}
	}
}

func benchmarkMultiPut(b *testing.B, tree *Tree, size int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tree.Clear()
		b.StartTimer()
		for n := 0; n < size; n++ {
			tree.Put(n, struct{}{})
		}
	}
}

func benchmarkMultiRemove(b *testing.B, tree *Tree, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.Remove(n)
		}
	}
}

func BenchmarkRedBlackTreeMultiGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewWithIntComparator(true)
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiGet(b, tree, size)
}

func BenchmarkRedBlackTreeMultiGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewWithIntComparator(true)
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiGet(b, tree, size)
}

func BenchmarkRedBlackTreeMultiGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewWithIntComparator(true)
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiGet(b, tree, size)
}

func BenchmarkRedBlackTreeMultiGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	tree := NewWithIntComparator(true)
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiGet(b, tree, size)
}

func BenchmarkRedBlackTreeMultiPut100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewWithIntComparator(true)
	b.StartTimer()
	benchmarkMultiPut(b, tree, size)
}

func BenchmarkRedBlackTreeMultiPut1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewWithIntComparator(true)
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiPut(b, tree, size)
}

func BenchmarkRedBlackTreeMultiPut10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewWithIntComparator(true)
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiPut(b, tree, size)
}

func BenchmarkRedBlackTreeMultiPut100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	tree := NewWithIntComparator(true)
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiPut(b, tree, size)
}

func BenchmarkRedBlackTreeMultiRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewWithIntComparator(true)
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiRemove(b, tree, size)
}

func BenchmarkRedBlackTreeMultiRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewWithIntComparator(true)
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiRemove(b, tree, size)
}

func BenchmarkRedBlackTreeMultiRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewWithIntComparator(true)
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiRemove(b, tree, size)
}

func BenchmarkRedBlackTreeMultiRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	tree := NewWithIntComparator(true)
	for n := 0; n < size; n++ {
		tree.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiRemove(b, tree, size)
}

func BenchmarkMultiIterator_Next(b *testing.B) {
	b.StopTimer()
	size := 100000
	tree := NewWithIntComparator(true)
	for n := 0; n < size; n++ {
		tree.Put(rand.Int(), struct{}{})
	}

	itr := tree.Iterator()
	b.StartTimer()
	for itr.Next() {

	}
	benchmarkMultiRemove(b, tree, size)
}
