package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/entity"
	. "github.com/eosspark/eos-go/exception"
	"os"

	"github.com/eosspark/container/sets/treeset"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"math"
)

/*type AccountForSet common.AccountName

func (f *AccountForSet) GetKey() uint64 {
	return uint64(*f)
}*/

type TransactionContext struct {
	Control               *Controller
	Trx                   *types.SignedTransaction
	ID                    common.TransactionIdType
	UndoSession           *database.Session
	Trace                 *types.TransactionTrace
	Start                 common.TimePoint
	Published             common.TimePoint
	Executed              []types.ActionReceipt
	BillToAccounts        treeset.Set
	ValidateRamUsage      treeset.Set
	InitialMaxBillableCpu uint64
	Delay                 common.Microseconds
	IsInput               bool
	ApplyContextFree      bool
	CanSubjectivelyFail   bool
	Deadline              common.TimePoint //c++ fc::time_point::maximum()
	Leeway                common.Microseconds
	BilledCpuTimeUs       int64
	ExplicitBilledCpuTime bool

	isInitialized         bool
	netLimit              uint64
	netLimitDueToBlock    bool
	netLimitDueToGreylist bool
	cpuLimitDueToGreylist bool
	eagerNetLimit         uint64
	netUsage              *uint64
	//initialObjectiveDurationLimit common.Microseconds //microseconds
	objectiveDurationLimit    common.Microseconds
	deadline                  common.TimePoint //maximum
	deadlineExceptionCode     int64
	billingTimerExceptionCode int64
	pseudoStart               common.TimePoint
	billedTime                common.Microseconds
	billingTimerDurationLimit common.Microseconds

	ilog log.Logger
}

func NewTransactionContext(c *Controller, t *types.SignedTransaction, trxId common.TransactionIdType, s common.TimePoint) *TransactionContext {

	tc := TransactionContext{
		Control:               c,
		Trx:                   t,
		Start:                 s,
		pseudoStart:           s,
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
		deadlineExceptionCode:     int64((BlockCpuUsageExceeded{}).Code()),
		billingTimerExceptionCode: int64((BlockCpuUsageExceeded{}).Code()),
		ValidateRamUsage:          *treeset.NewWith(common.CompareName),
	}

	//for testing
	//tc.ValidateRamUsage = make([]common.AccountName, 10)

	// if !c.SkipDbSessions() {
	// 	tc.UndoSession = c.DB.StartSession()
	// }
	tc.Trace = &types.TransactionTrace{
		ID:              trxId,
		BlockNum:        c.PendingBlockState().BlockNum,
		BlockTime:       types.BlockTimeStamp(c.PendingBlockTime()),
		ProducerBlockId: c.PendingProducerBlockId(),
		//ActionTraces: []types.ActionTrace{},
	}
	tc.netUsage = &tc.Trace.NetUsage
	tc.Executed = make([]types.ActionReceipt, tc.Trx.TotalActions())

	EosAssert(len(tc.Trx.TransactionExtensions) == 0, &UnsupportedFeature{}, "we don't support any extensions yet")

	tc.ilog = log.New("transaction_context")
	logHandler := log.StreamHandler(os.Stdout, log.TerminalFormat(true))
	//tc.ilog.SetHandler(log.LvlFilterHandler(log.LvlDebug, logHandler))
	tc.ilog.SetHandler(log.LvlFilterHandler(log.LvlInfo, logHandler))

	return &tc
}

