package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/programs/cli/utils"
	"gopkg.in/urfave/cli.v1"
	"time"
)

var tx_expiration time.Duration = 30 * time.Second
var (
	parse_expiration = flag.Duration("-x,--expiration", 30*time.Second, "set the time in seconds before a transaction expires, defaults to 30s")
	tx_force_unique  = flag.Bool("-f,--force-unique", false, "force the transaction to be unique. this will consume extra bandwidth and remove any "+
		"protections against accidently issuing the same transaction multiple times")
	tx_skip_sign = flag.Bool("-s,--skip-sign", false, "Specify if unlocked wallet keys should be used to sign transaction")

	tx_print_json          = flag.Bool("-j,--json", false, "print result as json")
	tx_dont_broadcast      = flag.Bool("-d,--dont-broadcast", false, "don't broadcast transaction to the network (just print to stdout)")
	tx_ref_block_num_or_id = flag.String("-r,--ref-block", "", "set the reference block num or block id used for TAPOS (Transaction as Proof-of-Stake)")
	//tx_permission = flag.String("-p,--permission","","An account and permission level to authorize, as in 'account@permission' (defaults to '" + default_permission + "')")
	tx_max_cpu_usage = flag.Uint("--max-cpu-usage-ms", 0, "set an upper limit on the milliseconds of cpu usage budget, for the execution of the transaction (defaults to 0 which means no limit)")
	tx_max_net_usage = flag.Uint("--max-net-usage", 0, "set an upper limit on the net usage budget, in bytes, for the transaction (defaults to 0 which means no limit)")
	delaysec         = flag.Uint("--delay-sec", 0, "set the delay_sec seconds, defaults to 0s")
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

var txPrintJson bool = false

func createKey(ctx *cli.Context) (err error) {
	prikey, err := ecc.NewRandomPrivateKey()
	if err != nil {
		return err
	}
	fmt.Println("Private Key:", prikey.String())
	fmt.Println("Public Key:", prikey.PublicKey().String())
	return
}

//vector<chain::permission_level> get_account_permissions(const vector<string>& permissions) {
//   auto fixedPermissions = permissions | boost::adaptors::transformed([](const string& p) {
//      vector<string> pieces;
//      split(pieces, p, boost::algorithm::is_any_of("@"));
//      if( pieces.size() == 1 ) pieces.push_back( "active" );
//      return chain::permission_level{ .actor = pieces[0], .permission = pieces[1] };
//   });
//   vector<chain::permission_level> accountPermissions;
//   boost::range::copy(fixedPermissions, back_inserter(accountPermissions));
//   return accountPermissions;
//}

//func getAccountPermission( permissions []string) (accountPermissions []common.PermissionLevel){
//
//}

func createAccount(ctx *cli.Context) (err error) {

	creator := ctx.String("creator")
	accountname := ctx.String("name")
	ownerkey := ctx.String("ownerkey")
	activekey := ctx.String("activekey")
	// stake_net :=
	// stake_cpu :=
	// var buy_ram_bytes_in_kbytes uint32 = 0
	// var buy_ram_eos string
	// var transfer bool = false
	var simple bool = true

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
	create := createNewAccount(creator, accountname, ownerKey, activeKey)
	createAction := []*types.Action{create}
	// storage, err := rlp.EncodeToBytes(create)
	// fmt.Println(create, "encode: ", storage)
	// aa, err := json.Marshal(create)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("decode:", string(aa))

	if !simple {
		fmt.Println("system create account")

	} else {
		fmt.Println("creat account in test net")
		sendActions(createAction, 1000, common.CompressionNone)
	}

	fmt.Println("New account:")
	return
}

// ./accountcmd create account -creator eosio -name walker -ownerkey EOS7vnBoERUwrqeRTfot79xhwFvWsTjhg1YU9KA5hinAYMETREWYT -activekey EOS7vnBoERUwrqeRTfot79xhwFvWsTjhg1YU9KA5hinAYMETREWYT
func createNewAccount(creatorstr, newaccountstr string, owner, active ecc.PublicKey) *types.Action {
	creator := common.S(creatorstr)
	accountName := common.S(newaccountstr)
	ownerauthority := &types.Authority{
		Threshold: 1,
		Keys:      []types.KeyWeight{{Key: owner, Weight: 1}},
	}
	activeauthority := &types.Authority{
		Threshold: 1,
		Keys:      []types.KeyWeight{{Key: active, Weight: 1}},
	}

	var auth = []types.PermissionLevel{{common.AccountName(creator), common.PermissionName(common.DefaultConfig.ActiveName)}} //TODO -p

	newaccount := &types.NewAccount{
		Creator: common.AccountName(creator),
		Name:    common.AccountName(accountName),
		Owner:   *ownerauthority,
		Active:  *activeauthority,
	}

	data, err := rlp.EncodeToBytes(newaccount)
	if err != nil {
		panic("create new account error")
	}

	return &types.Action{
		Account:       newaccount.GetAccount(),
		Name:          newaccount.GetName(),
		Authorization: auth,
		Data:          data,
	}
}
func sendActions(actions []*types.Action, extraKcpu int32, compression common.CompressionType) {
	fmt.Println("send action")
	result := pushActions(actions, extraKcpu, compression)

	if txPrintJson {
		fmt.Println("txPrintJson")
		// fmt.Println(string(result))
	} else {
		printResult(result)
	}

}
func pushActions(actions []*types.Action, extraKcpu int32, compression common.CompressionType) interface{} {
	fmt.Println("push actions")
	trx := &types.SignedTransaction{}
	trx.Actions = actions
	fmt.Println("trx.Actions")
	return pushTransaction(trx, extraKcpu, compression)
}
func pushTransaction(trx *types.SignedTransaction, extraKcpu int32, compression common.CompressionType) interface{} {
	fmt.Println("push transaction")
	info, err := getInfo()
	if err != nil {
		panic(err)
	}
	if len(trx.Signatures) == 0 { // #5445 can't change txn content if already signed

		// fmt.Println(info.HeadBlockTime.Totime(), info.HeadBlockTime.Totime().Add(tx_expiration))
		// trx.SetExpiration(tx_expiration)//now()
		// calculate expiration date
		trx.Expiration = common.JSONTime{info.HeadBlockTime.Totime().Add(tx_expiration)}
		// fmt.Println(trx.Expiration)

		// Set tapos, default to last irreversible block if it's not specified by the user
		refBlockID := info.LastIrreversibleBlockID
		if len(*tx_ref_block_num_or_id) > 0 {
			fmt.Println("tx_ref_block_num_or_id")
			var resp BlockResp
			variant, err := DoHttpCall(chainUrl, getBlockHeaderStateFunc, Variants{"block_num_or_id": tx_ref_block_num_or_id})
			if err != nil {
				panic(err)
			}
			if err := json.Unmarshal(variant, &resp); err != nil {
				return fmt.Errorf("Unmarshal: %s", err)
			}
			refBlockID = resp.ID
		}
		trx.SetReferenceBlock(refBlockID)

		if *tx_force_unique {
			// trx.ContextFreeActions. //TODO
		}
		trx.MaxCPUUsageMS = uint8(*tx_max_cpu_usage)
		trx.MaxNetUsageWords = (uint32(*tx_max_net_usage) + 7) / 8
		trx.DelaySec = uint32(*delaysec)
	}

	if !*tx_skip_sign {
		requiredKeys := determineRequiredKeys(trx)
		fmt.Println(requiredKeys)
		// signTransaction(trx, requiredKeys, info.ChainID)
	}
	if !*tx_dont_broadcast {
		fmt.Println("push transaction")
		return nil
	} else {
		return trx
	}

	return nil
}

func determineRequiredKeys(trx *types.SignedTransaction) Variants {
	var publicKeys []string
	variant, err := DoHttpCall(walletUrl, walletPublicKeys, nil)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(variant, &publicKeys); err != nil {
		return nil
	}
	fmt.Println("get public keys: ", publicKeys)

	var keys map[string][]string
	fmt.Println("action data:", trx.Actions[0])

	arg := &Variants{
		"transaction":    trx,
		"available_keys": publicKeys,
	}
	variant, err = DoHttpCall(chainUrl, getRequiredKeys, arg)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(variant, &keys); err != nil {
		return nil
	}

	return nil
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
	variant, err := DoHttpCall(chainUrl, getInfoFunc, nil)
	if err := json.Unmarshal(variant, &out); err != nil {
		return nil, fmt.Errorf("Unmarshal: %s", err)
	}
	return
}

func getBlock(ctx *cli.Context) (err error) {
	getBHS := ctx.Bool(utils.BlockHeadStateFlag.Name)
	blockarg := ctx.Args().First()
	// var resp BlockResp
	var variant []byte
	if getBHS {
		variant, err = DoHttpCall(chainUrl, getBlockHeaderStateFunc, Variants{"block_num_or_id": blockarg})
	} else {
		variant, err = DoHttpCall(chainUrl, getBlockFunc, Variants{"block_num_or_id": blockarg})
	}
	if err != nil {
		return
	}
	fmt.Println("resp: ", string(variant))

	// if err := json.Unmarshal(variant, &resp); err != nil {
	// 	return fmt.Errorf("Unmarshal: %s", err)
	// }
	// fmt.Println(resp.BlockNumber())

	return

}

// type M map[string]interface{}
// auto arg = fc::mutable_variant_object("block_num_or_id", blockArg);
func getBlockID(getbhs bool, blockarg string) (resp *BlockResp, err error) {
	var variant []byte
	if getbhs {
		variant, err = DoHttpCall(chainUrl, getBlockHeaderStateFunc, Variants{"block_num_or_id": blockarg})
	} else {
		variant, err = DoHttpCall(chainUrl, getBlockFunc, Variants{"block_num_or_id": blockarg})
	}
	if err := json.Unmarshal(variant, &resp); err != nil {
		return nil, fmt.Errorf("Unmarshal: %s", err)
	}
	fmt.Println(resp)

	return
}

func getAccount(ctx *cli.Context) (err error) {
	printJson := ctx.Bool("json")
	name := ctx.Args().First()
	var resp AccountResp
	variant, err := DoHttpCall(chainUrl, getAccountFunc, Variants{"account_name": name})
	if err := json.Unmarshal(variant, &resp); err != nil {
		return fmt.Errorf("Unmarshal: %s", err)
	}
	fmt.Println(string(variant)) //for display

	// fmt.Println(resp.AccountName)
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
	variant, err := DoHttpCall(chainUrl, getRawCodeAndAbiFunc, Variants{"account_name": name, "code_as_wasm": true})
	if err := json.Unmarshal(variant, &resp); err != nil {
		return fmt.Errorf("Unmarshal: %s", err)
	}
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

	variant, err := DoHttpCall(chainUrl, getAbiFunc, Variants{"account_name": name})
	if err != nil {
		return err
	}
	if err := json.Unmarshal(variant, &resp); err != nil {
		return fmt.Errorf("Unmarshal: %s", err)
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
	variant, err := DoHttpCall(chainUrl, getTableFunc,
		Variants{"json": !binary,
			"code":           code,
			"scope":          scope,
			"table":          table,
			"table_key":      tableKey,
			"lower_bound":    lower,
			"upper_bound":    upper,
			"limit":          limit,
			"key_type":       keyType,
			"index_position": indexPosition})
	if err != nil {
		return err
	}
	fmt.Println("resp: ", string(variant))

	if err := json.Unmarshal(variant, &resp); err != nil {
		return fmt.Errorf("Unmarshal: %s", err)
	}

	return nil
}

func getCurrencyBalance(ctx *cli.Context) (err error) {
	name := ctx.String("account")
	code := ctx.String("contract")
	symbol := ctx.String("symbol")
	params := Variants{"account_name": name, "code": code}

	if symbol != "" {
		params["symbol"] = symbol
	}
	var resp []common.Asset
	variant, err := DoHttpCall(chainUrl, getCurrencyBalanceFunc, params)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(variant, &resp); err != nil {
		return fmt.Errorf("Unmarshal: %s", err)
	}
	fmt.Println("resp: ", string(variant))

	for i := 0; i < len(resp); i++ {
		fmt.Println(resp[i])
	}
	return nil
}
func getCurrencyStats(ctx *cli.Context) (err error) {
	code := ctx.String("contract")
	symbol := ctx.String("symbol")
	var resp json.RawMessage //TODO
	variant, err := DoHttpCall(chainUrl, getCurrencyStatsFunc, Variants{"code": code, "symbol": symbol})
	if err != nil {
		return err
	}
	fmt.Println("resp: ", string(variant))
	if err := json.Unmarshal(variant, &resp); err != nil {
		return fmt.Errorf("Unmarshal: %s", err)
	}

	for i := 0; i < len(resp); i++ {
		fmt.Println(resp[i])
	}
	return nil
}

func printResult(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Println("%v\n", string(data))
	}
}

// storage, err := rlp.EncodeToBytes(create)
// fmt.Println("encode: ", storage)
// aa, err := json.Marshal(create)
// if err != nil {
// 	fmt.Println(err)
// }
// fmt.Println("decode:", string(aa))

// resp, err := getInfo()
// if err != nil {
// 	return err
// }
// fmt.Println(resp.HeadBlockID, resp.HeadBlockNum)
// display, err := json.Marshal(resp)
// if err != nil {
// 	return err
// }
// fmt.Println(string(display))
