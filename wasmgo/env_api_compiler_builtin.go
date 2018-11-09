package wasmgo

import (
	arithmetic "github.com/eosspark/eos-go/common/arithmetic_types"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"math"
	"unsafe"
	"github.com/eosspark/eos-go/exception/try"
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

//

// arithmetic long double
// void __addtf3( float128_t& ret, uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb ) {
//    float128_t a = {{ la, ha }};
//    float128_t b = {{ lb, hb }};
//    ret = f128_add( a, b );
// }

func addtf3(w *WasmGo, ret int, la, ha, lb, hb int64) {
	a := arithmetic.Float128{High: uint64(ha), Low: uint64(la)}
	b := arithmetic.Float128{High: uint64(hb), Low: uint64(lb)}

	re, _ := rlp.EncodeToBytes(a.Add(b))
	setMemory(w, ret, re, 0, len(re))
}

// void __subtf3( float128_t& ret, uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb ) {
//    float128_t a = {{ la, ha }};
//    float128_t b = {{ lb, hb }};
//    ret = f128_sub( a, b );
// }
func subtf3(w *WasmGo, ret int, la, ha, lb, hb int64) {
	a := arithmetic.Float128{High: uint64(ha), Low: uint64(la)}
	b := arithmetic.Float128{High: uint64(hb), Low: uint64(lb)}

	re := a.Sub(b).Bytes()
	setMemory(w, ret, re, 0, len(re))

}

// void __multf3( float128_t& ret, uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb ) {
//    float128_t a = {{ la, ha }};
//    float128_t b = {{ lb, hb }};
//    ret = f128_mul( a, b );
// }
func multf3(w *WasmGo, ret int, la, ha, lb, hb int64) {
	a := arithmetic.Float128{High: uint64(ha), Low: uint64(la)}
	b := arithmetic.Float128{High: uint64(hb), Low: uint64(lb)}

	re := a.Mul(b).Bytes()
	setMemory(w, ret, re, 0, len(re))
}

// void __divtf3( float128_t& ret, uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb ) {
//    float128_t a = {{ la, ha }};
//    float128_t b = {{ lb, hb }};
//    ret = f128_div( a, b );
// }
func divtf3(w *WasmGo, ret int, la, ha, lb, hb int64) {
	a := arithmetic.Float128{High: uint64(ha), Low: uint64(la)}
	b := arithmetic.Float128{High: uint64(hb), Low: uint64(lb)}

	re := a.Div(b).Bytes()
	setMemory(w, ret, re, 0, len(re))
}

// void __negtf2( float128_t& ret, uint64_t la, uint64_t ha ) {
//    ret = {{ la, (ha ^ (uint64_t)1 << 63) }};
// }
func negtf2(w *WasmGo, ret int, la, ha int64) {
	high := uint64(ha)
	high ^= uint64(1) << 63
	re := arithmetic.Float128{High: high, Low: uint64(la)}

	setMemory(w, ret, re.Bytes(), 0, len(re.Bytes()))
}

// // conversion long double
// void __extendsftf2( float128_t& ret, float f ) {
//    ret = f32_to_f128( softfloat_api::to_softfloat32(f) );
// }
func extendsftf2(w *WasmGo, ret int, f float32) { //TODO f float??
	f32 := arithmetic.Float32(math.Float32bits(f))
	re := arithmetic.F32ToF128(f32)
	setMemory(w, ret, re.Bytes(), 0, len(re.Bytes()))

}

//void __extenddftf2( float128_t& ret, double d ) {
//ret = f64_to_f128( softfloat_api::to_softfloat64(d) );
//}
func extenddftf2(w *WasmGo, ret int, d float64) { //TODO d double??
	f64 := arithmetic.Float64(math.Float64bits(d))
	re := arithmetic.F64ToF128(f64)
	setMemory(w, ret, re.Bytes(), 0, len(re.Bytes()))
}

//double __trunctfdf2( uint64_t l, uint64_t h ) {
//float128_t f = {{ l, h }};
//return softfloat_api::from_softfloat64(f128_to_f64( f ));
//}
func trunctfdf2(w *WasmGo, l, h int64) float64 { //TODO double??
	f128 := arithmetic.Float128{High: uint64(h), Low: uint64(l)}
	f64 := arithmetic.F128ToF64(f128)
	re := math.Float64frombits(uint64(f64))
	return re

}

//float __trunctfsf2( uint64_t l, uint64_t h ) {
//float128_t f = {{ l, h }};
//return softfloat_api::from_softfloat32(f128_to_f32( f ));
//}

func trunctfsf2(w *WasmGo, l, h int64) float32 { //TODO float??
	f128 := arithmetic.Float128{High: uint64(h), Low: uint64(l)}
	f32 := arithmetic.F128ToF32(f128)
	re := math.Float32frombits(uint32(f32))
	return re
}

//int32_t __fixtfsi( uint64_t l, uint64_t h ) {
//float128_t f = {{ l, h }};
//return f128_to_i32( f, 0, false );
//}
func fixtfsi(w *WasmGo, l, h int64) int {
	f128 := arithmetic.Float128{High: uint64(h), Low: uint64(l)}
	return int(arithmetic.F128ToI32(f128, 0, false))
}

//int64_t __fixtfdi( uint64_t l, uint64_t h ) {
//float128_t f = {{ l, h }};
//return f128_to_i64( f, 0, false );
//}
func fixtfdi(w *WasmGo, l, h int64) int64 {
	f128 := arithmetic.Float128{High: uint64(h), Low: uint64(l)}
	return arithmetic.F128ToI64(f128, 0, false)
}

//void __fixtfti( __int128& ret, uint64_t l, uint64_t h ) {
//float128_t f = {{ l, h }};
//ret = ___fixtfti( f );
//}

func fixtfti(w *WasmGo, ret int, l, h int64) {
	//f128 := arithmetic.Float128{High:uint64(h),Low:uint64(l)}
	//re :=

}

//uint32_t __fixunstfsi( uint64_t l, uint64_t h ) {
//float128_t f = {{ l, h }};
//return f128_to_ui32( f, 0, false );
//}

func fixunstfsi(w *WasmGo, l, h int64) int {
	f128 := arithmetic.Float128{High: uint64(h), Low: uint64(l)}
	return int(arithmetic.F128ToUi32(f128, 0, false))
}

//uint64_t __fixunstfdi( uint64_t l, uint64_t h ) {
//float128_t f = {{ l, h }};
//return f128_to_ui64( f, 0, false );
//}

func fixunstfdi(w *WasmGo, l, h int64) int64 {
	f128 := arithmetic.Float128{High: uint64(h), Low: uint64(l)}
	return int64(arithmetic.F128ToUi64(f128, 0, false))
}

//void __fixunstfti( unsigned __int128& ret, uint64_t l, uint64_t h ) {
//float128_t f = {{ l, h }};
//ret = ___fixunstfti( f );
//}
func fixunstfti(w *WasmGo, ret int, l, h int64) {

}

//void __fixsfti( __int128& ret, float a ) {
//ret = ___fixsfti( softfloat_api::to_softfloat32(a).v );
//}
func fixsfti(w *WasmGo, ret int, a float32) { //TODO float??

}

//void __fixdfti( __int128& ret, double a ) {
//ret = ___fixdfti( softfloat_api::to_softfloat64(a).v );
//}

func fixdfti(w *WasmGo, ret int, a float64) { //TODO double??

}

//void __fixunssfti( unsigned __int128& ret, float a ) {
//ret = ___fixunssfti( softfloat_api::to_softfloat32(a).v );
//}

func fixunssfti(w *WasmGo, ret int, a float32) { //TODO float??

}

//void __fixunsdfti( unsigned __int128& ret, double a ) {
//ret = ___fixunsdfti( softfloat_api::to_softfloat64(a).v );
//}

func fixunsdfti(w *WasmGo, ret int, a float64) { //TODO double??

}

//double __floatsidf( int32_t i ) {
//return softfloat_api::from_softfloat64(i32_to_f64(i));
//}
func floatsidf(w *WasmGo, i int) float64 { //TODO double??
	return 0
}

//void __floatsitf( float128_t& ret, int32_t i ) {
//ret = i32_to_f128(i);
//}

func floatsitf(w *WasmGo, ret int, i int) {
	re := arithmetic.I32ToF128(int32(i))
	setMemory(w, ret, re.Bytes(), 0, len(re.Bytes()))
}

//void __floatditf( float128_t& ret, uint64_t a ) {
//ret = i64_to_f128( a );
//}
func floatditf(w *WasmGo, ret int, a int64) {
	re := arithmetic.I64ToF128(a)
	setMemory(w, ret, re.Bytes(), 0, len(re.Bytes()))
}

//void __floatunsitf( float128_t& ret, uint32_t i ) {
//ret = ui32_to_f128(i);
//}

func floatunsitf(w *WasmGo, ret int, i int) {
	re := arithmetic.Ui32ToF128(uint32(i))
	setMemory(w, ret, re.Bytes(), 0, len(re.Bytes()))
}

//void __floatunditf( float128_t& ret, uint64_t a ) {
//ret = ui64_to_f128( a );
//}
func floatunditf(w *WasmGo, ret int, a int64) {
	re := arithmetic.Ui64ToF128(uint64(a))
	setMemory(w, ret, re.Bytes(), 0, len(re.Bytes()))
}

//double __floattidf( uint64_t l, uint64_t h ) {
//fc::uint128_t v(h, l);
//unsigned __int128 val = (unsigned __int128)v;
//return ___floattidf( *(__int128*)&val );
//}

func floattidf(w *WasmGo, l, h int64) float64 { //TODO double
	//v := arithmetic.Uint128{uint64(h), uint64(l)}
	return 0
}

//double __floatuntidf( uint64_t l, uint64_t h ) {
//fc::uint128_t v(h, l);
//return ___floatuntidf( (unsigned __int128)v );
//}

func floatuntidf(w *WasmGo, l, h int64) float64 { //TODO double
	//v := arithmetic.Uint128{h, l}
	//return floatuntidf(w,v)
	return 0
}

//int ___cmptf2( uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb, int return_value_if_nan ) {
//float128_t a = {{ la, ha }};
//float128_t b = {{ lb, hb }};
//if ( __unordtf2(la, ha, lb, hb) )
//return return_value_if_nan;
//if ( f128_lt( a, b ) )
//return -1;
//if ( f128_eq( a, b ) )
//return 0;
//return 1;
//}
func _cmptf2(w *WasmGo, la, ha, lb, hb int64, return_value_if_nan int) int { //TODO unsame with regist
	a := arithmetic.Float128{High: uint64(ha), Low: uint64(la)}
	b := arithmetic.Float128{High: uint64(hb), Low: uint64(lb)}
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

//int __eqtf2( uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb ) {
//return ___cmptf2(la, ha, lb, hb, 1);
//}
func eqtf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, 1)
}

