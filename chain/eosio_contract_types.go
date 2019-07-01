package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/eos_math"
)

type NewAccount struct {
	Creator common.AccountName `json:"creator"`
	Name    common.AccountName `json:"name"`
	Owner   types.Authority    `json:"owner"`
	Active  types.Authority    `json:"active"`
}

func (n *NewAccount) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (n *NewAccount) GetName() common.ActionName {
	return common.N("newaccount")
}

type SetCode struct {
	Account   common.AccountName `json:"account"`
	VmType    uint8              `json:"vmtype"`
	VmVersion uint8              `json:"vmversion"`
	Code      []byte             `json:"code"`
}

func (s *SetCode) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s *SetCode) GetName() common.ActionName {
	return common.N("setcode")
}

type SetAbi struct {
	Account common.AccountName `json:"account"`
	Abi     []byte             `json:"abi"`
}

func (s SetAbi) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s SetAbi) GetName() common.ActionName {
	return common.N("setabi")
}

type UpdateAuth struct {
	Account    common.AccountName    `json:"account"`
	Permission common.PermissionName `json:"permission"`
	Parent     common.PermissionName `json:"parent"`
	Auth       types.Authority       `json:"auth"`
}

func (u UpdateAuth) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (u UpdateAuth) GetName() common.ActionName {
	return common.N("updateauth")
}

type DeleteAuth struct {
	Account    common.AccountName    `json:""`
	Permission common.PermissionName `json:""`
}

func (d DeleteAuth) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (d DeleteAuth) GetName() common.ActionName {
	return common.N("deleteauth")
}

type LinkAuth struct {
	Account     common.AccountName    `json:"account"`
	Code        common.AccountName    `json:"code"`
	Type        common.ActionName     `json:"type"`
	Requirement common.PermissionName `json:"requirement"`
}

func (l LinkAuth) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (l LinkAuth) GetName() common.ActionName {
	return common.N("linkauth")
}

type UnLinkAuth struct {
	Account common.AccountName `json:"account"`
	Code    common.AccountName `json:"code"`
	Type    common.ActionName  `json:"type"`
}

func (u UnLinkAuth) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (u UnLinkAuth) GetName() common.ActionName {
	return common.N("unlinkauth")
}

type CancelDelay struct {
	CancelingAuth common.PermissionLevel   `json:""`
	TrxId         common.TransactionIdType `json:""`
}

func (c CancelDelay) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (c CancelDelay) GetName() common.ActionName {
	return common.N("canceldelay")
}

type OnError struct {
	SenderId eos_math.Uint128 `json:"sender_id"`
	SentTrx  []byte           `json:"sent_trx"`
}

func NewOnError(sid eos_math.Uint128, data []byte) *OnError {
	oe := OnError{SenderId: sid, SentTrx: data}
	return &oe
}
func (o OnError) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (o OnError) GetName() common.ActionName {
	return common.N("onerror")
}
