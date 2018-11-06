package arithmeticTypes

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

//func Test_toFloat(t *testing.T) {
//	f64 := math.Float64bits(float64(1000.0000043))
//	a := F64ToF128(Float64(f64))
//
//	f64_2 := math.Float64bits(float64(1.0000043))
//	b := F64ToF128(Float64(f64_2))
//
//	c := a.Add(b)
//	fmt.Printf("%#v,%#v,%#v\n", a, b, c)
//
//	c.High = 4614206648838792283
//	c.Low = 18442240474082181120
//	fmt.Printf("c :   %#v\n", c)
//	ref64 := F128ToF64(c)
//	re := math.Float64frombits(uint64(ref64))
//	fmt.Println("re  ", re)
//
//	c.High = 0x4008f4800048245b
//	c.Low = 0xfff0000000000000
//	ref64 = F128ToF64(c)
//	re = math.Float64frombits(uint64(ref64))
//	fmt.Println("re2 ", re)
//
//	f64 = math.Float64bits(float64(1000.0043042163287))
//	a = F64ToF128(Float64(f64))
//	fmt.Printf("%#v\n", a)
//
//	gg := uint64(16986580520311341056)
//	byteslice := make([]byte, 8)
//	binary.BigEndian.PutUint64(byteslice, gg)
//
//	fmt.Printf("%#v\n", byteslice)
//	gg = uint64(4614206101444564455)
//	byteslice = make([]byte, 8)
//	binary.BigEndian.PutUint64(byteslice, gg)
//
//	fmt.Printf("%#v\n", byteslice)
//
//	c.High = 4614206096716674601
//	c.Low = 377958988276582065
//	fmt.Printf("cccccccc :   %#v\n", c)
//	ref64 = F128ToF64(c)
//	re = math.Float64frombits(uint64(ref64))
//	fmt.Println("C++ re  ", re)
//
//}
