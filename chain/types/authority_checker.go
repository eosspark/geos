package types

import (
	"github.com/eosspark/eos-go/common"
)

type PermissionToAuthorityFunc interface {}
type AuthorityChecker struct {
	permissionToAuthority PermissionToAuthorityFunc
	CheckTime             func()
	ProvidedKeys          []common.PublicKeyType
	ProvidePermissions    []PermissionLevel
	UsedKeys              []bool
	ProvideDelay          common.Microseconds
	RecursionDepthLimit   uint16
}

func (ac *AuthorityChecker) SatisfiedLoc( permission *PermissionLevel,
										  overrideProvidedDelay common.Microseconds,
	    								  cacheType PermissionCacheType) bool {
	return true
}

func (ac *AuthorityChecker) SatisfiedLc(permission *PermissionLevel, cacheType PermissionCacheType) bool {
 	return true
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

type WeightTallyVisitor struct {
	Checker          AuthorityChecker
	CachePermissions PermissionCacheType
	RecursionDepth   uint16
	TotalWeight      uint32
}

func MakeAuthChecker(pta PermissionToAuthorityFunc,
	recursionDepthLimit uint16,
	providedKeys []common.PublicKeyType,
	providedPermission []PermissionLevel,
	providedDelay common.Microseconds,
	checkTime func()) AuthorityChecker {
	//noopChecktime := func() {}
	return AuthorityChecker{permissionToAuthority: pta, RecursionDepthLimit: recursionDepthLimit,
		ProvidedKeys: providedKeys, ProvidePermissions: providedPermission,
		ProvideDelay: providedDelay, CheckTime: checkTime,
	}
}
