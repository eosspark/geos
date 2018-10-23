package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"testing"
)

func TestResourceLimitsManager_UpdateAccountUsage(t *testing.T) {
	control := GetControllerInstance()
	rlm := control.ResourceLimists
	rlm.InitializeDatabase()
	a := common.AccountName(common.N("yuanchao"))
	account := []common.AccountName{a}
	rlm.InitializeAccount(a)
	rlm.AddTransactionUsage(account, 100, 100, 1)
	rlm.UpdateAccountUsage(account, 1)
	rlm.UpdateAccountUsage(account, 86401)
	rlm.UpdateAccountUsage(account, 172801)
	//结果value_ex应该为579/2 579/2/2
	rlm.UpdateAccountUsage(account, 1)
	rlm.UpdateAccountUsage(account, 172801)
	//结果value_ex为0
}

func TestResourceLimitsManager_SetAccountLimits(t *testing.T) {
	control := GetControllerInstance()
	rlm := newResourceLimitsManager(control)
	rlm.InitializeDatabase()
	a := common.AccountName(common.N("yuanchao"))
	rlm.InitializeAccount(a)
	rlm.SetAccountLimits(a, 100, 100, 100)
	var r, n, c int64
	rlm.GetAccountLimits(a, &r, &n, &c)
	fmt.Println(r, n, c)
}

func TestResourceLimitsManager_ProcessBlockUsage(t *testing.T) {

}