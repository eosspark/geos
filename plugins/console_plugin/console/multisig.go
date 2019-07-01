package console

import (
	"encoding/hex"
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

type multiSig struct {
	c *Console
}

//MultiSig contract commands
func newMultiSig(c *Console) *multiSig {
	return &multiSig{c: c}
}

//Propose proposes action
func (m *multiSig) Propose(call otto.FunctionCall) (response otto.Value) {
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
			accountPermissions = []common.PermissionLevel{{Actor: common.N(params.Proposer), Permission: common.DefaultConfig.ActiveName}}
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

	actionPropose := &types.Action{
		Account:       common.N("eosio.msig"),
		Name:          common.N("propose"),
		Authorization: accountPermissions,
		Data: variantToBin(
			common.N("eosio.msig"),
			common.N("propose"),
			&common.Variants{
				"proposer":      params.Proposer,
				"proposal_name": params.ProposalName,
				"requested":     reqperm,
				"trx":           trx,
			}),
	}
	sendActions([]*types.Action{actionPropose}, 1000, types.CompressionNone, &params)
	return
}

//ProposeTrx proposes transaction
func (m *multiSig) ProposeTrx(call otto.FunctionCall) (response otto.Value) {
	var params ProposeTrxParams
	readParams(&params, call)

	reqperm := make([]common.PermissionLevel, 0)
	err := json.Unmarshal([]byte(params.RequestedPerm), &reqperm)
	if err != nil {
		fmt.Println("Unmarshal requestedPerm is error: ", err)
		return otto.FalseValue()
	}

	trx := types.Transaction{}
	err = json.Unmarshal([]byte(params.TrxToPush), &trx)
	if err != nil {
		fmt.Println("Unmarshal requestedPerm is error: ", err)
		return otto.FalseValue()
	}
	accountPermissions := getAccountPermissions(params.TxPermission)
	if len(accountPermissions) == 0 {
		if len(params.Proposer) > 0 {
			accountPermissions = []common.PermissionLevel{{Actor: common.N(params.Proposer), Permission: common.DefaultConfig.ActiveName}}
		} else {
			try.EosThrow(&exception.MissingAuthException{}, "Authority is not provided (either by multisig parameter <proposer> or -p")
		}
	}
	if len(params.Proposer) == 0 {
		params.Proposer = accountPermissions[0].Actor.String()
	}

	actionPropose := &types.Action{
		Account:       common.N("eosio.msig"),
		Name:          common.N("propose"),
		Authorization: accountPermissions,
		Data: variantToBin(common.N("eosio.msig"), common.N("propose"),
			&common.Variants{
				"proposer":      params.Proposer,
				"proposal_name": params.ProposalName,
				"requested":     reqperm,
				"trx":           trx,
			}),
	}
	sendActions([]*types.Action{actionPropose}, 1000, types.CompressionNone, &params)
	return
}

//Review reviews transaction
func (m *multiSig) Review(call otto.FunctionCall) (response otto.Value) {
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
	if err != nil {
		return otto.FalseValue()
	}

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

	data, _ := hex.DecodeString(obj.Get("packed_transaction").String())
	trx := types.Transaction{}
	rlp.DecodeBytes(data, &trx)
	re, _ := common.Set(obj.String(), "transaction", trx)
	fmt.Println("{")
	fmt.Println("\tproposal_name: ", gjson.GetBytes([]byte(re), "proposal_name"))
	fmt.Println("\tpacked_transaction:", gjson.GetBytes([]byte(re), "packed_transaction"))
	fmt.Println("\ttransaction:", gjson.GetBytes([]byte(re), "transaction"))
	fmt.Println("}")
	return
}

//Approve approves proposed transaction
func (m *multiSig) Approve(call otto.FunctionCall) (response otto.Value) {
	var params ApproveAndUnapproveParams
	readParams(&params, call)

	m.approveOrUnapprove("approve", &params)
	return
}

//Unapprove unapproves proposed transaction
func (m *multiSig) Unapprove(call otto.FunctionCall) (response otto.Value) {
	var params ApproveAndUnapproveParams
	readParams(&params, call)

	m.approveOrUnapprove("unapprove", &params)
	return
}

func (m *multiSig) approveOrUnapprove(action string, p *ApproveAndUnapproveParams) {
	var permissions common.PermissionLevel
	err := json.Unmarshal([]byte(p.Perm), &permissions)
	if err != nil {
		fmt.Println("Fail to parse permissions JSON ", p.Perm, err)
		return
	}

	var accountPermissions []common.PermissionLevel
	if len(p.TxPermission) == 0 {
		accountPermissions = []common.PermissionLevel{{Actor: common.N(p.Proposer), Permission: common.DefaultConfig.ActiveName}}
	} else {
		accountPermissions = getAccountPermissions(p.TxPermission)
	}
	a := &types.Action{
		Account:       common.N("eosio.msig"),
		Name:          common.N(action),
		Authorization: accountPermissions,
		Data: variantToBin(common.N("eosio.msig"), common.N(action), &common.Variants{
			"proposer":      p.Proposer,
			"proposal_name": p.ProposalName,
			"level":         permissions,
		}),
	}
	sendActions([]*types.Action{a}, 1000, types.CompressionNone, p)
}

//Cancel cancels proposed transaction
func (m *multiSig) Cancel(call otto.FunctionCall) (response otto.Value) {
	var params CancelParams
	readParams(&params, call)

	accountPermissions := getAccountPermissions(params.TxPermission)
	if len(accountPermissions) == 0 {
		if len(params.Proposer) > 0 {
			accountPermissions = []common.PermissionLevel{{Actor: common.N(params.Canceler), Permission: common.DefaultConfig.ActiveName}}
		} else {
			try.EosThrow(&exception.MissingAuthException{}, "Authority is not provided (either by multisig parameter <proposer> or -p")
		}
	}
	if len(params.Canceler) == 0 {
		params.Canceler = accountPermissions[0].Actor.String()
	}

	action := &types.Action{
		Account:       common.N("eosio.msig"),
		Name:          common.N("cancel"),
		Authorization: accountPermissions,
		Data: variantToBin(common.N("eosio.msig"), common.N("cancel"), &common.Variants{
			"proposer":      params.Proposer,
			"proposal_name": params.ProposalName,
			"canceler":      params.Canceler,
		}),
	}
	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}

//Exec execute proposed transaction
func (m *multiSig) Exec(call otto.FunctionCall) (response otto.Value) {
	var params ExecuteParams
	readParams(&params, call)

	accountPermissions := getAccountPermissions(params.TxPermission)
	if len(accountPermissions) == 0 {
		if len(params.Proposer) > 0 {
			accountPermissions = []common.PermissionLevel{{Actor: common.N(params.Executer), Permission: common.DefaultConfig.ActiveName}}
		} else {
			try.EosThrow(&exception.MissingAuthException{}, "Authority is not provided (either by multisig parameter <proposer> or -p")
		}
	}
	if len(params.Executer) == 0 {
		params.Executer = accountPermissions[0].Actor.String()
	}

	action := &types.Action{
		Account:       common.N("eosio.msig"),
		Name:          common.N("exec"),
		Authorization: accountPermissions,
		Data: variantToBin(common.N("eosio.msig"), common.N("exec"), &common.Variants{
			"proposer":      params.Proposer,
			"proposal_name": params.ProposalName,
			"executer":      params.Executer,
		}),
	}
	sendActions([]*types.Action{action}, 1000, types.CompressionNone, &params)
	return
}
