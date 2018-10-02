package chain

import (
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

	idx64     Idx64
	idxDouble IdxDouble
	// IDX128        GenericIndex
	// IDX256        GenericIndex
	// IDXLongDouble GenericIndex

	//GenericIndex
	KeyvalCache          iteratorCache
	Notified             []common.AccountName
	PendingConsoleOutput string
}

// type itrObjectInterface interface {
// 	GetBillableSize() uint64
// 	// GetTableId() types.IdType
// 	// GetValue() common.HexBytes
// }

type pairTableIterator struct {
	tableIDObject *types.TableIdObject
	iterator      int
}

type iteratorCache struct {
	tableCache         map[types.IdType]*pairTableIterator
	endIteratorToTable []*types.TableIdObject
	iteratorToObject   []interface{}
	objectToIterator   map[interface{}]int
}

func NewIteratorCache() *iteratorCache {

	i := &iteratorCache{
		tableCache:         make(map[types.IdType]*pairTableIterator),
		endIteratorToTable: make([]*types.TableIdObject, 8),
		iteratorToObject:   make([]interface{}, 32),
		objectToIterator:   make(map[interface{}]int),
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
func (a *ApplyContext) RequireAuthorization(account int64) {
	for k, v := range a.Act.Authorization {
		if v.Actor == common.AccountName(account) {
			a.UsedAuthorizations[k] = true
			return
		}
	}
	// EOS_ASSERT( false, missing_auth_exception, "missing authority of ${account}/${permission}",
	//              ("account",account)("permission",permission) );
}
func (a *ApplyContext) HasAuthorization(account int64) bool {
	for _, v := range a.Act.Authorization {
		if v.Actor == common.AccountName(account) {
			return true
		}
	}
	return false
}
func (a *ApplyContext) RequireAuthorization2(account int64, permission int64) {
	for k, v := range a.Act.Authorization {
		if v.Actor == common.AccountName(account) && v.Permission == common.PermissionName(permission) {
			a.UsedAuthorizations[k] = true
			return
		}
	}

	//EOS_ASSERT( false, missing_auth_exception, "missing authority of ${account}", ("account",account));
}

//func (a *ApplyContext) RequireAuthorizations(account common.AccountName) {}
func (a *ApplyContext) RequireRecipient(recipient int64) {
	if a.HasReciptient(recipient) {
		a.Notified = append(a.Notified, common.AccountName(recipient))
	}
}
func (a *ApplyContext) IsAccount(n int64) bool {
	return true
	// account := types.AccountObject{Name: common.AccountName(n)}
	// err := a.DB.ByIndex("byName", &account)
	// if err == nil {
	// 	return true
	// }
	// return false
}
func (a *ApplyContext) HasReciptient(code int64) bool {
	for _, a := range a.Notified {
		if a == common.AccountName(code) {
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

	// tab := a.FindOrCreateTable(code, scope, table, payer)
	// tid := tab.ID

	// obj := &types.KeyValueObject{
	// 	TId:        tid,
	// 	PrimaryKey: id,
	// 	Value:      buffer,
	// 	Payer:      common.AccountName(payer),
	// }

	// a.DB.Insert(obj)
	// a.DB.Modify(tab, func(t *types.TableIDObject) {
	// 	t.Count++
	// })

	// // int64_t billable_size = (int64_t)(buffer_size + config::billable_size_v<key_value_object>);
	// billableSize := len(buffer) + types.BillableSizeV(obj.GetBillableSize())
	// UpdateDbUsage( payer, billableSize )
	// a.KeyvalCache.cacheTable(tab)
	// return a.KeyvalCache.add(obj)
}
func (a *ApplyContext) DbUpdateI64(iterator int, payer int64, buffer []byte) {

	// obj := (*types.KeyValueObject)(a.KeyvalCache.get(iterator))
	// objTable := a.KeyvalCache.getTable(obj.GetTableId())

	// //EOS_ASSERT( objTable.Code == a.Receiver, table_access_violation, "db access violation" );

	// // const int64_t overhead = config::billable_size_v<key_value_object>;
	// overhead = types.BillableSizeV(obj.GetBillableSize())
	// oldSize := len(obj.Value) + overhead
	// newSize := len(buffer) + overhead

	// payerAccount := common.AccountName(payer)
	// if payerAccount == common.AccountName{} { payerAccount = obj.Payer}

	// if obj.Payer == payerAccount {
	// 	a.UpdateDbUsage(obj.Payer, -(oldSize))
	// 	a.UpdateDbUsage(payerAccount, newSize)
	// } else if oldSize != newSize{
	// 	a.UpdateDbUsage(obj.Payer, newSize - oldSize)
	// }

	// a.DB.Modify(obj, func(obj *types.KeyValueObject) {
	// 	obj.Value = buffer
	// 	obj.Payer = payerAccount
	// })
}
func (a *ApplyContext) DbRemoveI64(iterator int) {
	// obj := (*types.KeyValueObject)(a.KeyvalCache.get(iterator))
	// tab := a.KeyvalCache.getTable(obj.ID)
	// // 	EOS_ASSERT( table_obj.code == receiver, table_access_violation, "db access violation" );
	// // //   require_write_lock( table_obj.scope );
	// billableSize := len(buffer) + types.BillableSizeV(obj.GetBillableSize())
	// UpdateDBUsage( obj.Payer,  - billableSize )
	// a.DB.Modify(tab, func(t *types.TableIdObject) {
	// 	t.Count --
	// })

	// a.DB.Remove(obj)
	// if tab.Count == 0 {
	// 	a.DB.Remove(tab)
	// }
	// a.KeyvalCache.remove(iterator)
}
func (a *ApplyContext) DbGetI64(iterator int, buffer []byte, bufferSize int) int {
	return 0
	// obj := (*types.KeyValueObject)(a.KeyvalCache.get(iterator))
	// s := len(obj.value)

	// if bufferSize == 0 {
	// 	return s
	// }

	// copySize = min(bufferSize, s)
	// copy(buffer[0:copySize], obj.value[0:copySize])
	// return copySize
}
func (a *ApplyContext) DbNextI64(iterator int, primary *uint64) int {

	return 0
	// if iterator < -1 {
	// 	return -1
	// }
	// obj := (*types.KeyValueObject)(a.KeyvalCache.get(iterator))
	// idx := a.DB.GetIndex("byScopePrimary", obj)

	// itr := idx.IteratorTo(obj)
	// itrNext := itr.Next()
	// objNext := ( *types.KeyValueObject )(itr.GetObject())
	// if itr == idx.End() || objNext.TId  != obj.TId  {
	// 	return a.KeyvalCache.getEndIteratorByTableID(obj.TId
	// }

	// *primary = objNext.primaryKey
	// return a.KeyvalCache.add(objNext)
}

func (a *ApplyContext) DbPreviousI64(iterator int, primary *uint64) int {
	return 0

	// idx := a.DB.GetIndex("byScopePrimary", &types.KeyValueObject{})

	// if iterator < -1 {
	//    tab = a.KeyvalCache.findTableByEndIterator(iterator)
	//    //EOS_ASSERT( tab, invalid_table_iterator, "not a valid end iterator" );

	//    itr := idx.Upperbound(tab.ID)
	//    if( idx.Begin() == idx.End() || itr == idx.Begin() ) return -1;

	//    itrPrev := itr.Prev()
	//    objPrev := (*types.KeyValueObject)(itr.GetObject())
	//    if( objPrev.TId != tab.ID ) return -1;

	//    *primary =  objPrev.PrimaryKey
	//    return a.KeyvalCache.add(objPrev)
	// }

	// obj := (*types.KeyValueObject)(a.KeyvalCache.get(iterator))
	// itr := idx.IteratorTo(obj)
	// itrPrev := itr.Prev()

	// objPrev := (*types.KeyValueObject)(itr.GetObject()) //return -1 for nil
	// if objPrev.TId != obj.TId  {return -1}

	// *primary = objPrev.primaryKey
	// return a.KeyvalCache.add(objPrev)
}
func (a *ApplyContext) DbFindI64(code int64, scope int64, table int64, id int64) int {
	return 0

	// tab := a.FindTable(code, scope, table)
	// if tab == nil {
	// 	return -1
	// }

	// tableEndItr := a.KeyvalCache.cacheTable(tab)

	// obj := &types.KeyValueObject{}
	// err := a.DB.Get("byScopePrimary", obj, obj.MakeTuple(tab.ID, id) )

	// if err == nil {return tableEndItr}
	// return a.KeyvalCache.add(obj)

}
func (a *ApplyContext) DbLowerboundI64(code int64, scope int64, table int64, id int64) int {
	return 0

	// tab := a.FindTable(code, scope, table)
	// if tab == nil {return -1}

	// tableEndItr := a.KeyvalCache.cacheTable(tab)

	// obj := &types.KeyValueObject{}
	// idx := a.DB.GetIndex("byScopePrimary", obj)

	// itr := idx.Lowerbound(obj.MakeTuple(tab.ID,id))
	// if itr == idx.End()  {return tableEndItr}

	// objLowerbound = (*types.KeyValueObject)(*itr.GetObject())
	// if objLowerbound.TId != tab.ID {return tableEndItr}

	// return keyval_cache.add(objLowerbound)

}
func (a *ApplyContext) DbUpperboundI64(code int64, scope int64, table int64, id int64) int {
	return 0

	// tab := a.FindTable(code, scope, table)
	// if tab == nil {return -1}

	// tableEndItr := a.KeyvalCache.cacheTable(tab)

	// obj := &types.KeyValueObject{}
	// idx := a.DB.GetIndex("byScopePrimary", &obj)

	// itr := idx.Upperbound(obj.MakeTuple(tab.ID,id))
	// if itr == idx.End()  {return tableEndItr}

	// objUpperbound = (*types.KeyValueObject)(*itr.GetObject())
	//    if objUpperbound.TId != tab.ID {return tableEndItr}

	// return keyval_cache.add(objUpperbound)

}
func (a *ApplyContext) DbEndI64(code int64, scope int64, table int64) int {
	return 0

	// tab := a.FindTable(code, scope, table)
	// if tab == nil {
	// 	return -1
	// }

	// return a.KeyvalCache.cacheTable(tab)
}

//index for sceondarykey
func (a *ApplyContext) Idx64Store(scope int64, table int64, payer int64, id int64, value *types.Uint64_t) int {
	return a.idx64.store(scope, table, payer, id, value)
}
func (a *ApplyContext) Idx64Remove(iterator int) {
	a.idx64.remove(iterator)
}
func (a *ApplyContext) Idx64Update(iterator int, payer int64, value *types.Uint64_t) {
	a.idx64.update(iterator, payer, value)
}
func (a *ApplyContext) Idx64FindSecondary(code int64, scope int64, table int64, secondary *types.Uint64_t, primary *uint64) int {
	//a.idx64.update(iterator, payer, value)
	return a.idx64.findSecondary(code, scope, table, secondary, primary)
}
func (a *ApplyContext) Idx64Lowerbound(code int64, scope int64, table int64, secondary *types.Uint64_t, primary *uint64) int {
	//a.idx64.update(iterator, payer, value)
	return a.idx64.lowerbound(code, scope, table, secondary, primary)
}
func (a *ApplyContext) Idx64Upperbound(code int64, scope int64, table int64, secondary *types.Uint64_t, primary *uint64) int {
	return a.idx64.upperbound(code, scope, table, secondary, primary)
}
func (a *ApplyContext) Idx64End(code int64, scope int64, table int64) int {
	return a.idx64.end(code, scope, table)
}
func (a *ApplyContext) Idx64Next(iterator int, primary *uint64) int {
	return a.idx64.next(iterator, primary)
}
func (a *ApplyContext) Idx64Previous(iterator int, primary *uint64) int {
	return a.idx64.previous(iterator, primary)
}
func (a *ApplyContext) Idx64FindPrimary(code int64, scope int64, table int64, secondary *types.Uint64_t, primary *uint64) int {
	//a.idx64.update(iterator, payer, value)
	return a.idx64.findPrimary(code, scope, table, secondary, primary)
}

func (a *ApplyContext) IdxDoubleStore(scope int64, table int64, payer int64, id int64, value *types.Float64_t) int {
	return a.idxDouble.store(scope, table, payer, id, value)
}
func (a *ApplyContext) IdxDoubleRemove(iterator int) {
	a.idxDouble.remove(iterator)
}
func (a *ApplyContext) IdxDoubleUpdate(iterator int, payer int64, value *types.Float64_t) {
	a.idxDouble.update(iterator, payer, value)
}
func (a *ApplyContext) IdxDoubleFindSecondary(code int64, scope int64, table int64, secondary *types.Float64_t, primary *uint64) int {
	return a.idxDouble.findSecondary(code, scope, table, secondary, primary)
}
func (a *ApplyContext) IdxDoubleLowerbound(code int64, scope int64, table int64, secondary *types.Float64_t, primary *uint64) int {
	return a.idxDouble.lowerbound(code, scope, table, secondary, primary)
}
func (a *ApplyContext) IdxDoubleUpperbound(code int64, scope int64, table int64, secondary *types.Float64_t, primary *uint64) int {
	return a.idxDouble.upperbound(code, scope, table, secondary, primary)
}
func (a *ApplyContext) IdxDoubleEnd(code int64, scope int64, table int64) int {
	return a.idxDouble.end(code, scope, table)
}
func (a *ApplyContext) IdxDoubleNext(iterator int, primary *uint64) int {
	return a.idxDouble.next(iterator, primary)
}
func (a *ApplyContext) IdxDoublePrevious(iterator int, primary *uint64) int {
	return a.idxDouble.previous(iterator, primary)
}
func (a *ApplyContext) IdxDoubleFindPrimary(code int64, scope int64, table int64, secondary *types.Float64_t, primary *uint64) int {
	return a.idxDouble.findPrimary(code, scope, table, secondary, primary)
}

func (a *ApplyContext) FindTable(code int64, scope int64, table int64) *types.TableIdObject {
	// // table := types.TableIdObject{Code: common.AccountName(code), Scope: common.ScopeName(scope), Table: common.TableName(table)}
	// table := types.TableIdObject{}
	// a.DB.Get("byCodeScopeTable", &table,  table.MakeTuple(code,scope,table))
	// return table
	return &types.TableIdObject{}
}
func (a *ApplyContext) FindOrCreateTable(code int64, scope int64, table int64, payer int64) *types.TableIdObject {

	return &types.TableIdObject{}
	// //table := types.TableIdObject{Code: common.AccountName(code), Scope: common.ScopeName(scope), Table: common.TableName(table), Payer: common.AccountName(payer)}
	// table := types.TableIdObject{}
	// err := a.DB.Get("byCodeScopeTable", &table, table.MakeTuple(code,scope,table))
	// if err == nil {
	// 	return &table
	// }
	// a.DB.Insert(&table)
	// return &table
}
func (a *ApplyContext) RemoveTable(tid types.TableIdObject) {
	// UpdateDBUsage(tid.Payer, -types.TableIdObject{}.GetBillableSize())
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
	copySize := common.Min(uint64(bufferSize), uint64(s))
	return int(copySize), trx.ContextFreeData[index][0:copySize]

}
