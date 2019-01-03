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
	. "github.com/eosspark/eos-go/exception/try"
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

var clog log.Logger

func init() {
	clog = log.New("console")
	clog.SetHandler(log.TerminalHandler)
	clog.SetHandler(log.DiscardHandler())
}

type eosgo struct {
	c *Console
}

func newEosgo(c *Console) *eosgo {
	e := &eosgo{
		c: c,
	}
	return e
}

func getAccountPermissions(permissions []string) []types.PermissionLevel {
	accountPermissions := make([]types.PermissionLevel, 0)

	for _, str := range permissions {
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
	FcAssert(!common.Empty(abis), fmt.Sprintf("No ABI found %s", account))

	actionType := abis.GetActionType(action)
	FcAssert(len(actionType) != 0, fmt.Sprintf("Unknown action %s in contract %s", action, account))
	return abis.VariantToBinary(actionType, actionArgsVar, abiSerializerMaxTime)
}

func binToVariant(account common.AccountName, action common.ActionName, actionArgs []byte) common.Variants {
	abis := abisSerializerResolver(account)
	FcAssert(!common.Empty(abis), fmt.Sprintf("No ABI found %s", account))
	actionType := abis.GetActionType(action)
	FcAssert(len(actionType) != 0, fmt.Sprintf("Unknown action %s in contract %s", action, account))
	return abis.BinaryToVariant(actionType, actionArgs, abiSerializerMaxTime, false)
}

func createBuyRam(creator common.Name, newaccount common.Name, quantity *common.Asset, txPermission []string) *types.Action {
	actPayload := common.Variants{
		"payer":    creator.String(),
		"receiver": newaccount.String(),
		"quant":    quantity.String(),
	}
	var auth []types.PermissionLevel
	if len(txPermission) == 0 {
		auth = []types.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}}
	} else {
		auth = getAccountPermissions(txPermission)
	}
	return createAction(auth, common.DefaultConfig.SystemAccountName, common.N("buyram"), &actPayload)
}

func createBuyRamBytes(creator common.Name, newaccount common.Name, numbytes uint32, txPermission []string) *types.Action {
	actPayload := common.Variants{
		"payer":    creator.String(),
		"receiver": newaccount.String(),
		"bytes":    numbytes,
	}
	var auth []types.PermissionLevel
	if len(txPermission) == 0 {
		auth = []types.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}}
	} else {
		auth = getAccountPermissions(txPermission)
	}
	return createAction(auth, common.DefaultConfig.SystemAccountName, common.N("buyrambytes"), &actPayload)
}

func createDelegate(from common.Name, receiver common.Name, net *common.Asset, cpu *common.Asset, transfer bool, txPermission []string) *types.Action {
	actPayLoad := common.Variants{
		"from":               from.String(),
		"receiver":           receiver.String(),
		"stake_net_quantity": net.String(),
		"stake_cpu_quantity": cpu.String(),
		"transfer":           transfer,
	}
	var auth []types.PermissionLevel
	if len(txPermission) == 0 {
		auth = []types.PermissionLevel{{Actor: from, Permission: common.DefaultConfig.ActiveName}}
	} else {
		auth = getAccountPermissions(txPermission)
	}
	return createAction(auth, common.DefaultConfig.SystemAccountName, common.N("delegatebw"), &actPayLoad)
}

func regProducerVariant(producer common.AccountName, key ecc.PublicKey, url string, location uint16) *common.Variants {
	return &common.Variants{
		"producer":     producer,
		"producer_key": key,
		"url":          url,
		"location":     location,
	}
}

