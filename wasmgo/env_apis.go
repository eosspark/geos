package wasmgo

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/eosspark/container/sets/treeset"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/eos_math"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"math"
	"strconv"
	"strings"
	"unsafe"
)

func readActionData(vm *VM) {
	w := vm.WasmGo

	bufferSize := int(vm.popUint64())
	memory := int(vm.popUint64())

	// if bufferSize > (1<<16) || memory+bufferSize > (1<<16) {
	// 	w.ilog.Error("access violation")
	// 	vm.pushUint64(uint64(0))
	// 	return
	// }

	data := w.context.GetActionData()
	s := len(data)
	if bufferSize == 0 {
		vm.pushUint64(uint64(s))
		w.ilog.Debug("action data size:%d", s)
		return
	}

	copySize := min(bufferSize, s)
	w.ilog.Debug("action data:%v size:%d", data, copySize)
	//w.ilog.Debug("action data right:%d", data, memory+copySize)
	setMemory(vm, memory, data, 0, copySize)

	vm.pushUint64(uint64(copySize))

}

func actionDataSize(vm *VM) {
	w := vm.WasmGo
	size := len(w.context.GetActionData())
	vm.pushUint64(uint64(size))

	w.ilog.Debug("actionDataSize:%d", size)

}

func currentReceiver(vm *VM) {
	w := vm.WasmGo
	receiver := w.context.GetReceiver()
	vm.pushUint64(uint64(receiver))
	w.ilog.Debug("currentReceiver:%v", receiver)
}

func requireAuthorization(vm *VM) {
	w := vm.WasmGo

	//w.ilog.Debug("ContextFreeAction:%v", w.context.ContextFreeAction())
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	account := int64(vm.popUint64())

	w.ilog.Debug("account:%v", common.AccountName(account))
	w.context.RequireAuthorization(account)

}

func hasAuthorization(vm *VM) {
	w := vm.WasmGo
	account := int64(vm.popUint64())
	ret := w.context.HasAuthorization(account)

	vm.pushUint64(uint64(b2i(ret)))
	w.ilog.Debug("account:%v authorization:%v", common.AccountName(account), ret)
}

func requireAuth2(vm *VM) {
	w := vm.WasmGo
	permission := int64(vm.popUint64())
	account := int64(vm.popUint64())

	w.ilog.Debug("account:%v permission:%v", common.AccountName(account), common.PermissionName(permission))
	w.context.RequireAuthorization2(account, permission)
}

func requireRecipient(vm *VM) {
	w := vm.WasmGo
	recipient := int64(vm.popUint64())

	w.ilog.Debug("recipient:%v ", common.AccountName(recipient))
	w.context.RequireRecipient(recipient)

}

func isAccount(vm *VM) {
	w := vm.WasmGo
	account := int64(vm.popUint64())
	ret := w.context.IsAccount(account)
	w.ilog.Debug("account:%v isAccount:%v", common.AccountName(account), ret)

	vm.pushUint64(uint64(b2i(ret)))
}

var count = 0

const SHIFT_WIDTH = uint32(unsafe.Sizeof(uint64(0)*8) - 1) //63

func ashlti3(vm *VM) {

	shift := int(vm.popUint64())
	high := int64(vm.popUint64())
	low := int64(vm.popUint64())
	ret := int(vm.popUint64())

	i := eos_math.Int128{Low: uint64(low), High: uint64(high)}
	i.LeftShifts(shift)

	re, _ := rlp.EncodeToBytes(i)
	setMemory(vm, ret, re, 0, len(re))
}

func ashrti3(vm *VM) {

	shift := int(vm.popUint64())
	high := int64(vm.popUint64())
	low := int64(vm.popUint64())
	ret := int(vm.popUint64())

	i := eos_math.Int128{Low: uint64(low), High: uint64(high)}
	i.RightShifts(shift)

	re, _ := rlp.EncodeToBytes(i)
	setMemory(vm, ret, re, 0, len(re))
}

func lshlti3(vm *VM) {

	shift := int(vm.popUint64())
	high := int64(vm.popUint64())
	low := int64(vm.popUint64())
	ret := int(vm.popUint64())

	i := eos_math.Int128{Low: uint64(low), High: uint64(high)}
	i.LeftShifts(shift)

	re, _ := rlp.EncodeToBytes(i)
	setMemory(vm, ret, re, 0, len(re))
}

func lshrti3(vm *VM) {

	shift := int(vm.popUint64())
	high := int64(vm.popUint64())
	low := int64(vm.popUint64())
	ret := int(vm.popUint64())

	i := eos_math.Uint128{Low: uint64(low), High: uint64(high)}
	i.RightShifts(shift)

	re, _ := rlp.EncodeToBytes(i)
	setMemory(vm, ret, re, 0, len(re))
}

func divti3(vm *VM) {

	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())
	ret := int(vm.popUint64())

	lhs := eos_math.Int128{Low: uint64(la), High: uint64(ha)}
	rhs := eos_math.Int128{Low: uint64(lb), High: uint64(hb)}

	EosAssert(!rhs.IsZero(), &ArithmeticException{}, "divide by zero")

	quotient, _ := lhs.Div(rhs)
	re, _ := rlp.EncodeToBytes(quotient)
	setMemory(vm, ret, re, 0, len(re))
}

func udivti3(vm *VM) {

	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())
	ret := int(vm.popUint64())

	lhs := eos_math.Uint128{Low: uint64(la), High: uint64(ha)}
	rhs := eos_math.Uint128{Low: uint64(lb), High: uint64(hb)}

	EosAssert(!rhs.IsZero(), &ArithmeticException{}, "divide by zero")
	quotient, _ := lhs.Div(rhs)

	re, _ := rlp.EncodeToBytes(quotient)
	setMemory(vm, ret, re, 0, len(re))
}

func multi3(vm *VM) {

	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())
	ret := int(vm.popUint64())

	lhs := eos_math.Int128{Low: uint64(la), High: uint64(ha)}
	rhs := eos_math.Int128{Low: uint64(lb), High: uint64(hb)}

	re, _ := rlp.EncodeToBytes(lhs.Mul(rhs))
	setMemory(vm, ret, re, 0, len(re))

}

func modti3(vm *VM) {

	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())
	ret := int(vm.popUint64())

	lhs := eos_math.Int128{High: uint64(ha), Low: uint64(la)}
	rhs := eos_math.Int128{High: uint64(hb), Low: uint64(lb)}
	EosAssert(!rhs.IsZero(), &ArithmeticException{}, "divide by zero")

	_, remainder := lhs.Div(rhs)
	re, _ := rlp.EncodeToBytes(remainder)
	setMemory(vm, ret, re, 0, len(re))
}

func umodti3(vm *VM) {

	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())
	ret := int(vm.popUint64())

	lhs := eos_math.Uint128{Low: uint64(la), High: uint64(ha)}
	rhs := eos_math.Uint128{Low: uint64(lb), High: uint64(hb)}

	EosAssert(!rhs.IsZero(), &ArithmeticException{}, "divide by zero")
	_, remainder := lhs.Div(rhs)
	re, _ := rlp.EncodeToBytes(remainder)
	setMemory(vm, ret, re, 0, len(re))
}

func addtf3(vm *VM) {

	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())
	ret := int(vm.popUint64())

	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}

	re, _ := rlp.EncodeToBytes(a.Add(b))
	setMemory(vm, ret, re, 0, len(re))
}

func subtf3(vm *VM) {

	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())
	ret := int(vm.popUint64())

	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}

	re, _ := rlp.EncodeToBytes(a.Sub(b))
	setMemory(vm, ret, re, 0, len(re))
}

func multf3(vm *VM) {

	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())
	ret := int(vm.popUint64())

	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}

	re, _ := rlp.EncodeToBytes(a.Mul(b))
	setMemory(vm, ret, re, 0, len(re))
}

func divtf3(vm *VM) {

	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())
	ret := int(vm.popUint64())

	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}

	re, _ := rlp.EncodeToBytes(a.Div(b))
	setMemory(vm, ret, re, 0, len(re))
}

func negtf2(vm *VM) {

	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())
	ret := int(vm.popUint64())

	high := uint64(ha)
	high ^= uint64(1) << 63
	f128 := eos_math.Float128{Low: uint64(la), High: high}

	re, _ := rlp.EncodeToBytes(f128)
	setMemory(vm, ret, re, 0, len(re))
}

