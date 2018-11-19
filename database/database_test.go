package database

import (
	"bytes"
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
	db, err := leveldb.OpenFile(fileName, nil)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		db.Close()
		reFn()
	}()

	objs, houses := Objects()
	if len(objs) != len(houses) {
		log.Fatalln("ERROR")
	}
	keys := [][]byte{}
	for i := 1; i <= 10; i++ {
		key := []byte(string(i))
		key = append(key, key...)
		keys = append(keys, key)
	}
	for _, v := range keys {
		db.Put(v, v, nil)
	}
	if bytes.HasPrefix(keys[0], []byte(string(1))) {
		//fmt.Println(keys[0],[]byte(string(1)))
	}

	it := db.NewIterator(&util.Range{Start: []byte(string(5)), Limit: nil}, nil)
	if it.Seek([]byte(string(2))) {
		//fmt.Println("---------")
	}
	//fmt.Println(it.Value(),it.Key())
	for it.Next() {

		//fmt.Println(it.Value(),it.Key())
	}
	it.Release()
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
	if len(objs) != len(houses) {
		log.Fatalln("ERROR")
	}

	saveObjs(objs, houses, db)
}

func insert_te() {

	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	if len(objs) != len(houses) {
		log.Fatalln("ERROR")
	}

	saveObjs(objs, houses, db)
}

func Benchmark_insert(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		insert_te()
	}
}

func Test_find(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	objs_, houses_ := saveObjs(objs, houses, db)

	findObjs(objs_, houses_, db)

	lowerAndUpper(objs_, houses_, db)
	//
	findInLineFieldObjs(objs_, houses_, db)

	getLessObjs(objs_, houses_, db)
}

func Test_modifyUndo(t *testing.T) {

	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	objs_, _ := saveObjs(objs, houses, db)

	idx, err := db.GetIndex("Code", DbTableIdObject{})
	if err != nil {
		log.Println(err)
	}
	it, err := idx.LowerBound(DbTableIdObject{Code: 11})
	if err != nil {
		log.Fatalln(err)
	}
	i := 0
	for it.Next() {
		tmp := DbTableIdObject{}
		it.Data(&tmp)
		if objs_[i] != tmp {
			logObj(tmp)
			logObj(objs_[i])
			log.Fatalln("error lower bound")
		}
		i++
	}
	it.Release()

	session := db.StartSession()
	defer session.Undo()
	obj := DbTableIdObject{ID: 4, Code: 21, Scope: 22, Table: 26, Payer: 27, Count: 25}
	newobj := DbTableIdObject{ID: 4, Code: 200, Scope: 22, Table: 26, Payer: 27, Count: 25}

	err = db.Modify(&obj, func(object *DbTableIdObject) {
		object.Code = 200
	})
	if err != nil {
		log.Fatalln(err)
	}
	session.Undo()
	obj = DbTableIdObject{}
	tmp := DbTableIdObject{}
	obj.ID = 4
	err = db.Find("id", obj, &tmp)
	if err != nil {
		log.Fatalln(err)
	}
	if tmp == newobj {
		logObj(newobj)
		log.Fatalln("modify test error")
	}
}

func Test_undoInsert(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	//////////////////////////////////////////////		Insert UNDO		///////////////////////////////////
	db.SetRevision(10)
	session := db.StartSession()
	objs, _ := Objects()
	for i := 0; i < 3; i++ {
		err := db.Insert(&objs[i])
		if err != nil {
			log.Println(err)
		}
	}

	session.Undo()
	idx, err := db.GetIndex("Code", DbTableIdObject{})
	if err != nil {
		log.Println(err)
	}

	_, err = idx.LowerBound(DbTableIdObject{Code: 11})
	if err != ErrNotFound {
		log.Fatalln(err)
	}

	//////////////////////////////////////////////		COMMIT		///////////////////////////////////

	session = db.StartSession()
	for i := 0; i < 3; i++ {
		err := db.Insert(&objs[i])
		if err != nil {
			log.Println(err)
		}
	}

	db.Commit(11)
	session.Undo()
	idx, err = db.GetIndex("Code", DbTableIdObject{})
	if err != nil {
		log.Println(err)
	}
	it, err := idx.LowerBound(DbTableIdObject{Code: 11})
	if err != nil {
		log.Fatalln(err)
	}
	i := 0
	for it.Next() {
		tmp := DbTableIdObject{}
		it.Data(&tmp)
		//logObj(tmp)
		if objs[i] != tmp {
			logObj(tmp)
			log.Fatalln("error lower bound")
		}
		i++
	}
	it.Release()

}

