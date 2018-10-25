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

	findObjs(objs_, houses_, db)

	findInLineFieldObjs(objs_, houses_, db)

	findAllNonUniqueFieldObjs(objs_, houses_, db);

	getErrStruct(db)

	getLessObjs(objs_, houses_, db)
}

//func Test_modifyUndo(t *testing.T) {
//
//	db, clo := openDb()
//	if db == nil {
//		log.Fatalln("db open failed")
//	}
//	defer clo()
//
//
//	objs, houses := Objects()
//	objs_,_:=saveObjs(objs, houses, db)
//
//	//session := db.StartSession()
//	//session.Commit()
//
//	idx, err := db.GetIndex("Code", TableIdObject{})
//	if err != nil{
//		log.Println(err)
//	}
//	it,err := idx.LowerBound(TableIdObject{Code:11})
//	if err != nil{
//		log.Fatalln(err)
//	}
//	i := 0
//	for it.Next(){
//		tmp := TableIdObject{}
//		it.Data(&tmp)
//		//logObj(tmp)
//		if objs_[i] != tmp{
//			logObj(tmp)
//			log.Fatalln("error lower bound")
//		}
//		i++
//	}
//	it.Release()
//
//	session := db.StartSession()
//	defer session.Undo()
//	obj := TableIdObject{ID: 4, Code: 21, Scope: 22, Table: 26, Payer: 27, Count: 25}
//	newobj := TableIdObject{ID: 4, Code: 200, Scope: 22, Table: 26, Payer: 27, Count: 25}
//
//	err = db.Modify(&obj, func(object *TableIdObject) {
//		object.Code = 200
//	})
//	if err != nil {
//		log.Fatalln(err)
//	}
//	session.Undo()
//	obj = TableIdObject{}
//	tmp := TableIdObject{}
//	obj.ID = 4
//	err = db.Find("id", obj, &tmp)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	if tmp == newobj {
//		logObj(newobj)
//		log.Fatalln("modify test error")
//	}
//}
//
//func Test_undoInsert(t *testing.T) {
//	db, clo := openDb()
//	if db == nil {
//		log.Fatalln("db open failed")
//	}
//	defer clo()
//
//	//////////////////////////////////////////////		Insert UNDO		///////////////////////////////////
//	db.SetRevision(10)
//	session := db.StartSession()
//	objs, _ := Objects()
//	for i:= 0;i < 3;i++{
//		err := db.Insert(&objs[i])
//		if err != nil{
//			log.Println(err)
//		}
//	}
//
//	session.Undo()
//	idx, err := db.GetIndex("Code", TableIdObject{})
//	if err != nil{
//		log.Println(err)
//	}
//
//	_,err = idx.LowerBound(TableIdObject{Code:11})
//	if err != ErrNotFound{
//		log.Fatalln(err)
//	}
//
//
//	//////////////////////////////////////////////		COMMIT		///////////////////////////////////
//
//	session = db.StartSession()
//	for i:= 0;i < 3;i++{
//		err := db.Insert(&objs[i])
//		if err != nil{
//			log.Println(err)
//		}
//	}
//
//	db.Commit(11)
//	session.Undo()
//	idx, err = db.GetIndex("Code", TableIdObject{})
//	if err != nil{
//		log.Println(err)
//	}
//	it,err := idx.LowerBound(TableIdObject{Code:11})
//	if err != nil{
//		log.Fatalln(err)
//	}
//	i := 0
//	for it.Next(){
//		tmp := TableIdObject{}
//		it.Data(&tmp)
//		//logObj(tmp)
//		if objs[i] != tmp{
//			logObj(tmp)
//			log.Fatalln("error lower bound")
//		}
//		i++
//	}
//	it.Release()
//
//}
//
//func Test_undoRemove(t *testing.T) {
//	db, clo := openDb()
//	if db == nil {
//		log.Fatalln("db open failed")
//	}
//	defer clo()
//
//	//////////////////////////////////////////////	ready
//	objs, _ := Objects()
//	for i:= 0;i < 3;i++{
//		err := db.Insert(&objs[i])
//		if err != nil{
//			log.Println(err)
//		}
//	}
//	idx, err := db.GetIndex("Code", TableIdObject{})
//	if err != nil{
//		log.Println(err)
//	}
//	it,err := idx.LowerBound(TableIdObject{Code:11})
//	if err != nil{
//		log.Fatalln(err)
//	}
//
//	table := TableIdObject{}
//	i := 0
//	for it.Next(){
//		it.Data(&table)
//		if objs[i] != table{
//			logObj(objs[i])
//			logObj(table)
//			log.Fatalln("undo failed")
//		}
//		i++
//	}
//	session := db.StartSession()
//
//	err = db.Remove(table)
//	if err != nil{
//		log.Fatalln(err)
//	}
//	////////////////////////////////////////// begin
//	beginUndo, err := db.GetIndex("Code", TableIdObject{})
//	if err != nil{
//		log.Println(err)
//	}
//	beginIt,err := beginUndo.LowerBound(TableIdObject{Code:11})
//	if err != nil{
//		log.Fatalln(err)
//	}
//	i = 0
//	for beginIt.Next(){
//		table := TableIdObject{}
//		beginIt.Data(&table)
//		logObj(table)
//		if objs[i] != table{
//			logObj(objs[i])
//			logObj(table)
//			log.Fatalln("undo failed")
//		}
//		i++
//	}
//	if i != 2{
//		log.Println(i)
//		log.Fatalln("undo failed")
//	}
//	session.Undo()// undo
//	/////////////////////////////////////////// end
//	endUndo, err := db.GetIndex("Code", TableIdObject{})
//	if err != nil{
//		log.Println(err)
//	}
//	endIt,err := endUndo.LowerBound(TableIdObject{Code:11})
//	if err != nil{
//		log.Fatalln(err)
//	}
//	i = 0
//	for endIt.Next(){
//		table := TableIdObject{}
//		endIt.Data(&table)
//		if objs[i] != table{
//			logObj(objs[i])
//			logObj(table)
//			log.Fatalln("undo failed")
//		}
//		i++
//	}
//}

