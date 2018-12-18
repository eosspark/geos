package console

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/console_plugin/console/js/eosapi"
	"github.com/robertkrimen/otto"
	"time"
)

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
func (e *eosgo) GetInfo(call otto.FunctionCall) (response otto.Value) {
	var result eosapi.InfoResp
	err := e.c.client.Call(&result, "/v1/chain/get_info", nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n\n", result)

	v, _ := call.Otto.ToValue(result.String())
	return v
}

type msg struct {
	Msg string
}

func (e *eosgo) CreateAccount(call otto.FunctionCall) (resonse otto.Value) {
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

	if len(activekey) == 0 {
		activekey = ownerkey
	}

	ownerKey, err := ecc.NewPublicKey(ownerkey)
	if err != nil {
		throwJSException(fmt.Sprintf("Invalid owner public key: %s", ownerkey))
	}
	activeKey, err := ecc.NewPublicKey(activekey)
	if err != nil {
		throwJSException(fmt.Sprintf("Invalid owner public key: %s", activekey))
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

	//e.SendActions([]types.Action{action})
	createAction := []*types.Action{action}
	//if !simple {
	//	fmt.Println("system create account")
	//
	//} else {
	e.log.Info("creat account in test net")
	e.sendActions(createAction, 1000, types.CompressionNone)
	//}

	v, _ := call.Otto.ToValue(&msg{})

	return v
}

func (e *eosgo) sendActions(actions []*types.Action, extraKcpu int32, compression types.CompressionType) {
	e.log.Info("send action")
	result := e.pushActions(actions, extraKcpu, compression)

	//if txPrintJson {
	//fmt.Println("txPrintJson")
	//fmt.Println(string(result))
	//} else {
	printResult(result)
	//}
}

func (e *eosgo) pushActions(actions []*types.Action, extraKcpu int32, compression types.CompressionType) interface{} {
	e.log.Info("push actions")
	trx := &types.SignedTransaction{}
	trx.Actions = actions
	e.log.Info("trx.Actions")
	return e.pushTransaction(trx, extraKcpu, compression)
}

var txExpiration time.Duration = 30 * time.Second
var abi_serializer_max_time = 10 * time.Second // No risk to client side serialization taking a long time
var tx_ref_block_num_or_id string
var tx_force_unique = false
var tx_dont_broadcast = false
var tx_return_packed = false
var tx_skip_sign = false
var tx_print_json = false
var print_request = false
var print_response = false
var no_auto_keosd = false

var tx_max_cpu_usage uint8 = 0
var tx_max_net_usage uint32 = 0

var delaysec uint32 = 0

type Variants map[string]interface{}

func (e *eosgo) pushTransaction(trx *types.SignedTransaction, extraKcpu int32, compression types.CompressionType) interface{} {
	e.log.Info("push transaction")

	var info eosapi.InfoResp
	err := e.c.client.Call(&info, "/v1/chain/get_info", nil)
	if err != nil {
		fmt.Println(err)
	}
	e.log.Info("receive getInfo()")

	if len(trx.Signatures) == 0 { // #5445 can't change txn content if already signed
		// calculate expiration date
		tx_expiration := info.HeadBlockTime.AddUs(common.Microseconds(txExpiration.Seconds()))
		trx.Expiration = common.TimePointSec(tx_expiration)
		// fmt.Println(trx.Expiration)

		// Set tapos, default to last irreversible block if it's not specified by the user
		refBlockID := info.LastIrreversibleBlockID
		if len(tx_ref_block_num_or_id) > 0 {
			fmt.Println("tx_ref_block_num_or_id")
			var refBlock eosapi.BlockResp
			err := e.c.client.Call(&refBlock, "/v1/chain/get_block", Variants{"block_num_or_id": tx_ref_block_num_or_id})
			if err != nil {
				fmt.Println(err)
				try.EosThrow(&exception.InvalidRefBlockException{}, "Invalid reference block num or id: %s", tx_ref_block_num_or_id)
			}
			e.log.Info("receive getInfo()")
			refBlockID = refBlock.ID
		}
		fmt.Println(refBlockID)
		trx.SetReferenceBlock(&refBlockID)

		if tx_force_unique {
			// trx.ContextFreeActions. //TODO
		}
		trx.MaxCpuUsageMS = uint8(tx_max_cpu_usage)
		trx.MaxNetUsageWords = (uint32(tx_max_net_usage) + 7) / 8
		trx.DelaySec = uint32(delaysec)
	}
	e.log.Info("end %#v", trx)

	if !tx_skip_sign {
		requiredKeys := e.determineRequiredKeys(trx)
		fmt.Println(requiredKeys)
		e.signTransaction(trx, requiredKeys, &info.ChainID)
	}
	if !tx_dont_broadcast {
		fmt.Println("push transaction")
		var re Variants
		packedTrx := types.NewPackedTransactionBySignedTrx(trx, compression)
		err := e.c.client.Call(&re, "/v1/chain/push_transaction", packedTrx)
		if err != nil {
			fmt.Println(err)
		}
		return re
	} else {
		if !tx_return_packed {
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
	err := e.c.client.Call(&publicKeys, "/v1/wallet/get_public_keys", nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("get public keys: ", publicKeys)

	var keys map[string][]string
	fmt.Println("action data:", trx.Actions[0])
	arg := &Variants{
		"transaction":    trx,
		"available_keys": publicKeys,
	}
	err = e.c.client.Call(&keys, "/v1/chain/get_required_keys", arg)
	if err != nil {
		fmt.Println(err)
	}
	return keys["required_keys"]
}

//void sign_transaction(signed_transaction& trx, fc::variant& required_keys, const chain_id_type& chain_id) {
//fc::variants sign_args = {fc::variant(trx), required_keys, fc::variant(chain_id)};
//const auto& signed_trx = call(wallet_url, wallet_sign_trx, sign_args);
//trx = signed_trx.as<signed_transaction>();
//}

func (e *eosgo) signTransaction(trx *types.SignedTransaction, requiredKeys []string, chainID *common.ChainIdType) {
	signedTrx := Variants{"signed_transaction": trx, "keys": requiredKeys, "id": chainID}
	err := e.c.client.Call(trx, "/v1/wallet/sign_transaction", signedTrx)
	if err != nil {
		fmt.Println(err)
	}
}

func printResult(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Println("%v\n", string(data))
	}
}

func (e *eosgo) PushTrx(call otto.FunctionCall) (resonse otto.Value) {
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

//{"ref_block_num":"101","ref_block_prefix":"4159312339","expiration":"2017-09-25T06:28:49","scope":["initb","initc"],"actions":[{"code":"currency","type":"transfer","recipients":["initb","initc"],"authorization":[{"account":"initb","permission":"active"}],"data":"000000000041934b000000008041934be803000000000000"}],"signatures":[],"authorizations":[]}, {"ref_block_num":"101","ref_block_prefix":"4159312339","expiration":"2017-09-25T06:28:49","scope":["inita","initc"],"actions":[{"code":"currency","type":"transfer","recipients":["inita","initc"],"authorization":[{"account":"inita","permission":"active"}],"data":"000000008040934b000000008041934be803000000000000"}],"signatures":[],"authorizations":[]}]

//func (e *eosgo) importWallet(call otto.FunctionCall) (response otto.Value) {
//	walletName, err := call.Argument(0).ToString()
//	if err != nil {
//		return otto.UndefinedValue()
//	}
//	walletKeyStr, err := call.Argument(0).ToString()
//	if err != nil {
//		return otto.UndefinedValue()
//	}
//
//	walletKey, err := ecc.NewPrivateKey(walletKeyStr)
//	if err != nil {
//		try.EosThrow(&exception.PrivateKeyTypeException{}, "Invalid private key: %s", walletKeyStr)
//	}
//
//	err := e.c.client.Call(trx, "/v1/wallet/import_key", Variants{"name":walletName,"key":walletKeyStr})
//	if err != nil {
//		fmt.Println(err)
//	}else{
//		pubkey := walletKey.PublicKey()
//		fmt.Println("imported private key for: ", pubkey.String())
//
//	}
//
//}
