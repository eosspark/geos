package chain

import (
	"github.com/eosspark/eos-go/common"
	"testing"
	"fmt"
	"math"
	"github.com/stretchr/testify/assert"
)

func initialize() *ResourceLimitsManager{
	control := GetControllerInstance()
	rlm := control.ResourceLimits
	rlm.InitializeDatabase()
	return rlm
}

func expectedElasticIterations(from uint64, to uint64, rateNum uint64, rateDen uint64) uint64 {
	result := uint64(0)
	cur := from
	for (from < to && cur < to) || (from > to && cur > to){
		cur = cur * rateNum / rateDen
		result ++
	}
	return result
}

func expectedExponentialAverageIterations(from uint64, to uint64, value uint64, windowSize uint64) uint64 {
	result := uint64(0)
	cur := from
	for (from < to && cur < to) || (from > to && cur > to){
		cur = cur * (windowSize - 1) / windowSize
		cur += value / windowSize
		result ++
	}
	return result
}

func TestElasticCpuRelaxContract(t *testing.T){
	rlm := initialize()
	desiredVirtualLimit := uint64(common.DefaultConfig.MaxBlockCpuUsage) * 1000
	expectedRelaxIteration := expectedElasticIterations(uint64(common.DefaultConfig.MaxBlockCpuUsage), desiredVirtualLimit, 1000, 999)

	expectedContractIteration := expectedExponentialAverageIterations(0, common.EosPercent(uint64(common.DefaultConfig.MaxBlockCpuUsage), common.DefaultConfig.TargetBlockCpuUsagePct),
		uint64(common.DefaultConfig.MaxBlockCpuUsage), uint64(common.DefaultConfig.BlockCpuUsageAverageWindowMs)/uint64(common.DefaultConfig.BlockIntervalMs)) +
		expectedElasticIterations(desiredVirtualLimit, uint64(common.DefaultConfig.MaxBlockCpuUsage), 99 ,100) - 1

	account := common.AccountName(common.N("1"))
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, -1, -1, -1)
	rlm.ProcessAccountLimitUpdates()
	f := common.FlatSet{}
	f.Insert(&account)

	iterations := uint32(0)
	for rlm.GetVirtualBlockCpuLimit() < desiredVirtualLimit && uint64(iterations) <= expectedRelaxIteration {
		rlm.AddTransactionUsage(&f,0,0,iterations)
		rlm.ProcessBlockUsage(iterations)
		iterations++
	}

	assert.Equal(t, expectedRelaxIteration, uint64(iterations))
	assert.Equal(t, desiredVirtualLimit, rlm.GetVirtualBlockCpuLimit())

	for rlm.GetVirtualBlockCpuLimit() > uint64(common.DefaultConfig.MaxBlockCpuUsage) && uint64(iterations) <= expectedRelaxIteration + expectedContractIteration{
		rlm.AddTransactionUsage(&f,uint64(common.DefaultConfig.MaxBlockCpuUsage),0,iterations)
		rlm.ProcessBlockUsage(iterations)
		iterations++
	}

	assert.Equal(t, expectedRelaxIteration + expectedContractIteration, uint64(iterations))
	assert.Equal(t, uint64(common.DefaultConfig.MaxBlockCpuUsage), rlm.GetVirtualBlockCpuLimit())
}

func TestElasticNetRelaxContract(t *testing.T){
	rlm := initialize()
	desiredVirtualLimit := uint64(common.DefaultConfig.MaxBlockNetUsage) * 1000
	expectedRelaxIteration := expectedElasticIterations(uint64(common.DefaultConfig.MaxBlockNetUsage), desiredVirtualLimit, 1000, 999)

	expectedContractIteration := expectedExponentialAverageIterations(0, common.EosPercent(uint64(common.DefaultConfig.MaxBlockNetUsage), common.DefaultConfig.TargetBlockNetUsagePct),
		uint64(common.DefaultConfig.MaxBlockNetUsage), uint64(common.DefaultConfig.BlockSizeAverageWindowMs)/uint64(common.DefaultConfig.BlockIntervalMs)) +
		expectedElasticIterations(desiredVirtualLimit, uint64(common.DefaultConfig.MaxBlockNetUsage), 99 ,100) - 1

	account := common.AccountName(common.N("1"))
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, -1, -1, -1)
	rlm.ProcessAccountLimitUpdates()

	f := common.FlatSet{}
	f.Insert(&account)

	iterations := uint32(0)
	for rlm.GetVirtualBlockNetLimit() < desiredVirtualLimit && uint64(iterations) <= expectedRelaxIteration {
		rlm.AddTransactionUsage(&f,0,0,iterations)
		rlm.ProcessBlockUsage(iterations)
		iterations++
	}

	assert.Equal(t, expectedRelaxIteration, uint64(iterations))
	assert.Equal(t, desiredVirtualLimit, rlm.GetVirtualBlockNetLimit())

	for rlm.GetVirtualBlockNetLimit() > uint64(common.DefaultConfig.MaxBlockNetUsage) && uint64(iterations) <= expectedRelaxIteration + expectedContractIteration{
		rlm.AddTransactionUsage(&f,0,uint64(common.DefaultConfig.MaxBlockNetUsage),iterations)
		rlm.ProcessBlockUsage(iterations)
		iterations++
	}

	assert.Equal(t, expectedRelaxIteration + expectedContractIteration, uint64(iterations))
	assert.Equal(t, uint64(common.DefaultConfig.MaxBlockNetUsage), rlm.GetVirtualBlockNetLimit())
}

