package chain

import (
	"testing"
	"github.com/eosspark/eos-go/db"
	"github.com/eosspark/eos-go/common"
	"fmt"
)

func Test_resourceSetGet(t *testing.T) {
	var rlm *ResourceLimitsManager
	db, _ := eosiodb.NewDataBase("./","eos.db", true)
	defer db.Close()
	rlm = NewResourceLimitsManager(db)
	rlm.AddIndices()
	rlm.InitializeDatabase()
	account := common.AccountName(common.StringToName("yc"))
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, 123, 456, 789)
	var r , n , c int64
	rlm.GetAccountLimits(account, &r, &n, &c)
	fmt.Println(r,n,c)
	var rlo []ResourceLimitsObject
	db.All(&rlo)
	fmt.Println(rlo)
	var ruo []ResourceUsageObject
	db.All(&ruo)
	fmt.Println(ruo)
	rlm.ProcessAccountLimitUpdates()
	//查看账户limits与usage的情况
	var rlo2 []ResourceLimitsObject
	db.All(&rlo2)
	fmt.Println(rlo2)
	var ruo2 []ResourceUsageObject
	db.All(&ruo2)
	fmt.Println(ruo2)

}
func Test_resourceFuncAdd(t *testing.T) {
	var rlm *ResourceLimitsManager
	db, _ := eosiodb.NewDataBase("./","eos.db", true)
	defer db.Close()
	rlm = NewResourceLimitsManager(db)
	rlm.AddIndices()
	rlm.InitializeDatabase()
	account1 := common.AccountName(common.StringToName("yc"))
	//account2 := common.AccountName(common.StringToName("sf"))
	//account3 := common.AccountName(common.StringToName("hn"))
	account := []common.AccountName{account1}
	for _, a := range account {
		rlm.InitializeAccount(a)
	}

	//查看账户limits与usage的情况
	var rlo []ResourceLimitsObject
	db.All(&rlo)
	fmt.Println(rlo)
	var ruo []ResourceUsageObject
	db.All(&ruo)
	fmt.Println(ruo)

	rlm.AddTransactionUsage(account, 100000, 100000, 100)
	var ruo2 []ResourceUsageObject
	db.All(&ruo2)
	fmt.Println(ruo2)

	rlm.AddTransactionUsage(account, 100000, 100000, 175600)
	var ruo3 []ResourceUsageObject
	db.All(&ruo3)
	fmt.Println(ruo3)

	rlm.SetAccountLimits(account1, 123, 456, 789)

	rlm.ProcessAccountLimitUpdates()
	//rlm.ProcessBlockUsage(300)
	var state ResourceLimitsStateObject
	db.Find("Id", ResourceLimitsState, &state)
	fmt.Println(state)

	var arl AccountResourceLimit
	arl = rlm.GetAccountCpuLimitEx(common.AccountName(common.StringToName("yc")), true)
	fmt.Println(arl)
}

func Test_(t *testing.T) {

}