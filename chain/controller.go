package chain

import (
	"fmt"
	"github.com/eosspark/container/sets/treeset"
	abi "github.com/eosspark/eos-go/chain/abi_serializer"
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
	"github.com/eosspark/eos-go/plugins/appbase/app/include"
	"github.com/eosspark/eos-go/wasmgo"
	"os"
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

func (d DBReadMode) String() string {
	switch d {
	case SPECULATIVE:
		return "speculative"
	case HEADER:
		return "header"
	case READONLY:
		return "readonly"
	case IRREVERSIBLE:
		return "irreversible"
	default:
		return ""
	}
}

func DBReadModeFromString(s string) (DBReadMode, bool) {
	switch s {
	case "SPECULATIVE", "speculative":
		return SPECULATIVE, true
	case "HEADER", "header":
		return HEADER, true
	case "READONLY", "readonly":
		return READONLY, true
	case "IRREVERSIBLE", "irreversible":
		return IRREVERSIBLE, true
	default:
		return -1, false
	}
}

type ValidationMode int8

const (
	FULL = ValidationMode(iota)
	LIGHT
)

func (v ValidationMode) String() string {
	switch v {
	case FULL:
		return "full"
	case LIGHT:
		return "light"
	default:
		return ""
	}
}

func ValidationModeFromString(s string) (ValidationMode, bool) {
	switch s {
	case "FULL", "full":
		return FULL, true
	case "LIGHT", "light":
		return LIGHT, true
	default:
		return -1, false
	}
}

type Config struct {
	ActorWhitelist          treeset.Set //common.AccountName
	ActorBlacklist          treeset.Set //common.AccountName
	ContractWhitelist       treeset.Set //common.AccountName
	ContractBlacklist       treeset.Set //common.AccountName]struct{}
	ActionBlacklist         treeset.Set //common.Pair //see actionBlacklist
	KeyBlacklist            treeset.Set
	ResourceGreylist        treeset.Set
	TrustedProducers        treeset.Set
	BlocksDir               string
	StateDir                string
	ReversibleDir           string
	StateSize               uint64
	StateGuardSize          uint64
	ReversibleCacheSize     uint64
	ReversibleGuardSize     uint64
	ReadOnly                bool
	ForceAllChecks          bool
	DisableReplayOpts       bool
	DisableReplay           bool
	ContractsConsole        bool
	AllowRamBillingInNotify bool
	Genesis                 *types.GenesisState
	VmType                  wasmgo.WasmGo
	ReadMode                DBReadMode
	BlockValidationMode     ValidationMode
}

func NewConfig() *Config {
	return &Config{
		BlocksDir:               common.DefaultConfig.DefaultBlocksDirName,
		StateDir:                common.DefaultConfig.DefaultStateDirName,
		StateSize:               common.DefaultConfig.DefaultStateSize,
		StateGuardSize:          common.DefaultConfig.DefaultStateGuardSize,
		ReversibleCacheSize:     common.DefaultConfig.DefaultReversibleCacheSize,
		ReversibleGuardSize:     common.DefaultConfig.DefaultReversibleGuardSize,
		ReadOnly:                false,
		ForceAllChecks:          false,
		DisableReplayOpts:       false,
		ContractsConsole:        false,
		AllowRamBillingInNotify: false,
		ReadMode:                SPECULATIVE,
		BlockValidationMode:     FULL,
		Genesis:                 types.NewGenesisState(),

		ActorWhitelist:    *treeset.NewWith(common.TypeName, common.CompareName),
		ActorBlacklist:    *treeset.NewWith(common.TypeName, common.CompareName),
		ContractWhitelist: *treeset.NewWith(common.TypeName, common.CompareName),
		ContractBlacklist: *treeset.NewWith(common.TypeName, common.CompareName),
		ActionBlacklist:   *treeset.NewWith(common.TypePair, common.ComparePair),
		KeyBlacklist:      *treeset.NewWith(ecc.TypePubKey, ecc.ComparePubKey),
		ResourceGreylist:  *treeset.NewWith(common.TypeName, common.CompareName),
		TrustedProducers:  *treeset.NewWith(common.TypeName, common.CompareName),
	}
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
	UnappliedTransactions          map[crypto.Sha256]types.TransactionMetadata
	GpoCache                       map[common.IdType]*entity.GlobalPropertyObject
	PreAcceptedBlock               include.Signal
	AcceptedBlockHeader            include.Signal
	AcceptedBlock                  include.Signal
	IrreversibleBlock              include.Signal
	AcceptedTransaction            include.Signal
	AppliedTransaction             include.Signal
	AcceptedConfirmation           include.Signal
	BadAlloc                       include.Signal
}

func GetControllerInstance() *Controller {
	if !isActiveController {
		validPath()
		instance = newController()
	}
	return instance
}

func validPath() {
	path := []string{common.DefaultConfig.DefaultStateDirName, common.DefaultConfig.DefaultBlocksDirName, common.DefaultConfig.DefaultReversibleBlocksDirName}
	for _, d := range path {
		_, err := os.Stat(d)
		if os.IsNotExist(err) {
			err := os.MkdirAll(d, os.ModePerm)
			if err != nil {
				log.Info("controller validPath mkdir failed:%s\n", err)
			} else {
				log.Info("controller validPath mkdir success:%s\n", d)
			}
		}
	}
}
func NewController(cfg *Config) *Controller {
	isActiveController = true //controller is active
	db, err := database.NewDataBase(cfg.StateDir)
	reversibleDB, err := database.NewDataBase(cfg.ReversibleDir)

	if err != nil {
		log.Error("newController create database is error :%s", err)
		return nil
	}
	con := &Controller{InTrxRequiringChecks: false, RePlaying: false, TrustedProducerLightValidation: false}
	con.DB = db
	con.ReversibleBlocks = reversibleDB

	con.Blog = NewBlockLog(common.DefaultConfig.DefaultBlocksDirName)

	con.ForkDB, err = newForkDatabase(cfg.BlocksDir, common.DefaultConfig.ForkDbName, true)

	con.ChainID = cfg.Genesis.ComputeChainID()

	con.ReadMode = cfg.ReadMode
	con.ApplyHandlers = make(map[string]v)
	con.WasmIf = wasmgo.NewWasmGo()
	con.ResourceLimits = newResourceLimitsManager(con)
	con.Authorization = newAuthorizationManager(con)
	con.UnappliedTransactions = make(map[crypto.Sha256]types.TransactionMetadata)
	con.Config = *cfg

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
		common.ActionName(common.N("linkauth")), applyEosioLinkauth)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("unlinkauth")), applyEosioUnlinkauth)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("canceldelay")), applyEosioCanceldalay)
	con.initialize()
	/*fork_db.irreversible.connect( [&]( auto b ) {
		on_irreversible(b);
	})*/

	return con
}

