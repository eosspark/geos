package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain"
)

type ResourceLimitsConfigObject struct {
	ID                           common.IdType `multiIndex:"id,increment,byId"`
	CpuLimitParameters           chain.ElasticLimitParameters
	NetLimitParameters           chain.ElasticLimitParameters
	AccountCpuUsageAverageWindow uint32
	AccountNetUsageAverageWindow uint32
}

func NewResourceLimitsConfigObject() *ResourceLimitsConfigObject{
	config := ResourceLimitsConfigObject{}
	config.CpuLimitParameters = chain.ElasticLimitParameters{Target:common.EosPercent(uint64(common.DefaultConfig.MaxBlockCpuUsage), common.DefaultConfig.TargetBlockCpuUsagePct),
		Max:uint64(common.DefaultConfig.MaxBlockCpuUsage),
		Periods:uint32(common.DefaultConfig.BlockCpuUsageAverageWindowMs) / uint32(common.DefaultConfig.BlockIntervalMs),
		MaxMultiplier:1000, ContractRate:chain.Ratio{Numerator:99, Denominator:100}, ExpandRate:chain.Ratio{Numerator:1000, Denominator:999},
	}

	config.NetLimitParameters = chain.ElasticLimitParameters{Target:common.EosPercent(uint64(common.DefaultConfig.MaxBlockNetUsage), common.DefaultConfig.TargetBlockNetUsagePct),
		Max:uint64(common.DefaultConfig.MaxBlockCpuUsage),
		Periods:uint32(common.DefaultConfig.BlockSizeAverageWindowMs) / uint32(common.DefaultConfig.BlockIntervalMs),
		MaxMultiplier:1000, ContractRate:chain.Ratio{Numerator:99, Denominator:100}, ExpandRate:chain.Ratio{Numerator:1000, Denominator:999},
	}
	config.AccountCpuUsageAverageWindow = common.DefaultConfig.AccountCpuUsageAverageWindowMs / uint32(common.DefaultConfig.BlockIntervalMs)
	config.AccountNetUsageAverageWindow = common.DefaultConfig.AccountNetUsageAverageWindowMs / uint32(common.DefaultConfig.BlockIntervalMs)
	return &config
}