//int __netf2( uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb ) {
//return ___cmptf2(la, ha, lb, hb, 1);
//}

func netf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, 1)
}

//int __getf2( uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb ) {
//return ___cmptf2(la, ha, lb, hb, -1);
//}

func getf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, -1)
}

//int __gttf2( uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb ) {
//return ___cmptf2(la, ha, lb, hb, 0);
//}

func gttf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, 0)
}

//int __letf2( uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb ) {
//return ___cmptf2(la, ha, lb, hb, 1);
//}
func letf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, 1)
}

//int __lttf2( uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb ) {
//return ___cmptf2(la, ha, lb, hb, 0);
//}
func lttf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, 0)
}

//int __cmptf2( uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb ) {
//return ___cmptf2(la, ha, lb, hb, 1);
//}
func cmptf2(w *WasmGo, la, ha, lb, hb int64) int {
	return _cmptf2(w, la, ha, lb, hb, 1)
}

//int __unordtf2( uint64_t la, uint64_t ha, uint64_t lb, uint64_t hb ) {
//float128_t a = {{ la, ha }};
//float128_t b = {{ lb, hb }};
//if ( softfloat_api::is_nan(a) || softfloat_api::is_nan(b) )
//return 1;
//return 0;
//}
//

func unordtf2(w *WasmGo, la, ha, lb, hb int64) int {
	a := arithmetic.Float128{uint64(ha), uint64(la)}
	b := arithmetic.Float128{uint64(hb), uint64(lb)}
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