func Test_undoRemove(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	//////////////////////////////////////////////	ready
	objs, _ := Objects()
	for i := 0; i < 3; i++ {
		err := db.Insert(&objs[i])
		if err != nil {
			log.Println(err)
		}
	}
	idx, err := db.GetIndex("Code", DbTableIdObject{})
	if err != nil {
		log.Println(err)
	}
	it, err := idx.LowerBound(DbTableIdObject{Code: 11})
	if err != nil {
		log.Fatalln(err)
	}

	table := DbTableIdObject{}
	i := 0
	for it.Next() {
		it.Data(&table)
		if objs[i] != table {
			logObj(objs[i])
			logObj(table)
			log.Fatalln("undo failed")
		}
		i++
	}
	session := db.StartSession()

	err = db.Remove(&table)
	if err != nil {
		log.Fatalln(err)
	}
	////////////////////////////////////////// begin
	beginUndo, err := db.GetIndex("Code", DbTableIdObject{})
	if err != nil {
		log.Println(err)
	}
	beginIt, err := beginUndo.LowerBound(DbTableIdObject{Code: 11})
	if err != nil {
		log.Fatalln(err)
	}
	i = 0
	for beginIt.Next() {
		table := DbTableIdObject{}
		beginIt.Data(&table)
		//logObj(table)
		if objs[i] != table {
			logObj(objs[i])
			logObj(table)
			log.Fatalln("undo failed")
		}
		i++
	}
	if i != 2 {
		log.Println(i)
		log.Fatalln("undo failed")
	}
	session.Undo() // undo
	/////////////////////////////////////////// end
	endUndo, err := db.GetIndex("Code", DbTableIdObject{})
	if err != nil {
		log.Println(err)
	}
	endIt, err := endUndo.LowerBound(DbTableIdObject{Code: 11})
	if err != nil {
		log.Fatalln(err)
	}
	i = 0
	for endIt.Next() {
		table := DbTableIdObject{}
		endIt.Data(&table)
		if objs[i] != table {
			logObj(objs[i])
			logObj(table)
			log.Fatalln("undo failed")
		}
		i++
	}
}

func Test_Squash(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	//////////////////////////////////////////////	ready
	objs, _ := Objects()
	for i := 0; i < 3; i++ {
		err := db.Insert(&objs[i])
		if err != nil {
			log.Println(err)
		}
	}
	idx, err := db.GetIndex("Code", DbTableIdObject{})
	if err != nil {
		log.Println(err)
	}
	it, err := idx.LowerBound(DbTableIdObject{Code: 12})
	if err != nil {
		log.Fatalln(err)
	}

	table := DbTableIdObject{}
	tmp :=  objs[1]
	tmp.ID = 2
	i := 1
	for it.Next() {
		it.Data(&table)
		if objs[i] != table {
			logObj(objs[i])
			logObj(table)
			log.Fatalln("undo failed")
		}
		i++
	}


	session := db.StartSession()
	err = db.Remove(&table) 	/*	id --> 3 */
	if err != nil {
		log.Fatalln(err)
	}

	session_ := db.StartSession()
	err = db.Remove(&tmp) 		/*	id --> 2*/
	if err != nil {
		log.Fatalln(err)
	}

	beginIt, err := idx.LowerBound(DbTableIdObject{Code: 11})
	if err != nil {
		log.Fatalln(err)
	}
	i = 0
	for beginIt.Next() {
		table := DbTableIdObject{}
		beginIt.Data(&table)
		if objs[i] != table {
			logObj(objs[i])
			logObj(table)
			log.Fatalln("undo failed")
		}
		i++
	}
	if i != 1 {
		log.Println(i)
		log.Fatalln("undo failed")
	}


	session.Squash()


	 defer session.Undo() 	// undo
	 session_.Undo() 		/* after squash undo all */
	/////////////////////////////////////////// end
	endIt, err := idx.LowerBound(DbTableIdObject{Code: 11})
	if err != nil {
		log.Fatalln(err)
	}
	i = 0
	for endIt.Next() {
		table := DbTableIdObject{}
		endIt.Data(&table)
		if objs[i] != table {
			logObj(objs[i])
			logObj(table)
			log.Fatalln("undo failed")
		}
		i++
	}
}

