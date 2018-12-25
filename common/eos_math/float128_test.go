package eos_math

import (
	"encoding/binary"
	"fmt"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/stretchr/testify/assert"
	"math"
	"math/big"
	"strconv"
	"testing"
	"unsafe"
)

func TestUiToFloat128(t *testing.T) {
	check := Float128{High: 4613251722985340928, Low: 0}
	a := Ui64ToF128(100)
	assert.Equal(t, check, a)

	b := Ui32ToF128(100)
	assert.Equal(t, check, b)

}

func TestItoFloat128(t *testing.T) {
	check := Float128{High: 13836623759840116736, Low: 0}
	a := I64ToF128(-100)
	assert.Equal(t, check, a)

	b := I32ToF128(-100)
	assert.Equal(t, check, b)
}

func TestF128ToUi32(t *testing.T) {
	f128 := Float128{High: 4613251722985340928, Low: 0}
	result := F128ToUi32(f128, 0, false)
	assert.Equal(t, uint32(100), result)
}

func TestF128ToUi64(t *testing.T) {
	f128 := Float128{High: 4613251722985340928, Low: 0}
	result := F128ToUi64(f128, 0, false)
	assert.Equal(t, uint64(100), result)
}

func TestF128ToI64(t *testing.T) {
	plusF128 := Float128{High: 4613251722985340928, Low: 0}
	re1 := F128ToI64(plusF128, 0, false)
	assert.Equal(t, int64(100), re1)

	minusF128 := Float128{High: 13836623759840116736, Low: 0}
	re2 := F128ToI64(minusF128, 0, false)
	assert.Equal(t, int64(-100), re2)
}

func TestF128ToI32(t *testing.T) {
	plusF128 := Float128{High: 4613251722985340928, Low: 0}
	re1 := F128ToI32(plusF128, 0, false)
	assert.Equal(t, int32(100), re1)

	minusF128 := Float128{High: 13836623759840116736, Low: 0}
	re2 := F128ToI32(minusF128, 0, false)
	assert.Equal(t, int32(-100), re2)
}

func TestF128ToF32(t *testing.T) {
	plusF128 := Float128{High: 4613251722985340928, Low: 0}
	a := F128ToF32(plusF128)
	assert.Equal(t, 1120403456, int(a))

	minusF128 := Float128{High: 13836623759840116736, Low: 0}
	b := F128ToF32(minusF128)
	assert.Equal(t, Float32(3267887104), b)

	int60 := int64(1) << 60
	f12860 := I64ToF128(int60)
	c := F128ToF32(f12860)
	assert.Equal(t, Float32(1568669696), c)
}

func TestF128ToF64(t *testing.T) {
	plusF128 := Float128{High: 4613251722985340928, Low: 0} //100
	a := F128ToF64(plusF128)
	assert.Equal(t, Float64(4636737291354636288), a)

	minusF128 := Float128{High: 13836623759840116736, Low: 0}
	b := F128ToF64(minusF128)
	assert.Equal(t, Float64(13860109328209412096), b)

	int60 := int64(1) << 60
	f12860 := I64ToF128(int60)
	c := F128ToF64(f12860)
	assert.Equal(t, Float64(4877398396442247168), c)

	test := Float128{High: 4629393042053316608, Low: 4629393042053316608}
	d := F128ToF64(test)
	assert.Equal(t, Float64(4894998396442247172), d)

}

func TestF32ToF128(t *testing.T) {
	plusF128 := Float128{High: 4613251722985340928, Low: 0}
	a := F128ToF32(plusF128)
	assert.Equal(t, Float32(1120403456), a)
	f128 := F32ToF128(a)
	assert.Equal(t, plusF128, f128)

	minusF128 := Float128{High: 13836623759840116736, Low: 0}
	b := F128ToF32(minusF128)
	assert.Equal(t, Float32(3267887104), b)
	f128minus := F32ToF128(b)
	assert.Equal(t, minusF128, f128minus)

	int60 := int64(1) << 60
	f12860 := I64ToF128(int60)
	c := F128ToF32(f12860)
	f60 := F32ToF128(c)
	assert.Equal(t, f12860, f60)
}

