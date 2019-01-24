package console

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/plugins/chain_plugin"
	"github.com/robertkrimen/otto"
	"github.com/tidwall/gjson"
)

type multisig struct {
	c *Console
}

//Multisig contract commands
func newMultisig(c *Console) *multisig {
	m := &multisig{
		c: c,
	}
	return m
}

type ProposeParams struct {
	ProposalName            string `json:"proposal_name"`
	RequestedPerm           string `json:"requested_permissions"`
	TransactionPerm         string `json:"trx_permissions"`
	ProposedContract        string `json:"contract"`
	ProposedAction          string `json:"action"`
	ProposedTransaction     string `json:"data"`
	Proposer                string `json:"proposer"`
	ProposalExpirationHours int    `json:"proposal_expiration"` //TODO default 24
	StandardTransactionOptions
}

//Propose proposes action
func (m *multisig) Propose(call otto.FunctionCall) (response otto.Value) {
	var params ProposeParams
	readParams(&params, call)

	reqperm := make([]common.PermissionLevel, 0)
	err := json.Unmarshal([]byte(params.RequestedPerm), &reqperm)
	if err != nil {
		fmt.Println("Unmarshal requestedPerm is error: ", err)
		return otto.FalseValue()
	}

	trxperm := make([]common.PermissionLevel, 0)
	err = json.Unmarshal([]byte(params.TransactionPerm), &trxperm)
	if err != nil {
		fmt.Println("Unmarshal TransactionPerm is error: ", err)
		return otto.FalseValue()
	}

	var proposedTrx types.Transaction
	err = json.Unmarshal([]byte(params.ProposedTransaction), &proposedTrx)
	if err != nil {
		fmt.Println("Unmarshal ProposedTransaction to transaction is error: ", err)
		return otto.FalseValue()
	}
	trxVar := common.Variants{}
	err = json.Unmarshal([]byte(params.ProposedTransaction), &trxVar)
	if err != nil {
		fmt.Println("Unmarshal ProposedTransaction to Variants is error: ", err)
		return otto.FalseValue()
	}
	proposedTrxSerialized := variantToBin(common.N(params.ProposedContract), common.N(params.ProposedAction), &trxVar)

	accountPermissions := getAccountPermissions(params.TxPermission)
	if len(accountPermissions) == 0 {
		if len(params.Proposer) > 0 {
			accountPermissions = []common.PermissionLevel{{common.N(params.Proposer), common.DefaultConfig.ActiveName}}
		} else {
			try.EosThrow(&exception.MissingAuthException{}, "Authority is not provided (either by multisig parameter <proposer> or -p")
		}
	}
	if len(params.Proposer) == 0 {
		params.Proposer = accountPermissions[0].Actor.String()
	}

	if params.ProposalExpirationHours == 0 {
		params.ProposalExpirationHours = 24 //default =24
	}
	trx := types.Transaction{}
	trx.Expiration = common.NewTimePointSecTp(common.Now().AddUs(common.Microseconds(params.ProposalExpirationHours * 60 * 60 * 1000 * 1000)))
	trx.RefBlockNum = 0
	trx.RefBlockPrefix = 0
	trx.MaxNetUsageWords = 0
	trx.MaxCpuUsageMS = 0
	trx.DelaySec = 0

	action := &types.Action{
		Account:       common.N(params.ProposedContract),
		Name:          common.N(params.ProposedAction),
		Authorization: trxperm,
		Data:          proposedTrxSerialized,
	}
	trx.Actions = append(trx.Actions, action)

	//fc::to_variant(trx, trx_var);//TODO
	bytes, err := json.Marshal(trx)
	if err != nil {
		fmt.Println("Marshal trx is error: ", err)
		return otto.FalseValue()
	}

	args := common.Variants{
		"proposer":      params.Proposer,
		"proposal_name": params.ProposalName,
		"requested":     params.RequestedPerm,
		"trx":           string(bytes), //trxVar,
	}

	actionPropose := &types.Action{
		Account:       common.N("eosio.msig"),
		Name:          common.N("propose"),
		Authorization: accountPermissions,
		Data:          variantToBin(common.N("eosio.msig"), common.N("propose"), &args),
	}
	sendActions([]*types.Action{actionPropose}, 1000, types.CompressionNone, &params)

	return otto.UndefinedValue()
}

