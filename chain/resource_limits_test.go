package chain

import (
	"fmt"
	"github.com/eosspark/container/sets/treeset"
	"github.com/eosspark/eos-go/common"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/exception"
)

func initializeResource() *ResourceLimitsManager {
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
	rlm := initializeResource()
	desiredVirtualLimit := uint64(common.DefaultConfig.MaxBlockCpuUsage) * 1000
	expectedRelaxIteration := expectedElasticIterations(uint64(common.DefaultConfig.MaxBlockCpuUsage), desiredVirtualLimit, 1000, 999)

	expectedContractIteration := expectedExponentialAverageIterations(0, common.EosPercent(uint64(common.DefaultConfig.MaxBlockCpuUsage), common.DefaultConfig.TargetBlockCpuUsagePct),
		uint64(common.DefaultConfig.MaxBlockCpuUsage), uint64(common.DefaultConfig.BlockCpuUsageAverageWindowMs)/uint64(common.DefaultConfig.BlockIntervalMs)) +
		expectedElasticIterations(desiredVirtualLimit, uint64(common.DefaultConfig.MaxBlockCpuUsage), 99, 100) - 1

	account := common.AccountName(common.N("1"))
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, -1, -1, -1)
	rlm.ProcessAccountLimitUpdates()
	f := treeset.NewWith(common.CompareName)
	f.AddItem(account)

	iterations := uint32(0)
	for rlm.GetVirtualBlockCpuLimit() < desiredVirtualLimit && uint64(iterations) <= expectedRelaxIteration {
		rlm.AddTransactionUsage(f, 0, 0, iterations)
		rlm.ProcessBlockUsage(iterations)
		iterations++
	}

	assert.Equal(t, expectedRelaxIteration, uint64(iterations))
	assert.Equal(t, desiredVirtualLimit, rlm.GetVirtualBlockCpuLimit())

	for rlm.GetVirtualBlockCpuLimit() > uint64(common.DefaultConfig.MaxBlockCpuUsage) && uint64(iterations) <= expectedRelaxIteration+expectedContractIteration {
		rlm.AddTransactionUsage(f, uint64(common.DefaultConfig.MaxBlockCpuUsage), 0, iterations)
		rlm.ProcessBlockUsage(iterations)
		iterations++
	}

	assert.Equal(t, expectedRelaxIteration+expectedContractIteration, uint64(iterations))
	assert.Equal(t, uint64(common.DefaultConfig.MaxBlockCpuUsage), rlm.GetVirtualBlockCpuLimit())
}

func TestElasticNetRelaxContract(t *testing.T) {
	rlm := initializeResource()
	desiredVirtualLimit := uint64(common.DefaultConfig.MaxBlockNetUsage) * 1000
	expectedRelaxIteration := expectedElasticIterations(uint64(common.DefaultConfig.MaxBlockNetUsage), desiredVirtualLimit, 1000, 999)

	expectedContractIteration := expectedExponentialAverageIterations(0, common.EosPercent(uint64(common.DefaultConfig.MaxBlockNetUsage), common.DefaultConfig.TargetBlockNetUsagePct),
		uint64(common.DefaultConfig.MaxBlockNetUsage), uint64(common.DefaultConfig.BlockSizeAverageWindowMs)/uint64(common.DefaultConfig.BlockIntervalMs)) +
		expectedElasticIterations(desiredVirtualLimit, uint64(common.DefaultConfig.MaxBlockNetUsage), 99, 100) - 1

	account := common.AccountName(common.N("1"))
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, -1, -1, -1)
	rlm.ProcessAccountLimitUpdates()

	f := treeset.NewWith(common.CompareName)
	f.AddItem(account)

	iterations := uint32(0)
	for rlm.GetVirtualBlockNetLimit() < desiredVirtualLimit && uint64(iterations) <= expectedRelaxIteration {
		rlm.AddTransactionUsage(f, 0, 0, iterations)
		rlm.ProcessBlockUsage(iterations)
		iterations++
	}
	assert.Equal(t, expectedRelaxIteration, uint64(iterations))
	assert.Equal(t, desiredVirtualLimit, rlm.GetVirtualBlockNetLimit())

	for rlm.GetVirtualBlockNetLimit() > uint64(common.DefaultConfig.MaxBlockNetUsage) && uint64(iterations) <= expectedRelaxIteration+expectedContractIteration {
		rlm.AddTransactionUsage(f, 0, uint64(common.DefaultConfig.MaxBlockNetUsage), iterations)
		rlm.ProcessBlockUsage(iterations)
		iterations++
	}
	assert.Equal(t, expectedRelaxIteration+expectedContractIteration, uint64(iterations))
	assert.Equal(t, uint64(common.DefaultConfig.MaxBlockNetUsage), rlm.GetVirtualBlockNetLimit())
}

