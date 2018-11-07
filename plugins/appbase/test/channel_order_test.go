package test

import (
	"testing"
	"github.com/eosspark/eos-go/plugins/appbase/app/include"
	"fmt"
)

type  Student struct {
	name string
}
func (u *Student) hello(t int) {
	fmt.Println("hello,5",t)
}
func (u *Student) sayHi(t int)  {
	fmt.Println("Hi! 8",t)
}

func (u *Student) hello1(t int) {
	fmt.Println("hello1,3",t)
}
func (u *Student) sayHi1(t int)  {
	fmt.Println("Hi!1,7",t)
}

func (u *Student) sayHi2(t int)  {
	fmt.Println("Hi! 9",t)
}

func Test_Channel_Order(t *testing.T){
	user := new(Student)
	user.name ="sf"
	test := include.GetChannel2()
	test.Connect2(user.sayHi,8)
	test.Connect2(user.hello,5)
	test.Connect2(user.hello1,3)
	test.Connect2(user.sayHi2,9)
	test.Connect2(user.sayHi1,7)
	for _,v := range test.Channels {
		fmt.Println(v.Level)
		fmt.Println(v.Channel)
	}

	test.Emit2(1)



}