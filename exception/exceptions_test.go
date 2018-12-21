package exception_test

import (
	"fmt"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/go-stack/stack"
	"testing"
)

func TestFcException(t *testing.T) {
	Try(func() {
		EosAssert(false, &AssertException{}, "test error")
	}).Catch(func(e Exception) {
		log.Warn(e.DetailMessage())
	})
}

func EOS_ASSERT() {

	caller := log.Record{
		Name: "",
		Call: stack.Caller(2),
	}
	fmt.Println(caller.Call.String())
	fmt.Println(stack.Caller(1).String())
	fmt.Println(stack.Caller(2).String())
	fmt.Println(stack.Caller(3).String())
	fmt.Println(stack.Caller(4).String())
}

func TestStack_line(t *testing.T) {
	EOS_ASSERT()
}
