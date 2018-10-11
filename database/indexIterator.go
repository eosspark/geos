
package database
import (
	"github.com/syndtr/goleveldb/leveldb"
	"regexp"
)

func newIndexIterator(typeName []byte,it Iterator,db *leveldb.DB,rege string,greater bool)(*indexIterator){
	ite := &indexIterator{ iterator{typeName:typeName,it:it,db:db,rege:rege,greater:greater}}

	if ite.greater{
		if ! ite.it.Last()	{
			return nil
		}
	}
	return ite
}

func (iterator *indexIterator)Next()bool{
	if iterator.greater{
		return iterator.prev()
	}
	return iterator.next()
	return false
}

func (iterator *indexIterator)Prev()bool{
	if iterator.greater{
		return iterator.next()
	}
	return iterator.prev()
	return false
}

func (iterator *indexIterator)next()bool{
	for iterator.it.Next(){
		reg := regexp.MustCompile(iterator.rege)
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


func (iterator *indexIterator)prev()bool{
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

func(iterator *indexIterator)Release(){
	iterator.it.Release()
}

func (iterator *indexIterator)Key()[]byte{
	return iterator.key
}

func (iterator *indexIterator)Value()[]byte{
	return iterator.value
}

func (iterator *indexIterator)Last() bool{// TODO
	return iterator.it.Last()
}

func (iterator *indexIterator)First() bool{// TODO
	return iterator.it.First()
}
//
//func (iterator *indexIterator)Error() error{// TODO
//	return iterator.it.Error()
//}
//
//func (iterator *indexIterator)Seek(key []byte) bool{// TODO
//	return iterator.it.Seek(key)
//}
