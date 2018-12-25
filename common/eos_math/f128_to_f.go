package eos_math

import "unsafe"

//THREAD_LOCAL uint_fast8_t softfloat_roundingMode = softfloat_round_near_even;
//THREAD_LOCAL uint_fast8_t softfloat_detectTininess = init_detectTininess;
//THREAD_LOCAL uint_fast8_t softfloat_exceptionFlags = 0;
//
//THREAD_LOCAL uint_fast8_t extF80_roundingPrecision = 80;

const (
	softfloat_roundingMode   = softfloat_round_near_even
	softfloat_detectTininess = softfloat_tininess_beforeRounding //ARM-VFPv2/ARM_VFPv2_defaultNaN
	//softfloat_detectTininess = softfloat_tininess_afterRounding//8086/RISCV/8086-SSE
	softfloat_exceptionFlags = 0
	extF80_roundingPrecision = 80
)

const (
	softfloat_tininess_beforeRounding uint8 = 0
	softfloat_tininess_afterRounding  uint8 = 1
)

/*----------------------------------------------------------------------------
| Default value for 'softfloat_detectTininess'.
*----------------------------------------------------------------------------*/

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
func softfloat_isSigNaNF128UI(uiA64, uiA0 uint64) bool {
	return (uiA64&uint64(0x7FFF800000000000)) == uint64(0x7FFF000000000000) && (uiA0 != 0 || (uiA64&uint64(0x00007FFFFFFFFFFF) != 0))
}

/*----------------------------------------------------------------------------
| Converts the common NaN pointed to by `aPtr' into a 32-bit floating-point
| NaN, and returns the bit pattern of this value as an unsigned integer.
*----------------------------------------------------------------------------*/
func softfloat_commonNaNToF32UI(aPtr *commonNaN) uint32 {
	if aPtr.sign {
		return uint32(1)<<31 | uint32(0x7FC00000) | uint32(aPtr.v64>>41)
	} else {
		return uint32(0x7FC00000) | uint32(aPtr.v64>>41)
	}
	//return (uint_fast32_t) aPtr->sign<<31 | 0x7FC00000 | aPtr->v64>>41

}

//#define signF32UI( a ) ((bool) ((uint32_t) (a)>>31))
//#define expF32UI( a ) ((int_fast16_t) ((a)>>23) & 0xFF)
//#define fracF32UI( a ) ((a) & 0x007FFFFF)
//#define packToF32UI( sign, exp, sig ) (((uint32_t) (sign)<<31) + ((uint32_t) (exp)<<23) + (sig))

func signF32UI(a uint32) bool {
	if a>>31 == 0 {
		return false
	}
	return true
}
func expF32UI(a uint32) int16 {
	return int16(a>>23) & 0xFF
}
func fracF32UI(a uint32) uint32 {
	return a & 0x007FFFFF
}

func packToF32UI(sign bool, exp, sig uint32) uint32 {
	if sign {
		return uint32(1)<<31 + uint32(exp<<23) + sig
	} else {
		return uint32(exp<<23) + sig
	}

}

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
			isTiny = (softfloat_detectTininess == softfloat_tininess_beforeRounding) || (exp < -1) || (sig+uint32(roundIncrement) < uint32(0x80000000))
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

func F128ToF64(a Float128) Float64 {
	var uA Float128
	var uiA64, uiA0 uint64
	var sign bool
	var exp int32
	var frac64, frac0 uint64
	var commonNaN commonNaN
	var uiZ uint64
	var frac128 Uint128

	uA = a
	uiA64 = uA.High
	uiA0 = uA.Low
	sign = signF128Ui64(uiA64)
	exp = expF128Ui64(uiA64)
	frac64 = fracF128Ui64(uiA64)
	frac0 = uiA0

	if exp == 0x7FFF {
		if (frac64 | frac0) != 0 {
			softfloat_f128UIToCommonNaN(uiA64, uiA0, &commonNaN)
			uiZ = softfloat_commonNaNToF64UI(&commonNaN)
		} else {
			uiZ = packToF64UI(sign, 0x7FF, 0)
		}
		goto uiZ
	}

	frac128 = softfloat_shortShiftLeft128(frac64, frac0, 14)
	if frac128.Low != 0 {
		frac64 = frac128.High | 1
	} else {
		frac64 = frac128.High
	}

	if int64(exp)|int64(frac64) == 0 {
		uiZ = packToF64UI(sign, 0, 0)
		goto uiZ
	}

	exp -= 0x3C01
	if unsafe.Sizeof(int16(0)) < unsafe.Sizeof(int32(0)) {
		if exp < -0x1000 {
			exp = -0x1000
		}
	}

	return softfloat_roundPackToF64(sign, int16(exp), frac64|uint64(0x4000000000000000))

uiZ:
	return Float64(uiZ)

}