//func extendsftf2(w *WasmGo, ret int, f float32) { //TODO f float??
func extendsftf2(vm *VM) {

	f := uint32(vm.popUint64())
	ret := int(vm.popUint64())

	//f32 := eos_math.Float32(math.Float32bits(f))
	f32 := eos_math.Float32(f)
	f128 := eos_math.F32ToF128(f32)

	re, _ := rlp.EncodeToBytes(f128)
	setMemory(vm, ret, re, 0, len(re))
}

//func extenddftf2(w *WasmGo, ret int, d float64) { //TODO d double??
func extenddftf2(vm *VM) { //TODO d double??

	d := vm.popUint64()
	ret := int(vm.popUint64())

	//f64 := eos_math.Float64(math.Float64bits(d))
	f64 := eos_math.Float64(d)
	f128 := eos_math.F64ToF128(f64)

	re, _ := rlp.EncodeToBytes(f128)
	setMemory(vm, ret, re, 0, len(re))
}

//func trunctfdf2(w *WasmGo, l, h int64) float64 { //TODO double??
func trunctfdf2(vm *VM) { //TODO double??

	h := int64(vm.popUint64())
	l := int64(vm.popUint64())

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	f64 := eos_math.F128ToF64(f128)
	//return math.Float64frombits(uint64(f64))
	vm.pushUint64(uint64(f64))
}

//func trunctfsf2(w *WasmGo, l, h int64) float32 { //TODO float??
func trunctfsf2(vm *VM) { //TODO float??
	h := int64(vm.popUint64())
	l := int64(vm.popUint64())

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	f32 := eos_math.F128ToF32(f128)
	//return math.Float32frombits(uint32(f32))
	vm.pushUint64(uint64(f32))
}

func fixtfsi(vm *VM) {
	h := int64(vm.popUint64())
	l := int64(vm.popUint64())

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	//return int(eos_math.F128ToI32(f128, 0, false))
	vm.pushUint64(uint64(eos_math.F128ToI32(f128, 0, false)))
}

func fixtfdi(vm *VM) {
	h := int64(vm.popUint64())
	l := int64(vm.popUint64())

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	//return eos_math.F128ToI64(f128, 0, false)
	vm.pushUint64(uint64(eos_math.F128ToI64(f128, 0, false)))
}

func fixunstfsi(vm *VM) {
	h := int64(vm.popUint64())
	l := int64(vm.popUint64())

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	//return int(eos_math.F128ToUi32(f128, 0, false))
	vm.pushUint64(uint64(eos_math.F128ToUi32(f128, 0, false)))
}

func fixunstfdi(vm *VM) {
	h := int64(vm.popUint64())
	l := int64(vm.popUint64())

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	//return int64(eos_math.F128ToUi64(f128, 0, false))
	vm.pushUint64(uint64(eos_math.F128ToUi64(f128, 0, false)))
}

func fixtfti(vm *VM) {

	h := int64(vm.popUint64())
	l := int64(vm.popUint64())
	ret := int(vm.popUint64())

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	int128 := eos_math.Fixtfti(f128)

	re, _ := rlp.EncodeToBytes(int128)
	setMemory(vm, ret, re, 0, len(re))
}

func fixunstfti(vm *VM) {
	h := int64(vm.popUint64())
	l := int64(vm.popUint64())
	ret := int(vm.popUint64())

	f := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	uint128 := eos_math.Fixunstfti(f)

	re, _ := rlp.EncodeToBytes(uint128)
	setMemory(vm, ret, re, 0, len(re))
}

func fixsfti(vm *VM) { //TODO float??
	a := uint32(vm.popUint64())
	ret := int(vm.popUint64())

	//int128 := eos_math.Fixsfti(math.Float32bits(a))
	int128 := eos_math.Fixsfti(a)
	re, _ := rlp.EncodeToBytes(int128)
	setMemory(vm, ret, re, 0, len(re))
}

func fixdfti(vm *VM) { //TODO double??
	a := vm.popUint64()
	ret := int(vm.popUint64())

	//int128 := eos_math.Fixdfti(math.Float64bits(a))
	int128 := eos_math.Fixdfti(a)
	re, _ := rlp.EncodeToBytes(int128)
	setMemory(vm, ret, re, 0, len(re))
}

func fixunssfti(vm *VM) { //TODO float??
	a := uint32(vm.popUint64())
	ret := int(vm.popUint64())

	//uint128 := eos_math.Fixunssfti(math.Float32bits(a))
	uint128 := eos_math.Fixunssfti(a)
	re, _ := rlp.EncodeToBytes(uint128)
	setMemory(vm, ret, re, 0, len(re))
}

func fixunsdfti(vm *VM) { //TODO double??
	a := vm.popUint64()
	ret := int(vm.popUint64())

	//uint128 := eos_math.Fixunsdfti(math.Float64bits(a))
	uint128 := eos_math.Fixunsdfti(a)
	re, _ := rlp.EncodeToBytes(uint128)
	setMemory(vm, ret, re, 0, len(re))
}

func floatsidf(vm *VM) { //TODO double??
	i := int(vm.popUint64())

	//return math.Float64frombits(uint64(eos_math.I32ToF64(int32(i))))
	vm.pushUint64(uint64(eos_math.I32ToF64(int32(i))))
}

func floatsitf(vm *VM) {
	i := int(vm.popUint64())
	ret := int(vm.popUint64())

	re := eos_math.I32ToF128(int32(i))
	setMemory(vm, ret, re.Bytes(), 0, len(re.Bytes()))
}

func floatditf(vm *VM) {
	a := int64(vm.popUint64())
	ret := int(vm.popUint64())

	re := eos_math.I64ToF128(a)
	setMemory(vm, ret, re.Bytes(), 0, len(re.Bytes()))
}

func floatunsitf(vm *VM) {
	i := int(vm.popUint64())
	ret := int(vm.popUint64())

	re := eos_math.Ui32ToF128(uint32(i))
	setMemory(vm, ret, re.Bytes(), 0, len(re.Bytes()))
}

func floatunditf(vm *VM) {
	a := vm.popUint64()
	ret := int(vm.popUint64())

	re := eos_math.Ui64ToF128(uint64(a))
	setMemory(vm, ret, re.Bytes(), 0, len(re.Bytes()))
}

func floattidf(vm *VM) { //TODO double
	h := vm.popUint64()
	l := vm.popUint64()

	v := eos_math.Int128{Low: uint64(l), High: uint64(h)}
	//return eos_math.Floattidf(v)
	vm.pushFloat64(eos_math.Floattidf(v))
}

func floatuntidf(vm *VM) { //TODO double
	h := vm.popUint64()
	l := vm.popUint64()

	v := eos_math.Uint128{Low: uint64(l), High: uint64(h)}
	// return eos_math.Floatuntidf(v)
	// return
	vm.pushFloat64(eos_math.Floatuntidf(v))
}

func cmptf2(vm *VM) { //TODO unsame with regist

	return_value_if_nan := int(vm.popUint64())
	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())

	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}
	if _unordtf2(la, ha, lb, hb) != 0 {
		vm.pushUint64(uint64(return_value_if_nan))
		return
		//return return_value_if_nan
	}
	if a.F128Lt(b) {
		vm.pushInt64(-1)
		return
		//return -1
	}
	if a.F128EQ(b) {
		vm.pushUint64(0)
		return
		//return 0
	}
	vm.pushUint64(1)
	//return 1
}

func eqtf2(vm *VM) {
	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())

	vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, 1)))
	//return _cmptf2(vm, la, ha, lb, hb, 1)
}

func netf2(vm *VM) {
	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())

	vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, 1)))
	//return _cmptf2(vm, la, ha, lb, hb, 1)
}

func getf2(vm *VM) {
	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())

	vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, -1)))
	//return _cmptf2(vm, la, ha, lb, hb, -1)
}

func gttf2(vm *VM) {
	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())

	vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, 0)))

	//return _cmptf2(vm, la, ha, lb, hb, 0)
}

func letf2(vm *VM) {
	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())

	vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, 1)))

	//return _cmptf2(vm, la, ha, lb, hb, 1)
}

func lttf2(vm *VM) {
	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())

	vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, 0)))
	//return _cmptf2(vm, la, ha, lb, hb, 0)
}

func __cmptf2(vm *VM) {
	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())

	vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, 1)))
	//return _cmptf2(vm, la, ha, lb, hb, 1)
}

