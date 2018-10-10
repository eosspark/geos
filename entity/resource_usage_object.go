package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain"
)

type ResourceUsageObject struct {
	ID       common.IdType      `multiIndex:"id,increment,byId"`
	Owner    common.AccountName `multiIndex:"orderedNonUnique,byOwner"`
	NetUsage chain.UsageAccumulator
	CpuUsage chain.UsageAccumulator
	RamUsage uint64
}
