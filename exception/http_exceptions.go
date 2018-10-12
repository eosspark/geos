package exception

type HttpException struct{ logMessage }

func (e *HttpException) ChainExceptions() {}
func (e *HttpException) HttpExceptions()  {}
func (e *HttpException) Code() ExcTypes   { return 3200000 }
func (e *HttpException) What() string {
	return "http exception"
}
