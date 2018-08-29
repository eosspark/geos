package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/cmd/cli/utils"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"github.com/eosspark/eos-go/rlp"
	"gopkg.in/urfave/cli.v1"
)

var (
	accountCommand = cli.Command{
		Name:        "create",
		Usage:       "Create various items, on and off the blockchain",
		ArgsUsage:   "SUBCOMMAND",
		Category:    "ACCOUNT COMMANDS",
		Description: `Create various items, on and off the blockchain`,
		Subcommands: []cli.Command{
			{
				Name:        "key",
				Usage:       "Create a new keypair and print the public and private keys",
				Action:      createKey,
				Category:    "ACCOUNT COMMANDS",
				Description: `Create a new keypair and print the public and private keys`,
			},
			{
				Name:      "account",
				Usage:     " Create an account, buy ram, stake for bandwidth for the account",
				ArgsUsage: "create account",
				Action:    createAccount,
				Category:  "ACCOUNT COMMANDS",
				Flags: []cli.Flag{
					utils.AccountcreateorFlag,
					utils.AccountNewaccountFlag,
					utils.AccountOwnerKeyFlag,
					utils.AccountActiveKeyFlag,
				},
				Description: `Create an account, buy ram, stake for bandwidth for the account`,
			},
		},
	}

	getCommand = cli.Command{
		Name:        "get",
		Usage:       "Retrieve various items and information from the blockchain",
		ArgsUsage:   "SUBCOMMAND",
		Category:    "GET COMMANDS",
		Description: `Retrieve various items and information from the blockchain`,
		Subcommands: []cli.Command{
			{
				Name:        "info",
				Usage:       "Get current blockchain information",
				Action:      getInfoCli,
				Category:    "GET COMMANDS",
				Description: `Get current blockchain information`,
			},
			{
				Name:      "block",
				Usage:     "Retrieve a full block from the blockchain",
				ArgsUsage: "block",
				Action:    getBlock,
				Category:  "GET COMMANDS",
				Flags: []cli.Flag{
					utils.BlockHeadStateFlag,
				},
				Description: `Retrieve a full block from the blockchain`,
			},
			{
				Name:      "account",
				Usage:     "Retrieve an account from the blockchain",
				ArgsUsage: "name",
				Action:    getAccount,
				Category:  "GET COMMANDS",
				Flags: []cli.Flag{
					utils.PrintJSONFlag,
				},
				Description: `Retrieve an account from the blockchain`,
			},
			{
				Name:      "code",
				Usage:     "Retrieve the code and ABI for an account",
				ArgsUsage: "name",
				Action:    getCode,
				Category:  "GET COMMANDS",
				Flags: []cli.Flag{
					utils.CodeFileNameFlag,
					utils.AbiFileNameFlag,
					utils.CodeAsWasmFlag,
				},
				Description: `Retrieve the code and ABI for an account`,
			},
			{
				Name:      "abi",
				Usage:     "Retrieve the ABI for an account",
				ArgsUsage: "name",
				Action:    getAbi,
				Category:  "GET COMMANDS",
				Flags: []cli.Flag{
					utils.AbiFileNameFlag,
				},
				Description: `Retrieve the code and ABI for an account`,
			},
			{
				Name:      "table",
				Usage:     "Retrieve the contents of a database table",
				ArgsUsage: "contract scope table",
				Action:    getTable,
				Category:  "GET COMMANDS",
				Flags: []cli.Flag{
					utils.ContractOwnerFlag,
					utils.ContractScopeFlag,
					utils.ContractTableFlag,
					utils.ContractBinaryFlag,
					utils.ContractLimitFlag,
					utils.ContractTableKeyFlag,
					utils.ContractLowerFlag,
					utils.ContractUpperFlag,
					utils.ContractIndexPositonFlag,
					utils.ContractKeyTypeFlag,
				},
				Description: `Retrieve the contents of a database table`,
			},
			{
				Name:      "currency",
				Usage:     "Retrieve information related to standard currencies",
				ArgsUsage: "SUBCOMMAND",
				Category:  "GET COMMANDS",
				Subcommands: []cli.Command{
					{
						Name:     "balance",
						Usage:    "Retrieve the balance of an account for a given currency",
						Action:   getCurrencyBalance,
						Category: "GET COMMANDS",
						Flags: []cli.Flag{
							utils.CurrencyCodeFlag,
							utils.CurrencyAccountNameFlag,
							utils.CurrencySymbolFlag,
						},
						Description: `Retrieve the balance of an account for a given currency`,
					},
					{
						Name:     "stats",
						Usage:    "Retrieve the stats of for a given currency",
						Action:   getCurrencyStats,
						Category: "GET COMMANDS",
						Flags: []cli.Flag{
							utils.CurrencyCodeFlag,
							utils.CurrencySymbolFlag,
						},
						Description: `Retrieve the stats of for a given currency`,
					},
				},
			},
		},
	}
)

