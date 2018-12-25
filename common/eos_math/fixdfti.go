package eos_math

/* ===-- fixdfti.c - Implement __fixdfti -----------------------------------===
 *
 *                     The LLVM Compiler Infrastructure
 *
 * This file is dual licensed under the MIT and the University of Illinois Open
 * Source Licenses. See LICENSE.TXT for details.
 *
 * ===----------------------------------------------------------------------===
 */
//#define significandBits 52
//#define typeWidth       64
//#define exponentBits    11
//#define maxExponent     0x7FF
//#define exponentBias    0x3FF
//
//#define implicitBit     0x10000000000000
//#define significandMask 0xFFFFFFFFFFFFF
//#define signBit         0x8000000000000000
//#define absMask         0x7FFFFFFFFFFFFFFF
//#define exponentMask    0x7FF0000000000000
//#define oneRep          0x3FF0000000000000
//#define infRep          0x7FF0000000000000
//#define quietBit        0x8000000000000
//#define qnanRep         0x7FF8000000000000

func Fixdfti(a uint64) Int128 {
	fixint_max := MaxInt128()
	fixint_min := MinInt128()

	// Break a into sign, exponent, significand
	aRep := a
	aAbs := aRep & 0x7FFFFFFFFFFFFFFF
	var sign Int128
	if aRep&0x8000000000000000 != 0 {
		sign = CreateInt128(-1)
	} else {
		sign = CreateInt128(1)
	}
	exponent := (aAbs >> 52) - 0x3FF
	significand := (aAbs & 0xFFFFFFFFFFFFF) | 0x10000000000000

	// If exponent is negative, the result is zero.
	if exponent < 0 {
		return CreateInt128(0)
	}

	// If the value is too large for the integer type, saturate.
	if exponent >= 128 {
		if sign == CreateInt128(1) {
			return fixint_max
		} else {
			return fixint_min
		}
	}

	//If 0 <= exponent < significandBits, right shift to get the result.
	//Otherwise, shift left.
	if exponent < 52 {
		return sign.Mul(Int128{Low: significand >> (52 - exponent)})
	} else {
		re := Int128{Low: significand}
		re.LeftShifts(int(exponent - 52))
		return sign.Mul(re)
	}

}
