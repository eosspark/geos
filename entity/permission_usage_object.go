package entity

import "github.com/eosspark/eos-go/common"

type PermissionUsageObject struct {
	ID                  common.IdType     `storm:"id"`
	LastUsed            common.TimePoint `json:"last_used"`
	ByAccountPermission common.Tuple     `storm:"index"`
}
