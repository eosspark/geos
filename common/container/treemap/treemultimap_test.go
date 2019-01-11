// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package treemap

import (
	"fmt"
	"github.com/eosspark/container/utils"
	. "github.com/eosspark/eos-go/common/container/treemap/example"
	"testing"
)

func TestMultiMapPut(t *testing.T) {
	m := NewMultiIntStringMap()
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") //overwrite

	if actualValue := m.Size(); actualValue != 8 {
		t.Errorf("Got %v expected %v", actualValue, 8)
	}
	utils.AssertTest(t, []int{1, 1, 2, 3, 4, 5, 6, 7}, m.Keys())
	utils.AssertTest(t, []string{"x", "a", "b", "c", "d", "e", "f", "g"}, m.Values())

	// key,expectedValue,expectedFound
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
		lower := m.Get(test[0].(int))
		if expect := test[1]; expect != nil {
			utils.AssertTest(t, test[1].(string), lower.Value())
		} else {
			utils.AssertTest(t, true, lower == m.End())
			utils.AssertTest(t, true, lower.IsEnd())
			utils.AssertTest(t, false, lower.HasNext())
		}
	}
}

func TestMultiMapRemove(t *testing.T) {
	m := NewMultiIntStringMap()
	m.Put(5, "e")
	m.Put(6, "f")
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(4, "d")
	m.Put(1, "x")
	m.Put(2, "b")
	m.Put(1, "a") //overwrite

	m.Remove(5)
	m.Remove(6)
	m.Remove(7)
	m.Remove(8)
	m.Remove(5)

	utils.AssertTest(t, []int{1, 1, 2, 3, 4}, m.Keys())
	utils.AssertTest(t, []interface{}{"x", "a", "b", "c", "d"}, m.Values())

	if actualValue := m.Size(); actualValue != 5 {
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
		actualValue := m.Get(test[0].(int))
		if (actualValue.HasNext() && actualValue.Value() != test[1]) || actualValue.HasNext() != test[2] {
			t.Errorf("Got %v expected %v", actualValue, test[1])
		}
	}

	m.Remove(1)
	m.Remove(4)
	m.Remove(2)
	m.Remove(3)
	m.Remove(2)
	m.Remove(2)

	if actualValue, expectedValue := fmt.Sprintf("%d", m.Keys()), "[]"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue, expectedValue := fmt.Sprintf("%s", m.Values()), "[]"; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if actualValue := m.Size(); actualValue != 0 {
		t.Errorf("Got %v expected %v", actualValue, 0)
	}
	if actualValue := m.Empty(); actualValue != true {
		t.Errorf("Got %v expected %v", actualValue, true)
	}
}

