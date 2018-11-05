package arithmeticTypes

func softfloat_shiftRightJam128(a64, a0 uint64, dist uint32) Uint128 {
	var u8NegDist uint8
	var z Uint128

	if dist < 64 {
		u8NegDist = uint8(-dist)
		z.High = a64 >> dist
		//z.Low = a64<<(u8NegDist & 63) | a0>>dist | ( uint64(a0<<(u8NegDist & 63)) != 0)
		if uint64(a0<<(u8NegDist&63)) != 0 {
			z.Low = a64<<(u8NegDist&63) | a0>>dist | 1
		} else {
			z.Low = a64<<(u8NegDist&63) | a0>>dist | 0
		}
	} else {
		z.High = 0

		if dist < 127 {
			//z.Low = a64>>(dist & 63) | (((a64 & (( uint64(1)<< (dist & 63)) - 1)) | a0)  !=0)
			if ((a64 & ((uint64(1) << (dist & 63)) - 1)) | a0) != 0 {
				z.Low = a64>>(dist&63) | 1
			} else {
				z.Low = a64>>(dist&63) | 0
			}
		} else {
			if (a64 | a0) != 0 {
				z.Low = 1
			} else {
				z.Low = 0
			}
		}
	}
	return z
}

/*----------------------------------------------------------------------------
| Returns the difference of the 128-bit integer formed by concatenating 'a64'
| and 'a0' and the 128-bit integer formed by concatenating 'b64' and 'b0'.
| The subtraction is modulo 2^128, so any borrow out (carry out) is lost.
*----------------------------------------------------------------------------*/
func softfloat_sub128(a64, a0, b64, b0 uint64) Uint128 {
	var z Uint128
	z.Low = a0 - b0
	z.High = a64 - b64
	if a0 < b0 {
		z.High -= 1
	}
	return z
}

func softfloat_normRoundPackToF128(sign bool, exp int32, sig64 uint64, sig0 uint64) Float128 {
	var shiftDist int8
	var sig128 Uint128
	var uZ Float128
	var sigExtra uint64
	var sig128Extra uint128Extra
	var signToUint64 uint64
	if sign {
		signToUint64 = 1
	} else {
		signToUint64 = 0
	}

	if sig64 == 0 {
		exp -= 64
		sig64 = sig0
		sig0 = 0
	}

	shiftDist = int8(softfloat_countLeadingZeros64(sig64) - 15)
	exp -= int32(shiftDist)
	if 0 <= shiftDist {
		if shiftDist != 0 {
			sig128 = softfloat_shortShiftLeft128(sig64, sig0, uint8(shiftDist))
			sig64 = sig128.High
			sig0 = sig128.Low
		}
		if uint32(exp) < 0x7FFD {
			if sig64|sig0 != 0 {
				uZ.High = packToF128UI64(signToUint64, uint64(exp), sig64)
			} else {
				uZ.High = packToF128UI64(signToUint64, 0, sig64)
			}
			uZ.Low = sig0
			return uZ
		}
		sigExtra = 0
	} else {
		sig128Extra =
			softfloat_shortShiftRightJam128Extra(sig64, sig0, 0, uint8(-shiftDist))
		sig64 = sig128Extra.v.High
		sig0 = sig128Extra.v.Low
		sigExtra = sig128Extra.extra
	}
	return softfloat_roundPackToF128(sign, exp, sig64, sig0, sigExtra)
}
