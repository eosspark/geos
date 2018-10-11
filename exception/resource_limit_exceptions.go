package exception

type ResourceLimitException struct{ logMessage }

func (e *ResourceLimitException) ChainExceptions()         {}
func (e *ResourceLimitException) ResourceLimitExceptions() {}
func (e *ResourceLimitException) Code() ExcTypes           { return 3210000 }
func (e *ResourceLimitException) What() string {
	return "Resource limit exception"
}
