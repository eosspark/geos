package chain

import (
	"github.com/eosspark/eos-go/chain/types"
)

type GenericIndex struct {
	//secondaryKey SecondaryKeyInterface
	context    *ApplyContext
	ObjectType *types.SecondaryObject
	itrCache   *iteratorCache
}

func NewGenericIndex(c *ApplyContext, o *types.SecondaryObject) *GenericIndex {
	return &GenericIndex{
		context:    c,
		ObjectType: o,
		itrCache:   NewIteratorCache(),
	}
}

func (i *GenericIndex) store(scope int64, table int64, payer int64, id int64, secondary types.SecondaryKeyInterface) int {
	return 0
	//EOS_ASSERT( payer != account_name(), invalid_table_payer, "must specify a valid account to pay for new record" );
	// tab := i.context.FindOrCreateTable(int64(i.context.Receiver), scope, table, payer)

	// obj := types.SecondaryObject{
	// 	  TId:          tab.ID,
	// 	  PrimaryKey:   id,
	//    SecondaryKey: *secondary,
	// 	  Payer:        payer,
	// }
	// i.context.DB.Insert(&obj)
	// i.context.DB.Modify(tab, func(t *types.TableIDObject) {
	// 	  t.Count++
	// })

	// BillableSize := secondary.Size()+secondary.GetBillableSize()
	// i.context.UpdateDbUsage(payer, secondary.Size()+overhead)

	// i.itrCache.cacheTable(&tab)
	// return i.itrCache.add(&obj)
}

func (i *GenericIndex) remove(iterator int) int {
	return 0

	// obj := &types.SecondaryObject(*i.itrCache.get(iterator))
	// tab := i.itrCache.getTable(obj.TId)

	// i.context.UpdateDbUsage(obj.payer, -obj.GetBillableSize())
	// i.context.DB.Modify(tab, func(t *types.TableIdObject) {
	// 	t.Count--
	// })

	// i.context.DB.Remove(&obj)
	// if tab.Count == 0 {
	// 	i.context.Remove(&tab)
	// }
	// i.itrCache.remove(iterator)
}

func (i *GenericIndex) update(iterator int, payer int64, secondary types.SecondaryKeyInterface) {

	// obj := &types.SecondaryObject(*i.itrCache.get(iterator))
	// objTable := i.itrCache.getTable(obj.TId)

	// //EOS_ASSERT( table_obj.code == i.context.Receiver, table_access_violation, "db access violation" )
	// if payer == common.AccountName{} payer = obj.Payer

	// billingSize := obj.GetBillableSize()
	// if obj.Payer != payer {
	// 	i.context.UpdateDbUsage(obj.Payer, - billingSize)
	// 	i.context.UpdateDbUsage(payer, + billingSize)
	// }

	// i.context.DB.Modify(obj, func(o *types.SecondaryObject){
	// 	o.SecondaryKey = *secondary
	// 	o.Payer = payer
	// })
}

func (i *GenericIndex) findSecondary(code int64, scope int64, table int64, secondary types.SecondaryKeyInterface, primary *uint64) int {
	return 0
	// tab := i.context.FindTable(code, scope, table)
	// if tab == nil {return -1}

	// tableEndItr := i.itrCache.cacheTable(&tab)

	// obj := &types.SecondaryObject{TId:tab.ID,SecondaryKey:*secondary}
	// err := i.context.DB.get("bySecondary", obj, obj.MakeTuple(tab.ID, *secondary)

	// *primary = obj.PrimaryKey

	// if err == nil {return tableEndItr}
	// return i.itrCache.add(&obj)
}

