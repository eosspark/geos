package arithmeticTypes

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
)

/*----------------------------------------------------------------------------
| Types used to pass 16-bit, 32-bit, 64-bit, and 128-bit floating-point
| arguments and results to/from functions.  These types must be exactly
| 16 bits, 32 bits, 64 bits, and 128 bits in size, respectively.  Where a
| platform has "native" support for IEEE-Standard floating-point formats,
| the types below may, if desired, be defined as aliases for the native types
| (typically 'float' and 'double', and possibly 'long double').
*----------------------------------------------------------------------------*/

type Float16 uint16
type Float32 uint32
type Float64 uint64

type Float128 struct {
	Low  uint64
	High uint64
}

type ExtFloat80M struct {
	signExp uint16
	signIf  uint64
}

type ExtFloat80_t ExtFloat80M

func (f Float32) String() string {
	return fmt.Sprintf("%e", math.Float32frombits(uint32(f)))
}

func (f Float64) String() string {
	return fmt.Sprintf("%e", math.Float64frombits(uint64(f)))
}

func (f Float128) String() string {
	// Same for Int128, Float128
	number := make([]byte, 16)
	binary.LittleEndian.PutUint64(number[:], f.Low)
	binary.LittleEndian.PutUint64(number[8:], f.High)
	fmt.Println(f128_to_extF80(f))
	return fmt.Sprintf("0x%s%s", hex.EncodeToString(number[:8]), hex.EncodeToString(number[8:]))

}

func (f Float128) Bytes() []byte {
	return []byte{}
}

func (f *Float128) IsNan() bool {
	return (^f.High&uint64(0x7FFF000000000000)) == 0 && (f.Low != 0 || ((f.High & uint64(0x0000FFFFFFFFFFFF)) != 0))
}

//#define signF128UI64( a64 ) ((bool) ((uint64_t) (a64)>>63))
//#define expF128UI64( a64 ) ((int_fast32_t) ((a64)>>48) & 0x7FFF)
//#define fracF128UI64( a64 ) ((a64) & UINT64_C( 0x0000FFFFFFFFFFFF ))
//#define packToF128UI64( sign, exp, sig64 ) (((uint_fast64_t) (sign)<<63) + ((uint_fast64_t) (exp)<<48) + (sig64))

func signF128UI64(a64 uint64) bool {
	return a64>>63 != 0
}

func expF128UI64(a64 uint64) uint32 {
	return uint32((a64 >> 48) & 0x7FFF)

}
func fracF128UI64(a64 uint64) uint64 {
	return a64 & 0x0000FFFFFFFFFFFF
}

/*----------------------------------------------------------------------------
| The bit pattern for a default generated 128-bit floating-point NaN.
*----------------------------------------------------------------------------*/

const (
	defaultNaNF128UI64 = uint64(0xFFFF800000000000)
	defaultNaNF128UI0  = uint64(0)
)

type extF80M_extF80 struct {
	fM ExtFloat80M
	f  ExtFloat80_t
}

type uint128Extra struct {
	extra uint64
	v     Uint128
}

//type uint64Extra struct{
//   extra uint64
//   v uint64
//}

func softfloat_shiftRightJam128Extra(a64, a0, extra uint64, dist uint32) uint128Extra {
	var u8NegDist uint8
	var z uint128Extra
	u8NegDist = uint8(-dist)
	if dist < 64 {
		z.v.High = a64 >> dist
		z.v.Low = a64<<(u8NegDist&63) | a0>>dist
		z.extra = a0 << (u8NegDist & 63)
	} else {
		z.v.High = 0
		if dist == 64 {
			z.v.High = a64
			z.extra = a0
		} else {
			extra |= a0
			if dist < 128 {
				z.v.Low = a64 >> (dist & 63)
				z.extra = a64 << (u8NegDist & 63)
			} else {
				z.v.Low = 0
				if dist == 128 {
					z.extra = a64
				} else {
					if a64 != 0 {
						z.extra = 1
					} else {
						z.extra = 0
					}
				}
			}
		}

	}
	if extra != 0 {
		z.extra |= 1
	} else {
		z.extra |= 0
	}

	return z

}
func softfloat_shortShiftRightJam128Extra(a64, a0, extra uint64, dist uint8) uint128Extra {
	negDist := uint8(-dist)
	var z uint128Extra

	z.v.High = a64 >> dist
	z.v.Low = a64<<(negDist&63) | a0>>dist
	if extra != 0 {
		z.extra = a0<<(negDist&63) | 1
	} else {
		z.extra = a0 << (negDist & 63)
	}

	return z
}

type exp32_sig128 struct {
	exp int32
	sig Uint128
}
