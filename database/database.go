package database

import (
	"fmt"
	"log"
	"math"
	"os"
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
	reversion int64
	nextId   map[string]int64
	logFlag  bool
	dbLog 	*log.Logger
	logFile *os.File
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
	nextId, err := typeIncrement(db)
	if err != nil {
		log.Panicln("database init failed : " + err.Error())
	}

	logFlag := false
	if len(flag) > 0 {
		logFlag = flag[0]
	}

	fileName := path + "/database.log"

	reFn := func() {
		errs := os.RemoveAll(fileName)
		if errs != nil {
			log.Fatalln(errs)
		}
	}
	_, exits := os.Stat(fileName)
	if exits == nil {
		reFn()
	}

	logFile, err := os.Create(fileName)
	if err != nil {
		log.Fatalln("open database log file failed : ",err)
	}

	// create log object
	dbLog := log.New(logFile, "[Debug]", log.LstdFlags|log.Lshortfile)
	return &LDataBase{db: db, stack: newDeque(), path: path, nextId: nextId, logFlag: logFlag,dbLog:dbLog,logFile:logFile}, nil
}

func typeIncrement(db *leveldb.DB) (map[string]int64, error) {
	nextId := make(map[string]int64)
	dbIncrement := dbIncrement
	key := []byte(dbIncrement)
	val, err := db.Get(key, nil)
	if err != nil && err != leveldb.ErrNotFound {
		return nil, err
	}
	if err == nil && len(val) != 0 {
		err = rlp.DecodeBytes(val, &nextId)
		if err != nil {
			panic("database init failed : " + err.Error())
		}
	}
	err = db.Delete(key,nil)
	if err != nil{
		panic("database init failed : " + err.Error())
	}
	return nextId, nil
}

func (ldb *LDataBase) Close() {
	err := ldb.WriteIncrement()
	if err != nil {
		ldb.dbLog.Panicln("database close failed : " + err.Error())
	}
	err = ldb.db.Close()
	if err != nil {
		ldb.dbLog.Panicln("database close failed : " + err.Error())
		// log
	} else {
		ldb.dbLog.Println("----------------- database close -----------------")
		// log
	}

	ldb.logFile.Close()
}

func (ldb *LDataBase) WriteIncrement() error {
	val, err := rlp.EncodeToBytes(ldb.nextId)
	if err != nil {
		ldb.dbLog.Println("WriteIncrement rlp EncodeToBytes failed is : ",err)
		return err
	}
	dbIncrement := dbIncrement
	key := []byte(dbIncrement)
	err = saveKey(key, val, ldb.db)
	if err != nil {
		ldb.dbLog.Println("WriteIncrement saveKey failed is : ",err)
		return err
	}
	return nil
}

func (ldb *LDataBase) Revision() int64 {
	ldb.dbLog.Println("ldb reversion is : ",ldb.reversion)
	return ldb.reversion
}

func (ldb *LDataBase) Undo() {
	//
	stack := ldb.getStack()
	if stack == nil {
		return
	}

	ldb.nextId = stack.oldIds

	for key, _ := range stack.OldValue {

		ldb.undoModifyKv(key)
	}
	for key, _ := range stack.NewValue {
		// db.remove
		ldb.remove(key)

	}
	for key, _ := range stack.RemoveValue {
		// db.insert
		ldb.insert(key, true)
		//save(key,ldb.db,true)
	}
	ldb.stack.Pop()
	ldb.reversion--
}

func (ldb *LDataBase) UndoAll() {
	// TODO wait Undo
	log.Fatalln("UndoAll do not work,Please wait")
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
	ldb.dbLog.Println("database StartSession reversion is : ",ldb.reversion,space,"next id is : ",ldb.nextId)
	return &Session{db: ldb, apply: true, reversion: ldb.reversion}
}

