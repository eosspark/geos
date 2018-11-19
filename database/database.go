package database

import (
	"bytes"
	"fmt"
	"github.com/eosspark/eos-go/log"
	"math"
	"reflect"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LDataBase struct {
	db        *leveldb.DB
	stack     *deque
	path      string
	reversion int64
	nextId    map[string]int64
	logFlag   bool
	log       log.Logger
	batch     *leveldb.Batch
	count     int64
}

/*

@param path 		--> 	database file (note:type-->d)

@return

success 			-->		database handle error is nil
error 				-->		database is nil ,error

*/
func NewDataBase(path string, flag ...bool) (DataBase, error) {

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
	/*	read every type increment	*/
	nextId, err := readIncrementFromDb(db)
	if err != nil {
		log.Error("database init failed : %s", err.Error())
		panic("open database file failed : " + err.Error())
	}

	logFlag := false
	if len(flag) > 0 {
		logFlag = flag[0]
	}


	dbLog := log.New("db")
	if logFlag {
		dbLog.SetHandler(log.TerminalHandler)
	} else {
		dbLog.SetHandler(log.DiscardHandler())
	}
	return &LDataBase{db: db, stack: newDeque(), path: path, nextId: nextId, logFlag: logFlag, log: dbLog, batch: new(leveldb.Batch)}, nil
}

func readIncrementFromDb(db *leveldb.DB) (map[string]int64, error) {
	nextId := make(map[string]int64)

	key := []byte(dbIncrement)
	val, err := db.Get(key, nil)
	if err != nil && err != leveldb.ErrNotFound {
		return nil, err
	}
	if err == nil && len(val) != 0 {
		err = DecodeBytes(val, &nextId)
		if err != nil {
			panic("database init failed : " + err.Error())
		}
	}
	err = db.Delete(key, nil)
	if err != nil {
		panic("database init failed : " + err.Error())
	}
	return nextId, nil
}

func (ldb *LDataBase) Close() {
	err := ldb.writeIncrementToDb()
	if err != nil {
		ldb.log.Error("database close failed : %s", err.Error())
	}
	err = ldb.db.Close()
	if err != nil {
		ldb.log.Error("database close failed : %s", err.Error())
	} else {
		ldb.log.Info("----------------- database close -----------------")
	}
}

func (ldb *LDataBase) writeIncrementToDb() error {

	if len(ldb.nextId) <=0{

		return nil
	}

	val, err := EncodeToBytes(ldb.nextId)
	if err != nil {
		ldb.log.Error("WriteIncrement rlp EncodeToBytes failed is : %s", err.Error())
		return err
	}
	key := []byte(dbIncrement)
	err = saveKey(key, val, ldb.db)
	if err != nil {
		ldb.log.Error("WriteIncrement saveKey failed is : %s", err.Error())
		return err
	}
	return nil
}

func (ldb *LDataBase) Revision() int64 {
	ldb.log.Info("ldb reversion is : %d", ldb.reversion)
	return ldb.reversion
}

func (ldb *LDataBase) Undo() {

	stack := ldb.getStack()
	if stack == nil {
		return
	}

	ldb.nextId = stack.oldIds
	for key, _ := range stack.OldValue {  		/* 	Undo edit */

		ldb.undoModifyKv(key)
	}
	for key, _ := range stack.NewValue {		/*	Undo new */
		ldb.remove(key)

	}
	for key, _ := range stack.RemoveValue {		/*	Undo remove */
		ldb.insert(key, true)
	}
	ldb.stack.Pop()
	ldb.reversion--
}

func (ldb *LDataBase) UndoAll() {
	// TODO wait Undo
	for ldb.enable() {
		ldb.Undo()
	}
}

func (ldb *LDataBase) squash() {

	if ldb.stack.Size() == 1 {
		ldb.stack.Pop()
		ldb.reversion--
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
	ldb.reversion--
}

func (ldb *LDataBase) StartSession() *Session {
	ldb.reversion++
	state := newUndoState(ldb.reversion, ldb.nextId)
	ldb.stack.Append(state)
	space := " "
	ldb.log.Info("database StartSession reversion is : %d,%s,  next id is : %v", ldb.reversion, space, ldb.nextId)
	return &Session{db: ldb, apply: true, reversion: ldb.reversion}
}

func (ldb *LDataBase) Commit(reversion int64) {

	for {
		if ldb.stack.Size() == 0 {
			ldb.log.Info("database Commit stack is empty")
			return
		}

		stack := ldb.getFirstStack()
		if stack == nil {
			ldb.log.Info("database Commit stack is nil")
			break
		}
		ldb.log.Info("database Commit stack reversion is : %d", stack.reversion)
		if stack.reversion > reversion {
			break
		}

		ldb.stack.PopFront()
	}
}

func (ldb *LDataBase) SetRevision(reversion int64) {
	if ldb.stack.Size() != 0 {
		panic("cannot set reversion while there is an existing undo stack")
		// throw
	}
	if reversion > math.MaxInt64 {
		//throw
		panic("reversion to set is too high")
	}
	ldb.log.Info("database  SetRevision  reversion is : ", reversion)
	ldb.reversion = reversion
}

/*
insert object to database
@param in 			--> 	object(pointer)

@return
success 			-->		nil
error 				-->		error

*/

func (ldb *LDataBase) Insert(in interface{}) error {
	err := ldb.insert(in)
	if err != nil {
		ldb.log.Error("error database insert failed : %s", err.Error())
		return err
	}

	ldb.undoInsert(in) // undo
	return nil
}

func (ldb *LDataBase) insert(in interface{}, flag ...bool) error { /* struct cfg --> KV struct --> kv to db --> undo db */

	cfg, err := parseObjectToCfg(in) /* (struct cfg) parse object tag */
	if err != nil {
		ldb.log.Error("error database insert  parseObjectToCfg failed : %s", err.Error())
		return err
	}
	ldb.log.Info("Info: %s", cfg.showCfg())
	undoFlag := false
	if len(flag) > 0 {
		undoFlag = flag[0]
	}

	if !undoFlag {
		err = ldb.setIncrement(cfg) /* (kv.id) set increment id */
		if err != nil {
			ldb.log.Error("error database insert cfg setIncrement failed : %s,  cfg is : %v", err.Error(), cfg)
			return err
		}
	}

	dbKV := &dbKeyValue{}
	structKV(in, dbKV, cfg) /* (kv.index) all key and value*/

	//dbKV.showDbKV()
	err = ldb.insertKvToBatch(dbKV) /* (kv to db) kv insert database (atomic) */
	if err != nil {
		ldb.log.Error("error database insert insertKvToDb failed : %s", err)
		return err
	}
	return nil
}

func (ldb *LDataBase) insertKvToBatch(dbKV *dbKeyValue) error {
	ldb.log.Info("Info kv database insertKvToDb dbKV is : %v", dbKV)
	ldb.batch.Reset()
	for _, v := range dbKV.index {
		ldb.batch.Put(v.key, v.value)
	}

	ldb.batch.Put(dbKV.id.key, dbKV.id.value)
	err := ldb.db.Write(ldb.batch, nil)
	if err != nil {
		return err
	}
	ldb.batch.Reset()
	return nil
}

func (ldb *LDataBase) setIncrement(cfg *structInfo) error {

	var id int64
	id = 1

	if _, ok := ldb.nextId[cfg.Name]; ok { // First insertion
		id = ldb.nextId[cfg.Name]
	}

	cfg.Id.Set(reflect.ValueOf(id).Convert(cfg.Id.Type()))
	ldb.nextId[cfg.Name] = id + 1
	ldb.log.Info("Info database setIncrement cfg id set is : %d,  next id is : %v", id, ldb.nextId)
	return nil
}

/*

remove object from database
@param in			--> 	object

@return
success 			-->		nil
error 				-->		error

*/

func (ldb *LDataBase) Remove(in interface{}) error {
	err := ldb.remove(in)
	if err != nil {
		ldb.log.Error("error database Remove failed : %s", err.Error())
		return err
	}
	ldb.undoRemove(in)
	return nil
}

func (ldb *LDataBase) remove(in interface{}) error {

	cfg, err := parseObjectToCfg(in)
	if err != nil {
		ldb.log.Error("error database remove parseObjectToCfg failed : %s", err.Error())
		return err
	}
	ldb.log.Info(cfg.showCfg())
	if isZero(cfg.Id) {
		ldb.log.Info("error database remove parseObjectToCfg cfg id isZero  : %v", cfg.Id)
		return ErrIncompleteStructure
	}

	dbKV := &dbKeyValue{}
	structKV(in, dbKV, cfg) /* (kv.index) all key and value*/

	//dbKV.showDbKV()

	err = ldb.removeKvToDb(dbKV)
	if err != nil {
		ldb.log.Error("error database remove removeKvToDb failed  : %s, dbKV is : %v", err, dbKV)
		return err
	}

	return nil
}

func (ldb *LDataBase) removeKvToDb(dbKV *dbKeyValue) error {
	undo := false
	undoKeyValue := &dbKeyValue{}

	defer func() {
		if !undo {
			return
		}
		/* 		assert(undo == true) */
		/* 		insert error --> remove undo 			*/
		ldb.undoRemoveKV(undoKeyValue)
	}()
	for _, v := range dbKV.index {
		err := removeKey(v.key, ldb.db)
		if err != nil {
			ldb.log.Error("error database  removeKvToDb removeKey failed  : %s, key is : %v", err.Error(), v.key)
			undo = true
			return err
		}
		undoKeyValue.index = append(undoKeyValue.index, v)
	}

	err := removeKey(dbKV.id.key, ldb.db)
	if err != nil {
		ldb.log.Error("error database  removeKvToDb removeKey failed  : %s, key is : %v", err, dbKV.id.key)
		undo = true
		return err
	}
	// assert(undo == false)
	return nil
}

/*

modify object from database
@param old 			--> 	object(pointer)
@param fn 			-->		function

@return
success 			-->		nil
error 				-->		error

*/

func (ldb *LDataBase) Modify(old interface{}, fn interface{}) error {
	copy_ := cloneInterface(old)
	err := ldb.modify(old, fn)
	if err != nil {
		ldb.log.Error("%s", err.Error())
		return err
	}
	ldb.undoModify(copy_)
	return nil
}

func (ldb *LDataBase) modify(data interface{}, fn interface{}) error {

	dataRef := reflect.ValueOf(data)
	if dataRef.Kind() != reflect.Ptr {
		return ErrPtrNeeded
	}

	oldInter := cloneInterface(data)
	dataType := dataRef.Type()

	fnRef := reflect.ValueOf(fn)
	fnType := fnRef.Type()
	if fnType.NumIn() != 1 {
		return errors.New("func Too many parameters : " + string(fnType.NumIn()))
	}

	pType := fnType.In(0)

	if pType.Kind() != dataType.Kind() {
		return errors.New("Parameter type does not match : " + pType.Kind().String() + " : " + dataType.Kind().String())
	}
	fnRef.Call([]reflect.Value{dataRef}) /*	call fn */
	// modify
	oldRef := reflect.ValueOf(oldInter)
	return ldb.modifyToDb(&oldRef, &dataRef)
}

func (ldb *LDataBase) modifyToDb(oldRef, newRef *reflect.Value) error {

	ldb.log.Info("Info database modifyKvToDb oldRef is : %v", oldRef)
	ldb.log.Info("Info database modifyKvToDb newRef is : %v", newRef)

	oldCfg, err := extractObjectTagInfo(oldRef)
	if err != nil {
		return errors.New("error database modifyKvToDb extractObjectTagInfo oldRef failed : " + err.Error())
	}
	newCfg, err := extractObjectTagInfo(newRef)
	if err != nil {
		return errors.New("error database modifyKvToDb extractObjectTagInfo newRef failed : " + err.Error())
	}

	ldb.log.Info("Info database modifyKvToDb oldCfg is : %s", oldCfg.showCfg())
	ldb.log.Info("Info database modifyKvToDb newCfg is : %s", newCfg.showCfg())

	if !reflect.DeepEqual(oldCfg.Id.Interface(), newCfg.Id.Interface()) {
		return errors.New(fmt.Sprintf("error database modifyKvToDb newCfg and oldCfg id failed,  newCfg id is :  %v,  oldCfg id is : %v", newCfg.Id, oldCfg.Id))
	}

	newKV := &dbKeyValue{}
	oldKV := &dbKeyValue{}
	structKV(oldRef.Interface(), oldKV, oldCfg)
	structKV(newRef.Interface(), newKV, newCfg)

	ldb.log.Info("Info database modifyKvToDb structKV oldKV is : %v", oldKV)
	ldb.log.Info("Info database modifyKvToDb structKV newKV is : %v", newKV)

	return ldb.modifyKvToDb(oldKV, newKV)
}

func (ldb *LDataBase) modifyKvToDb(oldKV, newKV *dbKeyValue) error {

	undoOldKV := &dbKeyValue{}
	undo := false
	defer func() {
		if !undo {
			return
		}
		ldb.undoRemoveKV(undoOldKV)
	}()

	for idx, _ := range oldKV.index {
		err := ldb.db.Delete(oldKV.index[idx].key, nil)
		if err != nil {
			undo = true
			return err
		}
	}

	err := ldb.db.Delete(oldKV.id.key, nil)
	if err != nil {
		undo = true
		return err
	}

	err = ldb.insertKvToBatch(newKV) // atomic
	if err != nil {
		undo = true
		return errors.New(fmt.Sprintf("error database modifyKvToDb insertKvToDb newKV is : %v", newKV))
	}

	return nil
}

/*

There is no modify fn in undo
so special processing is required when undo

*/

func (ldb *LDataBase) undoModifyKv(old interface{}) error {

	oldRef := reflect.ValueOf(old)
	if oldRef.Kind() != reflect.Ptr {
		return ErrPtrNeeded
	}

	ldb.log.Info("Info database undoModifyKv extractObjectTagInfo oldCfg is : %v", oldRef)
	oldCfg, err := extractObjectTagInfo(&oldRef)
	if err != nil {
		return errors.New(fmt.Sprintf("error database undoModifyKv extractObjectTagInfo oldCfg failed : %s", err.Error()))
	}
	id, err := EncodeToBytes(oldCfg.Id.Interface())
	if err != nil {
		return errors.New(fmt.Sprintf("error database undoModifyKv rlp EncodeToBytes id failed : %s, : %v", err.Error(), oldCfg.Id))
	}
	typeName := []byte(oldCfg.Name)
	key := splicingString(typeName,id)

	val, err := getDbKey(key, ldb.db)
	if err != nil {
		return errors.New(fmt.Sprintf("error database undoModifyKv getDbKey  failed : %s, key is : %v", err.Error(), key))
	}

	dst := reflect.New(reflect.Indirect(oldRef).Type())
	err = DecodeBytes(val, dst.Interface())
	if err != nil {
		return errors.New(fmt.Sprintf("error database undoModifyKv rlp DecodeBytes failed : %s, dst interface is %v", err, dst.Interface()))
	}

	return ldb.modifyToDb(&dst, &oldRef)
}

/*
find object from database
@param tagName 		--> 	tag in field tags
@param in 			--> 	object
@param out 			-->		output(pointer)

@return
success 			-->		nil 	(out valid)
error 				-->		error 	(out invalid)

*/

func (ldb *LDataBase) Find(tagName string, in interface{}, out interface{}) error {
	return ldb.find(tagName, in, out)
}

func (ldb *LDataBase) find(tagName string, value interface{}, to interface{}) error {
	ldb.log.Info("Info database find : tagName is: %s", tagName)
	fieldName := []byte(tagName)
	fields, err := getFieldInfo(tagName, value)
	if err != nil {
		return errors.New(fmt.Sprintf("error database find  getFieldInfo failed : %s", err.Error()))
	}

	ldb.log.Info("Info database find fields is : %v", fields)

	typeName := []byte(fields.typeName)

	suffix,err := fieldValueToByte(fields,true)
	if err != nil{
		return err
	}

	if len(suffix) == 0 {
		return errors.New(fmt.Sprintf("error database find nonUniqueValue failed : %s", err.Error()))
	}

	key := splicingString (typeName, fieldName)
	key = append(key, suffix...)

	end := keyEnd(key)
	it := ldb.db.NewIterator(&util.Range{Start: key, Limit: end}, nil)

	if !it.Next() {
		return leveldb.ErrNotFound
	}
	return ldb.findFields(it.Value(), typeName, to)
}

func (ldb *LDataBase) findFields(key, typeName []byte, to interface{}) error {
	ldb.log.Info("Info database findFields key is : %v", key)

	val := splicingString(typeName,key)

	v, err := getDbKey(val, ldb.db)
	err = DecodeBytes(v, to)
	if err != nil {
		return errors.New(fmt.Sprintf("error database findDbObject rlp DecodeBytes failed : %s", err.Error()))
	}
	return nil
}

/*
get MultiIndex from database
@param tagName 		--> 	tag in field tags
@param in 			--> 	object

@return
success 			-->		iterator
error 				-->		error

*/
func (ldb *LDataBase) GetIndex(tagName string, in interface{}) (*MultiIndex, error) {
	return ldb.getIndex(tagName, in)
}

func (ldb *LDataBase) getIndex(tagName string, value interface{}) (*MultiIndex, error) {

	// fieldName == tagName --> Just different nextId
	fieldName := []byte(tagName)
	fields, err := getFieldInfo(tagName, value)
	if err != nil {
		return nil, err
	}

	typeName := []byte(fields.typeName)
	begin := splicingString(typeName, fieldName)

	end := keyEnd(begin)

	ldb.log.Info("getIndex typeName : %v, fieldName : %v, begin: %v , end: %v, greater: %t", typeName, fieldName, begin, end, fields.greater)
	it := newMultiIndex(typeName, fieldName, begin, end, ldb)
	return it, nil
}

func (ldb *LDataBase) GetMutableIndex(fieldName string, in interface{}) (*MultiIndex, error) {
	return ldb.GetIndex(fieldName, in)
}

func (ldb *LDataBase) lowerBound(begin, end, fieldName []byte, data interface{}) (*DbIterator, error) {

	key, typeName := ldb.dbPrefix(begin, end, fieldName, data)
	if !bytes.HasPrefix(key, begin) {
		return nil, ErrNotFound
	}

	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
	if !it.Next() {
		return nil, ErrNotFound
	}

	if !it.Seek(key) {
		return nil, ErrNotFound
	}

	idx, err := newDbIterator(typeName, it, ldb.db)
	if err != nil {
		return nil, err
	}
	return idx, nil
	//begin ,typeName := ldb.dbPrefix(begin,end,fieldName,data,greater)
	//return ldb.dbIterator(begin,end,typeName,greater)
}

func (ldb *LDataBase) dbIterator(begin, end, typeName []byte) (*DbIterator, error) {
	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
	if !it.Next() {
		fmt.Println(begin)
		return nil, ErrNotFound
	}

	idx, err := newDbIterator(typeName, it, ldb.db)
	if err != nil {
		return nil, err
	}
	return idx, nil
}

func (ldb *LDataBase) EndIterator(begin, end, typeName []byte) (*DbIterator, error) {

	ldb.log.Info("begin : %v, end : %v, fieldName: %v , greater: %t", begin, end)

	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)

	it.Next()

	if it.Last() {
		it.Next()
		itr := &DbIterator{it: it, db: ldb.db, first: false, typeName: typeName, currentStatus: itEND}
		return itr, nil
	}

	return nil, ErrNotFound
}

func (ldb *LDataBase) dbPrefix(begin, end, fieldName []byte, data interface{}) ([]byte, []byte) {
	ldb.log.Info("begin : %v, end : %v, fieldName: %v , greater: %t", begin, end, fieldName)
	fields, err := getFieldInfo(string(fieldName), data)
	if err != nil {
		ldb.log.Error("db prefix getFieldInfo error %s", err.Error())
		return nil, nil
	}


	prefix,err := fieldValueToByte(fields,true)
	if err != nil{
		ldb.log.Error("db prefix fieldValueToByte error %s", err.Error())
		return nil,nil
	}

	ldb.log.Info("prefix is : %v", prefix)
	if len(prefix) == 0 {
		ldb.log.Error("db prefix getFieldValue error prefix failed")
		return nil, nil
	}

	prefix = append(begin, prefix...)

	return prefix, []byte(fields.typeName)
}

func (ldb *LDataBase) Empty(begin, end, fieldName []byte) bool {

	ldb.log.Info("begin : %v, end : %v, fieldName: %v ", begin, end, fieldName)
	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
	defer it.Release()
	if it.Next() {
		return false
	}
	return true
}

func (ldb *LDataBase) IteratorTo(begin, end, fieldName []byte, data interface{}) (*DbIterator, error) {

	ldb.log.Info("begin : %v, end : %v, fieldName: %v , greater: %t", begin, end, fieldName)
	fields, err := getFieldInfo(string(fieldName), data)
	if err != nil {
		return nil, err
	}
	prefix,err := fieldValueToByte(fields,true)
	if err != nil{
		ldb.log.Error("db IteratorTo fieldValueToByte error %s", err.Error())
		return nil,nil
	}
	if len(prefix) == 0 {
		return nil, errors.New("Get Field Value Failed")
	}

	key := []byte{}
	key = append(begin, prefix...)

	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
	if !it.Seek(key) {
		return nil, errors.New("Iterator To Not Found")
	}
	k := splicingString([]byte(fields.typeName),it.Value() )
	val, err := getDbKey(k, ldb.db)
	if err != nil {
		return nil, err
	}

	//
	itr := &DbIterator{it: it, db: ldb.db, first: false, value: val, typeName: []byte(fields.typeName)}
	return itr, nil
}

func (ldb *LDataBase) BeginIterator(begin, end, typeName []byte) (*DbIterator, error) {

	ldb.log.Info("begin : %v, end : %v, fieldName: %v , greater: %t", begin, end)
	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)

	if !it.Next() {
		return nil, errors.New("error DataBase BeginIterator : Greater False : Next Failed")
	}

	k := splicingString([]byte(typeName),it.Value() )
	val, err := getDbKey(k, ldb.db)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error DataBase BeginIterator : %s", err.Error()))
	}

	itr := &DbIterator{it: it, db: ldb.db, first: false, typeName: typeName, value: val, currentStatus: itBEGIN}

	return itr, nil
}

