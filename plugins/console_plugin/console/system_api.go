package console

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/chain/types/generated_containers"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/robertkrimen/otto"
	"sort"
	"strconv"
	"strings"
)

type system struct {
	c *Console
}

func newSystem(c *Console) *system {
	s := &system{
		c: c,
	}
	return s
}

//NewAccount creates a new account on the blockchain with initial resources
func (s *system) NewAccount(call otto.FunctionCall) (response otto.Value) {
	var params NewAccountParams
	readParams(&params, call)

	if len(params.ActiveKey) == 0 {
		params.ActiveKey = params.OwnerKey
	}

	ownerKey, err := ecc.NewPublicKey(params.OwnerKey)
	if err != nil {
		throwJSException(fmt.Sprintf("Invalid owner public key: %s", params.OwnerKey))
	}
	activeKey, err := ecc.NewPublicKey(params.ActiveKey)
	if err != nil {
		throwJSException(fmt.Sprintf("Invalid active public key: %s", params.ActiveKey))
	}

	create := createNewAccount(params.Creator, params.Name, ownerKey, activeKey, params.TxPermission)

	EosAssert(len(params.BuyRamEos) > 0 || params.BuyRamBytesInKbytes > 0 || params.BuyRamBytes > 0, &exception.ExplainedException{},
		"ERROR: One of  buy_ram, buy_ram_kbytes or buy_ram_bytes should have non_zero value")
	EosAssert(params.BuyRamBytesInKbytes == 0 || params.BuyRamBytes == 0, &exception.ExplainedException{},
		"ERROR: buy_ram_kbytes and buy_ram_bytes cannot be set at the same time")

	buyram := &types.Action{}
	if len(params.BuyRamEos) > 0 {
		buyram = createBuyRam(params.Creator, params.Name, toAssetFromString(params.BuyRamEos), params.TxPermission)
	} else {
		var numBytes uint32
		if params.BuyRamBytesInKbytes > 0 {
			numBytes = params.BuyRamBytesInKbytes * 1024
		} else {
			numBytes = params.BuyRamBytes
		}
		buyram = createBuyRamBytes(params.Creator, params.Name, numBytes, params.TxPermission)
	}
	net := toAssetFromString(params.StakeNet)
	cpu := toAssetFromString(params.StakeCpu)

	if net.Amount != 0 || cpu.Amount != 0 {
		delegate := createDelegate(params.Creator, params.Name, net, cpu, params.Transfer, params.TxPermission)
		sendActions([]*types.Action{create, buyram, delegate}, 1000, types.CompressionNone, &params)
	} else {
		sendActions([]*types.Action{create, buyram}, 1000, types.CompressionNone, &params)
	}

	return
}