func TestWeightedCapacityCpu(t *testing.T){
	rlm := initialize()
	weights := []int64{234, 511, 672, 800, 1213}
	total := int64(0)
	for _, w := range weights {
		total += w
	}
	expectedLimits := make([]int64,5)
	for i, w := range weights{
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
		assert.Equal(t, expectedLimits[idx], rlm.GetAccountCpuLimit(account,true))
		f := common.FlatSet{}
		f.Insert(&account)
		//s := rlm.db.StartSession()
		//rlm.AddTransactionUsage(&f, uint64(expectedLimits[idx]),0,  0)
		//s.Undo()
		//
		////expect txCpuUsageExceededFailure
		//rlm.AddTransactionUsage(&f, uint64(expectedLimits[idx]),0,  0)

	}
}

func TestWeightedCapacityNet(t *testing.T){
	rlm := initialize()
	weights := []int64{234, 511, 672, 800, 1213}
	total := int64(0)
	for _, w := range weights {
		total += w
	}
	expectedLimits := make([]int64,5)
	for i, w := range weights{
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
		assert.Equal(t, expectedLimits[idx], rlm.GetAccountNetLimit(account,true))
		f := common.FlatSet{}
		f.Insert(&account)
		//s := rlm.db.StartSession()
		//rlm.AddTransactionUsage(&f, 0, uint64(expectedLimits[idx]), 0)
		//s.Undo()
		//
		////expect txNetUsageExceededFailure
		//rlm.AddTransactionUsage(&f,0,uint64(expectedLimits[idx]),0)
	}
}

func TestEnforceBlockLimitsCpu(t *testing.T){
	rlm := initialize()
	account := common.AccountName(1)
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, -1, -1, -1)
	rlm.ProcessAccountLimitUpdates()

	f := common.FlatSet{}
	f.Insert(&account)

	increment := uint64(1000)
	expectedIterations := uint64(common.DefaultConfig.MaxBlockCpuUsage) / increment

	for idx := 0; uint64(idx) < expectedIterations; idx++ {
		rlm.AddTransactionUsage(&f, increment, 0, 0)
	}
	//expect blockResourceExhausted
	rlm.AddTransactionUsage(&f, increment, 0, 0)
}

func TestEnforceBlockLimitsNet(t *testing.T){
	rlm := initialize()
	account := common.AccountName(1)
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, -1, -1, -1)
	rlm.ProcessAccountLimitUpdates()

	f := common.FlatSet{}
	f.Insert(&account)

	increment := uint64(1000)
	expectedIterations := uint64(common.DefaultConfig.MaxBlockNetUsage) / increment

	for idx := 0; uint64(idx) < expectedIterations; idx++ {
		rlm.AddTransactionUsage(&f,0, increment, 0)
	}
	//expect blockResourceExhausted
	rlm.AddTransactionUsage(&f,0, increment, 0)
}

func TestEnforceAccountRamLimit(t *testing.T){
	rlm := initialize()
	limit := uint64(1000)
	increment := uint64(77)
	expectedIterations := (limit + increment - 1) / increment

	account := common.AccountName(1)
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, int64(limit), -1, -1)
	rlm.ProcessAccountLimitUpdates()

	for idx := 0; uint64(idx) < expectedIterations - 1; idx++ {
		rlm.AddPendingRamUsage(account, int64(increment))
		rlm.VerifyAccountRamUsage(account)
	}
	rlm.AddPendingRamUsage(account, int64(increment))
	//throw ramUsageExceeded
	rlm.VerifyAccountRamUsage(account)
}

