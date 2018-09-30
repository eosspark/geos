package exec

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	//"github.com/eosspark/eos-go/common"
)

// int db_store_i64( uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, array_ptr<const char> buffer, size_t buffer_size ) {
//    return context.db_store_i64( scope, table, payer, id, buffer, buffer_size );
// }
func dbStoreI64(w *WasmInterface, scope int64, table int64, payer int64, id int64, buffer int, bufferSize int) int {
	fmt.Println("db_store_i64")

	bytes := getMemory(w, buffer, bufferSize)
	return w.context.DbStoreI64(scope, table, payer, id, bytes)

}

// void db_update_i64( int itr, uint64_t payer, array_ptr<const char> buffer, size_t buffer_size ) {
//    context.db_update_i64( itr, payer, buffer, buffer_size );
// }
func dbUpdateI64(w *WasmInterface, itr int, payer int64, buffer int, bufferSize int) {
	fmt.Println("db_update_i64")

	bytes := getMemory(w, buffer, bufferSize)
	w.context.DbUpdateI64(itr, payer, bytes)
}

// void db_remove_i64( int itr ) {
//    context.db_remove_i64( itr );
// }
func dbRemoveI64(w *WasmInterface, itr int) {
	fmt.Println("db_remove_i64")

	w.context.DbRemoveI64(itr)
}

// int db_get_i64( int itr, array_ptr<char> buffer, size_t buffer_size ) {
//    return context.db_get_i64( itr, buffer, buffer_size );
// }
func dbGetI64(w *WasmInterface, itr int, buffer int, bufferSize int) int {
	fmt.Println("db_get_i64")

	bytes := make([]byte, bufferSize)
	return w.context.DbGetI64(itr, bytes, bufferSize)
}

// int db_next_i64( int itr, uint64_t& primary ) {
//    return context.db_next_i64(itr, primary);
// }
func dbNextI64(w *WasmInterface, itr int, primary int) int {
	fmt.Println("db_next_i64")

	var p uint64

	iterator := w.context.DbNextI64(itr, &p)
	setUint64(w, primary, p)

	return iterator
}

// int db_previous_i64( int itr, uint64_t& primary ) {
//    return context.db_previous_i64(itr, primary);
// }
func dbPreviousI64(w *WasmInterface, itr int, primary int) int {
	fmt.Println("db_previous_i64")

	var p uint64

	iterator := w.context.DbPreviousI64(itr, &p)
	setUint64(w, primary, p)
	return iterator
}

// int db_find_i64( uint64_t code, uint64_t scope, uint64_t table, uint64_t id ) {
//    return context.db_find_i64( code, scope, table, id );
// }
func dbFindI64(w *WasmInterface, code int64, scope int64, table int64, id int64) int {
	fmt.Println("db_find_i64")
	return w.context.DbFindI64(code, scope, table, id)
}

// int db_lowerbound_i64( uint64_t code, uint64_t scope, uint64_t table, uint64_t id ) {
//    return context.db_lowerbound_i64( code, scope, table, id );
// }
func dbLowerboundI64(w *WasmInterface, code int64, scope int64, table int64, id int64) int {
	fmt.Println("db_lowerbound_i64")

	return w.context.DbLowerboundI64(code, scope, table, id)
}

// int db_upperbound_i64( uint64_t code, uint64_t scope, uint64_t table, uint64_t id ) {
//    return context.db_upperbound_i64( code, scope, table, id );
// }
func dbUpperboundI64(w *WasmInterface, code int64, scope int64, table int64, id int64) int {
	fmt.Println("db_upperbound_i64")
	return w.context.DbUpperboundI64(code, scope, table, id)
}

// int db_end_i64( uint64_t code, uint64_t scope, uint64_t table ) {
//    return context.db_end_i64( code, scope, table );
// }
func dbEndI64(w *WasmInterface, code int64, scope int64, table int64) int {
	fmt.Println("db_end_i64")
	return w.context.DbEndI64(code, scope, table)
}

//secondaryKey Index
func dbIdx64Store(w *WasmInterface, scope int64, table int64, payer int64, id int64, pValue int) int {
	fmt.Println("db_idx64_store")

	secondaryKey := &types.Uint64_t{Value: getUint64(w, pValue)}
	//secondaryKey.SetValue(getUint64(w, pValue))
	return w.context.Idx64Store(scope, table, payer, id, secondaryKey)
}

