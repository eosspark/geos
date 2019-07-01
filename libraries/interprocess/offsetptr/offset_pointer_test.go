package offsetptr

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/eosspark/eos-go/libraries/interprocess/allocator/callocator"

	"github.com/stretchr/testify/assert"
)

func TestPointer_Get(t *testing.T) {
	p := Pointer{}
	pint := new(int)
	*pint = 123
	p.Set(unsafe.Pointer(pint))

	assert.Equal(t, 123, *(*int)(p.Get()))

	pp := Pointer{}
	pp.Set(unsafe.Pointer(pint))

	assert.Equal(t, 123, *(*int)(pp.Get()))

	p1 := Pointer{} //<Pointer>
	pr := new(int)
	*pr = 3

	p2 := NewPointer(unsafe.Pointer(pr))
	p1.Set(unsafe.Pointer(p2))

	assert.Equal(t, 3, *(*int)((*Pointer)(p1.Get()).Get()))
}

func TestPointer_ToRaw(t *testing.T) {
	type raw struct {
		a    int
		next Pointer
	}

	raw1 := raw{a: 10}
	raw1.next.Set(nil)

	raw2 := raw{a: 20}
	raw2.next.Set(unsafe.Pointer(&raw1))

	off := Pointer{}
	off.Set(unsafe.Pointer(&raw2))

	r2 := (*raw)(off.Get())
	r1 := (*raw)(r2.next.Get())

	assert.Equal(t, 10, r1.a)
	assert.Equal(t, 20, r2.a)
}

func TestPointer_Forward(t *testing.T) {
	type sd struct {
		p Pointer
	}

	pp := Pointer{}
	pp.Set(unsafe.Pointer(new(int)))
	*(*int)(pp.Get()) = 100

	s := sd{}
	s.p.Forward(&pp)

	assert.Equal(t, 100, *(*int)(s.p.Get()))
}

func TestPointer_Set(t *testing.T) {
	alloc := callocator.Instance
	pint := NewPointer(alloc.Allocate(unsafe.Sizeof(int(0))))
	//pint := NewPointer(unsafe.Pointer(new(int)))
	*(*int)(pint.Get()) = 100

	runtime.GC()

	assert.Equal(t, 100, *(*int)(pint.Get()))
}
