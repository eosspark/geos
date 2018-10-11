package exception

type ResourceExhaustedException struct{ logMessage }

func (e *ResourceExhaustedException) ChainExceptions()             {}
func (e *ResourceExhaustedException) ResourceExhaustedExceptions() {}
func (e *ResourceExhaustedException) Code() ExcTypes               { return 3080000 }
func (e *ResourceExhaustedException) What() string                 { return "Resource exhausted exception" }

type TxCpuUsageExceed struct{ logMessage }

func (e *TxCpuUsageExceed) ChainExceptions()             {}
func (e *TxCpuUsageExceed) ResourceExhaustedExceptions() {}
func (e *TxCpuUsageExceed) Code() ExcTypes               { return 3080004 }
func (e *TxCpuUsageExceed) What() string {
	return "Transaction exceeded the current CPU usage limit imposed on the transaction"
}

type BlockCpuUsageExceeded struct{ logMessage }

func (e *BlockCpuUsageExceeded) ChainExceptions()             {}
func (e *BlockCpuUsageExceeded) ResourceExhaustedExceptions() {}
func (e *BlockCpuUsageExceeded) Code() ExcTypes               { return 3080005 }
func (e *BlockCpuUsageExceeded) What() string {
	return "Transaction CPU usage is too much for the remaining allowable usage of the current block"
}

type DeadlineException struct{ logMessage }

func (e *DeadlineException) ChainExceptions()             {}
func (e *DeadlineException) ResourceExhaustedExceptions() {}
func (e *DeadlineException) Code() ExcTypes               { return 3080006 }
func (e *DeadlineException) What() string {
	return "Transaction took too long"
}

type GreylistNetUsageExceeded struct{ logMessage }

func (e *GreylistNetUsageExceeded) ChainExceptions()             {}
func (e *GreylistNetUsageExceeded) ResourceExhaustedExceptions() {}
func (e *GreylistNetUsageExceeded) Code() ExcTypes               { return 3080007 }
func (e *GreylistNetUsageExceeded) What() string {
	return "Transaction exceeded the current greylisted account network usage limit"
}

type GreylistCpuUsageExceeded struct{ logMessage }

func (e *GreylistCpuUsageExceeded) ChainExceptions()             {}
func (e *GreylistCpuUsageExceeded) ResourceExhaustedExceptions() {}
func (e *GreylistCpuUsageExceeded) Code() ExcTypes               { return 3080008 }
func (e *GreylistCpuUsageExceeded) What() string {
	return "Transaction exceeded the current greylisted account CPU usage limit"
}

type LeewayDeadlineException struct{ logMessage }

func (e *LeewayDeadlineException) ChainExceptions()             {}
func (e *LeewayDeadlineException) ResourceExhaustedExceptions() {}
func (e *LeewayDeadlineException) DeadlineExceptions()          {}
func (e *LeewayDeadlineException) Code() ExcTypes               { return 3081001 }
func (e *LeewayDeadlineException) What() string {
	return "Transaction reached the deadline set due to leeway on account CPU limits"
}
