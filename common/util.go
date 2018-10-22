package common

import "reflect"

func Empty(i interface{}) bool {
	current := reflect.ValueOf(i).Interface()
	empty := reflect.Zero(reflect.ValueOf(i).Type()).Interface()

	return reflect.DeepEqual(current, empty)
}

func CompareSlice(first interface{}, secend interface{}) bool {

	/*if len(first) != len(secend){
		return false
	}
	for i:=0; i<len(first);i++{
		if first[i]!=secend[i]{
			return false
		}
	}*/
	return true
}
