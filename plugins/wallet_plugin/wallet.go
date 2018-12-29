package wallet_plugin

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"io/ioutil"
	"os"
)

var (
	ErrWalletLocked          = errors.New("Unable to handle a locked wallet")
	ErrWalletNoPassword      = errors.New("No password")
	ErrWallerInvalidPassword = errors.New("Invalid password for wallet")
	ErrWalletKeyExist        = errors.New("Key already in wallet")
)

const (
	walletFilenameExtension string = ".wallet"
	defaultKeyType          string = "K1"
)

type CKeys []byte
type WalletData struct {
	CipherKeys CKeys `json:"cipher_keys"` /** encrypted keys */
}

func (w CKeys) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(w))
}

func (w *CKeys) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	out, err := hex.DecodeString(s)
	*w = CKeys(out)
	return err
}

type PlainKeys struct {
	CheckSum []byte
	Keys     map[ecc.PublicKey]ecc.PrivateKey
}

type Sprivate struct {
	Curve   ecc.CurveID
	PrivKey []byte
}

type SprivateKeys struct {
	CheckSum []byte
	Keys     map[ecc.PublicKey]Sprivate
}

type SoftWallet struct {
	walletFilename string
	wallet         WalletData
	Keys           map[ecc.PublicKey]ecc.PrivateKey
	checksum       []byte
}

func (w *SoftWallet) CopyWalletFile(password string) {

}

//GetWalletFilename Returns the current wallet filename.
func (w *SoftWallet) GetWalletFilename() string {
	return w.walletFilename
}

func (w *SoftWallet) isNew() bool {
	return len(w.wallet.CipherKeys) == 0
}

func (w *SoftWallet) isLocked() bool {
	return bytes.Compare(w.checksum, nil) == 0
}

func (w *SoftWallet) Lock() (err error) {
	if w.isLocked() {
		return ErrWalletLocked
	}
	err = w.encryptKeys()
	if err != nil {
		return err
	}

	for i := range w.Keys {
		w.Keys[i] = ecc.PrivateKey{}
	}
	w.Keys = nil //TODO to clear all data
	w.checksum = nil

	return nil
}

func (w *SoftWallet) UnLock(password string) {
	Try(func() {
		FcAssert(len(password) > 0, "No password")
		pw := hash512(password)
		decrypted, err := Decrypt(string(pw[:]), w.wallet.CipherKeys)
		if err != nil {
			fmt.Println("decrypt is error:", err)
			FcAssert(false)
		}

		var pk SprivateKeys
		err = rlp.DecodeBytes(decrypted, &pk)
		if err != nil {
			fmt.Println("decodeBytes is error:", err)
			FcAssert(false)
		}

		FcAssert(bytes.Compare(pw, pk.CheckSum) == 0, "Invalid password for wallet")
		keyMap := make(map[ecc.PublicKey]ecc.PrivateKey, len(pk.Keys))
		for pub, pri := range pk.Keys {
			privateKey, err := ecc.NewDeterministicPrivateKey(bytes.NewReader(pri.PrivKey)) //TODO
			if err != nil {
				fmt.Println("NewDeterministicPrivateKey is wrong: ", err.Error())
				FcAssert(false)
			}
			keyMap[pub] = *privateKey
		}

		w.Keys = keyMap
		w.checksum = pk.CheckSum
	}).EosRethrowExceptions(&WalletInvalidPasswordException{}, "Invalid password for wallet: %s", w.GetWalletFilename()).End()

}

func (w *SoftWallet) CheckPassword(password string) {
	Try(func() {
		FcAssert(len(password) > 0)
		pw := hash512(password)
		decrypted, err := Decrypt(string(pw[:]), w.wallet.CipherKeys)
		FcAssert(err == nil)

		var pk PlainKeys
		err = rlp.DecodeBytes(decrypted, &pk.CheckSum)
		FcAssert(err == nil)
		result := bytes.Compare(pw, pk.CheckSum)
		FcAssert(result == 0)
	}).EosRethrowExceptions(&WalletInvalidPasswordException{}, "Invalid password for wallet: %s", w.walletFilename).End()
}

