package walletPlugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"net/http"
	"os"
	"time"
)

// var (
//  timeOut   time.Time
//  tstampMax int64 = 10000000 //TODO
// )
var (
	ErrWalletNotAvaliable = errors.New("You don't have any wallet")
	ErrWalletNotUnlocked  = errors.New("You don't have any unlocked wallet!")
)

const (
	file_ext        string = ".wallet"
	password_prefix string = "pw"
)

func GenPassword() (password string, err error) {
	prikey, err := ecc.NewRandomPrivateKey()
	if err != nil {
		return "", err
	}
	password = "PW" + prikey.String()
	return
}

var timeOutTime time.Time
var aesEnc AesEncrypt
var wallets map[string]WalletData
var d WalletData

func init() {
	wallets = make(map[string]WalletData)

}

func SetTimeOut() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("set timeout")
		var inputs []int64
		_ = json.NewDecoder(r.Body).Decode(&inputs)

		// timeOut = 1 * time.Second
		// now := time.Now().Second()
		// timeOutTime = now + timeOut
		// fmt.Println("timeOutTime: ", timeOutTime)
		// if timeOutTime < now {
		//  resp := fmt.Sprintf("Overflow on timeout_time, specified %t, now %t, timeout_time %t", timeOut, now, timeOutTime)
		//  // w.WriteHeader(201)
		//  w.Write([]byte(resp))
		//  return
		// }
		for i := 0; i < len(inputs); i++ {
			fmt.Println(inputs[i])
		}
		w.WriteHeader(201)
		w.Write([]byte("{}"))

	}
	return http.HandlerFunc(fn)
}

func check_timeout() {
	// if timeOut != tstampMax{
	// now := time.Now().Second()
	// if now > timeOutTime {
	//  // lockAll()
	// }
	// timeOutTime = now + timeOut.Second()
	// }
}

const dir string = "."

func WalletCreate() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("wallet creating")

		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)
		name := inputs[0]

		password, err := GenPassword()
		if err != nil {
			fmt.Println(err)
		}

		file, err := os.Open(dir)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		}
		allwallets, err := file.Readdirnames(-1)
		if err != nil {
			fmt.Println(err)
		}

		walletName := name + file_ext
		for _, f := range allwallets {
			if f == walletName {
				errResp := fmt.Sprintf("Wallet with name: %s already exists at %s", walletName, dir)
				// w.WriteHeader(201)
				json.NewEncoder(w).Encode(errResp)
				return
			}
		}

		walletFileName := fmt.Sprintf("%s/%s%s", dir, name, file_ext)
		//walletFile, err := os.OpenFile(walletFileName, os.O_RDWR|os.O_CREATE, 0766)
		//defer walletFile.Close()
		if err != nil {
			fmt.Println(err)
		}

		// var d WalletData
		err = d.SetPassword(password)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("walletFileName: ", walletFileName)
		d.SetWalletFilename(walletFileName)
		err = d.UnLock(password)
		if err != nil {
			fmt.Println(err)
		}
		d.Lock()
		d.UnLock(password)
		SaveWalletFile()

		// wallets = make(map[string]WalletData)
		if _, ok := wallets[name]; ok {
			delete(wallets, "name")
		}
		wallets[name] = d
		fmt.Println("all wallets")
		for wallet := range wallets {
			fmt.Println(wallet, wallets[wallet])
		}
		// encdata, err := aesEnc.Encrypt(password, "strMsg")
		// if err != nil {
		// 	fmt.Println("data: ", encdata, "err: ", err)
		// }

		// walletFile.Write(encdata)

		// fmt.Println("encdata: ", encdata)
		// walletFile.Close()
		// time.Sleep(1 * time.Second)
		// file4, err := os.Open(walletFileName)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// var datalen int
		// buf := make([]byte, 1024)
		// datalen, err = file4.Read(buf)
		// fmt.Println(datalen, err)

		// fmt.Println("walletFile:", buf[:datalen])
		// strMsg, err := aesEnc.Decrypt(password, buf[:datalen])
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println("strMsg: ", strMsg)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(password)

	}
	return http.HandlerFunc(fn)
}

func WalletOpen() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// localAll()
		fmt.Println("Open wallet")
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)

		for i := 0; i < len(inputs); i++ {
			fmt.Println(inputs[i])
		}

		w.WriteHeader(201)
		w.Write([]byte("{PW5J6Y2prcCHz7xzDJ3asTJg5dCtDsYpt6xxtHhT2Fy4TAqyruZcz}"))

	}
	return http.HandlerFunc(fn)
}

func ListWallets() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("list wallets")
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)

		for i := 0; i < len(inputs); i++ {
			fmt.Println(inputs[i])
		}
		w.WriteHeader(201)
		w.Write([]byte("PW5J6Y2prcCHz7xzDJ3asTJg5dCtDsYpt6xxtHhT2Fy4TAqyruZcz"))

	}
	return http.HandlerFunc(fn)
}

