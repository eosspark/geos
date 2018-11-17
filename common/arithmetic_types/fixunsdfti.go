package arithmeticTypes

func Fixunsdfti(a uint64) Uint128 {
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
	if exponent < 0 || sign == CreateInt128(-1) {
		return CreateUint128(0)
	}

	// If the value is too large for the integer type, saturate.
	if uint64(exponent) >= 128 {
		return MaxUint128()
	}

	//If 0 <= exponent < significandBits, right shift to get the result.
	//Otherwise, shift left.
	if exponent < 52 {
		return Uint128{Low: significand >> (52 - exponent)}
	} else {
		re := Uint128{Low: significand}
		re.LeftShifts(int(exponent - 52))
		return re
	}

}
