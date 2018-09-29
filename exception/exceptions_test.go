package exception

import (
	"fmt"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExceptionCode(t *testing.T) {
	assert.Equal(t, ExcTypes(21), divideByZeroCode)
	assert.Equal(t, ExcTypes(10), assertExceptionCode)
	assert.Equal(t, ExcTypes(15), unknownHostExceptionCode)
}

func TestEosAssert(t *testing.T) {
	EosAssert(true, &BlockValidateException{}, "block #%s error :%s", "00000006367c1f4...", "msg")
}

func TestEosAssert_catch(t *testing.T) {
	var scopeExit int
	defer func() {
		assert.Equal(t, 1, scopeExit)
	}()

	try.Try(func() {
		EosAssert(false, &ChainException{}, "test")
	}).Catch(func(e ChainExceptions) {
		fmt.Println(e.What())
	}).End()

	scopeExit = 1

}

func TestExceptions(t *testing.T) {
	try.Try(func() {
		EosAssert(false, &ChainTypeException{}, "wrong chain type of type:%s", "abc")
	}).Catch(func(e ChainExceptions) {
		assert.Equal(t, "wrong chain type of type:abc", e.Message())
	}).End()

	try.Try(func() {
		EosAssert(false, &ChainException{}, "wrong chain id:%d", 12345)
	}).Catch(func(e Exception) {
		assert.Equal(t, "wrong chain id:12345", e.Message())
	}).End()

	try.Try(func() {
		EosAssert(false, &BlockValidateException{}, "test")
	}).Catch(func(e BlockValidateExceptions) {
		assert.Equal(t, "test", e.Message())
	}).End()

	try.Try(func() {
	}).Catch(func(e ChainExceptions) {}).End()
}
