package int

import (
	"fmt"
	"testing"
)

func TestIndex_Insert(t *testing.T) {
	mi := NewTestIndex()
	byIdIndex := mi.GetById()
	defer func() { mi.Free() }()

	mi.Insert(1)
	//
	mi.Insert(2)

	mi.Insert(3)
	mi.Insert(4)
	mi.Insert(5)

	print(mi)

	itr2 := byIdIndex.Find(2)


	byIdIndex.Erase(itr2)

	print(mi)

	itr5 := byIdIndex.Find(5)
	byIdIndex.Modify(itr5, func(i *int) {
		*i --
	})

	print(mi)
	//
	//byNum := mi.GetByNum()
	//
	//it := byNum.LowerBound(3)
	//byNum.Modify(it, func(i *int) {
	//	*i *= 3
	//})
	//
	//print(mi)

}

func print(mi *TestIndex) {
	for it := mi.GetById().Begin(); it.HasNext();it.Next() {
		fmt.Print(it.Value()," ")
	}
	fmt.Println()

	//for itr := mi.GetByNum().Begin(); itr.HasNext(); itr.Next() {
	//	fmt.Print(itr.Value(), " ")
	//}
	//fmt.Println()
	//
	//mi.GetByPrev().Each(func(key int, obj int) {
	//	fmt.Print(obj, " ")
	//})
	//
	//fmt.Println()
	//fmt.Println()
}

//
//func TestRollback(t *testing.T) {
//	m := NewTestIndex()
//
//	m.Insert(1)
//	m.Insert(1)
//	m.Insert(1)
//	m.Insert(2)
//	m.Insert(3)
//
//	fmt.Println(m.Size())
//	fmt.Println(m.GetById().Size())
//	fmt.Println(m.GetByNum().Size())
//	fmt.Println(m.GetByPrev().Size())
//}
