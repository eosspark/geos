package common

import (
	"fmt"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewName(t *testing.T) {
	name := "eosio.system"
	val := N(name)
	assert.Equal(t, val, Name(6138663591228101920))
	//fmt.Printf("%d\n", val)
	name2 := S(6138663591228101920)
	//fmt.Println(name2)
	assert.Equal(t, name2, name)
}

func TestNameStr(t *testing.T) {
	name := "eosio.systemabdxs"
	testflag := false
	var val Name
	try.Try(func() {
		val = N(name)
	}).Catch(func(ex *exception.NameTypeException) {
		//assert.Equal(t, "Invalid name", exception.What(), exception.Message())
		fmt.Println(exception.GetDetailMessage(ex))
		testflag = true
	})
	assert.Equal(t, true, testflag, "check name is wrong")

}

func TestNameSuffix(t *testing.T) {
	name := N("eosio.token")
	check := N("token")
	suffix := NameSuffix(uint64(name))

	assert.Equal(t, check, Name(suffix))
}
