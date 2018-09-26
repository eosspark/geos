package walletPlugin

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
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
	fmt.Println(data)
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	out, err := hex.DecodeString(s)
	fmt.Println("out: ", out)
	*w = CKeys(out)
	fmt.Println(w)
	return err
}

type PlainKeys struct {
	CheckSum []byte
	Keys     map[ecc.PublicKey]ecc.PrivateKey
}

type Keysmap struct {
	Pubkey ecc.PublicKey
	Prikey ecc.PrivateKey
}
type Keyspair struct {
	Prikey ecc.PrivateKey
	Pubkey ecc.PublicKey
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

// func (w *WalletData) GetPrivateKey(password string) ecc.PrivateKey {

// 	return nil
// }
// func (w *WalletData) GetPrivateKeyFromPassword(password string) Keyspair { //TODO
// 	return nil
// }

func (w *SoftWallet) isnew() bool {
	if len(w.wallet.CipherKeys) == 0 {
		return true
	}
	return false

}

func (w *SoftWallet) isLocked() bool {
	result := bytes.Compare(w.checksum, nil)
	return result == 0
	// if result := bytes.Compare(w.checksum, nil); result == 0 {
	// 	return true
	// }
	// return false //checksum 为 nil
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
func (w *SoftWallet) UnLock(password string) (err error) {
	if len([]rune(password)) == 0 {
		return ErrWalletNoPassword
	}
	pw := hash512(password)
	decrypted, err := Decrypt(string(pw[:]), w.wallet.CipherKeys)
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
	w.Keys = pk.Keys
	w.checksum = pk.CheckSum
	return nil
}

func (w *SoftWallet) CheckPassword(password string) (err error) {
	if len(password) > 0 {
		pw := hash512(password)
		decrypted, err := Decrypt(string(pw[:]), w.wallet.CipherKeys)
		if err != nil {
			return err
		}
		var pk PlainKeys
		err = rlp.DecodeBytes(decrypted, &pk)
		if err != nil {
			return err
		}
		if result := bytes.Compare(pw, pk.CheckSum); result == 0 {
			return nil
		}
	}
	return ErrWallerInvalidPassword

}

//SetPassword Sets a new password on the wallet
func (w *SoftWallet) SetPassword(password string) error {
	if !w.isnew() {
		fmt.Println("old ")
		if !w.isLocked() {
			return ErrWalletLocked
		}
	}

	w.checksum = hash512(password)
	w.Lock()
	return nil
}

func (w *SoftWallet) ListKeys(password string) []Keyspair {
	return nil
}
func (w *SoftWallet) ListPublicKeys(password string) []ecc.PublicKey {

	return nil
}
func (w *SoftWallet) LoadWalletFile() bool { //TODO need filename ?
	// TODO:  Merge imported wallet with existing wallet,
	//        instead of replacing it

	walletFile, err := os.Open(w.walletFilename)
	defer walletFile.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}
	buf := make([]byte, 1024) //TODO
	lenth, err := walletFile.Read(buf)
	if err != nil {
		fmt.Println(w.walletFilename, string(buf[:lenth]), err)
		return false
	}
	err = json.Unmarshal(buf[:lenth], &w.wallet)
	fmt.Println("wallet: ", w.wallet, err)
	return true
}

// if( wallet_filename == "" )
//    wallet_filename = _wallet_filename;

// if( ! fc::exists( wallet_filename ) )
//    return false;

// _wallet = fc::json::from_file( wallet_filename ).as< wallet_data >();

// return true;

// func (w *SoftWallet) SaveWalletFile(walletFilename string) {
// func (w *SoftWallet) SaveWalletFile() (err error) { //TODO need walletFilename ?
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

// func (w CipherKeys) MarshalJSON() ([]byte, error) {

// 	return json.Marshal(hex.EncodeToString(w))
// }

func (w *SoftWallet) SaveWalletFile() (err error) { //TODO need walletFilename ?
	w.encryptKeys()
	fmt.Printf("Saving wallet to file %s\n", w.walletFilename)
	fmt.Println(w.wallet)
	data, err := json.Marshal(w.wallet)
	if err != nil {
		fmt.Println(w.wallet, err)
		return err
	}
	fmt.Println("data: ", data, w.wallet)

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
		// fmt.Println("Keypair: ", keys[wifPubKey], wifPubKey)
		return true, nil
	} else {
		return false, ErrWalletKeyExist
	}
	return false, nil
}
func (w *SoftWallet) RemoveKey(key string) bool {
	return true
}
func (w *SoftWallet) CreateKey(keyType string) string {
	return "test"
}

// func (w *SoftWallet) TrySignDigest(digest []byte, pubkey ecc.PublicKey) ecc.Signature {
// 	return nil
// }

func (w *SoftWallet) encryptKeys() (err error) {
	if !w.isLocked() {
		data := PlainKeys{}
		data.Keys = w.Keys
		data.CheckSum = w.checksum
		PlainTxt, err := rlp.EncodeToBytes(data)

		if err != nil {
			fmt.Println("error while encoding wallet's key pair")
		}

		w.wallet.CipherKeys, err = Encrypt(string(data.CheckSum[:]), string(PlainTxt[:]))
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
	decresult, err := Decrypt(string(deckey[:]), wallet.CipherKeys)

	var data PlainKeys
	err = rlp.DecodeBytes(decresult, &data)
	for pub, priv := range data.Keys {
		fmt.Println("keypairs: ", pub, priv)
	}
	return data.Keys
}