func newController() *Controller {
	isActiveController = true //controller is active
	//init db
	db, err := database.NewDataBase(common.DefaultConfig.DefaultStateDirName)
	if err != nil {
		log.Error("newController is error :%s", err)
		return nil
	}
	//init ReversibleBlocks
	//reversibleDir := common.DefaultConfig.DefaultBlocksDirName + "/" + common.DefaultConfig.DefaultReversibleBlocksDirName
	reversibleDB, err := database.NewDataBase(common.DefaultConfig.DefaultReversibleBlocksDirName)
	if err != nil {
		log.Error("newController init reversibleDB is error:%s", err)
	}
	con := &Controller{InTrxRequiringChecks: false, RePlaying: false, TrustedProducerLightValidation: false}
	con.DB = db
	con.ReversibleBlocks = reversibleDB

	con.Blog = NewBlockLog(common.DefaultConfig.DefaultBlocksDirName)

	con.ForkDB, _ = newForkDatabase(common.DefaultConfig.DefaultBlocksDirName, common.DefaultConfig.ForkDbName, true)
	con.initConfig()
	con.ChainID = con.Config.Genesis.ComputeChainID()

	con.ReadMode = con.Config.ReadMode
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
		common.ActionName(common.N("linkauth")), applyEosioLinkauth)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("unlinkauth")), applyEosioUnlinkauth)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.AccountName(common.N("eosio")),
		common.ActionName(common.N("canceldelay")), applyEosioCanceldalay)

	//IrreversibleBlock.connect()
	//readycontroller = make(chan bool)
	//go initResource(con, readycontroller)
	//con.Pending = &PendingState{}
	con.ResourceLimits = newResourceLimitsManager(con)
	con.Authorization = newAuthorizationManager(con)
	con.UnappliedTransactions = make(map[crypto.Sha256]types.TransactionMetadata)
	con.initialize()
	return con
}

func (c *Controller) PopBlock() {
	prev := c.ForkDB.GetBlock(&c.Head.Header.Previous)
	r := entity.ReversibleBlockObject{}
	//r.BlockNum = c.Head.BlockNum
	EosAssert(common.Empty(prev), &BlockValidateException{}, "attempt to pop beyond last irreversible block")
	errs := c.ReversibleBlocks.Find("BlockNum", c.Head.BlockNum, &r)
	if errs != nil {
		log.Error("PopBlock ReversibleBlocks Find is error :%s", errs.Error())
	}
	if !common.Empty(r) {
		c.ReversibleBlocks.Remove(&r)
	}

	if c.ReadMode == SPECULATIVE {
		//version 1.4
		//EosAssert(c.Head.SignedBlock!=nil, &BlockValidateException{}, "attempting to pop a block that was sparsely loaded from a snapshot")
		for _, trx := range c.Head.Trxs {
			c.UnappliedTransactions[crypto.Sha256(trx.SignedID)] = *trx
		}
	}
	c.Head = prev
	c.UndoSession.Undo()
}

func (c *Controller) SetApplayHandler(receiver common.AccountName, contract common.AccountName, action common.ActionName, handler func(a *ApplyContext)) {
	handlerKey := receiver + contract + action
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
		log.Error("Controller OnIrreversible ReversibleBlocks.GetIndex is error:%s", err)
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
	//emit( self.irreversible_block, s )
}

func (c *Controller) AbortBlock() {
	if c.Pending != nil {
		if c.ReadMode == SPECULATIVE {
			if c.Pending.PendingBlockState != nil {
				for _, trx := range c.Pending.PendingBlockState.Trxs {
					c.UnappliedTransactions[crypto.Sha256(trx.SignedID)] = *trx
				}
			}
		}
		c.Pending = c.Pending.Reset()
	}
}
func (c *Controller) StartBlock(when types.BlockTimeStamp, confirmBlockCount uint16) {
	pbi := common.BlockIdType(*crypto.NewSha256Nil())
	c.startBlock(when, confirmBlockCount, types.Incomplete, &pbi)
	c.ValidateDbAvailableSize()
}
func (c *Controller) startBlock(when types.BlockTimeStamp, confirmBlockCount uint16, s types.BlockStatus, producerBlockId *common.BlockIdType) {

	EosAssert(c.Pending == nil, &BlockValidateException{}, "pending block already exists")
	defer func() {
		if c.Pending != nil {
			if c.Pending.PendingValid {
				c.Pending = c.Pending.Reset()
			}
		}
	}()
	if !c.SkipDbSession(s) {
		EosAssert(uint32(c.DB.Revision()) == c.Head.BlockNum, &DatabaseException{}, "db revision is not on par with head block,Revision %v,BlockNum %v,ForkDB.Header().BlockNum %v",
			c.DB.Revision(), c.Head.BlockNum, c.ForkDB.Header().BlockNum)
		c.Pending = NewPendingState(c.DB)
	} else {
		c.Pending = NewDefaultPendingState()
	}

	c.Pending.BlockStatus = s
	c.Pending.ProducerBlockId = *producerBlockId
	c.Pending.PendingBlockState = types.NewBlockState2(&c.Head.BlockHeaderState, when) //TODO std::make_shared<block_state>( *head, when ); // promotes pending schedule (if any) to active
	c.Pending.PendingBlockState.InCurrentChain = true
	c.Pending.PendingBlockState.SetConfirmed(confirmBlockCount)
	wasPendingPromoted := c.Pending.PendingBlockState.MaybePromotePending()

	if c.ReadMode == DBReadMode(SPECULATIVE) || c.Pending.BlockStatus != types.BlockStatus(types.Incomplete) {
		gpo := c.GetGlobalProperties()
		if (gpo.ProposedScheduleBlockNum != 0 && gpo.ProposedScheduleBlockNum <= c.Pending.PendingBlockState.DposIrreversibleBlocknum) &&
			(len(c.Pending.PendingBlockState.PendingSchedule.Producers) == 0) && (!wasPendingPromoted) {
			if !c.RePlaying {
				log.Info("promoting proposed schedule (set in block %d) to pending; current block: %d lib: %d schedule: %v ",
					gpo.ProposedScheduleBlockNum, c.Pending.PendingBlockState.BlockNum, c.Pending.PendingBlockState.DposIrreversibleBlocknum, gpo.ProposedSchedule)
			}
			tmp := gpo.ProposedSchedule.ProducerScheduleType()
			ps := types.SharedProducerScheduleType{}
			ps.Version = tmp.Version
			ps.Producers = tmp.Producers
			c.Pending.PendingBlockState.SetNewProducers(ps)
			c.DB.Modify(gpo, func(i *entity.GlobalPropertyObject) {
				i.ProposedScheduleBlockNum = 0
				i.ProposedSchedule.Clear()
			})
			c.GpoCache[gpo.ID] = gpo
		}

		Try(func() {
			signedTransaction := c.GetOnBlockTransaction()
			onbtrx := types.NewTransactionMetadataBySignedTrx(&signedTransaction, 0)
			onbtrx.Implicit = true
			defer func(b bool) {
				c.InTrxRequiringChecks = b
			}(c.InTrxRequiringChecks)
			c.InTrxRequiringChecks = true
			c.pushTransaction(onbtrx, common.MaxTimePoint(), gpo.Configuration.MinTransactionCpuUsage, true)
		}).Catch(func(e Exception) {
			log.Error("Controller StartBlock exception:%s", e.DetailMessage())
			Throw(e)
		}).Catch(func(i interface{}) {
			//c++ nothing
		}).End()

		c.clearExpiredInputTransactions()
		c.updateProducersAuthority()
	}
	//c.Pending.PendingValid = false
}

