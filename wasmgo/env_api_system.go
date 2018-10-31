package wasmgo

import (
	"fmt"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
)

func checkTime(w *WasmGo) {

	fmt.Println("checktime")
	w.context.CheckTime()

}

// uint64_t current_time() {
//    return static_cast<uint64_t>( context.control.pending_block_time().time_since_epoch().count() );
// }
func currentTime(w *WasmGo) int64 {
	fmt.Println("current_time")

	//return uint64(wasmInterface.Context.Controller.PendingBlockTime().TimeSinceEpoch().Count())
	return w.context.CurrentTime()
}

// uint64_t publication_time() {
//    return static_cast<uint64_t>( context.trx_context.published.time_since_epoch().count() );
// }
func publicationTime(w *WasmGo) int64 {
	fmt.Println("publication_time")

	return w.context.PublicationTime()
}

// void abort() {
//    edump(("abort() called"));
//    EOS_ASSERT( false, abort_called, "abort() called");
// }
func abort(w *WasmGo) {
	fmt.Println("abort")
	exception.EosAssert(false, &exception.AbortCalled{}, exception.AbortCalled{}.What())
}

// void eosio_assert( bool condition, null_terminated_ptr msg ) {
//    if( BOOST_UNLIKELY( !condition ) ) {
//       std::string message( msg );
//       edump((message));
//       EOS_THROW( eosio_assert_message_exception, "assertion failure with message: ${s}", ("s",message) );
//    }
// }
func eosioAssert(w *WasmGo, condition int, val int) {
	if debug {
		fmt.Println("eosio_assert")
	}

	if condition != 1 {
		message := getMemory(w, val, getStringLength(w, val))

		exception.EosAssert(condition != 1, &exception.EosioAssertMessageException{}, string(message))
		try.Throw(&exception.EosioAssertMessageException{})
		//fmt.Println(string(message))
		// edump(message)
		// E_THROW()
	}
}

// void eosio_assert_message( bool condition, array_ptr<const char> msg, size_t msg_len ) {
//    if( BOOST_UNLIKELY( !condition ) ) {
//       std::string message( msg, msg_len );
//       edump((message));
//       EOS_THROW( eosio_assert_message_exception, "assertion failure with message: ${s}", ("s",message) );
//    }
// }
func eosioAssertMessage(w *WasmGo, condition int, msg int, msgLen size_t) {
	fmt.Println("eosio_assert_message")
}

// void eosio_assert_code( bool condition, uint64_t error_code ) {
//    if( BOOST_UNLIKELY( !condition ) ) {
//       edump((error_code));
//       EOS_THROW( eosio_assert_code_exception,
//                  "assertion failure with error code: ${error_code}", ("error_code", error_code) );
//    }
// }
func eosioAssertCode(w *WasmGo, condition int, errorCode int64) {
	fmt.Println("eosio_assert_code")
}

// void eosio_exit(int32_t code) {
//    throw wasm_exit{code};
// }
func eosioExit(w *WasmGo, code int) {
	fmt.Println("eosio_exit")
}
