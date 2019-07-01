package eos_math

import (
	"math"
)

/* Returns: convert a to a double, rounding toward even.*/

/* Assumption: double is a IEEE 64 bit floating point type
 *            ti_int is a 128 bit integral type
 */

/* seee eeee eeee mmmm mmmm mmmm mmmm mmmm | mmmm mmmm mmmm mmmm mmmm mmmm mmmm mmmm */
const DBL_MANT_DIG = 53

func Floattidf(a Int128) float64 {
	if a.IsZero() {
		return 0.0
	}
	var N = 128

	s := a.High >> 63
	a.High = a.High ^ uint64(0)
	a.Low = a.Low ^ s - s

	sd := N - clzti2(a)
	//sd :=7 //TODO

	e := sd - 1 /* exponent */
	if sd > DBL_MANT_DIG {
		/* start:  0000000000000000000001xxxxxxxxxxxxxxxxxxxxxxPQxxxxxxxxxxxxxxxxxx
		 *  finish: 000000000000000000000000000000000000001xxxxxxxxxxxxxxxxxxxxxxPQR
		 *                                               12345678901234567890123456
		 * 1 = msb 1 bit
		 * P = bit DBL_MANT_DIG-1 bits to the right of 1
		 * Q = bit DBL_MANT_DIG bits to the right of 1
		 * R = "or" of all bits to the right of Q
		 */
		switch sd {
		case DBL_MANT_DIG + 1:
			a.LeftShift()
		case DBL_MANT_DIG + 2:
		default:
			temp := a
			temp.RightShifts(sd - (DBL_MANT_DIG + 2))
			minusone := CreateInt128(-1)
			minusone.RightShifts((N + DBL_MANT_DIG + 2) - sd)
			if a.High&minusone.High != 0 && a.Low&minusone.Low != 0 {
				a.High = temp.High | CreateInt128(1).High
				a.Low = temp.High | CreateInt128(1).Low
			} else {
				a.High = temp.High | CreateInt128(0).High
				a.Low = temp.High | CreateInt128(0).Low
			}

		}

		/* finish: */
		if a.Low&4 != 0 { /* Or P into R */
			a.Low = a.Low | 1
		}
		a = a.Add(CreateInt128(1)) /* round - this step may add a significant bit */
		a.RightShifts(2)           /* dump Q and R */
		/* a is now rounded to DBL_MANT_DIG or DBL_MANT_DIG+1 bits */
		plusone := CreateInt128(1)
		plusone.LeftShifts(DBL_MANT_DIG)
		a.High = a.High & plusone.High
		a.Low = a.Low & plusone.Low
		if !a.IsZero() {
			a.RightShift()
			e += 1
		}
		/* a is now rounded to DBL_MANT_DIG bits */
	} else {
		a.RightShifts(DBL_MANT_DIG - sd)
		/* a is now rounded to DBL_MANT_DIG bits */
	}
	var hi uint32
	var lo uint32

	tempf := a
	tempf.RightShifts(32)
	hi = (uint32(s) & 0x80000000) | /* sign */
		uint32((e+1023)<<20) | /* exponent */
		(uint32(tempf.Low) & 0x000FFFFF) /* mantissa-high */
	lo = uint32(a.Low) /* mantissa-low */

	var f64 uint64
	f64 = uint64(hi)<<32 + uint64(lo)
	return math.Float64frombits(f64)

}

func clzti2(a Int128) int {
	var i int
	if a.High == 0 {
		for a.Low > 0 {
			a.Low = a.Low >> 1
			i += 1
		}
		return 128 - i
	} else {
		for a.High > 0 {
			a.High = a.High >> 1
			i += 1
		}
		return 64 - i
	}
}
