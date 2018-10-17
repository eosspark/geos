package figure

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
