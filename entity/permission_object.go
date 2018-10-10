package entity

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type PermissionObject struct {
	ID          common.IdType         `multiIndex:"id,increment,byId,byParent,byName"`
	UsageId     common.IdType
	Parent      common.IdType         `multiIndex:"id,increment,byParent"`
	Owner       common.AccountName    `multiIndex:"id,increment,byOwner,byName"`
	Name        common.PermissionName `multiIndex:"id,increment,byOwner"`
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
