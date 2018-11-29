package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	arithmetic "github.com/eosspark/eos-go/common/arithmetic_types"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"strings"
)

func transactionIdToSenderId(tid common.TransactionIdType) *arithmetic.Uint128 {
	id := &arithmetic.Uint128{tid.Hash[3], tid.Hash[2]}
	return id
}

func validateAuthorityPrecondition(context *ApplyContext, auth *types.Authority) {

	for _, a := range auth.Accounts {

		Obj := &entity.AccountObject{Name: a.Permission.Actor}
		err := context.DB.Find("byName", Obj, &Obj)

		EosAssert(err == nil, &ActionValidateException{},
			"account '${account}' does not exist",
			common.S(uint64(a.Permission.Actor)))

		// account was already checked to exist, so its owner and active permissions should exist
		if a.Permission.Permission == common.PermissionName(common.DefaultConfig.OwnerName) ||
			a.Permission.Permission == common.PermissionName(common.DefaultConfig.ActiveName) {
			continue
		}

		// virtual eosio.code permission does not really exist but is allowed
		if a.Permission.Permission == common.PermissionName(common.DefaultConfig.EosioCodeName) {
			continue
		}

		Try(func() {
			context.Control.GetAuthorizationManager().GetPermission(&types.PermissionLevel{a.Permission.Actor, a.Permission.Permission})
		}).Catch(func(e *PermissionQueryException) {
			//   EOS_THROW( action_validate_exception,
			//             "permission '${perm}' does not exist",
			//             ("perm", a.Permission))
			// }
			EosThrow(&ActionValidateException{}, "permission '%s' does not exist", a.Permission.String())
		}).End()
	}

	if context.Control.IsProducingBlock() {
		for _, p := range auth.Keys {
			context.Control.CheckKeyList(&p.Key)
		}
	}

}

func empty(a uint64) bool {
	if a == 0 {
		return true
	}
	return false
}

func ApplyEosioNewaccount(context *ApplyContext) {

	create := &NewAccount{}
	rlp.DecodeBytes(context.Act.Data, create)

	//try.Try()
	context.RequireAuthorization(int64(create.Creator))

	EosAssert(types.Validate(create.Owner), &ActionValidateException{}, "Invalid owner authority")
	EosAssert(types.Validate(create.Active), &ActionValidateException{}, "Invalid owner authority")

	db := context.DB
	nameStr := common.S(uint64(create.Name))

	EosAssert(!empty(uint64(create.Name)), &ActionValidateException{}, "account name cannot be empty")
	EosAssert(len(nameStr) <= 12, &ActionValidateException{}, "account names can only be 12 chars long")

	// Check if the creator is privileged
	creator := entity.AccountObject{Name: create.Creator}
	err := context.DB.Find("byName", creator, &creator)
	if err != nil && !creator.Privileged {

		EosAssert(strings.Index(nameStr, "eosio.") != 0, &ActionValidateException{},
			"only privileged accounts can have names that start with 'eosio.'")

	}

	existingAccount := entity.AccountObject{Name: create.Name}
	err = db.Find("byName", existingAccount, &existingAccount)
	EosAssert(err != nil, &AccountNameExistsException{}, "Cannot create account named ${name}, as that name is already taken", common.S(uint64(create.Name)))

	newAccountObject := entity.AccountObject{Name: create.Name, CreationDate: types.BlockTimeStamp(context.Control.PendingBlockTime())}
	db.Insert(&newAccountObject)

	newAccountSequenceObj := entity.AccountSequenceObject{Name: create.Name}
	db.Insert(&newAccountSequenceObj)

	validateAuthorityPrecondition(context, &create.Owner)
	validateAuthorityPrecondition(context, &create.Active)

	authorization := context.Control.GetMutableAuthorizationManager()
	ownerPemission := authorization.CreatePermission(create.Name, common.DefaultConfig.OwnerName, 0, create.Owner, common.TimePoint(0))
	activePemission := authorization.CreatePermission(create.Name, common.DefaultConfig.ActiveName, PermissionIdType(ownerPemission.ID), create.Owner, common.TimePoint(0))

	context.Control.GetMutableResourceLimitsManager().InitializeAccount(create.Name)
	ramDelta := uint64(common.DefaultConfig.OverheadPerAccountRamBytes)
	ramDelta += 2 * common.BillableSizeV("permission_object")
	ramDelta += ownerPemission.Auth.GetBillableSize()
	ramDelta += activePemission.Auth.GetBillableSize()

	context.AddRamUsage(create.Name, int64(ramDelta))
	//}capture_and_rethrow(create)

}

