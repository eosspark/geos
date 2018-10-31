package database

import (
	"bytes"
	"errors"
	"github.com/eosspark/eos-go/crypto/rlp"
	"log"
)

type MultiIndex struct {
	begin     []byte
	end       []byte
	itBegin   []byte
	itEnd     []byte
	typeName  []byte
	fieldName []byte
	db        DataBase
	it        DbIterator
	greater   bool
}

func newMultiIndex(typeName, fieldName, begin, end []byte, greater bool, db DataBase) *MultiIndex {
	return &MultiIndex{typeName: typeName, fieldName: fieldName, begin: begin, end: end, greater: greater, db: db}
}

/*

@param in 			--> 	object

@return
success 			-->		nil 	(Iterator valid)
error 				-->		error 	(Iterator invalid)

*/
func (index *MultiIndex) LowerBound(in interface{}) (Iterator, error) {
	it, err := index.db.lowerBound(index.begin, index.end, index.fieldName, in, index.greater)
	if err != nil {
		return nil, err
	}
	index.itBegin = make([]byte, len(it.Begin()))
	copy(index.itBegin, it.Begin())
	index.itEnd = nil
	return it, nil
}

/*

@param in 			--> 	object

@return

success 			-->		nil 	(Iterator valid)
error 				-->		error 	(Iterator invalid)

*/

func (index *MultiIndex) UpperBound(in interface{}) (Iterator, error) {
	it, err := index.db.upperBound(index.begin, index.end, index.fieldName, in, index.greater)
	if err != nil {
		return nil, err
	}
	//fmt.Println(it.Begin())
	index.itBegin = make([]byte, len(it.Begin()))
	copy(index.itBegin, it.Begin())
	index.itEnd = nil
	return it, nil
}

/*

@param in 			--> 	object
@param out 			--> 	output object(pointer)

@return
success 			-->		nil
error 				-->		error

*/

func (index *MultiIndex) Find(in interface{}, out interface{}) error {
	return index.db.Find(string(index.fieldName), in, out)
}

/*

@param out 			--> 	output object(pointer)

@return
success 			-->		nil
error 				-->		error

*/

func (index *MultiIndex) BeginData(out interface{}) error {
	// TODO
	it := index.Begin()
	if it == nil{
		return errors.New("MultiIndex BeginData : iterator is nil")
	}
	err := rlp.DecodeBytes(it.Value(), out)
	if err != nil {
		return errors.New("MultiIndex BeginData : " + err.Error())
	}
	return nil
}

/*

--> it == idx.begin() <--

@param in 			--> 	Iterator

@return
success 			-->		true
error 				-->		false

*/

func (index *MultiIndex) CompareBegin(in Iterator) bool {
	it := index.Begin()
	//fmt.Println(it.Value())
	//fmt.Println(in.Value())
	return bytes.Compare(it.Value(), in.Value()) == 0
}

/*

--> it1 == it2 <--

@param in 			--> 	Iterator

@return
success 			-->		true
error 				-->		false

*/
func (index *MultiIndex) CompareIterator(it1 Iterator,it2 Iterator) bool {
	if it1  == nil || it2 == nil{
		return false
	}
	//fmt.Println(it1.Value())
	//fmt.Println(it2.Value())
	return bytes.Compare(it1.Value(), it2.Value()) == 0
}
/*

--> it == idx.end() <--

@param in 			--> 	Iterator

@return
success 			-->		true
error 				-->		false

*/
func (index *MultiIndex) CompareEnd(in Iterator) bool {
	return in.Value() == nil
}

/*

@return 			-->		Iterator

*/

func (index *MultiIndex) End() Iterator {
	// TODO
	return nil
}

/*

@param in 			--> 	object

@return
success 			-->		Iterator
error 				-->		nil

*/

func (index *MultiIndex) Begin() Iterator {
	it, err := index.db.BeginIterator(index.begin, index.end, index.fieldName, index.typeName, index.greater)
	if err != nil {
		log.Println("MultiIndex Begin Error : ",err)
		return nil
	}
	if it == nil {
		log.Println("MultiIndex Begin Iterator Is Empty ")
		return nil
	}
	return it
}

func (index *MultiIndex) IteratorTo(in interface{}) Iterator {
	it, err := index.db.IteratorTo(index.begin, index.end, index.fieldName, in, index.greater)
	if err != nil {
		panic(err)
		//log ?
		return nil
	}
	return it
}

func (index *MultiIndex) Empty() bool {
	return index.db.Empty(index.begin, index.end, index.fieldName)
}
