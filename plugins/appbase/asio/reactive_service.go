package asio

import "os"

type Reactor interface {
	run()
	stop()
	post(op interface{}, args ...interface{})
	notify(sig ...os.Signal)
}

// use GoroutineReactor for each operating system
func (i *IoContext) GetService () Reactor {
	if i.reactor == nil {
		i.reactor = NewGoroutineReactor()
	}
	return i.reactor
}

type ErrorCode struct {
	Valid bool
	Error error
}

func NewErrorCode(err error) ErrorCode {
	ec := ErrorCode{}
	if err != nil {
		ec.Valid = true
		ec.Error = err
	}
	return ec
}