func unordtf2(vm *VM) {
	hb := int64(vm.popUint64())
	lb := int64(vm.popUint64())
	ha := int64(vm.popUint64())
	la := int64(vm.popUint64())

	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}
	if a.IsNan() || b.IsNan() {
		vm.pushUint64(1)
		//return 1
	}
	vm.pushUint64(0)
	//return 0
}

func _cmptf2(la, ha, lb, hb int64, return_value_if_nan int) int { //TODO unsame with regist
	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}
	if _unordtf2(la, ha, lb, hb) != 0 {
		return return_value_if_nan
	}
	if a.F128Lt(b) {
		return -1
	}
	if a.F128EQ(b) {
		return 0
	}
	return 1
}

func _unordtf2(la, ha, lb, hb int64) int {
	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}
	if a.IsNan() || b.IsNan() {
		return 1
	}
	return 0
}

/*
REGISTER_INTRINSICS(compiler_builtins,
(__ashlti3,     void(int, int64_t, int64_t, int)               )
(__ashrti3,     void(int, int64_t, int64_t, int)               )
(__lshlti3,     void(int, int64_t, int64_t, int)               )
(__lshrti3,     void(int, int64_t, int64_t, int)               )
(__divti3,      void(int, int64_t, int64_t, int64_t, int64_t)  )
(__udivti3,     void(int, int64_t, int64_t, int64_t, int64_t)  )
(__modti3,      void(int, int64_t, int64_t, int64_t, int64_t)  )
(__umodti3,     void(int, int64_t, int64_t, int64_t, int64_t)  )
(__multi3,      void(int, int64_t, int64_t, int64_t, int64_t)  )
(__addtf3,      void(int, int64_t, int64_t, int64_t, int64_t)  )
(__subtf3,      void(int, int64_t, int64_t, int64_t, int64_t)  )
(__multf3,      void(int, int64_t, int64_t, int64_t, int64_t)  )
(__divtf3,      void(int, int64_t, int64_t, int64_t, int64_t)  )
(__eqtf2,       int(int64_t, int64_t, int64_t, int64_t)        )
(__netf2,       int(int64_t, int64_t, int64_t, int64_t)        )
(__getf2,       int(int64_t, int64_t, int64_t, int64_t)        )
(__gttf2,       int(int64_t, int64_t, int64_t, int64_t)        )
(__lttf2,       int(int64_t, int64_t, int64_t, int64_t)        )
(__letf2,       int(int64_t, int64_t, int64_t, int64_t)        )
(__cmptf2,      int(int64_t, int64_t, int64_t, int64_t)        )
(__unordtf2,    int(int64_t, int64_t, int64_t, int64_t)        )
(__negtf2,      void (int, int64_t, int64_t)                   )
(__floatsitf,   void (int, int)                                )
(__floatunsitf, void (int, int)                                )
(__floatditf,   void (int, int64_t)                            )
(__floatunditf, void (int, int64_t)                            )
(__floattidf,   double (int64_t, int64_t)                      )
(__floatuntidf, double (int64_t, int64_t)                      )
(__floatsidf,   double(int)                                    )
(__extendsftf2, void(int, float)                               )
(__extenddftf2, void(int, double)                              )
(__fixtfti,     void(int, int64_t, int64_t)                    )
(__fixtfdi,     int64_t(int64_t, int64_t)                      )
(__fixtfsi,     int(int64_t, int64_t)                          )
(__fixunstfti,  void(int, int64_t, int64_t)                    )
(__fixunstfdi,  int64_t(int64_t, int64_t)                      )
(__fixunstfsi,  int(int64_t, int64_t)                          )
(__fixsfti,     void(int, float)                               )
(__fixdfti,     void(int, double)                              )
(__fixunssfti,  void(int, float)                               )
(__fixunsdfti,  void(int, double)                              )
(__trunctfdf2,  double(int64_t, int64_t)                       )
(__trunctfsf2,  float(int64_t, int64_t)                        )
);
*/

func assertRecoverKey(vm *VM) {
	w := vm.WasmGo

	publen := int(vm.popUint64())
	pub := int(vm.popUint64())
	siglen := int(vm.popUint64())
	sig := int(vm.popUint64())
	digest := int(vm.popUint64())

	digBytes := getSha256(vm, digest)
	sigBytes := getMemory(vm, sig, siglen)
	pubBytes := getMemory(vm, pub, publen)

	w.ilog.Debug("digest:%v signature:%v publickey:%v", digBytes, sigBytes, pubBytes)

	s := ecc.NewSigNil()
	p := ecc.NewPublicKeyNil()

	d := digBytes
	rlp.DecodeBytes(sigBytes, s)
	rlp.DecodeBytes(pubBytes, p)

	check, err := s.PublicKey(d)
	EosAssert(err == nil, &CryptoApiException{}, "can not get the right publickey from digest")
	EosAssert(strings.Compare(check.String(), p.String()) == 0, &CryptoApiException{}, "Error expected key different than recovered key")

}

func recoverKey(vm *VM) {
	//fmt.Println("recover_key")
	w := vm.WasmGo

	publen := int(vm.popUint64())
	pub := int(vm.popUint64())
	siglen := int(vm.popUint64())
	sig := int(vm.popUint64())
	digest := int(vm.popUint64())

	digBytes := getSha256(vm, digest)
	sigBytes := getMemory(vm, sig, siglen)

	s := ecc.NewSigNil()
	rlp.DecodeBytes(sigBytes, s)
	check, _ := s.PublicKey(digBytes)

	p, err := rlp.EncodeToBytes(check)
	if err != nil {
		//return -1
		vm.pushInt64(-1)
		return
	}

	w.ilog.Debug("digest:%v signature:%v publickey:%v", digBytes, sigBytes, p)

	l := len(p)
	if l > publen {
		l = publen
	}
	setMemory(vm, pub, p, 0, l)
	vm.pushUint64(uint64(l))
}

type shaInterface interface {
	Write(p []byte) (nn int, err error)
	Sum(b []byte) []byte
}

func encode(w *WasmGo, s shaInterface, data []byte, dataLen int) []byte {

	bs := int(common.DefaultConfig.HashingChecktimeBlockSize)

	i := 0
	l := dataLen

	for i = 0; l > bs; i += bs {
		s.Write(data[i : i+bs])
		l -= bs
		w.context.CheckTime()
	}

	s.Write(data[i : i+l])

	return s.Sum(nil)

}

func assertSha256(vm *VM) {

	w := vm.WasmGo

	hashVal := int(vm.popUint64())
	dataLen := int(vm.popUint64())
	data := int(vm.popUint64())

	dataBytes := getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewSha256()
	hashEncode := encode(w, s, dataBytes, dataLen)
	hash := getSha256(vm, hashVal)

	//w.ilog.Debug("encoded:%#v hash:%#v data:%#v", hashEncode, hash, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)
	EosAssert(bytes.Compare(hashEncode, hash) == 0, &CryptoApiException{}, "sha256 hash mismatch")

}

func assertSha1(vm *VM) {

	w := vm.WasmGo

	hashVal := int(vm.popUint64())
	dataLen := int(vm.popUint64())
	data := int(vm.popUint64())

	dataBytes := getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewSha1()
	hashEncode := encode(w, s, dataBytes, dataLen)
	hash := getSha1(vm, hashVal)

	//w.ilog.Debug("encoded:%#v hash:%#v data:%#v", hashEncode, hash, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)
	EosAssert(bytes.Compare(hashEncode, hash) == 0, &CryptoApiException{}, "sha1 hash mismatch")
}

func assertSha512(vm *VM) {
	w := vm.WasmGo

	hashVal := int(vm.popUint64())
	dataLen := int(vm.popUint64())
	data := int(vm.popUint64())

	dataBytes := getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewSha512()
	hashEncode := encode(w, s, dataBytes, dataLen)
	hash := getSha512(vm, hashVal)

	//w.ilog.Debug("encoded:%#v hash:%#v data:%#v", hashEncode, hash, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)
	EosAssert(bytes.Compare(hashEncode, hash) == 0, &CryptoApiException{}, "sha512 hash mismatch")

}

