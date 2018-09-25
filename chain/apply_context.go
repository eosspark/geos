package chain

import (
	"github.com/eosspark/eos-go/base"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/rlp"
)

type ApplyContext struct {
	Controller *Controller

	DB                 *eosiodb.DataBase
	TrxContext         *TransactionContext
	Act                types.Action
	Receiver           common.AccountName
	UsedAuthorizations []bool
	RecurseDepth       uint32
	Privileged         bool
	ContextFree        bool
	UsedContestFreeApi bool
	Trace              types.ActionTrace

	//GenericIndex
	//_pending_console_output
	PendingConsoleOutput string
	Notified             []common.AccountName
}

func (a *ApplyContext) execOne() (trace types.ActionTrace) { return }
func (a *ApplyContext) Exec()                              {}

//context action api
func (a *ApplyContext) GetActionData() []byte           { return a.Act.Data }
func (a *ApplyContext) GetReceiver() common.AccountName { return a.Receiver }
func (a *ApplyContext) GetCode() common.AccountName     { return a.Act.Account }
func (a *ApplyContext) GetAct() common.ActionName       { return a.Act.Name }

//context authorization api
func (a *ApplyContext) RequireAuthorization(account common.AccountName) {
	for k, v := range a.Act.Authorization {
		if v.Actor == account {
			a.UsedAuthorizations[k] = true
			return
		}
	}
	// EOS_ASSERT( false, missing_auth_exception, "missing authority of ${account}/${permission}",
	//              ("account",account)("permission",permission) );
}
func (a *ApplyContext) HasAuthorization(account common.AccountName) bool {
	for _, v := range a.Act.Authorization {
		if v.Actor == account {
			return true
		}
	}
	return false
}
func (a *ApplyContext) RequireAuthorization2(account common.AccountName, permission common.PermissionName) {
	for k, v := range a.Act.Authorization {
		if v.Actor == account && v.Permission == permission {
			a.UsedAuthorizations[k] = true
			return
		}
	}

	//EOS_ASSERT( false, missing_auth_exception, "missing authority of ${account}", ("account",account));
}

//func (a *ApplyContext) RequireAuthorizations(account common.AccountName) {}
func (a *ApplyContext) RequireRecipient(recipient common.AccountName) {
	if a.HasReciptient(recipient) {
		a.Notified = append(a.Notified, recipient)
	}
}
func (a *ApplyContext) IsAccount(account common.AccountName) bool {
	return false
	//return nullptr != db.find<account_object,by_name>( account );
}
func (a *ApplyContext) HasReciptient(code common.AccountName) bool {
	for _, a := range a.Notified {
		if a == code {
			return true
		}
	}
	return false
}

//context console api
func (a *ApplyContext) ResetConsole()            { return }
func (a *ApplyContext) ContextAppend(str string) { a.PendingConsoleOutput += str }

//context database api
func (a *ApplyContext) DBStoreI64() int { return 0 }
func (a *ApplyContext) DBUpdateI64(
	iterator int,
	payer common.AccountName,
	buffer []byte,
	bufferSize int) {
}
func (a *ApplyContext) DBRemoveI64(iterator int)                                  {}
func (a *ApplyContext) DBGetI64(iterator int, buffer *[]byte, bufferSize int) int { return 0 }
func (a *ApplyContext) DBNextI64(iterator int, primary uint64) int                { return 0 }
func (a *ApplyContext) DBPreviousI64(iterator int, primary uint64) int            { return 0 }
func (a *ApplyContext) DBFindI64(iterator int, primary uint64) int                { return 0 }
func (a *ApplyContext) DBLowerboundI64(iterator int, primary uint64) int          { return 0 }
func (a *ApplyContext) UpdateDBUsage(payer common.AccountName, delta int64)       {}
func (a *ApplyContext) FindTable(
	code common.Name,
	scope common.Name,
	table common.Name) types.TableIDObject {
	return types.TableIDObject{}
}
func (a *ApplyContext) FindOrCreateTable(code common.Name,
	scope common.Name,
	table common.Name,
	payer *common.AccountName) types.TableIDObject {
	return types.TableIDObject{}
}
func (a *ApplyContext) RemoveTable(tid types.TableIDObject) {}

//context permission api
func (a *ApplyContext) GetPermissionLastUsed(account common.AccountName, permission common.PermissionName) int64 {
	//return 0
	//am := a.Controller.G
	return 0
}
func (a *ApplyContext) GetAccountCreateTime(account common.AccountName) int64 { return 0 }

