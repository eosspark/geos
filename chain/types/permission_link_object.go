package types

import "github.com/eosspark/eos-go/common"

type PermissionLinkObject struct {
	Id                 uint64
	Account            common.AccountName
	Code               common.AccountName
	MessageType        common.ActionName
	RequiredPermission common.PermissionName
}
