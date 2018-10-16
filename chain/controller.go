package chain

import (
	"fmt"
	"os"
	"time"

	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/cvm/exec"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/log"
)

type DBReadMode int8

const (
	SPECULATIVE = DBReadMode(iota)
	HEADER      //HEAD
	READONLY
	IRREVERSIBLE
)

type ValidationMode int8

const (
	FULL = ValidationMode(iota)
	LIGHT
)

type Config struct {
	ActorWhitelist      map[common.AccountName]struct{}
	ActorBlacklist      map[common.AccountName]struct{}
	ContractWhitelist   map[common.AccountName]struct{}
	ContractBlacklist   map[common.AccountName]struct{}
	ActionBlacklist     map[common.Pair]struct{} //see actionBlacklist
	KeyBlacklist        map[ecc.PublicKey]struct{}
	blocksDir           string
	stateDir            string
	stateSize           uint64
	stateGuardSize      uint64
	reversibleCacheSize uint64
	reversibleGuardSize uint64
	readOnly            bool
	forceAllChecks      bool
	disableReplayOpts   bool
	disableReplay       bool
	contractsConsole    bool
	genesis             types.GenesisState
	vmType              exec.WasmInterface
	readMode            DBReadMode
	blockValidationMode ValidationMode
	resourceGreylist    []common.AccountName
	trustedProducers    []common.AccountName
}

var isActiveController bool //default value false ;Does the process include control ;

var instance *Controller

type v func(ctx *ApplyContext)

//type HandlerKey common.Tuple
type Controller struct {
	DB                             *database.LDataBase
	DbSession                      *database.Session
	ReversibleBlocks               *database.LDataBase
	Blog                           string //TODO
	Pending                        *types.PendingState
	Head                           *types.BlockState
	ForkDB                         *types.ForkDatabase
	WasmIf                         *exec.WasmInterface
	ResourceLimists                *ResourceLimitsManager
	Authorization                  *AuthorizationManager
	Config                         Config //local	Config
	ChainID                        common.ChainIdType
	RePlaying                      bool
	ReplayHeadTime                 common.TimePoint //optional<common.Tstamp>
	ReadMode                       DBReadMode
	InTrxRequiringChecks           bool                //if true, checks that are normally skipped on replay (e.g. auth checks) cannot be skipped
	SubjectiveCupLeeway            common.Microseconds //optional<common.Tstamp>
	TrustedProducerLightValidation bool                //default value false
	ApplyHandlers                  map[common.AccountName]map[HandlerKey]v
	UnAppliedTransactions          map[crypto.Sha256]types.TransactionMetadata
}

func GetControllerInstance() *Controller {
	if !isActiveController {
		validPath()
		instance = newController()
		readycontroller <- true
		time.Sleep(2 * time.Second) //TODO for test case
	}
	return instance
}

//TODO tmp code

func validPath() {
	path := []string{common.DefaultConfig.DefaultStateDirName, common.DefaultConfig.DefaultBlocksDirName, common.DefaultConfig.DefaultReversibleBlocksDirName}
	for _, d := range path {
		_, err := os.Stat(d)
		if os.IsNotExist(err) {
			err := os.MkdirAll(d, os.ModePerm)
			if err != nil {
				fmt.Printf("controller validPath mkdir failed![%v]\n", err)
			} else {
				fmt.Printf("controller validPath mkdir success!\n", d)
			}
		}
	}
}
func newController() *Controller {
	isActiveController = true //controller is active
	//init db
	db, err := database.NewDataBase(common.DefaultConfig.DefaultStateDirName)
	if err != nil {
		fmt.Println("newController is error detail:", err)
		return nil
	}
	defer db.Close()

	//init ReversibleBlocks
	//reversibleDir := common.DefaultConfig.DefaultBlocksDirName + "/" + common.DefaultConfig.DefaultReversibleBlocksDirName
	reversibleDB, err := database.NewDataBase(common.DefaultConfig.DefaultReversibleBlocksDirName)
	if err != nil {
		fmt.Println("newController init reversibleDB is error", err)
	}
	con := &Controller{InTrxRequiringChecks: false, RePlaying: false, TrustedProducerLightValidation: false}
	con.DB = db
	con.ReversibleBlocks = reversibleDB
	con.ForkDB = types.GetForkDbInstance(common.DefaultConfig.DefaultBlocksDirName)

	con.ChainID = types.GetGenesisStateInstance().ComputeChainID()

	con.initConfig()
	con.ReadMode = con.Config.readMode
	con.ApplyHandlers = make(map[common.AccountName]map[HandlerKey]v)
	con.WasmIf = exec.NewWasmInterface()

	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("eosio")), applyEosioNewaccount)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("eosio")), applyEosioSetcode)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("eosio")), applyEosioSetabi)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("eosio")), applyEosioUpdateauth)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("eosio")), applyEosioDeleteauth)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("eosio")), applyEosioUnlinkauth)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("eosio")), applyEosioLinkauth)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("eosio")), applyEosioCanceldalay)

	//IrreversibleBlock.connect()
	readycontroller = make(chan bool)
	go initResource(con, readycontroller)

	return con
}

func condition(contract common.AccountName, action common.ActionName) string {
	c := capitalize(common.S(uint64(contract)))
	a := capitalize(common.S(uint64(action)))

	return c + a
}

