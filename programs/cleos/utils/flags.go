package utils

import (
	"github.com/eosspark/eos-go/common"
	"gopkg.in/urfave/cli.v1"
	"strings"
)

const (
	DefaultHTTPHost = "localhost" // Default host interface for the HTTP RPC server
	DefaultHTTPPort = 8888        // Default TCP port for the HTTP RPC server
	DefaultWSHost   = "localhost" // Default host interface for the websocket RPC server
	DefaultWSPort   = 8888        // Default TCP port for the websocket RPC server
)

var (
	WalletNameCreateFlag = cli.StringFlag{
		Name:  "name,n",
		Usage: "The name of the new wallet",
		Value: "default",
	}
	WalletNameOpenFlag = cli.StringFlag{
		Name:  "name,n",
		Usage: "The name of the wallet to open",
		Value: "default",
	}
	WalletNameLockFlag = cli.StringFlag{
		Name:  "name,n",
		Usage: "The name of the wallet to lock",
		Value: "default",
	}
	WalletNameUnlockFlag = cli.StringFlag{
		Name:  "name,n",
		Usage: "The name of the wallet to unlock",
	}
	WalletPasswordFlag = cli.StringFlag{
		Name:  "password",
		Usage: "The password returned by wallet create",
	}

	WalletNameImportFlag = cli.StringFlag{
		Name:  "name,n",
		Usage: "The name of the wallet to import key into",
	}
	WalletPriKeyFlag = cli.StringFlag{
		Name:  "prikey",
		Usage: "Private key in WIF format to import",
	}

	WalletNameRemoveKeyFlag = cli.StringFlag{
		Name:  "name,n",
		Usage: "The name of the wallet to remove key from",
	}
	WalletRemovePriKeyFlag = cli.StringFlag{
		Name:  "prikey",
		Usage: "Private key in WIF format to remove",
	}
	WalletNameListKeysFlag = cli.StringFlag{
		Name:  "name,n",
		Usage: "The name of the wallet to list keys from",
	}
)

var (
	AccountcreateorFlag = cli.StringFlag{
		Name:  "creator",
		Usage: "The name of the account creating the new account",
		// Value: "default",
	}
	AccountNewaccountFlag = cli.StringFlag{
		Name:  "name,n",
		Usage: "The name of the new account",
		// Value: "default",
	}
	AccountOwnerKeyFlag = cli.StringFlag{
		Name:  "ownerkey",
		Usage: "The owner public key for the new account",
		// Value: "default",
	}
	AccountActiveKeyFlag = cli.StringFlag{
		Name:  "activekey",
		Usage: "The active public key for the new account",
		// Value: "default",
	}
)

var (
	BlockHeadStateFlag = cli.BoolFlag{
		Name:  "header-state",
		Usage: "Get block header state from fork database instead",
	}
	PrintJSONFlag = cli.BoolFlag{
		Name:  "json,j",
		Usage: "Output in JSON format",
	}

	CodeFileNameFlag = cli.StringFlag{
		Name:  "code,c",
		Usage: "The name of the file to save the contract .wast/wasm to",
	}
	AbiFileNameFlag = cli.StringFlag{
		Name:  "abi,b",
		Usage: "The name of the file to save the contract .abi to",
	}
	CodeAsWasmFlag = cli.BoolFlag{
		Name:  "wasm",
		Usage: "Save contract as wasm",
	}
	ContractOwnerFlag = cli.StringFlag{
		Name:  "contract",
		Usage: "The contract who owns the table",
	}
	ContractScopeFlag = cli.StringFlag{
		Name:  "scope",
		Usage: "The scope within the contract in which the table is found",
	}
	ContractTableFlag = cli.StringFlag{
		Name:  "table",
		Usage: "The name of the table as specified by the contract abi",
	}
	ContractTableKeyFlag = cli.StringFlag{
		Name:  "key,k",
		Usage: "Deprecated",
	}
	ContractLowerFlag = cli.StringFlag{
		Name:  "lower,L",
		Usage: "JSON representation of lower bound value of key, defaults to first",
	}
	ContractUpperFlag = cli.StringFlag{
		Name:  "upper,U",
		Usage: "JSON representation of upper bound value value of key, defaults to last",
	}
	ContractIndexPositonFlag = cli.StringFlag{
		Name:  "index",
		Usage: "Index number, 1 - primary (first), 2 - secondary index (in order defined by multi_index), 3 - third index, etc.\n\t\tNumber or name of index can be specified, e.g. 'secondary' or '2'.",
	}
	ContractKeyTypeFlag = cli.StringFlag{
		Name:  "key-type",
		Usage: "The key type of --index, primary only supports (i64), all others support (i64, i128, i256, float64, float128).\n\t\tSpecial type 'name' indicates an account name.",
	}
	ContractLimitFlag = cli.IntFlag{
		Name:  "limt,l",
		Usage: "The maximum number of rows to return",
		Value: 10,
	}
	ContractBinaryFlag = cli.BoolFlag{
		Name:  "binary,b",
		Usage: "Return the value as BINARY rather than using abi to interpret as JSON",
	}
	CurrencyCodeFlag = cli.StringFlag{
		Name:  "contract",
		Usage: "The contract that operates the currency",
	}
	CurrencyAccountNameFlag = cli.StringFlag{
		Name:  "account",
		Usage: "The account to query balances for",
	}
	CurrencySymbolFlag = cli.StringFlag{
		Name:  "symbol",
		Usage: "The symbol for the currency if the contract operates multiple currencies",
	}
)