func Test_iteratorTo(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	//objs_, houses_ :=
	saveObjs(objs, houses, db)

	idx, err := db.GetIndex("id", DbTableIdObject{})
	if err != nil {
		log.Fatalln(err)
	}
	obj := DbTableIdObject{ID: 1}
	it, err := idx.LowerBound(obj)
	if err != nil {
		log.Fatalln(err)
	}

	tmp := DbTableIdObject{}
	for it.Next() {
		it.Data(&tmp)
		//logObj(tmp)
	}
	if !idx.CompareEnd(it) {
		log.Fatalln("CompareEnd failed")
	}
	it.Release()

	it = idx.IteratorTo(&tmp)
	if it == nil {
		log.Panicln("iterator to failed")
	}

	tmp = DbTableIdObject{}
	it.Data(&tmp)
	//logObj(tmp)
	for it.Prev() {
		it.Data(&tmp)
		//logObj(tmp)
	}
	it.Release()
}

func Test_begin(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	//objs_, houses_ :=
	saveObjs(objs, houses, db)

	idx, err := db.GetIndex("id", DbTableIdObject{})
	if err != nil {
		log.Fatalln(err)
	}
	obj := DbTableIdObject{ID: 1}
	it, err := idx.LowerBound(obj)
	if err != nil {
		log.Fatalln(err)
	}

	if !idx.CompareBegin(it) {
		log.Fatalln("begin failed")
	}
	//if !idx.CompareEnd(nil){
	//	log.Fatalln("end failed")
	//}

	it.Release()

	it = idx.Begin()
	if it == nil {
		log.Panicln("iterator to failed")
	}

	it1 := idx.Begin()
	if !idx.CompareBegin(it1) {
		log.Fatalln("begin failed")
	}
	if !idx.CompareIterator(it1, it) {
		log.Fatalln("begin failed")
	}

	it.Release()
	it1.Release()
}

func Test_end(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	//objs_, houses_ :=
	saveObjs(objs, houses, db)

	idx, err := db.GetIndex("id", DbTableIdObject{})
	if err != nil {
		log.Fatalln(err)
	}

	/*less end*/
	obj := DbTableIdObject{ID: 100}
	itLess, err := idx.LowerBound(obj)
	for itLess.Prev() {
		tmp := DbTableIdObject{}
		itLess.Data(&tmp)
		//logObj(tmp)
	}
	itLess.Release()

	/*greater end*/

	idxGreater, err := db.GetIndex("byTable", obj)
	if err != nil {
		log.Fatalln(err)
	}
	obj = DbTableIdObject{Scope: 20}
	itGreater, err := idxGreater.LowerBound(obj)
	if err != nil {
		log.Fatalln(err)
	}

	for itGreater.Prev() {
		tmp := DbTableIdObject{}
		itGreater.Data(&tmp)
		//logObj(tmp)
	}

	if !idxGreater.CompareBegin(itGreater) {
		t.Fatal("Greater compareBegin failed")
	}

	it := idxGreater.Begin()
	if !idxGreater.CompareIterator(it, itGreater) {
		t.Fatal("Greater compareBegin failed")
	}
	it.Release()

	for itGreater.Next() {
		tmp := DbTableIdObject{}
		itGreater.Data(&tmp)
		//logObj(tmp)
	}

	if !idxGreater.CompareEnd(itGreater) {
		t.Fatal("Greater compareEnd failed")
	}
	it = idxGreater.End()
	if !idxGreater.CompareIterator(it, itGreater) {
		t.Fatal("Greater compareEnd failed")
	}
	it.Release()
	itGreater.Release()

	itT := idxGreater.End()
	for itT.Prev() {
		tmp := DbTableIdObject{}
		itT.Data(&tmp)
		//logObj(tmp)
	}

	itT.Release()
}

