package database

import (
	"bytes"
	"github.com/eosspark/eos-go/crypto/rlp"
)

type multiIndex struct {
	begin 		[]byte
	end 		[]byte
	itBegin     []byte
	itEnd       []byte
	typeName 	[]byte
	fieldName 	[]byte
	db 			DataBase
	it 			DbIterator
	rege 		string
	greater		bool
}

func newMultiIndex(typeName,fieldName,begin,end []byte, greater bool,db DataBase) *multiIndex {
	return &multiIndex{typeName: typeName,fieldName:fieldName,begin:begin,end:end, greater: greater,db:db}
}

/*

@param in 			--> 	object

@return
success 			-->		nil 	(Iterator valid)
error 				-->		error 	(Iterator invalid)

*/
func (index *multiIndex)LowerBound(in interface{}) (Iterator,error) {
	it,err := index.db.lowerBound(index.begin,index.end,index.fieldName,in,index.greater)
	if err != nil{
		return nil,err
	}
	index.itBegin = make([]byte,len(it.Begin()))
	copy(index.itBegin,it.Begin())
	index.itEnd = nil
	return it,nil
}
/*

@param in 			--> 	object

@return

success 			-->		nil 	(Iterator valid)
error 				-->		error 	(Iterator invalid)

*/

func (index *multiIndex) UpperBound(in interface{}) (Iterator,error) {
	it,err := index.db.upperBound(index.begin,index.end,index.fieldName,in,index.greater)
	if err != nil{
		return nil,err
	}
	//fmt.Println(it.Begin())
	index.itBegin = make([]byte,len(it.Begin()))
	copy(index.itBegin,it.Begin())
	index.itEnd = nil
	return it,nil
}
/*

@param in 			--> 	object
@param out 			--> 	output object(pointer)

@return
success 			-->		nil
error 				-->		error

*/

func (index *multiIndex) Find( in interface{}, out interface{}) error {
	return index.db.Find(string(index.fieldName),in,out)
}

/*

@param out 			--> 	output object(pointer)

@return
success 			-->		nil
error 				-->		error

*/

func (index *multiIndex)Begin(out interface{}) error {
	// TODO
	err := rlp.DecodeBytes(index.itBegin,out)
	if err != nil{
		return err
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

func (index *multiIndex)CompareBegin(in Iterator)bool{
	//fmt.Println(in.Begin())
	//fmt.Println(index.itBegin)
	return bytes.Compare(in.Begin(),index.itBegin) == 0
}

/*

--> it == idx.end() <--

@param in 			--> 	Iterator

@return
success 			-->		true
error 				-->		false

*/
func (index *multiIndex)CompareEnd(in Iterator)bool{
	return in.Value() == nil
}
/*

@return 			-->		Iterator

*/
func (index *multiIndex)End() Iterator {
	// TODO
	return nil
}

/*

@param in 			--> 	object

@return
success 			-->		Iterator
error 				-->		nil

*/

func (index *multiIndex) BeginIterator() Iterator {
	// TODO
	if len(index.typeName) == 0{
		return nil
	}

	key := append(index.typeName, '_')
	key = append(key, '_')
	// typeName__tag__
	key = append(key, index.fieldName...)

	if index.it.Seek(key){
		return nil
	}
	return &index.it
}

func (index *multiIndex) IteratorTo  (in interface{}) Iterator {
	// TODO
	if len(index.typeName) == 0{
		return nil
	}
	fields, err := getFieldInfo(string(index.fieldName), in)
	if err != nil {
		return nil
	}
	key := append(index.typeName, '_')
	key = append(key, '_')
	// typeName__tag__
	key = append(key, index.fieldName...)
	key =getFieldValue(key, fields)

	if index.it.Seek(key){
		return nil
	}
	return &index.it
}

func (index *multiIndex)Empty() bool {
	//fmt.Println(index.begin," : ",index.end)
	//fmt.Println(string(index.begin),": : :",string(index.end))
	return index.db.Empty(index.begin,index.end,index.fieldName)
}

