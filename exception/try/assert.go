package try

import (
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/log"
	"fmt"
)

func EosAssert(expr bool, exception Exception, format string, args ...interface{}) {
	if !expr {
		FcThrowException(exception, format, args...)
	}
}

func EosThrow(exception Exception, format string, args ...interface{}) {
	exception.FcLogMessage(LvlError, format, args...)
	Throw(exception)
}

func (c *CatchOrFinally) EosRethrowExceptions(exception Exception, format string, args ...interface{}) *CatchOrFinally {
	return c.Catch(func(e ChainExceptions) {
		FcRethrowException(e, LvlWarn, format, args...)

	}).Catch(func(e Exception) {
		exception.FcLogMessage(LvlWarn, fmt.Sprintf("%s, %s", format, e.Message()), args...)
		Throw(exception)

	}).Catch(func(e error) {
		exception.FcLogMessage(LvlWarn, fmt.Sprintf("%s (%s)", format, e.Error()))
		Throw(exception)

	}).Catch(func(interface{}) {
		Throw(&UnHandledException{LogMessage{LvlWarn, format, args}})
	})
}


func FcAssert(test bool, args ...interface{}) {
	if !test {
		FcThrowException(&AssertException{}, "assert:", args...)
	}
}

func FcCaptureAndThrow(exception Exception, format string, args ...interface{}) {
	exception.FcLogMessage(LvlError, format, args...)
	Throw(exception)
}

func FcThrow(format string, args ...interface{}) {
	Throw(&FcException{LogMessage{LvlError, format, args}})
}

func FcThrowException(exception Exception, format string, args ...interface{}) {
	exception.FcLogMessage(LvlError, format, args...)
	Throw(exception)
}

func FcRethrowException(er Exception, logLevel LogLevel, format string, args ...interface{}) {
	er.FcLogMessage(LvlWarn, fmt.Sprintf("%s, %s", er.Message(), format), args...)
	Throw(er)
}

func (c *CatchOrFinally) FcLogAndRethrow() *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(er.Message())

		FcRethrowException(er, LvlWarn, "rethrow")
	}).Catch(func(e error) {
		fce := &FcException{}
		fce.FcLogMessage(LvlWarn, "rethrow %s: ", e.Error())
		Warn(fce.Message())
		Throw(fce)

	}).Catch(func(a interface{}) {
		e := UnHandledException{}
		e.FcLogMessage(LvlWarn, "rethrow", a)
		Warn(e.Message())
		Throw(e)
	})
}

func (c *CatchOrFinally) FcCaptureLogAndRethrow(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(er.Message())
		FcRethrowException(er, LvlWarn, "rethrow", args...)

	}).Catch(func(e error) {
		fce := &FcException{}
		fce.FcLogMessage(LvlWarn, "rethrow %s ", e.Error())
		fce.Message()
		Throw(fce)

	}).Catch(func(a interface{}) {
		e := &UnHandledException{}
		e.FcLogMessage(LvlWarn, "rethrow", a)
		Warn(e.Message())
		Throw(e)
	})
}

func (c *CatchOrFinally) FcCaptureAndLog(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(er.Message())

	}).Catch(func(e error) {
		fce := &FcException{}
		fce.FcLogMessage(LvlWarn, "rethrow %s: ", e.Error())
		Warn(fce.Message())

	}).Catch(func(a interface{}) {
		e := &UnHandledException{}
		e.FcLogMessage(LvlWarn, "rethrow", a)
		Warn(e.Message())
	})
}

func (c *CatchOrFinally) FcLogAndDrop(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(er.Message())

	}).Catch(func(e error) {
		fce := &FcException{}
		fce.FcLogMessage(LvlWarn, "rethrow %s: ", e.Error())
		Warn(fce.Message())

	}).Catch(func(a interface{}) {
		e := &UnHandledException{}
		e.FcLogMessage(LvlWarn, "rethrow", a)
		Warn(e.Message())
	})
}

func (c *CatchOrFinally) FcRethrowExceptions(logLevel LogLevel, format string, args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception){
		FcRethrowException(er, logLevel, format, args...)

	}).Catch(func(e error) {
		fce := &FcException{}
		fce.FcLogMessage(logLevel, "%s: ", e.Error())
		Throw(fce)

	}).Catch(func(interface{}) {
		e := &UnHandledException{}
		e.FcLogMessage(logLevel, format, args...)
		Throw(e)
	})
}

func (c *CatchOrFinally) FcCaptureAndRethrow(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		FcRethrowException(er, LvlWarn, "", args...)

	}).Catch(func(e error) {
		fce := &FcException{}
		fce.FcLogMessage(LvlWarn, "%s: ", e.Error())
		Throw(fce)

	}).Catch(func(interface{}) {
		e := &UnHandledException{}
		e.FcLogMessage(LvlWarn, "")
		Throw(e)
	})
}