func (t *TransactionContext) init(initialNetUsage uint64) {

	EosAssert(!t.isInitialized, &TransactionException{}, "cannot initialize twice")
	//const static int64_t large_number_no_overflow = std::numeric_limits<int64_t>::max()/2;

	cfg := t.Control.GetGlobalProperties().Configuration
	rl := t.Control.GetMutableResourceLimitsManager()
	t.netLimit = rl.GetBlockNetLimit()
	t.objectiveDurationLimit = common.Microseconds(rl.GetBlockCpuLimit())
	t.deadline = t.Start + common.TimePoint(t.objectiveDurationLimit)

	// Possibly lower net_limit to the maximum net usage a transaction is allowed to be billed
	mtn := uint64(cfg.MaxTransactionNetUsage)
	if mtn <= t.netLimit {
		t.netLimit = mtn
		t.netLimitDueToBlock = false
	}

	// Possibly lower objective_duration_limit to the maximum cpu usage a transaction is allowed to be billed
	mtcu := uint64(cfg.MaxTransactionCpuUsage)
	if mtcu <= uint64(t.objectiveDurationLimit.Count()) {
		t.objectiveDurationLimit = common.Microseconds(int64(cfg.MaxTransactionCpuUsage))
		t.billingTimerExceptionCode = int64(TxCpuUsageExceeded{}.Code())
		t.deadline = t.Start + common.TimePoint(t.objectiveDurationLimit)
	}

	// Possibly lower net_limit to optional limit set in the transaction header
	trxSpecifiedNetUsageLimit := uint64(t.Trx.MaxNetUsageWords * 8)
	if trxSpecifiedNetUsageLimit > 0 && trxSpecifiedNetUsageLimit <= t.netLimit {
		t.netLimit = trxSpecifiedNetUsageLimit
		t.netLimitDueToBlock = false
	}

	// Possibly lower objective_duration_limit to optional limit set in transaction header
	if t.Trx.MaxCpuUsageMS > 0 {
		trxSpecifiedCpuUsageLimit := common.Milliseconds(int64(t.Trx.MaxCpuUsageMS))
		if trxSpecifiedCpuUsageLimit <= t.objectiveDurationLimit {
			t.objectiveDurationLimit = trxSpecifiedCpuUsageLimit
			t.billingTimerExceptionCode = int64(TxCpuUsageExceeded{}.Code()) //TODO
			t.deadline = t.Start + common.TimePoint(t.objectiveDurationLimit)
		}
	}

	//t.initialObjectiveDurationLimit = t.objectiveDurationLimit

	if t.BilledCpuTimeUs > 0 { // could also call on explicit_billed_cpu_time but it would be redundant
		t.validateCpuUsageToBill(t.BilledCpuTimeUs, false) // Fail early if the amount to be billed is too high
	}
	t.BillToAccounts = *treeset.NewWith(common.CompareName)
	// Record accounts to be billed for network and CPU usage
	for _, act := range t.Trx.Actions {
		for _, auth := range act.Authorization {
			//t.BillToAccounts = append(t.BillToAccounts, auth.Actor)
			/*a := common.S(uint64(auth.Actor))
			fmt.Println(a)

			account := AccountForSet(auth.Actor)
			t.BillToAccounts.Insert(&account)*/
			t.BillToAccounts.AddItem(auth.Actor)
		}
	}

	//t.ValidateRamUsage = make([]common.AccountName, t.BillToAccounts.Len())

	// Update usage values of accounts to reflect new time
	rl.UpdateAccountUsage(&t.BillToAccounts, uint32(types.BlockTimeStamp(t.Control.PendingBlockTime())))

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
	if common.Microseconds(int64(accountCpuLimit))+t.Leeway <= common.Microseconds(t.deadline-t.Start) {
		t.deadline = t.Start + common.TimePoint(accountCpuLimit) + common.TimePoint(t.Leeway)
		t.billingTimerExceptionCode = int64(LeewayDeadlineException{}.Code())
	}

	t.billingTimerDurationLimit = common.Microseconds(t.deadline - t.Start)

	// Check if deadline is limited by caller-set deadline (only change deadline if billed_cpu_time_us is not set)
	if t.ExplicitBilledCpuTime || t.Deadline < t.deadline {
		t.deadline = t.Deadline
		t.deadlineExceptionCode = int64(DeadlineException{}.Code())
	} else {
		t.deadlineExceptionCode = t.billingTimerExceptionCode
	}

	t.eagerNetLimit = (t.eagerNetLimit / 8) * 8 // Round down to nearest multiple of word size (8 bytes) so check_net_usage can be efficient
	if initialNetUsage > 0 {
		t.AddNetUsage(initialNetUsage) // Fail early if current net usage is already greater than the calculated limit
	}

	t.CheckTime()
	t.isInitialized = true

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
		initialNetUsage += uint64(common.DefaultConfig.TransactionIdNetUsage)
	}

	t.Published = t.Control.PendingBlockTime()
	t.IsInput = true

	if !t.Control.SkipTrxChecks() {
		t.Control.ValidateExpiration(&t.Trx.Transaction)
		t.Control.ValidateTapos(&t.Trx.Transaction)
		t.Control.ValidateReferencedAccounts(&t.Trx.Transaction)
	}

	t.init(initialNetUsage)
	if !skipRecording {
		t.recordTransaction(&t.ID, t.Trx.Expiration)
	}

}

func (t *TransactionContext) InitForDeferredTrx(p common.TimePoint) {
	t.Published = p
	t.Trace.Scheduled = true
	t.ApplyContextFree = false
	t.init(0)
}

