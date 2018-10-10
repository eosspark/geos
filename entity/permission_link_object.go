package entity

import (
	"github.com/eosspark/eos-go/common"
)

type PermissionLinkObject struct {
	ID                 common.IdType          `multiIndex:"id,increment,byId"`
	Account            common.AccountName     `multiIndex:"orderedNonUnique,byActionName,byPermissionName"`
	Code               common.AccountName	  `multiIndex:"orderedNonUnique,byActionName,byPermissionName"`
	MessageType        common.ActionName	  `multiIndex:"orderedNonUnique,byActionName,byPermissionName"`
	RequiredPermission common.PermissionName  `multiIndex:"orderedNonUnique,byPermissionName"`
}