func (ldb *LDataBase) Commit(reversion int64) {

	for {
		if ldb.stack.Size() == 0 {
			ldb.dbLog.Println("database Commit stack is empty")
			return
		}

		stack := ldb.getFirstStack()
		if stack == nil {
			ldb.dbLog.Println("database Commit stack is nil")
			break
		}
		ldb.dbLog.Println("database Commit stack reversion is : ",stack.reversion)
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
	ldb.dbLog.Println("database  SetRevision  reversion is : ",reversion)
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
		ldb.dbLog.Println("error database insert failed : ",err)
		return err
	}

	ldb.undoInsert(in) // undo
	return nil
}

func (ldb *LDataBase) insert(in interface{}, flag ...bool) error { /* struct cfg --> KV struct --> kv to db --> undo db */

	cfg, err := parseObjectToCfg(in) /* (struct cfg) parse object tag */
	if err != nil {
		ldb.dbLog.Println("error database insert  parseObjectToCfg failed : ",err)
		return err
	}
	ldb.dbLog.Println("Info ",cfg.showCfg())
	undoFlag := false
	if len(flag) > 0 {
		undoFlag = flag[0]
	}

	if !undoFlag {
		err = ldb.setIncrement(cfg) /* (kv.id) set increment id */
		if err != nil {
			ldb.dbLog.Println("error database insert cfg setIncrement failed : ",err,"  ","cfg is : ",cfg)
			return err
		}
	}

	dbKV := &dbKeyValue{}
	structKV(in, dbKV, cfg) /* (kv.index) all key and value*/

	//dbKV.showDbKV()
	err = ldb.insertKvToDb(dbKV) /* (kv to db) kv insert database (atomic) */
	if err != nil {
		ldb.dbLog.Println("error database insert insertKvToDb failed : ",err)
		return err
	}
	return nil
}

func (ldb *LDataBase) insertKvToDb(dbKV *dbKeyValue) error {

	ldb.dbLog.Println("Info kv database insertKvToDb dbKV is : ",dbKV)

	undoKeyValue := dbKeyValue{} /* 	Record all operations before a db operation is completed		*/
	/* 	If you need to roll back before the operation is complete 		*/
	/* 	execute the undoKV function										*/
	undoKeyValue.typeName = dbKV.typeName
	undo := false
	space := " "
	defer func() { /*	undo  */
		if !undo {
			return
		}
		ldb.undoInsertKV(&undoKeyValue)
	}()
	for _, v := range dbKV.index {
		err := saveKey(v.key, v.value, ldb.db)
		if err != nil {
			ldb.dbLog.Println("error database insertKvToDb saveKey failed : ",err,space,"key is : ",v.key,space,"value is : ",v.value)
			undo = true
			return err
		}
		undoKeyValue.index = append(undoKeyValue.index, v)
	}

	err := saveKey(dbKV.id.key, dbKV.id.value, ldb.db)
	if err != nil {
		ldb.dbLog.Println("error database insertKvToDb saveKey failed : ",err,space,"key is : ",dbKV.id.key,space,"value is : ",dbKV.id.value)
		return err
	}

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
	ldb.dbLog.Println("Info database setIncrement cfg id set is : ",id,"  ","next id is : ",ldb.nextId)
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
		ldb.dbLog.Println("error database Remove failed : ",err)
		return err
	}
	ldb.undoRemove(in)
	return nil
}

