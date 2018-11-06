package try

import (
	. "github.com/eosspark/eos-go/exception"
	"fmt"
)


func FcAssert(test bool, args ...interface{}) {
	if !test {
		FcThrowException(&AssertException{}, "assert:", args...)
	}
}

func FcCaptureAndThrow(exception Exception, args ...interface{}) {
	Throw(exception)
	//TODO log message
}

func FcThrow(args ...interface{}) {
	Throw(&FcException{})
	//TODO log message
}

func FcThrowException(exception Exception, format string, args ...interface{}) {
	Throw(exception)
	//TODO log message
}

func FcRethrowException(exception Exception, logLevel string, format string, args ...interface{}) {
	//TODO appends a log message
	Throw(exception)
}

func (c *CatchOrFinally) FcLogAndRethrow() *CatchOrFinally {
	return c.Catch(func(er Exception) {
		//TODO: log
		FcRethrowException(er, "warn", "rethrow")
	}).Catch(func(e error) {
		//TODO log message
		fce := &FcException{}
		//TODO: log
		Throw(fce)
	}).Catch(func(interface{}) {
		//TODO log message
		e := UnHandledException{}
		//TODO: log
		Throw(e)
	})
}

func (c *CatchOrFinally) FcCaptureLogAndRethrow(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		//TODO wlog
		FcRethrowException(er, "warn", "rethrow", args...)
	}).Catch(func(e error) {
		//TODO log message
		fce := &FcException{}
		//TODO wlog
		Throw(fce)
	}).Catch(func(interface{}) {
		//TODO log message
		e := &UnHandledException{}
		//TODO wlog
		Throw(e)
	})
}

func (c *CatchOrFinally) FcCaptureAndLog(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		//TODO wlog
	}).Catch(func(e error) {
		//TODO log message
		fce := &FcException{}
		fmt.Println(fce.What()) //TODO wlog
	}).Catch(func(interface{}) {
		//TODO log message
		e := &UnHandledException{}
		fmt.Println(e.What()) //TODO wlog
	})
}

func (c *CatchOrFinally) FcLogAndDrop(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		//TODO wlog
		fmt.Println(er.Message())
	}).Catch(func(e error) {
		//TODO log message
		fce := &FcException{}
		fmt.Println(fce.What()) //TODO wlog
	}).Catch(func(interface{}) {
		//TODO log message
		e := &UnHandledException{}
		fmt.Println(e.What()) //TODO wlog
	})
}

func (c *CatchOrFinally) FcRethrowExceptions(logLevel string, format string, args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception){
		FcRethrowException(er, logLevel, format, args...)
	}).Catch(func(e error) {
		//TODO log message
		fce := &FcException{}
		Throw(fce)
	}).Catch(func(interface{}) {
		//TODO log message
		Throw(&UnHandledException{})
	})
}

func (c *CatchOrFinally) FcCaptureAndRethrow(args ...interface{}) *CatchOrFinally {
	return c.Catch(func(er Exception) {
		FcRethrowException(er, "warn", "", args...)
	}).Catch(func(e error) {
		//TODO log message
		fce := &FcException{}
		Throw(fce)
	}).Catch(func(interface{}) {
		//TODO log message
		Throw(&UnHandledException{})
	})
}
