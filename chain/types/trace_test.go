package types

/*func TestFlatSet_Append(t *testing.T) {

	before, end, midle := 0, 19, 8
	if before == 0 {
		f := common.FlatSet{}
		for i := 0; i < 20; i++ {
			if i != 0 {
				ad := AccountDelta{}
				ad.Account = common.AccountName(i)
				ad.Delta = int64(i)
				f.Data = append(f.Data, &ad)
			}

		}
		element := AccountDelta{}
		element.Account = common.AccountName(0)
		element.Delta = int64(0)
		result, _ := f.Insert(&element)
		assert.Equal(t, &element, result)
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
		element := AccountDelta{}
		element.Account = common.AccountName(19)
		element.Delta = int64(19)

		result, _ := common.Append(f, element)
		assert.Equal(t, element, (*result)[19].(AccountDelta))
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
		element := AccountDelta{}
		element.Account = common.AccountName(8)
		element.Delta = int64(8)

		result, _ := common.Append(f, element)
		assert.Equal(t, element, (*result)[8].(AccountDelta))
	}
}
*/
//assert.Equal(t, element, f.data[8000])

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
