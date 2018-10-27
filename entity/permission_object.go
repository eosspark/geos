package entity

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

// type PermissionObject struct {
// 	Parent      common.IdType         `multiIndex:"byParent,orderedNonUnique"`
// 	ID          common.IdType         `multiIndex:"id,increment,byParent,byName"`
// 	UsageId     common.IdType
// 	Owner       common.AccountName    `multiIndex:"byOwner,orderedNonUnique"`
// 	Name        common.PermissionName `multiIndex:"byOwner,orderedNonUnique:byName,orderedNonUnique"`
// 	LastUpdated common.TimePoint
// 	Auth        types.SharedAuthority
// }

type PermissionObject struct {
	Parent      common.IdType `multiIndex:"byParent,orderedUnique"`
	ID          common.IdType `multiIndex:"id,increment,byParent,byName"`
	UsageId     common.IdType
	Owner       common.AccountName    `multiIndex:"byOwner,orderedUnique"`
	Name        common.PermissionName `multiIndex:"byOwner,orderedUnique:byName,orderedUnique"`
	LastUpdated common.TimePoint
	Auth        types.SharedAuthority
}

func (po *PermissionObject) Satisfies(other PermissionObject) bool {
	if po.Owner != other.Owner {
		return false
	}
	if po.ID == other.ID || po.ID == other.Parent {
		return true
	}
	//TODO po.Parent
	return false
}