func applyEosioNewaccount(context *ApplyContext) {

	create := &newAccount{}
	rlp.DecodeBytes(context.Act.Data, create)

	//try.Try()
	context.RequireAuthorization(int64(create.Creator))

	EosAssert(types.Validate(create.Owner), &ActionValidateException{}, "Invalid owner authority")
	EosAssert(types.Validate(create.Active), &ActionValidateException{}, "Invalid owner authority")

	db := context.DB
	nameStr := common.S(uint64(create.Name))

	EosAssert(!empty(uint64(create.Name)), &ActionValidateException{}, "account name cannot be empty")
	EosAssert(len(nameStr) <= 12, &ActionValidateException{}, "account names can only be 12 chars long")

	// Check if the creator is privileged
	creator := entity.AccountObject{Name: create.Creator}
	err := context.DB.Find("byName", creator, &creator)
	if err != nil && !creator.Privileged {

		EosAssert(strings.Index(nameStr, "eosio.") != 0, &ActionValidateException{},
			"only privileged accounts can have names that start with 'eosio.'")

	}

	existingAccount := entity.AccountObject{Name: create.Name}
	err = db.Find("byName", existingAccount, &existingAccount)
	EosAssert(err != nil, &AccountNameExistsException{}, "Cannot create account named ${name}, as that name is already taken", common.S(uint64(create.Name)))

	newAccountObject := entity.AccountObject{Name: create.Name, CreationDate: types.BlockTimeStamp(context.Control.PendingBlockTime())}
	db.Insert(&newAccountObject)

	newAccountSequenceObj := entity.AccountSequenceObject{Name: create.Name}
	db.Insert(&newAccountSequenceObj)

	validateAuthorityPrecondition(context, &create.Owner)
	validateAuthorityPrecondition(context, &create.Active)

	authorization := context.Control.GetMutableAuthorizationManager()
	ownerPemission := authorization.CreatePermission(create.Name, common.DefaultConfig.OwnerName, 0, create.Owner, common.TimePoint(0))
	activePemission := authorization.CreatePermission(create.Name, common.DefaultConfig.ActiveName, PermissionIdType(ownerPemission.ID), create.Owner, common.TimePoint(0))

	context.Control.GetMutableResourceLimitsManager().InitializeAccount(create.Name)
	ramDelta := uint64(common.DefaultConfig.OverheadPerAccountRamBytes)
	ramDelta += 2 * common.BillableSizeV("permission_object")
	ramDelta += ownerPemission.Auth.GetBillableSize()
	ramDelta += activePemission.Auth.GetBillableSize()

	context.AddRamUsage(create.Name, int64(ramDelta))
	//}capture_and_rethrow(create)

}

func applyEosioSetcode(context *ApplyContext) {

	//cfg := context.Control.GetGlobalProperties()

	db := context.DB
	act := setCode{}
	rlp.DecodeBytes(context.Act.Data, &act)

	context.RequireAuthorization(int64(act.Account))

	EosAssert(act.VmType == 0, &InvalidContractVmType{}, "code should be 0")
	EosAssert(act.VmVersion == 0, &InvalidContractVmVersion{}, "version should be 0")

	var codeId *crypto.Sha256
	if len(act.Code) > 0 {
		codeId = crypto.NewSha256Byte(act.Code)
		//exec.validate(context.Control, act.Code)
	}

	accountObject := entity.AccountObject{Name: act.Account}
	db.Find("byName", accountObject, &accountObject)

	codeSize := len(act.Code)
	oldSize := len(accountObject.Code) * int(common.DefaultConfig.SetcodeRamBytesMultiplier)
	newSize := codeSize * int(common.DefaultConfig.SetcodeRamBytesMultiplier)

	EosAssert(accountObject.CodeVersion != *codeId, &SetExactCode{}, "contract is already running this version of code")

	db.Modify(&accountObject, func(a *entity.AccountObject) {
		a.LastCodeUpdate = context.Control.PendingBlockTime()
		a.CodeVersion = *codeId
		if codeSize > 0 {
			a.Code = act.Code
		}
	})

	accountSequenceObj := entity.AccountSequenceObject{Name: act.Account}
	db.Modify(&accountSequenceObj, func(aso *entity.AccountSequenceObject) {
		aso.CodeSequence += 1
	})

	if newSize != oldSize {
		context.AddRamUsage(act.Account, int64(newSize-oldSize))
	}

}

