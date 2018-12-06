package try

import (
	"reflect"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/log"
)

type CatchOrFinally struct {
	e         interface{}
	stackInfo []byte
	//StackTrace []StackInfo
}

//Catch call the exception handler. And return interface CatchOrFinally that
//can call Catch or Finally.
func (c *CatchOrFinally) Catch(f interface{}) (r *CatchOrFinally) {
	///*debug*/s := time.Now().Nanosecond()
	if c == nil || c.e == nil {
		///*debug*/fmt.Println("catch-none", time.Now().Nanosecond() - s, "ns")
		return nil
	}

	switch ft := f.(type) {

	/*
	 * catch exception interface
	 */
	case func(Exception):
		if et, ok := c.e.(Exception); ok {
			ft(et)
			///*debug*/fmt.Println("catch-special", time.Now().Nanosecond() - s, "ns")
			return nil
		}
		///*debug*/fmt.Println("catch-special-fail", time.Now().Nanosecond() - s, "ns")
		return c

	case func(GuardExceptions):
		if et, ok := c.e.(GuardExceptions); ok {
			ft(et)
			return nil
		}
		return c

	/*
	 * catch specific exception
	 */
	case func(*PermissionQueryException):
		if et, ok := c.e.(*PermissionQueryException); ok {
			ft(et)
			return nil
		}
		return c

	case func(*AssertException):
		if et, ok := c.e.(*AssertException); ok {
			ft(et)
			return nil
		}
		return c

	case func(*UnknownBlockException):
		if et, ok := c.e.(*UnknownBlockException); ok {
			ft(et)
			return nil
		}
		return c


	case func(error):
		if et, ok := c.e.(error); ok {
			ft(et)
			return nil
		}
		return c

	case func(interface{}):
		ft(c.e)
		return nil
	}


	// make sure all panic can be caught
	rf := reflect.ValueOf(f)
	ft := rf.Type()
	if ft.NumIn() > 0 {
		it := ft.In(0)
		ct := reflect.TypeOf(c.e)

		its, cts := it.String(), ct.String()

		if its == cts || (it.Kind() == reflect.Interface && ct.Implements(it)) {
			reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(c.e)})
			///*debug*/fmt.Println("catch-reflect", time.Now().Nanosecond() - s, "ns")
			return nil

		} else if ct.Kind() == reflect.Ptr && cts[1:] == its { // make pointer can be caught by its value type
			reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(reflect.ValueOf(c.e).Elem().Interface())})
			return nil

		}
		//else if cts == "runtime.errorString" && its == "try.RuntimeError" {
		//	var rte RuntimeError
		//	rte.Message = c.e.(error).Error()
		//	rte.stackInfo = c.stackInfo
		//	ev := reflect.ValueOf(rte)
		//	reflect.ValueOf(f).Call([]reflect.Value{ev})
		//	return nil
		//}

		//println(it.String(), ct.String())

	}

	///*debug*/fmt.Println("catch-fail", time.Now().Nanosecond() - s, "ns")
	return c
}

//Necessary to call at the end of try-catch block, to ensure panic uncaught exceptions
func (c *CatchOrFinally) End() {
	if c != nil && c.e != nil {
		if DEBUG {
			c.printStackInfo()
		}
		Throw(c.e)
	}
}

func (c *CatchOrFinally) printStackInfo() {
	switch e := c.e.(type) {
	case Exception:
		Error(GetDetailMessage(e))
	case error:
		Error(e.Error())
	}

	Error(string(c.stackInfo))
}

func (c *CatchOrFinally) CatchAndCall(Next func(interface{})) *CatchOrFinally {
	return c.Catch(func(err Exception) {
		Next(err)

	}).Catch(func(e error) {
		fce := &FcException{ELog: NewELog(FcLogMessage(LvlWarn, "rethrow %s: ", e.Error()))}
		Next(fce)

	}).Catch(func(interface{}) {
		e := &UnHandledException{ELog: NewELog(FcLogMessage(LvlWarn, "rethrow"))}
		Next(e)
	})
}

//Finally always be called if defined.
//func (c *CatchOrFinally) Finally(f interface{}) (r *OrThrowable) {
//	reflect.ValueOf(f).Call([]reflect.Value{})
//	if c == nil || c.e == nil {
//		return nil
//	}
//	return &OrThrowable{c.e}
//}

//OrThrow throw error then never catch block entered.

//OrThrow throw error then never catch block entered.
//func (c *OrThrowable) End() {
//	if c != nil && c.e != nil {
//		Throw(c.e)
//	}
//}

//Throw is wrapper of panic().