type ProposeTrxParams struct {
	ProposalName  string `json:"proposal_name"`
	RequestedPerm string `json:"requested_permissions"`
	Proposer      string `json:"proposer"`
	TrxToPush     string `json:"transaction"`
	StandardTransactionOptions
}

//ProposeTrx proposes transaction
func (m *multisig) ProposeTrx(call otto.FunctionCall) (response otto.Value) {
	var params ProposeTrxParams
	readParams(&params, call)

	reqperm := make([]common.PermissionLevel, 0)
	err := json.Unmarshal([]byte(params.RequestedPerm), &reqperm)
	if err != nil {
		fmt.Println("Unmarshal requestedPerm is error: ", err)
		return otto.FalseValue()
	}

	trx := types.Transaction{}
	err = json.Unmarshal([]byte(params.RequestedPerm), &trx)
	if err != nil {
		fmt.Println("Unmarshal requestedPerm is error: ", err)
		return otto.FalseValue()
	}
	accountPermissions := getAccountPermissions(params.TxPermission)
	if len(accountPermissions) == 0 {
		if len(params.Proposer) > 0 {
			accountPermissions = []common.PermissionLevel{{common.N(params.Proposer), common.DefaultConfig.ActiveName}}
		} else {
			try.EosThrow(&exception.MissingAuthException{}, "Authority is not provided (either by multisig parameter <proposer> or -p")
		}
	}
	if len(params.Proposer) == 0 {
		params.Proposer = accountPermissions[0].Actor.String()
	}

	args := common.Variants{
		"proposer":      params.Proposer,
		"proposa;_name": params.ProposalName,
		"requested":     params.RequestedPerm,
		"trx":           params.TrxToPush,
	}

	actionPropose := &types.Action{
		Account:       common.N("eosio.msig"),
		Name:          common.N("propose"),
		Authorization: accountPermissions,
		Data:          variantToBin(common.N("eosio.msig"), common.N("propose"), &args),
	}
	sendActions([]*types.Action{actionPropose}, 1000, types.CompressionNone, &params)

	return otto.UndefinedValue()
}

type ReviewParams struct {
	ProposalName string `json:"proposal_name"`
	Proposer     string `json:"proposer"`
}

//Review reviews transaction
func (m *multisig) Review(call otto.FunctionCall) (response otto.Value) {
	var params ReviewParams
	readParams(&params, call)

	var resp chain_plugin.GetTableRowsResult
	err := DoHttpCall(&resp, common.GetTableFunc, common.Variants{
		"json":        true,
		"code":        "eosio.msig",
		"scope":       params.Proposer,
		"table":       "proposal",
		"table_key":   "",
		"lower_bound": common.N(params.ProposalName),
		"upper_bound": "",
		"limit":       1,
	})

	result, _ := json.Marshal(resp)
	rows := gjson.GetBytes(result, "rows").Array()
	if len(rows) == 0 {
		fmt.Println("Proposal not found")
		return
	}
	obj := rows[0]
	if obj.Get("proposal_name").String() != params.ProposalName {
		fmt.Println("Proposal not found")
		return
	}
	trxHex := obj.Get("packed_transaction").String()
	trx := types.Transaction{}
	err = rlp.DecodeBytes([]byte(trxHex), &trx)
	if err != nil {
		fmt.Println("decode packed_transaction is error:", err)
		return
	}

	common.Set(obj.String(), "transaction", trx) //TODO
	fmt.Println(obj.String())

	return otto.UndefinedValue()
}

type ApproveAndUnapproveParams struct {
	Proposer     string `json:"proposer"`
	ProposalName string `json:"proposal_name"`
	Perm         string `json:"permissions"`
	StandardTransactionOptions
}

//Approve approves proposed transaction
func (m *multisig) Approve(call otto.FunctionCall) (response otto.Value) {
	var params ApproveAndUnapproveParams
	readParams(&params, call)

	approveOrunapprove("approve", &params)
	return otto.UndefinedValue()
}

