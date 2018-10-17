package figure

///* 7.18.1.2 Minimum-width integer types */
//type int_least8_t int8
//type int_least16_t int16
//type int_least32_t int32
//type int_least64_t int64
//type uint_least8_t uint8
//type uint_least16_t uint16
//type uint_least32_t uint32
//type uint_least64_t uint64
//
///* 7.18.1.3 Fastest-width integer types */
//type intFast8 int8
//type intFast16 int16
//type intFast32 int32
//type intFast64 int64
//type uintFast8 uint8
//type uintFast16 uint16
//type uintFast32 uint32
//type uintFast64 uint64

//#define packToF128UI64( sign, exp, sig64 ) (((uint_fast64_t) (sign)<<63) + ((uint_fast64_t) (exp)<<48) + (sig64))
func packToF128UI64(sign, exp, sig64 uint64) uint64 {
	return sign<<63 + exp<<48 + sig64
}

//float128_t ui64_to_f128( uint64_t a )
//{
//uint_fast64_t uiZ64, uiZ0;
//int_fast8_t shiftDist;
//struct uint128 zSig;
//union ui128_f128 uZ;
//
//if ( ! a ) {
//uiZ64 = 0;
//uiZ0  = 0;
//} else {
//shiftDist = softfloat_countLeadingZeros64( a ) + 49;
//if ( 64 <= shiftDist ) {
//zSig.v64 = a<<(shiftDist - 64);
//zSig.v0  = 0;
//} else {
//zSig = softfloat_shortShiftLeft128( 0, a, shiftDist );
//}
//uiZ64 = packToF128UI64( 0, 0x406E - shiftDist, zSig.v64 );
//uiZ0  = zSig.v0;
//}
//uZ.ui.v64 = uiZ64;
//uZ.ui.v0  = uiZ0;
//return uZ.f;
//
//}

func Ui64ToF128(a uint64) Float128 {
	var uiZ64, uiZ0 uint64
	var shiftDist uint8
	var zSig Uint128
	var uZ Float128

	if a == 0 {
		uiZ64 = 0
		uiZ0 = 0
	} else {
		shiftDist = softfloatCountLeadingZeros64(a) + 49
		if 64 <= shiftDist {
			zSig.High = a<<shiftDist - 64
			zSig.Low = 0

		} else {
			zSig = softfloatShortShiftLeft128(0, a, shiftDist)
		}
		uiZ64 = packToF128UI64(0, uint64(0x406E)-uint64(shiftDist), zSig.High)
		uiZ0 = zSig.Low
	}

	uZ[1] = uiZ64
	uZ[0] = uiZ0

	return uZ
}

//float128_t ui32_to_f128( uint32_t a )
//{
//uint_fast64_t uiZ64;
//int_fast8_t shiftDist;
//union ui128_f128 uZ;
//
//uiZ64 = 0;
//if ( a ) {
//shiftDist = softfloat_countLeadingZeros32( a ) + 17;
//uiZ64 =
//packToF128UI64(
//0, 0x402E - shiftDist, (uint_fast64_t) a<<shiftDist );
//}
//uZ.ui.v64 = uiZ64;
//uZ.ui.v0  = 0;
//return uZ.f;
//
//}

func Ui32ToF128(a uint32) Float128 {
	var uiZ64 uint64
	var shiftDist uint8
	var uZ Float128

	uiZ64 = 0
	if a != 0 {
		shiftDist = softfloatCountLeadingZeros32(a) + 17
		uiZ64 = packToF128UI64(0, uint64(0x402E)-uint64(shiftDist), uint64(a<<shiftDist))
	}
	uZ[1] = uiZ64
	uZ[0] = 0
	return uZ
}

//float128_t i64_to_f128( int64_t a )
//{
//uint_fast64_t uiZ64, uiZ0;
//bool sign;
//uint_fast64_t absA;
//int_fast8_t shiftDist;
//struct uint128 zSig;
//union ui128_f128 uZ;
//
//if ( ! a ) {
//uiZ64 = 0;
//uiZ0  = 0;
//} else {
//sign = (a < 0);
//absA = sign ? -(uint_fast64_t) a : (uint_fast64_t) a;
//shiftDist = softfloat_countLeadingZeros64( absA ) + 49;
//if ( 64 <= shiftDist ) {
//zSig.v64 = absA<<(shiftDist - 64);
//zSig.v0  = 0;
//} else {
//zSig = softfloat_shortShiftLeft128( 0, absA, shiftDist );
//}
//uiZ64 = packToF128UI64( sign, 0x406E - shiftDist, zSig.v64 );
//uiZ0  = zSig.v0;
//}
//uZ.ui.v64 = uiZ64;
//uZ.ui.v0  = uiZ0;
//return uZ.f;
//
//}

