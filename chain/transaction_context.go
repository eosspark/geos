package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/log"
)

type TransactionContext struct {
	Control               *Controller
	Trx                   *types.SignedTransaction
	ID                    common.TransactionIdType
	UndoSession           *eosiodb.Session
	Trace                 types.TransactionTrace
	Start                 common.TimePoint
	Published             common.TimePoint
	Executed              []types.ActionReceipt
	BillToAccounts        []common.AccountName
	ValidateRamUsage      []common.AccountName
	InitialMaxBillableCpu uint64
	Delay                 common.Microseconds
	IsInput               bool
	ApplyContextFree      bool
	CanSubjectivelyFail   bool
	Deadline              common.TimePoint //c++ fc::time_point::maximum()
	Leeway                common.Microseconds
	BilledCpuTimeUs       int64
	ExplicitBilledCpuTime bool

	isInitialized                 bool
	netLimit                      uint64
	netLimitDueToBlock            bool
	netLimitDueToGreylist         bool
	cpuLimitDueToGreylist         bool
	eagerNetLimit                 uint64
	netUsage                      *uint64
	initialObjectiveDurationLimit common.Microseconds //microseconds
	objectiveDurationLimit        common.Microseconds
	deadline                      common.TimePoint //maximum
	deadlineExceptionCode         int64
	billingTimerExceptionCode     int64
	pseudoStart                   common.TimePoint
	billedTime                    common.Microseconds
	billingTimerDurationLimit     common.Microseconds
}

func NewTransactionContext(
	c *Controller,
	t *types.SignedTransaction,
	trxId common.TransactionIdType,
	s common.TimePoint) *TransactionContext {

	tc := &TransactionContext{
		Control:     c,
		Trx:         t,
		Start:       s,
		pseudoStart: s,
		Trace:       types.TransactionTrace{ID: trxId},
		//Trace.I:trxId,

		InitialMaxBillableCpu: 0,
		IsInput:               false,
		ApplyContextFree:      true,
		CanSubjectivelyFail:   true,
		Deadline:              common.MaxTimePoint(),
		Leeway:                common.Microseconds(3000),
		BilledCpuTimeUs:       0,
		ExplicitBilledCpuTime: false,

		isInitialized:         false,
		netLimit:              0,
		netLimitDueToBlock:    true,
		netLimitDueToGreylist: false,
		cpuLimitDueToGreylist: false,
		eagerNetLimit:         0,

		deadline:                  common.MaxTimePoint(),
		deadlineExceptionCode:     int64((&exception.BlockCpuUsageExceeded{}).Code()),
		billingTimerExceptionCode: int64((&exception.BlockCpuUsageExceeded{}).Code()),
	}

	//tc.Trace.Id = trxId
	tc.netUsage = &tc.Trace.NetUsage

	if !c.SkipDbSessions() {
		tc.UndoSession = c.DB.StartSession()
	}
	//t.Deadline = common.MaxTimePoint()
	//assert(len(t.Trx.Extensions) == 0), unsupported_feature, "we don't support any extensions yet")

	return tc
}

func int64Max() int64 {
	return int64(^uint(0) >> 1)
}

