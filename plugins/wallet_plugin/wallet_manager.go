package wallet_plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"os"
	"strings"
	"time"
)

var (
	ErrWalletNotAvaliable = errors.New("You don't have any wallet")
	ErrWalletNotUnlocked  = errors.New("You don't have any unlocked wallet!")
)

const (
	fileExt        string = ".wallet"
	defaultKeyType string = "K1"
	passwordPrefix string = "pw"
	tstampMax             = 3600 * time.Second
)

type WalletManager struct {
	timeOut     time.Duration //senconds max //how long to wait before calling lock_all()
	timeOutTime time.Time     // when to call lock_all()
	dir         string
	lockPath    string
	Self        *WalletPlugin
	log         log.Logger

	Wallets map[string]*SoftWallet
}

func walletManager() *WalletManager {
	manager := &WalletManager{
		timeOut:  tstampMax,
		dir:      ".",
		lockPath: "./wallet.lock",
		Wallets:  make(map[string]*SoftWallet),
	}

	manager.log = log.New("wallet_plugin")
	manager.log.SetHandler(log.TerminalHandler)
	manager.log.SetHandler(log.DiscardHandler())
	return manager
}

func checkNum(r rune) bool {
	if r >= '0' && r <= '9' {
		return true
	}
	return false
}
func validFileName(name string) bool {
	if len(name) == 0 {
		return false
	}
	if strings.Contains(".", name) || strings.Contains("_", name) || strings.Contains("-", name) {
		return false
	}
	if strings.IndexFunc(name, checkNum) != -1 {
		return false
	}
	return true
}

func genPassword() string {
	key, _ := ecc.NewRandomPrivateKey()
	return passwordPrefix + key.String()
}

func (wm *WalletManager) Create(name string) string {
	wm.log.Debug("wallet creating")
	wm.checkTimeout()

	EosAssert(validFileName(name), &WalletException{}, "Invalid filename, path not allowed in wallet name %s", name)

	os.MkdirAll(wm.dir, os.ModePerm)
	file, err := os.Open(wm.dir)
	defer file.Close()
	if err != nil {
		EosThrow(&WalletException{}, "Invalid dir: %s", wm.dir)
	}
	allWallets, err := file.Readdirnames(-1)
	if err != nil {
		EosThrow(&WalletException{}, "Invalid dir: %s", wm.dir)
	}

	walletName := name + fileExt
	for _, f := range allWallets {
		if f == walletName {
			EosThrow(&WalletExistException{}, "Wallet with name: %s already exists at %s", name, wm.dir)
		}
	}

	password := genPassword()

	var wallet SoftWallet
	wallet.SetPassword(password)
	walletFileName := fmt.Sprintf("%s/%s%s", wm.dir, name, fileExt)
	wallet.SetWalletFilename(walletFileName)
	wallet.Unlock(password)
	wallet.Lock()
	wallet.Unlock(password)
	// Explicitly save the wallet file here, to ensure it now exists.
	wallet.SaveWalletFile()

	// If we have name in our map then remove it since we want the emplace below to replace.
	// This can happen if the wallet file is removed while eos-walletd is running.
	if _, ok := wm.Wallets[name]; ok {
		delete(wm.Wallets, name)
	}
	wm.Wallets[name] = &wallet
	return password
}

func (wm *WalletManager) Open(name string) {
	wm.checkTimeout()
	wm.log.Debug("Opening wallet :   wallet name: %s", name)
	EosAssert(validFileName(name), &WalletException{}, "Invalid filename, path not allowed in wallet name %s", name)

	var wallet SoftWallet
	walletFileName := fmt.Sprintf("%s/%s%s", wm.dir, name, fileExt)
	wallet.SetWalletFilename(walletFileName)
	if !wallet.LoadWalletFile() {
		EosThrow(&WalletNonexistentException{}, "Unable to open file: %s", walletFileName)
	}
	// If we have name in our map then remove it since we want the emplace below to replace.
	// This can happen if the wallet file is added while eos-walletd is running.
	if _, ok := wm.Wallets[name]; ok {
		delete(wm.Wallets, name)
	}
	wm.Wallets[name] = &wallet
}

