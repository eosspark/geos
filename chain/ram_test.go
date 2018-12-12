package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"testing"
)

func TestRamTests(t *testing.T){
	e := initEosioSystemTester()
	initRequestBytes := 80000
	//increment_contract_bytes must be less than table_allocation_bytes for this test setup to work
	//incrementContractBytes := 10000
	//tableAllocationBytes := 12000
	e.BuyRamBytes(&common.DefaultConfig.SystemAccountName,common.N("eosio"),7000)
	e.ProduceBlocks(10,false)
	e.CreateAccountWithResources(common.N("testram11111"),common.N("eosio"),
		common.Asset{Amount:int64(initRequestBytes + 40)},false, CoreFromString("10.0000"),CoreFromString("10.0000"))
	e.CreateAccountWithResources(common.N("testram22222"),common.N("eosio"),
		common.Asset{Amount:int64(initRequestBytes + 1190)},false, CoreFromString("10.0000"),CoreFromString("10.0000"))
	e.ProduceBlocks(10,false)
	e.ProduceBlocks(10,false)


}

func TestSimple(t *testing.T){
	fmt.Println(common.MaxMicroseconds().ToSeconds())
}
