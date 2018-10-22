package wasmgo

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	arithmetic "github.com/eosspark/eos-go/common/arithmetic_types"
)

type EnvContext interface {

	//action
	GetActionData() []byte
	GetReceiver() common.AccountName
	GetCode() common.AccountName
	GetAct() common.ActionName

	//authorization
	RequireAuthorization(account int64)
	HasAuthorization(account int64) bool
	RequireAuthorization2(account int64, permission int64)
	//RequireAuthorizations(account common.AccountName)
	RequireRecipient(recipient int64)
	IsAccount(n int64) bool
	HasReciptient(code int64) bool

	//console
	//ResetConsole()
	ContextAppend(str string)

	//database
	//primaryKey
	DbStoreI64(scope int64, table int64, payer int64, id int64, buffer []byte) int
	DbUpdateI64(iterator int, payer int64, buffer []byte)
	DbRemoveI64(iterator int)
	DbGetI64(iterator int, buffer []byte, bufferSize int) int
	DbNextI64(iterator int, primary *uint64) int
	DbPreviousI64(iterator int, primary *uint64) int
	DbFindI64(code int64, scope int64, table int64, id int64) int
	DbLowerboundI64(code int64, scope int64, table int64, id int64) int
	DbUpperboundI64(code int64, scope int64, table int64, id int64) int
	DbEndI64(code int64, scope int64, table int64) int

	//secondaryKey 64
	Idx64Store(scope int64, table int64, payer int64, id int64, value *types.Uint64_t) int
	Idx64Remove(iterator int)
	Idx64Update(iterator int, payer int64, value *types.Uint64_t)
	Idx64FindSecondary(code int64, scope int64, table int64, secondary *types.Uint64_t, primary *uint64) int
	Idx64Lowerbound(code int64, scope int64, table int64, secondary *types.Uint64_t, primary *uint64) int
	Idx64Upperbound(code int64, scope int64, table int64, secondary *types.Uint64_t, primary *uint64) int
	Idx64End(code int64, scope int64, table int64) int
	Idx64Next(iterator int, primary *uint64) int
	Idx64Previous(iterator int, primary *uint64) int
	Idx64FindPrimary(code int64, scope int64, table int64, secondary *types.Uint64_t, primary *uint64) int

	//secondaryKey Double
	IdxDoubleStore(scope int64, table int64, payer int64, id int64, value *types.Float64_t) int
	IdxDoubleRemove(iterator int)
	IdxDoubleUpdate(iterator int, payer int64, value *types.Float64_t)
	IdxDoubleFindSecondary(code int64, scope int64, table int64, secondary *types.Float64_t, primary *uint64) int
	IdxDoubleLowerbound(code int64, scope int64, table int64, secondary *types.Float64_t, primary *uint64) int
	IdxDoubleUpperbound(code int64, scope int64, table int64, secondary *types.Float64_t, primary *uint64) int
	IdxDoubleEnd(code int64, scope int64, table int64) int
	IdxDoubleNext(iterator int, primary *uint64) int
	IdxDoublePrevious(iterator int, primary *uint64) int
	IdxDoubleFindPrimary(code int64, scope int64, table int64, secondary *types.Float64_t, primary *uint64) int

	//permission
	GetPermissionLastUsed(account common.AccountName, permission common.PermissionName) int64
	GetAccountCreateTime(account common.AccountName) int64

	//privileged
	SetResourceLimits(account common.AccountName, ramBytes uint64, netWeight uint64, cpuWeigth uint64)
	GetResourceLimits(account common.AccountName, ramBytes *uint64, netWeight *uint64, cpuWeigth *uint64)
	SetBlockchainParametersPacked(parameters []byte)
	GetBlockchainParametersPacked() []byte
	IsPrivileged(n common.AccountName) bool
	SetPrivileged(n common.AccountName, isPriv bool)

	//producer
	SetProposedProducers(producers []byte)
	GetActiveProducersInBytes() []byte
	//GetActiveProducers() []common.AccountName

	//system
	CheckTime()
	CurrentTime() int64
	PublicationTime() int64

	//transaction
	ExecuteInline(action []byte)
	ExecuteContextFreeInline(action []byte)
	ScheduleDeferredTransaction(sendId *arithmetic.Uint128, payer common.AccountName, trx []byte, replaceExisting bool)
	CancelDeferredTransaction(sendId *arithmetic.Uint128) bool
	GetPackedTransaction() []byte
	Expiration() int
	TaposBlockNum() int
	TaposBlockPrefix() int
	GetAction(typ uint32, index int, bufferSize int) (int, []byte)
	GetContextFreeData(intdex int, bufferSize int) (int, []byte)
}
