package exec

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

// int db_store_i64( uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, array_ptr<const char> buffer, size_t buffer_size ) {
//    return context.db_store_i64( scope, table, payer, id, buffer, buffer_size );
// }
func db_store_i64(w *WasmInterface, scope int64, table int64, payer int64, id int64, buffer int, bufferSize int) int {
	fmt.Println("db_store_i64")

	bytes := getMemory(w, buffer, bufferSize)
	return w.context.DBStoreI64(scope, table, payer, id, bytes)

}

// void db_update_i64( int itr, uint64_t payer, array_ptr<const char> buffer, size_t buffer_size ) {
//    context.db_update_i64( itr, payer, buffer, buffer_size );
// }
func db_update_i64(w *WasmInterface,
	itr int, payer int64,
	buffer int, bufferSize int) {
	fmt.Println("db_update_i64")

	bytes := getMemory(w, buffer, bufferSize)
	w.context.DBUpdateI64(itr, common.AccountName(payer), bytes)
}

// void db_remove_i64( int itr ) {
//    context.db_remove_i64( itr );
// }
func db_remove_i64(w *WasmInterface, itr int) {
	fmt.Println("db_update_i64")

	w.context.DBRemoveI64(itr)
}

// int db_get_i64( int itr, array_ptr<char> buffer, size_t buffer_size ) {
//    return context.db_get_i64( itr, buffer, buffer_size );
// }
func db_get_i64(w *WasmInterface, itr int, buffer int, bufferSize int) int {
	fmt.Println("db_get_i64")

	bytes := make([]byte, bufferSize)
	return w.context.DBGetI64(itr, bytes, bufferSize)
}

// int db_next_i64( int itr, uint64_t& primary ) {
//    return context.db_next_i64(itr, primary);
// }
func db_next_i64(w *WasmInterface, itr int, primary int) int {
	fmt.Println("db_next_i64")

	var p uint64

	iterator := w.context.DBNextI64(itr, &p)
	setUint64(w, primary, p)

	return iterator
}

// int db_previous_i64( int itr, uint64_t& primary ) {
//    return context.db_previous_i64(itr, primary);
// }
func db_previous_i64(w *WasmInterface, itr int, primary int) int {
	fmt.Println("db_previous_i64")

	var p uint64

	iterator := w.context.DBPreviousI64(itr, &p)
	setUint64(w, primary, p)
	return iterator
}

// int db_find_i64( uint64_t code, uint64_t scope, uint64_t table, uint64_t id ) {
//    return context.db_find_i64( code, scope, table, id );
// }
func db_find_i64(w *WasmInterface, code int64, scope int64, table int64, id int64) int {
	fmt.Println("db_find_i64")
	return w.context.DBFindI64(code, scope, table, id)
}

// int db_lowerbound_i64( uint64_t code, uint64_t scope, uint64_t table, uint64_t id ) {
//    return context.db_lowerbound_i64( code, scope, table, id );
// }
func db_lowerbound_i64(w *WasmInterface, code int64, scope int64, table int64, id int64) int {
	fmt.Println("db_lowerbound_i64")

	return w.context.DBLowerBoundI64(code, scope, table, id)
}

// int db_upperbound_i64( uint64_t code, uint64_t scope, uint64_t table, uint64_t id ) {
//    return context.db_upperbound_i64( code, scope, table, id );
// }
func db_upperbound_i64(w *WasmInterface, code int64, scope int64, table int64, id int64) int {
	fmt.Println("db_upperbound_i64")
	return w.context.DBUpperBoundI64(code, scope, table, id)
}

// int db_end_i64( uint64_t code, uint64_t scope, uint64_t table ) {
//    return context.db_end_i64( code, scope, table );
// }
func db_end_i64(w *WasmInterface, code int64, scope int64, table int64) int {
	fmt.Println("db_end_i64")
	return w.context.DBEndI64(code, scope, table)
}

//secondaryKey Index
func db_idx64_store(w *WasmInterface, scope int64, table int64, payer int64, id int64, pValue int) int {
	fmt.Println("db_idx64_store")

	secondaryKey := &types.Uint64_t{Value: getUint64(w, pValue)}
	//secondaryKey.SetValue(getUint64(w, pValue))
	return w.context.IdxI64Store(scope, table, payer, id, secondaryKey)
}

func db_idx64_remove(w *WasmInterface, itr int) {
	fmt.Println("db_update_i64")
	w.context.IdxI64Remove(itr)
}

func db_idx64_update(w *WasmInterface, itr int, payer int64, pValue int) {
	fmt.Println("db_update_i64")

	secondaryKey := &types.Uint64_t{Value: getUint64(w, pValue)}
	//secondaryKey.SetValue(getUint64(w, pValue))
	w.context.IdxI64Update(itr, payer, secondaryKey)
}

func db_idx64_find_secondary(w *WasmInterface, code int64, scope int64, table int64, payer int64, pSecondary int, pPrimary int) {

	fmt.Println("db_update_i64")

	primaryKey := getUint64(w, pPrimary)
	secondaryKey := &types.Uint64_t{Value: getUint64(w, pSecondary)}

	w.context.IdxI64FindSecondary(code, scope, table, secondaryKey, &primaryKey)

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
