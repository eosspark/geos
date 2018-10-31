package wasmgo

import (
	"fmt"
)

// int read_action_data(array_ptr<char> memory, size_t buffer_size) {
//    auto s = context.act.data.size();
//    if( buffer_size == 0 ) return s;

//    auto copy_size = std::min( buffer_size, s );
//    memcpy( memory, context.act.data.data(), copy_size );

//    return copy_size;
// }
func readActionData(w *WasmGo, memory int, bufferSize int) int {
	// if debug {
	// 	fmt.Println("read_action_data")
	// }

	if bufferSize > (1<<16) || memory+bufferSize > (1<<16) {
		//assert
		fmt.Println("access violation")
		return 0
	}

	data := w.context.GetActionData()
	s := len(data)
	if bufferSize == 0 {
		return s
	}

	copySize := min(bufferSize, s)
	setMemory(w, memory, data, 0, copySize)
	return copySize

}

// int action_data_size() {
//    return context.act.data.size();
// }
func actionDataSize(w *WasmGo) int {
	// if debug {
	// 	fmt.Println("action_data_size")
	// }
	return len(w.context.GetActionData())
}

// name current_receiver() {
//    return context.receiver;
// }
func currentReceiver(w *WasmGo) int64 {
	// if debug {
	// 	fmt.Println("current_receiver")
	// }

	return int64(w.context.GetReceiver())
}