//RegProducer registers a new producer
func (s *system) RegProducer(call otto.FunctionCall) (response otto.Value) {
	var params RegisterProducer
	readParams(&params, call)

	producerKey, err := ecc.NewPublicKey(params.Key)
	if err != nil {
		Throw(err) // EOS_RETHROW_EXCEPTIONS(public_key_type_exception, "Invalid producer public key: ${public_key}", ("public_key", producer_key_str))
	}

	regprodVar := regProducerVariant(common.N(params.Producer), producerKey, params.Url, params.Loc)

	action := createAction([]common.PermissionLevel{{common.N(params.Producer), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("regproducer"), regprodVar)

	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}

//Unregprod unregisters an existing producer
func (s *system) Unregprod(call otto.FunctionCall) (response otto.Value) {
	var params UnregrodParams
	readParams(&params, call)

	actPayload := common.Variants{"producer": params.Producer}

	action := createAction([]common.PermissionLevel{{common.N(params.Producer), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("unregprod"), &actPayload)
	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}

//VoteproducerProxy votes your stake through a proxy
func (s *system) VoteproducerProxy(call otto.FunctionCall) (response otto.Value) {
	var params Proxy
	readParams(&params, call)

	actPayload := common.Variants{
		"voter":     params.Voter,
		"proxy":     params.Proxy,
		"producers": []common.AccountName{},
	}
	action := createAction([]common.PermissionLevel{{common.N(params.Voter), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("voteproduer"), &actPayload)
	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}

//VoteproducerProds votes for one or more producers
func (s *system) VoteproducerProds(call otto.FunctionCall) (response otto.Value) {
	var params Prods
	readParams(&params, call)

	sort.Sort(params.ProducerNames)
	fmt.Println("producerNames after:", params.ProducerNames)

	actPayload := common.Variants{
		"voter":     params.Voter,
		"proxy":     "",
		"producers": params.ProducerNames,
	}

	action := createAction([]common.PermissionLevel{{common.N(params.Voter), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("voteproducer"), &actPayload)
	sendActions([]*types.Action{action}, 10000, types.CompressionNone, &params)
	return
}

//Approveroducer adds one producer to list of voted producers
func (s *system) Approveroducer(call otto.FunctionCall) (response otto.Value) {
	var params Approve
	readParams(&params, call)

	var res chain_plugin.GetTableRowsResult
	err := DoHttpCall(&res, common.GetTableFunc, common.Variants{
		"json":        true,
		"code":        common.DefaultConfig.SystemAccountName.String(),
		"scope":       common.DefaultConfig.SystemAccountName.String(),
		"table":       "voters",
		"table_key":   "owner",
		"lower_bound": uint64(params.Voter),
		"limit":       1,
	})
	if err != nil {
		throwJSException(fmt.Sprintf("http err :%s", err))
	}

	if len(res.Rows) == 0 {
		throwJSException(fmt.Sprintf("Voter info not found for account %s", params.Voter))
	} else {
		rows0, ok := res.Rows[0]["owner"]
		if !ok {
			throwJSException(fmt.Sprintf("Voter info not found for account %s", params.Voter))
		} else {
			name := rows0.(string)
			if strings.Compare(name, params.Voter.String()) != 0 {
				throwJSException(fmt.Sprintf("Voter info not found for account %s", params.Voter))
			}
		}
	}
	EosAssert(1 == len(res.Rows), &exception.MultipleVoterInfo{}, "More than one voter_info for account")

	prodsInterface, ok := res.Rows[0]["producers"] //TODO
	if !ok {
		throwJSException(fmt.Sprintf("Voter info not found producers"))
	}
	var prods Names
	prodsVars := prodsInterface.(Names)

	for _, name := range prodsVars {
		if uint64(name) != uint64(params.ProducerName) {
			prods = append(prods, name)
		} else {
			throwJSException(fmt.Sprintf("Producer %s is already on the list.", params.ProducerName))
		}
	}
	prods = append(prods, params.ProducerName)

	sort.Sort(prods)
	clog.Debug("prods : %s", prods)

	actPayload := common.Variants{
		"voter":     params.Voter.String(),
		"proxy":     "",
		"producers": prods,
	}

	action := createAction([]common.PermissionLevel{{params.Voter, common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("voteproducer"), &actPayload)
	sendActions([]*types.Action{action}, 10000, types.CompressionNone, &params)
	return
}

//UnapproveProducer removes one producer from list of voted producers
func (s *system) UnapproveProducer(call otto.FunctionCall) (response otto.Value) {
	var params UnapproveProducer
	readParams(&params, call)

	var res chain_plugin.GetTableRowsResult
	err := DoHttpCall(&res, common.GetTableFunc, common.Variants{
		"json":        true,
		"code":        common.DefaultConfig.SystemAccountName.String(),
		"scope":       common.DefaultConfig.SystemAccountName.String(),
		"table":       "voters",
		"table_key":   "owner",
		"lower_bound": uint64(params.Voter),
		"limit":       1,
	})
	if err != nil {
		throwJSException(fmt.Sprintf("http err :%s", err))
	}

	if len(res.Rows) == 0 {
		throwJSException(fmt.Sprintf("Voter info not found for account %s", params.Voter))
	} else {
		rows0, ok := res.Rows[0]["owner"]
		if !ok {
			throwJSException(fmt.Sprintf("Voter info not found for account %s", params.Voter))
		} else {
			name := rows0.(string)
			if strings.Compare(name, params.Voter.String()) != 0 {
				throwJSException(fmt.Sprintf("Voter info not found for account %s", params.Voter))
			}
		}
	}
	EosAssert(1 == len(res.Rows), &exception.MultipleVoterInfo{}, "More than one voter_info for account")

	prodsInterface, ok := res.Rows[0]["producers"] //TODO
	if !ok {
		throwJSException(fmt.Sprintf("Voter info not found producers"))
	}

	var prods Names
	prodsVars := prodsInterface.(Names)

	for _, name := range prodsVars {
		if uint64(name) != uint64(params.ProducerName) {
			prods = append(prods, name)
		}
	}
	if len(prodsVars) == len(prods) {
		throwJSException(fmt.Sprintf("Cannot remove: producer %s is not on the list.", params.ProducerName))
	}

	actPayload := common.Variants{
		"voter":     params.Voter.String(),
		"proxy":     "",
		"producers": prods,
	}
	action := createAction([]common.PermissionLevel{{params.Voter, common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("voteproducer"), &actPayload)

	sendActions([]*types.Action{action}, 10000, types.CompressionNone, &params)
	return
}

//Listproducers lists producers
func (s *system) Listproducers(call otto.FunctionCall) (response otto.Value) {
	var params ListproducersParams
	readParams(&params, call)
	if params.Limit == 0 {
		params.Limit = 50 //default =50
	}
	var resp chain_plugin.GetProducersResult
	err := DoHttpCall(&resp, common.GetProducersFunc, &common.Variants{
		"json":        true,
		"lower_bound": params.Lower,
		"limit":       params.Limit,
	})
	if err != nil {
		throwJSException(fmt.Sprintf("http error : %s", err))
	}

	if params.PrintJson {
		return getJsResult(call, fmt.Sprintln(resp))
	}
	if len(resp.Rows) == 0 {
		return getJsResult(call, fmt.Sprintln("No producers found"))
	}
	weight := resp.TotalProducerVoteWeight
	if weight == 0 {
		weight = 1
	}

	type ShowStruct struct {
		Owner       string  `json:"owner"`
		ProducerKey string  `json:"producer_key"`
		Url         string  `json:"url"`
		TotalVotes  float64 `json:"total_votes"`
	}
	//todo better display information
	fmt.Printf("%-13s %-57s %-59s %s\n", "Producer", "Producer key", "Url", "Scaled votes")
	for _, row := range resp.Rows {
		var show ShowStruct
		bytes, _ := json.Marshal(row)
		err := json.Unmarshal(bytes, &show)
		if err != nil {
			fmt.Printf("resp.rows unmarhal is error:%s", err)
		}

		fmt.Printf("%-13.13s %-57.57s %-59.59s %1.4f\n",
			show.Owner,
			show.ProducerKey,
			show.Url,
			show.TotalVotes/weight)
	}

	if len(resp.More) > 0 {
		fmt.Printf("-L %s for more\n", resp.More)
	}

	return
}

//Delegatebw delegates bandwidth
func (s *system) Delegatebw(call otto.FunctionCall) (response otto.Value) {
	var params DelegatebwParams
	readParams(&params, call)

	EosAssert(len(params.BuyRamAmount) == 0 || params.BuyRamBytes == 0, &exception.ExplainedException{},
		"ERROR: buyram and buy_ram_bytes cannot be set at the same time")

	actPayload := common.Variants{
		"from":               params.From,
		"receiver":           params.Receiver,
		"stake_net_quantity": toAssetFromString(params.StakeNetAmount),
		"stake_cpu_quantity": toAssetFromString(params.StakeCpuAmount),
		"transfer":           params.Transfer,
	}

	action := createAction([]common.PermissionLevel{{common.N(params.From), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("delegatebw"), &actPayload)
	acts := []*types.Action{action}

	if len(params.BuyRamAmount) > 0 {
		acts = append(acts, createBuyRam(common.N(params.From), common.N(params.Receiver), toAssetFromString(params.BuyRamAmount), params.TxPermission))
	} else if params.BuyRamBytes > 0 {
		acts = append(acts, createBuyRamBytes(common.N(params.From), common.N(params.Receiver), params.BuyRamBytes, params.TxPermission))
	}

	sendActions(acts, 1000, types.CompressionNone, &params)
	return
}

//Undelegatebw undelegates bandwidth
func (s *system) Undelegatebw(call otto.FunctionCall) (response otto.Value) {
	var params UndelegatebwParams
	readParams(&params, call)

	actPayload := common.Variants{
		"from":                 params.From,
		"receiver":             params.Receive,
		"unstake_net_quantity": toAssetFromString(params.UnstakeNetAmount),
		"unstake_cpu_quantity": toAssetFromString(params.UnstakeCpuAmount),
	}
	action := createAction([]common.PermissionLevel{{common.N(params.From), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("undelegatebw"), &actPayload)
	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}

//Listbw lists delegated bandwidth
func (s *system) Listbw(call otto.FunctionCall) (response otto.Value) {
	var params ListbwParams
	readParams(&params, call)
	//get entire table in scope of user account
	var resp chain_plugin.GetTableRowsResult
	err := DoHttpCall(&resp, common.GetTableFunc, common.Variants{
		"json":  true,
		"code":  common.DefaultConfig.SystemAccountName.String(),
		"scope": params.Account.String(),
		"table": "delband",
	})
	if err != nil {
		throwJSException(fmt.Sprintf("http is error :%s", err))
	}

	if !params.PrintJson {
		if len(resp.Rows) != 0 { //Todo better display
			//std::cout << std::setw(13) << std::left << "Receiver" << std::setw(21) << std::left << "Net bandwidth"
			//	<< std::setw(21) << std::left << "CPU bandwidth" << std::endl;
			//	for ( auto& r : res.rows ){
			//	std::cout << std::setw(13) << std::left << r["to"].as_string()
			//		<< std::setw(21) << std::left << r["net_weight"].as_string()
			//		<< std::setw(21) << std::left << r["cpu_weight"].as_string()
			//		<< std::endl;
			//	}
		} else {
			return getJsResult(call, fmt.Sprintln("Delegated bandwitdth not found"))
		}
	} else {
		return getJsResult(call, fmt.Sprintln(resp))
	}

	return
}

//Bidname: Name bidding
func (s *system) Bidname(call otto.FunctionCall) (response otto.Value) {
	var params BidnameParams
	readParams(&params, call)

	actPayload := common.Variants{
		"bidder":  params.Bidder,
		"newname": params.NewName,
		"bid":     toAssetFromString(params.BidAmount),
	}
	action := createAction([]common.PermissionLevel{{common.N(params.Bidder), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("bidname"), &actPayload)

	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}

//Bidnameinfo gets bidname info
func (s *system) Bidnameinfo(call otto.FunctionCall) (response otto.Value) {
	var params BidNameinfoParams
	readParams(&params, call)

	var resp chain_plugin.GetTableRowsResult
	err := DoHttpCall(&resp, common.GetTableFunc, common.Variants{
		"json":        true,
		"code":        "eosio",
		"scope":       "eosio",
		"table":       "namebids",
		"lower_bound": common.N(params.Newname),
		"limit":       1,
	})
	if err != nil {
		throwJSException(fmt.Sprintf("http is error: %s\n", err))
	}

	if params.PrintJson {
		return getJsResult(call, fmt.Sprintln(resp))
	}
	if len(resp.Rows) == 0 {
		return getJsResult(call, fmt.Sprintln("No bidname record found"))
	}
	for _, row := range resp.Rows { //todo better display
		fmt.Println(row)
		//fc::time_point time(fc::microseconds(row["last_bid_time"].as_uint64()));
		//	int64_t bid = row["high_bid"].as_int64();
		//std::cout << std::left << std::setw(18) << "bidname:" << std::right << std::setw(24) << row["newname"].as_string() << "\n"
		//	<< std::left << std::setw(18) << "highest bidder:" << std::right << std::setw(24) << row["high_bidder"].as_string() << "\n"
		//	<< std::left << std::setw(18) << "highest bid:" << std::right << std::setw(24) << (bid > 0 ? bid : -bid) << "\n"
		//	<< std::left << std::setw(18) << "last bid time:" << std::right << std::setw(24) << ((std::string)time).c_str() << std::endl;
		//	if (bid < 0) std::cout << "This auction has already closed" << std::endl;
	}

	return
}

//Buyram buy RAM
func (s *system) Buyram(call otto.FunctionCall) (response otto.Value) {
	var params BuyramParams
	readParams(&params, call)
	EosAssert(!params.Kbytes || !params.BytesFlag, &exception.ExplainedException{}, "ERROR: kbytes and bytes cannot be set at the same time")

	action := &types.Action{}
	if params.Kbytes || params.BytesFlag {
		var unit uint64
		if params.Kbytes {
			unit = 1024
		} else {
			unit = 1
		}
		amount, err := strconv.ParseUint(params.Amount, 10, 64)
		if err != nil {
			throwJSException(fmt.Sprintf("parseUint is error: %s\n", err))
		}
		action = createBuyRamBytes(common.N(params.Payer), common.N(params.Receiver), uint32(amount*unit), params.TxPermission)
	} else {
		action = createBuyRam(common.N(params.Payer), common.N(params.Receiver), toAssetFromString(params.Amount), params.TxPermission)
	}

	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}

//Sellram sell RAM
func (s *system) Sellram(call otto.FunctionCall) (response otto.Value) {
	var params SellRamParams
	readParams(&params, call)
	actPayload := common.Variants{
		"account": params.Receiver,
		"bytes":   params.Amount,
	}
	action := createAction([]common.PermissionLevel{{common.N(params.Receiver), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("sellram"), &actPayload)
	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}

//Claimrewards claims producer rewards
func (s *system) Claimrewards(call otto.FunctionCall) (response otto.Value) {
	var params ClaimrewardsParams
	readParams(&params, call)

	actPayload := common.Variants{
		"owner": params.Owner,
	}
	action := createAction([]common.PermissionLevel{{common.N(params.Owner), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("claimrewards"), &actPayload)

	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}

//Regproxy Registers an account as a proxy (for voting)
func (s *system) Regproxy(call otto.FunctionCall) (response otto.Value) {
	var params RegproxyParams
	readParams(&params, call)
	actPayload := common.Variants{
		"proxy":   params.Proxy,
		"isproxy": true,
	}
	action := createAction([]common.PermissionLevel{{common.N(params.Proxy), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("regproxy"), &actPayload)
	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}

//Unregproxy unregisters an account as a proxy (for voting)
func (s *system) Unregproxy(call otto.FunctionCall) (response otto.Value) {
	var params RegproxyParams
	readParams(&params, call)

	actPayload := common.Variants{
		"proxy":   params.Proxy,
		"isproxy": false,
	}
	action := createAction([]common.PermissionLevel{{common.N(params.Proxy), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("regproxy"), &actPayload)
	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}

//Canceldelay cancels a delayed transaction
func (s *system) Canceldelay(call otto.FunctionCall) (response otto.Value) {
	var params CanceldelayParams
	readParams(&params, call)

	cancelingAuth := common.PermissionLevel{common.N(params.CancelingAccount), common.N(params.CanclingPermission)}
	actPayload := common.Variants{
		"canceling_auth": cancelingAuth,
		"trx_id":         params.TrxID,
	}

	action := createAction([]common.PermissionLevel{cancelingAuth}, common.DefaultConfig.SystemAccountName, common.N("canceldelay"), &actPayload)
	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}

func readParams(params interface{}, call otto.FunctionCall) {
	JSON, _ := call.Otto.Object("JSON")
	reqVal, err := JSON.Call("stringify", call.Argument(0))
	if err != nil {
		throwJSException(err.Error())
	}

	rawReq := reqVal.String()
	dec := json.NewDecoder(strings.NewReader(rawReq))
	dec.Decode(&params)
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

func createAction(authorization []common.PermissionLevel, code common.AccountName, act common.ActionName, args *common.Variants) *types.Action {
	return &types.Action{
		Account:       code,
		Name:          act,
		Data:          variantToBin(code, act, args),
		Authorization: authorization,
	}
}

func createBuyRam(creator common.Name, newaccount common.Name, quantity *common.Asset, txPermission []string) *types.Action {
	actPayload := common.Variants{
		"payer":    creator.String(),
		"receiver": newaccount.String(),
		"quant":    quantity.String(),
	}
	var auth []common.PermissionLevel
	if len(txPermission) == 0 {
		auth = []common.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}}
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
	var auth []common.PermissionLevel
	if len(txPermission) == 0 {
		auth = []common.PermissionLevel{{Actor: creator, Permission: common.DefaultConfig.ActiveName}}
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
	var auth []common.PermissionLevel
	if len(txPermission) == 0 {
		auth = []common.PermissionLevel{{Actor: from, Permission: common.DefaultConfig.ActiveName}}
	} else {
		auth = getAccountPermissions(txPermission)
	}
	return createAction(auth, common.DefaultConfig.SystemAccountName, common.N("delegatebw"), &actPayLoad)
}

func sendActions(actions []*types.Action, extraKcpu int32, compression types.CompressionType, c ConsoleInterface) interface{} {
	result := pushActions(actions, extraKcpu, compression, c)

	if c.getOptions().TxPrintJson {
		fmt.Println(result)
	} else {
		printResult(result)
	}

	return nil
}

func pushActions(actions []*types.Action, extraKcpu int32, compression types.CompressionType, c ConsoleInterface) interface{} {
	trx := types.NewSignedTransactionNil()
	trx.Actions = actions
	return pushTransaction(trx, extraKcpu, compression, c)
}

func pushTransaction(trx *types.SignedTransaction, extraKcpu int32, compression types.CompressionType, c ConsoleInterface) interface{} {
	var info chain_plugin.GetInfoResult
	err := DoHttpCall(&info, common.GetInfoFunc, nil)
	if err != nil {
		fmt.Println(err)
	}

	if len(trx.Signatures) == 0 { // #5445 can't change txn content if already signed
		// calculate expiration date
		var expiration common.Microseconds
		if c.getOptions().Expiration == 0 {
			expiration = common.Seconds(30)
		} else {
			expiration = common.Seconds(int64(c.getOptions().Expiration))
		}
		trx.Expiration = common.NewTimePointSecTp(info.HeadBlockTime.AddUs(expiration))

		// Set tapos, default to last irreversible block if it's not specified by the user
		refBlockID := info.LastIrreversibleBlockID
		if len(c.getOptions().TxRefBlockNumOrId) > 0 {
			//var refBlock GetBlockResult
			var refBlock chain_plugin.GetBlockResult
			err := DoHttpCall(&refBlock, common.GetBlockFunc, common.Variants{"block_num_or_id": c.getOptions().TxRefBlockNumOrId})
			if err != nil {
				fmt.Println(err)
				EosThrow(&exception.InvalidRefBlockException{}, "Invalid reference block num or id: %s", c.getOptions().TxRefBlockNumOrId)
			}
			refBlockID = refBlock.ID
		}
		trx.SetReferenceBlock(&refBlockID)

		if c.getOptions().TxForceUnique {
			trx.ContextFreeActions = append(trx.ContextFreeActions, generateNonceAction())
		}
		trx.MaxCpuUsageMS = uint8(c.getOptions().TxMaxCpuUsage)
		trx.MaxNetUsageWords = (common.Vuint32(c.getOptions().TxMaxNetUsage) + 7) / 8
		trx.DelaySec = common.Vuint32(c.getOptions().DelaySec)
	}
	if !c.getOptions().TxSkipSign {
		requiredKeys := determineRequiredKeys(trx)
		signTransaction(trx, requiredKeys, &info.ChainID)
	}
	if !c.getOptions().TxDontBroadcast {
		var result chain_plugin.PushTransactionResult
		packedTrx := types.NewPackedTransactionBySignedTrx(trx, compression)
		err = DoHttpCall(&result, common.PushTxnFunc, packedTrx)
		if err != nil {
			clog.Error(err.Error())
		}
		return result
	} else {
		if !c.getOptions().TxReturnPacked {
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

	publicKeySet := generated.NewPublicKeySet()
	for _, key := range publicKeys {
		pubKey, _ := ecc.NewPublicKey(key)
		publicKeySet.Add(pubKey)
	}
	var keys chain_plugin.GetRequiredKeysResult
	arg := &common.Variants{
		"transaction":    trx,
		"available_keys": publicKeySet,
	}
	err = DoHttpCall(&keys, common.GetRequiredKeys, arg)
	if err != nil {
		clog.Error(err.Error())
	}

	re := make([]string, 0, keys.RequiredKeys.Size())
	for _, key := range keys.RequiredKeys.Values() {
		re = append(re, key.String())
	}

	return re
}

func signTransaction(trx *types.SignedTransaction, requiredKeys []string, chainID *common.ChainIdType) {
	signedTrx := common.Variants{"signed_transaction": trx, "keys": requiredKeys, "id": chainID}
	err := DoHttpCall(trx, common.WalletSignTrx, signedTrx)
	if err != nil {
		clog.Error(err.Error())
	}
}

func generateNonceAction() *types.Action {
	t := common.Now().TimeSinceEpoch()
	data, _ := rlp.EncodeToBytes(t)

	return &types.Action{
		Account:       common.DefaultConfig.NullAccountName,
		Name:          common.N("nonce"),
		Authorization: []common.PermissionLevel{},
		Data:          data,
	}
}

func generateNonceString() string {
	return strconv.FormatInt(common.Now().TimeSinceEpoch().Count(), 10)
}
