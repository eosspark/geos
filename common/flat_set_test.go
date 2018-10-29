package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type AccountDeltaDemo struct {
	A AccountName
	B int64
}

func (f AccountDeltaDemo) Compare(first FlatSet, second FlatSet) bool {
	return first.(AccountDeltaDemo).A <= second.(AccountDeltaDemo).A

}

/*
func (f AccountDeltaDemo) GetData() []Flat{
	return f.Data
}

func (f AccountDeltaDemo) SetData(param []Flat){
	f.Data = param
}*/

func TestFlat(t *testing.T) {
	before, end, midle := 0, 19, 8
	if before == 0 {
		f := []FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 0 {
				ad := AccountDeltaDemo{}
				ad.A = AccountName(i)
				ad.B = int64(i)
				f = append(f, ad)
			}
		}
		param := AccountDeltaDemo{}
		param.A = AccountName(0)
		param.B = int64(0)
		result, _ := Append(f, param)
		assert.Equal(t, param, (*result)[0].(AccountDeltaDemo))
	}
	if end == 19 {
		f := []FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 19 {
				ad := AccountDeltaDemo{}
				ad.A = AccountName(i)
				ad.B = int64(i)
				f = append(f, ad)
			}
		}
		param := AccountDeltaDemo{}
		param.A = AccountName(19)
		param.B = int64(19)

		result, _ := Append(f, param)
		assert.Equal(t, param, (*result)[19].(AccountDeltaDemo))
	}
	if midle == 8 {
		f := []FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 8 {
				ad := AccountDeltaDemo{}
				ad.A = AccountName(i)
				ad.B = int64(i)
				f = append(f, ad)
			}
		}
		param := AccountDeltaDemo{}
		param.A = AccountName(8)
		param.B = int64(8)

		result, _ := Append(f, param)
		assert.Equal(t, param, (*result)[8].(AccountDeltaDemo))
	}
}