func TestF64ToF128(t *testing.T) {
	plusF128 := Float128{High: 4613251722985340928, Low: 0} //100
	a := F128ToF64(plusF128)
	assert.Equal(t, Float64(4636737291354636288), a)

	f128plus := F64ToF128(a)
	assert.Equal(t, plusF128, f128plus)

	minusF128 := Float128{High: 13836623759840116736, Low: 0}
	b := F128ToF64(minusF128)
	assert.Equal(t, Float64(13860109328209412096), b)
	f128minus := F64ToF128(b)
	assert.Equal(t, minusF128, f128minus)

	int60 := int64(1) << 60
	f12860 := I64ToF128(int60)
	c := F128ToF64(f12860)
	assert.Equal(t, Float64(4877398396442247168), c)

	test := Float128{High: 4629393042053316608, Low: 4629393042053316608}
	d := F128ToF64(test)
	assert.Equal(t, Float64(4894998396442247172), d)
	f128d := F64ToF128(d)
	assert.Equal(t, Float128{High: 4629393042053316608, Low: 4611686018427387904}, f128d)
}
func TestFloat128_IsNan(t *testing.T) {
	f128 := Float128{High: 0x7FFF000000000000, Low: 1}
	assert.Equal(t, true, f128.IsNan())

	f128 = Float128{0x7FFF000000000000, 0}
	assert.Equal(t, false, f128.IsNan())
}

func TestFloat128_Add(t *testing.T) {
	a := Float128{High: 4613251722985340928, Low: 0} //100
	b := Float128{High: 4613251722985340928, Low: 0} //100
	c := a.Add(b)
	assert.Equal(t, Float128{Low: 0x0, High: 0x4006900000000000}, c)

	a = Float128{High: 4613251722985340928, Low: 0}  //100
	b = Float128{High: 13836623759840116736, Low: 0} //-100
	c = a.Add(b)
	assert.Equal(t, Float128{Low: 0x0, High: 0x0}, c)

	a = Float128{High: 13836623759840116736, Low: 0} //-100
	b = Float128{High: 13836623759840116736, Low: 0} //-100
	c = a.Add(b)
	assert.Equal(t, Float128{Low: 0x0, High: 0xc006900000000000}, c)

	a = Float128{High: 0x4008f4000024122e, Low: 0}                  //1000.0000043
	b = Float128{High: 0x3fff000048245bff, Low: 0xe000000000000000} //1.0000043
	c = a.Add(b)
	assert.Equal(t, Float128{Low: 0xfff0000000000000, High: 0x4008f4800048245b}, c)

	//fmt.Printf("%#v\n",c)

}

func TestFloat128_Sub(t *testing.T) {
	a := Float128{High: 4613251722985340928, Low: 0} //100
	b := Float128{High: 4613251722985340928, Low: 0} //100
	c := a.Sub(b)
	assert.Equal(t, Float128{Low: 0x0, High: 0x0}, c)

	a = Float128{High: 4613251722985340928, Low: 0}  //100
	b = Float128{High: 13836623759840116736, Low: 0} //-100
	c = a.Sub(b)
	assert.Equal(t, Float128{Low: 0x0, High: 0x4006900000000000}, c)

	a = Float128{High: 13836623759840116736, Low: 0} //-100
	b = Float128{High: 13836623759840116736, Low: 0} //-100
	c = a.Sub(b)
	assert.Equal(t, Float128{Low: 0x0, High: 0x0}, c)

	a = Float128{High: 0x4008f4000024122e, Low: 0}                  //1000.0000043
	b = Float128{High: 0x3fff000048245bff, Low: 0xe000000000000000} //1.0000043
	c = a.Sub(b)
	assert.Equal(t, Float128{Low: 0x10000000000000, High: 0x4008f38000000000}, c)
}

func Test_mul128To256M(t *testing.T) {
	a := [4]uint64{1, 2, 3, 4}
	a = softfloat_mul128To256M(1, 2, 3, 4, a)
}