func (c *Controller) pushReceipt(trx interface{}, status types.TransactionStatus, cpuUsageUs uint64, netUsage uint64) *types.TransactionReceipt {
	trxReceipt := types.TransactionReceipt{}
	tr := types.TransactionWithID{}
	switch trx.(type) {
	case common.TransactionIdType:
		tr.TransactionID = trx.(common.TransactionIdType)
	case types.PackedTransaction:
		pt := trx.(types.PackedTransaction)
		tr.PackedTransaction = &pt
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
	var trace *types.TransactionTrace
	returning, trace := false, (*types.TransactionTrace)(nil)

	Try(func() {
		trxContext := *NewTransactionContext(c, trx.Trx, trx.ID, common.Now())
		if c.SubjectiveCupLeeway != 0 {
			if c.Pending.BlockStatus == types.BlockStatus(types.Incomplete) {
				trxContext.Leeway = c.SubjectiveCupLeeway
			}
		}
		trxContext.Deadline = deadLine
		trxContext.ExplicitBilledCpuTime = explicitBilledCpuTime
		trxContext.BilledCpuTimeUs = int64(billedCpuTimeUs)

		trace = trxContext.Trace
		Try(func() {
			if trx.Implicit {
				trxContext.InitForImplicitTrx(0) //default value 0
				trxContext.CanSubjectivelyFail = false
			} else {
				skipRecording := (c.ReplayHeadTime != 0) && (trx.Trx.Expiration.ToTimePoint() <= c.ReplayHeadTime)
				trxContext.InitForInputTrx(uint64(trx.PackedTrx.GetUnprunableSize()), uint64(trx.PackedTrx.GetPrunableSize()),
					uint32(len(trx.Trx.Signatures)), skipRecording)
			}
			if trxContext.CanSubjectivelyFail && c.Pending.BlockStatus == types.Incomplete {
				c.CheckActorList(&trxContext.BillToAccounts)
			}
			trxContext.Delay = common.Microseconds(trx.Trx.DelaySec)
			checkTime := func() {}
			set := treeset.NewWith(types.PermissionLevelType, types.ComparePermissionLevel)
			if !c.SkipAuthCheck() && !trx.Implicit {
				c.Authorization.CheckAuthorization(trx.Trx.Actions,
					trx.RecoverKeys(&c.ChainID),
					set,
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
				tr := c.pushReceipt(*trx.PackedTrx, s, uint64(trxContext.BilledCpuTimeUs), trace.NetUsage)
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
				delete(c.UnappliedTransactions, crypto.Hash256(trx.SignedID))
			}

			returning = true
		}).Catch(func(ex Exception) {
			trace.Except = ex
			trace.ExceptPtr = ex
		}).End()
		if returning {
			return
		}
		if !failureIsSubjective(trace.Except) {
			delete(c.UnappliedTransactions, crypto.Sha256(trx.SignedID))
		}
		/*emit( c.accepted_transaction, trx )
		emit( c.applied_transaction, trace )*/
		trxContext.Undo()
		return
	}).FcCaptureAndRethrow("trace:%v", trace).End()
	return trace
}

func (c *Controller) GetGlobalProperties() *entity.GlobalPropertyObject {

	gpo := entity.GlobalPropertyObject{}
	gpo.ID = 0
	if c.GpoCache[gpo.ID] != nil {
		return c.GpoCache[gpo.ID]
	}
	err := c.DB.Find("id", gpo, &gpo)
	if err != nil {
		log.Error("GetGlobalProperties is error detail:%s", err)
	}
	return &gpo
}

func (c *Controller) GetDynamicGlobalProperties() (r *entity.DynamicGlobalPropertyObject) {
	dgpo := entity.DynamicGlobalPropertyObject{}
	dgpo.ID = 0
	err := c.DB.Find("id", dgpo, &dgpo)
	if err != nil {
		log.Error("GetDynamicGlobalProperties is error detail:%s", err)
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
	in := c.PendingBlockTime().AddUs(common.Microseconds(999999))
	trx.Expiration = common.NewTimePointSecTp(in)
	return trx
}
func (c *Controller) SkipDbSession(bs types.BlockStatus) bool {
	considerSkipping := (bs == types.Irreversible)
	return considerSkipping && !c.Config.DisableReplayOpts && !c.InTrxRequiringChecks
}

func (c *Controller) SkipDbSessions() bool {
	if c.Pending != nil {
		return c.SkipDbSession(c.Pending.BlockStatus)
	} else {
		return false
	}
}

func (c *Controller) SkipTrxChecks() (b bool) {
	b = c.LightValidationAllowed(c.Config.DisableReplayOpts)
	return
}

func (c *Controller) IsProducingBlock() bool {
	if c.Pending == nil {
		return false
	}
	return c.Pending.BlockStatus == types.Incomplete
}

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
		log.Error("Node data has been emptied is error:%s", err)
	}
}

func (c *Controller) GetUnappliedTransactions() []*types.TransactionMetadata {
	result := []*types.TransactionMetadata{}
	if c.ReadMode == SPECULATIVE {
		for _, v := range c.UnappliedTransactions {
			result = append(result, &v)
		}
	} else {
		log.Info("not empty unapplied_transactions in non-speculative mode")
		EosAssert(common.Empty(c.UnappliedTransactions), &TransactionException{}, "not empty unapplied_transactions in non-speculative mode")
	}
	return result
}

func (c *Controller) DropUnappliedTransaction(metadata *types.TransactionMetadata) {
	delete(c.UnappliedTransactions, crypto.Sha256(metadata.SignedID))
}

func (c *Controller) DropAllUnAppliedTransactions() {
	c.UnappliedTransactions = nil
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
		log.Error("Controller GetScheduledTransactions is error:%s", err)
	}
	if itr != nil {
		itr.Release()
	} /*else {
		log.Info("Controller GetScheduledTransactions byDelay is not found data")
	}*/
	return result
}
func (c *Controller) PushScheduledTransaction(trxId *common.TransactionIdType, deadLine common.TimePoint, billedCpuTimeUs uint32) *types.TransactionTrace {
	c.ValidateDbAvailableSize()
	return c.pushScheduledTransactionById(trxId, deadLine, billedCpuTimeUs, billedCpuTimeUs > 0)

}
func (c *Controller) pushScheduledTransactionById(sheduled *common.TransactionIdType,
	deadLine common.TimePoint,
	billedCpuTimeUs uint32, explicitBilledCpuTime bool) *types.TransactionTrace {

	gto := entity.GeneratedTransactionObject{}
	gto.TrxId = *sheduled
	c.DB.Find("byTrxId", gto, &gto)

	EosAssert(&gto != nil, &UnknownTransactionException{}, "unknown transaction")
	return c.pushScheduledTransactionByObject(&gto, deadLine, billedCpuTimeUs, explicitBilledCpuTime)
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
		log.Error("PushScheduleTransaction1 DecodeBytes is error :%s", err)
	}

	//trx := types.NewTransactionMetadataBySignedTrx(&dtrx,0) //TODO emit

	trace := &types.TransactionTrace{}
	if gtrx.Expiration < c.PendingBlockTime() {
		trace.ID = gtrx.TrxId
		trace.BlockNum = c.PendingBlockState().BlockNum
		trace.BlockTime = types.NewBlockTimeStamp(c.PendingBlockTime())
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
	returning := false
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
		returning = true
		//return trace
	}).Catch(func(ex Exception) {
		log.Error("PushScheduledTransaction is error:%s", ex.DetailMessage())
		cpuTimeToBillUs = trxContext.UpdateBilledCpuTime(common.Now())
		trace.Except = ex
		trace.ExceptPtr = ex
		trace.Elapsed = (common.Now() - trxContext.Start).TimeSinceEpoch()
	}).End()
	if returning {
		return trace
	}
	trxContext.Undo()
	if !failureIsSubjective(trace.Except) && gtrx.Sender != 0 { /*gtrx.Sender != account_name()*/
		log.Info("%v", trace.Except.DetailMessage())
		errorTrace := c.applyOnerror(gtrx, deadLine, trxContext.pseudoStart, &cpuTimeToBillUs, billedCpuTimeUs, explicitBilledCpuTime)
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
			rl.UpdateAccountUsage(&trxContext.BillToAccounts, uint32(types.NewBlockTimeStamp(c.PendingBlockTime())) /*.slot*/)
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

func (c *Controller) applyOnerror(gtrx *entity.GeneratedTransaction, deadline common.TimePoint, start common.TimePoint,
	cpuTimeToBillUs *uint32, billedCpuTimeUs uint32, explicitBilledCpuTime bool) *types.TransactionTrace {
	etrx := types.SignedTransaction{}
	action := types.Action{}
	action.Authorization = []types.PermissionLevel{{gtrx.Sender, common.DefaultConfig.ActiveName}}
	action.Data = gtrx.PackedTrx
	onError := NewOnError(gtrx.SenderId, gtrx.PackedTrx)
	action.Account = onError.GetAccount()
	action.Name = onError.GetName()
	etrx.Actions = append(etrx.Actions, &action)
	in := c.PendingBlockTime().AddUs(common.Microseconds(999999))
	etrx.Expiration = common.NewTimePointSecTp(in)
	blockId := c.HeadBlockId()
	etrx.SetReferenceBlock(&blockId)
	trxContext := NewTransactionContext(c, &etrx, etrx.ID(), start)
	trxContext.deadline = deadline
	trxContext.ExplicitBilledCpuTime = explicitBilledCpuTime
	trxContext.BilledCpuTimeUs = int64(billedCpuTimeUs)
	trace := trxContext.Trace
	returning := false
	Try(func() {
		trxContext.InitForImplicitTrx(0) //default
		trxContext.Published = gtrx.Published
		//trxContext.Trace.ActionTraces
		at := types.ActionTrace{}
		trxContext.Trace.ActionTraces = append(trxContext.Trace.ActionTraces, at)
		tr := trxContext.Trace.ActionTraces[len(trxContext.Trace.ActionTraces)-1]
		last := etrx.Actions[len(etrx.Actions)-1]
		trxContext.DispatchAction(&tr, last, gtrx.Sender, false, 0) //default false 0
		trxContext.Finalize()
		defer func() {}() //TODO
		//pushReceipt(trx interface{}, status types.TransactionStatus, cpuUsageUs uint64, netUsage uint64) *types.TransactionReceipt
		t := c.pushReceipt(gtrx.TrxId, types.TransactionStatusSoftFail, uint64(trxContext.BilledCpuTimeUs), trace.NetUsage)
		trace.Receipt = t.TransactionReceiptHeader
		trxContext.Squash()
		//restore.cancel()

		returning = true
	}).Catch(func(e Exception) {
		t := trxContext.UpdateBilledCpuTime(common.Now())
		cpuTimeToBillUs = &t
		trace.Except = e
		trace.ExceptPtr = e
	}).End()

	if returning {
		return trace
	}

	return trace
}
func (c *Controller) RemoveScheduledTransaction(gto *entity.GeneratedTransactionObject) {
	c.ResourceLimits.AddPendingRamUsage(gto.Payer, int64(9)+int64(len(gto.PackedTrx)))
	c.DB.Remove(gto)
}

func failureIsSubjective(e Exception) bool {
	code := e.Code()
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

	EosAssert(c.Pending != nil, &BlockValidateException{}, "it is not valid to finalize when there is no pending block")

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
	cpu.ContractRate.Denominator = 100
	cpu.ExpandRate.Numerator = 999
	cpu.ExpandRate.Denominator = 1000

	net := types.ElasticLimitParameters{}
	netTarget := common.EosPercent(uint64(chainConfig.MaxBlockNetUsage), chainConfig.TargetBlockNetUsagePct)
	net.Target = netTarget
	net.Max = uint64(chainConfig.MaxBlockNetUsage)
	net.Periods = common.DefaultConfig.BlockSizeAverageWindowMs / uint32(common.DefaultConfig.BlockIntervalMs)
	net.MaxMultiplier = m

	net.ContractRate.Numerator = 99
	net.ContractRate.Denominator = 100
	net.ExpandRate.Numerator = 999
	net.ExpandRate.Denominator = 1000
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
	p.SignedBlock.SignedBlockHeader = p.Header
}

func (c *Controller) applyBlock(b *types.SignedBlock, s types.BlockStatus) {
	Try(func() {
		EosAssert(len(b.BlockExtensions) == 0, &BlockValidateException{}, "no supported extensions")
		producerBlockId := b.BlockID()
		c.startBlock(b.Timestamp, b.Confirmed, s, &producerBlockId)

		trace := &types.TransactionTrace{}
		for _, receipt := range b.Transactions {
			numPendingReceipts := len(c.Pending.PendingBlockState.SignedBlock.Transactions)
			if !common.Empty(receipt.Trx.PackedTransaction) {
				pt := receipt.Trx.PackedTransaction
				mtrx := types.NewTransactionMetadata(pt)
				trace = c.pushTransaction(mtrx, common.TimePoint(common.MaxMicroseconds()), receipt.CpuUsageUs, true)
			} else if !common.Empty(receipt.Trx.TransactionID) {
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
			EosAssert(trxReceipt.TransactionReceiptHeader == receipt.TransactionReceiptHeader, &BlockValidateException{}, "receipt does not match,producer_receipt:%v", receipt, "validator_receipt:%v", trxReceipt)
		}
		c.FinalizeBlock()

		EosAssert(producerBlockId == c.Pending.PendingBlockState.Header.BlockID(), &BlockValidateException{},
			"Block ID does not match,producer_block_id:%v", producerBlockId, "validator_block_id:%v", c.Pending.PendingBlockState.Header.BlockID())

		c.Pending.PendingBlockState.Header.ProducerSignature = b.ProducerSignature
		c.CommitBlock(false)
		return
	}).Catch(func(ex Exception) {
		log.Error("controller ApplyBlock is error:%s", ex.DetailMessage())
		c.AbortBlock()
	}).End()
}

func (c *Controller) CommitBlock(addToForkDb bool) {
	defer func() {
		if c.Pending != nil && c.Pending.PendingValid {
			c.Pending = c.Pending.Reset()
		}
	}()
	Try(func() {
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
	}).Catch(func(e interface{}) {
		c.AbortBlock()
		Throw(e)
	}).End()
	c.Pending.Push()
	c.Pending.PendingValid = true
	//log.Info("commitBlock success!")

}

func (c *Controller) PushBlock(b *types.SignedBlock, s types.BlockStatus) {
	EosAssert(c.Pending == nil, &BlockValidateException{}, "it is not valid to push a block when there is a pending block")
	defer func() {
		c.TrustedProducerLightValidation = false
	}()

	Try(func() {
		EosAssert(b != nil, &BlockValidateException{}, "trying to push empty block")
		EosAssert(s != types.Incomplete, &BlockLogException{}, "invalid block status for a completed block")
		//TODO: add to forkdb
		c.PreAcceptedBlock.Emit(b)
		//emit(self.pre_accepted_block, b )
		//trust := !c.Config.ForceAllChecks && (s == types.Irreversible || s == types.Validated)
		//newHeaderState := c.ForkDB.AddSignedBlockState(b, trust)
		exist, _ := c.Config.TrustedProducers.Find(func(index int, value interface{}) bool {
			return c.Config.TrustedProducers.GetComparator()(value, b.Producer) == 0
		})
		//exist, _ := c.Config.trustedProducers.Find(&b.Producer)
		if exist != -1 {
			c.TrustedProducerLightValidation = true
		}
		//c.AcceptedBlockHeader.Emit(newHeaderState)
		//emit( self.accepted_block_header, new_header_state )
		if c.ReadMode != IRREVERSIBLE {
			c.maybeSwitchForks(s)
		}
		if s == types.Irreversible {
			//c.IrreversibleBlock.Emit(newHeaderState)
			//emit( self.irreversible_block, new_header_state )
		}
	}).FcLogAndRethrow().End()

}

func (c *Controller) PushConfirmation(hc *types.HeaderConfirmation) {
	EosAssert(c.Pending == nil, &BlockValidateException{}, "it is not valid to push a confirmation when there is a pending block")
	c.ForkDB.Add(hc)
	//c.AcceptedConfirmation.Emit(hc)
	//emit( c.accepted_confirmation, hc )
	if c.ReadMode != IRREVERSIBLE {
		c.maybeSwitchForks(types.Complete)
	}
}

func (c *Controller) maybeSwitchForks(s types.BlockStatus) {
	newHead := c.ForkDB.Head
	if newHead.Header.Previous == c.Head.BlockId {
		Try(func() {
			c.applyBlock(newHead.SignedBlock, s)
			c.ForkDB.MarkInCurrentChain(newHead, true)
			c.ForkDB.SetValidity(newHead, true)
			c.Head = newHead
		}).Catch(func(e Exception) {
			c.ForkDB.SetValidity(newHead, false)
			EosThrow(e, "maybeSwitchForks is error:%v", e.DetailMessage())
		}).End()
	} else if newHead.BlockId != c.Head.BlockId {
		log.Info("switching forks from: %v (block number %v) to %v (block number %v)", c.Head.BlockId, c.Head.BlockNum, newHead.BlockId, newHead.BlockNum)
		branches := c.ForkDB.FetchBranchFrom(&newHead.BlockId, &c.Head.BlockId)

		for i := 0; i < len(branches.second); i++ {
			c.ForkDB.MarkInCurrentChain(&branches.second[i], false)
			c.PopBlock()
		}
		length := len(branches.second) - 1
		EosAssert(c.HeadBlockId() == branches.second[length].Header.Previous, &ForkDatabaseException{}, "loss of sync between fork_db and chainbase during fork switch")
		le := len(branches.first) - 1
		for i := le; i >= 0; i-- {
			itr := branches.first[i]
			var except Exception
			Try(func() {
				if itr.Validated {
					c.applyBlock(itr.SignedBlock, types.Validated)
				} else {
					c.applyBlock(itr.SignedBlock, types.Complete)
				}
				c.Head = &itr
				c.ForkDB.MarkInCurrentChain(&itr, true)
			}).Catch(func(e Exception) {
				except = e
			}).End()
			if except == nil {
				log.Error("exception thrown while switching forks :%s", except.DetailMessage())
				c.ForkDB.SetValidity(&itr, false)
				// pop all blocks from the bad fork
				// ritr base is a forward itr to the last block successfully applied
				for j := i; j < len(branches.first); j++ {
					c.ForkDB.MarkInCurrentChain(&branches.first[j], true)
					c.PopBlock()
				}
				EosAssert(c.HeadBlockId() == branches.second[len(branches.second)-1].Header.Previous, &ForkDatabaseException{}, "loss of sync between fork_db and chainbase during fork switch reversal")
				// re-apply good blocks
				l := len(branches.second) - 1
				for end := l; end >= 0; end-- {
					c.applyBlock(branches.second[end].SignedBlock, types.Validated)
					c.Head = &branches.second[end]
					c.ForkDB.MarkInCurrentChain(&branches.second[end], true)
				}
				EosThrow(except, "maybeSwitchForks is error:%v", except.DetailMessage())
			}
			log.Info("successfully switched fork to new head %v", newHead.BlockId)
		}
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
		log.Error("GetAccount is error :%s", err)
	}
	return &accountObj
}

func (c *Controller) GetAuthorizationManager() *AuthorizationManager { return c.Authorization }

func (c *Controller) GetMutableAuthorizationManager() *AuthorizationManager { return c.Authorization }

//c++ flat_set<account_name> map[common.AccountName]interface{}
func (c *Controller) GetActorWhiteList() *treeset.Set {
	return &c.Config.ActorWhitelist
}

func (c *Controller) GetActorBlackList() *treeset.Set {
	return &c.Config.ActorBlacklist
}

func (c *Controller) GetContractWhiteList() *treeset.Set {
	return &c.Config.ContractWhitelist
}

func (c *Controller) GetContractBlackList() *treeset.Set {
	return &c.Config.ContractBlacklist
}

func (c *Controller) GetActionBlockList() *treeset.Set { return &c.Config.ActionBlacklist }

func (c *Controller) GetKeyBlackList() *treeset.Set { return &c.Config.KeyBlacklist }

func (c *Controller) SetActorWhiteList(params *treeset.Set) {
	c.Config.ActorWhitelist = *params
}

func (c *Controller) SetActorBlackList(params *treeset.Set) {
	c.Config.ActorBlacklist = *params
}

func (c *Controller) SetContractWhiteList(params *treeset.Set) {
	c.Config.ContractWhitelist = *params
}

func (c *Controller) SetContractBlackList(params *treeset.Set) {
	c.Config.ContractBlacklist = *params
}

func (c *Controller) SetActionBlackList(params *treeset.Set) {
	c.Config.ActionBlacklist = *params
}

func (c *Controller) SetKeyBlackList(params *treeset.Set) {
	c.Config.KeyBlacklist = *params
}

func (c *Controller) HeadBlockNum() uint32 { return c.Head.BlockNum }

func (c *Controller) HeadBlockTime() common.TimePoint { return c.Head.Header.Timestamp.ToTimePoint() }

func (c *Controller) HeadBlockId() common.BlockIdType { return c.Head.BlockId }

func (c *Controller) HeadBlockProducer() common.AccountName { return c.Head.Header.Producer }

func (c *Controller) HeadBlockHeader() *types.BlockHeader { return &c.Head.Header.BlockHeader }

func (c *Controller) HeadBlockState() *types.BlockState { return c.Head }

func (c *Controller) ForkDbHeadBlockNum() uint32 { return c.ForkDB.Header().BlockNum }

func (c *Controller) ForkDbHeadBlockId() common.BlockIdType { return c.ForkDB.Head.BlockId }

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
	return nil
}

func (c *Controller) PendingBlockTime() common.TimePoint {
	EosAssert(c.Pending != nil, &BlockValidateException{}, "no pending block")
	return c.Pending.PendingBlockState.Header.Timestamp.ToTimePoint()
}

func (c *Controller) PendingProducerBlockId() common.BlockIdType {
	EosAssert(c.Pending != nil, &BlockValidateException{}, "no pending block")
	return c.Pending.ProducerBlockId
}

func (c *Controller) ActiveProducers() *types.ProducerScheduleType {
	if c.Pending == nil {
		return &c.Head.ActiveSchedule
	}
	return &c.Pending.PendingBlockState.ActiveSchedule
}

func (c *Controller) PendingProducers() *types.ProducerScheduleType {
	if c.Pending == nil {
		return &c.Head.PendingSchedule
	}
	return &c.Pending.PendingBlockState.PendingSchedule
}

func (c *Controller) ProposedProducers() types.ProducerScheduleType {
	gpo := c.GetGlobalProperties()
	if gpo.ProposedScheduleBlockNum == 0 {
		return types.ProducerScheduleType{}
	}
	return *gpo.ProposedSchedule.ProducerScheduleType()
}

func (c *Controller) LightValidationAllowed(dro bool) (b bool) {
	if c.Pending == nil || c.InTrxRequiringChecks {
		return false
	}

	pbStatus := c.Pending.BlockStatus
	considerSkippingOnReplay := (pbStatus == types.Irreversible || pbStatus == types.Validated) && !dro

	considerSkippingOnvalidate := (pbStatus == types.Complete && c.Config.BlockValidationMode == LIGHT)

	return considerSkippingOnReplay || considerSkippingOnvalidate
}

func (c *Controller) LastIrreversibleBlockNum() uint32 {
	if c.Head.BftIrreversibleBlocknum > c.Head.DposIrreversibleBlocknum {
		return c.Head.BftIrreversibleBlocknum
	} else {
		return c.Head.DposIrreversibleBlocknum
	}
}

func (c *Controller) LastIrreversibleBlockId() common.BlockIdType {
	libNum := c.LastIrreversibleBlockNum()
	bso := entity.BlockSummaryObject{}
	bso.Id = common.IdType(uint16(libNum))
	idx, err := c.DataBase().GetIndex("id", &entity.BlockSummaryObject{})
	err = idx.Find(bso, &bso)
	if err != nil {
		log.Error("controller LastIrreversibleBlockId is error:%s", err)
	}
	if types.NumFromID(&bso.BlockId) == libNum {
		return bso.BlockId
	}
	return c.FetchBlockByNumber(libNum).BlockID()
}

func (c *Controller) FetchBlockByNumber(blockNum uint32) *types.SignedBlock {
	r := (*types.SignedBlock)(nil)
	Try(func() {
		blkState := c.ForkDB.GetBlockInCurrentChainByNum(blockNum)
		if blkState != nil {
			r = blkState.SignedBlock
			return
		}

		r = c.Blog.ReadBlockByNum(blockNum)
		return

	}).FcCaptureAndRethrow("blockNum:%d", blockNum).End()

	return r
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
	return nil
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
	EosAssert(!common.Empty(signedBlk), &UnknownBlockException{}, "Could not find block: %d", blockNum)
	return signedBlk.BlockID()
}

func (c *Controller) CheckContractList(code common.AccountName) {
	if c.Config.ContractWhitelist.Size() > 0 {
		exist, _ := c.Config.ContractWhitelist.Find(func(index int, value interface{}) bool {
			return c.Config.ContractWhitelist.GetComparator()(value, &code) == 0
		})
		EosAssert(exist != -1, &ContractWhitelistException{}, "account %d is not on the contract whitelist", code)
	} else if c.Config.ContractBlacklist.Size() > 0 {
		exist, _ := c.Config.ContractBlacklist.Find(func(index int, value interface{}) bool {
			return c.Config.ContractBlacklist.GetComparator()(value, &code) == 0
		})
		EosAssert(exist == -1, &ContractBlacklistException{}, "account %d is on the contract blacklist", code)
	}
}

func (c *Controller) CheckActionList(code common.AccountName, action common.ActionName) {
	if c.Config.ActionBlacklist.Size() > 0 {
		abl := common.MakePair(code, action)
		//key := Hash(abl)
		exist, _ := c.Config.ActionBlacklist.Find(func(index int, value interface{}) bool {
			return c.Config.ActionBlacklist.GetComparator()(value, &abl) == 0
		})

		EosAssert(exist == -1, &ActionBlacklistException{}, "action '%v::%v' is on the action blacklist", code, action)
	}
}

func (c *Controller) CheckKeyList(key *ecc.PublicKey) {
	if c.Config.KeyBlacklist.Size() > 0 {
		exist, _ := c.Config.KeyBlacklist.Find(func(index int, value interface{}) bool {
			return c.Config.KeyBlacklist.GetComparator()(value, key) == 0
		})
		EosAssert(exist == -1, &KeyBlacklistException{}, "public key %v is on the key blacklist", key)
	}
}

func (c *Controller) IsRamBillingInNotifyAllowed() bool {
	return !c.IsProducingBlock() || c.Config.AllowRamBillingInNotify
}

func (c *Controller) AddResourceGreyList(name *common.AccountName) {
	c.Config.ResourceGreylist.AddItem(*name)
}

func (c *Controller) RemoveResourceGreyList(name *common.AccountName) {
	c.Config.ResourceGreylist.Remove(*name)
}

func (c *Controller) IsResourceGreylisted(name *common.AccountName) bool {
	exist, _ := c.Config.ResourceGreylist.Find(func(index int, value interface{}) bool {
		return c.Config.ResourceGreylist.GetComparator()(value, *name) == 0
	})
	if exist > -1 {
		return true
	} else {
		return false
	}
}
func (c *Controller) GetResourceGreyList() treeset.Set {
	return c.Config.ResourceGreylist
}

func (c *Controller) ValidateReferencedAccounts(t *types.Transaction) {
	ac := entity.AccountObject{}
	for _, a := range t.ContextFreeActions {
		ac.Name = a.Account
		err := c.DB.Find("byName", ac, &ac)
		EosAssert(err == nil, &TransactionException{}, "action's code account '%v' does not exist", a.Account)
		EosAssert(len(a.Authorization) == 0, &TransactionException{}, "context-free actions cannot have authorizations")
	}
	oneAuth := false
	for _, a := range t.Actions {
		ac.Name = a.Account
		err := c.DB.Find("byName", ac, &ac)
		EosAssert(err == nil, &TransactionException{}, "action's code account '%v' does not exist", a.Account)
		for _, auth := range a.Authorization {
			oneAuth = true
			ac.Name = auth.Actor
			err := c.DB.Find("byName", ac, &ac)
			EosAssert(err == nil, &TransactionException{}, "action's authorizing actor '%v' does not exist", auth.Actor)
			EosAssert(c.Authorization.FindPermission(&auth) != nil, &TransactionException{}, "action's authorizations include a non-existent permission: %v", auth)
		}
	}
	EosAssert(oneAuth, &TxNoAuths{}, "transaction must have at least one authorization")
}

func (c *Controller) ValidateExpiration(t *types.Transaction) {
	chainConfiguration := c.GetGlobalProperties().Configuration
	//log.Info("ValidateExpiration t.Expiration.ToTimePoint():%v,c.PendingBlockTime():%v", t.Expiration.ToTimePoint(), c.PendingBlockTime())
	EosAssert(t.Expiration.ToTimePoint() >= c.PendingBlockTime(),
		&ExpiredTxException{}, "transaction has expired, expiration is %v and pending block time is %v",
		t.Expiration, c.PendingBlockTime())
	EosAssert(t.Expiration.ToTimePoint() <= c.PendingBlockTime().AddUs(common.Seconds(int64(chainConfiguration.MaxTrxLifetime))),
		&TxExpTooFarException{}, "Transaction expiration is too far in the future relative to the reference time of %v, expiration is %v and the maximum transaction lifetime is %v seconds",
		t.Expiration, c.PendingBlockTime(), chainConfiguration.MaxTrxLifetime)
}

func (c *Controller) ValidateTapos(t *types.Transaction) {
	bso := entity.BlockSummaryObject{}
	bso.Id = common.IdType(t.RefBlockNum)
	err := c.DB.Find("id", bso, &bso)
	if err != nil {
		log.Error("ValidateTapos Is Error:%s", err)
	}
	EosAssert(t.VerifyReferenceBlock(&bso.BlockId), &InvalidRefBlockException{},
		"Transaction's reference block did not match. Is this transaction from a different fork? taposBlockSummary:%v", bso)
}

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
	t := entity.TransactionObject{}
	t.TrxID = *id
	return nil == c.DB.Find("byTrxId", t, &t)
}