func createNewAccount(creator common.Name, newaccount common.Name, owner ecc.PublicKey, active ecc.PublicKey, txPermission []string) *types.Action {
	a := chain.NewAccount{
		Creator: creator,
		Name:    newaccount,
		Owner: types.Authority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{Key: owner, Weight: 1}},
		},
		Active: types.Authority{
			Threshold: 1,
			Keys:      []types.KeyWeight{{Key: active, Weight: 1}},
		},
	}
	buffer, _ := rlp.EncodeToBytes(&a)
	action := &types.Action{
		Account: a.GetAccount(),
		Name:    a.GetName(),
		//Authorization:,
		Data: buffer,
	}
	if len(txPermission) == 0 {
		action.Authorization = []types.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}}
	} else {
		action.Authorization = getAccountPermissions(txPermission)
	}

	return action
}
func createAction(authorization []types.PermissionLevel, code common.AccountName, act common.ActionName, args *common.Variants) *types.Action {
	return &types.Action{
		Account:       code,
		Name:          act,
		Data:          variantToBin(code, act, args),
		Authorization: authorization,
	}
}

func createOpen(contract string, owner common.Name, sym common.Symbol, ramPayer common.Name, txPermission []string) *types.Action {
	open := common.Variants{
		"owner":     owner,
		"symbol":    sym,
		"ram_payer": ramPayer,
	}
	action := &types.Action{
		Account: common.N(contract),
		Name:    common.N("open"),
		Data:    variantToBin(common.N(contract), common.N("open"), &open),
	}
	if len(txPermission) == 0 {
		action.Authorization = []types.PermissionLevel{{Actor: ramPayer, Permission: common.DefaultConfig.ActiveName}}
	} else {
		action.Authorization = getAccountPermissions(txPermission)
	}
	return action
}

func createTransfer(contract string, sender common.Name, recipient common.Name, amount common.Asset, memo string, txPermission []string) *types.Action {
	transfer := common.Variants{
		"from":     sender,
		"to":       recipient,
		"quantity": amount,
		"memo":     memo,
	}
	action := &types.Action{
		Account: common.N(contract),
		Name:    common.N("transfer"),
		Data:    variantToBin(common.N(contract), common.N("transfer"), &transfer),
	}
	if len(txPermission) == 0 {
		action.Authorization = []types.PermissionLevel{{Actor: sender, Permission: common.DefaultConfig.ActiveName}}
	} else {
		action.Authorization = getAccountPermissions(txPermission)
	}
	return action
}

func createSetABI(account common.Name, abi []byte, txPermission []string) *types.Action {
	a := chain.SetAbi{
		Account: account,
		Abi:     abi,
	}
	buffer, _ := rlp.EncodeToBytes(a)
	action := &types.Action{
		Account: a.GetAccount(),
		Name:    a.GetName(),
		Data:    buffer,
	}
	if len(txPermission) == 0 {
		action.Authorization = []types.PermissionLevel{{Actor: account, Permission: common.DefaultConfig.ActiveName}}
	} else {
		action.Authorization = getAccountPermissions(txPermission)
	}
	return action
}

func createSetCode(account common.Name, code []byte, txPermission []string) *types.Action {
	a := chain.SetCode{
		Account:   account,
		VmType:    0,
		VmVersion: 0,
		Code:      code,
	}
	buffer, _ := rlp.EncodeToBytes(a)
	action := &types.Action{
		Account: a.GetAccount(),
		Name:    a.GetName(),
		Data:    buffer,
	}
	if len(txPermission) == 0 {
		action.Authorization = []types.PermissionLevel{{Actor: account, Permission: common.DefaultConfig.ActiveName}}
	} else {
		action.Authorization = getAccountPermissions(txPermission)
	}
	return action
}

func createUpdateAuth(account common.Name, permission common.Name, parent common.Name, auth types.Authority, txPermission []string) *types.Action {
	a := chain.UpdateAuth{
		Account:    account,
		Permission: permission,
		Parent:     parent,
		Auth:       auth,
	}
	buffer, _ := rlp.EncodeToBytes(a)
	action := &types.Action{
		Account: a.GetAccount(),
		Name:    a.GetName(),
		Data:    buffer,
	}
	if len(txPermission) == 0 {
		action.Authorization = []types.PermissionLevel{{Actor: account, Permission: common.DefaultConfig.ActiveName}}
	} else {
		action.Authorization = getAccountPermissions(txPermission)
	}
	return action
}

