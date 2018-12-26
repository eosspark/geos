package unittests

import (
	"testing"
	"github.com/stretchr/testify/assert"

	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
)

func CheckThrow(t *testing.T, f func(), exception Exception) {
	check := false
	Try(f).Catch(func(e Exception) {
		assert.Equal(t, exception.Code(), e.Code())
		check = true
	}).End()
	assert.Equal(t, true, check)
}

func CheckNoThrow(t *testing.T, f func()) {
	Try(f).Catch(func(e Exception) {
		assert.Fail(t, e.DetailMessage())
	}).Catch(func(interface{}) {
		assert.Fail(t, "check no throw failed")
	})
}

