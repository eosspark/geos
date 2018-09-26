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

// int is_feature_active( int64_t feature_name ) {
//          return false;
// }
func is_feature_active(wasmInterface *WasmInterface, feature_name int64) int {
	fmt.Println("is_feature_active")
	return false
}

// void activate_feature( int64_t feature_name ) {
//  EOS_ASSERT( false, unsupported_feature, "Unsupported Hardfork Detected" );
// }
func activate_feature(wasmInterface *WasmInterface, feature_name int64) {
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
func set_resource_limits(wasmInterface *WasmInterface, account AccountName, ramBytes int64, netWeight int64, cpuWeigth int64) {
	fmt.Println("set_resource_limits")

}

// void get_resource_limits( account_name account, int64_t& ram_bytes, int64_t& net_weight, int64_t& cpu_weight ) {
//  context.control.get_resource_limits_manager().get_account_limits( account, ram_bytes, net_weight, cpu_weight);
// }
func get_resource_limits(wasmInterface *WasmInterface, account AccountName, ramBytes *int, netWeight *int, cpuWeigth *int) {
	fmt.Println("get_resource_limits")
}

// int64_t set_proposed_producers( array_ptr<char> packed_producer_schedule, size_t datalen) {
//  datastream<const char*> ds( packed_producer_schedule, datalen );
//  vector<producer_key> producers;
//  fc::raw::unpack(ds, producers);
//  EOS_ASSERT(producers.size() <= config::max_producers, wasm_execution_error, "Producer schedule exceeds the maximum producer count for this chain");
//  // check that producers are unique
//  std::set<account_name> unique_producers;
//  for (const auto& p: producers) {
//     EOS_ASSERT( context.is_account(p.producer_name), wasm_execution_error, "producer schedule includes a nonexisting account" );
//     EOS_ASSERT( p.block_signing_key.valid(), wasm_execution_error, "producer schedule includes an invalid key" );
//     unique_producers.insert(p.producer_name);
//  }
//  EOS_ASSERT( producers.size() == unique_producers.size(), wasm_execution_error, "duplicate producer name in producer schedule" );
//  return context.control.set_proposed_producers( std::move(producers) );
// }
func set_proposed_producers(wasmInterface *WasmInterface, packed_producer_schedule int, datalen size_t) {
	fmt.Println("set_proposed_producers")
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
func get_blockchain_parameters_packed(wasmInterface *WasmInterface, packed_blockchain_parameters int, buffer_size size_t) int {
	fmt.Println("get_blockchain_parameters_packed")
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
func set_blockchain_parameters_packed(wasmInterface *WasmInterface, packed_blockchain_parameters int, datalen size_t) {
	fmt.Println("set_blockchain_parameters_packed")
}

// bool is_privileged( account_name n )const {
//  return context.db.get<account_object, by_name>( n ).privileged;
// }
func is_privileged(wasmInterface *WasmInterface, n AccountName) int {
	fmt.Println("is_privileged")
}

// void set_privileged( account_name n, bool is_priv ) {
//  const auto& a = context.db.get<account_object, by_name>( n );
//  context.db.modify( a, [&]( auto& ma ){
//     ma.privileged = is_priv;
//  });
// }
func set_privileged(wasmInterface *WasmInterface, n AccountName, is_priv int) {
	fmt.Println("set_privileged")
}
