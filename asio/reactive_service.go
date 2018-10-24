package asio

import "os"

type Reactor interface {
	run()
	stop()
	push(op interface{}, args ...interface{})
	notify(sig ...os.Signal)
}

func (i *IoContext) GetService () Reactor {
	if i.reactor == nil {
		i.reactor = NewGouroutineReactor()
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