package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestContract(t *testing.T) {

	name := "../wasmgo/testdata_context/eosio.token.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		control := GetControllerInstance()
		blockTimeStamp := common.NewBlockTimeStamp(common.Now())
		control.StartBlock(blockTimeStamp, 0)

		eosioToken := "eosio.token"
		account1 := "testapi1"
		account2 := "testapi2"

		CreateNewAccount(control, eosioToken)
		CreateNewAccount(control, account1)
		CreateNewAccount(control, account2)

		SetCode(control, eosioToken, code)

		control.Close()

	})

}

func SetCode(control *Controller, account string, code []byte) {

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
			{Actor: common.AccountName(common.N(account)), Permission: common.PermissionName(common.N("active"))},
		},
	}

	wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	privateKey, _ := ecc.NewPrivateKey(wif)

	trx := newTransaction(control, &action, privateKey)
	pushTransaction(control, trx)
}

func CreateNewAccount(control *Controller, name string) {

	//action for create a new account
	wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	privKey, _ := ecc.NewPrivateKey(wif)
	pubKey := privKey.PublicKey()

	creator := newAccount{
		Creator: common.AccountName(common.N("eosio")),
		Name:    common.AccountName(common.N(name)),
		Owner: types.Authority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{Key: pubKey, Weight: 1}},
		},
		Active: types.Authority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{Key: pubKey, Weight: 1}},
		},
	}

	buffer, _ := rlp.EncodeToBytes(&creator)

	action := types.Action{
		Account: common.AccountName(common.N("eosio")),
		Name:    common.ActionName(common.N("newaccount")),
		Data:    buffer,
		Authorization: []types.PermissionLevel{
			{Actor: common.AccountName(common.N("eosio")), Permission: common.PermissionName(common.N("active"))},
		},
	}

	privateKey, _ := ecc.NewPrivateKey(wif)

	trx := newTransaction(control, &action, privateKey)
	pushTransaction(control, trx)

}

func pushTransaction(control *Controller, trx *types.TransactionMetadata) {
	control.PushTransaction(*trx, common.TimePoint(common.MaxMicroseconds()), 1000, true)
}

func newTransaction(control *Controller, action *types.Action, privateKey *ecc.PrivateKey) *types.TransactionMetadata {

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
		Actions:               []*types.Action{action},
		TransactionExtensions: []*types.Extension{},
		RecoveryCache:         make(map[ecc.Signature]types.CachedPubKey),
	}
	signedTrx := types.NewSignedTransaction(&trx, []ecc.Signature{}, []common.HexBytes{})
	//privateKey, _ := ecc.NewRandomPrivateKey()
	//chainIdType := common.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
	chainIdType := control.GetChainId()
	signedTrx.Sign(privateKey, &chainIdType)

	metaTrx := types.NewTransactionMetadataBySignedTrx(signedTrx, common.CompressionNone)

	return metaTrx
}

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
		for i := 0; i < 100; i++ {
			trxContext := NewTransactionContext(control, signedTrx, trx.ID(), common.Now())
			trxContext.Deadline = common.Now() + common.TimePoint(200000)
			trxContext.ExplicitBilledCpuTime = false
			trxContext.BilledCpuTimeUs = 0
			//trace = trxContext.Trace

			trxContext.InitForInputTrx(uint64(metaTrx.PackedTrx.GetUnprunableSize()),
				uint64(metaTrx.PackedTrx.GetPrunableSize()),
				uint32(len(signedTrx.Signatures)),
				true)

			trxContext.Delay = common.Seconds(int64(metaTrx.Trx.DelaySec)) // seconds
			trxContext.Exec()
			trxContext.Finalize()

			//usage := entity.ResourceUsageObject{Owner: common.AccountName(common.N(account))}
			//control.DB.Find("byOwner", usage, &usage)
			//fmt.Println(i, ":", usage)
			//fmt.Println("No.", i)

		}

		// accountObject := entity.AccountObject{Name: action.Account}
		// control.DB.Find("byName", accountObject, &accountObject)
		// assert.Equal(t, accountObject.Code, common.HexBytes(code))

		control.Close()

	})

}
