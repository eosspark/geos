package entity

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type PermissionObject struct {
	Parent      common.IdType `multiIndex:"byParent,orderedUnique"`
	ID          common.IdType `multiIndex:"id,increment,byParent,byName"`
	UsageId     common.IdType
	Owner       common.AccountName    `multiIndex:"byOwner,orderedUnique"`
	Name        common.PermissionName `multiIndex:"byOwner,orderedUnique:byName,orderedUnique"`
	LastUpdated common.TimePoint
	Auth        types.SharedAuthority
}

func (po *PermissionObject) Satisfies(other PermissionObject/*, index *database.MultiIndex*/) bool {
	if po.Owner != other.Owner {
		return false
	}
	if po.ID == other.ID || po.ID == other.Parent {
		return true
	}
	//parent := permissionIndex.GetIndex("id")
	//for {
	//	if id == parent.parent{
	//		return true
	//	}
	//	if parent.parent.id == 0{
	//		return false
	//	}
	//	parent = permission.GetIndex("id").Find(parent.parent)
	//}
	return false
}
