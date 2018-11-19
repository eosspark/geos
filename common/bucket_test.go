package common

import (
	"fmt"
	"github.com/eosspark/eos-go/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

type A struct {
	id  int
	num int
}

func lessCompare(first ElementObject, second ElementObject) int {
	if first.(*A).id < second.(*A).id {
		return -1
	} else if first.(*A).id == second.(*A).id {
		return 0
	} else {
		return 1
	}
}

func greaterCompare(first ElementObject, second ElementObject) int {
	if first.(*A).id > second.(*A).id {
		return -1
	} else if first.(*A).id == second.(*A).id {
		return 0
	} else {
		return 1
	}
}

func (a *A) ElementObject() {} //default implements interface only taget

func Test_ADemo(t *testing.T) {
	f := Bucket{}
	a := A{1, 1}
	b := A{2, 2}
	f.Compare = lessCompare

	log.Info("%d", f.Compare(&a, &b))
}

func TestBucket_InsertLess(t *testing.T) {
	f := Bucket{}
	f.Compare = lessCompare
	for i := 0; i < 10; i++ {
		a := A{i, i}
		f.Insert(&a)
	}
	/*for i:=0;i<f.Len();i++{
		log.Info("%#v",f.Data[i])
	}*/
	assert.Equal(t, 9, f.Data[9].(*A).id)
}

func TestBucket_InsertGreater(t *testing.T) {
	f := Bucket{}
	f.Compare = greaterCompare
	for i := 0; i < 10; i++ {
		a := A{i, i}
		f.Insert(&a)
	}
	/*for i:=0;i<f.Len();i++{
		log.Info("%#v",f.Data[i])
	}*/
	assert.Equal(t, 0, f.Data[9].(*A).id)
}

func TestBucket_InsertRepeat_Greater(t *testing.T) {
	f := Bucket{}
	f.Compare = greaterCompare
	for i := 0; i < 10; i++ {
		if i != 5 {
			a := A{i, i}
			f.Insert(&a)
		}
	}
	tmp := &A{5, 5}
	f.Insert(tmp)
	f.Insert(tmp)
	f.Insert(tmp)
	f.Insert(tmp)
	/*for i:=0;i<f.Len();i++{
		log.Info("%#v",f.Data[i])
	}*/
	assert.Equal(t, tmp, f.Data[4].(*A))
}

func TestBucket_Find(t *testing.T) {
	f := Bucket{}
	var tmp A
	f.Compare = greaterCompare
	for i := 0; i < 10; i++ {
		a := A{i, i}
		if i == 5 {
			tmp = a
		}
		f.Insert(&a)
	}
	exist, i := f.Find(&tmp)
	//fmt.Println(exist,i)
	assert.Equal(t, true, exist)
	assert.Equal(t, 4, i)
}

func TestBucket_FindNil(t *testing.T) {
	f := Bucket{}
	tmp := A{10, 10}
	f.Compare = greaterCompare
	for i := 0; i < 10; i++ {
		a := A{i, i}
		f.Insert(&a)
	}

	exist, i := f.Find(&tmp)
	//fmt.Println(exist,i)
	assert.Equal(t, false, exist)
	assert.Equal(t, -1, i)
}

func TestBucket_GetData(t *testing.T) {
	f := Bucket{}
	f.Compare = greaterCompare
	for i := 0; i < 10; i++ {
		a := A{i, i}
		f.Insert(&a)
	}
	obj, _ := f.GetData(5)

	assert.Equal(t, &A{4, 4}, obj)
}

func TestBucket_GetDataNil(t *testing.T) {
	f := Bucket{}
	f.Compare = greaterCompare
	for i := 0; i < 10; i++ {
		a := A{i, i}
		f.Insert(&a)
	}
	obj, err := f.GetData(10)
	assert.Equal(t, nil, obj)
	assert.Error(t, err, "not found data")
}

func TestBucket_Eraser_Nil(t *testing.T) {
	f := Bucket{}
	f.Compare = greaterCompare
	for i := 0; i < 10; i++ {
		a := A{i, i}
		f.Insert(&a)
	}
	obj := f.Eraser(&A{10, 10})
	assert.Equal(t, false, obj)
	//assert.Error(t,err,"not found data")
}

func TestBucket_LowerBound_Greater(t *testing.T) {
	f := Bucket{}
	f.Compare = greaterCompare
	for i := 0; i < 10; i++ {
		if i != 5 {
			a := A{i, i}
			f.Insert(&a)
		}
	}
	tmp := &A{5, 5}
	f.Insert(tmp)
	f.Insert(tmp)
	f.Insert(tmp)
	f.Insert(tmp)
	fmt.Println("greater result:")
	/*for i:=0;i<f.Len();i++{
		log.Info("%#v",f.Data[i])
	}*/
	assert.Equal(t, tmp, f.Data[4].(*A))

	value, sub := f.LowerBound(tmp)
	log.Info("%v,%d", value, sub)
	assert.Equal(t, 4, sub)
}

func TestBucket_UpperBound_Greater(t *testing.T) {
	f := Bucket{}
	f.Compare = greaterCompare
	for i := 0; i < 10; i++ {
		if i != 5 {
			a := A{i, i}
			f.Insert(&a)
		}
	}
	tmp := &A{5, 5}
	f.Insert(tmp)
	f.Insert(tmp)
	f.Insert(tmp)
	f.Insert(tmp)
	/*for i:=0;i<f.Len();i++{
		log.Info("%#v",f.Data[i])
	}*/
	assert.Equal(t, tmp, f.Data[4].(*A))

	value, sub := f.UpperBound(tmp)
	log.Info("%v,%d", value, sub)
	assert.Equal(t, 7, sub)
}

func TestBucket_LowerBound_Less(t *testing.T) {
	f := Bucket{}
	f.Compare = lessCompare
	for i := 0; i < 10; i++ {
		if i != 5 {
			a := A{i, i}
			f.Insert(&a)
		}
	}
	tmp := &A{5, 5}
	f.Insert(tmp)
	f.Insert(tmp)
	f.Insert(tmp)
	f.Insert(tmp)

	assert.Equal(t, tmp, f.Data[5].(*A))

	value, sub := f.LowerBound(tmp)
	log.Info("%v,%d", value, sub)
	assert.Equal(t, 5, sub)
}

func TestBucket_UpperBound_Less(t *testing.T) {
	f := Bucket{}
	f.Compare = lessCompare
	for i := 0; i < 10; i++ {
		if i != 5 {
			a := A{i, i}
			f.Insert(&a)
		}
	}
	tmp := &A{5, 5}
	f.Insert(tmp)
	f.Insert(tmp)
	f.Insert(tmp)
	f.Insert(tmp)

	assert.Equal(t, tmp, f.Data[5].(*A))

	value, sub := f.UpperBound(tmp)
	fmt.Println(value, sub)
	assert.Equal(t, 8, sub)
}
