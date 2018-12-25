package eos_math

func Fixunssfti(a uint32) Uint128 {
	// Break a into sign, exponent, significand
	aRep := a //uint32
	aAbs := aRep & 0x7FFFFFFF

	var sign bool
	if aRep&0x80000000 != 0 {
		sign = false
	} else {
		sign = true
	}

	exponent := aAbs>>23 - 0x7F
	significand := (aAbs & 0x7FFFFF) | 0x800000

	// If either the value or the exponent is negative, the result is zero.
	if exponent < 0 || sign == false {
		return CreateUint128(0)
	}

	// If the value is too large for the integer type, saturate.
	if uint64(exponent) > 128 {
		return MaxUint128()
	}

	// If 0 <= exponent < significandBits, right shift to get the result.
	// Otherwise, shift left.
	if exponent < 23 {
		return CreateUint128(int(significand >> (23 - exponent)))
	} else {
		re := Uint128{Low: uint64(significand)}
		re.LeftShifts(int(exponent - 23))
		return re
	}
}
