package arithmeticTypes

import "unsafe"

//THREAD_LOCAL uint_fast8_t softfloat_roundingMode = softfloat_round_near_even;
//THREAD_LOCAL uint_fast8_t softfloat_detectTininess = init_detectTininess;
//THREAD_LOCAL uint_fast8_t softfloat_exceptionFlags = 0;
//
//THREAD_LOCAL uint_fast8_t extF80_roundingPrecision = 80;

const (
	softfloat_roundingMode = softfloat_round_near_even
)

const (
	softfloat_tininess_beforeRounding uint8 = 0
	softfloat_tininess_afterRounding  uint8 = 1

	softfloat_detectTininess = softfloat_tininess_afterRounding
)

/*----------------------------------------------------------------------------
| "Common NaN" structure, used to transfer NaN representations from one format
| to another.
*----------------------------------------------------------------------------*/
type commonNaN struct {
	sign    bool
	v0, v64 uint64
}

func F128ToF32(a Float128) Float32 {
	var uA Float128
	var uiA64, uiA0 uint64
	var sign bool
	var exp int32
	var frac64 uint64
	var commonNaN commonNaN
	var uiZ, frac32 uint32
	var uZ Float32

	uA = a
	uiA64 = uA.High
	uiA0 = uA.Low
	sign = signF128Ui64(uiA64)
	exp = expF128Ui64(uiA64)
	if uiA0 != 0 {
		frac64 = fracF128Ui64(uiA64) | 1
	} else {
		frac64 = fracF128Ui64(uiA64)
	}

	if exp == 0x7FFF {
		if frac64 != 0 {
			softfloat_f128UIToCommonNaN(uiA64, uiA0, &commonNaN)
			uiZ = softfloat_commonNaNToF32UI(&commonNaN)
		} else {
			uiZ = packToF32UI(sign, 0xFF, 0)
		}
		goto uiZ
	}

	frac32 = uint32(softfloat_shortShiftRightJam64(frac64, 18))
	if (exp | int32(frac32)) == 0 {
		uiZ = packToF32UI(sign, 0, 0)
		goto uiZ
	}

	exp -= 0x3F81

	if unsafe.Sizeof(int16(0)) < unsafe.Sizeof(int32(0)) {
		if exp < -0x1000 {
			exp = -0x1000
		}
	}
	return softfloat_roundPackToF32(sign, int16(exp), frac32|0x40000000)

uiZ:
	uZ = Float32(uiZ)
	return uZ
}

/*----------------------------------------------------------------------------
| Assuming the unsigned integer formed from concatenating `uiA64' and `uiA0'
| has the bit pattern of a 128-bit floating-point NaN, converts this NaN to
| the common NaN form, and stores the resulting common NaN at the location
| pointed to by `zPtr'.  If the NaN is a signaling NaN, the invalid exception
| is raised.
*----------------------------------------------------------------------------*/
//void
//softfloat_f128UIToCommonNaN(
//uint_fast64_t uiA64, uint_fast64_t uiA0, struct commonNaN *zPtr )
//{
//struct uint128 NaNSig;
//
//if ( softfloat_isSigNaNF128UI( uiA64, uiA0 ) ) {
//softfloat_raiseFlags( softfloat_flag_invalid );
//}
//NaNSig = softfloat_shortShiftLeft128( uiA64, uiA0, 16 );
//zPtr->sign = uiA64>>63;
//zPtr->v64  = NaNSig.v64;
//zPtr->v0   = NaNSig.v0;
//
//}

func softfloat_f128UIToCommonNaN(uiA64, uiA0 uint64, zPtr *commonNaN) {
	var NaNSig Uint128
	if softfloat_isSigNaNF128UI(uiA64, uiA0) {
		softfloat_raiseFlags(softfloat_flag_invalid)
	}
	NaNSig = softfloat_shortShiftLeft128(uiA64, uiA0, 16)
	if uiA64>>63 == 1 {
		zPtr.sign = true
	} else {
		zPtr.sign = false
	}

	zPtr.v64 = NaNSig.High
	zPtr.v0 = NaNSig.Low

}

/*----------------------------------------------------------------------------
| Returns true when the 128-bit unsigned integer formed from concatenating
| 64-bit 'uiA64' and 64-bit 'uiA0' has the bit pattern of a 128-bit floating-
| point signaling NaN.
| Note:  This macro evaluates its arguments more than once.
*----------------------------------------------------------------------------*/
//#define softfloat_isSigNaNF128UI( uiA64, uiA0 ) ((((uiA64) & UINT64_C( 0x7FFF800000000000 )) == UINT64_C( 0x7FFF000000000000 )) && ((uiA0) || ((uiA64) & UINT64_C( 0x00007FFFFFFFFFFF ))))

func softfloat_isSigNaNF128UI(uiA64, uiA0 uint64) bool {
	return (uiA64&uint64(0x7FFF800000000000)) == uint64(0x7FFF000000000000) && (uiA0 != 0 || (uiA64&uint64(0x00007FFFFFFFFFFF) != 0))
}

