package gosio

type Reactor interface {
	run()
	push(op interface{}, args ...interface{})
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