func (ldb *LDataBase) remove(in interface{}) error {

	cfg, err := parseObjectToCfg(in)
	if err != nil {
		ldb.dbLog.Println("error database remove parseObjectToCfg failed : ",err)
		return err
	}
	ldb.dbLog.Println("Info ",cfg.showCfg())
	if isZero(cfg.Id) {
		ldb.dbLog.Println("error database remove parseObjectToCfg cfg id isZero  : ",cfg.Id)
		return ErrIncompleteStructure
	}

	dbKV := &dbKeyValue{}
	structKV(in, dbKV, cfg) /* (kv.index) all key and value*/

	//dbKV.showDbKV()

	err = ldb.removeKvToDb(dbKV)
	if err != nil {
		ldb.dbLog.Println("error database remove removeKvToDb failed  : ",err,"  ","dbKV is : ",dbKV)
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
			ldb.dbLog.Println("error database  removeKvToDb removeKey failed  : ",err,"  ","key is : ",v.key)
			undo = true
			return err
		}
		undoKeyValue.index = append(undoKeyValue.index, v)
	}

	err := removeKey(dbKV.id.key, ldb.db)
	if err != nil {
		ldb.dbLog.Println("error database  removeKvToDb removeKey failed  : ",err,"  ","key is : ",dbKV.id.key)
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
		ldb.dbLog.Println("error database  Modify failed  : ",err)
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
		ldb.dbLog.Println("error database  modify failed  : func Too many parameters : ",fnType.NumIn())
		return errors.New("func Too many parameters : "+ string(fnType.NumIn()))
	}

	pType := fnType.In(0)

	if pType.Kind() != dataType.Kind() {
		ldb.dbLog.Println("error database  modify failed  : Parameter type does not match: ",pType.Kind().String()," : ",dataType.Kind().String())
		return errors.New("Parameter type does not match : " + pType.Kind().String() + " : " + dataType.Kind().String())
	}
	fnRef.Call([]reflect.Value{dataRef}) /*	call fn */
	// modify
	oldRef := reflect.ValueOf(oldInter)
	return ldb.modifyKvToDb(&oldRef, &dataRef)
}

