package chain_plugin

import (
	"testing"
	"github.com/eosspark/eos-go/common"
)

func TestParamAndResult(t *testing.T) {
	params := GetCurrencyBalanceParams{
		Code:    common.Name(common.N("eosio.token")),
		Account: common.Name(common.N("eosio")),
		Symbol:  "SYS",
	}
}
