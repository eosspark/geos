// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package treeset

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/eosspark/eos-go/log"
)

func TestMultiSetNew(t *testing.T) {
	set := NewMultiWithIntComparator(2, 1)
	if actualValue := set.Size(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	values := set.Values()
	if actualValue := values[0]; actualValue != 1 {
		t.Errorf("Got %v expected %v", actualValue, 1)
	}
	if actualValue := values[1]; actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
}

func TestMultiSetAdd(t *testing.T) {
	set := NewMultiWithIntComparator()
	set.Add()
	set.Add(1)
	set.Add(2)
	set.Add(2, 3)
	set.Add()
	if actualValue := set.Empty(); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
	if actualValue := set.Size(); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, 4)
	}
	if actualValue, expectedValue := fmt.Sprintf("%d%d%d%d", set.Values()...), "1223"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestMultiSetContains(t *testing.T) {
	set := NewMultiWithIntComparator()
	set.Add(3, 1, 2)
	if actualValue := set.Contains(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := set.Contains(1); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := set.Contains(1, 2, 3); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
	if actualValue := set.Contains(1, 2, 3, 4); actualValue != false {
		t.Errorf("Got %v expected %v", actualValue, false)
	}
}

func TestMultiSetRemove(t *testing.T) {
	set := NewMultiWithIntComparator()
	set.Add(3, 1, 2)
	set.Remove()
	if actualValue := set.Size(); actualValue != 3 {
		t.Errorf("Got %v expected %v", actualValue, 3)
	}
	set.Remove(1)
	if actualValue := set.Size(); actualValue != 2 {
		t.Errorf("Got %v expected %v", actualValue, 2)
	}
	set.Remove(3)
	set.Remove(3)
	set.Remove()
	set.Remove(2)
	if actualValue := set.Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
}

func TestMultiSetEach(t *testing.T) {
	set := NewMultiWithStringComparator()
	set.Add("c", "a", "b", "a")
	index := -1
	set.Each(func(value interface{}) {
		index ++

		switch index {
		case 0:
			if actualValue, expectedValue := value, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 3:
			if actualValue, expectedValue := value, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
	})
}

func TestMultiSetMap(t *testing.T) {
	set := NewMultiWithStringComparator()
	set.Add("c", "a", "b")
	mappedSet := set.Map(func(value interface{}) interface{} {
		return "mapped: " + value.(string)
	})
	if actualValue, expectedValue := mappedSet.Contains("mapped: a", "mapped: b", "mapped: c"), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := mappedSet.Contains("mapped: a", "mapped: b", "mapped: x"), false; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if mappedSet.Size() != 3 {
		t.Errorf("Got %v expected %v", mappedSet.Size(), 3)
	}
}

//func TestSetSelect(t *testing.T) {
//	set := NewWithStringComparator()
//	set.Add("c", "a", "b")
//	selectedSet := set.Select(func(index int, value interface{}) bool {
//		return value.(string) >= "a" && value.(string) <= "b"
//	})
//	if actualValue, expectedValue := selectedSet.Contains("a", "b"), true; actualValue != expectedValue {
//		fmt.Println("A: ", selectedSet.Contains("b"))
//		t.Errorf("Got %v (%v) expected %v (%v)", actualValue, selectedSet.Values(), expectedValue, "[a b]")
//	}
//	if actualValue, expectedValue := selectedSet.Contains("a", "b", "c"), false; actualValue != expectedValue {
//		t.Errorf("Got %v (%v) expected %v (%v)", actualValue, selectedSet.Values(), expectedValue, "[a b]")
//	}
//	if selectedSet.Size() != 2 {
//		t.Errorf("Got %v expected %v", selectedSet.Size(), 3)
//	}
//}
//
//func TestSetAny(t *testing.T) {
//	set := NewWithStringComparator()
//	set.Add("c", "a", "b")
//	any := set.Any(func(index int, value interface{}) bool {
//		return value.(string) == "c"
//	})
//	if any != true {
//		t.Errorf("Got %v expected %v", any, true)
//	}
//	any = set.Any(func(index int, value interface{}) bool {
//		return value.(string) == "x"
//	})
//	if any != false {
//		t.Errorf("Got %v expected %v", any, false)
//	}
//}
//
//func TestSetAll(t *testing.T) {
//	set := NewWithStringComparator()
//	set.Add("c", "a", "b")
//	all := set.All(func(index int, value interface{}) bool {
//		return value.(string) >= "a" && value.(string) <= "c"
//	})
//	if all != true {
//		t.Errorf("Got %v expected %v", all, true)
//	}
//	all = set.All(func(index int, value interface{}) bool {
//		return value.(string) >= "a" && value.(string) <= "b"
//	})
//	if all != false {
//		t.Errorf("Got %v expected %v", all, false)
//	}
//}

func TestMultiSetFind(t *testing.T) {
	set := NewMultiWithStringComparator()
	set.Add("c", "a", "b")
	foundValue, found := set.Find(func(value interface{}) bool {
		return value.(string) == "c"
	})
	if !found || foundValue != "c" {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, found, "c", 2)
	}
	foundValue, found = set.Find(func(value interface{}) bool {
		return value.(string) == "x"
	})
	if foundValue != nil || found {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, found, nil, nil)
	}
}

