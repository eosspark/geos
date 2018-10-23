package entity

import (
	"github.com/eosspark/eos-go/common"
)

type Object struct {
	TypeId uint16 `multiIndex:"id,increment"`
}

type TableIdObject struct {
	ID    common.IdType      `multiIndex:"id,increment"`
	Code  common.AccountName `multiIndex:"byCodeScopeTable,orderedUnique"`
	Scope common.ScopeName   `multiIndex:"byCodeScopeTable,orderedUnique"`
	Table common.TableName   `multiIndex:"byCodeScopeTable,orderedUnique"`
	Payer common.AccountName
	Count uint32
}

type KeyValueObject struct {
	ID         common.IdType `storm:"id,increment"`
	TId        common.IdType
	PrimaryKey uint64
	Payer      common.AccountName
	Value      common.HexBytes // c++ SharedString
}

// func (u *KeyValueObject) GetValue() IdType {
// 	return u.Value
// }
// func (u *KeyValueObject) GetTableId() IdType {
// 	return u.TId
// }

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
	ID           common.IdType `multiIndex:"id,increment"`
	TId          common.IdType `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	PrimaryKey   uint64        `multiIndex:"byName,orderedUnique,less:bySecondary,orderedUnique,less"`
	Payer        common.AccountName
	SecondaryKey Uint64_t `multiIndex:"bySecondary,orderedUnique"`
}

type SecondaryObjectDouble struct {
	ID           common.IdType `multiIndex:"id,increment"`
	TId          common.IdType `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	PrimaryKey   uint64        `multiIndex:"byName,orderedUnique,less:bySecondary,orderedUnique,less"`
	Payer        common.AccountName
	SecondaryKey Float64_t `multiIndex:"bySecondary,orderedUnique"`
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
