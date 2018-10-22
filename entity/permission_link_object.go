package entity

import (
	"github.com/eosspark/eos-go/common"
)

type PermissionLinkObject struct {
	ID                 common.IdType          `multiIndex:"id,increment"`
	Account            common.AccountName     `multiIndex:"byAction,orderedNonUnique:byPermissionName,orderedNonUnique"`
	Code               common.AccountName	  `multiIndex:"byAction,orderedNonUnique:byPermissionName,orderedNonUnique"`
	MessageType        common.ActionName	  `multiIndex:"byAction,orderedNonUnique:byPermissionName,orderedNonUnique"`
	RequiredPermission common.PermissionName  `multiIndex:"byPermissionName,orderedNonUnique"`
}
