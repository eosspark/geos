package eos_math

func (a Float128) Mul(b Float128) Float128 {
	var uiA64, uiA0 uint64
	var signA bool
	var expA int32
	var sigA Uint128
	var uiB64, uiB0 uint64
	var signB bool
	var expB int32
	var sigB Uint128
	var signZ bool
	var magBits uint64
	var normExpSig exp32_sig128
	var expZ int32
	var sig256Z [4]uint64
	var sigZExtra uint64
	var sigZ Uint128
	var sig128Extra uint128Extra
	var uiZ Uint128
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
		if (sigA.High|sigA.Low) != 0 || ((expB == 0x7FFF) && (sigB.High|sigB.Low) != 0) {
			goto propagateNaN
		}
		magBits = uint64(expB) | sigB.High | sigB.Low
		goto infArg
	}
	if expB == 0x7FFF {
		if sigB.High|sigB.Low != 0 {
			goto propagateNaN
		}
		magBits = uint64(expA) | sigA.High | sigA.Low
		goto infArg
	}
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	if expA == 0 {
		if (sigA.High | sigA.Low) == 0 {
			goto zero
		}
		normExpSig = softfloat_normSubnormalF128Sig(sigA.High, sigA.Low)
		expA = normExpSig.exp
		sigA = normExpSig.sig
	}
	if expB == 0 {
		if (sigB.High | sigB.Low) == 0 {
			goto zero
		}
		normExpSig = softfloat_normSubnormalF128Sig(sigB.High, sigB.Low)
		expB = normExpSig.exp
		sigB = normExpSig.sig
	}
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
	expZ = expA + expB - 0x4000
	sigA.High |= uint64(0x0001000000000000)
	sigB = softfloat_shortShiftLeft128(sigB.High, sigB.Low, 16)
	sig256Z = softfloat_mul128To256M(sigA.High, sigA.Low, sigB.High, sigB.Low, sig256Z) //TODO
	//sigZExtra = sig256Z[indexWord( 4, 1 )] | (sig256Z[indexWord( 4, 0 )] != 0)
	if sig256Z[indexWord(4, 0)] != 0 {
		sigZExtra = sig256Z[indexWord(4, 1)] | 1
	} else {
		sigZExtra = sig256Z[indexWord(4, 1)] | 0
	}

	sigZ = softfloat_add128(sig256Z[indexWord(4, 3)], sig256Z[indexWord(4, 2)], sigA.High, sigA.Low)
	if uint64(0x0002000000000000) <= sigZ.High {
		expZ += 1
		sig128Extra = softfloat_shortShiftRightJam128Extra(sigZ.High, sigZ.Low, sigZExtra, 1)
		sigZ = sig128Extra.v
		sigZExtra = sig128Extra.extra
	}
	return softfloat_roundPackToF128(signZ, expZ, sigZ.High, sigZ.Low, sigZExtra)
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
propagateNaN:
	uiZ = softfloat_propagateNaNF128UI(uiA64, uiA0, uiB64, uiB0)
	goto uiZ
	/*------------------------------------------------------------------------
	*------------------------------------------------------------------------*/
infArg:
	if magBits == 0 {
		softfloat_raiseFlags(softfloat_flag_invalid)
		uiZ.High = defaultNaNF128UI64
		uiZ.Low = defaultNaNF128UI0
		goto uiZ
	}
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
	//uiZ.High = packToF128UI64( signZ, 0, 0 );
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

//#define wordIncr 1
//#define indexWord( total, n ) (n)
//#define indexWordHi( total ) ((total) - 1)
//#define indexWordLo( total ) 0
//#define indexMultiword( total, m, n ) (n)
//#define indexMultiwordHi( total, n ) ((total) - (n))
//#define indexMultiwordLo( total, n ) 0
//#define indexMultiwordHiBut( total, n ) (n)
//#define indexMultiwordLoBut( total, n ) 0
//#define INIT_UINTM4( v3, v2, v1, v0 ) { v0, v1, v2, v3 }

func indexWord(total, n uint8) uint8 {
	return n
}

func softfloat_normSubnormalF128Sig(sig64, sig0 uint64) exp32_sig128 {
	var shiftDist int8
	var z exp32_sig128

	if sig64 == 0 {
		shiftDist = int8(softfloat_countLeadingZeros64(sig0)) - 15
		z.exp = int32(-63 - shiftDist)
		if shiftDist < 0 {
			z.sig.High = sig0 >> uint8(-shiftDist)
			z.sig.Low = sig0 << uint8(shiftDist&63)
		} else {
			z.sig.High = sig0 << uint8(shiftDist)
			z.sig.Low = 0
		}
	} else {
		shiftDist = int8(softfloat_countLeadingZeros64(sig64)) - 15
		z.exp = int32(1 - shiftDist)
		z.sig = softfloat_shortShiftLeft128(sig64, sig0, uint8(shiftDist))
	}
	return z
}

func softfloat_mul128To256M(a64, a0, b64, b0 uint64, zPtr [4]uint64) [4]uint64 {
	var p0, p64, p128 Uint128
	var z64, z128, z192 uint64

	p0 = softfloat_mul64To128(a0, b0)
	zPtr[indexWord(4, 0)] = p0.Low
	p64 = softfloat_mul64To128(a64, b0)
	z64 = p64.Low + p0.High
	//z128 = p64.High + (z64 < p64.Low)
	if z64 < p64.Low {
		z128 = p64.High + 1
	} else {
		z128 = p64.High
	}
	p128 = softfloat_mul64To128(a64, b64)
	z128 += p128.Low
	//z192 = p128.High + (z128 < p128.Low)
	if z128 < p128.Low {
		z192 = p128.High + 1
	} else {
		z192 = p128.High
	}
	p64 = softfloat_mul64To128(a0, b64)
	z64 += p64.Low
	zPtr[indexWord(4, 1)] = z64
	//p64.High += (z64 < p64.Low);
	if z64 < p64.Low {
		p64.High += 1
	}
	z128 += p64.High
	zPtr[indexWord(4, 2)] = z128
	//zPtr[indexWord( 4, 3 )] = z192 + (z128 < p64.High)
	if z128 < p64.High {
		zPtr[indexWord(4, 3)] = z192 + 1
	} else {
		zPtr[indexWord(4, 3)] = z192
	}
	return zPtr
}

func softfloat_mul64To128(a, b uint64) Uint128 {
	var a32, a0, b32, b0 uint32
	var z Uint128
	var mid1, mid uint64

	a32 = uint32(a >> 32)
	a0 = uint32(a)
	b32 = uint32(b >> 32)
	b0 = uint32(b)
	z.Low = uint64(a0) * uint64(b0)
	mid1 = uint64(a32) * uint64(b0)
	mid = mid1 + uint64(a0)*uint64(b32)
	z.High = uint64(a32) * uint64(b32)

	if mid < mid1 {
		z.High += uint64(1)<<32 | mid>>32
	} else {
		z.High += uint64(0)<<32 | mid>>32
	}
	mid <<= 32
	z.Low += mid
	//z.High += (z.Low < mid)
	if z.Low < mid {
		z.High += 1
	}
	return z
}
