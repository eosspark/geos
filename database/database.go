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

type DataBase interface {
	Insert(data interface{})

	Find(fieldName string, data interface{}) (DbIterator, error)

	Get(fieldName string, data interface{}) (Iterator, error)

	Modify(data interface{}, fn interface{}) error

	Remove(data interface{}) error
}

type LDataBase struct {
	db *leveldb.DB
	fn string
}

func NewDataBase(fileName string) (*LDataBase, error) {

	db, err := leveldb.OpenFile(fileName, &opt.Options{
		OpenFilesCacheCapacity: 16,
		BlockCacheCapacity:     16 / 2 * opt.MiB,
		WriteBuffer:            16 / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(fileName, nil)
	}
	if err != nil {
		return nil, err
	}

	return &LDataBase{db: db, fn: fileName}, nil
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
	return insert(data, ldb.db)
}

//////////////////////////////////////////////////////	find object from database //////////////////////////////////////////////////////
/*

@parameters
fieldName 	--> 	rule
value 		--> 	object

@return
success 	-->		iterator
error 		-->		error object

*/

func (ldb *LDataBase) Find(fieldName string, data interface{}) (DbIterator, error) {
	return find(fieldName, data, ldb.db)
}

func (ldb *LDataBase) Get(fieldName string, data interface{}) (DbIterator, error) {
	return find(fieldName, data, ldb.db)
}

func (ldb *LDataBase) Modify(data interface{}, fn interface{}) error {
	return update(data, fn, ldb.db)
}

func (ldb *LDataBase) Remove(data interface{}) error {
	return delete_(data, ldb.db)
}

func insert(data interface{}, db *leveldb.DB) error {

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

	err = incrementField(cfg, db)
	if err != nil {
		return err
	}
	//	cfg.showStructInfo()			// XXX
	id, err := numbertob(cfg.Id.Interface())
	typeName := []byte(cfg.Name)

	callBack := func(key, value []byte) error {
		return save(key, value, db)
	}
	err = fieldIndex(id, typeName, cfg, callBack)
	if err != nil {
		return err
	}
	key := idKey(id, typeName)
	value, err := rlp.EncodeToBytes(data)
	if err != nil {
		return err
	}

	return save(key, value, db)
}

func save(key, value []byte, db *leveldb.DB) error {

	if ok, _ := db.Has(key, nil); ok {
		return ErrAlreadyExists
	}

	err := db.Put(key, value, nil)
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

	id, err := numbertob(cfg.Id.Interface())

	typeName := []byte(cfg.Name)

	//fmt.Println(typeName)

	callBack := func(key, value []byte) error {
		exist, err := db.Has(key, nil)
		if err != nil {
			return nil
		}
		if !exist {
			return ErrIncompleteStructure
		}

		return remove(key, db)
	}

	err = fieldIndex(id, typeName, cfg, callBack)
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

func update(data interface{}, fn interface{}, db *leveldb.DB) error {
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
	err := modify(&oldRef, &dataRef, db)
	if err != nil {
		return err
	}

	return nil
}

func modify(old, new *reflect.Value, db *leveldb.DB) error {
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
		save(newKey, value, db)
		return nil
	}

	id, err := numbertob(newCfg.Id.Interface())
	typeName := []byte(newCfg.Name)
	key := idKey(id, typeName)
	val, err := rlp.EncodeToBytes(new.Interface())
	if err != nil {
		return err
	}
	modifyField(newCfg, oldCfg, callBack)
	err = db.Delete(key, nil)
	if err != nil {
		return err
	}
	err = save(key, val, db)
	if err != nil {
		return err
	}
	return nil
}

func modifyField(cfg, oldCfg *structInfo, callBack func(newKey, oldKey []byte) error) error {

	id, err := numbertob(cfg.Id.Interface())
	if err != nil {
		return err
	}
	typeName := []byte(cfg.Name)

	for tag, fieldCfg := range cfg.Fields {
		// typeName__
		key := append(typeName, '_')
		key = append(key, '_')
		// typeName__tag__
		key = append(key, tag...)

		oldKey := fieldKey(key, fieldCfg)
		newKey := fieldKey(key, oldCfg.Fields[tag])

		newKey = append(newKey, id...)
		oldKey = append(oldKey, id...)

		err := callBack(newKey, oldKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func fieldKey(key []byte, info *fieldInfo) []byte {

	for _, v := range info.fieldValue {
		// typeName__tag__fieldValue...
		key = append(key, '_')
		key = append(key, '_')
		value, err := rlp.EncodeToBytes(v.Interface())
		if err != nil {
			return nil
		}
		key = append(key, value...)
	}

	//fmt.Println("func fieldKey value is : ",string(key))
	//fmt.Println("func fieldKey value is : ",key)
	return key
}

func fieldIndex(id, typeName []byte, cfg *structInfo, callBack func(key, value []byte) error) error {
	for tag, fieldCfg := range cfg.Fields {
		// typeName__
		key := append(typeName, '_')
		key = append(key, '_')
		// typeName__tag__
		key = append(key, tag...)
		key = fieldKey(key, fieldCfg)
		key = append(key, id...)
		//fmt.Println("func fieldIndex value is : ",string(key))
		//fmt.Println("func fieldIndex value is : ",key)
		err := callBack(key, id)
		if err != nil {
			return err
		}
	}
	return nil
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

func get(fieldName string, value interface{}, db *leveldb.DB) {

}

func find(fieldName string, value interface{}, db *leveldb.DB) (DbIterator, error) {

	ref := reflect.ValueOf(value)
	if !ref.IsValid() || reflect.Indirect(ref).Kind() != reflect.Struct {
		return nil, ErrBadType
	}
	if ref.Kind() == reflect.Ptr {
		return nil, ErrStructNeeded
	}
	typeName := ref.Type().Name()
	cfg, err := extractStruct(&ref)
	if err != nil {
		return nil, err
	}

	fields, ok := cfg.Fields[fieldName]
	if !ok {
		return nil, ErrNotFound
	}
	rege := ""
	for _, v := range fields.fieldValue {
		rege += "__"
		if isZero(v) {
			rege += "(.*)"
			continue
		}
		re, err := rlp.EncodeToBytes(v.Interface())
		if err != nil {
			return nil, err
		}

		rege += string(re)
	}

	key := []byte(typeName)
	key = append(key, '_')
	key = append(key, '_')
	key = append(key, []byte(fieldName)...)
	key = append(key, '_')
	key = append(key, '_')

	if fields.unique {
		/*
			unique --> typename__fieldName__fieldValue
		*/
		rege = string(key)
		end := make([]byte, len(key))
		copy(end, key)
		end[len(end)-1] = end[len(end)-1] + 1

		iter := db.NewIterator(&util.Range{Start: key, Limit: end}, nil)

		value, err := rlp.EncodeToBytes(fields.fieldValue[0].Interface())
		if err != nil {
			return nil, nil
		}
		key = append(key, value...)

		ok := iter.Seek(key)
		if !ok {
			return nil, ErrNotFound
		}
		/*	fmt.Println("unique seek key is : ",key) */
		it := newUniqueIterator([]byte(typeName), iter, db, rege)
		return it, nil
	} else {
		/*
			index --> typename__fieldName__
		*/
		end := make([]byte, len(key))
		copy(end, key)
		end[len(end)-1] = end[len(end)-1] + 1

		iter := db.NewIterator(&util.Range{Start: key, Limit: end}, nil)
		//if iter.Next(){
		//	fmt.Println(iter.Key())
		//	fmt.Println(string(iter.Key()))
		//}
		it := newIndexIterator([]byte(typeName), iter, db, rege, fields.greater)
		return it, nil
	}

	return nil, nil
}

//////////////////////////////////////////////////////	save increment id to database  //////////////////////////////////////////////////////
/*
key 	 --> typeName__fieldName
value --> id
*/

func incrementField(cfg *structInfo, db *leveldb.DB) error {

	typeName := []byte(cfg.Name)
	fieldName := []byte(tagID)
	// typeName__fieldName
	key := append(typeName, '_')
	key = append(key, '_')
	key = append(key, fieldName...)

	valByte, err := db.Get(key, nil)
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
	return db.Put(key, value, nil)
}

/////////////////////////////////////////////////////// Session  //////////////////////////////////////////////////////////
type Session struct {
	db      *DataBase
	version uint64
	apply   bool
}

func (session *Session) Commit() {
	if !session.apply {
		// log ?
		return
	}
	//	version := session.version
	//	session.db.commit(version)
	session.apply = false
}

func (session *Session) Squash() {
	if !session.apply {
		return
	}
	//	session.db.squash()
	session.apply = false
}

func (session *Session) Undo() {
	if !session.apply {
		return
	}
	//	session.db.undo()
	session.apply = false
}
