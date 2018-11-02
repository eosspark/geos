package entity

import (
	"github.com/eosspark/eos-go/common"
)

type Object struct {
	TypeId uint16 `multiIndex:"id,increment"`
}

type TableIdObject struct {
	ID    common.IdType      `multiIndex:"id,increment"`
	Code  common.AccountName `multiIndex:"byCodeScopeTable,orderedUnique,less"`
	Scope common.ScopeName   `multiIndex:"byCodeScopeTable,orderedUnique,less"`
	Table common.TableName   `multiIndex:"byCodeScopeTable,orderedUnique,less"`
	Payer common.AccountName
	Count uint32
}

type KeyValueObject struct {
	ID         common.IdType `multiIndex:"id,increment"`
	TId        common.IdType `multiIndex:"byScopePrimary,orderedUnique,less"`
	PrimaryKey uint64        `multiIndex:"byScopePrimary,orderedUnique,less"`
	Payer      common.AccountName
	Value      common.HexBytes // c++ SharedString
}

// type Uint64_t struct {
// 	Value uint64
// }

// type Float64_t struct {
// 	Value float64
// }

// type SecondaryObjectI64 struct {
// 	ID           common.IdType `multiIndex:"id,increment"`
// 	TId          common.IdType `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
// 	PrimaryKey   uint64        `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
// 	Payer        common.AccountName
// 	SecondaryKey uint64 `multiIndex:"bySecondary,orderedUnique"`
// }

type SecondaryObjectI64 struct {
	ID           common.IdType `multiIndex:"id,increment"`
	TId          common.IdType `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	PrimaryKey   uint64        `multiIndex:"byPrimary,orderedUnique"`
	Payer        common.AccountName
	SecondaryKey uint64 `multiIndex:"bySecondary,orderedUnique"`
}

// type SecondaryObjectDouble struct {
// 	ID           common.IdType `multiIndex:"id,increment"`
// 	TId          common.IdType `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
// 	PrimaryKey   uint64        `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
// 	Payer        common.AccountName
// 	SecondaryKey float64 `multiIndex:"bySecondary,orderedUnique"`
// }
type SecondaryObjectDouble struct {
	ID           common.IdType `multiIndex:"id,increment"`
	TId          common.IdType `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	PrimaryKey   uint64        `multiIndex:"byPrimary,orderedUnique"`
	Payer        common.AccountName
	SecondaryKey float64 `multiIndex:"bySecondary,orderedUnique"`
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
