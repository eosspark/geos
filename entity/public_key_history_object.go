package entity

import "github.com/eosspark/eos-go/common"

type PublicKeyHistoryObject struct {
	ID         common.IdType         `storm:"id,increment"`
	PublicKey  common.PublicKeyType  `storm:"index"` //c++ publicKey+id unique
	Name       common.AccountName    `storm:"unique,ByAccountPermission"`	//c++ ByAccountPermission+id unique
	Permission common.PermissionName `storm:"unique,ByAccountPermission"`	//c++ ByAccountPermission+id unique
}
