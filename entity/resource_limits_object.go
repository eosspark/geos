package entity

import (
	"github.com/eosspark/eos-go/common"
)

type ResourceLimitsObject struct {
	ID        common.IdType      `multiIndex:"id,increment"`
	Pending   bool               `multiIndex:"byOwner,orderedUnique"`
	Owner     common.AccountName `multiIndex:"byOwner,orderedUnique"`
	NetWeight int64
	CpuWeight int64
	RamBytes  int64
}

func NewResourceLimitsObject() ResourceLimitsObject{
	rlo := ResourceLimitsObject{}
	rlo.Pending = false
	rlo.NetWeight = -1
	rlo.CpuWeight = -1
	rlo.RamBytes = -1
	return rlo
}