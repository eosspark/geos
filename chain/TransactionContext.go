package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"time"
)

type TransactionContext struct {
	Controller            *Controller
	Trx                   *types.SignedTransaction
	ID                    common.TransactionIDType
	UndoSession           eosiodb.Session
	Trace                 types.TransactionTrace
	Start                 common.Tstamp
	Publishe              common.Tstamp
	Executed              []types.ActionReceipt
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

func (trxCon *TransactionContext) NewTransactionContext(t types.SignedTransaction, trxId *common.TransactionIDType, s time.Time) (trx TransactionContext) {

	return
}
