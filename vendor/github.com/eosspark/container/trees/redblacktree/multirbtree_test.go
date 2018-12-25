// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package redblacktree

import (
	"fmt"
	"testing"
		"math/rand"
)

func TestRedBlackTreeMultiPut(t *testing.T) {
	tree := NewWithIntComparator()
	tree.MultiPut(5, "e")
	tree.MultiPut(6, "f")
	tree.MultiPut(7, "g")
	tree.MultiPut(3, "c")
	tree.MultiPut(4, "d")
	tree.MultiPut(1, "x")
	tree.MultiPut(2, "b")
	tree.MultiPut(1, "a") //overwrite

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
		actualValue, actualFound := tree.MultiGet(test[0])
		if actualFound != test[2] || (actualFound && actualValue.Value() != test[1]) {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}
}

func TestRedBlackTreeMultiRemove(t *testing.T) {
	tree := NewWithIntComparator()
	tree.MultiPut(5, "e")
	tree.MultiPut(6, "f")
	tree.MultiPut(7, "g")
	tree.MultiPut(3, "c")
	tree.MultiPut(4, "d")
	tree.MultiPut(1, "x")
	tree.MultiPut(2, "b")
	tree.MultiPut(1, "a") //overwrite

	tree.MultiRemove(5)
	tree.MultiRemove(6)
	tree.MultiRemove(7)
	tree.MultiRemove(8)
	tree.MultiRemove(5)

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
		actualValue, actualFound := tree.MultiGet(test[0])
		if actualFound != test[2] || (actualFound && actualValue.Value() != test[1]) {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}

	tree.MultiRemove(1)
	tree.MultiRemove(4)
	tree.MultiRemove(2)
	tree.MultiRemove(3)
	tree.MultiRemove(2)
	tree.MultiRemove(2)

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
	tree := NewWithIntComparator()

	if actualValue := tree.Left(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}
	if actualValue := tree.Right(); actualValue != nil {
		t.Errorf("Got %v expected %v", actualValue, nil)
	}

	tree.MultiPut(1, "a")
	tree.MultiPut(5, "e")
	tree.MultiPut(6, "f")
	tree.MultiPut(7, "g")
	tree.MultiPut(3, "c")
	tree.MultiPut(4, "d")
	tree.MultiPut(1, "x") // overwrite
	tree.MultiPut(2, "b")

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
	tree := NewWithIntComparator()

	if node, found := tree.Floor(0); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}
	if node, found := tree.Ceiling(0); node != nil || found {
		t.Errorf("Got %v expected %v", node, "<nil>")
	}

	tree.MultiPut(5, "e")
	tree.MultiPut(6, "f")
	tree.MultiPut(7, "g")
	tree.MultiPut(3, "c")
	tree.MultiPut(4, "d")
	tree.MultiPut(1, "x")
	tree.MultiPut(2, "b")

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
	tree := NewWithIntComparator()
	it := tree.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty tree")
	}
}

func TestRedBlackMultiTreeIteratorPrevOnEmpty(t *testing.T) {
	tree := NewWithIntComparator()
	it := tree.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty tree")
	}
}

