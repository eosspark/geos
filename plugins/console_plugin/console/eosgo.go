package console

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"strings"
)

var txExpiration = common.Seconds(30)         //30s
var abiSerializerMaxTime = common.Seconds(10) // No risk to client side serialization taking a long time
var txRefBlockNumOrID string
var txForceUnique = false
var txDontBroadcast = false
var txReturnPacked = false
var txSkipSign = false
var txPrintJson = false
var printRequest = false
var printResponse = false

var txMaxCpuUsage uint8 = 0
var txMaxNetUsage uint32 = 0

var delaySec uint32 = 0

type eosgo struct {
	c   *Console
	log log.Logger
}

func newEosgo(c *Console) *eosgo {
	e := &eosgo{
		c: c,
	}
	e.log = log.New("eosgo")
	e.log.SetHandler(log.TerminalHandler)
	return e
}

type Keys struct {
	Pri string `json:"Private Key"`
	Pub string `json:"Public Key"`
}

func (e *eosgo) CreateKey(call otto.FunctionCall) (response otto.Value) {
	privateKey, _ := ecc.NewRandomPrivateKey()
	v, _ := call.Otto.ToValue(&Keys{Pri: privateKey.String(), Pub: privateKey.PublicKey().String()})
	return v
}

func (e *eosgo) CreateAccount(call otto.FunctionCall) (response otto.Value) {
	creator, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	name, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	ownerkey, err := call.Argument(2).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	activekey, err := call.Argument(3).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	fmt.Println("eosgo params:  ", creator, name, ownerkey, activekey)

	if len(activekey) == 0 {
		activekey = ownerkey
	}

	ownerKey, err := ecc.NewPublicKey(ownerkey)
	if err != nil {
		throwJSException(fmt.Sprintf("Invalid owner public key: %s", ownerkey))
	}
	activeKey, err := ecc.NewPublicKey(activekey)
	if err != nil {
		throwJSException(fmt.Sprintf("Invalid active public key: %s", activekey))
	}

	c := chain.NewAccount{
		Creator: common.AccountName(common.N(creator)),
		Name:    common.AccountName(common.N(name)),
		Owner: types.Authority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{Key: ownerKey, Weight: 1}},
		},
		Active: types.Authority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{Key: activeKey, Weight: 1}},
		},
	}

	buffer, _ := rlp.EncodeToBytes(&c)

	action := &types.Action{
		Account: common.AccountName(common.N("eosio")),
		Name:    common.ActionName(common.N("newaccount")),
		Data:    buffer,
		Authorization: []types.PermissionLevel{
			//types.PermissionLevel{Actor: common.AccountName(common.N("eosio.token")), Permission: common.PermissionName(common.N("active"))},
			{Actor: common.AccountName(common.N("eosio")), Permission: common.PermissionName(common.N("active"))},
		},
	}

	createAction := []*types.Action{action}
	//if !simple {
	//	fmt.Println("system create account")
	//
	//} else {
	e.log.Info("creat account in test net")
	result := e.sendActions(createAction, 1000, types.CompressionNone)
	//}

	return getJsResult(call, result)
}

// // push action
// string contract_account;
// string action;
// string data;
// vector<string> permissions;
// auto actionsSubcommand = push->add_subcommand("action", localized("Push a transaction with a single action"));
// actionsSubcommand->fallthrough(false);
// actionsSubcommand->add_option("account", contract_account,
//                               localized("The account providing the contract to execute"), true)->required();
// actionsSubcommand->add_option("action", action,
//                               localized("A JSON string or filename defining the action to execute on the contract"), true)->required();
// actionsSubcommand->add_option("data", data, localized("The arguments to the contract"))->required();

// add_standard_transaction_options(actionsSubcommand);
// actionsSubcommand->set_callback([&] {
//    fc::variant action_args_var;
//    if( !data.empty() ) {
//       try {
//          action_args_var = json_from_file_or_string(data, fc::json::relaxed_parser);
//       } EOS_RETHROW_EXCEPTIONS(action_type_exception, "Fail to parse action JSON data='${data}'", ("data", data))
//    }
//    auto accountPermissions = get_account_permissions(tx_permission);

