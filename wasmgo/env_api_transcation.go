package wasmgo

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	arithmetic "github.com/eosspark/eos-go/common/arithmetic_types"
	"github.com/eosspark/eos-go/crypto/rlp"
)

// void send_inline( array_ptr<char> data, size_t data_len ) {
//    //TODO: Why is this limit even needed? And why is it not consistently checked on actions in input or deferred transactions
//    EOS_ASSERT( data_len < context.control.get_global_properties().configuration.max_inline_action_size, inline_action_too_big,
//               "inline action too big" );

//    action act;
//    fc::raw::unpack<action>(data, data_len, act);
//    context.execute_inline(std::move(act));
// }
func sendInline(w *WasmGo, data int, dataLen int) {
	fmt.Println("send_inline")

	action := getBytes(w, data, dataLen)
	w.context.ExecuteInline(action)

}

// void send_context_free_inline( array_ptr<char> data, size_t data_len ) {
//    //TODO: Why is this limit even needed? And why is it not consistently checked on actions in input or deferred transactions
//    EOS_ASSERT( data_len < context.control.get_global_properties().configuration.max_inline_action_size, inline_action_too_big,
//              "inline action too big" );

//    action act;
//    fc::raw::unpack<action>(data, data_len, act);
//    context.execute_context_free_inline(std::move(act));
// }
func sendContextFreeInline(w *WasmGo, data int, dataLen int) {
	fmt.Println("send_context_free_inline")

	action := getBytes(w, data, dataLen)
	w.context.ExecuteContextFreeInline(action)
}

// void send_deferred( const uint128_t& sender_id, account_name payer, array_ptr<char> data, size_t data_len, uint32_t replace_existing) {
//    try {
//       transaction trx;
//       fc::raw::unpack<transaction>(data, data_len, trx);
//       context.schedule_deferred_transaction(sender_id, payer, std::move(trx), replace_existing);
//    } FC_RETHROW_EXCEPTIONS(warn, "data as hex: ${data}", ("data", fc::to_hex(data, data_len)))
// }
func sendDeferred(w *WasmGo, senderId int, payer common.AccountName, data int, dataLen int, replaceExisting int32) {
	fmt.Println("send_deferred")

	//id := big.Int.SetBytes(w.vm.memory[sender_id : sender_id+32])
	//id, _ := common.DecodeIdTypeByte(w.vm.memory[sender_id : sender_id+32])
	bytes := getMemory(w, senderId, 16)
	id := &arithmetic.Uint128{}
	rlp.DecodeBytes(bytes, id)

	trx := getBytes(w, data, dataLen)
	w.context.ScheduleDeferredTransaction(id, payer, trx, i2b(int(replaceExisting)))
}

// bool cancel_deferred( const unsigned __int128& val ) {
//    fc::uint128_t sender_id(val>>64, uint64_t(val) );
//    return context.cancel_deferred_transaction( (unsigned __int128)sender_id );
// }
func cancelDeferred(w *WasmGo, senderId int) int {
	fmt.Println("cancel_deferred")

	//id, _ := common.DecodeIdTypeByte(w.vm.memory[senderId : senderId+32])

	bytes := getMemory(w, senderId, 16)
	id := &arithmetic.Uint128{}
	rlp.DecodeBytes(bytes, id)

	//return b2i(w.context.CancelDeferredTransaction(common.TransactionIdType{id}))
	return b2i(w.context.CancelDeferredTransaction(id))

}

// int read_transaction( array_ptr<char> data, size_t buffer_size ) {
//    bytes trx = context.get_packed_transaction();

//    auto s = trx.size();
//    if( buffer_size == 0) return s;

//    auto copy_size = std::min( buffer_size, s );
//    memcpy( data, trx.data(), copy_size );

//    return copy_size;
// }
func readTransaction(w *WasmGo, buffer int, bufferSize int) int {
	fmt.Println("read_transaction")

	trx := w.context.GetPackedTransaction()

	s := len(trx)
	if bufferSize == 0 {
		return s
	}

	copySize := min(bufferSize, s)
	setMemory(w, buffer, trx, 0, copySize)
	return copySize
}

// int transaction_size() {
//    return context.get_packed_transaction().size();
// }
func transactionSize(w *WasmGo) int {
	fmt.Println("transaction_size")

	return len(w.context.GetPackedTransaction())
}

// int expiration() {
//   return context.trx_context.trx.expiration.sec_since_epoch();
// }
func expiration(w *WasmGo) int {
	fmt.Println("expiration")

	return w.context.Expiration()
}

// int tapos_block_num() {
//   return context.trx_context.trx.ref_block_num;
// }
func taposBlockNum(w *WasmGo) int {
	fmt.Println("tapos_block_num")

	return w.context.TaposBlockNum()
}

// int tapos_block_prefix() {
//   return context.trx_context.trx.ref_block_prefix;
// }
func taposBlockPrefix(w *WasmGo) int {
	fmt.Println("tapos_block_prefix")

	return w.context.TaposBlockPrefix()
}

// int get_action( uint32_t type, uint32_t index, array_ptr<char> buffer, size_t buffer_size )const {
//    return context.get_action( type, index, buffer, buffer_size );
// }
func getAction(w *WasmGo, typ int, index int, buffer int, bufferSize int) int {
	w.ilog.Info("getAction index:%d buffer:%d bufferSize:%d", index, buffer, bufferSize)

	s, action := w.context.GetAction(uint32(typ), index, bufferSize)
	if bufferSize == 0 || action == nil {
		return s
	}
	setMemory(w, buffer, action, 0, s)
	return s

}

// int get_context_free_data( uint32_t index, array_ptr<char> buffer, size_t buffer_size )const {
//  return context.get_context_free_data( index, buffer, buffer_size );
// }
func getContextFreeData(w *WasmGo, index int, buffer int, bufferSize int) int {

	fmt.Println("get_context_free_data")

	s, data := w.context.GetContextFreeData(index, bufferSize)
	if bufferSize == 0 || s == -1 {
		return s
	}
	setMemory(w, buffer, data, 0, s)
	return s

}
