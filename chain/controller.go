package chain

import (
	"os"
	"time"

	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/entity"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/wasmgo"
)

//var readycontroller chan bool //TODO test code

/*var PreAcceptedBlock chan *types.SignedBlock
var AcceptedBlockdHeader chan *types.BlockState
var AcceptedBlock chan *types.BlockState
var IrreversibleBlock chan *types.BlockState
var AcceptedTransaction chan *types.TransactionMetadata
var AppliedTransaction chan *types.TransactionTrace
var AcceptedConfirmation chan *types.HeaderConfirmation
var BadAlloc chan *int*/

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
	ActorWhitelist          common.FlatSet //common.AccountName
	ActorBlacklist          common.FlatSet //common.AccountName
	ContractWhitelist       common.FlatSet //common.AccountName
	ContractBlacklist       common.FlatSet //common.AccountName]struct{}
	ActionBlacklist         common.FlatSet //common.Pair //see actionBlacklist
	KeyBlacklist            common.FlatSet
	blocksDir               string
	stateDir                string
	stateSize               uint64
	stateGuardSize          uint64
	reversibleCacheSize     uint64
	reversibleGuardSize     uint64
	readOnly                bool
	forceAllChecks          bool
	disableReplayOpts       bool
	disableReplay           bool
	contractsConsole        bool
	allowRamBillingInNotify bool
	genesis                 types.GenesisState
	vmType                  wasmgo.WasmGo
	readMode                DBReadMode
	blockValidationMode     ValidationMode
	resourceGreylist        map[common.AccountName]struct{}
	trustedProducers        map[common.AccountName]struct{}
}

var isActiveController bool //default value false ;Does the process include control ;

var instance *Controller

type v func(ctx *ApplyContext)

type Controller struct {
	DB                             database.DataBase
	UndoSession                    database.Session
	ReversibleBlocks               database.DataBase
	Blog                           *BlockLog
	Pending                        *PendingState
	Head                           *types.BlockState
	ForkDB                         *ForkDatabase
	WasmIf                         *wasmgo.WasmGo
	ResourceLimits                 *ResourceLimitsManager
	Authorization                  *AuthorizationManager
	Config                         Config //local	Config
	ChainID                        common.ChainIdType
	RePlaying                      bool
	ReplayHeadTime                 common.TimePoint //optional<common.Tstamp>
	ReadMode                       DBReadMode
	InTrxRequiringChecks           bool                //if true, checks that are normally skipped on replay (e.g. auth checks) cannot be skipped
	SubjectiveCupLeeway            common.Microseconds //optional<common.Tstamp>
	TrustedProducerLightValidation bool                //default value false
	ApplyHandlers                  map[string]v
	UnAppliedTransactions          map[crypto.Sha256]types.TransactionMetadata
}

func GetControllerInstance() *Controller {
	if !isActiveController {
		validPath()
		instance = newController()
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
				log.Error("controller validPath mkdir failed![%v]\n", err)
			} else {
				log.Error("controller validPath mkdir success![%v]\n", d)
			}
		}
	}
}
func newController() *Controller {
	isActiveController = true //controller is active
	//init db
	db, err := database.NewDataBase(common.DefaultConfig.DefaultStateDirName)
	if err != nil {
		log.Error("newController is error :", err)
		return nil
	}
	//defer db.Close()

	//init ReversibleBlocks
	//reversibleDir := common.DefaultConfig.DefaultBlocksDirName + "/" + common.DefaultConfig.DefaultReversibleBlocksDirName
	reversibleDB, err := database.NewDataBase(common.DefaultConfig.DefaultReversibleBlocksDirName)
	if err != nil {
		log.Error("newController init reversibleDB is error", err)
	}
	con := &Controller{InTrxRequiringChecks: false, RePlaying: false, TrustedProducerLightValidation: false}
	con.DB = db
	con.ReversibleBlocks = reversibleDB

	con.Blog = NewBlockLog(common.DefaultConfig.DefaultBlocksDirName)

	con.ForkDB, _ = newForkDatabase(common.DefaultConfig.DefaultBlocksDirName, common.DefaultConfig.ForkDBName, true)
	con.ChainID = types.GetGenesisStateInstance().ComputeChainID()

	con.initConfig()
	con.ReadMode = con.Config.readMode
	con.ApplyHandlers = make(map[string]v)
	con.WasmIf = wasmgo.NewWasmGo()

	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("newaccount")), applyEosioNewaccount)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("setcode")), applyEosioSetcode)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("setabi")), applyEosioSetabi)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("updateauth")), applyEosioUpdateauth)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("deleteauth")), applyEosioDeleteauth)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("linkauth")), applyEosioUnlinkauth)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("unlinkauth")), applyEosioLinkauth)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("canceldelay")), applyEosioCanceldalay)

	//IrreversibleBlock.connect()
	//readycontroller = make(chan bool)
	//go initResource(con, readycontroller)
	con.Pending = &PendingState{}
	con.ResourceLimits = newResourceLimitsManager(con)
	con.Authorization = newAuthorizationManager(con)
	con.initialize()
	return con
}

/*func initResource(c *Controller, ready chan bool) {
	<-ready
	//con.Blog
	//c.ForkDB = types.GetForkDbInstance(common.DefaultConfig.DefaultStateDirName)

	c.initialize()
}*/

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
				log.Info("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func (c *Controller) OnIrreversible(s *types.BlockState) {
	if !common.Empty(c.Blog.head) {
		c.Blog.ReadHead()
	}
	logHead := c.Blog.head
	EosAssert(common.Empty(logHead), &BlockLogException{}, "block log head can not be found")
	lhBlockNum := logHead.BlockNumber()
	c.DB.Commit(int64(s.BlockNum))
	if s.BlockNum <= lhBlockNum {
		return
	}
	EosAssert(s.BlockNum-1 == lhBlockNum, &UnlinkableBlockException{}, "unlinkable block", s.BlockNum, lhBlockNum)
	EosAssert(s.Header.Previous == logHead.BlockID(), &UnlinkableBlockException{}, "irreversible doesn't link to block log head")
	c.Blog.Append(s.SignedBlock)
	bs := types.BlockState{}
	ubi, err := c.ReversibleBlocks.GetIndex("byNum", &bs)
	if err != nil {
		log.Error("Controller OnIrreversible ReversibleBlocks.GetIndex is error:", err)
	}
	itr := ubi.Begin()
	tbs := types.BlockState{}
	err = itr.Data(&tbs)
	for itr != ubi.End() && tbs.BlockNum <= s.BlockNum {
		c.ReversibleBlocks.Remove(itr)
		itr = ubi.Begin()
	}
	if c.ReadMode == IRREVERSIBLE {
		c.applyBlock(s.SignedBlock, types.Complete)
		c.ForkDB.MarkInCurrentChain(s, true)
		c.ForkDB.SetValidity(s, true)
		c.Head = s
	}
	//emit( self.irreversible_block, s ) 	//TODO
}

func (c *Controller) PopBlock() {
	prev := c.ForkDB.GetBlock(&c.Head.Header.Previous)
	r := entity.ReversibleBlockObject{}
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
	c.UndoSession.Undo()
}

func (c *Controller) SetApplayHandler(receiver common.AccountName, contract common.AccountName, action common.ActionName, handler func(a *ApplyContext)) {
	handlerKey := receiver + contract + action
	//fmt.Println("handlerKey----:",handlerKey.String())
	c.ApplyHandlers[handlerKey.String()] = handler
}

func (c *Controller) FindApplyHandler(receiver common.AccountName,
	scope common.AccountName,
	act common.ActionName) func(*ApplyContext) {
	handlerKey := receiver + scope + act

	handler, ok := c.ApplyHandlers[handlerKey.String()]
	if ok {
		return handler
	}
	return nil
}