func TestFloat128_Mul(t *testing.T) {
	a := Float128{High: 4613251722985340928, Low: 0} //100
	b := Float128{High: 4613251722985340928, Low: 0} //100
	c := a.Mul(b)
	assert.Equal(t, Float128{Low: 0x0, High: 0x400c388000000000}, c)

	a = Float128{High: 4613251722985340928, Low: 0}  //100
	b = Float128{High: 13836623759840116736, Low: 0} //-100
	c = a.Mul(b)
	assert.Equal(t, Float128{Low: 0x0, High: 0xc00c388000000000}, c)

	a = Float128{High: 13836623759840116736, Low: 0} //-100
	b = Float128{High: 13836623759840116736, Low: 0} //-100
	c = a.Mul(b)
	assert.Equal(t, Float128{Low: 0x0, High: 0x400c388000000000}, c)

	a = Float128{High: 0x4008f4000024122e, Low: 0}                  //1000.0000043
	b = Float128{High: 0x3fff000048245bff, Low: 0xe000000000000000} //1.0000043
	c = a.Mul(b)
	assert.Equal(t, Float128{Low: 0xebbc74fc05ba4000, High: 0x4008f4008d0b15e7}, c)

}

func TestFloat128_Div(t *testing.T) {
	a := Float128{High: 4613251722985340928, Low: 0} //100
	b := Float128{High: 4613251722985340928, Low: 0} //100
	c := a.Div(b)
	assert.Equal(t, Float128{Low: 0x0, High: 0x3fff000000000000}, c) //1

	a = Float128{High: 4613251722985340928, Low: 0}  //100
	b = Float128{High: 13836623759840116736, Low: 0} //-100
	c = a.Div(b)
	assert.Equal(t, Float128{Low: 0x0, High: 0xbfff000000000000}, c) //-1

	a = Float128{High: 13836623759840116736, Low: 0} //-100
	b = Float128{High: 13836623759840116736, Low: 0} //-100
	c = a.Div(b)
	assert.Equal(t, Float128{Low: 0x0, High: 0x3fff000000000000}, c) //1

	a = Float128{High: 0x4008f40000000000, Low: 0} //1000
	b = Float128{High: 0x4002400000000000, Low: 0} //10
	c = a.Div(b)
	assert.Equal(t, Float128{Low: 0x0, High: 0x4005900000000000}, c) //100

	a = Float128{High: 0x4008f4000024122e, Low: 0}                  //1000.0000043
	b = Float128{High: 0x3fff000048245bff, Low: 0xe000000000000000} //1.0000043
	c = a.Div(b)
	assert.Equal(t, Float128{Low: 0x53ec7b283e446b1, High: 0x4008f3ff733d3629}, c)

}

func Test_toFloat(t *testing.T) {
	f64 := math.Float64bits(float64(1000.0000043))
	a := F64ToF128(Float64(f64))

	f64_2 := math.Float64bits(float64(1.0000043))
	b := F64ToF128(Float64(f64_2))

	c := a.Add(b)
	fmt.Printf("%#v,%#v,%#v\n", a, b, c)

	c.High = 4614206648838792283
	c.Low = 18442240474082181120
	fmt.Printf("c :   %#v\n", c)
	ref64 := F128ToF64(c)
	re := math.Float64frombits(uint64(ref64))
	fmt.Println("re  ", re)

	c.High = 0x4008f4800048245b
	c.Low = 0xfff0000000000000
	ref64 = F128ToF64(c)
	re = math.Float64frombits(uint64(ref64))
	fmt.Println("re2 ", re)

	f64 = math.Float64bits(float64(-3.75))
	a = F64ToF128(Float64(f64))
	fmt.Printf("aaaaaa    %#v\n", a)

	gg := uint64(16986580520311341056)
	byteslice := make([]byte, 8)
	binary.BigEndian.PutUint64(byteslice, gg)

	fmt.Printf("%#v\n", byteslice)
	gg = uint64(4614206101444564455)
	byteslice = make([]byte, 8)
	binary.BigEndian.PutUint64(byteslice, gg)

	fmt.Printf("%#v\n", byteslice)

	c.High = 4614206096716674601
	c.Low = 377958988276582065
	fmt.Printf("cccccccc :   %#v\n", c)
	ref64 = F128ToF64(c)
	re = math.Float64frombits(uint64(ref64))
	fmt.Println("C++ re  ", re)

	f6464 := float64(1000.0000043) / float64(1.0000043)
	fmt.Println(f6464)
}

