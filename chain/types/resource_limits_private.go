package types

import (
	"github.com/eosspark/eos-go/common"
	"math"
	"math/big"
	"github.com/eosspark/eos-go/log"
)

type Ratio struct {
	Numerator   uint64 `json:"numerator"`
	Denominator uint64 `json:"denominator"`
}

type ElasticLimitParameters struct {
	Target        uint64 `json:"target"`
	Max           uint64 `json:"max"`
	Periods       uint32 `json:"periods"`
	MaxMultiplier uint32 `json:"max_multiplier"`
	ContractRate  Ratio  `json:"contract_rate"`
	ExpandRate    Ratio  `json:"expand_rate"`
}

type AccountResourceLimit struct {
	Used      int64 `json:"used"`
	Available int64 `json:"available"`
	Max       int64 `json:"max"`
}

func IntegerDivideCeil(num *big.Int, den *big.Int) *big.Int {
	result := new(big.Int).Div(num, den)

	if new(big.Int).Mod(num, den).Int64() > 0 {
		result = new(big.Int).Add(result, big.NewInt(1))
	}
	return result
}

func IntegerDivideCeilUint64(num uint64, den uint64) uint64 {
	if num%den > 0 {
		return num/den + 1
	} else {
		return num / den
	}
}

type ExponentialMovingAverageAccumulator struct {
	LastOrdinal uint32 `json:"last_ordinal"`
	ValueEx     uint64 `json:"value_ex"`
	Consumed    uint64 `json:"consumed"`
}

func makeRatio(numerator uint64, denominator uint64) Ratio {
	return Ratio{numerator, denominator}
}

func MultiWithRatio(value uint64, ratio Ratio) uint64 {
	//eos.Asset{ratio.Denominator != 0 , "Usage exceeds maximum value representable after extending for precision"}
	return value * ratio.Numerator / ratio.Denominator
}

func DowngradeCast(val *big.Int) int64 {
	max := big.NewInt(math.MaxInt64)
	min := big.NewInt(math.MinInt64)

	if val.Cmp(max) == 1 || val.Cmp(min) == -1 {
		log.Error("Usage exceeds maximum value representable after extending for precision")
	}
	return val.Int64()
}

func (ema *ExponentialMovingAverageAccumulator) Average() uint64 {
	return IntegerDivideCeilUint64(ema.ValueEx, uint64(common.DefaultConfig.RateLimitingPrecision))
}

func (ema *ExponentialMovingAverageAccumulator) add(units uint64, ordinal uint32, windowSize uint32) {
	valueExContrib := IntegerDivideCeilUint64(units*uint64(common.DefaultConfig.RateLimitingPrecision), uint64(windowSize))
	if ema.LastOrdinal != ordinal {
		if ema.LastOrdinal+windowSize > ordinal {
			delta := ordinal - ema.LastOrdinal
			decay := makeRatio(uint64(windowSize-delta), uint64(windowSize))
			ema.ValueEx = MultiWithRatio(ema.ValueEx, decay)
		} else {
			ema.ValueEx = 0
		}
		ema.LastOrdinal = ordinal
		ema.Consumed = ema.Average()
	}
	ema.Consumed += units
	ema.ValueEx += valueExContrib
}

type UsageAccumulator struct {
	ExponentialMovingAverageAccumulator
}

type ResourceLimitsObject struct {
	Rlo       RloIndex           `storm:"id"`
	ID        ResourceObjectType `storm:"index"`
	Owner     common.AccountName `storm:"index"`
	Pending   bool               `storm:"index"`
	NetWeight int64              `json:"net_weight"`
	CpuWeight int64              `json:"cpu_weight"`
	RamBytes  int64              `json:"ram_bytes"`
}

func (rlo *ResourceLimitsObject) Init() {
	rlo.NetWeight = -1
	rlo.CpuWeight = -1
	rlo.RamBytes = -1
}

type RloIndex struct {
	ID      ResourceObjectType `json:"id"`
	Owner   common.AccountName `json:"owner"`
	Pending bool               `json:"pending"`
}

type ResourceUsageObject struct {
	Ruo      RuoIndex           `storm:"id"`
	ID       ResourceObjectType `storm:"index"`
	Owner    common.AccountName `storm:"index"`
	NetUsage UsageAccumulator   `json:"net_usage"`
	CpuUsage UsageAccumulator   `json:"cpu_usage"`
	RamUsage uint64             `json:"ram_usage"`
}

type RuoIndex struct {
	ID    ResourceObjectType `json:"id"`
	Owner common.AccountName `json:"owner"`
}

type ResourceLimitsConfigObject struct {
	ID                           ResourceObjectType     `storm:"id"`
	CpuLimitParameters           ElasticLimitParameters `json:"cpu_limit_parameters"`
	NetLimitParameters           ElasticLimitParameters `json:"net_limit_parameters"`
	AccountCpuUsageAverageWindow uint32                 `json:"account_cpu_usage_average_window"`
	AccountNetUsageAverageWindow uint32                 `json:"account_net_usage_average_window"`
}

func (config *ResourceLimitsConfigObject) Init() {
	config.CpuLimitParameters = ElasticLimitParameters{common.EosPercent(uint64(common.DefaultConfig.MaxBlockCpuUsage), common.DefaultConfig.TargetBlockCpuUsagePct),
		uint64(common.DefaultConfig.MaxBlockCpuUsage),
		uint32(common.DefaultConfig.BlockCpuUsageAverageWindowMs) / uint32(common.DefaultConfig.BlockIntervalMs),
		1000, Ratio{99, 100}, Ratio{1000, 999},
	}

	config.NetLimitParameters = ElasticLimitParameters{common.EosPercent(uint64(common.DefaultConfig.MaxBlockNetUsage), common.DefaultConfig.TargetBlockNetUsagePct),
		uint64(common.DefaultConfig.MaxBlockCpuUsage),
		uint32(common.DefaultConfig.BlockSizeAverageWindowMs) / uint32(common.DefaultConfig.BlockIntervalMs),
		1000, Ratio{99, 100}, Ratio{1000, 999},
	}
	config.AccountCpuUsageAverageWindow = common.DefaultConfig.AccountCpuUsageAverageWindowMs / uint32(common.DefaultConfig.BlockIntervalMs)
	config.AccountNetUsageAverageWindow = common.DefaultConfig.AccountNetUsageAverageWindowMs / uint32(common.DefaultConfig.BlockIntervalMs)
}

type ResourceLimitsStateObject struct {
	ID                   ResourceObjectType `storm:"id"`
	AverageBlockNetUsage UsageAccumulator   `json:"average_block_net_usage"`
	AverageBlockCpuUsage UsageAccumulator   `json:"average_block_cpu_usage"`
	PendingNetUsage      uint64             `json:"pending_net_usage"`
	PendingCpuUsage      uint64             `json:"pending_cpu_usage"`
	TotalNetWeight       uint64             `json:"total_net_weight"`
	TotalCpuWeight       uint64             `json:"total_cpu_weight"`
	TotalRamBytes        uint64             `json:"total_ram_bytes"`
	VirtualNetLimit      uint64             `json:"virtual_net_limit"`
	VirtualCpuLimit      uint64             `json:"virtual_cpu_limit"`
}

type ResourceObjectType uint64

const (
	_ ResourceObjectType = iota
	ResourceLimits
	ResourceUsage
	ResourceLimitsConfig
	ResourceLimitsState
)