func (c *Controller) AbortBlock() {
	if common.Empty(c.Pending) {
		if c.ReadMode == SPECULATIVE {
			trx := append(c.Pending.PendingBlockState.Trxs)
			step := 0
			for ; step < len(trx); step++ {
				c.UnAppliedTransactions[crypto.Sha256(trx[step].SignedID)] = *trx[step]
			}
		}
	}
}
func (c *Controller) StartBlock(when types.BlockTimeStamp, confirmBlockCount uint16) {
	pbi := common.BlockIdType(*crypto.NewSha256Nil())
	c.startBlock(when, confirmBlockCount, types.Incomplete, &pbi)
	c.ValidateDbAvailableSize()
}
func (c *Controller) startBlock(when types.BlockTimeStamp, confirmBlockCount uint16, s types.BlockStatus, producerBlockId *common.BlockIdType) {
	//fmt.Println(c.Config)
	EosAssert(nil != c.Pending, &BlockValidateException{}, "pending block already exists")
	defer func() {
		if PendingValid {
			c.Pending.Reset()
		}
	}()
	if c.SkipDbSession(s) {
		/*EosAssert( c.DB.revision() == head->block_num, database_exception, "db revision is not on par with head block",
		("db.revision()", db.revision())("controller_head_block", head->block_num)("fork_db_head_block", fork_db.head()->block_num) )*/
		c.Pending = NewPendingState(c.DB)
	} else {
		//c.Pending = types.GetInstance()
	}

	c.Pending.BlockStatus = s
	c.Pending.ProducerBlockId = *producerBlockId
	c.Pending.PendingBlockState = types.NewBlockState2(&c.Head.BlockHeaderState, when) //TODO std::make_shared<block_state>( *head, when ); // promotes pending schedule (if any) to active
	//c.Pending.PendingBlockState.SignedBlock.Timestamp = when
	c.Pending.PendingBlockState.InCurrentChain = true
	c.Pending.PendingBlockState.SetConfirmed(confirmBlockCount)
	wasPendingPromoted := c.Pending.PendingBlockState.MaybePromotePending()
	log.Info("wasPendingPromoted:%t", wasPendingPromoted)
	if c.ReadMode == DBReadMode(SPECULATIVE) || c.Pending.BlockStatus != types.BlockStatus(types.Incomplete) {
		gpo := c.GetGlobalProperties()
		if (!common.Empty(gpo.ProposedScheduleBlockNum) && gpo.ProposedScheduleBlockNum <= c.Pending.PendingBlockState.DposIrreversibleBlocknum) &&
			(len(c.Pending.PendingBlockState.PendingSchedule.Producers) == 0) &&
			(!wasPendingPromoted) {
			if !c.RePlaying {
				tmp := gpo.ProposedSchedule.ProducerScheduleType()
				ps := types.SharedProducerScheduleType{}
				ps.Version = tmp.Version
				ps.Producers = tmp.Producers
				c.Pending.PendingBlockState.SetNewProducers(ps)
			}

			c.DB.Modify(&gpo, func(i *entity.GlobalPropertyObject) {
				i.ProposedScheduleBlockNum = 1
				i.ProposedSchedule.Clear()
			})
		}
		//try.Try(func() {
		signedTransaction := c.GetOnBlockTransaction()
		onbtrx := types.TransactionMetadata{Trx: &signedTransaction}
		onbtrx.Implicit = true
		//TODO defer
		defer func(b bool) {
			c.InTrxRequiringChecks = b
		}(c.InTrxRequiringChecks)
		c.InTrxRequiringChecks = true
		c.pushTransaction(&onbtrx, common.MaxTimePoint(), gpo.Configuration.MinTransactionCpuUsage, true)
		/*}).Catch(func(e Exception) {
			//TODO
			fmt.Println("Controller StartBlock exception:",e.Message())
		})*/

		c.clearExpiredInputTransactions()
		c.updateProducersAuthority()
	}
	PendingValid = true

}

func (c *Controller) pushReceipt(trx interface{}, status types.TransactionStatus, cpuUsageUs uint64, netUsage uint64) *types.TransactionReceipt {
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
	EosAssert(netUsageWords*8 == netUsage, &TransactionException{}, "net_usage is not divisible by 8")
	c.Pending.PendingBlockState.SignedBlock.Transactions = append(c.Pending.PendingBlockState.SignedBlock.Transactions, trxReceipt)
	trxReceipt.CpuUsageUs = uint32(cpuUsageUs)
	trxReceipt.NetUsageWords = uint32(netUsageWords)
	trxReceipt.Status = types.TransactionStatus(status)
	return &trxReceipt
}

func (c *Controller) PushTransaction(trx *types.TransactionMetadata, deadLine common.TimePoint, billedCpuTimeUs uint32) *types.TransactionTrace {
	c.ValidateDbAvailableSize()
	EosAssert(c.GetReadMode() != READONLY, &TransactionTypeException{}, "push transaction not allowed in read-only mode")
	EosAssert(!common.Empty(trx) && !trx.Implicit && !trx.Scheduled, &TransactionTypeException{}, "Implicit/Scheduled transaction not allowed")
	return c.pushTransaction(trx, deadLine, billedCpuTimeUs, billedCpuTimeUs > 0)
}

func (c *Controller) pushTransaction(trx *types.TransactionMetadata, deadLine common.TimePoint, billedCpuTimeUs uint32, explicitBilledCpuTime bool) (trxTrace *types.TransactionTrace) {
	EosAssert(deadLine != common.TimePoint(0), &TransactionException{}, "deadline cannot be uninitialized")

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
		c.CheckActorList(&trxContext.BillToAccounts)
	}
	trxContext.Delay = common.Microseconds(trx.Trx.DelaySec)
	checkTime := func() {}
	if !c.SkipAuthCheck() && !trx.Implicit {
		c.Authorization.CheckAuthorization(trx.Trx.Actions,
			trx.RecoverKeys(&c.ChainID),
			&common.FlatSet{},
			trxContext.Delay,
			&checkTime,
			false)
	}
	trxContext.Exec()
	trxContext.Finalize()
	//TODO
	defer func(b bool) {
		c.InTrxRequiringChecks = b
	}(c.InTrxRequiringChecks)

	if !trx.Implicit {
		var s types.TransactionStatus
		if trxContext.Delay == common.Microseconds(0) {
			s = types.TransactionStatusExecuted
		} else {
			s = types.TransactionStatusDelayed
		}
		tr := c.pushReceipt(trx.PackedTrx.PackedTrx, s, uint64(trxContext.BilledCpuTimeUs), trace.NetUsage)
		trace.Receipt = tr.TransactionReceiptHeader
		c.Pending.PendingBlockState.Trxs = append(c.Pending.PendingBlockState.Trxs, trx)
	} else {
		r := types.TransactionReceiptHeader{}
		r.CpuUsageUs = uint32(trxContext.BilledCpuTimeUs)
		r.NetUsageWords = uint32(trace.NetUsage / 8)
		trace.Receipt = r
	}
	//fc::move_append(pending->_actions, move(trx_context.executed))
	c.Pending.Actions = append(c.Pending.Actions, trxContext.Executed...)
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
	if !failureIsSubjective(trace.Except) {
		delete(c.UnAppliedTransactions, crypto.Sha256(trx.SignedID))
	}
	/*emit( c.accepted_transa
	ction, trx )
	emit( c.applied_transaction, trace )*/
	return trace
}

func (c *Controller) GetGlobalProperties() *entity.GlobalPropertyObject {
	gpo := entity.GlobalPropertyObject{}
	gpo.ID = 1
	err := c.DB.Find("id", gpo, &gpo)
	if err != nil {
		log.Error("GetGlobalProperties is error detail:%s", err.Error())
	}
	return &gpo
}

func (c *Controller) GetDynamicGlobalProperties() (r *entity.DynamicGlobalPropertyObject) {
	dgpo := entity.DynamicGlobalPropertyObject{}
	dgpo.ID = 1
	err := c.DB.Find("id", dgpo, &dgpo)
	if err != nil {
		log.Error("GetDynamicGlobalProperties is error detail:%s", err.Error())
	}

	return &dgpo
}

func (c *Controller) GetMutableResourceLimitsManager() *ResourceLimitsManager {
	return c.ResourceLimits
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
	trx.SetReferenceBlock(&c.Head.BlockId)
	in := c.Pending.PendingBlockState.Header.Timestamp + 999999
	trx.Expiration = common.TimePointSec(in)
	log.Info("getOnBlockTransaction trx.Expiration:%#v", trx)
	return trx
}
func (c *Controller) SkipDbSession(bs types.BlockStatus) bool {
	considerSkipping := (bs == types.BlockStatus(IRREVERSIBLE))
	//log.Info("considerSkipping:", considerSkipping)
	return considerSkipping
}

