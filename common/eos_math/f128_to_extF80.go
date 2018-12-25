package eos_math

func f128_to_extF80(a Float128) ExtFloat80_t {
	var uiA64, uiA0 uint64
	var sign bool
	var exp int32
	var frac64, frac0 uint64
	var commonNaN commonNaN
	var uiZ Uint128
	var uiZ64 uint16
	var uiZ0 uint64
	var normExpSig exp32_sig128
	var sig128 Uint128
	var uZ ExtFloat80_t

	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/

	uiA64 = a.High
	uiA0 = a.Low
	sign = signF128UI64(uiA64)
	exp = int32(expF128UI64(uiA64))
	frac64 = fracF128UI64(uiA64)
	frac0 = uiA0
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	if exp == 0x7FFF {
		if (frac64 | frac0) != 0 {
			softfloat_f128UIToCommonNaN(uiA64, uiA0, &commonNaN)
			uiZ = softfloat_commonNaNToExtF80UI(&commonNaN)
			uiZ64 = uint16(uiZ.High)
			uiZ0 = uiZ.Low
		} else {
			uiZ64 = packToExtF80UI64(sign, 0x7FFF)
			uiZ0 = uint64(0x8000000000000000)
		}
		goto uiZ
	}
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	if exp == 0 {
		if (frac64 | frac0) == 0 {
			uiZ64 = packToExtF80UI64(sign, 0)
			uiZ0 = 0
			goto uiZ
		}
		normExpSig = softfloat_normSubnormalF128Sig(frac64, frac0)
		exp = normExpSig.exp
		frac64 = normExpSig.sig.High
		frac0 = normExpSig.sig.Low
	}
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	sig128 =
		softfloat_shortShiftLeft128(
			frac64|uint64(0x0001000000000000), frac0, 15)
	return softfloat_roundPackToExtF80(sign, exp, sig128.High, sig128.Low, 80)
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
uiZ:
	uZ.signExp = uiZ64
	uZ.signIf = uiZ0
	return uZ

}

/*----------------------------------------------------------------------------
| Converts the common NaN pointed to by `aPtr' into an 80-bit extended
| floating-point NaN, and returns the bit pattern of this value as an unsigned
| integer.
*----------------------------------------------------------------------------*/
func softfloat_commonNaNToExtF80UI(aPtr *commonNaN) Uint128 {
	var uiZ Uint128
	if aPtr.sign {
		uiZ.High = uint64(1)<<15 | 0x7FFF
	} else {
		uiZ.High = 0x7FFF
	}

	uiZ.Low = uint64(0xC000000000000000) | aPtr.v64>>1
	return uiZ
}

//#define signExtF80UI64( a64 ) ((bool) ((uint16_t) (a64)>>15))
//#define expExtF80UI64( a64 ) ((a64) & 0x7FFF)
//#define packToExtF80UI64( sign, exp ) ((uint_fast16_t) (sign)<<15 | (exp))
//
//#define isNaNExtF80UI( a64, a0 ) ((((a64) & 0x7FFF) == 0x7FFF) && ((a0) & UINT64_C( 0x7FFFFFFFFFFFFFFF )))

func packToExtF80UI64(sign bool, exp uint16) uint16 {
	var re uint16
	if sign {
		re = (uint16(1) << 15) | exp
	} else {
		re = exp
	}
	return re
}

