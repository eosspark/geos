package arithmeticTypes

import (
	"fmt"
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
	fmt.Println(d)
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
	fmt.Println(f128d)
	assert.Equal(t, Float128{High: 4629393042053316608, Low: 4611686018427387904}, f128d)
}