func (t *TransactionContext) init(initialNetUsage uint64) {
	//const          = int64Max() / 2

	cfg := t.Control.GetGlobalProperties().Configuration
	rl := t.Control.GetMutableResourceLimitsManager()
	t.netLimit = rl.GetBlockNetLimit()
	t.objectiveDurationLimit = common.Microseconds(rl.GetBlockCpuLimit())
	t.deadline = t.Start + common.TimePoint(t.objectiveDurationLimit)

	// Possibly lower net_limit to the maximum net usage a transaction is allowed to be billed
	_mtn := uint64(cfg.MaxTransactionNetUsage)
	if _mtn <= t.netLimit {
		t.netLimit = _mtn
		t.netLimitDueToBlock = false
	}

	// Possibly lower objective_duration_limit to the maximum cpu usage a transaction is allowed to be billed
	_mtcu := uint64(cfg.MaxTransactionCpuUsage)
	if _mtcu <= uint64(t.objectiveDurationLimit.Count()) {
		t.objectiveDurationLimit = common.Milliseconds(int64(cfg.MaxTransactionCpuUsage))
		t.billingTimerExceptionCode = int64((&exception.TxCpuUsageExceed{}).Code()) //TODO
		t.deadline = t.Start + common.TimePoint(t.objectiveDurationLimit)
	}

	// Possibly lower net_limit to optional limit set in the transaction header
	trxSpecifiedNetUsageLimit := uint64(t.Trx.MaxNetUsageWords * 8)
	if trxSpecifiedNetUsageLimit > 0 && trxSpecifiedNetUsageLimit <= t.netLimit {
		t.netLimit = trxSpecifiedNetUsageLimit
		t.netLimitDueToBlock = false
	}

	// Possibly lower objective_duration_limit to optional limit set in transaction header
	//TODO
	if t.Trx.MaxCpuUsageMS > 0 {
		trxSpecifiedCpuUsageLimit := common.Milliseconds(int64(t.Trx.MaxCpuUsageMS))
		if trxSpecifiedCpuUsageLimit <= t.objectiveDurationLimit {
			t.objectiveDurationLimit = trxSpecifiedCpuUsageLimit
			//t.billingTimerExceptionCode = excptionCode	//TODO
			t.deadline = t.Start + common.TimePoint(t.objectiveDurationLimit)
		}
	}

	t.initialObjectiveDurationLimit = t.objectiveDurationLimit

	if t.BilledCpuTimeUs > 0 {
		t.validateCpuUsageToBill(t.BilledCpuTimeUs, false)
	}

	// Record accounts to be billed for network and CPU usage
	for _, act := range t.Trx.Actions {
		for _, auth := range act.Authorization {
			t.BillToAccounts = append(t.BillToAccounts, auth.Actor)
		}
	}

	// Update usage values of accounts to reflect new time
	//bts := common.BlockTimeStamp(t.Control.PendingBlockTime())
	rl.UpdateAccountUsage(t.BillToAccounts, uint32(common.BlockTimeStamp(t.Control.PendingBlockTime())))

	// Calculate the highest network usage and CPU time that all of the billed accounts can afford to be billed
	accountNetLimit, accountCpuLimit, greylistedNet, greylistedCpu := t.MaxBandwidthBilledAccountsCanPay(false)
	t.netLimitDueToGreylist = t.netLimitDueToGreylist || greylistedNet
	t.cpuLimitDueToGreylist = t.cpuLimitDueToGreylist || greylistedCpu

	t.eagerNetLimit = t.netLimit

	// Possible lower eager_net_limit to what the billed accounts can pay plus some (objective) leeway
	newEagerNetLimit := common.Min(t.eagerNetLimit, uint64(accountNetLimit+int64(cfg.NetUsageLeeway)))
	if newEagerNetLimit < t.eagerNetLimit {
		t.eagerNetLimit = newEagerNetLimit
		t.netLimitDueToBlock = false
	}

	// Possibly limit deadline if the duration accounts can be billed for (+ a subjective leeway) does not exceed current delta
	if common.Milliseconds(accountCpuLimit)+t.Leeway <= common.Microseconds(t.deadline-t.Start) {
		t.deadline = t.Start + common.TimePoint(accountCpuLimit) + common.TimePoint(t.Leeway)
		t.billingTimerExceptionCode = int64((&exception.LeewayDeadlineException{}).Code())
	}

	t.billingTimerDurationLimit = common.Microseconds(t.deadline - t.Start)

	if t.ExplicitBilledCpuTime || t.Deadline < t.deadline {
		t.deadline = t.Deadline
		t.deadlineExceptionCode = int64((&exception.DeadlineException{}).Code())
	} else {
		t.deadlineExceptionCode = t.billingTimerExceptionCode
	}

	t.eagerNetLimit = (t.netLimit / 8) * 8 // Round down to nearest multiple of word size (8 bytes) so check_net_usage can be efficient
	if initialNetUsage > 0 {
		t.AddNetUsage(initialNetUsage) // Fail early if current net usage is already greater than the calculated limit
	}

	t.CheckTime()
	t.isInitialized = true
	//fmt.Println(cfg, rl, trxSpecifiedNetUsageLimit, t)

}

func (t *TransactionContext) InitForImplicitTrx(initialNetUsage uint64) {
	t.Published = t.Control.PendingBlockTime()
	t.init(initialNetUsage)
}

