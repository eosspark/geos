package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/common/eos_math"
)

type contractTypesInterface interface {
	GetAccount() common.AccountName
	GetName() common.ActionName
}

type NewAccount struct {
	Creator common.AccountName
	Name    common.AccountName
	Owner   types.Authority
	Active  types.Authority
}

func (n *NewAccount) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (n *NewAccount) GetName() common.ActionName {
	return common.ActionName(common.N("newaccount"))
}

type SetCode struct {
	Account   common.AccountName
	VmType    uint8
	VmVersion uint8
	Code      []byte
}

func (s *SetCode) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s *SetCode) GetName() common.ActionName {
	return common.ActionName(common.N("setcode"))
}

type SetAbi struct {
	Account common.AccountName
	Abi     []byte
}

func (s SetAbi) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s SetAbi) GetName() common.ActionName {
	return common.ActionName(common.N("setabi"))
}

type UpdateAuth struct {
	Account    common.AccountName
	Permission common.PermissionName
	Parent     common.PermissionName
	Auth       types.Authority
}

func (u UpdateAuth) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (u UpdateAuth) GetName() common.ActionName {
	return common.ActionName(common.N("updateauth"))
}

type DeleteAuth struct {
	Account    common.AccountName
	Permission common.PermissionName
}

func (d DeleteAuth) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (d DeleteAuth) GetName() common.ActionName {
	return common.ActionName(common.N("deleteauth"))
}

type LinkAuth struct {
	Account     common.AccountName
	Code        common.AccountName
	Type        common.ActionName
	Requirement common.PermissionName
}

func (l LinkAuth) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (l LinkAuth) GetName() common.ActionName {
	return common.ActionName(common.N("linkauth"))
}

type UnLinkAuth struct {
	Account common.AccountName
	Code    common.AccountName
	Type    common.ActionName
}

func (u UnLinkAuth) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (u UnLinkAuth) GetName() common.ActionName {
	return common.ActionName(common.N("unlinkauth"))
}

type CancelDelay struct {
	CancelingAuth types.PermissionLevel
	TrxId         common.TransactionIdType
}

func (c CancelDelay) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (c CancelDelay) GetName() common.ActionName {
	return common.ActionName(common.N("canceldelay"))
}

type OnError struct {
	SenderId eos_math.Uint128
	SentTrx  []byte
}

func NewOnError(sid eos_math.Uint128, data []byte) *OnError {
	oe := OnError{SenderId: sid, SentTrx: data}
	return &oe
}
func (o OnError) GetAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (o OnError) GetName() common.ActionName {
	return common.ActionName(common.N("onerror"))
}