func TestF128_to_extF80(t *testing.T) {
	f64 := math.Float64bits(float64(1.0003233))
	a := F64ToF128(Float64(f64))
	fmt.Printf("aaaaaa    %v\n", a)
	abyte, _ := rlp.EncodeToBytes(f64)
	fmt.Printf("rlp result :::  %v\n", abyte)

	f64_2 := math.Float64bits(float64(1.0003233))
	//a2 := F64ToF128(Float64(f64_2))
	//fmt.Printf("aaaaaa    %#v\n", a2)

	bbyte, _ := rlp.EncodeToBytes(f64_2)
	fmt.Printf("bbbyte:  %#v\n", bbyte)

	//b := f128_to_extF80(a)
	var b ExtFloat80_t
	f128M_to_extF80M(&a, &b)

	fmt.Println(a.High, a.Low, b)
	//fmt.Printf("%b\n",b.signIf)
	//bbbb := make([]byte,8)
	//binary.BigEndian.PutUint64(bbbb,b.signIf)
	//fmt.Printf("%#v\n",bbbb)
	////binary.LittleEndian.PutUint64(bbbb,b.signIf)
	////fmt.Printf("%#v",bbbb)
	//re64 :=binary.LittleEndian.Uint64(bbbb)
	//fmt.Println(re64)
	//
	//b4 := make([]byte,8)
	//
	//binary.BigEndian.PutUint16(b4,b.signExp)
	//re16 :=binary.LittleEndian.Uint16(b4)
	//fmt.Println(re16)

	var ef *ExtFloat = (*ExtFloat)(unsafe.Pointer(&b))

	fmt.Println(ef.signIf, ef.signExp)
	//
	//var out ExtFloat80_t
	//
	//out.signExp = uint16(b.signIf)
	//out.signIf = uint64(b.signExp)<<48 +b.signIf>>16
	//fmt.Println(out)

}

func f128M_to_extF80M(aPtr *Float128, zPtr *ExtFloat80_t) {
	*zPtr = f128_to_extF80(*aPtr)
}

type ExtFloat struct {
	signIf  uint64
	signExp uint16
}

func Test_ext(t *testing.T) {
	f64 := math.Float64bits(float64(1.003))
	a := F64ToF128(Float64(f64))
	fmt.Printf("aaaaaa    %#v\n", a)
	var b ExtFloat80_t
	f128M_to_extF80M(&a, &b)

	fmt.Println(a.High, a.Low, b)

	var ef *ExtFloat = (*ExtFloat)(unsafe.Pointer(&b))

	fmt.Println(ef.signIf, ef.signExp)
}

func Test_float32(t *testing.T) {
	f32 := math.Float32bits(float32(1.0003233))
	fmt.Println(f32)
	a := F32ToF128(Float32(f32))
	fmt.Printf("aaaaaa    %#v\n", a)
	abyte, _ := rlp.EncodeToBytes(f32)
	fmt.Printf("rlp result :::  %#v\n", abyte)
}

func Test_Fixtfti(t *testing.T) {

	f64 := math.Float64bits(float64(-99999999899.9999999))
	f128 := F64ToF128(Float64(f64))
	int128 := Fixtfti(f128)
	fmt.Printf("%#v\n", int128)

	fmt.Println(int128.Low)

}

//1 000 0000 0000 1011    00100000000000000000000000000000000000000000000

func Test_Fixunstfti(t *testing.T) {
	f64 := math.Float64bits(float64(1004432323.990))
	f128 := F64ToF128(Float64(f64))
	fmt.Println(f128.Low, f128.High)
	uint128 := Fixunstfti(f128)
	fmt.Printf("%#v\n", uint128)

	fmt.Println(uint128.Low)
}

