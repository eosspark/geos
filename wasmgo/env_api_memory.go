package wasmgo

import (
	"bytes"
	"fmt"
	"github.com/eosspark/eos-go/exception"
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

	// if abs(dest-src) < length {
	// 	fmt.Println("memcpy can only accept non-aliasing pointers")
	// 	//ASSERT(math.Abs(dest-src) >= length, "memcpy can only accept non-aliasing pointers")
	// 	return -1
	// }
	exception.EosAssert(abs(dest-src) >= length, &exception.OverlappingMemoryError{}, "memcpy with overlapping memeory")

	copy(w.vm.Memory()[dest:dest+length], w.vm.Memory()[src:src+length])

	return dest

}

// char* memmove( array_ptr<char> dest, array_ptr<const char> src, size_t length) {
//    return (char *)::memmove(dest, src, length);
// }
func memmove(w *WasmGo, dest int, src int, length int) int {
	fmt.Println("memmove")

	//ASSERT(math.Abs(dest-src) >= length, "memmove can only accept non-aliasing pointers")
	// if abs(dest-src) < length {
	// 	fmt.Println("memmove can only accept non-aliasing pointers")
	// 	//ASSERT(math.Abs(dest-src) >= length, "memcpy can only accept non-aliasing pointers")
	// 	return -1
	// }
	exception.EosAssert(abs(dest-src) >= length, &exception.OverlappingMemoryError{}, "memove with overlapping memeory")

	copy(w.vm.Memory()[dest:dest+length], w.vm.Memory()[src:src+length])

	return dest

}

func memcmp(w *WasmGo, dest int, src int, length int) int {
	fmt.Println("memcmp")

	return bytes.Compare(w.vm.Memory()[dest:dest+length], w.vm.Memory()[src:src+length])
}

// char* memset( array_ptr<char> dest, int value, size_t length ) {
//    return (char *)::memset( dest, value, length );
// }
func memset(w *WasmGo, dest int, value int, length int) int {
	fmt.Println("memset")

	cap := cap(w.vm.Memory())
	if cap < dest || cap < dest+length {
		exception.EosAssert(false, &exception.OverlappingMemoryError{}, "memset with heap memory out of bound")
	}

	b := bytes.Repeat([]byte{byte(value)}, length)
	copy(w.vm.Memory()[dest:dest+length], b[:])

	return dest
}
