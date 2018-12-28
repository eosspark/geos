package chain

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/eos_math"
	"github.com/eosspark/eos-go/entity"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
)

type IdxLongDouble struct {
	context  *ApplyContext
	itrCache *iteratorCache
}

func NewIdxLongDouble(c *ApplyContext) *IdxLongDouble {
	return &IdxLongDouble{
		context:  c,
		itrCache: NewIteratorCache(),
	}
}

func (i *IdxLongDouble) store(scope uint64, table uint64, payer uint64, id uint64, secondary *eos_math.Float128) int {

	EosAssert(common.AccountName(payer) != common.AccountName(0), &InvalidTablePayer{}, "must specify a valid account to pay for new record")
	tab := i.context.FindOrCreateTable(uint64(i.context.Receiver), scope, table, payer)

	obj := entity.IdxLongDoubleObject{
		TId:          tab.ID,
		PrimaryKey:   uint64(id),
		SecondaryKey: *secondary,
		Payer:        common.AccountName(payer),
	}

	i.context.DB.Insert(&obj)
	i.context.DB.Modify(tab, func(t *entity.TableIdObject) {
		t.Count++
	})

	i.context.UpdateDbUsage(common.AccountName(payer), int64(common.BillableSizeV("IdxLongDoubleObject")))

	i.itrCache.cacheTable(tab)
	iteratorOut := i.itrCache.add(&obj)
	i.context.ilog.Debug("object:%v iteratorOut:%d code:%v scope:%v table:%v payer:%v id:%d secondary:%v",
		obj, iteratorOut, i.context.Receiver, common.ScopeName(scope), common.TableName(table), common.AccountName(payer), id, *secondary)
	return iteratorOut

}

func (i *IdxLongDouble) remove(iterator int) {

	obj := (i.itrCache.get(iterator)).(*entity.IdxLongDoubleObject)
	i.context.UpdateDbUsage(obj.Payer, -int64(common.BillableSizeV("IdxLongDoubleObject")))

	tab := i.itrCache.getTable(obj.TId)
	EosAssert(tab.Code == i.context.Receiver, &TableAccessViolation{}, "db access violation")

	i.context.DB.Modify(tab, func(t *entity.TableIdObject) {
		t.Count--
	})

	i.context.ilog.Debug("object:%v iterator:%d ", *obj, iterator)

	i.context.DB.Remove(obj)
	if tab.Count == 0 {
		i.context.DB.Remove(tab)
	}
	i.itrCache.remove(iterator)

}

func (i *IdxLongDouble) update(iterator int, payer uint64, secondary *eos_math.Float128) {

	obj := (i.itrCache.get(iterator)).(*entity.IdxLongDoubleObject)
	objTable := i.itrCache.getTable(obj.TId)
	i.context.ilog.Debug("object:%v iterator:%d payer:%v secondary:%d", *obj, iterator, common.AccountName(payer), *secondary)
	EosAssert(objTable.Code == i.context.Receiver, &TableAccessViolation{}, "db access violation")

	accountPayer := common.AccountName(payer)
	if accountPayer == common.AccountName(0) {
		accountPayer = obj.Payer
	}

	billingSize := int64(common.BillableSizeV("IdxLongDoubleObject"))
	if obj.Payer != accountPayer {
		i.context.UpdateDbUsage(obj.Payer, -billingSize)
		i.context.UpdateDbUsage(accountPayer, +billingSize)
	}

	i.context.DB.Modify(obj, func(o *entity.IdxLongDoubleObject) {
		o.SecondaryKey = *secondary
		o.Payer = accountPayer
	})

}

func (i *IdxLongDouble) findSecondary(code uint64, scope uint64, table uint64, secondary *eos_math.Float128, primary *uint64) int {

	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}

	tableEndItr := i.itrCache.cacheTable(tab)

	obj := entity.IdxLongDoubleObject{TId: tab.ID, SecondaryKey: *secondary}
	err := i.context.DB.Find("bySecondary", obj, &obj)

	if err != nil {
		return tableEndItr
	}

	*primary = obj.PrimaryKey
	iteratorOut := i.itrCache.add(&obj)
	i.context.ilog.Debug("object:%v iteratorOut:%d code:%v scope:%v table:%v secondary:%d",
		obj, iteratorOut, common.AccountName(code), common.ScopeName(scope), common.TableName(table), *secondary)
	return iteratorOut
}