func (c *Controller) SetProposedProducers(producers []types.ProducerKey) int64 {
	gpo := c.GetGlobalProperties()
	curBlockNum := c.HeadBlockNum() + 1
	if gpo.ProposedScheduleBlockNum != 0 {
		if gpo.ProposedScheduleBlockNum != curBlockNum {
			return -1
		}
		if compare(producers, gpo.ProposedSchedule.Producers) {
			return -1
		}
	}
	sch := types.ProducerScheduleType{}
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
	c.DB.Modify(gpo, func(p *entity.GlobalPropertyObject) {
		p.ProposedScheduleBlockNum = curBlockNum
		tmp := p.ProposedSchedule.SharedProducerScheduleType(sch)
		p.ProposedSchedule = *tmp
	})
	c.GpoCache[gpo.ID] = gpo

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

func (c *Controller) SkipAuthCheck() bool { return c.LightValidationAllowed(c.Config.ForceAllChecks) }

func (c *Controller) ContractsConsole() bool { return c.Config.ContractsConsole }

func (c *Controller) GetChainId() common.ChainIdType { return c.ChainID }

func (c *Controller) GetReadMode() DBReadMode { return c.ReadMode }

func (c *Controller) GetValidationMode() ValidationMode { return c.Config.BlockValidationMode }

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
	account.CreationDate = types.NewBlockTimeStamp(c.Config.Genesis.InitialTimestamp)
	account.Privileged = isPrivileged
	if name == common.AccountName(common.DefaultConfig.SystemAccountName) {
		abiDef := abi.AbiDef{}
		account.SetAbi(EosioContractAbi(abiDef))
	}
	err := c.DB.Insert(&account)
	if err != nil {
		log.Error("CreateNativeAccount Insert Is Error:%s", err)
	}

	aso := entity.AccountSequenceObject{}
	aso.Name = name
	c.DB.Insert(&aso)

	ownerPermission := c.Authorization.CreatePermission(name, common.PermissionName(common.DefaultConfig.OwnerName), 0, owner, c.Config.Genesis.InitialTimestamp)

	activePermission := c.Authorization.CreatePermission(name, common.PermissionName(common.DefaultConfig.ActiveName), PermissionIdType(ownerPermission.ID), active, c.Config.Genesis.InitialTimestamp)

	c.ResourceLimits.InitializeAccount(name)
	ramDelta := uint64(common.DefaultConfig.OverheadPerAccountRamBytes)
	ramDelta += 2 * common.BillableSizeV("permission_object") //::billable_size_v<permission_object>
	ramDelta += ownerPermission.Auth.GetBillableSize()
	ramDelta += activePermission.Auth.GetBillableSize()
	c.ResourceLimits.AddPendingRamUsage(name, int64(ramDelta))
	c.ResourceLimits.VerifyAccountRamUsage(name)
}

