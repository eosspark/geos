package try

import (
	"reflect"
		"fmt"
	"runtime"
	)

//StackInfo store code informations when catched exception.
type StackInfo struct {
	PC   uintptr
	File string
	Line int
}

//RuntimeError is wrapper of runtime.errorString and stacktrace.
type RuntimeError struct {
	Message    string
	stackInfo  []byte
	//StackTrace []StackInfo
}

func (rte RuntimeError) String() string {
	return rte.Message
}

type CatchOrFinally struct {
	e 		  interface{}
	stackInfo []byte
	//StackTrace []StackInfo
}

type OrThrowable struct {
	e interface{}
}

//Try call the function. And return interface that can call Catch or Finally.
func Try(f func()) (r *CatchOrFinally) {
	defer func() {
		if e := recover(); e != nil {

			if rt, ok := e.(returnTypes); ok {
				panic(rt)
			}

			r = &CatchOrFinally{}
			r.e = e

			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]

			r.stackInfo = buf

			//i := 1
			//for {
			//	if p, f, l, o := runtime.Caller(i); o {
			//		r.StackTrace = append(r.StackTrace, StackInfo{p, f, l})
			//		i++
			//	} else {
			//		break
			//	}
			//}
		}
	}()
	reflect.ValueOf(f).Call([]reflect.Value{})
	return nil
}

//Catch call the exception handler. And return interface CatchOrFinally that
//can call Catch or Finally.
func (c *CatchOrFinally) Catch(f interface{}) (r *CatchOrFinally) {
	if c == nil || c.e == nil {
		return nil
	}

	rf := reflect.ValueOf(f)
	ft := rf.Type()
	if ft.NumIn() > 0 {
		it := ft.In(0)
		ct := reflect.TypeOf(c.e)

		its, cts := it.String(), ct.String()

		if its == cts || (it.Kind() == reflect.Interface && ct.Implements(it)) {
			reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(c.e)})
			return nil

		} else if ct.Kind() == reflect.Ptr && cts[1:] == its { // make pointer can be caught by its value type
			reflect.ValueOf(f).Call([]reflect.Value{reflect.ValueOf(reflect.ValueOf(c.e).Elem().Interface())})
			return nil

		} else if cts == "runtime.errorString" && its == "try.RuntimeError" {
			var rte RuntimeError
			rte.Message = c.e.(error).Error()
			rte.stackInfo = c.stackInfo
			ev := reflect.ValueOf(rte)
			reflect.ValueOf(f).Call([]reflect.Value{ev})
			return nil
		}

		//println(it.String(), ct.String())

	}
	return c
}

func Throw(e interface{}) {
	panic(e)
}

//Necessary to call at the end of try-catch block, to ensure panic uncaught exceptions
func (c *CatchOrFinally) End() {
	if c != nil && c.e != nil {
		c.printStackInfo()
		Throw(c.e)
	}
}

func (c *CatchOrFinally) printStackInfo() {
	fmt.Println(string(c.stackInfo))
}

type returnTypes struct{}

//Just use in try-catch block, you should update return-value before call it
func Return() {
	panic(returnTypes{})
}

//Use defer HandleReturn() before try-catch block when the block includes Return function
func HandleReturn() {
	if rv := recover(); rv != nil {
		if _, ok := rv.(returnTypes); !ok {
			panic(rv)
		}
	}
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