func (i *IdxLongDouble) lowerbound(code uint64, scope uint64, table uint64, secondary *eos_math.Float128, primary *uint64) int {

	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}

	tableEndItr := i.itrCache.cacheTable(tab)

	//for test
	//if *secondary == 0 {
	//	*secondary = 1
	//}

	obj := entity.IdxLongDoubleObject{TId: tab.ID, SecondaryKey: *secondary}

	idx, _ := i.context.DB.GetIndex("bySecondary", &obj)
	itr, _ := idx.LowerBound(&obj)
	if idx.CompareEnd(itr) {
		return tableEndItr
	}

	i.context.ilog.Debug("secondary:%v", *secondary)

	objLowerbound := entity.IdxLongDoubleObject{}
	itr.Data(&objLowerbound)
	if objLowerbound.TId != tab.ID {
		return tableEndItr
	}

	*primary = objLowerbound.PrimaryKey
	*secondary = objLowerbound.SecondaryKey

	iteratorOut := i.itrCache.add(&objLowerbound)
	i.context.ilog.Debug("object:%v iteratorOut:%d code:%v scope:%v table:%v secondary:%v",
		objLowerbound, iteratorOut, common.AccountName(code), common.ScopeName(scope), common.TableName(table), *secondary)
	return iteratorOut
}

func (i *IdxLongDouble) upperbound(code uint64, scope uint64, table uint64, secondary *eos_math.Float128, primary *uint64) int {

	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}

	tableEndItr := i.itrCache.cacheTable(tab)

	obj := entity.IdxLongDoubleObject{TId: tab.ID, SecondaryKey: *secondary}

	idx, _ := i.context.DB.GetIndex("bySecondary", &obj)
	itr, _ := idx.UpperBound(&obj)
	if idx.CompareEnd(itr) {
		return tableEndItr
	}

	objUpperbound := entity.IdxLongDoubleObject{}
	itr.Data(&objUpperbound)
	if objUpperbound.TId != tab.ID {
		return tableEndItr
	}

	*primary = objUpperbound.PrimaryKey
	*secondary = objUpperbound.SecondaryKey

	iteratorOut := i.itrCache.add(&objUpperbound)
	i.context.ilog.Debug("object:%v iteratorOut:%d code:%v scope:%v table:%v secondary:%d",
		objUpperbound, iteratorOut, common.AccountName(code), common.ScopeName(scope), common.TableName(table), *secondary)
	return iteratorOut
}

func (i *IdxLongDouble) end(code uint64, scope uint64, table uint64) int {

	i.context.ilog.Info("code:%v scope:%v table:%v ",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table))

	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}
	return i.itrCache.cacheTable(tab)
}

func (i *IdxLongDouble) next(iterator int, primary *uint64) int {

	if iterator < -1 {
		return -1
	}
	obj := (i.itrCache.get(iterator)).(*entity.IdxLongDoubleObject)

	idx, _ := i.context.DB.GetIndex("bySecondary", obj)
	itr := idx.IteratorTo(obj)

	itr.Next()
	objNext := entity.IdxLongDoubleObject{}
	itr.Data(&objNext)

	i.context.ilog.Debug("IdxLongDoubleObject objNext:%v", objNext)

	if idx.CompareEnd(itr) || objNext.TId != obj.TId {
		return i.itrCache.getEndIteratorByTableID(obj.TId)
	}

	*primary = objNext.PrimaryKey

	iteratorOut := i.itrCache.add(&objNext)
	i.context.ilog.Debug("object:%v iteratorIn:%d iteratorOut:%d", objNext, iterator, iteratorOut)
	return iteratorOut

}

func (i *IdxLongDouble) previous(iterator int, primary *uint64) int {

	idx, _ := i.context.DB.GetIndex("bySecondary", &entity.IdxLongDoubleObject{})

	if iterator < -1 {
		tab := i.itrCache.findTablebyEndIterator(iterator)
		EosAssert(tab != nil, &InvalidTableIterator{}, "not a valid end iterator")

		objTId := entity.IdxLongDoubleObject{TId: tab.ID}

		itr, _ := idx.UpperBound(&objTId)
		if idx.CompareIterator(idx.Begin(), idx.End()) || idx.CompareBegin(itr) {
			i.context.ilog.Info("iterator is the begin(nil) of index, iteratorIn:%d iteratorOut:%d", iterator, -1)
			return -1
		}

		itr.Prev()
		objPrev := entity.KeyValueObject{}
		itr.Data(&objPrev)

		if objPrev.TId != tab.ID {
			i.context.ilog.Info("previous iterator out of tid, iteratorIn:%d iteratorOut:%d", iterator, -1)
			return -1
		}

		*primary = objPrev.PrimaryKey
		iteratorOut := i.itrCache.add(&objPrev)
		i.context.ilog.Info("object:%v iteratorIn:%d iteratorOut:%d", objPrev, iterator, iteratorOut)
		return iteratorOut
	}

	obj := (i.itrCache.get(iterator)).(*entity.IdxLongDoubleObject)
	itr := idx.IteratorTo(obj)

	if idx.CompareBegin(itr) {
		return -1
	}

	itr.Prev()
	objPrev := entity.IdxLongDoubleObject{}
	itr.Data(&objPrev)
	i.context.ilog.Debug("IdxLongDoubleObject objPrev:%v", objPrev)
	if objPrev.TId != obj.TId {
		return -1
	}
	*primary = objPrev.PrimaryKey

	iteratorOut := i.itrCache.add(&objPrev)
	i.context.ilog.Debug("object:%v secondaryKey:%v iteratorIn:%d iteratorOut:%d", objPrev, objPrev.SecondaryKey, iterator, iteratorOut)
	return iteratorOut
}