//    send_actions({chain::action{accountPermissions, contract_account, action, variant_to_bin( contract_account, action, action_args_var ) }});
// });

func (e *eosgo) PushAction(call otto.FunctionCall) (response otto.Value) {
	contractAccount, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	actionName, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	data, err := call.Argument(2).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	permissonstr, err := call.Argument(3).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	actionArgsVar := &common.Variants{}
	json.Unmarshal([]byte(data), actionArgsVar)

	persmissions := getAccountPermissions(permissonstr)
	fmt.Println(persmissions)
	action := &types.Action{
		Account:       common.N(contractAccount),
		Name:          common.N(actionName),
		Authorization: persmissions,
		Data:          variantToBin(common.N(contractAccount), common.N(actionName), actionArgsVar),
	}
	//e.log.Warn("%#v", action)
	//e.log.Info("push action...")
	actions := []*types.Action{action}
	result := e.sendActions(actions, 1000, types.CompressionNone)
	return getJsResult(call, result)
}

func getAccountPermissions(in string) []types.PermissionLevel {
	accountPermissions := make([]types.PermissionLevel, 0)
	permissionStrs := strings.Split(in, ",")
	for _, str := range permissionStrs {
		pieces := strings.Split(str, "@")
		if len(pieces) == 1 {
			pieces = append(pieces, "active")
		}
		permission := types.PermissionLevel{
			Actor:      common.N(pieces[0]),
			Permission: common.N(pieces[1]),
		}
		accountPermissions = append(accountPermissions, permission)
	}
	return accountPermissions
}

var abiCache = make(map[common.AccountName]*abi_serializer.AbiSerializer)

//resolver for ABI serializer to decode actions in proposed transaction in multisig contract
func abisSerializerResolver(account common.AccountName) *abi_serializer.AbiSerializer {
	it, ok := abiCache[account]
	if ok {
		return it
	}
	var abiResult chain_plugin.GetAbiResult
	err := DoHttpCall(&abiResult, common.GetAbiFunc, common.Variants{"account_name": account})
	if err != nil {
		fmt.Println("get abi from chain is error: ", err.Error())
	}
	var abis *abi_serializer.AbiSerializer
	if !common.Empty(abiResult.Abi) {
		abis = abi_serializer.NewAbiSerializer(&abiResult.Abi, abiSerializerMaxTime)
		abiCache[account] = abis
	} else {
		fmt.Printf("ABI for contract %s not found. Action data will be shown in hex only.\n", account)
	}

	return abis
}

func variantToBin(account common.AccountName, action common.ActionName, actionArgsVar *common.Variants) []byte {
	abis := abisSerializerResolver(account)
	try.FcAssert(!common.Empty(abis), "No ABI found %s", account)

	actionType := abis.GetActionType(action)
	try.FcAssert(len(actionType) != 0, "Unknown action %s in contract %s", action, account)
	return abis.VariantToBinary(actionType, actionArgsVar, abiSerializerMaxTime)
}

func (e *eosgo) sendActions(actions []*types.Action, extraKcpu int32, compression types.CompressionType) interface{} {
	result := e.pushActions(actions, extraKcpu, compression)

	//if txPrintJson {
	//fmt.Println("txPrintJson")
	//fmt.Println(string(result))
	//} else {
	printResult(result)
	//}
	return result
}

func (e *eosgo) pushActions(actions []*types.Action, extraKcpu int32, compression types.CompressionType) interface{} {
	trx := &types.SignedTransaction{}
	trx.Actions = actions
	return e.pushTransaction(trx, extraKcpu, compression)
}