func TestWeightedCapacityCpu(t *testing.T) {
	rlm := initializeResource()
	weights := []int64{234, 511, 672, 800, 1213}
	total := int64(0)
	for _, w := range weights {
		total += w
	}
	expectedLimits := make([]int64, 5)
	for i, w := range weights {
		windowSize := int64(common.DefaultConfig.AccountCpuUsageAverageWindowMs) / int64(common.DefaultConfig.BlockIntervalMs)
		expectedLimits[i] = w * int64(common.DefaultConfig.MaxBlockCpuUsage) * windowSize / total
	}
	for idx := int(0); idx < len(weights); idx++ {
		account := common.AccountName(idx + 100)
		rlm.InitializeAccount(account)
		rlm.SetAccountLimits(account, -1, -1, weights[idx])
	}

	rlm.ProcessAccountLimitUpdates()

	for idx := int(0); idx < len(weights); idx++ {
		account := common.AccountName(idx + 100)
		assert.Equal(t, expectedLimits[idx], rlm.GetAccountCpuLimit(account, true))
		f := treeset.NewWith(common.CompareName)
		f.AddItem(account)
		//s := rlm.db.StartSession()
		//rlm.AddTransactionUsage(f, uint64(expectedLimits[idx]), 0, 0)
		//s.Undo()

		try.Try(func() {
			rlm.AddTransactionUsage(f, uint64(expectedLimits[idx]), 0, 0)
		}).Catch(func(e exception.TxCpuUsageExceeded) {
			fmt.Println(e)
		}).End()
	}
}

func TestWeightedCapacityNet(t *testing.T) {
	rlm := initializeResource()
	weights := []int64{234, 511, 672, 800, 1213}
	total := int64(0)
	for _, w := range weights {
		total += w
	}
	expectedLimits := make([]int64, 5)
	for i, w := range weights {
		windowSize := int64(common.DefaultConfig.AccountNetUsageAverageWindowMs) / int64(common.DefaultConfig.BlockIntervalMs)
		expectedLimits[i] = w * int64(common.DefaultConfig.MaxBlockNetUsage) * windowSize / total
	}
	for idx := int(0); idx < len(weights); idx++ {
		account := common.AccountName(idx + 100)
		rlm.InitializeAccount(account)
		rlm.SetAccountLimits(account, -1, weights[idx], -1)
	}

	rlm.ProcessAccountLimitUpdates()

	for idx := int(0); idx < len(weights); idx++ {
		account := common.AccountName(idx + 100)
		assert.Equal(t, expectedLimits[idx], rlm.GetAccountNetLimit(account, true))
		f := treeset.NewWith(common.NameComparator)
		f.AddItem(account)
		s := rlm.db.StartSession()
		rlm.AddTransactionUsage(f, 0, uint64(expectedLimits[idx]), 0)
		s.Undo()

		try.Try(func() {
			rlm.AddTransactionUsage(f, 0, uint64(expectedLimits[idx]), 0)
		}).Catch(func(e exception.TxCpuUsageExceeded) {
			fmt.Println(e)
		}).End()
	}
}

func TestEnforceBlockLimitsCpu(t *testing.T) {
	rlm := initializeResource()
	account := common.AccountName(1)
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, -1, -1, -1)
	rlm.ProcessAccountLimitUpdates()
	f := treeset.NewWith(common.CompareName)
	f.AddItem(account)

	increment := uint64(1000)
	expectedIterations := uint64(common.DefaultConfig.MaxBlockCpuUsage) / increment

	for idx := 0; uint64(idx) < expectedIterations; idx++ {
		rlm.AddTransactionUsage(f, increment, 0, 0)
	}

	try.Try(func() {
		rlm.AddTransactionUsage(f, increment, 0, 0)
	}).Catch(func(e exception.BlockResourceExhausted) {
		fmt.Println(e)
	}).End()
}

func TestEnforceBlockLimitsNet(t *testing.T) {
	rlm := initializeResource()
	account := common.AccountName(1)
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, -1, -1, -1)
	rlm.ProcessAccountLimitUpdates()

	f := treeset.NewWith(common.CompareName)
	f.AddItem(account)

	increment := uint64(1000)
	expectedIterations := uint64(common.DefaultConfig.MaxBlockNetUsage) / increment

	for idx := 0; uint64(idx) < expectedIterations; idx++ {
		rlm.AddTransactionUsage(f, 0, increment, 0)
	}

	try.Try(func() {
		rlm.AddTransactionUsage(f, 0, increment, 0)
	}).Catch(func(e exception.BlockResourceExhausted) {
		fmt.Println(e)
	}).End()
}

