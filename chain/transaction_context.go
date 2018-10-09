package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/database"
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
	DeadLine              common.TimePoint //c++ fc::time_point::maximum()
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

func (self *TransactionContext) NewTransactionContext(c *Controller, t *types.SignedTransaction, trxId common.TransactionIdType, s common.TimePoint) *TransactionContext {
	self.Control = c
	self.Trx = t
	self.Start = s
	self.netUsage = &self.Trace.NetUsage
	self.pseudoStart = s

	if !c.SkipDbSessions() {
		self.UndoSession = c.DB.StartSession()
	}
	self.Trace.ID = trxId
	self.Trace.BlockNum = c.PendingBlockState().BlockNum
	self.Trace.BlockTime = common.NewBlockTimeStamp(c.PendingBlockTime())
	self.Trace.ProducerBlockId = c.PendingProducerBlockId()

	/*	TODO
		EOS_ASSERT( trx.transaction_extensions.size() == 0, unsupported_feature, "we don't support any extensions yet" )
	*/

	self.IsInput = false
	self.ApplyContextFree=true
	self.DeadLine = common.MaxTimePoint()
	self.Leeway =common.Microseconds(3000)
	self.BilledCpuTimeUs = 0
	self.ExplicitBilledCpuTime=false
	self.isInitialized = false
	self.netLimit = 0
	self.netLimitDueToBlock = true
	self.netLimitDueToGreylist=false
	self.eagerNetLimit = 0
	self.cpuLimitDueToGreylist = false
	self.deadline = common.MaxTimePoint()
	self.deadlineExceptionCode = 0	//TODO exception code value
	self.billingTimerExceptionCode = 0 //TODO exception code value
	return self
}

func (trxCon *TransactionContext) init(initialNetUsage uint64) {
	const LargeNumberNoOverflow = int64(^uint(0)>>1) / 2

	cfg := trxCon.Control.GetGlobalProperties().Configuration
	rl := trxCon.Control.GetMutableResourceLimitsManager()
	trxCon.netLimit = rl.GetBlockNetLimit()
	trxCon.objectiveDurationLimit = common.Microseconds(rl.GetBlockCpuLimit())
	trxCon.deadline = trxCon.Start + common.TimePoint(trxCon.objectiveDurationLimit)
	_mtn := uint64(cfg.MaxTransactionNetUsage)
	if _mtn <= trxCon.netLimit {
		trxCon.netLimit = _mtn
		trxCon.netLimitDueToBlock = false
	}

	_mtcu := uint64(cfg.MaxTransactionCpuUsage)
	if _mtcu <= uint64(trxCon.objectiveDurationLimit) {
		trxCon.objectiveDurationLimit = common.Milliseconds(int64(cfg.MaxTransactionCpuUsage))
		//trxCon.billingTimerExceptionCode = excptionCode	//TODO
		trxCon.DeadLine = common.TimePoint(trxCon.Start + common.TimePoint(trxCon.objectiveDurationLimit))
	}

	trxSpecifiedNetUsageLimit := uint64(trxCon.Trx.MaxNetUsageWords * 8)

	if trxSpecifiedNetUsageLimit > 0 && trxSpecifiedNetUsageLimit <= trxCon.netLimit {
		trxCon.netLimit = trxSpecifiedNetUsageLimit
		trxCon.netLimitDueToBlock = false
	}

	if trxCon.Trx.MaxCpuUsageMS > 0 {
		trxSpecifiedCpuUsageLimit := common.Milliseconds(int64(trxCon.Trx.MaxCpuUsageMS))
		if trxSpecifiedCpuUsageLimit <= trxCon.objectiveDurationLimit {
			trxCon.objectiveDurationLimit = trxSpecifiedCpuUsageLimit
			//trxCon.billingTimerExceptionCode = excptionCode	//TODO
			trxCon.deadline = trxCon.Start + common.TimePoint(trxCon.objectiveDurationLimit)
		}
	}

	trxCon.initialObjectiveDurationLimit = trxCon.objectiveDurationLimit

	if trxCon.BilledCpuTimeUs > 0 {
		trxCon.validateCpuUsageToBill(trxCon.BilledCpuTimeUs, false)
	}

	for _, act := range trxCon.Trx.Actions {
		for _, auth := range act.Authorization {
			trxCon.BillToAccounts = append(trxCon.BillToAccounts, auth.Actor)
		}
	}

	bts := common.BlockTimeStamp(trxCon.Control.PendingBlockTime())
	rl.UpdateAccountUsage(trxCon.BillToAccounts, uint32(bts))

	t := trxCon.MaxBandwidthBilledAccountsCanPay(false) //default false

	if trxCon.netLimitDueToGreylist || t._gn {
		trxCon.netLimitDueToGreylist = true
	} else {
		trxCon.netLimitDueToGreylist = false
	}
	if trxCon.cpuLimitDueToGreylist || t._gc {
		trxCon.cpuLimitDueToGreylist = true
	} else {
		trxCon.cpuLimitDueToGreylist = false
	}

	trxCon.eagerNetLimit = (trxCon.netLimit / 8) * 8
	if initialNetUsage > 0 {
		trxCon.AddNetUsage(initialNetUsage)
	}

	trxCon.CheckTime()

	trxCon.isInitialized = true
	fmt.Println(cfg, rl, trxSpecifiedNetUsageLimit, t)

}