func (e *eosgo) pushTransaction(trx *types.SignedTransaction, extraKcpu int32, compression types.CompressionType) interface{} {
	var info chain_plugin.GetInfoResult
	err := DoHttpCall(&info, common.GetInfoFunc, nil)
	if err != nil {
		fmt.Println(err)
	}

	if len(trx.Signatures) == 0 { // #5445 can't change txn content if already signed
		// calculate expiration date
		trx.Expiration = common.NewTimePointSecTp(info.HeadBlockTime.AddUs(txExpiration))
		fmt.Println(trx.Expiration.String())

		// Set tapos, default to last irreversible block if it's not specified by the user
		refBlockID := info.LastIrreversibleBlockID
		if len(txRefBlockNumOrID) > 0 {
			//var refBlock GetBlockResult
			var refBlock chain_plugin.GetBlockResult
			err := DoHttpCall(&refBlock, common.GetBlockFunc, common.Variants{"block_num_or_id": txRefBlockNumOrID})
			if err != nil {
				fmt.Println(err)
				try.EosThrow(&exception.InvalidRefBlockException{}, "Invalid reference block num or id: %s", txRefBlockNumOrID)
			}

			refBlockID = refBlock.ID
		}
		trx.SetReferenceBlock(&refBlockID)

		if txForceUnique {
			// trx.ContextFreeActions. //TODO
		}
		trx.MaxCpuUsageMS = uint8(txMaxCpuUsage)
		trx.MaxNetUsageWords = (uint32(txMaxNetUsage) + 7) / 8
		trx.DelaySec = uint32(delaySec)
	}

	if !txSkipSign {
		requiredKeys := e.determineRequiredKeys(trx)
		fmt.Println(requiredKeys)
		e.signTransaction(trx, requiredKeys, &info.ChainID)
	}
	if !txDontBroadcast {
		var re common.Variant
		packedTrx := types.NewPackedTransactionBySignedTrx(trx, compression)
		err := DoHttpCall(&re, common.PushTxnFunc, packedTrx)
		if err != nil {
			e.log.Error(err.Error())
		}
		return re
	} else {
		if !txReturnPacked {
			out, _ := json.Marshal(trx)
			return out
		} else {
			out, _ := json.Marshal(types.NewPackedTransactionBySignedTrx(trx, compression))
			return out
		}
	}
}

//type GetRequiredKeysResult struct {
//	RequiredKeys []ecc.PublicKey `json:"required_keys"`
//}
func (e *eosgo) determineRequiredKeys(trx *types.SignedTransaction) []string {
	var publicKeys []string
	err := DoHttpCall(&publicKeys, common.WalletPublicKeys, nil)
	if err != nil {
		e.log.Error(err.Error())
	}
	fmt.Println("get public keys: ", publicKeys)

	var keys map[string][]string
	arg := &common.Variants{
		"transaction":    trx,
		"available_keys": publicKeys,
	}
	err = DoHttpCall(&keys, common.GetRequiredKeys, arg)
	if err != nil {
		e.log.Error(err.Error())
	}
	return keys["required_keys"]
}

func (e *eosgo) signTransaction(trx *types.SignedTransaction, requiredKeys []string, chainID *common.ChainIdType) {
	signedTrx := common.Variants{"signed_transaction": trx, "keys": requiredKeys, "id": chainID}
	err := DoHttpCall(trx, common.WalletSignTrx, signedTrx)
	if err != nil {
		e.log.Error(err.Error())
	}
}

func printResult(v interface{}) {
	data, _ := json.Marshal(v)
	fmt.Println(string(data))
}

