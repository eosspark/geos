package entity

import (
	"github.com/eosspark/eos-go/common"
)

type ResourceLimitsObject struct {
	ID        common.IdType `storm:"id,increment"`
	Owner     common.AccountName `storm:"unique,ByOwner"`
	Pending   bool               `storm:"unique,ByOwner"`
	NetWeight int64              `json:"net_weight"`
	CpuWeight int64              `json:"cpu_weight"`
	RamBytes  int64              `json:"ram_bytes"`
}

func NewResourceLimitsObject() *ResourceLimitsObject{
	rlo := ResourceLimitsObject{}
	rlo.Pending = false
	rlo.NetWeight = -1
	rlo.CpuWeight = -1
	rlo.RamBytes = -1
	return &rlo
}

type ExponentialMovingAverageAccumulator struct {
	LastOrdinal uint32 `json:"last_ordinal"`
	ValueEx     uint64 `json:"value_ex"`
	Consumed    uint64 `json:"consumed"`
}
