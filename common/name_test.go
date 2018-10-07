package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewName(t *testing.T) {
	name := "eosio.system"
	val := S(name)
	assert.Equal(t, val, uint64(6138663591228101920))
	//fmt.Printf("%d\n", val)
	name2 := N(6138663591228101920)
	//fmt.Println(name2)
	assert.Equal(t, name2, name)
}

func TestNameStr(t *testing.T) {
	name := "eosio.systemabdxs"
	val := S(name)
	fmt.Printf("%d\n", val)

}

func TestNameSuffix(t *testing.T) {
	name := S("eosio.token")
	check := S("token")
	test := NameSuffix(name)
	//fmt.Println(N(test))
	assert.Equal(t, test, check)
}
