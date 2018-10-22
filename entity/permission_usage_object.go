package entity

import "github.com/eosspark/eos-go/common"

type PermissionUsageObject struct {
	ID                  common.IdType    `multiIndex:"id,increment"`
	LastUsed            common.TimePoint
}