func TestMultiSetChaining(t *testing.T) {
	set := NewMultiWithStringComparator()
	set.Add("c", "a", "b")
}

func TestMultiSetIteratorNextOnEmpty(t *testing.T) {
	set := NewMultiWithStringComparator()
	it := set.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty set")
	}
}

func TestMultiSetIteratorPrevOnEmpty(t *testing.T) {
	set := NewMultiWithStringComparator()
	it := set.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty set")
	}
}

func TestMultiSetIteratorNext(t *testing.T) {
	set := NewMultiWithStringComparator()
	set.Add("c", "a", "b", "b")
	it := set.Iterator()
	count := 0
	index := -1
	for it.Next() {
		count++
		index++
		value := it.Value()
		switch index {
		case 0:
			if actualValue, expectedValue := value, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 3:
			if actualValue, expectedValue := value, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
		if actualValue, expectedValue := index, count-1; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}
	if actualValue, expectedValue := count, 4; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestMultiSetIteratorPrev(t *testing.T) {
	set := NewWithStringComparator()
	set.Add("c", "a", "b")
	it := set.Iterator()
	for it.Prev() {
	}
	count := 0
	for it.Next() {
		count++
		index := it.Index()
		value := it.Value()
		switch index {
		case 0:
			if actualValue, expectedValue := value, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 1:
			if actualValue, expectedValue := value, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := value, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
		if actualValue, expectedValue := index, count-1; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}
	if actualValue, expectedValue := count, 3; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}


func TestMultiSetIteratorBegin(t *testing.T) {
	set := NewMultiWithStringComparator()
	it := set.Iterator()
	it.Begin()
	set.Add("a", "b", "c")
	for it.Next() {
	}
	it.Begin()
	it.Next()
	if value := it.Value(); value != "a" {
		t.Errorf("Got %v,%v expected %v", value, 0, "a")
	}
}

func TestMultiSetIteratorEnd(t *testing.T) {
	set := NewMultiWithStringComparator()
	it := set.Iterator()

	set.Add("a", "b", "b", "c")
	it.End()
	it.Prev()
	if value := it.Value(); value != "c" {
		t.Errorf("Got %v,%v expected %v", value, set.Size()-1, "c")
	}
}

func TestMultiSetIteratorFirst(t *testing.T) {
	set := NewMultiWithStringComparator()
	set.Add("a", "a", "b", "c")
	it := set.Iterator()
	if actualValue, expectedValue := it.First(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if value := it.Value(); value != "a" {
		t.Errorf("Got %v,%v expected %v", value, 0, "a")
	}
}

func TestMultiSetIteratorLast(t *testing.T) {
	set := NewMultiWithStringComparator()
	set.Add("a", "a", "b", "c")
	it := set.Iterator()
	if actualValue, expectedValue := it.Last(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if value := it.Value(); value != "c" {
		t.Errorf("Got %v,%v expected %v", value, 2, "c")
	}
}

func TestMultiSetSerialization(t *testing.T) {
	set := NewWithStringComparator()
	set.Add("a", "b", "c")

	var err error
	assert := func() {
		if actualValue, expectedValue := set.Size(), 3; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		if actualValue := set.Contains("a", "b", "c"); actualValue != true {
			t.Errorf("Got %v expected %v", actualValue, true)
		}
		if err != nil {
			t.Errorf("Got error %v", err)
		}
	}

	assert()

	json, err := set.ToJSON()
	assert()

	err = set.FromJSON(json)
	assert()
}


func TestMultiSetIntersection(t *testing.T) {
	a := NewMultiWithIntComparator(1, 3, 5, 7, 9)
	b := NewMultiWithIntComparator(2, 3, 7, 10)
	res := make([]int, 0, 2)

	MultiSetIntersection(a, b, func(elem interface{}) {
		res = append(res, elem.(int))
	})

	if len(res) != 2 || res[0] != 3 || res[1] != 7 {
		t.Errorf("Got %v expected (3,7)", res)
	}

	fmt.Println(res)
}


func benchmarkMultiContains(b *testing.B, set *MultiSet, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			set.Contains(n)
		}
	}
}

func benchmarkMultiAdd(b *testing.B, set *MultiSet, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			set.Add(n)
		}
	}
}

func benchmarkMultiRemove(b *testing.B, set *MultiSet, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			set.Remove(n)
		}
	}
}

