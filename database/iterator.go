
package database

import (
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"reflect"
	"regexp"
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


	Data(data interface{})error

	Begin() []byte
}

//Do not use the functions in this file
type DbIterator struct{
	key 		[]byte
	value 		[]byte
	begin 		[]byte
	typeName 	[]byte
	db 			*leveldb.DB // TODO interface
	it 			iterator
	rege 		string // FIXME unused
	greater		bool
}

//Do not use the functions in this file
func newDbIterator (typeName []byte,it iterator,db *leveldb.DB,rege string,greater bool) (*DbIterator,error){
	if greater {
		if it.Last(){
			idx := &DbIterator{typeName:typeName,it:it,db:db,rege:rege,greater:greater}
			key := idKey(it.Value(),typeName)
			key, err := getDbKey(key,db)
			if err != nil{
				return nil,err
			}
			idx.begin = make([]byte,len(key))
			copy(idx.begin,key)
			return idx,nil
		}

		return nil, ErrNotFound
	}
	for it.Next(){
		idx := &DbIterator{typeName:typeName,it:it,db:db,rege:rege,greater:greater}
		idx.begin = make([]byte,len(it.Value()))
		copy(idx.begin,it.Value())
		return idx,nil
	}
	return nil, ErrNotFound
}

/* Do not use the functions in this file */
/////////////////////////////////////////////// iterator //////////////////////////////////
func (index *DbIterator) Next() bool {
	if index.greater {
		if len(index.key) == len(index.value) && len(index.key) == 0{
			value := index.it.Value()
			k := idKey(value, index.typeName)

			v, err := index.db.Get(k, nil)
			if err != nil {
				return false
			}
			index.key = k
			index.value = v
			return true
		}
		return index.prev()
	}
	if len(index.key) == len(index.value) && len(index.key) == 0{
		value := index.it.Value()
		k := idKey(value, index.typeName)

		v, err := index.db.Get(k, nil)
		if err != nil {
			return false
		}
		index.key = k
		index.value = v
		return true
	}
	return index.next()
}

func (index *DbIterator) Prev() bool {
	if index.greater {
		if len(index.key) == len(index.value) && len(index.key) == 0{
			value := index.it.Value()
			k := idKey(value, index.typeName)

			v, err := index.db.Get(k, nil)
			if err != nil {
				return false
			}
			index.key = k
			index.value = v
			return true
		}
		return index.next()
	}
	if len(index.key) == len(index.value) && len(index.key) == 0{
		value := index.it.Value()
		k := idKey(value, index.typeName)

		v, err := index.db.Get(k, nil)
		if err != nil {
			return false
		}
		index.key = k
		index.value = v
		return true
	}
	return index.prev()
}

func (index *DbIterator) Begin()[]byte{
	return index.begin
}

func (index *DbIterator) prefix()[]byte{
	for index.it.Next() {
		reg := regexp.MustCompile(index.rege)
		find := reg.Match(index.it.Key())
		if !find {
			continue
		}
		return index.it.Key()
	}
	return nil
}

func (index *DbIterator) next() bool {
	for index.it.Next() {
		value := index.it.Value()
		k := idKey(value, index.typeName)
		v, err := index.db.Get(k, nil)
		if err != nil {
			index.value = nil
			return false
		}
		index.key = k
		index.value = v
		return true
	}
	index.value = nil
	return false
}

func (index *DbIterator) prev() bool {
	for index.it.Prev() {
		value := index.it.Value()
		k := idKey(value, index.typeName)
		v, err := index.db.Get(k, nil)
		if err != nil {
			index.value = nil
			return false
		}
		index.key = k
		index.value = v
		return true
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

func (index *DbIterator)Seek(key []byte) bool{// TODO
	return index.it.Seek(key)
}