func assertRipemd160(vm *VM) {
	w := vm.WasmGo

	hashVal := int(vm.popUint64())
	dataLen := int(vm.popUint64())
	data := int(vm.popUint64())

	dataBytes := getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewRipemd160()
	hashEncode := encode(w, s, dataBytes, dataLen)
	hash := getRipemd160(vm, hashVal)

	//w.ilog.Debug("encoded:%#v hash:%#v data:%#v", hashEncode, hash, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)
	EosAssert(bytes.Compare(hashEncode, hash) == 0, &CryptoApiException{}, "ripemd160 hash mismatch")
}

func sha1(vm *VM) {
	w := vm.WasmGo

	hashVal := int(vm.popUint64())
	dataLen := int(vm.popUint64())
	data := int(vm.popUint64())

	dataBytes := getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewSha1()
	hashEncode := encode(w, s, dataBytes, dataLen)
	setSha1(vm, hashVal, hashEncode)

	//w.ilog.Debug("encoded:%#v data:%#v", hashEncode, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)
}

func sha256(vm *VM) {
	w := vm.WasmGo

	hashVal := int(vm.popUint64())
	dataLen := int(vm.popUint64())
	data := int(vm.popUint64())

	dataBytes := getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewSha256()

	hashEncode := encode(w, s, dataBytes, dataLen)
	setSha256(vm, hashVal, hashEncode)

	//w.ilog.Debug("encoded:%#v data:%#v", hashEncode, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hex.EncodeToString(hashEncode), dataBytes)
}

func sha512(vm *VM) {
	w := vm.WasmGo

	hashVal := int(vm.popUint64())
	dataLen := int(vm.popUint64())
	data := int(vm.popUint64())

	dataBytes := getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewSha512()

	hashEncode := encode(w, s, dataBytes, dataLen)
	setSha512(vm, hashVal, hashEncode)

	//w.ilog.Debug("encoded:%#v data:%#v", hashEncode, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)
}

func ripemd160(vm *VM) {
	w := vm.WasmGo

	hashVal := int(vm.popUint64())
	dataLen := int(vm.popUint64())
	data := int(vm.popUint64())

	dataBytes := getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewRipemd160()
	hashEncode := encode(w, s, dataBytes, dataLen)
	setRipemd160(vm, hashVal, hashEncode)

	//w.ilog.Debug("encoded:%#v data:%#v", hashEncode, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)
}

func dbStoreI64(vm *VM) {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	bufferSize := int(vm.popUint64())
	buffer := int(vm.popUint64())
	id := vm.popUint64()
	payer := vm.popUint64()
	table := vm.popUint64()
	scope := vm.popUint64()

	bytes := getMemory(vm, buffer, bufferSize)

	iterator := w.context.DbStoreI64(scope, table, payer, id, bytes)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("scope:%v table:%v payer:%v id:%d data:%v iterator:%d",
		common.ScopeName(scope), common.TableName(table), common.AccountName(payer), id, bytes, iterator)

}

func dbUpdateI64(vm *VM) {
	w := vm.WasmGo

	bufferSize := int(vm.popUint64())
	buffer := int(vm.popUint64())
	payer := vm.popUint64()
	iterator := int(vm.popUint64())

	bytes := getMemory(vm, buffer, bufferSize)
	w.context.DbUpdateI64(iterator, payer, bytes)

	w.ilog.Debug("data:%v iterator:%d payer:%v ", bytes, iterator, payer)

}

func dbRemoveI64(vm *VM) {
	w := vm.WasmGo
	iterator := int(vm.popUint64())

	w.context.DbRemoveI64(iterator)
	w.ilog.Debug("iterator:%d", iterator)

}

func dbGetI64(vm *VM) {
	w := vm.WasmGo

	bufferSize := int(vm.popUint64())
	buffer := int(vm.popUint64())
	iterator := int(vm.popUint64())

	bytes := make([]byte, bufferSize)
	size := w.context.DbGetI64(iterator, bytes, bufferSize)
	if bufferSize == 0 {
		vm.pushUint64(uint64(size))
		w.ilog.Debug("iterator:%d size:%d", iterator, size)
		return
	}
	setMemory(vm, buffer, bytes, 0, size)
	vm.pushUint64(uint64(size))

	w.ilog.Debug("iterator:%d data:%v size:%d", iterator, bytes, size)
}

func dbNextI64(vm *VM) {
	w := vm.WasmGo

	primary := int(vm.popUint64())
	itr := int(vm.popUint64())

	var p uint64
	iterator := w.context.DbNextI64(itr, &p)

	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return
	}
	setUint64(vm, primary, p)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d nextIterator:%d primary:%d", itr, iterator, p)

}

func dbPreviousI64(vm *VM) {
	w := vm.WasmGo

	primary := int(vm.popUint64())
	itr := int(vm.popUint64())

	var p uint64
	iterator := w.context.DbPreviousI64(itr, &p)
	w.ilog.Debug("dbNextI64 iterator:%d", iterator)
	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return
	}
	setUint64(vm, primary, p)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d priviousIterator:%d primary:%d", itr, iterator, p)
}

func dbFindI64(vm *VM) {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	id := vm.popUint64()
	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	iterator := w.context.DbFindI64(code, scope, table, id)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v id:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), id, iterator)

}

func dbLowerboundI64(vm *VM) {
	w := vm.WasmGo

	id := vm.popUint64()
	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	iterator := w.context.DbLowerboundI64(code, scope, table, id)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v id:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), id, iterator)
}

func dbUpperboundI64(vm *VM) {
	w := vm.WasmGo

	id := vm.popUint64()
	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	iterator := w.context.DbUpperboundI64(code, scope, table, id)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v id:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), id, iterator)
}

func dbEndI64(vm *VM) {
	w := vm.WasmGo

	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	iterator := w.context.DbEndI64(code, scope, table)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), iterator)
}

//secondaryKey Index
func dbIdx64Store(vm *VM) {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	pValue := int(vm.popUint64())
	id := vm.popUint64()
	payer := vm.popUint64()
	table := vm.popUint64()
	scope := vm.popUint64()

	secondaryKey := getUint64(vm, pValue)
	iterator := w.context.Idx64Store(scope, table, payer, id, &secondaryKey)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("scope:%v table:%v payer:%v id:%d secondaryKey:%d iterator:%d",
		common.ScopeName(scope), common.TableName(table), common.AccountName(payer), id, secondaryKey, iterator)
}

func dbIdx64Remove(vm *VM) {
	w := vm.WasmGo
	iterator := int(vm.popUint64())

	w.context.Idx64Remove(iterator)
	w.ilog.Debug("iterator:%d", iterator)
}

func dbIdx64Update(vm *VM) {
	w := vm.WasmGo

	pValue := int(vm.popUint64())
	payer := vm.popUint64()
	iterator := int(vm.popUint64())

	secondaryKey := getUint64(vm, pValue)
	w.context.Idx64Update(iterator, payer, &secondaryKey)

	w.ilog.Debug("payer:%v data:%v secondaryKey:%d", common.AccountName(payer), secondaryKey, iterator)

}

func dbIdx64findSecondary(vm *VM) {
	w := vm.WasmGo

	pPrimary := int(vm.popUint64())
	pSecondary := int(vm.popUint64())
	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	var primaryKey uint64
	secondaryKey := getUint64(vm, pSecondary)
	iterator := w.context.Idx64FindSecondary(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return
	}
	setUint64(vm, pPrimary, primaryKey)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)
}

func dbIdx64Lowerbound(vm *VM) {
	w := vm.WasmGo

	pPrimary := int(vm.popUint64())
	pSecondary := int(vm.popUint64())
	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	var primaryKey uint64

	secondaryKey := getUint64(vm, pSecondary)
	iterator := w.context.Idx64Lowerbound(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return
	}

	setUint64(vm, pPrimary, primaryKey)
	setUint64(vm, pSecondary, secondaryKey)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)
}

func dbIdx64Upperbound(vm *VM) {
	w := vm.WasmGo

	pPrimary := int(vm.popUint64())
	pSecondary := int(vm.popUint64())
	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	var primaryKey uint64
	secondaryKey := getUint64(vm, pSecondary)
	iterator := w.context.Idx64Upperbound(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return
	}
	setUint64(vm, pPrimary, primaryKey)
	setUint64(vm, pSecondary, secondaryKey)

	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)
}

func dbIdx64End(vm *VM) {
	w := vm.WasmGo

	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	iterator := w.context.Idx64End(code, scope, table)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), iterator)
}

