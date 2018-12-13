package try

import (
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/log"
	"runtime"
)

//StackInfo store code informations when catched exception.

const (
	DEBUG = false
)

var (
	stackInfo []byte = nil
	stackSize        = 65536
)

//RuntimeError is wrapper of runtime.errorString and stacktrace.
type RuntimeError struct {
	Message   string
	stackInfo []byte
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
	///*debug*/s := time.Now().Nanosecond()
	defer func() {
		if e := recover(); e != nil {

			r = &CatchOrFinally{e}

			if DEBUG {
				stackInfo = make([]byte, stackSize)
				stackInfo = stackInfo[:runtime.Stack(stackInfo, false)]
				printStackInfo(e)
			}

		}
	}()

	f()
	return nil
}

func Throw(e interface{}) {
	if e == nil {
		return
	}
	panic(e)
}

//Use defer HandleStackInfo() before main func panic
func printStackInfo(errorPtr interface{}) {
	if DEBUG && stackInfo != nil {
		switch e := errorPtr.(type) {
		case exception.Exception:
			log.Warn("%s: %s", exception.GetDetailMessage(e), string(stackInfo))
		case error:
			log.Warn("error %s: %s", e.Error(), string(stackInfo))
		default:
			log.Warn("panic %#v: %s", e, string(stackInfo))
		}
	}
}

//type returnTypes struct{}

//Just use in try-catch block, you should update return-value before call it
//Deprecated: use returning flag instead
//func Return() {
//	panic(returnTypes{})
//}

//Use defer HandleReturn() before try-catch block when the block includes Return function
//Deprecated: use returning flag instead
//func HandleReturn() {
//	if rv := recover(); rv != nil {
//		if _, ok := rv.(returnTypes); !ok {
//			panic(rv)
//		}
//	}
//}
