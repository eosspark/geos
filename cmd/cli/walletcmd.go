package main

import (
	"encoding/hex"
	// "bytes"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/cmd/cli/utils"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"gopkg.in/urfave/cli.v1"
)

var (
	walletCommand = cli.Command{
		Name:        "wallet",
		Usage:       "manage EOS presalse wallets",
		ArgsUsage:   "SUBCOMMAND",
		Category:    "WALLET COMMANDS",
		Description: `Interact with local wallet`,
		Subcommands: []cli.Command{
			{
				Name:     "create",
				Usage:    "Create a new wallet locally",
				Action:   createWallet,
				Category: "WALLET COMMANDS",
				Flags: []cli.Flag{
					utils.WalletNameCreateFlag,
				},
				Description: `create wallet`,
			},
			{
				Name:     "open",
				Usage:    "Open an existing wallet",
				Action:   openWallet,
				Category: "WALLET COMMANDS",
				Flags: []cli.Flag{
					utils.WalletNameOpenFlag,
				},
				Description: `Open an existing wallet`,
			},
			{
				Name:     "lock",
				Usage:    "Lock wallet",
				Action:   openWallet,
				Category: "WALLET COMMANDS",
				Flags: []cli.Flag{
					utils.WalletNameLockFlag,
				},
				Description: `Lock wallet`,
			},
			{
				Name:        "lock_all",
				Usage:       "Lock allwallet",
				Action:      lockAllWallet,
				Category:    "WALLET COMMANDS",
				Description: `Lock all unlocked wallets`,
			},
			{
				Name:     "unlock",
				Usage:    "Unlock Wallet",
				Action:   unlockWallet,
				Category: "WALLET COMMANDS",
				Flags: []cli.Flag{
					utils.WalletNameUnlockFlag,
					utils.WalletPasswordFlag,
				},
				Description: `Unlock wallet`,
			},
			{
				Name:     "import",
				Usage:    "Import private key into wallet",
				Action:   importWallet,
				Category: "WALLET COMMANDS",
				Flags: []cli.Flag{
					utils.WalletNameImportFlag,
					utils.WalletPriKeyFlag,
				},
				Description: `Unlock wallet`,
			},
			{
				Name:     "remove_key",
				Usage:    "Import private key into wallet",
				Action:   removeKey,
				Category: "WALLET COMMANDS",
				Flags: []cli.Flag{
					utils.WalletNameRemoveKeyFlag,
					utils.WalletRemovePriKeyFlag,
					utils.WalletPasswordFlag,
				},
				Description: `Remove key from wallet`,
			},
			{
				Name:        "list",
				Usage:       "List opened wallets, * = unlocked",
				Action:      listWallet,
				Category:    "WALLET COMMANDS",
				Description: `List opened wallets, * = unlocked"`,
			},
			{
				Name:        "keys",
				Usage:       "List of public keys from all unlocked wallets.",
				Action:      getPublicKeys,
				Category:    "WALLET COMMANDS",
				Description: `List of public keys from all unlocked wallets.`,
			},
			{
				Name:     "private_keys",
				Usage:    "List of private keys from an unlocked wallet in wif or PVT_R1 format.",
				Action:   listKeys,
				Category: "WALLET COMMANDS",
				Flags: []cli.Flag{
					utils.WalletNameListKeysFlag,
					utils.WalletPasswordFlag,
				},
				Description: `List of private keys from an unlocked wallet in wif or PVT_R1 format.`,
			},
		},
	}

	SignCommand = cli.Command{
		Name:        "sign",
		Usage:       "Sign a transaction",
		ArgsUsage:   "transaction",
		Category:    "SIGN COMMANDS",
		Description: `Sign a transaction`,
		Action:      signTransactionCli,
		Flags: []cli.Flag{
			utils.TrxJsonToSignFlag,
			utils.StrPrivateKeyFlag,
			utils.StrChainIDFlag,
			utils.PushTrxFlag,
		},
	}
)

func createWallet(ctx *cli.Context) (err error) {
	walletname := ctx.String("name")

	var resp string
	err = DoHttpCall(walletCreate, []string{walletname}, &resp)
	if err != nil {
		return
	}
	fmt.Println("Creating wallet: ", walletname)
	fmt.Println("Save password to use in the future to unlock this wallet.")
	fmt.Println("Without password imported keys will not be retrievable.")
	fmt.Println(resp)
	return
}

