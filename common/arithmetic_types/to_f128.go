package arithmeticTypes

func packToF128UI64(sign, exp, sig64 uint64) uint64 {
	return sign<<63 + exp<<48 + sig64
}

func Ui64ToF128(a uint64) Float128 {
	var uiZ64, uiZ0 uint64
	var shiftDist uint8
	var zSig Uint128
	var uZ Float128

	if a == 0 {
		uiZ64 = 0
		uiZ0 = 0
	} else {
		shiftDist = softfloatCountLeadingZeros64(a) + 49
		if 64 <= shiftDist {
			zSig.High = a << (shiftDist - 64)
			zSig.Low = 0

		} else {
			zSig = softfloatShortShiftLeft128(0, a, shiftDist)
		}
		uiZ64 = packToF128UI64(0, uint64(0x406E)-uint64(shiftDist), zSig.High)
		uiZ0 = zSig.Low
	}

	uZ.Low = uiZ0
	uZ.High = uiZ64
	return uZ
}

func Ui32ToF128(a uint32) Float128 {
	var uiZ64 uint64
	var shiftDist uint8
	var uZ Float128

	uiZ64 = 0
	if a != 0 {
		shiftDist = softfloatCountLeadingZeros32(a) + 17
		uiZ64 = packToF128UI64(0, uint64(0x402E)-uint64(shiftDist), uint64(a)<<shiftDist) //TODO
	}
	uZ.Low = 0
	uZ.High = uiZ64
	return uZ
}

func I64ToF128(a int64) Float128 {
	var uiZ64, uiZ0 uint64
	var sign uint64
	var absA uint64

	var shiftDist uint8
	var zSig Uint128
	var uZ Float128

	if a == 0 {
		uiZ64 = 0
		uiZ0 = 0
	} else {
		if a < 0 {
			sign = 1 //true
			absA = -uint64(a)
		} else {
			sign = 0 //false
			absA = uint64(a)
		}
		shiftDist = softfloatCountLeadingZeros64(absA) + 49
		if 64 <= shiftDist {
			zSig.High = absA << (shiftDist - 64)
			zSig.Low = 0
		} else {
			zSig = softfloatShortShiftLeft128(0, absA, shiftDist)
		}
		uiZ64 = packToF128UI64(sign, uint64(0x406E)-uint64(shiftDist), zSig.High)
		uiZ0 = zSig.Low
	}
	uZ.Low = uiZ0
	uZ.High = uiZ64
	return uZ
}

func I32ToF128(a int32) Float128 {
	var uiZ64 uint64
	var sign uint64
	var absA uint32
	var shiftDist uint8
	var uZ Float128

	uiZ64 = 0
	if a != 0 {
		if a < 0 {
			sign = 1
			absA = -uint32(a)
		} else {
			sign = 0
			absA = uint32(a)
		}
		shiftDist = softfloatCountLeadingZeros32(absA) + 17
		uiZ64 = packToF128UI64(sign, uint64(0x402E)-uint64(shiftDist), uint64(absA)<<shiftDist)
	}
	uZ.Low = 0
	uZ.High = uiZ64
	return uZ
}