func dbIdx64Remove(w *WasmInterface, itr int) {
	fmt.Println("db_idx64_remove")
	w.context.Idx64Remove(itr)
}

func dbIdx64Update(w *WasmInterface, itr int, payer int64, pValue int) {
	fmt.Println("db_idx64_update")

	secondaryKey := &types.Uint64_t{Value: getUint64(w, pValue)}
	//secondaryKey.SetValue(getUint64(w, pValue))
	w.context.Idx64Update(itr, payer, secondaryKey)
}

func dbIdx64findSecondary(w *WasmInterface, code int64, scope int64, table int64, payer int64, pSecondary int, pPrimary int) int {

	fmt.Println("db_idx64_find_secondary")

	var primaryKey uint64 //:= getUint64(w, pPrimary)
	secondaryKey := &types.Uint64_t{Value: getUint64(w, pSecondary)}
	itr := w.context.Idx64FindSecondary(code, scope, table, secondaryKey, &primaryKey)
	setUint64(w, pPrimary, primaryKey)

	return itr
}

func dbIdx64Lowerbound(w *WasmInterface, code int64, scope int64, table int64, pSecondary int, pPrimary int) int {

	fmt.Println("db_idx64_lowerbound")

	var primaryKey uint64 //:= getUint64(w, pPrimary)
	secondaryKey := types.Uint64_t{}
	itr := w.context.Idx64Lowerbound(code, scope, table, &secondaryKey, &primaryKey)
	setUint64(w, pPrimary, primaryKey)
	setUint64(w, pSecondary, secondaryKey.Value)

	return itr
}

func dbIdx64Upperbound(w *WasmInterface, code int64, scope int64, table int64, pSecondary int, pPrimary int) int {

	fmt.Println("db_idx64_upperbound")

	var primaryKey uint64 //:= getUint64(w, pPrimary)
	secondaryKey := types.Uint64_t{}
	itr := w.context.Idx64Lowerbound(code, scope, table, &secondaryKey, &primaryKey)
	setUint64(w, pPrimary, primaryKey)
	setUint64(w, pSecondary, secondaryKey.Value)

	return itr
}

func dbIdx64End(w *WasmInterface, code int64, scope int64, table int64) int {

	fmt.Println("db_idx64_end")

	return w.context.Idx64End(code, scope, table)
}

func dbIdx64Next(w *WasmInterface, itr int, primary int) int {
	fmt.Println("db_idx64_next")

	var p uint64

	iterator := w.context.Idx64Next(itr, &p)
	setUint64(w, primary, p)

	return iterator
}

func dbIdx64Previous(w *WasmInterface, itr int, primary int) int {
	fmt.Println("db_idx64_previous")

	var p uint64

	iterator := w.context.Idx64Previous(itr, &p)
	setUint64(w, primary, p)

	return iterator
}

func dbIdx64FindPrimary(w *WasmInterface, code int64, scope int64, table int64, payer int64, pSecondary int, pPrimary int) int {

	fmt.Println("db_idx64_find_primary")

	//var primaryKey uint64 //:= getUint64(w, pPrimary)

	primaryKey := getUint64(w, pPrimary)
	secondaryKey := types.Uint64_t{}
	itr := w.context.Idx64FindPrimary(code, scope, table, &secondaryKey, &primaryKey)
	setUint64(w, pSecondary, secondaryKey.Value)

	return itr
}

func dbIdxDoubleStore(w *WasmInterface, scope int64, table int64, payer int64, id int64, pValue int) int {
	fmt.Println("db_idx64_store")

	secondaryKey := &types.Float64_t{Value: getFloat64(w, pValue)}
	//secondaryKey.SetValue(getFloat64(w, pValue))
	return w.context.IdxDoubleStore(scope, table, payer, id, secondaryKey)
}

func dbIdxDoubleRemove(w *WasmInterface, itr int) {
	fmt.Println("db_idx64_remove")
	w.context.IdxDoubleRemove(itr)
}

func dbIdxDoubleUpdate(w *WasmInterface, itr int, payer int64, pValue int) {
	fmt.Println("db_idx64_update")

	secondaryKey := &types.Float64_t{Value: getFloat64(w, pValue)}
	//secondaryKey.SetValue(getFloat64(w, pValue))
	w.context.IdxDoubleUpdate(itr, payer, secondaryKey)
}

