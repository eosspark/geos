package console

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/robertkrimen/otto"
	"strings"
)

type chainAPI struct {
	c   *Console
	log log.Logger
}

func newchainAPI(c *Console) *chainAPI {
	e := &chainAPI{
		c: c,
	}
	e.log = log.New("chainAPI")
	e.log.SetHandler(log.TerminalHandler)
	return e
}

func getJsResult(call otto.FunctionCall, in interface{}) otto.Value {
	bytes, _ := json.Marshal(in)
	resps, _ := call.Otto.Object("new Array()")
	JSON, _ := call.Otto.Object("JSON")
	resultVal, _ := JSON.Call("parse", string(bytes))
	resp, _ := call.Otto.Object(`({"eosgo":"1.0"})`)
	resp.Set("result", resultVal)
	resps.Call("push", resp)

	return resps.Value()
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
		a.log.Error("get account is error: %s", err.Error())
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

	name, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	code, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	abi, err := call.Argument(2).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	wasm, err := call.Argument(3).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	a.log.Debug("%s,%s,%s,%s", name, code, abi, wasm)
	var resp chain_plugin.GetCodeResult
	err = DoHttpCall(&resp, common.GetCodeFunc, common.Variants{"account_name": name, "code_as_wasm": true})
	if err != nil {
		a.log.Error("get abi is error: %s", err.Error())
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
		a.log.Error("get abi is error: %s", err.Error())
	}
	return getJsResult(call, resp)
}

//func (a *chainAPI)GetRawCodeAndAbi(call otto.FunctionCall)(response otto.Value){
//
//}

// get table
// string scope;
// string code;
// string table;
// string lower;
// string upper;
// string table_key;
// string key_type;
// string encode_type{"dec"};
// bool binary = false;
// uint32_t limit = 10;
// string index_position;
// auto getTable = get->add_subcommand( "table", localized("Retrieve the contents of a database table"), false);
// getTable->add_option( "account", code, localized("The account who owns the table") )->required();
// getTable->add_option( "scope", scope, localized("The scope within the contract in which the table is found") )->required();
// getTable->add_option( "table", table, localized("The name of the table as specified by the contract abi") )->required();
// getTable->add_option( "-b,--binary", binary, localized("Return the value as BINARY rather than using abi to interpret as JSON") );
// getTable->add_option( "-l,--limit", limit, localized("The maximum number of rows to return") );
// getTable->add_option( "-k,--key", table_key, localized("Deprecated") );
// getTable->add_option( "-L,--lower", lower, localized("JSON representation of lower bound value of key, defaults to first") );
// getTable->add_option( "-U,--upper", upper, localized("JSON representation of upper bound value of key, defaults to last") );
// getTable->add_option( "--index", index_position,
//                       localized("Index number, 1 - primary (first), 2 - secondary index (in order defined by multi_index), 3 - third index, etc.\n"
//                                 "\t\t\t\tNumber or name of index can be specified, e.g. 'secondary' or '2'."));
// getTable->add_option( "--key-type", key_type,
//                       localized("The key type of --index, primary only supports (i64), all others support (i64, i128, i256, float64, float128, ripemd160, sha256).\n"
//                                 "\t\t\t\tSpecial type 'name' indicates an account name."));
// getTable->add_option( "--encode-type", encode_type,
//                       localized("The encoding type of key_type (i64 , i128 , float64, float128) only support decimal encoding e.g. 'dec'"
//                                  "i256 - supports both 'dec' and 'hex', ripemd160 and sha256 is 'hex' only\n"));

// getTable->set_callback([&] {
//    auto result = call(get_table_func, fc::mutable_variant_object("json", !binary)
//                       ("code",code)
//                       ("scope",scope)
//                       ("table",table)
//                       ("table_key",table_key) // not used
//                       ("lower_bound",lower)
//                       ("upper_bound",upper)
//                       ("limit",limit)
//                       ("key_type",key_type)
//                       ("index_position", index_position)
//                       ("encode_type", encode_type)
//                       );