func WalletImportKey() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("wallet import keys 218")
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)

		name := inputs[0]
		wifkey := inputs[1]
		fmt.Println(name, wifkey)

		wallet, ok := wallets[name]
		if !ok {
			errResp := fmt.Sprintf("Wallet not found: %s\n", name)
			// w.WriteHeader(201)
			json.NewEncoder(w).Encode(errResp)
			return
		}

		if wallet.isLocked() {
			errResp := fmt.Sprintf("Wallet is locked: %s\n", name)
			// w.WriteHeader(201)
			json.NewEncoder(w).Encode(errResp)
			return
		}
		ok, err := wallet.ImportKey(wifkey)
		if err != nil {
			// w.WriteHeader(201)
			json.NewEncoder(w).Encode(err)
			return
		}
		if ok {
			SaveWalletFile()
		}
	}
	return http.HandlerFunc(fn)
}

func ListKey() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("list keys")

		w.WriteHeader(201)
		w.Write([]byte("{}"))

	}
	return http.HandlerFunc(fn)
}

func GetPublicKeys() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get public keys")
		var out []string
		if len(wallets) == 0 {
			fmt.Println("You don't have any wallet")
			// w.WriteHeader(201)
			json.NewEncoder(w).Encode(ErrWalletNotAvaliable)
			return
		}

		isAllWalletLocked := true
		for name, wallet := range wallets {
			if !wallet.isLocked() {
				isAllWalletLocked = false
				fmt.Printf("wallet: %s is unlocked\n", name)
				for pubkey, _ := range keys {
					out = append(out, pubkey.String())
				}
			}
		}
		if isAllWalletLocked {
			fmt.Println("You don't have any unlocked wallet!")
			json.NewEncoder(w).Encode(ErrWalletNotUnlocked)
			return
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(out)

	}
	return http.HandlerFunc(fn)
}

func LockAllwallets() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("lock all wallets")
		LockAll()
	}
	return http.HandlerFunc(fn)
}
func LockAll() {
	fmt.Println("locak all wallets")
	// for i := range wallets {

	// }
}

func LockWallet() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("lock wallet")
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)

		for i := 0; i < len(inputs); i++ {
			fmt.Println(inputs[i])
		}
		w.WriteHeader(201)
		w.Write([]byte("{}"))

	}
	return http.HandlerFunc(fn)
}

func UnLockWallet() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("unlock wallet")
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)
		name := inputs[0]
		password := inputs[1]
		fmt.Println(name, password)
		// walletFileName := fmt.Sprintf("%s/%s%s", dir, name, file_ext)
		// result := getdata(walletFileName, password)
		// for pub, pri := range result {
		// 	fmt.Println(pub, pri)
		// }

		w.WriteHeader(201)
		w.Write([]byte("{}"))

	}
	return http.HandlerFunc(fn)
}

func ImportKey() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("wallet import keys")
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)

		for i := 0; i < len(inputs); i++ {
			fmt.Println(inputs[i])
		}
		w.WriteHeader(201)
		w.Write([]byte("{}"))

	}
	return http.HandlerFunc(fn)
}

func RemoveKey() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("remove key")
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)

		for i := 0; i < len(inputs); i++ {
			fmt.Println(inputs[i])
		}
		w.WriteHeader(201)
		w.Write([]byte("{}"))

	}
	return http.HandlerFunc(fn)
}

func CreateKey() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("create key")
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)

		for i := 0; i < len(inputs); i++ {
			fmt.Println(inputs[i])
		}
		w.WriteHeader(201)
		w.Write([]byte("{}"))

	}
	return http.HandlerFunc(fn)
}

func SignTransaction() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("sign transaction")
		var inputs []json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil {
			fmt.Println("sign_transaction: error:", err)
			http.Error(w, "couldn't decode input", 500)
			return
		}

		var tx *types.SignedTransaction
		var requiredKeys []ecc.PublicKey
		var chainID common.ChainIDType

		if len(inputs) != 3 {
			http.Error(w, "invalid length of message, should be 3 parameters", 500)
			return
		}

		err := json.Unmarshal(inputs[0], &tx)
		if err != nil {
			http.Error(w, "decoding transaction", 500)
			return
		}

		err = json.Unmarshal(inputs[1], &requiredKeys)
		if err != nil {
			http.Error(w, "decoding required keys", 500)
			return
		}

		err = json.Unmarshal(inputs[2], &chainID)
		if err != nil {
			http.Error(w, "decoding chain id", 500)
			return
		}
		//TODO
		// signed, err := keyBag.Sign(tx, chainID, requiredKeys...)
		// if err != nil {
		//  http.Error(w, fmt.Sprintf("error signing: %s", err), 500)
		//  return
		// }

		// w.WriteHeader(201)
		// _ = json.NewEncoder(w).Encode(signed)

	}
	return http.HandlerFunc(fn)
}
func SignDigest() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("sign digest")
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)

		for i := 0; i < len(inputs); i++ {
			fmt.Println(inputs[i])
		}
		w.WriteHeader(201)
		w.Write([]byte("{}"))

	}
	return http.HandlerFunc(fn)
}

func OwnAndUseWallet() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("own and use wallet")
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)

		for i := 0; i < len(inputs); i++ {
			fmt.Println(inputs[i])
		}
		w.WriteHeader(201)
		w.Write([]byte("{}"))

	}
	return http.HandlerFunc(fn)
}
