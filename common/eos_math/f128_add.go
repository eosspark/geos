package eos_math

func (a Float128) Add(b Float128) Float128 {
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
		return softfloat_addMagsF128(uiA64, uiA0, uiB64, uiB0, signA)
	} else {
		return softfloat_subMagsF128(uiA64, uiA0, uiB64, uiB0, signA)
	}
	return a
}

func softfloat_addMagsF128(uiA64, uiA0, uiB64, uiB0 uint64, signZ bool) Float128 {
	var expA uint32
	var sigA Uint128
	var expB uint32
	var sigB Uint128
	var expDiff uint32

	var uiZ, sigZ Uint128
	var expZ uint32
	var sigZExtra uint64
	var sig128Extra uint128Extra
	var uZ Float128

	expA = expF128UI64(uiA64)
	sigA.High = fracF128UI64(uiA64)
	sigA.Low = uiA0
	expB = expF128UI64(uiB64)
	sigB.High = fracF128UI64(uiB64)
	sigB.Low = uiB0
	expDiff = expA - expB
	if expDiff == 0 {
		if expA == 0x7FFF {
			if sigA.High|sigA.Low|sigB.High|sigB.Low != 0 {
				goto propagateNaN
			}
			uiZ.High = uiA64
			uiZ.Low = uiA0
			goto uiZ
		}
		sigZ = sigA.Add(sigB)
		if expA == 0 {

			if signZ {
				uiZ.High = packToF128UI64(1, 0, sigZ.High)
			} else {
				uiZ.High = packToF128UI64(0, 0, sigZ.High)
			}
			uiZ.Low = sigZ.Low
			goto uiZ

		}

		expZ = expA
		sigZ.High |= 0x0002000000000000
		sigZExtra = 0
		goto shiftRight1

	}

	if expDiff < 0 {
		if expB == 0x7FFF {
			if sigB.High|sigB.Low != 0 {
				goto propagateNaN
			}
			if signZ {
				uiZ.High = packToF128UI64(1, 0x7FFF, 0)
			} else {
				uiZ.High = packToF128UI64(0, 0x7FFF, 0)
			}

			uiZ.Low = 0
			goto uiZ
		}
		expZ = expB
		if expA != 0 {
			sigA.High |= 0x0001000000000000
		} else {
			expDiff++
			sigZExtra = 0
			if expDiff == 0 {
				goto newlyAligned
			}
			sig128Extra = softfloat_shiftRightJam128Extra(sigA.High, sigA.Low, 0, -expDiff)
			sigA = sig128Extra.v
			sigZExtra = sig128Extra.extra
		}
	} else {
		if expA == 0x7FFF {
			if sigA.High|sigA.Low != 0 {
				goto propagateNaN
			}
			uiZ.High = uiA64
			uiZ.Low = uiA0
			goto uiZ
		}
		expZ = expA
		if expB != 0 {
			sigB.High |= 0x0001000000000000
		} else {
			expDiff--
			sigZExtra = 0
			if expDiff == 0 {
				goto newlyAligned
			}
		}
		sig128Extra =
			softfloat_shiftRightJam128Extra(sigB.High, sigB.Low, 0, expDiff)
		sigB = sig128Extra.v
		sigZExtra = sig128Extra.extra
	}
newlyAligned:
	sigA.High |= 0x0001000000000000
	sigZ = sigA.Add(sigB)

	expZ--
	if sigZ.High < 0x0002000000000000 {
		goto roundAndPack
	}
	expZ++
shiftRight1:
	sig128Extra =
		softfloat_shortShiftRightJam128Extra(
			sigZ.High, sigZ.Low, sigZExtra, 1)
	sigZ = sig128Extra.v
	sigZExtra = sig128Extra.extra
roundAndPack:
	return softfloat_roundPackToF128(signZ, int32(expZ), sigZ.High, sigZ.Low, sigZExtra)
propagateNaN:
	uiZ = softfloat_propagateNaNF128UI(uiA64, uiA0, uiB64, uiB0)
uiZ:
	uZ.High = uiZ.High
	uZ.Low = uiZ.Low

	return uZ
}