func createDeleteAuth(account common.Name, permission common.Name, txPermission []string) *types.Action {
	a := chain.DeleteAuth{
		Account:    account,
		Permission: permission,
	}
	buffer, _ := rlp.EncodeToBytes(a)

	action := &types.Action{
		Account: a.GetAccount(),
		Name:    a.GetName(),
		Data:    buffer,
	}
	if len(txPermission) == 0 {
		action.Authorization = []types.PermissionLevel{{Actor: account, Permission: common.DefaultConfig.ActiveName}}
	} else {
		action.Authorization = getAccountPermissions(txPermission)
	}
	return action
}

func createLinkAuth(account common.Name, code common.Name, typeName common.Name, requirement common.Name, txPermission []string) *types.Action {
	a := chain.LinkAuth{
		Account:     account,
		Code:        code,
		Type:        typeName,
		Requirement: requirement,
	}
	buffer, _ := rlp.EncodeToBytes(a)

	action := &types.Action{
		Account: a.GetAccount(),
		Name:    a.GetName(),
		Data:    buffer,
	}
	if len(txPermission) == 0 {
		action.Authorization = []types.PermissionLevel{{Actor: account, Permission: common.DefaultConfig.ActiveName}}
	} else {
		action.Authorization = getAccountPermissions(txPermission)
	}
	return action
}

func createUnlinkAuth(account common.Name, code common.Name, typeName common.Name, txPermission []string) *types.Action {
	a := chain.UnLinkAuth{
		Account: account,
		Code:    code,
		Type:    typeName,
	}
	buffer, _ := rlp.EncodeToBytes(a)

	action := &types.Action{
		Account: a.GetAccount(),
		Name:    a.GetName(),
		Data:    buffer,
	}
	if len(txPermission) == 0 {
		action.Authorization = []types.PermissionLevel{{Actor: account, Permission: common.DefaultConfig.ActiveName}}
	} else {
		action.Authorization = getAccountPermissions(txPermission)
	}
	return action
}

func parseJsonAuthority(authorityJsonOrFile string) (auth types.Authority) {
	Try(func() {
		err := json.Unmarshal([]byte(authorityJsonOrFile), &auth)
		if err != nil {
			Throw(err.Error())
		}
	}).EosRethrowExceptions(&exception.AuthorityTypeException{}, "Fail to parse Authority JSON: %s", authorityJsonOrFile).End()
	return
}

func parseJsonAuthorityOrKey(authorityJsonOrFile string) (result types.Authority) {
	if strings.HasPrefix(authorityJsonOrFile, "EOS") || strings.HasPrefix(authorityJsonOrFile, "PUB_R1") {
		Try(func() {
			pubKey, err := ecc.NewPublicKey(authorityJsonOrFile)
			if err != nil {
				Throw("parse public key from authorityJsonOrFile is error")
			}
			result = types.NewAuthority(pubKey, 0)
		}).EosRethrowExceptions(&exception.PublicKeyTypeException{}, "Invalid public key:%s", authorityJsonOrFile)
	} else {
		result = parseJsonAuthority(authorityJsonOrFile)
		EosAssert(types.Validate(result), &exception.AuthorityTypeException{},
			"Authority failed validation! ensure that keys, accounts, and waits are sorted and that the threshold is valid and satisfiable!")
	}
	return
}

type AssetPair struct {
	Name       common.AccountName
	SymbolCode common.SymbolCode
}

var assetCache = make(map[AssetPair]common.Symbol)

func newAssetCache() {
	a := AssetPair{
		Name:       common.N("eosio.token"),
		SymbolCode: 5462355,
	}
	assetCache[a] = common.Symbol{4, "SYS"}
}

