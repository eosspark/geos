package arithmeticTypes

func (f Float128) F128Lt(b Float128) bool {
	var uiA64, uiA0 uint64
	var uiB64, uiB0 uint64
	var signA, signB bool

	uiA64 = f.High
	uiA0 = f.Low

	uiB64 = b.High
	uiB0 = b.Low

	if isNaNF128UI(uiA64, uiA0) || isNaNF128UI(uiB64, uiB0) {
		softfloat_raiseFlags(softfloat_flag_invalid)
		return false
	}
	signA = signF128Ui64(uiA64)
	signB = signF128Ui64(uiB64)
	if signA != signB {
		return signA && (uiA0|uiA64|((uiA64|uiB64)&uint64(0x7FFFFFFFFFFFFFFF))) != 0
	} else {
		re := softfloat_lt128(uiA64, uiA0, uiB64, uiB0)
		if (re && signA) || (!re && !signA) {
			return false
		} else {
			return uiA64 != uiB64 || uiA0 != uiB0
		}
	}
}

func (f Float128) F128EQ(b Float128) bool {
	var uiA64, uiA0 uint64
	var uiB64, uiB0 uint64

	uiA64 = f.High
	uiA0 = f.Low

	uiB64 = b.High
	uiB0 = b.Low
	if isNaNF128UI(uiA64, uiA0) || isNaNF128UI(uiB64, uiB0) {
		if softfloat_isSigNaNF128UI(uiA64, uiA0) || softfloat_isSigNaNF128UI(uiB64, uiB0) {
			softfloat_raiseFlags(softfloat_flag_invalid)
		}
		return false
	}
	return uiA0 == uiB0 && (uiA64 == uiB64 || (uiA0 == 0 && ((uiA64|uiB64)&uint64(0x7FFFFFFFFFFFFFFF)) == 0))
}