func (c *Controller) SkipDbSessions() bool {
	if !common.Empty(c.Pending) {
		return c.SkipDbSession(c.Pending.BlockStatus)
	} else {
		return false
	}
}

func (c *Controller) SkipTrxChecks() (b bool) {
	b = c.LightValidationAllowed(c.Config.disableReplayOpts)
	return
}

func (c *Controller) IsProducingBlock() bool {
	if common.Empty(c.Pending) {
		return false
	}
	return c.Pending.BlockStatus == types.Incomplete
}

/*func (c *Controller) IsResourceGreylisted(name *common.AccountName) bool {
	_,ok:=c.Config.resourceGreylist[*name]
	if ok {
		return true
	}
	return false
}*/

func (c *Controller) Close() {
	//session.close()
	c.ForkDB.Close()
	c.DB.Close()
	c.ReversibleBlocks.Close()

	//log.Info("Controller destory!")
	c.testClean()
	isActiveController = false

	c = nil
}

func (c *Controller) testClean() {
	err := os.RemoveAll("/tmp/data/")
	if err != nil {
		log.Error("Node data has been emptied is error:%s", err.Error())
	}
}

func (c *Controller) GetUnAppliedTransactions() *[]types.TransactionMetadata {
	result := []types.TransactionMetadata{}
	if c.ReadMode == SPECULATIVE {
		for _, v := range c.UnAppliedTransactions {
			result = append(result, v)
		}
	} else {
		log.Info("not empty unapplied_transactions in non-speculative mode")
	}
	return &result
}

func (c *Controller) DropUnAppliedTransaction(metadata *types.TransactionMetadata) {
	delete(c.UnAppliedTransactions, crypto.Sha256(metadata.SignedID))
}

func (c *Controller) DropAllUnAppliedTransactions() {
	c.UnAppliedTransactions = nil
}
func (c *Controller) GetScheduledTransactions() []common.TransactionIdType {

	result := []common.TransactionIdType{}
	gto := entity.GeneratedTransactionObject{}
	idx, err := c.DB.GetIndex("byDelay", &gto)
	itr := idx.Begin()
	for itr != idx.End() && gto.DelayUntil <= c.PendingBlockTime() {
		result = append(result, gto.TrxId)
		itr.Next()
		err = itr.Data(&gto)
	}
	if err != nil {
		log.Error("Controller GetScheduledTransactions is error:%s", err.Error())
	}
	itr.Release()
	return result
}
func (c *Controller) PushScheduledTransaction(trxId *common.TransactionIdType, deadLine common.TimePoint, billedCpuTimeUs uint32) *types.TransactionTrace {
	c.ValidateDbAvailableSize()
	return c.pushScheduledTransactionById(trxId, deadLine, billedCpuTimeUs, billedCpuTimeUs > 0)

}
func (c *Controller) pushScheduledTransactionById(sheduled *common.TransactionIdType,
	deadLine common.TimePoint,
	billedCpuTimeUs uint32, explicitBilledCpuTime bool) *types.TransactionTrace {

	in := entity.GeneratedTransactionObject{}
	in.TrxId = *sheduled
	out := entity.GeneratedTransactionObject{}
	c.DB.Find("byTrxId", in, &out)
	/*if err == nil {
		fmt.Println("unknown_transaction_exception", "unknown transaction")
	}*/
	EosAssert(&out != nil, &UnknownTransactionException{}, "unknown transaction")
	return c.pushScheduledTransactionByObject(&out, deadLine, billedCpuTimeUs, explicitBilledCpuTime)
}

func (c *Controller) pushScheduledTransactionByObject(gto *entity.GeneratedTransactionObject,
	deadLine common.TimePoint,
	billedCpuTimeUs uint32,
	explicitBilledCpuTime bool) *types.TransactionTrace {
	if !c.SkipDbSessions() {
		c.UndoSession = *c.DB.StartSession()
	}
	gtrx := entity.GeneratedTransactions(gto)
	c.RemoveScheduledTransaction(gto)
	EosAssert(gtrx.DelayUntil <= c.PendingBlockTime(), &TransactionException{}, "this transaction isn't ready,gtrx.DelayUntil:%s,PendingBlockTime:%s", gtrx.DelayUntil, c.PendingBlockTime())

	dtrx := types.SignedTransaction{}

	err := rlp.DecodeBytes(gtrx.PackedTrx, &dtrx)
	if err != nil {
		log.Error("PushScheduleTransaction1 DecodeBytes is error :%s", err.Error())
	}

	trx := &types.TransactionMetadata{}
	trx.Trx = &dtrx
	trx.Accepted = true
	trx.Scheduled = true

	trace := &types.TransactionTrace{}
	if gtrx.Expiration < c.PendingBlockTime() {
		trace.ID = gtrx.TrxId
		trace.BlockNum = c.PendingBlockState().BlockNum
		trace.BlockTime = types.BlockTimeStamp(c.PendingBlockTime())
		trace.ProducerBlockId = c.PendingProducerBlockId()
		trace.Scheduled = true
		trace.Receipt = (*c.pushReceipt(&gtrx.TrxId, types.TransactionStatusExecuted, uint64(billedCpuTimeUs), 0)).TransactionReceiptHeader
		//TODO
		/*emit( self.accepted_transaction, trx );
		emit( self.applied_transaction, trace );*/
		c.UndoSession.Squash()
		return trace
	}

	defer func(b bool) {
		c.InTrxRequiringChecks = b
	}(c.InTrxRequiringChecks)

	c.InTrxRequiringChecks = true
	cpuTimeToBillUs := billedCpuTimeUs
	trxContext := NewTransactionContext(c, &dtrx, gtrx.TrxId, common.Now())
	trxContext.Leeway = common.Milliseconds(0)
	trxContext.Deadline = deadLine
	trxContext.ExplicitBilledCpuTime = explicitBilledCpuTime
	trxContext.BilledCpuTimeUs = int64(billedCpuTimeUs)
	trace = trxContext.Trace
	Try(func() {
		trxContext.InitForDeferredTrx(gtrx.Published)
		trxContext.Exec()
		trxContext.Finalize()
		v := false
		defer func() {
			if v {
				log.Info("defer func exec")
			}
		}() //TODO
		trace.Receipt = (*c.pushReceipt(gtrx.TrxId, types.TransactionStatusExecuted, uint64(trxContext.BilledCpuTimeUs), trace.NetUsage)).TransactionReceiptHeader
		c.Pending.Actions = append(c.Pending.Actions, trxContext.Executed...)
		/*emit( self.accepted_transaction, trx );
		emit( self.applied_transaction, trace );*/
		trxContext.Squash()
		c.UndoSession.Squash()
		v = true
		//return trace
	}).Catch(func(ex Exception) {
		log.Error("PushScheduledTransaction is error:%s", ex.Message())
		cpuTimeToBillUs = trxContext.UpdateBilledCpuTime(common.Now())
		trace.Except = ex
		trace.ExceptPtr = &ex
		trace.Elapsed = (common.Now() - trxContext.Start).TimeSinceEpoch()
	}).End()

	trxContext.Undo()
	if !failureIsSubjective(trace.Except) && gtrx.Sender != 0 { /*gtrx.Sender != account_name()*/
		log.Info("%v", trace.Except.Message())
		errorTrace := applyOnerror(gtrx, deadLine, trxContext.pseudoStart, &cpuTimeToBillUs, billedCpuTimeUs, explicitBilledCpuTime)
		errorTrace.FailedDtrxTrace = trace
		trace = errorTrace
		if common.Empty(trace.ExceptPtr) {
			/*emit( self.accepted_transaction, trx );
			emit( self.applied_transaction, trace );*/
			c.UndoSession.Squash()
			return trace
		}
		trace.Elapsed = common.Now().TimeSinceEpoch() - trxContext.Start.TimeSinceEpoch()
	}

	subjective := false
	if explicitBilledCpuTime {
		subjective = failureIsSubjective(trace.Except)
	} else {
		subjective = scheduledFailureIsSubjective(trace.Except)
	}

	if !subjective {
		// hard failure logic

		if !explicitBilledCpuTime {
			rl := c.GetMutableResourceLimitsManager()
			rl.UpdateAccountUsage(&trxContext.BillToAccounts, uint32(types.BlockTimeStamp(c.PendingBlockTime())) /*.slot*/)
			//accountCpuLimit := 0
			accountNetLimit, accountCpuLimit, greylistedNet, greylistedCpu := trxContext.MaxBandwidthBilledAccountsCanPay(true)

			log.Info("test print: %v,%v,%v,%v", accountNetLimit, accountCpuLimit, greylistedNet, greylistedCpu) //TODO

			//cpuTimeToBillUs = cpuTimeToBillUs<accountCpuLimit:?trxContext.initialObjectiveDurationLimit.Count()
			tmp := uint32(0)
			if cpuTimeToBillUs < uint32(accountCpuLimit) {
				tmp = cpuTimeToBillUs
			} else {
				tmp = uint32(accountCpuLimit)
			}
			if tmp < uint32(trxContext.objectiveDurationLimit) {
				cpuTimeToBillUs = tmp
			}
		}

		c.ResourceLimits.AddTransactionUsage(&trxContext.BillToAccounts, uint64(cpuTimeToBillUs), 0,
			uint32(types.BlockTimeStamp(c.PendingBlockTime()))) // Should never fail

		receipt := *c.pushReceipt(gtrx.TrxId, types.TransactionStatusHardFail, uint64(cpuTimeToBillUs), 0)
		trace.Receipt = receipt.TransactionReceiptHeader
		/*emit( self.accepted_transaction, trx );
		emit( self.applied_transaction, trace );*/

		c.UndoSession.Squash()
	} else {
		/*emit( self.accepted_transaction, trx );
		emit( self.applied_transaction, trace );*/
	}
	trxContext.InitForDeferredTrx(gtrx.Published)
	//}
	return trace
}

