package walletPlugin

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/ecc"
	"github.com/eosspark/eos-go/rlp"
	"os"
)

var (
	ErrWalletLocked          = errors.New("Unable to handle a locked wallet")
	ErrWalletNoPassword      = errors.New("No password")
	ErrWallerInvalidPassword = errors.New("Invalid password for wallet")
	ErrWalletKeyExist        = errors.New("Key already in wallet")
)

var (
	walletFilename string
	wallet         WalletData
	keys           map[ecc.PublicKey]ecc.PrivateKey
	checksum       []byte
)

const (
	walletFilenameExtension string = ".wallet"
	defaultKeyType          string = "K1"
)

type Keysmap struct {
	Pubkey ecc.PublicKey
	Prikey ecc.PrivateKey
}
type Keyspair struct {
	Prikey ecc.PrivateKey
	Pubkey ecc.PublicKey
}

type WalletData struct {
	CipherKeys []byte /** encrypted keys */
}

type PlainKeys struct {
	CheckSum []byte
	Keys     map[ecc.PublicKey]ecc.PrivateKey
}

func (w *WalletData) CopyWalletFile(password string) {

}

//GetWalletFilename Returns the current wallet filename.
func (w *WalletData) GetWalletFilename() string {
	return "test"
}

// func (w *WalletData) GetPrivateKey(password string) ecc.PrivateKey {

// 	return nil
// }
// func (w *WalletData) GetPrivateKeyFromPassword(password string) Keyspair { //TODO
// 	return nil
// }

func (w *WalletData) isnew() bool {
	if len(w.CipherKeys) == 0 {
		fmt.Println(true)
		return true
	}
	return false

}

func (w *WalletData) isLocked() bool {
	if result := bytes.Compare(checksum, nil); result == 0 {
		return true
	}
	return false //checksum 为 nil
}

func (w *WalletData) Lock() (err error) {
	if w.isLocked() {
		return ErrWalletLocked
	}
	err = w.encryptKeys()
	if err != nil {
		return err
	}

	for i := range keys {
		keys[i] = ecc.PrivateKey{}
	}
	keys = nil
	checksum = nil

	return nil
}
func (w *WalletData) UnLock(password string) (err error) {
	if len([]rune(password)) == 0 {
		return ErrWalletNoPassword
	}
	pw := hash512(password)
	decrypted, err := aesEnc.Decrypt(string(pw[:]), w.CipherKeys)
	if err != nil {
		return err
	}
	var pk PlainKeys
	err = rlp.DecodeBytes(decrypted, &pk)
	if err != nil {
		return err
	}

	if result := bytes.Compare(pw, pk.CheckSum); result != 0 {
		return ErrWallerInvalidPassword
	}
	keys = pk.Keys
	checksum = pk.CheckSum
	return nil
}

func (w *WalletData) CheckPassword(password string) {

}

//SetPassword Sets a new password on the wallet
func (w *WalletData) SetPassword(password string) error {
	if !w.isnew() {
		fmt.Println("old ")
		if !w.isLocked() {
			return ErrWalletLocked
		}
	}

	checksum = hash512(password)
	w.Lock()
	return nil
}

func (w *WalletData) ListKeys(password string) []Keyspair {
	return nil
}
func (w *WalletData) ListPublicKeys(password string) []ecc.PublicKey {

	return nil
}
func (w *WalletData) LoadWalletFile(walletFilename string) {

}

// func (w *WalletData) SaveWalletFile(walletFilename string) {
// func (w *WalletData) SaveWalletFile() (err error) { //TODO need walletFilename ?
// 	w.encryptKeys()
// 	fmt.Printf("Saving wallet to file %s\n", walletFilename)

// 	data, err := json.Marshal(w)
// 	if err != nil {
// 		return err
// 	}
// 	walletFile, err := os.OpenFile(walletFilename, os.O_RDWR|os.O_CREATE, 0766)
// 	defer walletFile.Close()
// 	_, err = walletFile.Write(data)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func SaveWalletFile() (err error) { //TODO need walletFilename ?
	d.encryptKeys()
	fmt.Printf("Saving wallet to file %s\n", walletFilename)

	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	walletFile, err := os.OpenFile(walletFilename, os.O_RDWR|os.O_CREATE, 0766)
	defer walletFile.Close()
	_, err = walletFile.Write(data)
	if err != nil {
		return err
	}
	return nil
}
func (w *WalletData) SetWalletFilename(filename string) {
	walletFilename = filename
}
func (w *WalletData) ImportKey(wifKey string) (n bool, err error) {
	if w.isLocked() {
		return false, ErrWalletLocked
	}

	priv, err := ecc.NewPrivateKey(wifKey)
	if err != nil {
		return false, err
	}
	wifPubKey := priv.PublicKey()

	if _, find := keys[wifPubKey]; !find {
		keys[wifPubKey] = *priv
		// fmt.Println("Keypair: ", keys[wifPubKey], wifPubKey)
		return true, nil
	} else {
		return false, ErrWalletKeyExist
	}

	return false, nil
}
func (w *WalletData) RemoveKey(key string) bool {
	return true
}
func (w *WalletData) CreateKey(keyType string) string {
	return "test"
}

// func (w *WalletData) TrySignDigest(digest []byte, pubkey ecc.PublicKey) ecc.Signature {
// 	return nil
// }

func (w *WalletData) encryptKeys() (err error) {
	if !w.isLocked() {
		data := PlainKeys{}
		data.Keys = keys
		data.CheckSum = checksum
		PlainTxt, err := rlp.EncodeToBytes(data)

		if err != nil {
			fmt.Println("error while encoding wallet's key pair")
		}

		d.CipherKeys, err = aesEnc.Encrypt(string(data.CheckSum[:]), string(PlainTxt[:]))
		if err != nil {
			return err
		}
	}
	return nil
}

func hash512(str string) (s []byte) {
	h := sha512.New()
	_, _ = h.Write([]byte(str))
	s = h.Sum(nil)
	return
}

func getdata(walletname, password string) map[ecc.PublicKey]ecc.PrivateKey {
	file, err := os.Open(walletname)
	if err != nil {
		fmt.Println(err)
	}
	buf := make([]byte, 1024)
	lenth, _ := file.Read(buf) //TODO lenth > 1024?
	file.Close()

	wallet := WalletData{}
	json.Unmarshal(buf[:lenth], &wallet)

	deckey := hash512(password)
	decresult, err := aesEnc.Decrypt(string(deckey[:]), wallet.CipherKeys)

	var data PlainKeys
	err = rlp.DecodeBytes(decresult, &data)
	for pub, priv := range data.Keys {
		fmt.Println("keypairs: ", pub, priv)
	}
	return data.Keys
}