var (
	TrxJsonToSignFlag = cli.StringFlag{
		Name:  "transaction",
		Usage: "The JSON string or filename defining the transaction to sign",
	}
	StrPrivateKeyFlag = cli.StringFlag{
		Name:  "private-key,k",
		Usage: "The private key that will be used to sign the transaction",
	}
	StrChainIdFlag = cli.StringFlag{
		Name:  "chain-id,c",
		Usage: "The chain id that will be used to sign the transaction",
	}
	PushTrxFlag = cli.BoolFlag{
		Name:  "push-transaction,p",
		Usage: "Push transaction after signing",
	}
)

var (
	OpenFileFlag = cli.StringFlag{
		Name:  "x",
		Usage: "The symbol for the currency if the contract operates multiple currencies",
	}
)

var (
	// RPC settings
	RPCEnabledFlag = cli.BoolFlag{
		Name:  "rpc",
		Usage: "Enable the HTTP-RPC server",
	}
	RPCListenAddrFlag = cli.StringFlag{
		Name:  "rpcaddr",
		Usage: "HTTP-RPC server listening interface",
		Value: DefaultHTTPHost,
	}
	RPCPortFlag = cli.IntFlag{
		Name:  "rpcport",
		Usage: "HTTP-RPC server listening port",
		Value: DefaultHTTPPort,
	}
	RPCCORSDomainFlag = cli.StringFlag{
		Name:  "rpccorsdomain",
		Usage: "Comma separated list of domains from which to accept cross origin requests (browser enforced)",
		Value: "",
	}
	RPCVirtualHostsFlag = cli.StringFlag{
		Name:  "rpcvhosts",
		Usage: "Comma separated list of virtual hostnames from which to accept requests (server enforced). Accepts '*' wildcard.",
		//Value: strings.Join(node.DefaultConfig.HTTPVirtualHosts, ","),
	}
	RPCApiFlag = cli.StringFlag{
		Name:  "rpcapi",
		Usage: "API's offered over the HTTP-RPC interface",
		Value: "",
	}
	IPCDisabledFlag = cli.BoolFlag{
		Name:  "ipcdisable",
		Usage: "Disable the IPC-RPC server",
	}
	IPCPathFlag = DirectoryFlag{
		Name:  "ipcpath",
		Usage: "Filename for IPC socket/pipe within the datadir (explicit paths escape it)",
	}
	WSEnabledFlag = cli.BoolFlag{
		Name:  "ws",
		Usage: "Enable the WS-RPC server",
	}
	WSListenAddrFlag = cli.StringFlag{
		Name:  "wsaddr",
		Usage: "WS-RPC server listening interface",
		Value: DefaultWSHost,
	}
	WSPortFlag = cli.IntFlag{
		Name:  "wsport",
		Usage: "WS-RPC server listening port",
		Value: DefaultWSPort,
	}
	WSApiFlag = cli.StringFlag{
		Name:  "wsapi",
		Usage: "API's offered over the WS-RPC interface",
		Value: "",
	}
	WSAllowedOriginsFlag = cli.StringFlag{
		Name:  "wsorigins",
		Usage: "Origins from which to accept websockets requests",
		Value: "",
	}
	ExecFlag = cli.StringFlag{
		Name:  "exec",
		Usage: "Execute JavaScript statement",
	}
	PreloadJSFlag = cli.StringFlag{
		Name:  "preload",
		Usage: "Comma separated list of JavaScript files to preload into the console",
	}
)

var (
	// ATM the url is left to the user and deployment to
	JSpathFlag = cli.StringFlag{
		Name:  "jspath",
		Usage: "JavaScript root path for `loadScript`",
		Value: ".",
	}

	// General settings
	DataDirFlag = DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		//Value: DirectoryString{node.DefaultDataDir()},
	}
)

// MigrateFlags sets the global flag from a local flag when it's set.
// This is a temporary function used for migrating old command/flags to the
// new format.
//
// e.g. geth account new --keystore /tmp/mykeystore --lightkdf
//
// is equivalent after calling this method with:
//
// geth --keystore /tmp/mykeystore --lightkdf account new
//
// This allows the use of the existing configuration functionality.
// When all flags are migrated this function can be removed and the existing
// configuration functionality must be changed that is uses local flags
// func MigrateFlags(action func(ctx *cli.Context) error) func(*cli.Context) error {
// 	return func(ctx *cli.Context) error {
// 		for _, name := range ctx.FlagNames() {
// 			fmt.Println("util", name, ctx.String(name))
// 			if ctx.IsSet(name) {
// 				fmt.Println("util", name, ctx.String(name))
// 				ctx.GlobalSet(name, ctx.String(name))
// 			}
// 		}
// 		return action(ctx)
// 	}
// }

// MakeConsolePreloads retrieves the absolute paths for the console JavaScript
// scripts to preload before starting.
func MakeConsolePreloads(ctx *cli.Context) []string {
	// Skip preloading if there's nothing to preload
	if ctx.GlobalString(PreloadJSFlag.Name) == "" {
		return nil
	}
	// Otherwise resolve absolute paths and return them
	preloads := []string{}

	assets := ctx.GlobalString(JSpathFlag.Name)
	for _, file := range strings.Split(ctx.GlobalString(PreloadJSFlag.Name), ",") {
		preloads = append(preloads, common.AbsolutePath(assets, strings.TrimSpace(file)))
	}
	return preloads
}