/*----------------------------------------------------------------------------
| Converts the common NaN pointed to by `aPtr' into a 64-bit floating-point
| NaN, and returns the bit pattern of this value as an unsigned integer.
*----------------------------------------------------------------------------*/
func softfloat_commonNaNToF64UI(aPtr *commonNaN) uint64 {
	if aPtr.sign {
		return uint64(1)<<63 | uint64(0x7FF8000000000000) | aPtr.v64>>12
	}
	return uint64(0x7FF8000000000000) | aPtr.v64>>12
}

//#define signF64UI( a ) ((bool) ((uint64_t) (a)>>63))
//#define expF64UI( a ) ((int_fast16_t) ((a)>>52) & 0x7FF)
//#define fracF64UI( a ) ((a) & UINT64_C( 0x000FFFFFFFFFFFFF ))
//#define packToF64UI( sign, exp, sig ) ((uint64_t) (((uint_fast64_t) (sign)<<63) + ((uint_fast64_t) (exp)<<52) + (sig)))

func signF64UI(a uint64) bool {
	if a>>63 == 0 {
		return false
	}
	return true
}

func expF64UI(a uint64) int16 {
	return int16(a>>52) & 0x7FF
}

func fracF64UI(a uint64) uint64 {
	return a & uint64(0x000FFFFFFFFFFFFF)
}

func packToF64UI(sign bool, exp, sig uint64) uint64 {
	if sign {
		return uint64(1)<<63 + exp<<52 + sig
	}
	return exp<<52 + sig //return uint64(0)<< 63 + exp <<52 + sig
}

func softfloat_roundPackToF64(sign bool, exp int16, sig uint64) Float64 {
	var roundingMode uint8
	var roundNearEven bool
	var roundIncrement, roundBits uint16
	var isTiny bool
	var uiZ uint64

	roundingMode = softfloat_roundingMode
	roundNearEven = roundingMode == softfloat_round_near_even
	roundIncrement = 0x200
	if !(roundNearEven && (roundingMode != softfloat_round_near_maxMag)) {
		var mode uint8
		if sign {
			mode = softfloat_round_min
		} else {
			mode = softfloat_round_max
		}
		if roundingMode == mode {
			roundIncrement = 0x3FF
		} else {
			roundIncrement = 0
		}
	}
	roundBits = uint16(sig & 0x3FF)

	if 0x7FD <= uint16(exp) {
		if exp < 0 {
			/*----------------------------------------------------------------
			*----------------------------------------------------------------*/
			isTiny = (softfloat_detectTininess == softfloat_tininess_beforeRounding) || (exp < -1) || (sig+uint64(roundIncrement) < uint64(0x8000000000000000))
			sig = softfloat_shiftRightJam64(sig, uint32(-exp))
			exp = 0
			roundBits = uint16(sig & 0x3FF)
			if isTiny && roundBits != 0 {
				softfloat_raiseFlags(softfloat_flag_underflow)
			}
		} else if (0x7FD < exp) || (uint64(0x8000000000000000) <= sig+uint64(roundIncrement)) {
			/*----------------------------------------------------------------
			*----------------------------------------------------------------*/
			//softfloat_raiseFlags(softfloat_flag_overflow | softfloat_flag_inexact)

			if roundIncrement != 0 {
				uiZ = packToF64UI(sign, 0xFF, 0)
			} else {
				uiZ = packToF64UI(sign, 0xFF, 0) - 1
			}
			goto uiZ
		}
	}

	sig = (sig + uint64(roundIncrement)) >> 10
	if roundBits^0x200 != 0 {
		sig &= ^uint64(0)
	} else {
		if roundNearEven {
			sig &= ^uint64(1)
		} else {
			sig &= ^uint64(0)
		}
	}

	if sig == 0 {
		exp = 0
	}

	uiZ = packToF64UI(sign, uint64(exp), sig)

uiZ:
	return Float64(uiZ)

}
