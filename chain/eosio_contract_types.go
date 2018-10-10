package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

type contractTypesInterface interface {
	getAccount() common.AccountName
	getName() common.AccountName
}

type newAccount struct {
	Createor common.AccountName
	Name     common.AccountName
	Owner    types.Authority
	Active   types.Authority
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

func (s *setCode) getName() common.AccountName {
	return common.AccountName(common.N("setcode"))
}

type setAbi struct {
	Account common.AccountName
	Abi     []byte
}

func (s *setAbi) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s *setAbi) getName() common.AccountName {
	return common.AccountName(common.N("setabi"))
}

type updateAuth struct {
	Account    common.AccountName
	Permission types.Permission
	Parent     types.Permission
	Auth       types.Authority
}

func (s *updateAuth) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s *updateAuth) getName() common.AccountName {
	return common.AccountName(common.N("updateauth"))
}

type deleteAuth struct {
	Account    common.AccountName
	Permission common.PermissionName
}

func (s *deleteAuth) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s *deleteAuth) getName() common.AccountName {
	return common.AccountName(common.N("deleteauth"))
}

type linkAuth struct {
	Account     common.AccountName
	Code        common.AccountName
	Type        common.ActionName
	Requirement common.PermissionName
}

func (s *linkAuth) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s *linkAuth) getName() common.AccountName {
	return common.AccountName(common.N("linkauth"))
}

type unlinkAuth struct {
	Account common.AccountName
	Code    common.AccountName
	Type    common.ActionName
}

func (s *unlinkAuth) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s *unlinkAuth) getName() common.AccountName {
	return common.AccountName(common.N("unlinkauth"))
}

type cancelDelay struct {
	cancelingAuth types.PermissionLevel
	TrxId         common.TransactionIdType
}

func (s *cancelDelay) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s *cancelDelay) getName() common.AccountName {
	return common.AccountName(common.N("canceldelay"))
}

type onError struct {
	SenderId common.Uint128
	SentTrx  []byte
}

func (s *onError) getAccount() common.AccountName {
	return common.DefaultConfig.SystemAccountName
}

func (s *onError) getName() common.AccountName {
	return common.AccountName(common.N("onerror"))
}
