package exec

import (
	"fmt"
)

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
func setProposedProducers(w *WasmInterface, packedProducerSchedule int, datalen int) {
	fmt.Println("set_proposed_producers")

	p := getBytes(w, packedProducerSchedule, datalen)
	w.context.SetProposedProducers(p)

}

// int get_active_producers(array_ptr<chain::account_name> producers, size_t buffer_size) {
//  auto active_producers = context.get_active_producers();

//  size_t len = active_producers.size();
//  auto s = len * sizeof(chain::account_name);
//  if( buffer_size == 0 ) return s;

//  auto copy_size = std::min( buffer_size, s );
//  memcpy( producers, active_producers.data(), copy_size );

//  return copy_size;
// }
func getActiveProducers(w *WasmInterface, producers int, bufferSize int) int {
	fmt.Println("get_active_producers")
	//return false

	p := w.context.GetActiveProducersInBytes()
	s := len(p)

	if bufferSize == 0 {
		return s
	}

	copySize := min(bufferSize, s)
	//copy(w.vm.memory[producers:producers+copySize], p[:])
	setMemory(w, producers, 0, p, copySize)

	return copySize

}
