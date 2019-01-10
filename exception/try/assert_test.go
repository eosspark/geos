package try

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"strings"

	//. "github.com/eosspark/eos-go/exceptionx"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/log"
	"testing"
)

func TestStaticAssert(t *testing.T) {
	//Assert(1 != 1, "test assert")
}

func TestErrorThrow(t *testing.T) {
	Try(func() {
		Try(func() {
			var a *int
			*a ++
		}).Catch(func(e Exception) {
			FcRethrowException(e, LvlWarn, "rethrow error")
		}).End()
	}).Catch(func(e StdException) {
		detail := e.DetailMessage()
		assert.True(t, strings.Contains(detail, "rethrow error"))
		assert.True(t, strings.Contains(detail, "assert_test.go:22"))
		Error(detail)
	}).End()

}

func TestEosAssert(t *testing.T) {
	EosAssert(true, &BlockValidateException{}, "block #%s error :%s", "00000006367c1f4...", "msg")

	Try(func() {
		EosAssert(false, &BlockNetUsageExceeded{}, "tester exception %s", "BlockNetUsageExceeded")
	}).Catch(func(e Exception) {
		detail := e.DetailMessage()
		assert.True(t, strings.Contains(detail, "tester exception BlockNetUsageExceeded"))
		assert.True(t, strings.Contains(detail, "assert_test.go:39"))
		Error(detail)
	}).End()
}

func TestFcAssert(t *testing.T) {
	Try(func() {
		FcAssert(false, "tester exception %s", "BlockNetUsageExceeded")
	}).Catch(func(e Exception) {
		detail := e.DetailMessage()
		assert.True(t, strings.Contains(detail, "assert:"))
		assert.True(t, strings.Contains(detail, "tester exception BlockNetUsageExceeded"))
		assert.True(t, strings.Contains(detail, "assert_test.go:50"))
		Error(detail)
	}).End()
}

func TestEosThrow(t *testing.T) {
	Try(func() {
		Try(func() {
			EosThrow(&DatabaseGuardException{}, "tester exception %s", "DatabaseGuardException")
		}).Catch(func(e Exception) {
			assert.True(t, strings.Contains(e.DetailMessage(), "tester exception DatabaseGuardException"))
			assert.True(t, strings.Contains(e.DetailMessage(), "assert_test.go:63"))
			Throw(e)
		}).End()
	}).Catch(func(e DatabaseExceptions) {
		assert.True(t, strings.Contains(e.DetailMessage(), "tester exception DatabaseGuardException"))
		assert.True(t, strings.Contains(e.DetailMessage(), "assert_test.go:63"))
		Warn(e.DetailMessage())
	}).End()
}

func TestFcThrow(t *testing.T) {
	Try(func() {
		FcThrow("tester exception %s", "FcThrow")
	}).Catch(func(e Exception) {
		assert.True(t, strings.Contains(e.DetailMessage(), "FcException"))
		assert.True(t, strings.Contains(e.DetailMessage(), "tester exception FcThrow"))
		assert.True(t, strings.Contains(e.DetailMessage(), "assert_test.go:78"))
		Error(e.DetailMessage())
	}).End()
}

func TestFcRethrowException(t *testing.T) {
	Try(func() {
		Try(func() {
			FcThrow("tester exception %s", "FcThrow")
		}).Catch(func(e Exception) {
			assert.True(t, strings.Contains(e.DetailMessage(), "tester exception FcThrow"))
			assert.True(t, strings.Contains(e.DetailMessage(), "assert_test.go:90"))
			FcRethrowException(e, LvlWarn, "rethrow, FcRethrowException")
		}).End()
	}).Catch(func(e Exception) {
		assert.True(t, strings.Contains(e.DetailMessage(), "tester exception FcThrow"))
		assert.True(t, strings.Contains(e.DetailMessage(), "rethrow, FcRethrowException"))
		assert.True(t, strings.Contains(e.DetailMessage(), "assert_test.go:90"))
		assert.True(t, strings.Contains(e.DetailMessage(), "assert_test.go:94"))
		Warn(e.DetailMessage())
	}).End()

}