func applyEosioSetabi(context *ApplyContext) {

	db := context.DB
	act := setAbi{}
	rlp.DecodeBytes(context.Act.Data, &act)

	context.RequireAuthorization(int64(act.Account))

	accountObject := entity.AccountObject{Name: act.Account}
	db.Find("byName", accountObject, &accountObject)

	abiSize := len(act.Abi)
	oldSize := len(accountObject.Abi)
	newSize := abiSize

	db.Modify(&accountObject, func(a *entity.AccountObject) {
		if abiSize > 0 {
			a.Abi = act.Abi
		}
	})

	accountSequenceObj := entity.AccountSequenceObject{Name: act.Account}
	db.Modify(&accountSequenceObj, func(aso *entity.AccountSequenceObject) {
		aso.CodeSequence += 1
	})

	if newSize != oldSize {
		context.AddRamUsage(act.Account, int64(newSize-oldSize))
	}

}

func applyEosioUpdateauth(context *ApplyContext) {

	update := updateAuth{}
	rlp.DecodeBytes(context.Act.Data, &update)
	context.RequireAuthorization(int64(update.Account))

	db := context.DB

	EosAssert(!common.Empty(&update.Permission), &ActionValidateException{}, "Cannot create authority with empty name")
	EosAssert(strings.Index(common.S(uint64(update.Permission)), "eosio.") != 0,
		&ActionValidateException{},
		"Permission names that start with 'eosio.' are reserved")
	EosAssert(update.Permission != update.Parent, &ActionValidateException{}, "Cannot set an authority as its own parent")

	accountObject := entity.AccountObject{Name: update.Account}
	db.Find("byName", accountObject, &accountObject)
	EosAssert(types.Validate(update.Auth), &ActionValidateException{},
		"Invalid authority: %s", update.Auth)

	if update.Permission == common.DefaultConfig.ActiveName {
		EosAssert(update.Parent == common.DefaultConfig.OwnerName,
			&ActionValidateException{},
			"Cannot change active authority's parent from owner, update.parent %s", common.S(uint64(update.Parent)))
	}

	if update.Permission == common.DefaultConfig.OwnerName {
		EosAssert(common.Empty(&update.Parent),
			&ActionValidateException{},
			"Cannot change owner authority's parent")
	} else {
		EosAssert(!common.Empty(&update.Parent),
			&ActionValidateException{},
			"Only owner permission can have empty parent")
	}

	if len(update.Auth.Waits) > 0 {
		maxDelay := context.Control.GetGlobalProperties().Configuration.MaxTrxDelay
		EosAssert(update.Auth.Waits[len(update.Auth.Waits)-1].WaitSec <= maxDelay, &ActionValidateException{},
			"Cannot set delay longer than max_transacton_delay, which is %d seconds", maxDelay)
	}

	validateAuthorityPrecondition(context, &update.Auth)

	authorization := context.Control.GetMutableAuthorizationManager()
	permission := authorization.FindPermission(&types.PermissionLevel{update.Account, update.Permission})

	parentId := common.IdType(0)
	if update.Permission != common.PermissionName(common.DefaultConfig.OwnerName) {
		parent := authorization.GetPermission(&types.PermissionLevel{update.Account, update.Parent})
		parentId = parent.ID
	}

	if permission != nil {
		oldSize := common.BillableSizeV("permission_object") + permission.Auth.GetBillableSize()
		authorization.ModifyPermission(permission, &update.Auth)
		newSize := common.BillableSizeV("permission_object") + permission.Auth.GetBillableSize()
		context.AddRamUsage(permission.Owner, int64(newSize-oldSize))
	} else {

		p := authorization.CreatePermission(update.Account, update.Permission, PermissionIdType(parentId), update.Auth, common.TimePoint(0))
		newSize := common.BillableSizeV("permission_object") + p.Auth.GetBillableSize()

		context.AddRamUsage(update.Account, int64(newSize))
	}

}

