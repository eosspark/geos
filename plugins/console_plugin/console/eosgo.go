package console

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

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
	"github.com/tidwall/gjson"
)

var abiSerializerMaxTime = common.Seconds(10) // No risk to client side serialization taking a long time
var clog log.Logger

func init() {
	clog = log.New("console")
	clog.SetHandler(log.TerminalHandler)
	//clog.SetHandler(log.DiscardHandler())
}

type eosgo struct {
	c *Console
}

func newEosgo(c *Console) *eosgo {
	e := &eosgo{c: c}
	return e
}

//CreateKey creates a new keypair and print the public and private keys
func (e *eosgo) CreateKey(call otto.FunctionCall) (response otto.Value) {
	privateKey, _ := ecc.NewRandomPrivateKey()
	fmt.Println("Private key: ", privateKey.String())
	fmt.Println("Public key: ", privateKey.PublicKey().String())
	return
}

//CreateAccount creates a new account on the blockchain (assumes system contract does not restrict RAM usage)
func (e *eosgo) CreateAccount(call otto.FunctionCall) (response otto.Value) {
	var params CreateAccountParams
	readParams(&params, call)

	if len(params.ActiveKey) == 0 {
		params.ActiveKey = params.OwnerKey
	}

	ownerKey, err := ecc.NewPublicKey(params.OwnerKey)
	if err != nil {
		throwJSException(fmt.Sprintf("Invalid owner public key: %s\n", params.OwnerKey))
	}
	activeKey, err := ecc.NewPublicKey(params.ActiveKey)
	if err != nil {
		throwJSException(fmt.Sprintf("Invalid active public key: %s\n", params.OwnerKey))
	}

	action := createNewAccount(params.Creator, params.Name, ownerKey, activeKey, params.TxPermission)

	clog.Info("creat account in test net")
	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}

func (e *eosgo) SetCode(call otto.FunctionCall) (response otto.Value) {
	var params SetContractParams
	readParams(&params, call)

	e.setCodeCallBack(&params)
	return
}

func (e *eosgo) SetAbi(call otto.FunctionCall) (response otto.Value) {
	var params SetContractParams
	readParams(&params, call)

	e.setAbiCallBack(&params)
	return
}

//SetContract creates or update the code on an account
func (e *eosgo) SetContract(call otto.FunctionCall) (response otto.Value) {
	var params SetContractParams
	readParams(&params, call)

	e.setCodeCallBack(&params)
	e.setAbiCallBack(&params)

	return
}

func (e *eosgo) setAbiCallBack(params *SetContractParams) {
	abiFile, err := ioutil.ReadFile(params.AbiPath)
	if err != nil {
		clog.Error("get abi from file is error %s", err.Error())
		return
	}

	abiDef := &abi_serializer.AbiDef{}
	if json.Unmarshal(abiFile, abiDef) != nil {
		clog.Error("unmarshal abi from file is error ")
		return
	}

	abiContent, err := rlp.EncodeToBytes(abiDef)
	if err != nil {
		clog.Error("pack abi is error %s", err.Error())
		return
	}
	action := createSetABI(common.N(params.Account), abiContent, params.TxPermission)
	clog.Info("Setting ABI...")
	sendActions([]*types.Action{action}, 10000, types.CompressionZlib, params)

	abiCache[common.N(params.Account)] = abi_serializer.NewAbiSerializer(abiDef, abiSerializerMaxTime) //for resolve abi
}

func (e *eosgo) setCodeCallBack(params *SetContractParams) {
	codeContent, err := ioutil.ReadFile(params.ContractPath)
	if err != nil {
		clog.Error("get abi from file is error %s", err.Error())
		return
	}

	action := createSetCode(common.N(params.Account), codeContent, params.TxPermission)
	clog.Info("Setting Code...")
	sendActions([]*types.Action{action}, 10000, types.CompressionZlib, params)
}

