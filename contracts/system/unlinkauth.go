package system

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

// NewUnlinkAuth creates an action from the `eosio.system` contract
// called `unlinkauth`.
//
// `unlinkauth` detaches a previously set resouce from a
// `code::actionName`. See `linkauth`.
func NewUnlinkAuth(account, code common.AccountName, actionName common.ActionName) *types.Action {
	a := &types.Action{
		Account: common.AccountName(common.S("eosio")),
		Name:    common.ActionName(common.S("unlinkauth")),
		Authorization: []types.PermissionLevel{
			{account, common.PermissionName(common.S("active"))},
		},
		// Data: common.NewActionData(UnlinkAuth{//TODO
		// 	Account: account,
		// 	Code:    code,
		// 	Type:    actionName,
		// }),
	}

	return a
}

// UnlinkAuth represents the native `unlinkauth` action, through the
// system contract.
type UnlinkAuth struct {
	Account common.AccountName `json:"account"`
	Code    common.AccountName `json:"code"`
	Type    common.ActionName  `json:"type"`
}
