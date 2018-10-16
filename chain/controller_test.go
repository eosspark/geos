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
	"reflect"
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
	name := common.AccountName(common.N("eos"))

	owner := types.Authority{}
	owner.Threshold = 2
	active := types.Authority{}
	active.Threshold = 1
	control.CreateNativeAccount(name, owner, active, false)
	fmt.Println(name)
	result := types.AccountObject{}
	result.Name = name
	//control.DB.Find("name", result)

	fmt.Println("check account name:", strings.Compare(name.String(), "eos"))
}

func TestController_GetWasmInterface(t *testing.T) {
	control := GetControllerInstance()
	fmt.Println(control.WasmIf)
}

func test(atx *ApplyContext) {
	fmt.Println("this reflect call func :hello test")
}
func TestController_SetApplayHandler(t *testing.T) {
	control := GetControllerInstance()
	receiver := common.AccountName(common.N("reveiver"))
	scope := common.AccountName(common.N("scope"))
	action := common.ActionName(common.N("action"))
	control.SetApplayHandler(receiver, scope, action, test)

	fun := control.FindApplyHandler(receiver, scope, action)

	/*o:=reflect.TypeOf(fun)
	fmt.Println("=========================",o.Kind().String())
	fmt.Println("============1=============",o.MethodByName)
	fmt.Println("=============2============",reflect.ValueOf(fun).MethodByName("test").String())

	fmt.Println("=========================",reflect.ValueOf(fun).String())*/
	ac := ApplyContext{}
	v := []reflect.Value{
		reflect.ValueOf(&ac)}

	fmt.Println("-------address-------", fun)

	reflect.ValueOf(fun).Call(v)
	//a :="test"

	//fmt.Println(strings.Compare(a,o.Name()))
}
