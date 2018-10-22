package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
)

type PublicKeyHistoryObject struct {
	ID         common.IdType         `multiIndex:"id,increment,byPubKey,byAccountPermission"`
	PublicKey  ecc.PublicKey         `multiIndex:"byPubKey,orderedNonUnique"`                      //c++ publicKey+id unique
	Name       common.AccountName    `multiIndex:"byAccountPermission,orderedNonUnique"` //c++ ByAccountPermission+id unique
	Permission common.PermissionName `multiIndex:"byAccountPermission,orderedNonUnique"` //c++ ByAccountPermission+id unique
}
