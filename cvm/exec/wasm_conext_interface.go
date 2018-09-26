package exec

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type WasmContextInterface interface {

	//action api
	GetActionData() []byte
	GetReceiver() common.AccountName
	GetCode() common.AccountName
	GetAct() common.ActionName

	//context authorization api
	RequireAuthorization(account common.AccountName)
	HasAuthorization(account common.AccountName) bool
	RequireAuthorization2(account common.AccountName, permission common.PermissionName)
	//RequireAuthorizations(account common.AccountName)
	RequireRecipient(recipient common.AccountName)
	IsAccount(n common.AccountName) bool
	HasReciptient(code common.AccountName) bool

	//contet console text
	ResetConsole()
	ContextAppend(str string)

	//context database api
	DBStoreI64(scope int64, table int64, payer int64, id int64, buffer int, buffer_size int) int
	DBUpdateI64(iterator int, payer common.AccountName, buffer []byte, bufferSize int)
	DBRemoveI64(iterator int)
	DBGetI64(iterator int, buffer *[]byte, bufferSize int) int
	DBNextI64(iterator int, primary uint64) int
	DBPreviousI64(iterator int, primary uint64) int
	DBFindI64(iterator int, primary uint64) int
	DBLowerboundI64(iterator int, primary uint64) int
	UpdateDBUsage(pager common.AccountName, delta int64)
	FindTable(code common.Name, scope common.Name, table common.Name) types.TableIDObject
	//FindOrCreateTable(code common.Name, scope common.Name, table common.Name, payer *common.AccountName) types.TableIDObject
	RemoveTable(tid types.TableIDObject)

	//context permission api
	GetPermissionLastUsed(account common.AccountName, permission common.PermissionName) int64
	GetAccountCreateTime(account common.AccountName) int64

	//context privileged api
	SetResourceLimits(account common.AccountName, ramBytes uint64, netWeight uint64, cpuWeigth uint64)
	GetResourceLimits(account common.AccountName, ramBytes *uint64, netWeight *uint64, cpuWeigth *uint64)
	SetBlockchainParametersPacked(parameters []byte)
	GetBlockchainParametersPacked() []byte
	IsPrivileged(n common.AccountName) bool
	SetPrivileged(n common.AccountName, isPriv bool)

	//context producer api
	SetProposedProducers(producers []byte)
	GetActiveProducersInBytes() []byte
	//GetActiveProducers() []common.AccountName

	//context system api
	CheckTime()
	CurrentTime() int64
	PublicationTime() int64

	//context transaction api
	ExecuteInline(action []byte)
	ExecuteContextFreeInline(action []byte)
	ScheduleDeferredTransaction(sendId common.TransactionIDType, payer common.AccountName, trx []byte, replaceExisting bool)
	CancelDeferredTransaction(sendId common.TransactionIDType) bool
	GetPackedTransaction() []byte
	Expiration() int
	TaposBlockNum() int
	TaposBlockPrefix() int
	GetAction(typ uint32, index int, bufferSize int) (int, []byte)
	GetContextFreeData(intdex int, bufferSize int) (int, []byte)
}