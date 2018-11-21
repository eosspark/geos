package wallet_plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
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
	passwordPrefix string = "pw"
)

//var wallets map[string]SoftWallet
var timeOut time.Duration //senconds max //how long to wait before calling lock_all()
var timeOuttime time.Time // when to call lock_all()
var dir string = "."

const tstampMax = 3600 * time.Second

// const timepointMax =
func init() {
	timeOut = tstampMax
}

type RespKeys map[ecc.PublicKey]ecc.PrivateKey

func (k RespKeys) MarshalJSON() ([]byte, error) {
	out := map[string]string{}
	for pub, pri := range k {
		putstr := pub.String()
		pristr := pri.String()
		out[putstr] = pristr
	}
	return json.Marshal(out)
}

type WalletPlugin struct {
	Wallets map[string]SoftWallet
}

func NewWalletPlugin() *WalletPlugin {
	return &WalletPlugin{
		Wallets: make(map[string]SoftWallet),
	}
}

func (wp *WalletPlugin) Create(name string) (password string, err error) {
	log.Debug("wallet creating")
	checkTimeout()

	password, err = genPassword()

	file, err := os.Open(dir)
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
			err = fmt.Errorf("Wallet with name: %s already exists at %s", walletName, dir)
			return
		}
	}

	var wallet SoftWallet
	err = wallet.SetPassword(password)
	if err != nil {
		return
	}
	walletFileName := fmt.Sprintf("%s/%s%s", dir, name, fileExt)
	wallet.SetWalletFilename(walletFileName)
	err = wallet.UnLock(password)
	if err != nil {
		return
	}
	wallet.Lock()
	wallet.UnLock(password)
	wallet.SaveWalletFile()

	if _, ok := wp.Wallets[name]; ok {
		delete(wp.Wallets, name)
	}
	wp.Wallets[name] = wallet
	return password, nil
}

func (wp *WalletPlugin) Open(name string) error {
	checkTimeout()
	log.Debug("Opening wallet :   wallet name: ", name)
	var wallet SoftWallet
	walletFileName := fmt.Sprintf("%s/%s%s", dir, name, fileExt)
	wallet.SetWalletFilename(walletFileName)
	if !wallet.LoadWalletFile() {
		return fmt.Errorf("Unable to open file: %s", walletFileName)
	}
	if _, ok := wp.Wallets[name]; ok {
		delete(wp.Wallets, name)
	}
	wp.Wallets[name] = wallet
	// log.Debug(walletname, wallet.wallet.CipherKeys)
	return nil
}

func (wp *WalletPlugin) ListWallets() []string {
	log.Debug("list wallets")
	var result []string
	for name, wallet := range wp.Wallets {
		if wallet.isLocked() {
			result = append(result, name)
		} else {
			result = append(result, name+"*")
		}
	}
	return result
}

