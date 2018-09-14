package exec

import (
	//	"errors"
	"bytes"
	"errors"
	"fmt"
	"log"
	"reflect"

	//"math"
	//"os"
	"strings"

	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/cvm/wasm"
)

// void send_inline( array_ptr<char> data, size_t data_len ) {
//    //TODO: Why is this limit even needed? And why is it not consistently checked on actions in input or deferred transactions
//    EOS_ASSERT( data_len < context.control.get_global_properties().configuration.max_inline_action_size, inline_action_too_big,
//               "inline action too big" );

//    action act;
//    fc::raw::unpack<action>(data, data_len, act);
//    context.execute_inline(std::move(act));
// }
func send_inline(wasmInterface *WasmInterface, data int, data_len size_t) {
	fmt.Println("send_inline")
}

// void send_context_free_inline( array_ptr<char> data, size_t data_len ) {
//    //TODO: Why is this limit even needed? And why is it not consistently checked on actions in input or deferred transactions
//    EOS_ASSERT( data_len < context.control.get_global_properties().configuration.max_inline_action_size, inline_action_too_big,
//              "inline action too big" );

//    action act;
//    fc::raw::unpack<action>(data, data_len, act);
//    context.execute_context_free_inline(std::move(act));
// }
func send_context_free_inline(wasmInterface *WasmInterface, data int, data_len size_t) {
	fmt.Println("send_context_free_inline")
}

// void send_deferred( const uint128_t& sender_id, account_name payer, array_ptr<char> data, size_t data_len, uint32_t replace_existing) {
//    try {
//       transaction trx;
//       fc::raw::unpack<transaction>(data, data_len, trx);
//       context.schedule_deferred_transaction(sender_id, payer, std::move(trx), replace_existing);
//    } FC_RETHROW_EXCEPTIONS(warn, "data as hex: ${data}", ("data", fc::to_hex(data, data_len)))
// }
func send_deferred(wasmInterface *WasmInterface, sender_id int, payer AccountName, data int, data_len size_t, replace_existing uint32) {
	fmt.Println("send_deferred")
}

// bool cancel_deferred( const unsigned __int128& val ) {
//    fc::uint128_t sender_id(val>>64, uint64_t(val) );
//    return context.cancel_deferred_transaction( (unsigned __int128)sender_id );
// }
func cancel_deferred(val int) {
	fmt.Println("cancel_deferred")
}

// int read_transaction( array_ptr<char> data, size_t buffer_size ) {
//    bytes trx = context.get_packed_transaction();

//    auto s = trx.size();
//    if( buffer_size == 0) return s;

//    auto copy_size = std::min( buffer_size, s );
//    memcpy( data, trx.data(), copy_size );

//    return copy_size;
// }
func read_transaction(data int, buffer_size size_t) int {
	fmt.Println("read_transaction")
}

// int transaction_size() {
//    return context.get_packed_transaction().size();
// }
func transaction_size() int {
	fmt.Println("transaction_size")
}

// int expiration() {
//   return context.trx_context.trx.expiration.sec_since_epoch();
// }
func expiration() int {
	fmt.Println("expiration")
}

// int tapos_block_num() {
//   return context.trx_context.trx.ref_block_num;
// }
func tapos_block_num() int {
	fmt.Println("tapos_block_num")
}

// int tapos_block_prefix() {
//   return context.trx_context.trx.ref_block_prefix;
// }
func tapos_block_prefix() int {
	fmt.Println("tapos_block_prefix")
}

// int get_action( uint32_t type, uint32_t index, array_ptr<char> buffer, size_t buffer_size )const {
//    return context.get_action( type, index, buffer, buffer_size );
// }
func get_action(typ int, index int, buffer int, buffer_size size_t) int {
	fmt.Println("get_action")
}
