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

	//"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/cvm/wasm"
)

// int get_context_free_data( uint32_t index, array_ptr<char> buffer, size_t buffer_size )const {
//  return context.get_context_free_data( index, buffer, buffer_size );
// }

func get_context_free_data(wasmInterface *WasmInterface, index int, buffer int, buffer_size size_t) int {

	fmt.Println("get_context_free_data")
	//return wasmInterface.context.get_context_free_data( index, buffer, buffer_size );

}
