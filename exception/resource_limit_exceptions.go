package exception

import . "github.com/eosspark/eos-go/log"

type ResourceLimitException struct{ LogMessage }

func (ResourceLimitException) ChainExceptions()         {}
func (ResourceLimitException) ResourceLimitExceptions() {}
func (ResourceLimitException) Code() ExcTypes           { return 3210000 }
func (ResourceLimitException) What() string {
	return "Resource limit exception"
}
