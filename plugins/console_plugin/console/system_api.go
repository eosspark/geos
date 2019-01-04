package console

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
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

type ConsoleInterface interface {
	getOptions() *StandardTransactionOptions
}

type StandardTransactionOptions struct {
	Expiration        uint64   `json:"expiration"`
	TxForceUnique     bool     `json:"force_unique"`
	TxSkipSign        bool     `json:"skip_sign"`
	TxPrintJson       bool     `json:"json"`
	TxDontBroadcast   bool     `json:"dont_broadcast"`
	TxReturnPacked    bool     `json:"return_packed"`
	TxRefBlockNumOrId string   `json:"ref_block"`
	TxPermission      []string `json:"permission"`
	TxMaxCpuUsage     uint8    `json:"max_cpu_usage_ms"`
	TxMaxNetUsage     uint32   `json:"max_net_usage"`
	DelaySec          uint32   `json:"delay_sec"`
}

func (s *StandardTransactionOptions) getOptions() *StandardTransactionOptions {
	return s
}

type NewAccountParams struct {
	Creator             common.Name `json:"creator"`
	Name                common.Name `json:"name"`
	OwnerKey            string      `json:"owner"`
	ActiveKey           string      `json:"active"`
	StakeNet            string      `json:"stake_net"`
	StakeCpu            string      `json:"stake_cpu"`
	BuyRamBytesInKbytes uint32      `json:"buy_ram_kbytes"`
	BuyRamBytes         uint32      `json:"buy_ram_bytes"`
	BuyRamEos           string      `json:"buy_ram"`
	Transfer            bool        `json:"transfer"`
	StandardTransactionOptions
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

type system struct {
	c *Console
}

func newSystem(c *Console) *system {
	s := &system{
		c: c,
	}
	return s
}

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
	fmt.Println("net and cpu :   ", net, cpu)

	if net.Amount != 0 || cpu.Amount != 0 {
		delegate := createDelegate(params.Creator, params.Name, net, cpu, params.Transfer, params.TxPermission)
		sendActions([]*types.Action{create, buyram, delegate}, 1000, types.CompressionNone, &params)
	} else {
		sendActions([]*types.Action{create, buyram}, 1000, types.CompressionNone, &params)
	}

	return getJsResult(call, nil)
}

type RegisterProducer struct {
	Producer    string `json:"producer"`
	ProducerKey string `json:"producer_key"`
	Url         string `json:"url"`
	Loc         uint16 `json:"loc"`
	StandardTransactionOptions
}

