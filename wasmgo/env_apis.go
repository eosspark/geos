package wasmgo

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/eosspark/eos-go/chain/types"
	. "github.com/eosspark/eos-go/chain/types/generated_containers"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/eos_math"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func readActionData(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	bufferSize := int(frame.Locals[1])

	data := w.context.GetActionData()
	s := len(data)
	if bufferSize == 0 {
		w.ilog.Debug("action data size:%d", s)
		return int64(s)
	}

	copySize := min(bufferSize, s)
	w.ilog.Debug("action data:%v size:%d", data, copySize)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return int64(copySize)

}

func actionDataSize(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	size := len(w.context.GetActionData())
	w.ilog.Debug("actionDataSize:%d", size)
	return int64(size)

}

func currentReceiver(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	receiver := w.context.GetReceiver()
	w.ilog.Debug("currentReceiver:%v", receiver)
	return int64(receiver)
}

func requireAuthorization(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")
	account := int64(vm.GetCurrentFrame().Locals[0])
	w.ilog.Debug("account:%v", common.AccountName(account))
	w.context.RequireAuthorization(account)
	return 0
}

func hasAuthorization(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	account := vm.GetCurrentFrame().Locals[0]
	ret := w.context.HasAuthorization(account)
	w.ilog.Debug("account:%v authorization:%v", common.AccountName(account), ret)
	if ret {
		return 1
	} else {
		return 0
	}
}

func requireAuth2(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	frame := vm.GetCurrentFrame()
	account := frame.Locals[0]
	permission := frame.Locals[1]

	w.ilog.Debug("account:%v permission:%v", common.AccountName(account), common.PermissionName(permission))
	w.context.RequireAuthorization2(account, permission)
	return 0
}

func requireRecipient(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	recipient := vm.GetCurrentFrame().Locals[0]

	w.ilog.Debug("recipient:%v ", common.AccountName(recipient))
	w.context.RequireRecipient(recipient)

	return 0
}

func isAccount(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	account := vm.GetCurrentFrame().Locals[0]
	ret := w.context.IsAccount(account)
	w.ilog.Debug("account:%v isAccount:%v", common.AccountName(account), ret)
	if ret {
		return 1
	} else {
		return 0
	}
}

