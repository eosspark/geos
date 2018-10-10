package test

import (
	"fmt"
	"testing"
)

// 实现枚举例子

type State int

// iota 初始化后会自动递增
const (
	Running State = iota // value --> 0
	Stopped              // value --> 1
	Rebooting            // value --> 2
	Terminated           // value --> 3
)

func (this State) String() string {
	switch this {
	case Running:
		return "Running"
	case Stopped:
		return "Stopped"
	default:
		return "Unknow"
	}
}

func Test_Enum_Iota(t *testing.T) {
	state := Stopped
	fmt.Println("state", state)
}
// 输出 state Running