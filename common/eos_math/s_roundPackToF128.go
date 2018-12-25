package eos_math

func softfloat_roundPackToF128(sign bool, exp int32, sig64 uint64, sig0 uint64, sigExtra uint64) Float128 {
	var roundingMode uint8
	var roundNearEven, doIncrement, isTiny bool
	var sig128Extra uint128Extra
	var uiZ64, uiZ0 uint64
	var sig128 Uint128
	//union ui128_f128 uZ

	var f128 Float128

	var modeEnum uint8
	var signToUint64 uint64
	if sign {
		modeEnum = softfloat_round_min
		signToUint64 = 1
	} else {
		modeEnum = softfloat_round_max
		signToUint64 = 0
	}

	//
	//    /*------------------------------------------------------------------------
	//    *------------------------------------------------------------------------*/
	roundingMode = softfloat_roundingMode
	roundNearEven = roundingMode == softfloat_round_near_even
	doIncrement = uint64(0x8000000000000000) <= sigExtra
	if !roundNearEven && (roundingMode != softfloat_round_near_maxMag) {
		doIncrement = (roundingMode == modeEnum) && (sigExtra != 0)
	}
	//    /*------------------------------------------------------------------------
	//    *------------------------------------------------------------------------*/

	if 0x7FFD <= uint32(exp) {
		if exp < 0 {
			//			   /*----------------------------------------------------------------
			//	           *----------------------------------------------------------------*/
			isTiny =
				(softfloat_detectTininess == softfloat_tininess_beforeRounding) || (exp < -1) || !doIncrement ||
					softfloat_lt128(sig64, sig0, uint64(0x0001FFFFFFFFFFFF), uint64(0xFFFFFFFFFFFFFFFF))
			sig128Extra =
				softfloat_shiftRightJam128Extra(sig64, sig0, sigExtra, uint32(-exp))
			sig64 = sig128Extra.v.High
			sig0 = sig128Extra.v.Low
			sigExtra = sig128Extra.extra
			exp = 0
			if isTiny && (sigExtra != 0) {
				softfloat_raiseFlags(softfloat_flag_underflow)
			}
			doIncrement = uint64(0x8000000000000000) <= sigExtra
			if !roundNearEven && (roundingMode != softfloat_round_near_maxMag) {
				doIncrement = (roundingMode == modeEnum) && (sigExtra != 0)
			}
		} else if (0x7FFD < exp) || ((exp == 0x7FFD) && softfloat_eq128(sig64, sig0, uint64(0x0001FFFFFFFFFFFF), uint64(0xFFFFFFFFFFFFFFFF)) && doIncrement) {
			/*----------------------------------------------------------------
			 *----------------------------------------------------------------*/
			softfloat_raiseFlags(softfloat_flag_overflow | softfloat_flag_inexact)

			if roundNearEven || (roundingMode == softfloat_round_near_maxMag) || (roundingMode == modeEnum) {
				uiZ64 = packToF128UI64(signToUint64, 0x7FFF, 0)
				uiZ0 = 0
			} else {
				uiZ64 = packToF128UI64(signToUint64, 0x7FFE, uint64(0x0000FFFFFFFFFFFF))
				uiZ0 = uint64(0xFFFFFFFFFFFFFFFF)
			}
			goto uiZ
		}

	}

	/*------------------------------------------------------------------------
	 *------------------------------------------------------------------------*/
	if doIncrement {
		sig128 = softfloat_add128(sig64, sig0, 0, 1)
		sig64 = sig128.High
		//sig0 = sig128.Low & ~(uint64_t)(! (sigExtra & uint64( 0x7FFFFFFFFFFFFFFF )) & roundNearEven)
		if (sigExtra&uint64(0x7FFFFFFFFFFFFFFF) == 0) && roundNearEven {
			//sig0 = sig128.Low &  (^uint64(1))
			sig0 = sig128.Low & uint64(0xFEFFFFFFFFFFFFFF)
		} else {
			//sig0 = sig128.Low &  (^uint64(0))
			sig0 = sig128.Low & uint64(0xFFFFFFFFFFFFFFFF)
		}
	} else {
		if (sig64 | sig0) == 0 {
			exp = 0
		}
	}

	/*------------------------------------------------------------------------
	 *------------------------------------------------------------------------*/
	//packReturn:
	uiZ64 = packToF128UI64(signToUint64, uint64(exp), sig64)
	uiZ0 = sig0
uiZ:
	f128.High = uiZ64
	f128.Low = uiZ0

	return f128

}

/*----------------------------------------------------------------------------
| Returns true if the 128-bit unsigned integer formed by concatenating 'a64'
| and 'a0' is equal to the 128-bit unsigned integer formed by concatenating
| 'b64' and 'b0'.
*----------------------------------------------------------------------------*/
func softfloat_eq128(a64, a0, b64, b0 uint64) bool {
	return (a64 == b64) && (a0 == b0)
}

/*----------------------------------------------------------------------------
| Returns the sum of the 128-bit integer formed by concatenating 'a64' and
| 'a0' and the 128-bit integer formed by concatenating 'b64' and 'b0'.  The
| addition is modulo 2^128, so any carry out is lost.
*----------------------------------------------------------------------------*/
func softfloat_add128(a64, a0, b64, b0 uint64) Uint128 {
	var z Uint128

	z.Low = a0 + b0
	if z.Low < a0 {
		z.High = a64 + b64 + 1
	} else {
		z.High = a64 + b64
	}
	return z
}

/*----------------------------------------------------------------------------
| Interpreting the unsigned integer formed from concatenating `uiA64' and
| `uiA0' as a 128-bit floating-point value, and likewise interpreting the
| unsigned integer formed from concatenating `uiB64' and `uiB0' as another
| 128-bit floating-point value, and assuming at least on of these floating-
| point values is a NaN, returns the bit pattern of the combined NaN result.
| If either original floating-point value is a signaling NaN, the invalid
| exception is raised.
*----------------------------------------------------------------------------*/
func softfloat_propagateNaNF128UI(uiA64, uiA0, uiB64, uiB0 uint64) Uint128 {
	var isSigNaNA bool
	var uiZ Uint128

	isSigNaNA = softfloat_isSigNaNF128UI(uiA64, uiA0)
	if isSigNaNA || softfloat_isSigNaNF128UI(uiB64, uiB0) {
		softfloat_raiseFlags(softfloat_flag_invalid)
		if isSigNaNA {
			uiZ.High = uiA64
			uiZ.Low = uiA0
			uiZ.High |= uint64(0x0000800000000000)
			return uiZ
		}
	}

	if isNaNF128UI(uiA64, uiA0) {
		uiZ.High = uiA64
		uiZ.Low = uiA0
	} else {
		uiZ.High = uiB64
		uiZ.Low = uiB0
	}

	uiZ.High |= uint64(0x0000800000000000)
	return uiZ
}
