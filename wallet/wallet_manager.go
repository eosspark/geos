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

// For reference:
// https://github.com/EOSIO/eos/tree/master/plugins/wallet_plugin

var (
	ErrWalletNotAvaliable = errors.New("You don't have any wallet")
	ErrWalletNotUnlocked  = errors.New("You don't have any unlocked wallet!")
)

const (
	fileExt        string = ".wallet"
	passwordPrefix string = "pw"
)

var wallets map[string]SoftWallet
var timeOut time.Duration //senconds max //how long to wait before calling lock_all()
var timeOuttime time.Time // when to call lock_all()
var dir string = "."

const tstampMax = 3600 * time.Second

// const timepointMax =
func init() {
	wallets = make(map[string]SoftWallet)
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
	fmt.Println("dir: ", dir)
}

func SetTimeOut(t int64) {
	timeOut = time.Duration(t) * time.Second
	now := time.Now()
	timeOuttime = now.Add(timeOut)
	fmt.Println("timeOutTime: ", timeOut)
}

//checkTimeout verify timeout has not occurred and reset timeout if not, calls lock_all() if timeout has passed
func checkTimeout() {
	if timeOut != tstampMax {
		now := time.Now()
		if exp := now.After(timeOuttime); exp {
			// lockAll()
			fmt.Println("wallet has been locked,please unlock firstly") //TODO
		}
		timeOuttime = now.Add(timeOut)
	}
}

func OwnAndUseWallet() {
	fmt.Println("own and use wallet")
}

func Create() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("wallet creating")
		checkTimeout()

		var name string
		_ = json.NewDecoder(r.Body).Decode(&name)

		password, _ := genPassword()

		file, err := os.Open(dir)
		defer file.Close()
		if err != nil {
			fmt.Println(err)
		}
		allwallets, err := file.Readdirnames(-1)
		if err != nil {
			fmt.Println(err)
		}

		walletName := name + fileExt
		for _, f := range allwallets {
			if f == walletName {
				errResp := fmt.Sprintf("Wallet with name: %s already exists at %s", walletName, dir)
				http.Error(w, errResp, 500)
				return
			}
		}

		var wallet SoftWallet
		err = wallet.SetPassword(password)
		if err != nil {
			fmt.Println(err)
		}
		walletFileName := fmt.Sprintf("%s/%s%s", dir, name, fileExt)
		wallet.SetWalletFilename(walletFileName)
		err = wallet.UnLock(password)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		wallet.Lock()
		wallet.UnLock(password)
		wallet.SaveWalletFile()

		if _, ok := wallets[name]; ok {
			delete(wallets, name)
		}
		wallets[name] = wallet

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(password)

	}
	return http.HandlerFunc(fn)
}

func Open() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		checkTimeout()
		var walletname string
		_ = json.NewDecoder(r.Body).Decode(&walletname)

		fmt.Println("Opening wallet :   wallet name: ", walletname)
		var wallet SoftWallet
		walletFileName := fmt.Sprintf("%s/%s%s", dir, walletname, fileExt)
		wallet.SetWalletFilename(walletFileName)
		if !wallet.LoadWalletFile() {
			errResp := fmt.Sprintf("Unable to open file: %s", walletFileName)
			http.Error(w, errResp, 500)
			return
		}
		if _, ok := wallets[walletname]; ok {
			delete(wallets, walletname)
		}
		wallets[walletname] = wallet
		// fmt.Println(walletname, wallet.wallet.CipherKeys)
	}
	return http.HandlerFunc(fn)
}

func ListWallets() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("list wallets")
		var result []string
		for name, wallet := range wallets {
			if wallet.isLocked() {
				result = append(result, name)
			} else {
				result = append(result, name+"*")
			}
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(result)
	}
	return http.HandlerFunc(fn)
}

func ListKeys() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		checkTimeout()
		fmt.Println("list keys")
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)
		name := inputs[0]
		pw := inputs[1]

		if _, ok := wallets[name]; !ok {
			errResp := fmt.Sprintf("Wallet not found: %s", name)
			http.Error(w, errResp, 500)
			return
		}
		wallet := wallets[name]
		if wallet.isLocked() {
			errResp := fmt.Sprintf("Wallet is locked: %s", name)
			http.Error(w, errResp, 500)
			return
		}
		err := wallet.CheckPassword(pw)
		if err != nil {
			http.Error(w, "Invalid password for wallet", 500)
			return
		}

		Resp := RespKeys{}
		for pub, pri := range wallet.Keys {
			Resp[pub] = pri
			// fmt.Println(pub, wallet.Keys[pub])
		}

		w.WriteHeader(201)
		json.NewEncoder(w).Encode(Resp)

	}
	return http.HandlerFunc(fn)
}

