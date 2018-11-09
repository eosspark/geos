package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/log"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
)

var noopCheckTime *func()

type AuthorizationManager struct {
	control *Controller
	db      database.DataBase
}

func newAuthorizationManager(control *Controller) *AuthorizationManager {
	azInstance := &AuthorizationManager{}
	azInstance.control = control
	azInstance.db = control.DB
	return azInstance
}

type PermissionIdType common.IdType

func (a *AuthorizationManager) CreatePermission(account common.AccountName,
	name common.PermissionName,
	parent PermissionIdType,
	auth types.Authority,
	initialCreationTime common.TimePoint,
) *entity.PermissionObject {
	creationTime := initialCreationTime
	if creationTime == 1 {
		creationTime = a.control.PendingBlockTime()
	}

	permUsage := entity.PermissionUsageObject{}
	permUsage.LastUsed = creationTime
	err := a.db.Insert(&permUsage)
	if err != nil {
		log.Error("CreatePermission is error: %s", err)
	}

	perm := entity.PermissionObject{
		UsageId:     permUsage.ID,
		Parent:      common.IdType(parent),
		Owner:       account,
		Name:        name,
		LastUpdated: creationTime,
		Auth:        auth.ToSharedAuthority(),
	}
	err = a.db.Insert(&perm)
	if err != nil {
		log.Error("CreatePermission is error: %s", err)
	}
	return &perm
}

func (a *AuthorizationManager) ModifyPermission(permission *entity.PermissionObject, auth *types.Authority) {
	err := a.db.Modify(&permission, func(po *entity.PermissionObject) {
		po.Auth = (*auth).ToSharedAuthority()
		po.LastUpdated = a.control.PendingBlockTime()
	})
	if err != nil {
		log.Error("ModifyPermission is error: %s", err)
	}
}

func (a *AuthorizationManager) RemovePermission(permission *entity.PermissionObject) {
	index, err := a.db.GetIndex("byParent", entity.PermissionObject{})
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
	_, err = index.LowerBound(entity.PermissionObject{Parent: permission.ID})
	EosAssert(err == nil, &ActionValidateException{}, "Cannot remove a permission which has children. Remove the children first.")
	usage := entity.PermissionUsageObject{ID: permission.UsageId}
	err = a.db.Find("id", usage, &usage)
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
	err = a.db.Remove(usage)
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
	err = a.db.Remove(*permission)
	if err != nil {
		log.Error("RemovePermission is error: %s", err)
	}
}

func (a *AuthorizationManager) UpdatePermissionUsage(permission *entity.PermissionObject) {
	puo := entity.PermissionUsageObject{}
	puo.ID = permission.UsageId
	err := a.db.Find("id", puo, &puo)
	if err != nil {
		log.Error("UpdatePermissionUsage is error: %s", err)
	}
	err = a.db.Modify(&puo, func(p *entity.PermissionUsageObject) {
		puo.LastUsed = a.control.PendingBlockTime()
	})
	if err != nil {
		log.Error("UpdatePermissionUsage is error: %s", err)
	}
}

func (a *AuthorizationManager) GetPermissionLastUsed(permission *entity.PermissionObject) common.TimePoint {
	puo := entity.PermissionUsageObject{}
	puo.ID = permission.UsageId
	err := a.db.Find("id", puo, &puo)
	if err != nil {
		log.Error("GetPermissionLastUsed is error: %s", err)
	}
	return puo.LastUsed
}

func (a *AuthorizationManager) FindPermission(level *types.PermissionLevel) (p *entity.PermissionObject) { //TODO
	Try(func() {
		EosAssert(!level.Actor.Empty() && !level.Permission.Empty(), &InvalidPermission{}, "Invalid permission")
		po := entity.PermissionObject{}
		po.Owner = level.Actor
		po.Name = level.Permission
		err := a.db.Find("byOwner", po, &po)
		if err != nil {
			log.Error("FindPermission is error: %s", err)
		}
		p = &po
	}).Catch(func(e PermissionQueryException) {

	}).End()
	return p
}

func (a *AuthorizationManager) GetPermission(level *types.PermissionLevel) (p *entity.PermissionObject) {
	Try(func() {
		EosAssert(!level.Actor.Empty() && !level.Permission.Empty(), &InvalidPermission{}, "Invalid permission")
		po := entity.PermissionObject{}
		po.Owner = level.Actor
		po.Name = level.Permission
		err := a.db.Find("byOwner", po, &po)
		if err != nil {
			log.Error("GetPermission is error: %s", err)
		}
		p = &po
	}).Catch(func(e PermissionQueryException) {

	}).End()
	return p
}

