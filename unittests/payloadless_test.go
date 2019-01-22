package unittests

import (
	"testing"
	"github.com/eosspark/eos-go/common"
	"io/ioutil"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/chain"
	"github.com/stretchr/testify/assert"
	"fmt"
)

type PayloadlessTester struct {
	ValidatingTester
}

func NewPayloadlessTester() *PayloadlessTester {
	pt := &PayloadlessTester{}
	pt.DefaultExpirationDelta = 6
	pt.DefaultBilledCpuTimeUs = 2000
	pt.AbiSerializerMaxTime = 1000 * 1000
	pt.ChainTransactions = make(map[common.BlockIdType]types.TransactionReceipt)
	pt.LastProducedBlock = make(map[common.AccountName]common.BlockIdType)
	pt.VCfg = *newConfig(chain.SPECULATIVE)
	pt.ValidatingControl = chain.NewController(&pt.VCfg)
	pt.ValidatingControl.Startup()
	pt.init(true, chain.SPECULATIVE)
	return pt
}

func TestPayloadless (t *testing.T) {
	pt := NewPayloadlessTester()
	pt.CreateAccount(common.N("payloadless"),common.DefaultConfig.SystemAccountName,false,true)

	payloadless := common.N("payloadless")
	doit :=common.N("doit")

	wasm := "test_contracts/payloadless.wasm"
	abi := "test_contracts/payloadless.abi"
	code, _ := ioutil.ReadFile(wasm)
	abiCode,_ := ioutil.ReadFile(abi)
	pt.SetCode(payloadless,code,nil)
	pt.SetAbi(payloadless,abiCode,nil)
	data := common.Variants{}
	trace := pt.PushAction2(
		&payloadless,
		&doit,
		payloadless,
		&data,
		pt.DefaultExpirationDelta,
		0,
	)
	msg := trace.ActionTraces[0].Console
	assert.Equal(t,msg,"Im a payloadless action")
	pt.close()
}

// test GH#3916 - contract api action with no parameters fails when called from cleos
// abi_serializer was failing when action data was empty.
func TestAbiSerializer(t *testing.T) {
	pt := NewPayloadlessTester()
	pt.CreateAccount(common.N("payloadless"),common.DefaultConfig.SystemAccountName,false,true)

	payloadless := common.N("payloadless")
	doit :=common.N("doit")

	wasm := "test_contracts/payloadless.wasm"
	abi := "test_contracts/payloadless.abi"
	code, _ := ioutil.ReadFile(wasm)
	abiCode,_ := ioutil.ReadFile(abi)
	pt.SetCode(payloadless,code,nil)
	pt.SetAbi(payloadless,abiCode,nil)

	prettyTrx := &types.SignedTransaction{}
	prettyTrx.Actions = append(prettyTrx.Actions, &types.Action{
		Account: payloadless,
		Name :doit,
		Authorization: []types.PermissionLevel {
			{
				Actor: payloadless,
				Permission:common.DefaultConfig.ActiveName,
			},
		},
		Data:nil,
	})



	// from_variant is key to this test as abi_serializer was explicitly not allowing empty "data"
	//abi_serializer.FromVariant(&prettyTrx,trx,pt.GetResolver(),pt.AbiSerializerMaxTime)
	pt.SetTransactionHeaders(&prettyTrx.Transaction,pt.DefaultExpirationDelta,0)

	priKey,chainId := pt.getPrivateKey(payloadless,"active"),pt.Control.GetChainId()
	prettyTrx.Sign(&priKey,&chainId)

	trace := pt.PushTransaction(prettyTrx, common.MaxTimePoint(), pt.DefaultBilledCpuTimeUs)
	msg := trace.ActionTraces[len(trace.ActionTraces)-1].Console
	fmt.Println(msg)
	assert.Equal(t,msg,"Im a payloadless action")
	pt.close()
}
