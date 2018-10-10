package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain"
)

type ResourceLimitsConfigObject struct {
	ID                           common.IdType          `storm:"id"`
	CpuLimitParameters           chain.ElasticLimitParameters `json:"cpu_limit_parameters"`
	NetLimitParameters           chain.ElasticLimitParameters `json:"net_limit_parameters"`
	AccountCpuUsageAverageWindow uint32                 `json:"account_cpu_usage_average_window"`
	AccountNetUsageAverageWindow uint32                 `json:"account_net_usage_average_window"`
}
