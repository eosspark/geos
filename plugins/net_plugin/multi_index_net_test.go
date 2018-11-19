package net_plugin

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"testing"
)

type A struct {
	id  int
	num int
}

func lessCompare(first common.ElementObject, second common.ElementObject) int {
	if first.(*A).id < second.(*A).id {
		return -1
	} else if first.(*A).id == second.(*A).id {
		return 0
	} else {
		return 1
	}
}

func greaterCompare(first common.ElementObject, second common.ElementObject) int {
	if first.(*A).id > second.(*A).id {
		return -1
	} else if first.(*A).id == second.(*A).id {
		return 0
	} else {
		return 1
	}
}

func (a *A) ElementObject() {}

func initRepeatData() *common.Bucket {
	var tmp *A
	f := common.Bucket{}
	f.Compare = lessCompare
	for i := 0; i < 10; i++ {
		a := A{i, i}
		f.Insert(&a)
		if i == 5 {
			tmp = &a
		}
	}

	for i := 0; i < 4; i++ {
		f.Insert(tmp)
	}
	for i := 0; i < f.Len(); i++ {
		fmt.Print(f.Data[i])
	}
	return &f
}

func TestIndexNet_LowerBound(t *testing.T) {

}
