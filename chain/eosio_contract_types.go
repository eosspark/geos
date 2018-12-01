package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	arithmetic "github.com/eosspark/eos-go/common/arithmetic_types"
)

type contractTypesInterface interface {
	getAccount() common.AccountName
	getName() common.AccountName
}

type NewAccount struct {
	Creator common.AccountName
	Name    common.AccountName
	Owner   types.Authority
	Active  types.Authority
}

func (n *NewAccount) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (n *NewAccount) getName() common.AccountName {
	return common.AccountName(common.N("newaccount"))
}

type SetCode struct {
	Account   common.AccountName
	VmType    uint8
	VmVersion uint8
	Code      []byte
}

func (s *SetCode) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s *SetCode) getName() common.ActionName {
	return common.ActionName(common.N("setcode"))
}

type newAccount struct {
	Creator common.AccountName
	Name    common.AccountName
	Owner   types.Authority
	Active  types.Authority
}

func (n *newAccount) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (n *newAccount) getName() common.AccountName {
	return common.AccountName(common.N("newaccount"))
}

type setCode struct {
	Account   common.AccountName
	VmType    uint8
	VmVersion uint8
	Code      []byte
}

func (s *setCode) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s *setCode) getName() common.ActionName {
	return common.ActionName(common.N("setcode"))
}

type setAbi struct {
	Account common.AccountName
	Abi     []byte
}

func (s setAbi) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s setAbi) getName() common.ActionName {
	return common.ActionName(common.N("setabi"))
}

type updateAuth struct {
	Account    common.AccountName
	Permission common.PermissionName
	Parent     common.PermissionName
	Auth       types.Authority
}

func (u updateAuth) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (u updateAuth) getName() common.ActionName {
	return common.ActionName(common.N("updateauth"))
}

type deleteAuth struct {
	Account    common.AccountName
	Permission common.PermissionName
}

func (d deleteAuth) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (d deleteAuth) getName() common.ActionName {
	return common.ActionName(common.N("deleteauth"))
}

type linkAuth struct {
	Account     common.AccountName
	Code        common.AccountName
	Type        common.ActionName
	Requirement common.PermissionName
}

func (l linkAuth) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (l linkAuth) getName() common.ActionName {
	return common.ActionName(common.N("linkauth"))
}

type unlinkAuth struct {
	Account common.AccountName
	Code    common.AccountName
	Type    common.ActionName
}

func (u unlinkAuth) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (u unlinkAuth) getName() common.ActionName {
	return common.ActionName(common.N("unlinkauth"))
}

type cancelDelay struct {
	CancelingAuth types.PermissionLevel
	TrxId         common.TransactionIdType
}

func (c cancelDelay) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (c cancelDelay) getName() common.ActionName {
	return common.ActionName(common.N("canceldelay"))
}

type onError struct {
	SenderId arithmetic.Uint128
	SentTrx  []byte
}

func (o onError) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (o onError) getName() common.ActionName {
	return common.ActionName(common.N("onerror"))
}
