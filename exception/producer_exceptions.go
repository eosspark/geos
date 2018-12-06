package exception

import _ "github.com/eosspark/eos-go/log"

type ProducerException struct{ ELog }

func (ProducerException) ChainExceptions()    {}
func (ProducerException) ProducerExceptions() {}
func (ProducerException) Code() ExcTypes      { return 3170000 }
func (ProducerException) What() string {
	return "Producer exception"
}

type ProducerPrivKeyNotFound struct{ ELog }

func (ProducerPrivKeyNotFound) ChainExceptions()    {}
func (ProducerPrivKeyNotFound) ProducerExceptions() {}
func (ProducerPrivKeyNotFound) Code() ExcTypes      { return 3170001 }
func (ProducerPrivKeyNotFound) What() string {
	return "Producer private key is not available"
}

type MissingPendingBlockState struct{ ELog }

func (MissingPendingBlockState) ChainExceptions()    {}
func (MissingPendingBlockState) ProducerExceptions() {}
func (MissingPendingBlockState) Code() ExcTypes      { return 3170002 }
func (MissingPendingBlockState) What() string {
	return "Pending block state is missing"
}

type ProducerDoubleConfirm struct{ ELog }

func (ProducerDoubleConfirm) ChainExceptions()    {}
func (ProducerDoubleConfirm) ProducerExceptions() {}
func (ProducerDoubleConfirm) Code() ExcTypes      { return 3170003 }
func (ProducerDoubleConfirm) What() string {
	return "Producer is double confirming known range"
}

type ProducerScheduleException struct{ ELog }

func (ProducerScheduleException) ChainExceptions()    {}
func (ProducerScheduleException) ProducerExceptions() {}
func (ProducerScheduleException) Code() ExcTypes      { return 3170004 }
func (ProducerScheduleException) What() string {
	return "Producer schedule exception"
}

type ProducerNotInSchedule struct{ ELog }

func (ProducerNotInSchedule) ChainExceptions()    {}
func (ProducerNotInSchedule) ProducerExceptions() {}
func (ProducerNotInSchedule) Code() ExcTypes      { return 3170006 }
func (ProducerNotInSchedule) What() string {
	return "The producer is not part of current schedule"
}
