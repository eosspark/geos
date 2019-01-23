package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func NewActionData(action interface{}) []byte {
	bytes, _ := rlp.EncodeToBytes(action)
	return bytes
}

//type eosiotokenAccount struct {
//	balance common.Asset
//}

func TestContract(t *testing.T) {

	name := "../wasmgo/testdata_context/eosio.token.wasm"
	t.Run(filepath.Base(name), func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		control := NewController(NewConfig())
		blockTimeStamp := types.NewBlockTimeStamp(common.Now())
		control.StartBlock(blockTimeStamp, 0)

		eosioToken := "eosio.token"
		account1 := "testapi1"
		account2 := "testapi2"

		CreateNewAccount(control, eosioToken)
		CreateNewAccount(control, account1)
		CreateNewAccount(control, account2)

		SCode(control, eosioToken, code)

		createToken(control, account1, 1000000000, "BTC")
		issueToken(control, account1, account1, 20000, "BTC", "issue")
		issueTransfer(control, account1, account2, 10000, "BTC", "transfer")
		issueToken(control, account1, account2, 20000, "BTC", "issue")
		//
		tab := entity.TableIdObject{Code: common.AccountName(common.N(eosioToken)),
			Scope: common.ScopeName(common.N(account2)),
			Table: common.TableName(common.N("accounts")),
		}
		control.DB.Find("byCodeScopeTable", tab, &tab)
		obj := entity.KeyValueObject{
			TId:        tab.ID,
			PrimaryKey: uint64(4412482),
		}
		control.DB.Find("byScopePrimary", obj, &obj)
		asset := common.Asset{}
		rlp.DecodeBytes(obj.Value, &asset)
		assert.Equal(t, asset.Amount, int64(30000))

		control.Close()
	})

}

func issueTransfer(control *Controller, from string, to string, amount int64, symbol string, memo string) {

	action := NewTransfer(common.AccountName(common.N(from)), common.AccountName(common.N(to)), common.Asset{amount, common.Symbol{4, symbol}}, memo)

	wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	privateKey, _ := ecc.NewPrivateKey(wif)

	trx := newTransaction(control, action, privateKey)
	pushTransaction(control, trx)

}

func issueToken(control *Controller, issuer string, to string, amount int64, symbol string, memo string) {

	action := NewIssue(common.AccountName(common.N(issuer)), common.AccountName(common.N(to)), common.Asset{amount, common.Symbol{4, symbol}}, memo)

	wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	privateKey, _ := ecc.NewPrivateKey(wif)

	trx := newTransaction(control, action, privateKey)
	pushTransaction(control, trx)

}

func createToken(control *Controller, issuer string, amount int64, symbol string) {

	action := NewCreate(common.AccountName(common.N(issuer)), common.Asset{amount, common.Symbol{4, symbol}})

	wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
	privateKey, _ := ecc.NewPrivateKey(wif)

	trx := newTransaction(control, action, privateKey)
	pushTransaction(control, trx)

}

func NewTransfer(from, to common.AccountName, quantity common.Asset, memo string) *types.Action {
	return &types.Action{
		Account: common.AccountName(common.N("eosio.token")),
		Name:    common.ActionName(common.N("transfer")),
		Authorization: []common.PermissionLevel{
			{Actor: from, Permission: common.PermissionName(common.N("active"))},
		},
		Data: NewActionData(Transfer{
			From:     from,
			To:       to,
			Quantity: quantity,
			Memo:     memo,
		}),
	}
}

// Transfer represents the `transfer` struct on `eosio.token` contract.
type Transfer struct {
	From     common.AccountName `json:"from"`
	To       common.AccountName `json:"to"`
	Quantity common.Asset       `json:"quantity"`
	Memo     string             `json:"memo"`
}

