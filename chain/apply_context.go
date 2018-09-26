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
	Control *Controller

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
	KeyvalCache          iteratorCache
	Notified             []common.AccountName
	PendingConsoleOutput string
}

type pairTableIterator struct {
	tableIDObject *types.TableIDObject
	iterator      int
}
type iteratorCache struct {
	tableCache         map[types.IdType]*pairTableIterator
	endIteratorToTable []*types.TableIDObject
	iteratorToObject   []interface{}
	objectToIterator   map[interface{}]int
}

func NewIteratorCache() *iteratorCache {

	i := &iteratorCache{
		tableCache:         make(map[types.IdType]*pairTableIterator),
		endIteratorToTable: make([]*types.TableIDObject, 8),
		iteratorToObject:   make([]interface{}, 32),
		objectToIterator:   make(map[interface{}]int),
	}

	return i
}

func (i *iteratorCache) endIteratorToIndex(ei int) int   { return (-ei - 2) }
func (i *iteratorCache) IndexToEndIterator(indx int) int { return -(indx + 2) }
func (i *iteratorCache) cacheTable(tobj *types.TableIDObject) int {
	if itr, ok := i.tableCache[tobj.ID]; ok {
		return itr.iterator
	}
	ei := i.IndexToEndIterator(len(i.endIteratorToTable))
	i.endIteratorToTable = append(i.endIteratorToTable, tobj)

	pair := &pairTableIterator{
		tableIDObject: tobj,
		iterator:      ei,
	}
	i.tableCache[tobj.ID] = pair
	return ei
}
func (i *iteratorCache) getTable(id types.IdType) *types.TableIDObject {
	if itr, ok := i.tableCache[id]; ok {
		return itr.tableIDObject
	}

	return &types.TableIDObject{}
	//EOS_ASSERT( itr != _table_cache.end(), table_not_in_cache, "an invariant was broken, table should be in cache" );
}
func (i *iteratorCache) getEndIteratorByTableID(id types.IdType) int {
	if itr, ok := i.tableCache[id]; ok {
		return itr.iterator
	}
	//EOS_ASSERT( itr != _table_cache.end(), table_not_in_cache, "an invariant was broken, table should be in cache" );
	return -1
}
func (i *iteratorCache) findTablebyEndIterator(ei int) *types.TableIDObject {
	//EOS_ASSERT( ei < -1, invalid_table_iterator, "not an end iterator" );
	indx := i.endIteratorToIndex(ei)

	if indx >= len(i.endIteratorToTable) {
		return nil
	}
	return i.endIteratorToTable[indx]
}
func (i *iteratorCache) get(iterator int) interface{} {
	// EOS_ASSERT( iterator != -1, invalid_table_iterator, "invalid iterator" );
	// EOS_ASSERT( iterator >= 0, table_operation_not_permitted, "dereference of end iterator" );
	// EOS_ASSERT( iterator < _iterator_to_object.size(), invalid_table_iterator, "iterator out of range" );
	//auto result = _iterator_to_object[iterator];

	obj := i.iteratorToObject[iterator]
	return obj

	//return nil
	//EOS_ASSERT( result, table_operation_not_permitted, "dereference of deleted object" );
}
func (i *iteratorCache) remove(iterator int) {
	// EOS_ASSERT( iterator != -1, invalid_table_iterator, "invalid iterator" );
	// EOS_ASSERT( iterator >= 0, table_operation_not_permitted, "cannot call remove on end iterators" );
	// EOS_ASSERT( iterator < _iterator_to_object.size(), invalid_table_iterator, "iterator out of range" );

	obj := i.iteratorToObject[iterator]
	i.iteratorToObject[iterator] = nil
	delete(i.objectToIterator, obj)

	//EOS_ASSERT( result, table_operation_not_permitted, "dereference of deleted object" );
}

