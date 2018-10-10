package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain/types"
)

type ResourceLimitsStateObject struct {
	ID                   common.IdType `multiIndex:"id,increment,byId"`
	AverageBlockNetUsage types.UsageAccumulator
	AverageBlockCpuUsage types.UsageAccumulator
	PendingNetUsage      uint64
	PendingCpuUsage      uint64
	TotalNetWeight       uint64
	TotalCpuWeight       uint64
	TotalRamBytes        uint64
	VirtualNetLimit      uint64
	VirtualCpuLimit      uint64
}