//SetPassword Sets a new password on the wallet
func (w *SoftWallet) SetPassword(password string) {
	if !w.isNew() {
		EosAssert(!w.isLocked(), &WalletLockedException{}, "The wallet must be unlocked before the password can be set")
	}
	w.checksum = hash512(password)
	w.Lock()
}

func (w *SoftWallet) ListPublicKeys(password string) []ecc.PublicKey {

	return nil
}
func (w *SoftWallet) LoadWalletFile() bool {
	contents, err := ioutil.ReadFile(w.walletFilename)
	if err != nil {
		fmt.Printf("read file from %s   :%s\n", w.walletFilename, err)
		return false
	}
	err = json.Unmarshal(contents, &w.wallet)
	if err != nil {
		fmt.Println("Unmarshal wallet: ", err)
		return false
	}
	return true
}

func (w *SoftWallet) SaveWalletFile() (err error) { //TODO need walletFilename ?

	w.encryptKeys()

	log.Debug("Saving wallet to file %s\n", w.walletFilename)
	data, err := json.Marshal(w.wallet)
	if err != nil {
		fmt.Println(w.wallet, err)
		return err
	}

	walletFile, err := os.OpenFile(w.walletFilename, os.O_RDWR|os.O_CREATE, 0766)
	defer walletFile.Close()
	_, err = walletFile.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (w *SoftWallet) SetWalletFilename(filename string) {
	w.walletFilename = filename
}

func (w *SoftWallet) ImportKey(wifKey string) (n bool, err error) {
	if w.isLocked() {
		return false, ErrWalletLocked
	}
	priv, err := ecc.NewPrivateKey(wifKey)
	if err != nil {
		return false, err
	}
	wifPubKey := priv.PublicKey()
	if _, find := w.Keys[wifPubKey]; !find {
		w.Keys[wifPubKey] = *priv
		return true, nil
	} else {
		return false, ErrWalletKeyExist
	}
}

func (w *SoftWallet) RemoveKey(key string) bool {
	return true
}
func (w *SoftWallet) CreateKey(keyType string) string {
	if len(keyType) == 0 {
		keyType = defaultKeyType
	}
	var privKey *ecc.PrivateKey
	switch keyType {
	case "K1":
		privKey, _ = ecc.NewRandomPrivateKey()
	case "R1":
		privKey, _ = ecc.NewRandomPrivateKey() //TODO now not suppoted r1

	default:
		EosThrow(&UnsupportedKeyTypeException{}, "Key type %s not supported by software wallet", keyType)
	}

	w.ImportKey(privKey.String())
	return privKey.PublicKey().String()
}

func (w *SoftWallet) encryptKeys() (err error) {
	if !w.isLocked() {
		keymap := make(map[ecc.PublicKey]Sprivate, 0)
		for pub, pri := range w.Keys {
			keymap[pub] = Sprivate{Curve: pri.Curve, PrivKey: pri.Serialize()}
		}
		plainkeys := SprivateKeys{Keys: keymap, CheckSum: w.checksum}
		PlainTxt, err := rlp.EncodeToBytes(plainkeys)
		if err != nil {
			fmt.Println("error while encoding wallet's key pair")
		}

		w.wallet.CipherKeys, err = Encrypt(string(plainkeys.CheckSum[:]), string(PlainTxt[:]))
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *SoftWallet) trySignDigest(digest []byte, publicKey ecc.PublicKey) *ecc.Signature {
	it, ok := w.Keys[publicKey]
	if !ok {
		return ecc.NewSigNil()
	}

	sig, err := it.Sign(digest)
	if err != nil {
		fmt.Println(err)
		return ecc.NewSigNil()
	}
	return &sig
}

func hash512(str string) []byte {
	h := sha512.New()
	_, _ = h.Write([]byte(str))
	return h.Sum(nil)
}