func (c *Controller) initializeForkDB() {

	gs := c.Config.Genesis
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
	genHeader.ProducerToLastProduced = *types.NewAccountNameUint32Map()
	genHeader.ProducerToLastImpliedIrb = *types.NewAccountNameUint32Map()
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

	for i := 0; i < 0x10000; i++ {
		bso := entity.BlockSummaryObject{}
		err := c.DB.Insert(&bso)
		if err != nil {
			log.Error("Controller initializeDatabase Insert BlockSummaryObject is error:%s", err)
		}
	}
	b := entity.BlockSummaryObject{}
	b.Id = 1

	err := c.DB.Find("id", b, &b)
	c.DB.Modify(&b, func(bs *entity.BlockSummaryObject) {
		bs.BlockId = c.Head.BlockId
	})
	gi := c.Config.Genesis.Initial()
	gi.Validate() //check config
	gpo := entity.GlobalPropertyObject{}
	gpo.Configuration = gi
	err = c.DB.Insert(&gpo)
	c.GpoCache = make(map[common.IdType]*entity.GlobalPropertyObject)
	c.GpoCache[gpo.ID] = &gpo
	if err != nil {
		log.Error("Controller initializeDatabase insert GlobalPropertyObject is error:%s", err)
		EosAssert(err == nil, &DatabaseException{}, "Controller initializeDatabase is error :%s", err)
	}
	dgpo := entity.DynamicGlobalPropertyObject{}
	dgpo.ID = 0
	err = c.DB.Insert(&dgpo)
	if err != nil {
		log.Error("Controller initializeDatabase insert DynamicGlobalPropertyObject is error:%s", err)
	}

	c.ResourceLimits.InitializeDatabase()
	systemAuth := types.NewAuthority(c.Config.Genesis.InitialKey, 0)
	c.CreateNativeAccount(common.DefaultConfig.SystemAccountName, systemAuth, systemAuth, true)
	emptyAuthority := types.Authority{}
	emptyAuthority.Threshold = 1
	activeProducersAuthority := types.Authority{}
	activeProducersAuthority.Threshold = 1
	p := types.PermissionLevelWeight{types.PermissionLevel{common.DefaultConfig.SystemAccountName, common.DefaultConfig.ActiveName}, 1}
	activeProducersAuthority.Accounts = append(activeProducersAuthority.Accounts, p)
	c.CreateNativeAccount(common.DefaultConfig.NullAccountName, emptyAuthority, emptyAuthority, false)
	c.CreateNativeAccount(common.DefaultConfig.ProducersAccountName, emptyAuthority, activeProducersAuthority, false)
	activePermission := c.Authorization.GetPermission(&types.PermissionLevel{common.DefaultConfig.ProducersAccountName, common.DefaultConfig.ActiveName})

	majorityPermission := c.Authorization.CreatePermission(common.DefaultConfig.ProducersAccountName,
		common.DefaultConfig.MajorityProducersPermissionName,
		PermissionIdType(activePermission.ID),
		activeProducersAuthority,
		c.Config.Genesis.InitialTimestamp)

	c.Authorization.CreatePermission(common.DefaultConfig.ProducersAccountName,
		common.DefaultConfig.MinorityProducersPermissionName,
		PermissionIdType(majorityPermission.ID),
		activeProducersAuthority,
		c.Config.Genesis.InitialTimestamp)

	//log.Info("initializeDatabase print:%v,%v", majorityPermission.ID, minorityPermission.ID)
}

