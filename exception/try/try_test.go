package try_test

import (
	"github.com/eosspark/eos-go/exception/try"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestTry_RuntimeError(t *testing.T) {
	try.Try(func() {
		a, b := 1, 0
		println(a / b)
	}).Catch(func(n try.RuntimeError) {
		assert.Equal(t, "runtime error: integer divide by zero", n.String())
	}).End()
}

func TestFinally(t *testing.T) {

}
