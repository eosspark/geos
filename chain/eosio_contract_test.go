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
	"testing"
)

func TestApplyEosioNewAccount(t *testing.T) {

	t.Run("", func(t *testing.T) {

		//get a control
		control := GetControllerInstance()
		blockTimeStamp := types.NewBlockTimeStamp(common.Now())
		control.StartBlock(blockTimeStamp, 0)

		//action for create a new account
		wif := "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
		privKey, _ := ecc.NewPrivateKey(wif)
		pubKey := privKey.PublicKey()

		creator := NewAccount{
			Creator: common.AccountName(common.N("eosio")),
			Name:    common.AccountName(common.N("xiaoyu")),
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
		act := types.Action{
			Account: common.AccountName(common.N("eosio")),
			Name:    common.ActionName(common.N("newaccount")),
			Data:    buffer,
			Authorization: []types.PermissionLevel{
				//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
				{Actor: common.AccountName(common.N("eosio")), Permission: common.PermissionName(common.N("active"))},
			},
		}

		//pack a singedTrx
		trxHeader := types.TransactionHeader{
			Expiration:       common.MaxTimePointSec(),
			RefBlockNum:      4,
			RefBlockPrefix:   3832731038,
			MaxNetUsageWords: 0,
			MaxCpuUsageMS:    0,
			DelaySec:         0,
		}

		trx := types.Transaction{
			TransactionHeader:     trxHeader,
			ContextFreeActions:    []*types.Action{},
			Actions:               []*types.Action{&act},
			TransactionExtensions: []*types.Extension{},
		}
		signedTrx := types.NewSignedTransaction(&trx, []ecc.Signature{}, []common.HexBytes{})
		privateKey, _ := ecc.NewRandomPrivateKey()
		chainIdType := common.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
		signedTrx.Sign(privateKey, &chainIdType)
		trxContext := NewTransactionContext(control, signedTrx, trx.ID(), common.Now())

		//pack a applycontext from control, trxContext and act
		a := NewApplyContext(control, trxContext, &act, 0)

		applyEosioNewaccount(a)
		isAccount := a.IsAccount(int64(common.AccountName(common.N("xiaoyu"))))
		assert.Equal(t, isAccount, true)

		control.Close()

	})

}

func TestApplyEosioSetcode(t *testing.T) {

	name := "../wasmgo/testdata_context/hello.wasm"
	t.Run("", func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		//get a control
		control := GetControllerInstance()
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
			Authorization: []types.PermissionLevel{
				//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
				{Actor: common.AccountName(common.N(account)), Permission: common.PermissionName(common.N("active"))},
			},
		}

		a := newApplyContext(control, &action)
		applyEosioSetcode(a)

		accountObject := entity.AccountObject{Name: action.Account}
		control.DB.Find("byName", accountObject, &accountObject)
		assert.Equal(t, accountObject.Code, common.HexBytes(code))

		control.Close()

	})

}

func TestApplyEosioUpdateauth(t *testing.T) {

	t.Run("", func(t *testing.T) {

		//get a control
		control := GetControllerInstance()
		blockTimeStamp := types.NewBlockTimeStamp(common.Now())
		control.StartBlock(blockTimeStamp, 0)

		account1 := "xiaoyu"
		createNewAccount(control, account1)

		account2 := "michael"
		createNewAccount(control, account2)

		updateAuth := UpdateAuth{
			Account:    common.AccountName(common.N(account1)),
			Permission: common.PermissionName(common.N("active")),
			Parent:     common.PermissionName(common.N("owner")),
			Auth: types.Authority{
				Threshold: 2,
				Accounts: []types.PermissionLevelWeight{{
					Permission: types.PermissionLevel{Actor: common.AccountName(common.N(account2)), Permission: common.PermissionName(common.N("active"))},
					Weight:     2,
				}},
			},
		}

		buffer, _ := rlp.EncodeToBytes(&updateAuth)
		action := types.Action{
			Account: common.AccountName(common.N(account1)),
			Name:    common.ActionName(common.N("updateauth")),
			Data:    buffer,
			Authorization: []types.PermissionLevel{
				//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
				{common.AccountName(common.N(account1)), common.PermissionName(common.N("active"))},
			},
		}

		a := newApplyContext(control, &action)
		applyEosioUpdateauth(a)

		authorization := a.Control.GetMutableAuthorizationManager()
		permission := authorization.FindPermission(&types.PermissionLevel{updateAuth.Account, updateAuth.Permission})
		assert.Equal(t, permission.Auth.Threshold, updateAuth.Auth.Threshold)

		// accountObject := entity.AccountObject{Name: action.Account}
		// control.DB.Find("byName", accountObject, &accountObject)
		// assert.Equal(t, accountObject.Abi, common.HexBytes(abi))

		control.Close()

	})

}

