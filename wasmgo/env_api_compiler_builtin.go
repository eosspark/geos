package wasmgo

import (
	arithmetic "github.com/eosspark/eos-go/common/arithmetic_types"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"math"
	"unsafe"
)

var count = 0

const SHIFT_WIDTH = uint32(unsafe.Sizeof(uint64(0)*8) - 1) //63

func ashlti3(w *WasmGo, ret int, low, high int64, shift int) {
	i := arithmetic.Int128{Low: uint64(low), High: uint64(high)}
	i.LeftShifts(shift)

	re, _ := rlp.EncodeToBytes(i)
	setMemory(w, ret, re, 0, len(re))
}

func ashrti3(w *WasmGo, ret int, low, high int64, shift int) {
	i := arithmetic.Int128{Low: uint64(low), High: uint64(high)}
	i.RightShifts(shift)

	re, _ := rlp.EncodeToBytes(i)
	setMemory(w, ret, re, 0, len(re))
}

func lshlti3(w *WasmGo, ret int, low, high int64, shift int) {
	i := arithmetic.Int128{Low: uint64(low), High: uint64(high)}
	i.LeftShifts(shift)

	re, _ := rlp.EncodeToBytes(i)
	setMemory(w, ret, re, 0, len(re))
}

func lshrti3(w *WasmGo, ret int, low, high int64, shift int) {
	i := arithmetic.Uint128{Low: uint64(low), High: uint64(high)}
	i.RightShifts(shift)

	re, _ := rlp.EncodeToBytes(i)
	setMemory(w, ret, re, 0, len(re))
}

func divti3(w *WasmGo, ret int, la, ha, lb, hb int64) {
	lhs := arithmetic.Int128{Low: uint64(la), High: uint64(ha)}
	rhs := arithmetic.Int128{Low: uint64(lb), High: uint64(hb)}

	try.EosAssert(!rhs.IsZero(), &exception.ArithmeticException{}, "divide by zero")

	quotient, _ := lhs.Div(rhs)
	re, _ := rlp.EncodeToBytes(quotient)
	setMemory(w, ret, re, 0, len(re))
}

func udivti3(w *WasmGo, ret int, la, ha, lb, hb int64) {
	lhs := arithmetic.Uint128{Low: uint64(la), High: uint64(ha)}
	rhs := arithmetic.Uint128{Low: uint64(lb), High: uint64(hb)}

	try.EosAssert(!rhs.IsZero(), &exception.ArithmeticException{}, "divide by zero")
	quotient, _ := lhs.Div(rhs)

	re, _ := rlp.EncodeToBytes(quotient)
	setMemory(w, ret, re, 0, len(re))
}

func multi3(w *WasmGo, ret int, la, ha, lb, hb int64) {
	lhs := arithmetic.Int128{Low: uint64(la), High: uint64(ha)}
	rhs := arithmetic.Int128{Low: uint64(lb), High: uint64(hb)}

	re, _ := rlp.EncodeToBytes(lhs.Mul(rhs))
	setMemory(w, ret, re, 0, len(re))

}

func modti3(w *WasmGo, ret int, la, ha, lb, hb int64) {
	lhs := arithmetic.Int128{High: uint64(ha), Low: uint64(la)}
	rhs := arithmetic.Int128{High: uint64(hb), Low: uint64(lb)}
	try.EosAssert(!rhs.IsZero(), &exception.ArithmeticException{}, "divide by zero")

	_, remainder := lhs.Div(rhs)
	re, _ := rlp.EncodeToBytes(remainder)
	setMemory(w, ret, re, 0, len(re))
}

func umodti3(w *WasmGo, ret int, la, ha, lb, hb int64) {
	lhs := arithmetic.Uint128{Low: uint64(la), High: uint64(ha)}
	rhs := arithmetic.Uint128{Low: uint64(lb), High: uint64(hb)}

	try.EosAssert(!rhs.IsZero(), &exception.ArithmeticException{}, "divide by zero")
	_, remainder := lhs.Div(rhs)
	re, _ := rlp.EncodeToBytes(remainder)
	setMemory(w, ret, re, 0, len(re))
}