func BenchmarkTreeMultiSetContains100(b *testing.B) {
	b.StopTimer()
	size := 100
	set := NewMultiWithIntComparator()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkMultiContains(b, set, size)
}

func BenchmarkTreeMultiSetContains1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	set := NewMultiWithIntComparator()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkMultiContains(b, set, size)
}

func BenchmarkTreeMultiSetContains10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	set := NewMultiWithIntComparator()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkMultiContains(b, set, size)
}

func BenchmarkTreeMultiSetContains100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	set := NewMultiWithIntComparator()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkMultiContains(b, set, size)
}

func BenchmarkTreeMultiSetAdd100(b *testing.B) {
	b.StopTimer()
	size := 100
	set := NewMultiWithIntComparator()
	b.StartTimer()
	benchmarkMultiAdd(b, set, size)
}

func BenchmarkTreeMultiSetAdd1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	set := NewMultiWithIntComparator()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkMultiAdd(b, set, size)
}

func BenchmarkTreeMultiSetAdd10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	set := NewMultiWithIntComparator()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkMultiAdd(b, set, size)
}

func BenchmarkTreeMultiSetAdd100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	set := NewMultiWithIntComparator()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkMultiAdd(b, set, size)
}

func BenchmarkTreeMultiSetRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	set := NewMultiWithIntComparator()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkMultiRemove(b, set, size)
}

func BenchmarkTreeMultiSetRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	set := NewMultiWithIntComparator()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkMultiRemove(b, set, size)
}

func BenchmarkTreeMultiSetRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	set := NewMultiWithIntComparator()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkMultiRemove(b, set, size)
}

func BenchmarkTreeMultiSetRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	set := NewMultiWithIntComparator()
	for n := 0; n < size; n++ {
		set.Add(n)
	}
	b.StartTimer()
	benchmarkMultiRemove(b, set, size)
}

func TestMultiSet_UpperBound(t *testing.T) {
	set := NewMultiWithStringComparator()
	set.Add("c", "a", "b","c","c","d")

	u := set.UpperBound("a")
	if u!=nil{
		log.Info("%v",u.Value())
	}
	assert.Equal(t,"b",u.Value())
}

func TestMultiSet_LowerBound(t *testing.T) {
	set := NewMultiWithStringComparator()
	set.Add("c", "a", "b","c","c","b","d")
	//foundValue, found := set.Find(func(value interface{}) bool {
	//	return value.(string) == "c"
	//})
	//fmt.Println(foundValue,found)
	u := set.LowerBound("a")
	//fmt.Println(u.Next())
	sec := set.UpperBound("a")
	for u.Next(){

		if u.Equal(*sec){
			break
		}
		fmt.Print(u.Value())
	}

}

func TestMultiSet_easer(t *testing.T){
	set := NewMultiWithStringComparator()
	set.Add("c", "a", "b","c","c","b","d")
	lb := set.LowerBound("c")
	up:= set.UpperBound("c")

	for lb.Next(){
		if lb.iterator.Equal(up.iterator){
			break
		}
		fmt.Println("lower-upper:",lb.Value())
	}
	set.tree.MultiRemove("c")
	itr:=set.Iterator()
	itr.Begin()
	for itr.Next(){
		fmt.Println(itr.Value())
	}
	assert.Equal(t,4,set.Size())
}