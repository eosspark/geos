package exec

import (
	"fmt"
)

// int db_store_i64( uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, array_ptr<const char> buffer, size_t buffer_size ) {
//    return context.db_store_i64( scope, table, payer, id, buffer, buffer_size );
// }
func db_store_i64(w *WasmInterface,
	scope int64, table int64, payer int64, id int64,
	buffer int, buffer_size int) int {
	fmt.Println("db_store_i64")

	return w.context.DBStoreI64(scope, table, payer, id, buffer, buffer_size)

}

// void db_update_i64( int itr, uint64_t payer, array_ptr<const char> buffer, size_t buffer_size ) {
//    context.db_update_i64( itr, payer, buffer, buffer_size );
// }
func db_update_i64(w *WasmInterface,
	itr int, payer int64,
	buffer int, buffer_size int) {
	fmt.Println("db_update_i64")
}

// void db_remove_i64( int itr ) {
//    context.db_remove_i64( itr );
// }
func db_remove_i64(w *WasmInterface, itr int) {
	fmt.Println("db_update_i64")
}

// int db_get_i64( int itr, array_ptr<char> buffer, size_t buffer_size ) {
//    return context.db_get_i64( itr, buffer, buffer_size );
// }
func db_get_i64(w *WasmInterface, itr int, buffer int, buffer_size int) int {
	fmt.Println("db_get_i64")
	return 0
}

// int db_next_i64( int itr, uint64_t& primary ) {
//    return context.db_next_i64(itr, primary);
// }
func db_next_i64(w *WasmInterface, itr int, primary int) int {
	fmt.Println("db_next_i64")
	return 0
}

// int db_previous_i64( int itr, uint64_t& primary ) {
//    return context.db_previous_i64(itr, primary);
// }
func db_previous_i64(w *WasmInterface, itr int, primary int) int {
	fmt.Println("db_previous_i64")
	return 0
}

// int db_find_i64( uint64_t code, uint64_t scope, uint64_t table, uint64_t id ) {
//    return context.db_find_i64( code, scope, table, id );
// }
func db_find_i64(w *WasmInterface, code int64, scope int64, table int64, id int64) int {
	fmt.Println("db_find_i64")
	return 0
}

// int db_lowerbound_i64( uint64_t code, uint64_t scope, uint64_t table, uint64_t id ) {
//    return context.db_lowerbound_i64( code, scope, table, id );
// }
func db_lowerbound_i64(w *WasmInterface, code int64, scope int64, table int64, id int64) int {
	fmt.Println("db_lowerbound_i64")

	return 0
}

// int db_upperbound_i64( uint64_t code, uint64_t scope, uint64_t table, uint64_t id ) {
//    return context.db_upperbound_i64( code, scope, table, id );
// }
func db_upperbound_i64(w *WasmInterface, code int64, scope int64, table int64, id int64) int {
	fmt.Println("db_upperbound_i64")
	return 0
}

// int db_end_i64( uint64_t code, uint64_t scope, uint64_t table ) {
//    return context.db_end_i64( code, scope, table );
// }
func db_end_i64(w *WasmInterface, code int64, scope int64, table int64) int {
	fmt.Println("db_end_i64")
	return 0
}

// DB_API_METHOD_WRAPPERS_SIMPLE_SECONDARY(idx64,  uint64_t)
// DB_API_METHOD_WRAPPERS_SIMPLE_SECONDARY(idx128, uint128_t)
// DB_API_METHOD_WRAPPERS_ARRAY_SECONDARY(idx256, 2, uint128_t)
// DB_API_METHOD_WRAPPERS_FLOAT_SECONDARY(idx_double, float64_t)
// DB_API_METHOD_WRAPPERS_FLOAT_SECONDARY(idx_long_double, float128_t)