func (ldb *LDataBase) upperBound(begin, end, fieldName []byte, data interface{}) (*DbIterator, error) {

	begin, typeName := ldb.dbPrefix(begin, end, fieldName, data)
	begin[len(begin)-1] = begin[len(begin)-1] + 1

	return ldb.dbIterator(begin, end, typeName)
}

func (ldb *LDataBase) enable() bool { /* Whether the database enables the undo function*/
	return ldb.stack.Size() != 0
}

/* The following three functions are the undo functions provided by the database.*/

func (ldb *LDataBase) undoInsert(in interface{}) {
	if !ldb.enable() {
		return
	}

	stack := ldb.getStack()
	if stack == nil {
		ldb.log.Warn("undo session empty")
		return
	}
	copy_ := cloneInterface(in)
	stack.undoInsert(copy_)
}

func (ldb *LDataBase) undoModify(in interface{}) {
	if !ldb.enable() {
		return
	}

	stack := ldb.getStack()
	if stack == nil {
		ldb.log.Warn("undo session empty")
		return
	}
	copy_ := cloneInterface(in)
	stack.undoModify(copy_)
}

func (ldb *LDataBase) undoRemove(in interface{}) {
	if !ldb.enable() {
		return
	}

	stack := ldb.getStack()
	if stack == nil {
		ldb.log.Warn("undo session empty")
		return
	}
	copy_ := cloneInterface(in)
	stack.undoRemove(copy_)
}

