package database

import (
	"fmt"
	"reflect"
)

/////////////////////////////////////////////////////// UndoState  //////////////////////////////////////////////////////////
type undoState struct {
	NewValue    map[interface{}]int64 // hash or md5 ?
	RemoveValue map[interface{}]int64
	OldValue    map[interface{}]int64
	oldIds      map[string]int64
	reversion   int64
}

func newUndoState(reversion int64, oldIds map[string]int64) *undoState {

	oldIds_ := make(map[string]int64)
	for k, v := range oldIds {
		oldIds_[k] = v
	}
	return &undoState{
		NewValue:    make(map[interface{}]int64),
		RemoveValue: make(map[interface{}]int64),
		OldValue:    make(map[interface{}]int64),
		oldIds:      oldIds_,
		reversion:   reversion,
	}
}

func (stack *undoState) undoInsert(data interface{}) {
	stack.NewValue[data] = stack.reversion
}

func (stack *undoState) undoRemove(data interface{}) {
	_, ok := stack.NewValue[data]
	if ok {
		undoMapRemove(stack.NewValue, data)
		return
	}
	_, ok = stack.OldValue[data]
	if ok {
		stack.RemoveValue[data] = stack.reversion
		undoMapRemove(stack.OldValue, data)
		return
	}

	_, ok = stack.RemoveValue[data]
	if ok {
		return
	}
	stack.RemoveValue[data] = stack.reversion
}

func (stack *undoState) undoModify(data interface{}) {
	key := undoEqual(stack.NewValue, data)
	if key != nil {
		return
	}
	key = undoEqual(stack.OldValue, data)
	if key != nil {
		return
	}
	stack.reversion++
	stack.OldValue[data] = stack.reversion
}

func undoEqual(m map[interface{}]int64, data interface{}) interface{} {
	for key, value := range m {
		if reflect.DeepEqual(key, data) {
			fmt.Println(value)
			return key
		}
	}
	return nil
}

func undoMapRemove(m map[interface{}]int64, data interface{}) {
	key := undoEqual(m, data)
	delete(m, key)
}
