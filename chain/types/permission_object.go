package types

import (
	"github.com/eosspark/eos-go/common"
)

type PermissionUsageObject struct {
	ID       uint64           `storm:"id"`
	LastUsed common.TimePoint `json:"last_used"`
}
type PermissionObject struct {
	ID          uint64
	UsageId     uint64
	Parent      uint64
	Owner       common.AccountName
	Name        common.PermissionName
	LastUpdated common.TimePoint
	Auth        SharedAuthority
}

func (po PermissionObject) Satisfies(other PermissionObject) bool {
	return false
}