func applyEosioDeleteauth(context *ApplyContext) {

	remove := deleteAuth{}
	rlp.DecodeBytes(context.Act.Data, &remove)
	context.RequireAuthorization(int64(remove.Account))

	EosAssert(remove.Permission != common.PermissionName(common.DefaultConfig.OwnerName), &ActionValidateException{}, "Cannot delete active authority")
	EosAssert(remove.Permission != common.PermissionName(common.DefaultConfig.OwnerName), &ActionValidateException{}, "Cannot delete active authority")

	//db := context.DB

	{
		//permissionLinkIndexObject := entity.PermissionLinkObject{Account: remove.Account, RequiredPermission: remove.Permission}
		//index,_ := db.GetIndex("byPermissionName", &permissionLinkIndexObject)

		//r := index.EqualRange(permissionLinkIndexObject
		// assert(r.first == r.second,  action_validate_exception,
		//             "Cannot delete a linked authority. Unlink the authority first. This authority is linked to ${code}::${type}.",
		//             ("code", string(r.first.GetObject().code))("type", string(r.first.GetObject().message_type)))
	}

	authorization := context.Control.GetMutableAuthorizationManager()
	permission := authorization.GetPermission(&types.PermissionLevel{remove.Account, remove.Permission})
	oldSize := common.BillableSizeV("permission_object") + permission.Auth.GetBillableSize()
	authorization.RemovePermission(permission)

	context.AddRamUsage(remove.Account, -int64(oldSize))

}

func applyEosioLinkauth(context *ApplyContext) {

	requirement := linkAuth{}
	rlp.DecodeBytes(context.Act.Data, &requirement)

	EosAssert(!empty(uint64(requirement.Requirement)), &ActionValidateException{}, "Required permission cannot be empty")
	context.RequireAuthorization(int64(requirement.Account))

	db := context.DB

	accountObject := entity.AccountObject{Name: requirement.Account}
	err := db.Find("byName", accountObject, &accountObject)
	EosAssert(err == nil, &AccountQueryException{},
		"Failed to retrieve account: %s",
		common.S(uint64(requirement.Account)))

	codeObject := entity.AccountObject{Name: requirement.Code}
	err = db.Find("byName", codeObject, &codeObject)
	EosAssert(err == nil, &AccountQueryException{},
		"Failed to retrieve account: %s",
		common.S(uint64(requirement.Code)))

	if requirement.Requirement != common.PermissionName(common.DefaultConfig.EosioAnyName) {

		// permissionObject := entity.PermissionObject{Name: requirement.Requirement}
		// err = db.Find("byName", permissionObject, &permissionObject)
		permissionObject := entity.PermissionObject{Owner: requirement.Account, Name: requirement.Requirement}
		err = db.Find("byOwner", permissionObject, &permissionObject)
		EosAssert(err == nil, &PermissionQueryException{},
			"Failed to retrieve permission: %s",
			common.S(uint64(requirement.Requirement)))
	}

	permissionLinkObject := entity.PermissionLinkObject{
		Account:     requirement.Account,
		Code:        requirement.Code,
		MessageType: requirement.Type}
	err = db.Find("byActionName", permissionLinkObject, &permissionLinkObject)

	if err == nil {
		EosAssert(permissionLinkObject.RequiredPermission != requirement.Requirement, &ActionValidateException{},
			"Attempting to update required authority, but new requirement is same as old")

		db.Modify(&permissionLinkObject, func(link *entity.PermissionLinkObject) {
			link.RequiredPermission = common.PermissionName(requirement.Account)
			link.Code = requirement.Code
			link.MessageType = requirement.Type
			link.RequiredPermission = requirement.Requirement
		})

	} else {
		permissionLinkObject.RequiredPermission = requirement.Requirement
		db.Insert(&permissionLinkObject)

		context.AddRamUsage(permissionLinkObject.Account, int64(common.BillableSizeV("PermissionLinkObject")))
	}

}

func applyEosioUnlinkauth(context *ApplyContext) {

	db := context.DB

	unlink := unlinkAuth{}
	rlp.DecodeBytes(context.Act.Data, &unlink)

	context.RequireAuthorization(int64(unlink.Account))

	link := entity.PermissionLinkObject{Account: unlink.Account, Code: unlink.Code, MessageType: unlink.Type}
	err := db.Find("byActionName", link, &link)
	EosAssert(err == nil, &ActionValidateException{}, "Attempting to unlink authority, but no link found")

	context.AddRamUsage(link.Account, -int64(common.BillableSizeV("permission_link_object")))
	db.Remove(&link)
}

func applyEosioCanceldalay(context *ApplyContext) {

	cancel := &cancelDelay{}
	rlp.DecodeBytes(context.Act.Data, cancel)

	context.RequireAuthorization(int64(cancel.CancelingAuth.Actor))
	trxId := cancel.TrxId

	context.CancelDeferredTransaction2(transactionIdToSenderId(trxId), common.AccountName(0))

}
