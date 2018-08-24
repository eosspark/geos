package eosiodb

import (
	"fmt"
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
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()
}

func TestWriteOne(t *testing.T) {
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	obj := &AccountObject{Id: 10, Name: "linx", Tag: 99}
	err = db.Insert(obj)
	if err != nil {
		t.Error(err)
	}

	obj_ := &AccountObject{Id: 10, Name: "qieqie", Tag: 99}
	err = db.Insert(obj_)
	if err != nil {
		t.Error(err)
	}
}

func TestFind(t *testing.T) {
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	raw := AccountObject{Id: 10, Name: "qieqie", Tag: 99}
	account := AccountObject{Id: 10, Name: "garytone", Tag: 10}
	var obj AccountObject
	err = db.Find("Name", "qieqie", &obj)
	if err != nil {
		t.Error(err)
	}

	if obj != raw {
		fmt.Println("find error one")
	}
	if obj == account {
		fmt.Println("find error two")
	}

	err = db.Insert(&obj)
	if err != nil {
		t.Error(err)
	}
}

func TestInsertSome(t *testing.T) {
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	for key, value := range names {
		obj := &AccountObject{Id: uint64(key + 11), Name: value, Tag: uint64(10)}
		err = db.Insert(obj)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestGet(t *testing.T) {
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	var objs []AccountObject
	err = db.Get("Tag", 10, &objs)
	if err != nil {
		t.Error(err)
	}
	if len(objs) != 5 {
		fmt.Println("TestGet Failed")
		return
	}
}

func TestAll(t *testing.T) {
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	var objs []AccountObject
	err = db.All(&objs)
	if err != nil {
		t.Error(err)
	}
	if len(objs) != 6 {
		fmt.Println("TestAll Failed")
		return
	}
}

func TestUpdateItem(t *testing.T) {
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	var obj AccountObject
	err = db.Find("Name", "qieqie", &obj)
	if err != nil {
		t.Error(err)
	}
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
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateField(t *testing.T) {
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	//obj := &AccountObject{Id: 10, Name: "hello", Tag: 1000}
	//err = db.UpdateField(obj, "Tag", uint64(0))
	err = db.UpdateField(&AccountObject{Id: 10}, "Tag", uint64(0))
	if err != nil {
		t.Error(err)
	}
}

func TestRemove(t *testing.T) {
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	obj := &AccountObject{Id: 10, Name: "hello", Tag: 0}
	err = db.Remover(obj)
	if err != nil {
		t.Error(err)
	}

	var objs []AccountObject
	err = db.All(&objs)
	if err != nil {
		t.Error(err)
	}
	if len(objs) != 5 {
		fmt.Println("TestRemove Failed")
		return
	}
}

func TestUnique(t *testing.T) {
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	for key, value := range names {
		user := &User{Id: uint64(key + 11), Name: value, Tag: uint64(10)}
		err = db.Insert(user)
		if err != nil {
			t.Error(err)
		}
	}
	user := &User{Id: 1000, Name: "linx", Tag: uint64(12)}
	err = db.Insert(user)
	if err == nil {
		fmt.Println("TestUnique error")
		return
	}
	fmt.Println(err.Error())
}

///////////////////////////////////////////////////////// Inline Test ////////////////////////////////////////////////////////////////

type Base struct {
	Id  uint64 `storm:"id,increment"`
	Tag string `storm:"index"`
}

type Country struct {
	Base  Base   `storm:"inline,unique"`
	Name  string `storm:"unique"`
	users []User `storm:"inline"`
}

func TestInlineWriteOne(t *testing.T) {
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	country := &Country{
		Base: Base{
			Tag: "Big Country",
		},
		Name: "China",
	}
	err = db.Insert(country)
	if err != nil {
		t.Error(err)
	}

	err = db.Insert(country)
	if err != nil {
		t.Error(err)
	}
}

func TestInlineAll(t *testing.T) {
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	var countrys []Country
	err = db.All(&countrys)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(countrys)
}

func TestInlineFind(t *testing.T) {
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	var country Country
	err = db.Find("Base", Base{Id: 1, Tag: "Big Country"}, &country)
	if err != nil {
		t.Error(err)
	}
}

type UserNested struct {
	ID        int
	OneNested struct {
		A         string
		TwoNested struct {
			C int
			D int
		} `storm:"unique"`
	} `storm:"inline"`
}

func TestInlineNested(t *testing.T) {
	db, err := NewDatabase("./", "eos.db", true)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	var nested UserNested
	nested.ID = 10
	nested.OneNested.A = "A"
	nested.OneNested.TwoNested.C = 100
	nested.OneNested.TwoNested.D = 200

	err = db.Insert(&nested)
	if err != nil {
		t.Error(err)
	}

	err = db.Insert(&nested) // NOTE Repeat the same value of the same unique field, the database will ignore it and will not return failure
	if err != nil {
		t.Error(err)
	}

	nested.ID = 11
	var nesteds []UserNested
	err = db.All(&nesteds)
	if len(nesteds) != 1 {
		fmt.Println("TestInlineNested All Failed")
		return
	}

	nested.ID = 11
	err = db.Insert(&nested)
	if err == nil {
		fmt.Println("TestInlineNested Insert Failed")
	}
	fmt.Println(err.Error())
}