func (wm *WalletManager) ListWallets() []string {
	var result []string
	for name, wallet := range wm.Wallets {
		if wallet.IsLocked() {
			result = append(result, name)
		} else {
			result = append(result, name+"*")
		}
	}
	return result
}

type RespKeys map[ecc.PublicKey]ecc.PrivateKey

func (k RespKeys) MarshalJSON() ([]byte, error) {
	out := make(map[string]string, len(k))
	for pub, pri := range k {
		out[pub.String()] = pri.String()
	}
	return json.Marshal(out)
}

func (k *RespKeys) UnmarshalJSON(v []byte) (err error) {
	out := make(map[string]string)
	err = json.Unmarshal(v, &out)
	if err != nil {
		return err
	}
	keyMap := make(RespKeys, len(out))
	for pubStr, priStr := range out {
		priKey, _ := ecc.NewPrivateKey(priStr)
		pubKey, _ := ecc.NewPublicKey(pubStr)
		keyMap[pubKey] = *priKey
	}
	return nil
}

type ListKeysParams struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (wm *WalletManager) ListKeys(name, password string) RespKeys {
	wm.checkTimeout()
	//wm.log.Debug("all wallet: %v", wm.Wallets)

	wallet, ok := wm.Wallets[name]
	if !ok {
		EosThrow(&WalletNonexistentException{}, "Wallet not found: %s", name)
	}
	if wallet.IsLocked() {
		EosThrow(&WalletLockedException{}, "Wallet is locked: %s", name)
	}
	wallet.CheckPassword(password)

	return wallet.Keys
}

func (wm *WalletManager) GetPublicKeys() (re []string) {
	EosAssert(len(wm.Wallets) != 0, &WalletNotAvailableException{}, "You don't have any wallet!")
	isAllWalletLocked := true
	for name, wallet := range wm.Wallets {
		if !wallet.IsLocked() {
			isAllWalletLocked = false
			wm.log.Debug("wallet: %s is unlocked\n", name)
			for pubkey, _ := range wallet.Keys {
				re = append(re, pubkey.String())
			}
		}
	}

	EosAssert(!isAllWalletLocked, &WalletLockedException{}, "You don't have any unlocked wallet!")
	return re
}

func (wm *WalletManager) LockAllwallets() {
	wm.lockAll()
}

func (wm *WalletManager) lockAll() {
	// no call to check_timeout since we are locking all anyway
	for _, wallet := range wm.Wallets {
		if !wallet.IsLocked() {
			wallet.Lock()
		}
	}
}

func (wm *WalletManager) Lock(name string) {
	if _, ok := wm.Wallets[name]; !ok {
		EosThrow(&WalletNonexistentException{}, "Wallet not found:%s", name)
	}
	wallet := wm.Wallets[name]
	wallet.Lock()
}

type UnlockParams struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (wm *WalletManager) Unlock(name, password string) {
	wm.checkTimeout()

	for name, _ := range wm.Wallets {
		wm.log.Debug(name)
	}
	wm.log.Debug("all wallets: %v", wm.Wallets)
	if _, ok := wm.Wallets[name]; !ok {
		wm.Open(name)
	}

	wallet, ok := wm.Wallets[name]
	if !ok {
		EosThrow(&WalletException{}, "Wallet not found", name)
	}
	if !wallet.IsLocked() {
		EosThrow(&WalletUnlockedException{}, "Wallet is already unlocked:%s", name)
	}
	wallet.Unlock(password)
	wm.log.Debug("locked :%b", wallet.IsLocked())

}

