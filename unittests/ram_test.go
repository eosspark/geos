package unittests

import (
	"fmt"
	"github.com/docker/docker/pkg/testutil/assert"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	assert2 "github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestRamTests(t *testing.T) {
	e := initEosioSystemTester()
	initRequestBytes := 80000
	//increment_contract_bytes must be less than table_allocation_bytes for this test setup to work
	incrementContractBytes := 10000
	tableAllocationBytes := 12000
	e.BuyRamBytes(common.DefaultConfig.SystemAccountName, common.N("eosio"), 70000)
	e.ProduceBlocks(10, false)

	test1 := common.N("testram11111")
	test2 := common.N("testram22222")
	e.CreateAccountWithResources2(test1, common.N("eosio"),
		uint32(initRequestBytes+40))
	e.CreateAccountWithResources2(test2, common.N("eosio"),
		uint32(initRequestBytes+1190))
	e.ProduceBlocks(10, false)
	assert.Equal(t, e.Stake(common.N("eosio.stake"), common.N("testram11111"), CoreFromString("10.0000"), CoreFromString("10.0000")), e.Success())
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
			e.BuyRamBytes(common.N("eosio"), test1, uint32(incrementContractBytes))
			e.BuyRamBytes(common.N("eosio"), test2, uint32(incrementContractBytes))
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
			e.BuyRamBytes(common.N("eosio"), test1, uint32(incrementContractBytes))
			e.BuyRamBytes(common.N("eosio"), test2, uint32(incrementContractBytes))
		}).End()
		if skipLoop {
			break
		}
	}
	e.ProduceBlocks(10, false)

	e.SetCode(test2, code, nil)
	e.SetAbi(test2, abi, nil)
	e.ProduceBlocks(10, false)

	total := e.GetTotalStake(uint64(test1))
	fmt.Println(total)
	initBytes := uint64(total["ram_bytes"].(int64))

	rlm := e.Control.GetMutableResourceLimitsManager()
	initialRamUsage := rlm.GetAccountRamUsage(test1)

	moreRam := uint64(tableAllocationBytes) + initBytes - uint64(initRequestBytes)
	assert2.True(t, moreRam >= 0, "Underlying understanding changed, need to reduce size of init_request_bytes")
	log.Warn("init_bytes: %d, initial_ram_usage: %d, init_request_bytes: %d, more_ram: %d.", initBytes, initialRamUsage, initRequestBytes, moreRam)
	e.BuyRamBytes(common.N("eosio"), test1, uint32(moreRam))
	e.BuyRamBytes(common.N("eosio"), test2, uint32(moreRam))

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
	total = e.GetTotalStake(uint64(test1))
	ramBytes := uint64(total["ram_bytes"].(int64))
	log.Warn("ram_bytes: %d, ram_usage: %d, initial_ram_usage: %d, init_bytes: %d, ram_usage - initial_ram_usage: %d, init_bytes - ram_usage: %d.",
		ramBytes, ramUsage, initialRamUsage, initBytes, ramUsage-initialRamUsage, initBytes-uint64(ramUsage))

	// allocate just beyond the allocated bytes
	setentry = common.Variants{
		"payer": test1,
		"from":  1,
		"to":    10,
		"size":  1790,
	}
	try.Try(func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}).Catch(func(e exception.RamUsageExceeded) {
		fmt.Println("account testram11111 has insufficient ram", e.String())
	}).End()

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
	try.Try(func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}).Catch(func(e exception.RamUsageExceeded) {
		fmt.Println("account testram11111 has insufficient ram", e.String())
	}).End()
	e.ProduceBlocks(1, false)
	assert.Equal(t, rlm.GetAccountRamUsage(test1), ramUsage-1000)

	// verify the new entry's bytes minus the freed up bytes for existing entries still exceeds the allocation bytes limit
	setentry = common.Variants{
		"payer": test1,
		"from":  1,
		"to":    11,
		"size":  1760,
	}
	try.Try(func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}).Catch(func(e exception.RamUsageExceeded) {
		fmt.Println("account testram11111 has insufficient ram", e.String())
	}).End()
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
	try.Try(func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}).Catch(func(e exception.RamUsageExceeded) {
		fmt.Println("account testram11111 has insufficient ram", e.String())
	}).End()
	e.ProduceBlocks(1, false)

	// verify that the new entry is under the allocation bytes limit
	setentry = common.Variants{
		"payer": test1,
		"from":  12,
		"to":    12,
		"size":  1620,
	}
	try.Try(func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}).Catch(func(e exception.RamUsageExceeded) {
		fmt.Println("account testram11111 has insufficient ram", e.String())
	}).End()
	e.ProduceBlocks(1, false)

	// verify that anoth new entry will exceed the allocation bytes limit, to setup testing of new payer
	setentry = common.Variants{
		"payer": test1,
		"from":  13,
		"to":    13,
		"size":  1660,
	}
	try.Try(func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}).Catch(func(e exception.RamUsageExceeded) {
		fmt.Println("account testram11111 has insufficient ram", e.String())
	}).End()
	e.ProduceBlocks(1, false)

	// verify that the new entry is under the allocation bytes limit
	setentry = common.Variants{
		"payer": test1,
		"from":  12,
		"to":    12,
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

	// verify that another new entry that is too big will exceed the allocation bytes limit, to setup testing of new payer
	setentry = common.Variants{
		"payer": test1,
		"from":  13,
		"to":    13,
		"size":  1900,
	}
	try.Try(func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}).Catch(func(e exception.RamUsageExceeded) {
		fmt.Println("account testram11111 has insufficient ram", e.String())
	}).End()
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
		"payer": test1,
		"from":  12,
		"to":    21,
		"size":  1930,
	}
	try.Try(func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}).Catch(func(e exception.RamUsageExceeded) {
		fmt.Println("account testram11111 has insufficient ram", e.String())
	}).End()
	e.ProduceBlocks(1, false)

	// verify that new entries for testram22222 are under the allocation bytes limit
	setentry = common.Variants{
		"payer": test2,
		"from":  12,
		"to":    21,
		"size":  1910,
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

	// verify that new entry for testram22222 exceed the allocation bytes limit
	setentry = common.Variants{
		"payer": test2,
		"from":  22,
		"to":    22,
		"size":  1910,
	}
	try.Try(func() {
		e.PushAction2(
			&test1,
			&actSenName,
			test1,
			&setentry,
			e.DefaultExpirationDelta,
			0,
		)
	}).Catch(func(e exception.RamUsageExceeded) {
		fmt.Println("account testram22222 has insufficient ram", e.String())
	}).End()
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
	e.PushAction2(
		&test1,
		&actSenName,
		test1,
		&setentry,
		e.DefaultExpirationDelta,
		0,
	)
	e.ProduceBlocks(1, false)

	e.close()
}

func TestSimple(t *testing.T) {
	fmt.Printf("%d\n", '\'')
}