func (i *GenericIndex) lowerbound(code int64, scope int64, table int64, secondary types.SecondaryKeyInterface, primary *uint64) int {
	return 0
	// tab := i.context.FindTable(code, scope, table)
	// if tab == nil {
	// 	return -1
	// }

	// tableEndItr := i.itrCache.cacheTable(&tab)

	// obj := types.SecondaryObject{SecondaryKey:*secondary}

	// idx := i.context.DB.GetIndex("bySecondary", &obj)
	// itr := idx.LowerBound(obj.maketuple(tab.ID, *secondary))

	// itrObject := &types.SecondaryObject(*itr.GetObject())

	// *primary = itrObject.PrimaryKey
	// *secondary = itrObject.SecondaryKey

	// return i.itrCache.add(itrObject)
}

func (i *GenericIndex) upperbound(code int64, scope int64, table int64, secondary types.SecondaryKeyInterface, primary *uint64) int {
	return 0
	// tab := i.context.FindTable(code, scope, table)
	// if tab == nil {
	// 	return -1
	// }

	// tableEndItr := i.itrCache.cacheTable(&tab)

	// obj := types.SecondaryObject{SecondaryKey:*secondary}

	// idx := i.context.DB.GetIndex("bySecondary", &obj)
	// itr := idx.UpperBound(obj.maketuple(tab.ID, *secondary))
	// if itr == idx.End() {
	// 	return tableEndItr
	// }

	// o = &types.SecondaryObject(*itr.GetObject())
	// if o.TId != tab.ID {
	// 	return tableEndItr
	// }

	// *primary = o.PrimaryKey
	// *secondary = o.SecondaryKey

	// return i.itrCache.add(o)
}

func (i *GenericIndex) end(code int64, scope int64, table int64) int {
	return 0

	// tab := i.context.FindTable(code, scope, table)
	// if tab == nil {
	// 	return -1
	// }
	// return i.itrCache.cacheTable(&tab)
}

func (i *GenericIndex) next(iterator int, primary *uint64) int {
	return 0

	// if iterator < -1 {
	// 	return -1
	// }
	// obj := &types.SecondaryObject(*i.itrCache.get(iterator))

	// idx := i.context.DB.GetIndex("bySecondary", obj)
	// itr := idx.iteratorTo(obj)

	// itrNext := itr.Next()
	// objNext := &types.SecondaryObject(*itrNext.GetObject())

	// if itr == idx.End() || objNext.TId != obj.TId {
	// 	return i.itrCache.getEndIteratorByTableID(obj.TId)
	// }

	// *primary = objNext.PrimaryKey
	// return i.itrCache.add(objNext)

}

func (i *GenericIndex) previous(iterator int, primary *uint64) int {
	return 0

	// idx := i.context.DB.GetIndex("bySecondary", ObjectType)

	// if( iterator < -1) {
	//     tab = i.itrCache.findTablebyEndIterator(iterator)
	//    //EOS_ASSERT( tab, invalid_table_iterator, "not a valid end iterator" );

	//    itr := idx.UpperBound(tab.ID)
	//    if( idx.begin() == idx.end() || itr == idx.begin() ) return -1;

	//    itrPrev := itr.Prev()
	//    objPrev := itr.GetObject()
	//    if( objPrev.TId != tab->ID ) return -1;

	//    *primary = objPrev.PrimaryKey
	//    return a.KeyvalCache.add(objPrev)
	// }

	// obj := i.itrCache.get(iterator)
	// itr := idx.iteratorTo(obj)

	// if itr == idx.begin() {
	// 	return -1
	// }
	// itrNext := itr.Next()
	// objNext := itr.GetObject()

	// if objNext.TId != obj.TId {
	// 	return -1
	// }
	// *primary = objNext.PrimaryKey
	// return i.itrCache.add(objNext)
}

func (i *GenericIndex) findPrimary(code int64, scope int64, table int64, secondary types.SecondaryKeyInterface, primary *uint64) int {
	return 0

	// tab := i.context.FindTable(code, scope, table)
	// if tab == nil {
	// 	return -1
	// }

	// tableEndItr := i.itrCache.cacheTable(&tab)

	// obj := types.SecondaryObject{TId: tab.ID, PrimaryKey: *primary}
	// err := i.context.DB.get("byPrimary", &obj, ObjectType.makeTuple(tab.ID, *primary))

	// *secondary = obj.SecondaryKey

	// if err == nil {
	// 	return tableEndItr
	// }
	// return i.itrCache.add(&obj)
}

