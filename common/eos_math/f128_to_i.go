package eos_math

/*----------------------------------------------------------------------------
| Software floating-point rounding mode.  (Mode "odd" is supported only if
| SoftFloat is compiled with macro 'SOFTFLOAT_ROUND_ODD' defined.)
*----------------------------------------------------------------------------*/
const (
	softfloat_round_near_even   = 0
	softfloat_round_minMag      = 1
	softfloat_round_min         = 2
	softfloat_round_max         = 3
	softfloat_round_near_maxMag = 4
	softfloat_round_odd         = 6
)

/*----------------------------------------------------------------------------
| Software floating-point exception flags.
*----------------------------------------------------------------------------*/
const (
	softfloat_flag_inexact   = 1
	softfloat_flag_underflow = 2
	softfloat_flag_overflow  = 4
	softfloat_flag_infinite  = 8
	softfloat_flag_invalid   = 16
)

/*----------------------------------------------------------------------------
| The values to return on conversions to 32-bit integer formats that raise an
| invalid exception.
*----------------------------------------------------------------------------*/
const (
	ui32_fromPosOverflow = 0xFFFFFFFF
	ui32_fromNegOverflow = 0xFFFFFFFF
	ui32_fromNaN         = 0xFFFFFFFF
	i32_fromPosOverflow  = (-0x7FFFFFFF - 1)
	i32_fromNegOverflow  = (-0x7FFFFFFF - 1)
	i32_fromNaN          = (-0x7FFFFFFF - 1)
)

/*----------------------------------------------------------------------------
| The values to return on conversions to 64-bit integer formats that raise an
| invalid exception.
*----------------------------------------------------------------------------*/
const (
	ui64_fromPosOverflow = uint64(0xFFFFFFFFFFFFFFFF)
	ui64_fromNegOverflow = uint64(0xFFFFFFFFFFFFFFFF)
	ui64_fromNaN         = uint64(0xFFFFFFFFFFFFFFFF)

	i64_fromPosOverflow = (-int64(0x7FFFFFFFFFFFFFFF) - 1)
	i64_fromNegOverflow = (-int64(0x7FFFFFFFFFFFFFFF) - 1)
	i64_fromNaN         = (-int64(0x7FFFFFFFFFFFFFFF) - 1)
)

func signF128Ui64(a64 uint64) bool {
	if (a64 >> 63) != 0 {
		return true
	} else {
		return false
	}
}

func expF128Ui64(a64 uint64) int32 {
	return int32(a64>>48) & 0x7FFF
}

func fracF128Ui64(a64 uint64) uint64 {
	return a64 & uint64(0x0000FFFFFFFFFFFF)
}

func isNaNF128UI(a64, a0 uint64) bool {
	return (^a64&uint64(0x7FFF000000000000)) == 0 && (a0 != 0 || (a64&uint64(0x0000FFFFFFFFFFFF)) != 0)
}

func F128ToI32(a Float128, roundingMode uint8, exact bool) int32 {
	var uA Float128
	var uiA64, uiA0 uint64
	var sign bool
	var exp int32
	var sig0, sig64 uint64
	var shiftDist int32

	uA = a
	uiA64 = uA.High
	uiA0 = uA.Low

	sign = signF128Ui64(uiA64)
	exp = expF128Ui64(uiA64)
	sig64 = fracF128Ui64(uiA64)
	sig0 = uiA0

	if exp != 0 {
		sig64 |= uint64(0x0001000000000000)
	}
	if sig0 != 0 {
		sig64 |= 1
	}
	shiftDist = 0x4023 - exp
	if 0 < shiftDist {
		sig64 = softfloat_shiftRightJam64(sig64, uint32(shiftDist))
	}
	return softfloat_roundToI32(sign, sig64, roundingMode, exact)
}

