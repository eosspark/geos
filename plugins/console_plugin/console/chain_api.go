package console

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/robertkrimen/otto"
	"github.com/tidwall/gjson"
	"sort"
	"strings"
	"time"
)

type chainAPI struct {
	c *Console
}

func newchainAPI(c *Console) *chainAPI {
	e := &chainAPI{
		c: c,
	}
	return e
}

//GetInfo gets current blockchain information
func (a *chainAPI) GetInfo(call otto.FunctionCall) (response otto.Value) {
	var info chain_plugin.GetInfoResult
	err := DoHttpCall(&info, common.GetInfoFunc, nil)
	if err != nil {
		return getJsResult(call, err.Error())
	}
	return getJsResult(call, info)
}

//GetBlock retrieves a full block from the blockchain
func (a *chainAPI) GetBlock(call otto.FunctionCall) (response otto.Value) {
	txRefBlockNumOrID, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	getBHS, err := call.Argument(1).ToBoolean()
	if err != nil {
		return otto.UndefinedValue()
	}

	arg := common.Variants{"block_num_or_id": txRefBlockNumOrID}
	if getBHS {
		var resp types.BlockHeaderState
		err = DoHttpCall(&resp, common.GetBlockHeaderStateFunc, arg)
		if err == nil {
			return getJsResult(call, resp)
		}
	} else {
		var resp chain_plugin.GetBlockResult
		err = DoHttpCall(&resp, common.GetBlockFunc, arg)
		if err == nil {
			return getJsResult(call, resp)
		}
	}
	return getJsResult(call, err.Error())
}

//GetAccount retrieves an account from the blockchain
func (a *chainAPI) GetAccount(call otto.FunctionCall) (response otto.Value) {
	name, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	var resp chain_plugin.GetAccountResult
	err = DoHttpCall(&resp, common.GetAccountFunc, common.Variants{"account_name": name})
	if err != nil {
		clog.Error("get account is error: %s", err.Error())
		return otto.UndefinedValue()
	}
	PrintAccountResult(&resp)
	//return getJsResult(call, resp)
	return
}

//GetCode retrieves the code and ABI for an account
func (a *chainAPI) GetCode(call otto.FunctionCall) (response otto.Value) { //TODO save to file
	var params GetCodeParams
	readParams(&params, call)

	var resp chain_plugin.GetCodeResult
	err := DoHttpCall(&resp, common.GetCodeFunc, common.Variants{"account_name": params.AccountName, "code_as_wasm": params.CodeAsWasm})
	if err != nil {
		clog.Error("get abi is error: %s", err.Error())
	}
	return getJsResult(call, resp)
}

//GetAbi retrieves the ABI for an account
func (a *chainAPI) GetAbi(call otto.FunctionCall) (response otto.Value) { //TODO save to file
	name, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	var resp chain_plugin.GetAbiResult
	err = DoHttpCall(&resp, common.GetAbiFunc, common.Variants{"account_name": name})
	if err != nil {
		clog.Error("get abi is error: %s", err.Error())
	}
	return getJsResult(call, resp)
}

//GetTable retrieves the contents of a database table
func (a *chainAPI) GetTable(call otto.FunctionCall) (response otto.Value) {
	var params GetTableParams
	readParams(&params, call)

	if params.Limit == 0 {
		params.Limit = 10
	}
	if len(params.EncodeType) == 0 {
		params.EncodeType = "dec"
	}
	var resp chain_plugin.GetTableRowsResult
	err := DoHttpCall(&resp, common.GetTableFunc, common.Variants{
		"json":           !params.Binary,
		"code":           params.Code,
		"scope":          params.Scope,
		"table":          params.Table,
		"table_key":      params.TableKey,
		"lower_bound":    params.Lower,
		"upper_bound":    params.Upper,
		"limit":          params.Limit,
		"key_type":       params.KeyType,
		"index_position": params.IndexPosition})
	if err != nil {
		clog.Error("get abi is error: %s", err.Error())
	}
	return getJsResult(call, resp)
}

