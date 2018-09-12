package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
)

type TransactionContext struct {
	Controller            *Controller
	Trx                   *types.SignedTransaction
	ID                    common.TransactionIDType
	UndoSession           *eosiodb.Session
	Trace                 types.TransactionTrace
	Start                 common.TimePoint
	Publishe              common.TimePoint
	Executed              []types.ActionReceipt
	BillToAccounts        []common.AccountName
	ValidateRamUsage      []common.AccountName
	InitialMaxBillableCpu uint64
	Delay                 common.Microseconds
	IsInput               bool
	ApplyContextFree      bool
	CanSubjectivelyFail   bool
	DeadLine              common.TimePoint //c++ fc::time_point::maximum()
	Leeway                common.Microseconds
	BilledCpuTimeUs       int64
	ExplicitBilledCpuTime bool

	isInitialized                 bool
	netLimit                      uint64
	netLimitDueToBlock            bool
	netLimitDueToGreylist         bool
	eagerNetLimit                 uint64
	netUsage                      uint64
	initialObjectiveDurationLimit common.Microseconds //microseconds
	objectiveDurationLimit        common.Microseconds
	deadline                      common.TimePoint //maximum
	deadlineExceptionCode         int64
	billingTimerExceptionCode     int64
	pseudoStart                   common.TimePoint
	billedTime                    common.Microseconds
	billingTimerDurationLimit     common.Microseconds
}

func (trxCon *TransactionContext) NewTransactionContext(c *Controller, t *types.SignedTransaction, trxId common.TransactionIDType, s common.TimePoint) *TransactionContext {
	trxCon.Controller = c
	trxCon.Trx = t
	trxCon.Start = s
	trxCon.netUsage = trxCon.Trace.NetUsage
	trxCon.pseudoStart = s

	if !c.skipDBSessions() {
		trxCon.UndoSession = c.db.StartSession()
	}
	trxCon.Trace.Id = trxId

	return trxCon
}

func (trxCon *TransactionContext) init(initialNetUsage uint64) {
	const LargeNumberNoOverflow = int(^uint(0)>>1) / 2
	cfg := trxCon.Controller.GetGlobalProperties().Configuration

	fmt.Println(cfg)
}
