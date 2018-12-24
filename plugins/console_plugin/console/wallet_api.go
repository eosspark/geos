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

type walletapi struct {
	c   *Console
	log log.Logger
}

func newWalletapi(c *Console) *walletapi {
	w := &walletapi{
		c: c,
	}
	w.log = log.New("eosgo")
	w.log.SetHandler(log.TerminalHandler)
	return w
}

func (w *walletapi) CreateWallet(call otto.FunctionCall) (resonse otto.Value) {
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

func (w *walletapi) OpenWallet(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}
func (w *walletapi) ListWallets(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}
func (w *walletapi) ListKeys(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}

func (w *walletapi) GetPublicKeys(call otto.FunctionCall) (resonse otto.Value) {
	var resp []string
	err := DoHttpCall(&resp, common.WalletPublicKeys, nil)
	if err != nil {
		throwJSException(err)
	}
	v, _ := call.Otto.ToValue(resp)
	return v
}
func (w *walletapi) LockWallet(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}
func (w *walletapi) LockAllWallet(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}
func (w *walletapi) UnlockWallet(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}

func (w *walletapi) ImportKey(call otto.FunctionCall) (resonse otto.Value) {
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
func (w *walletapi) RemoveKey(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}
func (w *walletapi) CreateKeyByWallet(call otto.FunctionCall) (resonse otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}
func (w *walletapi) SignTransaction(call otto.FunctionCall) (resonse otto.Value) {
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