func addtf3(w *WasmGo, ret int, la, ha, lb, hb int64) {
	a := arithmetic.Float128{Low: uint64(la), High: uint64(ha)}
	b := arithmetic.Float128{Low: uint64(lb), High: uint64(hb)}

	re, _ := rlp.EncodeToBytes(a.Add(b))
	setMemory(w, ret, re, 0, len(re))
}

func subtf3(w *WasmGo, ret int, la, ha, lb, hb int64) {
	a := arithmetic.Float128{Low: uint64(la), High: uint64(ha)}
	b := arithmetic.Float128{Low: uint64(lb), High: uint64(hb)}

	re, _ := rlp.EncodeToBytes(a.Sub(b))
	setMemory(w, ret, re, 0, len(re))
}

func multf3(w *WasmGo, ret int, la, ha, lb, hb int64) {
	a := arithmetic.Float128{Low: uint64(la), High: uint64(ha)}
	b := arithmetic.Float128{Low: uint64(lb), High: uint64(hb)}

	re, _ := rlp.EncodeToBytes(a.Mul(b))
	setMemory(w, ret, re, 0, len(re))
}

func divtf3(w *WasmGo, ret int, la, ha, lb, hb int64) {
	a := arithmetic.Float128{Low: uint64(la), High: uint64(ha)}
	b := arithmetic.Float128{Low: uint64(lb), High: uint64(hb)}

	re, _ := rlp.EncodeToBytes(a.Div(b))
	setMemory(w, ret, re, 0, len(re))
}

func negtf2(w *WasmGo, ret int, la, ha int64) {
	high := uint64(ha)
	high ^= uint64(1) << 63
	f128 := arithmetic.Float128{Low: uint64(la), High: high}

	re, _ := rlp.EncodeToBytes(f128)
	setMemory(w, ret, re, 0, len(re))
}

func extendsftf2(w *WasmGo, ret int, f float32) { //TODO f float??
	f32 := arithmetic.Float32(math.Float32bits(f))
	f128 := arithmetic.F32ToF128(f32)

	re, _ := rlp.EncodeToBytes(f128)
	setMemory(w, ret, re, 0, len(re))
}

func extenddftf2(w *WasmGo, ret int, d float64) { //TODO d double??
	f64 := arithmetic.Float64(math.Float64bits(d))
	f128 := arithmetic.F64ToF128(f64)

	re, _ := rlp.EncodeToBytes(f128)
	setMemory(w, ret, re, 0, len(re))
}

func trunctfdf2(w *WasmGo, l, h int64) float64 { //TODO double??
	f128 := arithmetic.Float128{Low: uint64(l), High: uint64(h)}
	f64 := arithmetic.F128ToF64(f128)
	return math.Float64frombits(uint64(f64))
}

func trunctfsf2(w *WasmGo, l, h int64) float32 { //TODO float??
	f128 := arithmetic.Float128{Low: uint64(l), High: uint64(h)}
	f32 := arithmetic.F128ToF32(f128)
	return math.Float32frombits(uint32(f32))
}

func fixtfsi(w *WasmGo, l, h int64) int {
	f128 := arithmetic.Float128{Low: uint64(l), High: uint64(h)}
	return int(arithmetic.F128ToI32(f128, 0, false))
}

func fixtfdi(w *WasmGo, l, h int64) int64 {
	f128 := arithmetic.Float128{Low: uint64(l), High: uint64(h)}
	return arithmetic.F128ToI64(f128, 0, false)
}

func fixunstfsi(w *WasmGo, l, h int64) int {
	f128 := arithmetic.Float128{Low: uint64(l), High: uint64(h)}
	return int(arithmetic.F128ToUi32(f128, 0, false))
}

func fixunstfdi(w *WasmGo, l, h int64) int64 {
	f128 := arithmetic.Float128{Low: uint64(l), High: uint64(h)}
	return int64(arithmetic.F128ToUi64(f128, 0, false))
}

