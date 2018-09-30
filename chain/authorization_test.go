package chain

import (
	"testing"
	"fmt"
)

func tType (t interface{}){
	b, ok := t.(aa)
	fmt.Println(ok," ",b)
}
type aa struct {
	a string
	b int
}

type bb struct {
	a string
	b string
}

func Test_reflect(t *testing.T) {
	var b interface{}
	b = bb{"1","2"}

	switch v := b.(type){
	case aa:
		fmt.Println("aa",b,v)
	case bb:
		fmt.Println("bb",b,v)
	default:
		fmt.Println("other")
	}
}