func toAsset(code common.AccountName, s string) *common.Asset {
	var expectedSymbol common.Symbol
	a := common.Asset{}.FromString(&s)
	sym := a.Symbol.ToSymbolCode()
	symStr := a.Name()

	if len(assetCache) == 0 { //TODO get currency stats now is not ready!!!!
		newAssetCache()
	}

	asset, ok := assetCache[AssetPair{code, sym}]
	if !ok {
		var resp chain_plugin.GetCurrencyStatsResult
		err := DoHttpCall(&resp, common.GetCurrencyStatsFunc, common.Variants{
			"json":   false,
			"code":   code,
			"symbol": symStr,
		})

		if err != nil {
			fmt.Println(err)
			Throw("get currency stats is error")
		}
		objIt, ok := resp[symStr]
		if !ok {
			EosThrow(&exception.SymbolTypeException{}, "Symbol %s is not supported by token contract %s", symStr, code)
		}

		assetCache[AssetPair{code, sym}] = objIt.MaxSupply.Symbol
		expectedSymbol = objIt.MaxSupply.Symbol
	} else {
		expectedSymbol = asset
	}

	if a.Decimals() < expectedSymbol.Decimals() {
		factor := expectedSymbol.Precision / a.Precision
		a = *common.NewAssetWithCheck(a.Amount*int64(factor), expectedSymbol)
	} else if a.Decimals() > expectedSymbol.Decimals() {
		EosThrow(&exception.SymbolTypeException{}, "Too many decimal digits in %s, only %d supported", a, expectedSymbol.Decimals())
	}
	return &a
}

func toAssetFromString(s string) *common.Asset {
	return toAsset(common.N("eosio.token"), s)
}

func printResult(v interface{}) {
	data, _ := json.Marshal(v)
	fmt.Println(string(data))
}

func (e *eosgo) CreateKey(call otto.FunctionCall) (response otto.Value) {
	type Keys struct {
		Pri string `json:"Private Key"`
		Pub string `json:"Public Key"`
	}

	privateKey, _ := ecc.NewRandomPrivateKey()
	key := Keys{Pri: privateKey.String(), Pub: privateKey.PublicKey().String()}
	return getJsResult(call, key)
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
			{Actor: common.AccountName(common.N("eosio")), Permission: common.DefaultConfig.ActiveName},
		},
	}

	createAction := []*types.Action{action}
	//if !simple {
	//	fmt.Println("system create account")
	//
	//} else {
	clog.Info("creat account in test net")
	result := sendActions(createAction, 1000, types.CompressionNone)
	//}

	return getJsResult(call, result)
}

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
	permissonStr, err := call.Argument(3).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	actionArgsVar := &common.Variants{}
	json.Unmarshal([]byte(data), actionArgsVar)

	permissions := getAccountPermissions([]string{permissonStr})
	fmt.Println(permissions)
	action := &types.Action{
		Account:       common.N(contractAccount),
		Name:          common.N(actionName),
		Authorization: permissions,
		Data:          variantToBin(common.N(contractAccount), common.N(actionName), actionArgsVar),
	}
	actions := []*types.Action{action}
	result := sendActions(actions, 1000, types.CompressionNone)
	return getJsResult(call, result)
}

func sendActions(actions []*types.Action, extraKcpu int32, compression types.CompressionType) interface{} {
	fmt.Println("send actions...")
	result := pushActions(actions, extraKcpu, compression)

	//if txPrintJson {
	//fmt.Println("txPrintJson")
	//fmt.Println(string(result))
	//} else {
	printResult(result)
	//}
	return result
}

func pushActions(actions []*types.Action, extraKcpu int32, compression types.CompressionType) interface{} {
	trx := &types.SignedTransaction{}
	trx.Actions = actions
	return pushTransaction(trx, extraKcpu, compression)
}