//context privileged api
func (a *ApplyContext) SetResourceLimits(
	account common.AccountName,
	ramBytes uint64,
	netWeight uint64,
	cpuWeigth uint64) {

}
func (a *ApplyContext) GetResourceLimits(
	account common.AccountName,
	ramBytes *uint64,
	netWeight *uint64,
	cpuWeigth *uint64) {
}
func (a *ApplyContext) SetBlockchainParametersPacked(parameters []byte) {

	newGPO := types.GlobalPropertyObject{}
	rlp.DecodeBytes(parameters, &newGPO)
	oldGPO := a.Controller.GetGlobalProperties()
	a.Controller.db.UpdateObject(&oldGPO, &newGPO)

}

func (a *ApplyContext) GetBlockchainParametersPacked() []byte {
	gpo := a.Controller.GetGlobalProperties()
	bytes, err := rlp.EncodeToBytes(gpo)
	if err != nil {
		log.Error("EncodeToBytes is error detail:", err)
		return nil
	}
	return bytes
}
func (a *ApplyContext) IsPrivileged(n common.AccountName) bool {
	//return false
	account := types.AccountObject{Name: n}

	err := a.Controller.db.ByIndex("byName", &account)
	if err != nil {
		log.Error("getaAccount is error detail:", err)
		return false
	}
	return account.Privileged

}
func (a *ApplyContext) SetPrivileged(n common.AccountName, isPriv bool) {
	oldAccount := types.AccountObject{Name: n}
	err := a.Controller.db.ByIndex("byName", &oldAccount)
	if err != nil {
		log.Error("getaAccount is error detail:", err)
		return
	}

	newAccount := oldAccount
	newAccount.Privileged = isPriv
	a.Controller.db.UpdateObject(&oldAccount, &newAccount)
}

//context producer api
func (a *ApplyContext) SetProposedProducers(producers []byte) { return }
func (a *ApplyContext) GetActiveProducersInBytes() []byte {
	b := make([]byte, 256)
	return b
}

//func (a *ApplyContext) GetActiveProducers() []common.AccountName { return }

//context system api
func (a *ApplyContext) CheckTime() { return }
func (a *ApplyContext) CurrentTime() int64 {

	//return a.Controller.PendingBlockTime().TimeSinceEpoch().Count()
	return 0
}
func (a *ApplyContext) PublicationTime() int64 { return 0 }

//context transaction api
func (a *ApplyContext) ExecuteInline(action []byte)            {}
func (a *ApplyContext) ExecuteContextFreeInline(action []byte) {}
func (a *ApplyContext) ScheduleDeferredTransaction(sendId common.TransactionIDType, payer common.AccountName, trx []byte, replaceExisting bool) {
}
func (a *ApplyContext) CancelDeferredTransaction(sendId common.TransactionIDType) bool { return false }
func (a *ApplyContext) GetPackedTransaction() []byte                                   { return []byte{} }
func (a *ApplyContext) Expiration() int                                                { return 0 }
func (a *ApplyContext) TaposBlockNum() int                                             { return 0 }
func (a *ApplyContext) TaposBlockPrefix() int                                          { return 0 }
func (a *ApplyContext) GetAction(typ uint32, index int, bufferSize int) (int, []byte) {
	trx := a.TrxContext.Trx
	var a_ptr *types.Action
	if typ == 0 {
		if index >= len(trx.ContextFreeActions) {
			return -1, nil
		}
		a_ptr = trx.ContextFreeActions[index]
	} else if typ == 1 {
		if index >= len(trx.ContextFreeActions) {
			return -1, nil
		}
		a_ptr = trx.Actions[index]
	}
	if a_ptr == nil {
		return -1, nil
	}

	s, _ := rlp.EncodeSize(*a_ptr)
	if s <= bufferSize {
		bytes, _ := rlp.EncodeToBytes(*a_ptr)
		return s, bytes
	}
	return s, nil

}
func (a *ApplyContext) GetContextFreeData(index int, bufferSize int) (int, []byte) {

	trx := a.TrxContext.Trx
	if index >= len(trx.ContextFreeData) {
		return -1, nil
	}
	s := len(trx.ContextFreeData[index])
	if bufferSize == 0 {
		return s, nil
	}
	copySize := base.Min(uint64(bufferSize), uint64(s))
	return int(copySize), trx.ContextFreeData[index][0:copySize]

}