func softfloat_roundToI32(sign bool, sig uint64, roundingMode uint8, exact bool) int32 {
	var roundIncrement, roundBits uint16
	var sig32 uint32
	var uZ uint32
	var z int32
	var zminus bool
	var condition bool

	roundIncrement = 0x800
	if (roundingMode != softfloat_round_near_maxMag) && (roundingMode != softfloat_round_near_even) {
		roundIncrement = 0
		var signTrue bool
		if sign {
			signTrue = roundingMode == softfloat_round_min

		} else {
			signTrue = roundingMode == softfloat_round_max
		}
		if signTrue {
			roundIncrement = 0xFFF
		}
	}

	roundBits = uint16(sig & 0xFFF)
	sig += uint64(roundIncrement)
	if sig&uint64(0xFFFFF00000000000) != 0 {
		goto invalid
	}
	sig32 = uint32(sig >> 12)
	if (roundBits == 0x800) && (roundingMode == softfloat_round_near_even) {
		sig32 &= ^uint32(1)
	}

	if sign {
		uZ = -sig32
	} else {
		uZ = sig32
	}
	z = int32(uZ)

	zminus = z < 0

	if (zminus && sign) || (!zminus && !sign) {
		condition = false
	} else {
		condition = true
	}
	if z != 0 && condition {
		goto invalid
	}
	return z

invalid:
	//softfloat_raiseFlags( softfloat_flag_invalid )
	if sign {
		return i32_fromNegOverflow
	}
	return i32_fromPosOverflow
}

type uint64Extra struct {
	extra uint64
	v     uint64
}

func F128ToI64(a Float128, roundingMode uint8, exact bool) int64 {
	var uA Float128
	var uiA64, uiA0 uint64
	var sign bool
	var exp int32
	var sig64, sig0 uint64
	var shiftDist int32
	var sig128 Uint128
	var sigExtra uint64Extra

	uA = a
	uiA64 = uA.High
	uiA0 = uA.Low
	sign = signF128Ui64(uiA64)
	exp = expF128Ui64(uiA64)
	sig64 = fracF128Ui64(uiA64)
	sig0 = uiA0

	shiftDist = 0x402F - exp
	if shiftDist <= 0 {
		if shiftDist <= -15 {
			softfloat_raiseFlags(softfloat_flag_invalid)
			if (exp == 0x7FFF) && (sig64|sig0 != 0) {
				return i64_fromNaN
			} else {
				if sign {
					return i64_fromNegOverflow
				} else {
					return i64_fromPosOverflow
				}
			}
		}

		sig64 |= uint64(0x0001000000000000)
		if shiftDist != 0 {
			sig128 = softfloat_shortShiftLeft128(sig64, sig0, uint8(-shiftDist))
			sig64 = sig128.High
			sig0 = sig128.Low
		}
	} else {
		if exp != 0 {
			sig64 |= uint64(0x0001000000000000)
		}

		sigExtra = softfloat_shiftRightJam64Extra(sig64, sig0, uint32(shiftDist))
		sig64 = sigExtra.v
		sig0 = sigExtra.extra
	}
	return softfloat_roundToI64(sign, sig64, sig0, roundingMode, exact)
}