func GetPublicKeys() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get public keys")
		var out []string
		if len(wallets) == 0 {
			http.Error(w, "You don't have any wallet", 500)
			return
		}

		isAllWalletLocked := true
		for name, wallet := range wallets {
			if !wallet.isLocked() {
				isAllWalletLocked = false
				fmt.Printf("wallet: %s is unlocked\n", name)
				for pubkey, _ := range wallet.Keys {
					out = append(out, pubkey.String())
				}
			}
		}
		if isAllWalletLocked {
			http.Error(w, "You don't have any unlocked wallet!", 500)
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
		lockAll()
	}
	return http.HandlerFunc(fn)
}
func lockAll() {
	for _, wallet := range wallets {
		if !wallet.isLocked() {
			wallet.Lock()
		}
	}
}

func Lock() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("lock wallet")
		var name string
		_ = json.NewDecoder(r.Body).Decode(&name)

		if _, ok := wallets[name]; !ok {
			errResp := fmt.Sprintf("Wallet not found: %s", name)
			http.Error(w, errResp, 500)
			return
		}
		wallet := wallets[name]
		wallet.Lock()

		w.WriteHeader(201)
		w.Write([]byte("{TODO}"))

	}
	return http.HandlerFunc(fn)
}

func UnLock() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		checkTimeout()
		var wallet SoftWallet
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)
		walletname := inputs[0]
		password := inputs[1]
		fmt.Println("unlock wallet", walletname, password)

		if _, ok := wallets[walletname]; !ok {
			// open(){
			walletFileName := fmt.Sprintf("%s/%s%s", dir, walletname, fileExt)
			wallet.SetWalletFilename(walletFileName)
			if !wallet.LoadWalletFile() {
				errResp := fmt.Sprintf("Unable to open file: %s", walletFileName)
				http.Error(w, errResp, 500)
				return
			}
			if _, ok := wallets[walletname]; ok {
				delete(wallets, walletname)
			}
			wallets[walletname] = wallet
			// }
		}

		wallet = wallets[walletname]
		if !wallet.isLocked() {
			errResp := fmt.Sprintf("Wallet is already unlocked: %s", walletname)
			http.Error(w, errResp, 500)
			return
		}

		err := wallet.UnLock(password)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		delete(wallets, walletname)
		wallets[walletname] = wallet

		// for pub, pri := range wallet.Keys {
		// 	fmt.Println(pub, pri, wallet.Keys[pub])
		// }
	}
	return http.HandlerFunc(fn)
}

func ImportKey() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var inputs []string
		_ = json.NewDecoder(r.Body).Decode(&inputs)
		name := inputs[0]
		wifkey := inputs[1]

		fmt.Println("wallet import keys", name, wifkey)

		wallet, ok := wallets[name]
		if !ok {
			errResp := fmt.Sprintf("Wallet not found: %s\n", name)
			http.Error(w, errResp, 500)
			return
		}

		if wallet.isLocked() {
			errResp := fmt.Sprintf("Wallet is locked: %s\n", name)
			http.Error(w, errResp, 500)
			return
		}

		ok, err := wallet.ImportKey(wifkey)
		if err != nil {
			http.Error(w, "Unable to import key", 500)
			return
		}
		if ok {
			wallet.SaveWalletFile()
		}
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
		w.Write([]byte("{TODO}")) //TODO

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
		fmt.Println(string(inputs[0]), string(inputs[1]), string(inputs[2]))
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

		// for
		// 		//TODO
		// 		signed, err := keyBag.Sign(tx, chainID, requiredKeys...)
		// 		if err != nil {
		// 			http.Error(w, fmt.Sprintf("error signing: %s", err), 500)
		// 			return
		// 		}

		// 		w.WriteHeader(201)
		// 		_ = json.NewEncoder(w).Encode(signed)

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
		fmt.Println(10)

	}
	return http.HandlerFunc(fn)
}
