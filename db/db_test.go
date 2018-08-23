package main

import (
	"fmt"
	"github.com/db"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

// This test case only tests db, the session is not tested
// because they have the same usage, and subsequent test cases will be added to the session.

var names = []string{"linx", "garytone", "elk", "fox", "lion"}

type AccountObject struct {
	Id   uint64 `storm:"id"`
	Name string `storm:index`
	Tag  uint64
}

type User struct {
	Id   uint64 `storm:"id"`
	Name string `storm:"unique"`
	Tag  uint64 `storm:"index"`
}

func Test_Instance(t *testing.T) {
	db, err := eosiodb.NewDatabase("./", "eos.db", true)
	require.NoError(t, err)
	defer db.Close()
}

func TestWriteOne(t *testing.T) {
	db, err := eosiodb.NewDatabase("./", "eos.db", true)
	require.NoError(t, err)
	defer db.Close()

	obj := &AccountObject{Id: 10, Name: "linx", Tag: 99}
	err = db.Insert(obj)
	require.NoError(t, err)

	obj_ := &AccountObject{Id: 10, Name: "qieqie", Tag: 99}
	err = db.Insert(obj_)
	require.NoError(t, err)
}

func TestFind(t *testing.T) {
	db, err := eosiodb.NewDatabase("./", "eos.db", true)
	require.NoError(t, err)
	defer db.Close()

	raw := AccountObject{Id: 10, Name: "qieqie", Tag: 99}
	account := AccountObject{Id: 10, Name: "garytone", Tag: 10}
	var obj AccountObject
	err = db.Find("Name", "qieqie", &obj)
	require.NoError(t, err)
	require.Equal(t, raw, obj)
	require.NotEqual(t, account, obj)
}

func TestInsertSome(t *testing.T) {
	db, err := eosiodb.NewDatabase("./", "eos.db", true)
	require.NoError(t, err)
	defer db.Close()

	for key, value := range names {
		obj := &AccountObject{Id: uint64(key + 11), Name: value, Tag: uint64(10)}
		err = db.Insert(obj)
		require.NoError(t, err)
	}
}

func TestGet(t *testing.T) {
	db, err := eosiodb.NewDatabase("./", "eos.db", true)
	require.NoError(t, err)
	defer db.Close()

	var objs []AccountObject
	err = db.Get("Tag", 10, &objs)
	require.NoError(t, err)
	fmt.Println(len(objs))
	require.Equal(t, len(objs), 5)
}

func TestAll(t *testing.T) {
	db, err := eosiodb.NewDatabase("./", "eos.db", true)
	require.NoError(t, err)
	defer db.Close()

	var objs []AccountObject
	err = db.All(&objs)
	require.NoError(t, err)
	fmt.Println(len(objs))
	require.Equal(t, len(objs), 6)
}

func TestUpdateItem(t *testing.T) {
	db, err := eosiodb.NewDatabase("./", "eos.db", true)
	require.NoError(t, err)
	defer db.Close()

	var obj AccountObject
	err = db.Find("Name", "qieqie", &obj)
	require.NoError(t, err)
	fmt.Println(obj)

	fn := func(data interface{}) error {
		ref := reflect.ValueOf(data).Elem()
		if ref.CanSet() {
			ref.Field(1).SetString("hello")
			ref.Field(2).SetUint(1000)
		} else {
			// log ?
		}
		return nil
	}
	err = db.Update(&obj, fn)
	require.NoError(t, err)
}

func TestUpdateField(t *testing.T) {
	db, err := eosiodb.NewDatabase("./", "eos.db", true)
	require.NoError(t, err)
	defer db.Close()

	//obj := &AccountObject{Id: 10, Name: "hello", Tag: 1000}
	//err = db.UpdateField(obj, "Tag", uint64(0))
	err = db.UpdateField(&AccountObject{Id: 10}, "Tag", uint64(0))
	require.NoError(t, err)
}

func TestRemove(t *testing.T) {
	db, err := eosiodb.NewDatabase("./", "eos.db", true)
	require.NoError(t, err)
	defer db.Close()

	obj := &AccountObject{Id: 10, Name: "hello", Tag: 0}
	err = db.Remover(obj)
	require.NoError(t, err)

	var objs []AccountObject
	err = db.All(&objs)
	require.NoError(t, err)
	fmt.Println(len(objs))
	require.Equal(t, len(objs), 5)
}

func TestUnique(t *testing.T) {
	db, err := eosiodb.NewDatabase("./", "eos.db", true)
	require.NoError(t, err)
	defer db.Close()

	for key, value := range names {
		user := &User{Id: uint64(key + 11), Name: value, Tag: uint64(10)}
		err = db.Insert(user)
		require.NoError(t, err)
	}
	user := &User{Id: 1000, Name: "linx", Tag: uint64(12)}
	err = db.Insert(user)
	require.Error(t, err, string("already exists"))
}
