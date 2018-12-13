package try

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"errors"
	"github.com/eosspark/eos-go/exception"
)

func TestTry_int(t *testing.T) {
	Try(func() {
		panic(1)

	}).Catch(func(n int) {
		//assert.Equal(t, 1, n)

	}).End()
}

func TestTry_string(t *testing.T) {
	Try(func() {
		panic("try")

	}).Catch(func(n string) {
		assert.Equal(t, "try", n)

	}).End()
}

func TestTry_pointer(t *testing.T) {
	Try(func() {
		panic(&struct{}{})

	}).Catch(func(n struct{}) {
		assert.Equal(t, struct{}{}, n)

	}).End()

	// also catch pointer type
	Try(func() {
		panic(&struct{}{})

	}).Catch(func(n *struct{}) {
		assert.Equal(t, struct{}{}, *n)

	}).End()
}

func TestTry_RuntimeError(t *testing.T) {
	Try(func() {
		a, b := 1, 0
		println(a / b)

	}).Catch(func(n RuntimeError) {
		assert.Equal(t, "runtime error: integer divide by zero", n.String())

	}).End()
}

func TestCatch_all(t *testing.T) {
	Try(func() {
		panic("123")

	}).Catch(func(interface{}) {

	}).End()

	Try(func() {
		panic(1)

	}).Catch(func(interface{}) {

	}).End()

	Try(func() {
		panic(struct {
			a int
			b string
		}{})

	}).Catch(func(interface{}) {

	}).End()

	Try(func() {
		panic(errors.New(""))

	}).Catch(func(interface{}) {

	}).End()

}

func TestStackInfo(t *testing.T) {
	defer func() {
		assert.Equal(t, true, len(stackInfo) == 0)
	}()

	Try(func() {
		Throw("error")
	}).Catch(func(n string) {

	}).End()
}

func TestStackInfo_throw(t *testing.T) {
	defer func() {
		recover()
		assert.Equal(t, true, len(stackInfo) > 0)
	}()

	Try(func() {
		Throw("error")
	}).Catch(func(n int) {

	}).End()
}

func TestStackInfo_rethrow(t *testing.T) {
	defer func() {
		recover()
		assert.Equal(t, true, len(stackInfo) > 0)
	}()

	Try(func() {
		EosThrow(&exception.TransactionTypeException{}, "error")
	}).Catch(func(n string) {
		Throw(n)
	}).End()
}