func TestApplyEosioLinkauthAndUnlinkauth(t *testing.T) {

	name := "../wasmgo/testdata_context/hello.wasm"
	t.Run("", func(t *testing.T) {
		code, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}

		//get a control
		control := GetControllerInstance()
		blockTimeStamp := types.NewBlockTimeStamp(common.Now())
		control.StartBlock(blockTimeStamp, 0)

		account := "xiaoyu"
		createNewAccount(control, account)

		accountCode := "hello"
		createNewAccount(control, accountCode)

		setCode := SetCode{
			Account:   common.AccountName(common.N(accountCode)),
			VmType:    0,
			VmVersion: 0,
			Code:      code,
		}
		buffer, _ := rlp.EncodeToBytes(&setCode)
		action := types.Action{
			Account: common.AccountName(common.N(accountCode)),
			Name:    common.ActionName(common.N("setcode")),
			Data:    buffer,
			Authorization: []types.PermissionLevel{
				//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
				{Actor: common.AccountName(common.N(accountCode)), Permission: common.PermissionName(common.N("active"))},
			},
		}

		a := newApplyContext(control, &action)
		applyEosioSetcode(a)

		linkAuth := LinkAuth{
			common.AccountName(common.N(account)),
			common.AccountName(common.N(accountCode)),
			common.ActionName(common.N("hi")),
			common.PermissionName(common.N("active")),
		}
		buffer, _ = rlp.EncodeToBytes(&linkAuth)
		action = types.Action{
			Account: common.AccountName(common.N(account)),
			Name:    common.ActionName(common.N("linkauth")),
			Data:    buffer,
			Authorization: []types.PermissionLevel{
				//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
				{Actor: common.AccountName(common.N(account)), Permission: common.PermissionName(common.N("active"))},
			},
		}

		a = newApplyContext(control, &action)
		applyEosioLinkauth(a)

		permissionLinkObject := entity.PermissionLinkObject{
			Account:     linkAuth.Account,
			Code:        linkAuth.Code,
			MessageType: linkAuth.Type}
		a.DB.Find("byActionName", permissionLinkObject, &permissionLinkObject)
		assert.Equal(t, permissionLinkObject.RequiredPermission, linkAuth.Requirement)

		unlikAuth := UnLinkAuth{
			common.AccountName(common.N(account)),
			common.AccountName(common.N(accountCode)),
			common.ActionName(common.N("hi")),
		}
		buffer, _ = rlp.EncodeToBytes(&unlikAuth)
		action = types.Action{
			Account: common.AccountName(common.N(account)),
			Name:    common.ActionName(common.N("linkauth")),
			Data:    buffer,
			Authorization: []types.PermissionLevel{
				//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
				{Actor: common.AccountName(common.N(account)), Permission: common.PermissionName(common.N("active"))},
			},
		}
		a = newApplyContext(control, &action)
		applyEosioUnlinkauth(a)
		permissionLinkObject = entity.PermissionLinkObject{
			Account:     unlikAuth.Account,
			Code:        unlikAuth.Code,
			MessageType: unlikAuth.Type}
		a.DB.Find("byActionName", permissionLinkObject, &permissionLinkObject)
		assert.Equal(t, permissionLinkObject.RequiredPermission, common.PermissionName(0))

		control.Close()

	})

}

func newApplyContext(control *Controller, action *types.Action) *ApplyContext {

	//pack a singedTrx
	trxHeader := types.TransactionHeader{
		Expiration:       common.MaxTimePointSec(),
		RefBlockNum:      4,
		RefBlockPrefix:   3832731038,
		MaxNetUsageWords: 0,
		MaxCpuUsageMS:    0,
		DelaySec:         0,
	}

	trx := types.Transaction{
		TransactionHeader:     trxHeader,
		ContextFreeActions:    []*types.Action{},
		Actions:               []*types.Action{action},
		TransactionExtensions: []*types.Extension{},
	}
	signedTrx := types.NewSignedTransaction(&trx, []ecc.Signature{}, []common.HexBytes{})
	privateKey, _ := ecc.NewRandomPrivateKey()
	chainIdType := common.ChainIdType(*crypto.NewSha256String("cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f"))
	signedTrx.Sign(privateKey, &chainIdType)
	trxContext := NewTransactionContext(control, signedTrx, trx.ID(), common.Now())

	//pack a applycontext from control, trxContext and act
	a := NewApplyContext(control, trxContext, action, 0)
	return a
}

func createNewAccount(control *Controller, name string) {

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

	act := types.Action{
		Account: common.AccountName(common.N("eosio")),
		Name:    common.ActionName(common.N("newaccount")),
		Data:    buffer,
		Authorization: []types.PermissionLevel{
			//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
			{Actor: common.AccountName(common.N("eosio")), Permission: common.PermissionName(common.N("active"))},
		},
	}

	a := newApplyContext(control, &act)

	//create new account
	applyEosioNewaccount(a)
}