func Test_empty(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	saveObjs(objs, houses, db)

	idx, err := db.GetIndex("Code", TableIdObject{})
	if err != nil {
		log.Fatalln(err)
	}

	obj := TableIdObject{Code: 11}
	it ,err := idx.LowerBound(obj)
	if err != nil {
		log.Fatalln(err)
	}
	for it.Next(){
		tmp := TableIdObject{}
		it.Data(&tmp)
		err = db.Remove(&tmp)
		if err != nil{
			log.Fatalln(err)
		}

	}
	it.Release()

	if !idx.Empty(){
		log.Fatalln("empty error")
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

func Test_resourceLimitsObject(t *testing.T) {

	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}

	defer clo()

	limits := MakeResourceLimitsObjects()
	for _,v := range limits{
		err := db.Insert(&v)
		if err != nil{
			log.Fatalln(err)
		}
	}

	idx, err := db.GetIndex("byOwner", ResourceLimitsObject{})
	if err != nil {
		log.Fatalln(err)
	}


	for !idx.Empty() {
		tmp := ResourceLimitsObject{}
		obj := ResourceLimitsObject{Pending: false}
		it ,err := idx.LowerBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		idx.Begin(&tmp)
		//logObj(tmp)
		if idx.CompareEnd(it) || tmp.Pending == true{
			fmt.Println("db is empty")
		}

		err = db.Remove(&tmp)
		if err != nil{
			log.Fatalln(err)
		}
		it.Release()
	}

	if !idx.Empty(){
		log.Fatalln("empty error")
	}
	if idx.Empty(){
		//fmt.Println("empty successful !")
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
		obj := TableIdObject{Code: AccountName(number + 1), Scope: ScopeName(number + 2), Table: TableName(number + 3+ i + 1), Payer: AccountName(number + 4 + i + 1), Count: uint32(number + 5)}
		objs = append(objs, obj)
		house := House{Area: uint64(number + 7), Carnivore: Carnivore{number + 8, number + 8}}
		Houses = append(Houses, house)
		obj = TableIdObject{Code: AccountName(number + 1), Scope: ScopeName(number + 2), Table: TableName(number + 3 + i + 2), Payer: AccountName(number + 4 + i + 2), Count: uint32(number + 5)}
		objs = append(objs, obj)
		house = House{Area: uint64(number + 8), Carnivore: Carnivore{number + 8, number + 8}}
		Houses = append(Houses, house)

		obj = TableIdObject{Code: AccountName(number + 1), Scope: ScopeName(number + 2), Table: TableName(number + 3+ i + 3), Payer: AccountName(number + 4 +i + 3), Count: uint32(number + 5)}
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

		//logObj(v)
		err := db.Insert(&v)
		if err != nil {
			log.Fatalln(err)
			log.Fatalln("insert table object failed")
		}

		objs_ = append(objs_, v)
	}
	return objs_, houses_
}

func getErrStruct(db DataBase) {

	obj := TableIdObject{Scope: 12,Table:13}
	_, err := db.GetIndex("byTable", &obj)
	if err != ErrStructNeeded {
		log.Fatalln(err)
	}
}

func getGreaterObjs(objs []TableIdObject, houses []House, db DataBase) {

	obj := TableIdObject{Scope:22}
	idx, err := db.GetIndex("byTable", obj)
	if err != nil {
		log.Fatalln(err)
	}


	it,err := idx.LowerBound(obj)
	if err != nil{
		log.Fatalln(err)
	}
	defer it.Release()

	//for _,v := range objs{
	//	logObj(v)
	//}
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
	idx, err := db.GetIndex("Code", TableIdObject{})
	if err != nil {
		log.Fatalln(err)
	}

	it ,err := idx.LowerBound(obj)
	if err != nil {
		log.Fatalln(err)
	}
	i := 3
	for it.Next(){
		tmp := TableIdObject{}
		it.Data(&tmp)
		if tmp != objs[i]{
			logObj(objs[i])
			logObj(tmp)
		}
		i++
	}
	it.Release()
}

func modifyObjs(db DataBase) {

	obj := TableIdObject{ID: 4, Code: 21, Scope: 22, Table: 26, Payer: 27, Count: 25}
	newobj := TableIdObject{ID: 4, Code: 200, Scope: 22, Table: 26, Payer: 27, Count: 25}

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

func findObjs(objs []TableIdObject, houses []House, db DataBase) {
	obj := TableIdObject{ID: 4}
	tmp := TableIdObject{}
	err := db.Find("id", obj, &tmp)
	if err != nil {
		log.Fatalln(err)
	}
	//logObj(tmp)
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

	obj := TableIdObject{Scope:12,Table:15}

	err := db.Find("byTable",obj,&obj)
	if err != nil{
		log.Fatalln(err)
	}
	//logObj(obj)
}

func removeUnique(db DataBase) {
	obj := TableIdObject{Code: 21, Scope: 22, Table: 23, Payer: 24, Count: 25}
	err := db.Remove(&obj)
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

	err = db.Remove(&tmp)
	if err != nil {
		log.Fatalln(err)
	}

	tmp = TableIdObject{}
	err = db.Find("id", obj, &tmp)
	if err != ErrNotFound {
		log.Fatalln(err)
	}
}

func MakeResourceLimitsObjects()([]ResourceLimitsObject){
	//limits := make([]ResourceLimitsObject,0)
	limits := []ResourceLimitsObject{}

	for i := 1; i <= 13; i++ {
		number := 100
		obj := ResourceLimitsObject{Owner:AccountName(number + i)}
		limits = append(limits,obj)
	}
	return  limits
}

func Test_create(t *testing.T) {

	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}

	defer clo()

	//i := 1
	//number := 10
	//obj := TableIdObject{Code: AccountName(number + 1), Scope: ScopeName(number + 2), Table: TableName(number + 3+ i + 1), Payer: AccountName(number + 4 + i + 1), Count: uint32(number + 5)}
	//house := House{Area: uint64(number + 7), Carnivore: Carnivore{number + 8, number + 8}}

	//db.create(&obj)
	//db.create(&house)
}

