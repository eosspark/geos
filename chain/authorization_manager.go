package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/entity"
	. "github.com/eosspark/eos-go/exception"
	"log"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/crypto/rlp"
)

var noopCheckTime *func()

var IsActiveAz bool

type AuthorizationManager struct {
	control *Controller
	db      database.DataBase
}

func newAuthorizationManager(control *Controller) *AuthorizationManager {
	azInstance := &AuthorizationManager{}
	if !IsActiveAz {
		azInstance.control = control
		azInstance.db = control.DB
		IsActiveAz = true
	}
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
	a.db.Insert(&permUsage)

	perm := entity.PermissionObject{
		UsageId:     permUsage.ID,
		Parent:      common.IdType(parent),
		Owner:       account,
		Name:        name,
		LastUpdated: creationTime,
		Auth:        auth.ToSharedAuthority(),
	}
	a.db.Insert(&perm)
	return &perm
}

func (a *AuthorizationManager) ModifyPermission(permission *entity.PermissionObject, auth *types.Authority) {
	a.db.Modify(&permission, func(po *entity.PermissionObject) {
		po.Auth = (*auth).ToSharedAuthority()
		po.LastUpdated = a.control.PendingBlockTime()
	})
}

func (a *AuthorizationManager) RemovePermission(permission *entity.PermissionObject) {
	index, err := a.db.GetIndex("byParent", entity.PermissionObject{})
	if err != nil {
		log.Fatalln(err)
	}
	_, err = index.LowerBound(entity.PermissionObject{Parent:permission.ID})
	EosAssert(err == nil, &ActionValidateException{},"Cannot remove a permission which has children. Remove the children first.")
	usage := entity.PermissionUsageObject{ID: permission.UsageId}
	a.db.Find("id", usage, &usage)
	a.db.Remove(usage)
	a.db.Remove(*permission)
}

func (a *AuthorizationManager) UpdatePermissionUsage(permission *entity.PermissionObject) {
	puo := entity.PermissionUsageObject{}
	puo.ID = permission.UsageId
	a.db.Find("id", puo, &puo)
	a.db.Modify(&puo, func(p *entity.PermissionUsageObject) {
		puo.LastUsed = a.control.PendingBlockTime()
	})
}

func (a *AuthorizationManager) GetPermissionLastUsed(permission *entity.PermissionObject) common.TimePoint {
	puo := entity.PermissionUsageObject{}
	puo.ID = permission.UsageId
	a.db.Find("id", puo, &puo)
	return puo.LastUsed
}

func (a *AuthorizationManager) FindPermission(level *types.PermissionLevel) (p *entity.PermissionObject) { //TODO
	defer HandleReturn()
	Try(func(){
		EosAssert(!level.Actor.Empty() && !level.Permission.Empty(), &InvalidPermission{}, "Invalid permission")
		po := entity.PermissionObject{}
		po.Owner = level.Actor
		po.Name = level.Permission
		a.db.Find("byOwner", po, &po)
		p = &po
		Return()
	}).Catch(func (e PermissionQueryException){

	}).End()
	return
}

func (a *AuthorizationManager) GetPermission(level *types.PermissionLevel) (p *entity.PermissionObject) {
	defer HandleReturn()
	Try(func(){
		EosAssert(!level.Actor.Empty() && !level.Permission.Empty(), &InvalidPermission{}, "Invalid permission")
		po := entity.PermissionObject{}
		po.Owner = level.Actor
		po.Name = level.Permission
		a.db.Find("byOwner", po, &po)
		p = &po
		Return()
	}).Catch(func (e PermissionQueryException){

	}).End()
	return
}


func (a *AuthorizationManager) LookupLinkedPermission(authorizerAccount common.AccountName,
	scope common.AccountName,
	actName common.ActionName,
) (p common.PermissionName) {
	defer HandleReturn()
	Try(func() {         //TODO
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
			p = link.RequiredPermission
			Return()
		}
		p = common.PermissionName(common.N(""))
		Return()
	}).End()

	return
}

func (a *AuthorizationManager) LookupMinimumPermission(authorizerAccount common.AccountName,
	scope common.AccountName,
	actName common.ActionName,
) (pn common.PermissionName) {
	if scope == common.DefaultConfig.SystemAccountName {
		//TODO
	}
	defer HandleReturn()
	Try(func() {
		linkedPermission := a.LookupLinkedPermission(authorizerAccount, scope, actName)
		if linkedPermission == common.PermissionName(common.N("")) {
			pn = common.DefaultConfig.ActiveName
			Return()
		}

		if linkedPermission == common.PermissionName(common.DefaultConfig.EosioAnyName) {
			pn = common.PermissionName(common.N(""))
			Return()
		}

		pn = linkedPermission
		Return()
	}).End()
	return
}

func (a *AuthorizationManager) CheckUpdateauthAuthorization(update updateAuth, auths []types.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "updateauth action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == update.Account, &IrrelevantAuthException{}, "the owner of the affected permission needs to be the actor of the declared authorization")
	minPermission := a.FindPermission(&types.PermissionLevel{update.Account, update.Permission})
	if minPermission == nil {
		permission := a.GetPermission(&types.PermissionLevel{update.Account, update.Permission})
		minPermission = permission
	}
	//permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	EosAssert(a.GetPermission(&auth).Satisfies(*minPermission/*, permissionIndex*/), &IrrelevantAuthException{}, "") //TODO
}

