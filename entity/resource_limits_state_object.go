package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain"
)

type ResourceLimitsStateObject struct {
	ID                   common.IdType `storm:"id"`
	AverageBlockNetUsage chain.UsageAccumulator   `json:"average_block_net_usage"`
	AverageBlockCpuUsage chain.UsageAccumulator   `json:"average_block_cpu_usage"`
	PendingNetUsage      uint64             `json:"pending_net_usage"`
	PendingCpuUsage      uint64             `json:"pending_cpu_usage"`
	TotalNetWeight       uint64             `json:"total_net_weight"`
	TotalCpuWeight       uint64             `json:"total_cpu_weight"`
	TotalRamBytes        uint64             `json:"total_ram_bytes"`
	VirtualNetLimit      uint64             `json:"virtual_net_limit"`
	VirtualCpuLimit      uint64             `json:"virtual_cpu_limit"`
}