func softfloat_roundToI64(sign bool, sig uint64, sigExtra uint64, roundingMode uint8, exact bool) int64 {
	var z int64
	var uZ uint64
	var compare bool
	var condition bool
	var condition1 bool
	var zminus bool
	if (roundingMode == softfloat_round_near_maxMag) || (roundingMode == softfloat_round_near_even) {
		if uint64(0x8000000000000000) <= sigExtra {
			goto increment
		}
	} else {
		if sign {
			compare = roundingMode == softfloat_round_max
		} else {
			compare = roundingMode == softfloat_round_max
		}
		if sigExtra != 0 && compare {
			sig++
			if sig == 0 {
				goto invalid
			}
			if (sigExtra == uint64(0x8000000000000000)) && (roundingMode == softfloat_round_near_even) {
				sig &= ^uint64(1)
			}
		}
	}

	if sign {
		uZ = -sig
	} else {
		uZ = sig
	}
	z = int64(uZ)

	zminus = z < 0

	if (zminus && sign) || (!zminus && !sign) {
		condition = false
	} else {
		condition = true
	}
	if z != 0 && condition {
		goto invalid
	}

	//if sigExtra{
	//	if exact{
	//		softfloat_exceptionFlags |= softfloat_flag_inexact
	//	}
	//}
	return z

invalid:
	//softfloat_raiseFlags( softfloat_flag_invalid )
	if sign {
		return i64_fromNegOverflow
	}
	return i64_fromPosOverflow

increment:
	sig++
	if sig == 0 {
		goto invalid
	}
	if (sigExtra == uint64(0x8000000000000000)) && (roundingMode == softfloat_round_near_even) {
		sig &= ^uint64(1)
	}

	if sign {
		uZ = -sig
	} else {
		uZ = sig
	}
	z = int64(uZ)

	zminus = z < 0

	if (zminus && sign) || (!zminus && !sign) {
		condition1 = false
	} else {
		condition1 = true
	}
	if z != 0 && condition1 {
		goto invalid
	}
	return z
}

/*----------------------------------------------------------------------------
| Raises the exceptions specified by `flags'.  Floating-point traps can be
| defined here if desired.  It is currently not possible for such a trap
| to substitute a result value.  If traps are not implemented, this routine
| should be simply `softfloat_exceptionFlags |= flags;'.
*----------------------------------------------------------------------------*/
func softfloat_raiseFlags(flags uint8) {
	//softfloat_exceptionFlags |= flags
}

func F128ToUi32(a Float128, roundingMode uint8, exact bool) uint32 {
	var uA Float128
	var uiA64, uiA0 uint64
	var sign bool
	var exp int32
	var sig64 uint64
	var shiftDist int32

	uA = a
	uiA64 = uA.High
	uiA0 = uA.Low
	sign = signF128Ui64(uiA64)
	exp = expF128Ui64(uiA64)
	if uiA0 != 0 {
		sig64 = fracF128Ui64(uiA64) | 1
	} else {
		sig64 = fracF128Ui64(uiA64) //fracF128Ui64( uiA64 )|0
	}

	if exp != 0 {
		sig64 |= uint64(0x0001000000000000)
	}
	shiftDist = 0x4023 - exp
	if 0 < shiftDist {
		sig64 = softfloat_shiftRightJam64(sig64, uint32(shiftDist))
	}
	return softfloat_roundToUI32(sign, sig64, roundingMode, exact)

}

func softfloat_shiftRightJam64(a uint64, dist uint32) uint64 {
	if dist < 63 {
		if a<<(-dist&63) != 0 {
			return a>>dist | 1
		}
		return a >> dist
	} else {
		if a != 0 {
			return 1
		}
		return 0
	}
}

func softfloat_roundToUI32(sign bool, sig uint64, roundingMode uint8, exact bool) uint32 {
	var roundIncrement, roundBits uint16
	var z uint32

	roundIncrement = 0x800

	if roundingMode != softfloat_round_near_maxMag && roundingMode != softfloat_round_near_even {
		roundIncrement = 0
		if sign {
			if sig == 0 {
				return 0
			}
			if roundingMode == softfloat_round_min {
				goto invalid
			}
		} else {
			if roundingMode == softfloat_round_max {
				roundIncrement = 0xFFF
			}
		}
	}

	roundBits = uint16(sig & 0xFFF)
	sig += uint64(roundIncrement)
	if sig&uint64(0xFFFFF00000000000) != 0 {
		goto invalid
	}
	z = uint32(sig >> 12)

	if (roundBits == 0x800) && (roundingMode == softfloat_round_near_even) {
		//z &= ~(uint32) 1;
		z &= ^uint32(1)
	}

	if sign && z != 0 {
		goto invalid
	}
	//if roundBits != 0 {
	//	if exact {
	//		softfloat_exceptionFlags |= softfloat_flag_inexact
	//	}
	//}
	return z

invalid:
	//softfloat_raiseFlags(softfloat_flag_invalid)
	if sign {
		return ui32_fromNegOverflow
	} else {
		return ui32_fromPosOverflow
	}

}

