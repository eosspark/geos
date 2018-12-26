package database

import (
	"fmt"
	"reflect"
)

/////////////////////////////////////////////////////// test struct ///////////////////////////////////////////
type Carnivore struct {
	Lion  int `multiIndex:"orderedUnique,less"`
	Tiger int `multiIndex:"orderedUnique,less"`
}

type DbHouse struct {
	Id        uint64 `multiIndex:"id,increment"`
	Area      uint64 `multiIndex:"orderedUnique,less"`
	Name      string `multiIndex:"orderedUnique,less"`
	Carnivore Carnivore `multiIndex:"inline"`
}

type IdType int64
type Name uint64
type AccountName uint64
type PermissionName uint64
type ActionName uint64
type TableName uint64
type ScopeName uint64

type DbTableIdObject struct {
	ID    IdType      `multiIndex:"id,increment,byScope"`
	Code  AccountName `multiIndex:"orderedUnique,less"`
	Scope ScopeName   `multiIndex:"byTable,orderedUnique,less:byScope,orderedUnique,less"`
	Table TableName   `multiIndex:"byTable,orderedUnique,less"`
	Payer AccountName `multiIndex:"byScope,orderedUnique"`
	Count uint32
}

type DbResourceLimitsObject struct {
	ID        IdType      `multiIndex:"id,increment"`
	Pending   bool        `multiIndex:"byOwner,orderedUnique"`
	Owner     AccountName `multiIndex:"byOwner,orderedUnique"`
	NetWeight int64
	CpuWeight int64
	RamBytes  int64
}

func LogObj(data interface{}) {
	space := "	"
	ref := reflect.ValueOf(data)
	if !ref.IsValid() || reflect.Indirect(ref).Kind() != reflect.Struct {
		fmt.Println("log obj valid")
		return
	}

	s := &ref
	if s.Kind() == reflect.Ptr {
		e := s.Elem()
		s = &e
	}
	if s.Kind() != reflect.Struct {
		fmt.Println("log obj valid")
		return
	}
	typ := s.Type()

	num := s.NumField()
	for i := 0; i < num; i++ {
		v := s.Field(i)
		t := typ.Field(i)
		fmt.Print(t.Name, space, v, space)
	}
	fmt.Println("")
}