func (wp *WalletPlugin) ListKeys(name, password string) (re RespKeys, err error) {
	checkTimeout()
	log.Debug("list keys")

	if _, ok := wp.Wallets[name]; !ok {
		err = fmt.Errorf("Wallet not found: %s", name)
		return nil, err
	}
	wallet := wp.Wallets[name]
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

func (wp *WalletPlugin) GetPublicKeys() (re []string, err error) {
	log.Debug("get public keys")

	if len(wp.Wallets) == 0 {
		return nil, ErrWalletNotAvaliable
	}

	isAllWalletLocked := true
	for name, wallet := range wp.Wallets {
		if !wallet.isLocked() {
			isAllWalletLocked = false
			log.Debug("wallet: %s is unlocked\n", name)
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

func (wp *WalletPlugin) LockAllwallets() {
	log.Debug("lock all wallets")
	wp.lockAll()
}

func (wp *WalletPlugin) lockAll() {
	for _, wallet := range wp.Wallets {
		if !wallet.isLocked() {
			wallet.Lock()
		}
	}
}

func (wp *WalletPlugin) Lock(name string) error {
	log.Debug("lock wallet")
	if _, ok := wp.Wallets[name]; !ok {
		return fmt.Errorf("Wallet not found: %s", name)
	}
	wallet := wp.Wallets[name]
	wallet.Lock()
	return nil
}

func (wp *WalletPlugin) Unlock(name, password string) error {
	checkTimeout()
	var wallet SoftWallet

	log.Debug("unlock wallet", name, password)

	if _, ok := wp.Wallets[name]; !ok {
		// open(){
		walletFileName := fmt.Sprintf("%s/%s%s", dir, name, fileExt)
		wallet.SetWalletFilename(walletFileName)
		if !wallet.LoadWalletFile() {
			return fmt.Errorf("Unable to open file: %s", walletFileName)
		}
		if _, ok := wp.Wallets[name]; ok {
			delete(wp.Wallets, name)
		}
		wp.Wallets[name] = wallet
		// }
	}

	wallet = wp.Wallets[name]
	if !wallet.isLocked() {
		return fmt.Errorf("Wallet is already unlocked: %s", name)
	}

	err := wallet.UnLock(password)
	if err != nil {
		return err
	}
	delete(wp.Wallets, name)
	wp.Wallets[name] = wallet

	// for pub, pri := range wallet.Keys {
	// 	log.Debug(pub, pri, wallet.Keys[pub])
	// }
	return nil
}

func (wp *WalletPlugin) ImportKey(name, wifkey string) error {
	log.Debug("wallet import keys", name, wifkey)
	wallet, ok := wp.Wallets[name]
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

func (wp *WalletPlugin) RemoveKey(name, password, key string) error {
	log.Debug("remove key")
	checkTimeout()
	wallet, ok := wp.Wallets[name]
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

func (wp *WalletPlugin) CreateKey(name, keyType string) (string, error) {
	log.Debug("create key")
	checkTimeout()
	wallet, ok := wp.Wallets[name]
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

func (wp *WalletPlugin) SignTransaction(txn types.SignedTransaction, keys []ecc.PublicKey, chainID common.ChainIdType) (types.SignedTransaction, error) {
	checkTimeout()
	log.Debug("sign transaction")
	for _, key := range keys {
		found := false

		for _, wallet := range wp.Wallets {
			if !wallet.isLocked() {
				sig := wallet.trySignDigest(txn.SigDigest(&chainID, txn.ContextFreeData), key)
				if !common.Empty(sig) {
					txn.Signatures = append(txn.Signatures, *sig)
					found = true
					break
				}
			}
		}
		if !found {
			try.EosThrow(&exception.WalletMissingPubKeyException{}, "public key not found in unlocked wallets %s", key)
		}
	}
	return txn, nil

}

func (wp *WalletPlugin) SignDigest(digest common.DigestType, key ecc.PublicKey) (sig ecc.Signature) {
	checkTimeout()
	log.Debug("sign digest")
	try.Try(func() {
		for _, wallet := range wp.Wallets {
			if !wallet.isLocked() {
				sig = *wallet.trySignDigest(crypto.Sha256(digest).Bytes(), key)
				if !common.Empty(sig) {
					break
				}
			}
		}
	}).FcLogAndRethrow().End()

	if common.Empty(sig) {
		try.EosThrow(&exception.WalletMissingPubKeyException{}, "public key not found in unlocked wallets %s", key)
	}
	return

}

func genPassword() (password string, err error) {
	prikey, err := ecc.NewRandomPrivateKey()
	if err != nil {
		return "", err
	}
	password = "PW" + prikey.String()
	return
}

func SetDir(path string) {
	dir = path
	log.Debug("dir: %s", dir)
}

func SetTimeOut(t int64) {
	timeOut = time.Duration(t) * time.Second
	now := time.Now()
	timeOuttime = now.Add(timeOut)
	log.Debug("timeOutTime: %s", timeOut)
}

//checkTimeout verify timeout has not occurred and reset timeout if not, calls lock_all() if timeout has passed
func checkTimeout() {
	if timeOut != tstampMax {
		now := time.Now()
		if exp := now.After(timeOuttime); exp {
			// lockAll()
			log.Debug("wallet has been locked,please unlock firstly") //TODO
		}
		timeOuttime = now.Add(timeOut)
	}
}
