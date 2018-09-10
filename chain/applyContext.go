package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/db"
	types2 "go/types"
	"math/big"
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

}


func (applyContext *ApplyContext) execOne() (trace types.ActionTrace) {

	return
}

func (applyContext *ApplyContext) Exec(){

}

func (applyContext *ApplyContext) IsAccount(account common.AccountName) bool{

	return false
}

func (applyContext *ApplyContext) RequireAuthorization(account common.AccountName){

}

func (applyContext *ApplyContext) HasAuthorization(account common.AccountName) bool{

	return false
}

func (applyContext *ApplyContext) RequireAuthorizations(account common.AccountName){


}

func (applyContext *ApplyContext) HasReciptient(code common.AccountName) bool{

	return false
}

func (applyContext *ApplyContext) RequireRecipient(recipient common.AccountName){

}

func (applyContext *ApplyContext) ExecuteInline(a *common.Action){

}

func (applyContext *ApplyContext) ExecuteContextFreeInline(a *common.Action){

}

func (applyContext *ApplyContext) ScheduleDeferredTransaction(sendId *big.Int,payer common.AccountName,trx types.Transaction,replaceExisting bool){

}

func (applyContext *ApplyContext) CancelDeferredTransaction(sendId *big.Int,sender common.AccountName) bool{

	return false
}

func (applyContext *ApplyContext) FindTable(code common.Name,scope common.Name,table common.Name) (tid *types.TableIDObject){

	return
}

func (applyContext *ApplyContext) FindOrCreateTable(code common.Name,scope common.Name,table common.Name,payer *common.AccountName) (tid *types.TableIDObject){

	return
}

func (applyContext *ApplyContext) RemoveTable(tid types.TableIDObject){

}

func (applyContext *ApplyContext) GetActiveProducers() (accounts []common.AccountName){

	return
}

func (applyContext *ApplyContext) ResetConsole(){

}

func (applyContext *ApplyContext) GetPackedTransaction() (bytes []byte){

	return
}

func (applyContext *ApplyContext) UpdateDBUsage(payer common.AccountName,delta int64){

}

func (applyContext *ApplyContext) GetAction(typ uint32,index uint32,buffer []byte,bufferSize types2.Sizes) int{

	return 0
}

func (applyContext *ApplyContext) GetContextFreeData(index uint32,buffer *[]byte,bufferSize types2.Sizes) int {

	return 0
}

func (applyContext *ApplyContext) DBStoreI64() int{

	return 0
}

func (applyContext *ApplyContext) DBUpdateI64(iterator int,payer common.AccountName,buffer []byte,bufferSize int){//TODO bufferSize int

}

func (applyContext *ApplyContext) DBRemoveI64(iterator int){

}

func (applyContext *ApplyContext) DBGetI64(iterator int,buffer *[]byte,bufferSize int) int {//TODO bufferSize int

	return 0
}

func (applyContext *ApplyContext) DBNextI64(iterator int,primary uint64) int{

	return 0
}

func (applyContext *ApplyContext) DBPreviousI64(iterator int,primary uint64) int{

	return 0
}

func (applyContext *ApplyContext) DBFindI64(iterator int,primary uint64) int{

	return 0
}

func (applyContext *ApplyContext) DBLowerboundI64(iterator int,primary uint64) int{

	return 0
}