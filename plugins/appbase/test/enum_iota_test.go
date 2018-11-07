package test

import (
	"fmt"
	"testing"
)

// 实现枚举例子

type State1 int

// iota 初始化后会自动递增
const (
	Running    State1 = iota // value --> 0
	Stopped1                 // value --> 1
	Rebooting               // value --> 2
	Terminated              // value --> 3
)

func (this State1) String() string {
	switch this {
	case Running:
		return "Running"
	case Stopped1:
		return "Stopped"
	default:
		return "Unknow"
	}
}

func Test_Enum_Iota(t *testing.T) {
	state := Stopped1
	fmt.Println("state", state)
}

// 输出 state Running
