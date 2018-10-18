package database

import (
	"fmt"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"reflect"
)


type LDataBase struct {
	db   *leveldb.DB
	path string
}


/*

@param path 			--> 	database file (note:type-->d)

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

	return &LDataBase{db: db, path: path}, nil
}

func (ldb *LDataBase) Close() {
	err := ldb.db.Close()
	if err != nil {
		// log
	} else {
		// log
	}
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
	return nil
}

//////////////////////////////////////////////////////	find object from database //////////////////////////////////////////////////////
/*

@param fieldName 	--> 	rule
@param in 			--> 	object
@param out 			-->		output(pointer)

@return
success 			-->		nil 	(out valid)
error 				-->		error 	(out invalid)

*/

func (ldb *LDataBase) Find(fieldName string, in interface{}, out interface{}) error {
	return find(fieldName, in, out, ldb.db)
}

//////////////////////////////////////////////////////	get multiIndex from database //////////////////////////////////////////////////////
/*

@param fieldName 	--> 	rule
@param in 			--> 	object

@return
success 			-->		iterator
error 				-->		error

*/
func (ldb *LDataBase) GetIndex(fieldName string, in interface{}) (*multiIndex, error) {
	return getIndex(fieldName, in,ldb)
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

	err := modify(old, fn, ldb.db)
	if err != nil {
		// undo
		return err
	}
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
	return remove(in, ldb.db)
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
	//	cfg.showStructInfo()			// XXX id, err := numbertob(cfg.Id.Interface())
	id, err := numbertob(cfg.Id.Interface())
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
	id, err := numbertob(cfg.Id.Interface())

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

	oldInter := copyInterface(data)
	dataType := dataRef.Type()

	fnRef := reflect.ValueOf(fn)
	fnType := fnRef.Type()
	if fnType.NumIn() != 1 {
		return errors.New("func Too many parameters")
	}

	pType := fnType.In(0)

	if pType.Kind() != dataType.Kind() {
		fmt.Println(pType.String(), " <--> ", dataType.String())
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

	id, err := numbertob(newCfg.Id.Interface())
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

func find(fieldName string, value interface{}, to interface{}, db *leveldb.DB) error {

	fields, err := getFieldInfo(fieldName, value)
	if err != nil {
		return err
	}

	typeName := fields.typeName
	if !fields.unique {
		return ErrNotFound
	}

	suffix := nonUniqueValue(fields)
	if suffix == nil {
		return ErrNotFound
	}
	/*
		unique --> typename__fieldName__fieldValue
	*/
	key := typeNameFieldName(typeName, fieldName)

	key = append(key, suffix...)

	v, err := getDbKey(key, db)
	if err != nil {
		return err
	}

	id := idKey(v, []byte(typeName))

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

func getIndex(fieldName string, value interface{}, db DataBase) (*multiIndex, error) {

	fields, err := getFieldInfo(fieldName, value)
	if err != nil {
		return nil, err
	}

	typeName := fields.typeName

	begin := typeNameFieldName(typeName, fieldName)
	begin = append(begin, '_')
	begin = append(begin, '_')

	if !fields.unique {
		/*
			index --> typename__fieldName__
		*/
		endl := indexEnd(begin)
		//begin[len(begin)-1] = begin[len(begin)-1] - 1
		it := newMultiIndex([]byte(typeName),[]byte(fieldName),begin,endl,fields.greater,db)
		return it, nil
	}

	return nil, ErrNotFound
}

/*
key 	 --> typeName__fieldName
value --> id
*/

func incrementField(cfg *structInfo, tx *leveldb.DB) error {

	typeName := []byte(cfg.Name)
	fieldName := []byte(tagID)
	// typeName__fieldName
	key := append(typeName, '_')
	key = append(key, '_')
	key = append(key, fieldName...)

	valByte, err := tx.Get(key, nil)
	if err != nil && err != leveldb.ErrNotFound {
		return err
	}
	counter := cfg.IncrementStart
	if valByte != nil {
		counter, err = numberfromb(valByte)
		if err != nil {
			return err
		}
		//fmt.Println(key ,"  ","found id : ",counter)
		counter++
	}

	cfg.Id.Set(reflect.ValueOf(counter).Convert(cfg.Id.Type()))
	value, err := numbertob(cfg.Id.Interface())
	if value == nil && err == nil {
		return err
	}
	return tx.Put(key, value, nil)
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

func copyInterface(data interface{}) interface{} {

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

func (ldb *LDataBase) lowerBound(begin,end,fieldName []byte,data interface{},greater bool) (*DbIterator,error){
	//TODO
	fields, err := getFieldInfo(string(fieldName), data)
	if err != nil {
		return nil, err
	}

	reg ,prefix:= getNonUniqueFieldValue(fields)
	if reg == nil {
		return nil, ErrNoID
	}
	sift := string(reg)
//	fmt.Println(begin)
	if len(prefix) != 0{
		begin = append(begin,prefix...)
		it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
		//for it.Next(){
		//	fmt.Println(it.Key())
		//}
		idx,err := newDbIterator([]byte(fields.typeName),it,ldb.db,sift,greater)
		if err != nil{
			return nil,err
		}
		return idx,nil
	}

	return nil,ErrNotFound
}

func (ldb *LDataBase) upperBound(begin,end,fieldName []byte,data interface{},greater bool) (*DbIterator,error){
	//TODO
	fields, err := getFieldInfo(string(fieldName), data)
	if err != nil {
		return nil, err
	}

	reg ,prefix := getNonUniqueFieldValue(fields)
	if reg == nil {
		return nil, ErrNoID
	}
	if len(prefix) != 0{
		begin = append(begin,prefix...)
	}

	sift := string(reg)
	begin[len(begin) - 1] = begin[len(begin) - 1] + 1
	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)

	idx,err := newDbIterator([]byte(fields.typeName),it,ldb.db,sift,greater)
	if err != nil{
		return nil,err
	}
	return idx,nil
}