func ashlti3(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	low := frame.Locals[1]
	high := frame.Locals[2]
	shift := int(frame.Locals[3])

	i := eos_math.Int128{Low: uint64(low), High: uint64(high)}
	i.LeftShifts(shift)

	data, _ := rlp.EncodeToBytes(i)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func ashrti3(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	low := frame.Locals[1]
	high := frame.Locals[2]
	shift := int(frame.Locals[3])

	i := eos_math.Int128{Low: uint64(low), High: uint64(high)}
	i.RightShifts(shift)

	data, _ := rlp.EncodeToBytes(i)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func lshlti3(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	low := frame.Locals[1]
	high := frame.Locals[2]
	shift := int(frame.Locals[3])

	i := eos_math.Int128{Low: uint64(low), High: uint64(high)}
	i.LeftShifts(shift)

	data, _ := rlp.EncodeToBytes(i)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func lshrti3(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	low := frame.Locals[1]
	high := frame.Locals[2]
	shift := int(frame.Locals[3])

	i := eos_math.Uint128{Low: uint64(low), High: uint64(high)}
	i.RightShifts(shift)

	data, _ := rlp.EncodeToBytes(i)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func divti3(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	la := frame.Locals[1]
	ha := frame.Locals[2]
	lb := frame.Locals[3]
	hb := frame.Locals[4]

	lhs := eos_math.Int128{Low: uint64(la), High: uint64(ha)}
	rhs := eos_math.Int128{Low: uint64(lb), High: uint64(hb)}
	EosAssert(!rhs.IsZero(), &ArithmeticException{}, "divide by zero")
	quotient, _ := lhs.Div(rhs)

	data, _ := rlp.EncodeToBytes(quotient)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func udivti3(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	la := frame.Locals[1]
	ha := frame.Locals[2]
	lb := frame.Locals[3]
	hb := frame.Locals[4]

	lhs := eos_math.Uint128{Low: uint64(la), High: uint64(ha)}
	rhs := eos_math.Uint128{Low: uint64(lb), High: uint64(hb)}

	EosAssert(!rhs.IsZero(), &ArithmeticException{}, "divide by zero")
	quotient, _ := lhs.Div(rhs)

	data, _ := rlp.EncodeToBytes(quotient)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func multi3(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	la := frame.Locals[1]
	ha := frame.Locals[2]
	lb := frame.Locals[3]
	hb := frame.Locals[4]

	lhs := eos_math.Int128{Low: uint64(la), High: uint64(ha)}
	rhs := eos_math.Int128{Low: uint64(lb), High: uint64(hb)}

	data, _ := rlp.EncodeToBytes(lhs.Mul(rhs))
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0

}

func modti3(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	la := frame.Locals[1]
	ha := frame.Locals[2]
	lb := frame.Locals[3]
	hb := frame.Locals[4]

	lhs := eos_math.Int128{High: uint64(ha), Low: uint64(la)}
	rhs := eos_math.Int128{High: uint64(hb), Low: uint64(lb)}
	EosAssert(!rhs.IsZero(), &ArithmeticException{}, "divide by zero")

	_, remainder := lhs.Div(rhs)
	data, _ := rlp.EncodeToBytes(remainder)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func umodti3(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	la := frame.Locals[1]
	ha := frame.Locals[2]
	lb := frame.Locals[3]
	hb := frame.Locals[4]

	lhs := eos_math.Uint128{Low: uint64(la), High: uint64(ha)}
	rhs := eos_math.Uint128{Low: uint64(lb), High: uint64(hb)}

	EosAssert(!rhs.IsZero(), &ArithmeticException{}, "divide by zero")
	_, remainder := lhs.Div(rhs)
	data, _ := rlp.EncodeToBytes(remainder)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func addtf3(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	la := frame.Locals[1]
	ha := frame.Locals[2]
	lb := frame.Locals[3]
	hb := frame.Locals[4]

	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}

	data, _ := rlp.EncodeToBytes(a.Add(b))
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func subtf3(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	la := frame.Locals[1]
	ha := frame.Locals[2]
	lb := frame.Locals[3]
	hb := frame.Locals[4]

	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}

	data, _ := rlp.EncodeToBytes(a.Sub(b))
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func multf3(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	la := frame.Locals[1]
	ha := frame.Locals[2]
	lb := frame.Locals[3]
	hb := frame.Locals[4]

	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}

	data, _ := rlp.EncodeToBytes(a.Mul(b))
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func divtf3(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	la := frame.Locals[1]
	ha := frame.Locals[2]
	lb := frame.Locals[3]
	hb := frame.Locals[4]

	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}

	data, _ := rlp.EncodeToBytes(a.Div(b))
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func negtf2(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	la := frame.Locals[1]
	ha := frame.Locals[2]

	high := uint64(ha)
	high ^= uint64(1) << 63
	f128 := eos_math.Float128{Low: uint64(la), High: high}

	data, _ := rlp.EncodeToBytes(f128)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

//func extendsftf2(w *WasmGo, ret int, f float32) { //TODO f float??
func extendsftf2(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	f := int32(frame.Locals[1])

	//f32 := eos_math.Float32(math.Float32bits(f))
	f32 := eos_math.Float32(f)
	f128 := eos_math.F32ToF128(f32)

	data, _ := rlp.EncodeToBytes(f128)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

//func extenddftf2(w *WasmGo, ret int, d float64) { //TODO d double??
func extenddftf2(vm *VirtualMachine) int64 { //TODO d double??

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	d := uint64(frame.Locals[1])

	//f64 := eos_math.Float64(math.Float64bits(d))
	f64 := eos_math.Float64(d)
	f128 := eos_math.F64ToF128(f64)

	data, _ := rlp.EncodeToBytes(f128)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

//func trunctfdf2(w *WasmGo, l, h int64) float64 { //TODO double??
func trunctfdf2(vm *VirtualMachine) int64 { //TODO double??

	frame := vm.GetCurrentFrame()
	l := frame.Locals[0]
	h := frame.Locals[1]

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	f64 := eos_math.F128ToF64(f128)
	//return math.Float64frombits(uint64(f64))
	//vm.pushUint64(uint64(f64))

	return int64(f64)
}

//func trunctfsf2(w *WasmGo, l, h int64) float32 { //TODO float??
func trunctfsf2(vm *VirtualMachine) int64 { //TODO float??
	frame := vm.GetCurrentFrame()
	l := frame.Locals[0]
	h := frame.Locals[1]

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	f32 := eos_math.F128ToF32(f128)
	//return math.Float32frombits(uint32(f32))
	//vm.pushUint64(uint64(f32))
	return int64(f32)
}

func fixtfsi(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	l := frame.Locals[0]
	h := frame.Locals[1]

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	//return int(eos_math.F128ToI32(f128, 0, false))
	//vm.pushUint64(uint64(eos_math.F128ToI32(f128, 0, false)))

	return int64(eos_math.F128ToI32(f128, 0, false))
}

func fixtfdi(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	l := frame.Locals[0]
	h := frame.Locals[1]

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	//return eos_math.F128ToI64(f128, 0, false)
	//vm.pushUint64(uint64(eos_math.F128ToI64(f128, 0, false)))
	return int64(eos_math.F128ToI64(f128, 0, false))
}

func fixunstfsi(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	l := frame.Locals[0]
	h := frame.Locals[1]

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	//return int(eos_math.F128ToUi32(f128, 0, false))
	//vm.pushUint64(uint64(eos_math.F128ToUi32(f128, 0, false)))
	return int64(eos_math.F128ToUi32(f128, 0, false))
}

func fixunstfdi(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	l := frame.Locals[0]
	h := frame.Locals[1]

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	//return int64(eos_math.F128ToUi64(f128, 0, false))
	//vm.pushUint64(uint64(eos_math.F128ToUi64(f128, 0, false)))
	return int64(eos_math.F128ToUi64(f128, 0, false))
}

func fixtfti(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	l := frame.Locals[1]
	h := frame.Locals[2]

	f128 := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	int128 := eos_math.Fixtfti(f128)

	data, _ := rlp.EncodeToBytes(int128)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func fixunstfti(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	l := frame.Locals[1]
	h := frame.Locals[2]

	f := eos_math.Float128{Low: uint64(l), High: uint64(h)}
	uint128 := eos_math.Fixunstfti(f)

	data, _ := rlp.EncodeToBytes(uint128)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func fixsfti(vm *VirtualMachine) int64 { //TODO float??
	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	a := uint32(frame.Locals[1])

	//int128 := eos_math.Fixsfti(math.Float32bits(a))
	int128 := eos_math.Fixsfti(a)
	data, _ := rlp.EncodeToBytes(int128)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func fixdfti(vm *VirtualMachine) int64 { //TODO double??

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	a := uint64(frame.Locals[1])
	//int128 := eos_math.Fixdfti(math.Float64bits(a))
	int128 := eos_math.Fixdfti(a)
	data, _ := rlp.EncodeToBytes(int128)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func fixunssfti(vm *VirtualMachine) int64 { //TODO float??
	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	a := uint32(frame.Locals[1])

	//uint128 := eos_math.Fixunssfti(math.Float32bits(a))
	uint128 := eos_math.Fixunssfti(a)
	data, _ := rlp.EncodeToBytes(uint128)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func fixunsdfti(vm *VirtualMachine) int64 { //TODO double??
	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	a := uint64(frame.Locals[1])

	//uint128 := eos_math.Fixunsdfti(math.Float64bits(a))
	uint128 := eos_math.Fixunsdfti(a)
	data, _ := rlp.EncodeToBytes(uint128)
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func floatsidf(vm *VirtualMachine) int64 { //TODO double??
	frame := vm.GetCurrentFrame()
	i := int32(frame.Locals[0])

	//return math.Float64frombits(uint64(eos_math.I32ToF64(int32(i))))
	//vm.pushUint64(uint64(eos_math.I32ToF64(int32(i))))
	return int64(eos_math.I32ToF64(i))
}

func floatsitf(vm *VirtualMachine) int64 {

	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	i := int32(frame.Locals[1])

	data := eos_math.I32ToF128(i).Bytes()
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func floatditf(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	a := frame.Locals[1]

	data := eos_math.I64ToF128(a).Bytes()
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func floatunsitf(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	i := uint32(frame.Locals[1])

	data := eos_math.Ui32ToF128(i).Bytes()
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func floatunditf(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	ret := int(frame.Locals[0])
	a := uint64(frame.Locals[1])

	data := eos_math.Ui64ToF128(a).Bytes()
	bufferSize := len(data)
	copy(vm.Memory[ret:ret+bufferSize], data[0:bufferSize])

	return 0
}

func floattidf(vm *VirtualMachine) int64 { //TODO double
	frame := vm.GetCurrentFrame()
	l := uint64(frame.Locals[0])
	h := uint64(frame.Locals[1])

	v := eos_math.Int128{Low: l, High: h}
	//return eos_math.Floattidf(v)
	//vm.pushFloat64(eos_math.Floattidf(v))
	return int64(math.Float64bits(eos_math.Floattidf(v)))
}

func floatuntidf(vm *VirtualMachine) int64 { //TODO double
	frame := vm.GetCurrentFrame()
	l := uint64(frame.Locals[0])
	h := uint64(frame.Locals[1])

	v := eos_math.Uint128{Low: l, High: h}
	// return eos_math.Floatuntidf(v)
	//vm.pushFloat64(eos_math.Floatuntidf(v))
	return int64(math.Float64bits(eos_math.Floatuntidf(v)))
}

func cmptf2(vm *VirtualMachine) int64 { //TODO unsame with regist

	// return_value_if_nan := int(vm.popUint64())
	// hb := int64(vm.popUint64())
	// lb := int64(vm.popUint64())
	// ha := int64(vm.popUint64())
	// la := int64(vm.popUint64())

	frame := vm.GetCurrentFrame()
	la := frame.Locals[0]
	ha := frame.Locals[1]
	lb := frame.Locals[2]
	hb := frame.Locals[3]
	return_value_if_nan := int(frame.Locals[4])

	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}
	if _unordtf2(la, ha, lb, hb) != 0 {
		//vm.pushUint64(uint64(return_value_if_nan))
		return int64(return_value_if_nan)
		//return return_value_if_nan
	}
	if a.F128Lt(b) {
		// vm.pushInt64(-1)
		// return
		return -1
	}
	if a.F128EQ(b) {
		// vm.pushUint64(0)
		// return
		return 0
	}
	//vm.pushUint64(1)
	return 1
}

func eqtf2(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := frame.Locals[0]
	ha := frame.Locals[1]
	lb := frame.Locals[2]
	hb := frame.Locals[3]

	//vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, 1)))
	return _cmptf2(la, ha, lb, hb, 1)
}

func netf2(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := frame.Locals[0]
	ha := frame.Locals[1]
	lb := frame.Locals[2]
	hb := frame.Locals[3]

	//vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, 1)))
	return _cmptf2(la, ha, lb, hb, 1)
}

func getf2(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := frame.Locals[0]
	ha := frame.Locals[1]
	lb := frame.Locals[2]
	hb := frame.Locals[3]

	//vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, -1)))
	return _cmptf2(la, ha, lb, hb, -1)
}

func gttf2(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := frame.Locals[0]
	ha := frame.Locals[1]
	lb := frame.Locals[2]
	hb := frame.Locals[3]

	//vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, 0)))

	return _cmptf2(la, ha, lb, hb, 0)
}

func letf2(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := frame.Locals[0]
	ha := frame.Locals[1]
	lb := frame.Locals[2]
	hb := frame.Locals[3]

	//vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, 1)))

	return _cmptf2(la, ha, lb, hb, 1)
}

func lttf2(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := frame.Locals[0]
	ha := frame.Locals[1]
	lb := frame.Locals[2]
	hb := frame.Locals[3]

	//vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, 0)))
	return _cmptf2(la, ha, lb, hb, 0)
}

func __cmptf2(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := frame.Locals[0]
	ha := frame.Locals[1]
	lb := frame.Locals[2]
	hb := frame.Locals[3]

	//vm.pushUint64(uint64(_cmptf2(la, ha, lb, hb, 1)))
	return _cmptf2(la, ha, lb, hb, 1)
}

func unordtf2(vm *VirtualMachine) int64 {
	frame := vm.GetCurrentFrame()
	la := frame.Locals[0]
	ha := frame.Locals[1]
	lb := frame.Locals[2]
	hb := frame.Locals[3]

	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}
	if a.IsNan() || b.IsNan() {
		//vm.pushUint64(1)
		return 1
	}
	//vm.pushUint64(0)
	return 0
}

func _cmptf2(la, ha, lb, hb int64, return_value_if_nan int) int64 { //TODO unsame with regist
	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}
	if _unordtf2(la, ha, lb, hb) != 0 {
		return int64(return_value_if_nan)
	}
	if a.F128Lt(b) {
		return -1
	}
	if a.F128EQ(b) {
		return 0
	}
	return 1
}

func _unordtf2(la, ha, lb, hb int64) int64 {
	a := eos_math.Float128{Low: uint64(la), High: uint64(ha)}
	b := eos_math.Float128{Low: uint64(lb), High: uint64(hb)}
	if a.IsNan() || b.IsNan() {
		return 1
	}
	return 0
}

func assertRecoverKey(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	digest := int(frame.Locals[0])
	sig := int(frame.Locals[1])
	siglen := int(frame.Locals[2])
	pub := int(frame.Locals[3])
	publen := int(frame.Locals[4])

	digBytes := vm.Memory[digest : digest+32] //getSha256(vm, digest)
	sigBytes := vm.Memory[sig : sig+siglen]   //getMemory(vm, sig, siglen)
	pubBytes := vm.Memory[pub : pub+publen]   //getMemory(vm, pub, publen)

	w.ilog.Debug("digest:%v signature:%v publickey:%v", digBytes, sigBytes, pubBytes)

	s := ecc.NewSigNil()
	p := ecc.NewPublicKeyNil()

	d := digBytes
	rlp.DecodeBytes(sigBytes, s)
	rlp.DecodeBytes(pubBytes, p)

	check, err := s.PublicKey(d)
	EosAssert(err == nil, &CryptoApiException{}, "can not get the right publickey from digest")
	EosAssert(strings.Compare(check.String(), p.String()) == 0, &CryptoApiException{}, "Error expected key different than recovered key")

	return 0

}

func recoverKey(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	digest := int(frame.Locals[0])
	sig := int(frame.Locals[1])
	siglen := int(frame.Locals[2])
	pub := int(frame.Locals[3])
	publen := int(frame.Locals[4])

	//digBytes := getSha256(vm, digest)
	//sigBytes := getMemory(vm, sig, siglen)
	digBytes := vm.Memory[digest : digest+32] //getSha256(vm, digest)
	sigBytes := vm.Memory[sig : sig+siglen]   //getMemory(vm, sig, siglen)

	s := ecc.NewSigNil()
	rlp.DecodeBytes(sigBytes, s)
	check, _ := s.PublicKey(digBytes)

	data, err := rlp.EncodeToBytes(check)
	if err != nil {
		return -1
		//vm.pushInt64(-1)
		//return
	}

	w.ilog.Debug("digest:%v signature:%v publickey:%v", digBytes, sigBytes, data)

	bufferSize := len(data)
	if bufferSize > publen {
		bufferSize = publen
	}

	copy(vm.Memory[pub:pub+bufferSize], data[0:bufferSize])
	return int64(bufferSize)

	// setMemory(vm, pub, p, 0, l)
	// vm.pushUint64(uint64(l))
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

func assertSha256(vm *VirtualMachine) int64 {

	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	data := int(frame.Locals[0])
	dataLen := int(frame.Locals[1])
	hashVal := int(frame.Locals[2])

	dataBytes := vm.Memory[data : data+dataLen] //getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewSha256()
	hashEncode := encode(w, s, dataBytes, dataLen)
	hash := vm.Memory[hashVal : hashVal+32] //getSha256(vm, hashVal)

	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)
	EosAssert(bytes.Compare(hashEncode, hash) == 0, &CryptoApiException{}, "sha256 hash mismatch")

	return 0

}

func assertSha1(vm *VirtualMachine) int64 {

	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	data := int(frame.Locals[0])
	dataLen := int(frame.Locals[1])
	hashVal := int(frame.Locals[2])

	dataBytes := vm.Memory[data : data+dataLen] //getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewSha1()
	hashEncode := encode(w, s, dataBytes, dataLen)
	hash := vm.Memory[hashVal : hashVal+20] //hash := getSha1(vm, hashVal)

	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)
	EosAssert(bytes.Compare(hashEncode, hash) == 0, &CryptoApiException{}, "sha1 hash mismatch")

	return 0
}

func assertSha512(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	data := int(frame.Locals[0])
	dataLen := int(frame.Locals[1])
	hashVal := int(frame.Locals[2])

	dataBytes := vm.Memory[data : data+dataLen] //getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewSha512()
	hashEncode := encode(w, s, dataBytes, dataLen)
	hash := vm.Memory[hashVal : hashVal+64] //hash := getSha512(vm, hashVal)

	//w.ilog.Debug("encoded:%#v hash:%#v data:%#v", hashEncode, hash, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)
	EosAssert(bytes.Compare(hashEncode, hash) == 0, &CryptoApiException{}, "sha512 hash mismatch")

	return 0

}

func assertRipemd160(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	data := int(frame.Locals[0])
	dataLen := int(frame.Locals[1])
	hashVal := int(frame.Locals[2])

	dataBytes := vm.Memory[data : data+dataLen] //getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewRipemd160()
	hashEncode := encode(w, s, dataBytes, dataLen)
	hash := vm.Memory[hashVal : hashVal+20] //hash := getRipemd160(vm, hashVal)

	//w.ilog.Debug("encoded:%#v hash:%#v data:%#v", hashEncode, hash, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)
	EosAssert(bytes.Compare(hashEncode, hash) == 0, &CryptoApiException{}, "ripemd160 hash mismatch")

	return 0
}

func sha1(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	data := int(frame.Locals[0])
	dataLen := int(frame.Locals[1])
	hashVal := int(frame.Locals[2])

	dataBytes := vm.Memory[data : data+dataLen] //getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewSha1()
	hashEncode := encode(w, s, dataBytes, dataLen)
	copy(vm.Memory[hashVal:hashVal+20], hashEncode[0:20]) //setSha1(vm, hashVal, hashEncode)

	//w.ilog.Debug("encoded:%#v data:%#v", hashEncode, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)

	return 0
}

func sha256(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	data := int(frame.Locals[0])
	dataLen := int(frame.Locals[1])
	hashVal := int(frame.Locals[2])

	dataBytes := vm.Memory[data : data+dataLen] //getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewSha256()

	hashEncode := encode(w, s, dataBytes, dataLen)
	copy(vm.Memory[hashVal:hashVal+32], hashEncode[0:32]) //setSha256(vm, hashVal, hashEncode)

	//w.ilog.Debug("encoded:%#v data:%#v", hashEncode, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hex.EncodeToString(hashEncode), dataBytes)

	return 0
}

func sha512(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	data := int(frame.Locals[0])
	dataLen := int(frame.Locals[1])
	hashVal := int(frame.Locals[2])

	dataBytes := vm.Memory[data : data+dataLen] //getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewSha512()

	hashEncode := encode(w, s, dataBytes, dataLen)
	copy(vm.Memory[hashVal:hashVal+64], hashEncode[0:64]) //setSha512(vm, hashVal, hashEncode)

	//w.ilog.Debug("encoded:%#v data:%#v", hashEncode, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)

	return 0
}

func ripemd160(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	data := int(frame.Locals[0])
	dataLen := int(frame.Locals[1])
	hashVal := int(frame.Locals[2])

	dataBytes := vm.Memory[data : data+dataLen] //getMemory(vm, data, dataLen)
	// if dataBytes == nil {
	// 	return
	// }

	s := crypto.NewRipemd160()
	hashEncode := encode(w, s, dataBytes, dataLen)
	copy(vm.Memory[hashVal:hashVal+20], hashEncode[0:20]) //setRipemd160(vm, hashVal, hashEncode)

	//w.ilog.Debug("encoded:%#v data:%#v", hashEncode, dataBytes)
	w.ilog.Debug("encoded:%v data:%v", hashEncode, dataBytes)

	return 0
}

func dbStoreI64(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	frame := vm.GetCurrentFrame()
	scope := uint64(frame.Locals[0])
	table := uint64(frame.Locals[1])
	payer := uint64(frame.Locals[2])
	id := uint64(frame.Locals[3])
	buffer := int(frame.Locals[4])
	bufferSize := int(frame.Locals[5])

	bytes := vm.Memory[buffer : buffer+bufferSize] //getMemory(vm, data, dataLen)

	iterator := w.context.DbStoreI64(scope, table, payer, id, bytes)
	//iterator := 0
	//vm.pushUint64(uint64(iterator))
	w.ilog.Debug("scope:%v table:%v payer:%v id:%d data:%v iterator:%d",
		common.ScopeName(scope), common.TableName(table), common.AccountName(payer), id, bytes, iterator)

	return int64(iterator)
}

func dbUpdateI64(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	iterator := int(frame.Locals[0])
	payer := uint64(frame.Locals[1])
	buffer := int(frame.Locals[2])
	bufferSize := int(frame.Locals[3])

	bytes := vm.Memory[buffer : buffer+bufferSize] //getMemory(vm, data, dataLen)

	w.context.DbUpdateI64(iterator, payer, bytes)
	w.ilog.Debug("data:%v iterator:%d payer:%v ", bytes, iterator, common.AccountName(payer))

	return 0
}

func dbRemoveI64(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	iterator := int(frame.Locals[0])

	w.context.DbRemoveI64(iterator)
	w.ilog.Debug("iterator:%d", iterator)

	return 0
}

func dbGetI64(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	iterator := int(frame.Locals[0])
	buffer := int(frame.Locals[1])
	bufferSize := int(frame.Locals[2])

	bytes := make([]byte, bufferSize)
	size := w.context.DbGetI64(iterator, bytes, bufferSize)
	if bufferSize == 0 {
		//vm.pushUint64(uint64(size))
		w.ilog.Debug("iterator:%d size:%d", iterator, size)
		//return

		return int64(size)
	}
	//setMemory(vm, buffer, bytes, 0, size)
	//vm.pushUint64(uint64(size))

	copy(vm.Memory[buffer:buffer+size], bytes[0:size])
	w.ilog.Debug("iterator:%d data:%v size:%d", iterator, bytes, size)
	return int64(size)
}

func dbNextI64(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	// primary := int(vm.popUint64())
	// itr := int(vm.popUint64())

	frame := vm.GetCurrentFrame()
	itr := int(frame.Locals[0])
	primary := int(frame.Locals[1])

	if itr < -1 {
		//itr = -1
		//vm.pushUint64(uint64(itr))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, itr)
		return -1
	}

	var p uint64
	iterator := w.context.DbNextI64(itr, &p)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d primary:%d", itr, iterator, p)
		return int64(iterator)
	}

	//setUint64(vm, primary, p)
	//vm.pushUint64(uint64(iterator))
	data, _ := rlp.EncodeToBytes(p)
	copy(vm.Memory[primary:primary+8], data[0:8])
	w.ilog.Debug("iterator:%d nextIterator:%d primary:%d", itr, iterator, p)

	return int64(iterator)
}

func dbPreviousI64(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	// primary := int(vm.popUint64())
	// itr := int(vm.popUint64())

	frame := vm.GetCurrentFrame()
	itr := int(frame.Locals[0])
	primary := int(frame.Locals[1])

	var p uint64
	iterator := w.context.DbPreviousI64(itr, &p)
	w.ilog.Debug("dbNextI64 iterator:%d", iterator)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return int64(iterator)
	}

	//setUint64(vm, primary, p)
	//vm.pushUint64(uint64(iterator))
	data, _ := rlp.EncodeToBytes(p)
	copy(vm.Memory[primary:primary+8], data[0:8])
	w.ilog.Debug("iterator:%d priviousIterator:%d primary:%d", itr, iterator, p)
	return int64(iterator)
}

func dbFindI64(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	// id := vm.popUint64()
	// table := vm.popUint64()
	// scope := vm.popUint64()
	// code := vm.popUint64()

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	id := uint64(frame.Locals[3])

	iterator := w.context.DbFindI64(code, scope, table, id)
	//iterator := -1
	//vm.pushUint64(uint64(iterator))
	w.ilog.Debug("code:%v scope:%v table:%v id:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), id, iterator)

	return int64(iterator)

}

func dbLowerboundI64(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	id := uint64(frame.Locals[3])

	iterator := w.context.DbLowerboundI64(code, scope, table, id)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v id:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), id, iterator)

	return int64(iterator)
}

func dbUpperboundI64(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	id := uint64(frame.Locals[3])

	iterator := w.context.DbUpperboundI64(code, scope, table, id)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v id:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), id, iterator)

	return int64(iterator)
}

func dbEndI64(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])

	iterator := w.context.DbEndI64(code, scope, table)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), iterator)

	return int64(iterator)
}

//secondaryKey Index
func dbIdx64Store(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	frame := vm.GetCurrentFrame()
	scope := uint64(frame.Locals[0])
	table := uint64(frame.Locals[1])
	payer := uint64(frame.Locals[2])
	id := uint64(frame.Locals[3])
	pValue := int(frame.Locals[4])

	secondaryKey := getUint64(vm, pValue)
	// var secondaryKey uint64
	// c := vm.Memory[pValue:pValue+8]
	// rlp.DecodeBytes(c, &secondaryKey)

	iterator := w.context.Idx64Store(scope, table, payer, id, &secondaryKey)

	w.ilog.Debug("scope:%v table:%v payer:%v id:%d secondaryKey:%d iterator:%d",
		common.ScopeName(scope), common.TableName(table), common.AccountName(payer), id, secondaryKey, iterator)

	return int64(iterator)
}

func dbIdx64Remove(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	frame := vm.GetCurrentFrame()
	iterator := int(frame.Locals[0])

	w.context.Idx64Remove(iterator)
	w.ilog.Debug("iterator:%d", iterator)

	return 0
}

func dbIdx64Update(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	frame := vm.GetCurrentFrame()
	iterator := int(frame.Locals[0])
	payer := uint64(frame.Locals[1])
	pValue := int(frame.Locals[2])

	secondaryKey := getUint64(vm, pValue)
	w.context.Idx64Update(iterator, payer, &secondaryKey)

	w.ilog.Debug("payer:%v data:%v secondaryKey:%d", common.AccountName(payer), secondaryKey, iterator)

	return 0
}

func dbIdx64findSecondary(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	pPrimary := int(frame.Locals[4])

	var primaryKey uint64
	secondaryKey := getUint64(vm, pSecondary)
	iterator := w.context.Idx64FindSecondary(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}
	setUint64(vm, pPrimary, primaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)

	return int64(iterator)
}

func dbIdx64Lowerbound(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	pPrimary := int(frame.Locals[4])

	var primaryKey uint64

	secondaryKey := getUint64(vm, pSecondary)
	iterator := w.context.Idx64Lowerbound(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}

	setUint64(vm, pPrimary, primaryKey)
	setUint64(vm, pSecondary, secondaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)

	return int64(iterator)
}

func dbIdx64Upperbound(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	pPrimary := int(frame.Locals[4])

	var primaryKey uint64
	secondaryKey := getUint64(vm, pSecondary)
	iterator := w.context.Idx64Upperbound(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}
	setUint64(vm, pPrimary, primaryKey)
	setUint64(vm, pSecondary, secondaryKey)

	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)

	return int64(iterator)
}

func dbIdx64End(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])

	iterator := w.context.Idx64End(code, scope, table)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), iterator)

	return int64(iterator)
}

func dbIdx64Next(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	itr := int(frame.Locals[0])
	primary := int(frame.Locals[1])

	var p uint64
	iterator := w.context.Idx64Next(itr, &p)
	w.ilog.Debug("dbIdx64Next iterator:%d", iterator)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return int64(iterator)
	}
	setUint64(vm, primary, p)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d nextIterator:%d primary:%d", itr, iterator, p)
	return int64(iterator)
}

func dbIdx64Previous(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	itr := int(frame.Locals[0])
	primary := int(frame.Locals[1])

	var p uint64
	iterator := w.context.Idx64Previous(itr, &p)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return int64(iterator)
	}
	setUint64(vm, primary, p)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d previousIterator:%d primary:%d", itr, iterator, p)
	return int64(iterator)
}