func dbIdx64Next(vm *VM) {
	w := vm.WasmGo

	primary := int(vm.popUint64())
	itr := int(vm.popUint64())

	var p uint64
	iterator := w.context.Idx64Next(itr, &p)
	w.ilog.Debug("dbIdx64Next iterator:%d", iterator)
	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return
	}
	setUint64(vm, primary, p)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d nextIterator:%d primary:%d", itr, iterator, p)
}

func dbIdx64Previous(vm *VM) {
	w := vm.WasmGo

	primary := int(vm.popUint64())
	itr := int(vm.popUint64())

	var p uint64
	iterator := w.context.Idx64Previous(itr, &p)
	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return
	}
	setUint64(vm, primary, p)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d previousIterator:%d primary:%d", itr, iterator, p)
}

func dbIdx64FindPrimary(vm *VM) {
	w := vm.WasmGo

	primary := vm.popUint64()
	pSecondary := int(vm.popUint64())
	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	var secondaryKey uint64
	iterator := w.context.Idx64FindPrimary(code, scope, table, &secondaryKey, primary)
	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, iterator)
		return
	}
	setUint64(vm, pSecondary, secondaryKey)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d secondaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, secondaryKey, iterator)

}

func dbIdxDoubleStore(vm *VM) {
	w := vm.WasmGo

	pValue := int(vm.popUint64())
	id := vm.popUint64()
	payer := vm.popUint64()
	table := vm.popUint64()
	scope := vm.popUint64()

	secondaryKey := eos_math.Float64(getUint64(vm, pValue))
	//float := math.Float64frombits(getUint64(vm, pValue))
	//w.ilog.Info("float:%v", float)

	iterator := w.context.IdxDoubleStore(scope, table, payer, id, &secondaryKey)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("scope:%v table:%v payer:%v id:%d secondaryKey:%v iterator:%d",
		common.ScopeName(scope), common.TableName(table), common.AccountName(payer), id, secondaryKey, iterator)
}

func dbIdxDoubleRemove(vm *VM) {
	w := vm.WasmGo
	iterator := int(vm.popUint64())

	w.context.IdxDoubleRemove(iterator)
	w.ilog.Debug("iterator:%d", iterator)
}

func dbIdxDoubleUpdate(vm *VM) {
	w := vm.WasmGo

	pValue := int(vm.popUint64())
	payer := vm.popUint64()
	iterator := int(vm.popUint64())

	secondaryKey := eos_math.Float64(getUint64(vm, pValue))
	w.context.IdxDoubleUpdate(iterator, payer, &secondaryKey)

	w.ilog.Debug("payer:%v secondaryKey:%v iterator:%v", common.AccountName(payer), secondaryKey, iterator)

}

func dbIdxDoublefindSecondary(vm *VM) {

	w := vm.WasmGo

	pPrimary := int(vm.popUint64())
	pSecondary := int(vm.popUint64())
	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	var primaryKey uint64
	secondaryKey := eos_math.Float64(getUint64(vm, pSecondary))
	iterator := w.context.IdxDoubleFindSecondary(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return
	}
	setUint64(vm, pPrimary, primaryKey)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)
}

func dbIdxDoubleLowerbound(vm *VM) {
	w := vm.WasmGo

	pPrimary := int(vm.popUint64())
	pSecondary := int(vm.popUint64())
	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	var primaryKey uint64
	secondaryKey := eos_math.Float64(getUint64(vm, pSecondary))
	iterator := w.context.IdxDoubleLowerbound(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return
	}
	setUint64(vm, pPrimary, primaryKey)
	setUint64(vm, pSecondary, uint64(secondaryKey))
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)
}

func dbIdxDoubleUpperbound(vm *VM) {
	w := vm.WasmGo

	pPrimary := int(vm.popUint64())
	pSecondary := int(vm.popUint64())
	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	var primaryKey uint64
	secondaryKey := eos_math.Float64(getUint64(vm, pSecondary))
	iterator := w.context.IdxDoubleUpperbound(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return
	}
	setUint64(vm, pPrimary, primaryKey)
	setUint64(vm, pSecondary, uint64(secondaryKey))
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)
}

func dbIdxDoubleEnd(vm *VM) {
	w := vm.WasmGo

	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	iterator := w.context.IdxDoubleEnd(code, scope, table)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), iterator)
}

func dbIdxDoubleNext(vm *VM) {
	w := vm.WasmGo

	primary := int(vm.popUint64())
	itr := int(vm.popUint64())

	var p uint64
	iterator := w.context.IdxDoubleNext(itr, &p)
	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return
	}

	setUint64(vm, primary, p)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d nextIterator:%d primary:%d", itr, iterator, p)
}

func dbIdxDoublePrevious(vm *VM) {
	w := vm.WasmGo

	primary := int(vm.popUint64())
	itr := int(vm.popUint64())

	var p uint64
	iterator := w.context.IdxDoublePrevious(itr, &p)
	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d proviousIterator:%d", itr, iterator)
		return
	}
	setUint64(vm, primary, p)
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d previousIterator:%d primary:%d", itr, iterator, p)
}

func dbIdxDoubleFindPrimary(vm *VM) {
	w := vm.WasmGo

	primary := vm.popUint64()
	pSecondary := int(vm.popUint64())
	table := vm.popUint64()
	scope := vm.popUint64()
	code := vm.popUint64()

	var secondaryKey eos_math.Float64
	iterator := w.context.IdxDoubleFindPrimary(code, scope, table, &secondaryKey, primary)
	if iterator <= -1 {
		vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, iterator)
		return
	}
	setUint64(vm, pSecondary, uint64(secondaryKey))
	vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d secondaryKey:%v iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, secondaryKey, iterator)
}

// (db_##IDX##_remove,         void(int))\
// (db_##IDX##_update,         void(int,int64_t,int))\
// (db_##IDX##_find_primary,   int(int64_t,int64_t,int64_t,int,int64_t))\
// (db_##IDX##_find_secondary, int(int64_t,int64_t,int64_t,int,int))\
// (db_##IDX##_lowerbound,     int(int64_t,int64_t,int64_t,int,int))\
// (db_##IDX##_upperbound,     int(int64_t,int64_t,int64_t,int,int))\
// (db_##IDX##_end,            int(int64_t,int64_t,int64_t))\
// (db_##IDX##_next,           int(int, int))\
// (db_##IDX##_previous,       int(int, int))

// DB_API_METHOD_WRAPPERS_SIMPLE_SECONDARY(idx64,  uint64_t)
// DB_API_METHOD_WRAPPERS_SIMPLE_SECONDARY(idx128, uint128_t)
// DB_API_METHOD_WRAPPERS_ARRAY_SECONDARY(idx256, 2, uint128_t)
// DB_API_METHOD_WRAPPERS_FLOAT_SECONDARY(idx_double, float64_t)
// DB_API_METHOD_WRAPPERS_FLOAT_SECONDARY(idx_long_double, float128_t)

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func memcpy(vm *VM) {
	w := vm.WasmGo

	length := int(vm.popUint64())
	src := int(vm.popUint64())
	dest := int(vm.popUint64())

	w.ilog.Debug("dest:%d src:%d length:%d ", dest, src, length)

	EosAssert(abs(dest-src) >= length, &OverlappingMemoryError{}, "memcpy can only accept non-aliasing pointers")
	copy(vm.Memory()[dest:dest+length], vm.Memory()[src:src+length])
	vm.pushUint64(uint64(dest))
}

func memmove(vm *VM) {
	w := vm.WasmGo

	length := int(vm.popUint64())
	src := int(vm.popUint64())
	dest := int(vm.popUint64())

	w.ilog.Debug("dest:%d src:%d length:%d ", dest, src, length)

	//EosAssert(abs(dest-src) >= length, &OverlappingMemoryError{}, "memove with overlapping memeory")
	copy(vm.Memory()[dest:dest+length], vm.Memory()[src:src+length])
	vm.pushUint64(uint64(dest))

}

func memcmp(vm *VM) {
	w := vm.WasmGo

	length := int(vm.popUint64())
	src := int(vm.popUint64())
	dest := int(vm.popUint64())

	w.ilog.Debug("dest:%d src:%d length:%d ", dest, src, length)

	ret := bytes.Compare(vm.Memory()[dest:dest+length], vm.Memory()[src:src+length])
	vm.pushUint64(uint64(ret))
}