func (s *system) RegProducer(call otto.FunctionCall) (response otto.Value) {
	var params RegisterProducer
	readParams(&params, call)

	producerKey, err := ecc.NewPublicKey(params.ProducerKey)
	if err != nil {
		Throw(err) // EOS_RETHROW_EXCEPTIONS(public_key_type_exception, "Invalid producer public key: ${public_key}", ("public_key", producer_key_str))
	}

	regprodVar := regProducerVariant(common.N(params.Producer), producerKey, params.Url, params.Loc)

	action := createAction([]types.PermissionLevel{{common.N(params.Producer), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("regproducer"), regprodVar)

	re := sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)

	return getJsResult(call, re)
}

type UnregrodParams struct {
	Producer string `json:"producer"`
	StandardTransactionOptions
}

func (s *system) Unregprod(call otto.FunctionCall) (response otto.Value) {
	var params UnregrodParams
	readParams(&params, call)

	actPayload := common.Variants{"producer": params.Producer}

	action := createAction([]types.PermissionLevel{{common.N(params.Producer), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("unregprod"), &actPayload)
	re := sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)

	return getJsResult(call, re)
}

type Proxy struct {
	Voter string `json:"voter"`
	Proxy string `json:"proxy"`
	StandardTransactionOptions
}

func (s *system) VoteproducerProxy(call otto.FunctionCall) (response otto.Value) {
	var params Proxy
	readParams(&params, call)

	actPayload := common.Variants{
		"voter":     params.Voter,
		"proxy":     params.Proxy,
		"producers": []common.AccountName{},
	}
	action := createAction([]types.PermissionLevel{{common.N(params.Voter), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("voteproduer"), &actPayload)
	re := sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)

	return getJsResult(call, re)
}

type Producers []common.Name

func (p Producers) Len() int {
	return len(p)
}
func (p Producers) Less(i, j int) bool {
	return uint64(p[i]) < uint64(p[j])
}
func (p Producers) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type Prods struct {
	Vote          string    `json:"vote"`
	ProducerNames Producers `json:"producer_names"`
	StandardTransactionOptions
}

func (s *system) VoteproducerProds(call otto.FunctionCall) (response otto.Value) {
	var params Prods
	readParams(&params, call)

	sort.Sort(params.ProducerNames)
	fmt.Println("producerNames after:", params.ProducerNames)

	actPayload := common.Variants{
		"voter":     params.Vote,
		"proxy":     "",
		"producers": params.ProducerNames,
	}

	action := createAction([]types.PermissionLevel{{common.N(params.Vote), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("voteproducer"), &actPayload)
	re := sendActions([]*types.Action{action}, 10000, types.CompressionNone, &params)
	return getJsResult(call, re)
}

type Approve struct {
	Voter        common.Name `json:"voter"`
	ProducerName common.Name `json:"producer_name"`
	StandardTransactionOptions
}

func (s *system) VoteproducerApprove(call otto.FunctionCall) (response otto.Value) {
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
	var prods Producers
	prodsVars := prodsInterface.(Producers)

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

	action := createAction([]types.PermissionLevel{{params.Voter, common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("voteproducer"), &actPayload)
	re := sendActions([]*types.Action{action}, 10000, types.CompressionNone, &params)
	return getJsResult(call, re)
}

type UnapproveProducer struct {
	Voter        common.Name `json:"voter"`
	ProducerName common.Name `json:"producer_name"`
	StandardTransactionOptions
}

func (s *system) VoteproducerUnapproveProducer(call otto.FunctionCall) (response otto.Value) {
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

	var prods Producers
	prodsVars := prodsInterface.(Producers)

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
	action := createAction([]types.PermissionLevel{{params.Voter, common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("voteproducer"), &actPayload)

	re := sendActions([]*types.Action{action}, 10000, types.CompressionNone, &params)
	return getJsResult(call, re)
}

type ListproducersParams struct {
	PrintJson bool   `json:"print_json"`
	Limit     uint32 `json:"limit"`
	Lower     string `json:"lower"`
}

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

	return getJsResult(call, nil)
}

type DelegatebwParams struct {
	From               string `json:"from"`
	Receive            string `json:"receive"`
	StakeNetAmount     string `json:"stake_net_amount"`
	StakeCpuAmount     string `json:"stake_cpu_amount"`
	StakeStorageAmount string `json:"stake_storage_amount"`
	BuyRamAmount       string `json:"buy_ram_amount"`
	BuyRamBytes        uint32 `json:"buy_ram_bytes"`
	Transfer           bool   `json:"transfer"`
	StandardTransactionOptions
}

func (s *system) Delegatebw(call otto.FunctionCall) (response otto.Value) {
	var params DelegatebwParams
	readParams(&params, call)

	EosAssert(len(params.BuyRamAmount) == 0 || params.BuyRamBytes == 0, &exception.ExplainedException{},
		"ERROR: buyram and buy_ram_bytes cannot be set at the same time")

	actPayload := common.Variants{
		"from":               params.From,
		"receiver":           params.Receive,
		"stake_net_quantity": toAssetFromString(params.StakeNetAmount),
		"stake_cpu_quantity": toAssetFromString(params.StakeCpuAmount),
		"transfer":           params.Transfer,
	}

	action := createAction([]types.PermissionLevel{{common.N(params.From), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("delegatebw"), &actPayload)
	acts := []*types.Action{action}

	if len(params.BuyRamAmount) > 0 {
		acts = append(acts, createBuyRam(common.N(params.From), common.N(params.Receive), toAssetFromString(params.BuyRamAmount), params.TxPermission))
	} else if params.BuyRamBytes > 0 {
		acts = append(acts, createBuyRamBytes(common.N(params.From), common.N(params.Receive), params.BuyRamBytes, params.TxPermission))
	}

	re := sendActions(acts, 1000, types.CompressionNone, &params)
	return getJsResult(call, re)
}

type UndelegatebwParams struct {
	From             string `json:"from"`
	Receive          string `json:"receive"`
	UnstakeNetAmount string `json:"unstake_net_amount"`
	UnstakeCpuAmount string `json:"unstake_cpu_amount"`
	StandardTransactionOptions
}

func (s *system) Undelegatebw(call otto.FunctionCall) (response otto.Value) {
	var params UndelegatebwParams
	readParams(&params, call)

	actPayload := common.Variants{
		"from":                 params.From,
		"receiver":             params.Receive,
		"unstake_net_quantity": toAssetFromString(params.UnstakeNetAmount),
		"unstake_cpu_quantity": toAssetFromString(params.UnstakeCpuAmount),
	}
	action := createAction([]types.PermissionLevel{{common.N(params.From), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("undelegatebw"), &actPayload)
	re := sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return getJsResult(call, re)
}

type ListbwParams struct {
	Account   common.Name `json:"name"`
	PrintJson bool        `json:"print_json"`
}

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

	return getJsResult(call, nil)
}

type BidnameParams struct {
	Bidder    string `json:"bidder"`
	NewName   string `json:"newname"`
	BidAmount string `json:"bid_amount"`
	StandardTransactionOptions
}

func (s *system) Bidname(call otto.FunctionCall) (response otto.Value) {
	var params BidnameParams
	readParams(&params, call)

	actPayload := common.Variants{
		"bidder":  params.Bidder,
		"newname": params.NewName,
		"bid":     toAssetFromString(params.BidAmount),
	}
	action := createAction([]types.PermissionLevel{{common.N(params.Bidder), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("bidname"), &actPayload)

	re := sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return getJsResult(call, re)
}

type BidNameinfoParams struct {
	PrintJson bool   `json:"print_json"`
	Newname   string `json:"newname"`
}

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

	return getJsResult(call, nil)
}

type BuyramParams struct {
	From      string `json:"from"`
	Receiver  string `json:"receiver"`
	Amount    string `json:"amount"`
	Kbytes    bool   `json:"kbytes"`
	BytesFlag bool   `json:"bytes"`
	StandardTransactionOptions
}

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
		action = createBuyRamBytes(common.N(params.From), common.N(params.Receiver), uint32(amount*unit), params.TxPermission)
	} else {
		action = createBuyRam(common.N(params.From), common.N(params.Receiver), toAssetFromString(params.Amount), params.TxPermission)
	}

	re := sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return getJsResult(call, re)
}

type SellRamParams struct {
	From     string `json:"from"`
	Receiver string `json:"receiver"`
	Amount   uint64 `json:"amount"`
	StandardTransactionOptions
}

func (s *system) Sellram(call otto.FunctionCall) (response otto.Value) {
	var params SellRamParams
	readParams(&params, call)
	actPayload := common.Variants{
		"account": params.Receiver,
		"bytes":   params.Amount,
	}
	action := createAction([]types.PermissionLevel{{common.N(params.Receiver), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("sellram"), &actPayload)
	re := sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return getJsResult(call, re)
}

type ClaimrewardsParams struct {
	Owner string `json:"owner"`
	StandardTransactionOptions
}

func (s *system) Claimrewards(call otto.FunctionCall) (response otto.Value) {
	var params ClaimrewardsParams
	readParams(&params, call)

	actPayload := common.Variants{
		"owner": params.Owner,
	}
	action := createAction([]types.PermissionLevel{{common.N(params.Owner), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("claimrewards"), &actPayload)

	re := sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return getJsResult(call, re)
}

type RegproxyParams struct {
	Proxy string `json:"proxy"`
	StandardTransactionOptions
}

func (s *system) Regproxy(call otto.FunctionCall) (response otto.Value) {
	var params RegproxyParams
	readParams(&params, call)
	actPayload := common.Variants{
		"proxy":   params.Proxy,
		"isproxy": true,
	}
	action := createAction([]types.PermissionLevel{{common.N(params.Proxy), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("regproxy"), &actPayload)
	re := sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return getJsResult(call, re)
}

func (s *system) Unregproxy(call otto.FunctionCall) (response otto.Value) {
	var params RegproxyParams
	readParams(&params, call)

	actPayload := common.Variants{
		"proxy":   params.Proxy,
		"isproxy": false,
	}
	action := createAction([]types.PermissionLevel{{common.N(params.Proxy), common.DefaultConfig.ActiveName}},
		common.DefaultConfig.SystemAccountName, common.N("regproxy"), &actPayload)
	re := sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return getJsResult(call, re)
}

type CanceldelayParams struct {
	CancelingAccount   string `json:"canceling_account"`
	CanclingPermission string `json:"canceling_permission"`
	TrxID              string `json:"trx_id"`
	StandardTransactionOptions
}

func (s *system) Canceldelay(call otto.FunctionCall) (response otto.Value) {
	var params CanceldelayParams
	readParams(&params, call)

	cancelingAuth := types.PermissionLevel{common.N(params.CancelingAccount), common.N(params.CanclingPermission)}
	actPayload := common.Variants{
		"canceling_auth": cancelingAuth,
		"trx_id":         params.TrxID,
	}

	action := createAction([]types.PermissionLevel{cancelingAuth}, common.DefaultConfig.SystemAccountName, common.N("canceldelay"), &actPayload)
	re := sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return getJsResult(call, re)
}

func createAction(authorization []types.PermissionLevel, code common.AccountName, act common.ActionName, args *common.Variants) *types.Action {
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

func sendActions(actions []*types.Action, extraKcpu int32, compression types.CompressionType, c ConsoleInterface) interface{} {
	fmt.Println("send actions...", actions[0].Name)
	result := pushActions(actions, extraKcpu, compression, c)

	if c.getOptions().TxPrintJson {
		return fmt.Sprintln(result)
	}
	return result
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
		trx.MaxNetUsageWords = (uint32(c.getOptions().TxMaxNetUsage) + 7) / 8
		trx.DelaySec = c.getOptions().DelaySec
	}
	if !c.getOptions().TxSkipSign {
		requiredKeys := determineRequiredKeys(trx)
		signTransaction(trx, requiredKeys, &info.ChainID)
	}
	if !c.getOptions().TxDontBroadcast {
		var re common.Variant
		packedTrx := types.NewPackedTransactionBySignedTrx(trx, compression)
		err := DoHttpCall(&re, common.PushTxnFunc, packedTrx)
		if err != nil {
			clog.Error(err.Error())
		}
		return re
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

func generateNonceAction() *types.Action {
	t := common.Now().TimeSinceEpoch()
	data, _ := rlp.EncodeToBytes(t)

	return &types.Action{
		Account:       common.DefaultConfig.NullAccountName,
		Name:          common.N("nonce"),
		Authorization: []types.PermissionLevel{},
		Data:          data,
	}
}