func (t *TransactionContext) Exec() {

	//t.ilog.Info("Exec receiver:%s action:%s", t.Trx.Actions[0].Account, t.Trx.Actions[0].Name)

	EosAssert(t.isInitialized, &TransactionException{}, "must first initialize")

	if t.ApplyContextFree {
		for _, act := range t.Trx.ContextFreeActions {
			t.Trace.ActionTraces = append(t.Trace.ActionTraces, types.ActionTrace{})
			t.DispathAction(&t.Trace.ActionTraces[len(t.Trace.ActionTraces)-1], act, act.Account, true, 0)
		}
	}

	if t.Delay == common.Microseconds(0) {
		for _, act := range t.Trx.Actions {
			t.Trace.ActionTraces = append(t.Trace.ActionTraces, types.ActionTrace{})
			t.DispathAction(&t.Trace.ActionTraces[len(t.Trace.ActionTraces)-1], act, act.Account, false, 0)
		}
	} else {
		t.scheduleTransaction()
	}
}

func (t *TransactionContext) Finalize() {
	EosAssert(t.isInitialized, &TransactionException{}, "must first initialize")

	if t.IsInput {
		am := t.Control.GetMutableAuthorizationManager()
		for _, act := range t.Trx.Actions {
			for _, auth := range act.Authorization {
				am.UpdatePermissionUsage(am.GetPermission(&auth))
			}
		}
	}

	rl := t.Control.GetMutableResourceLimitsManager()
	//for _, a := range t.ValidateRamUsage.Data {
	itr := t.ValidateRamUsage.Iterator()
	for itr.Next() {
		account := itr.Value().(common.AccountName)
		rl.VerifyAccountRamUsage(common.AccountName(account))
	}

	// Calculate the highest network usage and CPU time that all of the billed accounts can afford to be billed
	accountNetLimit, accountCpuLimit, greylistedNet, greylistedCpu := t.MaxBandwidthBilledAccountsCanPay(false)
	t.netLimitDueToGreylist = t.netLimitDueToGreylist || greylistedNet
	t.cpuLimitDueToGreylist = t.cpuLimitDueToGreylist || greylistedCpu

	if accountNetLimit <= int64(t.netLimit) {
		t.netLimit = uint64(accountNetLimit)
		t.netLimitDueToBlock = false
	}

	if accountCpuLimit <= int64(t.objectiveDurationLimit.Count()) {
		t.objectiveDurationLimit = common.Microseconds(accountCpuLimit)
		t.billingTimerExceptionCode = int64((TxCpuUsageExceeded{}).Code())
	}

	*t.netUsage = ((*t.netUsage + 7) / 8) * 8
	t.eagerNetLimit = t.netLimit

	t.CheckNetUsage()
	now := common.Now()
	t.Trace.Elapsed = common.Microseconds(now - t.Start)

	t.UpdateBilledCpuTime(now)
	t.validateCpuUsageToBill(t.BilledCpuTimeUs, true)

	rl.AddTransactionUsage(&t.BillToAccounts, uint64(t.BilledCpuTimeUs), *t.netUsage, uint32(types.BlockTimeStamp(t.Control.PendingBlockTime())))

}

func (t *TransactionContext) Squash() {
	if t.UndoSession != nil {
		t.UndoSession.Squash()
	}
}

func (t *TransactionContext) Undo() {
	if t.UndoSession != nil {
		t.UndoSession.Undo()
	}
}

func (t *TransactionContext) CheckNetUsage() {
	if !t.Control.SkipTrxChecks() {
		if *t.netUsage > t.eagerNetLimit {
			//TODO Throw Exception
			if t.netLimitDueToBlock {
				EosAssert(true,
					&BlockNetUsageExceeded{},
					"not enough space left in block: %d > %d", *t.netUsage, t.eagerNetLimit)
			} else if t.netLimitDueToGreylist {
				EosAssert(true,
					&GreylistNetUsageExceeded{},
					"greylisted transaction net usage is too high: %d > %d", *t.netUsage, t.eagerNetLimit)
			} else {
				EosAssert(true,
					&TxNetUsageExceeded{},
					"greylisted transaction net usage is too high: %d > %d", *t.netUsage, t.eagerNetLimit)
			}
		}
	}
}