func dbIdx64FindPrimary(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	primary := uint64(frame.Locals[4])

	var secondaryKey uint64
	iterator := w.context.Idx64FindPrimary(code, scope, table, &secondaryKey, primary)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, iterator)
		return int64(iterator)
	}
	setUint64(vm, pSecondary, secondaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d secondaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, secondaryKey, iterator)
	return int64(iterator)

}

func dbIdx128Store(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	frame := vm.GetCurrentFrame()
	scope := uint64(frame.Locals[0])
	table := uint64(frame.Locals[1])
	payer := uint64(frame.Locals[2])
	id := uint64(frame.Locals[3])
	pValue := int(frame.Locals[4])

	secondaryKey := getUint128(vm, pValue)
	iterator := w.context.Idx128Store(scope, table, payer, id, secondaryKey)

	w.ilog.Debug("scope:%v table:%v payer:%v id:%d secondaryKey:%d iterator:%d",
		common.ScopeName(scope), common.TableName(table), common.AccountName(payer), id, secondaryKey, iterator)

	return int64(iterator)
}

func dbIdx128Remove(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")
	frame := vm.GetCurrentFrame()
	iterator := int(frame.Locals[0])

	w.context.Idx128Remove(iterator)
	w.ilog.Debug("iterator:%d", iterator)

	return 0
}

