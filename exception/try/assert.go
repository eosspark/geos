package try

import (
	//. "github.com/eosspark/eos-go/exceptionx"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/log"
	"os"
)

func Assert(expr bool, message string) {
	if !expr {
		println(message)
		os.Exit(1)
	}
}

const (
	logMessageSkip         = 2
	fcRethrowExceptionSkip = 3
)

func EosAssert(expr bool, exception Exception, format string, args ...interface{}) {
	if !expr {
		exception.AppendLog(LogMessage(LvlError, format, args, logMessageSkip))
		panic(exception)
	}
}

func FcAssert(test bool, args ...interface{}) {
	if !test {
		format, arg := FcFormatArgParams(args)
		panic(&AssertException{Elog: []Message{LogMessage(LvlError, "assert:"+format, arg, logMessageSkip)}})
	}
}

func EosThrow(exception Exception, format string, args ...interface{}) {
	exception.AppendLog(LogMessage(LvlError, format, args, logMessageSkip))
	Throw(exception)
}

func FcThrow(format string, args ...interface{}) {
	Throw(&FcException{Elog: []Message{LogMessage(LvlError, format, args, logMessageSkip)}})
}

func FcRethrowException(er Exception, logLevel Lvl, format string, args ...interface{}) {
	fcRethrowException(er, logLevel, format, args, fcRethrowExceptionSkip)
}

func fcRethrowException(er Exception, logLevel Lvl, format string, args []interface{}, skip int) {
	er.AppendLog(LogMessage(logLevel, format, args, skip))
	Throw(er)
}

const (
	catchAndLogMessageSkip         = 4
	catchAndfcRethrowExceptionSkip = 5
)

//noinspection GoStructInitializationWithoutFieldNames
func (c *CatchOrFinally) EosRethrowExceptions(exception Exception, format string, args ...interface{}) *CatchOrFinally {
	return c.Catch(func(e ChainExceptions) {
		fcRethrowException(e, LvlWarn, format, args, catchAndfcRethrowExceptionSkip)

	}).Catch(func(e Exception) {
		exception.AppendLog(LogMessage(LvlWarn, format, args, catchAndLogMessageSkip))
		for _, log := range e.GetLog() {
			exception.AppendLog(log)
		}
		Throw(exception)

	}).Catch(func(interface{}) {
		Throw(&UnHandledException{Elog: []Message{LogMessage(LvlWarn, format, args, 4)}})
	}).End()
}

func (c *CatchOrFinally) FcLogAndRethrow() *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(er.DetailMessage())
		fcRethrowException(er, LvlWarn, "rethrow", nil, catchAndfcRethrowExceptionSkip)

	}).Catch(func(a interface{}) {
		e := &UnHandledException{Elog: []Message{LogMessage(LvlWarn, "rethrow %v", []interface{}{a}, catchAndLogMessageSkip)}}
		Warn(e.DetailMessage())
		Throw(e)
	}).End()
}

func (c *CatchOrFinally) FcCaptureLogAndRethrow(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(er.DetailMessage())
		format, arg := FcFormatArgParams(args)
		fcRethrowException(er, LvlWarn, "rethrow "+format, arg, catchAndfcRethrowExceptionSkip)

	}).Catch(func(interface{}) {
		format, arg := FcFormatArgParams(args)
		e := &UnHandledException{Elog: []Message{LogMessage(LvlWarn, "rethrow "+format, arg, catchAndLogMessageSkip)}}
		Warn(e.DetailMessage())
		Throw(e)
	}).End()
}

func (c *CatchOrFinally) FcCaptureAndLog(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(er.DetailMessage())

	}).Catch(func(a interface{}) {
		format, arg := FcFormatArgParams(args)
		e := &UnHandledException{Elog: []Message{LogMessage(LvlWarn, "rethrow "+format, arg, catchAndLogMessageSkip)}}
		Warn(e.DetailMessage())
	}).End()
}

func (c *CatchOrFinally) FcLogAndDrop(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(er.DetailMessage())

	}).Catch(func(a interface{}) {
		format, arg := FcFormatArgParams(args)
		e := &UnHandledException{Elog: []Message{LogMessage(LvlWarn, "rethrow "+format, arg, catchAndLogMessageSkip)}}
		Warn(e.DetailMessage())
	}).End()
}

func (c *CatchOrFinally) FcRethrowExceptions(logLevel Lvl, format string, args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		fcRethrowException(er, logLevel, format, args, catchAndfcRethrowExceptionSkip)

	}).Catch(func(interface{}) {
		e := &UnHandledException{Elog: []Message{LogMessage(logLevel, format, args, catchAndLogMessageSkip)}}
		Throw(e)
	}).End()
}

func (c *CatchOrFinally) FcCaptureAndRethrow(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		format, arg := FcFormatArgParams(args)
		fcRethrowException(er, LvlWarn, format, arg, catchAndfcRethrowExceptionSkip)

	}).Catch(func(interface{}) {
		format, arg := FcFormatArgParams(args)
		e := &UnHandledException{Elog: []Message{LogMessage(LvlWarn, format, arg, catchAndLogMessageSkip)}}
		Throw(e)
	}).End()
}

func FcFormatArgParams(args []interface{}) (string, []interface{}) {
	switch len(args) {
	case 0:
		return "", nil
	case 1:
		return args[0].(string), nil
	default:
		return args[0].(string), args[1:]
	}

}