//Transfer transfers EOS from account to account
func (e *eosgo) Transfer(call otto.FunctionCall) (response otto.Value) {
	con := "eosio.token"
	var params TransferParams
	readParams(&params, call)

	if params.TxForceUnique && len(params.Memo) == 0 {
		// use the memo to add a nonce
		params.Memo = generateNonceString()
		params.TxForceUnique = false
	}

	transferAmount := toAsset(common.N(con), params.Amount)
	transfer := createTransfer(con, common.N(params.Sender), common.N(params.Recipient), *transferAmount, params.Memo, params.TxPermission)
	if !params.PayRam {
		sendActions([]*types.Action{transfer}, 1000, types.CompressionNone, &params)
	} else {
		open := createOpen(con, common.N(params.Recipient), transferAmount.Symbol, common.N(params.Sender), params.TxPermission)
		sendActions([]*types.Action{open}, 1000, types.CompressionNone, &params)
	}
	return
}

//TODO set
//SetAccountPermission sets parameters dealing with account permissions
func (e *eosgo) SetAccountPermission(call otto.FunctionCall) (response otto.Value) {
	var params SetAccountPermissionParams
	readParams(&params, call)

	account := common.N(params.Account)
	permission := common.N(params.Permission)
	isDelete := strings.Compare(params.AuthorityJsonOrFile, "null") == 0
	if isDelete {
		action := createDeleteAuth(account, permission, params.TxPermission)
		sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	} else {
		auth := parseJsonAuthorityOrKey(params.AuthorityJsonOrFile)
		var parent common.Name
		if len(params.Parent) == 0 && strings.Compare(params.Permission, "owner") != 0 {
			//see if we can auto-determine the proper parent
			var accountResult chain_plugin.GetAccountResult
			err := DoHttpCall(&accountResult, common.GetAccountFunc, common.Variants{"account_name": params.Account})
			if err != nil {
				Throw(err.Error())
			}

			var itr chain_plugin.Permission
			var i int
			for i, itr = range accountResult.Permissions {
				if itr.PermName == permission {
					break
				}
			}
			if i != len(accountResult.Permissions) {
				parent = itr.Parent
			} else {
				//if this is a new permission and there is no parent we default to "active"
				parent = common.DefaultConfig.ActiveName
			}
		} else {
			parent = common.N(params.Parent)
		}
		action := createUpdateAuth(account, permission, parent, auth, params.TxPermission)
		sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	}
	return
}

//SetActionPermission sets parameters dealing with account permissions
func (e *eosgo) SetActionPermission(call otto.FunctionCall) (response otto.Value) {
	var params SetActionPermissionParams
	readParams(&params, call)

	accountName := common.N(params.Account)
	codeName := common.N(params.Code)
	typeName := common.N(params.TypeStr)
	isDelete := strings.Compare(params.Requirement, "null") == 0
	if isDelete {
		action := createUnlinkAuth(accountName, codeName, typeName, params.TxPermission)
		sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	} else {
		requirementName := common.N(params.Requirement)
		action := createLinkAuth(accountName, codeName, typeName, requirementName, params.TxPermission)
		sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	}
	return
}

func getAccountPermissions(permissions []string) []common.PermissionLevel {
	accountPermissions := make([]common.PermissionLevel, 0)

	for _, str := range permissions {
		pieces := strings.Split(str, "@")
		if len(pieces) == 1 {
			pieces = append(pieces, "active")
		}
		permission := common.PermissionLevel{
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

func binToVariant(account common.AccountName, action common.ActionName, actionArgs []byte) interface{} {
	abis := abisSerializerResolver(account)
	FcAssert(!common.Empty(abis), fmt.Sprintf("No ABI found %s", account))
	actionType := abis.GetActionType(action)
	FcAssert(len(actionType) != 0, fmt.Sprintf("Unknown action %s in contract %s", action, account))
	return abis.BinaryToVariantPrint(actionType, actionArgs, abiSerializerMaxTime, false)
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
		Data:    buffer,
	}
	if len(txPermission) == 0 {
		action.Authorization = []common.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}}
	} else {
		action.Authorization = getAccountPermissions(txPermission)
	}
	return action
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
		action.Authorization = []common.PermissionLevel{{Actor: ramPayer, Permission: common.DefaultConfig.ActiveName}}
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
		action.Authorization = []common.PermissionLevel{{Actor: sender, Permission: common.DefaultConfig.ActiveName}}
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
		action.Authorization = []common.PermissionLevel{{Actor: account, Permission: common.DefaultConfig.ActiveName}}
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
		action.Authorization = []common.PermissionLevel{{Actor: account, Permission: common.DefaultConfig.ActiveName}}
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
		action.Authorization = []common.PermissionLevel{{Actor: account, Permission: common.DefaultConfig.ActiveName}}
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
		action.Authorization = []common.PermissionLevel{{Actor: account, Permission: common.DefaultConfig.ActiveName}}
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
		action.Authorization = []common.PermissionLevel{{Actor: account, Permission: common.DefaultConfig.ActiveName}}
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
		action.Authorization = []common.PermissionLevel{{Actor: account, Permission: common.DefaultConfig.ActiveName}}
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
	assetCache[a] = common.Symbol{Precision: 4, Symbol: "SYS"}
}