func (ldb *LDataBase) modifyKvToDb(oldRef, newRef *reflect.Value) error {

	ldb.dbLog.Println("Info database modifyKvToDb oldRef is : ",oldRef)
	ldb.dbLog.Println("Info database modifyKvToDb newRef is : ",newRef)

	oldCfg, err := extractObjectTagInfo(oldRef)
	if err != nil {
		ldb.dbLog.Println("error database modifyKvToDb extractObjectTagInfo oldRef failed : ",err)
		return err
	}
	newCfg, err := extractObjectTagInfo(newRef)
	if err != nil {
		ldb.dbLog.Println("error database modifyKvToDb extractObjectTagInfo newRef failed : ",err)
		return err
	}

	ldb.dbLog.Println("Info database modifyKvToDb oldCfg is : ",oldCfg.showCfg())
	ldb.dbLog.Println("Info database modifyKvToDb newCfg is : ",newCfg.showCfg())

	if !reflect.DeepEqual(oldCfg.Id.Interface(), newCfg.Id.Interface()) {
		ldb.dbLog.Println("error database modifyKvToDb newCfg and oldCfg id failed : ","newCfg id is : ",newCfg.Id,"  "," oldCfg id is : ",oldCfg.Id)
		return ErrNoID
	}

	newKV := &dbKeyValue{}
	oldKV := &dbKeyValue{}
	structKV(oldRef.Interface(), oldKV, oldCfg)
	structKV(newRef.Interface(), newKV, newCfg)

	ldb.dbLog.Println("Info database modifyKvToDb structKV oldKV is : ",oldKV)
	ldb.dbLog.Println("Info database modifyKvToDb structKV newKV is : ",newKV)

	err = ldb.removeKvToDb(oldKV)
	if err != nil {
		ldb.dbLog.Println("error database modifyKvToDb removeKvToDb oldKV is : ",oldKV)
		return err
	}

	err = ldb.insertKvToDb(newKV)
	if err != nil {
		ldb.dbLog.Println("error database modifyKvToDb insertKvToDb newKV is : ",newKV)
		return err
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

	ldb.dbLog.Println("Info database undoModifyKv extractObjectTagInfo oldCfg is : ",oldRef)
	oldCfg, err := extractObjectTagInfo(&oldRef)
	if err != nil {
		ldb.dbLog.Println("error database undoModifyKv extractObjectTagInfo oldCfg failed : ",err)
		return err
	}
	id, err := rlp.EncodeToBytes(oldCfg.Id.Interface())
	if err != nil {
		ldb.dbLog.Println("error database undoModifyKv rlp EncodeToBytes id failed : ",err," : ",oldCfg.Id)
		return err
	}
	typeName := []byte(oldCfg.Name)
	key := idKey(id, typeName)
	val, err := getDbKey(key, ldb.db)
	if err != nil {
		ldb.dbLog.Println("error database undoModifyKv getDbKey  failed : ",err,"  ","key is : ",key)
		return err
	}

	dst := reflect.New(reflect.Indirect(oldRef).Type())
	err = rlp.DecodeBytes(val, dst.Interface())
	if err != nil {
		ldb.dbLog.Println("error database undoModifyKv rlp DecodeBytes failed : ",err,"  ","dst interface is : ",dst.Interface())
		return err
	}

	return ldb.modifyKvToDb(&dst, &oldRef)
}

/*

The following are some specific implementations of the features that may change, please do not use

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
	ldb.dbLog.Println("Info database find : ","tagName is : ",tagName)
	fieldName := []byte(tagName)
	fields, err := getFieldInfo(tagName, value)
	if err != nil {
		ldb.dbLog.Println("error database find  getFieldInfo failed : ",err)
		return err
	}

	ldb.dbLog.Println("Info database find fields is : ",fields)

	typeName := []byte(fields.typeName)

	suffix := getFieldValue(fields)
	if len(suffix) == 0 {
		ldb.dbLog.Println("error database find nonUniqueValue failed : ",err)
		return ErrNotFound
	}

	key := typeNameFieldName(typeName, fieldName)
	key = append(key, suffix...)
	return ldb.findFields(key, typeName, to)
}

func (ldb *LDataBase) findFields(key, typeName []byte, to interface{}) error {
	v, err := getDbKey(key, ldb.db)
	if err != nil {
		ldb.dbLog.Println("error database findFields key is : ",key)
		return err
	}
	return ldb.findDbObject(v, typeName, to)
}

// only key is id can be called
func (ldb *LDataBase) findDbObject(key, typeName []byte, to interface{}) error {

	id := idKey(key, typeName)
	ldb.dbLog.Println("Info database findDbObject id is : ",id)
	val, err := getDbKey(id, ldb.db)
	if err != nil {
		ldb.dbLog.Println("error database findDbObject getDbKey failed : ",err)
		return err
	}
	ldb.dbLog.Println("Info database findDbObject val is : ",val)
	err = rlp.DecodeBytes(val, to)
	if err != nil {
		ldb.dbLog.Println("error database findDbObject rlp DecodeBytes failed : ",err)
		return err
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
	begin := typeNameFieldName(typeName, fieldName)


	end := getNonUniqueEnd(begin)
	space := "  "
	ldb.dbLog.Println("getIndex typeName : ",typeName,space,"fieldName : ",fieldName,space,"begin : ",begin,space,"end : ",end,space,"greater : ",fields.greater)
	it := newMultiIndex(typeName, fieldName, begin, end, fields.greater, ldb)
	return it, nil
}

/*

The value corresponding to the given key is retrieved
from the db and returned to the caller
which may not exist

*/
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

func (ldb *LDataBase) GetMutableIndex(fieldName string, in interface{}) (*MultiIndex, error) {
	return ldb.GetIndex(fieldName, in)
}

func (ldb *LDataBase) lowerBound(begin, end, fieldName []byte, data interface{}, greater bool) (*DbIterator, error) {
	//TODO

	fields, err := getFieldInfo(string(fieldName), data)
	if err != nil {
		return nil, err
	}

	prefix := getFieldValue(fields)

	if len(prefix) != 0 {
		begin = append(begin, prefix...)
		it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)

		idx, err := newDbIterator([]byte(fields.typeName), it, ldb.db, greater)
		if err != nil {
			return nil, err
		}
		return idx, nil
	}

	return nil, ErrNotFound
}

func (ldb *LDataBase) Empty(begin, end, fieldName []byte) bool {

	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
	defer it.Release()
	if it.Next() {
		return false
	}
	return true
}