func openWallet(ctx *cli.Context) (err error) {
	walletname := ctx.String("name")
	fmt.Println("Opened: ", walletname)

	var resp string
	err = DoHttpCall(walletOpen, []string{walletname}, &resp)
	if err != nil {
		return
	}
	fmt.Println(resp)
	return
}

func lockWallet(ctx *cli.Context) (err error) {
	walletname := ctx.String("name")
	fmt.Println("Locked: ", walletname)
	return nil
}

func lockAllWallet(ctx *cli.Context) (err error) {
	fmt.Println("Locked All Wallet")
	return nil
}

func unlockWallet(ctx *cli.Context) (err error) {
	walletname := ctx.String("name") //utils.WalletUnlockFlag.Name
	password := ctx.String(utils.WalletPasswordFlag.Name)
	fmt.Println(walletname, password)
	err = DoHttpCall(walletUnlock, []string{walletname, password}, nil)
	if err != nil {
		return
	}
	fmt.Println("Unlocked: ", walletname)
	return nil
}
func importWallet(ctx *cli.Context) (err error) {
	walletname := ctx.String("name") //utils.WalletUnlockFlag.Name
	keywif := ctx.String(utils.WalletPriKeyFlag.Name)

	err = DoHttpCall(walletImportKey, []string{walletname, keywif}, nil)
	if err != nil {
		return
	}
	prikey, err := ecc.NewPrivateKey(keywif)
	if err != nil {
		err = fmt.Errorf("Invalid private key: %s", keywif)
		return err
	}
	pubkey := prikey.PublicKey()
	fmt.Println("imported private key for: ", pubkey.String())
	return
}
func removeKey(ctx *cli.Context) (err error) {
	walletname := ctx.String("name")
	keywif := ctx.String("prikey")
	password := ctx.String("password")
	prikey, err := ecc.NewPrivateKey(keywif)
	if err != nil {
		err = fmt.Errorf("Invalid private key: %s", keywif)
		return err
	}

	pubkey := prikey.PublicKey()
	fmt.Println(walletname, password)
	fmt.Println("removed private key for: ", pubkey.String())
	return nil
}

func listWallet(ctx *cli.Context) (err error) {
	var resp []string
	err = DoHttpCall(walletList, nil, &resp)
	if err != nil {
		return
	}
	fmt.Println("wallets: ")
	for _, wallet := range resp {
		fmt.Println(wallet)
	}
	return
}
func getPublicKeys(ctx *cli.Context) (err error) {
	var resp []string
	err = DoHttpCall(walletPublicKeys, nil, &resp)
	if err != nil {
		return
	}
	fmt.Println(resp)
	return
}

func listKeys(ctx *cli.Context) (err error) {
	walletname := ctx.String("name")
	password := ctx.String("password")

	fmt.Println(walletname, password)
	var resp map[string]string
	err = DoHttpCall(walletListKeys, []string{walletname, password}, &resp)
	if err != nil {
		return
	}
	fmt.Println(resp)
	return
}

func signTransactionCli(ctx *cli.Context) error {
	fmt.Println("sign transaction")
	trx_json_to_sign := ctx.String("transaction")
	str_private_key := ctx.String("private-key")
	str_chain_id := ctx.String("chain-id")
	push_trx := ctx.Bool("push-transaction")

	fmt.Println("cli body: ", trx_json_to_sign, str_private_key, str_chain_id, push_trx)

	// SignTransaction()
	var out WalletSignTransactionResp
	err := DoHttpCall(walletSignTrx, []interface{}{
		trx_json_to_sign,
		str_private_key,
		str_chain_id,
	}, &out)
	fmt.Println(out)
	return err
}

func SignTransaction(tx *types.SignedTransaction, chainID common.ChainIDType, pubKeys ...ecc.PublicKey) (out *WalletSignTransactionResp, err error) {
	textKeys := make([]string, 1)
	for _, key := range pubKeys {
		textKeys = append(textKeys, key.String())
	}
	chainid, _ := chainID.MarshalJSON()
	err = DoHttpCall(walletSignTrx, []interface{}{
		tx,
		textKeys,
		hex.EncodeToString(chainid),
	}, &out)
	fmt.Println(out)
	return
}