func TestEnforceAccountRamLimitUnderflow(t *testing.T){
	rlm := initialize()
	account := common.AccountName(1)
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, 100, -1, -1)
	rlm.VerifyAccountRamUsage(account)
	rlm.ProcessAccountLimitUpdates()
	//throw transactionException
	rlm.AddPendingRamUsage(account, -101)
}

func TestEnforceAccountRamLimitOverflow(t *testing.T){
	rlm := initialize()
	account := common.AccountName(1)
	rlm.InitializeAccount(account)
	//rlm.SetAccountLimits(account, math.MaxUint64, -1, -1)
	rlm.VerifyAccountRamUsage(account)
	rlm.AddPendingRamUsage(account, math.MaxUint64/2)
	rlm.VerifyAccountRamUsage(account)
	rlm.AddPendingRamUsage(account, math.MaxUint64/2)
	rlm.VerifyAccountRamUsage(account)
	//throw transactionException
	rlm.AddPendingRamUsage(account, 2)
}

func TestEnforceAccountRamCommitment(t *testing.T){
	rlm := initialize()
	limit := uint64(1000)
	commit := uint64(600)
	increment := uint64(77)
	expectedIterations := (limit - commit + increment - 1) / increment

	account := common.AccountName(1)
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, int64(limit),-1,-1)
	rlm.ProcessAccountLimitUpdates()
	rlm.AddPendingRamUsage(account, int64(commit))
	rlm.VerifyAccountRamUsage(account)

	for idx := 0; uint64(idx) < expectedIterations - 1; idx++ {
		rlm.SetAccountLimits(account, int64(limit - increment * uint64(idx)), -1, -1)
		rlm.VerifyAccountRamUsage(account)
		rlm.ProcessAccountLimitUpdates()
	}

	rlm.SetAccountLimits(account, int64(limit - increment * expectedIterations), -1, -1)
	//throw ramUsageExceeded
	rlm.VerifyAccountRamUsage(account)
}

func TestSanityCheck(t *testing.T){
	rlm := initialize()
	totalStakedTokens := uint64(10000000000000)
	userStake := uint64(10000)
	maxBlockCpu := uint64(100000)
	blocksPerDay := uint64(2*60*60*23)
	totalCpuPerPeriod := maxBlockCpu * blocksPerDay * 3

	congestedCpuTimePerPeriod := totalCpuPerPeriod * userStake / totalStakedTokens
	unCongestedCpuTimePerPeriod := (1000*totalCpuPerPeriod) * userStake / totalStakedTokens
	fmt.Println(congestedCpuTimePerPeriod)
	fmt.Println(unCongestedCpuTimePerPeriod)

	dan := common.AccountName(common.N("dan"))
	everyone := common.AccountName(common.N("everyone"))
	rlm.InitializeAccount(dan)
	rlm.InitializeAccount(everyone)
	rlm.SetAccountLimits(dan,0,0,10000)
	rlm.SetAccountLimits(everyone,0,0,10000000000000 - 10000)
	rlm.ProcessAccountLimitUpdates()
	f := common.FlatSet{}
	f.Insert(&dan)
	rlm.AddTransactionUsage(&f,10,0,1)
	fmt.Println(rlm.GetAccountCpuLimit(dan, true))
}

func TestResourceLimitsManager_UpdateAccountUsage(t *testing.T) {
	control := GetControllerInstance()
	rlm := control.ResourceLimits
	rlm.InitializeDatabase()
	a := common.AccountName(common.N("yuanchao"))
	f := common.FlatSet{}
	f.Insert(&a)
	rlm.InitializeAccount(a)
	rlm.AddTransactionUsage(&f, 100, 100, 1)
	rlm.UpdateAccountUsage(&f, 1)
	rlm.UpdateAccountUsage(&f, 86401)
	rlm.UpdateAccountUsage(&f, 172801)
	//结果value_ex应该为579/2 579/2/2
	rlm.UpdateAccountUsage(&f, 1)
	rlm.UpdateAccountUsage(&f, 172801)
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
	account := []common.AccountName{a,b,c}
	for _,acc := range account {
		rlm.InitializeAccount(acc)
	}
	rlm.SetAccountLimits(a, 100,100,100)
	rlm.SetAccountLimits(b, 200,300,100)

	rlm.ProcessAccountLimitUpdates()
}