func pushTransaction(trx *types.SignedTransaction, extraKcpu int32, compression types.CompressionType) interface{} {
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
				EosThrow(&exception.InvalidRefBlockException{}, "Invalid reference block num or id: %s", txRefBlockNumOrID)
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
		requiredKeys := determineRequiredKeys(trx)
		signTransaction(trx, requiredKeys, &info.ChainID)
	}
	if !txDontBroadcast {
		var re common.Variant
		packedTrx := types.NewPackedTransactionBySignedTrx(trx, compression)
		err := DoHttpCall(&re, common.PushTxnFunc, packedTrx)
		if err != nil {
			clog.Error(err.Error())
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

func determineRequiredKeys(trx *types.SignedTransaction) []string {
	var publicKeys []string
	err := DoHttpCall(&publicKeys, common.WalletPublicKeys, nil)
	if err != nil {
		clog.Error(err.Error())
	}

	var keys map[string][]string
	arg := &common.Variants{
		"transaction":    trx,
		"available_keys": publicKeys,
	}
	err = DoHttpCall(&keys, common.GetRequiredKeys, arg)
	if err != nil {
		clog.Error(err.Error())
	}
	return keys["required_keys"]
}

func signTransaction(trx *types.SignedTransaction, requiredKeys []string, chainID *common.ChainIdType) {
	signedTrx := common.Variants{"signed_transaction": trx, "keys": requiredKeys, "id": chainID}
	err := DoHttpCall(trx, common.WalletSignTrx, signedTrx)
	if err != nil {
		clog.Error(err.Error())
	}
}

func (e *eosgo) PushTrx(call otto.FunctionCall) (response otto.Value) {
	//var signtrx types.SignedTransaction
	//
	//trx_var, err := call.Argument(0).ToString()
	//if err != nil {
	//	return otto.UndefinedValue()
	//}
	//fmt.Println("receive trx:", trx_var, err)
	//fmt.Println()
	//fmt.Println()
	//aa := "{\"ref_block_num\":\"101\",\"ref_block_prefix\":\"4159312339\",\"expiration\":\"2017-09-25T06:28:49\",\"scope\":[\"initb\",\"initc\"],\"actions\":[{\"code\":\"currency\",\"type\":\"transfer\",\"recipients\":[\"initb\",\"initc\"],\"authorization\":[{\"account\":\"initb\",\"permission\":\"active\"}],\"data\":\"000000000041934b000000008041934be803000000000000\"}],\"signatures\":[],\"authorizations\":[]}, {\"ref_block_num\":\"101\",\"ref_block_prefix\":\"4159312339\",\"expiration\":\"2017-09-25T06:28:49\",\"scope\":[\"inita\",\"initc\"],\"actions\":[{\"code\":\"currency\",\"type\":\"transfer\",\"recipients\":[\"inita\",\"initc\"],\"authorization\":[{\"account\":\"inita\",\"permission\":\"active\"}],\"data\":\"000000008040934b000000008041934be803000000000000\"}],\"signatures\":[],\"authorizations\":[]}]"
	//
	//err = json.Unmarshal([]byte(aa), &signtrx)
	//if err != nil {
	//	fmt.Println(err)
	//	//EOS_RETHROW_EXCEPTIONS(transaction_type_exception, "Fail to parse transaction JSON '${data}'", ("data",trx_to_push))
	//	//try.FcThrowException(&exception.TransactionTypeException{},"Fail to parse transaction JSON %s",trx_var)
	//}
	//
	//re := e.pushTransaction(&signtrx, 1000, types.CompressionNone)
	//printResult(re)

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
		clog.Error("get abi from file is error %s", err.Error())
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
			{Actor: common.N(account), Permission: common.DefaultConfig.ActiveName},
		},
		Data: buffer,
	}
	createAction := []*types.Action{setCode}
	clog.Info("Setting Code...")
	result := sendActions(createAction, 10000, types.CompressionZlib)

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
		clog.Error("get abi from file is error %s", err.Error())
		return otto.FalseValue()
	}

	abiDef := &abi_serializer.AbiDef{}
	if json.Unmarshal(abiFile, abiDef) != nil {
		clog.Error("unmarshal abi from file is error ")
		return otto.FalseValue()
	}

	abiContent, err := rlp.EncodeToBytes(abiDef)
	if err != nil {
		clog.Error("pack abi is error %s", err.Error())
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
			{Actor: common.N(account), Permission: common.DefaultConfig.ActiveName},
		},
		Data: buffer,
	}
	createAction := []*types.Action{setAbi}

	clog.Info("Setting ABI...")
	result := sendActions(createAction, 10000, types.CompressionZlib)
	return getJsResult(call, result)
}

