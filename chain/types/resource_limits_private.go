package types

import (
	"github.com/eosspark/eos-go/common"
	"math"
	"math/big"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/exception"
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

func UpdateElasticLimit(currentLimit uint64, averageUsage uint64, params ElasticLimitParameters) uint64 {
	result := currentLimit
	if averageUsage > params.Target {
		result = result * params.ContractRate.Numerator / params.ContractRate.Denominator
	} else {
		result = result * params.ExpandRate.Numerator / params.ExpandRate.Denominator
	}
	return common.Min(common.Min(result, params.Max), uint64(params.Max*uint64(params.MaxMultiplier)))
}

func (elp ElasticLimitParameters) Validate() {
	EosAssert(elp.Periods > 0, &ResourceLimitException{}, "elastic limit parameter 'periods' cannot be zero")
	EosAssert(elp.ContractRate.Denominator > 0, &ResourceLimitException{}, "elastic limit parameter 'contract_rate' is not a well-defined ratio")
	EosAssert(elp.ExpandRate.Denominator > 0, &ResourceLimitException{}, "elastic limit parameter 'expand_rate' is not a well-defined ratio")
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

func (ema *ExponentialMovingAverageAccumulator) Add(units uint64, ordinal uint32, windowSize uint32) {
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
