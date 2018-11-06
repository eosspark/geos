package arithmeticTypes

func (a Float128) Div(b Float128) Float128 {
	var uiA64, uiA0 uint64
	var signA bool
	var expA int32
	var sigA Uint128
	var uiB64, uiB0 uint64
	var signB bool
	var expB int32
	var sigB Uint128
	var signZ bool
	var normExpSig exp32_sig128
	var expZ int32
	var rem Uint128
	var recip32 uint32
	var ix int
	var q64 uint64
	var q uint32
	var term Uint128
	var qs [3]uint32
	var sigZExtra uint64
	var sigZ, uiZ Uint128
	var uZ Float128

	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	uiA64 = a.High
	uiA0 = a.Low
	signA = signF128UI64(uiA64)
	expA = int32(expF128UI64(uiA64))
	sigA.High = fracF128UI64(uiA64)
	sigA.Low = uiA0

	uiB64 = b.High
	uiB0 = b.Low
	signB = signF128UI64(uiB64)
	expB = int32(expF128UI64(uiB64))
	sigB.High = fracF128UI64(uiB64)
	sigB.Low = uiB0
	//signZ = signA ^ signB
	if (signA && signB) || (signA == false && signB == false) {
		signZ = false
	} else {
		signZ = true
	}

	/*------------------------------------------------------------------------
	 *------------------------------------------------------------------------*/
	if expA == 0x7FFF {
		if (sigA.High | sigA.Low) != 0 {
			goto propagateNaN
		}
		if expB == 0x7FFF {
			if (sigB.High | sigB.Low) != 0 {
				goto propagateNaN
			}
			goto invalid
		}
		goto infinity
	}
	if expB == 0x7FFF {
		if (sigB.High | sigB.Low) != 0 {
			goto propagateNaN
		}
		goto zero
	}
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	if expB == 0 {
		if (sigB.High | sigB.Low) == 0 {
			if (uint64(expA) | sigA.High | sigA.Low) == 0 {
				goto invalid
			}
			softfloat_raiseFlags(softfloat_flag_infinite)
			goto infinity
		}
		normExpSig = softfloat_normSubnormalF128Sig(sigB.High, sigB.Low)
		expB = normExpSig.exp
		sigB = normExpSig.sig
	}
	if expA == 0 {
		if (sigA.High | sigA.Low) == 0 {
			goto zero
		}
		normExpSig = softfloat_normSubnormalF128Sig(sigA.High, sigA.Low)
		expA = normExpSig.exp
		sigA = normExpSig.sig
	}
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	expZ = expA - expB + 0x3FFE
	sigA.High |= uint64(0x0001000000000000)
	sigB.High |= uint64(0x0001000000000000)
	rem = sigA
	if softfloat_lt128(sigA.High, sigA.Low, sigB.High, sigB.Low) {
		expZ -= 1
		rem = softfloat_add128(sigA.High, sigA.Low, sigA.High, sigA.Low)
	}
	recip32 = softfloat_approxRecip32_1(sigB.High >> 17)
	ix = 3
	for {
		q64 = rem.High >> 19 * uint64(recip32)
		q = uint32((q64 + 0x80000000) >> 32)
		ix -= 1
		if ix < 0 {
			break
		}

		rem = softfloat_shortShiftLeft128(rem.High, rem.Low, 29)
		term = softfloat_mul128By32(sigB.High, sigB.Low, q)
		rem = softfloat_sub128(rem.High, rem.Low, term.High, term.Low)
		if (rem.High & uint64(0x8000000000000000)) != 0 {
			q -= 1
			rem = softfloat_add128(rem.High, rem.Low, sigB.High, sigB.Low)
		}
		qs[ix] = q
	}

	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	if ((q + 1) & 7) < 2 {
		rem = softfloat_shortShiftLeft128(rem.High, rem.Low, 29)
		term = softfloat_mul128By32(sigB.High, sigB.Low, q)
		rem = softfloat_sub128(rem.High, rem.Low, term.High, term.Low)
		if (rem.High & uint64(0x8000000000000000)) != 0 {
			q -= 1
			rem = softfloat_add128(rem.High, rem.Low, sigB.High, sigB.Low)
		} else if softfloat_le128(sigB.High, sigB.Low, rem.High, rem.Low) {
			q += 1
			rem = softfloat_sub128(rem.High, rem.Low, sigB.High, sigB.Low)
		}
		if (rem.High | rem.Low) != 0 {
			q |= 1
		}
	}
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	sigZExtra = uint64(q) << 60
	term = softfloat_shortShiftLeft128(0, uint64(qs[1]), 54)
	sigZ = softfloat_add128(uint64(qs[2])<<19, (uint64(qs[0])<<25)+uint64(q>>4), term.High, term.Low)
	return softfloat_roundPackToF128(signZ, expZ, sigZ.High, sigZ.Low, sigZExtra)
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
propagateNaN:
	uiZ = softfloat_propagateNaNF128UI(uiA64, uiA0, uiB64, uiB0)
	goto uiZ
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
invalid:
	softfloat_raiseFlags(softfloat_flag_invalid)
	uiZ.High = defaultNaNF128UI64
	uiZ.Low = defaultNaNF128UI0
	goto uiZ
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
infinity:
	//uiZ.High = packToF128UI64( signZ, 0x7FFF, 0 )
	if signZ {
		uiZ.High = packToF128UI64(1, 0x7FFF, 0)
	} else {
		uiZ.High = packToF128UI64(0, 0x7FFF, 0)
	}
	goto uiZ0
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
zero:
	//uiZ.High = packToF128UI64( signZ, 0, 0 )
	if signZ {
		uiZ.High = packToF128UI64(1, 0, 0)
	} else {
		uiZ.High = packToF128UI64(0, 0, 0)
	}
uiZ0:
	uiZ.Low = 0
uiZ:
	uZ.High = uiZ.High
	uZ.Low = uiZ.Low
	return uZ
}