func (e *eosgo) SetContract(call otto.FunctionCall) (response otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}

//TODO set

func (e *eosgo) SetAccountPermission(call otto.FunctionCall) (response otto.Value) {
	accountStr, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	permissionStr, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	authorityJsonOrFile, err := call.Argument(2).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	parentStr, err := call.Argument(3).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	txPermissionStr, err := call.Argument(4).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	account := common.N(accountStr)
	permission := common.N(permissionStr)
	isDelete := strings.Compare(authorityJsonOrFile, "null") == 0
	if isDelete {
		action := createDeleteAuth(account, permission, []string{txPermissionStr})
		sendActions([]*types.Action{action}, 1000, types.CompressionNone)
	} else {
		auth := parseJsonAuthorityOrKey(authorityJsonOrFile)
		var parent common.Name
		if len(parentStr) == 0 && strings.Compare(permissionStr, "owner") != 0 {
			//see if we can auto-determine the proper parent
			var accountResult chain_plugin.GetAccountResult
			err := DoHttpCall(&accountResult, common.GetAccountFunc, common.Variants{"account_name": accountStr})
			if err != nil {
				Throw(err.Error())
			}

			var itr types.Permission
			var i int
			for i, itr = range accountResult.Permissions {
				if itr.PermName == permissionStr {
					break
				}
			}
			if i != len(accountResult.Permissions) {
				parent = common.N(itr.Parent)
			} else {
				//if this is a new permission and there is no parent we default to "active"
				parent = common.DefaultConfig.ActiveName
			}
		} else {
			parent = common.N(parentStr)
		}
		action := createUpdateAuth(account, permission, parent, auth, []string{txPermissionStr})
		sendActions([]*types.Action{action}, 1000, types.CompressionNone)
	}
	return getJsResult(call, nil)
}

func (e *eosgo) SetActionPermission(call otto.FunctionCall) (response otto.Value) {
	accountStr, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	codeStr, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	typeStr, err := call.Argument(2).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	requirementStr, err := call.Argument(3).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	txPermissionStr, err := call.Argument(4).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	accountName := common.N(accountStr)
	codeName := common.N(codeStr)
	typeName := common.N(typeStr)
	isDelete := strings.Compare(requirementStr, "null") == 0
	if isDelete {
		action := createUnlinkAuth(accountName, codeName, typeName, []string{txPermissionStr})
		sendActions([]*types.Action{action}, 1000, types.CompressionNone)
	} else {
		requirementName := common.N(requirementStr)
		action := createLinkAuth(accountName, codeName, typeName, requirementName, []string{txPermissionStr})
		sendActions([]*types.Action{action}, 1000, types.CompressionNone)
	}
	return getJsResult(call, nil)
}

func (e *eosgo) RegisterProducer(call otto.FunctionCall) (response otto.Value) {
	producerStr, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	producerKeyStr, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	url, err := call.Argument(2).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	loc, err := call.Argument(3).ToInteger()
	if err != nil {
		return otto.UndefinedValue()
	}

	var producerKey ecc.PublicKey
	Try(func() {
		producerKey, err = ecc.NewPublicKey(producerKeyStr)
		if err != nil {
			Throw(err.Error())
		}
	}).EosRethrowExceptions(&exception.PublicKeyTypeException{}, "Invalid producer public key: %s", producerKeyStr).End()

	regprodVar := regProducerVariant(common.N(producerStr), producerKey, url, uint16(loc))
	action := createAction([]types.PermissionLevel{{common.N(producerStr), common.DefaultConfig.ActiveName}}, common.DefaultConfig.SystemAccountName, common.N("regproducer"), regprodVar)
	sendActions([]*types.Action{action}, 1000, types.CompressionNone)

	return getJsResult(call, nil)
}
