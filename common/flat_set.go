package common

type Element interface {
	GetKey() uint64
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

func (f *FlatSet) searchLessSub(key uint64) int {
	length := len(f.Data)
	i, j := 0, length-1
	for i < j {
		h := int(uint(i+j) >> 1)
		if i <= h && h < j {
			if f.Data[h].GetKey() <= key {
				i = h + 1
			} else {
				j = h
			}
		}
	}
	return i
}

func (f *FlatSet) searchEqualSub(key uint64) int {
	length := len(f.Data)
	i, j := 0, length-1
	for i < j {
		h := int(uint(i+j) >> 1)
		if i <= h && h < j {
			if f.Data[h].GetKey() == key {
				i = h + 1
			} else {
				j = h
			}
		}
	}
	return i
}
func (f *FlatSet) FindData(key uint64) (Element, int) {
	r := f.searchEqualSub(key)

	if r >= 0 && key == f.Data[r].GetKey() {
		return f.Data[r], r
	}

	return nil, -1
}

func (f *FlatSet) Find(element Element) (bool, int) {
	r := f.searchEqualSub(element.GetKey())
	if r >= 0 && element.GetKey() == f.Data[r].GetKey() {
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
		r := f.searchLessSub(element.GetKey())
		if target[0].GetKey() < element.GetKey() && element.GetKey() < target[length-1].GetKey() {
			//Insert middle
			if element.GetKey() == target[r-1].GetKey() {
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
			if element.GetKey() < target[0].GetKey() {
				elemnts := []Element{}
				elemnts = append(elemnts, element)
				elemnts = append(elemnts, target...)
				f.Data = elemnts
				result = elemnts[0]
			} else if element.GetKey() > target[length-1].GetKey() { //target append
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