func dbIdxDoublefindSecondary(w *WasmInterface, code int64, scope int64, table int64, payer int64, pSecondary int, pPrimary int) int {

	fmt.Println("db_idx64_find_secondary")

	var primaryKey uint64 //:= getFloat64(w, pPrimary)
	secondaryKey := &types.Float64_t{Value: getFloat64(w, pSecondary)}
	itr := w.context.IdxDoubleFindSecondary(code, scope, table, secondaryKey, &primaryKey)
	setUint64(w, pPrimary, primaryKey)

	return itr
}

func dbIdxDoubleLowerbound(w *WasmInterface, code int64, scope int64, table int64, pSecondary int, pPrimary int) int {

	fmt.Println("db_idx64_lowerbound")

	var primaryKey uint64 //:= getFloat64(w, pPrimary)
	secondaryKey := types.Float64_t{}
	itr := w.context.IdxDoubleLowerbound(code, scope, table, &secondaryKey, &primaryKey)
	setUint64(w, pPrimary, primaryKey)
	setFloat64(w, pSecondary, secondaryKey.Value)

	return itr
}

func dbIdxDoubleUpperbound(w *WasmInterface, code int64, scope int64, table int64, pSecondary int, pPrimary int) int {

	fmt.Println("db_idx64_upperbound")

	var primaryKey uint64 //:= getFloat64(w, pPrimary)
	secondaryKey := types.Float64_t{}
	itr := w.context.IdxDoubleLowerbound(code, scope, table, &secondaryKey, &primaryKey)
	setUint64(w, pPrimary, primaryKey)
	setFloat64(w, pSecondary, secondaryKey.Value)

	return itr
}

func dbIdxDoubleEnd(w *WasmInterface, code int64, scope int64, table int64) int {

	fmt.Println("db_idx64_end")

	return w.context.IdxDoubleEnd(code, scope, table)
}

func dbIdxDoubleNext(w *WasmInterface, itr int, primary int) int {
	fmt.Println("db_idx64_next")

	var p uint64

	iterator := w.context.IdxDoubleNext(itr, &p)
	setUint64(w, primary, p)

	return iterator
}

func dbIdxDoublePrevious(w *WasmInterface, itr int, primary int) int {
	fmt.Println("db_idx64_previous")

	var p uint64

	iterator := w.context.IdxDoublePrevious(itr, &p)
	setUint64(w, primary, p)

	return iterator
}

func dbIdxDoubleFindPrimary(w *WasmInterface, code int64, scope int64, table int64, payer int64, pSecondary int, pPrimary int) int {

	fmt.Println("db_idx64_find_primary")

	//var primaryKey uint64 //:= getFloat64(w, pPrimary)

	primaryKey := getUint64(w, pPrimary)
	secondaryKey := types.Float64_t{}
	itr := w.context.IdxDoubleFindPrimary(code, scope, table, &secondaryKey, &primaryKey)
	setFloat64(w, pSecondary, secondaryKey.Value)

	return itr
}

// (db_##IDX##_remove,         void(int))\
// (db_##IDX##_update,         void(int,int64_t,int))\
// (db_##IDX##_find_primary,   int(int64_t,int64_t,int64_t,int,int64_t))\
// (db_##IDX##_find_secondary, int(int64_t,int64_t,int64_t,int,int))\
// (db_##IDX##_lowerbound,     int(int64_t,int64_t,int64_t,int,int))\
// (db_##IDX##_upperbound,     int(int64_t,int64_t,int64_t,int,int))\
// (db_##IDX##_end,            int(int64_t,int64_t,int64_t))\
// (db_##IDX##_next,           int(int, int))\
// (db_##IDX##_previous,       int(int, int))

// DB_API_METHOD_WRAPPERS_SIMPLE_SECONDARY(idx64,  uint64_t)
// DB_API_METHOD_WRAPPERS_SIMPLE_SECONDARY(idx128, uint128_t)
// DB_API_METHOD_WRAPPERS_ARRAY_SECONDARY(idx256, 2, uint128_t)
// DB_API_METHOD_WRAPPERS_FLOAT_SECONDARY(idx_double, float64_t)
// DB_API_METHOD_WRAPPERS_FLOAT_SECONDARY(idx_long_double, float128_t)
