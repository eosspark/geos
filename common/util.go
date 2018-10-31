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
