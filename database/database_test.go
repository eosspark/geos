package database

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"log"
	"os"
	"testing"
)

var logFlag = false
//var logFlag = true

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
	for i := 1; i <= 720000; i++ {
		key := []byte(string(i))
		key = append(key, key...)
		keys = append(keys, key)
	}
	for _, v := range keys {
		db.Put(v, v, nil)
	}
	h := []byte("hello")
	w := []byte("world")
	db.Put(h,h,nil)
	db.Put(w,w,nil)
	it := db.NewIterator(util.BytesPrefix([]byte(string("linx"))),nil)
	if it.Next(){
		fmt.Println(it.Key() ," : ",it.Value())
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

	objs, houses := multiObjects()
	if len(objs) != len(houses) {
		log.Fatalln("ERROR")
	}

	saveObjs(objs, houses, db)

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

	tmp := DbTableIdObject{}
	//fmt.Println("---------------------------")

	db.Find("id",DbTableIdObject{ID:254},&tmp)
	//logObj(tmp)

	db.Find("id",DbTableIdObject{ID:255},&tmp)
	//logObj(tmp)

	tmp = DbTableIdObject{}
	db.Find("id",DbTableIdObject{ID:25590000},&tmp)
	//logObj(tmp)

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
	{
		// Code 11 12 13  21 22 23  31 32 33
		it, err := idx.LowerBound(DbTableIdObject{Code: 11})
		if err != nil {
			log.Fatalln(err)
		}
		i := 0
		tmp := DbTableIdObject{}
		it.Data(&tmp)
		if tmp != objs_[i] {
			logObj(objs_[i])
			logObj(tmp)
			log.Fatalln("failed")
		}
		for it.Next() {
			i++
			tmp := DbTableIdObject{}
			it.Data(&tmp)
			if tmp != objs_[i] {
				logObj(objs_[i])
				logObj(tmp)
				log.Fatalln("failed")
			}
		}
		it.Release()
	}


	session := db.StartSession()
	defer session.Undo()
	obj := DbTableIdObject{ID: 4, Code: 21, Scope: 22, Table: 26, Payer: 27, Count: 25}
	newObj := DbTableIdObject{ID: 4, Code: 200, Scope: 22, Table: 26, Payer: 27, Count: 25}

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
	if tmp == newObj {
		logObj(newObj)
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
	// Insert three elements
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

	// Code 11 12 13
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
	tmp := DbTableIdObject{}
	it.Data(&tmp)
	if tmp != objs[i] {
		logObj(objs[i])
		logObj(tmp)
		log.Fatalln("failed")
	}
	for it.Next() {
		i++
		tmp := DbTableIdObject{}
		it.Data(&tmp)
		if tmp != objs[i] {
			logObj(objs[i])
			logObj(tmp)
			log.Fatalln("failed")
		}
	}
	it.Release()
}

func Test_undoRemove(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	// Insert three elements
	objs, _ := Objects()
	for i := 0; i < 3; i++ {
		err := db.Insert(&objs[i])
		if err != nil {
			log.Println(err)
		}
	}
	table := DbTableIdObject{}
	// review
	{
		idx, err := db.GetIndex("Code", DbTableIdObject{})
		if err != nil {
			log.Println(err)
		}
		// Code 11 12 13
		it, err := idx.LowerBound(DbTableIdObject{Code: 11})
		if err != nil {
			log.Fatalln(err)
		}

		i := 0
		it.Data(&table)
		if objs[i] != table {
			logObj(objs[i])
			logObj(table)
			log.Fatalln("undo failed")
		}
		for it.Next() {
			i++
			table := DbTableIdObject{}
			it.Data(&table)
			if objs[i] != table {
				logObj(objs[i])
				logObj(table)
				log.Fatalln("undo failed")
			}
		}
		it.Release()
	}
	// remove and undo
	{
		session := db.StartSession()

		err := db.Remove(&table) // code 11 remove
		if err != nil {
			log.Fatalln(err)
		}
		/*  Code 11 12 13   11 is begin */
		beginUndo, err := db.GetIndex("Code", DbTableIdObject{})
		if err != nil {
			log.Println(err)
		}
		beginIt, err := beginUndo.LowerBound(DbTableIdObject{Code: 11})
		if err != nil {
			log.Fatalln(err)
		}
		i := 1
		table = DbTableIdObject{}
		beginIt.Data(&table)
		if objs[i] != table {
			logObj(objs[i])
			logObj(table)
			log.Fatalln("undo failed")
		}
		for beginIt.Next() {
			i++
			table := DbTableIdObject{}
			beginIt.Data(&table)
			if objs[i] != table {
				logObj(objs[i])
				logObj(table)
				log.Fatalln("undo failed")
			}
		}
		if i != 2 {
			log.Println(i)
			log.Fatalln("undo failed")
		}
		beginIt.Release()
		session.Undo() // undo
	}
	// review again
	{
		endUndo, err := db.GetIndex("Code", DbTableIdObject{})
		if err != nil {
			log.Println(err)
		}
		endIt, err := endUndo.LowerBound(DbTableIdObject{Code: 11})
		if err != nil {
			log.Fatalln(err)
		}
		i := 0
		table = DbTableIdObject{}
		endIt.Data(&table)
		if objs[i] != table {
			logObj(objs[i])
			logObj(table)
			log.Fatalln("undo failed")
		}
		for endIt.Next() {
			i++
			table := DbTableIdObject{}
			endIt.Data(&table)
			if objs[i] != table {
				logObj(objs[i])
				logObj(table)
				log.Fatalln("undo failed")
			}
		}
		endIt.Release()
	}

}

func Test_squash(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	// Insert three elements
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
	// Code 11 12 13
	it, err := idx.LowerBound(DbTableIdObject{Code: 12})
	if err != nil {
		log.Fatalln(err)
	}


	i := 1
	table := DbTableIdObject{}
	it.Data(&table)
	if objs[i] != table {
		logObj(objs[i])
		logObj(table)
		log.Fatalln("undo failed")
	}
	for it.Next() {
		i++
		it.Data(&table)
		if objs[i] != table {
			logObj(objs[i])
			logObj(table)
			log.Fatalln("undo failed")
		}
	}


	session := db.StartSession()
	err = db.Remove(&table) 	/*	id --> 3 */
	if err != nil {
		log.Fatalln(err)
	}

	session_ := db.StartSession()
	tmp :=  objs[1]
	err = db.Remove(&tmp) 		/*	id --> 2*/
	if err != nil {
		log.Fatalln(err)
	}
	// There is only one element in the database at this time
	beginIt, err := idx.LowerBound(DbTableIdObject{Code: 11})
	if err != nil {
		log.Fatalln(err)
	}

	if !idx.CompareBegin(beginIt){
		log.Fatalln("iterator is not begin")
	}

	if beginIt.Next() {
		log.Fatalln("iterator next")
	}
	if !idx.CompareEnd(beginIt){
		log.Fatalln("iterator is not end")
	}

	session.Squash()


	defer session.Undo() 	// undo
	session_.Undo() 		/* after squash undo all */

	endIt, err := idx.LowerBound(DbTableIdObject{Code: 11})
	if err != nil {
		log.Fatalln(err)
	}
	i = 0
	table = DbTableIdObject{}
	endIt.Data(&table)
	if objs[i] != table {
		logObj(objs[i])
		logObj(table)
		log.Fatalln("undo failed")
	}
	for endIt.Next() {
		i++
		table := DbTableIdObject{}
		endIt.Data(&table)
		if objs[i] != table {
			logObj(objs[i])
			logObj(table)
			log.Fatalln("undo failed")
		}
	}
	endIt.Release()
}

func Test_undoAll(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, _ := Objects()
	// Insert three elements
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

	i := 1
	table := DbTableIdObject{}
	it.Data(&table)
	if objs[i] != table {
		logObj(objs[i])
		logObj(table)
		log.Fatalln("undo failed")
	}
	for it.Next() {
		i++
		it.Data(&table)
		if objs[i] != table {
			logObj(objs[i])
			logObj(table)
			log.Fatalln("undo failed")
		}
	}

	session := db.StartSession()
	defer session.Undo() 	// undo
	err = db.Remove(&table) 	/*	id --> 3 */
	if err != nil {
		log.Fatalln(err)
	}

	session_ := db.StartSession()
	defer session_.Undo() 		/* after squash undo all */

	tmp :=  objs[1]
	err = db.Remove(&tmp) 		/*	id --> 2*/
	if err != nil {
		log.Fatalln(err)
	}


	session.Squash()
	db.UndoAll()

	endIt, err := idx.LowerBound(DbTableIdObject{Code: 11})
	if err != nil {
		log.Fatalln(err)
	}
	i = 0
	table = DbTableIdObject{}
	endIt.Data(&table)
	if objs[i] != table {
		logObj(objs[i])
		logObj(table)
		log.Fatalln("undo failed")
	}
	for endIt.Next() {
		i++
		table := DbTableIdObject{}
		endIt.Data(&table)
		if objs[i] != table {
			logObj(objs[i])
			logObj(table)
			log.Fatalln("undo failed")
		}
	}
	endIt.Release()
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

func Test_uniqueNoIterator(t *testing.T) {

	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	objs_,_ := saveObjs(objs, houses, db)
	idx, err := db.GetIndex("byTable", DbTableIdObject{})
	if err != nil {
		log.Fatalln(err)
	}

	// scope 11 lower prev iterator
	{
		obj := DbTableIdObject{Scope: 11}
		scope11, err := idx.LowerBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		if idx.CompareEnd(scope11){
			log.Fatalln("scope 11 is end iterator")
		}

		if !idx.CompareBegin(scope11)	{
			log.Fatalln("scope compare end failed")
		}

		i := 0
		tmp := DbTableIdObject{}
		scope11.Data(&tmp)
		if objs_[i] != tmp{
			logObj(objs_[i])
			logObj(tmp)
			log.Fatalln("compare failed")
		}

		for scope11.Prev() {
			log.Fatalln("compare failed")
		}
	}
	// scope 11 lower next iterator
	{
		obj := DbTableIdObject{Scope: 11}
		scope11, err := idx.LowerBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		if idx.CompareEnd(scope11){
			log.Fatalln("scope 11 is end iterator")
		}

		if !idx.CompareBegin(scope11)	{
			log.Fatalln("scope compare end failed")
		}
		i := 0
		tmp := DbTableIdObject{}
		scope11.Data(&tmp)
		if objs_[i] != tmp{
			logObj(objs_[i])
			logObj(tmp)
			log.Fatalln("compare failed")
		}

		for scope11.Next(){
			i++
			tmp := DbTableIdObject{}
			scope11.Data(&tmp)
			if objs_[i] != tmp{
				logObj(objs_[i])
				logObj(tmp)
				log.Fatalln("compare failed")
			}
		}
	}
	// scope 11 upper prev iterator
	{
		obj := DbTableIdObject{Scope: 11}
		scope11, err := idx.UpperBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		if idx.CompareEnd(scope11){
			log.Fatalln("scope 11 is end iterator")
		}

		if !idx.CompareBegin(scope11) {
			log.Fatalln("scope compare end failed")
		}
		i := 0
		tmp := DbTableIdObject{}
		scope11.Data(&tmp)
		if objs_[i] != tmp{
			logObj(objs_[i])
			logObj(tmp)
			log.Fatalln("compare failed")
		}

		for scope11.Prev(){
			log.Fatalln("compare failed")
		}
	}
	// scope 11 upper next iterator
	{
		obj := DbTableIdObject{Scope: 11}
		scope11, err := idx.UpperBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		if idx.CompareEnd(scope11){
			log.Fatalln("scope 11 is end iterator")
		}

		if !idx.CompareBegin(scope11)	{
			log.Fatalln("scope compare end failed")
		}
		i := 0
		tmp := DbTableIdObject{}
		scope11.Data(&tmp)
		if objs_[i] != tmp{
			logObj(objs_[i])
			logObj(tmp)
			log.Fatalln("compare failed")
		}

		for scope11.Next(){
			i++
			tmp := DbTableIdObject{}
			scope11.Data(&tmp)
			if objs_[i] != tmp{
				logObj(objs_[i])
				logObj(tmp)
				log.Fatalln("compare failed")
			}
		}
	}


	/* scope 30 lower prev iterator */
	{
		obj := DbTableIdObject{Scope: 30}
		scope11, err := idx.LowerBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		if idx.CompareEnd(scope11){
			log.Fatalln("scope is end")
		}
		if idx.CompareBegin(scope11){
			log.Fatalln("scope is begin")
		}
		i := 6
		tmp := DbTableIdObject{}
		scope11.Data(&tmp)
		if objs_[i] != tmp{
			logObj(objs_[i])
			logObj(tmp)
			log.Fatalln("compare failed")
		}

		for scope11.Prev() {
			i--
			tmp := DbTableIdObject{}
			scope11.Data(&tmp)
			if objs_[i] != tmp{
				logObj(objs_[i])
				logObj(tmp)
				log.Fatalln("compare failed")
			}
		}
	}
	/* scope 30 lower next iterator */
	{
		obj := DbTableIdObject{Scope: 30}
		scope11, err := idx.LowerBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		if idx.CompareEnd(scope11){
			log.Fatalln("scope is end")
		}
		if idx.CompareBegin(scope11){
			log.Fatalln("scope is begin")
		}
		tmp := DbTableIdObject{}
		scope11.Data(&tmp)
		i := 6
		if objs_[i] != tmp{
			logObj(objs_[i])
			logObj(tmp)
			log.Fatalln("compare failed")
		}
		for scope11.Next() {
			i++
			tmp := DbTableIdObject{}
			scope11.Data(&tmp)
			if objs_[i] != tmp{
				logObj(objs_[i])
				logObj(tmp)
				log.Fatalln("compare failed")
			}
		}
	}
	/* scope 30 upper prev iterator */
	{
		obj := DbTableIdObject{Scope: 30}
		scope11, err := idx.UpperBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		if idx.CompareEnd(scope11){
			log.Fatalln("scope is end")
		}
		if idx.CompareBegin(scope11){
			log.Fatalln("scope is begin")
		}

		tmp := DbTableIdObject{}
		scope11.Data(&tmp)
		i := 6
		if objs_[i] != tmp{
			logObj(objs_[i])
			logObj(tmp)
			log.Fatalln("compare failed")
		}

		for scope11.Prev() {
			i--
			tmp := DbTableIdObject{}
			scope11.Data(&tmp)
			if objs_[i] != tmp{
				logObj(objs_[i])
				logObj(tmp)
				log.Fatalln("compare failed")
			}
		}
	}
	/* scope 30 upper next iterator */
	{
		obj := DbTableIdObject{Scope: 30}
		scope11, err := idx.UpperBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		if idx.CompareEnd(scope11){
			log.Fatalln("scope is end")
		}
		if idx.CompareBegin(scope11){
			log.Fatalln("scope is begin")
		}

		tmp := DbTableIdObject{}
		scope11.Data(&tmp)
		i := 6
		if objs_[i] != tmp{
			logObj(objs_[i])
			logObj(tmp)
			log.Fatalln("compare failed")
		}

		for scope11.Next() {
			i++
			tmp := DbTableIdObject{}
			scope11.Data(&tmp)
			if objs_[i] != tmp{
				logObj(objs_[i])
				logObj(tmp)
				log.Fatalln("compare failed")
			}
		}
	}

	// scope 22 lower prev iterator
	{
		obj := DbTableIdObject{Scope: 22}
		scope11, err := idx.LowerBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		if idx.CompareEnd(scope11){
			log.Fatalln("scope is end")
		}
		if idx.CompareBegin(scope11){
			log.Fatalln("scope is begin")
		}
		tmp := DbTableIdObject{}
		scope11.Data(&tmp)
		i := 3
		if objs_[i] != tmp{
			logObj(objs_[i])
			logObj(tmp)
			log.Fatalln("compare failed")
		}
		for scope11.Prev() {
			i--
			tmp := DbTableIdObject{}
			scope11.Data(&tmp)
			if objs_[i] != tmp{
				logObj(objs_[i])
				logObj(tmp)
				log.Fatalln("compare failed")
			}
		}
	}
	// scope 22 lower next iterator
	{
		obj := DbTableIdObject{Scope: 22}
		scope11, err := idx.LowerBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		if idx.CompareEnd(scope11){
			log.Fatalln("scope is end")
		}
		if idx.CompareBegin(scope11){
			log.Fatalln("scope is begin")
		}
		tmp := DbTableIdObject{}
		scope11.Data(&tmp)
		i := 3
		if objs_[i] != tmp{
			logObj(objs_[i])
			logObj(tmp)
			log.Fatalln("compare failed")
		}
		for scope11.Next() {
			i++
			tmp := DbTableIdObject{}
			scope11.Data(&tmp)
			if objs_[i] != tmp{
				logObj(objs_[i])
				logObj(tmp)
				log.Fatalln("compare failed")
			}
		}
	}
	// scope 22 upper prev iterator
	{
		obj := DbTableIdObject{Scope: 22}
		scope11, err := idx.UpperBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		if idx.CompareEnd(scope11){
			log.Fatalln("scope is end")
		}
		if idx.CompareBegin(scope11){
			log.Fatalln("scope is begin")
		}
		tmp := DbTableIdObject{}
		scope11.Data(&tmp)
		i := 6
		if objs_[i] != tmp{
			logObj(objs_[i])
			logObj(tmp)
			log.Fatalln("compare failed")
		}
		for scope11.Prev() {
			i--
			tmp := DbTableIdObject{}
			scope11.Data(&tmp)
			if objs_[i] != tmp{
				logObj(objs_[i])
				logObj(tmp)
				log.Fatalln("compare failed")
			}
		}
	}
	// scope 22 upper next iterator
	{
		obj := DbTableIdObject{Scope: 22}
		scope11, err := idx.UpperBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		if idx.CompareEnd(scope11){
			log.Fatalln("scope is end")
		}
		if idx.CompareBegin(scope11){
			log.Fatalln("scope is begin")
		}
		tmp := DbTableIdObject{}
		scope11.Data(&tmp)
		i := 6
		if objs_[i] != tmp{
			logObj(objs_[i])
			logObj(tmp)
			log.Fatalln("compare failed")
		}
		for scope11.Next() {
			i++
			tmp := DbTableIdObject{}
			scope11.Data(&tmp)
			if objs_[i] != tmp{
				logObj(objs_[i])
				logObj(tmp)
				log.Fatalln("compare failed")
			}
		}
	}

	/* compare iterator test */
	{
		obj := DbTableIdObject{Scope: 20}
		scope20, err := idx.LowerBound(obj)
		if err != nil {
			log.Fatalln(err)
		}
		if idx.CompareEnd(scope20) {
			log.Fatalln("iterator is end")
		}

		for scope20.Prev() {
			/* go begin */
		}

		if !idx.CompareBegin(scope20) {
			log.Fatalln("iterator compare failed")
		}

		it := idx.Begin()
		if !idx.CompareIterator(it, scope20) {
			log.Fatalln("iterator compare failed")
		}
		it.Release()

		for scope20.Next() {
			/* go end */
		}

		if !idx.CompareEnd(scope20) {
			log.Fatalln("iterator compare failed")
		}
		it = idx.End()
		if !idx.CompareIterator(it, scope20) {
			log.Fatalln("iterator compare failed")
		}
		it.Release()
		scope20.Release()

		itT := idx.End()
		i := 8
		for itT.Prev() {
			tmp := DbTableIdObject{}
			itT.Data(&tmp)
			if objs_[i] != tmp{
				logObj(objs_[i])
				logObj(tmp)
				log.Fatalln("compare failed")
			}
			i--
		}

		itT.Release()
	}

}

func Test_uniqueIterator(t *testing.T) {
	db, clo := openDb()
	if db == nil {
		log.Fatalln("db open failed")
	}
	defer clo()

	objs, houses := Objects()
	objs_, _ := saveObjs(objs, houses, db)

	{
		idx, err := db.GetIndex("id", DbTableIdObject{})
		if err != nil {
			log.Fatalln(err)
		}
		obj := DbTableIdObject{ID: 5}
		// id 5 upper bound pre iterator
		{
			upIt, err := idx.UpperBound(obj)
			if err != nil {
				log.Fatalln(err)
			}

			i := 5
			tmp := DbTableIdObject{}
			upIt.Data(&tmp)
			if objs_[i]	 != tmp{
				logObj(tmp)
				logObj(objs_[i])
				log.Fatalln(err)
			}
			for upIt.Prev() {
				i--
				tmp := DbTableIdObject{}
				upIt.Data(&tmp)
				if objs_[i]	 != tmp{
					logObj(tmp)
					logObj(objs_[i])
					log.Fatalln(err)
				}
			}
			upIt.Release()
		}
		// id 5 upper bound next iterator
		{
			upIt, err := idx.UpperBound(obj)
			if err != nil {
				log.Fatalln(err)
			}

			tmp := DbTableIdObject{}
			upIt.Data(&tmp)
			i := 5
			if objs_[i]	 != tmp{
				logObj(tmp)
				logObj(objs_[i])
				log.Fatalln(err)
			}

			for upIt.Next() {
				i++
				tmp := DbTableIdObject{}
				upIt.Data(&tmp)
				if objs_[i]	 != tmp{
					logObj(tmp)
					logObj(objs_[i])
					log.Fatalln(err)
				}
			}
			upIt.Release()

		}
		// id 5 lower bound pre iterator
		{
			lowIt, err := idx.LowerBound(obj)
			if err != nil{
				log.Fatalln(err)
			}

			tmp := DbTableIdObject{}
			lowIt.Data(&tmp)
			i := 4
			if objs_[i]	 != tmp{
				logObj(tmp)
				logObj(objs_[i])
				log.Fatalln(err)
			}
			for lowIt.Prev() {
				i--
				tmp := DbTableIdObject{}
				lowIt.Data(&tmp)
				if objs_[i]	 != tmp{
					logObj(tmp)
					logObj(objs_[i])
					log.Fatalln(err)
				}
			}
		}
		// id 5 lower bound next iterator
		{
			lowIt, err := idx.LowerBound(obj)
			if err != nil{
				log.Fatalln(err)
			}

			tmp := DbTableIdObject{}
			lowIt.Data(&tmp)
			i := 4
			if objs_[i]	 != tmp{
				logObj(tmp)
				logObj(objs_[i])
				log.Fatalln(err)
			}

			for lowIt.Next() {
				i++
				tmp := DbTableIdObject{}
				lowIt.Data(&tmp)
				if objs_[i]	 != tmp{
					logObj(tmp)
					logObj(objs_[i])
					log.Fatalln(err)
				}
			}
		}
	}
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

	tmp := DbTableIdObject{}
	it.Data(&tmp)
	err = db.Remove(&tmp)
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

	db, err := NewDataBase(fileName, logFlag)
	if err != nil {

		log.Fatalln("new database failed : ", err)
		return nil, reFn
	}

	return db, func() {
		db.Close()
		reFn()
	}
}

func multiObjects() ([]DbTableIdObject, []DbHouse) {
	objs := []DbTableIdObject{}
	DbHouses := []DbHouse{}
	for i := 1; i <= 30000; i++ {
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
		if v.ID == 253 {
			//fmt.Println("go")
		}

		objs_ = append(objs_, v)
	}
	return objs_, houses_
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
	tmp := DbTableIdObject{}
	it.Data(&tmp)
	i := 2
	if tmp != objs[i] {
		logObj(objs[i])
		logObj(tmp)
		log.Fatalln("failed")
	}
	for it.Next() {
		i++
		tmp := DbTableIdObject{}
		it.Data(&tmp)
		if tmp != objs[i] {
			logObj(objs[i])
			logObj(tmp)
			log.Fatalln("failed")
		}
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
	tmp = DbTableIdObject{}
	it.Data(&tmp)
	i = 0
	if tmp != objs[i] {
		logObj(objs[i])
		logObj(tmp)
		log.Fatalln("failed")
	}
	for it.Next() {
		i++
		tmp := DbTableIdObject{}
		it.Data(&tmp)
		if tmp != objs[i] {
			logObj(objs[i])
			logObj(tmp)
			log.Fatalln("failed")
		}
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
	// find id
	{
		/* id  1 2 3 4 5 6 7 8 */
		err := db.Find("id", obj, &tmp)
		if err != nil {
			log.Fatalln(err)
		}
		if tmp != objs[3] {
			log.Fatalln("Find Object")
		}
	}
	// find Area
	{
		/* Area 17 18 19  	27 28 29 		37 38 39*/
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

	// getIndex Area
	{
		/* Area 17 18 19  	27 28 29 		37 38 39*/
		idx, err := db.GetIndex("Area", DbHouse{})
		if err != nil {
			log.Fatalln(err)
		}
		hou := DbHouse{Area: 18}
		it, err := idx.LowerBound(hou)
		if err != nil {
			log.Fatalln(err)
		}
		tmp := DbHouse{}
		it.Data(&tmp)
		i := 1
		if houses[i] != tmp {
			logObj(tmp)
			logObj(houses[i])
			log.Fatalln("Find Object")
		}

		for it.Next() {
			i++
			tmp := DbHouse{}
			it.Data(&tmp)
			if houses[i] != tmp {
				logObj(tmp)
				logObj(houses[i])
				log.Fatalln("Find Object")
			}
		}
		it.Release()
	}
}

func findInLineFieldObjs(objs []DbTableIdObject, houses []DbHouse, db DataBase) {

	/* Lion 17 18 19  	27 28 29 		37 38 39*/
	hou := DbHouse{Carnivore: Carnivore{28, 38}}
	idx, err := db.GetIndex("Lion", hou)
	if err != nil {
		log.Fatalln(err)
	}

	it, err := idx.LowerBound(hou)
	if err != nil {
		log.Fatalln(err)
	}
	i := 4
	tmp := DbHouse{}
	it.Data(&tmp)
	if tmp != houses[i] {
		logObj(houses[i])
		logObj(tmp)
		log.Fatalln("findInLineFieldObjs")
	}
	for it.Next() {
		i++
		tmp := DbHouse{}
		it.Data(&tmp)
		if tmp != houses[i] {
			logObj(houses[i])
			logObj(tmp)
			log.Fatalln("findInLineFieldObjs")
		}
	}
	it.Release()
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
	if err != ErrNotFound {
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