func (t *TransactionContext) CheckTime() {

	//return
	//if !t.Control.SkipTrxChecks() {
	now := common.Now()
	if now > t.deadline {
		if t.ExplicitBilledCpuTime || t.deadlineExceptionCode == int64(DeadlineException{}.Code()) { //|| deadline_exception_code TODO
			EosAssert(false,
				&DeadlineException{},
				"deadline exceeded, now %d deadline %d start %d",
				now, t.deadline, t.Start)

		} else if t.deadlineExceptionCode == int64(BlockCpuUsageExceeded{}.Code()) {
			EosAssert(false,
				&BlockCpuUsageExceeded{},
				"not enough time left in block to complete executing transaction, now %d deadline %d start %d billing_timer %d",
				now, t.deadline, t.Start, now-t.pseudoStart)
		} else if t.deadlineExceptionCode == int64(TxCpuUsageExceeded{}.Code()) {
			if t.cpuLimitDueToGreylist {
				EosAssert(false,
					&GreylistCpuUsageExceeded{},
					"greylisted transaction was executing for too long, now %d deadline %d start %d billing_timer %d",
					now, t.deadline, t.Start, now-t.pseudoStart)

			} else {
				EosAssert(false,
					&TxCpuUsageExceeded{},
					"transaction was executing for too long, now %d deadline %d start %d billing_timer %d",
					now, t.deadline, t.Start, now-t.pseudoStart)
			}

		} else if t.deadlineExceptionCode == int64(LeewayDeadlineException{}.Code()) {
			EosAssert(false,
				&LeewayDeadlineException{},
				"the transaction was unable to complete by deadline, ",
				"but it is possible it could have succeeded if it were allowed to run to completion, now %d deadline %d start %d billing_timer %d",
				now, t.deadline, t.Start, now-t.pseudoStart)

		}
		EosAssert(false,
			&TransactionException{},
			"unexpected deadline exception code")

	}
	//}
}

//added to deadline means delete time comsume from PauseBillingTimer to ResumeBillingTimer
func (t *TransactionContext) PauseBillingTimer() {

	if t.ExplicitBilledCpuTime || t.pseudoStart == common.TimePoint(0) {
		return // either irrelevant or already paused
	}

	now := common.Now()
	t.billedTime = common.Microseconds(now - t.pseudoStart)
	t.deadlineExceptionCode = int64((DeadlineException{}).Code()) // Other timeout exceptions cannot be thrown while billable timer is paused.
	t.pseudoStart = common.TimePoint(0)
}

func (t *TransactionContext) ResumeBillingTimer() {
	if t.ExplicitBilledCpuTime || t.pseudoStart != common.TimePoint(0) {
		return // either irrelevant or already running
	}

	now := common.Now()
	t.pseudoStart = now - common.TimePoint(t.billedTime)
	if t.pseudoStart+common.TimePoint(t.billingTimerDurationLimit) <= t.Deadline {
		t.deadline = t.pseudoStart + common.TimePoint(t.billingTimerDurationLimit)
		t.deadlineExceptionCode = t.billingTimerExceptionCode

	} else {
		t.deadline = t.Deadline
		t.deadlineExceptionCode = int64(DeadlineException{}.Code())
	}
}

func (t *TransactionContext) validateCpuUsageToBill(billedUs int64, checkMinimum bool) {
	//return
	if !t.Control.SkipTrxChecks() {
		if checkMinimum {
			cfg := t.Control.GetGlobalProperties().Configuration
			EosAssert(billedUs >= int64(cfg.MinTransactionCpuUsage),
				&TransactionException{},
				"cannot bill CPU time less than the minimum of %d us, billed_cpu_time_us %", cfg.MinTransactionCpuUsage, billedUs)
		}
		if t.billingTimerExceptionCode == int64(BlockCpuUsageExceeded{}.Code()) { //TODO
			EosAssert(billedUs <= t.objectiveDurationLimit.Count(),
				&BlockCpuUsageExceeded{},
				"billed CPU time (%d us) is greater than the billable CPU time left in the block (%d us)",
				billedUs, t.objectiveDurationLimit.Count())
		} else {
			if t.cpuLimitDueToGreylist {
				EosAssert(billedUs <= t.objectiveDurationLimit.Count(),
					&GreylistCpuUsageExceeded{},
					"billed CPU time (%d us) is greater than the maximum greylisted billable CPU time for the transaction (%d us)",
					billedUs, t.objectiveDurationLimit.Count())
			} else {
				EosAssert(billedUs <= t.objectiveDurationLimit.Count(),
					&TxCpuUsageExceeded{},
					"billed CPU time (%d us) is greater than the maximum billable CPU time for the transaction (%d us)",
					billedUs, t.objectiveDurationLimit.Count())
			}
		}
	}
}
func (t *TransactionContext) AddNetUsage(u uint64) {
	*t.netUsage = *t.netUsage + u
	t.CheckNetUsage()
}

