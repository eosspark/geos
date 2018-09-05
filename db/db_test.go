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
	Id   uint64 `storm:"id"`
	Name string `storm:index`
	Tag  uint64
}

type User struct {
	Id   uint64 `storm:"id"`
	Name string `storm:"unique"`
	Tag  uint64 `storm:"index"`
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