func (c *Controller) initialize() {
	if common.Empty(c.Head) {
		c.initializeForkDB()
		end := c.Blog.ReadHead()
		if common.Empty(end) && end.BlockNumber() > 1 {
			endTime := end.Timestamp.ToTimePoint()
			c.RePlaying = true
			c.ReplayHeadTime = endTime
			log.Info("existing block log, attempting to replay :%d blocks", end.BlockNumber())
			for next := c.Blog.ReadBlockByNum(c.Head.BlockNum + 1); next != nil; {
				c.PushBlock(next, types.Irreversible)
				if next.BlockNumber()%100 == 0 {
					log.Info("%d blocks replayed", next.BlockNumber())
				}
			}
			log.Info("%d blocks replayed", c.Head.BlockNum)
			c.DB.SetRevision(int64(c.Head.BlockNum))
			rev := 0
			r := entity.ReversibleBlockObject{}
			for {
				rev++
				r.BlockNum = c.HeadBlockNum() + 1
				err := c.ReversibleBlocks.Find("blockNum", r, &r)
				if err != nil {
					break
				}
				c.PushBlock(r.GetBlock(), types.Validated)
			}
			log.Info("%d reversible blocks replayed", rev)

			c.RePlaying = false
			c.ReplayHeadTime = common.TimePoint(0)
		} else if !common.Empty(end) {
			c.Blog.ResetToGenesis(c.Config.Genesis, c.Head.SignedBlock)
		}
	}
	rbi := entity.ReversibleBlockObject{}
	ubi, err := c.ReversibleBlocks.GetIndex("byNum", &rbi)
	if err != nil {
		fmt.Errorf("initialize database ReversibleBlocks GetIndex is error :%s", err)
		EosAssert(err == nil, &DatabaseException{}, "Controller initialize is error :%s", err)
	}
	//c++ rbegin and rend
	objitr := ubi.End()
	if objitr != ubi.Begin() {
		objitr.Prev()
		r := entity.ReversibleBlockObject{}
		objitr.Data(&r)
		EosAssert(r.BlockNum == c.Head.BlockNum, &ForkDatabaseException{},
			"reversible block database is inconsistent with fork database, replay blockchain", c.Head.BlockNum, r.BlockNum)
	} else {
		end := c.Blog.ReadHead()
		EosAssert(end != nil && end.BlockNumber() == c.Head.BlockNum, &ForkDatabaseException{},
			"fork database exists but reversible block database does not, replay blockchain", end.BlockNumber(), c.Head.BlockNum)
	}
	EosAssert(uint32(c.DB.Revision()) >= c.Head.BlockNum, &ForkDatabaseException{}, "fork database is inconsistent with shared memory", c.DB.Revision(), c.Head.BlockNum)
	for uint32(c.DB.Revision()) > c.Head.BlockNum {
		c.DB.Undo()
	}
	//log.Info("controller initialize finished")
}

