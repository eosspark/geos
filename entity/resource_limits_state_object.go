package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain/types"
)

var DefaultResourceLimitsStateObject ResourceLimitsStateObject

func init() {
	DefaultResourceLimitsStateObject.ID = 1
}

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

func (state *ResourceLimitsStateObject) UpdateVirtualCpuLimit(cfg ResourceLimitsConfigObject) {
	state.VirtualCpuLimit = types.UpdateElasticLimit(state.VirtualCpuLimit, state.AverageBlockCpuUsage.Average(), cfg.CpuLimitParameters)
}

func (state *ResourceLimitsStateObject) UpdateVirtualNetLimit(cfg ResourceLimitsConfigObject) {
	state.VirtualNetLimit = types.UpdateElasticLimit(state.VirtualNetLimit, state.AverageBlockNetUsage.Average(), cfg.NetLimitParameters)
}