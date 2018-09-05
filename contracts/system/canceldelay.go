package system

import (
	"github.com/eosspark/eos-go/common"
)

// NewCancelDelay creates an action from the `eosio.system` contract
// called `canceldelay`.
//
// `canceldelay` allows you to cancel a deferred transaction,
// previously sent to the chain with a `delay_sec` larger than 0.  You
// need to sign with cancelingAuth, to cancel a transaction signed
// with that same authority.
func NewCancelDelay(cancelingAuth common.PermissionLevel, transactionID common.SHA256Bytes) *common.Action {
	a := &common.Action{
		Account: common.AccountName(common.StringToName("eosio")),
		Name:    common.ActionName(common.StringToName("canceldelay")),
		Authorization: []common.PermissionLevel{
			cancelingAuth,
		},
		ActionData: common.NewActionData(CancelDelay{
			CancelingAuth: cancelingAuth,
			TransactionID: transactionID,
		}),
	}

	return a
}

// CancelDelay represents the native `canceldelay` action, through the
// system contract.
type CancelDelay struct {
	CancelingAuth common.PermissionLevel `json:"canceling_auth"`
	TransactionID common.SHA256Bytes     `json:"trx_id"`
}
