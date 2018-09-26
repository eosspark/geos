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

// int get_active_producers(array_ptr<chain::account_name> producers, size_t buffer_size) {
//  auto active_producers = context.get_active_producers();

//  size_t len = active_producers.size();
//  auto s = len * sizeof(chain::account_name);
//  if( buffer_size == 0 ) return s;

//  auto copy_size = std::min( buffer_size, s );
//  memcpy( producers, active_producers.data(), copy_size );

//  return copy_size;
// }
func get_active_producers(wasmInterface *WasmInterface, producers int, buffer_size size_t) int {
	fmt.Println("get_active_producers")
	//return false
}
