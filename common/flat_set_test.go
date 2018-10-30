package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type AccountDeltaDemo struct {
	A AccountName
	B int64
}

func (f AccountDeltaDemo) Compare(first Element, second Element) bool {
	return first.(AccountDeltaDemo).A <= second.(AccountDeltaDemo).A
}

func (f AccountDeltaDemo) Equal(first Element, second Element) bool {
	return first.(AccountDeltaDemo).A == second.(AccountDeltaDemo).A
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
				f.Data = append(f.Data, ad)
				//f.SetData(f.Data)
			}
		}
		param := AccountDeltaDemo{}
		param.A = AccountName(0)
		param.B = int64(0)
		result, p := f.Insert(param)

		assert.Equal(t, param, *result)
		assert.Equal(t, false, p)
	}
	if end == 19 {
		f := FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 19 {
				ad := AccountDeltaDemo{}
				ad.A = AccountName(i)
				ad.B = int64(i)

				f.Data = append(f.Data, ad)
			}
		}
		param := AccountDeltaDemo{19, 19}

		result, p := f.Insert(param)
		//fmt.Println(result)
		assert.Equal(t, param, *result)
		assert.Equal(t, false, p)
	}
	if midle == 8 {
		f := FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 8 {
				ad := AccountDeltaDemo{}
				ad.A = AccountName(i)
				ad.B = int64(i)
				f.Data = append(f.Data, ad)
			}
		}
		param := AccountDeltaDemo{8, 8}

		result, p := f.Insert(param)
		//fmt.Println(result,p)
		assert.Equal(t, param, *result)
		assert.Equal(t, false, p)
	}
	if eq == 10 {
		f := FlatSet{}
		for i := 0; i < 20; i++ {
			ad := AccountDeltaDemo{}
			ad.A = AccountName(i)
			ad.B = int64(i)
			f.Data = append(f.Data, ad)
		}
		param := AccountDeltaDemo{8, 8}

		result, p := f.Insert(param)
		//fmt.Println(result,p)
		assert.Equal(t, param, *result)
		assert.Equal(t, true, p)
	}
}

func TestFlatSet_GetData(t *testing.T) {
	f := FlatSet{}
	for i := 0; i < 20; i++ {

		ad := AccountDeltaDemo{}
		ad.A = AccountName(i)
		ad.B = int64(i)
		f.Data = append(f.Data, ad)

	}
	param := AccountDeltaDemo{8, 8}

	result := f.GetData(8)
	r := f.GetData(20)
	assert.Equal(t, true, Empty(r))
	assert.Equal(t, param, *result)

}