func (ldb *LDataBase) IteratorTo(begin, end, fieldName []byte, in interface{}, greater bool) (*DbIterator, error) {
	fields, err := getFieldInfo(string(fieldName), in)
	if err != nil {
		return nil, err
	}
	prefix := getFieldValue(fields)

	if len(prefix) == 0 {
		return nil, errors.New("Get Field Value Failed")
	}

	key := []byte{}
	key = append(begin, prefix...)
	fmt.Println(key)
	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
	if !it.Seek(key) {
		return nil, errors.New("Iterator To Not Found")
	}
	k := idKey(it.Value(), []byte(fields.typeName))
	val, err := getDbKey(k, ldb.db)
	if err != nil {
		return nil, err
	}

	//
	itr := &DbIterator{it: it, greater: fields.greater, db: ldb.db, first: false, value: val, typeName: []byte(fields.typeName)}
	return itr, nil
}

func (ldb *LDataBase) BeginIterator(begin, end, fieldName, typeName []byte, greater bool) (*DbIterator, error) {

	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
	if greater {
		if !it.Last() {
			ldb.dbLog.Println("error DataBase BeginIterator : Greater True : Last Failed")
			return nil,ErrNotFound
		}
	} else {
		if !it.Next() {
			ldb.dbLog.Println("error DataBase BeginIterator : Greater False : Next Failed")
			return nil,ErrNotFound
		}
	}

	k := idKey(it.Value(), []byte(typeName))
	val, err := getDbKey(k, ldb.db)
	if err != nil {
		ldb.dbLog.Println("error DataBase BeginIterator : " ,err)
		return nil, errors.New("DataBase BeginIterator : " + err.Error())
	}

	itr := &DbIterator{it: it, greater: greater, db: ldb.db, first: false, typeName: typeName, value: val}

	return itr, nil
}

func (ldb *LDataBase) upperBound(begin, end, fieldName []byte, data interface{}, greater bool) (*DbIterator, error) {
	//TODO
	fields, err := getFieldInfo(string(fieldName), data)
	if err != nil {
		return nil, err
	}

	prefix := getFieldValue(fields)

	if len(prefix) != 0 {
		begin = append(begin, prefix...)
	}

	begin[len(begin)-1] = begin[len(begin)-1] + 1
	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)

	idx, err := newDbIterator([]byte(fields.typeName), it, ldb.db, greater)
	if err != nil {
		return nil, err
	}
	return idx, nil
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
		log.Println("undo session empty")
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
		log.Println("undo session empty")
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
		log.Println("undo session empty")
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

func (ldb *LDataBase) undoInsertKV(undoKeyValue *dbKeyValue) { /*	insert kv to db error --> undo kv */

	ldb.dbLog.Println("info database undoInsertKV is : ",undoKeyValue)
	if len(undoKeyValue.index) == 0 {
		return
	}
	space := " "
	for _, v := range undoKeyValue.index { /* 	undo index*/
		err := removeKey(v.key, ldb.db)
		if err != nil {
			ldb.dbLog.Println("error database undoInsertKV removeKey failed : ",err,space,"key is : ",v.key)
			return
		}
	}
	if len(undoKeyValue.id.key) == 0 {
		return
	}
	err := removeKey(undoKeyValue.id.key, ldb.db) /* 	undo id*/
	if err != nil {
		ldb.dbLog.Println("error database undoInsertKV removeKey failed : ",err,space,"key is : ",undoKeyValue.id.key)
		return
	}
}

func (ldb *LDataBase) undoRemoveKV(undoKeyValue *dbKeyValue) { /*	remove kv from db error --> undo kv */

	ldb.dbLog.Println("info database undoRemoveKV  undoKeyValue is ",undoKeyValue)

	if len(undoKeyValue.index) == 0 {
		return
	}
	for _, v := range undoKeyValue.index { /* 	undo index*/
		err := saveKey(v.key, v.value, ldb.db)
		if err != nil {
			ldb.dbLog.Println("error database undoRemoveKV saveKey failed : ",err,"  ","key is : ",v.key,"  ","value is : ",v.value)
			return
		}
	}
	if len(undoKeyValue.id.key) == 0 {
		return
	}

	err := saveKey(undoKeyValue.id.key, undoKeyValue.id.value, ldb.db) /* undo id*/
	if err != nil {
		ldb.dbLog.Println("error database undoRemoveKV saveKey failed : ",err,"  ","key is : ",undoKeyValue.id.key,"  ","value is : ",undoKeyValue.id.value)
		return
	}
}