/*

The three functions here are the implementation of
the database in order to operate the double-ended queue
and the subsequent queues will provide corresponding operations

*/

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
	}
	return nil
}

/*

When deleting or inserting an object
if an error occurs before completion
the following two functions are called to restore the database state
The internal undo operation of the database is not used externally

*/

func (ldb *LDataBase) undoRemoveKV(undoKeyValue *dbKeyValue) { /*	remove kv from db error --> undo kv */

	ldb.log.Info("info database undoRemoveKV  undoKeyValue is %v", undoKeyValue)

	if len(undoKeyValue.index) == 0 {
		return
	}
	for _, v := range undoKeyValue.index { /* 	undo index*/
		err := saveKey(v.key, v.value, ldb.db)
		if err != nil {
			//TODO assert
			ldb.log.Error("error database undoRemoveKV saveKey failed : %s, key is %v, value is: %v ", err.Error(), v.key, v.value)
			return
		}
	}
	if len(undoKeyValue.id.key) == 0 {
		return
	}

	err := saveKey(undoKeyValue.id.key, undoKeyValue.id.value, ldb.db) /* undo id*/
	if err != nil {
		//TODO assert
		ldb.log.Error("error database undoRemoveKV saveKey failed : %s, key is: %v, value is %v", err.Error(), undoKeyValue.id.key, undoKeyValue.id.value)
		return
	}
}

/*

The value corresponding to the given key is retrieved
from the db and returned to the caller
which may not exist

*/

func saveKey(key, value []byte, tx *leveldb.DB) error {
	if ok, _ := tx.Has(key, nil); ok {
		return ErrAlreadyExists
	}
	err := tx.Put(key, value, nil)
	if err != nil {
		return err
	}
	return nil
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
