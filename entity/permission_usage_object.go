package entity

import "github.com/eosspark/eos-go/common"

type PermissionUsageObject struct {
	ID                  common.IdType    `multiIndex:"id,increment,byId"`
	LastUsed            common.TimePoint `multiIndex:"orderedNonUnique"`
}
