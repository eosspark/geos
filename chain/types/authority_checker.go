package types

import (
	"github.com/eosspark/eos-go/common"

)

type PermissionToAuthorityFunc interface {}
type AuthorityChecker struct {
	permissionToAuthority PermissionToAuthorityFunc
	CheckTime             func()
	ProvidedKeys          []common.PublicKeyType
	ProvidedPermissions   []PermissionLevel
	UsedKeys              []bool
	ProvidedDelay         common.Microseconds
	RecursionDepthLimit   uint16
	Visitor               WeightTallyVisitor
}

func (ac *AuthorityChecker) SatisfiedLoc( permission *PermissionLevel,
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

func (ac *AuthorityChecker) AllKeysUsed() bool {
	for _, usedKey := range ac.UsedKeys {
		if usedKey == false {
			return false
		}
	}
	return true
}

func (ac *AuthorityChecker) GetUsedKeys() []common.PublicKeyType {
	return nil
}

type PermissionCacheStatus uint64

const (
	_ PermissionCacheStatus = iota
	BeingEvaluated
	PermissionUnsatisfied
	PermissionSatisfied
)

type PermissionCacheType map[PermissionLevel]PermissionCacheStatus

func (ac *AuthorityChecker) initializePermissionCache( cachedPermission *PermissionCacheType ) *PermissionCacheType {
	//for _, p := range ac.ProvidedPermissions {
	//
	//}
	return cachedPermission
}


type WeightTallyVisitor struct {
	Checker          *AuthorityChecker
	CachePermissions *PermissionCacheType
	RecursionDepth   uint16
	TotalWeight      uint32
}

func (wtv *WeightTallyVisitor) Visit (permission interface{}) uint32 {
	switch v := permission.(type) {
	case WaitWeight :
		return wtv.VisitWaitWeight(v)
	case KeyWeight :
		return wtv.VisitKeyWeight(v)
	case PermissionLevelWeight :
		return wtv.VisitPermissionLevelWeight(v)
	default:
		return 0
	}
}

func (wtv *WeightTallyVisitor) VisitWaitWeight (permission WaitWeight) uint32 {
	if wtv.Checker.ProvidedDelay >= common.Seconds(int64(permission.WaitSec)) {
		wtv.TotalWeight += uint32(permission.Weight)
	}
	return wtv.TotalWeight
}

func (wtv *WeightTallyVisitor) VisitKeyWeight (permission KeyWeight) uint32 {
	return 1
}

func (wtv *WeightTallyVisitor) VisitPermissionLevelWeight (permission PermissionLevelWeight) uint32 {
	return 1
}

func MakeAuthChecker(pta PermissionToAuthorityFunc,
	recursionDepthLimit uint16,
	providedKeys       []common.PublicKeyType,
	providedPermission []PermissionLevel,
	providedDelay      common.Microseconds,
	checkTime          func()) AuthorityChecker {
	//noopChecktime := func() {}
	return AuthorityChecker{permissionToAuthority: pta, RecursionDepthLimit: recursionDepthLimit,
		ProvidedKeys: providedKeys, ProvidedPermissions: providedPermission,
		ProvidedDelay: providedDelay, CheckTime: checkTime,
	}
}
