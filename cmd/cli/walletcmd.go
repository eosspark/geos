package main

import (
	"fmt"
	"github.com/eosspark/eos-go/cmd/cli/utils"
	"github.com/eosspark/eos-go/ecc"
	"gopkg.in/urfave/cli.v1"
)

var (
	walletCommand = cli.Command{
		Name:        "wallet",
		Usage:       "manage EOS presalse wallets",
		ArgsUsage:   "",
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
				Name: "list",
				// Usage:       "List opened wallets, * = unlocked",
				Action:      listWallet,
				Category:    "WALLET COMMANDS",
				Description: `List opened wallets, * = unlocked"`,
			},
			{
				Name: "keys",
				// Usage:       "List of public keys from all unlocked wallets.",
				Action:      listKeys,
				Category:    "WALLET COMMANDS",
				Description: `List of public keys from all unlocked wallets.`,
			},
		},
	}
)

func createWallet(ctx *cli.Context) error {
	walletname := ctx.String("name")
	fmt.Println("Creating wallet: ", walletname)
	fmt.Println("Save password to use in the future to unlock this wallet.")
	fmt.Println("Without password imported keys will not be retrievable.")
	return nil

}

func openWallet(ctx *cli.Context) error {
	walletname := ctx.String("name")
	fmt.Println("Opened: ", walletname)
	return nil
}

func lockWallet(ctx *cli.Context) error {
	walletname := ctx.String("name")
	fmt.Println("Locked: ", walletname)
	return nil
}

func lockAllWallet(ctx *cli.Context) error {
	fmt.Println("Locked All Wallet")
	return nil
}

func unlockWallet(ctx *cli.Context) error {
	walletname := ctx.String("name") //utils.WalletUnlockFlag.Name
	password := ctx.String(utils.WalletPasswordFlag.Name)
	fmt.Println(walletname, password)
	fmt.Println("Unlocked: ", walletname)
	return nil
}
func importWallet(ctx *cli.Context) error {
	walletname := ctx.String("name") //utils.WalletUnlockFlag.Name
	keywif := ctx.String(utils.WalletPriKeyFlag.Name)
	prikey, err := ecc.NewPrivateKey(keywif)
	if err != nil {
		err = fmt.Errorf("Invalid private key: %s", keywif)
		return err
	}

	pubkey := prikey.PublicKey()
	fmt.Println(walletname)
	fmt.Println("imported private key for: ", pubkey.String())
	return nil
}
func removeKey(ctx *cli.Context) error {
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

func listWallet(ctx *cli.Context) error {
	fmt.Println("wallets: ")
	return nil
}
func listKeys(ctx *cli.Context) error {
	fmt.Println("Keys: ")
	return nil
}