// char* memset( array_ptr<char> dest, int value, size_t length ) {
//    return (char *)::memset( dest, value, length );
// }
func memset(vm *VM) {
	w := vm.WasmGo

	length := int(vm.popUint64())
	value := int(vm.popUint64())
	dest := int(vm.popUint64())

	w.ilog.Debug("dest:%d value:%d length:%d ", dest, value, length)

	cap := cap(vm.Memory())
	if cap < dest || cap < dest+length {
		EosAssert(false, &OverlappingMemoryError{}, "memset with heap memory out of bound")
	}

	b := bytes.Repeat([]byte{byte(value)}, length)
	copy(vm.Memory()[dest:dest+length], b[:])

	vm.pushUint64(uint64(dest))

}

// func free(w *WasmGo, index int) {
// 	fmt.Println("free")

// }

// bool check_transaction_authorization( array_ptr<char> trx_data,     size_t trx_size,
//                                             array_ptr<char> pubkeys_data, size_t pubkeys_size,
//                                             array_ptr<char> perms_data,   size_t perms_size
//                                           )
//       {
//          transaction trx = fc::raw::unpack<transaction>( trx_data, trx_size );

//          flat_set<public_key_type> provided_keys;
//          unpack_provided_keys( provided_keys, pubkeys_data, pubkeys_size );

//          flat_set<permission_level> provided_permissions;
//          unpack_provided_permissions( provided_permissions, perms_data, perms_size );

//          try {
//             context.control
//                    .get_authorization_manager()
//                    .check_authorization( trx.actions,
//                                          provided_keys,
//                                          provided_permissions,
//                                          fc::seconds(trx.delay_sec),
//                                          std::bind(&transaction_context::checktime, &context.trx_context),
//                                          false
//                                        );
//             return true;
//          } catch( const authorization_exception& e ) {}

//          return false;
//       }
//func checkTransactionAuthorization(w *WasmGo, trx_data int, trx_size size_t,
//	pubkeys_data int, pubkeys_size size_t,
//	perms_data int, perms_size size_t) int {
//	fmt.Println("check_transaction_authorization")
//	return 0
//}
func checkTransactionAuthorization(vm *VM) {
	fmt.Println("check_transaction_authorization")
	//return 0
}

//       bool check_permission_authorization( account_name account, permission_name permission,
//                                            array_ptr<char> pubkeys_data, size_t pubkeys_size,
//                                            array_ptr<char> perms_data,   size_t perms_size,
//                                            uint64_t delay_us
//                                          )
//       {
//          EOS_ASSERT( delay_us <= static_cast<uint64_t>(std::numeric_limits<int64_t>::max()),
//                      action_validate_exception, "provided delay is too large" );

//          flat_set<public_key_type> provided_keys;
//          unpack_provided_keys( provided_keys, pubkeys_data, pubkeys_size );

//          flat_set<permission_level> provided_permissions;
//          unpack_provided_permissions( provided_permissions, perms_data, perms_size );

//          try {
//             context.control
//                    .get_authorization_manager()
//                    .check_authorization( account,
//                                          permission,
//                                          provided_keys,
//                                          provided_permissions,
//                                          fc::microseconds(delay_us),
//                                          std::bind(&transaction_context::checktime, &context.trx_context),
//                                          false
//                                        );
//             return true;
//          } catch( const authorization_exception& e ) {}

//          return false;
//       }
//func checkPermissionAuthorization(w *WasmGo, permission common.PermissionName,
//	pubkeys_data int, pubkeys_size size_t,
//	perms_data int, perms_size size_t,
//	delay_us int64) int {
//	fmt.Println("check_permission_authorization")
//	return 0
//}
func checkPermissionAuthorization(vm *VM) {
	fmt.Println("check_permission_authorization")

	w := vm.WasmGo

	delayUS := vm.popUint64()
	permsSize := int(vm.popUint64())
	permsData := int(vm.popUint64())
	pubkeysSize := int(vm.popUint64())
	pubkeysData := int(vm.popUint64())
	permission := common.PermissionName(vm.popUint64())
	account := common.AccountName(vm.popUint64())

	pubkeysDataBytes := getMemory(vm, pubkeysData, pubkeysSize)
	permsDataBytes := getMemory(vm, permsData, permsSize)

	providedKeys := treeset.NewWith(ecc.TypePubKey, ecc.ComparePubKey)
	providedPermissions := treeset.NewWith(types.PermissionLevelType, types.ComparePermissionLevel)

	unpackProvidedKeys(providedKeys, &pubkeysDataBytes)
	unpackProvidedPermissions(providedPermissions, &permsDataBytes)

	w.ilog.Debug("account:%v permission:%v providedKeys:%v providedPermissions:%v", account, permission, providedKeys, providedPermissions)

	returning := false
	Try(func() {
		w.context.CheckAuthorization(account,
			permission,
			providedKeys,
			providedPermissions,
			delayUS)
	}).Catch(func(e Exception) {
		returning = true
	}).End()

	if returning {
		vm.pushUint64(0)
		return
	}

	vm.pushUint64(1)
}

func unpackProvidedKeys(ps *treeset.Set, pubkeysData *[]byte) {
	if len(*pubkeysData) == 0 {
		return
	}

	providedKey := []ecc.PublicKey{}
	rlp.DecodeBytes(*pubkeysData, &providedKey)

	for _, pk := range providedKey {
		ps.AddItem(pk)
	}

}

func unpackProvidedPermissions(ps *treeset.Set, permsData *[]byte) {
	if len(*permsData) == 0 {
		return
	}

	permissions := []types.PermissionLevel{}
	rlp.DecodeBytes(*permsData, &permissions)

	for _, permission := range permissions {
		ps.AddItem(permission)
	}

}

func getPermissionLastUsed(vm *VM) {
	w := vm.WasmGo
	account := common.AccountName(vm.popUint64())
	permission := common.PermissionName(vm.popUint64())

	ret := w.context.GetPermissionLastUsed(account, permission)
	vm.pushUint64(uint64(ret.TimeSinceEpoch().Count()))

	w.ilog.Debug("account:%v permission:%v LastUsed:%v", account, permission, ret)
}

func getAccountCreationTime(vm *VM) {
	w := vm.WasmGo
	account := common.AccountName(vm.popUint64())
	//w.ilog.Debug("account:%v ", account)

	ret := w.context.GetAccountCreateTime(account)
	vm.pushUint64(uint64(ret.TimeSinceEpoch().Count()))

	w.ilog.Debug("account:%v creationTime:%v", account, ret)

}

//    private:
//       void unpack_provided_keys( flat_set<public_key_type>& keys, const char* pubkeys_data, size_t pubkeys_size ) {
//          keys.clear();
//          if( pubkeys_size == 0 ) return;

//          keys = fc::raw::unpack<flat_set<public_key_type>>( pubkeys_data, pubkeys_size );
//       }

//       void unpack_provided_permissions( flat_set<permission_level>& permissions, const char* perms_data, size_t perms_size ) {
//          permissions.clear();
//          if( perms_size == 0 ) return;

//          permissions = fc::raw::unpack<flat_set<permission_level>>( perms_data, perms_size );
//       }

func prints(vm *VM) {
	if !ignore {
		w := vm.WasmGo
		strIndex := int(vm.popUint64())

		str := string(getMemory(vm, strIndex, getStringLength(vm, strIndex)))
		w.context.ContextAppend(str)

		w.ilog.Debug("prints:%v", str)
	}
}

func printsl(vm *VM) {
	if !ignore {
		w := vm.WasmGo
		strLen := int(vm.popUint64())
		strIndex := int(vm.popUint64())
		str := string(getMemory(vm, strIndex, strLen))
		w.context.ContextAppend(str)

		w.ilog.Debug("prints_l:%v", str)
	}
}

func printi(vm *VM) {
	if !ignore {
		w := vm.WasmGo
		//w.ilog.Debug("printi")
		val := int64(vm.popUint64())
		str := strconv.FormatInt(val, 10)
		w.context.ContextAppend(str)

		w.ilog.Debug("printi:%v", str)
	}
}

func printui(vm *VM) {
	if !ignore {
		w := vm.WasmGo
		//w.ilog.Debug("printui")
		val := vm.popUint64()
		str := strconv.FormatUint(val, 10)
		w.context.ContextAppend(str)

		w.ilog.Debug("printui:%v", str)
	}
}

func printi128(vm *VM) {
	if !ignore {
		w := vm.WasmGo
		val := int(vm.popUint64())

		bytes := getMemory(vm, val, 16)
		var v eos_math.Int128
		rlp.DecodeBytes(bytes, &v)
		str := v.String()
		w.context.ContextAppend(str)

		w.ilog.Debug("printi128:%v", str)
	}

}

