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
)

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
		Auth:        a.AuthToShared(auth),
	}
	a.db.Insert(&perm)
	return &perm
}

func (a *AuthorizationManager) ModifyPermission(permission *entity.PermissionObject, auth *types.Authority) {
	a.db.Modify(&permission, func(po *entity.PermissionObject) {
		po.Auth = a.AuthToShared(*auth)
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
	Try(func(){
		defer Return()
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
		defer Return()
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
		defer Return()
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
		defer Return()
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

func (a *AuthorizationManager) CheckUpdateauthAuthorization(update types.UpdateAuth, auths []types.PermissionLevel) {
	EosAssert(len(auths) == 1, &IrrelevantAuthException{}, "updateauth action should only have one declared authorization")
	auth := auths[0]
	EosAssert(auth.Actor == update.Account, &IrrelevantAuthException{}, "the owner of the affected permission needs to be the actor of the declared authorization")
	minPermission := a.FindPermission(&types.PermissionLevel{update.Account, update.Permission})
	if minPermission == nil {
		permission := a.GetPermission(&types.PermissionLevel{update.Account, update.Permission})
		minPermission = permission
	}
	EosAssert(a.GetPermission(&auth).Satisfies(*minPermission), &IrrelevantAuthException{}, "") //TODO
}

//func (am *AuthorizationManager) CheckDeleteauthAuthorization(del system.DeleteAuth, auths []types.PermissionLevel) {
//	if len(auths) != 1 {
//		fmt.Println("error")
//		return
//	}
//	auth := auths[0]
//	if auth.Actor != del.Account {
//		fmt.Println("error")
//		return
//	}
//	minPermission := am.GetPermission(&types.PermissionLevel{del.Account, del.Permission})
//	if am.GetPermission(&auth).Satisfies(minPermission) == false {
//		fmt.Println("error")
//		return
//	}
//}
//
//func (am *AuthorizationManager) CheckLinkauthAuthorization(link system.LinkAuth, auths []types.PermissionLevel) {
//	if len(auths) != 1 {
//		fmt.Println("error")
//		return
//	}
//	auth := auths[0]
//	if auth.Actor != link.Account {
//		fmt.Println("error")
//		return
//	}
//	//TODO
//	linkPermissionName := am.LookupMinimumPermission(link.Account, link.Code, link.Type)
//	if &linkPermissionName == nil {
//		return
//	}
//	//TODO
//}
//
//func (am *AuthorizationManager) CheckUnlinkauthAuthorization(unlink system.UnlinkAuth, auths []types.PermissionLevel) {
//	if len(auths) != 1 {
//		fmt.Println("error")
//		return
//	}
//	auth := auths[0]
//	if auth.Actor != unlink.Account {
//		fmt.Println("error")
//		return
//	}
//	//TODO
//}
//
//func (am *AuthorizationManager) CheckCanceldelayAuthorization(canceldelay system.CancelDelay, auths []types.PermissionLevel) {
//	if len(auths) != 1 {
//		fmt.Println("error")
//		return
//	}
//	auth := auths[0]
//	if am.GetPermission(&auth).Satisfies(am.GetPermission(&canceldelay.CancelingAuth)) == false {
//		fmt.Println("error")
//		return
//	}
//	//TODO
//}
//
func (am *AuthorizationManager) CheckAuthorization(actions []*types.Action,
	providedKeys []*ecc.PublicKey,
	providedPermission []*types.PermissionLevel,
	providedDelay common.Microseconds,
	checkTime *func(),
	allowUnusedKeys bool,
) {
	//delayMaxLimit := common.Seconds(int64(am.control.GetGlobalProperties().Configuration.MaxTrxDelay))
	//var effectiveProvidedDelay common.Microseconds
	//if providedDelay >= delayMaxLimit {
	//	effectiveProvidedDelay = common.MaxMicroseconds()
	//} else {
	//	effectiveProvidedDelay = providedDelay
	//}
	//checker := types.MakeAuthChecker(func(p *types.PermissionLevel) types.SharedAuthority { return am.GetPermission(p).Auth },
	//	am.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
	//	providedKeys,
	//	providedPermission,
	//	effectiveProvidedDelay,
	//	checkTime)
	//permissionToSatisfy := make(map[types.PermissionLevel]common.Microseconds)
	//for _, act := range actions {
	//	specialCase := false
	//	delay := effectiveProvidedDelay
	//
	//	if act.Account == common.DefaultConfig.SystemAccountName {
	//		specialCase = true
	//
	//		//TODO
	//	}
	//
	//	for _, declaredAuth := range act.Authorization {
	//		checkTime()
	//		if !specialCase {
	//			minPermissionName := am.LookupMinimumPermission(declaredAuth.Actor, act.Account, act.Name)
	//			if minPermissionName != common.PermissionName(0) {
	//				minPermission := am.GetPermission(&types.PermissionLevel{declaredAuth.Actor, minPermissionName})
	//				//EOS_ASSERT
	//				if !am.GetPermission(&declaredAuth).Satisfies(minPermission) {
	//					fmt.Println("error")
	//					return
	//				}
	//			}
	//		}
	//		permissionToSatisfy[declaredAuth] = delay
	//		//TODO
	//	}
	//}
	//for p, q := range permissionToSatisfy {
	//	checkTime()
	//	if !checker.SatisfiedLoc(&p, q, nil) {
	//		fmt.Println("error")
	//		return
	//	}
	//}
	//if !allowUnusedKeys {
	//	if !checker.AllKeysUsed() {
	//		fmt.Println("error")
	//		return
	//	}
	//}
}

func (am *AuthorizationManager) CheckAuthorization2(account common.AccountName,
	permission common.PermissionName,
	providedKeys []*ecc.PublicKey,
	providedPermission []*types.PermissionLevel,
	providedDelay common.Microseconds,
	checkTime *func(),
	allowUnusedKeys bool,
) {
	//delayMaxLimit := common.Seconds(int64(am.control.GetGlobalProperties().Configuration.MaxTrxDelay))
	//var effectiveProvidedDelay common.Microseconds
	//if providedDelay >= delayMaxLimit {
	//	effectiveProvidedDelay = common.MaxMicroseconds()
	//} else {
	//	effectiveProvidedDelay = providedDelay
	//}
	//checker := types.MakeAuthChecker(func(p *types.PermissionLevel) types.SharedAuthority { return am.GetPermission(p).Auth },
	//	am.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
	//	providedKeys,
	//	providedPermission,
	//	effectiveProvidedDelay,
	//	checkTime)
	////TODO
	//if !allowUnusedKeys {
	//	if !checker.AllKeysUsed() {
	//		fmt.Println("error")
	//		return
	//	}
	//}
}

func (am *AuthorizationManager) GetRequiredKeys(trx *types.Transaction,
	candidateKeys []*ecc.PublicKey,
	providedDelay common.Microseconds) []ecc.PublicKey {
	checker := types.AuthorityChecker{}
	//checker := types.MakeAuthChecker(func(p *types.PermissionLevel) types.SharedAuthority { return am.GetPermission(p).Auth },
	//	am.control.GetGlobalProperties().Configuration.MaxAuthorityDepth,
	//	candidateKeys,
	//	nil,
	//	providedDelay,
	//	func() {})
	return checker.GetUsedKeys()
}

func (am *AuthorizationManager) AuthToShared(auth types.Authority) types.SharedAuthority {
	return types.SharedAuthority{auth.Threshold, auth.Keys, auth.Accounts, auth.Waits}
}
