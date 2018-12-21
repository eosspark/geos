package unittests

import (
	"github.com/eosspark/eos-go/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuySell(t *testing.T){
	e := initEosioSystemTester()
	alice := common.N("alice1111111")
	assert.Equal(t, e.GetBalance(alice), CoreFromString("0.0000"))
}