func createKey(ctx *cli.Context) (err error) {
	prikey, err := ecc.NewRandomPrivateKey()
	if err != nil {
		return err
	}
	fmt.Println("Private Key:", prikey.String())
	fmt.Println("Public Key:", prikey.PublicKey().String())
	return
}

func createAccount(ctx *cli.Context) (err error) {

	creator := ctx.String("creator")
	newname := ctx.String("name")
	ownerkey := ctx.String("ownerkey")
	activekey := ctx.String("activekey")

	if len(activekey) == 0 {
		activekey = ownerkey
	}

	ownerKey, err := ecc.NewPublicKey(ownerkey)
	if err != nil {
		return fmt.Errorf("Invalid owner public key: %s", ownerkey)
	}
	activeKey, err := ecc.NewPublicKey(activekey)
	if err != nil {
		return fmt.Errorf("Invalid active public key: %s", activekey)
	}
	create := createNewAccount(creator, newname, ownerKey, activeKey)

	storage, err := rlp.EncodeToBytes(create)
	fmt.Println("encode: ", storage)
	aa, err := json.Marshal(create)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("decode:", string(aa))

	resp, err := getInfo()
	if err != nil {
		return err
	}
	display, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	fmt.Println(string(display))

	fmt.Println("New account:")
	return
}

// ./accountcmd create account -creator eosio -name walker -ownerkey EOS7vnBoERUwrqeRTfot79xhwFvWsTjhg1YU9KA5hinAYMETREWYT -activekey EOS7vnBoERUwrqeRTfot79xhwFvWsTjhg1YU9KA5hinAYMETREWYT
func createNewAccount(creatorstr, newaccountstr string, owner, active ecc.PublicKey) *types.NewAccount {
	creator := common.StringToName(creatorstr)
	newaccount := common.StringToName(newaccountstr)
	ownerauthority := &common.Authority{
		Threshold: 1,
		Keys:      []common.KeyWeight{{PublicKey: owner, Weight: 1}},
	}
	activeauthority := &common.Authority{
		Threshold: 1,
		Keys:      []common.KeyWeight{{PublicKey: active, Weight: 1}},
	}
	// {PublicKey: active, Weight: 1}
	return &types.NewAccount{
		Creator: common.AccountName(creator),
		Name:    common.AccountName(newaccount),
		Owner:   *ownerauthority,
		Active:  *activeauthority,
	}
}

func getInfoCli(ctx *cli.Context) (err error) {
	resp, err := getInfo()
	if err != nil {
		return err
	}
	display, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	fmt.Println(string(display))
	return
}

func getInfo() (out *InfoResp, err error) {
	_, err = DoHttpCall(getInfoFunc, nil, &out)
	return
}

func getBlock(ctx *cli.Context) (err error) {
	getBHS := ctx.Bool(utils.BlockHeadStateFlag.Name)
	blockarg := ctx.Args().First()
	var resp BlockResp
	var data []byte
	if getBHS {
		data, err = DoHttpCall(getBlockHeaderStateFunc, M{"block_num_or_id": blockarg}, &resp)
	} else {
		data, err = DoHttpCall(getBlockFunc, M{"block_num_or_id": blockarg}, &resp)
	}

	fmt.Println("json: ", bytes.NewBuffer(data).String())
	// display, err := json.Marshal(resp)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("解构以后： ", string(display))
	return

}

