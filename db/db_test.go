package eosiodb

import (
	"fmt"
	//"reflect"
	"testing"
)

// This test case only tests db, the session is not tested
// because they have the same usage, and subsequent test cases will be added to the session.

var names = []string{"linx", "garytone", "elk", "fox", "lion"}

type AccountObject struct {
	Id   uint64
	Name string
}

type User struct {
	Id   uint64        `storm:"id"`
	Name string        `storm:"unique"`
	Tag  uint64        `storm:"index"`
	Acc  AccountObject `storm:"unique"`
}

func TestNewBase(t *testing.T) {
	db, err := NewDataBase("./", "test.db", true)
	if err != nil {
		fmt.Println("NewDataBase failed")
		return
	}
	defer db.Close()
	fmt.Println("database successful")
}

type House struct {
	Id   uint64 `storm:"id,increment"`
	Area uint64
	Name string
}

func Test_Lower_Upper(t *testing.T) {
	db, err := NewDataBase("./", "test.db", true)
	if err != nil {
		fmt.Println("NewDataBase failed")
		return
	}
	defer db.Close()
	fmt.Println("database successful")

	for key, value := range names {
		house := House{Area: (uint64(key + 1)), Name: value}
		err = db.Insert(&house)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	/////////////////// Lower_Bound(///////////////////////////
	var houses []House
	err = db.All(&houses)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(houses)
	houses = nil
	err = db.LowerBound("Area", 3, &houses)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(houses)
	/////////////////// Upper_Bound(///////////////////////////
	houses = nil
	err = db.UpperBound("Area", 3, &houses)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(houses)
}

func TestInser(t *testing.T) {
	db, err := NewDataBase("./", "test.db", true)
	if err != nil {
		fmt.Println("NewDataBase failed")
		return
	}
	defer db.Close()
	fmt.Println("database successful")

	var user User
	user.Id = 10
	user.Name = "10"
	user.Tag = 10
	user.Acc.Id = 11
	user.Acc.Name = "11"

	err = db.Insert(&user)
	if err != nil {
		fmt.Println("Insert error : ", err.Error())
		return
	}
	var tmp User
	err = db.Find("Acc", AccountObject{11, "11"}, &tmp)
	if err != nil {
		fmt.Println("Insert error : ", err.Error())
		return
	}
	fmt.Println(tmp)
}
