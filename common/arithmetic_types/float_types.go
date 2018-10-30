package arithmeticTypes

//
////
/////*----------------------------------------------------------------------------
////| Types used to pass 16-bit, 32-bit, 64-bit, and 128-bit floating-point
////| arguments and results to/from functions.  These types must be exactly
////| 16 bits, 32 bits, 64 bits, and 128 bits in size, respectively.  Where a
////| platform has "native" support for IEEE-Standard floating-point formats,
////| the types below may, if desired, be defined as aliases for the native types
////| (typically 'float' and 'double', and possibly 'long double').
////*----------------------------------------------------------------------------*/
//
type Float16 uint16
type Float32 uint32
type Float64 uint64
type Float128 struct {
	Low uint64
	High uint64
}

type ExtFloat80M struct {
	signExp uint16
	signIf  uint64
}

type ExtFloat80M_t ExtFloat80M

func (f Float128) Add(b Float128) Float128 {

	return Float128{}
}
func (f Float128) Sub(b Float128) Float128 {
	return Float128{}
}

func (f Float128) Mul(b Float128) Float128 {
	return Float128{}
}

func (f Float128) Div(b Float128) Float128 {
	return Float128{}
}
func (f Float128) String() string {
	return ""
}
func (f Float128) Bytes() []byte {
	return []byte{}
}

func (f *Float128) IsNan() bool {
	return (^f.High&uint64(0x7FFF000000000000)) == 0 && (f.Low != 0 || ((f.High & uint64(0x0000FFFFFFFFFFFF)) != 0))
}

//func (a Float128) Add(b Float128) Float128 {
//var uA ui128_f128
//var uiA64,uiA0 uint64
//var signA bool
//
//var uB ui128_f128
//var uiB64,uiB0 uint64
//var signB bool
//
//uA.f = a
//uiA64 = uA.ui.High
//uiA0 = uA.ui.Low
//signA = signF128UI64(uiA64)
//
//uB.f=b
//uiB64 = uB.ui.High
//uiB0 = uB.ui.Low
//signB = signF128UI64(uiB64)
//if signA ==signB{
//	return softFloatAddMagsF128(uiA64,uiA0,uiB64,uiB0,signA)
//}else{
//	return softFloatSubMagsF128(uiA64,uiA0,uiB64,uiB0,signA)
//}
//	return a
//}

