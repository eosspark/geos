package try_test

import (
	"github.com/eosspark/eos-go/exception/try"
	"github.com/stretchr/testify/assert"
	"testing"
	"errors"
)

func TestTry_int(t *testing.T) {
	try.Try(func() {
		panic(1)

	}).Catch(func(n int) {
		assert.Equal(t, 1, n)

	}).End()
}

func TestTry_string(t *testing.T) {
	try.Try(func() {
		panic("try")

	}).Catch(func(n string) {
		assert.Equal(t, "try", n)

	}).End()
}

func TestTry_pointer(t *testing.T) {
	try.Try(func() {
		panic(&struct {}{})

	}).Catch(func(n struct{}) {
		assert.Equal(t, struct{}{}, n)

	}).End()

	// also catch pointer type
	try.Try(func() {
		panic(&struct {}{})

	}).Catch(func(n *struct{}) {
		assert.Equal(t, struct{}{}, *n)

	}).End()
}

func TestTry_RuntimeError(t *testing.T) {
	try.Try(func() {
		a, b := 1, 0
		println(a / b)

	}).Catch(func(n try.RuntimeError) {
		assert.Equal(t, "runtime error: integer divide by zero", n.String())

	}).End()
}

func TestCatch_all(t *testing.T) {
	try.Try(func() {
		panic("123")

	}).Catch(func(interface{}) {

	}).End()

	try.Try(func() {
		panic(1)

	}).Catch(func(interface{}) {

	}).End()

	try.Try(func() {
		panic(struct {
			a int
			b string
		}{})

	}).Catch(func(interface{}) {

	}).End()

	try.Try(func() {
		panic(errors.New(""))

	}).Catch(func(interface{}) {

	}).End()

}

//func TestFinally(t *testing.T) {
//	dofinal := false
//
//	defer func() {
//		recover()
//		assert.Equal(t, true, dofinal)
//	}()
//
//	try.Try(func() {
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