func Test_Floattidf(t *testing.T) {
	in := CreateInt128(1)
	in.LeftShifts(60)
	fmt.Printf("%#b", in)
	out := Floattidf(in)
	fmt.Println(out)

}

func TestFloatuntidf(t *testing.T) {
	in := CreateUint128(1)
	//in.LeftShifts(60)
	fmt.Println(in)
	out := Floatuntidf(in)
	fmt.Println(out)
}

func TestCount(t *testing.T) {
	in := CreateInt128(1)
	in.LeftShifts(80)
	count := clzti2(in)
	fmt.Println(count)
}

func TestFloat64String(t *testing.T) {
	f := float64(1.0023)
	f64 := math.Float64bits(f)
	var ff Float64
	ff = Float64(f64)

	assert.Equal(t, ff.String(), "1.002300e+00")

	f = float64(1.111111e100)
	f64 = math.Float64bits(f)
	ff = Float64(f64)

	assert.Equal(t, ff.String(), "1.111111e+100")
}

func TestFloat32String(t *testing.T) {
	f := float32(1.0023)
	f32 := math.Float32bits(f)
	ff := Float32(f32)
	assert.Equal(t, ff.String(), "1.002300e+00")

	f = float32(1.111111e33)
	f32 = math.Float32bits(f)
	ff = Float32(f32)
	assert.Equal(t, ff.String(), "1.111111e+33")
}

func TestUint256String(t *testing.T) {
	a := new(big.Int).SetUint64(math.MaxUint64)
	b := new(big.Int).Mul(a, a)
	c := new(big.Int).Mul(b, b)
	var d uint64
	d = math.MaxUint64
	e := MulUint64(uint64(d), uint64(d))

	f := Uint256{Low: e}
	g := f.Mul(f)
	fmt.Println(g)
	assert.Equal(t, g.String(), c.String())
}
func TestFloat128_String(t *testing.T) {
	plusF128 := Float128{High: 4613251722985340928, Low: 0} //100
	a := plusF128.String()
	fmt.Println(a)
	fmt.Println(plusF128.String2())

	ff, _ := strconv.ParseFloat("6.666666666666666667e-07", 64)
	fmt.Println(ff)

	//0.5,-3.75

}

func TestFloat128EXp(t *testing.T) {
	//f := float64(0.5)
	f := float64(1)
	f64 := math.Float64bits(f)
	f128 := F64ToF128(Float64(f64))

	ext80 := f128_to_extF80(f128)
	fmt.Printf("%b", ext80)
	fmt.Println(ext80)
	if ext80.signExp&0x8000 != 0 {
		fmt.Println("-")
	} else {
		fmt.Println("+")
	}
	exp := ext80.signExp & 0x7FFF
	fmt.Printf("exp:%b\n", exp)

	fmt.Println("two: ")
	fmt.Printf("%b\n%b\n", 0x3FFE, exp)
	E := exp - 0x3FFF
	fmt.Println("E:", E)
	fmt.Printf("%b\n", exp)
	str := strconv.FormatFloat(f, 'e', 18, 64)
	fmt.Println(str)
	fmt.Printf("%b\n", 0x3FFF-1)
}

//111 1111  1111  1111
// 11 1111  1111  1110

//assert.Equal(t, r[0], "5.000000000000000000e-01")
////-3.750000000000000000e+00
//assert.Equal(t, r[1], "-3.750000000000000000e+00")
//
//assert.Equal(t, r[2], "6.666666666666666667e-07")
//111 1111 1111 1110

func TestF64(t *testing.T) {
	f := float64(10.111)
	f64 := math.Float64bits(f)
	fmt.Printf("%b, %b\n%b\n", f64>>52, f64&0xFFFFFFFFFFFF, f64)
	fmt.Println(strconv.FormatFloat(f, 'e', 18, 64))
	fmt.Printf("%b\n", 0x3FE)
	expdiff := f64>>52 - 0x3FE
	fmt.Println(expdiff)
}

//1111 1111 11 0000000000000000000000000000000000000000000000000000