func Test_empty(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	saveObjs(objs, houses, db)

	idx, err := db.GetIndex("Code", DbTableIdObject{})
	if err != nil {
		log.Fatalln(err)
	}

	obj := DbTableIdObject{Code: 11}
	it, err := idx.LowerBound(obj)
	if err != nil {
		log.Fatalln(err)
	}
	for it.Next() {
		tmp := DbTableIdObject{}
		it.Data(&tmp)
		err = db.Remove(&tmp)
		if err != nil {
			log.Fatalln(err)
		}

	}
	it.Release()

	if !idx.Empty() {
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

func Test_compare(t *testing.T) {

	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	objs_, _ := saveObjs(objs, houses, db)

	idx, err := db.GetIndex("id", DbTableIdObject{})
	if err != nil {
		log.Fatalln(err)
	}

	i := 2
	{
		obj := DbTableIdObject{ID: 3}
		itLess, _ := idx.LowerBound(obj) // note  return err
		if !idx.CompareEnd(itLess) {
			for itLess.Next() {
				tmp := DbTableIdObject{}
				itLess.Data(&tmp)
				if tmp != objs_[i] {
					logObj(objs_[i])
					logObj(tmp)
					t.Fatal("compare error")
				}
				i++
			}
		}
		i--
		if idx.CompareEnd(itLess) {
			for itLess.Prev() {
				tmp := DbTableIdObject{}
				itLess.Data(&tmp)
				if tmp != objs_[i] {
					logObj(objs_[i])
					logObj(tmp)
					t.Fatal("compare error")
				}

				i--
			}
		}
		itLess.Release()
	}
	i = 8
	/*compare less end*/
	obj := DbTableIdObject{ID: 100}
	itLess, _ := idx.LowerBound(obj)
	if idx.CompareEnd(itLess) {
		for itLess.Prev() {
			tmp := DbTableIdObject{}
			itLess.Data(&tmp)
			if tmp != objs_[i] {
				logObj(tmp)
				logObj(objs_[i])
				t.Fatal("compare error")
			}
			i--
		}
	}
	itLess.Release()
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
	for _, v := range limits {
		err := db.Insert(&v)
		if err != nil {
			log.Fatalln(err)
		}
	}

	idx, err := db.GetIndex("byOwner", DbResourceLimitsObject{})
	if err != nil {
		log.Fatalln(err)
	}

	for !idx.Empty() {
		tmp := DbResourceLimitsObject{}
		obj := DbResourceLimitsObject{Pending: false}
		it, err := idx.LowerBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		idx.BeginData(&tmp)
		//logObj(tmp)
		if idx.CompareEnd(it) || tmp.Pending == true {
			log.Fatalln("db is empty")
		}

		err = db.Remove(&tmp)
		if err != nil {
			log.Fatalln(err)
		}
		it.Release()
	}

	if !idx.Empty() {
		log.Fatalln("empty error")
	}
	if idx.Empty() {
		//fmt.Println("empty successful !")
	}
}

func Test_Increment(t *testing.T) {

	fileName := "./increment"

	reFn := func() {
		errs := os.RemoveAll(fileName)
		if errs != nil {
			log.Fatalln(errs)
		}
	}
	defer reFn()
	_, exits := os.Stat(fileName)
	if exits == nil {
		reFn()
	}

	db, err := NewDataBase(fileName, false)
	if err != nil {
		log.Panicln("new database failed : ", err)
	}
	defer db.Close()

	objS, houses := Objects()
	saveObjs(objS, houses, db)

	obj := DbTableIdObject{Scope: 22}
	idx, err := db.GetIndex("byTable", obj)
	if err != nil {
		log.Fatalln(err)
	}

	it, err := idx.LowerBound(obj)
	if err != nil {
		log.Fatalln(err)
	}
	defer it.Release()

	for it.Next() {
		tmp := DbTableIdObject{}
		it.Data(&tmp)
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

	db, err := NewDataBase(fileName, false)
	if err != nil {

		log.Fatalln("new database failed : ", err)
		return nil, reFn
	}

	return db, func() {
		db.Close()
		reFn()
	}
}

func Objects() ([]DbTableIdObject, []DbHouse) {
	objs := []DbTableIdObject{}
	DbHouses := []DbHouse{}
	for i := 1; i <= 3; i++ {
		number := i * 10
		obj := DbTableIdObject{Code: AccountName(number + 1), Scope: ScopeName(number + 2), Table: TableName(number + 3 + i + 1), Payer: AccountName(number + 4 + i + 1), Count: uint32(number + 5)}
		objs = append(objs, obj)
		house := DbHouse{Area: uint64(number + 7), Carnivore: Carnivore{number + 7, number + 7}}
		DbHouses = append(DbHouses, house)
		obj = DbTableIdObject{Code: AccountName(number + 2), Scope: ScopeName(number + 2), Table: TableName(number + 3 + i + 2), Payer: AccountName(number + 4 + i + 2), Count: uint32(number + 5)}
		objs = append(objs, obj)
		house = DbHouse{Area: uint64(number + 8), Carnivore: Carnivore{number + 8, number + 8}}
		DbHouses = append(DbHouses, house)

		obj = DbTableIdObject{Code: AccountName(number + 3), Scope: ScopeName(number + 2), Table: TableName(number + 3 + i + 3), Payer: AccountName(number + 4 + i + 3), Count: uint32(number + 5)}
		objs = append(objs, obj)
		house = DbHouse{Area: uint64(number + 9), Carnivore: Carnivore{number + 9, number + 9}}
		DbHouses = append(DbHouses, house)
	}
	return objs, DbHouses
}

func saveObjs(objs []DbTableIdObject, houses []DbHouse, db DataBase) ([]DbTableIdObject, []DbHouse) {
	objs_ := []DbTableIdObject{}
	houses_ := []DbHouse{}

	for _, v := range houses {
		//logObj(v)
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

func lowerAndUpper(objs []DbTableIdObject, houses []DbHouse, db DataBase) {

	obj := DbTableIdObject{Scope: 22}
	idx, err := db.GetIndex("byTable", obj)
	if err != nil {
		log.Fatalln(err)
	}

	it, err := idx.LowerBound(obj)
	if err != nil {
		log.Fatalln(err)
	}
	defer it.Release()

	i := 3
	for it.Next() {
		tmp := DbTableIdObject{}
		it.Data(&tmp)
		if tmp != objs[i] {
			logObj(tmp)
			logObj(objs[i])
			log.Fatalln("getGreaterObjs ")
		}
		i++
	}
	if !idx.CompareEnd(it) {
		log.Fatalln("CompareEnd")
	}
	it.Release()

	it, err = idx.UpperBound(obj)
	if err != nil {
		log.Fatalln(err)
	}
	i = 6
	for it.Next() {
		tmp := DbTableIdObject{}
		it.Data(&tmp)
		if tmp != objs[i] {
			logObj(tmp)
			logObj(objs[i])
			log.Fatalln("getGreaterObjs ")
		}
		i++
	}

	it.Release()

	obj = DbTableIdObject{Scope: 202}
	it1, err := idx.LowerBound(obj)
	if err != nil {
		log.Fatalln(err)
	}

	if !idx.CompareEnd(it1) {
		log.Fatalln("getGreaterObjs ")
	}
	i = 8
	for it1.Prev() {
		tmp := DbTableIdObject{}
		it1.Data(&tmp)
		if tmp != objs[i] {
			logObj(tmp)
			logObj(objs[i])
			log.Fatalln("getGreaterObjs ")
		}
		i--
		//logObj(tmp)
	}

	if !idx.CompareBegin(it1) {
		log.Fatalln("getGreaterObjs ")
	}

	it1.Release()
}

func getLessObjs(objs []DbTableIdObject, houses []DbHouse, db DataBase) {
	obj := DbTableIdObject{Code: 13}
	idx, err := db.GetIndex("Code", DbTableIdObject{})
	if err != nil {
		log.Fatalln(err)
	}

	it, err := idx.LowerBound(obj)
	if err != nil {
		log.Fatalln(err)
	}
	i := 2
	for it.Next() {
		tmp := DbTableIdObject{}
		it.Data(&tmp)
		if tmp != objs[i] {
			logObj(objs[i])
			logObj(tmp)
		}
		i++
	}
	it.Release()

	idx, err = db.GetIndex("id", DbTableIdObject{})
	if err != nil {
		log.Fatalln(err)
	}
	obj = DbTableIdObject{ID: 1}
	it, err = idx.LowerBound(obj)
	if err != nil {
		log.Fatalln(err)
	}
	i = 0
	for it.Next() {
		tmp := DbTableIdObject{}
		it.Data(&tmp)
		//logObj(tmp)
		if tmp != objs[i] {
			logObj(objs[i])
			logObj(tmp)
		}
		i++
	}
	it.Release()
}

func modifyObjs(db DataBase) {

	obj := DbTableIdObject{ID: 4, Code: 21, Scope: 22, Table: 26, Payer: 27, Count: 25}
	newobj := DbTableIdObject{ID: 4, Code: 10199, Scope: 22, Table: 26, Payer: 27, Count: 25}
	for i := 0; i < 10000; i++ {
		err := db.Modify(&obj, func(object *DbTableIdObject) {
			object.Code = AccountName(200 + i)
		})
		if err != nil {
			log.Fatalln(err)
		}
	}

	obj = DbTableIdObject{}
	tmp := DbTableIdObject{}
	obj.ID = 4
	err := db.Find("id", obj, &tmp)
	if err != nil {
		log.Fatalln(err)
	}
	if tmp != newobj {
		logObj(tmp)
		log.Fatalln("modify test error")
	}
}

func findObjs(objs []DbTableIdObject, houses []DbHouse, db DataBase) {
	obj := DbTableIdObject{ID: 4}
	tmp := DbTableIdObject{}
	err := db.Find("id", obj, &tmp)
	if err != nil {
		log.Fatalln(err)
	}
	if tmp != objs[3] {
		log.Fatalln("Find Object")
	}
	//logObj(tmp)
	{
		hou := DbHouse{Area: 18}
		tmp := DbHouse{}
		err := db.Find("Area", hou, &tmp)
		if err != nil {
			log.Fatalln(err)
		}
		if houses[1] != tmp {
			logObj(tmp)
			logObj(houses[1])
			log.Fatalln("Find Object")
		}
	}

	{
		idx, err := db.GetIndex("Area", DbHouse{})
		if err != nil {
			log.Fatalln(err)
		}

		hou := DbHouse{Area: 18}
		it, err := idx.LowerBound(hou)
		if err != nil {
			log.Fatalln(err)
		}
		i := 1
		for it.Next() {
			tmp := DbHouse{}
			it.Data(&tmp)
			if houses[i] != tmp {
				logObj(tmp)
				logObj(houses[i])
				log.Fatalln("Find Object")
			}
			i++
		}
		it.Release()
	}
}

func findInLineFieldObjs(objs []DbTableIdObject, houses []DbHouse, db DataBase) {
	hou := DbHouse{Carnivore: Carnivore{28, 38}}
	//idx,err := db.GetIndex("Tiger", hou)
	idx, err := db.GetIndex("Lion", hou)
	if err != nil {
		log.Fatalln(err)
	}

	it, err := idx.LowerBound(hou)
	if err != nil {
		log.Fatalln(err)
	}
	i := 4
	defer it.Release()
	for it.Next() {
		tmp := DbHouse{}
		it.Data(&tmp)
		if tmp != houses[i] {
			logObj(houses[i])
			logObj(tmp)
			log.Fatalln("findInLineFieldObjs")
		}
		i++
	}
}

func removeUnique(db DataBase) {
	obj := DbTableIdObject{Code: 21, Scope: 22, Table: 23, Payer: 24, Count: 25}
	err := db.Remove(&obj)
	if err != ErrIncompleteStructure {
		log.Fatalln(err)
	}

	obj = DbTableIdObject{}
	obj.ID = 4

	tmp := DbTableIdObject{}
	err = db.Find("id", obj, &tmp)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Remove(&tmp)
	if err != nil {
		log.Fatalln(err)
	}

	tmp = DbTableIdObject{}
	err = db.Find("id", obj, &tmp)
	if err == nil {
		log.Fatalln("test remove failed")
	}
}

func MakeResourceLimitsObjects() []DbResourceLimitsObject {
	//limits := make([]DbResourceLimitsObject,0)
	limits := []DbResourceLimitsObject{}

	for i := 1; i <= 13; i++ {
		number := 100
		obj := DbResourceLimitsObject{Owner: AccountName(number + i)}
		limits = append(limits, obj)
	}
	return limits
}
