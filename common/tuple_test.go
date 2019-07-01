package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeTuple(t *testing.T) {
	var (
		a = int16(999)
		b = int32(4)
		c = int(10000000001)
		d = uint64(99998899889988)
		e = uint64(1200230230230000230)
		f = uint64(121)
	)

	test2 := MakeTuple(d, e)
	str2 := fmt.Sprintf("%v", test2)
	assert.Equal(t, str2, "[99998899889988 1200230230230000230]")

	test3 := MakeTuple(d, e, f)
	str3 := fmt.Sprintf("%v", test3)
	assert.Equal(t, str3, "[99998899889988 1200230230230000230 121]")

	test4 := MakeTuple(a, b, c, d)
	str4 := fmt.Sprintf("%v", test4)
	assert.Equal(t, str4, "[999 4 10000000001 99998899889988]")
}
