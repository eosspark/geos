package database

import (
	"github.com/syndtr/goleveldb/leveldb"
	"reflect"
	"fmt"
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
	End() bool
	Begin() bool
}

const (
	itBEGIN   = "iterator begin"
	itCURRENT = "iterator current"
	itEND     = "iterator end"
)

//Do not use the functions in this file
type DbIterator struct {
	key           []byte
	value         []byte
	begin         []byte
	currentStatus string
	typeName      []byte
	db            *leveldb.DB // TODO interface
	it            iterator
	first         bool
}

//Do not use the functions in this file
func newDbIterator(typeName []byte, it iterator, db *leveldb.DB) (*DbIterator, error) {

	idx := &DbIterator{typeName: typeName, it: it, db: db}

	return idx, nil
}

/* Do not use the functions in this file */

func (index *DbIterator) clearKV() {
	index.key = nil
	index.value = nil
}

func (index *DbIterator) setKV(k,v[]byte) {
	index.key = k
	index.value = v
}

func (index *DbIterator) keyValue(key []byte) error {

	k := splicingString(index.typeName,key)
	v,err := getDbKey(k,index.db)
	if err != nil {
		return err
	}

	index.clearKV()
	index.setKV(k,v)
	return nil
}

func (index *DbIterator) Next() bool {


	return index.next()
}

func (index *DbIterator) Prev() bool {


	return index.prev()
}

func (index *DbIterator) next() bool {
	for index.it.Next() {
		index.currentStatus = itCURRENT
		return index.keyValue(index.it.Value()) == nil
	}

	index.currentStatus = itEND
	return false
}

func (index *DbIterator) prev() bool {
	for index.it.Prev() {
		index.currentStatus = itCURRENT
		return index.keyValue(index.it.Value()) == nil
	}

	index.currentStatus = itBEGIN
	return false
}

func (index *DbIterator) End() bool {
	if index.it == nil {
		return false
	}
	if index.db == nil {
		return false
	}
	return index.currentStatus == itEND
}

func (index *DbIterator) Begin() bool {
	if index.it == nil {
		return false
	}
	if index.db == nil {
		return false
	}
	return index.currentStatus == itBEGIN
}

func (index *DbIterator) Release() {
	index.it.Release()
	index.clearKV()
	index.typeName = nil
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

	index.keyValue(index.it.Value())
	key := index.Value()
	err := DecodeBytes(key, data)
	if err != nil{
		//
		fmt.Println(err)
	}
	return nil
}

func (index *DbIterator) Key() []byte {
	return index.key
}

func (index *DbIterator) Value() []byte {
	index.keyValue(index.it.Value())
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
