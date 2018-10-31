/*
 *  @Time : 2018/8/29 下午5:47
 *  @Author : xueyahui
 *  @File : controller_test.go
 *  @Software: GoLand
 */

package chain

import (
	"crypto/md5"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/entity"
	"github.com/stretchr/testify/assert"
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

func CallBackApplayHandler(p *ApplyContext) {
	fmt.Println("SetApplyHandler CallBack")
}
func TestSetApplyHandler(t *testing.T) {
	con := GetControllerInstance()
	fmt.Println(con)
	//applyCon := ApplyContext{}
	con.SetApplayHandler(111, 111, 111, CallBackApplayHandler)
}

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
	result := entity.AccountObject{}
	result.Name = name
	//control.DB.Find("name", result)

	fmt.Println("check account name:", strings.Compare(name.String(), "eos"))
	assert.Equal(t, "eos", name.String())
}

func TestController_GetWasmInterface(t *testing.T) {
	control := GetControllerInstance()
	fmt.Println(control.WasmIf)
	assert.Equal(t, nil, control.WasmIf)
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

func TestController_GetGlobalProperties(t *testing.T) {
	c := GetControllerInstance()
	result := c.GetGlobalProperties()
	gp := entity.GlobalPropertyObject{}
	gp.ID = common.IdType(1)
	err := c.DB.Find("ID", gp, &gp)
	if err != nil {
		assert.Error(t, err, gp)
	}
	assert.Equal(t, false, common.Empty(result)) //GlobalProperties not initialized
	assert.Equal(t, false, result == &gp)
	c.Close()
}
func TestController_StartBlock(t *testing.T) {
	c := GetControllerInstance()
	w := common.NewBlockTimeStamp(common.Now())
	s := types.Irreversible
	c.StartBlock(w, uint16(s))
	c.Close()
}

func TestController_Clean(t *testing.T) {
	c := GetControllerInstance()
	c.Clean()
}

func TestController_UpdateProducersAuthority(t *testing.T) {
	c := GetControllerInstance()
	c.updateProducersAuthority()
}

func Test(t *testing.T) {
	str := "abc123asdfasdfasdfasdfasdfasdf"

	//方法一
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has)
	fmt.Println(md5str1)
}
