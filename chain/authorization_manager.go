package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/database"
	"github.com/eosspark/eos-go/entity"
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
	//puo := entity.PermissionUsageObject{}
	//am.db.Find("ID", permission.UsageId, &puo)
	//return puo.LastUsed
	return 0
}

func (am *AuthorizationManager) FindPermission(level *types.PermissionLevel) *types.PermissionObject {
	//po := types.PermissionObject{}
	//am.db.Find("ByOwner", common.Tuple{level.Actor, level.Permission}, &po)
	//poo := []types.PermissionObject{}
	//am.db.All(&poo)
	//fmt.Println(poo)
	//return &po
	return &types.PermissionObject{}
}

func (am *AuthorizationManager) GetPermission(level *types.PermissionLevel) *entity.PermissionObject {
	po := entity.PermissionObject{}
	//am.db.Find("ByOwner", common.Tuple{level.Actor, level.Permission}, &po)
	return &po
}

//
//func (am *AuthorizationManager) LookupLinkedPermission(authorizerAccount common.AccountName,
//	scope common.AccountName,
//	actName common.ActionName,
//) common.PermissionName {
//
//	key := common.MakeTuple(authorizerAccount, scope, actName)
//	link := types.PermissionLinkObject{}
//	err := am.db.Find("ByActionName", key, &link)
//	if err != nil {
//		key = common.MakeTuple(authorizerAccount, scope, common.AccountName(common.N("")))
//		err = am.db.Find("ByActionName", key, &link)
//	}
//	if err == nil {
//		return link.RequiredPermission
//	}
//	var pn common.PermissionName
//	return pn
//}
//
//func (am *AuthorizationManager) LookupMinimumPermission(authorizerAccount common.AccountName,
//	scope common.AccountName,
//	actName common.ActionName,
//) common.PermissionName {
//	if scope == common.DefaultConfig.SystemAccountName {
//		//EOS_ASSERT
//	}
//	// if !linkPermission
//	linkedPermission := am.LookupLinkedPermission(authorizerAccount, scope, actName)
//	if linkedPermission == common.PermissionName(common.DefaultConfig.EosioAnyName) {
//		var pn common.PermissionName
//		return pn
//	}
//	return linkedPermission
//}
//
//func (am *AuthorizationManager) CheckUpdateauthAuthorization(update system.UpdateAuth, auths []types.PermissionLevel) {
//	if len(auths) != 1 {
//		fmt.Println("error")
//		return
//	}
//	auth := auths[0]
//	if auth.Actor != update.Account {
//		fmt.Println("error")
//		return
//	}
//	minPermission := am.FindPermission(&types.PermissionLevel{update.Account, update.Permission})
//	if minPermission == nil {
//		permission := am.GetPermission(&types.PermissionLevel{update.Account, update.Permission})
//		minPermission = &permission
//	}
//	if am.GetPermission(&auth).Satisfies(*minPermission) == false {
//		fmt.Println("error")
//		return
//	}
//}
//
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
