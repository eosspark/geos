package console

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/robertkrimen/otto"
	"strings"
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

func (a *chainAPI) GetInfo(call otto.FunctionCall) (response otto.Value) {
	var info chain_plugin.GetInfoResult
	err := DoHttpCall(&info, common.GetInfoFunc, nil)
	if err != nil {
		fmt.Println(err)
	}

	return getJsResult(call, info)
}

func (a *chainAPI) GetBlock(call otto.FunctionCall) (response otto.Value) {
	txRefBlockNumOrID, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	var refBlock chain_plugin.GetBlockResult
	err = DoHttpCall(&refBlock, common.GetBlockFunc, common.Variants{"block_num_or_id": txRefBlockNumOrID})
	if err != nil {
		fmt.Println(err)
		try.EosThrow(&exception.InvalidRefBlockException{}, "Invalid reference block num or id: %s", txRefBlockNumOrID)
	}

	return getJsResult(call, refBlock)
}

func (a *chainAPI) GetBlockHeaderState(call otto.FunctionCall) (response otto.Value) {
	txRefBlockNumOrID, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}

	var resp types.BlockHeaderState
	err = DoHttpCall(&resp, common.GetBlockHeaderStateFunc, common.Variants{"block_num_or_id": txRefBlockNumOrID})
	if err != nil {
		fmt.Println(err)
		try.EosThrow(&exception.InvalidRefBlockException{}, "Invalid reference block num or id: %s", txRefBlockNumOrID)
	}

	return getJsResult(call, resp)
}

func (a *chainAPI) GetAccount(call otto.FunctionCall) otto.Value {
	name, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	var resp chain_plugin.GetAccountResult
	err = DoHttpCall(&resp, common.GetAccountFunc, common.Variants{"account_name": name})
	if err != nil {
		clog.Error("get account is error: %s", err.Error())
	}
	return getJsResult(call, resp)
}

// // get code
// string codeFilename;
// string abiFilename;
// bool code_as_wasm = false;
// auto getCode = get->add_subcommand("code", localized("Retrieve the code and ABI for an account"), false);
// getCode->add_option("name", accountName, localized("The name of the account whose code should be retrieved"))->required();
// getCode->add_option("-c,--code",codeFilename, localized("The name of the file to save the contract .wast/wasm to") );
// getCode->add_option("-a,--abi",abiFilename, localized("The name of the file to save the contract .abi to") );
// getCode->add_flag("--wasm", code_as_wasm, localized("Save contract as wasm"));
// getCode->set_callback([&] {
//    string code_hash, wasm, wast, abi;
//    try {
//       const auto result = call(get_raw_code_and_abi_func, fc::mutable_variant_object("account_name", accountName));
//       const std::vector<char> wasm_v = result["wasm"].as_blob().data;
//       const std::vector<char> abi_v = result["abi"].as_blob().data;

//       fc::sha256 hash;
//       if(wasm_v.size())
//          hash = fc::sha256::hash(wasm_v.data(), wasm_v.size());
//       code_hash = (string)hash;

//       wasm = string(wasm_v.begin(), wasm_v.end());
//       if(!code_as_wasm && wasm_v.size())
//          wast = wasm_to_wast((const uint8_t*)wasm_v.data(), wasm_v.size(), false);

//       abi_def abi_d;
//       if(abi_serializer::to_abi(abi_v, abi_d))
//          abi = fc::json::to_pretty_string(abi_d);
//    }
//    catch(chain::missing_chain_api_plugin_exception&) {
//       //see if this is an old nodeos that doesn't support get_raw_code_and_abi
//       const auto old_result = call(get_code_func, fc::mutable_variant_object("account_name", accountName)("code_as_wasm",code_as_wasm));
//       code_hash = old_result["code_hash"].as_string();
//       if(code_as_wasm) {
//          wasm = old_result["wasm"].as_string();
//          std::cout << localized("Warning: communicating to older nodeos which returns malformed binary wasm") << std::endl;
//       }
//       else
//          wast = old_result["wast"].as_string();
//       abi = fc::json::to_pretty_string(old_result["abi"]);
//    }

//    std::cout << localized("code hash: ${code_hash}", ("code_hash", code_hash)) << std::endl;

//    if( codeFilename.size() ){
//       std::cout << localized("saving ${type} to ${codeFilename}", ("type", (code_as_wasm ? "wasm" : "wast"))("codeFilename", codeFilename)) << std::endl;

//       std::ofstream out( codeFilename.c_str() );
//       if(code_as_wasm)
//          out << wasm;
//       else
//          out << wast;
//    }
//    if( abiFilename.size() ) {
//       std::cout << localized("saving abi to ${abiFilename}", ("abiFilename", abiFilename)) << std::endl;
//       std::ofstream abiout( abiFilename.c_str() );
//       abiout << abi;
//    }
// });

func (a *chainAPI) GetCode(call otto.FunctionCall) otto.Value { //TODO save to file
	var params GetCodeParams
	readParams(&params, call)

	var resp chain_plugin.GetCodeResult
	err := DoHttpCall(&resp, common.GetCodeFunc, common.Variants{"account_name": params.AccountName, "code_as_wasm": params.CodeAsWasm})
	if err != nil {
		clog.Error("get abi is error: %s", err.Error())
	}
	return getJsResult(call, resp)
}

// // get abi
// string filename;
// auto getAbi = get->add_subcommand("abi", localized("Retrieve the ABI for an account"), false);
// getAbi->add_option("name", accountName, localized("The name of the account whose abi should be retrieved"))->required();
// getAbi->add_option("-f,--file",filename, localized("The name of the file to save the contract .abi to instead of writing to console") );
// getAbi->set_callback([&] {
//    auto result = call(get_abi_func, fc::mutable_variant_object("account_name", accountName));
//    auto abi  = fc::json::to_pretty_string( result["abi"] );
//    if( filename.size() ) {
//       std::cerr << localized("saving abi to ${filename}", ("filename", filename)) << std::endl;
//       std::ofstream abiout( filename.c_str() );
//       abiout << abi;
//    } else {
//       std::cout << abi << "\n";
//    }
// });

func (a *chainAPI) GetAbi(call otto.FunctionCall) otto.Value { //TODO save to file
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

//func (a *chainAPI)GetRawCodeAndAbi(call otto.FunctionCall)(response otto.Value){
//
//}

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

func (a *chainAPI) PushTransaction(call otto.FunctionCall) (response otto.Value) {

	return getJsResult(call, nil)
}

func (a *chainAPI) PushTransactions(call otto.FunctionCall) (response otto.Value) {

	return getJsResult(call, nil)
}

func (a *chainAPI) AbiJsonToBin(call otto.FunctionCall) (response otto.Value) {

	return getJsResult(call, nil)
}

func (a *chainAPI) GetRawAbi(call otto.FunctionCall) (response otto.Value) {

	return getJsResult(call, nil)
}

func (a *chainAPI) GetRawCodeAndAbi(call otto.FunctionCall) (response otto.Value) {

	return getJsResult(call, nil)
}

func (a *chainAPI) GetProducers(call otto.FunctionCall) (response otto.Value) {

	return getJsResult(call, nil)
}

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
