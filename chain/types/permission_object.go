package types

import (
	"github.com/eosspark/eos-go/common"
)

type PermissionUsageObject struct {
	ID                  IdType           `storm:"id"`
	LastUsed            common.TimePoint `json:"last_used"`
	ByAccountPermission common.Pair      `storm:"index"`
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
	ByParent    common.Pair `storm:"index"`
	/*Owner、name*/
	ByOwner     common.Pair `storm:"index"`
	/*Name、ID*/
	ByName      common.Pair `storm:"index"`
}

func (po PermissionObject) Satisfies(other PermissionObject) bool {
	return false
}