func (i *IdxLongDouble) findPrimary(code uint64, scope uint64, table uint64, secondary *eos_math.Float128, primary uint64) int {

	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}

	tableEndItr := i.itrCache.cacheTable(tab)

	obj := entity.IdxLongDoubleObject{TId: tab.ID, PrimaryKey: primary}
	err := i.context.DB.Find("byPrimary", obj, &obj)
	if err != nil {
		return tableEndItr
	}

	*secondary = obj.SecondaryKey

	iteratorOut := i.itrCache.add(&obj)
	i.context.ilog.Debug("object:%v iteratorOut:%d code:%v scope:%v table:%v secondary:%d primary:%d ",
		obj, iteratorOut, common.AccountName(code), common.ScopeName(scope), common.TableName(table), *secondary, primary)
	return iteratorOut
}

func (i *IdxLongDouble) lowerboundPrimary(code uint64, scope uint64, table uint64, primary *uint64) int {

	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}

	tableEndItr := i.itrCache.cacheTable(tab)

	obj := entity.IdxLongDoubleObject{TId: tab.ID, PrimaryKey: *primary}
	idx, _ := i.context.DB.GetIndex("byPrimary", &obj)

	itr, _ := idx.LowerBound(&obj)
	if idx.CompareEnd(itr) {
		return tableEndItr
	}

	objLowerbound := entity.IdxLongDoubleObject{}
	itr.Data(&objLowerbound)

	if objLowerbound.TId != tab.ID {
		return tableEndItr
	}

	return i.itrCache.add(&objLowerbound)
}

func (i *IdxLongDouble) upperboundPrimary(code uint64, scope uint64, table uint64, primary *uint64) int {

	tab := i.context.FindTable(code, scope, table)
	if tab == nil {
		return -1
	}

	tableEndItr := i.itrCache.cacheTable(tab)

	obj := entity.IdxLongDoubleObject{TId: tab.ID, PrimaryKey: *primary}
	idx, _ := i.context.DB.GetIndex("byPrimary", &obj)
	itr, _ := idx.UpperBound(&obj)
	if idx.CompareEnd(itr) {
		return tableEndItr
	}
	//objUpperbound := (*types.IdxLongDouble)(itr.GetObject())
	objUpperbound := entity.IdxLongDoubleObject{}
	itr.Data(&objUpperbound)
	if objUpperbound.TId != tab.ID {
		return tableEndItr
	}

	i.itrCache.cacheTable(tab)
	return i.itrCache.add(&objUpperbound)
}

func (i *IdxLongDouble) nextPrimary(iterator int, primary *uint64) int {

	if iterator < -1 {
		return -1
	}
	obj := (i.itrCache.get(iterator)).(*entity.IdxLongDoubleObject)
	idx, _ := i.context.DB.GetIndex("byPrimary", obj)

	itr := idx.IteratorTo(obj)

	itr.Next()
	objNext := entity.IdxLongDoubleObject{}
	itr.Data(&objNext)

	if idx.CompareEnd(itr) || objNext.TId != obj.TId {
		return i.itrCache.getEndIteratorByTableID(obj.TId)
	}

	*primary = objNext.PrimaryKey
	return i.itrCache.add(&objNext)

}

func (i *IdxLongDouble) previousPrimary(iterator int, primary *uint64) int {

	idx, _ := i.context.DB.GetIndex("byPrimary", &entity.IdxLongDoubleObject{})

	if iterator < -1 {
		tab := i.itrCache.findTablebyEndIterator(iterator)
		EosAssert(tab != nil, &InvalidTableIterator{}, "not a valid end iterator")

		objTId := entity.IdxLongDoubleObject{TId: tab.ID}

		itr, _ := idx.UpperBound(&objTId)
		if idx.CompareIterator(idx.Begin(), idx.End()) || idx.CompareBegin(itr) {
			return -1
		}

		itr.Prev()
		objPrev := entity.IdxLongDoubleObject{}
		itr.Data(&objPrev)

		if objPrev.TId != tab.ID {
			return -1
		}

		*primary = objPrev.PrimaryKey
		return i.itrCache.add(&objPrev)
	}

	obj := (i.itrCache.get(iterator)).(*entity.IdxLongDoubleObject)
	itr := idx.IteratorTo(obj)

	if idx.CompareBegin(itr) {
		return -1
	}

	itr.Prev()
	objNext := entity.IdxLongDoubleObject{}
	itr.Data(&objNext)

	if objNext.TId != obj.TId {
		return -1
	}
	*primary = objNext.PrimaryKey
	return i.itrCache.add(&objNext)
}

func (i *IdxLongDouble) get(iterator int, secondary *eos_math.Float128, primary *uint64) {
	obj := (i.itrCache.get(iterator)).(*entity.IdxLongDoubleObject)

	*primary = obj.PrimaryKey
	*secondary = obj.SecondaryKey
}
