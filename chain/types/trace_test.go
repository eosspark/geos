package types

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFlatSet_Append(t *testing.T) {
	f := FlatSet{}
	for i := 0; i < 10000; i++ {
		if i != 8000 {
			ad := AccountDelta{}
			ad.Account = common.AccountName(i)
			ad.Delta = int64(i)
			f.data = append(f.data, ad)
		}
	}
	param := AccountDelta{}
	param.Account = common.AccountName(8000)
	param.Delta = int64(8000)
	f.Append(common.AccountName(8000), int64(8000))
	assert.Equal(t, param, f.data[8000])
}

func TestFlatSet_Append2(t *testing.T) {
	f := FlatSet{}
	for i := 0; i < 10; i++ {
		if i != 8 {
			ad := AccountDelta{}
			ad.Account = common.AccountName(i)
			ad.Delta = int64(i)
			f.data = append(f.data, ad)
		}
	}

	param := AccountDelta{}
	param.Account = common.AccountName(8)
	param.Delta = int64(8)
	f.Append(common.AccountName(8), int64(8))

	assert.Equal(t, param, f.data[8])
}

func Test_Sear(t *testing.T) {
	f := FlatSet{}
	for i := 0; i < 10; i++ {
		if i != 8 {
			ad := AccountDelta{}
			ad.Account = common.AccountName(i)
			ad.Delta = int64(i)
			f.data = append(f.data, ad)
		}
	}

	param := AccountDelta{}
	param.Account = common.AccountName(8)
	param.Delta = int64(8)

	sear(len(f.data), func(i int) bool {
		r := false
		if f.data[i].Account < param.Account && f.data[i+1].Account > param.Account {
			fmt.Println(i)
			fmt.Println(f.data[i])
			fmt.Println(f.data[i+1])
			r = true
		}
		return r
	})
}

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
