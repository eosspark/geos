package common

import (
	"errors"
)

type ElementObject interface {
	ElementObject() //default implements interface only taget
}

type Bucket struct {
	Data    []ElementObject
	Compare func(first ElementObject, second ElementObject) int
}

func (b *Bucket) Len() int {
	return len(b.Data)
}

func (b *Bucket) GetData(i int) (ElementObject, error) {
	if len(b.Data)-1 >= i {
		return b.Data[i], nil
	}
	return nil, errors.New("not found data")
}

func (b *Bucket) Clear() {
	if len(b.Data) > 0 {
		b.Data = nil
	}
}

func (b *Bucket) Find(element ElementObject) (bool, int) {
	r := b.searchSub(element)
	if r >= 0 && b.Compare(element, b.Data[r]) == 0 {
		return true, r
	}
	return false, -1
}

func (b *Bucket) Eraser(element ElementObject) bool {
	result := false
	if b.Len() == 0 {
		return result
	}
	exist, sub := b.Find(element)
	if exist /* && f.Len()>=1*/ {
		b.Data = append(b.Data[:sub], b.Data[sub+1:]...)
		result = true
	}
	return result
}

func (b *Bucket) searchSub(obj ElementObject) int {
	length := b.Len()
	if length == 0 {
		return -1
	}
	i, j := 0, length-1
	for i < j {
		h := int(uint(i+j) >> 1)
		if i <= h && h < j {
			if b.Compare(b.Data[h], obj) == -1 {
				i = h + 1
			} else if b.Compare(b.Data[h], obj) == 0 {
				return h
			} else {
				j = h
			}
		}
	}
	return i
}

func (b *Bucket) Insert(obj ElementObject) (*ElementObject, error) {
	if b.Compare == nil {
		return nil, errors.New("Bucket Compare is nil")
	}
	var result ElementObject
	length := b.Len()
	target := b.Data
	if length == 0 {
		b.Data = append(b.Data, obj)
		result = b.Data[0]
	} else {
		r := b.searchSub(obj)
		start := b.Compare(target[0], obj)
		end := b.Compare(obj, target[length-1])
		if (start == -1 || start == 0) && (end == -1 || end == 0) {
			//Insert middle
			elemnts := []ElementObject{}
			first := target[:r]
			second := target[r:length]
			elemnts = append(elemnts, first...)
			elemnts = append(elemnts, obj)
			elemnts = append(elemnts, second...)
			b.Data = elemnts
			result = elemnts[r]
		} else {
			//insert target before
			if b.Compare(obj, target[0]) == -1 {
				elemnts := []ElementObject{}
				elemnts = append(elemnts, obj)
				elemnts = append(elemnts, target...)
				b.Data = elemnts
				result = elemnts[0]
			} else if b.Compare(obj, target[length-1]) == 1 { //target append
				target = append(target, obj)
				result = target[length]
				b.Data = target
			}
		}
	}
	return &result, nil
}

func (b *Bucket) LowerBound(eo ElementObject) (ElementObject, int) {
	first := 0
	if b.Len() > 0 {
		ext := b.searchSub(eo)
		first = ext
		for i := first; i >= 0; i-- {
			if b.Compare(b.Data[i], eo) == -1 {
				value := b.Data[i+1]
				currentSub := i + 1
				return value, currentSub
			} else if i == 0 && b.Compare(b.Data[i], eo) == 0 {
				value := b.Data[i]
				currentSub := i
				return value, currentSub
			}
		}
	}
	return nil, -1
}

func (b *Bucket) UpperBound(eo ElementObject) (ElementObject, int) {
	if b.Len() > 0 {
		ext := b.searchSub(eo)
		for i := ext; i < b.Len(); i++ {
			if b.Compare(b.Data[i], eo) > 0 {
				value := b.Data[i-1]
				currentSub := i - 1
				return value, currentSub
			} else if i == b.Len()-1 && b.Compare(eo, b.Data[i]) == 0 {
				value := b.Data[i]
				currentSub := i
				return value, currentSub
			}
		}
	}
	return nil, -1
}
