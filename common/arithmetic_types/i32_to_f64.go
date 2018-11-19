package arithmeticTypes

func I32ToF64(a int32) Float64 {
	var uiZ uint64
	var sign bool
	var absA uint32
	var shiftDist uint8

	if a == 0 {
		uiZ = 0
	} else {
		sign = a < 0
		if sign {
			absA = uint32(^a + 1)
		} else {
			absA = uint32(a)
		}
		shiftDist = softfloat_countLeadingZeros32(absA) + 21
		uiZ = packToF64UI(sign, 0x432-uint64(shiftDist), uint64(absA)<<shiftDist)
	}
	return Float64(uiZ)
}