//
//
//
////#define signF128UI64( a64 ) ((bool) ((uint64_t) (a64)>>63))
////#define expF128UI64( a64 ) ((int_fast32_t) ((a64)>>48) & 0x7FFF)
////#define fracF128UI64( a64 ) ((a64) & UINT64_C( 0x0000FFFFFFFFFFFF ))
////#define packToF128UI64( sign, exp, sig64 ) (((uint_fast64_t) (sign)<<63) + ((uint_fast64_t) (exp)<<48) + (sig64))
//
//func signF128UI64(a64 uint64)bool{
//    return a64>>63 !=0
//}
//
//func expF128UI64(a64 uint64) uint32{
//   return uint32((a64>>48) &0x7FFF)
//
//}
//func fracF128UI64(a64 uint64) uint64{
//    return a64&0x0000FFFFFFFFFFFF
//}
//func packToF128UI64(sign bool,exp,sig64 uint64) uint64{
//    if sign{
//        return uint64(1 <<63)+uint64(exp<<48)+sig64
//    }else{
//        return uint64(exp<<48)+sig64
//    }
////return uint64(uint8(sign) <<63)+uint64(exp<<48)+sig64
//}
//
//func softFloatAddMagsF128(uiA64,uiA0,uiB64,uiB0 uint64,signZ bool)Float128_t {
//    var expA uint32
//    var sigA Uint128
//    var expB uint32
//    var sigB Uint128
//    var expDiff uint32
//
//    var uiZ, sigZ Uint128
//    var expZ uint32
//    var sigZExtra uint64
//    var sig128Extra uint128Extra
//    var uZ ui128_f128
//
//    expA = expF128UI64(uiA64)
//    sigA.High = fracF128UI64(uiA64)
//    sigA.Low = uiA0
//    expB = expF128UI64(uiB64)
//    sigB.High = fracF128UI64(uiB64)
//    sigB.Low = uiB0
//    expDiff = expA - expB
//    if expDiff == 0 {
//        if expA == 0x7FFF {
//            if sigA.High | sigA.Low | sigB.High | sigB.Low {
//                goto propagateNaN
//            }
//            uiZ.High = uiA64
//            uiZ.Low = uiA0
//            goto uiZ
//        }
//        sigZ = sigA.Add(sigB)
//        if expA == 0 {
//            uiZ.High = packToF128UI64(signZ, 0, sigZ.High)
//            uiZ.Low = sigZ.Low
//            goto uiZ
//
//        }
//
//        expZ = expA
//        sigZ.High |= 0x0002000000000000
//        sigZExtra = 0
//        goto shiftRight1
//
//    }
//
//    if expDiff < 0 {
//        if expB == 0x7FFF {
//            if sigB.High|sigB.Low != 0 {
//                goto propagateNaN
//            }
//            uiZ.High = packToF128UI64(signZ, 0x7FFF, 0)
//            uiZ.Low = 0
//            goto uiZ
//        }
//        expZ = expB
//        if expA != 0 {
//            sigA.High |= 0x0001000000000000
//        } else {
//            expDiff++
//            sigZExtra = 0
//            if expDiff == 0 {
//                goto newlyAligned
//            }
//            sig128Extra = softfloat_shiftRightJam128Extra(sigA.High, sigA.Low, 0, -expDiff)
//            sigA = sig128Extra.v
//            sigZExtra = sig128Extra.extra
//        }
//    } else {
//        if  expA == 0x7FFF  {
//            if ( sigA.High | sigA.Low ) goto propagateNaN
//            uiZ.High = uiA64
//            uiZ.Low  = uiA0
//            goto uiZ
//        }
//        expZ = expA
//        if expB !=0  {
//            sigB.High |=  0x0001000000000000
//        } else {
//            expDiff--
//            sigZExtra = 0
//            if expDiff==0 {
//                goto newlyAligned
//            }
//        }
//        sig128Extra =
//            softfloat_shiftRightJam128Extra( sigB.High, sigB.Low, 0, expDiff )
//        sigB = sig128Extra.v
//        sigZExtra = sig128Extra.extra
//    }
//newlyAligned:
//    sigA.High |= 0x0001000000000000
//    sigZ =sigA.Add(sigB)
//
//    expZ--
//    if  sigZ.High < 0x0002000000000000 {
//        goto roundAndPack
//    }
//    expZ++
//shiftRight1:
//    sig128Extra =
//        softfloat_shortShiftRightJam128Extra(
//            sigZ.High, sigZ.Low, sigZExtra, 1 );
//    sigZ = sig128Extra.v
//    sigZExtra = sig128Extra.extra
//roundAndPack:
//    return softfloat_roundPackToF128( signZ, expZ, sigZ.High, sigZ.Low, sigZExtra );
//propagateNaN:
//    uiZ = softfloat_propagateNaNF128UI( uiA64, uiA0, uiB64, uiB0 );
//uiZ:
//    uZ.ui = uiZ
//    return uZ.f
//}
//
//
//
//func softFloatSubMagsF128(uiA64,uiA0,uiB64,uiB0 uint64,signZ bool)Float128_t{
//
//}
//
//
//
//type ui16_f16 struct{
//  ui uint16
//  f Float16_t
//}
//type ui32_f32 struct{
//  ui uint32
//  f Float32_t
//}
//type ui64_f64 struct{
//  ui uint64
//  f Float64_t
//}
//
//
//type extF80M_extF80 struct{
//  fM ExtFloat80M
//  f ExtFloat80M_t
//}
//
//
//type ui128_f128 struct{
//  ui Uint128
//  f Float128_t
//}
//
//
//
//
//
//type uint128Extra struct{
//    extra uint64
//    v Uint128
//}
//type uint64Extra struct{
//    extra uint64
//    v uint64
//}
//
//
//
//
//func softfloat_shiftRightJam128Extra(a64,a0,extra uint64,dist uint32)uint128Extra{
//    var u8NegDist uint8
//    var z uint128Extra
//    u8NegDist = -dist
//    if dist <64{
//        z.v.High = a64 >>dist
//        z.v.Low = a64 <<(u8NegDist &63) |a0 >>dist
//        z.extra= a0 <<(u8NegDist&63)
//    }else{
//        z.v.High =0
//        if dist ==64{
//            z.v.High =a64
//            z.extra = a0
//        }else{
//            extra |=a0
//            if dist <128{
//                z.v.Low = a64 >>(dist&63)
//                z.extra = a64 <<(u8NegDist &63)
//            }else{
//                z.v.Low = 0
//                if dist ==128{
//                    z.extra = a64
//                }else{
//                    if a64 !=0{
//                        z.extra = 1
//                    }else{
//                        z.extra = 0
//                    }
//                }
//            }
//        }
//
//    }
//    if extra!=0{
//        z.extra |=1
//    }else{
//        z.extra |=0
//    }
//
//return z
//
//}
//func softfloat_shortShiftRightJam128Extra(a64,  a0,  extra uint64,  dist uint8)uint128Extra {
// negDist := uint8(-dist)
// var z uint128Extra
//
//z.v.High = a64>>dist
//z.v.Low = a64<<(negDist & 63) | a0>>dist
//if extra !=0{
//    z.extra = a0<<(negDist & 63) | 1
//}else{
//    z.extra = a0<<(negDist & 63)
//}
//
//return z
//}
