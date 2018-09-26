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

var (
	ignore bool
)

// uint64_t current_time() {
//    return static_cast<uint64_t>( context.control.pending_block_time().time_since_epoch().count() );
// }
func current_time(wasmInterface *WasmInterface) int64 {
	fmt.Println("prints")
}

// uint64_t publication_time() {
//    return static_cast<uint64_t>( context.trx_context.published.time_since_epoch().count() );
// }
func publication_time(wasmInterface *WasmInterface) int64 {
	fmt.Println("prints_l")
}

// void abort() {
//    edump(("abort() called"));
//    EOS_ASSERT( false, abort_called, "abort() called");
// }
func abort(wasmInterface *WasmInterface) {
	fmt.Println("abort")
}

// void eosio_assert( bool condition, null_terminated_ptr msg ) {
//    if( BOOST_UNLIKELY( !condition ) ) {
//       std::string message( msg );
//       edump((message));
//       EOS_THROW( eosio_assert_message_exception, "assertion failure with message: ${s}", ("s",message) );
//    }
// }
func eosio_assert(wasmInterface *WasmInterface, condition int, val int) {
	fmt.Println("eosio_assert")
}

// void eosio_assert_message( bool condition, array_ptr<const char> msg, size_t msg_len ) {
//    if( BOOST_UNLIKELY( !condition ) ) {
//       std::string message( msg, msg_len );
//       edump((message));
//       EOS_THROW( eosio_assert_message_exception, "assertion failure with message: ${s}", ("s",message) );
//    }
// }
func eosio_assert_message(wasmInterface *WasmInterface, condition int, msg int, msg_len size_t) {
	fmt.Println("eosio_assert_message")
}

// void eosio_assert_code( bool condition, uint64_t error_code ) {
//    if( BOOST_UNLIKELY( !condition ) ) {
//       edump((error_code));
//       EOS_THROW( eosio_assert_code_exception,
//                  "assertion failure with error code: ${error_code}", ("error_code", error_code) );
//    }
// }
func eosio_assert_code(wasmInterface *WasmInterface, condition int, error_code int64) {
	fmt.Println("eosio_assert_code")
}

// void eosio_exit(int32_t code) {
//    throw wasm_exit{code};
// }
func eosio_exit(wasmInterface *WasmInterface, code int) {
	fmt.Println("eosio_exit")
}
