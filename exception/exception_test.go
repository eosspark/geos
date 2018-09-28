package exception

import (
	"testing"
)

func assert() {
	EosAssert(false, ChainException, "wrong chain error : %s", 123)

	//defer func() {
	//	exp := recover()
	//	if exp, ok := recover().(Exception); ok {
	//		fmt.Println(exp.code, exp.message)
	//	}
	//}()
}

func TestEosAssert(t *testing.T) {
	assert()
}
