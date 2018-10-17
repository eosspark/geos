package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	//"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/entity"
	//"strings"
)

func transactionIdToSenderId(tid common.TransactionIdType) *common.Uint128 {
	id := &common.Uint128{tid.Hash[3], tid.Hash[2]}
	return id
}

func validateAuthorityPrecondition(context *ApplyContext, auth *types.Authority) {

	//Obj := &entity.AccountObject{}
	for _, a := range auth.Accounts {

		//Obj = &entity.AccountObject{Name: a.Permission.Actor}
		//err := context.DB.Get("byName", Obj)

		// EOS_ASSERT( err != nil, action_validate_exception,
		//           "account '${account}' does not exist",
		//           ("account", a.Permission.Actor))

		// account was already checked to exist, so its owner and active permissions should exist
		if a.Permission.Permission == common.PermissionName(common.DefaultConfig.OwnerName) ||
			a.Permission.Permission == common.PermissionName(common.DefaultConfig.ActiveName) {
			continue
		}

		// virtual eosio.code permission does not really exist but is allowed
		if a.Permission.Permission == common.PermissionName(common.DefaultConfig.EosioCodeName) {
			continue
		}

		//try{
		context.Control.GetAuthorizationManager().GetPermission(&types.PermissionLevel{a.Permission.Actor, a.Permission.Permission})
		// } catch(e *permission_query_exception) {

		//   EOS_THROW( action_validate_exception,
		//             "permission '${perm}' does not exist",
		//             ("perm", a.Permission))
		// }
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

func applyEosioNewaccount(context *ApplyContext) {

	create := &newAccount{}
	//try{
	rlp.DecodeBytes(context.Act.Data, create)

	context.RequireAuthorization(int64(create.Createor))

	//assert(types.Validate(create.Owner), action_validate_exception, "Invalid owner authority")
	//assert(types.Validate(create.Active), action_validate_exception, "Invalid owner authority")

	//db := context.DB
	//nameStr := common.S(uint64(create.Name))

	//assert(empty(uint64(create.Name)), action_validate_exception, "account name cannot be empty")
	//assert(len(nameStr) <= 12, action_validate_exception, "account names can only be 12 chars long")

	//creator := &entity.AccountObject{Name: create.Createor}
	//err := context.DB.Get("byName", creator)
	//if err != nil && !creator.Privileged {
	//
	//	assert(strings.Index(nameStr, "eosio.") != 0, action_validate_exception,
	//		"only privileged accounts can have names that start with 'eosio.'")
	//
	//}

	//newAccountObject := &entity.AccountObject{Name: create.Name, CreationDate: common.BlockTimeStamp(context.Control.PendingBlockTime())}
	//db.Insert(newAccountObject)
	//
	//newAccountSequenceObj := &entity.AccountSequenceObject{Name: create.Name}
	//db.Insert(newAccountSequenceObj)

	validateAuthorityPrecondition(context, &create.Owner)
	validateAuthorityPrecondition(context, &create.Active)

	//authorization := context.Control.GetMutableAuthorizationManager()
	//ownerPemission := authorization.CreatePermission(create.Name, common.DefaultConfig.OwnerName, 0, create.Owner)
	//activePemission := authorization.CreatePermission(create.Name, common.DefaultConfig.ActiveName, ownerPemission.ID, create.Owner)
	//
	//context.Control.GetMutableResourceLimitsManager().InitializeAccount(create.Name)
	//ramDelta := common.DefaultConfig.OverheadPerAccountRamBytes
	//ramDelta += uint32(2 * common.BillableSize["permission_object"].Value)
	//ramDelta += ownerPemission.Auth.GetBillableSize()
	//ramDelta += activePemission.Auth.GetBillableSize()

	//context.AddRamUsage(create.Name, int64(ramDelta))
	//}capture_and_rethrow(create)

}

func applyEosioSetcode(context *ApplyContext) {

	//cfg := context.Control.GetGlobalProperties()

	//db := context.DB
	act := &setCode{}
	rlp.DecodeBytes(context.Act.Data, act)

	context.RequireAuthorization(int64(act.Account))

	//assert(act.vmType == 0, invalid_contract_vm_type, "code should be 0")
	//assert(act.vmVersion == 0, invalid_contract_vm_version, "version should be 0")

	//codeId := &crypto.Sha256{}
	//if len(act.Code) > 0 {
	//	codeId = crypto.NewSha256Byte(act.Code)
	//	//exec.validate(context.Control, act.Code)
	//}

	accountObject := &entity.AccountObject{Name: act.Account}
	//db.Get("byName", accountObject)
	codeSize := len(act.Code)
	oldSize := len(accountObject.Code) * int(common.DefaultConfig.SetcodeRamBytesMultiplier)
	newSize := codeSize * int(common.DefaultConfig.SetcodeRamBytesMultiplier)

	//assert(accountObject.CodeVersion != codeId, set_exact_code, "contract is already running this version of code")

	//db.modify(accountObject, func(a *entity.AccountObject) {
	//	a.LastCodeUpdate = context.Control.PendingBlockTime()
	//	a.CodeVersion = *codeId
	//	if codeSize > 0 {
	//		a.Code = act.Code
	//	}
	//})

	//accountSequenceObj := &entity.AccountSequenceObject{Name: act.Account}
	//db.modify(accountSequenceObj, func(aso *entity.AccountSequenceObject) {
	//	aso.CodeSequence += 1
	//})

	if newSize != oldSize {
		context.AddRamUsage(act.Account, int64(newSize-oldSize))
	}

}

func applyEosioSetabi(context *ApplyContext) {

	//db := context.DB
	act := &setAbi{}
	rlp.DecodeBytes(context.Act.Data, act)

	context.RequireAuthorization(int64(act.Account))

	accountObject := &entity.AccountObject{Name: act.Account}
	//db.Get("byName", accountObject)

	abiSize := len(act.Abi)
	oldSize := len(accountObject.Abi)
	newSize := abiSize

	//db.Modify(accountObject, func(a *entity.AccountObject) {
	//	if abiSize > 0 {
	//		a.Abi = act.Abi
	//	}
	//})

	//accountSequenceObj := &entity.AccountSequenceObject{Name: act.Account}
	//db.modify(accountSequenceObj, func(aso *entity.AccountSequenceObject) {
	//	aso.CodeSequence += 1
	//})

	if newSize != oldSize {
		context.AddRamUsage(act.Account, int64(newSize-oldSize))
	}

}

func applyEosioUpdateauth(context *ApplyContext) {

	update := &updateAuth{}
	rlp.DecodeBytes(context.Act.Data, update)
	context.RequireAuthorization(int64(update.Account))

	//db := context.DB

	// EOS_ASSERT(!update.permission.empty(), action_validate_exception, "Cannot create authority with empty name");
	// EOS_ASSERT( update.permission.to_string().find( "eosio." ) != 0, action_validate_exception,
	//             "Permission names that start with 'eosio.' are reserved" );
	// EOS_ASSERT(update.permission != update.parent, action_validate_exception, "Cannot set an authority as its own parent");
	// db.get<account_object, by_name>(update.account);
	// EOS_ASSERT(validate(update.auth), action_validate_exception,
	//            "Invalid authority: ${auth}", ("auth", update.auth));
	// if( update.permission == config::active_name )
	//    EOS_ASSERT(update.parent == config::owner_name, action_validate_exception, "Cannot change active authority's parent from owner", ("update.parent", update.parent) );
	// if (update.permission == config::owner_name)
	//    EOS_ASSERT(update.parent.empty(), action_validate_exception, "Cannot change owner authority's parent");
	// else
	//    EOS_ASSERT(!update.parent.empty(), action_validate_exception, "Only owner permission can have empty parent" );

	// if( update.auth.waits.size() > 0 ) {
	//    auto max_delay = context.control.get_global_properties().configuration.max_transaction_delay;
	//    EOS_ASSERT( update.auth.waits.back().wait_sec <= max_delay, action_validate_exception,
	//                "Cannot set delay longer than max_transacton_delay, which is ${max_delay} seconds",
	//                ("max_delay", max_delay) );
	// }

	validateAuthorityPrecondition(context, &update.Auth)

	//authorization := context.Control.GetMutableAuthorizationManager()
	//permission := authorization.FindPermission(&types.PermissionLevel{update.Account, update.Parent})
	//
	//parentId := entity.IdType(0)
	//if update.Permission != common.PermissionName(common.DefaultConfig.OwnerName) {
	//	parent := authorization.GetPermission(&types.PermissionLevel{update.Account, update.Parent})
	//	parentId = parent.ID
	//}
	//
	//if permission != nil {
	//	oldSize := common.BillableSize["permission_object"].Value + permission.Auth.GetBillableSize()
	//	authorization.ModifyPermission(*permission, update.Auth)
	//	newSize := common.BillableSize["permission_object"].Value + permission.Auth.GetBillableSize()
	//	context.AddRamUsage(permission.Owner, newSize-oldSize)
	//} else {
	//
	//	p := authorization.CreatePermission(update, Account, update.Permission, parentId, update.Auth)
	//	newSize := common.BillableSize["permission_object"].Value + p.Auth.GetBillableSize()
	//
	//	context.AddRamUsage(update.Account, newSize)
	//}

}

func applyEosioDeleteauth(context *ApplyContext) {

	remove := &deleteAuth{}
	rlp.DecodeBytes(context.Act.Data, remove)
	context.RequireAuthorization(int64(remove.Account))

	//assert(remove.Permission != common.PermissionName(common.DefaultConfig.OwnerName), action_validate_exception, "Cannot delete active authority")
	//assert(remove.Permission != common.PermissionName(common.DefaultConfig.ActiveName), action_validate_exception, "Cannot delete owner authority")

	//db := context.DB

	{
		permissionLinkIndexObject := &entity.PermissionLinkObject{}
		//index := db.GetIndex("byPermissionName", permissionLinkIndexObject)

		permissionLinkIndexObject.Account = remove.Account
		permissionLinkIndexObject.RequiredPermission = remove.Permission
		//r := index.EqualRange(permissionLinkIndexObject
		// assert(r.first == r.second,  action_validate_exception,
		//             "Cannot delete a linked authority. Unlink the authority first. This authority is linked to ${code}::${type}.",
		//             ("code", string(r.first.GetObject().code))("type", string(r.first.GetObject().message_type)))
	}

	//authorization := context.Control.GetMutableAuthorizationManager()
	//permission := authorization.GetPermission(&types.PermissionLevel{remove.Account, remove.Permission})
	//
	//oldSize := common.DefaultConfig.BillableSize["permission_object"] + permission.Auth.GetBillableSize()
	//authorization.RemovePermission(permission)

	//context.AddRamUsage(remove.Account, -oldSize)

}

func applyEosioLinkauth(context *ApplyContext) {

	requirement := &linkAuth{}
	rlp.DecodeBytes(context.Act.Data, requirement)

	//	assert(!empty(uint64(requirement.Requirement)), action_validate_exception, "Required permission cannot be empty")
	context.RequireAuthorization(int64(requirement.Account))
	//db := context.DB

	//accountObject := &entity.AccountObject{Name: requirement.Account}
	//err := db.Get("byName", accountObject)
	// assert(err != nil, account_query_exception,
	//                 "Failed to retrieve account: ${account}", ("account", requirement.Account))

	//codeObject := &entity.AccountObject{Name: requirement.Code}
	//err = db.Get("byName", codeObject)
	// assert(err != nil, account_query_exception,
	//                 "Failed to retrieve account: ${account}", ("account", requirement.Account))

	if requirement.Requirement != common.PermissionName(common.DefaultConfig.EosioAnyName) {

		//permissionObject := &entity.PermissionObject{Name: requirement.Requirement}
		//err = db.Get("byName", permissionObject)
		// assert(err != nil, permission_query_exception,
		//                "Failed to retrieve permission: ${permission}", ("permission", requirement.requirement))
	}

	//permissionLinkObject := &entity.PermissionLinkObject{Account: requirement.Account, Code: requirement.Code, MessageType: requirement.Type}
	//err = db.Get("byActionName", permissionLinkObject)

	//if err != nil {
	//	// assert(permissionLinkObject.RequiredPermission != requirement.Requirement, action_validate_exception,
	//	// 	"Attempting to update required authority, but new requirement is same as old")
	//
	//	// db.modify(*link, [requirement = requirement.requirement](permission_link_object& link) {
	//	//           link.required_permission = requirement;
	//	//       });
	//	db.Modify(permissionLinkObject, func(link *entity.PermissionLinkObject) {
	//		link.RequiredPermission = requirement
	//	})
	//
	//} else {
	//	permissionLinkObject.RequiredPermission = requirement.Requirement
	//	db.Insert(permissionLinkObject)
	//
	//	context.AddRamUsage(permissionLinkObject.Account, common.BillableSize("PermissionLinkObject"))
	//}

}

func applyEosioUnlinkauth(context *ApplyContext) {

	//db := context.DB

	unlink := &unlinkAuth{}
	rlp.DecodeBytes(context.Act.Data, unlink)

	context.RequireAuthorization(int64(unlink.Account))

	link := &entity.PermissionLinkObject{Account: unlink.Account, Code: unlink.Code, MessageType: unlink.Type}
	//err := db.Get("byActionName", link)

	//assert(err != nil, action_validate_exception, "Attempting to unlink authority, but no link found")
	context.AddRamUsage(link.Account, -int64(common.BillableSizeV("permission_link_object")))
	//db.Remove(link)
}

func applyEosioCanceldalay(context *ApplyContext) {

	cancel := &cancelDelay{}
	rlp.DecodeBytes(context.Act.Data, cancel)

	context.RequireAuthorization(int64(cancel.cancelingAuth.Actor))
	trxId := cancel.TrxId

	context.CancelDeferredTransaction2(transactionIdToSenderId(trxId), common.AccountName(0))

}
