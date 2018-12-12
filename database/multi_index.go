package database

import (
	"bytes"
	"errors"
)

type MultiIndex struct {
	begin     []byte
	end       []byte
	typeName  []byte
	fieldName []byte
	db        DataBase
	it        DbIterator
}

func newMultiIndex(typeName, fieldName, begin, end []byte, db DataBase) *MultiIndex {
	return &MultiIndex{typeName: typeName, fieldName: fieldName, begin: begin, end: end, db: db}
}

/*

@param in 			--> 	object

@return
success 			-->		nil 	(Iterator valid)
error 				-->		error 	(Iterator invalid)

*/

func (index *MultiIndex) LowerBound(in interface{},skip... SkipSuffix) (Iterator, error) {
	it, err := index.db.lowerBound(index.begin, index.end, index.fieldName, in,skip...)
	if err != nil {
		return index.db.EndIterator(index.begin, index.end, index.typeName)
	}
	return it, nil
}

/*

@param in 			--> 	object

@return

success 			-->		nil 	(Iterator valid)
error 				-->		error 	(Iterator invalid)

*/

func (index *MultiIndex) UpperBound(in interface{},skip... SkipSuffix) (Iterator, error) {
	it, err := index.db.upperBound(index.begin, index.end, index.fieldName, in,skip...)
	if err != nil {
		return index.db.EndIterator(index.begin, index.end, index.typeName)
	}
	return it, nil
}

/*

@param in 			--> 	object
@param out 			--> 	output object(pointer)

@return
success 			-->		nil
error 				-->		error

*/

func (index *MultiIndex) Find(in interface{}, out interface{},skip... SkipSuffix) error {
	return index.db.Find(string(index.fieldName), in, out,skip...)
}

/*

--> it == idx.begin() <--

@param in 			--> 	Iterator

@return
success 			-->		true
error 				-->		false

*/

func (index *MultiIndex) CompareBegin(in Iterator) bool {
	if in == nil {
		return false
	}
	it := index.Begin()
	return bytes.Compare(in.Value(), it.Value()) == 0
	//return it.Begin() == in.Begin()
}

/*

--> it1 == it2 <--

@param it1 			--> 	Iterator
@param it2 			--> 	Iterator

@return
success 			-->		true
error 				-->		false

*/

func (index *MultiIndex) CompareIterator(it1 Iterator, it2 Iterator) bool {
	if it1 == nil && it2 == nil {
		return true
	}
	if it1 == nil || it2 == nil {
		return false
	}
	if it1.Begin() && it2.Begin() {
		return true
	}
	if it1.End() && it2.End() {
		return true
	}

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
	if in == nil {
		return false
	}

	end := index.End()

	return end.End() && in.End() // FIXME only end
}

/*

@return
success 			-->		Iterator
error 				-->		nil

*/

func (index *MultiIndex) End() Iterator {
	// TODO
	it, err := index.db.EndIterator(index.begin, index.end, index.typeName)
	if err != nil {
		return nil
	}
	if it == nil {
		return nil
	}
	return it
}

/*

@return
success 			-->		Iterator
error 				-->		nil

*/

func (index *MultiIndex) Begin() Iterator {
	it, err := index.db.BeginIterator(index.begin, index.end, index.typeName)
	if err != nil || it == nil {
		it, err = index.db.EndIterator(index.begin, index.end, index.typeName)
		if err != nil {
			//panic()
			return nil
		}
	}
	return it
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
	if it == nil {
		return errors.New("MultiIndex BeginData : iterator is nil")
	}
	err := it.Data(out)
	if err != nil {
		return errors.New("MultiIndex BeginData : " + err.Error())
	}
	return nil
}

func (index *MultiIndex) IteratorTo(in interface{}) Iterator {
	it, err := index.db.IteratorTo(index.begin, index.end, index.fieldName, in)
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