func getBlockID(getbhs bool, blockarg string) (resp *BlockResp, err error) {

	var data []byte
	if getbhs {
		data, err = DoHttpCall(getBlockHeaderStateFunc, M{"block_num_or_id": blockarg}, &resp)
	} else {
		data, err = DoHttpCall(getBlockFunc, M{"block_num_or_id": blockarg}, &resp)
	}

	fmt.Println("json: ", bytes.NewBuffer(data).String())
	return
}

func getAccount(ctx *cli.Context) (err error) {
	printJson := ctx.Bool("json")
	name := ctx.Args().First()
	var resp AccountResp
	_, err = DoHttpCall(getAccountFunc, M{"account_name": name}, &resp)

	display, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	fmt.Println(string(display))
	if !printJson {
		//TODO
	} else {

	}
	return
}

func getCode(ctx *cli.Context) (err error) {
	code := ctx.String("code")
	abi := ctx.String("abi")
	wasm := ctx.Bool("wasm")
	name := ctx.Args().First()
	for i := 0; i < ctx.NArg(); i++ {
		fmt.Println(ctx.Args()[i])
	}
	fmt.Println(code, abi, wasm, name)
	var resp GetCodeResp
	_, err = DoHttpCall(getRawCodeAndAbiFunc, M{"account_name": name, "code_as_wasm": true}, &resp)

	if err != nil {
		return err
	}
	display, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	fmt.Println("resp: ", string(display)) //TODO save to file
	return
}

func getAbi(ctx *cli.Context) (err error) {
	abi := ctx.String("abi")
	name := ctx.Args().First()
	for i := 0; i < ctx.NArg(); i++ {
		fmt.Println(ctx.Args()[i])
	}
	fmt.Println(abi, name)
	var resp GetABIResp

	_, err = DoHttpCall(getAbiFunc, M{"account_name": name}, &resp)
	if err != nil {
		return err
	}
	display, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	fmt.Println("resp: ", string(display)) //TODO save to file
	return
}

func getTable(ctx *cli.Context) (err error) {
	binary := ctx.Bool("binary")
	code := ctx.String("contract")
	scope := ctx.String("scope")
	table := ctx.String("table")
	tableKey := ctx.String("key")
	lower := ctx.String("lower")
	upper := ctx.String("upper")
	indexPosition := ctx.String("index")
	keyType := ctx.String("key-type")
	limit := ctx.Int("limt")

	var resp GetTableRowsResp
	_, err = DoHttpCall(getTableFunc,
		M{"json": !binary,
			"code":           code,
			"scope":          scope,
			"table":          table,
			"table_key":      tableKey,
			"lower_bound":    lower,
			"upper_bound":    upper,
			"limit":          limit,
			"key_type":       keyType,
			"index_position": indexPosition},
		&resp)
	if err != nil {
		return err
	}
	display, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	fmt.Println("resp: ", string(display))
	return nil
}
func getCurrencyBalance(ctx *cli.Context) (err error) {
	name := ctx.String("account")
	code := ctx.String("contract")
	symbol := ctx.String("symbol")
	params := M{"account_name": name, "code": code}

	if symbol != "" {
		params["symbol"] = symbol
	}
	var resp []common.Asset
	_, err = DoHttpCall(getCurrencyBalanceFunc, params, &resp)
	if err != nil {
		return err
	}
	display, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	for i := 0; i < len(resp); i++ {
		fmt.Println(resp[i])
	}
	fmt.Println("resp: ", string(display))
	return nil
}
func getCurrencyStats(ctx *cli.Context) (err error) {
	code := ctx.String("contract")
	symbol := ctx.String("symbol")
	var resp json.RawMessage //TODO
	data, err := DoHttpCall(getCurrencyStatsFunc, M{"code": code, "symbol": symbol}, &resp)
	if err != nil {
		return err
	}
	display, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	for i := 0; i < len(resp); i++ {
		fmt.Println(resp[i])
	}
	fmt.Println("resp: ", string(display))
	fmt.Println("json: ", bytes.NewBuffer(data).String())
	return nil
}
