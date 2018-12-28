package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/eos_math"
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

type Idx64Object struct {
	ID           common.IdType `multiIndex:"id,increment"`
	TId          common.IdType `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	SecondaryKey uint64        `multiIndex:"bySecondary,orderedUnique,less"`
	PrimaryKey   uint64        `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	Payer        common.AccountName
}

type Idx128Object struct {
	ID           common.IdType    `multiIndex:"id,increment"`
	TId          common.IdType    `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	SecondaryKey eos_math.Uint128 `multiIndex:"bySecondary,orderedUnique,less"`
	PrimaryKey   uint64           `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	Payer        common.AccountName
}

type Idx256Object struct {
	ID           common.IdType    `multiIndex:"id,increment"`
	TId          common.IdType    `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	SecondaryKey eos_math.Uint256 `multiIndex:"bySecondary,orderedUnique,less"`
	PrimaryKey   uint64           `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	Payer        common.AccountName
}

type IdxDoubleObject struct {
	ID           common.IdType    `multiIndex:"id,increment"`
	TId          common.IdType    `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	SecondaryKey eos_math.Float64 `multiIndex:"bySecondary,orderedUnique,less"`
	PrimaryKey   uint64           `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	Payer        common.AccountName
}

type IdxLongDoubleObject struct {
	ID           common.IdType     `multiIndex:"id,increment"`
	TId          common.IdType     `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	SecondaryKey eos_math.Float128 `multiIndex:"bySecondary,orderedUnique,less"`
	PrimaryKey   uint64            `multiIndex:"byPrimary,orderedUnique,less:bySecondary,orderedUnique,less"`
	Payer        common.AccountName
}
