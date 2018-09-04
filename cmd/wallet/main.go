package main

import (
	"flag"
	"fmt"
	"github.com/eosspark/eos-go/wallet"
	"log"
	"net/http"
)

const (
	walletFuncBase       string = "/v1/wallet"
	walletCreateFunc     string = walletFuncBase + "/create"
	walletOpenFunc       string = walletFuncBase + "/open"
	walletListFunc       string = walletFuncBase + "/list_wallets"
	walletListKeysFunc   string = walletFuncBase + "/list_keys"
	walletPublicKeysFunc string = walletFuncBase + "/get_public_keys"
	walletLockFunc       string = walletFuncBase + "/lock"
	walletLockAllFunc    string = walletFuncBase + "/lock_all"
	walletUnlockFunc     string = walletFuncBase + "/unlock"
	walletImportKeyFunc  string = walletFuncBase + "/import_key"
	walletRemoveKeyFunc  string = walletFuncBase + "/remove_key"
	walletCreateKeyFunc  string = walletFuncBase + "/create_key"
	walletSignTrxFunc    string = walletFuncBase + "/sign_transaction"

	walletSignDigestFunc string = walletFuncBase + "/sign_digest"
	walletSetTimeOutFunc string = walletFuncBase + "/set_timeout"
)

var walletlistenAddress = flag.String("wallet-listen-address", "127.0.0.1:8000", "The local IP and port to listen for incoming http connections;")

func main() {
	flag.Parse()
	done := make(chan bool)
	wallet := http.NewServeMux()
	wallet.Handle(walletSetTimeOutFunc, walletPlugin.SetTimeOut())
	wallet.Handle(walletSignTrxFunc, walletPlugin.SignTransaction())
	wallet.Handle(walletSignDigestFunc, walletPlugin.SignDigest())
	wallet.Handle(walletCreateFunc, walletPlugin.WalletCreate())
	wallet.Handle(walletOpenFunc, walletPlugin.WalletOpen())
	wallet.Handle(walletLockAllFunc, walletPlugin.LockAllwallets())
	wallet.Handle(walletLockFunc, walletPlugin.LockWallet())
	wallet.Handle(walletUnlockFunc, walletPlugin.UnLockWallet())
	wallet.Handle(walletImportKeyFunc, walletPlugin.WalletImportKey())
	wallet.Handle(walletRemoveKeyFunc, walletPlugin.RemoveKey())
	wallet.Handle(walletCreateKeyFunc, walletPlugin.CreateKey())
	wallet.Handle(walletListFunc, walletPlugin.ListWallets())
	wallet.Handle(walletListKeysFunc, walletPlugin.ListKey())
	wallet.Handle(walletPublicKeysFunc, walletPlugin.GetPublicKeys())

	fmt.Printf("Listening for wallet operations on %s\n", *walletlistenAddress)
	err := http.ListenAndServe(*walletlistenAddress, wallet)
	if err != nil {
		log.Println("Litsening failed:", err)
	}
	<-done
}