func TestRedBlackMultiTreeIterator1Next(t *testing.T) {
	tree := NewWithIntComparator()
	tree.MultiPut(5, "e")
	tree.MultiPut(6, "f")
	tree.MultiPut(7, "g")
	tree.MultiPut(3, "c")
	tree.MultiPut(4, "d")
	tree.MultiPut(2, "b")
	tree.MultiPut(1, "a") //overwrite
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
	tree := NewWithIntComparator()
	tree.MultiPut(5, "e")
	tree.MultiPut(6, "f")
	tree.MultiPut(7, "g")
	tree.MultiPut(3, "c")
	tree.MultiPut(4, "d")
	tree.MultiPut(2, "b")
	tree.MultiPut(1, "a") //overwrite
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
	tree := NewWithIntComparator()
	tree.MultiPut(3, "c")
	tree.MultiPut(1, "a")
	tree.MultiPut(2, "b")
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
	tree := NewWithIntComparator()
	tree.MultiPut(3, "c")
	tree.MultiPut(1, "a")
	tree.MultiPut(2, "b")
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
	tree := NewWithIntComparator()
	tree.MultiPut(1, "a")
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
	tree := NewWithIntComparator()
	tree.MultiPut(1, "a")
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
	tree := NewWithIntComparator()
	tree.MultiPut(13, 5)
	tree.MultiPut(8, 3)
	tree.MultiPut(17, 7)
	tree.MultiPut(1, 1)
	tree.MultiPut(11, 4)
	tree.MultiPut(15, 6)
	tree.MultiPut(25, 9)
	tree.MultiPut(6, 2)
	tree.MultiPut(22, 8)
	tree.MultiPut(27, 10)
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
	tree := NewWithIntComparator()
	tree.MultiPut(13, 5)
	tree.MultiPut(8, 3)
	tree.MultiPut(17, 7)
	tree.MultiPut(1, 1)
	tree.MultiPut(11, 4)
	tree.MultiPut(15, 6)
	tree.MultiPut(25, 9)
	tree.MultiPut(6, 2)
	tree.MultiPut(22, 8)
	tree.MultiPut(27, 10)
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

func TestRedBlackMultiTreeIteratorBegin(t *testing.T) {
	tree := NewWithIntComparator()
	tree.MultiPut(3, "c")
	tree.MultiPut(1, "a")
	tree.MultiPut(2, "b")
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
	tree := NewWithIntComparator()
	it := tree.Iterator()

	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	it.End()
	if it.node != nil {
		t.Errorf("Got %v expected %v", it.node, nil)
	}

	tree.MultiPut(3, "c")
	tree.MultiPut(1, "a")
	tree.MultiPut(2, "b")
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
	tree := NewWithIntComparator()
	tree.MultiPut(3, "c")
	tree.MultiPut(1, "a")
	tree.MultiPut(2, "b")
	it := tree.Iterator()
	if actualValue, expectedValue := it.First(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
	}
}

func TestRedBlackMultiTreeIteratorLast(t *testing.T) {
	tree := NewWithIntComparator()
	tree.MultiPut(3, "c")
	tree.MultiPut(1, "a")
	tree.MultiPut(2, "b")
	it := tree.Iterator()
	if actualValue, expectedValue := it.Last(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
	}
}

func TestRedBlackMultiTreeSerialization(t *testing.T) {
	tree := NewWithStringComparator()
	tree.MultiPut("c", "3")
	tree.MultiPut("b", "2")
	tree.MultiPut("a", "1")

	var err error
	assert := func() {
		if actualValue, expectedValue := tree.Size(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if actualValue := tree.Keys(); actualValue[0].(string) != "a" || actualValue[1].(string) != "b" || actualValue[2].(string) != "c" {
			t.Errorf("Got %v expected %v", actualValue, "[a,b,c]")
		}
		if actualValue := tree.Values(); actualValue[0].(string) != "1" || actualValue[1].(string) != "2" || actualValue[2].(string) != "3" {
			t.Errorf("Got %v expected %v", actualValue, "[1,2,3]")
		}
		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	json, err := tree.ToJSON()
	assert()

	err = tree.FromJSON(json)
	assert()
}

func benchmarkMultiGet(b *testing.B, tree *Tree, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.MultiGet(n)
		}
	}
}

func benchmarkMultiPut(b *testing.B, tree *Tree, size int) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		tree.Clear()
		b.StartTimer()
		for n := 0; n < size; n++ {
			tree.MultiPut(n, struct{}{})
		}
	}
}

func benchmarkMultiRemove(b *testing.B, tree *Tree, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			tree.MultiRemove(n)
		}
	}
}

func BenchmarkRedBlackTreeMultiGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.MultiPut(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiGet(b, tree, size)
}

func BenchmarkRedBlackTreeMultiGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.MultiPut(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiGet(b, tree, size)
}

func BenchmarkRedBlackTreeMultiGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.MultiPut(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiGet(b, tree, size)
}

func BenchmarkRedBlackTreeMultiGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.MultiPut(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiGet(b, tree, size)
}

func BenchmarkRedBlackTreeMultiPut100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewWithIntComparator()
	b.StartTimer()
	benchmarkMultiPut(b, tree, size)
}

func BenchmarkRedBlackTreeMultiPut1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.MultiPut(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiPut(b, tree, size)
}

func BenchmarkRedBlackTreeMultiPut10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.MultiPut(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiPut(b, tree, size)
}

func BenchmarkRedBlackTreeMultiPut100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.MultiPut(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiPut(b, tree, size)
}

func BenchmarkRedBlackTreeMultiRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.MultiPut(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiRemove(b, tree, size)
}

func BenchmarkRedBlackTreeMultiRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.MultiPut(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiRemove(b, tree, size)
}

func BenchmarkRedBlackTreeMultiRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.MultiPut(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiRemove(b, tree, size)
}

func BenchmarkRedBlackTreeMultiRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.MultiPut(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiRemove(b, tree, size)
}

func BenchmarkMultiIterator_Next(b *testing.B) {
	b.StopTimer()
	size := 100000
	tree := NewWithIntComparator()
	for n := 0; n < size; n++ {
		tree.MultiPut(rand.Int(), struct{}{})
	}

	itr := tree.Iterator()
	b.StartTimer()
	for itr.Next() {

	}
	benchmarkMultiRemove(b, tree, size)
}
