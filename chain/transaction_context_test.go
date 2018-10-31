package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	"io/ioutil"
	"testing"
)

func TestTransactionContextTest(t *testing.T) {

	name := "../wasmgo/testdata_context/hello.wasm"
	t.Run("", func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		//set code
		control := GetControllerInstance()
		blockTimeStamp := common.NewBlockTimeStamp(common.Now())
		control.StartBlock(blockTimeStamp, 0)

		account := "hello"
		createNewAccount(control, account)

		setCode := setCode{
			Account:   common.AccountName(common.N(account)),
			VmType:    0,
			VmVersion: 0,
			Code:      code,
		}
		buffer, _ := rlp.EncodeToBytes(&setCode)
		action := types.Action{
			Account: common.AccountName(common.N(account)),
			Name:    common.ActionName(common.N("setcode")),
			Data:    buffer,
			Authorization: []types.PermissionLevel{
				//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
				{Actor: common.AccountName(common.N(account)), Permission: common.PermissionName(common.N("active"))},
			},
		}

		a := newApplyContext(control, &action)
		applyEosioSetcode(a)

		//execute contract hello.hi
		buffer, _ = rlp.EncodeToBytes(common.N("walker"))
		action = types.Action{
			Account: common.AccountName(common.N(account)),
			Name:    common.ActionName(common.N("hi")),
			Data:    buffer,
			Authorization: []types.PermissionLevel{
				//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
				{Actor: common.AccountName(common.N(account)), Permission: common.PermissionName(common.N("active"))},
			},
		}
		trxHeader := types.TransactionHeader{
			Expiration:       common.MaxTimePointSec(),
			RefBlockNum:      4,
			RefBlockPrefix:   3832731038,
			MaxNetUsageWords: 100000,
			MaxCpuUsageMS:    200,
			DelaySec:         0,
		}

		trx := types.Transaction{
			TransactionHeader:     trxHeader,
			ContextFreeActions:    []*types.Action{},
			Actions:               []*types.Action{&action},
			TransactionExtensions: []*types.Extension{},
		}
		signedTrx := types.NewSignedTransaction(&trx, []ecc.Signature{}, []common.HexBytes{})
		privateKey, _ := ecc.NewRandomPrivateKey()
		chainIdType := common.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
		signedTrx.Sign(privateKey, &chainIdType)

		metaTrx := types.NewTransactionMetadataBySignedTrx(signedTrx, common.CompressionNone)

		//var trace *types.TransactionTrace
		for i := 0; i < 10; i++ {
			trxContext := NewTransactionContext(control, signedTrx, trx.ID(), common.Now())
			trxContext.Deadline = common.Now() + common.TimePoint(100000)
			trxContext.ExplicitBilledCpuTime = true
			trxContext.BilledCpuTimeUs = 150000
			//trace = trxContext.Trace

			trxContext.InitForInputTrx(uint64(metaTrx.PackedTrx.GetUnprunableSize()),
				uint64(metaTrx.PackedTrx.GetPrunableSize()),
				uint32(len(signedTrx.Signatures)),
				true)

			trxContext.Delay = common.Seconds(int64(metaTrx.Trx.DelaySec)) // seconds
			trxContext.Exec()
			trxContext.Finalize()

			usage := entity.ResourceUsageObject{Owner: common.AccountName(common.N(account))}
			control.DB.Find("byOwner", usage, &usage)
			fmt.Println(i, ":", usage)

		}

		// accountObject := entity.AccountObject{Name: action.Account}
		// control.DB.Find("byName", accountObject, &accountObject)
		// assert.Equal(t, accountObject.Code, common.HexBytes(code))

		control.Close()
		control.Clean()

	})

}
