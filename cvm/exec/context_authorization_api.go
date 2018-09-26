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

// void require_authorization( const account_name& account ) {
//   context.require_authorization( account );
// }
func require_authorization(wasmInterface *WasmInterface, account AccountName) {
	fmt.Println("require_authorization")
}

// bool has_authorization( const account_name& account )const {
//   return context.has_authorization( account );
// }
func has_authorization(wasmInterface *WasmInterface, account AccountName) int {
	fmt.Println("has_authorization")
}

// void require_authorization(const account_name& account,
//                                              const permission_name& permission) {
//   context.require_authorization( account, permission );
// }
func require_authorization(wasmInterface *WasmInterface, account int64, permission int64) int {
	fmt.Println("require_authorization")
}

// void require_recipient( account_name recipient ) {
//   context.require_recipient( recipient );
// }
func require_recipient(wasmInterface *WasmInterface, recipient AccountName) {
	fmt.Println("require_recipient")
}

// bool is_account( const account_name& account )const {
//   return context.is_account( account );
// }
func is_account(wasmInterface *WasmInterface, account int64) int {
	fmt.Println("is_account")
}