func (t *TransactionContext) AddRamUsage(account common.AccountName, ramDelta int64) {
	rl := t.Control.GetMutableResourceLimitsManager()
	rl.AddPendingRamUsage(account, ramDelta)
	if ramDelta > 0 {
		//a := AccountForSet(account)
		t.ValidateRamUsage.AddItem(account)
	}
}

func (t *TransactionContext) UpdateBilledCpuTime(now common.TimePoint) uint32 {
	if t.ExplicitBilledCpuTime {
		return uint32(t.BilledCpuTimeUs)
	}
	cfg := t.Control.GetGlobalProperties().Configuration
	t.BilledCpuTimeUs = int64(common.Max(uint64(now-t.pseudoStart), uint64(cfg.MinTransactionCpuUsage)))

	//t.ilog.Info("BilledCpuTimeUs:%v", t.BilledCpuTimeUs)
	return uint32(t.BilledCpuTimeUs)
}

//func Min(x, y int64) int64 {
//	if x < y {
//		return x
//	} else {
//		return y
//	}
//}

func (t *TransactionContext) MaxBandwidthBilledAccountsCanPay(forceElasticLimits bool) (int64, int64, bool, bool) {
	rl := t.Control.GetMutableResourceLimitsManager()
	largeNumberNoOverflow := int64(math.MaxInt64 / 2)
	accountNetLimit := largeNumberNoOverflow
	accountCpuLimit := largeNumberNoOverflow
	greylistedNet := false
	greylistedCpu := false
	//for _, a := range t.BillToAccounts.Data {
	itr := t.BillToAccounts.Iterator()
	for itr.Next() {
		accountName := itr.Value().(common.AccountName)
		account := common.AccountName(accountName)

		elastic := forceElasticLimits || !(t.Control.IsProducingBlock()) && t.Control.IsResourceGreylisted(&account)
		netLimit := uint64(rl.GetAccountNetLimit(account, elastic))
		if netLimit >= 0 {
			accountNetLimit = int64(common.Min(uint64(accountNetLimit), uint64(netLimit)))
			if !elastic {
				greylistedNet = true
			}
		}
		cpuLimit := uint64(rl.GetAccountCpuLimit(account, elastic))
		if cpuLimit >= 0 {
			accountCpuLimit = int64(common.Min(uint64(accountCpuLimit), uint64(cpuLimit)))
			if !elastic {
				greylistedCpu = true
			}
		}
	}
	return accountNetLimit, accountCpuLimit, greylistedNet, greylistedCpu
}

func (t *TransactionContext) DispathAction(trace *types.ActionTrace, action *types.Action, receiver common.AccountName, contextFree bool, recurseDepth uint32) {

	applyContext := NewApplyContext(t.Control, t, action, recurseDepth)
	applyContext.ContextFree = contextFree
	applyContext.Receiver = receiver

	//try.Try(func() {
	applyContext.Exec(trace)
	//}).Catch(func(e Exception) {
	//	*trace = applyContext.Trace
	//	try.Throw(e)
	//}).End()

	//*trace = applyContext.Trace
}

func (t *TransactionContext) scheduleTransaction() {

	if t.Trx.DelaySec == 0 {
		cfg := t.Control.GetGlobalProperties().Configuration
		t.AddNetUsage(uint64(cfg.BasePerTransactionNetUsage + common.DefaultConfig.TransactionIdNetUsage))
	}

	firstAuth := t.Trx.FirstAuthorizor()
	var trxSize uint32 = 0

	gto := entity.GeneratedTransactionObject{
		TrxId:     t.ID,
		Payer:     firstAuth,
		Sender:    common.AccountName(0),
		Published: t.Control.PendingBlockTime(),
	}
	//gto.SenderId = transactionIdToSenderId(gto.TrxId)
	gto.DelayUntil = gto.Published + common.TimePoint(t.Delay)
	gto.Expiration = gto.DelayUntil + common.TimePoint(common.Seconds(int64(t.Control.GetGlobalProperties().Configuration.DeferredTrxExpirationWindow)))
	trxSize = 0 //gto.set(t.Trx) //TODO
	t.Control.DB.Insert(&gto)

	t.AddRamUsage(gto.Payer, int64(common.BillableSizeV("generated_transaction_object")+uint64(trxSize)))

}

func (t *TransactionContext) recordTransaction(id *common.TransactionIdType, expire common.TimePointSec) {

	obj := entity.TransactionObject{Expiration: expire, TrxID: *id}
	t.Control.DB.Insert(&obj)
}
