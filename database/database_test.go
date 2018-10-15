package database

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func openDb() (*LDataBase, func()) {

	fileName := "./hello"
	reFn := func() {
		errs := os.RemoveAll(fileName)
		if errs != nil {
			log.Fatalln(errs)
		}
	}
	_, exits := os.Stat(fileName)
	if exits == nil {
		reFn()
	}

	db, err := NewDataBase(fileName)
	if err != nil {
		fmt.Println("new database failed")
		return nil, reFn
	}

	return db, func() {
		db.Close()
		reFn()
	}
}

func Objects() ([]TableIdObject, []House) {
	objs := []TableIdObject{}
	Houses := []House{}
	for i := 1; i <= 3; i++ {
		number := i * 10
		obj := TableIdObject{Code: AccountName(number + 1), Scope: ScopeName(number + 2), Table: TableName(number + 3), Payer: AccountName(number + 4), Count: uint32(number + 5)}
		objs = append(objs, obj)
		house := House{Area: uint64(number + 7), Carnivore: Carnivore{number + 8, number + 8}}
		Houses = append(Houses, house)
		obj = TableIdObject{Code: AccountName(number + 1), Scope: ScopeName(number + 2), Table: TableName(number + 3), Payer: AccountName(number + 4), Count: uint32(number + 5)}
		objs = append(objs, obj)
		house = House{Area: uint64(number + 8), Carnivore: Carnivore{number + 8, number + 8}}
		Houses = append(Houses, house)

		obj = TableIdObject{Code: AccountName(number + 1), Scope: ScopeName(number + 2), Table: TableName(number + 3), Payer: AccountName(number + 4), Count: uint32(number + 5)}
		objs = append(objs, obj)
		house = House{Area: uint64(number + 9), Carnivore: Carnivore{number + 8, number + 8}}
		Houses = append(Houses, house)
	}
	return objs, Houses
}

func saveObjs(objs []TableIdObject, houses []House, db *LDataBase) ([]TableIdObject, []House) {
	objs_ := []TableIdObject{}
	houses_ := []House{}
	for _, v := range objs {
		err := db.Insert(&v)
		if err != nil {
			log.Fatalln("insert table object failed")
		}
		objs_ = append(objs_, v)
	}

	for _, v := range houses {
		err := db.Insert(&v)
		if err != nil {
			log.Fatalln("insert house object failed")
		}
		houses_ = append(houses_, v)
	}
	return objs_, houses_
}

func Test_open(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()
}

func Test_insert(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()

	saveObjs(objs, houses, db)
}

func Test_find(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	objs_, houses_ := saveObjs(objs, houses, db)
	getGreaterObjs(objs_, houses_, db)

	getLessObjs(objs_, houses_, db)

	findObjs(objs_, houses_, db)

	getErrStruct(db)
}

func getErrStruct(db *LDataBase) {

	obj := TableIdObject{Table: 13}
	_, err := db.Get("byTable", &obj)
	if err != ErrStructNeeded {
		log.Fatalln(err)
	}
}

func getGreaterObjs(objs []TableIdObject, houses []House, db *LDataBase) {

	obj := TableIdObject{Table: 13}
	it, err := db.Get("byTable", obj)
	if err != nil {
		log.Fatalln(err)
	}

	/*                                                         */
	i := 2
	for it.Next() {
		obj = TableIdObject{}
		err = it.Data(&obj)
		if err != nil {
			log.Fatalln(err)
		}
		if obj != objs[i] {
			logObj(obj)
			logObj(objs[i])
			log.Fatalln("find next failed")
		}
		i--
	}
	it.Release()
}

