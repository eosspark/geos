package entity

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type ResourceUsageObject struct {
	ID       common.IdType      `multiIndex:"id,increment"`
	Owner    common.AccountName `multiIndex:"byOwner,orderedUnique"`
	NetUsage types.UsageAccumulator
	CpuUsage types.UsageAccumulator
	RamUsage uint64
}
