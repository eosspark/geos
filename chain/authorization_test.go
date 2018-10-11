package chain

import (
	"testing"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/chain/types"
)

func TestAuthorizationManager_CreatePermission(t *testing.T) {
	az := GetAuthorizationManager()

	az.CreatePermission(common.AccountName(common.N("yc")),
						common.PermissionName(common.N("active")),
						PermissionIdType(1),
						types.Authority{},2)

	//Need control.pending
	//az.CreatePermission(common.AccountName(common.N("yc")),
	//					common.PermissionName(common.N("active")),
	//					PermissionIdType(1),
	//					types.Authority{},1)

}

func TestAuthorizationManager_ModifyPermission(t *testing.T) {
	//am := GetAuthorizationManager()
	////am.CreatePermission(common.AccountName(common.N("yc")),
	////					common.PermissionName(common.N("active")),
	////					PermissionIdType(1),
	////					types.Authority{},2)
	//permUsage := types.PermissionUsageObject{}
	//permUsage.LastUsed = 2
	//perm := types.PermissionObject{
	//	UsageId:     permUsage.ID,
	//	Parent:      types.IdType(1),
	//	Owner:       common.AccountName(common.N("yc")),
	//	Name:        common.PermissionName(common.N("active")),
	//	LastUpdated: 2,
	//	Auth:        am.AuthToShared(types.Authority{}),
	//}
	//fmt.Println(perm)
	//am.db.Insert(&perm)
	//po := am.FindPermission(&types.PermissionLevel{common.AccountName(common.N("yc")),common.PermissionName(common.N("active"))})
	//fmt.Println(*po)
	//am.ModifyPermission(po, &types.Authority{2,[]types.KeyWeight{},[]types.PermissionLevelWeight{},[]types.WaitWeight{}})
}