func capitalize(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 {
				vv[i] -= 32
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

//TODO wait append block_log
func (c *Controller) OnIrreversible(b *types.BlockState) {

}

func (c *Controller) PopBlock() {
	prev := c.ForkDB.GetBlock(&c.Head.Header.Previous)
	r := types.ReversibleBlockObject{}
	//r.BlockNum = c.Head.BlockNum
	errs := c.ReversibleBlocks.Find("NUM", c.Head.BlockNum, r)

	if errs != nil {
		log.Error("PopBlock ReversibleBlocks Find is error,detail:", errs)
	}
	c.ReversibleBlocks.Remove(&r)

	if c.ReadMode == SPECULATIVE {
		trx := c.Head.Trxs
		step := 0
		for ; step < len(trx); step++ {
			c.UnAppliedTransactions[crypto.Sha256(trx[step].SignedID)] = *trx[step]
		}
	}
	c.Head = prev
	c.DbSession.Undo() //TODO
}

func (c *Controller) SetApplayHandler(receiver common.AccountName, contract common.AccountName, action common.ActionName, handler func(a *ApplyContext)) {
	hk := NewHandlerKey(common.ScopeName(contract), action)
	first := make(map[common.AccountName]map[HandlerKey]v)
	secend := make(map[HandlerKey]v)
	secend[hk] = handler
	first[receiver] = secend
	c.ApplyHandlers[receiver] = secend
}

func (c *Controller) AbortBlock() {
	if c.Pending != nil {
		if c.ReadMode == SPECULATIVE {
			trx := append(c.Pending.PendingBlockState.Trxs)
			step := 0
			for ; step < len(trx); step++ {
				c.UnAppliedTransactions[crypto.Sha256(trx[step].SignedID)] = *trx[step]
			}
		}
	}
}

func (c *Controller) StartBlock(when common.BlockTimeStamp, confirmBlockCount uint16, s types.BlockStatus, producerBlockId *common.BlockIdType) {
	if c.Pending != nil {
		fmt.Println("pending block already exists")
		return
	}
	//defer c.Pending.reset()
	if c.SkipDbSession(s) {
		c.Pending = types.NewPendingState(c.DB)
	} else {
		c.Pending = types.GetInstance()
	}

	c.Pending.BlockStatus = s
	c.Pending.ProducerBlockId = *producerBlockId
	c.Pending.PendingBlockState = c.Head //TODO std::make_shared<block_state>( *head, when ); // promotes pending schedule (if any) to active
	c.Pending.PendingBlockState.SignedBlock.Timestamp = when
	c.Pending.PendingBlockState.InCurrentChain = true
	c.Pending.PendingBlockState.SetConfirmed(confirmBlockCount)
	wasPendingPromoted := c.Pending.PendingBlockState.MaybePromotePending()
	log.Info("wasPendingPromoted", wasPendingPromoted)
	if c.ReadMode == DBReadMode(SPECULATIVE) || c.Pending.BlockStatus != types.BlockStatus(types.Incomplete) {
		gpo := types.GlobalPropertyObject{}
		itr, err := c.DB.Get("ID", gpo)
		if err != nil {
			log.Error("Controller StartBlock find GlobalPropertyObject is error :", err)
			return
		}

		err = itr.Data(gpo)
		//gpo.ProposedScheduleBlockNum.valid() //if there is a proposed schedule that was proposed in a block ...
		if ( /*gpo.ProposedScheduleBlockNum.valid() &&*/ gpo.ProposedScheduleBlockNum <= c.Pending.PendingBlockState.DposIrreversibleBlocknum) &&
			(len(c.Pending.PendingBlockState.PendingSchedule.Producers) == 0) &&
			(!wasPendingPromoted) {
			if !c.RePlaying {
				tmp := gpo.ProposedSchedule.ProducerScheduleType()
				ps := types.SharedProducerScheduleType{}
				ps.Version = tmp.Version
				ps.Producers = tmp.Producers
				c.Pending.PendingBlockState.SetNewProducers(ps)
			}

			c.DB.Modify(&gpo, func(i *types.GlobalPropertyObject) error {
				i.ProposedScheduleBlockNum = 1
				i.ProposedSchedule.Clear()
				return nil
			})
		}

		signedTransaction := c.GetOnBlockTransaction()
		onbtrx := types.TransactionMetadata{Trx: &signedTransaction}
		onbtrx.Implicit = true
		//TODO defer
		c.InTrxRequiringChecks = true
		c.PushTransaction(onbtrx, common.MaxTimePoint(), c.GetGlobalProperties().Configuration.MinTransactionCpuUsage, true)
		fmt.Println(onbtrx)

		c.clearExpiredInputTransactions()
		c.UpdateProducersAuthority()

	}

}

func (c *Controller) PushReceipt(trx interface{}, status types.BlockStatus, cpuUsageUs uint64, netUsage uint64) *types.TransactionReceipt {
	trxReceipt := types.TransactionReceipt{}
	tr := types.TransactionWithID{}
	switch trx.(type) {
	case common.TransactionIdType:
		tr.TransactionID = trx.(common.TransactionIdType)
	case types.PackedTransaction:
		tr.PackedTransaction = trx.(*types.PackedTransaction)
	}
	trxReceipt.Trx = tr
	netUsageWords := netUsage / 8
	//EOS_ASSERT( net_usage_words*8 == net_usage, transaction_exception, "net_usage is not divisible by 8" )
	c.Pending.PendingBlockState.SignedBlock.Transactions = append(c.Pending.PendingBlockState.SignedBlock.Transactions, trxReceipt)
	trxReceipt.CpuUsageUs = uint32(cpuUsageUs)
	trxReceipt.NetUsageWords = uint32(netUsageWords)
	trxReceipt.Status = types.TransactionStatus(status)
	return &trxReceipt
}
func (c *Controller) PushTransaction(trx types.TransactionMetadata, deadLine common.TimePoint, billedCpuTimeUs uint32, explicitBilledCpuTime bool) (trxTrace types.TransactionTrace) {
	exception.EosAssert(deadLine != common.TimePoint(0), &exception.TransactionException{}, "deadline cannot be uninitialized")

	trxContext := *NewTransactionContext(c, trx.Trx, trx.ID, common.Now())

	if c.SubjectiveCupLeeway != 0 {
		if c.Pending.BlockStatus == types.BlockStatus(types.Incomplete) {
			trxContext.Leeway = c.SubjectiveCupLeeway
		}
	}
	trxContext.Deadline = deadLine
	trxContext.ExplicitBilledCpuTime = explicitBilledCpuTime
	trxContext.BilledCpuTimeUs = int64(billedCpuTimeUs)

	trace := trxContext.Trace
	//try{
	if trx.Implicit {
		trxContext.InitForImplicitTrx(0) //default value 0
		trxContext.CanSubjectivelyFail = false
	} else {
		skipRecording := (c.ReplayHeadTime != 0) && (common.TimePoint(trx.Trx.Expiration) <= c.ReplayHeadTime)
		trxContext.InitForInputTrx(uint64(trx.PackedTrx.GetUnprunableSize()), uint64(trx.PackedTrx.GetPrunableSize()),
			uint32(len(trx.Trx.Signatures)), skipRecording)
	}
	if trxContext.CanSubjectivelyFail && c.Pending.BlockStatus == types.Incomplete {
		c.CheckActorList(trxContext.BillToAccounts)
	}
	trxContext.Delay = common.Microseconds(trx.Trx.DelaySec)
	if !c.SkipAuthCheck() && !trx.Implicit {
		c.Authorization.CheckAuthorization(trx.Trx.Actions,
			trx.RecoverKeys(c.ChainID),
			nil,
			trxContext.Delay,
			nil,
			false)
	}
	trxContext.Exec()
	trxContext.Finalize()

	//restore := c.MakeBlockRestorePoint()
	if !trx.Implicit {
		var s types.TransactionStatus
		if trxContext.Delay == common.Microseconds(0) {
			s = types.TransactionStatusExecuted
		} else {
			s = types.TransactionStatusDelayed
		}
		fmt.Println(trace, s)
		tr := c.PushReceipt(trx.PackedTrx.PackedTrx, types.BlockStatus(s), uint64(trxContext.BilledCpuTimeUs), trace.NetUsage)
		trace.Receipt = tr.TransactionReceiptHeader
		c.Pending.PendingBlockState.Trxs = append(c.Pending.PendingBlockState.Trxs, &trx)
	} else {
		r := types.TransactionReceiptHeader{}
		r.CpuUsageUs = uint32(trxContext.BilledCpuTimeUs)
		r.NetUsageWords = uint32(trace.NetUsage / 8)
		trace.Receipt = r
	}
	//fc::move_append(pending->_actions, move(trx_context.executed))
	if !trx.Accepted {
		trx.Accepted = true
		//emit( c.accepted_transaction, trx)
	}

	//emit(c.applied_transaction, trace)
	if c.ReadMode != SPECULATIVE && c.Pending.BlockStatus == types.Incomplete {
		trxContext.Undo()
	} else {
		//restore.cancel()
		trxContext.Squash()
	}

	if !trx.Implicit {
		delete(c.UnAppliedTransactions, crypto.Hash256(trx.SignedID))
	}

	//return trace
	/*}catch(Exception{}){

	}*/
	if !failureIsSubjective( /**trace.Except*/ ) {
		delete(c.UnAppliedTransactions, crypto.Sha256(trx.SignedID))
	}
	/*emit( c.accepted_transaction, trx )
	emit( c.applied_transaction, trace )*/
	return trace
}

func (c *Controller) GetGlobalProperties() (gp *types.GlobalPropertyObject) {
	gpo := types.GlobalPropertyObject{}
	itr, err := c.DB.Get("ID", gpo)
	if err != nil {
		log.Error("GetGlobalProperties is error detail:", err)
	}
	if itr.Next() {
		err = itr.Data(gp)
	}
	return gp
}

func (c *Controller) GetDynamicGlobalProperties() (r *types.DynamicGlobalPropertyObject) {
	dgpo := types.DynamicGlobalPropertyObject{}
	itr, err := c.DB.Get("ID", &dgpo)
	if err != nil {
		log.Error("GetDynamicGlobalProperties is error detail:", err)
	}
	if itr.Next() {
		err = itr.Data(r)
		if err != nil {
			fmt.Println("GetDynamicGlobalProperties Data is error:", err)
		}
	}
	return &dgpo
}

func (c *Controller) GetMutableResourceLimitsManager() *ResourceLimitsManager {
	return c.ResourceLimists
}

func (c *Controller) GetOnBlockTransaction() types.SignedTransaction {
	onBlockAction := types.Action{}
	onBlockAction.Account = common.AccountName(common.DefaultConfig.SystemAccountName)
	onBlockAction.Name = common.ActionName(common.N("onblock"))
	onBlockAction.Authorization = []types.PermissionLevel{{common.AccountName(common.DefaultConfig.SystemAccountName), common.PermissionName(common.DefaultConfig.ActiveName)}}

	data, err := rlp.EncodeToBytes(c.Head.Header)
	if err != nil {
		onBlockAction.Data = data
	}
	trx := types.SignedTransaction{}
	trx.Actions = append(trx.Actions, &onBlockAction)
	trx.SetReferenceBlock(c.Head.ID)
	in := c.Pending.PendingBlockState.Header.Timestamp + 999999
	trx.Expiration = common.TimePointSec(in)
	log.Error("getOnBlockTransaction trx.Expiration:", trx)
	return trx
}
func (c *Controller) SkipDbSession(bs types.BlockStatus) bool {
	considerSkipping := (bs == types.BlockStatus(IRREVERSIBLE))
	//log.Info("considerSkipping:", considerSkipping)
	return considerSkipping
}

func (c *Controller) SkipDbSessions() bool {
	if c.Pending == nil {
		return c.SkipDbSession(c.Pending.BlockStatus)
	} else {
		return false
	}
}

func (c *Controller) SkipTrxChecks() (b bool) {
	b = c.LightValidationAllowed(c.Config.disableReplayOpts)
	return
}

func (c *Controller) LightValidationAllowed(dro bool) (b bool) {
	if c.Pending != nil || c.InTrxRequiringChecks {
		return false
	}

	pbStatus := c.Pending.BlockStatus
	considerSkippingOnReplay := (pbStatus == types.Irreversible || pbStatus == types.Validated) && !dro

	considerSkippingOnvalidate := (pbStatus == types.Complete && c.Config.blockValidationMode == LIGHT)

	return considerSkippingOnReplay || considerSkippingOnvalidate
}

func (c *Controller) IsProducingBlock() bool {
	if c.Pending == nil {
		return false
	}
	return c.Pending.BlockStatus == types.Incomplete
}

func (c *Controller) IsResourceGreylisted(name *common.AccountName) bool {
	for _, account := range c.Config.resourceGreylist {
		if &account == name {
			return true
		}
	}
	return false
}

func (c *Controller) PendingBlockState() *types.BlockState {
	if c.Pending != nil {
		return c.Pending.PendingBlockState
	}
	return &types.BlockState{}
}

func (c *Controller) PendingBlockTime() common.TimePoint {
	if c.Pending == nil {
		log.Error("PendingBlockTime is error", "no pending block")
	}
	return c.Pending.PendingBlockState.Header.Timestamp.ToTimePoint()
}

func Close(db *database.LDataBase, session *database.Session) {
	//session.close()
	db.Close()
}

func (c *Controller) initConfig() *Controller {
	c.Config = Config{
		blocksDir:           common.DefaultConfig.DefaultBlocksDirName,
		stateDir:            common.DefaultConfig.DefaultStateDirName,
		stateSize:           common.DefaultConfig.DefaultStateSize,
		stateGuardSize:      common.DefaultConfig.DefaultStateGuardSize,
		reversibleCacheSize: common.DefaultConfig.DefaultReversibleCacheSize,
		reversibleGuardSize: common.DefaultConfig.DefaultReversibleGuardSize,
		readOnly:            false,
		forceAllChecks:      false,
		disableReplayOpts:   false,
		contractsConsole:    false,
		//vmType:              common.DefaultConfig.DefaultWasmRuntime, //TODO
		readMode:            SPECULATIVE,
		blockValidationMode: FULL,
	}
	return c
}

func (c *Controller) GetUnAppliedTransactions() *[]types.TransactionMetadata {
	result := []types.TransactionMetadata{}
	if c.ReadMode == SPECULATIVE {
		for _, v := range c.UnAppliedTransactions {
			result = append(result, v)
		}
	} else {
		fmt.Println("not empty unapplied_transactions in non-speculative mode")
	}
	return &result
}

func (c *Controller) DropUnAppliedTransaction(metadata *types.TransactionMetadata) {
	delete(c.UnAppliedTransactions, crypto.Sha256(metadata.SignedID))
}

func (c *Controller) DropAllUnAppliedTransactions() {
	c.UnAppliedTransactions = nil
}
func (c *Controller) GetScheduledTransactions() *[]common.TransactionIdType {
	//TODO add generated_transaction_object
	//c.Db.Find("",)
	/*const auto& idx = db().get_index<generated_transaction_multi_index,by_delay>();

	vector<transaction_id_type> result;

	static const size_t max_reserve = 64;
	result.reserve(std::min(idx.size(), max_reserve));

	auto itr = idx.begin();
	while( itr != idx.end() && itr->delay_until <= pending_block_time() ) {
		result.emplace_back(itr->trx_id);
		++itr;
	}*/
	return nil
}

func (c *Controller) PushScheduledTransactionById(sheduled common.TransactionIdType,
	deadLine common.TimePoint,
	billedCpuTimeUs uint32) *types.TransactionTrace {

	gto := types.GetGTOByTrxId(c.DB, sheduled)

	if gto == nil {
		fmt.Println("unknown_transaction_exception", "unknown transaction")
	}

	return c.PushScheduledTransactionByObject(*gto, deadLine, billedCpuTimeUs, false)
}

func (c *Controller) PushScheduledTransactionByObject(gto types.GeneratedTransactionObject,
	deadLine common.TimePoint,
	billedCpuTimeUs uint32,
	explicitBilledCpuTime bool) *types.TransactionTrace {
	/*var undoSession database.Session
	if !c.SkipDbSessions(){
		undoSession = c.DB.StartSession()
	}
	err := c.DB.Find("ByExpiration", ,&gto)
	if err != nil {
		fmt.Println("GetGeneratedTransactionObjectByExpiration is error :", err.Error())
	}*/

	//undo_session := c.DB.StartSession()
	gtrx := types.GeneratedTransactions(&gto)

	c.RemoveScheduledTransaction(&gto)

	if gtrx.DelayUntil <= c.PendingBlockTime() {
		fmt.Println("this transaction isn't ready")
		return nil
	}
	dtrx := types.SignedTransaction{}

	err := rlp.DecodeBytes(gtrx.PackedTrx, &dtrx)
	if err != nil {
		fmt.Println("PushScheduleTransaction1 DecodeBytes is error :", err.Error())
	}

	trx := &types.TransactionMetadata{}
	trx.Trx = &dtrx
	trx.Accepted = true
	trx.Scheduled = true

	trace := &types.TransactionTrace{}
	fmt.Println(trace)
	/*if( gtrx.expiration < c.pending_block_time() ) {
		trace = std::make_shared<transaction_trace>();
		trace->id = gtrx.trx_id;
		trace->block_num = c.pending_block_state()->block_num;
		trace->block_time = c.pending_block_time();
		trace->producer_block_id = c.pending_producer_block_id();
		trace->scheduled = true;
		trace->receipt = push_receipt( gtrx.trx_id, transaction_receipt::expired, billed_cpu_time_us, 0 ); // expire the transaction
		emit( c.accepted_transaction, trx );
		emit( c.applied_transaction, trace );
		undo_session.squash();
		return trace;
	}*/
	c.InTrxRequiringChecks = true
	//cpuTimeToBillUs := billedCpuTimeUs
	trxContext := NewTransactionContext(c, &dtrx, gtrx.TrxId, common.Now())
	trxContext.Leeway = common.Milliseconds(0)
	trxContext.Deadline = deadLine
	trxContext.ExplicitBilledCpuTime = explicitBilledCpuTime
	//trxContext.BilledCpuTimeUs = billedCpuTimeUs
	trace = &trxContext.Trace

	//try.CatchOrFinally{
	trxContext.InitForDeferredTrx(gtrx.Published)
	//}
	//TODO 2018-10-13

	fmt.Println(dtrx, trx) //TODO
	return nil
}

func (c *Controller) RemoveScheduledTransaction(gto *types.GeneratedTransactionObject) {
	c.ResourceLimists.AddPendingRamUsage(gto.Payer, int64(9)+int64(len(gto.PackedTrx))) //TODO billable_size_v
	c.DB.Remove(gto)
}

func failureIsSubjective() bool {
	/*code := e.code()
	return    (code == subjective_block_production_exception::code_value)
	|| (code == block_net_usage_exceeded::code_value)
	|| (code == greylist_net_usage_exceeded::code_value)
	|| (code == block_cpu_usage_exceeded::code_value)
	|| (code == greylist_cpu_usage_exceeded::code_value)
	|| (code == deadline_exception::code_value)
	|| (code == leeway_deadline_exception::code_value)
	|| (code == actor_whitelist_exception::code_value)
	|| (code == actor_blacklist_exception::code_value)
	|| (code == contract_whitelist_exception::code_value)
	|| (code == contract_blacklist_exception::code_value)
	|| (code == action_blacklist_exception::code_value)
	|| (code == key_blacklist_exception::code_value)*/
	return false
}

func (c *Controller) setActionMaerkle() {
	actionDigests := make([]crypto.Sha256, len(c.Pending.Actions))
	for _, b := range c.Pending.Actions {
		actionDigests = append(actionDigests, crypto.Hash256(b.ActDigest))
	}
	c.Pending.PendingBlockState.Header.ActionMRoot = common.CheckSum256Type(types.Merkle(actionDigests))
}

func (c *Controller) setTrxMerkle() {
	actionDigests := make([]crypto.Sha256, len(c.Pending.Actions))
	for _, b := range c.Pending.PendingBlockState.SignedBlock.Transactions {
		actionDigests = append(actionDigests, crypto.Hash256(b.Digest()))
	}
	c.Pending.PendingBlockState.Header.TransactionMRoot = common.CheckSum256Type(types.Merkle(actionDigests))
}
func (c *Controller) FinalizeBlock() {

	exception.EosAssert(c.Pending != nil, &exception.BlockValidateException{}, "it is not valid to finalize when there is no pending block")
	c.ResourceLimists.ProcessAccountLimitUpdates()
	chainConfig := c.GetGlobalProperties().Configuration
	cpuTarget := common.EosPercent(uint64(chainConfig.MaxBlockCpuUsage), chainConfig.TargetBlockCpuUsagePct)
	m := uint32(1000)
	cpu := types.ElasticLimitParameters{}
	cpu.Target = cpuTarget
	cpu.Max = uint64(chainConfig.MaxBlockCpuUsage)
	cpu.Periods = common.DefaultConfig.BlockCpuUsageAverageWindowMs / uint32(common.DefaultConfig.BlockIntervalMs)
	cpu.MaxMultiplier = m

	cpu.ContractRate.Numerator = 99
	cpu.ExpandRate.Denominator = 100

	net := types.ElasticLimitParameters{}
	netTarget := common.EosPercent(uint64(chainConfig.MaxBlockNetUsage), chainConfig.TargetBlockNetUsagePct)
	net.Target = netTarget
	net.Max = uint64(chainConfig.MaxBlockNetUsage)
	net.Periods = common.DefaultConfig.BlockSizeAverageWindowMs / uint32(common.DefaultConfig.BlockIntervalMs)
	net.MaxMultiplier = m

	net.ContractRate.Numerator = 99
	net.ExpandRate.Denominator = 100
	c.ResourceLimists.SetBlockParameters(cpu, net)

	c.setActionMaerkle()
}

func (c *Controller) SignBlock(callBack interface{}) {}

func (c *Controller) CommitBlock(addToForkDb bool) {
	//defer c.MakeBlockRestorePoint()
	//try{
	if addToForkDb {
		c.Pending.PendingBlockState.Validated = true
		newBsp := c.ForkDB.AddBlockState(c.Pending.PendingBlockState)
		//emit(self.accepted_block_header, pending->_pending_block_state)
		c.Head = c.ForkDB.Header()
		exception.EosAssert(newBsp == c.Head, &exception.ForkDatabaseException{}, "committed block did not become the new head in fork database")
	}

	if !c.RePlaying {
		ubo := entity.ReversibleBlockObject{}
		ubo.BlockNum = c.Pending.PendingBlockState.BlockNum
		ubo.SetBlock(c.Pending.PendingBlockState.SignedBlock)
		c.DB.Insert(ubo)
	}
	//emit( self.accepted_block, pending->_pending_block_state )
	//catch(){
	// reset_pending_on_exit.cancel();
	//         abort_block();
	//throw;
	// }
}

func (c *Controller) MakeBlockRestorePoint() {

}
func (c *Controller) PushBlock(sbp *types.SignedBlock, status types.BlockStatus) {} //status default value block_status s = block_status::complete

func (c *Controller) PushConfirnation(hc types.HeaderConfirmation) {}

func (c *Controller) DataBase() *database.LDataBase {
	return c.DB
}

func (c *Controller) ForkDataBase() *types.ForkDatabase {
	return c.ForkDB
}

func (c *Controller) GetAccount(name common.AccountName) *types.AccountObject {
	accountObj := types.AccountObject{}
	//accountObj.Name = name
	err := c.DB.Find("Name", name, accountObj)
	if err != nil {
		fmt.Println("GetAccount is error :", err)
	}
	return &accountObj
}

func (c *Controller) GetAuthorizationManager() *AuthorizationManager { return c.Authorization }

func (c *Controller) GetMutableAuthorizationManager() *AuthorizationManager { return c.Authorization }

//c++ flat_set<account_name> map[common.AccountName]interface{}
func (c *Controller) GetActorWhiteList() *map[common.AccountName]struct{} {
	return &c.Config.ActorWhitelist
}

func (c *Controller) GetActorBlackList() *map[common.AccountName]struct{} {
	return &c.Config.ActorBlacklist
}

func (c *Controller) GetContractWhiteList() *map[common.AccountName]struct{} {
	return &c.Config.ContractWhitelist
}

func (c *Controller) GetContractBlackList() *map[common.AccountName]struct{} {
	return &c.Config.ContractBlacklist
}

func (c *Controller) GetActionBlockList() *map[common.Pair]struct{} { return &c.Config.ActionBlacklist }

func (c *Controller) GetKeyBlackList() *map[ecc.PublicKey]struct{} { return &c.Config.KeyBlacklist }

func (c *Controller) SetActorWhiteList(params *map[common.AccountName]struct{}) {
	c.Config.ActorWhitelist = *params
}

func (c *Controller) SetActorBlackList(params *map[common.AccountName]struct{}) {
	c.Config.ActorBlacklist = *params
}

func (c *Controller) SetContractWhiteList(params *map[common.AccountName]struct{}) {
	c.Config.ContractWhitelist = *params
}

func (c *Controller) SetContractBlackList(params *map[common.AccountName]struct{}) {
	c.Config.ContractBlacklist = *params
}

func (c *Controller) SetActionBlackList(params *map[common.Pair]struct{}) {
	c.Config.ActionBlacklist = *params
}

func (c *Controller) SetKeyBlackList(params *map[ecc.PublicKey]struct{}) {
	c.Config.KeyBlacklist = *params
}

func (c *Controller) HeadBlockNum() uint32 { return 0 }

func (c *Controller) HeadBlockTime() common.TimePoint { return 0 }

func (c *Controller) HeadBlockId() common.BlockIdType { return common.BlockIdType{} }

func (c *Controller) HeadBlockProducer() common.AccountName { return 0 }

func (c *Controller) HeadBlockHeader() *types.BlockHeader { return nil }

func (c *Controller) HeadBlockState() types.BlockState { return types.BlockState{} }

func (c *Controller) ForkDbHeadBlockNum() uint32 { return 0 }

func (c *Controller) ForkDbHeadBlockId() common.BlockIdType { return common.BlockIdType{} }

func (c *Controller) ForkDbHeadBlockTime() common.TimePoint { return 0 }

func (c *Controller) ForkDbHeadBlockProducer() common.AccountName { return 0 }

func (c *Controller) ActiveProducers() *types.ProducerScheduleType { return nil }

func (c *Controller) PendingProducers() *types.ProducerScheduleType { return nil }

func (c *Controller) ProposedProducers() types.ProducerScheduleType {
	return types.ProducerScheduleType{}
}

func (c *Controller) LastIrreversibleBlockNum() uint32 { return 0 }

func (c *Controller) LastIrreversibleBlockId() common.BlockIdType { return common.BlockIdType{} }

func (c *Controller) FetchBlockByNumber(blockNum uint32) *types.SignedBlock {
	blkState := c.ForkDB.GetBlockInCurrentChainByNum(blockNum)
	if blkState != nil {
		return blkState.SignedBlock
	}
	//TODO blog
	return &types.SignedBlock{}
}

func (c *Controller) FetchBlockById(id common.BlockIdType) *types.SignedBlock {
	state := c.ForkDB.GetBlock(&id)
	if state != nil {
		return state.SignedBlock
	}
	bptr := c.FetchBlockByNumber(types.NumFromID(id))
	if bptr != nil && bptr.BlockID() == id {
		return bptr
	}
	return &types.SignedBlock{}
}

func (c *Controller) FetchBlockStateByNumber(blockNum uint32) *types.BlockState {
	return c.ForkDB.GetBlockInCurrentChainByNum(blockNum)
}

func (c *Controller) FetchBlockStateById(id common.BlockIdType) *types.BlockState {
	return c.ForkDB.GetBlock(&id)
}

func (c *Controller) GetBlcokIdForNum(blockNum uint32) common.BlockIdType { return common.BlockIdType{} }

func (c *Controller) CheckContractList(code common.AccountName) {
	if len(c.Config.ContractWhitelist) > 0 {
		if _, ok := c.Config.ContractWhitelist[code]; !ok {
			fmt.Println("account is not on the contract whitelist", code)
			return
		}
		/*EOS_ASSERT( conf.contract_whitelist.find( code ) != conf.contract_whitelist.end(),
			contract_whitelist_exception,
			"account '${code}' is not on the contract whitelist", ("code", code)
		);*/
	} else if len(c.Config.ContractBlacklist) > 0 {
		if _, ok := c.Config.ContractBlacklist[code]; ok {
			fmt.Println("account is on the contract blacklist", code)
			return
		}
		/*EOS_ASSERT( conf.contract_blacklist.find( code ) == conf.contract_blacklist.end(),
			contract_blacklist_exception,
			"account '${code}' is on the contract blacklist", ("code", code)
		);*/
	}
}

func (c *Controller) CheckActionList(code common.AccountName, action common.ActionName) {
	if len(c.Config.ActionBlacklist) > 0 {
		abl := common.MakePair(code, action)
		//key := Hash(abl)
		if _, ok := c.Config.ActionBlacklist[abl]; ok {
			fmt.Println("action '${code}::${action}' is on the action blacklist")
			return
		}
		/*EOS_ASSERT( conf.action_blacklist.find( std::make_pair(code, action) ) == conf.action_blacklist.end(),
			action_blacklist_exception,
			"action '${code}::${action}' is on the action blacklist",
			("code", code)("action", action)
		);*/
	}
}

func (c *Controller) CheckKeyList(key *ecc.PublicKey) {
	if len(c.Config.KeyBlacklist) > 0 {
		if _, ok := c.Config.KeyBlacklist[*key]; ok {
			fmt.Println("public key '${key}' is on the key blacklist", key)
			return
		}
		/*EOS_ASSERT( conf.key_blacklist.find( key ) == conf.key_blacklist.end(),
			key_blacklist_exception,
			"public key '${key}' is on the key blacklist",
			("key", key)
		);*/
	}
}

func (c *Controller) IsProducing() bool { return false }

func (c *Controller) IsRamBillingInNotifyAllowed() bool { return false }

func (c *Controller) AddResourceGreyList(name *common.AccountName) {}

func (c *Controller) RemoveResourceGreyList(name *common.AccountName) {}

func (c *Controller) IsResourceGreyListed(name *common.AccountName) bool { return false }

func (c *Controller) GetResourceGreyList() *map[common.AccountName]interface{} { return nil }

func (c *Controller) ValidateReferencedAccounts(t *types.Transaction) {}

func (c *Controller) ValidateExpiration(t *types.Transaction) {}

func (c *Controller) ValidateTapos(t *types.Transaction) {}

func (c *Controller) ValidateDbAvailableSize() {}

func (c *Controller) ValidateReversibleAvailableSize() {}

func (c *Controller) IsKnownUnexpiredTransaction(id *common.TransactionIdType) bool { return false }

func (c *Controller) SetProposedProducers(producers []types.ProducerKey) int64 { return 0 }

func (c *Controller) SkipAuthCheck() bool { return false }

func (c *Controller) ContractsConsole() bool { return false }

func (c *Controller) GetChainId() common.ChainIdType { return common.ChainIdType{} }

func (c *Controller) GetReadMode() DBReadMode { return 0 }

func (c *Controller) GetValidationMode() ValidationMode { return 0 }

func (c *Controller) SetSubjectiveCpuLeeway(leeway common.Microseconds) {}

func (c *Controller) PendingProducerBlockId() common.BlockIdType {
	//EOS_ASSERT( my->pending, block_validate_exception, "no pending block" )
	return c.Pending.ProducerBlockId
}

func (c *Controller) FindApplyHandler(receiver common.AccountName,
	scope common.AccountName,
	act common.ActionName) func(*ApplyContext) {

	handlerKey := NewHandlerKey(common.ScopeName(scope), act)
	secend, ok := c.ApplyHandlers[receiver]
	if ok {
		handler, success := secend[handlerKey]
		fmt.Println("find secend:", success)
		if success {
			fmt.Println("-=-=-=-=-=-=-=-==-=-=-=-=-=-=", handler)
			return handler
		}
	}
	return nil
}

func (c *Controller) GetWasmInterface() *exec.WasmInterface {
	return c.WasmIf
}

func (c *Controller) GetAbiSerializer(name common.AccountName,
	maxSerializationTime common.Microseconds) types.AbiSerializer {
	return types.AbiSerializer{}
}

func (c *Controller) ToVariantWithAbi(obj interface{}, maxSerializationTime common.Microseconds) {}

var readycontroller chan bool //TODO test code

func (c *Controller) CreateNativeAccount(name common.AccountName, owner types.Authority, active types.Authority, isPrivileged bool) {
	account := types.AccountObject{}
	account.Name = name
	account.CreationDate = common.BlockTimeStamp(c.Config.genesis.InitialTimestamp)
	account.Privileged = isPrivileged
	if name == common.AccountName(common.DefaultConfig.SystemAccountName) {
		abiDef := types.AbiDef{}
		account.SetAbi(EosioContractAbi(abiDef))
	}
	c.DB.Insert(account)

	aso := types.AccountSequenceObject{}
	aso.Name = name
	c.DB.Insert(aso)

	ownerPermission := c.Authorization.CreatePermission(name, common.PermissionName(common.DefaultConfig.OwnerName), 0, owner, c.Config.genesis.InitialTimestamp)

	activePermission := c.Authorization.CreatePermission(name, common.PermissionName(common.DefaultConfig.ActiveName), PermissionIdType(ownerPermission.ID), active, c.Config.genesis.InitialTimestamp)

	c.ResourceLimists.InitializeAccount(name)
	ramDelta := uint64(common.DefaultConfig.OverheadPerRowPerIndexRamBytes) //TODO c++ reference int64 but statement uint32
	ramDelta += 2 * common.BillableSizeV("permission_object")               //::billable_size_v<permission_object>
	ramDelta += ownerPermission.Auth.GetBillableSize()
	ramDelta += activePermission.Auth.GetBillableSize()
	c.ResourceLimists.AddPendingRamUsage(name, int64(ramDelta))
	c.ResourceLimists.VerifyAccountRamUsage(name)
}

func initResource(c *Controller, ready chan bool) {
	<-ready
	//con.Blog
	c.ForkDB = types.GetForkDbInstance(common.DefaultConfig.DefaultStateDirName)
	c.ResourceLimists = GetResourceLimitsManager()
	c.Authorization = GetAuthorizationManager()
}

/*    about chain

signal<void(const signed_block_ptr&)>         pre_accepted_block;
signal<void(const block_state_ptr&)>          accepted_block_header;
signal<void(const block_state_ptr&)>          accepted_block;
signal<void(const block_state_ptr&)>          irreversible_block;
signal<void(const transaction_metadata_ptr&)> accepted_transaction;
signal<void(const transaction_trace_ptr&)>    applied_transaction;
signal<void(const header_confirmation&)>      accepted_confirmation;
signal<void(const int&)>                      bad_alloc;*/

//c++ pair<scope_name,action_name>
type HandlerKey struct {
	//handMap map[common.AccountName]common.ActionName
	scope  common.ScopeName
	action common.ActionName
}

func NewHandlerKey(scopeName common.ScopeName, actionName common.ActionName) HandlerKey {
	hk := HandlerKey{scopeName, actionName}
	return hk
}

func (c *Controller) clearExpiredInputTransactions() {
	/*
		auto& transaction_idx = db.get_mutable_index<transaction_multi_index>();
		const auto& dedupe_index = transaction_idx.indices().get<by_expiration>();
		auto now = c.pending_block_time();
		while( (!dedupe_index.empty()) && ( now > fc::time_point(dedupe_index.begin()->expiration) ) ) {
		transaction_idx.remove(*dedupe_index.begin());
		}
	*/
}

func (c *Controller) CheckActorList(actors []common.AccountName) {
	if len(c.Config.ActorWhitelist) > 0 {
		//excluded :=make(map[common.AccountName]struct{})
		//set
		for _, an := range actors {
			if _, ok := c.Config.ActorWhitelist[an]; !ok {
				fmt.Println("authorizing actor(s) in transaction are not on the actor whitelist:", an)
				return
			}
		}
		/*EOS_ASSERT( excluded.size() == 0, actor_whitelist_exception,
			"authorizing actor(s) in transaction are not on the actor whitelist: ${actors}",
			("actors", excluded)
		)*/
	} else if len(c.Config.ActorBlacklist) > 0 {
		//black :=make(map[common.AccountName]struct{})
		//set
		for _, an := range actors {
			if _, ok := c.Config.ActorBlacklist[an]; ok {
				fmt.Println("authorizing actor(s) in transaction are not on the actor blacklist:", an)
				return
			}
		}
		/*EOS_ASSERT( blacklisted.size() == 0, actor_blacklist_exception,
			"authorizing actor(s) in transaction are on the actor blacklist: ${actors}",
			("actors", blacklisted)
		)*/
	}
}
func (c *Controller) UpdateProducersAuthority() {
	/*producers := c.Pending.PendingBlockState.ActiveSchedule.Producers
	updatePermission :=*/
}

/*
//for ActionBlacklist
type ActionBlacklistParam struct {
	AccountName common.AccountName
	ActionName  common.ActionName
}

func Hash(abp ActionBlacklistParam) string {
	return crypto.Hash256(abp).String()
}





type applyCon struct {
	handlerKey   map[common.AccountName]common.AccountName //c++ pair<scope_name,action_name>
	applyContext func(*ApplyContext)
}

//apply_context
type ApplyHandler struct {
	applyHandler map[common.AccountName]applyCon
	receiver     common.AccountName
}*/
