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

func (f *FlatSet) GetData(i int) *Element {
	if len(f.Data)-1 >= i {
		return &f.Data[i]
	}
	return nil
}

func (f *FlatSet) Clear() {
	if len(f.Data) > 0 {
		f.Data = nil
	}
}

func (f *FlatSet) Insert(element Element) (Element, bool) {
	length := f.Len()
	result := []Element{}
	target := f.Data
	exist := false
	r := 0
	if length == 0 {
		f.Data = append(f.Data, element)
		return f.Data[0], exist
	} else {
		i, j := 0, length-1
		for i < j {
			h := int(uint(i+j) >> 1)
			if element.GetKey() <= target[h].GetKey() {
				i = h + 1
			} else {
				j = h
			}
			r = h
		}
		if i <= j {
			if i == 0 || i == length-1 {
				//insert target before
				if element.GetKey() <= target[0].GetKey() {
					target = append(target, element)
					target = append(target, target...)
				} else if element.GetKey() >= target[length-1].GetKey() { //target append
					target = append(target, target...)
					target = append(target, element)
				}
			} else {
				//Insert middle
				if element.GetKey() == target[r].GetKey() {
					element = target[r]
					result = target
					exist = true
				} else {
					first := target[:r]
					second := target[r:length]
					result = append(result, first...)
					result = append(result, element)
					result = append(result, second...)
				}
			}
		}
	}
	return result[r], exist
}
