package system

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

// NewDeleteAuth creates an action from the `eosio.system` contract
// called `deleteauth`.
//
// You cannot delete the `owner` or `active` permissions.  Also, if a
// resouce is still linked through a previous `updatelink` action,
// you will need to `unlinkauth` first.
func NewDeleteAuth(account common.AccountName, permission common.PermissionName) *types.Action {
	a := &types.Action{
		Account: common.AccountName(common.N("eosio")),
		Name:    common.ActionName(common.N("deleteauth")),
		Authorization: []types.PermissionLevel{
			{Actor: account, Permission: common.PermissionName(common.N("active"))},
		},
		// Data: common.NewActionData(DeleteAuth{//TODO
		// 	Account:    account,
		// 	Permission: permission,
		// }),
	}

	return a
}

// DeleteAuth represents the native `deleteauth` action, reachable
// through the `eosio.system` contract.
type DeleteAuth struct {
	Account    common.AccountName    `json:"account"`
	Permission common.PermissionName `json:"resouce"`
}
