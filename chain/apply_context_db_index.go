package chain

import (
	"github.com/eosspark/eos-go/chain/types"
)

type DBGenericIndex struct {
	//secondaryKey SecondaryKeyInterface
	context  *ApplyContext
	itrCache *iteratorCache
}

func NewDBIndex(c *ApplyContext) *DBGenericIndex {

	return &DBGenericIndex{
		context:  c,
		itrCache: NewIteratorCache(),
	}

}

func (i *DBGenericIndex) store(scope int64, table int64, payer int64, id int64, secondary types.SecondaryKeyInterface) int {
	return 0
	//EOS_ASSERT( payer != account_name(), invalid_table_payer, "must specify a valid account to pay for new record" );

	// tab := i.context.FindOrCreateTable(i.context.Receiver, scope, table, payer)

	// obj := types.SecondaryObject{
	// 	TId:          tab.ID,
	// 	PrimaryKey:   id,
	// 	SecondaryKey: *secondary,
	// 	Payer:        payer,
	// }
	// i.context.DB.Insert(&obj)
	// i.context.DB.Modify(tab, func(t *types.TableIDObject) {
	// 	t.Count++
	// })

	// overhead := 0 //config::billable_size_v<key_value_object>)
	// i.context.UpdateDBUsage(payer, secondary.Size()+overhead)

	// i.itrCache.cacheTable(&tab)
	// return i.itrCache.add(&obj)
}

func (i *DBGenericIndex) remove(iterator int) int {
	return 0

	// obj := i.itrCache.get(iterator)
	//    tab := i.itrCache.getTable(obj.ID)

	// i.context.UpdateDBUsage( obj.payer, - obj.GetBillableSize() );
	// i.context.DB.Modify(tab, func(t *types.TableIDObject) {
	// 	t.Count--
	// })

	// i.context.DB.Remove(&obj)
	// if( tab.Count == 0){
	// 	i.context.Remove(&tab)
	// }
	// i.itrCache.remove(iterator)
}

func (i *DBGenericIndex) update(iterator int, payer int64, secondary types.SecondaryKeyInterface) {

	// obj := i.itrCache.get(iterator)
	// objTable := i.itrCache.getTable(obj.TId)

	// //EOS_ASSERT( table_obj.code == i.context.Receiver, table_access_violation, "db access violation" )
	// if payer == common.AccountName{} payer = obj.Payer

	// billingSize := obj.GetBillableSize()
	//    if obj.Payer != payer {
	//    	i.context.UpdateDBUsage(obj.Payer, - billingSize)
	//    	i.context.UpdateDBUsage(payer, + billingSize)
	//    }

	//    i.context.DB.Modify(obj,func(o *types.SecondaryKeyInterface){
	//    	o.SecondaryKey = *secondary
	//    	o.Payer = payer
	//    })
}

func (i *DBGenericIndex) findSecondary(code int64, scope int64, table int64, secondary types.SecondaryKeyInterface, primary *uint64) int {
	return 0
	// tab := i.context.FindTable(code, scope, table)
	// if tab == nil {return -1}

	// tableEndItr := i.itrCache.tableCache(&tab)

	// obj := types.SecondaryObject{TId:tab.ID,SecondaryKey:secondary}
	// err := i.context.DB.get("bySecondary", &obj)//,obj.makeTuple())

	// if err == nil {return tableEndItr}
	// return i.itrCache.add(&obj)
}
