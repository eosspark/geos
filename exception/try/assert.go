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
	Throw(&FcException{Elog: []Message{FcLogMessage(LvlError, format, args...)}})
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

		//replaced by StdException
		//}).Catch(func(e error) {
		//	exception.AppendLog(FcLogMessage(LvlWarn, fmt.Sprintf("%s (%s)", format, e.Error()), args...))
		//	Throw(exception)

	}).Catch(func(interface{}) {
		Throw(&UnHandledException{Elog: []Message{FcLogMessage(LvlWarn, format, args...)}})
	}).End()
}

func (c *CatchOrFinally) FcLogAndRethrow() *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(er.DetailMessage())
		FcRethrowException(er, LvlWarn, "rethrow")

		//replaced by StdException
		//}).Catch(func(e error) {
		//	fce := &FcException{Elog: []Message{FcLogMessage(LvlWarn, "rethrow: %s", e.Error())}}
		//	Warn(fce.DetailMessage())
		//	Throw(fce)

	}).Catch(func(a interface{}) {
		e := &UnHandledException{Elog: []Message{FcLogMessage(LvlWarn, "rethrow %v", a)}}
		Warn(e.DetailMessage())
		Throw(e)
	}).End()
}

func (c *CatchOrFinally) FcCaptureLogAndRethrow(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(er.DetailMessage())
		format, arg := FcFormatArgParams(args)
		FcRethrowException(er, LvlWarn, "rethrow "+format, arg...)

		//replaced by StdException
		//}).Catch(func(e error) {
		//	format, arg := FcFormatArgParams(args)
		//	fce := &FcException{Elog: []Message{FcLogMessage(LvlWarn, fmt.Sprintf("rethrow %s %s", e.Error(), format), arg...)}}
		//	Warn(fce.DetailMessage())
		//	Throw(fce)

	}).Catch(func(interface{}) {
		format, arg := FcFormatArgParams(args)
		e := &UnHandledException{Elog: []Message{FcLogMessage(LvlWarn, "rethrow "+format, arg...)}}
		Warn(e.DetailMessage())
		Throw(e)
	}).End()
}

func (c *CatchOrFinally) FcCaptureAndLog(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(er.DetailMessage())

		//replaced by StdException
		//}).Catch(func(e error) {
		//	format, arg := FcFormatArgParams(args)
		//	fce := &FcException{Elog: []Message{FcLogMessage(LvlWarn, fmt.Sprintf("rethrow %s: %s", e.Error(), format), arg...)}}
		//	Warn(fce.DetailMessage())

	}).Catch(func(a interface{}) {
		format, arg := FcFormatArgParams(args)
		e := &UnHandledException{Elog: []Message{FcLogMessage(LvlWarn, "rethrow "+format, arg...)}}
		Warn(e.DetailMessage())
	}).End()
}

func (c *CatchOrFinally) FcLogAndDrop(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		Warn(er.DetailMessage())

		//replaced by StdException
		//}).Catch(func(e error) {
		//	format, arg := FcFormatArgParams(args)
		//	fce := &FcException{Elog: []Message{FcLogMessage(LvlWarn, fmt.Sprintf("rethrow %s: %s", e.Error(), format), arg...)}}
		//	Warn(fce.DetailMessage())

	}).Catch(func(a interface{}) {
		format, arg := FcFormatArgParams(args)
		e := &UnHandledException{Elog: []Message{FcLogMessage(LvlWarn, "rethrow "+format, arg...)}}
		Warn(e.DetailMessage())
	}).End()
}

func (c *CatchOrFinally) FcRethrowExceptions(logLevel Lvl, format string, args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		FcRethrowException(er, logLevel, format, args...)

		//replaced by StdException
		//}).Catch(func(e error) {
		//	fce := &FcException{Elog: []Message{FcLogMessage(logLevel, fmt.Sprintf("%s: %s", e.Error(), format), args...)}}
		//	Throw(fce)

	}).Catch(func(interface{}) {
		e := &UnHandledException{Elog: []Message{FcLogMessage(logLevel, format, args...)}}
		Throw(e)
	}).End()
}

//noinspection ALL
func (c *CatchOrFinally) FcCaptureAndRethrow(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		format, arg := FcFormatArgParams(args)
		FcRethrowException(er, LvlWarn, format, arg...)

		//replaced by StdException
		//}).Catch(func(e error) {
		//	format, arg := FcFormatArgParams(args)
		//	fce := &FcException{Elog: []Message{FcLogMessage(LvlWarn, fmt.Sprintf("%s: %s", e.Error(), format), arg...)}}
		//	Throw(fce)

	}).Catch(func(interface{}) {
		format, arg := FcFormatArgParams(args)
		e := &UnHandledException{Elog: []Message{FcLogMessage(LvlWarn, format, arg...)}}
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