func applyOnerror(gtrx *entity.GeneratedTransaction, deadline common.TimePoint, start common.TimePoint,
	cpuTimeToBillUs *uint32, billedCpuTimeUs uint32, explicitBilledCpuTime bool) *types.TransactionTrace {
	/*etrx :=types.SignedTransaction{}
	pl := types.PermissionLevel{gtrx.Sender, common.DefaultConfig.ActiveName}
	etrx.Actions = append(etrx.Actions, {})*/

	return nil
}
func (c *Controller) RemoveScheduledTransaction(gto *entity.GeneratedTransactionObject) {
	c.ResourceLimits.AddPendingRamUsage(gto.Payer, int64(9)+int64(len(gto.PackedTrx)))
	c.DB.Remove(gto)
}

func failureIsSubjective(e Exception) bool {
	code := e.Code()
	//fmt.Println(code == SubjectiveBlockProductionException{}.Code())
	return code == SubjectiveBlockProductionException{}.Code() ||
		code == BlockNetUsageExceeded{}.Code() ||
		code == GreylistNetUsageExceeded{}.Code() ||
		code == BlockCpuUsageExceeded{}.Code() ||
		code == GreylistCpuUsageExceeded{}.Code() ||
		code == DeadlineException{}.Code() ||
		code == LeewayDeadlineException{}.Code() ||
		code == ActorWhitelistException{}.Code() ||
		code == ActorBlacklistException{}.Code() ||
		code == ContractWhitelistException{}.Code() ||
		code == ContractBlacklistException{}.Code() ||
		code == ActionBlacklistException{}.Code() ||
		code == KeyBlacklistException{}.Code()

}

func scheduledFailureIsSubjective(e Exception) bool {
	code := e.Code()
	return (code == TxCpuUsageExceeded{}.Code()) || failureIsSubjective(e)
}
func (c *Controller) setActionMerkle() {
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

	EosAssert(!common.Empty(c.Pending), &BlockValidateException{}, "it is not valid to finalize when there is no pending block")

	c.ResourceLimits.ProcessAccountLimitUpdates()
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
	c.ResourceLimits.SetBlockParameters(cpu, net)

	c.setActionMerkle()

	c.setTrxMerkle()

	p := c.Pending.PendingBlockState
	p.BlockId = p.Header.BlockID()

	c.createBlockSummary(&p.BlockId)
}

func (c *Controller) SignBlock(signerCallback func(sha256 crypto.Sha256) ecc.Signature) {
	p := c.Pending.PendingBlockState
	p.Sign(signerCallback)
	//p.SignedBlock
	(*p.SignedBlock).SignedBlockHeader = p.Header
}

func (c *Controller) applyBlock(b *types.SignedBlock, s types.BlockStatus) {
	Try(func() {
		EosAssert(len(b.BlockExtensions) == 0, &BlockValidateException{}, "no supported extensions")
		producerBlockId := b.BlockID()
		c.startBlock(b.Timestamp, b.Confirmed, s, &producerBlockId)

		trace := &types.TransactionTrace{}
		for _, receipt := range b.Transactions {
			numPendingReceipts := len(c.Pending.PendingBlockState.SignedBlock.Transactions)
			if common.Empty(receipt.Trx.PackedTransaction) {
				pt := receipt.Trx.PackedTransaction
				mtrx := types.TransactionMetadata{}
				mtrx.PackedTrx = pt
				trace = c.pushTransaction(&mtrx, common.TimePoint(common.MaxMicroseconds()), receipt.CpuUsageUs, true)
			} else if common.Empty(receipt.Trx.TransactionID) {
				trace = c.PushScheduledTransaction(&receipt.Trx.TransactionID, common.TimePoint(common.MaxMicroseconds()), receipt.CpuUsageUs)
			} else {
				EosAssert(false, &BlockValidateException{}, "encountered unexpected receipt type")
			}
			transactionFailed := common.Empty(trace) && common.Empty(trace.Except)
			transactionCanFail := receipt.Status == types.TransactionStatusHardFail && common.Empty(receipt.Trx.TransactionID)
			if transactionFailed && !transactionCanFail {
				/*edump((*trace));
				throw *trace->except;*/
			}
			EosAssert(len(c.Pending.PendingBlockState.SignedBlock.Transactions) > 0,
				&BlockValidateException{}, "expected a block:%v,expected_receipt:%v", *b, receipt)

			EosAssert(len(c.Pending.PendingBlockState.SignedBlock.Transactions) == numPendingReceipts+1,
				&BlockValidateException{}, "expected receipt was not added:%v,expected_receipt:%v", *b, receipt)

			var trxReceipt types.TransactionReceipt
			length := len(c.Pending.PendingBlockState.SignedBlock.Transactions)
			if length > 0 {
				trxReceipt = c.Pending.PendingBlockState.SignedBlock.Transactions[length-1]
			}
			//r := trxReceipt.TransactionReceiptHeader
			EosAssert(trxReceipt == receipt, &BlockValidateException{}, "receipt does not match,producer_receipt:%#v", receipt, "validator_receipt:%#v", trxReceipt)
		}
		c.FinalizeBlock()

		EosAssert(producerBlockId == c.Pending.PendingBlockState.Header.BlockID(), &BlockValidateException{},
			"Block ID does not match,producer_block_id:%#v", producerBlockId, "validator_block_id:%#v", c.Pending.PendingBlockState.Header.BlockID())

		c.Pending.PendingBlockState.Header.ProducerSignature = b.ProducerSignature
		c.CommitBlock(false)
		return
	}).Catch(func(ex Exception) {
		log.Error("controller ApplyBlock is error:%s", ex.Message())
		c.AbortBlock()
	})
}

func (c *Controller) CommitBlock(addToForkDb bool) {
	defer func() {
		if PendingValid {
			c.Pending.Reset()
		}
	}()
	//try{
	if addToForkDb {
		c.Pending.PendingBlockState.Validated = true
		newBsp := c.ForkDB.AddBlockState(c.Pending.PendingBlockState)
		//emit(self.accepted_block_header, pending->_pending_block_state)
		c.Head = c.ForkDB.Header()
		EosAssert(newBsp == c.Head, &ForkDatabaseException{}, "committed block did not become the new head in fork database")
	}

	if !c.RePlaying {
		ubo := entity.ReversibleBlockObject{}
		ubo.BlockNum = c.Pending.PendingBlockState.BlockNum
		ubo.SetBlock(c.Pending.PendingBlockState.SignedBlock)
		c.DB.Insert(&ubo)
	}
	//emit( self.accepted_block, pending->_pending_block_state )
	//catch(){
	// reset_pending_on_exit.cancel();
	PendingValid = true //TODO
	//         abort_block();
	//throw;
	// }
	c.Pending.Push()
}

