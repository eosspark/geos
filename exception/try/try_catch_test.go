package try

import (
	"errors"
	"fmt"
	//"github.com/eosspark/eos-go/exceptionx"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/log"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestCatch_message(t *testing.T) {
	Try(func() {
		a, b := 1, 0
		c := a / b
		fmt.Println(c)
	}).Catch(func(e exception.Exception) {
		log.Warn(e.DetailMessage())
	}).End()

	Try(func() {
		Try(func() {
			a := &struct {
				x int
			}{}
			a = nil
			fmt.Println(a.x)
		}).EosRethrowExceptions(&exception.TransactionTypeException{}, "eos exception re")
	}).Catch(func(e *exception.TransactionTypeException) {
		log.Warn(e.DetailMessage())
	})

	Try(func() {
		var m map[int]int
		m[1] = 2
	}).FcLogAndDrop("log and drop ")
}
