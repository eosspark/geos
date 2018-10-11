package database

import (
	"fmt"
	"github.com/eosspark/eos-go/crypto/rlp"
	"log"
	"os"
	"testing"
)


func Test_open(t *testing.T) {
	db,clo := openDb()

	if db == nil{
		log.Fatalln("db open failed")
	}
	defer clo()
}

func Test_insert(t *testing.T) {
	db,clo := openDb()

	if db == nil{
		log.Fatalln("db open failed")
	}
	defer clo()

	objs,houses := Objects()

	saveObjs(objs,houses,db)
}

func Test_find(t *testing.T) {
	db,clo := openDb()

	if db == nil{
		log.Fatalln("db open failed")
	}
	defer clo()

	objs,houses := Objects()
	objs_,houses_ := saveObjs(objs,houses,db)

	findGreaterObjs(objs_,houses_,db)

	findLessObjs(objs_,houses_,db)

	findIdObjs(objs_,houses_,db)
}

func openDb()(*LDatabase,func()){

	fileName := "./hello"
	reFn :=  func() {
		errs := os.RemoveAll(fileName)
		if errs != nil{
			log.Fatalln(errs)
		}
	}
	_,exits := os.Stat(fileName)
	if exits == nil{
		reFn()
	}

	db,err := NewLDatabase(fileName)
	if err != nil{
		fmt.Println("new database failed")
		return nil,reFn
	}

	return db,func(){
		db.Close()
		reFn()
	}
}

func Objects()([]TableIdObject,[]House){
	objs 	:= []TableIdObject{}
	Houses 	:= []House{}
	for i := 1; i <= 3;i++{
		number := i * 10
		obj := TableIdObject{Code:AccountName(number + 1),Scope:ScopeName(number + 2),Table:TableName(number + 3),Payer:AccountName(number + 4),Count:uint32(number + 5)}
		objs = append(objs, obj)
		house := House{Area:uint64(number + 9),Carnivore:Carnivore{number + 8,number + 8}}
		Houses = append(Houses,house)

		obj = TableIdObject{Code:AccountName(number + 1),Scope:ScopeName(number + 2),Table:TableName(number + 3),Payer:AccountName(number + 4),Count:uint32(number + 5)}
		objs = append(objs, obj)
		house = House{Area:uint64(number + 9),Carnivore:Carnivore{number + 8,number + 8}}
		Houses = append(Houses,house)

		obj = TableIdObject{Code:AccountName(number + 1),Scope:ScopeName(number + 2),Table:TableName(number + 3),Payer:AccountName(number + 4),Count:uint32(number + 5)}
		objs = append(objs, obj)
		house = House{Area:uint64(number + 9),Carnivore:Carnivore{number + 8,number + 8}}
		Houses = append(Houses,house)
	}
	return objs,Houses
}

func saveObjs(objs []TableIdObject,houses []House,db *LDatabase) ([]TableIdObject,[]House) {
	objs_ := []TableIdObject{}
	houses_ :=  []House{}
	for _,v:= range objs{
		err := db.Insert(&v)
		if err != nil{
			log.Fatalln("insert table object failed")
		}
		objs_ = append(objs_,v)
	}

	for _,v:= range houses{
		err := db.Insert(&v)
		if err != nil{
			log.Fatalln("insert house object failed")
		}
		houses_ = append(houses_,v)
	}
	return objs_,houses_
}

func findGreaterObjs(objs []TableIdObject,houses []House,db *LDatabase) {
	obj:= TableIdObject{Table:13}
	it,err := db.Find("byTable",obj)
	if err != nil{
		log.Fatalln(err)
	}
	if it == nil{
		log.Fatalln("iterator failed")
	}
	/*                                                         */
	i := 2
	for it.Next(){
		obj = TableIdObject{}
		err = rlp.DecodeBytes(it.Value(),&obj)
		if err != nil {
			log.Fatalln(err)
		}
		if obj != objs[i]{
			logObj(obj)
			logObj(objs[i])
			log.Fatalln("find next failed")
		}
		i--
	}
	it.Release()
}


func findLessObjs(objs []TableIdObject,houses []House,db *LDatabase) {
	i := 0
	obj := TableIdObject{Code:11}
	it,err := db.Find("Code",obj)
	if err != nil{
		log.Fatalln(err)
	}

	for it.Next(){
		obj = TableIdObject{}
		err = rlp.DecodeBytes(it.Value(),&obj)
		if err != nil {
			log.Fatalln(err)
		}

		if obj != objs[i]{
			logObj(obj)
			logObj(objs[i])
			log.Fatalln("find failed")
		}
		i++
	}
	i--
	for it.Prev(){
		obj = TableIdObject{}
		err = rlp.DecodeBytes(it.Value(),&obj)
		if err != nil {
			log.Fatalln(err)
		}
		if obj != objs[i]{
			logObj(obj)
			logObj(objs[i])
			log.Fatalln("find failed")
		}
		i--
	}
	it.Release()

}

func findIdObjs(objs []TableIdObject,houses []House,db *LDatabase){

	i := 2
	obj := TableIdObject{ID:4}
	it,err := db.Find("id",obj)
	if err != nil{
		fmt.Println(err)
	}

	for it.Prev(){
		obj = TableIdObject{}
		err = rlp.DecodeBytes(it.Value(),&obj)
		if err != nil {
			log.Fatalln(err)
		}

		if obj != objs[i]{
			logObj(obj)
			logObj(objs[i])
			log.Fatalln("find failed")
		}
		i--
	}
	i++

	for it.Next(){
		obj = TableIdObject{}
		err = rlp.DecodeBytes(it.Value(),&obj)
		if err != nil {
			log.Fatalln(err)
		}
		if obj != objs[i]{
			logObj(obj)
			logObj(objs[i])
			log.Fatalln("find failed")
		}
		i++
	}
	it.Release()
}
