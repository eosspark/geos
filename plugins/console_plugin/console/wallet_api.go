package console

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/robertkrimen/otto"
)

type walletApi struct {
	c   *Console
	log log.Logger
}

func newWalletApi(c *Console) *walletApi {
	w := &walletApi{
		c: c,
	}
	w.log = log.New("eosgo")
	w.log.SetHandler(log.TerminalHandler)
	return w
}

func (w *walletApi) CreateWallet(call otto.FunctionCall) (response otto.Value) {
	walletName, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	var resp string
	err = DoHttpCall(&resp, common.WalletCreate, walletName)
	if err != nil {
		return
	}

	//re := fmt.Sprintf(
	//	"Creating wallet: %s\n"+
	//		" Save password to use in the future to unlock this wallet.\n"+
	//		" Without password imported keys will not be retrievable.\n"+
	//		"%s", walletName, resp)

	resps, _ := call.Otto.Object("console")
	re := fmt.Sprintf("Creating wallet: %s", walletName)
	resps.Call("log", re, "\nSave password to use in the future to unlock this wallet.\nWithout password imported keys will not be retrievable.\n", resp)

	v, _ := call.Otto.ToValue("")
	return v
}

func (w *walletApi) OpenWallet(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}
func (w *walletApi) ListWallets(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}
func (w *walletApi) ListKeys(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}

func (w *walletApi) GetPublicKeys(call otto.FunctionCall) (resonse otto.Value) {
	var resp []string
	err := DoHttpCall(&resp, common.WalletPublicKeys, nil)
	if err != nil {
		throwJSException(err)
	}
	v, _ := call.Otto.ToValue(resp)
	return v
}
func (w *walletApi) LockWallet(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}
func (w *walletApi) LockAllWallet(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}
func (w *walletApi) UnlockWallet(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}

func (w *walletApi) ImportKey(call otto.FunctionCall) (resonse otto.Value) {
	walletName, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	walletKeyStr, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	fmt.Println(walletName, walletKeyStr)
	walletKey, err := ecc.NewPrivateKey(walletKeyStr)
	if err != nil {
		try.EosThrow(&exception.PrivateKeyTypeException{}, "Invalid private key: %s", walletKeyStr)
	}

	err = DoHttpCall(nil, common.WalletImportKey, common.Variants{"name": walletName, "key": walletKeyStr})
	if err != nil {
		throwJSException(err)
	}

	v, _ := call.Otto.ToValue(fmt.Sprintf("imported private key for: %s", walletKey.PublicKey().String()))
	return v
}
func (w *walletApi) RemoveKey(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}
func (w *walletApi) CreateKeyByWallet(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}
func (w *walletApi) SignTransaction(call otto.FunctionCall) (resonse otto.Value) {
	fmt.Println("sign transaction")

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

	fmt.Println(resp)

	v, _ := call.Otto.ToValue(resp)
	return v
}