//GetScope retrieves a list of scopes and tables owned by a contract
func (a *chainAPI) GetScope(call otto.FunctionCall) (response otto.Value) {
	var params GetScopeParams
	readParams(&params, call)

	if params.Limit == 0 {
		params.Limit = 10
	}

	var resp chain_plugin.GetTableByScopeResultRow
	err := DoHttpCall(&resp, common.GetTableByScopeFunc, common.Variants{
		"code":        params.Code,
		"table":       params.Table,
		"lower_bound": params.Lower,
		"upper_bound": params.Upper,
		"limit":       params.Limit})
	if err != nil {
		clog.Error("get abi is error: %s", err.Error())
	}
	return getJsResult(call, resp)
}

//GetCurrencyBalance retrieves information related to standard currencies
func (a *chainAPI) GetCurrencyBalance(call otto.FunctionCall) (response otto.Value) {
	var code, accountName, symbol string
	code, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	accountName, err = call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	if len(call.ArgumentList) == 2 {
		symbol = ""
	} else {
		symbol, err = call.Argument(2).ToString()
		if err != nil {
			return otto.UndefinedValue()
		}
	}
	var resp []common.Asset
	err = DoHttpCall(&resp, common.GetCurrencyBalanceFunc, common.Variants{"account_name": accountName, "code": code, "symbol": symbol})
	if err != nil {
		clog.Error("GetCurrencyBalance is error: %s", err.Error())
	}

	for i := 0; i < len(resp); i++ {
		fmt.Println(resp[i])
	}
	return getJsResult(call, resp)
}

//GetCurrencyStats retrieve the stats of for a given currency
func (a *chainAPI) GetCurrencyStats(call otto.FunctionCall) (response otto.Value) {
	code, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	symbol, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	var resp map[string]chain_plugin.GetCurrencyStatsResult
	err = DoHttpCall(&resp, common.GetCurrencyStatsFunc, common.Variants{"code": code, "symbol": symbol})
	if err != nil {
		clog.Error("GetCurrencyBalance is error: %s", err.Error())
	}
	return getJsResult(call, resp)
}

////TODO history_plugin
////GetAccounts retrieves accounts associated with a public key
//func (a *chainAPI) GetAccounts(call otto.FunctionCall) (response otto.Value) {
//	return otto.FalseValue()
//}
//
////GetServants retrieves accounts which are servants of a given account
//func (a *chainAPI) GetServants(call otto.FunctionCall) (response otto.Value) {
//	return otto.FalseValue()
//}
//
////GetTransaction retrieves a transaction from the blockchain
//func (a *chainAPI) GetTransaction(call otto.FunctionCall) (response otto.Value) {
//	return otto.FalseValue()
//}
//
////GetActions retrieves all actions with specific account name referenced in authorization or receiver
//func (a *chainAPI) GetActions(call otto.FunctionCall) (response otto.Value) {
//	return otto.FalseValue()
//}

//PushAction pushs a transaction with a single action
func (e *eosgo) PushAction(call otto.FunctionCall) (response otto.Value) {
	var params PushAction
	readParams(&params, call)

	actionArgsVar := &common.Variants{}
	err := json.Unmarshal([]byte(params.Data), actionArgsVar)
	if err != nil {
		throwJSException(fmt.Sprintln("Fail to parse action JSON data = ", params.Data))
	}

	permissions := getAccountPermissions(params.TxPermission)
	action := &types.Action{
		Account:       common.N(params.ContractAccount),
		Name:          common.N(params.Action),
		Authorization: permissions,
		Data:          variantToBin(common.N(params.ContractAccount), common.N(params.Action), actionArgsVar),
	}
	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return otto.UndefinedValue()
}