func (a *AuthorizationManager) LookupLinkedPermission(authorizerAccount common.AccountName,
	scope common.AccountName,
	actName common.ActionName,
) (p *common.PermissionName) {
	Try(func() { //TODO
		link := entity.PermissionLinkObject{}
		link.Account = authorizerAccount
		link.Code = scope
		link.MessageType = actName
		err := a.db.Find("byActionName", link, &link)
		if err != nil {
			link.Code = common.AccountName(common.N(""))
			err = a.db.Find("byActionName", link, &link)
		}
		if err == nil {
			p = &link.RequiredPermission
			return
		}
	}).End()

	return p
}

func (a *AuthorizationManager) LookupMinimumPermission(authorizerAccount common.AccountName,
	scope common.AccountName,
	actName common.ActionName,
) (p *common.PermissionName) {
	if scope == common.DefaultConfig.SystemAccountName {
		EosAssert(actName != updateAuth{}.getName() &&
			actName != deleteAuth{}.getName() &&
			actName != linkAuth{}.getName() &&
			actName != unlinkAuth{}.getName() &&
			actName != cancelDelay{}.getName(),
			&UnlinkableMinPermissionAction{},"cannot call lookup_minimum_permission on native actions that are not allowed to be linked to minimum permissions")
	}
	Try(func() {
		linkedPermission := a.LookupLinkedPermission(authorizerAccount, scope, actName)
		if common.Empty(linkedPermission) {
			p = &common.DefaultConfig.ActiveName
			return
		}

		if *linkedPermission == common.PermissionName(common.DefaultConfig.EosioAnyName) {
			return
		}

		p = linkedPermission
		return
	}).End()
	return p
}

func (a *AuthorizationManager) CheckUpdateauthAuthorization(update updateAuth, auths []types.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "updateauth action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == update.Account, &IrrelevantAuthException{}, "the owner of the affected permission needs to be the actor of the declared authorization")
	minPermission := a.FindPermission(&types.PermissionLevel{Actor:update.Account, Permission:update.Permission})
	if minPermission == nil {
		permission := a.GetPermission(&types.PermissionLevel{Actor:update.Account, Permission:update.Permission})
		minPermission = permission
	}
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckUpdateauthAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*minPermission , permissionIndex), &IrrelevantAuthException{},
	"updateauth action declares irrelevant authority '%v'; minimum authority is %v", auth, types.PermissionLevel{update.Account,minPermission.Name})
}

func (a *AuthorizationManager) CheckDeleteauthAuthorization(del deleteAuth, auths []types.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "deleteauth action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == del.Account, &IrrelevantAuthException{}, "the owner of the affected permission needs to be the actor of the declared authorization")
	minPermission := a.GetPermission(&types.PermissionLevel{Actor:del.Account, Permission:del.Permission})
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckDeleteauthAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*minPermission , permissionIndex), &IrrelevantAuthException{},
	"deleteauth action declares irrelevant authority '%v'; minimum authority is %v", auth, types.PermissionLevel{minPermission.Owner,minPermission.Name})
}

func (a *AuthorizationManager) CheckLinkauthAuthorization(link linkAuth, auths []types.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "link action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == link.Account, &IrrelevantAuthException{}, "the owner of the affected permission needs to be the actor of the declared authorization")

	EosAssert(link.Type != updateAuth{}.getName(), &ActionValidateException{}, "Cannot link eosio::updateauth to a minimum permission")
	EosAssert(link.Type != deleteAuth{}.getName(), &ActionValidateException{}, "Cannot link eosio::deleteauth to a minimum permission")
	EosAssert(link.Type != linkAuth{}.getName(), &ActionValidateException{}, "Cannot link eosio::linkauth to a minimum permission")
	EosAssert(link.Type != unlinkAuth{}.getName(), &ActionValidateException{}, "Cannot link eosio::unlinkauth to a minimum permission")
	EosAssert(link.Type != cancelDelay{}.getName(), &ActionValidateException{}, "Cannot link eosio::canceldelay to a minimum permission")

	linkedPermissionName := a.LookupMinimumPermission(link.Account, link.Code, link.Type)
	if common.Empty(&linkedPermissionName) {
		return
	}
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckLinkauthAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*a.GetPermission(&types.PermissionLevel{link.Account, *linkedPermissionName}) , permissionIndex), &IrrelevantAuthException{},
	"linkauth action declares irrelevant authority '%v'; minimum authority is %v", auth, types.PermissionLevel{link.Account,*linkedPermissionName})
}

