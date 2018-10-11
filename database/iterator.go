
package database



import (
	"github.com/syndtr/goleveldb/leveldb"
)

type Iterator interface {

	First() bool

	Last() bool

	Next() bool

	Prev() bool

	Key() []byte

	Value() []byte

	Release()
}

type iterator struct {
	key 		[]byte
	value 		[]byte
	typeName 	[]byte
	db 			*leveldb.DB
	it 			Iterator
	rege 		string
	greater		bool
}

type uniqueIterator struct{
	iterator
}

type indexIterator struct{
	 iterator
}
