package wallet_plugin

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/chain/types/generated_containers"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
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

var DefaultWalletFilename = ""

type SoftWalletImpl struct {
	WalletFilename string
	Wallet         WalletData
	Keys           map[ecc.PublicKey]ecc.PrivateKey
	Checksum       []byte
}

func (w *SoftWalletImpl) EncryptKeys() {
	if !w.IsLocked() {
		keymap := make(map[ecc.PublicKey]Sprivate, 0)
		for pub, pri := range w.Keys {
			keymap[pub] = Sprivate{Curve: pri.Curve, PrivKey: pri.Serialize()}
		}
		plainKeys := SprivateKeys{Keys: keymap, CheckSum: w.Checksum}
		PlainTxt, err := rlp.EncodeToBytes(plainKeys)
		if err != nil {
			fmt.Println("error while encoding wallet's key pair")
		}

		w.Wallet.CipherKeys, err = Encrypt(string(plainKeys.CheckSum[:]), string(PlainTxt[:]))
	}
}

func (w *SoftWalletImpl) CopyWalletFile(password string) bool {
	return true
}

func (w *SoftWalletImpl) IsLocked() bool {
	return bytes.Compare(w.Checksum, nil) == 0
}

func (w *SoftWalletImpl) GetWalletFilename() string {
	return w.WalletFilename
}

func (w *SoftWalletImpl) TryGetPrivateKey(id ecc.PublicKey) *ecc.PrivateKey {
	key, ok := w.Keys[id]
	if ok {
		return &key
	} else {
		//TODO
		return nil
	}
}

