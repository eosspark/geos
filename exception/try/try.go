package try

import (
	"reflect"
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


func Throw(e interface{}) {
	panic(e)
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