func (t *TransactionContext) InitForInputTrx(packeTrxUnprunableSize uint64, packeTrxPrunableSize uint64, nunSignatures uint32, skipRecording bool) {
	cfg := t.Control.GetGlobalProperties().Configuration
	discountedSizeForPrunedData := packeTrxPrunableSize
	if cfg.ContextFreeDiscountNetUsageDen > 0 && cfg.ContextFreeDiscountNetUsageNum < cfg.ContextFreeDiscountNetUsageDen {
		discountedSizeForPrunedData *= uint64(cfg.ContextFreeDiscountNetUsageNum)
		discountedSizeForPrunedData = (discountedSizeForPrunedData + uint64(cfg.ContextFreeDiscountNetUsageDen) - 1) / uint64(cfg.ContextFreeDiscountNetUsageDen)
	}

	initialNetUsage := uint64(cfg.BasePerTransactionNetUsage) + packeTrxUnprunableSize + discountedSizeForPrunedData
	if t.Trx.DelaySec > 0 {
		initialNetUsage += uint64(cfg.BasePerTransactionNetUsage)
		initialNetUsage += uint64(cfg.TransactionIdNetUsage)
	}

	t.Published = t.Control.PendingBlockTime()
	t.IsInput = true

	if t.Control.SkipTrxChecks() {
		t.Control.ValidateExpiration(&t.Trx.Transaction)
		t.Control.ValidateTapos(&t.Trx.Transaction)
		t.Control.ValidateReferencedAccounts(&t.Trx.Transaction)
	}

	t.init(initialNetUsage)
	if !skipRecording {
		//t.recordTransaction(t.ID, t.Trx.Expiration)
	}

}

func (t *TransactionContext) validateCpuUsageToBill(bctu int64, checkMinimum bool) {
	if !t.Control.SkipTrxChecks() {
		if checkMinimum {
			cfg := t.Control.GetGlobalProperties().Configuration
			fmt.Println(cfg)
			/*EOS_ASSERT( billed_us >= cfg.min_transaction_cpu_usage, transaction_exception,
				"cannot bill CPU time less than the minimum of ${min_billable} us",
				("min_billable", cfg.min_transaction_cpu_usage)("billed_cpu_time_us", billed_us)
			)*/
		}
		//if t.billingTimerExceptionCode == exceptionCode {//TODO
		/*EOS_ASSERT( billed_us <= objective_duration_limit.count(),
			block_cpu_usage_exceeded,
			"billed CPU time (${billed} us) is greater than the billable CPU time left in the block (${billable} us)",
			("billed", billed_us)("billable", objective_duration_limit.count())
		)
		} else {
			if t.CpuLimitDueToGreylist {
				EOS_ASSERT( billed_us <= objective_duration_limit.count(),
					greylist_cpu_usage_exceeded,
					"billed CPU time (${billed} us) is greater than the maximum greylisted billable CPU time for the transaction (${billable} us)",
					("billed", billed_us)("billable", objective_duration_limit.count())
				);
			} else {
				EOS_ASSERT( billed_us <= objective_duration_limit.count(),
					tx_cpu_usage_exceeded,
					"billed CPU time (${billed} us) is greater than the maximum billable CPU time for the transaction (${billable} us)",
					("billed", billed_us)("billable", objective_duration_limit.count())
				);
			}
		}*/
	}
}

func (t *TransactionContext) CheckTime() {

	if !t.Control.SkipTrxChecks() {
		_now := common.Now()
		if _now > t.deadline {
			if t.ExplicitBilledCpuTime { //|| deadline_exception_code TODO

			}
		}
	}
	/*if (!control.skip_t_checks()) {
		auto now = fc::time_point::now();
		if( BOOST_UNLIKELY( now > _deadline ) ) {
			// edump((now-start)(now-pseudo_start));
			if( explicit_billed_cpu_time || deadline_exception_code == deadline_exception::code_value ) {
				EOS_THROW( deadline_exception, "deadline exceeded", ("now", now)("deadline", _deadline)("start", start) );
			} else if( deadline_exception_code == block_cpu_usage_exceeded::code_value ) {
				EOS_THROW( block_cpu_usage_exceeded,
				"not enough time left in block to complete executing transaction",
				("now", now)("deadline", _deadline)("start", start)("billing_timer", now - pseudo_start) );
			} else if( deadline_exception_code == tx_cpu_usage_exceeded::code_value ) {
			if (cpu_limit_due_to_greylist) {
				EOS_THROW( greylist_cpu_usage_exceeded,
				"greylisted transaction was executing for too long",
				("now", now)("deadline", _deadline)("start", start)("billing_timer", now - pseudo_start) );
			} else {
				EOS_THROW( tx_cpu_usage_exceeded,
				"transaction was executing for too long",
				("now", now)("deadline", _deadline)("start", start)("billing_timer", now - pseudo_start) );
				}
			} else if( deadline_exception_code == leeway_deadline_exception::code_value ) {
				EOS_THROW( leeway_deadline_exception,
				"the transaction was unable to complete by deadline, "
				"but it is possible it could have succeeded if it were allowed to run to completion",
				("now", now)("deadline", _deadline)("start", start)("billing_timer", now - pseudo_start) );
			}
				EOS_ASSERT( false,  transaction_exception, "unexpected deadline exception code" );
		}
	}*/
}
func (t *TransactionContext) AddNetUsage(u uint64) {
	*t.netUsage = *t.netUsage + u
	t.CheckNetUsage()
}