func (a *AuthorizationManager) CheckDeleteauthAuthorization(del deleteAuth, auths []types.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "deleteauth action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == del.Account, &IrrelevantAuthException{}, "the owner of the affected permission needs to be the actor of the declared authorization")
	minPermission := a.GetPermission(&types.PermissionLevel{del.Account, del.Permission})
	//permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	EosAssert(a.GetPermission(&auth).Satisfies(*minPermission/*, permissionIndex*/), &IrrelevantAuthException{}, "") //TODO
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
	if &linkedPermissionName == nil {
		return
	}
	//permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	EosAssert(a.GetPermission(&auth).Satisfies(*a.GetPermission(&types.PermissionLevel{link.Account, linkedPermissionName})/*, permissionIndex*/), &IrrelevantAuthException{}, "") //TODO
}

func (a *AuthorizationManager) CheckUnlinkauthAuthorization(unlink unlinkAuth, auths []types.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "unlink action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == unlink.Account, &IrrelevantAuthException{},
	"the owner of the affected permission needs to be the actor of the declared authorization")

	unlinkedPermissionName := a.LookupLinkedPermission(unlink.Account, unlink.Code, unlink.Type)
	EosAssert(&unlinkedPermissionName != nil, &TransactionException{},
	"cannot unlink non-existent permission link of account '${account}' for actions matching '${code}::${action}'")//TODO

	if unlinkedPermissionName == common.DefaultConfig.EosioAnyName {
		return
	}
	//permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	EosAssert(a.GetPermission(&auth).Satisfies(*a.GetPermission(&types.PermissionLevel{unlink.Account, unlinkedPermissionName})/*, permissionIndex*/), &IrrelevantAuthException{}, "") //TODO
}

func (a *AuthorizationManager) CheckCanceldelayAuthorization(cancel cancelDelay, auths []types.PermissionLevel) common.Microseconds {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "canceldelay action should only have one declared authorization")
	auth := auths[0]
	EosAssert(a.GetPermission(&auth).Satisfies(*a.GetPermission(&cancel.CancelingAuth)), &IrrelevantAuthException{}, "") //TODO

	generatedTrx := entity.GeneratedTransactionObject{}
	trxId := cancel.TrxId
	generatedIndex, err := a.control.DB.GetIndex("byTrxId", entity.GeneratedTransactionObject{})
	if err != nil {
		log.Fatalln(err)
	}
	itr, err := generatedIndex.LowerBound(entity.GeneratedTransactionObject{TrxId:trxId})
	if err != nil {
		log.Fatalln(err)
	}

	generatedIndex.Begin(&generatedTrx)
	EosAssert(!generatedIndex.CompareEnd(itr)&&generatedTrx.TrxId == trxId, &TxNotFound{},
	"cannot cancel trx_id=${tid}, there is no deferred transaction with that transaction id")//TODO

	trx := types.Transaction{}
	rlp.DecodeBytes(generatedTrx.PackedTrx, &trx)
	found := false
	for _, act := range trx.Actions{
		for _, auth := range act.Authorization {
			if auth == cancel.CancelingAuth{
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
	providedKeys []*ecc.PublicKey,
	providedPermission []*types.PermissionLevel,
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
		providedPermission,
		effectiveProvidedDelay,
		checkTime)
	permissionToSatisfy := make(map[types.PermissionLevel]common.Microseconds)

	for _, act := range actions {
		specialCase := false
		delay := effectiveProvidedDelay

		if act.Account == common.DefaultConfig.SystemAccountName {
			specialCase = true
			switch act.Name{
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
				if minPermissionName != common.PermissionName(0) {
					minPermission := a.GetPermission(&types.PermissionLevel{declaredAuth.Actor, minPermissionName}) //TODO
					//permissionIndex, err := a.db.GetIndex("id", entity.PermissionObject{})
					//if err != nil {
					//	log.Fatalln(err)
					//}
					EosAssert(a.GetPermission(&declaredAuth).Satisfies(*minPermission/*, permissionIndex*/), &IrrelevantAuthException{} ,
					"action declares irrelevant authority '${auth}'; minimum authority is ${min}" ) //TODO
				}
			}
			permissionToSatisfy[declaredAuth] = delay
			//TODO
		}
	}
	for p, q := range permissionToSatisfy {
		(*checkTime)()
		EosAssert(checker.SatisfiedLoc(&p, q, nil),  &UnsatisfiedAuthorization{},
		"transaction declares authority '${auth}', " +
		"but does not have signatures for it under a provided delay of ${provided_delay} ms, " +
		"provided permissions ${provided_permissions}, and provided keys ${provided_keys}") //TODO
	}
	if !allowUnusedKeys {
		EosAssert(checker.AllKeysUsed(), &TxIrrelevantSig{}, "transaction bears irrelevant signatures from these keys: ${keys}")
	}
}

func (a *AuthorizationManager) CheckAuthorization2(account common.AccountName,
	permission common.PermissionName,
	providedKeys []*ecc.PublicKey,
	providedPermission []*types.PermissionLevel,
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
		providedPermission,
		effectiveProvidedDelay,
		checkTime)
	EosAssert(checker.SatisfiedLc(&types.PermissionLevel{account, permission}, nil), &UnsatisfiedAuthorization{},
	"permission '${auth}' was not satisfied under a provided delay of ${provided_delay} ms, " +
	"provided permissions ${provided_permissions}, and provided keys ${provided_keys}") //TODO

	if !allowUnusedKeys {
		EosAssert(checker.AllKeysUsed(), &TxIrrelevantSig{}, "irrelevant keys provided: ${keys}") //TODO
	}
}

func (a *AuthorizationManager) GetRequiredKeys(trx *types.Transaction,
	candidateKeys []*ecc.PublicKey,
	providedDelay common.Microseconds) []ecc.PublicKey {
	checker := types.MakeAuthChecker(func(p *types.PermissionLevel) types.SharedAuthority { return a.GetPermission(p).Auth },
		a.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
		candidateKeys,
		nil,
		providedDelay,
		noopCheckTime)
	return checker.GetUsedKeys()
}
