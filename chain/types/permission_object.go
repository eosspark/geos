package types

import (
	"github.com/eosspark/eos-go/common"
	"time"
)

type PermissionUsageObject struct {
	Id       uint64        `storm:"id"`
	LastUsed time.Duration `json:"last_used"`
}
type PermissionObject struct {
	Id          uint64
	UsageId     uint64
	Parent      uint64
	Owner       common.AccountName
	Name        common.PermissionName
	LastUpdated time.Duration
	Auth        SharedAuthority
}

func (po PermissionObject) Satisfies(other PermissionObject) bool {
	return false
}