func softfloat_subMagsF128(uiA64, uiA0, uiB64, uiB0 uint64, signZ bool) Float128 {
	var (
		expA          int32
		sigA          Uint128
		expB          int32
		sigB, sigZ    Uint128
		expDiff, expZ int32
		uiZ           Uint128
		uZ            Float128
	)

	expA = int32(expF128UI64(uiA64))
	sigA.High = fracF128UI64(uiA64)
	sigA.Low = uiA0
	expB = int32(expF128UI64(uiB64))
	sigB.High = fracF128UI64(uiB64)
	sigB.Low = uiB0
	sigA = softfloat_shortShiftLeft128(sigA.High, sigA.Low, 4)
	sigB = softfloat_shortShiftLeft128(sigB.High, sigB.Low, 4)
	expDiff = expA - expB
	if 0 < expDiff {
		goto expABigger
	}
	if expDiff < 0 {
		goto expBBigger
	}

	if expA == 0x7FFF {
		if sigA.High|sigA.Low|sigB.High|sigB.Low != 0 {
			goto propagateNaN
		}
		softfloat_raiseFlags(softfloat_flag_invalid)
		uiZ.High = defaultNaNF128UI64
		uiZ.Low = defaultNaNF128UI0
		goto uiZ
	}
	expZ = expA
	if expZ == 0 {
		expZ = 1
	}
	if sigB.High < sigA.High {
		goto aBigger
	}
	if sigA.High < sigB.High {
		goto bBigger
	}
	if sigB.Low < sigA.Low {
		goto aBigger
	}
	if sigA.Low < sigB.Low {
		goto bBigger
	}
	if softfloat_roundingMode == softfloat_round_min {
		uiZ.High =
			packToF128UI64(1, 0, 0)
	} else {
		uiZ.High =
			packToF128UI64(0, 0, 0)
	}
	uiZ.Low = 0
	goto uiZ

expBBigger:
	if expB == 0x7FFF {
		if sigB.High|sigB.Low != 0 {
			goto propagateNaN
		}
		if signZ {
			uiZ.High = packToF128UI64(0, 0x7FFF, 0)
		} else {
			uiZ.High = packToF128UI64(1, 0x7FFF, 0)
		}
		uiZ.Low = 0
		goto uiZ
	}
	if expA != 0 {
		sigA.High |= uint64(0x0010000000000000)
	} else {
		expDiff += 1
		if expDiff == 0 {
			goto newlyAlignedBBigger
		}
	}
	sigA = softfloat_shiftRightJam128(sigA.High, sigA.Low, uint32(-expDiff))
newlyAlignedBBigger:
	expZ = expB
	sigB.High |= uint64(0x0010000000000000)
bBigger:
	signZ = !signZ
	sigZ = softfloat_sub128(sigB.High, sigB.Low, sigA.High, sigA.Low)
	goto normRoundPack
expABigger:
	if expA == 0x7FFF {
		if sigA.High|sigA.Low != 0 {
			goto propagateNaN
		}
		uiZ.High = uiA64
		uiZ.Low = uiA0
		goto uiZ
	}
	if expB != 0 {
		sigB.High |= uint64(0x0010000000000000)
	} else {
		expDiff -= 1
		if expDiff == 0 {
			goto newlyAlignedABigger
		}

	}
	sigB = softfloat_shiftRightJam128(sigB.High, sigB.Low, uint32(expDiff))
newlyAlignedABigger:
	expZ = expA
	sigA.High |= uint64(0x0010000000000000)
aBigger:
	sigZ = softfloat_sub128(sigA.High, sigA.Low, sigB.High, sigB.Low)
normRoundPack:
	return softfloat_normRoundPackToF128(signZ, expZ-5, sigZ.High, sigZ.Low)
propagateNaN:
	uiZ = softfloat_propagateNaNF128UI(uiA64, uiA0, uiB64, uiB0)
uiZ:
	uZ.High = uiZ.High
	uZ.Low = uiZ.Low

	return uZ
}
