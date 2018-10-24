package types

import "github.com/eosspark/eos-go/common"

// NewAccount represents the newaccount
type NewAccount struct {
	Creator common.AccountName `json:"creator"`
	Name    common.AccountName `json:"name"`
	Owner   Authority          `json:"owner"`
	Active  Authority          `json:"active"`
}

func (n *NewAccount) GetAccount() common.AccountName {
	return common.AccountName(common.DefaultConfig.SystemAccountName)
}

func (n *NewAccount) GetName() common.ActionName {
	name := common.N("newaccount")
	return common.ActionName(name)
}

// SetCode represents the hard-coded `setcode` action.
type SetCode struct {
	Account   common.AccountName `json:"account"`
	VMType    byte               `json:"vmtype"`
	VMVersion byte               `json:"vmversion"`
	Code      common.HexBytes    `json:"bytes"`
}

func (n *SetCode) GetAccount() common.AccountName {
	return common.AccountName(common.DefaultConfig.SystemAccountName)
}

func (n *SetCode) GetName() common.ActionName {
	name := common.N("setcode")
	return common.ActionName(name)
}

// SetABI represents the hard-coded `setabi` action.
type SetABI struct {
	Account common.AccountName `json:"account"`
	ABI     AbiDef             `json:"abi"`
}

func (n *SetABI) GetAccount() common.AccountName {
	return common.AccountName(common.DefaultConfig.SystemAccountName)
}

func (n *SetABI) GetName() common.ActionName {
	name := common.N("setabi")
	return common.ActionName(name)
}

type UpdateAuth struct {
	Account    common.AccountName
	Permission common.PermissionName
	Parent     common.PermissionName
	Auth       Authority
}

func (n *UpdateAuth) GetAccount() common.AccountName {
	return common.AccountName(common.DefaultConfig.SystemAccountName)
}

func (n *UpdateAuth) GetName() common.ActionName {
	name := common.N("updateauth")
	return common.ActionName(name)
}

type DeleteAuth struct {
	Account    common.AccountName
	Permission common.PermissionName
}

func (n *DeleteAuth) GetAccount() common.AccountName {
	return common.AccountName(common.DefaultConfig.SystemAccountName)
}

func (n *DeleteAuth) GetName() common.ActionName {
	name := common.N("deleteauth")
	return common.ActionName(name)
}

type LinkAuth struct {
	Account     common.AccountName
	Code        common.AccountName
	Type        common.ActionName
	Requirement common.PermissionName
}

func (n *LinkAuth) GetAccount() common.AccountName {
	return common.AccountName(common.DefaultConfig.SystemAccountName)
}

func (n *LinkAuth) GetName() common.ActionName {
	name := common.N("linkauth")
	return common.ActionName(name)
}

type UnlinkAuth struct {
	Account     common.AccountName
	Code        common.AccountName
	Type        common.ActionName
}

func (n *UnlinkAuth) GetAccount() common.AccountName {
	return common.AccountName(common.DefaultConfig.SystemAccountName)
}

func (n *UnlinkAuth) GetName() common.ActionName {
	name := common.N("unlinkauth")
	return common.ActionName(name)
}

type CancelDelay struct {
	CancelingAuth PermissionLevel
	TrxId         common.TransactionIdType
}

func (n *CancelDelay) GetAccount() common.AccountName {
	return common.AccountName(common.DefaultConfig.SystemAccountName)
}

func (n *CancelDelay) GetName() common.ActionName {
	name := common.N("unlinkauth")
	return common.ActionName(name)
}

type OnError struct {
	CancelingAuth PermissionLevel
	TrxId         common.TransactionIdType
}

func (n *OnError) GetAccount() common.AccountName {
	return common.AccountName(common.DefaultConfig.SystemAccountName)
}

func (n *OnError) GetName() common.ActionName {
	name := common.N("onError")
	return common.ActionName(name)
}