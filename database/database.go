package database

import (
	"log"
	"math"
	"reflect"

	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LDataBase struct {
	db       *leveldb.DB
	stack    *deque
	path     string
	revision int64
}

/*

@param path 		--> 	database file (note:type-->d)

@return

success 			-->		database handle error is nil
error 				-->		database is nil ,error

*/
func NewDataBase(path string) (DataBase, error) {

	db, err := leveldb.OpenFile(path, &opt.Options{
		OpenFilesCacheCapacity: 16,
		BlockCacheCapacity:     16 / 2 * opt.MiB,
		WriteBuffer:            16 / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(path, nil)
	}
	if err != nil {
		return nil, err
	}

	return &LDataBase{db: db, stack:newDeque(),path: path}, nil
}

func (ldb *LDataBase) Close() {
	err := ldb.db.Close()
	if err != nil {
		// log
	} else {
		// log
	}
}

func (ldb *LDataBase) Revision() int64 {
	return ldb.revision
}

func (ldb *LDataBase) Undo() {
	//
	stack := ldb.getStack()
	if stack == nil {
		return
	}
	for key, _ := range stack.OldValue {
		log.Fatalln("modify do not work")
		// db.modify
		ldb.Modify(key,nil)
	}
	for key, _ := range stack.NewValue {
		// db.remove
		ldb.Remove(key)
	}
	for key, _ := range stack.RemoveValue {
		// db.insert
		ldb.Insert(key)
	}
	ldb.stack.Pop()
	ldb.revision--
}

func (ldb *LDataBase) UndoAll() {
	// TODO wait Undo
	log.Fatalln("UndoAll do not work,Please call linx")
}

func (ldb *LDataBase) squash() {

	if ldb.stack.Size() == 1 {
		ldb.stack.Pop()
		ldb.revision--
		return
	}
	stack := ldb.getStack()
	preStack := ldb.getSecond()
	for key, value := range stack.OldValue {
		if _, ok := preStack.NewValue[key]; ok {
			continue
		}
		if _, ok := preStack.OldValue[key]; ok {
			continue
		}
		if _, ok := preStack.RemoveValue[key]; ok {
			//fmt.Println("squash failed")
			// panic ?
		}
		preStack.OldValue[key] = value
	}

	for key, value := range stack.NewValue {
		preStack.NewValue[key] = value
	}

	for key, value := range stack.RemoveValue {
		//fmt.Println(key, " --> ", value)
		if _, ok := preStack.NewValue[key]; ok {
			k := undoEqual(preStack.NewValue, key)
			delete(preStack.NewValue, k)
		}
		if _, ok := preStack.OldValue[key]; ok {
			preStack.RemoveValue[key] = value
			k := undoEqual(preStack.OldValue, key)
			delete(preStack.OldValue, k)
		}
		preStack.RemoveValue[key] = value
	}
	ldb.stack.Pop()
	ldb.revision--
}

func (ldb *LDataBase) StartSession() *Session {
	ldb.revision++
	state := newUndoState(ldb.revision)
	ldb.stack.Append(state)
	return &Session{db: ldb, apply: true, revision: ldb.revision}
}

func (ldb *LDataBase) Commit(revision int64) {

	for {
		stack := ldb.getFirstStack()
		if stack == nil {
			break
		}
		if stack.reversion <= ldb.revision {
			ldb.stack.PopFront()
			ldb.revision--
		} else {
			break
		}
	}
}

func (ldb *LDataBase) SetRevision(revision int64) {
	if ldb.stack.Size() != 0 {
		// throw
	}
	if revision > math.MaxInt64 {
		//throw
	}
	ldb.revision = revision
}

//////////////////////////////////////////////////////	insert object to database //////////////////////////////////////////////////////
/*

@param in 			--> 	object(pointer)

@return
success 			-->		nil
error 				-->		error

*/

func (ldb *LDataBase) Insert(in interface{}) error {
	err := save(in, ldb.db)
	if err != nil {
		// undo
		return err
	}
	ldb.undoInsert(in)	// undo
	return nil
}

func(ldb *LDataBase)GetMutableIndex(fieldName string, in interface{}) (*multiIndex, error){
	return ldb.GetIndex(fieldName,in)
}

//////////////////////////////////////////////////////	find object from database //////////////////////////////////////////////////////
/*

@param tagName 		--> 	tag in field tags
@param in 			--> 	object
@param out 			-->		output(pointer)

@return
success 			-->		nil 	(out valid)
error 				-->		error 	(out invalid)

*/

func (ldb *LDataBase) Find(tagName string, in interface{}, out interface{}) error {
	return find(tagName, in, out, ldb.db)
}

//////////////////////////////////////////////////////	get multiIndex from database //////////////////////////////////////////////////////
/*

@param tagName 		--> 	tag in field tags
@param in 			--> 	object

@return
success 			-->		iterator
error 				-->		error

*/
func (ldb *LDataBase) GetIndex(tagName string, in interface{}) (*multiIndex, error) {
	return getIndex(tagName, in, ldb)
}

//////////////////////////////////////////////////////	modify object from database //////////////////////////////////////////////////////
/*

@param old 			--> 	object(pointer)
@param fn 			-->		function

@return
success 			-->		nil
error 				-->		error

*/
func (ldb *LDataBase) Modify(old interface{}, fn interface{}) error {
	copy_ := cloneInterface(old)
	err := modify(old, fn, ldb.db)
	if err != nil {
		return err
	}
	ldb.undoModify(copy_)
	return nil
}

//////////////////////////////////////////////////////	remove object from database //////////////////////////////////////////////////////
/*

@param in			--> 	object

@return
success 			-->		nil
error 				-->		error

*/
func (ldb *LDataBase) Remove(in interface{}) error {
	err := remove(in, ldb.db)
	if err != nil {
		return err
	}
	ldb.undoRemove(in)
	return  nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func save(data interface{}, tx *leveldb.DB) error {

	ref := reflect.ValueOf(data)
	if !ref.IsValid() || ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
		return ErrStructPtrNeeded
	}

	cfg, err := extractStruct(&ref)
	if err != nil {
		return err
	}
	//	cfg.showStructInfo() 			// XXX
	if _, ok := cfg.Fields[tagID]; !ok {
		return ErrNoID
	}

	err = incrementField(cfg, tx)
	if err != nil {
		return err
	}
	id,err := rlp.EncodeToBytes(cfg.Id.Interface())
	if err != nil{
		return err
	}
	typeName := []byte(cfg.Name)

	callBack := func(key, value []byte) error {
		return saveKey(key, value, tx)
	}
	err = doCallBack(id, typeName, cfg, callBack)
	if err != nil {
		return err
	}
	key := idKey(id, typeName)
	value, err := rlp.EncodeToBytes(data)
	if err != nil {
		return err
	}

	return saveKey(key, value, tx)
}

func saveKey(key, value []byte, tx *leveldb.DB) error {

	if ok, _ := tx.Has(key, nil); ok {
		return ErrAlreadyExists
	}
	//fmt.Println(key)
	err := tx.Put(key, value, nil)
	if err != nil {
		return err
	}
	return nil
}

func remove(data interface{}, db *leveldb.DB) error {

	ref := reflect.ValueOf(data)
	if !ref.IsValid() || reflect.Indirect(ref).Kind() != reflect.Struct {
		return ErrBadType
	}

	if ref.Kind() == reflect.Ptr {
		return ErrStructNeeded
	}

	cfg, err := extractStruct(&ref)
	if err != nil {
		return err
	}

	//	cfg.showStructInfo()

	if isZero(cfg.Id) {
		return ErrIncompleteStructure
	}
	id,err := rlp.EncodeToBytes(cfg.Id.Interface())
	if err != nil{
		return err
	}
	typeName := []byte(cfg.Name)

	//fmt.Println(typeName)

	removeField := func(key, value []byte) error {
		exist, err := db.Has(key, nil)
		if err != nil {
			return nil
		}
		if !exist {
			return ErrNotFound
		}

		return removeKey(key, db)
	}

	err = doCallBack(id, typeName, cfg, removeField) // FIXME --> id --> obj --> cfg
	if err != nil {
		return err
	}
	key := idKey(id, typeName)

	return removeKey(key, db)
}

func removeKey(key []byte, db *leveldb.DB) error {

	if ok, _ := db.Has(key, nil); !ok {
		return ErrNotFound
	}
	err := db.Delete(key, nil)
	if err != nil {
		return err
	}
	return nil
}

func modify(data interface{}, fn interface{}, db *leveldb.DB) error {
	// ready
	dataRef := reflect.ValueOf(data)
	if dataRef.Kind() != reflect.Ptr {
		return ErrPtrNeeded
	}

	oldInter := cloneInterface(data)
	dataType := dataRef.Type()

	fnRef := reflect.ValueOf(fn)
	fnType := fnRef.Type()
	if fnType.NumIn() != 1 {
		return errors.New("func Too many parameters")
	}

	pType := fnType.In(0)

	if pType.Kind() != dataType.Kind() {
		//fmt.Println(pType.String(), " <--> ", dataType.String())
		return errors.New("Parameter type does not match")
	}

	fnRef.Call([]reflect.Value{dataRef})
	// modify
	oldRef := reflect.ValueOf(oldInter)
	err := modifyKey(&oldRef, &dataRef, db)
	if err != nil {
		return err
	}

	return nil
}

func modifyKey(old, new *reflect.Value, db *leveldb.DB) error {
	newCfg, err := extractStruct(new)
	if err != nil {
		return err
	}
	oldCfg, err := extractStruct(old)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(oldCfg.Id.Interface(), newCfg.Id.Interface()) {
		return ErrNoID
	}

	callBack := func(newKey, oldKey []byte) error {
		find, err := db.Has(oldKey, nil)
		if err != nil {
			return err
		}
		if !find {
			return ErrNotFound
		}
		value, err := db.Get(oldKey, nil)
		if err != nil {
			return err
		}
		err = db.Delete(oldKey, nil)
		if err != nil {
			return err
		}
		//fmt.Println(value)
		return saveKey(newKey, value, db)
	}

	id,err := rlp.EncodeToBytes(newCfg.Id.Interface())
	if err != nil{
		return err
	}
	typeName := []byte(newCfg.Name)
	key := idKey(id, typeName)
	val, err := rlp.EncodeToBytes(new.Interface())
	if err != nil {
		return err
	}

	err = modifyField(newCfg, oldCfg, callBack)
	if err != nil {
		return err
	}

	err = db.Delete(key, nil)
	if err != nil {
		return err
	}

	return saveKey(key, val, db)
}

func find(tagName string, value interface{}, to interface{}, db *leveldb.DB) error {
	// fieldName == tagName --> Just different types
	fieldName := []byte(tagName)
	fields, err := getFieldInfo(tagName, value)
	if err != nil {
		return err
	}

	typeName := []byte(fields.typeName)

	suffix := nonUniqueValue(fields)
	if suffix == nil {
		return ErrNotFound
	}

	key := typeNameFieldName(typeName, fieldName)
	key = append(key, suffix...)
	//fmt.Println("key is : ",key)
	if !fields.unique {
		/*
			non unique --> typename__tagName__fieldValue[0]__fieldValue[1]...
		*/
		return findNonUniqueFields(key, typeName, to, db)
	} else {
		/*
			unique --> typename__tagName__fieldValue
		*/

		return findUniqueFields(key, typeName, to, db)
	}
	return nil
}

func findNonUniqueFields(key, typeName []byte, to interface{}, db *leveldb.DB) error {
	end := make([]byte, len(key))
	copy(end, key)
	end[len(end)-1] = end[len(end)-1] + 1
	it := db.NewIterator(&util.Range{Start: key, Limit: end}, nil)
	if !it.Next() {
		return ErrNotFound
	}
	//fmt.Println(it.Value())

	return findDbObject(it.Value(), []byte(typeName), to, db)
}

func findUniqueFields(key, typeName []byte, to interface{}, db *leveldb.DB) error {
	v, err := getDbKey(key, db)
	if err != nil {
		return err
	}
	return findDbObject(v, typeName, to, db)
}

// only key is id can be called
func findDbObject(key, typeName []byte, to interface{}, db *leveldb.DB) error {

	id := idKey(key, typeName)

	val, err := getDbKey(id, db)
	if err != nil {
		return err
	}
	err = rlp.DecodeBytes(val, to)
	if err != nil {
		return err
	}
	return nil
}

func getIndex(tagName string, value interface{}, db DataBase) (*multiIndex, error) {

	// fieldName == tagName --> Just different types
	fieldName := []byte(tagName)
	fields, err := getFieldInfo(tagName, value)
	if err != nil {
		return nil, err
	}

	if fields.unique {
		return nil, ErrNotFound
	}

	typeName := []byte(fields.typeName)
	begin := typeNameFieldName(typeName, fieldName)
	begin = append(begin, '_')
	begin = append(begin, '_')

	/*
		non unique --> typename__fieldName__
	*/

	end := getNonUniqueEnd(begin)
	it := newMultiIndex(typeName, fieldName, begin, end, fields.greater, db)
	return it, nil
}

/*
key 	 	--> typeName__tagName
value 		--> id
*/

func incrementField(cfg *structInfo, tx *leveldb.DB) error {

	typeName := []byte(cfg.Name)
	tagName := []byte(tagID)
	// typeName__tagName
	key := append(typeName, '_')
	key = append(key, '_')
	key = append(key, tagName...)

	counter, err := getIncrementId(key, cfg, tx)
	if err != nil {
		return err
	}
	return setIncrementId(counter, key, cfg, tx)
}

func setIncrementId(counter int64, key []byte, cfg *structInfo, tx *leveldb.DB) error {
	cfg.Id.Set(reflect.ValueOf(counter).Convert(cfg.Id.Type()))
	value ,err := rlp.EncodeToBytes(cfg.Id.Interface())
	if value == nil && err == nil {
		return err
	}
	err = tx.Delete(key, nil)
	if err != nil {
		return ErrIdTagIncrement
	}
	return saveKey(key, value, tx)
}

func getIncrementId(key []byte, cfg *structInfo, tx *leveldb.DB) (int64, error) {

	valByte, err := tx.Get(key, nil)
	if err != nil && err != leveldb.ErrNotFound {
		return 0, ErrIdTagIncrement
	}

	counter := cfg.IncrementStart
	if valByte != nil {
		err := rlp.DecodeBytes(valByte,&counter)
		if err != nil {
			return 0, ErrIdTagIncrement
		}
		//fmt.Println(key ,"  ","found id : ",counter)
		counter++
	}
	return counter, nil
}

func getDbKey(key []byte, db *leveldb.DB) ([]byte, error) {
	exits, err := db.Has(key, nil)
	if err != nil {
		return nil, err
	}
	if !exits {
		return nil, ErrNotFound
	}

	val, err := db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return val, err
}

func cloneInterface(data interface{}) interface{} {

	src := reflect.ValueOf(data)
	dst := reflect.New(reflect.Indirect(src).Type())

	srcElem := src.Elem()
	dstElem := dst.Elem()
	NumField := srcElem.NumField()
	for i := 0; i < NumField; i++ {
		sf := srcElem.Field(i)
		df := dstElem.Field(i)
		df.Set(sf)
	}
	return dst.Interface()
}

func cloneByte(dst, src []byte) int {
	dst = make([]byte, len(src))
	return copy(dst, src)
}

func (ldb *LDataBase) lowerBound(begin, end, fieldName []byte, data interface{}, greater bool) (*DbIterator, error) {
	//TODO

	fields, err := getFieldInfo(string(fieldName), data)
	if err != nil {
		return nil, err
	}

	reg, prefix := getNonUniqueFieldValue(fields)
	if reg == nil {
		return nil, ErrNoID
	}
	sift := string(reg)
	//	fmt.Println(begin)
	if len(prefix) != 0 {
		begin = append(begin, prefix...)
		it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
		//for it.Next(){
		//	fmt.Println(it.Key())
		//}
		idx, err := newDbIterator([]byte(fields.typeName), it, ldb.db, sift, greater)
		if err != nil {
			return nil, err
		}
		return idx, nil
	}

	return nil, ErrNotFound
}

func (ldb *LDataBase) upperBound(begin, end, fieldName []byte, data interface{}, greater bool) (*DbIterator, error) {
	//TODO
	fields, err := getFieldInfo(string(fieldName), data)
	if err != nil {
		return nil, err
	}

	reg, prefix := getNonUniqueFieldValue(fields)
	if reg == nil {
		return nil, ErrNoID
	}
	if len(prefix) != 0 {
		begin = append(begin, prefix...)
	}

	sift := string(reg)
	begin[len(begin)-1] = begin[len(begin)-1] + 1
	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)

	idx, err := newDbIterator([]byte(fields.typeName), it, ldb.db, sift, greater)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

func (ldb *LDataBase) getFirstStack() *undoState {
	stack := ldb.stack.First()
	switch typ := stack.(type) {
	case *undoState:
		return typ
	default:
		return nil
		//panic(TYPE_NOT_FOUND)
	}
	return nil
}

func (ldb *LDataBase) getSecond() *undoState {
	stack := ldb.stack.LastSecond()
	switch typ := stack.(type) {
	case *undoState:
		return typ
	default:
		return nil
		//panic(TYPE_NOT_FOUND)
	}
	return nil
}

func (ldb *LDataBase) getStack() *undoState {
	stack := ldb.stack.Last()
	if stack == nil {
		return nil
	}
	switch typ := stack.(type) {
	case *undoState:
		return typ
	default:
		return nil
		//panic(TYPE_NOT_FOUND)
	}
	return nil
}

///////////////

func (ldb *LDataBase)enable()bool{
	return ldb.stack.Size()	 != 0
}

func (ldb *LDataBase) undoInsert(in interface{}){
	if !ldb.enable(){
		return
	}

	stack := ldb.getStack()
	if stack == nil {
		log.Println("undo session empty")
		return
	}
	copy_ := cloneInterface(in)
	stack.undoInsert(copy_)
}

func (ldb *LDataBase) undoModify(in interface{}){
	if !ldb.enable(){
		return
	}

	stack := ldb.getStack()
	if stack == nil {
		log.Println("undo session empty")
		return
	}
	copy_ := cloneInterface(in)
	stack.undoModify(copy_)
}

func (ldb *LDataBase) undoRemove(in interface{}){
	if !ldb.enable(){
		return
	}

	stack := ldb.getStack()
	if stack == nil {
		log.Println("undo session empty")
		return
	}
	copy_ := cloneInterface(in)
	stack.undoRemove(copy_)
}