func getLessObjs(objs []TableIdObject, houses []House, db *LDataBase) {
	i := 0
	obj := TableIdObject{Code: 11}
	it, err := db.Get("Code", obj)
	if err != nil {
		log.Fatalln(err)
	}

	for it.Next() {
		obj = TableIdObject{}
		err = it.Data(&obj)
		if err != nil {
			log.Fatalln(err)
		}

		if obj != objs[i] {
			logObj(obj)
			logObj(objs[i])
			log.Fatalln("find failed")
		}
		i++
	}
	i--

	tobjs := []TableIdObject{}
	err = db.GetObjects("Code",obj,&tobjs)
	if err != nil{
		log.Fatalln(err)
	}

	for j := 0;j < 2;j++{
		if tobjs[j] != objs[j]{
			log.Fatalln("get objects failed")
		}
	}


	for it.Prev() {
		obj = TableIdObject{}
		err = it.Data(&obj)
		if err != nil {
			log.Fatalln(err)
		}
		if obj != objs[i] {
			logObj(obj)
			logObj(objs[i])
			log.Fatalln("find failed")
		}
		i--
	}
	it.Release()

}

func findObjs(objs []TableIdObject, houses []House, db *LDataBase) {
	obj := TableIdObject{ID: 4}
	tmp := TableIdObject{}
	err := db.Find("id", obj, &tmp)
	if err != nil {
		log.Fatalln(err)
	}

	{
		hou := House{Area: 17}
		tmp := House{}
		err := db.Find("Area", hou, &tmp)
		if err != nil {
			log.Fatalln(err)
		}
	}

}

func Test_modify(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	saveObjs(objs, houses, db)
	modifyObjs(db)
}

func modifyObjs(db *LDataBase) {

	obj := TableIdObject{ID: 4, Code: 21, Scope: 22, Table: 23, Payer: 24, Count: 25}
	newobj := TableIdObject{ID: 4, Code: 200, Scope: 22, Table: 23, Payer: 24, Count: 25}

	err := db.Modify(&obj, func(object *TableIdObject) {
		object.Code = 200
	})
	if err != nil {
		log.Fatalln(err)
	}

	obj = TableIdObject{}
	tmp := TableIdObject{}
	obj.ID = 4
	err = db.Find("id", obj, &tmp)
	if err != nil {
		log.Fatalln(err)
	}
	if tmp != newobj {
		logObj(tmp)
		log.Fatalln("modify test error")
	}
}

func Test_remove(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	saveObjs(objs, houses, db)
	removeUnique(db)
}

func removeUnique(db *LDataBase) {

	obj := TableIdObject{Code: 21, Scope: 22, Table: 23, Payer: 24, Count: 25}
	err := db.Remove(obj)
	if err != ErrIncompleteStructure {
		log.Fatalln(err)
	}

	obj = TableIdObject{}
	obj.ID = 4

	tmp := TableIdObject{}
	err = db.Find("id", obj, &tmp)
	if err != nil {
		log.Fatalln(err)
	}

	// logObj(tmp)

	err = db.Remove(&obj)
	if err != ErrStructNeeded {
		log.Fatalln(err)
	}

	err = db.Remove(tmp)
	if err != nil {
		log.Fatalln(err)
	}

	tmp = TableIdObject{}
	err = db.Find("id", obj, &tmp)
	if err != ErrNotFound {
		log.Fatalln(err)
	}
}

func Test_undo(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	objs_, houses_ := saveObjs(objs, houses, db)
	undoObjs(objs_, houses_, db)
}

func undoObjs(objs []TableIdObject, houses []House, db *LDataBase) {
	//i := 0
	obj := House{Carnivore: Carnivore{38, 38}}
	it, err := db.Get("Carnivore", obj)
	if err != nil {
		log.Fatalln(err)
	}

	for it.Next() {
		obj = House{}
		err = it.Data(&obj)
		if err != nil {
			log.Fatalln(err)
		}

		// logObj(obj)

	}

	it.Release()

	h := House{Area: 38, Carnivore: Carnivore{100, 100}}
	err = db.Insert(&h)
	if err != ErrAlreadyExists {
		log.Fatalln(err)
	}

	// fmt.Println("-------------")
	obj = House{Id: 10}
	tmp := House{}
	err = db.Find("id", obj, &tmp)
	if err != ErrNotFound {
		log.Fatalln(err)
	}

	it.Release()
}