func dbIdx128Update(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	frame := vm.GetCurrentFrame()
	iterator := int(frame.Locals[0])
	payer := uint64(frame.Locals[1])
	pValue := int(frame.Locals[2])

	secondaryKey := getUint128(vm, pValue)
	w.context.Idx128Update(iterator, payer, secondaryKey)

	w.ilog.Debug("payer:%v data:%v secondaryKey:%d", common.AccountName(payer), secondaryKey, iterator)

	return 0
}

func dbIdx128findSecondary(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	pPrimary := int(frame.Locals[4])

	var primaryKey uint64
	secondaryKey := getUint128(vm, pSecondary)
	iterator := w.context.Idx128FindSecondary(code, scope, table, secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}
	setUint64(vm, pPrimary, primaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)

	return int64(iterator)
}

func dbIdx128Lowerbound(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	pPrimary := int(frame.Locals[4])

	var primaryKey uint64

	secondaryKey := getUint128(vm, pSecondary)
	iterator := w.context.Idx128Lowerbound(code, scope, table, secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}

	setUint64(vm, pPrimary, primaryKey)
	setUint128(vm, pSecondary, secondaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)

	return int64(iterator)
}

func dbIdx128Upperbound(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	pPrimary := int(frame.Locals[4])

	var primaryKey uint64
	secondaryKey := getUint128(vm, pSecondary)
	iterator := w.context.Idx128Upperbound(code, scope, table, secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}
	setUint64(vm, pPrimary, primaryKey)
	setUint128(vm, pSecondary, secondaryKey)

	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)

	return int64(iterator)
}

func dbIdx128End(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])

	iterator := w.context.Idx128End(code, scope, table)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), iterator)

	return int64(iterator)
}

func dbIdx128Next(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	itr := int(frame.Locals[0])
	primary := int(frame.Locals[1])

	var p uint64
	iterator := w.context.Idx128Next(itr, &p)
	w.ilog.Debug("dbIdx128Next iterator:%d", iterator)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return int64(iterator)
	}
	setUint64(vm, primary, p)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d nextIterator:%d primary:%d", itr, iterator, p)
	return int64(iterator)
}

func dbIdx128Previous(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	itr := int(frame.Locals[0])
	primary := int(frame.Locals[1])

	var p uint64
	iterator := w.context.Idx128Previous(itr, &p)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return int64(iterator)
	}
	setUint64(vm, primary, p)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d previousIterator:%d primary:%d", itr, iterator, p)
	return int64(iterator)
}

func dbIdx128FindPrimary(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	primary := uint64(frame.Locals[4])

	var secondaryKey eos_math.Uint128
	iterator := w.context.Idx128FindPrimary(code, scope, table, &secondaryKey, primary)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, iterator)
		return int64(iterator)
	}
	setUint128(vm, pSecondary, &secondaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d secondaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, secondaryKey, iterator)
	return int64(iterator)
}

