package histtory_plugin

import (
	"fmt"
	"github.com/eosspark/eos-go/db"
	"testing"
	"time"

	"github.com/eosspark/eos-go/common"
)

func Test_AddAccountHistoryObject(t *testing.T) {
	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	if err != nil {
		fmt.Println("Test_AddAccountHistoryObject is error detail:", err.Error())
	}
	defer db.Close()

	aho := AccountHistoryObject{}
	aho.Account = common.AccountName(common.StringToName("tuanhuo"))
	aho.AccountSequenceNum = int32(2)
	aho.ActionSequenceNum = uint64(aho.AccountSequenceNum)
	err = db.Insert(&aho)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Test_AddAccountHistoryObject finished")
}
func Test_AddAccountHistoryIndex(t *testing.T) {

	fmt.Println(time.Now())
	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	if err != nil {
		fmt.Println("Test_AddAccountHistoryIndex is error detail:", err.Error())
	}
	defer db.Close()

	aho := AccountHistoryObject{}
	err = db.Get("ID", 1, &aho)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(aho)
	AddAccountHistoryIndex(db, &aho)

	/*for i:=0;i<1000000;i++ {
		aho := AccountHistoryObject{}
		aho.Account = common.AccountName(common.StringToName("tuanhuo"))
		aho.AccountSequenceNum = int32(i+1)
		aho.ActionSequenceNum = uint64(aho.AccountSequenceNum)
		AddAccountHistoryIndex(db, aho)
	}*/

	fmt.Println(time.Now())
}

func Test_GetAccountHistoryIndex(t *testing.T) {
	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	if err != nil {
		fmt.Println("Test_GetAccountHistoryIndex is error detail:", err.Error())
	}
	defer db.Close()

	/*result := AccountHistoryIndex{}
	result.ById = 0
	err = db.Find("ById",result,result)
	if err != nil {
		fmt.Println("Test_GetAccountHistoryIndex is error detail:",err.Error())
	}
	fmt.Println(result)*/

	//results := []AccountHistoryIndex{}
	var results []AccountHistoryIndex

	err = db.All(&results)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(results)
}
