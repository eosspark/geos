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

func TestException_catch_same(t *testing.T) {
	try.Try(func() {
		EosAssert(false, &NameTypeException{}, "name error")

	}).Catch(func(e NameTypeException) {
		assert.Equal(t, "name error", e.Message())

	}).End()
}

func TestException_catch_same_pointer(t *testing.T) {
	try.Try(func() {
		EosAssert(false, &NameTypeException{}, "name error")

	}).Catch(func(e *NameTypeException) {
		assert.Equal(t, "name error", e.Message())
		assert.Equal(t, ExcTypes(3010001), e.Code())

	}).End()
}

func TestException_catch_diff(t *testing.T) {
	try.Try(func() {
		try.Try(func() {
			EosAssert(false, &NameTypeException{}, "name error")

		}).Catch(func(e BlockValidateException) {
			// BlockValidateException is not conclude NameTypeException, can't be caught

		}).End()

	}).Catch(func(e Exception) {
		assert.Equal(t, "name error", e.Message())
		assert.Equal(t, ExcTypes(3010001), e.Code())

	}).End()
}

func TestException_catch_diff_pointer(t *testing.T) {
	try.Try(func() {
		try.Try(func() {
			EosAssert(false, &NameTypeException{}, "name error")

		}).Catch(func(e *BlockValidateException) {
			// BlockValidateException is not conclude NameTypeException, can't be caught

		}).End()

	}).Catch(func(e Exception) {
		assert.Equal(t, "name error", e.Message())
		assert.Equal(t, ExcTypes(3010001), e.Code())

	}).End()
}

func TestException_catch_interface(t *testing.T) {
	try.Try(func() {
		EosAssert(false, &NameTypeException{}, "name error")

	}).Catch(func(e ChainTypeExceptions) {
		assert.Equal(t, "name error", e.Message())
		assert.Equal(t, ExcTypes(3010001), e.Code())

	}).End()
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
		EosAssert(false, &ChainTypeException{}, "test")
	}).Catch(func(e ChainTypeException) {
		fmt.Println(e.Message())
	}).End()

	//TODO more exceptions
}

func TestReThrow(t *testing.T) {
	try.Try(func() {
		try.Try(func() {
			EosAssert(false, &ChainTypeException{}, "wrong chain type of type:%s", "abc")
		}).Catch(func(e Exception) {
			try.Throw(e) // always == panic(e)
		}).End()

	}).Catch(func(e ChainTypeExceptions) {

		assert.Equal(t, "wrong chain type of type:abc", e.Message())
	}).End()

}
