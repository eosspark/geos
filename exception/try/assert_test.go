package try

import (
	"testing"
	. "github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/log"
	"errors"
	"github.com/docker/docker/pkg/testutil/assert"
)

func TestStaticAssert(t *testing.T) {
	Assert(1 != 1, "test assert")
}

func TestEosAssert(t *testing.T) {
	EosAssert(true, &BlockValidateException{}, "block #%s error :%s", "00000006367c1f4...", "msg")

	Try(func() {
		EosAssert(false, &BlockNetUsageExceeded{}, "tester exception %s", "BlockNetUsageExceeded")
	}).Catch(func(e Exception) {
		log.Error(GetDetailMessage(e))
	}).End()
}

func TestCatchOrFinally_EosRethrowExceptions(t *testing.T) {
	defer HandleStackInfo()
	Try(func() {
		Try(func() {
			EosAssert(false, &AssertException{}, "tester exception %s", "AssertException")
		}).EosRethrowExceptions(&ChainTypeException{}, "block #%d assert", 100).End()
	}).Catch(func(e *ChainTypeException) {
		detail := GetDetailMessage(e)
		log.Error(detail)
		assert.Equal(t, detail, "3010000 *exception.ChainTypeException: chain type exception\n"+
			"block #100 assert\n"+
			"tester exception AssertException\n")
	}).End()
}

func TestCatchOrFinally_FcLogAndRethrow(t *testing.T) {
	Try(func() {
		Try(func() {
			EosAssert(false, &BlockNetUsageExceeded{}, "tester exception %s", "BlockNetUsageExceeded")
		}).FcLogAndRethrow().End()
	}).Catch(func(e Exception) {
		detail := GetDetailMessage(e)
		log.Error(detail)
		assert.Equal(t, detail, "3080003 *exception.BlockNetUsageExceeded: Transaction network usage is too much for the remaining allowable usage of the current block\n"+
			"tester exception BlockNetUsageExceeded\n"+
			"rethrow\n")
	}).End()

	Try(func() {
		Try(func() {
			Throw(errors.New("tester error"))
		}).FcLogAndRethrow().End()
	}).Catch(func(e Exception) {
		detail := GetDetailMessage(e)
		log.Error(detail)
		assert.Equal(t, detail, "0 *exception.FcException: unspecified\n"+
			"rethrow: tester error\n")
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
		}).FcRethrowExceptions(log.LvlWarn, "block %s", "001bac").End()
	}).Catch(func(e Exception) {
		detail := GetDetailMessage(e)
		log.Error(detail)
		assert.Equal(t, detail, "3080003 *exception.BlockNetUsageExceeded: Transaction network usage is too much for the remaining allowable usage of the current block\n"+
			"tester exception BlockNetUsageExceeded\n"+
			"block 001bac\n")
	}).End()

	Try(func() {
		Try(func() {
			Throw(errors.New("test error"))
		}).FcRethrowExceptions(log.LvlWarn, "block %s", "001bac").End()
	}).Catch(func(e Exception) {
		log.Error(GetDetailMessage(e))
	}).End()

	Try(func() {
		Try(func() {
			Throw(0)
		}).FcRethrowExceptions(log.LvlWarn, "block %s", "001bac").End()
	}).Catch(func(e Exception) {
		log.Error(GetDetailMessage(e))
	}).End()

}

func TestCatchOrFinally_FcCaptureLogAndRethrow(t *testing.T) {
	Try(func() {
		Try(func() {
			EosAssert(false, &BlockNetUsageExceeded{}, "tester exception %s", "BlockNetUsageExceeded")
		}).FcCaptureLogAndRethrow("rethrow %s", "exception").End()
	}).Catch(func(e Exception) {
		detail := GetDetailMessage(e)
		log.Error(detail)
		assert.Equal(t, detail, "3080003 *exception.BlockNetUsageExceeded: Transaction network usage is too much for the remaining allowable usage of the current block\n"+
			"tester exception BlockNetUsageExceeded\n"+
			"rethrow rethrow exception\n")
	}).End()
}

func TestCatchOrFinally_FcCaptureAndRethrow(t *testing.T) {
	Try(func() {
		Try(func() {
			EosAssert(false, &BlockNetUsageExceeded{}, "tester exception %s", "BlockNetUsageExceeded")
		}).FcCaptureAndRethrow("rethrow %s", "exception").End()
	}).Catch(func(e Exception) {
		detail := GetDetailMessage(e)
		log.Error(detail)
		assert.Equal(t, detail, "3080003 *exception.BlockNetUsageExceeded: Transaction network usage is too much for the remaining allowable usage of the current block\n"+
			"tester exception BlockNetUsageExceeded\n"+
			"rethrow exception\n")
	}).End()

	Try(func() {
		Try(func() {
			Throw(errors.New("tester error"))
		}).FcCaptureAndRethrow("rethrow %s", "error").End()
	}).Catch(func(e Exception) {
		detail := GetDetailMessage(e)
		log.Error(detail)
		assert.Equal(t, detail, "0 *exception.FcException: unspecified\n"+
			"tester error: rethrow error\n")
	}).End()

	Try(func() {
		Try(func() {
			Throw(0)
		}).FcCaptureAndRethrow("rethrow %s", "any").End()
	}).Catch(func(e Exception) {
		log.Error(GetDetailMessage(e))
	}).End()
}