func F128ToUi64(a Float128, roundingMode uint8, exact bool) uint64 {
	var uA Float128
	var uiA64, uiA0 uint64
	var sign bool
	var exp int32
	var sig64, sig0 uint64
	var shiftDist int32
	var sig128 Uint128
	var sigExtra uint64Extra

	uA = a
	uiA64 = uA.High
	uiA0 = uA.Low
	sign = signF128Ui64(uiA64)
	exp = expF128Ui64(uiA64)
	sig64 = fracF128Ui64(uiA64)
	sig0 = uiA0

	shiftDist = 0x402F - exp
	if shiftDist <= 0 {
		if shiftDist <= -15 {
			softfloat_raiseFlags(softfloat_flag_invalid)
			if (exp == 0x7FFF) && (sig64|sig0 != 0) {
				return ui64_fromNaN
			} else {
				if sign {
					return ui64_fromNegOverflow
				} else {
					return ui64_fromPosOverflow
				}
			}
		}

		sig64 |= uint64(0x0001000000000000)
		if shiftDist != 0 {
			sig128 = softfloat_shortShiftLeft128(sig64, sig0, uint8(-shiftDist))
			sig64 = sig128.High
			sig0 = sig128.Low
		}
	} else {
		if exp != 0 {
			sig64 |= uint64(0x0001000000000000)
		}

		sigExtra = softfloat_shiftRightJam64Extra(sig64, sig0, uint32(shiftDist))
		sig64 = sigExtra.v
		sig0 = sigExtra.extra
	}
	return softfloat_roundToUI64(sign, sig64, sig0, roundingMode, exact)
}

func softfloat_shortShiftLeft128(a64 uint64, a0 uint64, dist uint8) Uint128 {
	var z Uint128
	z.High = a64<<dist | a0>>(-dist&63)
	z.Low = a0 << dist
	return z

}

func softfloat_shiftRightJam64Extra(a, extra uint64, dist uint32) uint64Extra {
	var z uint64Extra
	if dist < 64 {
		z.v = a >> dist
		z.extra = a << (dist & 63)
	} else {
		z.v = 0
		if dist == 64 {
			z.extra = a
		} else {
			if a != 0 {
				z.extra = 1
			} else {
				z.extra = 0
			}
		}
	}
	//z.extra |= (extra != 0)
	if extra != 0 {
		z.extra |= 1
	}
	return z
}

func softfloat_roundToUI64(sign bool, sig, sigExtra uint64, roundingMode uint8, exact bool) uint64 {
	if roundingMode == softfloat_round_near_maxMag || roundingMode == softfloat_round_near_even {
		if uint64(0x8000000000000000) <= sigExtra {
			goto increment
		}
	} else {
		if sign {
			if sig|sigExtra == 0 {
				return 0
			}
			if roundingMode == softfloat_round_min {
				goto invalid
			}
		} else {
			if (roundingMode == softfloat_round_max) && sigExtra != 0 {
				sig++
				if sig == 0 {
					goto invalid
				}
				if (sigExtra == uint64(0x8000000000000000)) && (roundingMode == softfloat_round_near_even) {
					sig &= ^uint64(1)
				}
			}
		}
	}
	if sign && sig != 0 {
		goto invalid
	}
	return sig

invalid:
	if sign {
		return ui64_fromNegOverflow
	}
	return ui64_fromPosOverflow

increment:
	sig++
	if sig == 0 {
		goto invalid
	}
	if (sigExtra == uint64(0x8000000000000000)) && (roundingMode == softfloat_round_near_even) {
		sig &= ^uint64(1)
	}
	if sign && sig != 0 {
		goto invalid
	}
	return sig
}
