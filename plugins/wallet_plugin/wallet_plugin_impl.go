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
	"github.com/eosspark/eos-go/plugins/appbase/asio"
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

type WalletPluginImpl struct {
	timeOut     time.Duration //senconds max //how long to wait before calling lock_all()
	timeOutTime time.Time     // when to call lock_all()
	dir         string
	lockPath    string

	Wallets map[string]SoftWallet

	Self *WalletPlugin
	log  log.Logger
}

const tstampMax = 3600 * time.Second

func NewWalletPluginImpl(io *asio.IoContext) *WalletPluginImpl {

	impl := &WalletPluginImpl{
		timeOut:  tstampMax,
		dir:      ".",
		lockPath: "./wallet.lock",
		Wallets:  make(map[string]SoftWallet),
	}

	impl.log = log.New("wallet_plugin")
	impl.log.SetHandler(log.TerminalHandler)
	//impl.log.SetHandler(log.DiscardHandler())
	return impl
}

func genPassword() (password string, err error) {
	prikey, err := ecc.NewRandomPrivateKey()
	if err != nil {
		return "", err
	}
	password = "PW" + prikey.String()
	return
}

func (impl *WalletPluginImpl) Create(name string) (password string, err error) {
	impl.log.Debug("wallet creating")
	impl.checkTimeout()

	if name == "" {
		name = "default"
	}

	password, err = genPassword()

	file, err := os.Open(impl.dir)
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
			err = fmt.Errorf("Wallet with name: %s already exists at %s", walletName, impl.dir)
			return
		}
	}

	var wallet SoftWallet
	err = wallet.SetPassword(password)
	if err != nil {
		return
	}
	walletFileName := fmt.Sprintf("%s/%s%s", impl.dir, name, fileExt)
	wallet.SetWalletFilename(walletFileName)
	err = wallet.UnLock(password)
	if err != nil {
		return
	}
	wallet.Lock()
	wallet.UnLock(password)
	wallet.SaveWalletFile()

	if _, ok := impl.Wallets[name]; ok {
		delete(impl.Wallets, name)
	}
	impl.Wallets[name] = wallet
	return password, nil
}

func (impl *WalletPluginImpl) Open(name string) error {
	impl.checkTimeout()
	impl.log.Debug("Opening wallet :   wallet name: ", name)
	var wallet SoftWallet
	walletFileName := fmt.Sprintf("%s/%s%s", impl.dir, name, fileExt)
	wallet.SetWalletFilename(walletFileName)
	if !wallet.LoadWalletFile() {
		return fmt.Errorf("Unable to open file: %s", walletFileName)
	}
	if _, ok := impl.Wallets[name]; ok {
		delete(impl.Wallets, name)
	}
	impl.Wallets[name] = wallet
	// log.Debug(walletname, wallet.wallet.CipherKeys)
	return nil
}

