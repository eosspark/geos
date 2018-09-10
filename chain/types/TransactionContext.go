package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"time"
)

type TransactionContext struct {
	/*Control               *chain.Controller*/
	Trx                   *SignedTransaction
	ID                    common.TransactionIDType
	UndoSession           eosiodb.Session
	Trace                 TransactionTrace
	Start                 common.Tstamp
	Publishe              common.Tstamp
	Executed              []ActionReceipt
	BillToAccounts        []common.AccountName
	ValidateRamUsage      []common.AccountName
	InitialMaxBillableCpu uint64
	Delay                 common.Tstamp //microseconds
	IsInput               bool
	ApplyContextFree      bool
	CanSubjectivelyFail   bool
	DeadLine              common.Tstamp //c++ fc::time_point::maximum()
	Leeway                common.Tstamp
	BilledCpuTimeUs       int64
	ExplicitBilledCpuTime bool

	isInitialized                 bool
	netLimit                      uint64
	netLimitDueToBlock            bool
	netLimitDueToGreylist         bool
	eagerNetLimit                 uint64
	netUsage                      *uint64
	initialObjectiveDurationLimit common.Tstamp //microseconds
	objectiveDurationLimit        common.Tstamp
	deadline                      common.Tstamp
	deadlineExceptionCode         int64
	billingTimerExceptionCode     int64
	pseudoStart                   common.Tstamp
	billedTime                    common.Tstamp
	billingTimerDurationLimit     common.Tstamp
}

func (trxCon TransactionContext) NewTransactionContext(t SignedTransaction, trxId *common.TransactionIDType, s time.Time) (trx TransactionContext) {

	return
}