/*----------------------------------------------------------------------------
| Returns an approximation to the reciprocal of the number represented by 'a',
| where 'a' is interpreted as an unsigned fixed-point number with one integer
| bit and 31 fraction bits.  The 'a' input must be "normalized", meaning that
| its most-significant bit (bit 31) must be 1.  Thus, if A is the value of
| the fixed-point interpretation of 'a', then 1 <= A < 2.  The returned value
| is interpreted as a pure unsigned fraction, having no integer bits and 32
| fraction bits.  The approximation returned is never greater than the true
| reciprocal 1/A, and it differs from the true reciprocal by at most 2.006 ulp
| (units in the last place).
*----------------------------------------------------------------------------*/
//#define softfloat_approxRecip32_1( a ) ((uint32_t) (UINT64_C( 0x7FFFFFFFFFFFFFFF ) / (uint32_t) (a)))
func softfloat_approxRecip32_1(a uint64) uint32 {
	return uint32(uint64(0x7FFFFFFFFFFFFFFF) / a)
}

/*----------------------------------------------------------------------------
| Returns the product of the 128-bit integer formed by concatenating 'a64' and
| 'a0', multiplied by 'b'.  The multiplication is modulo 2^128; any overflow
| bits are discarded.
*----------------------------------------------------------------------------*/
func softfloat_mul128By32(a64, a0 uint64, b uint32) Uint128 {
	var z Uint128
	var mid uint64
	var carry uint32

	z.Low = a0 * uint64(b)
	mid = (a0 >> 32) * uint64(b)
	carry = uint32(z.Low>>32) - uint32(mid)
	z.High = a64*uint64(b) + uint64((mid+uint64(carry))>>32)
	return z
}

/*----------------------------------------------------------------------------
| Returns true if the 128-bit unsigned integer formed by concatenating 'a64'
| and 'a0' is less than or equal to the 128-bit unsigned integer formed by
| concatenating 'b64' and 'b0'.
*----------------------------------------------------------------------------*/
func softfloat_le128(a64, a0, b64, b0 uint64) bool {
	return (a64 < b64) || ((a64 == b64) && (a0 <= b0))
}
