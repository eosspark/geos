package common

import (
	"bytes"
)

type Element interface {
	GetKey() []byte
}

type FlatSet struct {
	Data []Element
}

func (f *FlatSet) Len() int {
	return len(f.Data)
}

func (f *FlatSet) GetData(i int) Element {
	if len(f.Data)-1 >= i {
		return f.Data[i]
	}
	return nil
}

func (f *FlatSet) Clear() {
	if len(f.Data) > 0 {
		f.Data = nil
	}
}

func (f *FlatSet) searchSub(key []byte) int {
	length := f.Len()
	if length == 0 {
		return -1
	}
	i, j := 0, length-1
	for i < j {
		h := int(uint(i+j) >> 1)
		if i <= h && h < j {
			if bytes.Compare(f.Data[h].GetKey(), key) == -1 {
				i = h + 1
			} else if bytes.Compare(f.Data[h].GetKey(), key) == 0 {
				return h
			} else {
				j = h
			}
		}
	}
	return i
}
func (f *FlatSet) FindData(key []byte) (Element, int) {
	r := f.searchSub(key)

	if r >= 0 && bytes.Compare(f.Data[r].GetKey(), key) == 0 {
		return f.Data[r], r
	}

	return nil, -1
}

func (f *FlatSet) Find(element Element) (bool, int) {
	r := f.searchSub(element.GetKey())
	if r >= 0 && bytes.Compare(element.GetKey(), f.Data[r].GetKey()) == 0 {
		return true, r
	}
	return false, -1
}

func (f *FlatSet) Insert(element Element) (Element, bool) {
	var result Element
	length := f.Len()
	target := f.Data
	exist := false
	if length == 0 {
		f.Data = append(f.Data, element)
		result = f.Data[0]
	} else {
		r := f.searchSub(element.GetKey())
		if bytes.Compare(target[0].GetKey(), element.GetKey()) == -1 &&
			bytes.Compare(element.GetKey(), target[length-1].GetKey()) == -1 {
			//Insert middle
			if bytes.Compare(element.GetKey(), target[r-1].GetKey()) == 0 {
				element = target[r-1]
				result = target[r-1]
				exist = true
			} else {
				elemnts := []Element{}
				first := target[:r]
				second := target[r:length]
				elemnts = append(elemnts, first...)
				elemnts = append(elemnts, element)
				elemnts = append(elemnts, second...)
				f.Data = elemnts
				result = elemnts[r]
			}
		} else {
			//insert target before
			if bytes.Compare(element.GetKey(), target[0].GetKey()) == -1 {
				elemnts := []Element{}
				elemnts = append(elemnts, element)
				elemnts = append(elemnts, target...)
				f.Data = elemnts
				result = elemnts[0]
			} else if bytes.Compare(element.GetKey(), target[length-1].GetKey()) == 1 { //target append
				target = append(target, element)
				result = target[length]
				f.Data = target
			}
		}
	}
	return result, exist
}

func (f *FlatSet) Update(element Element) bool {
	result := false
	if f.Len() == 0 {
		_, result = f.Insert(element)
		return true
	}
	if f.Len() > 0 {
		_, sub := f.Find(element)
		if sub == -1 {
			f.Insert(element)
			result = true
		} else {
			f.Data[sub] = element
			result = true
		}
	}
	return result
}

func (f *FlatSet) Remove(element Element) bool {

	result := false
	if f.Len() == 0 {
		return false
	}
	_, sub := f.FindData(element.GetKey())
	if sub >= 0 /* && f.Len()>=1*/ {
		f.Data = append(f.Data[:sub], f.Data[sub+1:]...)
		result = true
	}
	return result
}

func (f *FlatSet) Reserve(size int) {
	f.Data = make([]Element, 0, size)
}

func (f *FlatSet) Copy(oldFS FlatSet) {
	f.Data = make([]Element, oldFS.Len())
	copy(f.Data, oldFS.Data)
}
