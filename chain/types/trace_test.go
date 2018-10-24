package types

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestFlatSet_Less(t *testing.T) {
	f := FlatSet{}
	for i := 0; i < 10; i++ {
		ad := AccountDelta{}
		ad.Account = common.AccountName(i)
		ad.Delta = int64(i)
		f.data = append(f.data, ad)
	}

	param := AccountDelta{}
	param.Account = common.AccountName(20)
	param.Delta = int64(20)
	f.data = append(f.data, param)

	param1 := AccountDelta{}

	param1.Account = common.AccountName(15)
	param1.Delta = int64(15)
	f.data = append(f.data, param1)
	fmt.Println("sort before:", f)
	f.less = less
	sort.Sort(&f)
	fmt.Println(f)
}

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
	fmt.Println(f.Len())
	fmt.Println(f)
	fmt.Println(f.data[8000])
	assert.True(t, f.data[8000] == param)
	//assert.Equal(t,f.data[7000]==param,"this is error message")
}