func (w *SoftWalletImpl) TrySignDigest(digest []byte, publicKey ecc.PublicKey) *ecc.Signature {
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

func (w *SoftWalletImpl) GetPrivateKey(pubkey ecc.PublicKey) ecc.PrivateKey {
	hasKey := w.TryGetPrivateKey(pubkey)
	EosAssert(hasKey != nil, &KeyNonexistentException{}, "Key doesn't exist!")
	return *hasKey
}

func (w *SoftWalletImpl) ImportKey(wifKey string) bool {
	priv, err := ecc.NewPrivateKey(wifKey)
	if err != nil {
		log.Error("Wrong to NewPrivateKey", err)
		return false
	}
	wifPubKey := priv.PublicKey()
	if _, find := w.Keys[wifPubKey]; !find {
		w.Keys[wifPubKey] = *priv
		return true
	}
	EosThrow(&KeyExistException{}, "Key already in wallet")
	return false
}

func (w *SoftWalletImpl) RemoveKey(key string) bool {
	pub, err := ecc.NewPublicKey(key)
	if err != nil {
		log.Error("Wrong to NewPublicKey", err)
		return false
	}
	if _, find := w.Keys[pub]; find {
		delete(w.Keys, pub)
		return true
	}
	EosThrow(&KeyNonexistentException{}, "Key not in wallet")
	return false
}

func (w *SoftWalletImpl) CreateKey(keyType string) string {
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

func (w *SoftWalletImpl) LoadWalletFile(walletFilename string) bool {
	if walletFilename == "" {
		walletFilename = w.WalletFilename
	}
	contents, err := ioutil.ReadFile(w.WalletFilename)
	if err != nil {
		fmt.Printf("read file from %s   :%s\n", w.WalletFilename, err)
		return false
	}
	err = json.Unmarshal(contents, &w.Wallet)
	if err != nil {
		fmt.Println("Unmarshal wallet: ", err)
		return false
	}
	return true
}

func (w *SoftWalletImpl) SaveWalletFile(walletFilename string) {
	w.EncryptKeys()

	data, err := json.Marshal(w.Wallet)
	if err != nil {
		fmt.Println(w.Wallet, err)
	}
	walletFile, err := os.OpenFile(w.WalletFilename, os.O_RDWR|os.O_CREATE, 0766)
	defer walletFile.Close()
	_, err = walletFile.Write(data)
}

type SoftWallet struct {
	my *SoftWalletImpl
}

func (w *SoftWallet) CopyWalletFile(destinationFilename string) bool {
	return w.my.CopyWalletFile(destinationFilename)
}

func (w *SoftWallet) GetWalletFilename() string {
	return w.my.GetWalletFilename()
}

func (w *SoftWallet) ImportKey(wifKey string) bool {
	EosAssert(!w.IsLocked(), &WalletLockedException{}, "Unable to import key on a locked wallet")
	if w.my.ImportKey(wifKey) {
		w.SaveWalletFile(DefaultWalletFilename)
		return true
	}
	return false
}

func (w *SoftWallet) RemoveKey(key string) bool {
	EosAssert(!w.IsLocked(), &WalletLockedException{}, "Unable to remove key on a locked wallet")
	if w.my.RemoveKey(key) {
		w.SaveWalletFile(DefaultWalletFilename)
		return true
	}
	return false
}

func (w *SoftWallet) CreateKey(keyType string) string {
	EosAssert(!w.IsLocked(), &WalletLockedException{}, "Unable to create key on a locked wallet")
	ret := w.my.CreateKey(keyType)
	w.SaveWalletFile(DefaultWalletFilename)
	return ret
}

func (w *SoftWallet) LoadWalletFile(walletFilename string) bool {
	return w.my.LoadWalletFile(walletFilename)
}

func (w *SoftWallet) SaveWalletFile(walletFilename string) {
	w.my.SaveWalletFile(walletFilename)
}

func (w *SoftWallet) IsLocked() bool {
	return w.my.IsLocked()
}

func (w *SoftWallet) IsNew() bool {
	return len(w.my.Wallet.CipherKeys) == 0
}

func (w *SoftWallet) EncryptKeys() {
	w.my.EncryptKeys()
}

func (w *SoftWallet) Lock() {
	Try(func() {
		EosAssert(!w.IsLocked(), &WalletLockedException{}, "Unable to lock a locked wallet")
		w.EncryptKeys()
		for i := range w.my.Keys {
			w.my.Keys[i] = ecc.PrivateKey{}
		}
		w.my.Keys = map[ecc.PublicKey]ecc.PrivateKey{}
		w.my.Checksum = []byte{}
	}).EosRethrowExceptions(&WalletInvalidPasswordException{}, "Invalid password for wallet: \"%v\"", w.GetWalletFilename()).End()
}

func (w *SoftWallet) Unlock(password string) {
	Try(func() {
		FcAssert(len(password) > 0)
		pw := hash512(password)
		decrypted, err := Decrypt(string(pw[:]), w.my.Wallet.CipherKeys)
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

		w.my.Keys = keyMap
		w.my.Checksum = pk.CheckSum
	}).EosRethrowExceptions(&WalletInvalidPasswordException{}, "Invalid password for wallet: %s", w.GetWalletFilename()).End()

}

func (w *SoftWallet) CheckPassword(password string) {
	Try(func() {
		FcAssert(len(password) > 0)
		pw := hash512(password)
		decrypted, err := Decrypt(string(pw[:]), w.my.Wallet.CipherKeys)
		FcAssert(err == nil)

		var pk PlainKeys
		err = rlp.DecodeBytes(decrypted, &pk.CheckSum)
		FcAssert(err == nil)
		result := bytes.Compare(pw, pk.CheckSum)
		FcAssert(result == 0)
	}).EosRethrowExceptions(&WalletInvalidPasswordException{}, "Invalid password for wallet: %s", w.my.WalletFilename).End()
}

//SetPassword Sets a new password on the wallet
func (w *SoftWallet) SetPassword(password string) {
	if !w.IsNew() {
		EosAssert(!w.IsLocked(), &WalletLockedException{}, "The wallet must be unlocked before the password can be set")
	}
	w.my.Checksum = hash512(password)
	w.Lock()
}

func (w *SoftWallet) ListKeys() map[ecc.PublicKey]ecc.PrivateKey {
	EosAssert(!w.IsLocked(), &WalletLockedException{}, "Unable to list public keys of a locked wallet")
	return w.my.Keys
}

func (w *SoftWallet) ListPublicKeys() generated.PublicKeySet {
	EosAssert(!w.IsLocked(), &WalletLockedException{}, "Unable to list private keys of a locked wallet")
	keys := generated.PublicKeySet{}
	for pk := range w.my.Keys {
		keys.Add(pk)
	}
	return keys
}

func (w *SoftWallet) GetPrivateKey(pubkey ecc.PublicKey) ecc.PrivateKey {
	return w.my.GetPrivateKey(pubkey)
}

func (w *SoftWallet) TrySignDigest(digest []byte, publicKey ecc.PublicKey) *ecc.Signature {
	it, ok := w.my.Keys[publicKey]
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

func (w *SoftWallet) GetPrivateKeyFromPassword(account string, role string, password string) common.Pair {
	seed := account + role + password
	EosAssert(len(seed) != 0, &WalletException{}, "seed should not be empty")
	secret := crypto.Hash256(seed).Bytes()
	g := bytes.NewReader(secret)
	pk, _ := ecc.NewDeterministicPrivateKey(g)
	return common.MakePair(pk.PublicKey(), pk)
}

func (w *SoftWallet) SetWalletFilename(filename string) {
	w.my.WalletFilename = filename
}

func hash512(str string) []byte {
	h := sha512.New()
	_, _ = h.Write([]byte(str))
	return h.Sum(nil)
}
