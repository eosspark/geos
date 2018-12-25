package eos_math

func (a Float128) Sub(b Float128) Float128 {
	var uiA64, uiA0 uint64
	var signA bool

	var uiB64, uiB0 uint64
	var signB bool

	uiA64 = a.High
	uiA0 = a.Low
	signA = signF128UI64(uiA64)

	uiB64 = b.High
	uiB0 = b.Low
	signB = signF128UI64(uiB64)

	if signA == signB {
		return softfloat_subMagsF128(uiA64, uiA0, uiB64, uiB0, signA)
	} else {
		return softfloat_addMagsF128(uiA64, uiA0, uiB64, uiB0, signA)
	}
}