func (i *iteratorCache) add(obj interface{}) int {
	if itr, ok := i.objectToIterator[obj]; ok {
		return itr
	}

	i.iteratorToObject = append(i.iteratorToObject, obj)
	i.objectToIterator[obj] = len(i.iteratorToObject) - 1
	return len(i.iteratorToObject) - 1
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
func (a *ApplyContext) IsAccount(n common.AccountName) bool {
	//return nullptr != db.find<account_object,by_name>( account );
	account := types.AccountObject{Name: n}
	err := a.Control.db.ByIndex("byName", &account)
	if err == nil {
		return true
	}
	return false

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
func (a *ApplyContext) DBStoreI64(scope int64, table int64, payer int64, id int64, buffer []byte) int {
	return a.dbStoreI64(int64(a.Receiver), scope, table, payer, id, buffer)
}
func (a *ApplyContext) dbStoreI64(code int64, scope int64, table int64, payer int64, id int64, buffer []byte) int {

	return 0
	//tab := a.FindOrCreateTable(common.Name(code), common.Name(scope), common.Name(table), common.AccountName(payer))
	//tid := tab.ID
	//
	//obj := types.KeyValueObject{
	//	TId:        tid,
	//	PrimaryKey: id,
	//	Value:      buffer,
	//	Payer:      payer,
	//	ID:         id,
	//}
	//a.Control.db.Insert(&obj)
	//
	//newTab := tab
	//newTab.Count++
	//a.Control.db.UpdateObject(&tab, &newTab)
	//
	//// int64_t billable_size = (int64_t)(buffer_size + config::billable_size_v<key_value_object>);
	////    update_db_usage( payer, billable_size);
	//
	//a.KeyvalCache.cacheTable(&newTab)
	//return a.KeyvalCache.add(&obj)

}
func (a *ApplyContext) DBUpdateI64(iterator int, payer common.AccountName, buffer []byte) {

	//obj := a.KeyvalCache.get(iterator)
	//objTable := a.KeyvalCache.getTable(obj.ID)
	//
	////EOS_ASSERT( table_obj.code == receiver, table_access_violation, "db access violation" );
	//// const int64_t overhead = config::billable_size_v<key_value_object>;
	////    int64_t old_size = (int64_t)(obj.value.size() + overhead);
	////    int64_t new_size = (int64_t)(buffer_size + overhead);
	//
	////    if( payer == account_name() ) payer = obj.payer;
	//
	////    if( account_name(obj.payer) != payer ) {
	////      // refund the existing payer
	////       update_db_usage( obj.payer,  -(old_size) );
	////      // charge the new payer
	////       update_db_usage( payer,  (new_size));
	////    } else if(old_size != new_size) {
	////      // charge/refund the existing payer the difference
	////       update_db_usage( obj.payer, new_size - old_size);
	////    }
	//
	//a.Control.db.ByIndex("ID", &obj)
	//objNew := obj
	//objNew.Value = buffer
	//objNew.Payer = payer
	//
	//a.Control.db.UpdateObject(&obj, &objNew)

}
func (a *ApplyContext) DBRemoveI64(iterator int) {
	//obj := a.KeyvalCache.get(iterator)
	//objTable := a.KeyvalCache.getTable(obj.ID)
	//
	//// 	EOS_ASSERT( table_obj.code == receiver, table_access_violation, "db access violation" );
	//// //   require_write_lock( table_obj.scope );
	////     update_db_usage( obj.payer,  -(obj.value.size() + config::billable_size_v<key_value_object>) );
	//a.Control.db.ByIndex("ID", &objTable)
	//
	//newTable := objTable
	//newTable.Count--
	//a.Control.db.UpdateObject(&objTable, &newTable)
	//
	//a.Control.db.Remove(&obj)
	//
	//if newTable.Count == 0 {
	//	a.Control.db.Remove(&newTable)
	//}
	//
	//a.KeyvalCache.remove(iterator)

}
func (a *ApplyContext) DBGetI64(iterator int, buffer []byte, bufferSize int) int {
	return 0
	//obj := a.KeyvalCache.get(iterator)
	//s := len(obj.value)
	//
	//if bufferSize == 0 {
	//	return s
	//}
	//
	//copySize = min(bufferSize, s)
	//copy(buffer[0:copySize], obj.value[:])
	//return copySize
}
func (a *ApplyContext) DBNextI64(iterator int, primary uint64) int {

	return 0
	// if iterator < -1 {return -1}
	// obj := a.KeyvalCache.get(iterator)

	// idx := a.Control.db.GetIndex("byScopePrimary",obj)
	// itr := idx.IteratorTo(obj)
	// itr ++
	// if itr == idx.end() || itr.TId != obj.TId {
	// 	return a.KeyvalCache.getEndIteratorByTableID(obj.TId)
	// }

	//setUint64(itr.primaryKey)
	// return a.KeyvalCache.add(*itr)
}

func (a *ApplyContext) DBPreviousI64(iterator int, primary uint64) int {
	return 0
	// idx := a.Control.db.GetIndex("byScopePrimary",obj)

	//     if iterator < -1 {
	//        tab = a.KeyvalCache.findTablebyEndIterator(iterator)
	//        //EOS_ASSERT( tab, invalid_table_iterator, "not a valid end iterator" );

	//        itr := idx.upperBound(tab.ID)
	//        if( idx.begin() == idx.end() || itr == idx.begin() ) return -1;

	//        itr --
	//        if( itr->TId != tab->ID ) return -1;

	//        setUint32(itr.PrimaryKey)
	//        return a.KeyvalCache.add(*itr)
	//     }

	// obj := a.KeyvalCache.get(iterator)
	// itr := idx.IteratorTo(obj)
	// itr --
	// if itr.TId != obj.TId {return -1}
	// setUint64(itr.primaryKey)

	// return keyval_cache.add(*itr);
}
func (a *ApplyContext) DBFindI64(code int64, scope int64, table int64, id int64) int {
	return 0

	// tab := a.FindTable(code, scope, table)
	// if tab == nil {
	// 	return -1
	// }

	// tableEndItr := a.KeyvalCache.cacheTable(tab)

	// objTable := tab
	// err := a.Control.db.ByIndex("ID", &objTable)

	// if err == nil {return tableEndItr}
	// return a.KeyvalCache.add(&objTable)

}
func (a *ApplyContext) DBLowerboundI64(iterator int, primary uint64) int    { return 0 }
func (a *ApplyContext) UpdateDBUsage(payer common.AccountName, delta int64) {}
func (a *ApplyContext) FindTable(
	code common.Name,
	scope common.Name,
	table common.Name) types.TableIDObject {
	return types.TableIDObject{}
}
func (a *ApplyContext) FindOrCreateTable(code common.Name,
	scope common.Name,
	table common.Name,
	payer common.AccountName) types.TableIDObject {
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
	oldGPO := a.Control.GetGlobalProperties()
	a.Control.db.UpdateObject(&oldGPO, &newGPO)

}

func (a *ApplyContext) GetBlockchainParametersPacked() []byte {
	gpo := a.Control.GetGlobalProperties()
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

	err := a.Control.db.ByIndex("byName", &account)
	if err != nil {
		log.Error("getaAccount is error detail:", err)
		return false
	}
	return account.Privileged

}
func (a *ApplyContext) SetPrivileged(n common.AccountName, isPriv bool) {
	oldAccount := types.AccountObject{Name: n}
	err := a.Control.db.ByIndex("byName", &oldAccount)
	if err != nil {
		log.Error("getaAccount is error detail:", err)
		return
	}

	newAccount := oldAccount
	newAccount.Privileged = isPriv
	a.Control.db.UpdateObject(&oldAccount, &newAccount)
}

//context producer api
func (a *ApplyContext) SetProposedProducers(data []byte) {

	// producers []types.ProducerKey
	// rlp.DecodeBytes(data, &producers)

	// uniqueProducers map[common.AccountName]bool
	// for _,v := range producers {
	// 	//assert(a.IsAccount(v), "producer schedule includes a nonexisting account")
	// 	has = uniqueProducers[v]
	// 	if has == nil {
	// 		uniqueProducers[v] = true
	// 	}
	// }

	// //assert(len(producer) == len(uniqueProducers),"duplicate producer name in producer schedule")
	// a.Controller.SetProposed_Producers(producers)
}

func (a *ApplyContext) GetActiveProducersInBytes() []byte {

	// ap := a.Controller.ActiveProducers()
	// accounts := make([]types.ProducerKey,len(ap.Producers))
	// for _,producer := range ap.Producers {
	// 	accounts = append(accounts,producer)
	// }

	// bytes,_ := rlp.EncodeToBytes(accounts)
	// return bytes
	return []byte{}
}

//func (a *ApplyContext) GetActiveProducers() []common.AccountName { return }

//context system api
func (a *ApplyContext) CheckTime() {
	a.TrxContext.CheckTime()
}
func (a *ApplyContext) CurrentTime() int64 {
	return a.Control.PendingBlockTime().TimeSinceEpoch().Count()
}
func (a *ApplyContext) PublicationTime() int64 {
	return a.TrxContext.Published.TimeSinceEpoch().Count()
}

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