func (t *TransactionContext) CheckNetUsage() {
	if !t.Control.SkipTrxChecks() {
		if *t.netUsage > t.eagerNetLimit {
			//TODO Throw Exception
			if t.netLimitDueToBlock {
				log.Error("not enough space left in block:${net_usage} > ${net_limit}", t.netUsage, t.netLimit)
			} else if t.netLimitDueToGreylist {
				log.Error("greylisted transaction net usage is too high: ${net_usage} > ${net_limit}", t.netUsage, t.netLimit)
			} else {
				log.Error("transaction net usage is too high: ${net_usage} > ${net_limit}", t.netUsage, t.netLimit)
			}
		}
	}
}
func (t *TransactionContext) AddRamUsage(account common.AccountName, ramDelta int64) {
	rl := t.Control.GetMutableResourceLimitsManager()
	rl.AddPendingRamUsage(account, ramDelta)
	if ramDelta > 0 {
		if len(t.ValidateRamUsage) == 0 {
			t.ValidateRamUsage = []common.AccountName{5}
			t.ValidateRamUsage = append(t.ValidateRamUsage, account)
		} else {
			t.ValidateRamUsage = append(t.ValidateRamUsage, account)
		}
	}
}

func (t *TransactionContext) UpdateBilledCpuTime(now common.TimePoint) uint32 {
	if t.ExplicitBilledCpuTime {
		return uint32(t.BilledCpuTimeUs)
	}
	cfg := t.Control.GetGlobalProperties().Configuration
	first := common.Microseconds(now - t.pseudoStart)
	second := common.Microseconds(cfg.MinTransactionCpuUsage)
	if first > second {
		t.BilledCpuTimeUs = int64(first)
	} else {
		t.BilledCpuTimeUs = int64(second)
	}
	return uint32(t.BilledCpuTimeUs)
}

func (t *TransactionContext) MaxBandwidthBilledAccountsCanPay(forceElasticLimits bool) (int64, int64, bool, bool) {
	rl := t.Control.GetMutableResourceLimitsManager()
	_largeNumberNoOverflow := int64Max() / 2
	_accountNetLimit := _largeNumberNoOverflow
	_accountCpuLimit := _largeNumberNoOverflow
	_greylistedNet := false
	_greylistedCpu := false
	for _, a := range t.BillToAccounts {
		elastic := forceElasticLimits || !(t.Control.IsProducingBlock()) && t.Control.IsResourceGreylisted(&a)
		netLimit := rl.GetAccountNetLimit(a, elastic)
		if netLimit >= 0 {
			if _accountNetLimit > netLimit {
				_accountNetLimit = netLimit
				if !elastic {
					_greylistedCpu = true
				}
			}
		}
		cpuLimit := rl.GetAccountCpuLimit(a, elastic)
		if cpuLimit >= 0 {
			if _accountCpuLimit > cpuLimit {
				_accountCpuLimit = cpuLimit
				if !elastic {
					_greylistedCpu = true
				}
			}
		}
	}

	return _accountNetLimit, _accountCpuLimit, _greylistedNet, _greylistedCpu
}

func (t *TransactionContext) DispathAction(trace *types.ActionTrace, action *types.Action, receiver common.AccountName, contextFree bool, recurseDepth uint32) {

	applyContext := NewApplyContext(t.Control, t, action, recurseDepth)
	applyContext.ContextFree = contextFree
	applyContext.Receiver = receiver

	// try {
	applyContext.Exec()
	//   } catch( ... ) {
	//      *trace = applyContext.Trace
	//      throw
	//   }

	*trace = applyContext.Trace
}

//func (t *TransactionContext) recordTransaction(common.TransactionIdType id, common.TimePointSec expire){
//
//	transaction := types.TransactionObject{}
//	t.Control.DB.Insert()
//}