func dbIdx256Store(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	frame := vm.GetCurrentFrame()
	scope := uint64(frame.Locals[0])
	table := uint64(frame.Locals[1])
	payer := uint64(frame.Locals[2])
	id := uint64(frame.Locals[3])
	pValue := int(frame.Locals[4])
	dataLen := int(frame.Locals[5])

	EosAssert(dataLen == 2, &DbApiException{},
		"invalid size of secondary key array for Idx256: given %d bytes but expected %d bytes", dataLen, 2)

	secondaryKey := getUint256(vm, pValue)
	w.ilog.Debug("scope:%v table:%v payer:%v id:%d secondaryKey:%d",
		common.ScopeName(scope), common.TableName(table), common.AccountName(payer), id, secondaryKey)

	iterator := w.context.Idx256Store(scope, table, payer, id, secondaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("scope:%v table:%v payer:%v id:%d secondaryKey:%d iterator:%d",
		common.ScopeName(scope), common.TableName(table), common.AccountName(payer), id, secondaryKey, iterator)

	return int64(iterator)
}

func dbIdx256Remove(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")
	frame := vm.GetCurrentFrame()
	iterator := int(frame.Locals[0])

	w.context.Idx256Remove(iterator)
	w.ilog.Debug("iterator:%d", iterator)

	return 0
}

func dbIdx256Update(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	frame := vm.GetCurrentFrame()
	iterator := int(frame.Locals[0])
	payer := uint64(frame.Locals[1])
	pValue := int(frame.Locals[2])
	dataLen := int(frame.Locals[3])
	EosAssert(dataLen == 2, &DbApiException{},
		"invalid size of secondary key array for Idx256: given %d bytes but expected %d bytes", dataLen, 2)

	secondaryKey := getUint256(vm, pValue)
	w.context.Idx256Update(iterator, payer, secondaryKey)

	w.ilog.Debug("payer:%v data:%v secondaryKey:%d", common.AccountName(payer), secondaryKey, iterator)

	return 0
}

func dbIdx256findSecondary(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	dataLen := int(frame.Locals[4])
	pPrimary := int(frame.Locals[5])
	EosAssert(dataLen == 2, &DbApiException{},
		"invalid size of secondary key array for Idx256: given %d bytes but expected %d bytes", dataLen, 2)

	var primaryKey uint64
	secondaryKey := getUint256(vm, pSecondary)
	iterator := w.context.Idx256FindSecondary(code, scope, table, secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}
	setUint64(vm, pPrimary, primaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)
	return int64(iterator)
}

func dbIdx256Lowerbound(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	// pPrimary := int(vm.popUint64())
	// dataLen := int(vm.popUint64())
	// EosAssert(dataLen == 2, &DbApiException{},
	// 	"invalid size of secondary key array for Idx256: given %d bytes but expected %d bytes", dataLen, 2)
	// pSecondary := int(vm.popUint64())
	// table := vm.popUint64()
	// scope := vm.popUint64()
	// code := vm.popUint64()

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	dataLen := int(frame.Locals[4])
	pPrimary := int(frame.Locals[5])
	EosAssert(dataLen == 2, &DbApiException{},
		"invalid size of secondary key array for Idx256: given %d bytes but expected %d bytes", dataLen, 2)

	var primaryKey uint64

	secondaryKey := getUint256(vm, pSecondary)
	iterator := w.context.Idx256Lowerbound(code, scope, table, secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}

	setUint64(vm, pPrimary, primaryKey)
	setUint256(vm, pSecondary, secondaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)

	return int64(iterator)
}

func dbIdx256Upperbound(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	dataLen := int(frame.Locals[4])
	pPrimary := int(frame.Locals[5])
	EosAssert(dataLen == 2, &DbApiException{},
		"invalid size of secondary key array for Idx256: given %d bytes but expected %d bytes", dataLen, 2)

	var primaryKey uint64
	secondaryKey := getUint256(vm, pSecondary)
	iterator := w.context.Idx256Upperbound(code, scope, table, secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}
	setUint64(vm, pPrimary, primaryKey)
	setUint256(vm, pSecondary, secondaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%d primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)

	return int64(iterator)
}

func dbIdx256End(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])

	iterator := w.context.Idx256End(code, scope, table)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), iterator)

	return int64(iterator)
}

func dbIdx256Next(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	itr := int(frame.Locals[0])
	primary := int(frame.Locals[1])

	var p uint64
	iterator := w.context.Idx256Next(itr, &p)
	w.ilog.Debug("dbIdx256Next iterator:%d", iterator)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return int64(iterator)
	}
	setUint64(vm, primary, p)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d nextIterator:%d primary:%d", itr, iterator, p)
	return int64(iterator)
}

func dbIdx256Previous(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	itr := int(frame.Locals[0])
	primary := int(frame.Locals[1])

	var p uint64
	iterator := w.context.Idx256Previous(itr, &p)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return int64(iterator)
	}
	setUint64(vm, primary, p)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d previousIterator:%d primary:%d", itr, iterator, p)
	return int64(iterator)
}

func dbIdx256FindPrimary(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	dataLen := int(frame.Locals[4])
	primary := uint64(frame.Locals[5])
	EosAssert(dataLen == 2, &DbApiException{},
		"invalid size of secondary key array for Idx256: given %d bytes but expected %d bytes", dataLen, 2)

	var secondaryKey eos_math.Uint256
	iterator := w.context.Idx256FindPrimary(code, scope, table, &secondaryKey, primary)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, iterator)
		return int64(iterator)
	}
	setUint256(vm, pSecondary, &secondaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d secondaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, secondaryKey, iterator)
	return int64(iterator)

}

func dbIdxDoubleStore(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	scope := uint64(frame.Locals[0])
	table := uint64(frame.Locals[1])
	payer := uint64(frame.Locals[2])
	id := uint64(frame.Locals[3])
	pValue := int(frame.Locals[4])

	secondaryKey := eos_math.Float64(getUint64(vm, pValue))

	f := math.Float64frombits(uint64(secondaryKey))
	EosAssert(!math.IsNaN(f), &TransactionException{}, "NaN is not an allowed value for a secondary key")

	iterator := w.context.IdxDoubleStore(scope, table, payer, id, &secondaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("scope:%v table:%v payer:%v id:%d secondaryKey:%v iterator:%d",
		common.ScopeName(scope), common.TableName(table), common.AccountName(payer), id, secondaryKey, iterator)

	return int64(iterator)
}

func dbIdxDoubleRemove(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	frame := vm.GetCurrentFrame()
	iterator := int(frame.Locals[0])

	w.context.IdxDoubleRemove(iterator)
	w.ilog.Debug("iterator:%d", iterator)

	return 0
}

func dbIdxDoubleUpdate(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	iterator := int(frame.Locals[0])
	payer := uint64(frame.Locals[1])
	pValue := int(frame.Locals[2])

	secondaryKey := eos_math.Float64(getUint64(vm, pValue))
	f := math.Float64frombits(uint64(secondaryKey))
	EosAssert(!math.IsNaN(f), &TransactionException{}, "NaN is not an allowed value for a secondary key")

	w.context.IdxDoubleUpdate(iterator, payer, &secondaryKey)
	w.ilog.Debug("payer:%v secondaryKey:%v iterator:%v", common.AccountName(payer), secondaryKey, iterator)

	return 0

}

func dbIdxDoublefindSecondary(vm *VirtualMachine) int64 {

	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	pPrimary := int(frame.Locals[4])

	var primaryKey uint64
	secondaryKey := eos_math.Float64(getUint64(vm, pSecondary))
	f := math.Float64frombits(uint64(secondaryKey))
	EosAssert(!math.IsNaN(f), &TransactionException{}, "NaN is not an allowed value for a secondary key")

	iterator := w.context.IdxDoubleFindSecondary(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}
	setUint64(vm, pPrimary, primaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)

	return int64(iterator)
}

func dbIdxDoubleLowerbound(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	pPrimary := int(frame.Locals[4])

	var primaryKey uint64
	secondaryKey := eos_math.Float64(getUint64(vm, pSecondary))
	f := math.Float64frombits(uint64(secondaryKey))
	EosAssert(!math.IsNaN(f), &TransactionException{}, "NaN is not an allowed value for a secondary key")

	iterator := w.context.IdxDoubleLowerbound(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}
	setUint64(vm, pPrimary, primaryKey)
	setUint64(vm, pSecondary, uint64(secondaryKey))
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)

	return int64(iterator)
}

func dbIdxDoubleUpperbound(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	pPrimary := int(frame.Locals[4])

	var primaryKey uint64
	secondaryKey := eos_math.Float64(getUint64(vm, pSecondary))
	f := math.Float64frombits(uint64(secondaryKey))
	EosAssert(!math.IsNaN(f), &TransactionException{}, "NaN is not an allowed value for a secondary key")

	iterator := w.context.IdxDoubleUpperbound(code, scope, table, &secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}
	setUint64(vm, pPrimary, primaryKey)
	setUint64(vm, pSecondary, uint64(secondaryKey))
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)

	return int64(iterator)
}

func dbIdxDoubleEnd(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])

	iterator := w.context.IdxDoubleEnd(code, scope, table)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), iterator)

	return int64(iterator)
}

func dbIdxDoubleNext(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	itr := int(frame.Locals[0])
	primary := int(frame.Locals[1])

	var p uint64
	iterator := w.context.IdxDoubleNext(itr, &p)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return int64(iterator)
	}

	setUint64(vm, primary, p)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d nextIterator:%d primary:%d", itr, iterator, p)
	return int64(iterator)
}

func dbIdxDoublePrevious(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	itr := int(frame.Locals[0])
	primary := int(frame.Locals[1])

	var p uint64
	iterator := w.context.IdxDoublePrevious(itr, &p)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d proviousIterator:%d", itr, iterator)
		return int64(iterator)
	}
	setUint64(vm, primary, p)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d previousIterator:%d primary:%d", itr, iterator, p)
	return int64(iterator)
}

func dbIdxDoubleFindPrimary(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	primary := uint64(frame.Locals[4])

	var secondaryKey eos_math.Float64
	iterator := w.context.IdxDoubleFindPrimary(code, scope, table, &secondaryKey, primary)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, iterator)
		return int64(iterator)
	}
	setUint64(vm, pSecondary, uint64(secondaryKey))
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d secondaryKey:%v iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, secondaryKey, iterator)
	return int64(iterator)
}

