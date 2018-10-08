/*
 *  @Time : 2018/9/17 下午4:15
 *  @Author : xueyahui
 *  @File : contractTO_test.go
 *  @Software: GoLand
 */

package types

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"github.com/eosspark/eos-go/log"
	"testing"
)

func Test_Add_TableIdObject(t *testing.T) {

	code := common.AccountName(common.N("eostest"))
	scope := common.ScopeName(common.N("eostest"))
	table := common.TableName(common.N("eostest"))
	tid := TableIdObject{}
	tid.Code = code
	tid.Scope = scope
	tid.Payer = code
	tid.Table = table
	tid.Count = 2

	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	defer db.Close()
	if err != nil {
		log.Error("Test_Add_TableIdObject is error detail:", err)
	}
	db.Insert(&tid)
	fmt.Println(tid)
}

func Test_m(t *testing.T) {

}

func Test_Add_TableIdMeltiIndex(t *testing.T) {
	code := common.AccountName(common.N("eosio.token"))
	scope := common.ScopeName(common.N("xiaoyu"))
	table := common.TableName(common.N("accounts"))
	tid := TableIdObject{}
	tid.Code = code
	tid.Scope = scope
	tid.Payer = code
	tid.Table = table
	tid.Count = tid.Count + 1

	ti := TableIdMultiIndex{}
	ti.TableIdObject = tid
	ti.Bst.Code = tid.Code
	ti.Bst.Scope = tid.Scope
	ti.Bst.Table = tid.Table
	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	defer db.Close()
	if err != nil {
		log.Error("Test_Add_TableIdMeltiIndex is error detail:", err)
	}
	fmt.Println(ti)
	db.Insert(&ti)
	fmt.Println(ti)
}

func Test_Get_TableIdMultiIndex(t *testing.T) {
	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	if err != nil {
		log.Error("Test_Add_TableIdMeltiIndex is error detail:", err)
	}
	defer db.Close()
	ti := TableIdMultiIndex{}
	code := common.AccountName(common.N("eostest"))
	scope := common.ScopeName(common.N("eostest"))
	table := common.TableName(common.N("eostest"))
	ti.Bst.Code = code
	ti.Bst.Scope = scope
	ti.Bst.Table = table
	tmp := TableIdMultiIndex{}
	err = db.Find("Bst", ti.Bst, &tmp)
	if err != nil {
		log.Error("Test_Get_TableIdMeltiIndex byCodeScopeTable is error detail:", err)
	}
	fmt.Println(&tmp)
	log.Info("find table id multi index,info:", tmp)
	var tis []TableIdMultiIndex
	db.All(&tis)
	fmt.Println(tis)
}

func Test_GetById(t *testing.T) {
	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	if err != nil {
		log.Error("Test_Add_TableIdMeltiIndex is error detail:", err)
	}
	defer db.Close()

	tt := TableIdObject{}

	tt.ID = 1
	db.Find("ID", tt.ID, &tt)
	fmt.Println(tt)
	var tmp []TableIdObject
	db.All(&tmp)
	fmt.Println(tmp)
}

func Test_GetByIndexId(t *testing.T) {
	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	if err != nil {
		log.Error("Test_Add_TableIdMeltiIndex is error detail:", err)
	}
	defer db.Close()

	tt := TableIdMultiIndex{}
	//var tt TableIdMultiIndex
	tt.Id = 1
	err = db.Find("ID", tt.Id, &tt)

	if err != nil {
		log.Error("find error detail:", err)
		fmt.Println(err.Error())
	}
	fmt.Println(&tt)
}

func Test_GetByCodeScopeTable(t *testing.T) {
	db, err := eosiodb.NewDataBase("./", "shared_memory.bin", true)
	if err != nil {
		log.Error("Test_Add_TableIdMeltiIndex is error detail:", err)
	}
	//defer db.Close()

	cst := ByCodeScopeTable{}
	cst.Code = common.AccountName(common.N("eosio.token"))
	cst.Scope = common.ScopeName(common.N("xiaoyu"))
	cst.Table = common.TableName(common.N("accounts"))

	/*fmt.Println(cst)
	tmp:=GetByCodeScopeTable(db,cst)
	fmt.Println(tmp)*/

	tmi := TableIdMultiIndex{}
	err = db.Find("Bst", cst, &tmi)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(tmi)
}
