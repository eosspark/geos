package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewName(t *testing.T) {
	name := "eosio.system"
	val := StringToName(name)
	assert.Equal(t, val, uint64(6138663591228101920))
	fmt.Printf("%d\n", val)
	name2 := NameToString(6138663591228101920)
	fmt.Println(name2)
	assert.Equal(t, name2, name)
}
