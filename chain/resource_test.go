package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"testing"
)

func Test_resource(t *testing.T) {
	var rlm *ResourceLimitsManager
	db, _ := eosiodb.NewDataBase("./", "eos.db", true)
	defer db.Close()
	rlm = NewResourceLimitsManager(db)
	rlm.AddIndices()
	rlm.InitializeDatabase()
	account := common.AccountName(123)
	rlm.InitializeAccount(account)
	rlm.SetAccountLimits(account, 123, 456, 789)
	var r, n, c int64
	rlm.GetAccountLimits(account, &r, &n, &c)
	fmt.Println(r, n, c)
}
