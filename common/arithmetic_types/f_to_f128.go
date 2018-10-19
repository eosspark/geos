package arithmeticTypes

type exp16_sig32 struct {
	exp int
	sig uint32
}

func F32ToF128(a Float32) Float128 {
	var uA Float32
	var uiA uint32
	var sign bool
	var exp int16
	var frac uint32
	var commonNaN commonNaN
	var uiZ Uint128
	var normExpSig exp16_sig32

	uA = a
	uiA = uint32(uA)
	sign = signF32UI(uiA)
	exp = expF32UI(uiA)
	frac = fracF32UI(uiA)

	if exp == 0xFF {
		if frac != 0 {
			softfloat_f32UIToCommonNaN(uiA, &commonNaN)
			uiZ = softfloat_commonNaNToF128UI(&commonNaN)
		} else {
			if sign {
				uiZ.High = packToF128UI64(1, 0x7FFF, 0)
			} else {
				uiZ.High = packToF128UI64(0, 0x7FFF, 0)
			}
			uiZ.Low = 0
		}
		goto uiZ
	}

	if exp == 0 {
		if frac == 0 {
			if sign {
				uiZ.High = packToF128UI64(1, 0, 0)
			} else {
				uiZ.High = packToF128UI64(0, 0, 0)
			}
			uiZ.Low = 0
			goto uiZ
		}
		normExpSig = softfloat_normSubnormalF32Sig(frac)
		exp = int16(normExpSig.exp - 1)
		frac = normExpSig.sig
	}
	if sign {
		uiZ.High = packToF128UI64(1, uint64(exp+0x3F80), uint64(frac)<<25)
	} else {
		uiZ.High = packToF128UI64(0, uint64(exp+0x3F80), uint64(frac)<<25)
	}

	uiZ.Low = 0
uiZ:
	return Float128(uiZ)
}

/*----------------------------------------------------------------------------
| Assuming `uiA' has the bit pattern of a 32-bit floating-point NaN, converts
| this NaN to the common NaN form, and stores the resulting common NaN at the
| location pointed to by `zPtr'.  If the NaN is a signaling NaN, the invalid
| exception is raised.
*----------------------------------------------------------------------------*/
func softfloat_f32UIToCommonNaN(uiA uint32, zPtr *commonNaN) {
	if uiA>>31 == 0 {
		zPtr.sign = false
	} else {
		zPtr.sign = true
	}
	zPtr.v64 = uint64(uiA) << 41
	zPtr.v0 = 0
}

/*----------------------------------------------------------------------------
| Converts the common NaN pointed to by `aPtr' into a 128-bit floating-point
| NaN, and returns the bit pattern of this value as an unsigned integer.
*----------------------------------------------------------------------------*/
func softfloat_commonNaNToF128UI(aPtr *commonNaN) Uint128 {

	uiZ := softfloat_shortShiftRight128(aPtr.v64, aPtr.v0, 16)
	if aPtr.sign {
		uiZ.High |= uint64(1)<<63 | uint64(0x7FFF800000000000)
	} else {
		uiZ.High |= uint64(0x7FFF800000000000)
	}
	return uiZ
}

func softfloat_shortShiftRight128(a64, a0 uint64, dist uint8) Uint128 {
	var z Uint128
	z.High = a64 >> dist
	z.Low = a64<<(-dist&63) | a0>>dist
	return z
}

func softfloat_normSubnormalF32Sig(sig uint32) exp16_sig32 {
	var shiftDist uint8
	var z exp16_sig32

	shiftDist = softfloat_countLeadingZeros32(sig) - 8
	z.exp = int(1 - shiftDist)
	z.sig = sig << shiftDist
	return z
}

func softfloat_countLeadingZeros32(a uint32) uint8 {
	count := uint8(0)
	if a < 0x10000 {
		count = 16
		a <<= 16
	}
	if a < 0x1000000 {
		count += 8
		a <<= 8
	}
	count += softfloatCountLeadingZeros8[a>>24]
	return count
}

type exp16_sig64 struct {
	exp int16
	sig uint64
}

func F64ToF128(a Float64) Float128 {
	var uiA uint64
	var sign bool
	var exp int16
	var frac uint64
	var commonNaN commonNaN
	var uiZ Uint128
	var normExpSig exp16_sig64
	var frac128 Uint128

	uiA = uint64(a)
	sign = signF64UI(uiA)
	exp = expF64UI(uiA)
	frac = fracF64UI(uiA)

	if exp == 0x7FF {
		if frac != 0 {
			softfloat_f64UIToCommonNaN(uiA, &commonNaN)
			uiZ = softfloat_commonNaNToF128UI(&commonNaN)
		} else {
			if sign {
				uiZ.High = packToF128UI64(1, 0x7FFF, 0)
			} else {
				uiZ.High = packToF128UI64(0, 0x7FFF, 0)
			}
			uiZ.Low = 0
		}
		goto uiZ
	}

	if exp == 0 {
		if frac == 0 {
			if sign {
				uiZ.High = packToF128UI64(1, 0, 0)
			} else {
				uiZ.High = packToF128UI64(0, 0, 0)
			}
			uiZ.Low = 0
			goto uiZ
		}
		normExpSig = softfloat_normSubnormalF64Sig(frac)
		exp = normExpSig.exp - 1
		frac = normExpSig.sig
	}

	frac128 = softfloat_shortShiftLeft128(0, frac, 60)
	if sign {
		uiZ.High = packToF128UI64(1, uint64(exp+0x3C00), frac128.High)
	} else {
		uiZ.High = packToF128UI64(0, uint64(exp+0x3C00), frac128.High)
	}

	uiZ.Low = frac128.Low

uiZ:
	return Float128(uiZ)

}

/*----------------------------------------------------------------------------
| Assuming `uiA' has the bit pattern of a 64-bit floating-point NaN, converts
| this NaN to the common NaN form, and stores the resulting common NaN at the
| location pointed to by `zPtr'.  If the NaN is a signaling NaN, the invalid
| exception is raised.
*----------------------------------------------------------------------------*/
func softfloat_f64UIToCommonNaN(uiA uint64, zPtr *commonNaN) {
	if uiA>>63 == 0 {
		zPtr.sign = false
	} else {
		zPtr.sign = true
	}
	zPtr.v64 = uiA << 12
	zPtr.v0 = 0
}

func softfloat_normSubnormalF64Sig(sig uint64) exp16_sig64 {
	var shiftDist int8
	var z exp16_sig64

	shiftDist = int8(softfloat_countLeadingZeros64(sig) - 11)
	z.exp = int16(1 - shiftDist)
	z.sig = sig << uint8(shiftDist)
	return z
}

func softfloat_countLeadingZeros64(a uint64) uint8 {
	var count uint8
	var a32 uint32

	count = 0
	a32 = uint32(a >> 32)
	if a32 == 0 {
		count = 32
		a32 = uint32(a)
	}
	/*------------------------------------------------------------------------
	| From here, result is current count + count leading zeros of `a32'.
	*------------------------------------------------------------------------*/
	if a32 < 0x10000 {
		count += 16
		a32 <<= 16
	}
	if a32 < 0x1000000 {
		count += 8
		a32 <<= 8
	}
	count += softfloatCountLeadingZeros8[a32>>24]
	return count
}
