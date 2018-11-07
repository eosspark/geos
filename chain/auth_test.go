package chain

import (
	"testing"
)

func initializeAuth() *AuthorizationManager{
	control := GetControllerInstance()
	am := control.Authorization
	return am
}

func TestMissingSigs(t *testing.T){
	am := initializeAuth()
	createNewAccount(am.control, "Alice")
	//BaseTester{}.ProduceBlock()



}