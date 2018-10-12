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
	DB                   *database.LDataBase
	DbSession            *database.Session
	ReversibleBlocks     *database.LDataBase
	Blog                 string //TODO
	Pending              *types.PendingState
	Head                 types.BlockState
	ForkDB               *types.ForkDatabase
	WasmIf               *exec.WasmInterface
	ResourceLimists      *ResourceLimitsManager
	Authorization        *AuthorizationManager
	Config               Config //local	Config
	ChainID              common.ChainIdType
	RePlaying            bool
	ReplayHeadTime       common.TimePoint //optional<common.Tstamp>
	ReadMode             DBReadMode
	InTrxRequiringChecks bool                //if true, checks that are normally skipped on replay (e.g. auth checks) cannot be skipped
	SubjectiveCupLeeway  common.Microseconds //optional<common.Tstamp>

	ApplyHandlers         map[common.AccountName]map[HandlerKey]v
	UnAppliedTransactions map[crypto.Sha256]types.TransactionMetadata
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
	con := &Controller{InTrxRequiringChecks: false, RePlaying: false}
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
func (self *Controller) OnIrreversible(b *types.BlockState) {

}

func (self *Controller) PopBlock() {
	prev := self.ForkDB.GetBlock(&self.Head.Header.Previous)
	r := types.ReversibleBlockObject{}
	r.BlockNum = self.Head.BlockNum
	itr, errs := self.ReversibleBlocks.Find("NUM", r)

	if errs != nil {
		log.Error("PopBlock ReversibleBlocks Find is error,detail:", errs)
	}
	defer itr.Release()

	if itr.Next() {
		self.ReversibleBlocks.Remove(itr.First())
	}
	if self.ReadMode == SPECULATIVE {
		var trx []*types.TransactionMetadata = self.Head.Trxs
		step := 0
		for ; step < len(trx); step++ {
			self.UnAppliedTransactions[crypto.Sha256(trx[step].SignedID)] = *trx[step]
		}
	}
	self.Head = *prev
	self.DbSession.Undo() //TODO
}

func (self *Controller) SetApplayHandler(receiver common.AccountName, contract common.AccountName, action common.ActionName, handler func(a *ApplyContext)) {
	hk := NewHandlerKey(common.ScopeName(contract), action)
	first := make(map[common.AccountName]map[HandlerKey]v)
	secend := make(map[HandlerKey]v)
	secend[hk] = handler
	first[receiver] = secend
	self.ApplyHandlers[receiver] = secend
}

func (self *Controller) AbortBlock() {
	if self.Pending != nil {
		if self.ReadMode == SPECULATIVE {
			trx := append(self.Pending.PendingBlockState.Trxs)
			step := 0
			for ; step < len(trx); step++ {
				self.UnAppliedTransactions[crypto.Sha256(trx[step].SignedID)] = *trx[step]
			}
		}
	}
}

