package common

/*func Test_Reflect(t *testing.T){
	o :=ObjectDemo{}
	o.ID = 234
	o.Prev = 123
	o.BlockNum = 1

	//templet := TempleteProperty{}
	rv:=reflect.ValueOf(o)
	fmt.Println(rv.Interface().(ObjectDemo).ID.String())
	fmt.Println(rv.Field(2))
	fmt.Println("----")


	rt:=reflect.TypeOf(o)
	fmt.Println(rt.Name())
}*/

import (
	"fmt"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/stretchr/testify/assert"
	"testing"
)

type AccountDeltaDemo struct {
	A AccountName
	B int64
}

/*func (a AccountDeltaDemo) GetKey() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}*/

func (a AccountDeltaDemo) GetKey() uint64 {
	/*b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)*/
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
				//f.Data = append(f.Data, ad)
				f.Insert(ad)
			}
		}
		param := AccountDeltaDemo{}
		param.A = AccountName(0)
		param.B = int64(0)
		result, p := f.Insert(param)
		fmt.Println(result, p)
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
		fmt.Println(result, p)
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
		fmt.Println(result, p)
		//fmt.Println("**********",f)
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
		fmt.Println(result, p)
		assert.Equal(t, param, result)
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
		f.Data = append(f.Data, ad)
	}
	f.Clear()

	param := AccountDeltaDemo{0, 0}

	result, b := f.Insert(param)
	assert.Equal(t, param, result)
	assert.Equal(t, false, b)
}

func Test(t *testing.T) {
	f := FlatSet{}
	for i := 0; i < 20; i++ {
		//if i != 18 {
		ad := AccountDeltaDemo{}
		ad.A = AccountName(i)
		ad.B = int64(i)
		f.Data = append(f.Data, ad)
		//}
	}
	element := AccountDeltaDemo{20, 19}

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
	fmt.Println("r", r)
	fmt.Println("f:", f)
	fmt.Println("f:", f.Len())
}

func Test_test(t *testing.T) {
	/*array := []byte{0x00, 0x01, 0x08, 0x00, 0x08, 0x01, 0xab, 0x01}
	num := binary.LittleEndian.Uint64(array)
	fmt.Printf("%v, %x", array, num)*/
	fmt.Println(AccountName(N("yuanchao")) == AccountName(N("yuanchao")))
}

func TestFlatSet_Find(t *testing.T) {
	key, _ := ecc.NewPublicKey("EOS5fxEptrpsG2QTjRgi8Gf9EConFDH3jeUc24YFSemcW3bBDhuoW")
	fs := FlatSet{}
	k, b := fs.Insert(&key)
	fmt.Println("insert key:", k, "exist:", b)
	boo, i := fs.Find(&key)
	fmt.Println("insert key:", boo, "exist:", i)
	assert.Equal(t, true, boo)
}

func Test1(t *testing.T) {
	f := FlatSet{}
	for i := 0; i < 5; i++ {
		if i != 2 {
			ad := AccountDeltaDemo{}
			ad.A = AccountName(i)
			ad.B = int64(i)
			f.Insert(ad)
		}
	}
	fmt.Println(f)
	param := AccountDeltaDemo{2, 2}

	result, p := f.Insert(param)
	fmt.Println(result, p, f)
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
	fmt.Println(f)
	assert.Equal(t, true, p)
}