////PushTransaction pushes an arbitrary JSON transaction
//func (a *chainAPI) PushTransaction(call otto.FunctionCall) (response otto.Value) {
//	//var signtrx types.SignedTransaction
//	//
//	//trx_var, err := call.Argument(0).ToString()
//	//if err != nil {
//	//	return otto.UndefinedValue()
//	//}
//	//fmt.Println("receive trx:", trx_var, err)
//	//fmt.Println()
//	//fmt.Println()
//	//aa := "{\"ref_block_num\":\"101\",\"ref_block_prefix\":\"4159312339\",\"expiration\":\"2017-09-25T06:28:49\",\"scope\":[\"initb\",\"initc\"],\"actions\":[{\"code\":\"currency\",\"type\":\"transfer\",\"recipients\":[\"initb\",\"initc\"],\"authorization\":[{\"account\":\"initb\",\"permission\":\"active\"}],\"data\":\"000000000041934b000000008041934be803000000000000\"}],\"signatures\":[],\"authorizations\":[]}, {\"ref_block_num\":\"101\",\"ref_block_prefix\":\"4159312339\",\"expiration\":\"2017-09-25T06:28:49\",\"scope\":[\"inita\",\"initc\"],\"actions\":[{\"code\":\"currency\",\"type\":\"transfer\",\"recipients\":[\"inita\",\"initc\"],\"authorization\":[{\"account\":\"inita\",\"permission\":\"active\"}],\"data\":\"000000008040934b000000008041934be803000000000000\"}],\"signatures\":[],\"authorizations\":[]}]"
//	//
//	//err = json.Unmarshal([]byte(aa), &signtrx)
//	//if err != nil {
//	//	fmt.Println(err)
//	//	//EOS_RETHROW_EXCEPTIONS(transaction_type_exception, "Fail to parse transaction JSON '${data}'", ("data",trx_to_push))
//	//	//try.FcThrowException(&exception.TransactionTypeException{},"Fail to parse transaction JSON %s",trx_var)
//	//}
//	//
//	//re := e.pushTransaction(&signtrx, 1000, types.CompressionNone)
//	//printResult(re)
//
//	v, _ := call.Otto.ToValue(nil)
//	return v
//}

////PushTransactions pushes an array of arbitrary JSON transactions
//func (a *chainAPI) PushTransactions(call otto.FunctionCall) (response otto.Value) {
//
//	return getJsResult(call, nil)
//}
//
//func (a *chainAPI) AbiJsonToBin(call otto.FunctionCall) (response otto.Value) {
//
//	return getJsResult(call, nil)
//}
//
//func (a *chainAPI) GetRawAbi(call otto.FunctionCall) (response otto.Value) {
//
//	return getJsResult(call, nil)
//}
//
//func (a *chainAPI) GetRawCodeAndAbi(call otto.FunctionCall) (response otto.Value) {
//
//	return getJsResult(call, nil)
//}
//
//func (a *chainAPI) GetProducers(call otto.FunctionCall) (response otto.Value) {
//
//	return getJsResult(call, nil)
//}

//GetSchedule retrieves the producer schedule
func (a *chainAPI) GetSchedule(call otto.FunctionCall) (response otto.Value) {
	printJSON := false
	str, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	if strings.Contains(str, "j") {
		printJSON = true
	}

	var resp chain_plugin.GetProducerScheduleResult
	err = DoHttpCall(&resp, common.GetScheduleFunc, nil)
	if err != nil {
		clog.Error("GetCurrencyBalance is error: %s", err.Error())
	}
	if printJSON {
		return getJsResult(call, nil)
	}

	print("active", resp.Active)
	print("pending", resp.Pending)
	print("proposed", resp.Proposed)

	return getJsResult(call, nil)
}

func print(name string, schedule common.Variant) {

	producerST, ok := schedule.(types.ProducerScheduleType)
	if !ok {
		fmt.Println("schedule is not producerScheduleType")
	}
	if schedule == nil {
		fmt.Printf("%s schedule empty\n", name)
		return
	}
	fmt.Printf("%s schedule version %d\n", name, producerST.Version)
	fmt.Printf("    %-13s %s\n", "producer", "Producer key")
	fmt.Printf("    %-13s %s\n", "=============", "==================")
	for _, row := range producerST.Producers {
		fmt.Printf("    %-13s %s\n", row.ProducerName.String(), row.BlockSigningKey.String())
	}
	fmt.Printf("\n")

}