func toAsset(code common.AccountName, s string) *common.Asset {
	var expectedSymbol common.Symbol
	a := common.Asset{}.FromString(&s)
	sym := a.Symbol.ToSymbolCode()
	symStr := a.Name()

	if len(assetCache) == 0 {
		newAssetCache()
	}

	asset, ok := assetCache[AssetPair{code, sym}]
	if !ok {
		var resp map[string]chain_plugin.GetCurrencyStatsResult
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
	var status string
	var net int64 = -1
	var cpu int64 = -1

	data, _ := json.Marshal(v)
	processed := gjson.GetBytes(data, "Processed")
	if processed.Exists() {
		transactionID := gjson.GetBytes(data, "Processed.ID").String()
		receipt := gjson.GetBytes(data, "Processed.Receipt")
		if receipt.IsObject() {
			status = gjson.GetBytes(data, "Processed.Receipt.status").String()
			net = gjson.GetBytes(data, "Processed.Receipt.net_usage_words").Int() * 8
			cpu = gjson.GetBytes(data, "Processed.Receipt.cpu_usage_us").Int()
		} else {
			status = "failed"
			return
		}

		first := fmt.Sprint(status, " transaction: ", transactionID, "  ", net, "bytes  ", cpu, " us")
		fmt.Println(first)

		if status == "failed" {
			softExcept := gjson.GetBytes(data, "Processed.Except")
			if softExcept.Exists() {
				fmt.Println("soft_except:", softExcept.String())
			}
		} else {
			actions := gjson.GetBytes(data, "Processed.ActionTraces")
			if actions.IsArray() {
				results := actions.Array()
				for _, a := range results {
					printActionTree(a)
				}
			}
		}
		fmt.Println("\x1b[1;33m \rwarning: transaction executed locally, but may not be confirmed by the network yet \x1b[0m")
	} else {
		fmt.Println(string(data))
	}
}

func printActionTree(action gjson.Result) {
	printAction(action)
	inlineTraces := action.Get("InlineTraces")
	if inlineTraces.IsArray() {
		re := inlineTraces.Array()
		for _, inlineTrace := range re {
			printActionTree(inlineTrace)
		}
	}
}

func printAction(action gjson.Result) {
	receipt := action.Get("Receipt")
	receiver := receipt.Get("receiver").String()

	act := action.Get("Act")
	codeName := act.Get("account").String()
	funcName := act.Get("name").String()
	data, _ := hex.DecodeString(act.Get("data").String())

	a := binToVariant(common.N(codeName), common.N(funcName), data)
	args, _ := a.([]byte)

	//TODO Parameters should not be sorted ！！
	if len(args) > 100 {
		args = append(args[0:100], []byte("...")...)
	}
	actionName := fmt.Sprintf("%14s <= %-28s", receiver, codeName+"::"+funcName)

	second := fmt.Sprint("#", actionName, string(args))
	fmt.Println(second)

	console := act.Get("console").String()
	if len(console) > 0 {
		re := strings.Fields(console)
		third := fmt.Sprint(">>", re)
		fmt.Println(third)
	}

}
