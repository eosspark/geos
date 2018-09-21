package exec

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
)

// void require_authorization( const account_name& account ) {
//   context.require_authorization( account );
// }
func requireAuthorization(w *WasmInterface, account common.AccountName) {
	fmt.Println("require_authorization")
	w.context.RequireAuthorization(account)
}

// bool has_authorization( const account_name& account )const {
//   return context.has_authorization( account );
// }
func hasAuthorization(w *WasmInterface, account common.AccountName) int {
	fmt.Println("has_authorization")
	return b2i(w.context.HasAuthorization(account))
}

// void require_authorization(const account_name& account,
//                                              const permission_name& permission) {
//   context.require_authorization( account, permission );
// }
// func require_authorization(wasmInterface *WasmInterface, account int64, permission int64) int {
// 	fmt.Println("require_authorization")
// }
func requireAuth2(w *WasmInterface, account common.AccountName, permission common.PermissionName) {
	fmt.Println("require_authorization")
	w.context.RequireAuthorization2(account, permission)
}

// void require_recipient( account_name recipient ) {
//   context.require_recipient( recipient );
// }
func requireRecipient(w *WasmInterface, recipient common.AccountName) {
	fmt.Println("require_recipient")
	w.context.RequireRecipient(recipient)

}

// bool is_account( const account_name& account )const {
//   return context.is_account( account );
// }
func isAccount(w *WasmInterface, account common.AccountName) int {
	fmt.Println("is_account")
	return b2i(w.context.IsAccount(account))
}
