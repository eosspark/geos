package entity

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type PermissionLinkObject struct {
	Id                 types.IdType `storm:"id,increment"`
	Account            common.AccountName
	Code               common.AccountName
	MessageType        common.ActionName
	RequiredPermission common.PermissionName
	/*Account、Code、MessageType*/
	ByActionName common.Tuple `storm:"index"`
	/*Account、RequiredPermission、Code、MessageType*/
	ByPermissionName common.Tuple `storm:"index"`
}
