package eos_math

func Fixtfti(a Float128) Int128 {
	var ui64, ui0 uint64
	var sign bool
	var significand Int128
	ui64 = a.High
	ui0 = a.Low

	// Break a into sign, exponent, significand
	sign = signF128UI64(ui64)
	exponent := expF128UI64(ui64) - 0x3FFF
	significand.High = fracF128UI64(ui64) | uint64(1)<<48
	significand.Low = ui0

	// If exponent is negative, the result is zero.
	if exponent < 0 {
		return CreateInt128(0)
	}
	// If the value is too large for the integer type, saturate.
	if exponent >= 128 {
		if sign {
			return MaxInt128()
		} else {
			return MinInt128()
		}
	}

	// If 0 <= exponent < significandBits, right shift to get the result.
	// Otherwise, shift left.
	if exponent < 112 {
		significand.RightShifts(int(112 - exponent))
		if sign {
			significand.High = significand.High | uint64(1)<<63
		}
		return significand
	} else {
		significand.LeftShifts(int(exponent - 112))
		if sign {
			significand.High = significand.High | uint64(1)<<63
		}
		return significand
	}
}
