package exception

type ProducerException struct{ logMessage }

func (e *ProducerException) ChainExceptions()    {}
func (e *ProducerException) ProducerExceptions() {}
func (e *ProducerException) Code() ExcTypes      { return 3170000 }
func (e *ProducerException) What() string {
	return "Producer exception"
}

type ProducerPrivKeyNotFound struct{ logMessage }

func (e *ProducerPrivKeyNotFound) ChainExceptions()    {}
func (e *ProducerPrivKeyNotFound) ProducerExceptions() {}
func (e *ProducerPrivKeyNotFound) Code() ExcTypes      { return 3170001 }
func (e *ProducerPrivKeyNotFound) What() string {
	return "Producer private key is not available"
}

type MissingPendingBlockState struct{ logMessage }

func (e *MissingPendingBlockState) ChainExceptions()    {}
func (e *MissingPendingBlockState) ProducerExceptions() {}
func (e *MissingPendingBlockState) Code() ExcTypes      { return 3170002 }
func (e *MissingPendingBlockState) What() string {
	return "Pending block state is missing"
}

type ProducerDoubleConfirm struct{ logMessage }

func (e *ProducerDoubleConfirm) ChainExceptions()    {}
func (e *ProducerDoubleConfirm) ProducerExceptions() {}
func (e *ProducerDoubleConfirm) Code() ExcTypes      { return 3170003 }
func (e *ProducerDoubleConfirm) What() string {
	return "Producer is double confirming known range"
}

type ProducerScheduleException struct{ logMessage }

func (e *ProducerScheduleException) ChainExceptions()    {}
func (e *ProducerScheduleException) ProducerExceptions() {}
func (e *ProducerScheduleException) Code() ExcTypes      { return 3170004 }
func (e *ProducerScheduleException) What() string {
	return "Producer schedule exception"
}

type ProducerNotInSchedule struct{ logMessage }

func (e *ProducerNotInSchedule) ChainExceptions()    {}
func (e *ProducerNotInSchedule) ProducerExceptions() {}
func (e *ProducerNotInSchedule) Code() ExcTypes      { return 3170006 }
func (e *ProducerNotInSchedule) What() string {
	return "The producer is not part of current schedule"
}