func dbIdxLongDoubleStore(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	scope := uint64(frame.Locals[0])
	table := uint64(frame.Locals[1])
	payer := uint64(frame.Locals[2])
	id := uint64(frame.Locals[3])
	pValue := int(frame.Locals[4])

	secondaryKey := getFloat128(vm, pValue)

	// f := math.Float64frombits(uint64(secondaryKey))
	// EosAssert(!math.IsNaN(f), &TransactionException{}, "NaN is not an allowed value for a secondary key")

	iterator := w.context.IdxLongDoubleStore(scope, table, payer, id, secondaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("scope:%v table:%v payer:%v id:%d secondaryKey:%v iterator:%d",
		common.ScopeName(scope), common.TableName(table), common.AccountName(payer), id, secondaryKey, iterator)

	return int64(iterator)
}

func dbIdxLongDoubleRemove(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	frame := vm.GetCurrentFrame()
	iterator := int(frame.Locals[0])

	w.context.IdxLongDoubleRemove(iterator)
	w.ilog.Debug("iterator:%d", iterator)

	return 0
}

func dbIdxLongDoubleUpdate(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	iterator := int(frame.Locals[0])
	payer := uint64(frame.Locals[1])
	pValue := int(frame.Locals[2])

	// secondaryKey := eos_math.Float64(getUint64(vm, pValue))
	// f := math.Float64frombits(uint64(secondaryKey))
	// EosAssert(!math.IsNaN(f), &TransactionException{}, "NaN is not an allowed value for a secondary key")
	secondaryKey := getFloat128(vm, pValue)

	w.context.IdxLongDoubleUpdate(iterator, payer, secondaryKey)
	w.ilog.Debug("payer:%v secondaryKey:%v iterator:%v", common.AccountName(payer), secondaryKey, iterator)

	return 0

}

func dbIdxLongDoublefindSecondary(vm *VirtualMachine) int64 {

	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	pPrimary := int(frame.Locals[4])

	var primaryKey uint64
	// secondaryKey := eos_math.Float64(getUint64(vm, pSecondary))
	// f := math.Float64frombits(uint64(secondaryKey))
	// EosAssert(!math.IsNaN(f), &TransactionException{}, "NaN is not an allowed value for a secondary key")
	secondaryKey := getFloat128(vm, pSecondary)

	iterator := w.context.IdxLongDoubleFindSecondary(code, scope, table, secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}
	setUint64(vm, pPrimary, primaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)
	return int64(iterator)
}

func dbIdxLongDoubleLowerbound(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	pPrimary := int(frame.Locals[4])

	var primaryKey uint64
	// secondaryKey := eos_math.Float64(getUint64(vm, pSecondary))
	// f := math.Float64frombits(uint64(secondaryKey))
	// EosAssert(!math.IsNaN(f), &TransactionException{}, "NaN is not an allowed value for a secondary key")
	secondaryKey := getFloat128(vm, pSecondary)

	iterator := w.context.IdxLongDoubleLowerbound(code, scope, table, secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}
	setUint64(vm, pPrimary, primaryKey)
	setFloat128(vm, pSecondary, secondaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)
	return int64(iterator)
}

func dbIdxLongDoubleUpperbound(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	pPrimary := int(frame.Locals[4])

	var primaryKey uint64
	// secondaryKey := eos_math.Float64(getUint64(vm, pSecondary))
	// f := math.Float64frombits(uint64(secondaryKey))
	// EosAssert(!math.IsNaN(f), &TransactionException{}, "NaN is not an allowed value for a secondary key")
	secondaryKey := getFloat128(vm, pSecondary)

	iterator := w.context.IdxLongDoubleUpperbound(code, scope, table, secondaryKey, &primaryKey)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, iterator)
		return int64(iterator)
	}
	setUint64(vm, pPrimary, primaryKey)
	setFloat128(vm, pSecondary, secondaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v secondaryKey:%v primaryKey:%d iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), secondaryKey, primaryKey, iterator)

	return int64(iterator)
}

func dbIdxLongDoubleEnd(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])

	iterator := w.context.IdxLongDoubleEnd(code, scope, table)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), iterator)

	return int64(iterator)
}

func dbIdxLongDoubleNext(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	itr := int(frame.Locals[0])
	primary := int(frame.Locals[1])

	var p uint64
	iterator := w.context.IdxLongDoubleNext(itr, &p)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d nextIterator:%d", itr, iterator)
		return int64(iterator)
	}

	setUint64(vm, primary, p)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d nextIterator:%d primary:%d", itr, iterator, p)
	return int64(iterator)
}

func dbIdxLongDoublePrevious(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	itr := int(frame.Locals[0])
	primary := int(frame.Locals[1])

	var p uint64
	iterator := w.context.IdxLongDoublePrevious(itr, &p)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("iterator:%d proviousIterator:%d", itr, iterator)
		return int64(iterator)
	}
	setUint64(vm, primary, p)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("iterator:%d previousIterator:%d primary:%d", itr, iterator, p)
	return int64(iterator)
}

func dbIdxLongDoubleFindPrimary(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	code := uint64(frame.Locals[0])
	scope := uint64(frame.Locals[1])
	table := uint64(frame.Locals[2])
	pSecondary := int(frame.Locals[3])
	primary := uint64(frame.Locals[4])

	var secondaryKey eos_math.Float128
	iterator := w.context.IdxLongDoubleFindPrimary(code, scope, table, &secondaryKey, primary)
	if iterator <= -1 {
		//vm.pushUint64(uint64(iterator))
		w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d iterator:%d",
			common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, iterator)
		return int64(iterator)
	}

	setFloat128(vm, pSecondary, &secondaryKey)
	//vm.pushUint64(uint64(iterator))

	w.ilog.Debug("code:%v scope:%v table:%v primaryKey:%d secondaryKey:%v iterator:%d",
		common.AccountName(code), common.ScopeName(scope), common.TableName(table), primary, secondaryKey, iterator)
	return int64(iterator)
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func memcpy(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	dest := frame.Locals[0]
	src := frame.Locals[1]
	length := frame.Locals[2]

	w.ilog.Debug("dest:%d src:%d length:%d ", dest, src, length)
	EosAssert(abs(dest-src) >= length, &OverlappingMemoryError{}, "memcpy can only accept non-aliasing pointers")
	copy(vm.Memory[dest:dest+length], vm.Memory[src:src+length])
	return int64(dest)
}

func memmove(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	dest := frame.Locals[0]
	src := frame.Locals[1]
	length := frame.Locals[2]

	w.ilog.Debug("dest:%d src:%d length:%d ", dest, src, length)

	//EosAssert(abs(dest-src) >= length, &OverlappingMemoryError{}, "memmove with overlapping memory")
	copy(vm.Memory[dest:dest+length], vm.Memory[src:src+length])
	//vm.pushUint64(uint64(dest))

	return int64(dest)

}

func memcmp(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	dest := frame.Locals[0]
	src := frame.Locals[1]
	length := frame.Locals[2]

	w.ilog.Debug("dest:%d src:%d length:%d ", dest, src, length)

	ret := bytes.Compare(vm.Memory[dest:dest+length], vm.Memory[src:src+length])
	//vm.pushUint64(uint64(ret))

	return int64(ret)
}

func memset(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	// length := int(vm.popUint64())
	// value := int(vm.popUint64())
	// dest := int(vm.popUint64())

	frame := vm.GetCurrentFrame()
	dest := int(frame.Locals[0])
	value := byte(frame.Locals[1])
	length := int(frame.Locals[2])

	w.ilog.Debug("dest:%d value:%d length:%d ", dest, value, length)

	cap := cap(vm.Memory)
	if cap < dest || cap < dest+length {
		EosAssert(false, &OverlappingMemoryError{}, "memset with heap memory out of bound")
		return 0
	}

	b := bytes.Repeat([]byte{value}, length)
	copy(vm.Memory[dest:dest+length], b[:])

	//vm.pushUint64(uint64(dest))
	return int64(dest)
}

func checkTransactionAuthorization(vm *VirtualMachine) int64 {
	fmt.Println("check_transaction_authorization")

	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	trxData := int(frame.Locals[0])
	trxSize := int(frame.Locals[1])
	pubkeysData := int(frame.Locals[2])
	pubkeysSize := int(frame.Locals[3])
	permsData := int(frame.Locals[4])
	permsSize := int(frame.Locals[5])

	trxDataBytes := getMemory(vm, trxData, trxSize)
	pubkeysDataBytes := getMemory(vm, pubkeysData, pubkeysSize)
	permsDataBytes := getMemory(vm, permsData, permsSize)

	providedKeys := NewPublicKeySet()
	providedPermissions := NewPermissionLevelSet()
	trx := types.Transaction{}

	unpackProvidedKeys(providedKeys, &pubkeysDataBytes)
	unpackProvidedPermissions(providedPermissions, &permsDataBytes)
	rlp.DecodeBytes(trxDataBytes, &trx)

	w.ilog.Debug("actions:%v permission:%v providedKeys:%v providedPermissions:%v", trx.Actions, providedKeys, providedPermissions)

	returning := false
	Try(func() {
		w.context.CheckAuthorization(
			trx.Actions,
			providedKeys,
			providedPermissions,
			uint64(trx.DelaySec))
	}).Catch(func(e Exception) {
		returning = true
	}).End()

	if returning {
		//vm.pushUint64(0)
		return 0
	}

	//vm.pushUint64(1)
	return 1
}

func checkPermissionAuthorization(vm *VirtualMachine) int64 {
	fmt.Println("check_permission_authorization")

	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	account := common.AccountName(frame.Locals[0])
	permission := common.PermissionName(frame.Locals[1])
	pubkeysData := int(frame.Locals[2])
	pubkeysSize := int(frame.Locals[3])
	permsData := int(frame.Locals[4])
	permsSize := int(frame.Locals[5])
	delayUS := uint64(frame.Locals[6])

	pubkeysDataBytes := getMemory(vm, pubkeysData, pubkeysSize)
	permsDataBytes := getMemory(vm, permsData, permsSize)

	providedKeys := NewPublicKeySet()
	providedPermissions := NewPermissionLevelSet()

	unpackProvidedKeys(providedKeys, &pubkeysDataBytes)
	unpackProvidedPermissions(providedPermissions, &permsDataBytes)

	w.ilog.Debug("account:%v permission:%v providedKeys:%v providedPermissions:%v", account, permission, providedKeys, providedPermissions)

	returning := false
	Try(func() {
		w.context.CheckAuthorization2(account,
			permission,
			providedKeys,
			providedPermissions,
			delayUS)
	}).Catch(func(e Exception) {
		returning = true
	}).End()

	if returning {
		//vm.pushUint64(0)
		return 0
	}

	//vm.pushUint64(1)
	return 1
}

func unpackProvidedKeys(ps *PublicKeySet, pubkeysData *[]byte) {
	if len(*pubkeysData) == 0 {
		return
	}

	providedKey := []ecc.PublicKey{}
	rlp.DecodeBytes(*pubkeysData, &providedKey)

	for _, pk := range providedKey {
		ps.AddItem(pk)
	}

}

