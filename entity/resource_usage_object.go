package entity

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain/types"
)

type ResourceUsageObject struct {
	ID       common.IdType      `multiIndex:"id,increment,byId"`
	Owner    common.AccountName `multiIndex:"orderedNonUnique,byOwner"`
	NetUsage types.UsageAccumulator
	CpuUsage types.UsageAccumulator
	RamUsage uint64
}