func (c *Controller) clearExpiredInputTransactions() {
	transactionIdx, err := c.DB.GetIndex("byExpiration", &entity.TransactionObject{})
	now := c.PendingBlockTime()
	t := entity.TransactionObject{}
	for !transactionIdx.Empty() {
		err = transactionIdx.Begin().Data(&t)
		if err != nil {
			log.Error("controller clearExpiredInputTransactions transactionIdx.Begin() is error: %s", err)
			EosAssert(err == nil, &DatabaseException{}, "Controller clearExpiredInputTransactions is error :%s", err)
			return
		}
		if now > t.Expiration.ToTimePoint() {
			c.DB.Remove(&t)
		} else {
			break
		}
	}
}

func (c *Controller) CheckActorList(actors *treeset.Set) {
	if c.Config.ActorWhitelist.Size() > 0 {
		itr := actors.Iterator()
		for itr.Next() {
			exist, _ := c.Config.ActorWhitelist.Find(func(index int, value interface{}) bool {
				return c.Config.ActorWhitelist.GetComparator()(value, itr.Value()) == 0
			})
			EosAssert(exist != -1, &ActorWhitelistException{},
				"authorizing actor(s) in transaction are not on the actor whitelist: %v", actors)
		}
	} else if c.Config.ActorBlacklist.Size() > 0 {
		itr := actors.Iterator()
		for itr.Next() {
			exist, _ := c.Config.ActionBlacklist.Find(func(index int, value interface{}) bool {
				return c.Config.ActionBlacklist.GetComparator()(value, itr.Value()) == 0
			})
			EosAssert(exist == -1, &ActorBlacklistException{},
				"authorizing actor(s) in transaction are on the actor blacklist: %v", actors)
		}
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
			c.DB.Modify(permission, func(param *entity.PermissionObject) {
				param.Auth = auth.ToSharedAuthority()
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
	err := c.DB.Find("id", bso, &bso)
	if err != nil {
		log.Error("Controller createBlockSummary is error:%s", err)
		EosAssert(err == nil, &DatabaseException{}, "Controller createBlockSummary is error :%s", err)
	}
	c.DB.Modify(&bso, func(b *entity.BlockSummaryObject) {
		b.BlockId = *id
	})
}

func (c *Controller) initConfig() *Controller {
	c.Config = Config{
		BlocksDir:               common.DefaultConfig.DefaultBlocksDirName,
		StateDir:                common.DefaultConfig.DefaultStateDirName,
		ReversibleDir:           common.DefaultConfig.DefaultReversibleBlocksDirName,
		StateSize:               common.DefaultConfig.DefaultStateSize,
		StateGuardSize:          common.DefaultConfig.DefaultStateGuardSize,
		ReversibleCacheSize:     common.DefaultConfig.DefaultReversibleCacheSize,
		ReversibleGuardSize:     common.DefaultConfig.DefaultReversibleGuardSize,
		ReadOnly:                false,
		ForceAllChecks:          false,
		DisableReplayOpts:       false,
		ContractsConsole:        false,
		AllowRamBillingInNotify: false,
		//vmType:              common.DefaultConfig.DefaultWasmRuntime, //TODO
		ReadMode:            SPECULATIVE,
		BlockValidationMode: FULL,
		Genesis:             types.NewGenesisState(),
		//TODO tmp code
		ActorWhitelist:    *treeset.NewWith(common.TypeName, common.CompareName),
		ActorBlacklist:    *treeset.NewWith(common.TypeName, common.CompareName),
		ContractWhitelist: *treeset.NewWith(common.TypeName, common.CompareName),
		ContractBlacklist: *treeset.NewWith(common.TypeName, common.CompareName),
		ActionBlacklist:   *treeset.NewWith(common.TypePair, common.ComparePair),
		KeyBlacklist:      *treeset.NewWith(ecc.TypePubKey, ecc.ComparePubKey),
		ResourceGreylist:  *treeset.NewWith(common.TypeName, common.CompareName),
		TrustedProducers:  *treeset.NewWith(common.TypeName, common.CompareName),
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
