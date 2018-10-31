package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"testing"
)

func initialize() *ResourceLimitsManager {
	control := GetControllerInstance()
	rlm := control.ResourceLimits
	rlm.InitializeDatabase()
	return rlm
}

func expectedElasticIterations(from uint64, to uint64, rateNum uint64, rateDen uint64) uint64 {
	result := uint64(0)
	cur := from
	for (from < to && cur < to) || (from > to && cur > to) {
		cur = cur * rateNum / rateDen
		result++
	}
	return result
}

func expectedExponentialAverageIterations(from uint64, to uint64, value uint64, windowSize uint64) uint64 {
	result := uint64(0)
	cur := from
	for (from < to && cur < to) || (from > to && cur > to) {
		cur = cur * (windowSize - 1) / windowSize
		cur += value / windowSize
		result++
	}
	return result
}

func TestElasticCpuRelaxContract(t *testing.T) {
	rlm := initialize()
	desiredVirtualLimit := uint64(common.DefaultConfig.DefaultMaxBlockCpuUsage) * 1000
	expectedRelaxIteration := expectedElasticIterations(uint64(common.DefaultConfig.DefaultMaxBlockCpuUsage), desiredVirtualLimit, 1000, 999)

	expectedContractIteration := expectedExponentialAverageIterations(0, common.EosPercent(uint64(common.DefaultConfig.DefaultMaxBlockCpuUsage), common.DefaultConfig.DefaultTargetBlockCpuUsagePct),
		uint64(common.DefaultConfig.DefaultMaxBlockCpuUsage), uint64(common.DefaultConfig.BlockCpuUsageAverageWindowMs)/uint64(common.DefaultConfig.BlockIntervalMs)) +
		expectedElasticIterations(desiredVirtualLimit, uint64(common.DefaultConfig.DefaultMaxBlockCpuUsage), 99, 100) - 1

	account := common.AccountName(common.N("1"))
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, -1, -1, -1)
	rlm.ProcessAccountLimitUpdates()

	iterations := uint32(0)
	for rlm.GetVirtualBlockCpuLimit() < desiredVirtualLimit && uint64(iterations) <= expectedRelaxIteration {
		an := []common.Element{&account}
		fs := common.FlatSet{an}
		rlm.AddTransactionUsage(&fs, 0, 0, iterations)
		rlm.ProcessBlockUsage(iterations)
		iterations++
	}

	for rlm.GetVirtualBlockCpuLimit() > uint64(common.DefaultConfig.DefaultMaxBlockCpuUsage) && uint64(iterations) <= expectedRelaxIteration+expectedContractIteration {
		an := []common.Element{&account}
		fs := common.FlatSet{an}
		rlm.AddTransactionUsage(&fs, uint64(common.DefaultConfig.DefaultMaxBlockCpuUsage), 0, iterations)
		rlm.ProcessBlockUsage(iterations)
		iterations++
	}

}

func TestResourceLimitsManager_UpdateAccountUsage(t *testing.T) {
	control := GetControllerInstance()
	rlm := control.ResourceLimits
	rlm.InitializeDatabase()
	a := common.AccountName(common.N("yuanchao"))
	account := []common.Element{&a}
	fs := common.FlatSet{account}
	rlm.InitializeAccount(a)
	rlm.AddTransactionUsage(&fs, 100, 100, 1)
	rlm.UpdateAccountUsage(&fs, 1)
	rlm.UpdateAccountUsage(&fs, 86401)
	rlm.UpdateAccountUsage(&fs, 172801)
	//结果value_ex应该为579/2 579/2/2
	rlm.UpdateAccountUsage(&fs, 1)
	rlm.UpdateAccountUsage(&fs, 172801)
	//结果value_ex为0
}

func TestResourceLimitsManager_SetAccountLimits(t *testing.T) {
	control := GetControllerInstance()
	rlm := control.ResourceLimits
	rlm.InitializeDatabase()
	fmt.Println(rlm.GetBlockCpuLimit())
	a := common.AccountName(common.N("yuanchao"))
	rlm.InitializeAccount(a)
	rlm.SetAccountLimits(a, 100, 100, 100)
	var r, n, c int64
	rlm.GetAccountLimits(a, &r, &n, &c)
	fmt.Println(r, n, c)
}

func TestResourceLimitsManager_ProcessBlockUsage(t *testing.T) {
	control := GetControllerInstance()
	rlm := control.ResourceLimits
	rlm.InitializeDatabase()
	a := common.AccountName(common.N("yuanchao"))
	b := common.AccountName(common.N("shengfeng"))
	c := common.AccountName(common.N("haonan"))
	account := []common.AccountName{a, b, c}
	for _, acc := range account {
		rlm.InitializeAccount(acc)
	}
	rlm.SetAccountLimits(a, 100, 100, 100)
	rlm.SetAccountLimits(b, 200, 300, 100)

	rlm.ProcessAccountLimitUpdates()
}
