package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakePair(t *testing.T) {
	test := MakePair(100, 99999)
	assert.Equal(t, test.First, 100)
	assert.Equal(t, test.Second, 99999)
}

func TestNewPair(t *testing.T) {
	test := Pair{First: 100, Second: 8}
	check := struct {
		First  int
		Second int
	}{100, 8}

	assert.Equal(t, fmt.Sprintf("%v\n", test), fmt.Sprintf("%v\n", check))
	assert.Equal(t, test.First, 100)
	assert.Equal(t, test.Second, 8)

}

func TestGetIndex(t *testing.T) {
	test := Pair{First: 300, Second: 77778}
	check := []byte{0x2c, 0x1, 0x0, 0x0, 0xd2, 0x2f, 0x1, 0x0}
	b := test.GetIndex()
	fmt.Println(b)
	assert.Equal(t, b, check)
	//fmt.Printf("%#v\n", b)
}

func TestMakeTuple(t *testing.T) {
	test := MakeTuple(100, 99999, 80000)
	assert.Equal(t, test.First, 100)
	assert.Equal(t, test.Second, 99999)
	assert.Equal(t, test.Third, 80000)
}

func TestNewTuple(t *testing.T) {
	test := Tuple{First: 100, Second: 8, Third: 9999}
	check := struct {
		First  int
		Second int
		Third  int
	}{100, 8, 9999}

	assert.Equal(t, fmt.Sprintf("%v\n", test), fmt.Sprintf("%v\n", check))
	assert.Equal(t, test.First, 100)
	assert.Equal(t, test.Second, 8)
	assert.Equal(t, test.Third, 9999)

}
func TestTupleGetIndex(t *testing.T) {
	test := Tuple{First: 300, Second: 77778, Third: 9000}
	check := []byte{0x2c, 0x1, 0x0, 0x0, 0xd2, 0x2f, 0x1, 0x0, 0x28, 0x23, 0x0, 0x0}
	b := test.GetIndex()
	assert.Equal(t, b, check)
	//fmt.Printf("%#v\n", b)
}

func TestNewTuple4(t *testing.T) {
	test := Tuple4{First: 100, Second: 8, Third: 9999, Fourth: 1000000001}
	check := struct {
		First  int
		Second int
		Third  int
		Fourth int
	}{100, 8, 9999, 1000000001}

	assert.Equal(t, fmt.Sprintf("%v\n", test), fmt.Sprintf("%v\n", check))
	assert.Equal(t, test.First, 100)
	assert.Equal(t, test.Second, 8)
	assert.Equal(t, test.Third, 9999)
	assert.Equal(t, test.Fourth, 1000000001)

}
