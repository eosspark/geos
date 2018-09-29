package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/contracts/system"
	"github.com/eosspark/eos-go/db"
	"time"
)

var IsActiveAz bool

var azInstance *AuthorizationManager

type AuthorizationManager struct {
	control *Controller
	db      *eosiodb.DataBase
}

func GetAuthorizationManager() *AuthorizationManager {
	if !IsActiveAz {
		azInstance = newAuthorizationManager()
	}
	return azInstance
}

func newAuthorizationManager() *AuthorizationManager {
	control := GetControllerInstance()
	db := control.DataBase()
	return &AuthorizationManager{control: control, db: db}
}

type PermissionIdType types.IdType

//func (am *AuthorizationManager) AddIndices() {
//	am.db.Insert(&types.PermissionObject{})
//	am.db.Insert(&types.PermissionUsageObject{})
//	am.db.Insert(&types.PermissionLinkObject{})
//}

func (am *AuthorizationManager) InitializeDataBase() {

}

func (am *AuthorizationManager) CreatePermission(account common.AccountName,
												 name common.PermissionName,
												 parent PermissionIdType,
												 auth types.Authority,
												 initialCreationTime common.TimePoint,
											    ) types.PermissionObject {
	creationTime := initialCreationTime
	if creationTime == 1 {
		creationTime = am.control.PendingBlockTime()
	}

	permUsage := types.PermissionUsageObject{}
	permUsage.LastUsed = creationTime
	am.db.Insert(&permUsage)

	perm := types.PermissionObject{
		UsageId:     permUsage.ID,
		Parent:      types.IdType(parent),
		Owner:       account,
		Name:        name,
		LastUpdated: creationTime,
		Auth:        am.AuthToShared(auth),
	}
	am.db.Insert(&perm)
	return perm
}

func (am *AuthorizationManager) ModifyPermission(permission *types.PermissionObject, auth *types.Authority) {
	am.db.Update(&permission, func(data interface{}) error {
		permission.Auth = am.AuthToShared(*auth)
		permission.LastUpdated = am.control.PendingBlockTime()
		return nil
	})
}

func (am *AuthorizationManager) RemovePermission(permission *types.PermissionObject) {

}

func (am *AuthorizationManager) UpdatePermissionUsage(permission *types.PermissionObject) {
	puo := types.PermissionUsageObject{}
	am.db.Find("ID", permission.UsageId, &puo)
	am.db.Update(&puo, func(data interface{}) error {
		puo.LastUsed = am.control.PendingBlockTime()
		return nil
	})
}

func (am *AuthorizationManager) GetPermissionLastUsed(permission *types.Permission) common.TimePoint {
	puo := types.PermissionUsageObject{}

	return puo.LastUsed
}

func (am *AuthorizationManager) FindPermission(level types.PermissionLevel) *types.PermissionObject {
	po := types.PermissionObject{}
	am.db.Find("", 0, &po)
	return &po
}

func (am *AuthorizationManager) GetPermission(level types.PermissionLevel) types.PermissionObject {
	po := types.PermissionObject{}
	am.db.Find("", 0, &po)
	return po
}

func (am *AuthorizationManager) LookupLinkedPermission(authorizerAccount common.AccountName,
	scope common.AccountName,
	actName common.ActionName,
) common.PermissionName {
	var pn common.PermissionName
	return pn
}

func (am *AuthorizationManager) LookupMinimumPermission(authorizerAccount common.AccountName,
	scope common.AccountName,
	actName common.ActionName,
) common.PermissionName {
	var pn common.PermissionName
	return pn
}

func (am *AuthorizationManager) CheckUpdateauthAuthorization(update system.UpdateAuth, auths []types.PermissionLevel) {
	if len(auths) != 1 {
		fmt.Println("error")
		return
	}
	auth := auths[0]
	if auth.Actor != update.Account {
		fmt.Println("error")
		return
	}
	minPermission := am.FindPermission(types.PermissionLevel{update.Account, update.Permission})
	if minPermission == nil {
		permission := am.GetPermission(types.PermissionLevel{update.Account, update.Permission})
		minPermission = &permission
	}
	if am.GetPermission(auth).Satisfies(*minPermission) == false {
		fmt.Println("error")
		return
	}
}

func (am *AuthorizationManager) CheckDeleteauthAuthorization(del system.DeleteAuth, auths []types.PermissionLevel) {
	if len(auths) != 1 {
		fmt.Println("error")
		return
	}
	auth := auths[0]
	if auth.Actor != del.Account {
		fmt.Println("error")
		return
	}
	minPermission := am.GetPermission(types.PermissionLevel{del.Account, del.Permission})
	if am.GetPermission(auth).Satisfies(minPermission) == false {
		fmt.Println("error")
		return
	}
}

func (am *AuthorizationManager) CheckLinkauthAuthorization(link system.LinkAuth, auths []types.PermissionLevel) {
	if len(auths) != 1 {
		fmt.Println("error")
		return
	}
	auth := auths[0]
	if auth.Actor != link.Account {
		fmt.Println("error")
		return
	}
	//待完善
	linkPermissionName := am.LookupMinimumPermission(link.Account, link.Code, link.Type)
	if &linkPermissionName == nil {
		return
	}
	//待完善
}

func (am *AuthorizationManager) CheckUnlinkauthAuthorization(unlink system.UnlinkAuth, auths []types.PermissionLevel) {
	if len(auths) != 1 {
		fmt.Println("error")
		return
	}
	auth := auths[0]
	if auth.Actor != unlink.Account {
		fmt.Println("error")
		return
	}
	//待完善
}

func (am *AuthorizationManager) CheckCanceldelayAuthorization(canceldelay system.CancelDelay, auths []types.PermissionLevel) {
	if len(auths) != 1 {
		fmt.Println("error")
		return
	}
	auth := auths[0]
	if am.GetPermission(auth).Satisfies(am.GetPermission(canceldelay.CancelingAuth)) == false {
		fmt.Println("error")
		return
	}
	//待完善
}

func (am *AuthorizationManager) CheckAuthorization(actions []types.Action,
	providedKeys []common.PublicKeyType,
	providedPermission []types.PermissionLevel,
	providedDelay time.Time,
	allowUnusedKeys bool,
) {
	//delayMaxLimit := am.control

}

func (am *AuthorizationManager) GetRequiredKeys(trx types.Transaction,
	candidateKeys []common.PublicKeyType,
	providedDelay time.Time) {
	//check := MakeAuthChecker()
}

func (am *AuthorizationManager) AuthToShared(auth types.Authority) types.SharedAuthority{
	return types.SharedAuthority{auth.Threshold,auth.Keys,auth.Accounts,auth.Waits}
}