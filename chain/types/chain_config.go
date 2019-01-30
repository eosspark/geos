package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
)

type ChainConfig struct {
	MaxBlockNetUsage               uint64
	TargetBlockNetUsagePct         uint32
	MaxTransactionNetUsage         uint32
	BasePerTransactionNetUsage     uint32
	NetUsageLeeway                 uint32
	ContextFreeDiscountNetUsageNum uint32
	ContextFreeDiscountNetUsageDen uint32

	MaxBlockCpuUsage       uint32
	TargetBlockCpuUsagePct uint32
	MaxTransactionCpuUsage uint32
	MinTransactionCpuUsage uint32

	MaxTrxLifetime              uint32
	DeferredTrxExpirationWindow uint32
	MaxTrxDelay                 uint32
	MaxInlineActionSize         uint32
	MaxInlineActionDepth        uint16
	MaxAuthorityDepth           uint16
}

func (c *ChainConfig) Validate() {
	try.EosAssert(c.TargetBlockNetUsagePct <= uint32(common.DefaultConfig.Percent_100), &exception.ActionValidateException{},
		"target block net usage percentage cannot exceed 100%")
	try.EosAssert(c.TargetBlockNetUsagePct >= uint32(common.DefaultConfig.Percent_1/10), &exception.ActionValidateException{},
		"target block net usage percentage must be at least 0.1%")
	try.EosAssert(c.TargetBlockCpuUsagePct <= uint32(common.DefaultConfig.Percent_100), &exception.ActionValidateException{},
		"target block cpu usage percentage cannot exceed 100%")
	try.EosAssert(c.TargetBlockCpuUsagePct >= uint32(common.DefaultConfig.Percent_1/10), &exception.ActionValidateException{},
		"target block cpu usage percentage must be at least 0.1%")

	try.EosAssert(uint64(c.MaxTransactionNetUsage) < common.DefaultConfig.MaxBlockNetUsage, &exception.ActionValidateException{},
		"max transaction net usage must be less than max block net usage")
	try.EosAssert(c.MaxTransactionCpuUsage < common.DefaultConfig.MaxBlockCpuUsage, &exception.ActionValidateException{},
		"max transaction cpu usage must be less than max block cpu usage")

	try.EosAssert(c.BasePerTransactionNetUsage < common.DefaultConfig.MaxTransactionNetUsage, &exception.ActionValidateException{},
		"base net usage per transaction must be less than the max transaction net usage")
	try.EosAssert((c.MaxTransactionNetUsage-common.DefaultConfig.BasePerTransactionNetUsage) >= common.DefaultConfig.MinNetUsageDeltaBetweenBaseAndMaxForTrx,
		&exception.ActionValidateException{},
		"max transaction net usage must be at least: %s bytes larger than base net usage per transaction",
		common.DefaultConfig.MinNetUsageDeltaBetweenBaseAndMaxForTrx)
	try.EosAssert(c.ContextFreeDiscountNetUsageDen > 0, &exception.ActionValidateException{},
		"net usage discount ratio for context free data cannot have a 0 denominator")
	try.EosAssert(c.ContextFreeDiscountNetUsageNum <= common.DefaultConfig.ContextFreeDiscountNetUsageDen, &exception.ActionValidateException{},
		"net usage discount ratio for context free data cannot exceed 1")

	try.EosAssert(c.MinTransactionCpuUsage <= common.DefaultConfig.MaxTransactionCpuUsage, &exception.ActionValidateException{},
		"min transaction cpu usage cannot exceed max transaction cpu usage")
	try.EosAssert(c.MaxTransactionCpuUsage < (common.DefaultConfig.MaxBlockCpuUsage-common.DefaultConfig.MinTransactionCpuUsage), &exception.ActionValidateException{},
		"max transaction cpu usage must be at less than the difference between the max block cpu usage and the min transaction cpu usage")

	try.EosAssert(1 <= c.MaxAuthorityDepth, &exception.ActionValidateException{},
		"max authority depth should be at least 1")
}

func (c ChainConfig) IsEmpty() bool {
	return c.MaxBlockNetUsage == 0 && c.TargetBlockNetUsagePct == 0 && c.MaxTransactionNetUsage == 0 &&
		c.BasePerTransactionNetUsage == 0 && c.NetUsageLeeway == 0 && c.ContextFreeDiscountNetUsageNum == 0 &&
		c.ContextFreeDiscountNetUsageDen == 0 && c.MaxBlockCpuUsage == 0 && c.TargetBlockCpuUsagePct == 0 &&
		c.MaxTransactionCpuUsage == 0 && c.MinTransactionCpuUsage == 0 && c.MaxTrxLifetime == 0 &&
		c.DeferredTrxExpirationWindow == 0 && c.MaxTrxDelay == 0 && c.MaxInlineActionSize == 0 &&
		c.MaxInlineActionDepth == 0 && c.MaxAuthorityDepth == 0
}