//Unapprove unapproves proposed transaction
func (m *multisig) Unapprove(call otto.FunctionCall) (response otto.Value) {
	var params ApproveAndUnapproveParams
	readParams(&params, call)

	approveOrunapprove("unapprove", &params)
	return otto.UndefinedValue()
}

func approveOrunapprove(action string, p *ApproveAndUnapproveParams) {
	var permissions []common.PermissionLevel
	err := json.Unmarshal([]byte(p.Perm), &permissions)
	if err != nil {
		fmt.Println("Fail to parse permissions JSON ", p.Perm, err)
		return
	}
	args := common.Variants{
		"proposer":      p.Proposer,
		"proposal_name": p.ProposalName,
		"level":         p.Perm,
	}
	var accountPermissions []common.PermissionLevel
	if len(p.TxPermission) == 0 {
		accountPermissions = []common.PermissionLevel{{common.N(p.Proposer), common.DefaultConfig.ActiveName}}
	} else {
		accountPermissions = getAccountPermissions(p.TxPermission)
	}
	a := &types.Action{
		Account:       common.N("eosio.msig"),
		Name:          common.N(action),
		Authorization: accountPermissions,
		Data:          variantToBin(common.N("eosio.msig"), common.N(action), &args),
	}

	sendActions([]*types.Action{a}, 1000, types.CompressionNone, p)
}

type CancelParams struct {
	Proposer     string `json:"proposer"`
	ProposalName string `json:"proposal_name"`
	Canceler     string `json:"canceler"`
	StandardTransactionOptions
}

//Cancel cancels proposed transaction
func (m *multisig) Cancel(call otto.FunctionCall) (response otto.Value) {
	var params CancelParams
	readParams(&params, call)

	accountPermissions := getAccountPermissions(params.TxPermission)
	if len(accountPermissions) == 0 {
		if len(params.Proposer) > 0 {
			accountPermissions = []common.PermissionLevel{{common.N(params.Canceler), common.DefaultConfig.ActiveName}}
		} else {
			try.EosThrow(&exception.MissingAuthException{}, "Authority is not provided (either by multisig parameter <proposer> or -p")
		}
	}
	if len(params.Canceler) == 0 {
		params.Canceler = accountPermissions[0].Actor.String()
	}

	args := common.Variants{
		"proposer":      params.Proposer,
		"proposal_name": params.ProposalName,
		"canceler":      params.Canceler,
	}

	action := &types.Action{
		Account:       common.N("eosio.msig"),
		Name:          common.N("cancel"),
		Authorization: accountPermissions,
		Data:          variantToBin(common.N("eosio.msig"), common.N("cancel"), &args),
	}
	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)

	return otto.UndefinedValue()
}

type ExecuteParams struct {
	Proposer     string `json:"proposer"`
	ProposalName string `json:"proposal_name"`
	Executer     string `json:"executer"`
	StandardTransactionOptions
}

//Exec execute proposed transaction
func (m *multisig) Exec(call otto.FunctionCall) (response otto.Value) {
	var params ExecuteParams
	readParams(&params, call)

	accountPermissions := getAccountPermissions(params.TxPermission)
	if len(accountPermissions) == 0 {
		if len(params.Proposer) > 0 {
			accountPermissions = []common.PermissionLevel{{common.N(params.Executer), common.DefaultConfig.ActiveName}}
		} else {
			try.EosThrow(&exception.MissingAuthException{}, "Authority is not provided (either by multisig parameter <proposer> or -p")
		}
	}
	if len(params.Executer) == 0 {
		params.Executer = accountPermissions[0].Actor.String()
	}

	args := common.Variants{
		"proposer":      params.Proposer,
		"proposal_name": params.ProposalName,
		"executer":      params.Executer,
	}

	action := &types.Action{
		Account:       common.N("eosio.msig"),
		Name:          common.N("exec"),
		Authorization: accountPermissions,
		Data:          variantToBin(common.N("eosio.msig"), common.N("exec"), &args),
	}
	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)

	return otto.UndefinedValue()
}