/*----------------------------------------------------------------------------
| Converts the common NaN pointed to by `aPtr' into a 32-bit floating-point
| NaN, and returns the bit pattern of this value as an unsigned integer.
*----------------------------------------------------------------------------*/
//uint_fast32_t softfloat_commonNaNToF32UI( const struct commonNaN *aPtr )
//{
//
//return (uint_fast32_t) aPtr->sign<<31 | 0x7FC00000 | aPtr->v64>>41;
//
//}

func softfloat_commonNaNToF32UI(aPtr *commonNaN) uint32 {
	if aPtr.sign {
		return uint32(1)<<31 | uint32(0x7FC00000) | uint32(aPtr.v64>>41)
	} else {
		return uint32(0x7FC00000) | uint32(aPtr.v64>>41)
	}
	//return (uint_fast32_t) aPtr->sign<<31 | 0x7FC00000 | aPtr->v64>>41

}

//#define packToF32UI( sign, exp, sig ) (((uint32_t) (sign)<<31) + ((uint32_t) (exp)<<23) + (sig))

func packToF32UI(sign bool, exp, sig uint32) uint32 {

	if sign {
		return uint32(1)<<31 + uint32(exp<<23) + sig
	} else {
		return uint32(exp<<23) + sig
	}

}

//uint64_t softfloat_shortShiftRightJam64( uint64_t a, uint_fast8_t dist )
//{ return a>>dist | ((a & (((uint_fast64_t) 1<<dist) - 1)) != 0); }

func softfloat_shortShiftRightJam64(a uint64, dist uint8) uint64 {
	if (a & ((uint64(1) << dist) - 1)) != 0 {
		return a>>dist | 1
	}
	return a >> dist

}

func softfloat_roundPackToF32(sign bool, exp int16, sig uint32) Float32 {
	var roundingMode uint8
	var roundNearEven bool
	var roundIncrement, roundBits uint8
	var isTiny bool
	var uiZ uint32

	roundingMode = softfloat_roundingMode
	roundNearEven = roundingMode == softfloat_round_near_even
	roundIncrement = 0x40
	if !(roundNearEven && (roundingMode != softfloat_round_near_maxMag)) {
		var mode uint8
		if sign {
			mode = softfloat_round_min
		} else {
			mode = softfloat_round_max
		}
		if roundingMode == mode {
			roundIncrement = 0x7F
		} else {
			roundIncrement = 0
		}
	}
	roundBits = uint8(sig & 0x7F)

	if 0xFD <= uint(exp) {
		if exp < 0 {
			/*----------------------------------------------------------------
			*----------------------------------------------------------------*/
			isTiny = (softfloat_detectTininess == softfloat_tininess_beforeRounding) || (exp < -1) || (sig+uint32(roundIncrement) < 0x80000000)
			sig = softfloat_shiftRightJam32(sig, uint16(-exp))
			exp = 0
			roundBits = uint8(sig & 0x7F)
			if isTiny && roundBits != 0 {
				softfloat_raiseFlags(softfloat_flag_underflow)
			}
		} else if (0xFD < exp) || (0x80000000 <= sig+uint32(roundIncrement)) {
			/*----------------------------------------------------------------
			*----------------------------------------------------------------*/
			softfloat_raiseFlags(softfloat_flag_overflow | softfloat_flag_inexact)

			if roundIncrement != 0 {
				uiZ = packToF32UI(sign, 0xFF, 0)
			} else {
				uiZ = packToF32UI(sign, 0xFF, 0) - 1
			}

			goto uiZ
		}
	}

	sig = (sig + uint32(roundIncrement)) >> 7
	//if  roundBits !=0 {
	//	//softfloat_exceptionFlags |= softfloat_flag_inexact
	//	#ifdef SOFTFLOAT_ROUND_ODD
	//}
	//sig &= ~(uint_fast32_t) (! (roundBits ^ 0x40) & roundNearEven);
	if roundBits^0x40 != 0 {
		sig &= ^uint32(0)
	} else {
		if roundNearEven {
			sig &= ^uint32(1)
		} else {
			sig &= ^uint32(0)
		}
	}

	if sig == 0 {
		exp = 0
	}

	uiZ = packToF32UI(sign, uint32(exp), sig)
uiZ:
	return Float32(uiZ)

}

//INLINE uint32_t softfloat_shiftRightJam32( uint32_t a, uint_fast16_t dist )
//{
//return
//(dist < 31) ? a>>dist | ((uint32_t) (a<<(-dist & 31)) != 0) : (a != 0);
//}

func softfloat_shiftRightJam32(a uint32, dist uint16) uint32 {
	if dist < 31 {
		if a<<(-dist&31) != 0 {
			return a>>dist | 1
		} else {
			return a >> dist
		}
	}
	if a != 0 {
		return 1
	} else {
		return 0
	}

}
