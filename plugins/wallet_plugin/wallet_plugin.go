package wallet_plugin

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	. "github.com/eosspark/eos-go/plugins/appbase/app"
	"github.com/eosspark/eos-go/libraries/asio"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"runtime"
)

const WalletPlug = PluginTypeName("WalletPlugin")

var walletPlugin Plugin = App().RegisterPlugin(WalletPlug, NewWalletPlugin(App().GetIoService()))

type WalletPlugin struct {
	AbstractPlugin
	//ConfirmedBlock Signal //TODO signal ConfirmedBlock
	walletManager *WalletManager
}

func NewWalletPlugin(io *asio.IoContext) *WalletPlugin {
	return &WalletPlugin{}
}

func (w *WalletPlugin) SetProgramOptions(options *[]cli.Flag) {
	*options = append(*options,
		cli.StringFlag{
			Name:  "wallet-dir",
			Usage: "The path of the wallet files (absolute path or relative to application data dir)",
			Value: ".",
		},
		cli.IntFlag{
			Name: "unlock-timeout",
			Usage: "Timeout for unlocked wallet in seconds (default 900 (15 minutes))." +
				"Wallets will automatically lock after specified number of seconds of inactivity.Activity is defined as any wallet command e.g. list-wallets.",
			Value: 900,
		},
		cli.StringFlag{
			Name:  "yubihsm-url",
			Usage: "Override default URL of http://localhost:12345 for connecting to yubihsm-connector)",
			Value: "URL",
		},
		cli.UintFlag{
			Name:  "yubihsm-authkey",
			Usage: "Enables YubiHSM support using given Authkey",
			//Value: nil, //"key_num" TODO
		},
	)

}

func (w *WalletPlugin) PluginInitialize(c *cli.Context) {
	Try(func() {
		w.walletManager = walletManager()

		walletDir := common.AbsolutePath(getWalletDir(), c.String("wallet-dir"))
		w.walletManager.SetDir(walletDir)

		if c.IsSet("unlock-timeout") {
			timeout := c.Int64("unlock-timeout")
			EosAssert(timeout > 0, &exception.InvalidLockTimeoutException{}, "Please specify a positive timeout %d", timeout)
			w.walletManager.SetTimeOut(timeout)
		}

		//if c.IsSet("yubihsm-authkey") {
		//	key := uint16(c.Uint("yubihsm-authkey"))
		//	connectorEndpoint := "http://localhost:12345"
		//	if c.IsSet("yubihsm-url") {
		//		connectorEndpoint = c.String("yubihsm-url")
		//	}
		//	Try(func() {
		//		//w.my.ownAndUseWallet("YubiHSM",)
		//	}).FcLogAndRethrow().End()
		//}

	}).FcLogAndRethrow().End()
}

func (w *WalletPlugin) PluginStartup() {

}

func (w *WalletPlugin) PluginShutdown() {

}

func (w *WalletPlugin) GetWalletManager() *WalletManager {
	return w.walletManager
}

func getWalletDir() string {
	home := os.Getenv("HOME")
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "eosgo_wallet")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "eosgo_wallet")
		} else {
			return filepath.Join(home, "eosgo_wallet")
		}
	}
	return "."
}
