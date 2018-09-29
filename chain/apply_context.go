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

	IDX64         GenericIndex
	IDX128        GenericIndex
	IDX256        GenericIndex
	IDXDouble     GenericIndex
	IDXLongDouble GenericIndex

	//GenericIndex
	//_pending_console_output
	KeyvalCache          iteratorCache
	Notified             []common.AccountName
	PendingConsoleOutput string
}

type itrObjectInterface interface {
	GetBillableSize() uint64
}

type pairTableIterator struct {
	tableIDObject *types.TableIdObject
	iterator      int
}

type iteratorCache struct {
	tableCache         map[types.IdType]*pairTableIterator
	endIteratorToTable []*types.TableIdObject
	iteratorToObject   []itrObjectInterface
	objectToIterator   map[itrObjectInterface]int
}

func NewIteratorCache() *iteratorCache {

	i := &iteratorCache{
		tableCache:         make(map[types.IdType]*pairTableIterator),
		endIteratorToTable: make([]*types.TableIdObject, 8),
		iteratorToObject:   make([]itrObjectInterface, 32),
		objectToIterator:   make(map[itrObjectInterface]int),
	}

	return i
}

func (i *iteratorCache) endIteratorToIndex(ei int) int   { return (-ei - 2) }
func (i *iteratorCache) IndexToEndIterator(indx int) int { return -(indx + 2) }
func (i *iteratorCache) cacheTable(tobj *types.TableIdObject) int {
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
func (i *iteratorCache) getTable(id types.IdType) *types.TableIdObject {
	if itr, ok := i.tableCache[id]; ok {
		return itr.tableIDObject
	}

	return &types.TableIdObject{}
	//EOS_ASSERT( itr != _table_cache.end(), table_not_in_cache, "an invariant was broken, table should be in cache" );
}
func (i *iteratorCache) getEndIteratorByTableID(id types.IdType) int {
	if itr, ok := i.tableCache[id]; ok {
		return itr.iterator
	}
	//EOS_ASSERT( itr != _table_cache.end(), table_not_in_cache, "an invariant was broken, table should be in cache" );
	return -1
}
func (i *iteratorCache) findTablebyEndIterator(ei int) *types.TableIdObject {
	//EOS_ASSERT( ei < -1, invalid_table_iterator, "not an end iterator" );
	indx := i.endIteratorToIndex(ei)

	if indx >= len(i.endIteratorToTable) {
		return nil
	}
	return i.endIteratorToTable[indx]
}
func (i *iteratorCache) get(iterator int) itrObjectInterface {
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

func (i *iteratorCache) add(obj itrObjectInterface) int {
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
	err := a.DB.ByIndex("byName", &account)
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
func (a *ApplyContext) DbStoreI64(scope int64, table int64, payer int64, id int64, buffer []byte) int {
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
	//a.DB.Insert(&obj)

	// a.DB.Modify(tab, func(t *types.TableIDObject) {
	// 	t.Count++
	// })

	//
	//// int64_t billable_size = (int64_t)(buffer_size + config::billable_size_v<key_value_object>);
	////    UpdateDBUsage( payer, billable_size);
	// UpdateDBUsage( payer, len(buffer) + obj.GetBillableSize());
	//a.KeyvalCache.cacheTable(&tab)
	//return a.KeyvalCache.add(&obj)

}
func (a *ApplyContext) DbUpdateI64(iterator int, payer common.AccountName, buffer []byte) {

	// obj := a.KeyvalCache.get(iterator)
	// objTable := a.KeyvalCache.getTable(obj.ID)

	// //EOS_ASSERT( table_obj.code == receiver, table_access_violation, "db access violation" );

	// // const int64_t overhead = config::billable_size_v<key_value_object>;
	// overhead = obj.GetBillableSize()
	// oldSize := len(obj.Value) + overhead
	// newSize := len(buffer) + overhead

	//    if payer == common.AccountName{} { payer = obj.Payer}

	//    if obj.Payer == payer {
	//    	a.UpdateDBUsage(obj.Payer, -(oldSize))
	//    	a.UpdateDBUsage(payer, newSize)
	//    } else if oldSize != newSize{
	//    	a.UpdateDBUsage(obj.Payer, newSize - oldSize)
	//    }

	// a.DB.Modify(obj, func(t *types.KeyValueObject) {
	// 	t.Count++

	// 	obj.Value = buffer
	// 	obj.Payer = payer
	// })

}
func (a *ApplyContext) DbRemoveI64(iterator int) {
	// obj := a.KeyvalCache.get(iterator)
	// tab := a.KeyvalCache.getTable(obj.ID)

	// // 	EOS_ASSERT( table_obj.code == receiver, table_access_violation, "db access violation" );
	// // //   require_write_lock( table_obj.scope );
	// overhead := 0//config::billable_size_v<key_value_object>)
	// UpdateDBUsage( obj.Payer,  -(len(obj.Value) + overhead) )
	// a.DB.Get("ID", &tab)
	// a.DB.Modify(tab, func(t *types.TableIDObject) {
	// 	t.Count--
	// })

	// a.DB.Remove(&obj)

	// if tab.Count == 0 {
	// 	a.DB.Remove(&tab)
	// }
	// a.KeyvalCache.remove(iterator)

}
func (a *ApplyContext) DbGetI64(iterator int, buffer []byte, bufferSize int) int {
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
func (a *ApplyContext) DbNextI64(iterator int, primary *uint64) int {

	return 0
	// if iterator < -1 {
	// 	return -1
	// }
	// obj := a.KeyvalCache.get(iterator)

	// idx := a.DB.GetIndex("byScopePrimary", obj)
	// itr := idx.IteratorTo(obj)
	// itrNext := itr.Next()
	// objNext := types.KeyValueObject(itr.GetObject()) //return -1 for nil
	// if itr == idx.end() || objNext.TId != obj.TId {
	// 	return a.KeyvalCache.getEndIteratorByTableID(obj.TId)
	// }

	// *primary = itr.primaryKey
	// return a.KeyvalCache.add(objNext)
}

func (a *ApplyContext) DbPreviousI64(iterator int, primary *uint64) int {
	return 0
	// idx := a.DB.GetIndex("byScopePrimary",obj)

	// if iterator < -1 {
	//    tab = a.KeyvalCache.findTablebyEndIterator(iterator)
	//    //EOS_ASSERT( tab, invalid_table_iterator, "not a valid end iterator" );

	//    itr := idx.UpperBound(tab.ID)
	//    if( idx.begin() == idx.end() || itr == idx.begin() ) return -1;

	//    itrPrev := itr.Prev()
	//    objPrev := types.KeyValueObject(itr.GetObject())
	//    if( objPrev->TId != tab->ID ) return -1;

	//    setUint32(objPrev.PrimaryKey)
	//    return a.KeyvalCache.add(objPrev)
	// }

	// obj := a.KeyvalCache.get(iterator)
	// itr := idx.IteratorTo(obj)
	// itrPrev := itr.Prev()

	//    objPrev := types.KeyValueObject(itr.GetObject()) //return -1 for nil
	// if objPrev.TId != obj.TId {return -1}

	// *primary = objPrev.primaryKey
	// return keyval_cache.add(objPrev)
}
func (a *ApplyContext) DbFindI64(code int64, scope int64, table int64, id int64) int {
	return 0

	// tab := a.FindTable(code, scope, table)
	// if tab == nil {
	// 	return -1
	// }

	// tableEndItr := a.KeyvalCache.cacheTable(tab)

	// obj := types.KeyValueObject{TId:tab.ID,Primary:id}
	// err := a.DB.Get("byScopePrimary", &obj ) //, makeTupe(tab.ID,id))

	// if err == nil {return tableEndItr}
	// return a.KeyvalCache.add(&obj)

}
func (a *ApplyContext) DbLowerBoundI64(code int64, scope int64, table int64, id int64) int {
	return 0

	// tab := a.FindTable(code, scope, table)
	// if tab == nil {return -1}

	// tableEndItr := a.KeyvalCache.cacheTable(tab)

	// Obj := types.KeyValueObject{}
	// idx := a.DB.GetIndex("byScopePrimary",&Obj)

	// itr := idx.LowerBound(makeTupe(tab.ID,id))
	// if itr == idx.End()  {return tableEndItr}

	// obj := types.KeyValueObject(itr.GetObject())
	// return keyval_cache.add(types.KeyValueObject(itr.GetObject()))

}
func (a *ApplyContext) DbUpperBoundI64(code int64, scope int64, table int64, id int64) int {
	return 0

	// tab := a.FindTable(code, scope, table)
	// if tab == nil {return -1}

	// tableEndItr := a.KeyvalCache.cacheTable(tab)

	// Obj := types.KeyValueObject{}
	// idx := a.DB.GetIndex("byScopePrimary",&Obj)

	// itr := idx.UpperBound(makeTupe(tab.ID,id))
	// if itr == idx.End()  {return tableEndItr}

	// obj := types.KeyValueObject(itr.GetObject())
	// if obj.ID != tab.ID {return tableEndItr}

	// return keyval_cache.add(obj)

}
func (a *ApplyContext) DbEndI64(code int64, scope int64, table int64) int {
	return 0

	tab := a.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}

	return a.KeyvalCache.cacheTable(tab)
}

//index for sceondarykey
func (a *ApplyContext) IdxI64Store(scope int64, table int64, payer int64, id int64, value *types.Uint64_t) int {
	return a.IDX64.store(scope, table, payer, id, value)
}
func (a *ApplyContext) IdxI64Remove(iterator int) {
	a.IDX64.remove(iterator)
}
func (a *ApplyContext) IdxI64Update(iterator int, payer int64, value *types.Uint64_t) {
	a.IDX64.update(iterator, payer, value)
}
func (a *ApplyContext) IdxI64FindSecondary(code int64, scope int64, table int64, secondary *types.Uint64_t, primary *uint64) int {
	//a.IDX64.update(iterator, payer, value)
	return a.IDX64.findSecondary(code, scope, table, secondary, primary)
}
func (a *ApplyContext) IdxI64LowerBound(code int64, scope int64, table int64, secondary *types.Uint64_t, primary *uint64) int {
	//a.IDX64.update(iterator, payer, value)
	return a.IDX64.lowerbound(code, scope, table, secondary, primary)
}
func (a *ApplyContext) IdxI64UpperBound(code int64, scope int64, table int64, secondary *types.Uint64_t, primary *uint64) int {
	//a.IDX64.update(iterator, payer, value)
	return a.IDX64.upperbound(code, scope, table, secondary, primary)
}
func (a *ApplyContext) IdxI64End(code int64, scope int64, table int64) int {
	//a.IDX64.update(iterator, payer, value)
	return a.IDX64.end(code, scope, table)
}

func (a *ApplyContext) IdxI64Next(iterator int, primary *uint64) int {
	return a.IDX64.next(iterator, primary)
}
func (a *ApplyContext) IdxI64Previous(iterator int, primary *uint64) int {
	return a.IDX64.previous(iterator, primary)
}

func (a *ApplyContext) FindTable(code int64, scope int64, table int64) *types.TableIdObject {
	// table := types.TableIDObject{Code: common.AccountName(code), Scope: common.ScopeName(scope), Table: common.TableName(table)}
	// a.DB.Get("byCodeScopeTable", &table)
	// return table
	return &types.TableIdObject{}
}
func (a *ApplyContext) FindOrCreateTable(code int64, scope int64, table int64, payer int64) types.TableIdObject {

	return types.TableIdObject{}
	// table := types.TableIDObject{Code: common.AccountName(code), Scope: common.ScopeName(scope), Table: common.TableName(table), Payer: common.AccountName(payer)}
	// err := a.DB.Get("byCodeScopeTable", &table)
	// if err == nil {
	// 	return table
	// }
	// a.DB.Insert(&table)
	// return table
}
func (a *ApplyContext) RemoveTable(tid types.TableIdObject) {
	// overhead := 0 //config::billable_size_v<table_id_object>

	// UpdateDBUsage(tid.Payer, -overhead)
	// a.DB.remove(tid)
}

func (a *ApplyContext) UpdateDbUsage(payer common.AccountName, delta int64) {

}

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

	cfg := common.Config{}
	rlp.DecodeBytes(parameters, &cfg)

	// a.DB.modify(a.Control.GetGlobalProperties(), func(gpo *types.GlobalPropertyObject){
	//       gpo.Configuration = cfg
	// })

}

func (a *ApplyContext) GetBlockchainParametersPacked() []byte {
	gpo := a.Control.GetGlobalProperties()
	bytes, err := rlp.EncodeToBytes(gpo.Configuration)
	if err != nil {
		log.Error("EncodeToBytes is error detail:", err)
		return nil
	}
	return bytes
}
func (a *ApplyContext) IsPrivileged(n common.AccountName) bool {
	//return false
	account := types.AccountObject{Name: n}

	err := a.DB.ByIndex("byName", &account)
	if err != nil {
		log.Error("getaAccount is error detail:", err)
		return false
	}
	return account.Privileged

}
func (a *ApplyContext) SetPrivileged(n common.AccountName, isPriv bool) {
	oldAccount := types.AccountObject{Name: n}
	err := a.DB.ByIndex("byName", &oldAccount)
	if err != nil {
		log.Error("getaAccount is error detail:", err)
		return
	}

	newAccount := oldAccount
	newAccount.Privileged = isPriv
	a.DB.UpdateObject(&oldAccount, &newAccount)
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
func (a *ApplyContext) ScheduleDeferredTransaction(sendId common.TransactionIdType, payer common.AccountName, trx []byte, replaceExisting bool) {
}
func (a *ApplyContext) CancelDeferredTransaction(sendId common.TransactionIdType) bool { return false }
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