func fixtfti(w *WasmGo, ret int, l, h int64) {
	f128 := arithmetic.Float128{Low: uint64(l), High: uint64(h)}
	int128 := arithmetic.Fixtfti(f128)

	re, _ := rlp.EncodeToBytes(int128)
	setMemory(w, ret, re, 0, len(re))
}

func fixunstfti(w *WasmGo, ret int, l, h int64) {
	f := arithmetic.Float128{Low: uint64(l), High: uint64(h)}
	uint128 := arithmetic.Fixunstfti(f)

	re, _ := rlp.EncodeToBytes(uint128)
	setMemory(w, ret, re, 0, len(re))
}

func fixsfti(w *WasmGo, ret int, a float32) { //TODO float??
	int128 := arithmetic.Fixsfti(math.Float32bits(a))
	re, _ := rlp.EncodeToBytes(int128)
	setMemory(w, ret, re, 0, len(re))
}

func fixdfti(w *WasmGo, ret int, a float64) { //TODO double??
	int128 := arithmetic.Fixdfti(math.Float64bits(a))
	re, _ := rlp.EncodeToBytes(int128)
	setMemory(w, ret, re, 0, len(re))
}

func fixunssfti(w *WasmGo, ret int, a float32) { //TODO float??
	uint128 := arithmetic.Fixunssfti(math.Float32bits(a))
	re, _ := rlp.EncodeToBytes(uint128)
	setMemory(w, ret, re, 0, len(re))
}

func fixunsdfti(w *WasmGo, ret int, a float64) { //TODO double??
	uint128 := arithmetic.Fixunsdfti(math.Float64bits(a))
	re, _ := rlp.EncodeToBytes(uint128)
	setMemory(w, ret, re, 0, len(re))
}

func floatsidf(w *WasmGo, i int) float64 { //TODO double??
	return math.Float64frombits(uint64(arithmetic.I32ToF64(int32(i))))
}

func floatsitf(w *WasmGo, ret int, i int) {
	re := arithmetic.I32ToF128(int32(i))
	setMemory(w, ret, re.Bytes(), 0, len(re.Bytes()))
}

func floatditf(w *WasmGo, ret int, a int64) {
	re := arithmetic.I64ToF128(a)
	setMemory(w, ret, re.Bytes(), 0, len(re.Bytes()))
}

func floatunsitf(w *WasmGo, ret int, i int) {
	re := arithmetic.Ui32ToF128(uint32(i))
	setMemory(w, ret, re.Bytes(), 0, len(re.Bytes()))
}

func floatunditf(w *WasmGo, ret int, a int64) {
	re := arithmetic.Ui64ToF128(uint64(a))
	setMemory(w, ret, re.Bytes(), 0, len(re.Bytes()))
}

func floattidf(w *WasmGo, l, h int64) float64 { //TODO double
	v := arithmetic.Int128{Low: uint64(l), High: uint64(h)}
	return arithmetic.Floattidf(v)
}

func floatuntidf(w *WasmGo, l, h int64) float64 { //TODO double
	v := arithmetic.Uint128{Low: uint64(l), High: uint64(h)}
	return arithmetic.Floatuntidf(v)
	return 0
}

func _cmptf2(w *WasmGo, la, ha, lb, hb int64, return_value_if_nan int) int { //TODO unsame with regist
	a := arithmetic.Float128{Low: uint64(la), High: uint64(ha)}
	b := arithmetic.Float128{Low: uint64(lb), High: uint64(hb)}
	if unordtf2(w, la, ha, lb, hb) != 0 {
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

func eqtf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, 1)
}

func netf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, 1)
}

func getf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, -1)
}

func gttf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, 0)
}

func letf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, 1)
}

func lttf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, 0)
}

func cmptf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, 1)
}

func unordtf2(w *WasmGo, la, ha, lb, hb int64) int {
	a := arithmetic.Float128{Low: uint64(la), High: uint64(ha)}
	b := arithmetic.Float128{Low: uint64(lb), High: uint64(hb)}
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