func (a *chainAPI) GetTransactionID(call otto.FunctionCall) (response otto.Value) {
	var trx types.Transaction
	JSON, _ := call.Otto.Object("JSON")
	reqVal, err := JSON.Call("stringify", call.Argument(0))
	if err != nil {
		throwJSException(fmt.Sprintf("Fail to parse transaction JSON %s", reqVal.String()))
	}
	err = json.NewDecoder(strings.NewReader(reqVal.String())).Decode(&trx)
	if err != nil {
		throwJSException(fmt.Sprintf("Fail to parse transaction JSON %s", reqVal.String()))
	}
	id := trx.ID()

	return getJsResult(call, id)
}

/*
TODO  convert
pack_transaction            From plain signed json to packed form
unpack_transaction          From packed to plain signed json form
pack_action_data            From json action data to packed form
unpack_action_data          From packed to json action data form
*/

//string plain_signed_transaction_json;
//bool pack_action_data_flag = false;
func (a *chainAPI) ConvertPackTransaction(call otto.FunctionCall) (response otto.Value) {
	plainSignedTransactionJson, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	packActionDataFlag, err := call.Argument(1).ToBoolean()
	if err != nil {
		return otto.UndefinedValue()
	}
	var trx types.SignedTransaction
	var trxVar common.Variants
	err = json.Unmarshal([]byte(plainSignedTransactionJson), &trxVar)
	if err != nil {
		return throwJSException(err.Error())
	}

	packedTrx := &types.PackedTransaction{}
	if packActionDataFlag {
		abi_serializer.FromVariant(&trxVar, &trx, abisSerializerResolver, abiSerializerMaxTime)
		packedTrx = types.NewPackedTransactionBySignedTrx(&trx, types.CompressionNone)
	} else {
		err = json.Unmarshal([]byte(plainSignedTransactionJson), &trx)
		if err != nil {
			return throwJSException(err.Error())
		}
		packedTrx = types.NewPackedTransactionBySignedTrx(&trx, types.CompressionNone)

	}
	return getJsResult(call, packedTrx)

}

var indent string = strings.Repeat(" ", 5)

