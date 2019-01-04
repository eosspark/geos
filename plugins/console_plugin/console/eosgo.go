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

//var txExpiration = common.Seconds(30)         //30s
var abiSerializerMaxTime = common.Seconds(10) // No risk to client side serialization taking a long time
//var txRefBlockNumOrID string
//var txForceUnique = false
//var txDontBroadcast = false
//var txReturnPacked = false
//var txSkipSign = false
//var txPrintJson = false
//var printRequest = false
//var printResponse = false

//var txMaxCpuUsage uint8 = 0
//var txMaxNetUsage uint32 = 0

//var delaySec uint32 = 0

var clog log.Logger

func init() {
	clog = log.New("console")
	clog.SetHandler(log.TerminalHandler)
	clog.SetHandler(log.DiscardHandler())
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
		action.Authorization = []types.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}}
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

type eosgo struct {
	c *Console
}

func newEosgo(c *Console) *eosgo {
	e := &eosgo{
		c: c,
	}
	return e
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

type CreateAccountParams struct {
	Creator   common.Name `json:"creator"`
	Name      common.Name `json:"name"`
	OwnerKey  string      `json:"owner"`
	ActiveKey string      `json:"active"`
	StandardTransactionOptions
}

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
	re := sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return getJsResult(call, re)
}

type PushAction struct {
	ContractAccount string `json:"account"`
	Action          string `json:"action"`
	Data            string `json:"data"`
	StandardTransactionOptions
}

func (e *eosgo) PushAction(call otto.FunctionCall) (response otto.Value) {
	var params PushAction
	readParams(&params, call)

	actionArgsVar := &common.Variants{}
	err := json.Unmarshal([]byte(params.Data), actionArgsVar)
	if err != nil {
		throwJSException(fmt.Sprintln("Fail to parse action JSON data = ", params.Data))
	}

	permissions := getAccountPermissions(params.TxPermission)
	fmt.Println(permissions)

	action := &types.Action{
		Account:       common.N(params.ContractAccount),
		Name:          common.N(params.Action),
		Authorization: permissions,
		Data:          variantToBin(common.N(params.ContractAccount), common.N(params.Action), actionArgsVar),
	}
	result := sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return getJsResult(call, result)
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

type SetCodeParams struct {
	Account                string `json:"account"`
	ContractPath           string `json:"code_file"`
	ContractClear          bool   `json:"clear"`
	SuppressDuplicateCheck bool   `json:"suppress_duplicate_check"`
	StandardTransactionOptions
}

func (e *eosgo) SetCode(call otto.FunctionCall) (response otto.Value) {
	var params SetCodeParams
	readParams(&params, call)

	codeContent, err := ioutil.ReadFile(params.ContractPath)
	if err != nil {
		clog.Error("get abi from file is error %s", err.Error())
		return otto.FalseValue()
	}

	action := createSetCode(common.N(params.Account), codeContent, params.TxPermission)
	clog.Info("Setting Code...")
	re := sendActions([]*types.Action{action}, 10000, types.CompressionZlib, &params)
	return getJsResult(call, re)
}

type SetAbiParams struct {
	Account                string `json:"account"`
	AbiPath                string `json:"abi_file"`
	ContractClear          bool   `json:"clear"`
	SuppressDuplicateCheck bool   `json:"suppress_duplicate_check"`
	StandardTransactionOptions
}

func (e *eosgo) SetAbi(call otto.FunctionCall) (response otto.Value) {
	var params SetAbiParams
	readParams(&params, call)

	abiFile, err := ioutil.ReadFile(params.AbiPath)
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
	action := createSetABI(common.N(params.Account), abiContent, params.TxPermission)
	clog.Info("Setting ABI...")
	result := sendActions([]*types.Action{action}, 10000, types.CompressionZlib, &params)
	return getJsResult(call, result)
}

func (e *eosgo) SetContract(call otto.FunctionCall) (response otto.Value) {

	v, _ := call.Otto.ToValue(nil)
	return v
}

//TODO set
type SetAccountPermissionParams struct {
	Account             string `json:"account"`
	Permission          string `json:"permission"`
	AuthorityJsonOrFile string `json:"authority"`
	Parent              string `json:"parent"`
	StandardTransactionOptions
}

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

			var itr types.Permission
			var i int
			for i, itr = range accountResult.Permissions {
				if itr.PermName == params.Permission {
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
			parent = common.N(params.Parent)
		}
		action := createUpdateAuth(account, permission, parent, auth, params.TxPermission)
		sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	}
	return getJsResult(call, nil)
}

type SetActionPermissionParams struct {
	Account     string `json:"account"`
	Code        string `json:"code"`
	TypeStr     string `json:"type"`
	Requirement string `json:"requirement"`
	StandardTransactionOptions
}

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
	return getJsResult(call, nil)
}