func (impl *WalletPluginImpl) ListWallets() []string {
	impl.log.Debug("list wallets")
	var result []string
	for name, wallet := range impl.Wallets {
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

func (impl *WalletPluginImpl) ListKeys(name, password string) (re RespKeys, err error) {
	impl.checkTimeout()
	impl.log.Debug("list keys")

	if _, ok := impl.Wallets[name]; !ok {
		err = fmt.Errorf("Wallet not found: %s", name)
		return nil, err
	}
	wallet := impl.Wallets[name]
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

func (impl *WalletPluginImpl) GetPublicKeys() (re []string, err error) {
	impl.log.Debug("get public keys")

	if len(impl.Wallets) == 0 {
		return nil, ErrWalletNotAvaliable
	}

	isAllWalletLocked := true
	for name, wallet := range impl.Wallets {
		if !wallet.isLocked() {
			isAllWalletLocked = false
			impl.log.Debug("wallet: %s is unlocked\n", name)
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

func (impl *WalletPluginImpl) LockAllwallets() {
	impl.log.Debug("lock all wallets")
	impl.lockAll()
}

func (impl *WalletPluginImpl) lockAll() {
	for _, wallet := range impl.Wallets {
		if !wallet.isLocked() {
			wallet.Lock()
		}
	}
}

func (impl *WalletPluginImpl) Lock(name string) error {
	impl.log.Debug("lock wallet")
	if _, ok := impl.Wallets[name]; !ok {
		return fmt.Errorf("Wallet not found: %s", name)
	}
	wallet := impl.Wallets[name]
	wallet.Lock()
	return nil
}

func (impl *WalletPluginImpl) Unlock(name, password string) error {
	impl.checkTimeout()
	var wallet SoftWallet

	impl.log.Debug("unlock wallet", name, password)

	if _, ok := impl.Wallets[name]; !ok {
		// open(){
		walletFileName := fmt.Sprintf("%s/%s%s", impl.dir, name, fileExt)
		wallet.SetWalletFilename(walletFileName)
		if !wallet.LoadWalletFile() {
			return fmt.Errorf("Unable to open file: %s", walletFileName)
		}
		if _, ok := impl.Wallets[name]; ok {
			delete(impl.Wallets, name)
		}
		impl.Wallets[name] = wallet
		// }
	}

	wallet = impl.Wallets[name]
	if !wallet.isLocked() {
		return fmt.Errorf("Wallet is already unlocked: %s", name)
	}

	err := wallet.UnLock(password)
	if err != nil {
		return err
	}
	delete(impl.Wallets, name)
	impl.Wallets[name] = wallet

	// for pub, pri := range wallet.Keys {
	// 	log.Debug(pub, pri, wallet.Keys[pub])
	// }
	return nil
}

func (impl *WalletPluginImpl) ImportKey(name, wifkey string) error {
	impl.log.Debug("wallet import keys", name, wifkey)
	wallet, ok := impl.Wallets[name]
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

func (impl *WalletPluginImpl) RemoveKey(name, password, key string) error {
	impl.log.Debug("remove key")
	impl.checkTimeout()
	wallet, ok := impl.Wallets[name]
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

func (impl *WalletPluginImpl) CreateKey(name, keyType string) (string, error) {
	impl.log.Debug("create key")
	impl.checkTimeout()
	wallet, ok := impl.Wallets[name]
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

func (impl *WalletPluginImpl) SignTransaction(txn types.SignedTransaction, keys []ecc.PublicKey, chainID common.ChainIdType) (types.SignedTransaction, error) {
	impl.checkTimeout()
	impl.log.Debug("sign transaction")
	for _, key := range keys {
		found := false

		for _, wallet := range impl.Wallets {
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
			EosThrow(&WalletMissingPubKeyException{}, "public key not found in unlocked wallets %s", key)
		}
	}
	return txn, nil

}

func (impl *WalletPluginImpl) SignDigest(digest common.DigestType, key ecc.PublicKey) (sig ecc.Signature) {
	impl.checkTimeout()
	impl.log.Debug("sign digest")
	Try(func() {
		for _, wallet := range impl.Wallets {
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

func (impl *WalletPluginImpl) SetDir(path string) {
	impl.dir = path
	log.Debug("dir: %s", impl.dir)
}

func (impl *WalletPluginImpl) SetTimeOut(t int64) {
	impl.timeOut = time.Duration(t) * time.Second
	now := time.Now()
	impl.timeOutTime = now.Add(impl.timeOut)
	log.Debug("timeOutTime: %s", impl.timeOut)
}

//checkTimeout verify timeout has not occurred and reset timeout if not, calls lock_all() if timeout has passed
func (impl *WalletPluginImpl) checkTimeout() {
	if impl.timeOut != tstampMax {
		now := time.Now()
		if exp := now.After(impl.timeOutTime); exp {
			// lockAll()
			log.Debug("wallet has been locked,please unlock firstly") //TODO
		}
		impl.timeOutTime = now.Add(impl.timeOut)
	}
}

//func (impl *WalletPluginImpl) ownAndUseWallet(name string)
