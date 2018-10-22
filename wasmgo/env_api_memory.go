package wasmgo

import (
	"bytes"
	"fmt"
)

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// char* memcpy( array_ptr<char> dest, array_ptr<const char> src, size_t length) {
//    EOS_ASSERT((std::abs((ptrdiff_t)dest.value - (ptrdiff_t)src.value)) >= length,
//          overlapping_memory_error, "memcpy can only accept non-aliasing pointers");
//    return (char *)::memcpy(dest, src, length);
// }
func memcpy(w *WasmGo, dest int, src int, length int) int {
	fmt.Println("memcpy")

	if abs(dest-src) < length {
		fmt.Println("memcpy can only accept non-aliasing pointers")
		//ASSERT(math.Abs(dest-src) >= length, "memcpy can only accept non-aliasing pointers")
		return -1
	}
	copy(w.vm.Memory()[dest:dest+length], w.vm.Memory()[src:src+length])

	return dest

}

// char* memmove( array_ptr<char> dest, array_ptr<const char> src, size_t length) {
//    return (char *)::memmove(dest, src, length);
// }
func memmove(w *WasmGo, dest int, src int, length int) int {
	fmt.Println("memmove")

	//ASSERT(math.Abs(dest-src) >= length, "memmove can only accept non-aliasing pointers")
	if abs(dest-src) < length {
		fmt.Println("memmove can only accept non-aliasing pointers")
		//ASSERT(math.Abs(dest-src) >= length, "memcpy can only accept non-aliasing pointers")
		return -1
	}

	copy(w.vm.Memory()[dest:dest+length], w.vm.Memory()[src:src+length])

	return dest

}

func memcmp(w *WasmGo, dest int, src int, length int) int {
	fmt.Println("memcmp")

	// for i := length - 1; i >= 0; i-- {
	// 	if wasmInterface.vm.memory[dest+i] > wasmInterface.vm.memory[src+i] {
	// 		return 1
	// 	} else if wasmInterface.vm.memory[dest+i] < wasmInterface.vm.memory[src+i] {
	// 		return -1
	// 	}
	// }

	return bytes.Compare(w.vm.Memory()[dest:dest+length], w.vm.Memory()[src:src+length])
}

// char* memset( array_ptr<char> dest, int value, size_t length ) {
//    return (char *)::memset( dest, value, length );
// }
func memset(w *WasmGo, dest int, value int, length int) int {
	fmt.Println("memset")

	//copy[wasmInterface.vm.memory[dest:dest+length-1], byte(value))
	// for i := 0; i < length; i++ {
	// 	wasmInterface.vm.memory[dest + i] = byte(value)
	// }
	cap := cap(w.vm.Memory())
	if cap < dest || cap < dest+length {
		//assert()
		fmt.Println("memset heap memory out of bound")
		return -1
	}

	b := bytes.Repeat([]byte{byte(value)}, length)
	copy(w.vm.Memory()[dest:dest+length], b[:])

	return dest
}
