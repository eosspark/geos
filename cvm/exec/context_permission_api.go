package exec

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
)

// bool check_transaction_authorization( array_ptr<char> trx_data,     size_t trx_size,
//                                             array_ptr<char> pubkeys_data, size_t pubkeys_size,
//                                             array_ptr<char> perms_data,   size_t perms_size
//                                           )
//       {
//          transaction trx = fc::raw::unpack<transaction>( trx_data, trx_size );

//          flat_set<public_key_type> provided_keys;
//          unpack_provided_keys( provided_keys, pubkeys_data, pubkeys_size );

//          flat_set<permission_level> provided_permissions;
//          unpack_provided_permissions( provided_permissions, perms_data, perms_size );

//          try {
//             context.control
//                    .get_authorization_manager()
//                    .check_authorization( trx.actions,
//                                          provided_keys,
//                                          provided_permissions,
//                                          fc::seconds(trx.delay_sec),
//                                          std::bind(&transaction_context::checktime, &context.trx_context),
//                                          false
//                                        );
//             return true;
//          } catch( const authorization_exception& e ) {}

//          return false;
//       }
func checkTransactionAuthorization(w *WasmInterface, trx_data int, trx_size size_t,
	pubkeys_data int, pubkeys_size size_t,
	perms_data int, perms_size size_t) int {
	fmt.Println("check_transaction_authorization")
	return 0
}

//       bool check_permission_authorization( account_name account, permission_name permission,
//                                            array_ptr<char> pubkeys_data, size_t pubkeys_size,
//                                            array_ptr<char> perms_data,   size_t perms_size,
//                                            uint64_t delay_us
//                                          )
//       {
//          EOS_ASSERT( delay_us <= static_cast<uint64_t>(std::numeric_limits<int64_t>::max()),
//                      action_validate_exception, "provided delay is too large" );

//          flat_set<public_key_type> provided_keys;
//          unpack_provided_keys( provided_keys, pubkeys_data, pubkeys_size );

//          flat_set<permission_level> provided_permissions;
//          unpack_provided_permissions( provided_permissions, perms_data, perms_size );

//          try {
//             context.control
//                    .get_authorization_manager()
//                    .check_authorization( account,
//                                          permission,
//                                          provided_keys,
//                                          provided_permissions,
//                                          fc::microseconds(delay_us),
//                                          std::bind(&transaction_context::checktime, &context.trx_context),
//                                          false
//                                        );
//             return true;
//          } catch( const authorization_exception& e ) {}

//          return false;
//       }
func checkPermissionAuthorization(w *WasmInterface, permission common.PermissionName,
	pubkeys_data int, pubkeys_size size_t,
	perms_data int, perms_size size_t,
	delay_us int64) int {
	fmt.Println("check_permission_authorization")
	return 0
}

//       int64_t get_permission_last_used( account_name account, permission_name permission ) {
//          const auto& am = context.control.get_authorization_manager();
//          return am.get_permission_last_used( am.get_permission({account, permission}) ).time_since_epoch().count();
//       };
func getPermissionLastUsed(w *WasmInterface, account common.AccountName, permission common.PermissionName) int64 {
	fmt.Println("get_permission_last_used")

	return w.context.GetPermissionLastUsed(account, permission)
}

//       int64_t get_account_creation_time( account_name account ) {
//          auto* acct = context.db.find<account_object, by_name>(account);
//          EOS_ASSERT( acct != nullptr, action_validate_exception,
//                      "account '${account}' does not exist", ("account", account) );
//          return time_point(acct->creation_date).time_since_epoch().count();
//       }
func getAccountCreationTime(w *WasmInterface, account common.AccountName) int64 {
	fmt.Println("get_account_creation_time")

	return w.context.GetAccountCreateTime(account)
}

//    private:
//       void unpack_provided_keys( flat_set<public_key_type>& keys, const char* pubkeys_data, size_t pubkeys_size ) {
//          keys.clear();
//          if( pubkeys_size == 0 ) return;

//          keys = fc::raw::unpack<flat_set<public_key_type>>( pubkeys_data, pubkeys_size );
//       }

//       void unpack_provided_permissions( flat_set<permission_level>& permissions, const char* perms_data, size_t perms_size ) {
//          permissions.clear();
//          if( perms_size == 0 ) return;

//          permissions = fc::raw::unpack<flat_set<permission_level>>( perms_data, perms_size );
//       }
