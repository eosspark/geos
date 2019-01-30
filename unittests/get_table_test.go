package unittests

import (
	. "github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math"
	"testing"
)

func TestGetScope(t *testing.T) {
	tester := newValidatingTester(true, SPECULATIVE)
	tester.ProduceBlocks(2, false)
	tester.CreateAccounts([]common.AccountName{eosioToken, eosioRam, eosioRamFee, eosioStake, eosioBpay, eosioVpay, eosioSaving, eosioName}, false, true)
	accs := []common.AccountName{common.N("inita"), common.N("initb"), common.N("initc"), common.N("initd")}
	tester.CreateAccounts(accs, false, true)
	tester.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)
	wasmName := "test_contracts/eosio.token.wasm"
	code, _ := ioutil.ReadFile(wasmName)
	tester.SetCode(eosioToken, code, nil)
	abiName := "test_contracts/eosio.token.abi"
	abi, _ := ioutil.ReadFile(abiName)
	tester.SetAbi(eosioToken, abi, nil)
	tester.ProduceBlocks(1, false)

	// create currency
	{
		data := common.Variants{
			"issuer":         eosio,
			"maximum_supply": fromString("1000000000.0000 SYS"),
		}
		act := common.N("create")
		tester.PushAction2(&eosioToken, &act, eosioToken, &data, tester.DefaultExpirationDelta, 0)
	}
	// issue
	for _, a := range accs {
		data := common.Variants{
			"to":       a,
			"quantity": fromString("999.0000 SYS"),
			"memo":     "",
		}
		act := common.N("issue")
		tester.PushAction2(&eosioToken, &act, eosio, &data, tester.DefaultExpirationDelta, 0)
	}
	tester.ProduceBlocks(1, false)

	// iterate over scope
	readOnly := chain_plugin.NewReadOnly(tester.Control, common.Microseconds(math.MaxInt32))
	params := chain_plugin.GetTableByScopeParams{Code: eosioToken, Table: common.N("accounts"), LowerBound: "inita", UpperBound: "", Limit: 10}
	result := readOnly.GetTableByScope(params)
	assert.Equal(t, int(4), len(result.Rows))
	assert.Equal(t, string(""), result.More)

	if len(result.Rows) >= 4 {
		assert.Equal(t, eosioToken, result.Rows[0].Code)
		assert.Equal(t, common.N("inita"), result.Rows[0].Scope)
		assert.Equal(t, common.N("accounts"), result.Rows[0].Table)
		assert.Equal(t, eosio, result.Rows[0].Payer)
		assert.Equal(t, uint32(1), result.Rows[0].Count)

		assert.Equal(t, common.N("initb"), result.Rows[1].Scope)
		assert.Equal(t, common.N("initc"), result.Rows[2].Scope)
		assert.Equal(t, common.N("initd"), result.Rows[3].Scope)
	}

	params.LowerBound = "initb"
	params.UpperBound = "initd"
	result = readOnly.GetTableByScope(params)
	assert.Equal(t, int(2), len(result.Rows))
	assert.Equal(t, string(""), result.More)
	if len(result.Rows) >= 2 {
		assert.Equal(t, common.N("initb"), result.Rows[0].Scope)
		assert.Equal(t, common.N("initc"), result.Rows[1].Scope)
	}

	params.Limit = 1
	result = readOnly.GetTableByScope(params)
	assert.Equal(t, int(1), len(result.Rows))
	assert.Equal(t, string("initc"), result.More)

	params.Table = common.N(common.S(0))
	result = readOnly.GetTableByScope(params)
	assert.Equal(t, int(1), len(result.Rows))
	assert.Equal(t, string("initc"), result.More)

	params.Table = common.N("invalid")
	result = readOnly.GetTableByScope(params)
	assert.Equal(t, int(0), len(result.Rows))
	assert.Equal(t, string(""), result.More)

	tester.close()
}