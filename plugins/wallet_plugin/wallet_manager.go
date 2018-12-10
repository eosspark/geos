package wallet_plugin

import (
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

	"errors"
)

const (
	fileExt        string = ".wallet"
	passwordPrefix string = "pw"
)

var (
	ErrWalletNotAvaliable = errors.New("You don't have any wallet")
	ErrWalletNotUnlocked  = errors.New("You don't have any unlocked wallet!")
)

type WalletManager struct {
	timeOut     time.Duration //senconds max //how long to wait before calling lock_all()
	timeOutTime time.Time     // when to call lock_all()
	dir         string
	lockPath    string

	Wallets map[string]SoftWallet

	Self *WalletPlugin
	log  log.Logger
}

const tstampMax = 3600 * time.Second

func walletManager() *WalletManager {

	manager := &WalletManager{
		timeOut:  tstampMax,
		dir:      ".",
		lockPath: "./wallet.lock",
		Wallets:  make(map[string]SoftWallet),
	}

	manager.log = log.New("wallet_plugin")
	manager.log.SetHandler(log.TerminalHandler)
	//manager.log.SetHandler(log.DiscardHandler())
	return manager
}

func genPassword() (password string, err error) {
	prikey, err := ecc.NewRandomPrivateKey()
	if err != nil {
		return "", err
	}
	password = "PW" + prikey.String()
	return
}

func (wm *WalletManager) Create(name string) (password string, err error) {
	wm.log.Debug("wallet creating")
	wm.checkTimeout()

	if name == "" {
		name = "default"
	}

	password, err = genPassword()

	file, err := os.Open(wm.dir)
	defer file.Close()
	if err != nil {
		return "", err
	}
	allwallets, err := file.Readdirnames(-1)
	if err != nil {
		return "", err
	}

	walletName := name + fileExt
	for _, f := range allwallets {
		if f == walletName {
			err = fmt.Errorf("Wallet with name: %s already exists at %s", walletName, wm.dir)
			return
		}
	}

	var wallet SoftWallet
	err = wallet.SetPassword(password)
	if err != nil {
		return
	}
	walletFileName := fmt.Sprintf("%s/%s%s", wm.dir, name, fileExt)
	wallet.SetWalletFilename(walletFileName)
	err = wallet.UnLock(password)
	if err != nil {
		return
	}
	wallet.Lock()
	wallet.UnLock(password)
	wallet.SaveWalletFile()

	if _, ok := wm.Wallets[name]; ok {
		delete(wm.Wallets, name)
	}
	wm.Wallets[name] = wallet
	wm.log.Info("wallets: ")
	for name := range wm.Wallets {
		wm.log.Info(name)
	}

	return password, nil
}

func (wm *WalletManager) Open(name string) error {
	wm.checkTimeout()
	wm.log.Debug("Opening wallet :   wallet name: ", name)
	var wallet SoftWallet
	walletFileName := fmt.Sprintf("%s/%s%s", wm.dir, name, fileExt)
	wallet.SetWalletFilename(walletFileName)
	if !wallet.LoadWalletFile() {
		return fmt.Errorf("Unable to open file: %s", walletFileName)
	}
	if _, ok := wm.Wallets[name]; ok {
		delete(wm.Wallets, name)
	}
	wm.Wallets[name] = wallet
	// log.Debug(walletname, wallet.wallet.CipherKeys)
	return nil
}

func (wm *WalletManager) ListWallets() []string {
	wm.log.Debug("list wallets")
	var result []string
	for name, wallet := range wm.Wallets {
		if wallet.isLocked() {
			result = append(result, name)
		} else {
			result = append(result, name+"*")
		}
	}
	return result
}

type RespKeys map[ecc.PublicKey]ecc.PrivateKey

//func (k RespKeys) MarshalJSON() ([]byte, error) {
//	out := map[string]string{}
//	for pub, pri := range k {
//		putstr := pub.String()
//		pristr := pri.String()
//		out[putstr] = pristr
//	}
//	return json.Marshal(out)
//}

func (wm *WalletManager) ListKeys(name, password string) (re RespKeys, err error) {
	wm.checkTimeout()
	wm.log.Debug("list keys")

	if _, ok := wm.Wallets[name]; !ok {
		err = fmt.Errorf("Wallet not found: %s", name)
		return nil, err
	}
	wallet := wm.Wallets[name]
	if wallet.isLocked() {
		err = fmt.Errorf("Wallet is locked: %s", name)
		return nil, err
	}
	err = wallet.CheckPassword(password)
	if err != nil {
		return nil, err
	}

	for pub, pri := range wallet.Keys {
		re[pub] = pri
		// log.Debug(pub, wallet.Keys[pub])
	}
	return
}

func (wm *WalletManager) GetPublicKeys() (re []string, err error) {
	wm.log.Debug("get public keys")

	if len(wm.Wallets) == 0 {
		return nil, ErrWalletNotAvaliable
	}

	isAllWalletLocked := true
	for name, wallet := range wm.Wallets {
		if !wallet.isLocked() {
			isAllWalletLocked = false
			wm.log.Debug("wallet: %s is unlocked\n", name)
			for pubkey, _ := range wallet.Keys {
				re = append(re, pubkey.String())
			}
		}
	}
	if isAllWalletLocked {
		return nil, ErrWalletNotUnlocked
	}
	return re, nil
}

