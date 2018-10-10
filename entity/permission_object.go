package entity

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type PermissionUsageObject struct {
	ID                  common.IdType     `storm:"id"`
	LastUsed            common.TimePoint `json:"last_used"`
	ByAccountPermission common.Tuple     `storm:"index"`
}
type PermissionObject struct {
	ID          common.IdType `storm:"id,increment"`
	UsageId     common.IdType
	Parent      common.IdType
	Owner       common.AccountName
	Name        common.PermissionName
	LastUpdated common.TimePoint
	Auth        types.SharedAuthority
	/*ID、Parent*/
	ByParent common.Tuple `storm:"index"`
	/*Owner、name*/
	ByOwner common.Tuple `storm:"index"`
	/*Name、ID*/
	ByName common.Tuple `storm:"index"`
}

func (po PermissionObject) Satisfies(other PermissionObject) bool {
	return false
}