func (a *AuthorizationManager) CheckUnlinkauthAuthorization(unlink unlinkAuth, auths []types.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "unlink action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == unlink.Account, &IrrelevantAuthException{},
		"the owner of the affected permission needs to be the actor of the declared authorization")

	unlinkedPermissionName := a.LookupLinkedPermission(unlink.Account, unlink.Code, unlink.Type)
	EosAssert(!common.Empty(&unlinkedPermissionName), &TransactionException{},
		"cannot unlink non-existent permission link of account '%v' for actions matching '%v::%v", unlink.Account, unlink.Code, unlink.Type)

	if *unlinkedPermissionName == common.DefaultConfig.EosioAnyName {
		return
	}
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckUnlinkauthAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*a.GetPermission(&types.PermissionLevel{unlink.Account, *unlinkedPermissionName}) , permissionIndex), &IrrelevantAuthException{},
	"unlink action declares irrelevant authority '%v'; minimum authority is %v", auth, types.PermissionLevel{unlink.Account, *unlinkedPermissionName})
}

func (a *AuthorizationManager) CheckCanceldelayAuthorization(cancel cancelDelay, auths []types.PermissionLevel) common.Microseconds {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "canceldelay action should only have one declared authorization")
	auth := auths[0]
	permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	if err != nil {
		log.Error("CheckCanceldelayAuthorization is error: %s", err)
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*a.GetPermission(&cancel.CancelingAuth), permissionIndex), &IrrelevantAuthException{},
	"canceldelay action declares irrelevant authority '%v'; specified authority to satisfy is %v", auth, cancel.CancelingAuth)

	generatedTrx := entity.GeneratedTransactionObject{}
	trxId := cancel.TrxId
	generatedIndex, err := a.control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
	if err != nil {
		log.Error("CheckCanceldelayAuthorization is error: %s", err)
	}
	itr, err := generatedIndex.LowerBound(entity.GeneratedTransactionObject{TrxId: trxId})
	if err != nil {
		log.Error("CheckCanceldelayAuthorization is error: %s", err)
	}

	generatedIndex.BeginData(&generatedTrx)
	EosAssert(!generatedIndex.CompareEnd(itr) && generatedTrx.TrxId == trxId, &TxNotFound{},
		"cannot cancel trx_id=%v, there is no deferred transaction with that transaction id", trxId)

	trx := types.Transaction{}
	rlp.DecodeBytes(generatedTrx.PackedTrx, &trx)
	found := false
	for _, act := range trx.Actions {
		for _, auth := range act.Authorization {
			if auth == cancel.CancelingAuth {
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	EosAssert(found, &ActionValidateException{}, "canceling_auth in canceldelay action was not found as authorization in the original delayed transaction")
	return common.Milliseconds(int64(generatedTrx.DelayUntil) - int64(generatedTrx.Published))
}

func (a *AuthorizationManager) CheckAuthorization(actions []*types.Action,
	providedKeys *common.FlatSet,
	providedPermissions *common.FlatSet,
	providedDelay common.Microseconds,
	checkTime *func(),
	allowUnusedKeys bool,
) {
	delayMaxLimit := common.Seconds(int64(a.control.GetGlobalProperties().Configuration.MaxTrxDelay))
	var effectiveProvidedDelay common.Microseconds
	if providedDelay >= delayMaxLimit {
		effectiveProvidedDelay = common.MaxMicroseconds()
	} else {
		effectiveProvidedDelay = providedDelay
	}
	checker := types.MakeAuthChecker(func(p *types.PermissionLevel) types.SharedAuthority { return a.GetPermission(p).Auth },
		a.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
		providedKeys,
		providedPermissions,
		effectiveProvidedDelay,
		checkTime)
	permissionToSatisfy := make(map[types.PermissionLevel]common.Microseconds)

	for _, act := range actions {
		specialCase := false
		delay := effectiveProvidedDelay

		if act.Account == common.DefaultConfig.SystemAccountName {
			specialCase = true
			switch act.Name {
			case updateAuth{}.getName():
				updateAuth := updateAuth{}
				rlp.DecodeBytes(act.Data, &updateAuth)
				a.CheckUpdateauthAuthorization(updateAuth, act.Authorization)

			case deleteAuth{}.getName():
				deleteAuth := deleteAuth{}
				rlp.DecodeBytes(act.Data, &deleteAuth)
				a.CheckDeleteauthAuthorization(deleteAuth, act.Authorization)

			case linkAuth{}.getName():
				linkAuth := linkAuth{}
				rlp.DecodeBytes(act.Data, &linkAuth)
				a.CheckLinkauthAuthorization(linkAuth, act.Authorization)

			case unlinkAuth{}.getName():
				unlinkAuth := unlinkAuth{}
				rlp.DecodeBytes(act.Data, &unlinkAuth)
				a.CheckUnlinkauthAuthorization(unlinkAuth, act.Authorization)

			case cancelDelay{}.getName():
				cancelDelay := cancelDelay{}
				rlp.DecodeBytes(act.Data, &cancelDelay)
				a.CheckCanceldelayAuthorization(cancelDelay, act.Authorization)

			default:
				specialCase = false
			}
		}

		for _, declaredAuth := range act.Authorization {
			(*checkTime)()
			if !specialCase {
				minPermissionName := a.LookupMinimumPermission(declaredAuth.Actor, act.Account, act.Name)
				if *minPermissionName != common.PermissionName(0) {
					minPermission := a.GetPermission(&types.PermissionLevel{Actor:declaredAuth.Actor, Permission:*minPermissionName})
					permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
					if err != nil {
						log.Error("CheckAuthorization is error: %s", err)
					}
					EosAssert(a.GetPermission(&declaredAuth).Satisfies(*minPermission , permissionIndex), &IrrelevantAuthException{},
						"action declares irrelevant authority '%v'; minimum authority is %v", declaredAuth, types.PermissionLevel{minPermission.Owner, minPermission.Name})
				}
			}

			isExist := false
			for first, second := range permissionToSatisfy{
				if first == declaredAuth {
					if second > delay {
						second = delay
						isExist = true
						break
					}
				}
			}
			if !isExist{
				permissionToSatisfy[declaredAuth] = delay
			}
		}
	}
	for p, q := range permissionToSatisfy {
		(*checkTime)()
		EosAssert(checker.SatisfiedLoc(&p, q, nil), &UnsatisfiedAuthorization{},
			"transaction declares authority '%v', "+
				"but does not have signatures for it under a provided delay of %v ms, "+
				"provided permissions %v, and provided keys %v", p, providedDelay.Count()/1000, providedPermissions, providedKeys)
	}
	if !allowUnusedKeys {
		EosAssert(checker.AllKeysUsed(), &TxIrrelevantSig{}, "transaction bears irrelevant signatures from these keys: ${keys}")
	}
}

func (a *AuthorizationManager) CheckAuthorization2(account common.AccountName,
	permission common.PermissionName,
	providedKeys *common.FlatSet, //flat_set<public_key_type>
	providedPermissions *common.FlatSet, //flat_set<permission_level>
	providedDelay common.Microseconds,
	checkTime *func(),
	allowUnusedKeys bool,
) {
	delayMaxLimit := common.Seconds(int64(a.control.GetGlobalProperties().Configuration.MaxTrxDelay))
	var effectiveProvidedDelay common.Microseconds
	if providedDelay >= delayMaxLimit {
		effectiveProvidedDelay = common.MaxMicroseconds()
	} else {
		effectiveProvidedDelay = providedDelay
	}
	checker := types.MakeAuthChecker(func(p *types.PermissionLevel) types.SharedAuthority { return a.GetPermission(p).Auth },
		a.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
		providedKeys,
		providedPermissions,
		effectiveProvidedDelay,
		checkTime)
	EosAssert(checker.SatisfiedLc(&types.PermissionLevel{account, permission}, nil), &UnsatisfiedAuthorization{},
		"permission '%v' was not satisfied under a provided delay of %v ms, provided permissions %v, and provided keys %v",
		types.PermissionLevel{account, permission}, providedDelay.Count()/1000, providedPermissions, providedKeys)

	if !allowUnusedKeys {
		EosAssert(checker.AllKeysUsed(), &TxIrrelevantSig{}, "irrelevant keys provided: %v", checker.GetUnusedKeys())
	}
}

func (a *AuthorizationManager) GetRequiredKeys(trx *types.Transaction,
	candidateKeys *common.FlatSet,
	providedDelay common.Microseconds) common.FlatSet {
	checker := types.MakeAuthChecker(func(p *types.PermissionLevel) types.SharedAuthority { return a.GetPermission(p).Auth },
		a.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
		candidateKeys,
		nil,
		providedDelay,
		noopCheckTime)
	return checker.GetUsedKeys()
}
