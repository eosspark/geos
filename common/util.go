package common

import "reflect"

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
	/*for i, j := 0, 0; i < a.Len() && j < b.Len(); {
		if a.Data[i].GetKey() == b.Data[j].GetKey() {
			callback(a.Data[i], i, j)
			i++
			j++
		} else if a.Data[i].GetKey() > b.Data[j].GetKey() {
			j++
		} else if a.Data[i].GetKey() < b.Data[j].GetKey() {
			i++
		}
	}*/
}
