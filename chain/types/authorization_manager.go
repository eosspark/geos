package types

import (
	"github.com/eosspark/eos-go/db"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/contracts/system"
	"time"
	"fmt"
)

type AuthorizationManager struct {
	db *eosiodb.Session
}

type PermissionIdType uint64

func (am *AuthorizationManager) AddIndices() {
	am.db.Insert(&PermissionObject{})
	am.db.Insert(&PermissionUsageObject{})
	am.db.Insert(&PermissionLinkObject{})
}

func (am *AuthorizationManager) InitializeDataBase() {

}

func (am *AuthorizationManager) CreatePermission( account common.AccountName,
	                                              name    common.PermissionName,
		                                          parent  PermissionIdType,
		                                          auth    common.Authority,
		                                          initialCreationTime time.Duration,
		                                         ) PermissionObject {
	creationTime := initialCreationTime
	if creationTime == 1 {
		//createTime = pendingBlockTime
	}

	var permUsage PermissionUsageObject
	permUsage.LastUsed = creationTime
	am.db.Insert(&permUsage)

	perm := PermissionObject{
		UsageId:     permUsage.Id,
		Parent:      uint64(parent),
		Owner:       account,
		Name:        name,
		LastUpdated: creationTime,
		//Auth:        SharedAuthority(),
	}
	am.db.Insert(&perm)
	return perm
}

func (am *AuthorizationManager) ModifyPermission (permission PermissionObject, auth common.Authority) {
	am.db.Update( &permission, func(data interface{}) error {
		//permission.Auth = auth
		//permission.LastUpdated = pendingBlockTime
		return nil
	})
}

func (am *AuthorizationManager) RemovePermission (){

}

func (am *AuthorizationManager) UpdatePermissionUsage (){
	var puo PermissionUsageObject
	am.db.Update(&puo, func(data interface{}) error {
		//puo.LastUsed = pendingBlockTime
		return nil
	})
}

func (am *AuthorizationManager) GetPermissionLastUsed (permission common.Permission) time.Duration {
	var puo PermissionUsageObject
	return puo.LastUsed
}

func (am *AuthorizationManager) FindPermission (level common.PermissionLevel) PermissionObject {
	var po PermissionObject
	am.db.Find("", 0, &po)
	return po
}

func (am *AuthorizationManager) GetPermission (level common.PermissionLevel) PermissionObject {
	var po PermissionObject
	am.db.Find("", 0, &po)
	return po
}

func (am *AuthorizationManager) LookupLinkedPermission (authorizerAccount common.AccountName,
	                                                    scope             common.AccountName,
	                                                    actName           common.ActionName,
														) common.PermissionName {
	var pn common.PermissionName
	return pn
}

func (am *AuthorizationManager) LookupMinimumPermission (authorizerAccount common.AccountName,
														 scope             common.AccountName,
														 actName           common.ActionName,
													     ) common.PermissionName {
	var pn common.PermissionName
	return pn
}

func (am *AuthorizationManager) CheckUpdateauthAuthorization (update system.UpdateAuth, auths []common.PermissionLevel) {
	if len(auths) != 1 {
		fmt.Println("error")
		return
	}
	auth := auths[0]
	if auth.Actor != update.Account {
		fmt.Println("error")
		return
	}
	minPermission := am.FindPermission(common.PermissionLevel{update.Account, update.Permission})
	
}

//func (am *AuthorizationManager) CheckDeleteauthAuthorization ( ) {
//
//}
//
//func (am *AuthorizationManager) CheckLinkauthAuthorization ( ) {
//
//}
//
//func (am *AuthorizationManager) CheckUnlickauthAuthorization ( ) {
//
//}
//
//func (am *AuthorizationManager) CheckCanceldelayAuthorization ( ) {
//
//}
//
//func (am *AuthorizationManager) CheckAuthorization ( ) {
//
//}
//
//func (am *AuthorizationManager) GetRequiredKeys ( ) {
//	check := MakeAuthChecker()
//}