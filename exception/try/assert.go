package try

import (
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/log"
	"fmt"
	"os"
)

func Assert(expr bool, message string) {
	if !expr {
		println(message)
		os.Exit(1)
	}
}

func EosAssert(expr bool, exception Exception, format string, args ...interface{}) {
	if !expr {
		FcThrowException(exception, format, args...)
	}
}

func FcAssert(test bool, args ...interface{}) {
	if !test {
		FcThrowException(&AssertException{}, "assert:", args...)
	}
}

func EosThrow(exception Exception, format string, args ...interface{}) {
	exception.AppendLog(FcLogMessage(LvlError, format, args...))
	Throw(exception)
}

func FcThrow(format string, args ...interface{}) {
	Throw(&FcException{ELog: NewELog(FcLogMessage(LvlError, format, args...))})
}

func FcCaptureAndThrow(exception Exception, format string, args ...interface{}) {
	exception.AppendLog(FcLogMessage(LvlError, format, args...))
	Throw(exception)
}

func FcThrowException(exception Exception, format string, args ...interface{}) {
	exception.AppendLog(FcLogMessage(LvlError, format, args...))
	Throw(exception)
}

func FcRethrowException(er Exception, logLevel Lvl, format string, args ...interface{}) {
	er.AppendLog(FcLogMessage(logLevel, format, args...))
	Throw(er)
}

//noinspection GoStructInitializationWithoutFieldNames
func (c *CatchOrFinally) EosRethrowExceptions(exception Exception, format string, args ...interface{}) *CatchOrFinally {
	return c.Catch(func(e ChainExceptions) {
		FcRethrowException(e, LvlWarn, format, args...)

	}).Catch(func(e Exception) {
		exception.AppendLog(FcLogMessage(LvlWarn, format, args...))
		for _, log := range e.GetLog() {
			exception.AppendLog(log)
		}
		Throw(exception)

	}).Catch(func(e error) {
		exception.AppendLog(FcLogMessage(LvlWarn, fmt.Sprintf("%s (%s)", format, e.Error()), args...))
		Throw(exception)

	}).Catch(func(interface{}) {
		Throw(&UnHandledException{NewELog(FcLogMessage(LvlWarn, format, args...))})
	})
}

func (c *CatchOrFinally) FcLogAndRethrow() *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(GetDetailMessage(er))
		FcRethrowException(er, LvlWarn, "rethrow")

	}).Catch(func(e error) {
		fce := &FcException{ELog: NewELog(FcLogMessage(LvlWarn, "rethrow: %s", e.Error()))}
		Warn(GetDetailMessage(fce))
		Throw(fce)

	}).Catch(func(a interface{}) {
		e := &UnHandledException{ELog: NewELog(FcLogMessage(LvlWarn, "rethrow %v", a))}
		Warn(GetDetailMessage(e))
		Throw(e)
	})
}

func (c *CatchOrFinally) FcCaptureLogAndRethrow(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(GetDetailMessage(er))
		format, arg := FcFormatArgParams(args)
		FcRethrowException(er, LvlWarn, "rethrow "+format, arg...)

	}).Catch(func(e error) {
		format, arg := FcFormatArgParams(args)
		fce := &FcException{ELog: NewELog(FcLogMessage(LvlWarn, fmt.Sprintf("rethrow %s %s", e.Error(), format), arg...))}
		Warn(GetDetailMessage(fce))
		Throw(fce)

	}).Catch(func(interface{}) {
		format, arg := FcFormatArgParams(args)
		e := &UnHandledException{ELog: NewELog(FcLogMessage(LvlWarn, "rethrow "+format, arg...))}
		Warn(GetDetailMessage(e))
		Throw(e)
	})
}

func (c *CatchOrFinally) FcCaptureAndLog(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(GetDetailMessage(er))

	}).Catch(func(e error) {
		format, arg := FcFormatArgParams(args)
		fce := &FcException{ELog: NewELog(FcLogMessage(LvlWarn, fmt.Sprintf("rethrow %s: %s", e.Error(), format), arg...))}
		Warn(GetDetailMessage(fce))

	}).Catch(func(a interface{}) {
		format, arg := FcFormatArgParams(args)
		e := &UnHandledException{ELog: NewELog(FcLogMessage(LvlWarn, "rethrow "+format, arg...))}
		Warn(GetDetailMessage(e))
	})
}

func (c *CatchOrFinally) FcLogAndDrop(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(GetDetailMessage(er))

	}).Catch(func(e error) {
		format, arg := FcFormatArgParams(args)
		fce := &FcException{ELog: NewELog(FcLogMessage(LvlWarn, fmt.Sprintf("rethrow %s: %s", e.Error(), format), arg...))}
		Warn(GetDetailMessage(fce))

	}).Catch(func(a interface{}) {
		format, arg := FcFormatArgParams(args)
		e := &UnHandledException{ELog: NewELog(FcLogMessage(LvlWarn, "rethrow "+format, arg...))}
		Warn(GetDetailMessage(e))
	})
}

func (c *CatchOrFinally) FcRethrowExceptions(logLevel Lvl, format string, args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		FcRethrowException(er, logLevel, format, args...)

	}).Catch(func(e error) {
		fce := &FcException{ELog: NewELog(FcLogMessage(logLevel, fmt.Sprintf("%s: %s", e.Error(), format), args...))}
		Throw(fce)

	}).Catch(func(interface{}) {
		e := &UnHandledException{ELog: NewELog(FcLogMessage(logLevel, format, args...))}
		Throw(e)
	})
}

//noinspection ALL
func (c *CatchOrFinally) FcCaptureAndRethrow(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		format, arg := FcFormatArgParams(args)
		FcRethrowException(er, LvlWarn, format, arg...)

	}).Catch(func(e error) {
		format, arg := FcFormatArgParams(args)
		fce := &FcException{NewELog(FcLogMessage(LvlWarn, fmt.Sprintf("%s: %s", e.Error(), format), arg...))}
		Throw(fce)

	}).Catch(func(interface{}) {
		format, arg := FcFormatArgParams(args)
		e := &UnHandledException{NewELog(FcLogMessage(LvlWarn, format, arg...))}
		Throw(e)
	})
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
