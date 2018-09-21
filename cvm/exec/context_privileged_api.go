package exec

import (
	"fmt"

	//"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
)

// int is_feature_active( int64_t feature_name ) {
//          return false;
// }
func is_feature_active(w *WasmInterface, feature_name int64) int {
	fmt.Println("is_feature_active")
	return b2i(false)
}

// void activate_feature( int64_t feature_name ) {
//  EOS_ASSERT( false, unsupported_feature, "Unsupported Hardfork Detected" );
// }
func activate_feature(w *WasmInterface, feature_name int64) {
	fmt.Println("activate_feature")
	//EOS_ASSERT( false, unsupported_feature, "Unsupported Hardfork Detected" );
}

// void set_resource_limits( account_name account, int64_t ram_bytes, int64_t net_weight, int64_t cpu_weight) {
//  EOS_ASSERT(ram_bytes >= -1, wasm_execution_error, "invalid value for ram resource limit expected [-1,INT64_MAX]");
//  EOS_ASSERT(net_weight >= -1, wasm_execution_error, "invalid value for net resource weight expected [-1,INT64_MAX]");
//  EOS_ASSERT(cpu_weight >= -1, wasm_execution_error, "invalid value for cpu resource weight expected [-1,INT64_MAX]");
//  if( context.control.get_mutable_resource_limits_manager().set_account_limits(account, ram_bytes, net_weight, cpu_weight) ) {
//     context.trx_context.validate_ram_usage.insert( account );
//  }
// }
func set_resource_limits(w *WasmInterface, account common.AccountName, ramBytes uint64, netWeight uint64, cpuWeigth uint64) {
	fmt.Println("set_resource_limits")

	w.context.SetResourceLimits(account, ramBytes, netWeight, cpuWeigth)

}

// void get_resource_limits( account_name account, int64_t& ram_bytes, int64_t& net_weight, int64_t& cpu_weight ) {
//  context.control.get_resource_limits_manager().get_account_limits( account, ram_bytes, net_weight, cpu_weight);
// }
func get_resource_limits(w *WasmInterface, account common.AccountName, ramBytes int, netWeight int, cpuWeigth int) {
	fmt.Println("get_resource_limits")

	var r, n, c uint64
	w.context.GetResourceLimits(account, &r, &n, &c)

	setUint64(w, ramBytes, r)
	setUint64(w, netWeight, n)
	setUint64(w, cpuWeigth, c)

}

// uint32_t get_blockchain_parameters_packed( array_ptr<char> packed_blockchain_parameters, size_t buffer_size) {
//  auto& gpo = context.control.get_global_properties();

//  auto s = fc::raw::pack_size( gpo.configuration );
//  if( buffer_size == 0 ) return s;

//  if ( s <= buffer_size ) {
//     datastream<char*> ds( packed_blockchain_parameters, s );
//     fc::raw::pack(ds, gpo.configuration);
//     return s;
//  }
//  return 0;
// }
func get_blockchain_parameters_packed(w *WasmInterface, packed_blockchain_parameters int, buffer_size int) int {
	fmt.Println("get_blockchain_parameters_packed")

	p := w.context.GetBlockchainParametersPacked()
	s := len(p)

	if s <= buffer_size {
		copy(w.vm.memory[packed_blockchain_parameters:packed_blockchain_parameters+s], p[0:s])
		return s
	}

	return 0
}

// void set_blockchain_parameters_packed( array_ptr<char> packed_blockchain_parameters, size_t datalen) {
//  datastream<const char*> ds( packed_blockchain_parameters, datalen );
//  chain::chain_config cfg;
//  fc::raw::unpack(ds, cfg);
//  cfg.validate();
//  context.db.modify( context.control.get_global_properties(),
//     [&]( auto& gprops ) {
//          gprops.configuration = cfg;
//  });
// }
func set_blockchain_parameters_packed(w *WasmInterface, packed_blockchain_parameters int, datalen int) {
	fmt.Println("set_blockchain_parameters_packed")

	p := make([]byte, datalen)
	copy(p[0:datalen], w.vm.memory[packed_blockchain_parameters:packed_blockchain_parameters+datalen])

	w.context.SetBlockchainParametersPacked(p)
}

// bool is_privileged( account_name n )const {
//  return context.db.get<account_object, by_name>( n ).privileged;
// }
func is_privileged(w *WasmInterface, n common.AccountName) int {
	fmt.Println("is_privileged")

	return b2i(w.context.IsPrivileged(n))
}

// void set_privileged( account_name n, bool is_priv ) {
//  const auto& a = context.db.get<account_object, by_name>( n );
//  context.db.modify( a, [&]( auto& ma ){
//     ma.privileged = is_priv;
//  });
// }
func set_privileged(w *WasmInterface, n common.AccountName, is_priv int) {
	fmt.Println("set_privileged")

	w.context.SetPrivileged(n, i2b(is_priv))
}