func (self *Controller) StartBlock(when common.BlockTimeStamp, confirmBlockCount uint16, s types.BlockStatus, producerBlockId *common.BlockIdType) {
	if self.Pending != nil {
		fmt.Println("pending block already exists")
		return
	}
	//defer self.Pending.reset()
	if self.SkipDbSession(s) {
		self.Pending = types.NewPendingState(self.DB)
	} else {
		self.Pending = types.GetInstance()
	}

	self.Pending.BlockStatus = s
	self.Pending.ProducerBlockId = *producerBlockId
	self.Pending.PendingBlockState = self.Head //TODO std::make_shared<block_state>( *head, when ); // promotes pending schedule (if any) to active
	self.Pending.PendingBlockState.SignedBlock.Timestamp = when
	self.Pending.PendingBlockState.InCurrentChain = true
	self.Pending.PendingBlockState.SetConfirmed(confirmBlockCount)
	wasPendingPromoted := self.Pending.PendingBlockState.MaybePromotePending()
	log.Info("wasPendingPromoted", wasPendingPromoted)
	if self.ReadMode == DBReadMode(SPECULATIVE) || self.Pending.BlockStatus != types.BlockStatus(types.Incomplete) {
		gpo := types.GlobalPropertyObject{}
		itr, err := self.DB.Get("ID", gpo)
		if err != nil {
			log.Error("Controller StartBlock find GlobalPropertyObject is error :", err)
			return
		}
		err = itr.Data(gpo)
		//gpo.ProposedScheduleBlockNum.valid() //if there is a proposed schedule that was proposed in a block ...
		if ( /*gpo.ProposedScheduleBlockNum.valid() &&*/ gpo.ProposedScheduleBlockNum <= self.Pending.PendingBlockState.DposIrreversibleBlocknum) &&
			(len(self.Pending.PendingBlockState.PendingSchedule.Producers) == 0) &&
			(!wasPendingPromoted) {
			if !self.RePlaying {
				tmp := gpo.ProposedSchedule.ProducerScheduleType()
				ps := types.SharedProducerScheduleType{}
				ps.Version = tmp.Version
				ps.Producers = tmp.Producers
				self.Pending.PendingBlockState.SetNewProducers(ps)
			}

			self.DB.Modify(&gpo, func(i *types.GlobalPropertyObject) error {
				i.ProposedScheduleBlockNum = 1
				i.ProposedSchedule.Clear()
				return nil
			})
		}

		signedTransaction := self.GetOnBlockTransaction()
		onbtrx := types.TransactionMetadata{Trx: signedTransaction}
		onbtrx.Implicit = true
		//TODO defer
		self.InTrxRequiringChecks = true
		self.PushTransaction(onbtrx, common.MaxTimePoint(), self.GetGlobalProperties().Configuration.MinTransactionCpuUsage, true)
		fmt.Println(onbtrx)

		self.clearExpiredInputTransactions()
		self.UpdateProducersAuthority()

	}

}

func (self *Controller) PushTransaction(trx types.TransactionMetadata, deadLine common.TimePoint, billedCpuTimeUs uint32, explicitBilledCpuTime bool) (trxTrace types.TransactionTrace) {
	if deadLine == 0 {
		log.Error("deadline cannot be uninitialized")
		return
	}

	trxContext := TransactionContext{}
	trxContext = *NewTransactionContext(self, &trx.Trx, trx.ID, common.Now())

	if self.SubjectiveCupLeeway != 0 {
		if self.Pending.BlockStatus == types.BlockStatus(types.Incomplete) {
			trxContext.Leeway = self.SubjectiveCupLeeway
		}
	}
	trxContext.Deadline = deadLine
	trxContext.ExplicitBilledCpuTime = explicitBilledCpuTime
	trxContext.BilledCpuTimeUs = int64(billedCpuTimeUs)

	trace := trxContext.Trace
	if trx.Implicit {
		trxContext.InitForImplicitTrx(0) //default value 0
		trxContext.CanSubjectivelyFail = false
	} else {
		/*skipRecording := (self.replayHeadTime !=0) && (common.TimePoint(trx.Trx.Expiration) <= self.replayHeadTime)
		trxContext.InitForInputTrx(uint64(trx.PackedTrx.GetUnprunableSize()),uint64(trx.PackedTrx.GetPrunableSize()), uint32(len(trx.Trx.Signatures)),skipRecording)*/
	}

	fmt.Println(trace)

	return
}

func (self *Controller) GetGlobalProperties() (gp *types.GlobalPropertyObject) {
	gpo := types.GlobalPropertyObject{}
	itr, err := self.DB.Get("ID", gpo)
	if err != nil {
		log.Error("GetGlobalProperties is error detail:", err)
	}
	if itr.Next() {
		err = itr.Data(gp)
	}
	return gp
}

