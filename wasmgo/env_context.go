package wasmgo

import (
	"github.com/eosspark/container/sets/treeset"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/eos_math"
)

type EnvContext interface {

	//action
	GetActionData() []byte
	GetReceiver() common.AccountName
	GetCode() common.AccountName
	GetAct() common.ActionName
	ContextFreeAction() bool

	//authorization
	RequireAuthorization(account int64)
	HasAuthorization(account int64) bool
	RequireAuthorization2(account int64, permission int64)
	RequireRecipient(recipient int64)
	IsAccount(n int64) bool
	HasRecipient(code int64) bool

	//console
	//ResetConsole()
	ContextAppend(str string)

	//database
	//primaryKey
	DbStoreI64(scope uint64, table uint64, payer uint64, id uint64, buffer []byte) int
	DbUpdateI64(iterator int, payer uint64, buffer []byte)
	DbRemoveI64(iterator int)
	DbGetI64(iterator int, buffer []byte, bufferSize int) int
	DbNextI64(iterator int, primary *uint64) int
	DbPreviousI64(iterator int, primary *uint64) int
	DbFindI64(code uint64, scope uint64, table uint64, id uint64) int
	DbLowerboundI64(code uint64, scope uint64, table uint64, id uint64) int
	DbUpperboundI64(code uint64, scope uint64, table uint64, id uint64) int
	DbEndI64(code uint64, scope uint64, table uint64) int

	//index64 secondaryKey
	Idx64Store(scope uint64, table uint64, payer uint64, id uint64, value *uint64) int
	Idx64Remove(iterator int)
	Idx64Update(iterator int, payer uint64, value *uint64)
	Idx64FindSecondary(code uint64, scope uint64, table uint64, secondary *uint64, primary *uint64) int
	Idx64Lowerbound(code uint64, scope uint64, table uint64, secondary *uint64, primary *uint64) int
	Idx64Upperbound(code uint64, scope uint64, table uint64, secondary *uint64, primary *uint64) int
	Idx64End(code uint64, scope uint64, table uint64) int
	Idx64Next(iterator int, primary *uint64) int
	Idx64Previous(iterator int, primary *uint64) int
	Idx64FindPrimary(code uint64, scope uint64, table uint64, secondary *uint64, primary uint64) int

	//index128 secondaryKey
	Idx128Store(scope uint64, table uint64, payer uint64, id uint64, value *eos_math.Uint128) int
	Idx128Remove(iterator int)
	Idx128Update(iterator int, payer uint64, value *eos_math.Uint128)
	Idx128FindSecondary(code uint64, scope uint64, table uint64, secondary *eos_math.Uint128, primary *uint64) int
	Idx128Lowerbound(code uint64, scope uint64, table uint64, secondary *eos_math.Uint128, primary *uint64) int
	Idx128Upperbound(code uint64, scope uint64, table uint64, secondary *eos_math.Uint128, primary *uint64) int
	Idx128End(code uint64, scope uint64, table uint64) int
	Idx128Next(iterator int, primary *uint64) int
	Idx128Previous(iterator int, primary *uint64) int
	Idx128FindPrimary(code uint64, scope uint64, table uint64, secondary *eos_math.Uint128, primary uint64) int

	//index256 secondaryKey
	Idx256Store(scope uint64, table uint64, payer uint64, id uint64, value *eos_math.Uint256) int
	Idx256Remove(iterator int)
	Idx256Update(iterator int, payer uint64, value *eos_math.Uint256)
	Idx256FindSecondary(code uint64, scope uint64, table uint64, secondary *eos_math.Uint256, primary *uint64) int
	Idx256Lowerbound(code uint64, scope uint64, table uint64, secondary *eos_math.Uint256, primary *uint64) int
	Idx256Upperbound(code uint64, scope uint64, table uint64, secondary *eos_math.Uint256, primary *uint64) int
	Idx256End(code uint64, scope uint64, table uint64) int
	Idx256Next(iterator int, primary *uint64) int
	Idx256Previous(iterator int, primary *uint64) int
	Idx256FindPrimary(code uint64, scope uint64, table uint64, secondary *eos_math.Uint256, primary uint64) int

	//index Double secondaryKey
	IdxDoubleStore(scope uint64, table uint64, payer uint64, id uint64, value *eos_math.Float64) int
	IdxDoubleRemove(iterator int)
	IdxDoubleUpdate(iterator int, payer uint64, value *eos_math.Float64)
	IdxDoubleFindSecondary(code uint64, scope uint64, table uint64, secondary *eos_math.Float64, primary *uint64) int
	IdxDoubleLowerbound(code uint64, scope uint64, table uint64, secondary *eos_math.Float64, primary *uint64) int
	IdxDoubleUpperbound(code uint64, scope uint64, table uint64, secondary *eos_math.Float64, primary *uint64) int
	IdxDoubleEnd(code uint64, scope uint64, table uint64) int
	IdxDoubleNext(iterator int, primary *uint64) int
	IdxDoublePrevious(iterator int, primary *uint64) int
	IdxDoubleFindPrimary(code uint64, scope uint64, table uint64, secondary *eos_math.Float64, primary uint64) int

	//index LongDouble secondaryKey
	IdxLongDoubleStore(scope uint64, table uint64, payer uint64, id uint64, value *eos_math.Float128) int
	IdxLongDoubleRemove(iterator int)
	IdxLongDoubleUpdate(iterator int, payer uint64, value *eos_math.Float128)
	IdxLongDoubleFindSecondary(code uint64, scope uint64, table uint64, secondary *eos_math.Float128, primary *uint64) int
	IdxLongDoubleLowerbound(code uint64, scope uint64, table uint64, secondary *eos_math.Float128, primary *uint64) int
	IdxLongDoubleUpperbound(code uint64, scope uint64, table uint64, secondary *eos_math.Float128, primary *uint64) int
	IdxLongDoubleEnd(code uint64, scope uint64, table uint64) int
	IdxLongDoubleNext(iterator int, primary *uint64) int
	IdxLongDoublePrevious(iterator int, primary *uint64) int
	IdxLongDoubleFindPrimary(code uint64, scope uint64, table uint64, secondary *eos_math.Float128, primary uint64) int

	//permission
	GetPermissionLastUsed(account common.AccountName, permission common.PermissionName) common.TimePoint
	GetAccountCreateTime(account common.AccountName) common.TimePoint

	//privileged
	// SetResourceLimits(account common.AccountName, ramBytes uint64, netWeight uint64, cpuWeigth uint64) bool
	// GetResourceLimits(account common.AccountName, ramBytes *uint64, netWeight *uint64, cpuWeigth *uint64)

	SetAccountLimits(account common.AccountName, ramBytes int64, netWeight int64, cpuWeigth int64) bool
	GetAccountLimits(account common.AccountName, ramBytes *int64, netWeight *int64, cpuWeigth *int64)
	SetBlockchainParametersPacked(parameters []byte)
	GetBlockchainParametersPacked() []byte
	GetBlockchainParameters() *types.ChainConfig
	SetBlockchainParameters(cfg *types.ChainConfig)

	IsPrivileged(n common.AccountName) bool
	SetPrivileged(n common.AccountName, isPriv bool)
	ValidateRamUsageInsert(account common.AccountName)

	//producer
	SetProposedProducers(producers []byte) int64
	GetActiveProducersInBytes() []byte
	//GetActiveProducers() []common.AccountName

	//system
	CheckTime()
	//CurrentTime() int64
	CurrentTime() common.TimePoint
	//PublicationTime() int64
	PublicationTime() common.TimePoint

	//transaction
	// ExecuteInline(action []byte)
	// ExecuteContextFreeInline(action []byte)
	InlineActionTooBig(dataLen int) bool
	ExecuteInline(act *types.Action)
	ExecuteContextFreeInline(act *types.Action)
	//ScheduleDeferredTransaction(sendId *eos_math.Uint128, payer common.AccountName, trx []byte, replaceExisting bool)
	ScheduleDeferredTransaction(sendId *eos_math.Uint128, payer common.AccountName, trx *types.Transaction, replaceExisting bool)
	CancelDeferredTransaction(sendId *eos_math.Uint128) bool
	//GetPackedTransaction() []byte
	//GetPackedTransaction() *types.SignedTransaction
	GetPackedTransaction() *types.Transaction
	//Expiration() int
	Expiration() common.TimePointSec
	TaposBlockNum() int
	TaposBlockPrefix() int
	//GetAction(typ uint32, index int, bufferSize int) (int, []byte)
	GetAction(typ uint32, index int) *types.Action
	GetContextFreeData(intdex int, bufferSize int) (int, []byte)

	PauseBillingTimer()
	ResumeBillingTimer()

	CheckAuthorization(actions []*types.Action, providedKeys *treeset.Set, providedPermissions *treeset.Set, delayUS uint64)
	CheckAuthorization2(n common.AccountName, permission common.PermissionName, providedKeys *treeset.Set, providedPermissions *treeset.Set, delayUS uint64)
	//GetLogger() *log.Logger
}
