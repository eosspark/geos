package try_test

import (
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/stretchr/testify/assert"
	"testing"
	"errors"
	"fmt"
)

func TestTry_int(t *testing.T) {
	Try(func() {
		panic(1)

	}).Catch(func(n int) {
		assert.Equal(t, 1, n)

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

func returnFunc() (r int, flag bool) {


	defer HandleReturn()
	Try(func() {
		panic(1)
		panic("one error")

		r, flag = 1, false
		Return()

	}).Catch(func(a int) {
		//r = 2
		//Return()
	}).Catch(func(s string) {
		//r = 3
		//Return()
	}).End()

	return 0, true
}

func TestReturn(t *testing.T) {
	fmt.Println(returnFunc())
}

//func TestFinally(t *testing.T) {
//	dofinal := false
//
//	defer func() {
//		recover()
//		assert.Equal(t, true, dofinal)
//	}()
//
//	Try(func() {
//		panic(1)
//
//	}).Catch(func(e string) {
//		// not caught
//
//	}).Finally(func() {
//		dofinal = true
//
//	}).End()
//}
