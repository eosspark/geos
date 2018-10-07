/*
 *  @Time : 2018/8/29 下午5:47
 *  @Author : xueyahui
 *  @File : controller_test.go
 *  @Software: GoLand
 */

package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"strings"
	"testing"
)

func TestPopBlock(t *testing.T) {
	con := GetControllerInstance()
	con.PopBlock()
	fmt.Println(con)
}

func TestAbortBlock(t *testing.T) {
	con := GetControllerInstance()
	con.AbortBlock()
	fmt.Println(con)
}

/*func TestSetApplayHandler(t *testing.T) {
	con := GetControllerInstance()
	fmt.Println(con)
	applyCon := ApplyContext{}
	con.SetApplayHandler(111, 111, 111, applyCon)
}*/

func Test_ControllerDB(t *testing.T) {
	control := GetControllerInstance() //chain.GetControllerInstance()
	db := control.DataBase()
	fmt.Println(db)
}

var IrreversibleBlock chan types.BlockState = make(chan types.BlockState)

func TestController_CreateNativeAccount(t *testing.T) {
	//CreateNativeAccount(name common.AccountName,owner types.Authority,active types.Authority,isPrivileged bool)
	control := GetControllerInstance()
	name := common.AccountName(common.S("eos"))

	owner := types.Authority{}
	owner.Threshold = 2
	active := types.Authority{}
	active.Threshold = 1
	control.CreateNativeAccount(name, owner, active, false)
	fmt.Println(name)
	result := types.AccountObject{}
	control.DB.Find("name", name, result)

	fmt.Println("check account name:", strings.Compare(name.String(), "eos"))
}
