package database
import (
	"github.com/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"regexp"
)

func newUniqueIterator(typeName []byte,it Iterator,db *leveldb.DB,rege string)*uniqueIterator{
	return &uniqueIterator{ iterator{typeName:typeName,it:it,db:db,rege:rege}}
}
func (iterator *uniqueIterator) Next() bool {

	for iterator.it.Next(){
		reg := regexp.MustCompile(iterator.rege)
		//fmt.Println(iterator.it.Key())
		find := reg.Match(iterator.it.Key())
		if !find{
			continue
		}
		value := iterator.it.Value()
		k := idKey(value,iterator.typeName)
		v ,err := iterator.db.Get(k,nil)
		if err != nil{
			return false
		}

		iterator.key = k
		iterator.value = v
		return true
	}

	return false
}

func (iterator *uniqueIterator) Prev() bool {

	for iterator.it.Prev() {
		reg := regexp.MustCompile(iterator.rege)
		//fmt.Println(iterator.it.Key())
		find := reg.Match(iterator.it.Key())
		if !find{
			continue
		}

		value := iterator.it.Value()
		k := idKey(value,iterator.typeName)
		v ,err := iterator.db.Get(k,nil)
		if err != nil{
			return false
		}

		iterator.key = k
		iterator.value = v
		return true
	}

	return false
}

func (iterator *uniqueIterator) Release() {
	iterator.it.Release()
}

func (iterator *uniqueIterator) Key() []byte {
	return iterator.key
}

func(iterator *uniqueIterator)Data(data interface{})error{
	return rlp.DecodeBytes(iterator.Value(),data)
}

func (iterator *uniqueIterator) Value() []byte {
	return iterator.value
}

func (iterator *uniqueIterator) Last() bool {// TODO
	return iterator.it.Last()
}

func (iterator *uniqueIterator) First() bool {// TODO
	return iterator.it.First()
}

//
//func (iterator *uniqueIterator)Error() error{// TODO
//	return iterator.it.Error()
//}
//
//func (iterator *uniqueIterator)Seek(key []byte) bool{// TODO
//	return iterator.it.Seek(key)
//}