func (e *eosgo) PushTrx(call otto.FunctionCall) (response otto.Value) {
	var signtrx types.SignedTransaction

	trx_var, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	fmt.Println("receive trx:", trx_var, err)
	fmt.Println()
	fmt.Println()
	aa := "{\"ref_block_num\":\"101\",\"ref_block_prefix\":\"4159312339\",\"expiration\":\"2017-09-25T06:28:49\",\"scope\":[\"initb\",\"initc\"],\"actions\":[{\"code\":\"currency\",\"type\":\"transfer\",\"recipients\":[\"initb\",\"initc\"],\"authorization\":[{\"account\":\"initb\",\"permission\":\"active\"}],\"data\":\"000000000041934b000000008041934be803000000000000\"}],\"signatures\":[],\"authorizations\":[]}, {\"ref_block_num\":\"101\",\"ref_block_prefix\":\"4159312339\",\"expiration\":\"2017-09-25T06:28:49\",\"scope\":[\"inita\",\"initc\"],\"actions\":[{\"code\":\"currency\",\"type\":\"transfer\",\"recipients\":[\"inita\",\"initc\"],\"authorization\":[{\"account\":\"inita\",\"permission\":\"active\"}],\"data\":\"000000008040934b000000008041934be803000000000000\"}],\"signatures\":[],\"authorizations\":[]}]"

	err = json.Unmarshal([]byte(aa), &signtrx)
	if err != nil {
		fmt.Println(err)
		//EOS_RETHROW_EXCEPTIONS(transaction_type_exception, "Fail to parse transaction JSON '${data}'", ("data",trx_to_push))
		//try.FcThrowException(&exception.TransactionTypeException{},"Fail to parse transaction JSON %s",trx_var)
	}

	re := e.pushTransaction(&signtrx, 1000, types.CompressionNone)
	printResult(re)

	v, _ := call.Otto.ToValue(nil)
	return v
}

func (e *eosgo) SetCode(call otto.FunctionCall) (response otto.Value) {
	account, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	wasmPath, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	codeContent, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		e.log.Error("get abi from file is error %s", err.Error())
		return otto.FalseValue()
	}

	c := chain.SetCode{
		Account:   common.N(account),
		VmType:    0,
		VmVersion: 0,
		Code:      codeContent,
	}
	buffer, _ := rlp.EncodeToBytes(&c)

	setCode := &types.Action{
		Account: common.N("eosio"),
		Name:    common.N("setcode"),
		Authorization: []types.PermissionLevel{
			{Actor: common.N(account), Permission: common.N("active")},
		},
		Data: buffer,
	}
	createAction := []*types.Action{setCode}
	e.log.Info("Setting Code...")
	result := e.sendActions(createAction, 10000, types.CompressionZlib)

	return getJsResult(call, result)

}
func (e *eosgo) SetAbi(call otto.FunctionCall) (response otto.Value) {
	account, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	abiPath, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	abiFile, err := ioutil.ReadFile(abiPath)
	if err != nil {
		e.log.Error("get abi from file is error %s", err.Error())
		return otto.FalseValue()
	}

	fmt.Println(string(abiFile))
	abiDef := &abi_serializer.AbiDef{}
	if json.Unmarshal(abiFile, abiDef) != nil {
		e.log.Error("unmarshal abi from file is error ")
		return otto.FalseValue()
	}

	//if !abi_serializer.ToABI(abiFile,abiDef){
	//	e.log.Error("get abi from file is error ")
	//	return otto.FalseValue()
	//}

	abiContent, err := rlp.EncodeToBytes(abiDef)
	if err != nil {
		e.log.Error("pack abi is error %s", err.Error())
		return otto.FalseValue()
	}

	c := chain.SetAbi{
		Account: common.N(account),
		Abi:     abiContent,
	}
	buffer, _ := rlp.EncodeToBytes(&c)
	setAbi := &types.Action{
		Account: common.N("eosio"),
		Name:    common.N("setabi"),
		Authorization: []types.PermissionLevel{
			{Actor: common.N(account), Permission: common.N("active")},
		},
		Data: buffer,
	}
	createAction := []*types.Action{setAbi}

	e.log.Info("Setting ABI...")
	result := e.sendActions(createAction, 10000, types.CompressionZlib)
	return getJsResult(call, result)
}

func (e *eosgo) SetContract(call otto.FunctionCall) (response otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}
