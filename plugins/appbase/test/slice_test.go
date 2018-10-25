package test

import (
	"testing"
	"fmt"
)

var slice_1 []string
var map_2 = map[int]string{1:"sf",2:"sf2"}

func TestSlice (t *testing.T) {

	//for k,v := range map_2 {
	//	fmt.Println(k)
	//	fmt.Println(v)
	//}
	//
	slice_1 = append(slice_1, "sf")
	slice_1 = append(slice_1, "sf2")
	for i := 0;i< len(slice_1);i++ {
		fmt.Println(slice_1[i])
	}
	////slice_1 = append(slice_1,slice_2...)
	//fmt.Println(slice_1)
	//var s []string
	//fmt.Print(s == nil)
}

func findPlugin (name string) {

}