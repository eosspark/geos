package database

/////////////////////////////////////////////////////// UndoState  //////////////////////////////////////////////////////////
type modifyValue struct {
	Id    int64
	OldKv *dbKeyValue
	NewKv *dbKeyValue
}

type undoState struct {
	NewValue    map[int64]*modifyValue // hash or md5 ?
	RemoveValue map[int64]*modifyValue
	OldValue    map[int64]*modifyValue
}

type undoContainer struct {
	Undo      map[string]*undoState
	OldIds    map[string]int64
	Reversion int64
}

func newUndoContainer(reversion int64, oldIds map[string]int64) *undoContainer {

	oldIds_ := make(map[string]int64)
	for k, v := range oldIds {
		oldIds_[k] = v
	}
	return &undoContainer{Undo: make(map[string]*undoState), OldIds: oldIds, Reversion: reversion}
}

func newUndoState() *undoState {
	//nv := new (modifyValue)
	//rv := new (modifyValue)
	//ov := new (modifyValue)
	return &undoState{
		NewValue:    make(map[int64]*modifyValue),
		RemoveValue: make(map[int64]*modifyValue),
		OldValue:    make(map[int64]*modifyValue),
	}

}

func (container *undoContainer) checkType(typeName string) {
	if _, ok := container.Undo[typeName]; ok {
		return
	}
	undo := newUndoState()
	container.Undo[typeName] = undo
}

func (container *undoContainer) undoContainerInsert(typeName string, value *modifyValue) {
	container.checkType(typeName)

	_, ok := container.Undo[typeName]
	if !ok {
		panic("undo container insert error : " + typeName)
	}
	container.Undo[typeName].undoStackInsert(value)
}

func (container *undoContainer) undoContainerRemove(typeName string, value *modifyValue) {
	container.checkType(typeName)

	_, ok := container.Undo[typeName]
	if !ok {
		panic("undo container remove error : " + typeName)
	}
	container.Undo[typeName].undoStackRemove(value)
}

func (container *undoContainer) undoContainerModify(typeName string, value *modifyValue) {
	container.checkType(typeName)

	_, ok := container.Undo[typeName]
	if !ok {
		panic("undo container modify error : " + typeName)
	}
	container.Undo[typeName].undoStackModify(value)
}

func (stack *undoState) undoStackInsert(value *modifyValue) {
	id := value.Id
	stack.NewValue[id] = value
}

func (stack *undoState) undoStackRemove(value *modifyValue) {
	id := value.Id
	_, ok := stack.NewValue[id]
	if ok {
		delete(stack.NewValue, id)
		return
	}
	_, ok = stack.OldValue[id]
	if ok {
		stack.RemoveValue[id] = value
		delete(stack.OldValue, id)
	}

	_, ok = stack.RemoveValue[id]
	if ok {
		return
	}

	stack.RemoveValue[id] = value
}

func (stack *undoState) undoStackModify(value *modifyValue) {
	id := value.Id
	if _, ok := stack.NewValue[id]; ok {
		return
	}
	if _, ok := stack.OldValue[id]; ok {
		return
	}
	stack.OldValue[id] = value
}
