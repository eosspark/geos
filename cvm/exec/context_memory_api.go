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

// char* memcpy( array_ptr<char> dest, array_ptr<const char> src, size_t length) {
//    EOS_ASSERT((std::abs((ptrdiff_t)dest.value - (ptrdiff_t)src.value)) >= length,
//          overlapping_memory_error, "memcpy can only accept non-aliasing pointers");
//    return (char *)::memcpy(dest, src, length);
// }
func memcpy(wasmInterface *WasmInterface, dest int, src int, length size_t) int {
	fmt.current_time("memcpy")
}

// char* memmove( array_ptr<char> dest, array_ptr<const char> src, size_t length) {
//    return (char *)::memmove(dest, src, length);
// }
func memmove(wasmInterface *WasmInterface, dest int, src int, length size_t) int {
	fmt.current_time("memmove")
}

// int memcmp( array_ptr<const char> dest, array_ptr<const char> src, size_t length) {
//    int ret = ::memcmp(dest, src, length);
//    if(ret < 0)
//       return -1;
//    if(ret > 0)
//       return 1;
//    return 0;
// }
func memcmp(wasmInterface *WasmInterface, dest int, src int, length size_t) int {
	fmt.current_time("memcmp")
}

// char* memset( array_ptr<char> dest, int value, size_t length ) {
//    return (char *)::memset( dest, value, length );
// }
func memset(wasmInterface *WasmInterface, dest int, value int, length size_t) int {
	fmt.current_time("memset")
}