func (c *Controller) PushBlock(b *types.SignedBlock, s types.BlockStatus) {
	EosAssert(c.Pending != nil, &BlockValidateException{}, "it is not valid to push a block when there is a pending block")
	//resetProdLightValidation := c.makeBlockRestorePoint()

	EosAssert(b == nil, &BlockValidateException{}, "trying to push empty block")
	EosAssert(s != types.Incomplete, &BlockLogException{}, "invalid block status for a completed block")
	//emit(self.pre_accepted_block, b )
	//trust := !c.Config.forceAllChecks && (s== types.Irreversible || s== types.Validated)
	//newHeader := c.ForkDB.AddSignedBlockState(b,trust)

	if _, ok := c.Config.trustedProducers[b.Producer]; ok {
		//	resetProdLightValidation = true
	}
	//emit( self.accepted_block_header, new_header_state )
	if c.ReadMode != IRREVERSIBLE {
		//maybe_switch_forks( s )
	}

	if s == types.Irreversible {
		//emit( self.irreversible_block, new_header_state )
	}

} //status default value block_status s = block_status::complete

func (c *Controller) PushConfirmation(hc *types.HeaderConfirmation) {
	EosAssert(c.Pending != nil, &BlockValidateException{}, "it is not valid to push a confirmation when there is a pending block")
	c.ForkDB.Add(hc)
	//emit( c.accepted_confirmation, hc )
	if c.ReadMode != IRREVERSIBLE {
		c.maybeSwitchForks(types.Complete)
	}
}

func (c *Controller) maybeSwitchForks(s types.BlockStatus) {
	//TODO
	newHead := c.ForkDB.Head
	if newHead.Header.Previous == c.Head.BlockId {
		//try{

		c.applyBlock(newHead.SignedBlock, s)
		c.ForkDB.MarkInCurrentChain(newHead, true)
		c.ForkDB.SetValidity(newHead, true)
		c.Head = newHead

		//}catch(){
		c.ForkDB.SetValidity(newHead, false)
		//try.Throw()
		//}
	} else if newHead.ID != c.Head.ID {
		//branches := c.ForkDB.FetchBranchFrom( newHead.ID, c.Head.ID )
		/*for( auto itr = branches.second.begin(); itr != branches.second.end(); ++itr ) {
			fork_db.mark_in_current_chain( *itr , false );
			pop_block();
		}*/
		//exception.EosAssert( c.HeadBlockId() == branches.second.back()->header.previous, &exception.ForkDatabaseException{}, "loss of sync between fork_db and chainbase during fork switch" )
		/*for( auto ritr = branches.first.rbegin(); ritr != branches.first.rend(); ++ritr) {
			optional<fc::exception> except;
			try {
				apply_block( (*ritr)->block, (*ritr)->validated ? controller::block_status::validated : controller::block_status::complete );
				head = *ritr;
				fork_db.mark_in_current_chain( *ritr, true );
				(*ritr)->validated = true;
			}
			catch (const fc::exception& e) { except = e; }
			if (except) {
				elog("exception thrown while switching forks ${e}", ("e",except->to_detail_string()));

				// ritr currently points to the block that threw
				// if we mark it invalid it will automatically remove all forks built off it.
				fork_db.set_validity( *ritr, false );

				// pop all blocks from the bad fork
				// ritr base is a forward itr to the last block successfully applied
				auto applied_itr = ritr.base();
				for( auto itr = applied_itr; itr != branches.first.end(); ++itr ) {
					fork_db.mark_in_current_chain( *itr , false );
					pop_block();
				}
				EOS_ASSERT( self.head_block_id() == branches.second.back()->header.previous, fork_database_exception,
				"loss of sync between fork_db and chainbase during fork switch reversal" ); // _should_ never fail

				// re-apply good blocks
				for( auto ritr = branches.second.rbegin(); ritr != branches.second.rend(); ++ritr ) {
				apply_block( (*ritr)->block, controller::block_status::validated  );
				head = *ritr;
				fork_db.mark_in_current_chain( *ritr, true );
				}
				throw *except;
			}
		}*/
	}

}

func (c *Controller) DataBase() database.DataBase {
	return c.DB
}

func (c *Controller) ForkDataBase() *ForkDatabase {
	return c.ForkDB
}

func (c *Controller) GetAccount(name common.AccountName) *entity.AccountObject {
	accountObj := entity.AccountObject{}
	accountObj.Name = name
	err := c.DB.Find("byName", accountObj, &accountObj)
	if err != nil {
		log.Error("GetAccount is error :%s", err.Error())
	}
	return &accountObj
}

func (c *Controller) GetAuthorizationManager() *AuthorizationManager { return c.Authorization }

func (c *Controller) GetMutableAuthorizationManager() *AuthorizationManager { return c.Authorization }

//c++ flat_set<account_name> map[common.AccountName]interface{}
func (c *Controller) getActorWhiteList() *common.FlatSet {
	return &c.Config.ActorWhitelist
}

func (c *Controller) getActorBlackList() *common.FlatSet {
	return &c.Config.ActorBlacklist
}

func (c *Controller) getContractWhiteList() *common.FlatSet {
	return &c.Config.ContractWhitelist
}

func (c *Controller) getContractBlackList() *common.FlatSet {
	return &c.Config.ContractBlacklist
}

func (c *Controller) getActionBlockList() *common.FlatSet { return &c.Config.ActionBlacklist }

func (c *Controller) getKeyBlackList() *common.FlatSet { return &c.Config.KeyBlacklist }

func (c *Controller) SetActorWhiteList(params *common.FlatSet) {
	c.Config.ActorWhitelist = *params
}

func (c *Controller) SetActorBlackList(params *common.FlatSet) {
	c.Config.ActorBlacklist = *params
}

func (c *Controller) SetContractWhiteList(params *common.FlatSet) {
	c.Config.ContractWhitelist = *params
}

func (c *Controller) SetContractBlackList(params *common.FlatSet) {
	c.Config.ContractBlacklist = *params
}

func (c *Controller) SetActionBlackList(params *common.FlatSet) {
	c.Config.ActionBlacklist = *params
}

func (c *Controller) SetKeyBlackList(params *common.FlatSet) {
	c.Config.KeyBlacklist = *params
}

func (c *Controller) HeadBlockNum() uint32 { return c.Head.BlockNum }

func (c *Controller) HeadBlockTime() common.TimePoint { return c.Head.Header.Timestamp.ToTimePoint() }

func (c *Controller) HeadBlockId() common.BlockIdType { return common.BlockIdType{} }

func (c *Controller) HeadBlockProducer() common.AccountName { return c.Head.Header.Producer }

func (c *Controller) HeadBlockHeader() *types.BlockHeader { return &c.Head.Header.BlockHeader }

func (c *Controller) HeadBlockState() types.BlockState { return types.BlockState{} }

func (c *Controller) ForkDbHeadBlockNum() uint32 { return c.ForkDB.Header().BlockNum }

func (c *Controller) ForkDbHeadBlockId() common.BlockIdType { return common.BlockIdType{} }

func (c *Controller) ForkDbHeadBlockTime() common.TimePoint {
	return c.ForkDB.Header().Header.Timestamp.ToTimePoint()
}

func (c *Controller) ForkDbHeadBlockProducer() common.AccountName {
	return c.ForkDB.Header().Header.Producer
}

func (c *Controller) PendingBlockState() *types.BlockState {
	if c.Pending != nil {
		return c.Pending.PendingBlockState
	}
	return &types.BlockState{}
}

func (c *Controller) PendingBlockTime() common.TimePoint {
	EosAssert(!common.Empty(c.Pending), &BlockValidateException{}, "no pending block")
	return c.Pending.PendingBlockState.Header.Timestamp.ToTimePoint()
}

