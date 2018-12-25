package eos_math

/* ===-- fixsfti.c - Implement __fixsfti -----------------------------------===
 *
 *                     The LLVM Compiler Infrastructure
 *
 * This file is dual licensed under the MIT and the University of Illinois Open
 * Source Licenses. See LICENSE.TXT for details.
 *
 * ===----------------------------------------------------------------------===
 */

//#define significandBits 23
//#define typeWidth       32
//#define exponentBits    8
//#define maxExponent     0xFF
//#define exponentBias    0x7F
//
//#define implicitBit     0x800000
//#define significandMask 0x7FFFFF
//#define signBit         0x80000000
//#define absMask         0x7FFFFFFF
//#define exponentMask    0x7F800000
//#define oneRep          0x3F800000
//#define infRep          0x7F800000
//#define quietBit        0x400000
//#define qnanRep         0x7FC00000

func Fixsfti(a uint32) Int128 {
	fixint_max := MaxInt128()
	fixint_min := MinInt128()
	// Break a into sign, exponent, significand
	aRep := a //uint32
	aAbs := aRep & 0x7FFFFFFF

	var sign Int128
	if aRep&0x80000000 != 0 {
		sign = CreateInt128(-1)
	} else {
		sign = CreateInt128(1)
	}

	exponent := aAbs>>23 - 0x7F
	significand := (aAbs & 0x7FFFFF) | 0x800000

	// If exponent is negative, the result is zero.
	if exponent < 0 {
		return CreateInt128(0)
	}

	// If the value is too large for the integer type, saturate.
	if exponent > 128 {
		if sign == CreateInt128(1) {
			return fixint_max
		} else {
			return fixint_min
		}
	}

	// If 0 <= exponent < significandBits, right shift to get the result.
	// Otherwise, shift left.
	if exponent < 23 {
		return sign.Mul(CreateInt128(int(significand >> (23 - exponent))))
	} else {
		re := Int128{Low: uint64(significand)}
		re.LeftShifts(int(exponent - 23))
		return sign.Mul(re)
	}
}
