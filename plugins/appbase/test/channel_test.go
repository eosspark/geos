package test

import (
	"testing"
	"fmt"
	. "github.com/eosspark/eos-go/plugins/appbase/app/include"
)

type  User struct {
	name string
}
func (u *User) hello(t int) {
	fmt.Println("hello",t)
}
func (u *User) sayHi(t int)  {
	fmt.Println("Hi!",t)
}

func (u *User) hello1(t int) {
	fmt.Println("hello1",t)
}
func (u *User) sayHi1(t int)  {
	fmt.Println("Hi!1",t)
}

func (u *User) sayHi2(t int)  {
	fmt.Println("Hi!2",t)
}

var  preAcceptedBlock *Signal;
var  acceptedBlockHeader *Signal;
var  acceptedBlock *Signal;
var  irreversibleBlock *Signal;
var  acceptedTransaction *Signal;
var  appliedTransaction *Signal;
var  acceptedConfirmation *Signal;
var  badAlloc *Signal;


func Test_Channel(t *testing.T) {
	user :=new(User)
	user.name="user"
	preAcceptedBlock = GetChannel()
	preAcceptedBlock.Subscribe(user.hello)
	preAcceptedBlock.Subscribe(user.sayHi)
	preAcceptedBlock.Subscribe(user.hello1)
	preAcceptedBlock.Subscribe(user.sayHi1)
	preAcceptedBlock.Subscribe(user.sayHi2)
	preAcceptedBlock.Publish(1)


}
