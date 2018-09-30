package types

import (
	"github.com/eosspark/eos-go/common"
)

type PermissionUsageObject struct {
	ID                  IdType           `storm:"id"`
	LastUsed            common.TimePoint `json:"last_used"`
	ByAccountPermission common.Tuple     `storm:"index"`
}
type PermissionObject struct {
	ID          IdType `storm:"id,increment"`
	UsageId     IdType
	Parent      IdType
	Owner       common.AccountName
	Name        common.PermissionName
	LastUpdated common.TimePoint
	Auth        SharedAuthority
	/*ID、Parent*/
	ByParent    common.Tuple `storm:"index"`
	/*Owner、name*/
	ByOwner     common.Tuple `storm:"index"`
	/*Name、ID*/
	ByName      common.Tuple `storm:"index"`
}

func (po PermissionObject) Satisfies(other PermissionObject) bool {
	return false
}