type ImportKeyParams struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (wm *WalletManager) ImportKey(name, wifkey string) {
	wallet, ok := wm.Wallets[name]
	if !ok {
		EosThrow(&WalletNonexistentException{}, "Wallet nor found: %s", name)
	}

	if wallet.IsLocked() {
		EosThrow(&WalletLockedException{}, "Wallet is locked: %s\n", name)
	}

	ok = wallet.ImportKey(wifkey)
	//if err != nil {
	//	EosThrow(&KeyExistException{}, "Key already in wallet")
	//}
	if ok {
		wallet.SaveWalletFile()
	}
}

type RemoveKeyParams struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Key      string `json:"key"`
}

func (wm *WalletManager) RemoveKey(name, password, key string) {
	wm.checkTimeout()
	wallet, ok := wm.Wallets[name]
	if !ok {
		EosThrow(&WalletNonexistentException{}, "Wallet not found:%s\n", name)
	}
	if wallet.IsLocked() {
		EosThrow(&WalletLockedException{}, "Wallet is locked: %s\n", name)
	}
	wallet.CheckPassword(password) //throws if bad password
	wallet.RemoveKey(key)
}

type CreateKeyParams struct {
	Name    string `json:"name"`
	KeyType string `json:"keytype"`
}

func (wm *WalletManager) CreateKey(name, keyType string) string {
	wm.checkTimeout()
	wallet, ok := wm.Wallets[name]
	if !ok {
		EosThrow(&WalletNonexistentException{}, "Wallet not found:%s\n", name)
	}
	if wallet.IsLocked() {
		EosThrow(&WalletLockedException{}, "Wallet is locked: %s\n", name)
	}

	re := wallet.CreateKey(strings.ToUpper(keyType))
	return re
}

type SignTrxParams struct {
	Txn     *types.SignedTransaction `json:"signed_transaction"`
	Keys    []ecc.PublicKey          `json:"keys"`
	ChainID common.ChainIdType       `json:"id"`
}

func (wm *WalletManager) SignTransaction(txn *types.SignedTransaction, keys []ecc.PublicKey, chainID common.ChainIdType) *types.SignedTransaction {
	wm.checkTimeout()

	for _, key := range keys {
		found := false

		for _, wallet := range wm.Wallets {
			if !wallet.IsLocked() {
				sig := wallet.TrySignDigest(txn.SigDigest(&chainID, txn.ContextFreeData).Bytes(), key)
				if !common.Empty(sig) {
					txn.Signatures = append(txn.Signatures, *sig)
					found = true
					break
				}
			}
		}
		if !found {
			EosThrow(&WalletMissingPubKeyException{}, "public key not found in unlocked wallets %s", key)
		}
	}
	return txn
}

func (wm *WalletManager) SignDigest(digest common.DigestType, key ecc.PublicKey) (sig ecc.Signature) {
	wm.checkTimeout()
	Try(func() {
		for _, wallet := range wm.Wallets {
			if !wallet.IsLocked() {
				sig = *wallet.TrySignDigest(crypto.Sha256(digest).Bytes(), key)
				if !common.Empty(sig) {
					break
				}
			}
		}
	}).FcLogAndRethrow().End()

	if common.Empty(sig) {
		EosThrow(&WalletMissingPubKeyException{}, "public key not found in unlocked wallets %s", key)
	}
	return

}

func (wm *WalletManager) SetDir(path string) {
	wm.dir = path
}

func (wm *WalletManager) SetTimeOut(t int64) {
	wm.timeOut = time.Duration(t) * time.Second
	now := time.Now()
	wm.timeOutTime = now.Add(wm.timeOut)
	log.Debug("timeOutTime: %s", wm.timeOut)
}

//checkTimeout verify timeout has not occurred and reset timeout if not, calls lock_all() if timeout has passed
func (wm *WalletManager) checkTimeout() {
	if wm.timeOut != tstampMax {
		now := time.Now()
		if exp := now.After(wm.timeOutTime); exp {
			// lockAll()
			log.Debug("wallet has been locked,please unlock firstly") //TODO
		}
		wm.timeOutTime = now.Add(wm.timeOut)
	}
}

//func (wm *WalletManager) ownAndUseWallet(name string)
