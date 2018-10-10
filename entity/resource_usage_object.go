package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain"
)

type ResourceUsageObject struct {
	ID common.IdType
	Owner common.AccountName
	NetUsage chain.UsageAccumulator   `json:"net_usage"`
	CpuUsage chain.UsageAccumulator   `json:"cpu_usage"`
	RamUsage uint64             `json:"ram_usage"`
}