func unpackProvidedPermissions(ps *PermissionLevelSet, permsData *[]byte) {
	if len(*permsData) == 0 {
		return
	}

	permissions := []common.PermissionLevel{}
	rlp.DecodeBytes(*permsData, &permissions)

	for _, permission := range permissions {
		ps.AddItem(permission)
	}

}

func getPermissionLastUsed(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	permission := common.PermissionName(frame.Locals[0])
	account := common.AccountName(frame.Locals[1])

	ret := w.context.GetPermissionLastUsed(account, permission)
	//vm.pushUint64(uint64(ret.TimeSinceEpoch().Count()))

	w.ilog.Debug("account:%v permission:%v LastUsed:%v", account, permission, ret)
	return ret.TimeSinceEpoch().Count()
}

func getAccountCreationTime(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	account := common.AccountName(frame.Locals[0])

	ret := w.context.GetAccountCreateTime(account)
	//vm.pushUint64(uint64(ret.TimeSinceEpoch().Count()))

	w.ilog.Debug("account:%v creationTime:%v", account, ret)
	return ret.TimeSinceEpoch().Count()

}

func prints(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	ptr := int(vm.GetCurrentFrame().Locals[0])

	var size int
	for i := 0; i < 512; i++ {
		if vm.Memory[ptr+i] == 0 {
			break
		}
		size++
	}

	str := vm.Memory[ptr : ptr+size]
	w.context.ContextAppend(string(str))

	w.ilog.Debug("prints:%v", str)

	return 0
}

func printsl(vm *VirtualMachine) int64 {

	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	strIndex := int(frame.Locals[0])
	strLen := int(frame.Locals[1])

	str := string(getMemory(vm, strIndex, strLen))
	w.context.ContextAppend(str)

	w.ilog.Debug("prints_l:%v", str)

	return 0
}

func printi(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	val := vm.GetCurrentFrame().Locals[0]

	str := strconv.FormatInt(val, 10)
	w.context.ContextAppend(str)

	w.ilog.Debug("printi:%v", str)

	return 0
}

func printui(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	val := uint64(vm.GetCurrentFrame().Locals[0])

	str := strconv.FormatUint(val, 10)
	w.context.ContextAppend(str)

	w.ilog.Debug("printui:%v", str)
	return 0
}

func printi128(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	val := int(vm.GetCurrentFrame().Locals[0])

	bytes := getMemory(vm, val, 16)
	var v eos_math.Int128
	rlp.DecodeBytes(bytes, &v)
	str := v.String()
	w.context.ContextAppend(str)

	w.ilog.Debug("printi128:%v", str)
	return 0

}

func printui128(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	val := int(vm.GetCurrentFrame().Locals[0])

	bytes := getMemory(vm, val, 16)
	var v eos_math.Uint128
	rlp.DecodeBytes(bytes, &v)
	str := v.String()
	w.context.ContextAppend(str)

	w.ilog.Debug("printui128:%v", str)
	return 0
}

func printsf(vm *VirtualMachine) int64 {

	w := vm.WasmGo

	val := math.Float32frombits(uint32(vm.GetCurrentFrame().Locals[0]))
	str := strconv.FormatFloat(float64(val), 'e', 6, 32)
	// val := math.Float64frombits(vm.popUint64())
	// str := strconv.FormatFloat(val, 'e', 6, 32)

	w.context.ContextAppend(str)
	w.ilog.Debug("printsf:%v", str)

	return 0

}

func printdf(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	val := math.Float64frombits(uint64(vm.GetCurrentFrame().Locals[0]))
	str := strconv.FormatFloat(val, 'e', 15, 64)

	w.context.ContextAppend(str)
	w.ilog.Debug("printdf:%v", str)

	return 0
}

func printqf(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	val := int(vm.GetCurrentFrame().Locals[0])

	bytes := getMemory(vm, val, 16)
	var v eos_math.Float128
	rlp.DecodeBytes(bytes, &v)
	str := v.String()
	w.context.ContextAppend(str)

	w.ilog.Debug("printqf:%v", str)
	return 0
}

func printn(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	val := uint64(vm.GetCurrentFrame().Locals[0])
	str := common.S(val)
	w.context.ContextAppend(str)
	w.ilog.Debug("printn:%v", str)

	return 0
}

func printhex(vm *VirtualMachine) int64 {

	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	data := int(frame.Locals[0])
	dataLen := int(frame.Locals[1])

	str := hex.EncodeToString(getMemory(vm, data, dataLen))
	w.context.ContextAppend(str)

	w.ilog.Debug("printhex:%v", str)
	return 0

}

func isFeatureActive(vm *VirtualMachine) int64 {

	w := vm.WasmGo
	featureName := vm.GetCurrentFrame().Locals[0]
	//vm.pushUint64(uint64(b2i(false)))

	w.ilog.Debug("featureName:%v", common.S(uint64(featureName)))

	return 0

}

func activateFeature(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	featureName := vm.GetCurrentFrame().Locals[0]

	EosAssert(false, &UnsupportedFeature{}, "Unsupported Hardfork Detected")
	w.ilog.Debug("featureName:%v", common.S(uint64(featureName)))

	return 0
}

func setResourceLimits(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	account := common.AccountName(frame.Locals[0])
	ramBytes := frame.Locals[1]
	netWeight := frame.Locals[2]
	cpuWeight := frame.Locals[3]

	EosAssert(ramBytes >= -1, &WasmExecutionError{}, "invalid value for ram resource limit expected [-1,INT64_MAX]")
	EosAssert(netWeight >= -1, &WasmExecutionError{}, "invalid value for net resource limit expected [-1,INT64_MAX]")
	EosAssert(cpuWeight >= -1, &WasmExecutionError{}, "invalid value for cpu resource limit expected [-1,INT64_MAX]")

	if w.context.SetAccountLimits(account, ramBytes, netWeight, cpuWeight) {
		w.context.ValidateRamUsageInsert(account)
	}
	w.ilog.Debug("account:%v ramBytes:%d netWeight:%d cpuWeight:%d", account, ramBytes, netWeight, cpuWeight)

	return 0
}

func getResourceLimits(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	frame := vm.GetCurrentFrame()
	account := common.AccountName(frame.Locals[0])
	ramBytes := int(frame.Locals[1])
	netWeight := int(frame.Locals[2])
	cpuWeight := int(frame.Locals[3])

	var r, n, c int64
	w.context.GetAccountLimits(account, &r, &n, &c)

	setUint64(vm, ramBytes, uint64(r))
	setUint64(vm, netWeight, uint64(n))
	setUint64(vm, cpuWeight, uint64(c))

	w.ilog.Debug("account:%v ramBytes:%d netWeight:%d cpuWeigth:%d", account, ramBytes, netWeight, cpuWeight)

	return 0

}

func getBlockchainParametersPacked(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	packedBlockchainParameters := int(frame.Locals[0])
	bufferSize := int(frame.Locals[1])

	configuration := w.context.GetBlockchainParameters()
	p, _ := rlp.EncodeToBytes(configuration)
	//p := w.context.GetBlockchainParametersPacked()
	size := len(p)
	w.ilog.Debug("BlockchainParameters:%v bufferSize:%d size:%d", configuration, bufferSize, size)

	if bufferSize == 0 {
		//vm.pushUint64(uint64(size))
		return int64(size)
	}

	if size <= bufferSize {
		setMemory(vm, packedBlockchainParameters, p, 0, size)
		//vm.pushUint64(uint64(size))
		//w.ilog.Debug("BlockchainParameters:%v", configuration)
		return int64(size)
	}
	//vm.pushUint64(0)
	return 0

}

func setBlockchainParametersPacked(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	packedBlockchainParameters := int(frame.Locals[0])
	dataLen := int(frame.Locals[1])

	// p := make([]byte, datalen)
	// getMemory(vm,packedBlockchainParameters, 0, p, datalen)
	p := getMemory(vm, packedBlockchainParameters, dataLen)

	cfg := types.ChainConfig{}
	rlp.DecodeBytes(p, &cfg)

	//w.context.SetBlockchainParametersPacked(p)
	w.context.SetBlockchainParameters(&cfg)

	w.ilog.Debug("BlockchainParameters:%v ", cfg)

	return 0

}

func isPrivileged(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	account := common.AccountName(vm.GetCurrentFrame().Locals[0])

	ret := w.context.IsPrivileged(account)
	//vm.pushUint64(uint64(b2i(ret)))

	w.ilog.Debug("account:%v privileged:%v", account, ret)

	if ret {
		return 1
	} else {
		return 0
	}
}

func setPrivileged(vm *VirtualMachine) int64 {
	//fmt.Println("set_privileged")

	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	account := common.AccountName(frame.Locals[0])
	isPriv := int(frame.Locals[1])

	w.context.SetPrivileged(account, i2b(isPriv))

	w.ilog.Debug("account:%v privileged:%v", account, i2b(isPriv))

	return 0
}

func setProposedProducers(vm *VirtualMachine) int64 {
	//fmt.Println("set_proposed_producers")

	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	packedProducerSchedule := int(frame.Locals[0])
	dataLen := int(frame.Locals[1])

	p := getBytes(vm, packedProducerSchedule, dataLen)
	ret := w.context.SetProposedProducers(p)
	//vm.pushUint64(uint64(ret))

	producers := []types.ProducerKey{}
	rlp.DecodeBytes(p, &producers)
	w.ilog.Debug("packedProducerSchedule:%v ", producers)

	return int64(ret)
}

func getActiveProducers(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	frame := vm.GetCurrentFrame()
	producers := int(frame.Locals[0])
	bufferSize := int(frame.Locals[1])

	p := w.context.GetActiveProducersInBytes()
	s := len(p)

	if bufferSize == 0 {
		//vm.pushUint64(uint64(s))
		w.ilog.Debug("size:%d", s)
		return int64(s)
	}

	copySize := min(bufferSize, s)
	setMemory(vm, producers, p, 0, copySize)

	//vm.pushUint64(uint64(copySize))

	accounts := []common.AccountName{}
	rlp.DecodeBytes(p, &accounts)
	w.ilog.Debug("producers:%v", accounts)

	return int64(copySize)

}

