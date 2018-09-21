package exec

import (
	"bytes"
	"fmt"
)

// char* memcpy( array_ptr<char> dest, array_ptr<const char> src, size_t length) {
//    EOS_ASSERT((std::abs((ptrdiff_t)dest.value - (ptrdiff_t)src.value)) >= length,
//          overlapping_memory_error, "memcpy can only accept non-aliasing pointers");
//    return (char *)::memcpy(dest, src, length);
// }
func memcpy(w *WasmInterface, dest int, src int, length int) int {
	fmt.Println("memcpy")

	//ASSERT(math.Abs(dest-src) >= length, "memcpy can only accept non-aliasing pointers")
	copy(w.vm.memory[dest:dest+length], w.vm.memory[src:src+length])

	return dest

}

// char* memmove( array_ptr<char> dest, array_ptr<const char> src, size_t length) {
//    return (char *)::memmove(dest, src, length);
// }
func memmove(w *WasmInterface, dest int, src int, length int) int {
	fmt.Println("memmove")

	//ASSERT(math.Abs(dest-src) >= length, "memmove can only accept non-aliasing pointers")
	copy(w.vm.memory[dest:dest+length], w.vm.memory[src:src+length])

	return dest

}

// int memcmp( array_ptr<const char> dest, array_ptr<const char> src, size_t length) {
//    int ret = ::memcmp(dest, src, length);
//    if(ret < 0)
//       return -1;
//    if(ret > 0)
//       return 1;
//    return 0;
// }
func memcmp(w *WasmInterface, dest int, src int, length int) int {
	fmt.Println("memcmp")

	// for i := length - 1; i >= 0; i-- {
	// 	if wasmInterface.vm.memory[dest+i] > wasmInterface.vm.memory[src+i] {
	// 		return 1
	// 	} else if wasmInterface.vm.memory[dest+i] < wasmInterface.vm.memory[src+i] {
	// 		return -1
	// 	}
	// }

	return bytes.Compare(w.vm.memory[dest:dest+length], w.vm.memory[src:src+length])
}

// char* memset( array_ptr<char> dest, int value, size_t length ) {
//    return (char *)::memset( dest, value, length );
// }
func memset(w *WasmInterface, dest int, value int, length int) int {
	fmt.Println("memset")

	//copy[wasmInterface.vm.memory[dest:dest+length-1], byte(value))
	// for i := 0; i < length; i++ {
	// 	wasmInterface.vm.memory[dest + i] = byte(value)
	// }
	b := bytes.Repeat([]byte{byte(value)}, length)
	copy(w.vm.memory[dest:dest+length], b[:])

	return dest
}
