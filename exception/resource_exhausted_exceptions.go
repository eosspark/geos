package exception

import _ "github.com/eosspark/eos-go/log"

type ResourceExhaustedException struct{ ELog }

func (ResourceExhaustedException) ChainExceptions()             {}
func (ResourceExhaustedException) ResourceExhaustedExceptions() {}
func (ResourceExhaustedException) Code() ExcTypes               { return 3080000 }
func (ResourceExhaustedException) What() string                 { return "Resource exhausted exception" }

type RamUsageExceeded struct{ ELog }

func (RamUsageExceeded) ChainExceptions()             {}
func (RamUsageExceeded) ResourceExhaustedExceptions() {}
func (RamUsageExceeded) Code() ExcTypes               { return 3080001 }
func (RamUsageExceeded) What() string                 { return "Account using more than allotted RAM usage" }

type TxNetUsageExceeded struct{ ELog }

func (TxNetUsageExceeded) ChainExceptions()             {}
func (TxNetUsageExceeded) ResourceExhaustedExceptions() {}
func (TxNetUsageExceeded) Code() ExcTypes               { return 3080002 }
func (TxNetUsageExceeded) What() string {
	return "Transaction exceeded the current network usage limit imposed on the transaction"
}

type BlockNetUsageExceeded struct{ ELog }

func (BlockNetUsageExceeded) ChainExceptions()             {}
func (BlockNetUsageExceeded) ResourceExhaustedExceptions() {}
func (BlockNetUsageExceeded) Code() ExcTypes               { return 3080003 }
func (BlockNetUsageExceeded) What() string {
	return "Transaction network usage is too much for the remaining allowable usage of the current block"
}

type TxCpuUsageExceeded struct{ ELog }

func (TxCpuUsageExceeded) ChainExceptions()             {}
func (TxCpuUsageExceeded) ResourceExhaustedExceptions() {}
func (TxCpuUsageExceeded) Code() ExcTypes               { return 3080004 }
func (TxCpuUsageExceeded) What() string {
	return "Transaction exceeded the current CPU usage limit imposed on the transaction"
}

type BlockCpuUsageExceeded struct{ ELog }

func (BlockCpuUsageExceeded) ChainExceptions()             {}
func (BlockCpuUsageExceeded) ResourceExhaustedExceptions() {}
func (BlockCpuUsageExceeded) Code() ExcTypes               { return 3080005 }
func (BlockCpuUsageExceeded) What() string {
	return "Transaction CPU usage is too much for the remaining allowable usage of the current block"
}

type DeadlineException struct{ ELog }

func (DeadlineException) ChainExceptions()             {}
func (DeadlineException) ResourceExhaustedExceptions() {}
func (DeadlineException) Code() ExcTypes               { return 3080006 }
func (DeadlineException) What() string {
	return "Transaction took too long"
}

type GreylistNetUsageExceeded struct{ ELog }

func (GreylistNetUsageExceeded) ChainExceptions()             {}
func (GreylistNetUsageExceeded) ResourceExhaustedExceptions() {}
func (GreylistNetUsageExceeded) Code() ExcTypes               { return 3080007 }
func (GreylistNetUsageExceeded) What() string {
	return "Transaction exceeded the current greylisted account network usage limit"
}

type GreylistCpuUsageExceeded struct{ ELog }

func (GreylistCpuUsageExceeded) ChainExceptions()             {}
func (GreylistCpuUsageExceeded) ResourceExhaustedExceptions() {}
func (GreylistCpuUsageExceeded) Code() ExcTypes               { return 3080008 }
func (GreylistCpuUsageExceeded) What() string {
	return "Transaction exceeded the current greylisted account CPU usage limit"
}

// implements DeadlineExceptions
type LeewayDeadlineException struct{ ELog }

func (LeewayDeadlineException) ChainExceptions()             {}
func (LeewayDeadlineException) ResourceExhaustedExceptions() {}
func (LeewayDeadlineException) DeadlineExceptions()          {}
func (LeewayDeadlineException) Code() ExcTypes               { return 3081001 }
func (LeewayDeadlineException) What() string {
	return "Transaction reached the deadline set due to leeway on account CPU limits"
}
