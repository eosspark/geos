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
	db *leveldb.DB
	path string
}

func NewDataBase(path string) (*LDataBase, error) {

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

	return &LDataBase{db: db, path:path}, nil
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

@parameters
data 		--> 	object

@return
success 	-->		nil
error 		-->		error object

*/

func (ldb *LDataBase) Insert(data interface{}) error {
	tx, err := ldb.db.OpenTransaction()
	if err != nil {
		return err
	}
	defer tx.Commit()

	err = insert(data, tx)
	if err != nil {
		tx.Discard()
		return err
	}
	return nil
}


//////////////////////////////////////////////////////	find object from database //////////////////////////////////////////////////////
/*

@parameters
fieldName 	--> 	rule
value 		--> 	object
out 		-->		outPut parameters
@return
success 	-->		nil 	(out valid)
error 		-->		error 	(out invalid)

*/

func (ldb *LDataBase) Find(fieldName string, data interface{},out interface{}) error {
	return find(fieldName, data, out,ldb.db)
}


//////////////////////////////////////////////////////	get object from database //////////////////////////////////////////////////////
/*

@parameters
fieldName 	--> 	rule
value 		--> 	object

@return
success 	-->		iterator
error 		-->		error

*/
func (ldb *LDataBase) Get(fieldName string, data interface{}) (DbIterator, error) {
	return get(fieldName, data, ldb.db)
}


//////////////////////////////////////////////////////	get object from database //////////////////////////////////////////////////////
/*

@parameters
data 		--> 	old object
fn 			--> 	rule

@return
success 	-->		nil
error 		-->		error

*/
func (ldb *LDataBase) Modify(data interface{}, fn interface{}) error {

	tx, err := ldb.db.OpenTransaction()
	if err != nil {
		return err
	}
	defer tx.Commit()
	err = update(data, fn, tx)
	if err != nil {
		tx.Discard()
		return err
	}
	return nil
}

//////////////////////////////////////////////////////	get object from database //////////////////////////////////////////////////////
/*

@parameters
data 		--> 	object

@return
success 	-->		nil
error 		-->		error

*/
func (ldb *LDataBase) Remove(data interface{}) error {
	return delete_(data, ldb.db)
}


////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
func insert(data interface{}, tx *leveldb.Transaction) error {

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
	//	cfg.showStructInfo()			// XXX
	id, err := numbertob(cfg.Id.Interface())
	typeName := []byte(cfg.Name)

	callBack := func(key, value []byte) error {
		return save(key, value, tx)
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

	return save(key, value, tx)
}

func save(key, value []byte, tx *leveldb.Transaction) error {

	if ok, _ := tx.Has(key, nil); ok {
		return ErrAlreadyExists
	}

	err := tx.Put(key, value, nil)
	if err != nil {
		return err
	}
	return nil
}

func delete_(data interface{}, db *leveldb.DB) error {

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

	if isZero(cfg.Id){
		return ErrIncompleteStructure
	}
	id, err := numbertob(cfg.Id.Interface())

	typeName := []byte(cfg.Name)

	//fmt.Println(typeName)

	callBack := func(key, value []byte) error {
		exist, err := db.Has(key, nil)
		if err != nil {
			return nil
		}
		if !exist {
			return ErrNotFound
		}

		return remove(key, db)
	}

	err = doCallBack(id, typeName, cfg, callBack)// FIXME --> id --> obj --> cfg
	if err != nil {
		return err
	}
	key := idKey(id, typeName)

	return remove(key, db)
}

func remove(key []byte, db *leveldb.DB) error {

	if ok, _ := db.Has(key, nil); !ok {
		return ErrNotFound
	}
	err := db.Delete(key, nil)
	if err != nil {
		return err
	}
	return nil
}

func update(data interface{}, fn interface{}, tx *leveldb.Transaction) error {
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
	err := modify(&oldRef, &dataRef, tx)
	if err != nil {
		return err
	}

	return nil
}

func modify(old, new *reflect.Value, tx *leveldb.Transaction) error {
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
		find, err := tx.Has(oldKey, nil)
		if err != nil {
			return err
		}
		if !find {
			return ErrNotFound
		}
		value, err := tx.Get(oldKey, nil)
		if err != nil {
			return err
		}
		err = tx.Delete(oldKey, nil)
		if err != nil {
			return err
		}
		return save(newKey, value, tx)
	}

	id, err := numbertob(newCfg.Id.Interface())
	typeName := []byte(newCfg.Name)
	key := idKey(id, typeName)
	val, err := rlp.EncodeToBytes(new.Interface())
	if err != nil {
		return err
	}

	err = modifyField(newCfg, oldCfg, callBack)
	if err != nil{
		return err
	}

	err = tx.Delete(key, nil)
	if err != nil {
		return err
	}

	return save(key, val, tx)
}


func find(fieldName string, value interface{},to interface{}, db *leveldb.DB) error{

	fields,err := getFieldInfo(fieldName,value)
	if err != nil{
		return err
	}

	typeName := fields.typeName
	if !fields.unique{
		return ErrNotFound
	}

	suffix := nonUniqueValue(fields)
	if suffix == nil{
		return ErrNotFound
	}
/*
	unique --> typename__fieldName__fieldValue
*/
	key := typeNameFieldName(typeName,fieldName)

	key = append(key, suffix...)

	v ,err:= getDbKey(key,db)
	if err != nil {
		return  err
	}

	id := idKey(v,[]byte(typeName))

	val ,err:= getDbKey(id,db)
	if err != nil {
		return  err
	}
	err = rlp.DecodeBytes(val,to)
	if err != nil {
		return  err
	}
	return nil
}


func get(fieldName string, value interface{}, db *leveldb.DB) (DbIterator, error) {

	fields,err := getFieldInfo(fieldName,value)
	if err != nil{
		return nil,err
	}

	typeName := fields.typeName

	rege := getNonUniqueFieldValue(fields)
	if rege == nil{
		return nil,ErrNoID
	}

	key :=typeNameFieldName(typeName,fieldName)
	key = append(key, '_')
	key = append(key, '_')

	if !fields.unique {
		/*
			index --> typename__fieldName__
		*/
		end := indexEnd(key)
		iter := db.NewIterator(&util.Range{Start: key, Limit: end}, nil)

		it := newIndexIterator([]byte(typeName), iter, db, string(rege), fields.greater)
		return it, nil
	}

	return nil, ErrNotFound
}

//////////////////////////////////////////////////////	save increment id to database  //////////////////////////////////////////////////////
/*
key 	 --> typeName__fieldName
value --> id
*/

func incrementField(cfg *structInfo, tx *leveldb.Transaction) error {

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

func getDbKey(key []byte,db *leveldb.DB)([]byte,error){
	exits,err := db.Has(key,nil)
	if err != nil {
		return  nil,err
	}
	if !exits{
		return nil,ErrNotFound
	}

	val ,err:= db.Get(key,nil)
	if err != nil {
		return  nil,err
	}
	return val,err
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
