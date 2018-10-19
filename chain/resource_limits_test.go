package chain

import (
	"testing"
	"github.com/eosspark/eos-go/common"
)

func TestResourceLimitsManager_UpdateAccountUsage(t *testing.T) {
	rlm := GetResourceLimitsManager()
	rlm.InitializeDatabase()
	a := common.AccountName(common.N("yuanchao"))
	account := []common.AccountName{a}
	//fmt.Println(account)
	rlm.InitializeAccount(a)
	rlm.AddTransactionUsage(account, 100, 100, 1)
	rlm.UpdateAccountUsage(account,100)
}