package arithmeticTypes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUiToFloat128(t *testing.T) {
	check := Float128{High: 4613251722985340928, Low: 0}
	a := Ui64ToF128(100)
	//fmt.Println(a.High, a.Low)
	assert.Equal(t, check, a)

	b := Ui32ToF128(100)
	//fmt.Println(b.High, b.Low)
	assert.Equal(t, check, b)

}

func TestItoFloat128(t *testing.T) {
	check := Float128{High: 13836623759840116736, Low: 0}
	a := I64ToF128(-100)
	//fmt.Println(a.High, a.Low)
	assert.Equal(t, check, a)

	b := I32ToF128(-100)
	//fmt.Println(b.High, b.Low)
	assert.Equal(t, check, b)
}

func TestF128ToUi32(t *testing.T) {
	f128 := Float128{High: 4613251722985340928, Low: 0}
	result := F128ToUi32(f128, 0, false)
	//fmt.Println(result)
	assert.Equal(t, uint32(100), result)
}

func TestF128ToUi64(t *testing.T) {
	f128 := Float128{High: 4613251722985340928, Low: 0}
	result := F128ToUi64(f128, 0, false)
	//fmt.Println(result)
	assert.Equal(t, uint64(100), result)
}

func TestF128ToI64(t *testing.T) {
	plusF128 := Float128{High: 4613251722985340928, Low: 0}

	re1 := F128ToI64(plusF128, 0, false)
	//fmt.Println(re1)
	assert.Equal(t, int64(100), re1)

	minusF128 := Float128{High: 13836623759840116736, Low: 0}
	re2 := F128ToI64(minusF128, 0, false)
	//fmt.Println(re2)
	assert.Equal(t, int64(-100), re2)
}

func TestF128ToI32(t *testing.T) {
	plusF128 := Float128{High: 4613251722985340928, Low: 0}

	re1 := F128ToI32(plusF128, 0, false)
	//fmt.Println(re1)
	assert.Equal(t, int32(100), re1)

	minusF128 := Float128{High: 13836623759840116736, Low: 0}
	re2 := F128ToI32(minusF128, 0, false)
	//fmt.Println(re2)
	assert.Equal(t, int32(-100), re2)
}

func TestF128ToF32(t *testing.T) {
	plusF128 := Float128{High: 4613251722985340928, Low: 0}

	a := F128ToF32(plusF128)
	//fmt.Println(a)
	assert.Equal(t, 1120403456, int(a))

	minusF128 := Float128{High: 13836623759840116736, Low: 0}
	b := F128ToF32(minusF128)
	//fmt.Println(b)
	assert.Equal(t, 3267887104, int(b))

	int60 := int64(1) << 60
	f12860 := I64ToF128(int60)
	//F12860 := Float128{High: 4628293042053316608, Low: 0}//
	c := F128ToF32(f12860)
	//fmt.Println(c)
	assert.Equal(t, 1568669696, int(c))
}
