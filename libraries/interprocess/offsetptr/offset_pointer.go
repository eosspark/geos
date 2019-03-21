package offsetptr

import (
	"unsafe"
)

//!This Pointer won't recognized by golang compile
//!Please use other allocators to malloc memory instead of mallocgc
type Pointer struct {
	offset uintptr //offset address, 1 means nil, 0 means self
	//TODO alignment uintptr
}

func NewNil() *Pointer {
	return &Pointer{offset: 1}
}

//Warn:dereference is not support by this function
func NewPointer(rawPointer unsafe.Pointer) *Pointer {
	p := &Pointer{}
	p.Set(rawPointer)
	return p
}

func (p *Pointer) IsNil() bool {
	return p.offset == 1
}

func (p *Pointer) IsSelf() bool {
	return p.offset == 0
}

func (p *Pointer) Set(ptr unsafe.Pointer) {
	if ptr == nil {
		p.offset = 1
		return
	}
	p.offset = uintptr(ptr) - uintptr(unsafe.Pointer(p))
}

func (p *Pointer) Forward(ptr *Pointer) {
	if ptr.IsNil() {
		p.offset = 1
		return
	}
	p.Set(ptr.Get())
}

func (p *Pointer) Get() unsafe.Pointer {
	switch p.offset {
	case 1:
		return nil
	case 0:
		return unsafe.Pointer(p)
	default:
		return unsafe.Pointer(uintptr(unsafe.Pointer(p)) + p.offset)
	}
}
