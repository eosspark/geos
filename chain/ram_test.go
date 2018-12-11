package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"testing"
)

func TestRamTests(t *testing.T){
	//initRequestBytes := 80000
	////increment_contract_bytes must be less than table_allocation_bytes for this test setup to work
	//incrementContractBytes := 10000
	//tableAllocationBytes := 12000
	initEosioSystemTester()
}

func TestSimple(t *testing.T){
	fmt.Println(common.S(6138663591592764928))
}
