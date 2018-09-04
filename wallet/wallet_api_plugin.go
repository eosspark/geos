package walletPlugin

import (
	"flag"
)

var (
	walletDir = flag.String("wallet-dir", ".", "The path of the wallet files (absolute path or relative to application data dir)")

	unlockTimeOut = flag.Int("unlock-timeout", 900, "Timeout for unlocked wallet in seconds (default 900 (15 minutes)).Wallets will automatically lock after specified number of seconds of inactivity.Activity is defined as any wallet command e.g. list-wallets.")

	yubihsmUrl = flag.String("yubihsm-url", "URL", "Override default URL of http://localhost:12345 for connecting to yubihsm-connector")

	yubihsmAuthKey = flag.String("yubihsm-authkey", "key_num", "Enables YubiHSM support using given Authkey")
)

type TransactionHandleType uint16

func PluginInitialize() {
	flag.Parse()
	// dir := *walletDir

	// timeout := *unlockTimeOut
	// key := *yubihsmAuthKey
	// connectorEndpoint := "http://localhost:12345"

}
