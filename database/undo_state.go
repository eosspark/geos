package database

/////////////////////////////////////////////////////// UndoState  //////////////////////////////////////////////////////////
type modifyValue struct{
	id int64
	oldKv	*dbKeyValue
	newKv *dbKeyValue
}
type undoState struct {
	NewValue    map[int64]*modifyValue// hash or md5 ?
	RemoveValue map[int64]*modifyValue
	OldValue    map[int64]*modifyValue
	oldIds      map[string]int64
	reversion   int64
}

func newUndoState(reversion int64, oldIds map[string]int64) *undoState {

	oldIds_ := make(map[string]int64)
	for k, v := range oldIds {
		oldIds_[k] = v
	}
	return &undoState{
		NewValue:    make(map[int64]*modifyValue),
		RemoveValue: make(map[int64]*modifyValue),
		OldValue:    make(map[int64]*modifyValue),
		oldIds:      oldIds_,
		reversion:   reversion,
	}
}

func (stack *undoState) undoInsert(value *modifyValue) {
	id := value.id
	stack.NewValue[id] = value
}

func (stack *undoState) undoRemove(value *modifyValue) {
	id := value.id
	_, ok := stack.NewValue[id]
	if ok {
		delete(stack.NewValue,id)
		return
	}
	_, ok = stack.OldValue[id]
	if ok {
		stack.RemoveValue[id] = value
		delete(stack.OldValue,id)
	}

	_, ok = stack.RemoveValue[id]
	if ok {
		return
	}

	stack.RemoveValue[id] = value
}

func (stack *undoState) undoModify(value *modifyValue) {
	id := value.id
	if _, ok := stack.NewValue[id];ok{
		return
	}
	if _, ok := stack.OldValue[id];ok{
		return
	}
	stack.reversion++
	stack.OldValue[id] = value
}
