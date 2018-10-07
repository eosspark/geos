package system

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

// NewLinkAuth creates an action from the `eosio.system` contract
// called `linkauth`.
//
// `linkauth` allows you to attach certain resouce to the given
// `code::actionName`. With this set on-chain, you can use the
// `requiredPermission` to sign transactions for `code::actionName`
// and not rely on your `active` (which might be more sensitive as it
// can sign anything) for the given operation.
func NewLinkAuth(account, code common.AccountName, actionName common.ActionName, requiredPermission common.PermissionName) *types.Action {
	a := &types.Action{
		Account: common.AccountName(common.N("eosio")),
		Name:    common.ActionName(common.N("linkauth")),
		Authorization: []types.PermissionLevel{
			{account, common.PermissionName(common.N("active"))},
		},
		// Data: common.NewActionData(LinkAuth{//TODO
		// 	Account:     account,
		// 	Code:        code,
		// 	Type:        actionName,
		// 	Requirement: requiredPermission,
		// }),
	}

	return a
}

// LinkAuth represents the native `linkauth` action, through the
// system contract.
type LinkAuth struct {
	Account     common.AccountName    `json:"account"`
	Code        common.AccountName    `json:"code"`
	Type        common.ActionName     `json:"type"`
	Requirement common.PermissionName `json:"requirement"`
}
