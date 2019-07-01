package console

import (
	"fmt"

	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"

	"github.com/robertkrimen/otto"
)

type walletApi struct {
	c *Console
}

func newWalletApi(c *Console) *walletApi {
	return &walletApi{c: c}
}

//Create creates a new wallet locally
func (w *walletApi) Create(call otto.FunctionCall) (response otto.Value) {
	walletName, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	if len(walletName) == 0 {
		walletName = "default"
	}

	var resp string
	err = DoHttpCall(&resp, common.WalletCreate, walletName)
	if err != nil {
		return getJsResult(call, err.Error())
	}

	fmt.Println("Creating wallet: ", walletName)
	fmt.Println("Save password to use in the future to unlock this wallet.")
	fmt.Println("Without password imported keys will not be retrievable.")
	fmt.Println("Password: ", resp)
	return
}

//Open opens an existing wallet
func (w *walletApi) Open(call otto.FunctionCall) (response otto.Value) {
	walletName, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	err = DoHttpCall(nil, common.WalletOpen, walletName)
	if err != nil {
		return
	}
	result := fmt.Sprintf("Opened: %s", walletName)
	return getJsResult(call, result)
}

//List lists wallets
func (w *walletApi) List(call otto.FunctionCall) (response otto.Value) {
	var resp []string
	err := DoHttpCall(&resp, common.WalletList, nil)
	if err != nil {
		throwJSException(err)
	}
	return getJsResult(call, resp)
}

//PublicKeys lists of public keys from all unlocked wallets
func (w *walletApi) PublicKeys(call otto.FunctionCall) (response otto.Value) {
	var resp []string
	err := DoHttpCall(&resp, common.WalletPublicKeys, nil)
	if err != nil {
		throwJSException(err)
	}
	return getJsResult(call, resp)
}

//PrivateKeys lists of private keys from an unlocked wallet in wif or PVT_R1 format.
func (w *walletApi) PrivateKeys(call otto.FunctionCall) (response otto.Value) {
	walletName, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	password, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	vs := common.Variants{"name": walletName, "password": password}
	var resp map[string]string
	err = DoHttpCall(&resp, common.WalletListKeys, vs)
	if err != nil {
		throwJSException(err)
	}
	return getJsResult(call, resp)
}

//Lock locks wallet
func (w *walletApi) Lock(call otto.FunctionCall) (response otto.Value) {
	walletName, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	err = DoHttpCall(nil, common.WalletLock, walletName)
	if err != nil {
		return
	}
	result := fmt.Sprintf("Locked: %s", walletName)
	return getJsResult(call, result)
}

//LockAll locks all wallets
func (w *walletApi) LockAll(call otto.FunctionCall) (response otto.Value) {
	err := DoHttpCall(nil, common.WalletLockAll, nil)
	if err != nil {
		return
	}
	result := fmt.Sprintf("Locked All Wallets")
	return getJsResult(call, result)
}

//Unlock unlocks wallet
func (w *walletApi) Unlock(call otto.FunctionCall) (response otto.Value) {
	walletName, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	password, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	vs := common.Variants{"name": walletName, "password": password}

	err = DoHttpCall(nil, common.WalletUnlock, vs)
	if err != nil {
		return
	}
	result := fmt.Sprintf("Unlocked: %s", walletName)
	return getJsResult(call, result)
}

//ImportKey imports private key into wallet
func (w *walletApi) ImportKey(call otto.FunctionCall) (response otto.Value) {
	walletName, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	walletKeyStr, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	walletKey, err := ecc.NewPrivateKey(walletKeyStr)
	if err != nil {
		try.EosThrow(&exception.PrivateKeyTypeException{}, "Invalid private key: %s", walletKeyStr)
	}

	err = DoHttpCall(nil, common.WalletImportKey, common.Variants{"name": walletName, "key": walletKeyStr})
	if err != nil {
		throwJSException(err)
	}
	fmt.Println("imported private key for: ", walletKey.PublicKey().String())
	return otto.UndefinedValue()
	//re := fmt.Sprintf("imported private key for: %s", walletKey.PublicKey().String())
	//return getJsResult(call, re)
}

//RemoveKey removes keys from wallet wallet.RemoveKey('walletName','wallet_pw','pubkey')
func (w *walletApi) RemoveKey(call otto.FunctionCall) (response otto.Value) {
	walletName, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	walletPw, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	walletRmKeyStr, err := call.Argument(2).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	_, err = ecc.NewPublicKey(walletRmKeyStr)
	if err != nil {
		return throwJSException(fmt.Sprintf("Invalid public key: %s", walletRmKeyStr))
	}
	vs := common.Variants{
		"name":     walletName,
		"password": walletPw,
		"key":      walletRmKeyStr,
	}
	err = DoHttpCall(nil, common.WalletRemoveKey, vs)
	if err != nil {
		throwJSException(err)
	}

	result := fmt.Sprintf("removed private key for: %s", walletRmKeyStr)
	return getJsResult(call, result)
}

//CreateKey creates private key within wallet, wallet.CreateKey('walletName','k1')
func (w *walletApi) CreateKey(call otto.FunctionCall) (response otto.Value) {
	walletName, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	walletCreateKeyType, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	var resp string
	err = DoHttpCall(&resp, common.WalletCreateKey, common.Variants{"name": walletName, "keyType": walletCreateKeyType})
	if err != nil {
		throwJSException(err)
	}
	result := fmt.Sprintf("Created new private key with a public key of: %s", resp)
	return getJsResult(call, result)
}

//SignTransaction signs a transaction
func (w *walletApi) SignTransaction(call otto.FunctionCall) (response otto.Value) {
	trxJsonToSign, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	strPrivateKey, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	strChainID, err := call.Argument(2).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	var resp types.SignedTransaction
	err = DoHttpCall(&resp, common.WalletSignTrx, []interface{}{
		trxJsonToSign,
		strPrivateKey,
		strChainID,
	})
	if err != nil {
		return throwJSException("signedTransactoin is err")
	}

	return getJsResult(call, resp)
}
