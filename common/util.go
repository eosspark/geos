package common

import (
	"reflect"
	"bytes"
	)

func Empty(i interface{}) bool {
	if i == nil {
		return true
	}
	current := reflect.ValueOf(i).Interface()
	empty := reflect.Zero(reflect.ValueOf(i).Type()).Interface()

	return reflect.DeepEqual(current, empty)
}

//use callback to handle the same element
//<int>i,j indexes the same element in FlatSet a and b
func SetIntersection(a FlatSet, b FlatSet, callback func(e Element, i int, j int)) {
	for i, j := 0, 0; i < a.Len() && j < b.Len(); {
		if bytes.Compare(a.Data[i].GetKey(), b.Data[j].GetKey()) == 0 {
			callback(a.Data[i], i, j)
			i++
			j++
		} else if bytes.Compare(a.Data[i].GetKey(), b.Data[j].GetKey()) == 1 {
			j++
		} else if bytes.Compare(a.Data[i].GetKey(), b.Data[j].GetKey()) == -1 {
			i++
		}
	}
}

var NameComparator = func(a, b interface{}) int {
	aAsserted := uint64(a.(Name))
	bAsserted := uint64(b.(Name))
	switch {
	case aAsserted > bAsserted:
		return 1
	case aAsserted < bAsserted:
		return -1
	default:
		return 0
	}
}