package exception_test

import (
		. "github.com/eosspark/eos-go/exception/try"
	. "github.com/eosspark/eos-go/exception"
	"github.com/stretchr/testify/assert"
	"testing"
	)

func TestExceptionCode(t *testing.T) {
	assert.Equal(t, ExcTypes(21), DivideByZeroCode)
	assert.Equal(t, ExcTypes(10), AssertExceptionCode)
	assert.Equal(t, ExcTypes(15), UnknownHostExceptionCode)
}

func TestException_catch_same_pointer(t *testing.T) {
	Try(func() {
		EosAssert(false, &NameTypeException{}, "name error")

	}).Catch(func(e NameTypeException) {
		assert.Equal(t, ExcTypes(3010001), e.Code())

	}).End()
}

func TestException_catch_diff(t *testing.T) {
	Try(func() {
		Try(func() {
			EosAssert(false, &NameTypeException{}, "name error")

		}).Catch(func(e BlockValidateException) {
			// BlockValidateException is not conclude NameTypeException, can't be caught

		}).End()

	}).Catch(func(e Exception) {
		assert.Equal(t, ExcTypes(3010001), e.Code())

	}).End()
}

func TestException_catch_diff_pointer(t *testing.T) {
	Try(func() {
		Try(func() {
			EosAssert(false, &NameTypeException{}, "name error")

		}).Catch(func(e BlockValidateException) {
			// BlockValidateException is not conclude NameTypeException, can't be caught

		}).End()

	}).Catch(func(e Exception) {
		assert.Equal(t, ExcTypes(3010001), e.Code())

	}).End()
}

func TestException_catch_interface(t *testing.T) {
	Try(func() {
		EosAssert(false, &NameTypeException{}, "name error")

	}).Catch(func(e ChainTypeExceptions) {
		assert.Equal(t, ExcTypes(3010001), e.Code())

	}).End()
}


func TestReThrow(t *testing.T) {
	var s string
	Try(func() {
		Try(func() {
			EosAssert(false, &ChainTypeException{}, "wrong chain type of type:%s", "abc")
		}).Catch(func(e Exception) {
			s = GetDetailMessage(e)
			Throw(e) // always == panic(e)
		}).End()

	}).Catch(func(e ChainTypeExceptions) {
		assert.Equal(t, s, GetDetailMessage(e))
	}).End()

}

func TestCatchExceptionsRethrow(t *testing.T) {
	Try(func() {
		EosAssert(false, &ChainTypeException{}, "")
	}).Catch(func(e Exception) {
		Try(func() {
			Throw(e)
		}).Catch(func(e ChainTypeException) {
			//shouldn't throw
		}).End()
	}).End()
}
