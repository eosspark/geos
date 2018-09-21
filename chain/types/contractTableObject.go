package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"github.com/eosspark/eos-go/log"
	"fmt"
)

type IdType uint16
type KeyType uint64

type Object struct {
	TypeId uint16 `storm:"id,increment"`
}

type TableIdObject struct {
	//Object
	ID IdType		`storm:"id,increment"`
	Code   common.AccountName
	Scope  common.ScopeName
	Table  common.TableName
	Payer  common.AccountName
	Count  uint32
}

type ByCodeScopeTable  struct {
	Code  common.AccountName
	Scope common.ScopeName
	Table common.TableName
}

type TableIdMultiIndex struct {
	TableIdObject
	Id IdType `storm:"id,increment"`
	Bst ByCodeScopeTable `strom:"unique"`
	/*Bst struct{
		Code  common.AccountName
		Scope common.ScopeName
		Table common.TableName
	} `strom:"unique"`*/
}

type KeyValueObject struct {
	//Object
	KeyType	KeyType
	NumberOfKeys	int	//default 1
	ID IdType
	TId	IdType
	PrimaryKey	uint64
	Payer	common.AccountName
	Value *string// c++ SharedString TODO
}

type KeyValueIndex struct {
	KeyValueObject
	ID IdType	`strom:"id unique"`
	ByScopePrimaryKey struct{
		TID 	IdType
		PrimaryKey uint64
	} `strom:"unique"`
}



func AddTableIdObjectIndex(dbs *eosiodb.DataBase,tio TableIdObject){
	ti := TableIdMultiIndex{}
	ti.TableIdObject = tio
	ti.Bst.Code = tio.Code
	ti.Bst.Scope = tio.Scope
	ti.Bst.Table = tio.Table
	err := dbs.Insert(&ti)
	if err != nil{
		log.Error("Insert is error detail:",err)
		return
	}
}

func GetTableObjectById(dbs *eosiodb.DataBase,id IdType) (*TableIdMultiIndex){
	tmi := TableIdMultiIndex{}
	err := dbs.Find("ID",id, &tmi)
	if err != nil{
		fmt.Println(err.Error())
	}
	return &tmi
}

func GetByCodeScopeTable(dbs *eosiodb.DataBase,bst ByCodeScopeTable) (*TableIdMultiIndex){
	tmi := TableIdMultiIndex{}
	err := dbs.Find("Bst", bst, &tmi)
	if err !=nil{
		log.Error("GetByCodeScopeTable is error",err)
		return nil
	}

	/*fmt.Println("*************************************")
	fmt.Println(tmi)
	fmt.Println("*************************************")*/
	return &tmi
}