func (wm *WalletManager) LockAllwallets() {
	wm.log.Debug("lock all wallets")
	wm.lockAll()
}

func (wm *WalletManager) lockAll() {
	for _, wallet := range wm.Wallets {
		if !wallet.isLocked() {
			wallet.Lock()
		}
	}
}

func (wm *WalletManager) Lock(name string) error {
	wm.log.Debug("lock wallet")
	if _, ok := wm.Wallets[name]; !ok {
		return fmt.Errorf("Wallet not found: %s", name)
	}
	wallet := wm.Wallets[name]
	wallet.Lock()
	return nil
}

func (wm *WalletManager) Unlock(name, password string) error {
	wm.checkTimeout()
	var wallet SoftWallet

	wm.log.Debug("unlock wallet", name, password)

	if _, ok := wm.Wallets[name]; !ok {
		// open(){
		walletFileName := fmt.Sprintf("%s/%s%s", wm.dir, name, fileExt)
		wallet.SetWalletFilename(walletFileName)
		if !wallet.LoadWalletFile() {
			return fmt.Errorf("Unable to open file: %s", walletFileName)
		}
		if _, ok := wm.Wallets[name]; ok {
			delete(wm.Wallets, name)
		}
		wm.Wallets[name] = wallet
		// }
	}

	wallet = wm.Wallets[name]
	if !wallet.isLocked() {
		return fmt.Errorf("Wallet is already unlocked: %s", name)
	}

	err := wallet.UnLock(password)
	if err != nil {
		return err
	}
	delete(wm.Wallets, name)
	wm.Wallets[name] = wallet

	// for pub, pri := range wallet.Keys {
	// 	log.Debug(pub, pri, wallet.Keys[pub])
	// }
	return nil
}

func (wm *WalletManager) ImportKey(name, wifkey string) error {
	wm.log.Debug("wallet import keys %s,%s", name, wifkey)
	wallet, ok := wm.Wallets[name]
	if !ok {
		return fmt.Errorf("Wallet not found: %s\n", name)
	}

	if wallet.isLocked() {
		return fmt.Errorf("Wallet is locked: %s\n", name)
	}

	ok, err := wallet.ImportKey(wifkey)
	if err != nil {
		return fmt.Errorf("Unable to import key")
	}
	if ok {
		wallet.SaveWalletFile()
	}
	return nil
}

func (wm *WalletManager) RemoveKey(name, password, key string) error {
	wm.log.Debug("remove key")
	wm.checkTimeout()
	wallet, ok := wm.Wallets[name]
	if !ok {
		//EOS_THROW(chain::wallet_nonexistent_exception, "Wallet not found: ${w}", ("w", name));
		return fmt.Errorf("Wallet not found: %s", name)
	}
	if wallet.isLocked() {
		//EOS_THROW(chain::wallet_locked_exception, "Wallet is locked: ${w}", ("w", name));
		return fmt.Errorf("Wallet is locked:%s", name)
	}
	wallet.CheckPassword(password) //throws if bad password
	wallet.RemoveKey(key)
	return nil
}

func (wm *WalletManager) CreateKey(name, keyType string) (string, error) {
	wm.log.Debug("create key")
	wm.checkTimeout()
	wallet, ok := wm.Wallets[name]
	if !ok {
		//EOS_THROW(chain::wallet_nonexistent_exception, "Wallet not found: ${w}", ("w", name));
		return "", fmt.Errorf("Wallet not found: %s", name)
	}
	if wallet.isLocked() {
		//EOS_THROW(chain::wallet_locked_exception, "Wallet is locked: ${w}", ("w", name));
		return "", fmt.Errorf("Wallet is locked:%s", name)
	}

	re := wallet.CreateKey(strings.ToUpper(keyType))
	return re, nil
}

func (wm *WalletManager) SignTransaction(txn *types.SignedTransaction, keys []ecc.PublicKey, chainID common.ChainIdType) (*types.SignedTransaction, error) {
	wm.checkTimeout()
	wm.log.Debug("sign transaction")
	wm.log.Debug("%#v", txn)
	wm.log.Debug("%s,%s", keys, chainID)

	for _, key := range keys {
		found := false

		for _, wallet := range wm.Wallets {
			if !wallet.isLocked() {
				sig := wallet.trySignDigest(txn.SigDigest(&chainID, txn.ContextFreeData), key)
				wm.log.Error("sig :   %#v", sig)
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

	wm.log.Debug("%#v", txn)
	return txn, nil

}

func (wm *WalletManager) SignDigest(digest common.DigestType, key ecc.PublicKey) (sig ecc.Signature) {
	wm.checkTimeout()
	wm.log.Debug("sign digest")
	Try(func() {
		for _, wallet := range wm.Wallets {
			if !wallet.isLocked() {
				sig = *wallet.trySignDigest(crypto.Sha256(digest).Bytes(), key)
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
	log.Debug("dir: %s", wm.dir)
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