//    std::cout << fc::json::to_pretty_string(result)
//              << std::endl;
// });
//const resp = await rpc.get_table_rows({
//json: true,              // Get the response as json
//code: 'eosio.token',     // Contract that we target
//scope: 'testacc'         // Account that owns the data
//table: 'accounts'        // Table name
//limit: 10,               // maximum number of rows that we want to get
//});

func (a *chainAPI) GetTable(call otto.FunctionCall) (response otto.Value) {

	call.Argument(0).IsObject()
	code, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	scope, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	table, err := call.Argument(2).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	//tableKey ,err :=call.Argument(3).ToString()
	//if err !=nil{
	//	return otto.UndefinedValue()
	//}
	//lower ,err :=call.Argument(4).ToString()
	//if err !=nil{
	//	return otto.UndefinedValue()
	//}
	//upper ,err :=call.Argument(5).ToString()
	//if err !=nil{
	//	return otto.UndefinedValue()
	//}
	//indexPosition ,err :=call.Argument(6).ToString()
	//if err !=nil{
	//	return otto.UndefinedValue()
	//}
	//keyType ,err :=call.Argument(7).ToString()
	//if err !=nil{
	//	return otto.UndefinedValue()
	//}
	//binary,err :=call.Argument(8).ToBoolean()
	//if err !=nil{
	//	return otto.UndefinedValue()
	//}
	//limit,err := call.Argument(9).ToInteger()
	//if err !=nil{
	//	return otto.UndefinedValue()
	//}

	//if limit ==0{
	//	limit =10
	//}

	var resp chain_plugin.GetTableRowsResult
	err = DoHttpCall(&resp, common.GetTableFunc, common.Variants{"json": true,
		"code":           code,
		"scope":          scope,
		"table":          table,
		"table_key":      "",
		"lower_bound":    "",
		"upper_bound":    "",
		"limit":          10,
		"key_type":       "",
		"index_position": 1})
	if err != nil {
		a.log.Error("get abi is error: %s", err.Error())
	}
	return getJsResult(call, resp)
}

func (a *chainAPI) GetScope(call otto.FunctionCall) (response otto.Value) {
	code, err := call.Argument(0).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	table, err := call.Argument(1).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	lowerBound, err := call.Argument(2).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	upBound, err := call.Argument(3).ToString()
	if err != nil {
		return otto.UndefinedValue()
	}
	limit, err := call.Argument(4).ToInteger()
	if err != nil {
		return otto.UndefinedValue()
	}

	var resp chain_plugin.GetTableByScopeResultRow
	err = DoHttpCall(&resp, common.GetTableByScopeFunc, common.Variants{
		"code":        code,
		"table":       table,
		"lower_bound": lowerBound,
		"upper_bound": upBound,
		"limit":       limit})
	if err != nil {
		a.log.Error("get abi is error: %s", err.Error())
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
		a.log.Error("GetCurrencyBalance is error: %s", err.Error())
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
	var resp chain_plugin.GEtCurrencyStatsResult
	err = DoHttpCall(&resp, common.GetCurrencyStatsFunc, common.Variants{"code": code, "symbol": symbol})
	if err != nil {
		a.log.Error("GetCurrencyBalance is error: %s", err.Error())
	}
	return getJsResult(call, nil)
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
		a.log.Error("GetCurrencyBalance is error: %s", err.Error())
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
	fmt.Printf("%s schedule version %s\n", name, producerST.Version)
	fmt.Printf("    %-13s %s\n", "producer", "Producer key")
	fmt.Printf("    %-13s %s\n", "=============", "==================")
	for _, row := range producerST.Producers {
		fmt.Printf("    %-13s %s\n", row.ProducerName.String(), row.BlockSigningKey.String())
	}
	fmt.Printf("\n")

}
