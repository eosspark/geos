package unittests

import (
	"testing"
	"fmt"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"bytes"
	"io/ioutil"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/chain/types"
	. "github.com/eosspark/eos-go/common"
	. "github.com/eosspark/eos-go/exception/try"
	. "github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/log"
)

func TestGetBlockWithInvalidAbi(t *testing.T) {
	AssertAbi, err := ioutil.ReadFile("test_constracts/asserter.abi")
	assert.NoError(t, err)
	AssertWast, err := ioutil.ReadFile("test_constracts/asserter.wast")
	assert.NoError(t, err)

	Try(func() {
		tester := newValidatingTester(true, chain.SPECULATIVE)

		tester.ProduceBlocks(2, false)

		tester.CreateAccounts([]AccountName{N("asserter")}, false, true)
		tester.ProduceBlock(Milliseconds(DefaultConfig.BlockIntervalMs), 0)

		// setup contract and abi
		tester.SetCode(N("asserter"), AssertWast, nil)
		tester.SetAbi(N("asserter"), AssertAbi, nil)

		tester.ProduceBlocks(1, false)

		resolver := func(name AccountName) (r *abi_serializer.AbiSerializer) {
			Try(func() {
				accnt := entity.AccountObject{}
				accnt.Name = name
				if err := tester.Control.DataBase().Find("byName", accnt, &accnt); err != nil {
					EosThrow(&DatabaseException{}, err.Error())
				}
				var abi abi_serializer.AbiDef
				if abi_serializer.ToABI(accnt.Abi, &abi) {
					r = abi_serializer.NewAbiSerializer(&abi, tester.AbiSerializerMaxTime)
					return
				}
				r = nil
			}).FcRethrowExceptions(log.LvlError, "resolver failed at chain_plugin_tests::abi_invalid_type")
			return r
		}

		// abi should be resolved
		assert.NotNil(t, resolver(N("asserter")) != nil)

		prettyTrx := Variants{
			"actions": Variants{
				"account": "asserter",
				"name":    "procassert",
				"authorization": Variants{
					"actor":     "asserter",
					"permisson": DefaultConfig.ActiveName.String(),
				},
				"data": Variants{
					"condition": 1,
					"message":   "Should Not Assert!",
				},
			},
		}

		trx := types.SignedTransaction{}
		err := FromVariant(prettyTrx, &trx)
		assert.NoError(t, err)
		tester.SetTransactionHeaders(&trx.Transaction, tester.DefaultExpirationDelta, 0)
		priKey, chainId := tester.getPrivateKey(N("asserter"), "active"), tester.Control.GetChainId()
		trx.Sign(&priKey, &chainId)
		tester.PushTransaction(&trx, MaxTimePoint(), tester.DefaultBilledCpuTimeUs)
		tester.ProduceBlocks(1, false)

		// retrieve block num
		headNum := tester.Control.HeadBlockNum()
		headNumStr := fmt.Sprintf("%d", headNum)
		param := chain_plugin.GetBlockParams{headNumStr}
		plugin := chain_plugin.NewReadOnly(tester.Control, MaxMicroseconds())

		// block should be decoded successfully
		blockStr, err := json.Marshal(plugin.GetBlock(param))
		assert.NoError(t, err)

		// block should be decoded successfully
		assert.Equal(t, true, bytes.Contains(blockStr, []byte("procassert")))
		assert.Equal(t, true, bytes.Contains(blockStr, []byte("condition")))
		assert.Equal(t, true, bytes.Contains(blockStr, []byte("Should Not Assert!")))
		assert.Equal(t, true, bytes.Contains(blockStr, []byte("011253686f756c64204e6f742041737365727421"))) //action data

		// set an invalid abi (int8->xxxx)
		abi2 := AssertAbi
		pos := bytes.Index(abi2, []byte("int8"))
		assert.Equal(t, true, pos > 0)
		copy(abi2[pos:pos+4], []byte("xxxx"))
		tester.SetAbi(N("asserter"), abi2, nil)
		tester.ProduceBlocks(1, false)

		// resolving the invalid abi result in exception
		CheckThrow(t, func() { resolver(N("asserter")) }, &InvalidTypeInsideAbi{})

		// get the same block as string, results in decode failed(invalid abi) but not exception
		blockStr2, err := json.Marshal(plugin.GetBlock(param))
		assert.Equal(t, true, bytes.Contains(blockStr2, []byte("procassert")))
		assert.Equal(t, false, bytes.Contains(blockStr2, []byte("condition")))
		assert.Equal(t, false, bytes.Contains(blockStr2, []byte("Should Not Assert!")))
		assert.Equal(t, true, bytes.Contains(blockStr2, []byte("011253686f756c64204e6f742041737365727421"))) //action data

	}).Catch(func(e Exception) {
		t.Fatal(e.DetailMessage())
	}) // get_block_with_invalid_abi
}