func TestCatchOrFinally_EosRethrowExceptions(t *testing.T) {
	Try(func() {
		Try(func() {
			EosAssert(false, &ActionTypeException{}, "tester exception %s", "ActionTypeException")
		}).EosRethrowExceptions(&ChainTypeException{}, "eos rethrow")
	}).Catch(func(e ActionTypeException) {
		detail := e.DetailMessage()
		Error(detail)
	}).End()

	Try(func() {
		Try(func() {
			EosAssert(false, &AssertException{}, "tester exception %s", "AssertException")
		}).EosRethrowExceptions(&ChainTypeException{}, "eos rethrow non-chainException")
	}).Catch(func(e *ChainTypeException) {
		detail := e.DetailMessage()
		Error(detail)
	}).End()

	Try(func() {
		Try(func() {
			Try(func() {
				EosAssert(false, &AssertException{}, "tester exception %s", "AssertException")
			}).EosRethrowExceptions(&ChainTypeException{}, "eos rethrow non-chainException")

		}).EosRethrowExceptions(&BlockCpuUsageExceeded{}, "eos rethrow twice non-chainException")

	}).Catch(func(e ChainTypeException) {
		Error(e.DetailMessage())
	}).End()

}

func TestCatchOrFinally_FcLogAndRethrow(t *testing.T) {
	Try(func() {
		Try(func() {
			EosAssert(false, &BlockNetUsageExceeded{}, "tester exception %s", "BlockNetUsageExceeded")
		}).FcLogAndRethrow().End()
	}).Catch(func(e Exception) {
		Error(e.DetailMessage())
	}).End()

	Try(func() {
		Try(func() {
			Throw(errors.New("tester error"))
		}).FcLogAndRethrow().End()
	}).Catch(func(e Exception) {
		detail := e.DetailMessage()
		Error(detail)
	}).End()
}

func TestCatchOrFinally_FcCaptureAndLog(t *testing.T) {
	Try(func() {
		EosAssert(false, &BlockNetUsageExceeded{}, "tester exception %s", "BlockNetUsageExceeded")
	}).FcCaptureAndLog("blocknum %d", 1).End()

	Try(func() {
		Throw(errors.New("tester error"))
	}).FcCaptureAndLog("blocknum %d", 1).End()

	Try(func() {
		Throw(0)
	}).FcCaptureAndLog().End()
}

func TestCatchOrFinally_FcLogAndDrop(t *testing.T) {
	Try(func() {
		EosAssert(false, &BlockNetUsageExceeded{}, "tester exception %s", "BlockNetUsageExceeded")
	}).FcLogAndDrop("blocknum %d", 1).End()

	Try(func() {
		Throw(errors.New("tester error"))
	}).FcLogAndDrop("blocknum %d", 1).End()

	Try(func() {
		Throw(0)
	}).FcLogAndDrop().End()
}

func TestCatchOrFinally_FcRethrowExceptions(t *testing.T) {
	Try(func() {
		Try(func() {
			EosAssert(false, &BlockNetUsageExceeded{}, "tester exception %s", "BlockNetUsageExceeded")
		}).FcRethrowExceptions(LvlWarn, "block %s", "001bac").End()
	}).Catch(func(e Exception) {
		detail := e.DetailMessage()
		Error(detail)
	}).End()

	Try(func() {
		Try(func() {
			Throw(errors.New("test error"))
		}).FcRethrowExceptions(LvlWarn, "block %s", "001bac").End()
	}).Catch(func(e Exception) {
		Error(e.DetailMessage())
	}).End()
}

func TestCatchOrFinally_FcCaptureLogAndRethrow(t *testing.T) {
	Try(func() {
		Try(func() {
			EosAssert(false, &BlockNetUsageExceeded{}, "tester exception %s", "BlockNetUsageExceeded")
		}).FcCaptureLogAndRethrow("rethrow %s", "exception").End()
	}).Catch(func(e Exception) {
		detail := e.DetailMessage()
		Error(detail)
	}).End()
}

func TestCatchOrFinally_FcCaptureAndRethrow(t *testing.T) {
	Try(func() {
		Try(func() {
			EosAssert(false, &BlockNetUsageExceeded{}, "tester exception %s", "BlockNetUsageExceeded")
		}).FcCaptureAndRethrow("rethrow %s", "exception").End()
	}).Catch(func(e Exception) {
		detail := e.DetailMessage()
		Error(detail)
	}).End()

	Try(func() {
		Try(func() {
			Throw("tester throw")
		}).FcCaptureAndRethrow("rethrow").End()
	}).Catch(func(e Exception) {
		detail := e.DetailMessage()
		Error(detail)
	}).End()
}