/*
func TestMultiMapFloor(t *testing.T) {
	m := NewMultiWithIntComparator(utils.TypeString)
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(1, "a")

	// key,expectedKey,expectedValue,expectedFound
	tests1 := [][]interface{}{
		{-1, nil, nil, false},
		{0, nil, nil, false},
		{1, 1, "a", true},
		{2, 1, "a", true},
		{3, 3, "c", true},
		{4, 3, "c", true},
		{7, 7, "g", true},
		{8, 7, "g", true},
	}

	for _, test := range tests1 {
		// retrievals
		actualKey, actualValue := m.Floor(test[0])
		actualFound := actualKey != nil && actualValue != nil
		if actualKey != test[1] || actualValue != test[2] || actualFound != test[3] {
			t.Errorf("Got %v, %v, %v, expected %v, %v, %v", actualKey, actualValue, actualFound, test[1], test[2], test[3])
		}
	}
}

func TestMultiMapCeiling(t *testing.T) {
	m := NewMultiWithIntComparator(utils.TypeString)
	m.Put(7, "g")
	m.Put(3, "c")
	m.Put(1, "a")

	// key,expectedKey,expectedValue,expectedFound
	tests1 := [][]interface{}{
		{-1, 1, "a", true},
		{0, 1, "a", true},
		{1, 1, "a", true},
		{2, 3, "c", true},
		{3, 3, "c", true},
		{4, 7, "g", true},
		{7, 7, "g", true},
		{8, nil, nil, false},
	}

	for _, test := range tests1 {
		// retrievals
		actualKey, actualValue := m.Ceiling(test[0])
		actualFound := actualKey != nil && actualValue != nil
		if actualKey != test[1] || actualValue != test[2] || actualFound != test[3] {
			t.Errorf("Got %v, %v, %v, expected %v, %v, %v", actualKey, actualValue, actualFound, test[1], test[2], test[3])
		}
	}
}

func TestMultiMapEach(t *testing.T) {
	m := NewMultiWithStringComparator(utils.TypeString)
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)
	count := 0
	m.Each(func(key interface{}, value interface{}) {
		count++
		if actualValue, expectedValue := count, value; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		switch value {
		case 1:
			if actualValue, expectedValue := key, "a"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 2:
			if actualValue, expectedValue := key, "b"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case 3:
			if actualValue, expectedValue := key, "c"; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
	})
}

func TestMultiMapMap(t *testing.T) {
	m := NewWithStringComparator(utils.TypeString)
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)
	mappedMap := m.Map(func(key1 interface{}, value1 interface{}) (key2 interface{}, value2 interface{}) {
		return key1, value1.(int) * value1.(int)
	})
	if actualValue, _ := mappedMap.Get("a"); actualValue != 1 {
		t.Errorf("Got %v expected %v", actualValue, "mapped: a")
	}
	if actualValue, _ := mappedMap.Get("b"); actualValue != 4 {
		t.Errorf("Got %v expected %v", actualValue, "mapped: b")
	}
	if actualValue, _ := mappedMap.Get("c"); actualValue != 9 {
		t.Errorf("Got %v expected %v", actualValue, "mapped: c")
	}
	if mappedMap.Size() != 3 {
		t.Errorf("Got %v expected %v", mappedMap.Size(), 3)
	}
}

//func TestMapSelect(t *testing.T) {
//	m := NewWithStringComparator()
//	m.Put("c", 3)
//	m.Put("a", 1)
//	m.Put("b", 2)
//	selectedMap := m.Select(func(key interface{}, value interface{}) bool {
//		return key.(string) >= "a" && key.(string) <= "b"
//	})
//	if actualValue, _ := selectedMap.Get("a"); actualValue != 1 {
//		t.Errorf("Got %v expected %v", actualValue, "value: a")
//	}
//	if actualValue, _ := selectedMap.Get("b"); actualValue != 2 {
//		t.Errorf("Got %v expected %v", actualValue, "value: b")
//	}
//	if selectedMap.Size() != 2 {
//		t.Errorf("Got %v expected %v", selectedMap.Size(), 2)
//	}
//}
//
//func TestMapAny(t *testing.T) {
//	m := NewWithStringComparator()
//	m.Put("c", 3)
//	m.Put("a", 1)
//	m.Put("b", 2)
//	any := m.Any(func(key interface{}, value interface{}) bool {
//		return value.(int) == 3
//	})
//	if any != true {
//		t.Errorf("Got %v expected %v", any, true)
//	}
//	any = m.Any(func(key interface{}, value interface{}) bool {
//		return value.(int) == 4
//	})
//	if any != false {
//		t.Errorf("Got %v expected %v", any, false)
//	}
//}
//
//func TestMapAll(t *testing.T) {
//	m := NewWithStringComparator()
//	m.Put("c", 3)
//	m.Put("a", 1)
//	m.Put("b", 2)
//	all := m.All(func(key interface{}, value interface{}) bool {
//		return key.(string) >= "a" && key.(string) <= "c"
//	})
//	if all != true {
//		t.Errorf("Got %v expected %v", all, true)
//	}
//	all = m.All(func(key interface{}, value interface{}) bool {
//		return key.(string) >= "a" && key.(string) <= "b"
//	})
//	if all != false {
//		t.Errorf("Got %v expected %v", all, false)
//	}
//}

func TestMultiMapFind(t *testing.T) {
	m := NewMultiWithStringComparator(utils.TypeString)
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)
	foundKey, foundValue := m.Find(func(key interface{}, value interface{}) bool {
		return key.(string) == "c"
	})
	if foundKey != "c" || foundValue != 3 {
		t.Errorf("Got %v -> %v expected %v -> %v", foundKey, foundValue, "c", 3)
	}
	foundKey, foundValue = m.Find(func(key interface{}, value interface{}) bool {
		return key.(string) == "x"
	})
	if foundKey != nil || foundValue != nil {
		t.Errorf("Got %v at %v expected %v at %v", foundValue, foundKey, nil, nil)
	}
}

//func TestMultiMapChaining(t *testing.T) {
//	m := NewMultiWithStringComparator(utils.TypeString)
//	m.Put("c", 3)
//	m.Put("a", 1)
//	m.Put("b", 2)
//	chainedMap := m.Select(func(key interface{}, value interface{}) bool {
//		return value.(int) > 1
//	}).Map(func(key interface{}, value interface{}) (interface{}, interface{}) {
//		return key.(string) + key.(string), value.(int) * value.(int)
//	})
//	if actualValue := chainedMap.Size(); actualValue != 2 {
//		t.Errorf("Got %v expected %v", actualValue, 2)
//	}
//	if actualValue, found := chainedMap.Get("aa"); actualValue != nil || found {
//		t.Errorf("Got %v expected %v", actualValue, nil)
//	}
//	if actualValue, found := chainedMap.Get("bb"); actualValue != 4 || !found {
//		t.Errorf("Got %v expected %v", actualValue, 4)
//	}
//	if actualValue, found := chainedMap.Get("cc"); actualValue != 9 || !found {
//		t.Errorf("Got %v expected %v", actualValue, 9)
//	}
//}

func TestMultiMapIteratorNextOnEmpty(t *testing.T) {
	m := NewMultiWithStringComparator(utils.TypeString)
	it := m.Iterator()
	it = m.Iterator()
	for it.Next() {
		t.Errorf("Shouldn't iterate on empty map")
	}
}

func TestMultiMapIteratorPrevOnEmpty(t *testing.T) {
	m := NewMultiWithStringComparator(utils.TypeString)
	it := m.Iterator()
	it = m.Iterator()
	for it.Prev() {
		t.Errorf("Shouldn't iterate on empty map")
	}
}

func TestMultiMapIteratorNext(t *testing.T) {
	m := NewMultiWithStringComparator(utils.TypeString)
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)

	it := m.Iterator()
	count := 0
	for it.Next() {
		count++
		key := it.Key()
		value := it.Value()
		switch key {
		case "a":
			if actualValue, expectedValue := value, 1; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case "b":
			if actualValue, expectedValue := value, 2; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case "c":
			if actualValue, expectedValue := value, 3; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
		if actualValue, expectedValue := value, count; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
	}
	if actualValue, expectedValue := count, 3; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestMultiMapIteratorPrev(t *testing.T) {
	m := NewMultiWithStringComparator(utils.TypeString)
	m.Put("c", 3)
	m.Put("a", 1)
	m.Put("b", 2)

	it := m.Iterator()
	for it.Next() {
	}
	countDown := m.Size()
	for it.Prev() {
		key := it.Key()
		value := it.Value()
		switch key {
		case "a":
			if actualValue, expectedValue := value, 1; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case "b":
			if actualValue, expectedValue := value, 2; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		case "c":
			if actualValue, expectedValue := value, 3; actualValue != expectedValue {
				t.Errorf("Got %v expected %v", actualValue, expectedValue)
			}
		default:
			t.Errorf("Too many")
		}
		if actualValue, expectedValue := value, countDown; actualValue != expectedValue {
			t.Errorf("Got %v expected %v", actualValue, expectedValue)
		}
		countDown--
	}
	if actualValue, expectedValue := countDown, 0; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
}

func TestMultiMapIteratorBegin(t *testing.T) {
	m := NewMultiWithIntComparator(utils.TypeString)
	it := m.Iterator()
	it.Begin()
	m.Put(3, "c")
	m.Put(1, "a")
	m.Put(2, "b")
	for it.Next() {
	}
	it.Begin()
	it.Next()
	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
	}
}

func TestMultiMapTreeIteratorEnd(t *testing.T) {
	m := NewMultiWithIntComparator(utils.TypeString)
	it := m.Iterator()
	m.Put(3, "c")
	m.Put(1, "a")
	m.Put(2, "b")
	it.End()
	it.Prev()
	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
	}
}

func TestMultiMapIteratorFirst(t *testing.T) {
	m := NewMultiWithIntComparator(utils.TypeString)
	m.Put(3, "c")
	m.Put(1, "a")
	m.Put(2, "b")
	it := m.Iterator()
	if actualValue, expectedValue := it.First(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if key, value := it.Key(), it.Value(); key != 1 || value != "a" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 1, "a")
	}
}

func TestMultiMapIteratorLast(t *testing.T) {
	m := NewMultiWithIntComparator(utils.TypeString)
	m.Put(3, "c")
	m.Put(1, "a")
	m.Put(2, "b")
	it := m.Iterator()
	if actualValue, expectedValue := it.Last(), true; actualValue != expectedValue {
		t.Errorf("Got %v expected %v", actualValue, expectedValue)
	}
	if key, value := it.Key(), it.Value(); key != 3 || value != "c" {
		t.Errorf("Got %v,%v expected %v,%v", key, value, 3, "c")
	}
}

func TestMultiMapSerialization(t *testing.T) {
	for i := 0; i < 10; i++ {
		original := NewMultiWithStringComparator(utils.TypeString)
		original.Put("d", "4")
		original.Put("e", "5")
		original.Put("c", "3")
		original.Put("b", "2")
		original.Put("a", "1")

		assertMultiSerialization(original, "A", t)

		serialized, err := original.ToJSON()
		if err != nil {
			t.Errorf("Got error %v", err)
		}
		assertMultiSerialization(original, "B", t)

		deserialized := NewWithStringComparator(utils.TypeString)
		err = deserialized.FromJSON(serialized)
		if err != nil {
			t.Errorf("Got error %v", err)
		}
		assertSerialization(deserialized, "C", t)
	}
}

//noinspection GoBoolExpressions
func assertMultiSerialization(m *MultiMap, txt string, t *testing.T) {
	if actualValue := m.Keys(); false ||
		actualValue[0].(string) != "a" ||
		actualValue[1].(string) != "b" ||
		actualValue[2].(string) != "c" ||
		actualValue[3].(string) != "d" ||
		actualValue[4].(string) != "e" {
		t.Errorf("[%s] Got %v expected %v", txt, actualValue, "[a,b,c,d,e]")
	}
	if actualValue := m.Values(); false ||
		actualValue[0].(string) != "1" ||
		actualValue[1].(string) != "2" ||
		actualValue[2].(string) != "3" ||
		actualValue[3].(string) != "4" ||
		actualValue[4].(string) != "5" {
		t.Errorf("[%s] Got %v expected %v", txt, actualValue, "[1,2,3,4,5]")
	}
	if actualValue, expectedValue := m.Size(), 5; actualValue != expectedValue {
		t.Errorf("[%s] Got %v expected %v", txt, actualValue, expectedValue)
	}
}

func benchmarkMultiGet(b *testing.B, m *MultiMap, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			m.Get(n)
		}
	}
}

func benchmarkMultiPut(b *testing.B, m *MultiMap, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			m.Put(n, struct{}{})
		}
	}
}

func benchmarkMultiRemove(b *testing.B, m *MultiMap, size int) {
	for i := 0; i < b.N; i++ {
		for n := 0; n < size; n++ {
			m.Remove(n)
		}
	}
}

func BenchmarkTreeMultiMapGet100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := NewMultiWithIntComparator(utils.TypeString)
	for n := 0; n < size; n++ {
		m.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiGet(b, m, size)
}

func BenchmarkTreeMultiMapGet1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := NewMultiWithIntComparator(utils.TypeString)
	for n := 0; n < size; n++ {
		m.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiGet(b, m, size)
}

func BenchmarkTreeMultiMapGet10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := NewMultiWithIntComparator(utils.TypeString)
	for n := 0; n < size; n++ {
		m.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiGet(b, m, size)
}

func BenchmarkTreeMultiMapGet100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	m := NewMultiWithIntComparator(utils.TypeString)
	for n := 0; n < size; n++ {
		m.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiGet(b, m, size)
}

func BenchmarkTreeMultiMapPut100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := NewMultiWithIntComparator(utils.TypeString)
	b.StartTimer()
	benchmarkMultiPut(b, m, size)
}

func BenchmarkTreeMultiMapPut1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := NewMultiWithIntComparator(utils.TypeString)
	for n := 0; n < size; n++ {
		m.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiPut(b, m, size)
}

func BenchmarkTreeMultiMapPut10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := NewMultiWithIntComparator(utils.TypeString)
	for n := 0; n < size; n++ {
		m.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiPut(b, m, size)
}

func BenchmarkTreeMultiMapPut100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	m := NewMultiWithIntComparator(utils.TypeString)
	for n := 0; n < size; n++ {
		m.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiPut(b, m, size)
}

func BenchmarkTreeMultiMapRemove100(b *testing.B) {
	b.StopTimer()
	size := 100
	m := NewMultiWithIntComparator(utils.TypeString)
	for n := 0; n < size; n++ {
		m.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiRemove(b, m, size)
}

func BenchmarkTreeMultiMapRemove1000(b *testing.B) {
	b.StopTimer()
	size := 1000
	m := NewMultiWithIntComparator(utils.TypeString)
	for n := 0; n < size; n++ {
		m.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiRemove(b, m, size)
}

func BenchmarkTreeMultiMapRemove10000(b *testing.B) {
	b.StopTimer()
	size := 10000
	m := NewMultiWithIntComparator(utils.TypeString)
	for n := 0; n < size; n++ {
		m.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiRemove(b, m, size)
}

func BenchmarkTreeMultiMapRemove100000(b *testing.B) {
	b.StopTimer()
	size := 100000
	m := NewMultiWithIntComparator(utils.TypeString)
	for n := 0; n < size; n++ {
		m.Put(n, struct{}{})
	}
	b.StartTimer()
	benchmarkMultiRemove(b, m, size)
}
*/