func checkTime(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	//w.context.CheckTime()

	w.ilog.Debug("time:%v", common.Now())

	return 0
}

func currentTime(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")
	ret := w.context.CurrentTime()
	w.ilog.Debug("time:%v", ret)
	return ret.TimeSinceEpoch().Count()
}

func publicationTime(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	ret := w.context.PublicationTime()
	//vm.pushUint64(uint64(ret.TimeSinceEpoch().Count()))

	w.ilog.Debug("time:%v", ret)
	return ret.TimeSinceEpoch().Count()

}

func abort(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	EosAssert(false, &AbortCalled{}, AbortCalled{}.What())
	w.ilog.Debug("abort")

	return 0
}

func eosioAssert(vm *VirtualMachine) int64 {
	//w := vm.WasmGo
	frame := vm.GetCurrentFrame()

	condition := int(frame.Locals[0])
	ptr := int(frame.Locals[1])

	// message := string(getMemory(vm, ptr, getStringLength(vm, ptr)))
	// w.ilog.Debug("message:%v", string(message))

	if condition != 1 {
		var size int
		for i := 0; i < 512; i++ {
			if vm.Memory[ptr+i] == 0 {
				break
			}
			size++
		}
		message := vm.Memory[ptr : ptr+size]

		EosAssert(false, &EosioAssertMessageException{}, "assertion failure with message: %s", message)
		//Throw(&EosioAssertMessageException{})
		//Throw(&EosioAssertMessageException{}, "assertion failure with message: %v", message)
	}

	return 0
}

func eosioAssertMessage(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	condition := int(frame.Locals[0])
	msg := int(frame.Locals[1])
	msgLen := int(frame.Locals[2])

	message := string(getMemory(vm, msg, msgLen))
	w.ilog.Debug("message:%v", string(message))

	if condition != 1 {
		EosAssert(false, &EosioAssertMessageException{}, "assertion failure with message: %s", message)
		//Throw(&EosioAssertMessageException{}, "assertion failure with message: %v", message)
		//Throw(&EosioAssertMessageException{})
	}

	return 0

}

func eosioAssertCode(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	condition := int(frame.Locals[0])
	errorCode := frame.Locals[1]

	w.ilog.Debug("error code:%d", errorCode)
	if condition != 1 {
		EosAssert(false, &EosioAssertMessageException{}, "assertion failure with error code: %d", errorCode)
		//Throw(&EosioAssertMessageException{},"assertion failure with error code: %d", errorCode)
		//Throw(&EosioAssertMessageException{})
	}

	return 0

}

// void eosio_exit(int32_t code) {
//    throw wasm_exit{code};
// }
func eosioExit(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	errorCode := int(vm.GetCurrentFrame().Locals[0])

	w.ilog.Debug("error code:%d", errorCode)

	//Throw(wasmExit(code))
	return 0

}

func sendInline(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	frame := vm.GetCurrentFrame()
	data := int(frame.Locals[0])
	dataLen := int(frame.Locals[1])

	EosAssert(!w.context.InlineActionTooBig(dataLen), &InlineActionTooBig{}, "inline action too big")

	action := vm.Memory[data : data+dataLen] //action := getBytes(vm, data, dataLen)
	act := types.Action{}
	rlp.DecodeBytes(action, &act)
	w.context.ExecuteInline(&act)

	w.ilog.Debug("action:%v", act)
	return 0

}

func sendContextFreeInline(vm *VirtualMachine) int64 {

	w := vm.WasmGo
	frame := vm.GetCurrentFrame()
	data := int(frame.Locals[0])
	dataLen := int(frame.Locals[1])

	EosAssert(!w.context.InlineActionTooBig(dataLen), &InlineActionTooBig{}, "inline action too big")

	action := getBytes(vm, data, dataLen)
	act := types.Action{}
	rlp.DecodeBytes(action, &act)
	w.context.ExecuteContextFreeInline(&act)

	w.ilog.Debug("action:%v", act)
	return 0

}

// void send_deferred( const uint128_t& sender_id, account_name payer, array_ptr<char> data, size_t data_len, uint32_t replace_existing) {
//    try {
//       transaction trx;
//       fc::raw::unpack<transaction>(data, data_len, trx);
//       context.schedule_deferred_transaction(sender_id, payer, std::move(trx), replace_existing);
//    } FC_RETHROW_EXCEPTIONS(warn, "data as hex: ${data}", ("data", fc::to_hex(data, data_len)))
// }
func sendDeferred(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	EosAssert(!w.context.ContextFreeAction(), &UnaccessibleApi{}, "only context free api's can be used in this context")

	frame := vm.GetCurrentFrame()
	senderId := int(frame.Locals[0])
	payer := common.AccountName(frame.Locals[1])
	data := int(frame.Locals[2])
	dataLen := int(frame.Locals[3])
	replaceExisting := int(frame.Locals[4])

	bytes := getMemory(vm, senderId, 16)
	id := &eos_math.Uint128{}
	rlp.DecodeBytes(bytes, id)

	trx := getBytes(vm, data, dataLen)
	transaction := types.Transaction{}
	rlp.DecodeBytes(trx, &transaction)
	w.context.ScheduleDeferredTransaction(id, payer, &transaction, i2b(replaceExisting))

	w.ilog.Debug("id:%v transaction:%v", id, transaction)
	return 0
}

func cancelDeferred(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	//senderId := int(vm.popUint64())
	senderId := int(vm.GetCurrentFrame().Locals[0])

	bytes := getMemory(vm, senderId, 16)
	id := &eos_math.Uint128{}
	rlp.DecodeBytes(bytes, id)

	ret := w.context.CancelDeferredTransaction(id)
	//vm.pushUint64(uint64(b2i(ret)))

	w.ilog.Debug("id:%v", id)

	if ret {
		return 1
	} else {
		return 0
	}
}

// int read_transaction( array_ptr<char> data, size_t buffer_size ) {
//    bytes trx = context.get_packed_transaction();

//    auto s = trx.size();
//    if( buffer_size == 0) return s;

//    auto copy_size = std::min( buffer_size, s );
//    memcpy( data, trx.data(), copy_size );

//    return copy_size;
// }
func readTransaction(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	buffer := int(frame.Locals[0])
	bufferSize := int(frame.Locals[1])

	transaction := w.context.GetPackedTransaction()
	trx, _ := rlp.EncodeToBytes(transaction)

	s := len(trx)
	if bufferSize == 0 {
		w.ilog.Debug("transaction size:%d", s)
		//vm.pushUint64(uint64(s))
		return int64(s)
	}

	copySize := min(bufferSize, s)
	setMemory(vm, buffer, trx, 0, copySize)
	//vm.pushUint64(uint64(copySize))

	w.ilog.Debug("transaction:%v", transaction)
	return int64(copySize)

}

func transactionSize(vm *VirtualMachine) int64 {
	//fmt.Println("transaction_size")
	w := vm.WasmGo

	transaction := w.context.GetPackedTransaction()
	trx, _ := rlp.EncodeToBytes(transaction)
	s := len(trx)
	//vm.pushUint64(uint64(s))

	w.ilog.Debug("transaction size:%d", s)

	return int64(s)
}

func expiration(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	expiration := w.context.Expiration()
	//vm.pushUint64(uint64(expiration))

	w.ilog.Debug("expiration:%v", expiration)

	return int64(expiration)
}

func taposBlockNum(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	taposBlockNum := w.context.TaposBlockNum()
	//vm.pushUint64(uint64(taposBlockNum))

	w.ilog.Debug("taposBlockNum:%v", taposBlockNum)
	return int64(taposBlockNum)
}

func taposBlockPrefix(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	taposBlockPrefix := w.context.TaposBlockPrefix()
	//vm.pushUint64(uint64(taposBlockPrefix))

	w.ilog.Debug("taposBlockPrefix:%v", taposBlockPrefix)

	return int64(taposBlockPrefix)
}

func getAction(vm *VirtualMachine) int64 {
	w := vm.WasmGo

	frame := vm.GetCurrentFrame()
	typ := int(frame.Locals[0])
	index := int(frame.Locals[1])
	buffer := int(frame.Locals[2])
	bufferSize := int(frame.Locals[3])

	action := w.context.GetAction(uint32(typ), index)
	s, _ := rlp.EncodeSize(action)
	if bufferSize == 0 || bufferSize < s {
		//vm.pushUint64(uint64(s))
		w.ilog.Debug("action size:%d", s)
		return int64(s)
	}

	bytes, _ := rlp.EncodeToBytes(action)
	setMemory(vm, buffer, bytes, 0, s)
	//vm.pushUint64(uint64(s))
	w.ilog.Debug("action :%v size:%d", *action, s)
	return int64(s)

}

func getContextFreeData(vm *VirtualMachine) int64 {
	w := vm.WasmGo
	frame := vm.GetCurrentFrame()
	index := int(frame.Locals[0])
	buffer := int(frame.Locals[1])
	bufferSize := int(frame.Locals[2])

	EosAssert(w.context.ContextFreeAction(), &UnaccessibleApi{}, "this API may only be called from context_free apply")

	s, data := w.context.GetContextFreeData(index, bufferSize)
	if bufferSize == 0 || s == -1 {
		//vm.pushUint64(uint64(s))
		w.ilog.Debug("context free data size:%d", s)
		return int64(s)
	}
	setMemory(vm, buffer, data, 0, s)
	//vm.pushUint64(uint64(s))
	w.ilog.Debug("context free data :%v size:%d", data, s)
	return int64(s)

}
