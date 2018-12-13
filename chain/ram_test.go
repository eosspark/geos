package chain

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

func TestRamTests(t *testing.T){
	e := initEosioSystemTester()
	initRequestBytes := 80000
	//increment_contract_bytes must be less than table_allocation_bytes for this test setup to work
	incrementContractBytes := 10000
	tableAllocationBytes := 12000
	e.BuyRamBytes(common.DefaultConfig.SystemAccountName,common.N("eosio"),7000)
	e.ProduceBlocks(10,false)

	test1 := common.N("testram11111")
	test2 := common.N("testram22222")
	e.CreateAccountWithResources(test1,common.N("eosio"),
		common.Asset{Amount:int64(initRequestBytes + 40)},false, CoreFromString("10.0000"),CoreFromString("10.0000"))
	e.CreateAccountWithResources(test2,common.N("eosio"),
		common.Asset{Amount:int64(initRequestBytes + 1190)},false, CoreFromString("10.0000"),CoreFromString("10.0000"))
	e.ProduceBlocks(10,false)
	assert.Equal(t, e.Stake(common.N("eosio.stake"),common.N("testram11111"),CoreFromString("10.0000"),CoreFromString("10.0000")), e.Success())
	e.ProduceBlocks(10,false)

	//test_ram_limit
	wasmName := "../wasmgo/testdata_context/test_ram_limit.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	abiName := "../wasmgo/testdata_context/test_ram_limit.abi"
	abi, _ := ioutil.ReadFile(abiName)

	skipLoop := false
	for i := 0; i < 10; i++ {
		try.Try(func() {
			e.SetCode(test1, code, nil)
			skipLoop = true
		}).Catch(func(ex exception.RamUsageExceeded) {
			initRequestBytes += incrementContractBytes
			e.BuyRamBytes(common.N("eosio"),test1,uint32(incrementContractBytes))
			e.BuyRamBytes(common.N("eosio"),test2,uint32(incrementContractBytes))
		}).End()
		if skipLoop {
			break
		}
	}
	e.ProduceBlocks(10,false)

	skipLoop = false
	for i := 0; i < 10; i++ {
		try.Try(func() {
			e.SetAbi(test1, abi, nil)
			skipLoop = true
		}).Catch(func(ex exception.RamUsageExceeded) {
			initRequestBytes += incrementContractBytes
			e.BuyRamBytes(common.N("eosio"),test1,uint32(incrementContractBytes))
			e.BuyRamBytes(common.N("eosio"),test2,uint32(incrementContractBytes))
		}).End()
		if skipLoop {
			break
		}
	}
	e.ProduceBlocks(10,false)

	e.SetCode(test2, code, nil)
	e.SetAbi(test2, abi, nil)
	e.ProduceBlocks(10,false)

	total := e.GetTotalStake(test1)
	initBytes := total["ram_bytes"].(uint64)

	rlm := e.Control.GetMutableResourceLimitsManager()
	initialRamUsage := rlm.GetAccountRamUsage(test1)

	moreRam := uint64(tableAllocationBytes) + initBytes - uint64(initRequestBytes)
	assert2.True(t, moreRam >= 0, "Underlying understanding changed, need to reduce size of init_request_bytes")
	log.Warn("init_bytes: %d, initial_ram_usage: %d, init_request_bytes: %d, more_ram: %d.", initBytes, initialRamUsage, initRequestBytes, moreRam)
	e.BuyRamBytes(common.N("eosio"),test1,uint32(moreRam))
	e.BuyRamBytes(common.N("eosio"),test2,uint32(moreRam))

	actName := common.N("setentry")
	setentry := VariantsObject{
		"payer": test1,
		"from":  1,
		"to":    10,
		"size":  1780,
	}

	e.PushAction2(
		&test1,
		&actName,
		test1,
		&setentry,
		e.DefaultExpirationDelta,
		0,
	)
	e.ProduceBlocks(1,false)
	ramUsage := rlm.GetAccountRamUsage(test1)
	total = e.GetTotalStake(test1)
	ramBytes := total["ram_bytes"].(uint64)
	log.Warn("ram_bytes: %d, ram_usage: %d, initial_ram_usage: %d, init_bytes: %d, ram_usage - initial_ram_usage: %d, init_bytes - ram_usage: %d.",
		ramBytes, ramUsage, initialRamUsage, initBytes, ramUsage - initialRamUsage, initBytes - uint64(ramUsage))
	log.Warn("------ram_tests 1------")

}

func TestSimple(t *testing.T){
	fmt.Println(common.MaxMicroseconds().ToSeconds())
}
