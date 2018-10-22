package database

import (
	"fmt"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"os"
	"testing"
)

func Test_rawDb(t *testing.T) {
	//f, err := os.Create("./cpu.txt")
   //if err != nil {
   //    log.Fatal(err)
   //}
   //pprof.StartCPUProfile(f)
   //defer pprof.StopCPUProfile()

	fileName := "./eosspark"
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
	db,err := leveldb.OpenFile(fileName,nil)
	if err != nil{
		log.Fatalln(err)
	}
	defer func() {
		db.Close()
		reFn()
	}()


	objs, houses := Objects()
	if len(objs) != len(houses){
		log.Fatalln("ERROR")
	}
	for i := 1; i <= 10; i++ {
		db.Put([]byte(string(i)),[]byte(string(i)),nil)
	}
	it := db.NewIterator(nil,nil)
	for it.Next(){
		//fmt.Println(it.Key())
	}

	it = db.NewIterator(&util.Range{Start:[]byte(string(3)),Limit:[]byte(string(11))},nil)

	for it.Next(){
		//fmt.Println(it.Key())
	}
	i := 0
	for index,v := range houses{
		b,err := rlp.EncodeToBytes(v)
		if err != nil{
			log.Fatalln(err)
		}
		db.Put([]byte(string(index + 10)),b,nil)
		if err != nil{
			log.Fatalln(err)
		}
		i++
	}
}

func openDb() (DataBase, func()) {

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

func saveObjs(objs []TableIdObject, houses []House, db DataBase) ([]TableIdObject, []House) {
	objs_ := []TableIdObject{}
	houses_ := []House{}

	for _, v := range houses {
		err := db.Insert(&v)
		if err != nil {
			log.Fatalln(err)
		}
		houses_ = append(houses_, v)
	}

	for _, v := range objs {

		err := db.Insert(&v)
		if err != nil {
			log.Fatalln(err)
			log.Fatalln("insert table object failed")
		}

		objs_ = append(objs_, v)
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
	if len(objs) != len(houses){
		log.Fatalln("ERROR")
	}

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
	//
	findObjs(objs_, houses_, db)

	findInLineFieldObjs(objs_, houses_, db)

	findAllNonUniqueFieldObjs(objs_, houses_, db);

	getErrStruct(db)
}

func getErrStruct(db DataBase) {

	obj := TableIdObject{Scope: 12,Table:13}
	_, err := db.GetIndex("byTable", &obj)
	if err != ErrStructNeeded {
		log.Fatalln(err)
	}
}

func getGreaterObjs(objs []TableIdObject, houses []House, db DataBase) {

	obj := TableIdObject{Scope:23}
	idx, err := db.GetIndex("byTable", obj)
	if err != nil {
		log.Fatalln(err)
	}


	it,err := idx.LowerBound(obj)
	if err != nil{
		log.Fatalln(err)
	}
	defer it.Release()

	if idx.CompareBegin(it){
		tmp := TableIdObject{}
		idx.Begin(&tmp)
		if tmp != objs[8]{
			logObj(objs[8])
			logObj(tmp)
		}
	}
	i := 8
	for it.Next(){
		tmp := TableIdObject{}
		it.Data(&tmp)
		if tmp != objs[i] {
			logObj(tmp)
			logObj(objs[i])
		}
		i--
	}
	if !idx.CompareEnd(it){
		log.Fatalln("CompareEnd")
	}
	it.Release()

	it ,err = idx.UpperBound(obj)
	if err != nil{
		log.Fatalln(err)
	}
	i = 8
	for it.Next(){
		tmp := TableIdObject{}
		it.Data(&tmp)
		if tmp != objs[i] {
			logObj(tmp)
			logObj(objs[i])
		}
		i--
	}

	it.Release()
}

func getLessObjs(objs []TableIdObject, houses []House, db DataBase) {
	obj := TableIdObject{Code: 13}
	idx, err := db.GetIndex("Code", obj)
	if err != nil {
		log.Fatalln(err)
	}

	it ,err := idx.LowerBound(obj)
	if err != nil {
		log.Fatalln(err)
	}
	i := 3
	defer it.Release()
	for it.Next(){
		tmp := TableIdObject{}
		it.Data(&tmp)
		if tmp != objs[i]{
			logObj(objs[i])
			logObj(tmp)
		}
		i++
	}
}

func findObjs(objs []TableIdObject, houses []House, db DataBase) {
	obj := TableIdObject{ID: 4}
	tmp := TableIdObject{}
	err := db.Find("id", obj, &tmp)
	if err != nil {
		log.Fatalln(err)
	}

	{
		hou := House{Area: 18}
		tmp := House{}
		err := db.Find("Area", hou, &tmp)
		if err != nil {
			log.Fatalln(err)
		}
		if houses[1] != tmp{
			logObj(tmp)
			logObj(houses[1])
			log.Fatalln("Find Object")
		}
	}
}

func findInLineFieldObjs(objs []TableIdObject, houses []House, db DataBase) {
	hou := House{Carnivore:Carnivore{28,38}}
	//idx,err := db.GetIndex("Tiger", hou)
	idx,err := db.GetIndex("Lion", hou)
	if err != nil {
		log.Fatalln(err)
	}

	it ,err := idx.LowerBound(hou)
	if err != nil {
		log.Fatalln(err)
	}
	i := 8
	defer it.Release()
	for it.Next(){
		tmp := House{}
		it.Data(&tmp)
		if tmp != houses[i]{
			logObj(houses[i])
			logObj(tmp)
		}
		i--
	}
}

func findAllNonUniqueFieldObjs(objs []TableIdObject, houses []House, db DataBase) {

	obj := TableIdObject{Scope:12,Table:13}

	err := db.Find("byTable",obj,&obj)
	if err != nil{
		log.Fatalln(err)
	}
	//logObj(obj)
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

func modifyObjs(db DataBase) {

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

func removeUnique(db DataBase) {

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

