package wasmgo

import (
	"fmt"
	arithmetic "github.com/eosspark/eos-go/common/arithmetic_types"
	//"github.com/eosspark/eos-go/chain/types"
	//"github.com/eosspark/eos-go/entity"
	//"github.com/eosspark/eos-go/common"
)

// int db_store_i64( uint64_t scope, uint64_t table, uint64_t payer, uint64_t id, array_ptr<const char> buffer, size_t buffer_size ) {
//    return context.db_store_i64( scope, table, payer, id, buffer, buffer_size );
// }
func dbStoreI64(w *WasmGo, scope uint64, table uint64, payer uint64, id uint64, buffer int, bufferSize int) int {
	//fmt.Println("db_store_i64")

	bytes := getMemory(w, buffer, bufferSize)
	return w.context.DbStoreI64(scope, table, payer, id, bytes)

}

// func dbStoreI64(w *WasmGo, scope uint64, table uint64, payer uint64, id uint64, buffer int) int {
// 	fmt.Println("db_store_i64")

// 	bytes := getUint64(w, buffer)
// 	return w.context.DbStoreI64(scope, table, payer, id, bytes)
// }

// void db_update_i64( int itr, uint64_t payer, array_ptr<const char> buffer, size_t buffer_size ) {
//    context.db_update_i64( itr, payer, buffer, buffer_size );
// }
func dbUpdateI64(w *WasmGo, itr int, payer uint64, buffer int, bufferSize int) {
	//fmt.Println("db_update_i64")

	bytes := getMemory(w, buffer, bufferSize)
	w.context.DbUpdateI64(itr, payer, bytes)
}

// void db_remove_i64( int itr ) {
//    context.db_remove_i64( itr );
// }
func dbRemoveI64(w *WasmGo, itr int) {
	//fmt.Println("db_remove_i64")

	w.context.DbRemoveI64(itr)
}

// int db_get_i64( int itr, array_ptr<char> buffer, size_t buffer_size ) {
//    return context.db_get_i64( itr, buffer, buffer_size );
// }
func dbGetI64(w *WasmGo, itr int, buffer int, bufferSize int) int {
	//fmt.Println("db_get_i64")

	bytes := make([]byte, bufferSize)
	size := w.context.DbGetI64(itr, bytes, bufferSize)
	if bufferSize == 0 {
		return size
	}
	setMemory(w, buffer, bytes, 0, size)
	return size
}

// int db_next_i64( int itr, uint64_t& primary ) {
//    return context.db_next_i64(itr, primary);
// }
func dbNextI64(w *WasmGo, itr int, primary int) int {
	//fmt.Println("db_next_i64")

	var p uint64
	iterator := w.context.DbNextI64(itr, &p)
	w.ilog.Info("dbNextI64 iterator:%d", iterator)

	if iterator <= -1 {
		return iterator
	}
	setUint64(w, primary, p)
	return iterator
}

// int db_previous_i64( int itr, uint64_t& primary ) {
//    return context.db_previous_i64(itr, primary);
// }
func dbPreviousI64(w *WasmGo, itr int, primary int) int {
	//fmt.Println("db_previous_i64")

	var p uint64
	iterator := w.context.DbPreviousI64(itr, &p)
	w.ilog.Info("dbNextI64 iterator:%d", iterator)
	if iterator <= -1 {
		return iterator
	}
	setUint64(w, primary, p)
	return iterator
}

// int db_find_i64( uint64_t code, uint64_t scope, uint64_t table, uint64_t id ) {
//    return context.db_find_i64( code, scope, table, id );
// }
func dbFindI64(w *WasmGo, code uint64, scope uint64, table uint64, id uint64) int {
	//fmt.Println("db_find_i64")
	return w.context.DbFindI64(code, scope, table, id)
}

// int db_lowerbound_i64( uint64_t code, uint64_t scope, uint64_t table, uint64_t id ) {
//    return context.db_lowerbound_i64( code, scope, table, id );
// }
func dbLowerboundI64(w *WasmGo, code uint64, scope uint64, table uint64, id uint64) int {
	//fmt.Println("db_lowerbound_i64")
	return w.context.DbLowerboundI64(code, scope, table, id)
}

