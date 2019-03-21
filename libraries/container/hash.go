package container

import (
	"hash/fnv"
	"unsafe"
)

func init() {
	const _64bit = unsafe.Sizeof(uintptr(0)) == 8
	if _64bit {
		hash = hash64
	} else {
		hash = hash32
	}
}

var hash func([]byte) uintptr

func hash32(b []byte) uintptr {
	if b == nil {
		return 0
	}

	h := fnv.New32a()
	_, err := h.Write(b)
	if err != nil {
		panic(err)
	}
	return uintptr(h.Sum32())
}

func hash64(b []byte) uintptr {
	if b == nil {
		return 0
	}

	h := fnv.New64a()
	_, err := h.Write(b)
	if err != nil {
		panic(err)
	}
	return uintptr(h.Sum64())
}

func Hash(p interface{}) uintptr {
	switch pt := p.(type) {
	case uintptr:
		return uintptr(pt)
	case int:
		return uintptr(pt)
	case uint:
		return uintptr(pt)
	case int8:
		return uintptr(pt)
	case uint8:
		return uintptr(pt)
	case int16:
		return uintptr(pt)
	case uint16:
		return uintptr(pt)
	case int32:
		return uintptr(pt)
	case uint32:
		return uintptr(pt)
	case int64:
		return uintptr(pt)
	case uint64:
		return uintptr(pt)
	case bool:
		if pt {
			return 1
		}
		return 0
	case nil:
		return 0
	case string:
		return hash([]byte(pt))
	case Hashable:
		return pt.Hashcode()
	}

	return hash(ifaceMem(p))
}

type Hashable interface {
	Hashcode() uintptr
}

func ifaceMem(i interface{}) []byte {
	if i == nil {
		return nil
	}

	typeSize := *(*uintptr)(unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&i))))
	address := *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&i)) + unsafe.Sizeof(uintptr(0))))
	if address == 0 {
		return nil
	}

	bytes := struct {
		addr uintptr
		len  int
		cap  int
	}{address, int(typeSize), int(typeSize)}

	return *(*[]byte)(unsafe.Pointer(&bytes))
}
