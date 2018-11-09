package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
)

type PermissionToAuthorityFunc func(*PermissionLevel) SharedAuthority
type AuthorityChecker struct {
	permissionToAuthority PermissionToAuthorityFunc
	CheckTime             *func()
	ProvidedKeys          []ecc.PublicKey
	ProvidedPermissions   common.FlatSet
	UsedKeys              []bool
	ProvidedDelay         common.Microseconds
	RecursionDepthLimit   uint16
	Visitor               WeightTallyVisitor
}

func (ac *AuthorityChecker) SatisfiedLoc(permission *PermissionLevel,
	overrideProvidedDelay common.Microseconds,
	cachedPerms *PermissionCacheType) bool {
	ac.ProvidedDelay = overrideProvidedDelay
	return ac.SatisfiedLc(permission, cachedPerms)
}

func (ac *AuthorityChecker) SatisfiedLc(permission *PermissionLevel, cachedPerms *PermissionCacheType) bool {
	var cachedPermissions PermissionCacheType
	if cachedPerms == nil {
		cachedPerms = ac.initializePermissionCache(&cachedPermissions)
	}
	Visitor := WeightTallyVisitor{ac, cachedPerms, 0, 0}
	return Visitor.Visit(PermissionLevelWeight{*permission, 1}) > 0
}

func (ac *AuthorityChecker) SatisfiedAcd(authority *SharedAuthority, cachedPermissions *PermissionCacheType, depth uint16) bool {
	//TODO: lack of canceled
	var metaPermission []interface{}
	metaPermission = append(metaPermission, authority.Waits)
	metaPermission = append(metaPermission, authority.Keys)
	metaPermission = append(metaPermission, authority.Accounts)
	visitor := WeightTallyVisitor{ac, cachedPermissions, depth, 0}
	for _, permission := range metaPermission {
		if visitor.Visit(permission) >= authority.Threshold {
			return true
		}
	}
	return false
}

func (ac *AuthorityChecker) AllKeysUsed() bool {
	for _, usedKey := range ac.UsedKeys {
		if usedKey == false {
			return false
		}
	}
	return true
}

func (ac *AuthorityChecker) GetUsedKeys() common.FlatSet {
	f := common.FlatSet{}
	for i, usedKey := range ac.UsedKeys {
		if usedKey == true {
			f.Insert(&ac.ProvidedKeys[i])
		}
	}
	return f
}

func (ac *AuthorityChecker) GetUnusedKeys() common.FlatSet {
	f := common.FlatSet{}
	for i, usedKey := range ac.UsedKeys {
		if usedKey == false {
			f.Insert(&ac.ProvidedKeys[i])
		}
	}
	return f
}

type PermissionCacheStatus uint64

const (
	_ PermissionCacheStatus = iota
	BeingEvaluated
	PermissionUnsatisfied
	PermissionSatisfied
)

type PermissionCacheType map[PermissionLevel]PermissionCacheStatus

func (ac *AuthorityChecker) PermissionStatusInCache(permissions PermissionCacheType, level *PermissionLevel) PermissionCacheStatus {
	itr, ok := map[PermissionLevel]PermissionCacheStatus(permissions)[*level]
	if ok {
		return itr
	}
	itr2, ok := map[PermissionLevel]PermissionCacheStatus(permissions)[PermissionLevel{level.Actor, common.PermissionName(common.N(""))}]
	if ok {
		return itr2
	}
	return 0
}

func (ac *AuthorityChecker) initializePermissionCache(cachedPermission *PermissionCacheType) *PermissionCacheType {
	for i := 0; i < ac.ProvidedPermissions.Len(); i++{
		map[PermissionLevel]PermissionCacheStatus(*cachedPermission)[*(ac.ProvidedPermissions.Data[i].(*PermissionLevel))] = PermissionSatisfied
	}
	return cachedPermission
}

type WeightTallyVisitor struct {
	Checker           *AuthorityChecker
	CachedPermissions *PermissionCacheType
	RecursionDepth    uint16
	TotalWeight       uint32
}

func (wtv *WeightTallyVisitor) Visit(permission interface{}) uint32 {
	switch v := permission.(type) {
	case WaitWeight:
		return wtv.VisitWaitWeight(v)
	case KeyWeight:
		return wtv.VisitKeyWeight(v)
	case PermissionLevelWeight:
		return wtv.VisitPermissionLevelWeight(v)
	default:
		return 0
	}
}

func (wtv *WeightTallyVisitor) VisitWaitWeight(permission WaitWeight) uint32 {
	if wtv.Checker.ProvidedDelay >= common.Seconds(int64(permission.WaitSec)) {
		wtv.TotalWeight += uint32(permission.Weight)
	}
	return wtv.TotalWeight
}

func (wtv *WeightTallyVisitor) VisitKeyWeight(permission KeyWeight) uint32 {
	var itr int
	for _, key := range wtv.Checker.ProvidedKeys {
		if key.Content == permission.Key.Content && key.Curve == permission.Key.Curve {
			wtv.Checker.UsedKeys[itr] = true
			wtv.TotalWeight += uint32(permission.Weight)
			break
		}
		itr++
	}
	return wtv.TotalWeight
}

func (wtv *WeightTallyVisitor) VisitPermissionLevelWeight(permission PermissionLevelWeight) uint32 {
	status := wtv.Checker.PermissionStatusInCache(*wtv.CachedPermissions, &permission.Permission)
	if status != 0 {
		if wtv.RecursionDepth < wtv.Checker.RecursionDepthLimit {
			r := false
			propagateError := false
			//try_catch
			auth := wtv.Checker.permissionToAuthority(&permission.Permission)
			propagateError = true
			map[PermissionLevel]PermissionCacheStatus(*wtv.CachedPermissions)[permission.Permission] = BeingEvaluated
			r = wtv.Checker.SatisfiedAcd(&auth, wtv.CachedPermissions, wtv.RecursionDepth+1)
			if propagateError {
			} else {

			}

			if r {
				wtv.TotalWeight += uint32(permission.Weight)
				map[PermissionLevel]PermissionCacheStatus(*wtv.CachedPermissions)[permission.Permission] = PermissionSatisfied
			} else {
				map[PermissionLevel]PermissionCacheStatus(*wtv.CachedPermissions)[permission.Permission] = PermissionUnsatisfied
			}
		}
	} else if status == PermissionSatisfied {
		wtv.TotalWeight += uint32(permission.Weight)
	}
	return wtv.TotalWeight
}

func MakeAuthChecker(pta PermissionToAuthorityFunc,
	recursionDepthLimit uint16,
	providedKeys *common.FlatSet, //[]*ecc.PublicKey,
	providedPermission *common.FlatSet, //[]*PermissionLevel,
	providedDelay common.Microseconds,
	checkTime *func()) AuthorityChecker {
	//noopChecktime := func() {}
	providedKeysArray := make([]ecc.PublicKey, providedKeys.Len())
	for i := 0; i < providedKeys.Len(); i++{
		element := providedKeys.Data[i].(*ecc.PublicKey)
		providedKeysArray[i] = *element
	}
	return AuthorityChecker{permissionToAuthority: pta, RecursionDepthLimit: recursionDepthLimit,
		ProvidedKeys: providedKeysArray, ProvidedPermissions: *providedPermission,
		ProvidedDelay: providedDelay, CheckTime: checkTime,
	}
}