func (c *Controller) PendingProducerBlockId() common.BlockIdType {
	EosAssert(c.Pending != nil, &BlockValidateException{}, "no pending block")
	return c.Pending.ProducerBlockId
}

func (c *Controller) ActiveProducers() *types.ProducerScheduleType {
	if c.Pending != nil {
		return &c.Head.ActiveSchedule
	}
	return &c.Pending.PendingBlockState.ActiveSchedule
}

func (c *Controller) PendingProducers() *types.ProducerScheduleType {
	if c.Pending != nil {
		return &c.Head.PendingSchedule
	}
	return &c.Pending.PendingBlockState.ActiveSchedule
}

func (c *Controller) ProposedProducers() types.ProducerScheduleType {
	gpo := c.GetGlobalProperties()
	if common.Empty(gpo.ProposedScheduleBlockNum) {
		return types.ProducerScheduleType{}
	}
	return *gpo.ProposedSchedule.ProducerScheduleType()
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

func (c *Controller) LastIrreversibleBlockNum() uint32 { return 0 }

func (c *Controller) LastIrreversibleBlockId() common.BlockIdType { return common.BlockIdType{} }

func (c *Controller) FetchBlockByNumber(blockNum uint32) *types.SignedBlock {
	blkState := c.ForkDB.GetBlockInCurrentChainByNum(blockNum)
	if blkState != nil {
		return blkState.SignedBlock
	}
	return c.Blog.ReadBlockByNum(blockNum)
}

func (c *Controller) FetchBlockById(id common.BlockIdType) *types.SignedBlock {
	state := c.ForkDB.GetBlock(&id)
	if state != nil {
		return state.SignedBlock
	}
	bptr := c.FetchBlockByNumber(types.NumFromID(&id))
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

func (c *Controller) GetBlockIdForNum(blockNum uint32) common.BlockIdType {
	blkState := c.ForkDB.GetBlockInCurrentChainByNum(blockNum)
	if blkState != nil {
		return blkState.BlockId
	}

	signedBlk := c.Blog.ReadBlockByNum(blockNum)
	EosAssert(common.Empty(signedBlk), &UnknownBlockException{}, "Could not find block: %d", blockNum)
	return signedBlk.BlockID()
}

func (c *Controller) CheckContractList(code common.AccountName) {
	if len(c.Config.ContractWhitelist.Data) > 0 {
		exist, _ := c.Config.ContractWhitelist.Find(&code)
		EosAssert(!exist, &ContractWhitelistException{}, "account %d is not on the contract whitelist", code)
	} else if len(c.Config.ContractBlacklist.Data) > 0 {
		//EosAssert(exist, &ContractBlacklistException{}, "account %d is on the contract blacklist", code)
		/*EOS_ASSERT( conf.contract_blacklist.find( code ) == conf.contract_blacklist.end(),
			contract_blacklist_exception,
			"account '${code}' is on the contract blacklist", ("code", code)
		);*/
	}
}

func (c *Controller) CheckActionList(code common.AccountName, action common.ActionName) {
	if len(c.Config.ActionBlacklist.Data) > 0 {
		//abl := common.MakePair(code, action)
		//key := Hash(abl)
		/*if _, ok := c.Config.ActionBlacklist[abl]; ok {
			fmt.Println("action '${code}::${action}' is on the action blacklist")
			return
		}*/
		/*EOS_ASSERT( conf.action_blacklist.find( std::make_pair(code, action) ) == conf.action_blacklist.end(),
			action_blacklist_exception,
			"action '${code}::${action}' is on the action blacklist",
			("code", code)("action", action)
		);*/
	}
}

func (c *Controller) CheckKeyList(key *ecc.PublicKey) {
	if len(c.Config.KeyBlacklist.Data) > 0 {
		exist, _ := c.Config.KeyBlacklist.Find(key)
		EosAssert(exist, &KeyBlacklistException{}, "public key %v is on the key blacklist", key)
	}
}

func (c *Controller) IsProducing() bool {
	if !common.Empty(c.Pending) {
		return false
	}
	return c.Pending.BlockStatus == types.Incomplete
}

func (c *Controller) IsRamBillingInNotifyAllowed() bool {
	return !c.IsProducingBlock() || c.Config.allowRamBillingInNotify
}

func (c *Controller) AddResourceGreyList(name *common.AccountName) {
	c.Config.resourceGreylist[*name] = struct{}{}
}

func (c *Controller) RemoveResourceGreyList(name *common.AccountName) {
	delete(c.Config.resourceGreylist, *name)
}

func (c *Controller) IsResourceGreylisted(name *common.AccountName) bool {
	_, ok := c.Config.resourceGreylist[*name]
	if ok {
		return true
	}
	return false
}
func (c *Controller) GetResourceGreyList() map[common.AccountName]struct{} {
	return c.Config.resourceGreylist
}

//TODO
func (c *Controller) ValidateReferencedAccounts(t *types.Transaction) {
	/*for _,a := range t.ContextFreeActions{
		c.DB.f
	}*/
}

func (c *Controller) ValidateExpiration(t *types.Transaction) {
	chainConfiguration := c.GetGlobalProperties().Configuration
	EosAssert(common.TimePoint(t.Expiration) >= c.PendingBlockTime(),
		&ExpiredTxException{}, "transaction has expired, expiration is %v and pending block time is %v",
		t.Expiration, c.PendingBlockTime())
	EosAssert(common.TimePoint(t.Expiration) <= c.PendingBlockTime()+common.TimePoint(common.Seconds(int64(chainConfiguration.MaxTrxLifetime))),
		&TxExpTooFarException{}, "Transaction expiration is too far in the future relative to the reference time of %v, expiration is %v and the maximum transaction lifetime is %v seconds",
		t.Expiration, c.PendingBlockTime(), chainConfiguration.MaxTrxLifetime)
}

func (c *Controller) ValidateTapos(t *types.Transaction) {
	in := entity.BlockSummaryObject{}
	in.Id = common.IdType(t.RefBlockNum)
	taposBlockSummary := entity.BlockSummaryObject{}
	err := c.DB.Find("", in, &taposBlockSummary)
	if err != nil {
		log.Error("ValidateTapos Is Error:%s", err.Error())
	}
	EosAssert(t.VerifyReferenceBlock(&taposBlockSummary.BlockId), &InvalidRefBlockException{},
		"Transaction's reference block did not match. Is this transaction from a different fork? taposBlockSummary:%v", taposBlockSummary)
}

/* c++ 1.4.1
void controller::validate_tapos( const transaction& trx )const { try {
const auto& tapos_block_summary = db().get<block_summary_object>((uint16_t)trx.ref_block_num);

//Verify TaPoS block summary has correct ID prefix, and that this block's time is not past the expiration
EOS_ASSERT(trx.verify_reference_block(tapos_block_summary.block_id), invalid_ref_block_exception,
"Transaction's reference block did not match. Is this transaction from a different fork?",
("tapos_summary", tapos_block_summary));
} FC_CAPTURE_AND_RETHROW() }
*/
func (c *Controller) ValidateDbAvailableSize() {
	/*const auto free = db().get_segment_manager()->get_free_memory();
	const auto guard = my->conf.state_guard_size;
	EOS_ASSERT(free >= guard, database_guard_exception, "database free: ${f}, guard size: ${g}", ("f", free)("g",guard));*/
}

func (c *Controller) ValidateReversibleAvailableSize() {
	/*const auto free = my->reversible_blocks.get_segment_manager()->get_free_memory();
	const auto guard = my->conf.reversible_guard_size;
	EOS_ASSERT(free >= guard, reversible_guard_exception, "reversible free: ${f}, guard size: ${g}", ("f", free)("g",guard));*/
}

func (c *Controller) IsKnownUnexpiredTransaction(id *common.TransactionIdType) bool {
	result := entity.TransactionObject{}
	in := entity.TransactionObject{}
	in.TrxID = *id
	err := c.DB.Find("byTrxId", in, &result)
	if err != nil {
		log.Error("IsKnownUnexpiredTransaction Is Error:%s", err.Error())
	}
	return common.Empty(result)
}

func (c *Controller) SetProposedProducers(producers []types.ProducerKey) int64 {

	gpo := c.GetGlobalProperties()
	curBlockNum := c.HeadBlockNum() + 1
	if common.Empty(gpo.ProposedScheduleBlockNum) {
		if gpo.ProposedScheduleBlockNum != curBlockNum {
			return -1
		}

		if compare(producers, gpo.ProposedSchedule.Producers) {
			return -1
		}
	}
	sch := types.ProducerScheduleType{}
	/*begin :=types.ProducerKey{}
	end :=types.ProducerKey{}*/
	if len(c.Pending.PendingBlockState.PendingSchedule.Producers) == 0 {
		activeSch := c.Pending.PendingBlockState.ActiveSchedule
		if compare(producers, activeSch.Producers) {
			return -1
		}
		sch.Version = activeSch.Version + 1
	} else {
		pendingSch := c.Pending.PendingBlockState.PendingSchedule
		if compare(producers, pendingSch.Producers) {
			return -1
		}
		sch.Version = pendingSch.Version + 1
	}

	sch.Producers = producers
	version := sch.Version
	c.DB.Modify(&gpo, func(p *entity.GlobalPropertyObject) {
		p.ProposedScheduleBlockNum = curBlockNum
		tmp := p.ProposedSchedule.SharedProducerScheduleType(sch)
		p.ProposedSchedule = *tmp
	})
	return int64(version)
}

//for SetProposedProducers
func compare(first []types.ProducerKey, second []types.ProducerKey) bool {
	if len(first) != len(second) {
		return false
	}
	for i := 0; i < len(first); i++ {
		if first[i] != second[i] {
			return false
		}
	}
	return true
}

func (c *Controller) SkipAuthCheck() bool { return c.LightValidationAllowed(c.Config.forceAllChecks) }

func (c *Controller) ContractsConsole() bool { return c.Config.contractsConsole }

func (c *Controller) GetChainId() common.ChainIdType { return c.ChainID }

func (c *Controller) GetReadMode() DBReadMode { return c.ReadMode }

func (c *Controller) GetValidationMode() ValidationMode { return c.Config.blockValidationMode }

func (c *Controller) SetSubjectiveCpuLeeway(leeway common.Microseconds) {
	c.SubjectiveCupLeeway = leeway
}

func (c *Controller) GetWasmInterface() *wasmgo.WasmGo {
	return c.WasmIf
}

/*func (c *Controller) GetAbiSerializer(name common.AccountName,
	maxSerializationTime common.Microseconds) types.AbiSerializer {
	return types.AbiSerializer{}
}*/

/*func (c *Controller) ToVariantWithAbi(obj interface{}, maxSerializationTime common.Microseconds) {}*/

func (c *Controller) CreateNativeAccount(name common.AccountName, owner types.Authority, active types.Authority, isPrivileged bool) {
	account := entity.AccountObject{}
	account.Name = name
	account.CreationDate = types.BlockTimeStamp(c.Config.genesis.InitialTimestamp)
	account.Privileged = isPrivileged
	if name == common.AccountName(common.DefaultConfig.SystemAccountName) {
		abiDef := types.AbiDef{}
		account.SetAbi(EosioContractAbi(abiDef))
	}
	err := c.DB.Insert(&account)
	if err != nil {
		log.Error("CreateNativeAccount Insert Is Error:%s", err.Error())
	}

	aso := entity.AccountSequenceObject{}
	aso.Name = name
	c.DB.Insert(&aso)

	ownerPermission := c.Authorization.CreatePermission(name, common.PermissionName(common.DefaultConfig.OwnerName), 0, owner, c.Config.genesis.InitialTimestamp)

	activePermission := c.Authorization.CreatePermission(name, common.PermissionName(common.DefaultConfig.ActiveName), PermissionIdType(ownerPermission.ID), active, c.Config.genesis.InitialTimestamp)

	c.ResourceLimits.InitializeAccount(name)
	ramDelta := uint64(common.DefaultConfig.OverheadPerAccountRamBytes)
	ramDelta += 2 * common.BillableSizeV("permission_object") //::billable_size_v<permission_object>
	ramDelta += ownerPermission.Auth.GetBillableSize()
	ramDelta += activePermission.Auth.GetBillableSize()
	c.ResourceLimits.AddPendingRamUsage(name, int64(ramDelta))
	c.ResourceLimits.VerifyAccountRamUsage(name)
}

func (c *Controller) initializeForkDB() {

	gs := types.GetGenesisStateInstance()
	pst := types.ProducerScheduleType{0, []types.ProducerKey{
		{common.DefaultConfig.SystemAccountName, gs.InitialKey}}}
	genHeader := types.BlockHeaderState{}
	genHeader.ActiveSchedule = pst
	genHeader.PendingSchedule = pst
	genHeader.PendingScheduleHash = crypto.Hash256(pst)
	genHeader.Header.Timestamp = types.NewBlockTimeStamp(gs.InitialTimestamp)
	genHeader.Header.ActionMRoot = common.CheckSum256Type(gs.ComputeChainID())
	genHeader.BlockId = genHeader.Header.BlockID()
	genHeader.BlockNum = genHeader.Header.BlockNumber()
	c.Head = types.NewBlockState(&genHeader)
	signedBlock := types.SignedBlock{}
	signedBlock.SignedBlockHeader = genHeader.Header
	c.Head.SignedBlock = &signedBlock
	//log.Info("Controller initializeForkDB:%v", c.ForkDB.DB)

	c.ForkDB.SetHead(c.Head)
	c.DB.SetRevision(int64(c.Head.BlockNum))
	c.initializeDatabase()
}

func (c *Controller) initializeDatabase() {

	/*for (int i = 0; i < 0x10000; i++)
	db.create<block_summary_object>([&](block_summary_object&) {});

	const auto& tapos_block_summary = db.get<block_summary_object>(1);
	db.modify( tapos_block_summary, [&]( auto& bs ) {
		bs.block_id = head->id;
	})*/
	gi := c.Config.genesis.Initial()
	//gi.Validate()	//check config
	gpo := entity.GlobalPropertyObject{}
	gpo.Configuration = gi
	err := c.DB.Insert(&gpo)
	if err != nil {
		log.Error("Controller initializeDatabase insert GlobalPropertyObject is error:%s", err)
	}
	dgpo := entity.DynamicGlobalPropertyObject{}
	dgpo.ID = 1
	//dgpo.GlobalActionSequence = 10000
	err = c.DB.Insert(&dgpo)
	if err != nil {
		log.Error("Controller initializeDatabase insert DynamicGlobalPropertyObject is error:%s", err)
	}

	c.ResourceLimits.InitializeDatabase()
	systemAuth := types.Authority{}
	kw := types.KeyWeight{}
	kw.Key = c.Config.genesis.InitialKey
	systemAuth.Keys = []types.KeyWeight{kw}
	c.CreateNativeAccount(common.DefaultConfig.SystemAccountName, systemAuth, systemAuth, true)
	emptyAuthority := types.Authority{}
	emptyAuthority.Threshold = 1
	activeProducersAuthority := types.Authority{}
	activeProducersAuthority.Threshold = 1
	//plw:=types.PermissionLevelWeight{}
	p := types.PermissionLevelWeight{types.PermissionLevel{common.DefaultConfig.SystemAccountName, common.DefaultConfig.ActiveName}, 1}
	activeProducersAuthority.Accounts = append(activeProducersAuthority.Accounts, p)
	c.CreateNativeAccount(common.DefaultConfig.NullAccountName, emptyAuthority, emptyAuthority, false)
	c.CreateNativeAccount(common.DefaultConfig.ProducersAccountName, emptyAuthority, activeProducersAuthority, false)
	activePermission := c.Authorization.GetPermission(&types.PermissionLevel{common.DefaultConfig.ProducersAccountName, common.DefaultConfig.ActiveName})

	majorityPermission := c.Authorization.CreatePermission(common.DefaultConfig.ProducersAccountName,
		common.DefaultConfig.MajorityProducersPermissionName,
		PermissionIdType(activePermission.ID),
		activeProducersAuthority,
		c.Config.genesis.InitialTimestamp)

	minorityPermission := c.Authorization.CreatePermission(common.DefaultConfig.ProducersAccountName,
		common.DefaultConfig.MinorityProducersPermissionName,
		PermissionIdType(majorityPermission.ID),
		activeProducersAuthority,
		c.Config.genesis.InitialTimestamp)

	log.Info("initializeDatabase print:%#v,%#v", majorityPermission, minorityPermission)
}

func (c *Controller) initialize() {
	if common.Empty(c.Head) {
		c.initializeForkDB()
		end := c.Blog.ReadHead()
		if common.Empty(end) && end.BlockNumber() > 1 {
			endTime := end.Timestamp.ToTimePoint()
			replaying := true
			replayHeadTime := endTime
			//ilog( "existing block log, attempting to replay ${n} blocks", ("n",end->block_num()) )
			start := common.Now()
			next := c.Blog.ReadBlockByNum(c.Head.BlockNum + 1)
			/*while( auto next = blog.read_block_by_num( head->block_num + 1 ) ) {
				self.push_block( next, controller::block_status::irreversible );
				if( next->block_num() % 100 == 0 ) {
					std::cerr << std::setw(10) << next->block_num() << " of " << end->block_num() <<"\r";
				}
			}*/
			//ilog( "${n} blocks replayed", ("n", head->block_num) )
			//c.DB.set_revision(head->block_num)
			rev := 0
			//c.ReversibleBlocks.Find("",)
			r := entity.ReversibleBlockObject{}
			for {
				r.BlockNum = c.HeadBlockNum() + 1
				err := c.ReversibleBlocks.Find("blockNum", r, r)
				if err != nil {
					break
				}
				c.PushBlock(r.GetBlock(), types.Validated)
			}
			//ilog( "${n} reversible blocks replayed", ("n",rev) )
			end := time.Now()
			/*ilog( "replayed ${n} blocks in ${duration} seconds, ${mspb} ms/block",
				("n", head->block_num)("duration", (end-start).count()/1000000)
			("mspb", ((end-start).count()/1000.0)/head->block_num)        )*/
			c.RePlaying = false
			//c.ReplayHeadTime = nil

			log.Info("test print:", replaying, replayHeadTime, start, next, rev, end)
		} else if !common.Empty(end) {
			c.Blog.ResetToGenesis(&c.Config.genesis, c.Head.SignedBlock)
		}
		//TODO	wait append
		/*rbi := entity.ReversibleBlockObject{}
		ubi,err := c.ReversibleBlocks.GetIndex("byNum",&rbi)
		if err!= nil{
			fmt.Errorf("initialize database is error :",err)
		}
		objitr := ubi.Begin()*/
	}

}

//c++ pair<scope_name,action_name>
type HandlerKey struct {
	scope  common.ScopeName
	action common.ActionName
}

func NewHandlerKey(scopeName common.ScopeName, actionName common.ActionName) HandlerKey {
	hk := HandlerKey{scopeName, actionName}
	return hk
}

func (c *Controller) clearExpiredInputTransactions() {
	/*aa :=&entity.TransactionObject{}
	aa.Expiration = common.TimePointSec(common.Now())
	err := c.DB.Insert(aa)
	if err != nil{
		fmt.Println("insert success")
	}*/
	transactionIdx, err := c.DB.GetIndex("byExpiration", &entity.TransactionObject{})

	now := c.PendingBlockTime()
	t := &entity.TransactionObject{}

	for !transactionIdx.Empty() && now > common.TimePoint(t.Expiration) {
		tmp := &entity.TransactionObject{}
		itr := transactionIdx.Begin()
		if itr != nil {
			err = itr.Data(tmp)
			if err != nil {
				log.Error("TransactionIdx.Begin Is Error:%s", err.Error())
			}
		}
		c.DB.Remove(tmp)
	}
}

func (c *Controller) CheckActorList(actors *common.FlatSet) {
	if c.Config.ActorWhitelist.Len() > 0 {
		//excluded :=make(map[common.AccountName]struct{})

		//set
		/*for an := range actors.Data {
			if c.Config.ActorWhitelist.Find(an.(*common.AccountName)){

			}
		}*/
		/*EOS_ASSERT( excluded.size() == 0, actor_whitelist_exception,
			"authorizing actor(s) in transaction are not on the actor whitelist: ${actors}",
			("actors", excluded)
		)*/
		/*} else if len(c.Config.ActorBlacklist) > 0 {
		//black :=make(map[common.AccountName]struct{})
		//set
		for _, an := range actors.Data {
			if _, ok := c.Config.ActorBlacklist[*an.(*common.AccountName)]; ok {
				fmt.Println("authorizing actor(s) in transaction are not on the actor blacklist:", an)
				return
			}
		}*/
		/*EOS_ASSERT( blacklisted.size() == 0, actor_blacklist_exception,
			"authorizing actor(s) in transaction are on the actor blacklist: ${actors}",
			("actors", blacklisted)
		)*/
	}
}
func (c *Controller) updateProducersAuthority() {
	producers := c.Pending.PendingBlockState.ActiveSchedule.Producers
	updatePermission := func(permission *entity.PermissionObject, threshold uint32) {
		auth := types.Authority{threshold, []types.KeyWeight{}, []types.PermissionLevelWeight{}, []types.WaitWeight{}}
		for _, p := range producers {
			auth.Accounts = append(auth.Accounts, types.PermissionLevelWeight{types.PermissionLevel{p.ProducerName, common.DefaultConfig.ActiveName}, 1})
		}
		if !permission.Auth.Equals(auth.ToSharedAuthority()) {
			c.DB.Modify(permission, func(param *types.Permission) {
				param.RequiredAuth = auth
			})
		}
	}

	numProducers := len(producers)
	calculateThreshold := func(numerator uint32, denominator uint32) uint32 {
		return ((uint32(numProducers) * numerator) / denominator) + 1
	}
	updatePermission(c.Authorization.GetPermission(&types.PermissionLevel{common.DefaultConfig.ProducersAccountName, common.DefaultConfig.ActiveName}), calculateThreshold(2, 3))

	updatePermission(c.Authorization.GetPermission(&types.PermissionLevel{common.DefaultConfig.ProducersAccountName, common.DefaultConfig.MajorityProducersPermissionName}), calculateThreshold(1, 2))

	updatePermission(c.Authorization.GetPermission(&types.PermissionLevel{common.DefaultConfig.ProducersAccountName, common.DefaultConfig.MinorityProducersPermissionName}), calculateThreshold(1, 3))

}

func (c *Controller) createBlockSummary(id *common.BlockIdType) {
	blockNum := types.NumFromID(id)
	sid := blockNum & 0xffff
	bso := entity.BlockSummaryObject{}
	bso.Id = common.IdType(sid)
	err := c.DB.Find("id", bso, bso)
	if err != nil {
		log.Error("Controller createBlockSummary is error:%s", err.Error())
	}
	c.DB.Modify(bso, func(b *entity.BlockSummaryObject) {
		b.BlockId = *id
	})
}

func (c *Controller) initConfig() *Controller {
	c.Config = Config{
		blocksDir:               common.DefaultConfig.DefaultBlocksDirName,
		stateDir:                common.DefaultConfig.DefaultStateDirName,
		stateSize:               common.DefaultConfig.DefaultStateSize,
		stateGuardSize:          common.DefaultConfig.DefaultStateGuardSize,
		reversibleCacheSize:     common.DefaultConfig.DefaultReversibleCacheSize,
		reversibleGuardSize:     common.DefaultConfig.DefaultReversibleGuardSize,
		readOnly:                false,
		forceAllChecks:          false,
		disableReplayOpts:       false,
		contractsConsole:        false,
		allowRamBillingInNotify: false,
		//vmType:              common.DefaultConfig.DefaultWasmRuntime, //TODO
		readMode:            SPECULATIVE,
		blockValidationMode: FULL,
	}
	return c
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

/*    about chain

signal<void(const signed_block_ptr&)>         pre_accepted_block;
signal<void(const block_state_ptr&)>          accepted_block_header;
signal<void(const block_state_ptr&)>          accepted_block;
signal<void(const block_state_ptr&)>          irreversible_block;
signal<void(const transaction_metadata_ptr&)> accepted_transaction;
signal<void(const transaction_trace_ptr&)>    applied_transaction;
signal<void(const header_confirmation&)>      accepted_confirmation;
signal<void(const int&)>                      bad_alloc;*/