// int db_upperbound_i64( uint64_t code, uint64_t scope, uint64_t table, uint64_t id ) {
//    return context.db_upperbound_i64( code, scope, table, id );
// }
func dbUpperboundI64(w *WasmGo, code uint64, scope uint64, table uint64, id uint64) int {
	//fmt.Println("db_upperbound_i64")
	return w.context.DbUpperboundI64(code, scope, table, id)
}

// int db_end_i64( uint64_t code, uint64_t scope, uint64_t table ) {
//    return context.db_end_i64( code, scope, table );
// }
func dbEndI64(w *WasmGo, code uint64, scope uint64, table uint64) int {
	//fmt.Println("db_end_i64")
	return w.context.DbEndI64(code, scope, table)
}

//secondaryKey Index
func dbIdx64Store(w *WasmGo, scope uint64, table uint64, payer uint64, id uint64, pValue int) int {
	//fmt.Println("db_idx64_store")

	secondaryKey := getUint64(w, pValue)
	return w.context.Idx64Store(scope, table, payer, id, &secondaryKey)
}

func dbIdx64Remove(w *WasmGo, itr int) {
	//fmt.Println("db_idx64_remove")
	w.context.Idx64Remove(itr)
}

func dbIdx64Update(w *WasmGo, itr int, payer uint64, pValue int) {
	//fmt.Println("db_idx64_update")

	secondaryKey := getUint64(w, pValue)
	w.context.Idx64Update(itr, payer, &secondaryKey)
}

func dbIdx64findSecondary(w *WasmGo, code uint64, scope uint64, table uint64, pSecondary int, pPrimary int) int {

	//fmt.Println("db_idx64_find_secondary")

	var primaryKey uint64
	secondaryKey := getUint64(w, pSecondary)
	iterator := w.context.Idx64FindSecondary(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		return iterator
	}
	setUint64(w, pPrimary, primaryKey)

	return iterator
}

func dbIdx64Lowerbound(w *WasmGo, code uint64, scope uint64, table uint64, pSecondary int, pPrimary int) int {

	//fmt.Println("db_idx64_lowerbound")

	var primaryKey uint64

	secondaryKey := getUint64(w, pSecondary)
	iterator := w.context.Idx64Lowerbound(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		return iterator
	}
	setUint64(w, pPrimary, primaryKey)
	setUint64(w, pSecondary, secondaryKey)

	return iterator
}

func dbIdx64Upperbound(w *WasmGo, code uint64, scope uint64, table uint64, pSecondary int, pPrimary int) int {

	//fmt.Println("db_idx64_upperbound")

	var primaryKey uint64
	secondaryKey := getUint64(w, pSecondary)
	iterator := w.context.Idx64Upperbound(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		return iterator
	}
	setUint64(w, pPrimary, primaryKey)
	setUint64(w, pSecondary, secondaryKey)

	return iterator
}

func dbIdx64End(w *WasmGo, code uint64, scope uint64, table uint64) int {

	fmt.Println("db_idx64_end")

	return w.context.Idx64End(code, scope, table)
}

func dbIdx64Next(w *WasmGo, itr int, primary int) int {
	//fmt.Println("db_idx64_next")

	var p uint64
	iterator := w.context.Idx64Next(itr, &p)
	w.ilog.Info("dbIdx64Next iterator:%d", iterator)
	if iterator <= -1 {
		return iterator
	}
	setUint64(w, primary, p)

	return iterator
}

func dbIdx64Previous(w *WasmGo, itr int, primary int) int {
	//fmt.Println("db_idx64_previous")

	var p uint64
	iterator := w.context.Idx64Previous(itr, &p)
	w.ilog.Info("dbIdx64Previous iterator:%d", iterator)
	if iterator <= -1 {
		return iterator
	}
	setUint64(w, primary, p)

	return iterator
}

func dbIdx64FindPrimary(w *WasmGo, code uint64, scope uint64, table uint64, pSecondary int, primary uint64) int {

	//fmt.Println("db_idx64_find_primary")

	//primaryKey := getUint64(w, pPrimary)
	var secondaryKey uint64
	iterator := w.context.Idx64FindPrimary(code, scope, table, &secondaryKey, primary)
	if iterator <= -1 {
		return iterator
	}
	setUint64(w, pSecondary, secondaryKey)

	return iterator
}

