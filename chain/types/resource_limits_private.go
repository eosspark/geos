package types

import (
	"math"
	"fmt"
	"github.com/eosspark/eos-go/chain/config"
	"github.com/eosspark/eos-go/common"
)

func IntegerDivideCeil(num uint64, den uint64) uint64 {
	if num % den > 0 {
		return num / den + 1
	} else {
		return num / den
	}
}

type ExponentialMovingAverageAccumulator struct {
	LastOrdinal uint32 `json:"last_ordinal"`
	ValueEx     uint64 `json:"value_ex"`
	Consumed    uint64 `json:"consumed"`
}

func makeRatio(numerator uint64, denominator uint64) Ratio{
	return Ratio{numerator, denominator}
}

func MultiWithRatio(value uint64, ratio Ratio) uint64{
	//eos.Asset{ratio.Denominator != 0 , "Usage exceeds maximum value representable after extending for precision"}
	return value * ratio.Numerator / ratio.Denominator
}

func DowngradeCast(val uint64) int64{
	var max, min uint64
	max = math.MaxInt64
	min = 0
	//min = math.MinInt64
	if val > max || val < min {
		fmt.Println("error")
	}
	return int64(val)
}

func (ema *ExponentialMovingAverageAccumulator) Average() uint64{
	return IntegerDivideCeil(ema.ValueEx, uint64(config.RateLimitingPrecision))
}

func (ema *ExponentialMovingAverageAccumulator) add(units uint64, ordinal uint32, windowSize uint32){
	valueExContrib := IntegerDivideCeil(units * uint64(config.RateLimitingPrecision), uint64(windowSize))
	if ema.LastOrdinal != ordinal {

		if ema.LastOrdinal + windowSize > ordinal {
			delta := ordinal - ema.LastOrdinal
			decay := makeRatio(uint64(windowSize - delta), uint64(windowSize))
			ema.ValueEx = MultiWithRatio(ema.ValueEx, decay)
		} else {
			ema.ValueEx = 0
		}
		ema.LastOrdinal = ordinal
		ema.Consumed = ema.Average()
	}
	ema.Consumed += units
	ema.Consumed += valueExContrib
}

type UsageAccumulator struct{
	ExponentialMovingAverageAccumulator
}

type ResourceLimitsObject struct {
	Rlo		  RloIndex           `storm:"id"`
	Id        ResourceObjectType `storm:"index"`
	Owner     common.AccountName `storm:"index"`
	Pending   bool               `storm:"index"`
	NetWeight int64              `json:"net_weight"`
	CpuWeight int64              `json:"cpu_weight"`
	RamBytes  int64              `json:"ram_bytes"`
}

type RloIndex struct {
	Id        ResourceObjectType `json:"id"`
	Owner     common.AccountName `json:"owner"`
	Pending   bool               `json:"pending"`
}

type ResourceUsageObject struct {
	Ruo		 RuoIndex           `storm:"id"`
	Id       ResourceObjectType `storm:"index"`
	Owner    common.AccountName `storm:"index"`
	NetUsage UsageAccumulator   `json:"net_usage"`
	CpuUsage UsageAccumulator   `json:"cpu_usage"`
	RamUsage uint64             `json:"ram_usage"`
}

type RuoIndex struct {
	Id       ResourceObjectType `json:"id"`
	Owner    common.AccountName `json:"owner"`
}

type ResourceLimitsConfigObject struct {
	Id                           ResourceObjectType     `storm:"id"`
	CpuLimitParameters           ElasticLimitParameters `json:"cpu_limit_parameters"`
	NetLimitParameters           ElasticLimitParameters `json:"net_limit_parameters"`
	AccountCpuUsageAverageWindow uint32                 `json:"account_cpu_usage_average_window"`
	AccountNetUsageAverageWindow uint32                 `json:"account_net_usage_average_window"`
}

type ResourceLimitsStateObject struct {
	Id                   ResourceObjectType `storm:"id"`
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

const(
	_ ResourceObjectType = iota
	ResourceLimits
	ResourceUsage
	ResourceLimitsConfig
	ResourceLimitsState
)
