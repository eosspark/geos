package common

import (
	"reflect"
	)

type CheckEmpty interface {
	IsEmpty() bool
}

func empty(i interface{}) bool {
	switch t := i.(type) {
	case nil:
		return true
	case uint8:
		return t == 0
	case uint16:
		return t == 0
	case uint32:
		return t == 0
	case uint64:
		return t == 0
	case int32:
		return t == 0
	case int64:
		return t == 0
	case int:
		return t == 0
	case string:
		return t == ""
	case bool:
		return !t
	case *CheckEmpty:
		return t == nil
	case CheckEmpty:
		return t.IsEmpty()
	default:
		return false
	}
}

func Empty(i interface{}) bool {
	if i == nil {
		return true
	}
	current := reflect.ValueOf(i).Interface()
	empty := reflect.Zero(reflect.ValueOf(i).Type()).Interface()

	return reflect.DeepEqual(current, empty)
}