func (self *Controller) GetDynamicGlobalProperties() (r *types.DynamicGlobalPropertyObject) {
	dgpo := types.DynamicGlobalPropertyObject{}
	itr, err := self.DB.Get("ID", &dgpo)
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

func (self *Controller) GetMutableResourceLimitsManager() *ResourceLimitsManager {
	return self.ResourceLimists
}

func (self *Controller) GetOnBlockTransaction() types.SignedTransaction {
	var onBlockAction = types.Action{}
	onBlockAction.Account = common.AccountName(common.DefaultConfig.SystemAccountName)
	onBlockAction.Name = common.ActionName(common.N("onblock"))
	onBlockAction.Authorization = []types.PermissionLevel{{common.AccountName(common.DefaultConfig.SystemAccountName), common.PermissionName(common.DefaultConfig.ActiveName)}}

	data, err := rlp.EncodeToBytes(self.Head.Header)
	if err != nil {
		onBlockAction.Data = data
	}
	var trx = types.SignedTransaction{}
	trx.Actions = append(trx.Actions, &onBlockAction)
	trx.SetReferenceBlock(self.Head.ID)
	in := self.Pending.PendingBlockState.Header.Timestamp + 999999
	trx.Expiration = common.TimePointSec(in)
	log.Error("getOnBlockTransaction trx.Expiration:", trx)
	return trx
}
func (self *Controller) SkipDbSession(bs types.BlockStatus) bool {
	var considerSkipping = (bs == types.BlockStatus(IRREVERSIBLE))
	log.Info("considerSkipping:", considerSkipping)
	return considerSkipping
}

func (self *Controller) SkipDbSessions() bool {
	if self.Pending == nil {
		return self.SkipDbSession(self.Pending.BlockStatus)
	} else {
		return false
	}
}

func (self *Controller) SkipTrxChecks() (b bool) {
	b = self.LightValidationAllowed(self.Config.disableReplayOpts)
	return
}

func (self *Controller) LightValidationAllowed(dro bool) (b bool) {
	if self.Pending != nil || self.InTrxRequiringChecks {
		return false
	}

	pbStatus := self.Pending.BlockStatus
	considerSkippingOnReplay := (pbStatus == types.Irreversible || pbStatus == types.Validated) && !dro

	considerSkippingOnvalidate := (pbStatus == types.Complete && self.Config.blockValidationMode == LIGHT)

	return considerSkippingOnReplay || considerSkippingOnvalidate
}

func (self *Controller) IsProducingBlock() bool {
	if self.Pending == nil {
		return false
	}
	return self.Pending.BlockStatus == types.Incomplete
}

func (self *Controller) IsResourceGreylisted(name *common.AccountName) bool {
	for _, account := range self.Config.resourceGreylist {
		if &account == name {
			return true
		}
	}
	return false
}

func (self *Controller) PendingBlockState() *types.BlockState {
	if self.Pending != nil {
		return &self.Pending.PendingBlockState
	}
	return &types.BlockState{}
}

func (self *Controller) PendingBlockTime() common.TimePoint {
	if self.Pending == nil {
		log.Error("PendingBlockTime is error", "no pending block")
	}
	return self.Pending.PendingBlockState.Header.Timestamp.ToTimePoint()
}

func Close(db *database.LDataBase, session *database.Session) {
	//session.close()
	db.Close()
}

func (self *Controller) initConfig() *Controller {
	self.Config = Config{
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
	return self
}

func (self *Controller) GetUnAppliedTransactions() *[]types.TransactionMetadata {
	result := []types.TransactionMetadata{}
	if self.ReadMode == SPECULATIVE {
		for _, v := range self.UnAppliedTransactions {
			result = append(result, v)
		}
	} else {
		fmt.Println("not empty unapplied_transactions in non-speculative mode")
	}
	return &result
}

func (self *Controller) DropUnAppliedTransaction(metadata *types.TransactionMetadata) {
	delete(self.UnAppliedTransactions, crypto.Sha256(metadata.SignedID))
}

func (self *Controller) DropAllUnAppliedTransactions() {
	self.UnAppliedTransactions = nil
}
func (self *Controller) GetScheduledTransactions() *[]common.TransactionIdType {
	//TODO add generated_transaction_object
	//self.Db.Find("",)
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

func (self *Controller) PushScheduledTransaction(sheduled common.TransactionIdType,
	deadLine common.TimePoint,
	billedCpuTimeUs uint32) *types.TransactionTrace {

	gto := types.GetGTOByTrxId(self.DB, sheduled)

	if gto == nil {
		fmt.Println("unknown_transaction_exception", "unknown transaction")
	}

	return self.PushScheduledTransaction1(*gto, deadLine, billedCpuTimeUs)
}

func (self *Controller) PushScheduledTransaction1(gto types.GeneratedTransactionObject,
	deadLine common.TimePoint,
	billedCpuTimeUs uint32) *types.TransactionTrace {

	itr, err := self.DB.Find("ByExpiration", gto)
	if err != nil {
		fmt.Println("GetGeneratedTransactionObjectByExpiration is error :", err.Error())
	}
	defer itr.Release()

	//undo_session := self.DB.StartSession()
	gtrx := types.GeneratedTransactions(&gto)

	self.RemoveScheduledTransaction(&gto)
	if gtrx.DelayUntil <= self.PendingBlockTime() {
		fmt.Println("this transaction isn't ready")
		return nil
	}
	dtrx := types.SignedTransaction{}

	err = rlp.DecodeBytes(gtrx.PackedTrx, &dtrx)
	if err != nil {
		fmt.Println("PushScheduleTransaction1 DecodeBytes is error :", err.Error())
	}

	trx := types.TransactionMetadata{}
	fmt.Println(dtrx, trx) //TODO
	return nil
}

func (self *Controller) RemoveScheduledTransaction(gto *types.GeneratedTransactionObject) {
	self.ResourceLimists.AddPendingRamUsage(gto.Payer, int64(9)+int64(len(gto.PackedTrx))) //TODO billable_size_v
	self.DB.Remove(gto)
}

func (self *Controller) FinalizeBlock() {

}

func (self *Controller) SignBlock(callBack interface{}) {}

func (self *Controller) CommitBlock() {}

func (self *Controller) PushBlock(sbp *types.SignedBlock, status types.BlockStatus) {} //status default value block_status s = block_status::complete

func (self *Controller) PushConfirnation(hc types.HeaderConfirmation) {}

func (self *Controller) DataBase() *database.LDataBase {
	return self.DB
}

func (self *Controller) ForkDataBase() *types.ForkDatabase {
	return self.ForkDB
}

func (self *Controller) GetAccount(name common.AccountName) *types.AccountObject {
	accountObj := types.AccountObject{}
	accountObj.Name = name
	_, err := self.DB.Find("Name", accountObj)
	if err != nil {
		fmt.Println("GetAccount is error :", err)
	}
	return &accountObj
}

func (self *Controller) GetAuthorizationManager() *AuthorizationManager { return self.Authorization }

func (self *Controller) GetMutableAuthorizationManager() *AuthorizationManager {
	return self.Authorization
}

//c++ flat_set<account_name> map[common.AccountName]interface{}
func (self *Controller) GetActorWhiteList() *map[common.AccountName]struct{} {

	return &self.Config.ActorWhitelist
}

func (self *Controller) GetActorBlackList() *map[common.AccountName]struct{} {
	return &self.Config.ActorBlacklist
}

func (self *Controller) GetContractWhiteList() *map[common.AccountName]struct{} {
	return &self.Config.ContractWhitelist
}

func (self *Controller) GetContractBlackList() *map[common.AccountName]struct{} {
	return &self.Config.ContractBlacklist
}

func (self *Controller) GetActionBlockList() *map[common.Pair]struct{} {

	return &self.Config.ActionBlacklist
}

func (self *Controller) GetKeyBlackList() *map[ecc.PublicKey]struct{} {
	return &self.Config.KeyBlacklist
}

func (self *Controller) SetActorWhiteList(params *map[common.AccountName]struct{}) {
	self.Config.ActorWhitelist = *params
}

func (self *Controller) SetActorBlackList(params *map[common.AccountName]struct{}) {
	self.Config.ActorBlacklist = *params
}

func (self *Controller) SetContractWhiteList(params *map[common.AccountName]struct{}) {
	self.Config.ContractWhitelist = *params
}

func (self *Controller) SetContractBlackList(params *map[common.AccountName]struct{}) {
	self.Config.ContractBlacklist = *params
}

func (self *Controller) SetActionBlackList(params *map[common.Pair]struct{}) {
	self.Config.ActionBlacklist = *params
}

func (self *Controller) SetKeyBlackList(params *map[ecc.PublicKey]struct{}) {
	self.Config.KeyBlacklist = *params
}

func (self *Controller) HeadBlockNum() uint32 { return 0 }

func (self *Controller) HeadBlockTime() common.TimePoint { return 0 }

func (self *Controller) HeadBlockId() common.BlockIdType { return common.BlockIdType{} }

func (self *Controller) HeadBlockProducer() common.AccountName { return 0 }

func (self *Controller) HeadBlockHeader() *types.BlockHeader { return nil }

func (self *Controller) HeadBlockState() types.BlockState { return types.BlockState{} }

func (self *Controller) ForkDbHeadBlockNum() uint32 { return 0 }

func (self *Controller) ForkDbHeadBlockId() common.BlockIdType { return common.BlockIdType{} }

func (self *Controller) ForkDbHeadBlockTime() common.TimePoint { return 0 }

func (self *Controller) ForkDbHeadBlockProducer() common.AccountName { return 0 }

func (self *Controller) ActiveProducers() *types.ProducerScheduleType { return nil }

func (self *Controller) PendingProducers() *types.ProducerScheduleType { return nil }

func (self *Controller) ProposedProducers() types.ProducerScheduleType {
	return types.ProducerScheduleType{}
}

func (self *Controller) LastIrreversibleBlockNum() uint32 { return 0 }

func (self *Controller) LastIrreversibleBlockId() common.BlockIdType { return common.BlockIdType{} }

func (self *Controller) FetchBlockByNumber(blockNum uint32) *types.SignedBlock {
	blkState := self.ForkDB.GetBlockInCurrentChainByNum(blockNum)
	if blkState != nil {
		return blkState.SignedBlock
	}
	//TODO blog
	return &types.SignedBlock{}
}

func (self *Controller) FetchBlockById(id common.BlockIdType) *types.SignedBlock {
	state := self.ForkDB.GetBlock(&id)
	if state != nil {
		return state.SignedBlock
	}
	bptr := self.FetchBlockByNumber(types.NumFromID(id))
	if bptr != nil && bptr.BlockID() == id {
		return bptr
	}
	return &types.SignedBlock{}
}

func (self *Controller) FetchBlockStateByNumber(blockNum uint32) *types.BlockState {
	return self.ForkDB.GetBlockInCurrentChainByNum(blockNum)
}

func (self *Controller) FetchBlockStateById(id common.BlockIdType) *types.BlockState {
	return self.ForkDB.GetBlock(&id)
}

func (self *Controller) GetBlcokIdForNum(blockNum uint32) common.BlockIdType {
	return common.BlockIdType{}
}

func (self *Controller) CheckContractList(code common.AccountName) {
	if len(self.Config.ContractWhitelist) > 0 {
		if _, ok := self.Config.ContractWhitelist[code]; !ok {
			fmt.Println("account is not on the contract whitelist", code)
			return
		}
		/*EOS_ASSERT( conf.contract_whitelist.find( code ) != conf.contract_whitelist.end(),
			contract_whitelist_exception,
			"account '${code}' is not on the contract whitelist", ("code", code)
		);*/
	} else if len(self.Config.ContractBlacklist) > 0 {
		if _, ok := self.Config.ContractBlacklist[code]; ok {
			fmt.Println("account is on the contract blacklist", code)
			return
		}
		/*EOS_ASSERT( conf.contract_blacklist.find( code ) == conf.contract_blacklist.end(),
			contract_blacklist_exception,
			"account '${code}' is on the contract blacklist", ("code", code)
		);*/
	}
}

func (self *Controller) CheckActionList(code common.AccountName, action common.ActionName) {
	if len(self.Config.ActionBlacklist) > 0 {
		abl := common.MakePair(code, action)
		//key := Hash(abl)
		if _, ok := self.Config.ActionBlacklist[abl]; ok {
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

func (self *Controller) CheckKeyList(key *ecc.PublicKey) {
	if len(self.Config.KeyBlacklist) > 0 {
		if _, ok := self.Config.KeyBlacklist[*key]; ok {
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

func (self *Controller) IsProducing() bool { return false }

func (self *Controller) IsRamBillingInNotifyAllowed() bool { return false }

func (self *Controller) AddResourceGreyList(name *common.AccountName) {}

func (self *Controller) RemoveResourceGreyList(name *common.AccountName) {}

func (self *Controller) IsResourceGreyListed(name *common.AccountName) bool { return false }

func (self *Controller) GetResourceGreyList() *map[common.AccountName]interface{} { return nil }

func (self *Controller) ValidateReferencedAccounts(t *types.Transaction) {}

func (self *Controller) ValidateExpiration(t *types.Transaction) {}

func (self *Controller) ValidateTapos(t *types.Transaction) {}

func (self *Controller) ValidateDbAvailableSize() {}

func (self *Controller) ValidateReversibleAvailableSize() {}

func (self *Controller) IsKnownUnexpiredTransaction(id *common.TransactionIdType) bool { return false }

func (self *Controller) SetProposedProducers(producers []types.ProducerKey) int64 { return 0 }

func (self *Controller) SkipAuthCheck() bool { return false }

func (self *Controller) ContractsConsole() bool { return false }

func (self *Controller) GetChainId() common.ChainIdType { return common.ChainIdType{} }

func (self *Controller) GetReadMode() DBReadMode { return 0 }

func (self *Controller) GetValidationMode() ValidationMode { return 0 }

func (self *Controller) SetSubjectiveCpuLeeway(leeway common.Microseconds) {}

func (self *Controller) PendingProducerBlockId() common.BlockIdType {
	//EOS_ASSERT( my->pending, block_validate_exception, "no pending block" )
	return self.Pending.ProducerBlockId
}

func (self *Controller) FindApplyHandler(receiver common.AccountName,
	scope common.AccountName,
	act common.ActionName) func(*ApplyContext) {

	handlerKey := NewHandlerKey(common.ScopeName(scope), act)
	secend, ok := self.ApplyHandlers[receiver]
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

func (self *Controller) GetWasmInterface() *exec.WasmInterface {
	return self.WasmIf
}

func (self *Controller) GetAbiSerializer(name common.AccountName,
	maxSerializationTime common.Microseconds) types.AbiSerializer {
	return types.AbiSerializer{}
}

func (self *Controller) ToVariantWithAbi(obj interface{}, maxSerializationTime common.Microseconds) {}

var readycontroller chan bool //TODO test code

func (self *Controller) CreateNativeAccount(name common.AccountName, owner types.Authority, active types.Authority, isPrivileged bool) {
	account := types.AccountObject{}
	account.Name = name
	account.CreationDate = common.BlockTimeStamp(self.Config.genesis.InitialTimestamp)
	account.Privileged = isPrivileged
	if name == common.AccountName(common.DefaultConfig.SystemAccountName) {
		abiDef := types.AbiDef{}
		account.SetAbi(EosioContractAbi(abiDef))
	}
	self.DB.Insert(account)

	aso := types.AccountSequenceObject{}
	aso.Name = name
	self.DB.Insert(aso)

	ownerPermission := self.Authorization.CreatePermission(name, common.PermissionName(common.DefaultConfig.OwnerName), 0, owner, self.Config.genesis.InitialTimestamp)

	activePermission := self.Authorization.CreatePermission(name, common.PermissionName(common.DefaultConfig.ActiveName), PermissionIdType(ownerPermission.ID), active, self.Config.genesis.InitialTimestamp)

	self.ResourceLimists.InitializeAccount(name)
	ramDelta := uint64(common.DefaultConfig.OverheadPerRowPerIndexRamBytes) //TODO c++ reference int64 but statement uint32
	ramDelta += 2 * common.BillableSizeV("permission_object")               //::billable_size_v<permission_object>
	ramDelta += ownerPermission.Auth.GetBillableSize()
	ramDelta += activePermission.Auth.GetBillableSize()
	self.ResourceLimists.AddPendingRamUsage(name, int64(ramDelta))
	self.ResourceLimists.VerifyAccountRamUsage(name)
}

func initResource(self *Controller, ready chan bool) {
	<-ready
	//con.Blog
	self.ForkDB = types.GetForkDbInstance(common.DefaultConfig.DefaultStateDirName)
	self.ResourceLimists = GetResourceLimitsManager()
	self.Authorization = GetAuthorizationManager()
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

/*func main(){
	c := new(Controller)

	fmt.Println("asdf",c)
}*/

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

func (self *Controller) clearExpiredInputTransactions() {
	/*
		auto& transaction_idx = db.get_mutable_index<transaction_multi_index>();
		const auto& dedupe_index = transaction_idx.indices().get<by_expiration>();
		auto now = self.pending_block_time();
		while( (!dedupe_index.empty()) && ( now > fc::time_point(dedupe_index.begin()->expiration) ) ) {
		transaction_idx.remove(*dedupe_index.begin());
		}
	*/
}

func (self *Controller) UpdateProducersAuthority() {
	/*producers := self.Pending.PendingBlockState.ActiveSchedule.Producers
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
