package database

import (
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"reflect"
)

//Do not use the functions in this file
type iterator interface {
	First() bool

	Last() bool

	Next() bool

	Prev() bool

	Key() []byte

	Value() []byte

	Seek(key []byte) bool

	Release()
}

//Do not use the functions in this file
type Iterator interface {
	iterator

	Data(data interface{}) error

	Begin() []byte
}

//Do not use the functions in this file
type DbIterator struct {
	key      []byte
	value    []byte
	begin    []byte
	typeName []byte
	db       *leveldb.DB // TODO interface
	it       iterator
	greater  bool
	first    bool
}

//Do not use the functions in this file
func newDbIterator(typeName []byte, it iterator, db *leveldb.DB, greater bool) (*DbIterator, error) {
	if greater {
		if it.Last() {
			idx := &DbIterator{typeName: typeName, it: it, db: db, greater: greater}
			key := idKey(it.Value(), typeName)
			key, err := getDbKey(key, db)
			if err != nil {
				return nil, err
			}

			idx.copyBeginValue(key)

			return idx, nil
		}

		return nil, ErrNotFound
	}
	for it.Next() {
		idx := &DbIterator{typeName: typeName, it: it, db: db,greater: greater}

		key := idKey(it.Value(), typeName)
		key, err := getDbKey(key, db)
		if err != nil {
			return nil, err
		}

		idx.copyBeginValue(key)
		return idx, nil
	}

	return nil, ErrNotFound
}

/* Do not use the functions in this file */


func (index *DbIterator)copyBeginValue(key []byte)error{
	index.begin = make([]byte, len(key)) /* compare begin 	*/
	copy(index.begin, key)
	index.value = make([]byte, len(key)) /* begin data  	*/
	copy(index.value, key)
	index.first = true
	return nil
}

func (index *DbIterator)keyValue(key []byte)error{
	k := idKey(key, index.typeName)

	v, err := index.db.Get(k, nil)
	if err != nil {
		index.key = nil
		index.value = nil
		return err
	}
	index.key = k
	index.value = v
	return nil
}

func (index *DbIterator) Next() bool {
	if index.greater {
		if index.first == true {
			index.first = false

			return  index.keyValue(index.it.Value())  == nil
		}
		return index.prev()
	}
	if index.first == true {
		index.first = false

		return  index.keyValue(index.it.Value())  == nil
	}
	return index.next()
}

func (index *DbIterator) Prev() bool {
	if index.greater {
		if index.first == true {

			index.first = false

			return  index.keyValue(index.it.Value())  == nil
		}
		return index.next()
	}

	if index.first == true {

		index.first = false

		return  index.keyValue(index.it.Value())  == nil
	}
	return index.prev()
}

func (index *DbIterator) Begin() []byte {
	return index.begin
}

func (index *DbIterator) next() bool {
	for index.it.Next() {

		return index.keyValue(index.it.Value()) == nil
	}
	index.value = nil
	return false
}

func (index *DbIterator) prev() bool {
	for index.it.Prev() {

		return index.keyValue(index.it.Value()) == nil
	}
	index.value = nil
	return false
}

func (index *DbIterator) Release() {
	index.it.Release()
}

func (index *DbIterator) Data(data interface{}) error {
	ref := reflect.ValueOf(data)
	if !ref.IsValid() || reflect.Indirect(ref).Kind() != reflect.Struct {
		return ErrStructPtrNeeded
	}

	rv := reflect.Indirect(ref)
	if !rv.CanAddr() {
		return ErrPtrNeeded
	}

	return rlp.DecodeBytes(index.Value(), data)
}

func (index *DbIterator) Key() []byte {
	return index.key
}

func (index *DbIterator) Value() []byte {
	return index.value
}

func (index *DbIterator) Last() bool { // TODO
	return index.it.Last()
}

func (index *DbIterator) First() bool { // TODO
	return index.it.First()
}

func (index *DbIterator) Seek(key []byte) bool { // TODO
	return index.it.Seek(key)
}