func NewIssue(issuer common.AccountName, to common.AccountName, quantity common.Asset, memo string) *types.Action {
	return &types.Action{
		Account: common.AccountName(common.N("eosio.token")),
		Name:    common.ActionName(common.N("issue")),
		Authorization: []common.PermissionLevel{
			{Actor: issuer, Permission: common.PermissionName(common.N("active"))},
		},
		Data: NewActionData(Issue{
			To:       to,
			Quantity: quantity,
			Memo:     memo,
		}),
	}
}

// Issue represents the `issue` struct on the `eosio.token` contract.
type Issue struct {
	To       common.AccountName `json:"to"`
	Quantity common.Asset       `json:"quantity"`
	Memo     string             `json:"memo"`
}

func NewCreate(issuer common.AccountName, maxSupply common.Asset) *types.Action {
	return &types.Action{
		Account: common.AccountName(common.N("eosio.token")),
		Name:    common.ActionName(common.N("create")),
		Authorization: []common.PermissionLevel{
			{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
		},
		Data: NewActionData(Create{
			Issuer:        issuer,
			MaximumSupply: maxSupply,
		}),
	}
}

// Create represents the `create` struct on the `eosio.token` contract.
type Create struct {
	Issuer        common.AccountName `json:"issuer"`
	MaximumSupply common.Asset       `json:"maximum_supply"`
}

func SCode(control *Controller, account string, code []byte) {

	a := SetCode{
		Account:   common.AccountName(common.N(account)),
		VmType:    0,
		VmVersion: 0,
		Code:      code,
	}
	buffer, _ := rlp.EncodeToBytes(&a)
	action := types.Action{
		Account: common.AccountName(common.N("eosio")),
		Name:    common.ActionName(common.N("setcode")),
		Data:    buffer,
		Authorization: []common.PermissionLevel{
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

	creator := NewAccount{
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
		Authorization: []common.PermissionLevel{
			{Actor: common.AccountName(common.N("eosio")), Permission: common.PermissionName(common.N("active"))},
		},
	}

	privateKey, _ := ecc.NewPrivateKey(wif)

	trx := newTransaction(control, &action, privateKey)
	pushTransaction(control, trx)

}

func pushTransaction(control *Controller, trx *types.TransactionMetadata) {
	control.PushTransaction(trx, common.TimePoint(common.MaxMicroseconds()), 0)
}

func newTransaction(control *Controller, action *types.Action, privateKey *ecc.PrivateKey) *types.TransactionMetadata {

	trxHeader := types.TransactionHeader{
		Expiration: common.NewTimePointSecTp(control.PendingBlockTime()).AddSec(60),
		// RefBlockNum:      4,
		// RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 100000,
		MaxCpuUsageMS:    200,
		DelaySec:         0,
	}

	trx := types.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*types.Action{},
		Actions:               []*types.Action{action},
		TransactionExtensions: []*types.Extension{},
		//RecoveryCache:         make(map[ecc.Signature]types.CachedPubKey),
	}
	headBlockId := control.HeadBlockId()
	trx.SetReferenceBlock(&headBlockId)
	signedTrx := types.NewSignedTransaction(&trx, []ecc.Signature{}, []common.HexBytes{})
	//privateKey, _ := ecc.NewRandomPrivateKey()
	//chainIdType := common.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
	chainIdType := control.GetChainId()
	signedTrx.Sign(privateKey, &chainIdType)

	metaTrx := types.NewTransactionMetadataBySignedTrx(signedTrx, types.CompressionNone)

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
		control := NewController(NewConfig())
		blockTimeStamp := types.NewBlockTimeStamp(common.Now())
		control.StartBlock(blockTimeStamp, 0)

		account := "hello"
		createNewAccount(control, account)

		setCode := SetCode{
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
			Authorization: []common.PermissionLevel{
				//common.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
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
			Authorization: []common.PermissionLevel{
				//common.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
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

		metaTrx := types.NewTransactionMetadataBySignedTrx(signedTrx, types.CompressionNone)

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
