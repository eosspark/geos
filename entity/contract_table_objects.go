package entity

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/log"
)

type IdType uint16
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

func (u *TableIdObject) GetBillableSizeOverhead() uint64 {
	return overhead_per_row_per_index_ram_bytes * 2
}
func (u *TableIdObject) GetBillableSizeValue() uint64 {
	return 44 + u.GetBillableSizeOverhead()
}
func (s *TableIdObject) MakeTuple(values ...interface{}) *common.Tuple {
	//return s.SecondaryKey.MakeTupe(values...)
	tuple := common.Tuple{}
	return &tuple
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

// func (u *KeyValueObject) GetValue() IdType {
// 	return u.Value
// }
// func (u *KeyValueObject) GetTableId() IdType {
// 	return u.TId
// }
func (u *KeyValueObject) GetBillableSizeOverhead() uint64 {
	return overhead_per_row_per_index_ram_bytes * 2
}
func (u *KeyValueObject) GetBillableSizeValue() uint64 {
	return 32 + 8 + 4 + u.GetBillableSizeOverhead()
}
func (s *KeyValueObject) MakeTuple(values ...interface{}) *common.Tuple {
	//return s.SecondaryKey.MakeTupe(values...)
	tuple := common.Tuple{}
	return &tuple
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

// type SecondaryKeyInterface interface {
// 	SetValue(value interface{})
// 	GetValue() interface{}
// 	GetBillSize() int64
// 	Size() int64
// 	MakeTuple(values ...interface{}) *common.Tuple
// }

type Uint64_t struct {
	Value uint64
}

// func (u *Uint64_t) SetValue(value interface{}) {
// 	u.Value = reflect.ValueOf(value).Uint()
// }
// func (u *Uint64_t) GetValue() interface{} {
// 	return u.Value
// }
// func (u *Uint64_t) GetBillSize() int64 {
// 	return 24 + 16 + overhead_per_row_per_index_ram_bytes*3
// }
// func (u *Uint64_t) Size() int64 {
// 	return 8
// }
func (s *Uint64_t) MakeTuple(values ...interface{}) *common.Tuple {
	//return s.SecondaryKey.MakeTupe(values...)
	tuple := common.Tuple{}
	return &tuple
}

type Float64_t struct {
	Value float64
}

// func (u *Float64_t) SetValue(value interface{}) {
// 	u.Value = reflect.ValueOf(value).Float()
// }
// func (u *Float64_t) GetValue() interface{} {
// 	return Value
// }
// func (u *Float64_t) GetBillSize() int64 {
// 	return 24 + 8 + overhead_per_row_per_index_ram_bytes*3
// }
// func (u *Float64_t) Size() int64 {
// 	return 8
// }
func (s *Float64_t) MakeTuple(values ...interface{}) *common.Tuple {
	//return s.SecondaryKey.MakeTupe(values...)
	tuple := common.Tuple{}
	return &tuple
}

// type Uint128_t struct {
// 	Value [16]byte
// }
// type Uint256_t struct {
// 	Value [32]byte
// }
// type Float128_t struct {
// 	Value [16]byte
// }
// type SecondaryObjectInterface struct {
// 	GetTId() IdType
// 	SetTId(id IdType)

// 	GetPrimaryKey() uint64
// 	SetPrimaryKey(primaryKey uint64)

// 	GetPayer() common.AccountName
// 	SetPayer(payer common.AccountName)
// }

type SecondaryObjectI64 struct {
	ID           IdType `storm:"id,increment"`
	TId          IdType
	PrimaryKey   uint64
	Payer        common.AccountName
	SecondaryKey Uint64_t
}

func (s *SecondaryObjectI64) GetBillableSizeOverHead() int64 {
	return overhead_per_row_per_index_ram_bytes * 3
}

func (s *SecondaryObjectI64) GetBillableSizeValue() int64 {
	return 24 + 8 + s.GetBillableSizeOverHead()
}

func (s *SecondaryObjectI64) MakeTuple(values ...interface{}) *common.Tuple {
	return s.SecondaryKey.MakeTuple(values...)
}

type SecondaryObjectDouble struct {
	ID           IdType `storm:"id,increment"`
	TId          IdType
	PrimaryKey   uint64
	Payer        common.AccountName
	SecondaryKey Float64_t
}

func (s *SecondaryObjectDouble) GetBillableSizeOverHead() int64 {
	return overhead_per_row_per_index_ram_bytes * 3
}

func (s *SecondaryObjectDouble) GetBillableSizeValue() int64 {
	return 24 + 8 + s.GetBillableSizeOverHead()
}

func (s *SecondaryObjectDouble) MakeTuple(values ...interface{}) *common.Tuple {
	return s.SecondaryKey.MakeTuple(values...)
}

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

// type SecondaryObjectLongDouble struct {
// 	ID           IdType `storm:"id,increment"`
// 	TId          IdType
// 	PrimaryKey   uint64
// 	Payer        common.AccountName
// 	SecondaryKey float128_t
// }
const billableAlignment uint64 = 16

func BillableSizeV(value uint64) uint64 {
	return ((value + billableAlignment - 1) / billableAlignment) * billableAlignment
}
