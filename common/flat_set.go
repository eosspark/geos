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

func (f *FlatSet) searchSub(key uint64) int {
	length := len(f.Data)
	r := -1
	if length != 0 {
		i, j := 0, length-1
		for i <= j {
			if i+j > 0 {
				h := int(uint(i+j) >> 1)
				if f.Data[h].GetKey() <= key {
					i = h + 1
				} else {
					j = h
				}
				r = h
			} else {
				r = 0
				break
			}
		}
	}
	return r
}
func (f *FlatSet) FindData(key uint64) (Element, int) {
	r := f.searchSub(key)

	if key == f.Data[r].GetKey() {
		return f.Data[r], r
	}

	return nil, -1
}

func (f *FlatSet) Find(element Element) (bool, int) {
	r := f.searchSub(element.GetKey())

	if element.GetKey() == f.Data[r].GetKey() {
		return true, r
	}

	return false, -1
}

func (f *FlatSet) Insert(element Element) (Element, bool) {
	var result Element
	length := f.Len()
	target := f.Data
	exist := false
	r := 0
	if length == 0 {
		f.Data = append(f.Data, element)
		result = f.Data[0]
	} else {
		i, j := 0, length-1
		for i < j {
			h := int(uint(i+j) >> 1)
			if target[h].GetKey() <= element.GetKey() {
				i = h + 1
			} else {
				j = h
			}
			r = h
		}
		//r := f.searchSub(element.GetKey())
		if target[0].GetKey() < element.GetKey() && element.GetKey() < target[length-1].GetKey() {
			//Insert middle
			if element.GetKey() == target[r].GetKey() {
				element = target[r]
				result = target[r]
				exist = true
			} else {
				elemnts := []Element{}
				first := target[:r+1]
				second := target[r+1 : length]
				elemnts = append(elemnts, first...)
				elemnts = append(elemnts, element)
				elemnts = append(elemnts, second...)
				f.Data = elemnts
				result = elemnts[r+1]
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
		//}
	}
	return result, exist
}
