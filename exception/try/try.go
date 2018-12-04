package try

import (
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

const DEBUG = true

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

			if rt, ok := e.(returnTypes); ok {
				panic(rt)
			}

			r = &CatchOrFinally{}
			r.e = e

			if DEBUG {
				const size = 65536
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]

				r.stackInfo = buf
			}

		}
///*debug*/fmt.Println("try", time.Now().Nanosecond() - s, "ns")
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

type returnTypes struct{}

//Just use in try-catch block, you should update return-value before call it
//Deprecated: use returning flag instead
func Return() {
	panic(returnTypes{})
}

//Use defer HandleReturn() before try-catch block when the block includes Return function
//Deprecated: use returning flag instead
func HandleReturn() {
	if rv := recover(); rv != nil {
		if _, ok := rv.(returnTypes); !ok {
			panic(rv)
		}
	}
}



