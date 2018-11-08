package common

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"testing"
)

type AccountDeltaDemo struct {
	A AccountName
	B int64
}

func (a AccountDeltaDemo) GetKey() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(a.A))
	return b
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
				f.Insert(ad)
			}
		}
		param := AccountDeltaDemo{}
		param.A = AccountName(0)
		param.B = int64(0)
		result, p := f.Insert(param)
		assert.Equal(t, param, result)
		assert.Equal(t, false, p)
	}
	if end == 19 {
		f := FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 19 {
				ad := AccountDeltaDemo{}
				ad.A = AccountName(i)
				ad.B = int64(i)
				f.Insert(ad)
			}
		}
		param := AccountDeltaDemo{19, 19}

		result, p := f.Insert(param)
		assert.Equal(t, param, result)
		assert.Equal(t, false, p)
	}
	if midle == 8 {
		f := FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 8 {
				ad := AccountDeltaDemo{}
				ad.A = AccountName(i)
				ad.B = int64(i)
				f.Insert(&ad)
			}
		}
		param := AccountDeltaDemo{8, 8}

		result, p := f.Insert(param)
		assert.Equal(t, param, result)
		assert.Equal(t, false, p)
	}
	if eq == 10 {
		f := FlatSet{}
		for i := 0; i < 20; i++ {
			ad := AccountDeltaDemo{}
			ad.A = AccountName(i)
			ad.B = int64(i)
			f.Insert(ad)
		}
		param := AccountDeltaDemo{8, 8}

		result, p := f.Insert(param)
		assert.Equal(t, param, result)
		assert.Equal(t, false, p)
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
		f.Data = append(f.Data, ad)
	}
	f.Clear()

	param := AccountDeltaDemo{0, 0}

	result, b := f.Insert(param)
	assert.Equal(t, param, result)
	assert.Equal(t, false, b)
}

func Test_Insert(t *testing.T) {
	f := FlatSet{}
	for i := 0; i < 5; i++ {
		if i != 2 {
			ad := AccountDeltaDemo{}
			ad.A = AccountName(i)
			ad.B = int64(i)
			f.Insert(ad)
		}
	}
	param := AccountDeltaDemo{2, 2}
	result, p := f.Insert(param)
	assert.Equal(t, param, result)
	assert.Equal(t, false, p)
}

func TestFlatSet_Update(t *testing.T) {
	f := FlatSet{}
	ad := AccountDeltaDemo{10, 10}
	b := f.Update(ad)
	assert.Equal(t, true, b)
	for i := 0; i < 5; i++ {
		if i != 2 {
			ad := AccountDeltaDemo{}
			ad.A = AccountName(i)
			ad.B = int64(i)
			f.Insert(ad)
		}
	}

	add := AccountDeltaDemo{2, 10}
	p := f.Update(add)
	assert.Equal(t, true, p)
}

func TestFlatSet_RemoveOne(t *testing.T) {
	f := FlatSet{}
	ad := AccountDeltaDemo{10, 10}

	result, i := f.Insert(ad)
	assert.Equal(t, result, ad)
	assert.Equal(t, false, i)
	re := f.Remove(ad.GetKey())
	assert.Equal(t, true, re)
}

func TestFlatSet_Remove(t *testing.T) {
	f := FlatSet{}
	ad := AccountDeltaDemo{2, 2}
	ad1 := AccountDeltaDemo{1, 1}

	ad2 := AccountDeltaDemo{0, 0}

	ad3 := AccountDeltaDemo{3, 3}
	f.Insert(ad)
	f.Insert(ad1)
	f.Insert(ad2)
	re := f.Remove(ad1.GetKey())
	assert.Equal(t, true, re)
	r := f.Remove(ad3.GetKey())
	assert.Equal(t, false, r)
}

func Test_Find(t *testing.T) {
	f := FlatSet{}
	a := AccountName(6138663577826885632)
	b := AccountName(17765913279651119101)

	ele1, b1 := f.Insert(&a)
	ele2, b2 := f.Insert(&b)

	assert.Equal(t, &a, ele1)
	assert.Equal(t, &b, ele2)

	assert.Equal(t, false, b1)
	assert.Equal(t, false, b2)
}

func Test_Nil(t *testing.T) {
	f := FlatSet{}
	ele, i := f.FindData([]byte("1"))
	assert.Equal(t, nil, ele)
	assert.Equal(t, -1, i)
}