func PrintAccountResult(res *chain_plugin.GetAccountResult) {
	var staked, unstaking common.Asset
	if res.CoreLiquidBalance.Valid() {
		unstaking = common.Asset{0, res.CoreLiquidBalance.Symbol}
		staked = common.Asset{0, res.CoreLiquidBalance.Symbol}
	}
	fmt.Println("created: ", res.Created.String())
	if res.Privileged {
		fmt.Println("privileged: true")
	}

	fmt.Println("permissions: ")
	tree := make(map[common.Name]Names)
	var roots Names //we don't have multiple roots, but we can easily handle them here, so let's do it just in case
	cache := make(map[common.Name]chain_plugin.Permission)
	for _, perm := range res.Permissions {
		if perm.Parent > 0 {
			tree[perm.Parent] = append(tree[perm.Parent], perm.PermName)
		} else {
			roots = append(roots, perm.PermName)
		}
		cache[perm.PermName] = perm
	}

	dfsPrint := func(name common.AccountName, depth int) {
		p, _ := cache[name]
		fmt.Printf("%s%s%s%s%d%s", indent, strings.Repeat(" ", depth*3), name, " ", p.RequiredAuth.Threshold, ":    ")

		sep := ""
		for _, it := range p.RequiredAuth.Keys {
			fmt.Print(sep, it.Weight, " ", it.Key.String())
			sep = ", "
		}
		for _, acc := range p.RequiredAuth.Accounts {
			fmt.Print(sep, acc.Weight, " ", acc.Permission.Actor.String(), "@", acc.Permission.Permission.String())
			sep = ", "
		}
		fmt.Println()
	}

	sort.Sort(roots)
	for _, r := range roots {
		dfsPrint(r, 0)
		it, ok := tree[r]
		if ok {
			children := it
			sort.Sort(children)
			for i, n := range children {
				dfsPrint(n, 1+i)
			}
		} // else it's a leaf node
	}

	toPrettyNet := func(nbytes int64) string {
		if nbytes == -1 {
			return "unlimited"
		}

		unit := "bytes"
		bytes := float64(nbytes)
		if bytes >= 1024*1024*1024*1024 {
			unit = "TiB"
			bytes /= 1024 * 1024 * 1024 * 1024
		} else if bytes >= 1024*1024*1024 {
			unit = "GiB"
			bytes /= 1024 * 1024 * 1024
		} else if bytes >= 1024*1024 {
			unit = "MiB"
			bytes /= 1024 * 1024
		} else if bytes >= 1024 {
			unit = "KiB"
			bytes /= 1024
		}

		return fmt.Sprintf("%.4g ", bytes) + fmt.Sprintf("%-5s", unit)
	}

	fmt.Println("memory: ")
	fmt.Printf("%s%-15s%s%-15s%s\n\n", indent, "quota: ", toPrettyNet(res.RAMQuota), "  used: ", toPrettyNet(res.RAMUsage))

	fmt.Println("net bandwidth: ")
	if res.TotalResources != nil {
		jsonResources, _ := json.Marshal(res.TotalResources)
		netWeightStr := gjson.GetBytes(jsonResources, "net_weight").String()
		netTotal := common.Asset{}.FromString(&netWeightStr)
		if netTotal.Symbol != unstaking.Symbol {
			// Core symbol of nodeos responding to the request is different than core symbol built into cleos
			unstaking = common.Asset{0, netTotal.Symbol}
			staked = common.Asset{0, netTotal.Symbol}
		}

		if res.SelfDelegatedBandwidth != nil {
			jsonSelfDelegatedBandwidth, _ := json.Marshal(res.SelfDelegatedBandwidth)

			netOwnStr := gjson.GetBytes(jsonSelfDelegatedBandwidth, "net_weight").String()
			netOwn := common.Asset{}.FromString(&netOwnStr)
			staked = netOwn
			netOthers := netTotal.Sub(netOwn)

			fmt.Printf("%s%s%20s%s%s\n", indent, "staked:", netOwn.String(), strings.Repeat(" ", 11), "(total stake delegated from account to self)")
			fmt.Printf("%s%s%17s%s%s\n", indent, "delegated:", netOthers.String(), strings.Repeat(" ", 11), "(total staked delegated to account from others)")
		} else {
			netOthers := netTotal
			fmt.Printf("%s%s%17s%s%s\n", indent, "delegated:", netOthers.String(), strings.Repeat(" ", 11), "(total staked delegated to account from others)")

		}
	}
	fmt.Printf("%s%-11s%18s\n", indent, "used:", toPrettyNet(res.NetLimit.Used))
	fmt.Printf("%s%-11s%18s\n", indent, "available:", toPrettyNet(res.NetLimit.Available))
	fmt.Printf("%s%-11s%18s\n\n", indent, "limit:", toPrettyNet(res.NetLimit.Max))

	fmt.Println("cpu bandwidth:")
	if res.TotalResources != nil {
		jsonResources, _ := json.Marshal(res.TotalResources)
		cpuTotalStr := gjson.GetBytes(jsonResources, "cpu_weight").String()
		cpuTotal := common.Asset{}.FromString(&cpuTotalStr)

		if res.SelfDelegatedBandwidth != nil {
			jsonSelfDelegatedBandwidth, _ := json.Marshal(res.SelfDelegatedBandwidth)

			cpuOwnStr := gjson.GetBytes(jsonSelfDelegatedBandwidth, "cpu_weight").String()
			cpuOwn := common.Asset{}.FromString(&cpuOwnStr)
			staked = staked.Add(cpuOwn)
			cpuOthers := cpuTotal.Sub(cpuOwn)

			fmt.Printf("%s%s%20s%s%s\n", indent, "staked:", cpuOwn.String(), strings.Repeat(" ", 11), "(total stake delegated from account to self)")
			fmt.Printf("%s%s%17s%s%s\n", indent, "delegated:", cpuOthers.String(), strings.Repeat(" ", 11), "(total staked delegated to account from others)")
		} else {
			cpuOthers := cpuTotal
			fmt.Printf("%s%s%17s%s%s\n", indent, "delegated:", cpuOthers.String(), strings.Repeat(" ", 11), "(total staked delegated to account from others)")

		}
	}
	toPrettyTime := func(nmicro int64, widthForUnits uint8) string {
		if nmicro == -1 {
			// special case. Treat it as unlimited
			return "unlimited"
		}

		unit := "us"
		micro := float64(nmicro)
		if micro > 1000000*60*60 {
			micro /= 1000000 * 60 * 60
			unit = "hr"
		} else if micro > 1000000*60 {
			micro /= 1000000 * 60
			unit = "min"
		} else if micro > 1000000 {
			micro /= 1000000
			unit = "sec"
		} else if micro > 1000 {
			micro /= 1000
			unit = "ms"
		}

		if widthForUnits > 0 {
			return fmt.Sprintf("%.4g ", micro) + fmt.Sprintf("%-5s", unit)
		}
		return fmt.Sprintf("%.4g ", micro) + fmt.Sprintf("%s", unit)

	}
	fmt.Printf("%s%-11s%18s\n", indent, "used:", toPrettyTime(res.CpuLimit.Used, 5))
	fmt.Printf("%s%-11s%18s\n", indent, "available:", toPrettyTime(res.CpuLimit.Available, 5))
	fmt.Printf("%s%-11s%18s\n\n", indent, "limit:", toPrettyTime(res.CpuLimit.Max, 5))

	if res.RefundRequest != nil {
		jsonRefundRequest, _ := json.Marshal(res.RefundRequest)
		requestTimeStr := gjson.GetBytes(jsonRefundRequest, "request_time").String()
		requestTime, _ := common.FromIsoStringSec(requestTimeStr)
		refundTime := requestTime.AddSec(uint32(3 * 24 * time.Hour.Seconds())) // +fc::days(3)
		now := res.HeadBlockTime
		netAmountStr := gjson.GetBytes(jsonRefundRequest, "net_amount").String()
		net := common.Asset{}.FromString(&netAmountStr)
		cpuAmountStr := gjson.GetBytes(jsonRefundRequest, "cpu_amount").String()
		cpu := common.Asset{}.FromString(&cpuAmountStr)

		unstaking = net.Add(cpu)
		if unstaking.Amount > 0 {
			fmt.Println("unstaking tokens:")
			fmt.Printf("%s%-25s%20s\n", indent, "time of unstake request:", requestTime.String())

			if now.SubTps(refundTime) > 0 {
				fmt.Println(" (available to claim now with 'eosio::refund' action)")
			} else {
				fmt.Printf(" (funds will be available in %s", toPrettyTime(int64(refundTime.SubUs(now.TimeSinceEpoch())), 0))
			}
			fmt.Printf("%s%-25s%18s\n", indent, "from net bandwidth:", net)
			fmt.Printf("%s%-25s%18s\n", indent, "from cpu bandwidth:", cpu)
			fmt.Printf("%s%-25s%18s\n", indent, "total:", unstaking)
		}
	}

	if res.CoreLiquidBalance.Valid() {
		fmt.Println(res.CoreLiquidBalance.Symbol.Symbol, "balances: ")
		fmt.Printf("%s%-11s%18s\n", indent, "liquid:", res.CoreLiquidBalance.String())
		fmt.Printf("%s%-11s%18s\n", indent, "staked:", staked)
		fmt.Printf("%s%-11s%18s\n", indent, "unstaking:", unstaking)
		fmt.Printf("%s%-11s%18s\n\n", indent, "total:", res.CoreLiquidBalance.Add(staked).Add(unstaking).String())
	}

	if res.VoterInfo != nil {
		jsonVoterInfo, _ := json.Marshal(res.VoterInfo)
		proxyStr := gjson.GetBytes(jsonVoterInfo, "proxy").String()
		if len(proxyStr) == 0 {
			fmt.Print("producers:")
			jsonProds := gjson.GetBytes(jsonVoterInfo, "producers")
			if jsonProds.IsArray() {
				prods := jsonProds.Array()
				if len(prods) != 0 {
					for i, _ := range prods {
						if 1%3 == 0 {
							fmt.Printf("\n%s", indent)
						}
						fmt.Printf("%-16s", prods[i].String())
					}
					//fmt.Println()
				} else {
					fmt.Println(indent, "<not voted>")
				}
			}
		} else {
			fmt.Println("proxy:", indent, proxyStr)
		}
		//fmt.Println()
	}
}
