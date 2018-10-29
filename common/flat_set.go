package common

type FlatSet interface {
	Compare(first FlatSet, second FlatSet) bool

	Equal(first FlatSet, second FlatSet) bool
}

func Append(target []FlatSet, param FlatSet) (*[]FlatSet, *FlatSet) {
	length := len(target)
	result := []FlatSet{}

	if length == 0 {
		target = append(target, param)
		return &target, &param
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
					result = append(result, param)
					result = append(result, target...)
				} else if param.Compare(target[length-1], param) { //target append
					result = append(result, target...)
					result = append(result, param)
				}
			} else {
				//Insert middle
				if param.Equal(target[r], param) {
					param = target[r]
					result = target
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
	return &result, &param
}
