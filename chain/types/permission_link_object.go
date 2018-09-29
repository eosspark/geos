package types

import "github.com/eosspark/eos-go/common"

type PermissionLinkObject struct {
	Id                 IdType  `storm:"id,increment"`
	Account            common.AccountName
	Code               common.AccountName
	MessageType        common.ActionName
	RequiredPermission common.PermissionName

	//

}
//func s(){
//	common.MakeTuple(a,c,b,d)
//}