func printui128(vm *VM) {
	if !ignore {
		w := vm.WasmGo
		val := int(vm.popUint64())

		bytes := getMemory(vm, val, 16)
		var v eos_math.Uint128
		rlp.DecodeBytes(bytes, &v)
		str := v.String()
		w.context.ContextAppend(str)

		w.ilog.Debug("printui128:%v", str)
	}

}

func printsf(vm *VM) {
	//fmt.Println("printsf")
	if !ignore {
		w := vm.WasmGo

		val := math.Float32frombits(uint32(vm.popUint64()))
		str := strconv.FormatFloat(float64(val), 'e', 6, 32)
		// val := math.Float64frombits(vm.popUint64())
		// str := strconv.FormatFloat(val, 'e', 6, 32)

		w.context.ContextAppend(str)
		w.ilog.Debug("printsf:%v", str)
	}

}

func printdf(vm *VM) {
	if !ignore {
		w := vm.WasmGo
		val := math.Float64frombits(vm.popUint64())
		str := strconv.FormatFloat(val, 'e', 15, 64)

		w.context.ContextAppend(str)
		w.ilog.Debug("printdf:%v", str)
	}
}

func printqf(vm *VM) {
	if !ignore {
		w := vm.WasmGo
		val := int(vm.popUint64())

		bytes := getMemory(vm, val, 16)
		var v eos_math.Float128
		rlp.DecodeBytes(bytes, &v)
		str := v.String()
		w.context.ContextAppend(str)

		w.ilog.Debug("printqf:%v", str)
	}

}

func printn(vm *VM) {
	if !ignore {
		w := vm.WasmGo
		val := vm.popUint64()
		str := common.S(val)
		w.context.ContextAppend(str)

		w.ilog.Debug("printn:%v", str)
	}
}

func printhex(vm *VM) {
	if !ignore {
		w := vm.WasmGo
		dataLen := int(vm.popUint64())
		data := int(vm.popUint64())
		str := hex.EncodeToString(getMemory(vm, data, dataLen))
		w.context.ContextAppend(str)

		w.ilog.Debug("printhex:%v", str)
	}
}

func isFeatureActive(vm *VM) {

	w := vm.WasmGo
	featureName := int64(vm.popUint64())
	vm.pushUint64(uint64(b2i(false)))

	w.ilog.Debug("featureName:%v", common.S(uint64(featureName)))

}

func activateFeature(vm *VM) {
	w := vm.WasmGo
	featureName := int64(vm.popUint64())

	EosAssert(false, &UnsupportedFeature{}, "Unsupported Hardfork Detected")
	w.ilog.Debug("featureName:%v", common.S(uint64(featureName)))
}

func setResourceLimits(vm *VM) {
	w := vm.WasmGo
	cpuWeight := vm.popUint64()
	netWeight := vm.popUint64()
	ramBytes := vm.popUint64()
	account := common.AccountName(vm.popUint64())

	EosAssert(int64(ramBytes) >= -1, &WasmExecutionError{}, "invalid value for ram resource limit expected [-1,INT64_MAX]")
	EosAssert(int64(netWeight) >= -1, &WasmExecutionError{}, "invalid value for net resource limit expected [-1,INT64_MAX]")
	EosAssert(int64(cpuWeight) >= -1, &WasmExecutionError{}, "invalid value for cpu resource limit expected [-1,INT64_MAX]")

	if w.context.SetAccountLimits(account, int64(ramBytes), int64(netWeight), int64(cpuWeight)) {
		w.context.ValidateRamUsageInsert(account)
	}
	w.ilog.Debug("account:%v ramBytes:%d netWeight:%d cpuWeight:%d", account, ramBytes, netWeight, cpuWeight)
}

func getResourceLimits(vm *VM) {
	w := vm.WasmGo
	cpuWeight := int(vm.popUint64())
	netWeight := int(vm.popUint64())
	ramBytes := int(vm.popUint64())
	account := common.AccountName(vm.popUint64())

	var r, n, c int64
	w.context.GetAccountLimits(account, &r, &n, &c)

	setUint64(vm, ramBytes, uint64(r))
	setUint64(vm, netWeight, uint64(n))
	setUint64(vm, cpuWeight, uint64(c))

	w.ilog.Debug("account:%v ramBytes:%d netWeight:%d cpuWeigth:%d", account, ramBytes, netWeight, cpuWeight)

}

func getBlockchainParametersPacked(vm *VM) {
	w := vm.WasmGo
	bufferSize := int(vm.popUint64())
	packedBlockchainParameters := int(vm.popUint64())

	configuration := w.context.GetBlockchainParameters()
	p, _ := rlp.EncodeToBytes(configuration)
	//p := w.context.GetBlockchainParametersPacked()
	size := len(p)
	w.ilog.Debug("BlockchainParameters:%v bufferSize:%d size:%d", configuration, bufferSize, size)

	if bufferSize == 0 {
		vm.pushUint64(uint64(size))
		return
	}

	if size <= bufferSize {
		setMemory(vm, packedBlockchainParameters, p, 0, size)
		vm.pushUint64(uint64(size))
		//w.ilog.Debug("BlockchainParameters:%v", configuration)
		return
	}
	vm.pushUint64(0)

}

func setBlockchainParametersPacked(vm *VM) {
	w := vm.WasmGo
	dataLen := int(vm.popUint64())
	packedBlockchainParameters := int(vm.popUint64())

	// p := make([]byte, datalen)
	// getMemory(vm,packedBlockchainParameters, 0, p, datalen)
	p := getMemory(vm, packedBlockchainParameters, dataLen)

	cfg := types.ChainConfig{}
	rlp.DecodeBytes(p, &cfg)

	//w.context.SetBlockchainParametersPacked(p)
	w.context.SetBlockchainParameters(&cfg)

	w.ilog.Debug("BlockchainParameters:%v ", cfg)

}

func isPrivileged(vm *VM) {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	account := common.AccountName(vm.popUint64())

	ret := w.context.IsPrivileged(account)
	vm.pushUint64(uint64(b2i(ret)))

	w.ilog.Debug("account:%v privileged:%v", account, ret)
}

func setPrivileged(vm *VM) {
	//fmt.Println("set_privileged")

	w := vm.WasmGo
	isPriv := int(vm.popUint64())
	account := common.AccountName(vm.popUint64())
	w.context.SetPrivileged(account, i2b(isPriv))

	w.ilog.Debug("account:%v privileged:%v", account, i2b(isPriv))
}