func dbIdxDoubleStore(w *WasmGo, scope uint64, table uint64, payer uint64, id uint64, pValue int) int {
	//fmt.Println("db_double_store")

	secondaryKey := arithmetic.Float64(getUint64(w, pValue))
	//float := math.Float64frombits(getUint64(w, pValue))
	//w.ilog.Info("float:%v", float)

	return w.context.IdxDoubleStore(scope, table, payer, id, &secondaryKey)
}

func dbIdxDoubleRemove(w *WasmGo, itr int) {
	fmt.Println("db_double_remove")
	w.context.IdxDoubleRemove(itr)
}

func dbIdxDoubleUpdate(w *WasmGo, itr int, payer uint64, pValue int) {
	//fmt.Println("db_double_update")

	secondaryKey := arithmetic.Float64(getUint64(w, pValue))
	w.context.IdxDoubleUpdate(itr, payer, &secondaryKey)
}

func dbIdxDoublefindSecondary(w *WasmGo, code uint64, scope uint64, table uint64, pSecondary int, pPrimary int) int {

	//fmt.Println("db_double_find_secondary")

	var primaryKey uint64
	secondaryKey := arithmetic.Float64(getUint64(w, pSecondary))
	iterator := w.context.IdxDoubleFindSecondary(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		return iterator
	}
	setUint64(w, pPrimary, primaryKey)

	return iterator
}

func dbIdxDoubleLowerbound(w *WasmGo, code uint64, scope uint64, table uint64, pSecondary int, pPrimary int) int {

	//fmt.Println("db_double_lowerbound")

	var primaryKey uint64
	secondaryKey := arithmetic.Float64(getUint64(w, pSecondary))
	iterator := w.context.IdxDoubleLowerbound(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		return iterator
	}
	setUint64(w, pPrimary, primaryKey)
	setUint64(w, pSecondary, uint64(secondaryKey))

	return iterator
}

func dbIdxDoubleUpperbound(w *WasmGo, code uint64, scope uint64, table uint64, pSecondary int, pPrimary int) int {

	//fmt.Println("db_double_upperbound")

	var primaryKey uint64
	secondaryKey := arithmetic.Float64(getUint64(w, pSecondary))
	iterator := w.context.IdxDoubleUpperbound(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		return iterator
	}
	setUint64(w, pPrimary, primaryKey)
	setUint64(w, pSecondary, uint64(secondaryKey))

	return iterator
}

func dbIdxDoubleEnd(w *WasmGo, code uint64, scope uint64, table uint64) int {

	fmt.Println("db_double_end")

	return w.context.IdxDoubleEnd(code, scope, table)
}

func dbIdxDoubleNext(w *WasmGo, itr int, primary int) int {
	fmt.Println("db_double_next")

	var p uint64
	iterator := w.context.IdxDoubleNext(itr, &p)
	w.ilog.Info("dbIdxDoubleNext iterator:%d", iterator)
	if iterator <= -1 {
		return iterator
	}

	setUint64(w, primary, p)

	return iterator
}

func dbIdxDoublePrevious(w *WasmGo, itr int, primary int) int {
	fmt.Println("db_double_previous")

	var p uint64
	iterator := w.context.IdxDoublePrevious(itr, &p)
	w.ilog.Info("dbIdxDoublePrevious iterator:%d", iterator)
	if iterator <= -1 {
		return iterator
	}
	setUint64(w, primary, p)

	return iterator
}

func dbIdxDoubleFindPrimary(w *WasmGo, code uint64, scope uint64, table uint64, pSecondary int, primary uint64) int {

	//fmt.Println("db_double_find_primary")

	// primaryKey := getUint64(w, pPrimary)
	var secondaryKey arithmetic.Float64
	iterator := w.context.IdxDoubleFindPrimary(code, scope, table, &secondaryKey, primary)
	if iterator <= -1 {
		return iterator
	}
	setUint64(w, pSecondary, uint64(secondaryKey))

	return iterator
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
