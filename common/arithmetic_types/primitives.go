package arithmeticTypes

var softfloatCountLeadingZeros8 = [256]uint8{
	8, 7, 6, 6, 5, 5, 5, 5, 4, 4, 4, 4, 4, 4, 4, 4,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

/*----------------------------------------------------------------------------
| Returns the number of leading 0 bits before the most-significant 1 bit of
| 'a'.  If 'a' is zero, 32 is returned.
*----------------------------------------------------------------------------*/
//#if defined INLINE_LEVEL && (3 <= INLINE_LEVEL)
//INLINE uint_fast8_t softfloat_countLeadingZeros32( uint32_t a )
//{
//    uint_fast8_t count = 0;
//    if ( a < 0x10000 ) {
//        count = 16;
//        a <<= 16;
//    }
//    if ( a < 0x1000000 ) {
//        count += 8;
//        a <<= 8;
//    }
//    count += softfloat_countLeadingZeros8[a>>24];
//    return count;
//}

func softfloatCountLeadingZeros32(a uint32) uint8 {
	count := uint8(0)
	if a < 0x10000 {
		count = 16
		a <<= 16
	}
	if a < 0x1000000 {
		count += 8
		a <<= 8
	}
	count += softfloatCountLeadingZeros8[a>>24]
	return count
}

//uint_fast8_t softfloat_countLeadingZeros64( uint64_t a )
//{
//uint_fast8_t count;
//uint32_t a32;
//
//count = 0;
//a32 = a>>32;
//if ( ! a32 ) {
//count = 32;
//a32 = a;
//}
///*------------------------------------------------------------------------
//| From here, result is current count + count leading zeros of `a32'.
//*------------------------------------------------------------------------*/
//if ( a32 < 0x10000 ) {
//count += 16;
//a32 <<= 16;
//}
//if ( a32 < 0x1000000 ) {
//count += 8;
//a32 <<= 8;
//}
//count += softfloat_countLeadingZeros8[a32>>24];
//return count;
//
//}

func softfloatCountLeadingZeros64(a uint64) uint8 {
	var a32 uint32
	count := uint8(0)
	a32 = uint32(a >> 32)
	if a32 == 0 {
		count = 32
		a32 = uint32(a)
	}

	/*------------------------------------------------------------------------
	  | From here, result is current count + count leading zeros of `a32'.
	  *------------------------------------------------------------------------*/
	if a32 < 0x10000 {
		count += 16
		a32 <<= 16
	}
	if a32 < 0x1000000 {
		count += 8
		a32 <<= 8
	}
	count += softfloatCountLeadingZeros8[a32>>24]
	return count
}

///*----------------------------------------------------------------------------
//| Shifts the 128 bits formed by concatenating 'a64' and 'a0' left by the
//| number of bits given in 'dist', which must be in the range 1 to 63.
//*----------------------------------------------------------------------------*/
//#if defined INLINE_LEVEL && (2 <= INLINE_LEVEL)
//INLINE
//struct uint128
//softfloat_shortShiftLeft128( uint64_t a64, uint64_t a0, uint_fast8_t dist )
//{
//struct uint128 z;
//z.v64 = a64<<dist | a0>>(-dist & 63);
//z.v0 = a0<<dist;
//return z;
//}

/*----------------------------------------------------------------------------
| Shifts the 128 bits formed by concatenating 'a64' and 'a0' left by the
| number of bits given in 'dist', which must be in the range 1 to 63.
*----------------------------------------------------------------------------*/
func softfloatShortShiftLeft128(a64, a0 uint64, dist uint8) Uint128 {
	var z Uint128
	z.High = a64<<dist | a0>>(-dist&63)
	z.Low = a0 << dist
	return z
}