func (i *GenericIndex) lowerboundPrimary(code int64, scope int64, table int64, primary uint64) int {
	return 0

	// tab := i.context.FindTable(code, scope, table)
	// if tab == nil {
	// 	return -1
	// }

	// tableEndItr := i.itrCache.cacheTable(&tab)

	// //obj := types.SecondaryObject{TId: tab.ID, PrimaryKey: *primary}
	// idx := i.context.DB.GetIndex("byPrimary", ObjectType)
	// itr := idex.LowerBound(ObjectType.MakeTuple(tab.ID, primary))

	// *secondary = obj.SecondaryKey

	// if err == nil {
	// 	return tableEndItr
	// }
	// return i.itrCache.add(&obj)
}

func (i *GenericIndex) upperboundPrimary(code int64, scope int64, table int64, primary uint64) int {
	return 0

	// tab := i.context.FindTable(code, scope, table)
	// if tab == nil {
	// 	return -1
	// }

	// tableEndItr := i.itrCache.cacheTable(&tab)

	// //obj := types.SecondaryObject{TId: tab.ID, PrimaryKey: *primary}
	// idx := i.context.DB.GetIndex("byPrimary", ObjectType)
	// itr := idx.UpperBound(ObjectType.MakeTuple(tab.ID, primary))
	// if itr == idx.End() {
	// 	return tableEndItr
	// }

	// obj := itr.GetObject()
	// if obj.TId != tab.ID {
	// 	return tableEndItr
	// }

	// return i.itrCache.add(obj)
}

func (i *GenericIndex) nextPrimary(iterator int, primary *uint64) int {
	return 0

	// if iterator < -1 {
	// 	return -1
	// }
	// obj := i.itrCache.get(iterator)
	// idx := i.context.DB.GetIndex("byPrimary", ObjectType)

	// //idx := i.context.DB.GetIndex("bySecondary", obj)
	// itr := idx.iteratorTo(obj)

	// itrNext := itr.Next()
	// objNext := itrNext.GetObject()

	// if itr == idx.End() || objNext.TId != obj.TId {
	// 	return i.itrCache.getEndIteratorByTableID(obj.TId)
	// }

	// *primary = objNext.PrimaryKey
	// return i.itrCache.add(objNext)

}

func (i *GenericIndex) previousPrimary(iterator int, primary *uint64) int {
	return 0

	// idx := i.context.DB.GetIndex("byPrimary", ObjectType)

	// if( iterator < -1) {
	//     tab = i.itrCache.findTablebyEndIterator(iterator)
	//    //EOS_ASSERT( tab, invalid_table_iterator, "not a valid end iterator" );

	//    itr := idx.UpperBound(tab.ID)
	//    if( idx.begin() == idx.end() || itr == idx.begin() ) return -1;

	//    itrPrev := itr.Prev()
	//    objPrev := itr.GetObject()
	//    if( objPrev.TId != tab->ID ) return -1;

	//    *primary = objPrev.PrimaryKey
	//    return a.KeyvalCache.add(objPrev)
	// }

	// obj := i.itrCache.get(iterator)
	// itr := idx.iteratorTo(obj)

	// if itr == idx.begin() {
	// 	return -1
	// }
	// itrNext := itr.Next()
	// objNext := itr.GetObject()

	// if objNext.TId != obj.TId {
	// 	return -1
	// }
	// *primary = objNext.PrimaryKey
	// return i.itrCache.add(objNext)
}

func (i *GenericIndex) get(iterator int, secondary types.SecondaryKeyInterface, primary *uint64) {
	// obj := i.itrCache.get(iterator)

	// *secondary = obj.SecondaryKey
	// *primary = obj.PrimaryKey
}
