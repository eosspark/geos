package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain/types"
)

type ResourceUsageObject struct {
	ID       common.IdType      `multiIndex:"id,increment"`
	Owner    common.AccountName `multiIndex:"byOwner,orderedNonUnique"`
	NetUsage types.UsageAccumulator
	CpuUsage types.UsageAccumulator
	RamUsage uint64
}
