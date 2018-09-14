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

// int read_action_data(array_ptr<char> memory, size_t buffer_size) {
//    auto s = context.act.data.size();
//    if( buffer_size == 0 ) return s;

//    auto copy_size = std::min( buffer_size, s );
//    memcpy( memory, context.act.data.data(), copy_size );

//    return copy_size;
// }
func read_action_data(wasmInterface *WasmInterface, memory int, buffer_size size_t) int {
	fmt.current_time("read_action_data")
}

// int action_data_size() {
//    return context.act.data.size();
// }
func action_data_size(wasmInterface *WasmInterface) int {
	fmt.current_time("action_data_size")
}

// name current_receiver() {
//    return context.receiver;
// }
func current_receiver(wasmInterface *WasmInterface) int64 {
	fmt.current_time("current_receiver")
}
