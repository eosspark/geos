package common

type Element interface {
	Compare(first Element, second Element) bool
	Equal(first Element, second Element) bool
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

func (f *FlatSet) Insert(param Element) (*Element, bool) {
	length := f.Len()
	result := []Element{}
	target := f.Data
	exist := false
	if length == 0 {
		f.Data = append(f.Data, param)
		return &param, exist
	} else {
		r, i, j := 0, 0, length-1
		for i < j {
			h := int(uint(i+j) >> 1)
			if param.Compare(target[h], param) {
				i = h + 1
			} else {
				j = h
			}
			r = h
		}
		if i <= j {
			if i == 0 || i == length-1 {
				//insert target before
				if param.Compare(param, target[0]) {
					target = append(target, param)
					target = append(target, target...)
				} else if param.Compare(target[length-1], param) { //target append
					target = append(target, target...)
					target = append(target, param)
				}
			} else {
				//Insert middle
				if param.Equal(target[r], param) {
					param = target[r]
					result = target
					exist = true
				} else {
					first := target[:r]
					second := target[r:length]
					result = append(result, first...)
					result = append(result, param)
					result = append(result, second...)
				}
			}
		}
	}
	return &param, exist
}
