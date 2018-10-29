package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlatSet_Append(t *testing.T) {

	before, end, midle := 0, 19, 8
	if before == 0 {
		f := []common.FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 0 {
				ad := AccountDelta{}
				ad.Account = common.AccountName(i)
				ad.Delta = int64(i)
				f = append(f, ad)
			}

		}
		param := AccountDelta{}
		param.Account = common.AccountName(0)
		param.Delta = int64(0)
		result, _ := common.Append(f, param)
		assert.Equal(t, param, (*result)[0].(AccountDelta))
	}
	if end == 19 {
		f := []common.FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 19 {
				ad := AccountDelta{}
				ad.Account = common.AccountName(i)
				ad.Delta = int64(i)
				f = append(f, ad)
			}
		}
		param := AccountDelta{}
		param.Account = common.AccountName(19)
		param.Delta = int64(19)

		result, _ := common.Append(f, param)
		assert.Equal(t, param, (*result)[19].(AccountDelta))
	}
	if midle == 8 {
		f := []common.FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 8 {
				ad := AccountDelta{}
				ad.Account = common.AccountName(i)
				ad.Delta = int64(i)
				f = append(f, ad)
			}
		}
		param := AccountDelta{}
		param.Account = common.AccountName(8)
		param.Delta = int64(8)

		result, _ := common.Append(f, param)
		assert.Equal(t, param, (*result)[8].(AccountDelta))
	}
}

//assert.Equal(t, param, f.data[8000])

func sear(n int, f func(int) bool) int {
	i, j := 0, n
	for i < j {
		h := int(uint(i+j) >> 1)
		//fmt.Println("exec count：")
		if !f(h) {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