func TestEnforceAccountRamLimit(t *testing.T) {
	rlm := initializeResource()
	limit := uint64(1000)
	increment := uint64(77)
	expectedIterations := (limit + increment - 1) / increment

	account := common.AccountName(1)
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, int64(limit), -1, -1)
	rlm.ProcessAccountLimitUpdates()

	for idx := 0; uint64(idx) < expectedIterations-1; idx++ {
		rlm.AddPendingRamUsage(account, int64(increment))
		rlm.VerifyAccountRamUsage(account)
	}
	rlm.AddPendingRamUsage(account, int64(increment))

	try.Try(func() {
		rlm.VerifyAccountRamUsage(account)
	}).Catch(func(e exception.RamUsageExceeded) {
		fmt.Println(e)
	}).End()
}

func TestEnforceAccountRamLimitUnderflow(t *testing.T) {
	rlm := initializeResource()
	account := common.AccountName(1)
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, 100, -1, -1)
	rlm.VerifyAccountRamUsage(account)
	rlm.ProcessAccountLimitUpdates()

	try.Try(func() {
		rlm.AddPendingRamUsage(account, -101)
	}).Catch(func(e exception.TransactionException) {
		fmt.Println(e)
	}).End()
}

func TestEnforceAccountRamLimitOverflow(t *testing.T) {
	rlm := initializeResource()
	account := common.AccountName(1)
	rlm.InitializeAccount(account)
	//rlm.SetAccountLimits(account, math.MaxUint64, -1, -1)
	rlm.VerifyAccountRamUsage(account)
	rlm.AddPendingRamUsage(account, math.MaxUint64/2)
	rlm.VerifyAccountRamUsage(account)
	rlm.AddPendingRamUsage(account, math.MaxUint64/2)
	rlm.VerifyAccountRamUsage(account)

	try.Try(func() {
		rlm.AddPendingRamUsage(account, 2)
	}).Catch(func(e exception.TransactionException) {
		fmt.Println(e)
	}).End()
}

func TestEnforceAccountRamCommitment(t *testing.T) {
	rlm := initializeResource()
	limit := uint64(1000)
	commit := uint64(600)
	increment := uint64(77)
	expectedIterations := (limit - commit + increment - 1) / increment

	account := common.AccountName(1)
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, int64(limit), -1, -1)
	rlm.ProcessAccountLimitUpdates()
	rlm.AddPendingRamUsage(account, int64(commit))
	rlm.VerifyAccountRamUsage(account)

	for idx := 0; uint64(idx) < expectedIterations-1; idx++ {
		rlm.SetAccountLimits(account, int64(limit-increment*uint64(idx)), -1, -1)
		rlm.VerifyAccountRamUsage(account)
		rlm.ProcessAccountLimitUpdates()
	}

	rlm.SetAccountLimits(account, int64(limit-increment*expectedIterations), -1, -1)

	try.Try(func() {
		rlm.VerifyAccountRamUsage(account)
	}).Catch(func(e exception.RamUsageExceeded) {
		fmt.Println(e)
	}).End()
}

func TestSanityCheck(t *testing.T) {
	rlm := initializeResource()
	totalStakedTokens := uint64(10000000000000)
	userStake := uint64(10000)
	maxBlockCpu := uint64(100000)
	blocksPerDay := uint64(2 * 60 * 60 * 23)
	totalCpuPerPeriod := maxBlockCpu * blocksPerDay * 3

	congestedCpuTimePerPeriod := totalCpuPerPeriod * userStake / totalStakedTokens
	unCongestedCpuTimePerPeriod := (1000 * totalCpuPerPeriod) * userStake / totalStakedTokens
	fmt.Println(congestedCpuTimePerPeriod)
	fmt.Println(unCongestedCpuTimePerPeriod)

	dan := common.AccountName(common.N("dan"))
	everyone := common.AccountName(common.N("everyone"))
	rlm.InitializeAccount(dan)
	rlm.InitializeAccount(everyone)
	rlm.SetAccountLimits(dan, 0, 0, 10000)
	rlm.SetAccountLimits(everyone, 0, 0, 10000000000000-10000)
	rlm.ProcessAccountLimitUpdates()
	f := treeset.NewWith(common.CompareName)
	f.AddItem(dan)
	rlm.AddTransactionUsage(f, 10, 0, 1)
	fmt.Println(rlm.GetAccountCpuLimit(dan, true))
}