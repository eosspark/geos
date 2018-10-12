package exception

type ResourceExhaustedException struct{ logMessage }

func (ResourceExhaustedException) ChainExceptions()             {}
func (ResourceExhaustedException) ResourceExhaustedExceptions() {}
func (ResourceExhaustedException) Code() ExcTypes               { return 3080000 }
func (ResourceExhaustedException) What() string                 { return "Resource exhausted exception" }

type RamUsageExceeded struct{ logMessage }

func (RamUsageExceeded) ChainExceptions()             {}
func (RamUsageExceeded) ResourceExhaustedExceptions() {}
func (RamUsageExceeded) Code() ExcTypes               { return 3080001 }
func (RamUsageExceeded) What() string                 { return "Account using more than allotted RAM usage" }

type TxNetUsageExceeded struct{ logMessage }

func (TxNetUsageExceeded) ChainExceptions()             {}
func (TxNetUsageExceeded) ResourceExhaustedExceptions() {}
func (TxNetUsageExceeded) Code() ExcTypes               { return 3080002 }
func (TxNetUsageExceeded) What() string {
	return "Transaction exceeded the current network usage limit imposed on the transaction"
}

type BlockNetUsageExceeded struct{ logMessage }

func (BlockNetUsageExceeded) ChainExceptions()             {}
func (BlockNetUsageExceeded) ResourceExhaustedExceptions() {}
func (BlockNetUsageExceeded) Code() ExcTypes               { return 3080003 }
func (BlockNetUsageExceeded) What() string {
	return "Transaction network usage is too much for the remaining allowable usage of the current block"
}

type TxCpuUsageExceed struct{ logMessage }

func (TxCpuUsageExceed) ChainExceptions()             {}
func (TxCpuUsageExceed) ResourceExhaustedExceptions() {}
func (TxCpuUsageExceed) Code() ExcTypes               { return 3080004 }
func (TxCpuUsageExceed) What() string {
	return "Transaction exceeded the current CPU usage limit imposed on the transaction"
}

type BlockCpuUsageExceeded struct{ logMessage }

func (BlockCpuUsageExceeded) ChainExceptions()             {}
func (BlockCpuUsageExceeded) ResourceExhaustedExceptions() {}
func (BlockCpuUsageExceeded) Code() ExcTypes               { return 3080005 }
func (BlockCpuUsageExceeded) What() string {
	return "Transaction CPU usage is too much for the remaining allowable usage of the current block"
}

type DeadlineException struct{ logMessage }

func (DeadlineException) ChainExceptions()             {}
func (DeadlineException) ResourceExhaustedExceptions() {}
func (DeadlineException) Code() ExcTypes               { return 3080006 }
func (DeadlineException) What() string {
	return "Transaction took too long"
}

type GreylistNetUsageExceeded struct{ logMessage }

func (GreylistNetUsageExceeded) ChainExceptions()             {}
func (GreylistNetUsageExceeded) ResourceExhaustedExceptions() {}
func (GreylistNetUsageExceeded) Code() ExcTypes               { return 3080007 }
func (GreylistNetUsageExceeded) What() string {
	return "Transaction exceeded the current greylisted account network usage limit"
}

type GreylistCpuUsageExceeded struct{ logMessage }

func (GreylistCpuUsageExceeded) ChainExceptions()             {}
func (GreylistCpuUsageExceeded) ResourceExhaustedExceptions() {}
func (GreylistCpuUsageExceeded) Code() ExcTypes               { return 3080008 }
func (GreylistCpuUsageExceeded) What() string {
	return "Transaction exceeded the current greylisted account CPU usage limit"
}

// implements DeadlineExceptions
type LeewayDeadlineException struct{ logMessage }

func (LeewayDeadlineException) ChainExceptions()             {}
func (LeewayDeadlineException) ResourceExhaustedExceptions() {}
func (LeewayDeadlineException) DeadlineExceptions()          {}
func (LeewayDeadlineException) Code() ExcTypes               { return 3081001 }
func (LeewayDeadlineException) What() string {
	return "Transaction reached the deadline set due to leeway on account CPU limits"
}
