package callocator

import "C"
import "unsafe"

//#include <stdlib.h>
//#include <string.h>
import "C"

//allocator memory from cgo, can not freed by GC
type cAllocator struct{}

var Instance = &cAllocator{}

func (a *cAllocator) Allocate(size uintptr) unsafe.Pointer {
	return C.malloc(C.size_t(size))
}

func (a *cAllocator) DeAllocate(p unsafe.Pointer) {
	C.free(p)
}
