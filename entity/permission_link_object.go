package entity

import (
	"github.com/eosspark/eos-go/common"
)

type PermissionLinkObject struct {
	ID                 common.IdType         `multiIndex:"id,increment"`
	Account            common.AccountName    `multiIndex:"byActionName,orderedUnique:byPermissionName,orderedUnique"`
	Code               common.AccountName    `multiIndex:"byActionName,orderedUnique:byPermissionName,orderedUnique"`
	MessageType        common.ActionName     `multiIndex:"byActionName,orderedUnique:byPermissionName,orderedUnique"`
	RequiredPermission common.PermissionName `multiIndex:"byPermissionName,orderedUnique"`
}