func (trx *TransactionContext) validateCpuUsageToBill(bctu int64, checkMinimum bool) {
	if !trx.Control.SkipTrxChecks() {
		if checkMinimum {
			cfg := trx.Control.GetGlobalProperties().Configuration
			fmt.Println(cfg)
			/*EOS_ASSERT( billed_us >= cfg.min_transaction_cpu_usage, transaction_exception,
				"cannot bill CPU time less than the minimum of ${min_billable} us",
				("min_billable", cfg.min_transaction_cpu_usage)("billed_cpu_time_us", billed_us)
			)*/
		}
		//if trx.billingTimerExceptionCode == exceptionCode {//TODO
		/*EOS_ASSERT( billed_us <= objective_duration_limit.count(),
			block_cpu_usage_exceeded,
			"billed CPU time (${billed} us) is greater than the billable CPU time left in the block (${billable} us)",
			("billed", billed_us)("billable", objective_duration_limit.count())
		)
		} else {
			if trx.CpuLimitDueToGreylist {
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

func (trx *TransactionContext) CheckTime() {

	if !trx.Control.SkipTrxChecks() {
		_now := common.Now()
		if _now > trx.deadline {
			if trx.ExplicitBilledCpuTime { //|| deadline_exception_code TODO

			}
		}
	}
	/*if (!control.skip_trx_checks()) {
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
func (trx *TransactionContext) AddNetUsage(u uint64) {
	nu:=*trx.netUsage + u
	trx.netUsage = &nu
	trx.CheckNetUsage()
}

func (trx *TransactionContext) CheckNetUsage() {
	if !trx.Control.SkipTrxChecks() {
		if *trx.netUsage > trx.eagerNetLimit {
			//TODO Throw Exception
			if trx.netLimitDueToBlock {
				log.Error("not enough space left in block:${net_usage} > ${net_limit}", trx.netUsage, trx.netLimit)
			} else if trx.netLimitDueToGreylist {
				log.Error("greylisted transaction net usage is too high: ${net_usage} > ${net_limit}", trx.netUsage, trx.netLimit)
			} else {
				log.Error("transaction net usage is too high: ${net_usage} > ${net_limit}", trx.netUsage, trx.netLimit)
			}
		}
	}
}
func (trx *TransactionContext) AddRamUsage(account common.AccountName, ramDelta int64) {
	rl := trx.Control.GetMutableResourceLimitsManager()
	rl.AddPendingRamUsage(account, ramDelta)
	if ramDelta > 0 {
		if len(trx.ValidateRamUsage) == 0 {
			trx.ValidateRamUsage = []common.AccountName{5}
			trx.ValidateRamUsage = append(trx.ValidateRamUsage, account)
		} else {
			trx.ValidateRamUsage = append(trx.ValidateRamUsage, account)
		}
	}
}

func (trx *TransactionContext) UpdateBilledCpuTime(now common.TimePoint) uint32 {
	if trx.ExplicitBilledCpuTime {
		return uint32(trx.BilledCpuTimeUs)
	}
	cfg := trx.Control.GetGlobalProperties().Configuration
	first := common.Microseconds(now - trx.pseudoStart)
	second := common.Microseconds(cfg.MinTransactionCpuUsage)
	if first > second {
		trx.BilledCpuTimeUs = int64(first)
	} else {
		trx.BilledCpuTimeUs = int64(second)
	}
	return uint32(trx.BilledCpuTimeUs)
}

type tmp struct {
	_anlt int64
	_aclt int64
	_gn   bool
	_gc   bool
}

func (trx *TransactionContext) MaxBandwidthBilledAccountsCanPay(forceElasticLimits bool) tmp {
	rl := trx.Control.GetMutableResourceLimitsManager()
	_largeNumberNoOverflow := int64(^uint(0)>>1) / 2
	_accountNetLimit := _largeNumberNoOverflow
	_accountCpuLimit := _largeNumberNoOverflow
	_greylistedNet := false
	_greylistedCpu := false
	for _, a := range trx.BillToAccounts {
		elastic := forceElasticLimits || !(trx.Control.IsProducingBlock()) && trx.Control.IsResourceGreylisted(&a)
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
	_makeTuple := tmp{_accountNetLimit, _accountCpuLimit, _greylistedNet, _greylistedCpu}

	return _makeTuple
}

func (trx *TransactionContext) InitForImplicitTrx(initialNetUsage uint64) {
	trx.Published = trx.Control.PendingBlockTime()
	trx.init(initialNetUsage)
}

func (trx *TransactionContext) InitForInputTrx(packedTrxUnprunableSize uint64, packedTrxPrunableSize uint64, nunSignatures uint32, skipRecording bool) {
	cfg := trx.Control.GetGlobalProperties().Configuration
	dsfpd := packedTrxPrunableSize
	if cfg.ContextFreeDiscountNetUsageDen > 0 && cfg.ContextFreeDiscountNetUsageNum < cfg.ContextFreeDiscountNetUsageDen {
		dsfpd *= uint64(cfg.ContextFreeDiscountNetUsageNum)
		dsfpd = (dsfpd + uint64(cfg.ContextFreeDiscountNetUsageDen) - 1) / uint64(cfg.ContextFreeDiscountNetUsageDen)
	}
	//TODO append
	initialNetUsage :=uint64(cfg.BasePerTransactionNetUsage)+packedTrxUnprunableSize+dsfpd
	if trx.Trx.DelaySec >0{
		initialNetUsage += uint64(cfg.BasePerTransactionNetUsage+common.DefaultConfig.TransactionIdNetUsage)
	}

	trx.Published = trx.Control.PendingBlockTime()
	trx.IsInput = true
	if !trx.Control.SkipTrxChecks(){
		trx.Control.ValidateExpiration(&trx.Trx.Transaction)
		trx.Control.ValidateTapos(&trx.Trx.Transaction)
		trx.Control.ValidateReferencedAccounts(&trx.Trx.Transaction)
	}
	trx.init(initialNetUsage)
	if !skipRecording{
		trx.recordTransaction( &trx.ID, trx.Trx.Expiration ); /// checks for dupes
	}
}

func (trx *TransactionContext) recordTransaction(id *common.TransactionIdType,expire common.TimePointSec){
	//TODO wait modify callback
	to:=entity.TransactionObject{}
	to.TrxID = *id
	to.Expiration = expire
	trx.Control.DataBase().Insert(&to)
}

func (trx *TransactionContext) DispathAction(trace *types.ActionTrace, action *types.Action, receiver common.AccountName, contextFree bool, recurseDepth uint32) {

	applyContext := NewApplyContext(trx.Control, trx, action, recurseDepth)
	applyContext.ContextFree = contextFree
	applyContext.Receiver = receiver

	// try {
	//      applyContext.exec()
	//   } catch( ... ) {
	//      *trace = applyContext.Trace
	//      throw
	//   }

	trace = &applyContext.Trace

}
