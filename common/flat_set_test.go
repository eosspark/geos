package common

import (
	"encoding/binary"
	"fmt"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/stretchr/testify/assert"
	"testing"
)

type AccountDeltaDemo struct {
	A AccountName
	B int64
}

func (a *AccountDeltaDemo) GetKey() uint64 {
	return uint64(a.A)
}

func TestFlat(t *testing.T) {
	before, end, midle, eq := 0, 19, 8, 10
	if before == 0 {
		f := FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 0 {
				ad := AccountDeltaDemo{}
				ad.A = AccountName(i)
				ad.B = int64(i)
				f.Data = append(f.Data, &ad)
			}
		}
		param := AccountDeltaDemo{}
		param.A = AccountName(0)
		param.B = int64(0)
		result, p := f.Insert(&param)
		fmt.Println(result, p)
		assert.Equal(t, &param, result)
		assert.Equal(t, false, p)
	}
	if end == 19 {
		f := FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 19 {
				ad := AccountDeltaDemo{}
				ad.A = AccountName(i)
				ad.B = int64(i)
				f.Data = append(f.Data, &ad)
			}
		}
		param := AccountDeltaDemo{19, 19}

		result, p := f.Insert(&param)
		fmt.Println(result, p)
		assert.Equal(t, &param, result)
		assert.Equal(t, false, p)
	}
	if midle == 8 {
		f := FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 8 {
				ad := AccountDeltaDemo{}
				ad.A = AccountName(i)
				ad.B = int64(i)
				f.Data = append(f.Data, &ad)
			}
		}
		param := AccountDeltaDemo{8, 8}

		result, p := f.Insert(&param)
		fmt.Println(result, p)
		assert.Equal(t, &param, result)
		assert.Equal(t, false, p)
	}
	if eq == 10 {
		f := FlatSet{}
		for i := 0; i < 20; i++ {
			ad := AccountDeltaDemo{}
			ad.A = AccountName(i)
			ad.B = int64(i)
			f.Data = append(f.Data, &ad)
		}
		param := AccountDeltaDemo{8, 8}

		result, p := f.Insert(&param)
		fmt.Println(result, p)
		assert.Equal(t, &param, result)
		assert.Equal(t, true, p)
	}
}

func TestFlatSet_GetData(t *testing.T) {
	f := FlatSet{}
	for i := 0; i < 20; i++ {

		ad := AccountDeltaDemo{}
		ad.A = AccountName(i)
		ad.B = int64(i)
		f.Data = append(f.Data, &ad)

	}
	param := AccountDeltaDemo{8, 8}

	result := f.GetData(8)
	r := f.GetData(20)
	assert.Equal(t, &param, result)
	assert.Equal(t, true, Empty(r))

}

func TestFlatSet_Clear(t *testing.T) {
	f := FlatSet{}
	for i := 0; i < 20; i++ {
		ad := AccountDeltaDemo{}
		ad.A = AccountName(i)
		ad.B = int64(i)
		f.Data = append(f.Data, &ad)
	}
	f.Clear()

	param := AccountDeltaDemo{0, 0}

	result, b := f.Insert(&param)
	assert.Equal(t, &param, result)
	assert.Equal(t, false, b)
}

func Test(t *testing.T) {
	f := FlatSet{}
	for i := 0; i < 20; i++ {
		if i != 18 {
			ad := AccountDeltaDemo{}
			ad.A = AccountName(i)
			ad.B = int64(i)
			f.Data = append(f.Data, &ad)
		}
	}
	element := AccountDeltaDemo{8, 8}

	length := len(f.Data)
	r, i, j := 0, 0, length-1
	for i < j {
		h := int(uint(i+j) >> 1)
		if f.Data[h].GetKey() <= element.GetKey() {
			i = h + 1
		} else {
			j = h
		}
		r = h
	}
	fmt.Println(r)
}

func Test_test(t *testing.T) {
	array := []byte{0x00, 0x01, 0x08, 0x00, 0x08, 0x01, 0xab, 0x01}
	num := binary.LittleEndian.Uint64(array)
	fmt.Printf("%v, %x", array, num)
}

func TestFlatSet_Find(t *testing.T) {
	key, _ := ecc.NewPublicKey("EOS5fxEptrpsG2QTjRgi8Gf9EConFDH3jeUc24YFSemcW3bBDhuoW")
	fs := FlatSet{}
	k, b := fs.Insert(&key)
	fmt.Println("insert key:", k, "exist:", b)
	boo := fs.Find(&key)
	assert.Equal(t, true, boo)
}