func softfloat_roundPackToExtF80(sign bool, exp int32, sig uint64, sigExtra uint64, roundingPrecision uint8) ExtFloat80_t {
	var roundingMode uint8
	var roundNearEven bool
	var roundIncrement, roundMask, roundBits uint64
	var isTiny, doIncrement bool
	var sig64Extra uint64Extra
	var uZ ExtFloat80_t

	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	var signMode uint8
	if sign {
		signMode = softfloat_round_min
	} else {
		signMode = softfloat_round_max
	}
	roundingMode = softfloat_roundingMode
	roundNearEven = roundingMode == softfloat_round_near_even
	if roundingPrecision == 80 {
		goto precision80
	}
	if roundingPrecision == 64 {
		roundIncrement = uint64(0x0000000000000400)
		roundMask = uint64(0x00000000000007FF)
	} else if roundingPrecision == 32 {
		roundIncrement = uint64(0x0000008000000000)
		roundMask = uint64(0x000000FFFFFFFFFF)
	} else {
		goto precision80
	}

	if sigExtra != 0 {
		sig |= 1
	}
	if !roundNearEven && (roundingMode != softfloat_round_near_maxMag) {
		if roundingMode == signMode {
			roundIncrement = roundMask
		} else {
			roundIncrement = 0
		}
	}
	roundBits = sig & roundMask
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	if 0x7FFD <= uint32(exp-1) {
		if exp <= 0 {
			/*----------------------------------------------------------------
			*----------------------------------------------------------------*/
			isTiny = (softfloat_detectTininess == softfloat_tininess_beforeRounding) || (exp < 0) || (sig <= uint64(sig+roundIncrement))
			sig = softfloat_shiftRightJam64(sig, uint32(1-exp))
			roundBits = sig & roundMask
			if roundBits != 0 {
				if isTiny {
					softfloat_raiseFlags(softfloat_flag_underflow)
				}
				//softfloat_exceptionFlags |= softfloat_flag_inexact
			}
			sig += roundIncrement
			if (sig & uint64(0x8000000000000000)) != 0 {
				exp = 1
			} else {
				exp = 0
			}
			roundIncrement = roundMask + 1
			if roundNearEven && (roundBits<<1 == roundIncrement) {
				roundMask |= roundIncrement
			}
			sig &= ^roundMask
			goto packReturn
		}
		if (0x7FFE < exp) || ((exp == 0x7FFE) && (uint64(sig+roundIncrement) < sig)) {
			//goto overflow
			softfloat_raiseFlags(softfloat_flag_overflow | softfloat_flag_inexact)
			if roundNearEven || (roundingMode == softfloat_round_near_maxMag) || (roundingMode == signMode) {
				exp = 0x7FFF
				sig = uint64(0x8000000000000000)
			} else {
				exp = 0x7FFE
				sig = ^roundMask
			}
			goto packReturn
		}
	}

	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/

	sig = sig + roundIncrement
	if sig < roundIncrement {
		exp += 1
		sig = uint64(0x8000000000000000)
	}
	roundIncrement = roundMask + 1
	if roundNearEven && (roundBits<<1 == roundIncrement) {
		roundMask |= roundIncrement
	}
	sig &= ^roundMask
	goto packReturn
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/

precision80:
	doIncrement = uint64(0x8000000000000000) <= sigExtra
	if !roundNearEven && (roundingMode != softfloat_round_near_maxMag) {
		doIncrement = (roundingMode == signMode) && sigExtra != 0
	}
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	if 0x7FFD <= uint32(exp-1) {
		if exp <= 0 {
			/*----------------------------------------------------------------
			*----------------------------------------------------------------*/
			isTiny = (softfloat_detectTininess == softfloat_tininess_beforeRounding) || (exp < 0) || !doIncrement || (sig < uint64(0xFFFFFFFFFFFFFFFF))
			sig64Extra = softfloat_shiftRightJam64Extra(sig, sigExtra, uint32(1-exp))
			exp = 0
			sig = sig64Extra.v
			sigExtra = sig64Extra.extra
			if sigExtra != 0 {
				if isTiny {
					softfloat_raiseFlags(softfloat_flag_underflow)
				}
			}

			doIncrement = uint64(0x8000000000000000) <= sigExtra
			if !roundNearEven && (roundingMode != softfloat_round_near_maxMag) {
				doIncrement = roundingMode == signMode && sigExtra != 0
			}
			if doIncrement {
				sig += 1
				//sig &= ~(uint64)(! (sigExtra & uint64(0x7FFFFFFFFFFFFFFF)) & roundNearEven)
				if sigExtra&uint64(0x7FFFFFFFFFFFFFFF) == 0 && roundNearEven {
					sig &= ^uint64(1)
				} else {
					sig &= ^uint64(0)
				}
				//exp = (sig & uint64(0x8000000000000000)) != 0
				if (sig & uint64(0x8000000000000000)) != 0 {
					exp = 1
				} else {
					exp = 0
				}
			}
			goto packReturn
		}
		if (0x7FFE < exp) || ((exp == 0x7FFE) && (sig == uint64(0xFFFFFFFFFFFFFFFF)) && doIncrement) {
			/*----------------------------------------------------------------
			*----------------------------------------------------------------*/
			roundMask = 0
			//overflow:
			softfloat_raiseFlags(softfloat_flag_overflow | softfloat_flag_inexact)
			if roundNearEven || (roundingMode == softfloat_round_near_maxMag) || (roundingMode == signMode) {
				exp = 0x7FFF
				sig = uint64(0x8000000000000000)
			} else {
				exp = 0x7FFE
				sig = ^roundMask
			}
			goto packReturn
		}
	}

	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	if doIncrement {
		sig += 1
		if sig == 0 {
			exp += 1
			sig = uint64(0x8000000000000000)
		} else {
			if (sigExtra&uint64(0x7FFFFFFFFFFFFFFF)) == 0 && roundNearEven {
				sig &= ^uint64(1)
			} else {
				sig &= ^uint64(0)
			}
		}
	}
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
packReturn:
	uZ.signExp = packToExtF80UI64(sign, uint16(exp))
	uZ.signIf = sig
	return uZ
}
