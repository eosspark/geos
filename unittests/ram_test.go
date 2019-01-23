package unittests

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestRamTests(t *testing.T) {
	e := initEosioSystemTester()
	initRequestBytes := 80000
	//increment_contract_bytes must be less than table_allocation_bytes for this test setup to work
	incrementContractBytes := 10000
	tableAllocationBytes := 12000
	//TODO: asset, not in C++.
	//e.BuyRamBytes(common.DefaultConfig.SystemAccountName, common.N("eosio"), 70000)
	e.ProduceBlocks(10, false)

	e.CreateAccountWithResources2(test1, eosio, uint32(initRequestBytes+40))
	e.CreateAccountWithResources2(test2, eosio, uint32(initRequestBytes+1190))
	e.ProduceBlocks(10, false)
	assert.Equal(t, e.Success(), e.Stake(eosioStake, common.N("testram11111"), CoreFromString("10.0000"), CoreFromString("5.0000")))
	e.ProduceBlocks(10, false)

	//test_ram_limit
	wasmName := "test_contracts/test_ram_limit.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	abiName := "test_contracts/test_ram_limit.abi"
	abi, _ := ioutil.ReadFile(abiName)

	skipLoop := false
	for i := 0; i < 10; i++ {
		try.Try(func() {
			e.SetCode(test1, code, nil)
			skipLoop = true
		}).Catch(func(ex exception.RamUsageExceeded) {
			initRequestBytes += incrementContractBytes
			e.BuyRamBytes(eosio, test1, uint32(incrementContractBytes))
			e.BuyRamBytes(eosio, test2, uint32(incrementContractBytes))
		}).End()
		if skipLoop {
			break
		}
	}
	e.ProduceBlocks(10, false)

	skipLoop = false
	for i := 0; i < 10; i++ {
		try.Try(func() {
			e.SetAbi(test1, abi, nil)
			skipLoop = true
		}).Catch(func(ex exception.RamUsageExceeded) {
			initRequestBytes += incrementContractBytes
			e.BuyRamBytes(eosio, test1, uint32(incrementContractBytes))
			e.BuyRamBytes(eosio, test2, uint32(incrementContractBytes))
		}).End()
		if skipLoop {
			break
		}
	}
	e.ProduceBlocks(10, false)

	e.SetCode(test2, code, nil)
	e.SetAbi(test2, abi, nil)
	e.ProduceBlocks(10, false)

	total := e.GetTotalStake(test1)
	initBytes := total["ram_bytes"].(uint64)

	rlm := e.Control.GetMutableResourceLimitsManager()
	initialRamUsage := rlm.GetAccountRamUsage(test1)

	moreRam := uint64(tableAllocationBytes) + initBytes - uint64(initRequestBytes)
	assert.True(t, moreRam >= 0, "Underlying understanding changed, need to reduce size of init_request_bytes")
	log.Warn("init_bytes: %d, initial_ram_usage: %d, init_request_bytes: %d, more_ram: %d.", initBytes, initialRamUsage, initRequestBytes, moreRam)
	e.BuyRamBytes(eosio, test1, uint32(moreRam))
	e.BuyRamBytes(eosio, test2, uint32(moreRam))

	// allocate just under the allocated bytes
	actSenName := common.N("setentry")
	setentry := common.Variants{
		"payer": test1,
		"from":  1,
		"to":    10,
		"size":  1780,
	}

	e.PushAction2(
		&test1,
		&actSenName,
		test1,
		&setentry,
		e.DefaultExpirationDelta,
		0,
	)
	e.ProduceBlocks(1, false)
	ramUsage := rlm.GetAccountRamUsage(test1)
	total = e.GetTotalStake(test1)
	ramBytes := total["ram_bytes"].(uint64)
	log.Warn("ram_bytes: %d, ram_usage: %d, initial_ram_usage: %d, init_bytes: %d, ram_usage - initial_ram_usage: %d, init_bytes - ram_usage: %d.",
		ramBytes, ramUsage, initialRamUsage, initBytes, ramUsage-initialRamUsage, initBytes-uint64(ramUsage))

	// allocate just beyond the allocated bytes
	setentry = common.Variants{
		"payer": test1,
		"from":  1,
		"to":    10,
		"size":  1790,
	}

	stFunc := func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}
	CheckThrowExceptionAndMsg(t, &exception.RamUsageExceeded{}, "account testram11111 has insufficient ram", stFunc)
	ramUsage = rlm.GetAccountRamUsage(test1)

	e.ProduceBlocks(1, false)
	assert.Equal(t, rlm.GetAccountRamUsage(test1), ramUsage)

	// update the entries with smaller allocations so that we can verify space is freed and new allocations can be made
	setentry = common.Variants{
		"payer": test1,
		"from":  1,
		"to":    10,
		"size":  1680,
	}
	e.PushAction2(
		&test1,
		&actSenName,
		test1,
		&setentry,
		e.DefaultExpirationDelta,
		0,
	)
	e.ProduceBlocks(1, false)
	assert.Equal(t, rlm.GetAccountRamUsage(test1), ramUsage-1000)

	// verify the added entry is beyond the allocation bytes limit
	setentry = common.Variants{
		"payer": test1,
		"from":  1,
		"to":    11,
		"size":  1680,
	}
	stFunc = func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}
	CheckThrowExceptionAndMsg(t, &exception.RamUsageExceeded{}, "account testram11111 has insufficient ram", stFunc)
	e.ProduceBlocks(1, false)
	assert.Equal(t, rlm.GetAccountRamUsage(test1), ramUsage-1000)

	// verify the new entry's bytes minus the freed up bytes for existing entries still exceeds the allocation bytes limit
	setentry = common.Variants{
		"payer": test1,
		"from":  1,
		"to":    11,
		"size":  1760,
	}
	stFunc = func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}
	CheckThrowExceptionAndMsg(t, &exception.RamUsageExceeded{}, "account testram11111 has insufficient ram", stFunc)
	e.ProduceBlocks(1, false)
	assert.Equal(t, rlm.GetAccountRamUsage(test1), ramUsage-1000)

	// verify the new entry's bytes minus the freed up bytes for existing entries are under the allocation bytes limit
	setentry = common.Variants{
		"payer": test1,
		"from":  1,
		"to":    11,
		"size":  1600,
	}
	e.PushAction2(
		&test1,
		&actSenName,
		test1,
		&setentry,
		e.DefaultExpirationDelta,
		0,
	)
	e.ProduceBlocks(1, false)

	actRmenName := common.N("rmentry")
	rmentry := common.Variants{
		"from": 3,
		"to":   3,
	}
	e.PushAction2(
		&test1,
		&actRmenName,
		test1,
		&rmentry,
		e.DefaultExpirationDelta,
		0,
	)

	// verify that the new entry will exceed the allocation bytes limit
	setentry = common.Variants{
		"payer": test1,
		"from":  12,
		"to":    12,
		"size":  1780,
	}
	stFunc = func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}
	CheckThrowExceptionAndMsg(t, &exception.RamUsageExceeded{}, "account testram11111 has insufficient ram", stFunc)

	e.ProduceBlocks(1, false)

	// verify that the new entry is under the allocation bytes limit
	setentry = common.Variants{
		"payer": test1,
		"from":  12,
		"to":    12,
		"size":  1620,
	}
	e.PushAction2(
		&test1,
		&actSenName,
		test1,
		&setentry,
		e.DefaultExpirationDelta,
		0,
	)
	e.ProduceBlocks(1, false)

	// verify that anoth new entry will exceed the allocation bytes limit, to setup testing of new payer
	setentry = common.Variants{
		"payer": test1,
		"from":  13,
		"to":    13,
		"size":  1660,
	}
	stFunc = func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}
	CheckThrowExceptionAndMsg(t, &exception.RamUsageExceeded{}, "account testram11111 has insufficient ram", stFunc)
	e.ProduceBlocks(1, false)

	// verify that the new entry is under the allocation bytes limit
	setentry = common.Variants{
		"payer": test2,
		"from":  12,
		"to":    12,
		"size":  1720,
	}
	e.PushAction3(
		&test1,
		&actSenName,
		[]*common.AccountName{&test1, &test2},
		&setentry,
		e.DefaultExpirationDelta,
		0,
	)
	e.ProduceBlocks(1, false)

	// verify that another new entry that is too big will exceed the allocation bytes limit, to setup testing of new payer
	setentry = common.Variants{
		"payer": test1,
		"from":  13,
		"to":    13,
		"size":  1900,
	}
	stFunc = func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}
	CheckThrowExceptionAndMsg(t, &exception.RamUsageExceeded{}, "account testram11111 has insufficient ram", stFunc)
	e.ProduceBlocks(1, false)

	// verify that the new entry is under the allocation bytes limit, because entry 12 is now charged to testram22222
	setentry = common.Variants{
		"payer": test1,
		"from":  13,
		"to":    13,
		"size":  1720,
	}
	e.PushAction2(
		&test1,
		&actSenName,
		test1,
		&setentry,
		e.DefaultExpirationDelta,
		0,
	)
	e.ProduceBlocks(1, false)

	// verify that new entries for testram22222 exceed the allocation bytes limit
	setentry = common.Variants{
		"payer": test2,
		"from":  12,
		"to":    21,
		"size":  1930,
	}
	stFunc = func() {
		e.PushAction3(
			&test1,
			&actSenName,
			[]*common.AccountName{&test1, &test2},
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}
	CheckThrowExceptionAndMsg(t, &exception.RamUsageExceeded{}, "account testram22222 has insufficient ram", stFunc)
	e.ProduceBlocks(1, false)

	// verify that new entries for testram22222 are under the allocation bytes limit
	setentry = common.Variants{
		"payer": test2,
		"from":  12,
		"to":    21,
		"size":  1910,
	}
	e.PushAction3(
		&test1,
		&actSenName,
		[]*common.AccountName{&test1, &test2},
		&setentry,
		e.DefaultExpirationDelta,
		0,
	)
	e.ProduceBlocks(1, false)

	// verify that new entry for testram22222 exceed the allocation bytes limit
	setentry = common.Variants{
		"payer": test2,
		"from":  22,
		"to":    22,
		"size":  1910,
	}
	stFunc = func() {
		e.PushAction3(
			&test1,
			&actSenName,
			[]*common.AccountName{&test1, &test2},
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}
	CheckThrowExceptionAndMsg(t, &exception.RamUsageExceeded{}, "account testram22222 has insufficient ram", stFunc)
	e.ProduceBlocks(1, false)
	rmentry = common.Variants{
		"from": 20,
		"to":   20,
	}
	e.PushAction2(
		&test1,
		&actRmenName,
		test1,
		&rmentry,
		e.DefaultExpirationDelta,
		0,
	)
	e.ProduceBlocks(1, false)

	// verify that new entry for testram22222 are under the allocation bytes limit
	setentry = common.Variants{
		"payer": test2,
		"from":  22,
		"to":    22,
		"size":  1910,
	}
	e.PushAction3(
		&test1,
		&actSenName,
		[]*common.AccountName{&test1, &test2},
		&setentry,
		e.DefaultExpirationDelta,
		0,
	)
	e.ProduceBlocks(1, false)

	e.close()
}