func setProposedProducers(vm *VM) {
	//fmt.Println("set_proposed_producers")

	w := vm.WasmGo
	dataLen := int(vm.popUint64())
	packedProducerSchedule := int(vm.popUint64())

	p := getBytes(vm, packedProducerSchedule, dataLen)
	ret := w.context.SetProposedProducers(p)
	vm.pushUint64(uint64(ret))

	producers := []types.ProducerKey{}
	rlp.DecodeBytes(p, &producers)
	w.ilog.Debug("packedProducerSchedule:%v ", producers)
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
func getActiveProducers(vm *VM) {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	bufferSize := int(vm.popUint64())
	producers := int(vm.popUint64())

	p := w.context.GetActiveProducersInBytes()
	s := len(p)

	if bufferSize == 0 {
		vm.pushUint64(uint64(s))
		w.ilog.Debug("size:%d", s)
		return
	}

	copySize := min(bufferSize, s)
	setMemory(vm, producers, p, 0, copySize)

	vm.pushUint64(uint64(copySize))

	accounts := []common.AccountName{}
	rlp.DecodeBytes(p, &accounts)
	w.ilog.Debug("producers:%v", accounts)

}

func checkTime(vm *VM) {
	w := vm.WasmGo
	w.context.CheckTime()

	w.ilog.Debug("time:%v", common.Now())

}

func currentTime(vm *VM) {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	ret := w.context.CurrentTime()
	vm.pushUint64(uint64(ret.TimeSinceEpoch().Count()))

	//w.ilog.Debug("time:%v", uint64(ret.TimeSinceEpoch().Count()))
	w.ilog.Debug("time:%v", ret)

}

func publicationTime(vm *VM) {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	ret := w.context.PublicationTime()
	vm.pushUint64(uint64(ret.TimeSinceEpoch().Count()))

	w.ilog.Debug("time:%v", ret)

}

func abort(vm *VM) {
	w := vm.WasmGo

	EosAssert(false, &AbortCalled{}, AbortCalled{}.What())
	w.ilog.Debug("abort")
}

func eosioAssert(vm *VM) {
	w := vm.WasmGo
	val := int(vm.popUint64())
	condition := int(vm.popUint64())

	message := string(getMemory(vm, val, getStringLength(vm, val)))

	b, err := strconv.Atoi(message)
	if err == nil {
		if b == 999 {
			getStringLength(vm, val)
		}
	}

	w.ilog.Debug("message:%v", string(message))

	if condition != 1 {
		EosAssert(false, &EosioAssertMessageException{}, "assertion failure with message: %v", message)
		Throw(&EosioAssertMessageException{})
	}
}

func eosioAssertMessage(vm *VM) {
	w := vm.WasmGo
	msgLen := int(vm.popUint64())
	msg := int(vm.popUint64())
	condition := int(vm.popUint64())

	message := string(getMemory(vm, msg, msgLen))
	w.ilog.Debug("message:%v", string(message))

	if condition != 1 {
		EosAssert(false, &EosioAssertMessageException{}, "assertion failure with message: %v", message)
		//Throw(&EosioAssertMessageException{},"assertion failure with message: %v", message)
		Throw(&EosioAssertMessageException{})
	}

}

func eosioAssertCode(vm *VM) {
	w := vm.WasmGo
	errorCode := int64(vm.popUint64())
	condition := int(vm.popUint64())

	w.ilog.Debug("error code:%d", errorCode)
	if condition != 1 {
		EosAssert(false, &EosioAssertMessageException{}, "assertion failure with error code: %d", errorCode)
		//Throw(&EosioAssertMessageException{},"assertion failure with error code: %d", errorCode)
		Throw(&EosioAssertMessageException{})
	}

}

// void eosio_exit(int32_t code) {
//    throw wasm_exit{code};
// }
func eosioExit(vm *VM) {
	w := vm.WasmGo
	errorCode := int(vm.popUint64())
	w.ilog.Debug("error code:%d", errorCode)

	//Throw(wasmExit(code))

}

func sendInline(vm *VM) {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	dataLen := int(vm.popUint64())
	data := int(vm.popUint64())

	EosAssert(!w.context.InlineActionTooBig(dataLen), &InlineActionTooBig{}, "inline action too big")

	action := getBytes(vm, data, dataLen)
	act := types.Action{}
	rlp.DecodeBytes(action, &act)
	w.context.ExecuteInline(&act)

	w.ilog.Debug("action:%v", act)

}

func sendContextFreeInline(vm *VM) {

	w := vm.WasmGo
	dataLen := int(vm.popUint64())
	data := int(vm.popUint64())
	//w.ilog.Debug("action data:%d size:%d", data, dataLen)

	EosAssert(!w.context.InlineActionTooBig(dataLen), &InlineActionTooBig{}, "inline action too big")

	action := getBytes(vm, data, dataLen)
	act := types.Action{}
	rlp.DecodeBytes(action, &act)

	w.ilog.Debug("action:%v", act)

	w.context.ExecuteContextFreeInline(&act)

}

// void send_deferred( const uint128_t& sender_id, account_name payer, array_ptr<char> data, size_t data_len, uint32_t replace_existing) {
//    try {
//       transaction trx;
//       fc::raw::unpack<transaction>(data, data_len, trx);
//       context.schedule_deferred_transaction(sender_id, payer, std::move(trx), replace_existing);
//    } FC_RETHROW_EXCEPTIONS(warn, "data as hex: ${data}", ("data", fc::to_hex(data, data_len)))
// }
func sendDeferred(vm *VM) {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	replaceExisting := int(vm.popUint64())
	dataLen := int(vm.popUint64())
	data := int(vm.popUint64())
	payer := common.AccountName(vm.popUint64())
	senderId := int(vm.popUint64())

	bytes := getMemory(vm, senderId, 16)
	id := &eos_math.Uint128{}
	rlp.DecodeBytes(bytes, id)

	trx := getBytes(vm, data, dataLen)
	transaction := types.Transaction{}
	rlp.DecodeBytes(trx, &transaction)
	w.context.ScheduleDeferredTransaction(id, payer, &transaction, i2b(replaceExisting))

	w.ilog.Debug("id:%v transaction:%v", id, transaction)
}

func cancelDeferred(vm *VM) {
	w := vm.WasmGo
	senderId := int(vm.popUint64())

	bytes := getMemory(vm, senderId, 16)
	id := &eos_math.Uint128{}
	rlp.DecodeBytes(bytes, id)

	//return b2i(w.context.CancelDeferredTransaction(common.TransactionIdType{id}))
	ret := w.context.CancelDeferredTransaction(id)
	vm.pushUint64(uint64(b2i(ret)))

	w.ilog.Debug("id:%v", id)
}

// int read_transaction( array_ptr<char> data, size_t buffer_size ) {
//    bytes trx = context.get_packed_transaction();

//    auto s = trx.size();
//    if( buffer_size == 0) return s;

//    auto copy_size = std::min( buffer_size, s );
//    memcpy( data, trx.data(), copy_size );

//    return copy_size;
// }
func readTransaction(vm *VM) {
	w := vm.WasmGo
	bufferSize := int(vm.popUint64())
	buffer := int(vm.popUint64())

	transaction := w.context.GetPackedTransaction()
	trx, _ := rlp.EncodeToBytes(transaction)

	s := len(trx)
	if bufferSize == 0 {
		w.ilog.Debug("transaction size:%d", s)
		vm.pushUint64(uint64(s))
		return
	}

	copySize := min(bufferSize, s)
	setMemory(vm, buffer, trx, 0, copySize)
	vm.pushUint64(uint64(copySize))

	w.ilog.Debug("transaction:%v", transaction)

}

func transactionSize(vm *VM) {
	//fmt.Println("transaction_size")
	w := vm.WasmGo

	transaction := w.context.GetPackedTransaction()
	trx, _ := rlp.EncodeToBytes(transaction)
	s := len(trx)
	vm.pushUint64(uint64(s))

	w.ilog.Debug("transaction size:%d", s)
}

func expiration(vm *VM) {
	w := vm.WasmGo

	expiration := w.context.Expiration()
	vm.pushUint64(uint64(expiration))

	w.ilog.Debug("expiration:%v", expiration)
}

func taposBlockNum(vm *VM) {
	w := vm.WasmGo

	taposBlockNum := w.context.TaposBlockNum()
	vm.pushUint64(uint64(taposBlockNum))

	w.ilog.Debug("taposBlockNum:%v", taposBlockNum)
}

func taposBlockPrefix(vm *VM) {
	w := vm.WasmGo

	taposBlockPrefix := w.context.TaposBlockPrefix()
	vm.pushUint64(uint64(taposBlockPrefix))

	w.ilog.Debug("taposBlockPrefix:%v", taposBlockPrefix)
}

func getAction(vm *VM) {
	w := vm.WasmGo
	bufferSize := int(vm.popUint64())
	buffer := int(vm.popUint64())
	index := int(vm.popUint64())
	typ := int(vm.popUint64())

	action := w.context.GetAction(uint32(typ), index)
	s, _ := rlp.EncodeSize(action)
	if bufferSize == 0 || bufferSize < s {
		vm.pushUint64(uint64(s))
		w.ilog.Debug("action size:%d", s)
		return
		//return s
	}

	bytes, _ := rlp.EncodeToBytes(action)
	setMemory(vm, buffer, bytes, 0, s)
	vm.pushUint64(uint64(s))
	w.ilog.Debug("action :%v size:%d", *action, s)

}

func getContextFreeData(vm *VM) {
	w := vm.WasmGo
	bufferSize := int(vm.popUint64())
	buffer := int(vm.popUint64())
	index := int(vm.popUint64())

	EosAssert(w.context.ContextFreeAction(), &UnaccessibleApi{}, "this API may only be called from context_free apply")

	s, data := w.context.GetContextFreeData(index, bufferSize)
	if bufferSize == 0 || s == -1 {
		vm.pushUint64(uint64(s))
		w.ilog.Debug("context free data size:%d", s)
		return
	}
	setMemory(vm, buffer, data, 0, s)
	vm.pushUint64(uint64(s))
	w.ilog.Debug("context free data :%v size:%d", data, s)

}
