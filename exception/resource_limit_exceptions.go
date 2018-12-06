package exception

import _ "github.com/eosspark/eos-go/log"

type ResourceLimitException struct{ ELog }

func (ResourceLimitException) ChainExceptions()         {}
func (ResourceLimitException) ResourceLimitExceptions() {}
func (ResourceLimitException) Code() ExcTypes           { return 3210000 }
func (ResourceLimitException) What() string {
	return "Resource limit exception"
}
