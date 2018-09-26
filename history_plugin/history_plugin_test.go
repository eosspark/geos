package history_plugin

import (
	"fmt"
	"testing"

	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
)

func Test_AddAccountHistoryObject(t *testing.T) {
	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	if err != nil {
		fmt.Println("Test_AddAccountHistoryObject is error detail:", err.Error())
	}
	defer db.Close()

	aho := AccountHistoryObject{}
	acount := common.AccountName(common.StringToName("tuanhuo"))
	num := int32(66)
	aho.Account = acount
	aho.AccountSequenceNum = num
	aho.ByAccountActionSeq.Account = acount
	aho.ByAccountActionSeq.AccountSequenceNum = num
	aho.ActionSequenceNum = uint64(aho.ByAccountActionSeq.AccountSequenceNum)
	c := ByAccountActionSeq{}
	c.Account = aho.ByAccountActionSeq.Account
	c.AccountSequenceNum = aho.ByAccountActionSeq.AccountSequenceNum
	err = db.Insert(&aho)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Test_AddAccountHistoryObject finished")

	fmt.Println("*******************************************")
	result := []AccountHistoryObject{}
	db.All(&result)
	fmt.Print("all context:", result)
}

/*func Test_AddAccountHistorys(t *testing.T) {

	fmt.Println(time.Now())
	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	if err != nil {
		fmt.Println("Test_AddAccountHistoryIndex is error detail:",err.Error())
	}
	defer db.Close()

	results := []AccountHistoryObject{}
	err = db.All(&results)
	if err != nil{
		fmt.Println("db.all() is error",err.Error())
	}
	fmt.Println("db.all() result:",results)
	aho := AccountHistoryObject{}
	aho.ID = 1
	var h AccountHistoryObject
	err = db.Find("ID",&aho.ID,&h)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println(h)
	AddAccountHistoryObject(db, &h)

	for i :=0;i<1000;i++ {
		aho := AccountHistoryObject{}
		aho.Account = (h.Account+common.AccountName(i))
		aho.AccountSequenceNum = int32(i+1)
		aho.ActionSequenceNum = uint64(aho.AccountSequenceNum)
		aho.condition.account = aho.Account
		aho.condition.accountSequenceNum = aho.AccountSequenceNum
		AddAccountHistoryObject(db, &aho)
	}


	fmt.Println(time.Now())
}*/

func Test_GetAccountHistoryObject(t *testing.T) {
	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	if err != nil {
		fmt.Println("Test_GetAccountHistoryObject is error detail:", err.Error())
	}
	defer db.Close()

	result := AccountHistoryObject{}
	result.ID = 2
	err = db.Find("ID", result.ID, &result)
	if err != nil {
		fmt.Println("Test_GetAccountHistoryObject is error detail:", err.Error())
		return
	}
	fmt.Println("object:", result)

	//results := []AccountHistoryIndex{}
	var results []AccountHistoryObject

	err = db.All(&results)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("query all data:", results)
}

func Test_UpdateAccountHistoryObject(t *testing.T) {
	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	if err != nil {
		fmt.Println("Test_UpdateAccountHistoryObject is error detail:", err.Error())
	}
	defer db.Close()

	var results []AccountHistoryObject
	err = db.All(&results)
	if err != nil {
		fmt.Println("update is err :", err.Error())
		return
	}

	//updateObject All properties can not be empty.

	//param :=AccountHistoryObject{}
	param := results[0]
	aho := AccountHistoryObject{}
	aho.ID = param.ID
	//aho.Condition.AccountSequenceNum = param.ByAccountActionSeq.AccountSequenceNum+5
	//aho.Condition.Account = param.ByAccountActionSeq.Account
	fmt.Println("modify object param:", param)
	err = db.UpdateObject(&param, &aho)
	if err != nil {
		fmt.Println("modify object is error:", err.Error())
		return
	}
}

func Test_GetAccountHistoryObjectsByAccount(t *testing.T) {
	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	if err != nil {
		fmt.Println("Test_GetAccountHistoryObjectByAccount is error detail:", err.Error())
		return
	}
	defer db.Close()

	result := GetAccountHistoryObjectByAccount(db, common.AccountName(common.StringToName("tuanhuo2")))
	fmt.Print("Query many data :")
	fmt.Println(result)
}

func Test_GetActions(t *testing.T){
	param :=GetActionParam{}
	param.Pos = 0
	param.Offset = 3
	param.AccountName=common.AccountName(common.StringToName("tuanhuo"))

	result := GetActions(&param)
	fmt.Println("Test_GetActions result :",result)
}