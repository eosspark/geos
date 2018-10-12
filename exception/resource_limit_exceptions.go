package exception

type ResourceLimitException struct{ logMessage }

func (ResourceLimitException) ChainExceptions()         {}
func (ResourceLimitException) ResourceLimitExceptions() {}
func (ResourceLimitException) Code() ExcTypes           { return 3210000 }
func (ResourceLimitException) What() string {
	return "Resource limit exception"
}
