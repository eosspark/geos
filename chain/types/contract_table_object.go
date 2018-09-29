package types

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"github.com/eosspark/eos-go/log"
)

type IdType uint64
type KeyType uint64

type Object struct {
	TypeId uint16 `storm:"id,increment"`
}

type TableIdObject struct {
	ID    IdType `storm:"id,increment"`
	Code  common.AccountName
	Scope common.ScopeName
	Table common.TableName
	Payer common.AccountName
	Count uint32
}

func (u *TableIdObject) GetBillableSize() uint64 {
	return 44 + overhead_per_row_per_index_ram_bytes*2
}

type ByCodeScopeTable struct {
	Code  common.AccountName
	Scope common.ScopeName
	Table common.TableName
}

type TableIdMultiIndex struct {
	TableIdObject
	Id  IdType           `storm:"id,increment"`
	Bst ByCodeScopeTable `strom:"unique"`
}

type KeyValueObject struct {
	// KeyType      KeyType
	// NumberOfKeys int //default 1
	ID                IdType `storm:"id,increment"`
	TId               IdType
	PrimaryKey        uint64
	Payer             common.AccountName
	Value             common.HexBytes // c++ SharedString
	ByScopePrimaryKey struct {
		TID        IdType
		PrimaryKey uint64
	} `strom:"unique"`
}

func (u *KeyValueObject) GetBillableSize() uint64 {
	return 32 + 8 + 4 + overhead_per_row_per_index_ram_bytes*2
}

func AddTableIdObjectIndex(dbs *eosiodb.DataBase, tio TableIdObject) {
	ti := TableIdMultiIndex{}
	ti.TableIdObject = tio
	ti.Bst.Code = tio.Code
	ti.Bst.Scope = tio.Scope
	ti.Bst.Table = tio.Table
	err := dbs.Insert(&ti)
	if err != nil {
		log.Error("Insert is error detail:", err)
		return
	}
}

func GetTableObjectById(dbs *eosiodb.DataBase, id IdType) *TableIdMultiIndex {
	tmi := TableIdMultiIndex{}
	err := dbs.Find("ID", id, &tmi)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &tmi
}

func GetByCodeScopeTable(dbs *eosiodb.DataBase, bst ByCodeScopeTable) *TableIdMultiIndex {
	tmi := TableIdMultiIndex{}
	err := dbs.Find("Bst", bst, &tmi)
	if err != nil {
		log.Error("GetByCodeScopeTable is error", err)
		return nil
	}

	/*fmt.Println("*************************************")
	fmt.Println(tmi)
	fmt.Println("*************************************")*/
	return &tmi
}

const overhead_per_row_per_index_ram_bytes = 32

type SecondaryKeyInterface interface {
	//SetValue(value interface{})
	//GetValue() interface{}
	GetBillSize() int64
	Size() int64
}

type Uint64_t struct {
	Value uint64
}

//func (u *Uint64_t) SetValue(value interface{}) {
//	u.Value = value
//}
//func (u *Uint64_t) GetValue() interface{} {
//	return Value
//}
func (u *Uint64_t) GetBillSize() int64 {
	return 24 + 16 + overhead_per_row_per_index_ram_bytes*3
}
func (u *Uint64_t) Size() int64 {
	return 8
}

type Uint128_t struct {
	Value [16]byte
}

type Uint256_t struct {
	Value [32]byte
}

type Float64_t struct {
	Value float64
}

//func (u *Float64_t) SetValue(value interface{}) {
//	u.Value = value
//}
//func (u *Float64_t) GetValue() interface{} {
//	return Value
//}
func (u *Float64_t) GetBillSize() int64 {
	return 24 + 8 + overhead_per_row_per_index_ram_bytes*3
}
func (u *Float64_t) Size() int64 {
	return 8
}

type Float128_t struct {
	Value [16]byte
}

type SecondaryObject struct {
	ID           IdType `storm:"id,increment"`
	TId          IdType
	PrimaryKey   uint64
	Payer        common.AccountName
	SecondaryKey SecondaryKeyInterface
}

func (s *SecondaryObject) GetBillableSize() int64 {
	return s.SecondaryKey.GetBillSize()
}

// type SecondaryObjectI64 struct {
// 	ID           IdType `storm:"id,increment"`
// 	TId          IdType
// 	PrimaryKey   uint64
// 	Payer        common.AccountName
// 	SecondaryKey uint64
// }

// type SecondaryObjectI128 struct {
// 	ID           IdType `storm:"id,increment"`
// 	TId          IdType
// 	PrimaryKey   uint64
// 	Payer        common.AccountName
// 	SecondaryKey uint128_t
// }

// type SecondaryObjectI128 struct {
// 	ID           IdType `storm:"id,increment"`
// 	TId          IdType
// 	PrimaryKey   uint64
// 	Payer        common.AccountName
// 	SecondaryKey uint256_t
// }

// type SecondaryObjectDouble struct {
// 	ID           IdType `storm:"id,increment"`
// 	TId          IdType
// 	PrimaryKey   uint64
// 	Payer        common.AccountName
// 	SecondaryKey float64
// }

// type SecondaryObjectLongDouble struct {
// 	ID           IdType `storm:"id,increment"`
// 	TId          IdType
// 	PrimaryKey   uint64
// 	Payer        common.AccountName
// 	SecondaryKey float128_t
// }