func I64ToF128(a int64) Float128 {
	var uiZ64, uiZ0 uint64
	var sign uint64
	var absA uint64

	var shiftDist uint8
	var zSig Uint128
	var uZ Float128

	if a == 0 {
		uiZ64 = 0
		uiZ0 = 0
	} else {
		if a < 0 {
			sign = 1 //true
			absA = -uint64(a)
		} else {
			sign = 0 //false
			absA = uint64(a)
		}
		shiftDist = softfloatCountLeadingZeros64(absA) + 49
		if 64 <= shiftDist {
			zSig.High = absA << (shiftDist - 64)
			zSig.Low = 0
		} else {
			zSig = softfloatShortShiftLeft128(0, absA, shiftDist)
		}
		uiZ64 = packToF128UI64(sign, uint64(0x406E)-uint64(shiftDist), zSig.High)
		uiZ0 = zSig.Low
	}
	uZ[1] = uiZ64
	uZ[0] = uiZ0
	return uZ
}

//float128_t i32_to_f128( int32_t a )
//{
//uint_fast64_t uiZ64;
//bool sign;
//uint_fast32_t absA;
//int_fast8_t shiftDist;
//union ui128_f128 uZ;
//
//uiZ64 = 0;
//if ( a ) {
//sign = (a < 0);
//absA = sign ? -(uint_fast32_t) a : (uint_fast32_t) a;
//shiftDist = softfloat_countLeadingZeros32( absA ) + 17;
//uiZ64 =
//packToF128UI64(
//sign, 0x402E - shiftDist, (uint_fast64_t) absA<<shiftDist );
//}
//uZ.ui.v64 = uiZ64;
//uZ.ui.v0  = 0;
//return uZ.f;
//
//}

func I32ToF128(a int32) Float128 {
	var uiZ64 uint64
	var sign uint64
	var absA uint32
	var shiftDist uint8
	var uZ Float128

	uiZ64 = 0
	if a != 0 {
		if a < 0 {
			absA = -uint32(a)
		} else {
			absA = uint32(a)
		}
		shiftDist = softfloatCountLeadingZeros32(absA) + 17
		uiZ64 = packToF128UI64(sign, uint64(0x402E)-uint64(shiftDist), uint64(absA<<shiftDist))
	}
	uZ[1] = uiZ64
	uZ[0] = 0
	return uZ
}

//float128_t f32_to_f128( float32_t a )
//{
//union ui32_f32 uA;
//uint_fast32_t uiA;
//bool sign;
//int_fast16_t exp;
//uint_fast32_t frac;
//struct commonNaN commonNaN;
//struct uint128 uiZ;
//struct exp16_sig32 normExpSig;
//union ui128_f128 uZ;
//
///*------------------------------------------------------------------------
//*------------------------------------------------------------------------*/
//uA.f = a;
//uiA = uA.ui;
//sign = signF32UI( uiA );
//exp  = expF32UI( uiA );
//frac = fracF32UI( uiA );
///*------------------------------------------------------------------------
//*------------------------------------------------------------------------*/
//if ( exp == 0xFF ) {
//if ( frac ) {
//softfloat_f32UIToCommonNaN( uiA, &commonNaN );
//uiZ = softfloat_commonNaNToF128UI( &commonNaN );
//} else {
//uiZ.v64 = packToF128UI64( sign, 0x7FFF, 0 );
//uiZ.v0  = 0;
//}
//goto uiZ;
//}
///*------------------------------------------------------------------------
//*------------------------------------------------------------------------*/
//if ( ! exp ) {
//if ( ! frac ) {
//uiZ.v64 = packToF128UI64( sign, 0, 0 );
//uiZ.v0  = 0;
//goto uiZ;
//}
//normExpSig = softfloat_normSubnormalF32Sig( frac );
//exp = normExpSig.exp - 1;
//frac = normExpSig.sig;
//}
///*------------------------------------------------------------------------
//*------------------------------------------------------------------------*/
//uiZ.v64 = packToF128UI64( sign, exp + 0x3F80, (uint_fast64_t) frac<<25 );
//uiZ.v0  = 0;
//uiZ:
//uZ.ui = uiZ;
//return uZ.f;
//
//}

func F32ToF128(a Float32_t) Float128 {
	out := Float128{}
	return out
}
