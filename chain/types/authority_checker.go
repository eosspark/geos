package types

import (
	"github.com/eosspark/container/sets/treeset"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
)

type PermissionToAuthorityFunc func(*PermissionLevel) SharedAuthority
type AuthorityChecker struct {
	permissionToAuthority PermissionToAuthorityFunc
	CheckTime             *func()
	ProvidedKeys          []ecc.PublicKey
	ProvidedPermissions   treeset.Set
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
	cachedPermissions := make(PermissionCacheType)
	if cachedPerms == nil {
		cachedPerms = ac.initializePermissionCache(&cachedPermissions)
	}
	Visitor := WeightTallyVisitor{ac, cachedPerms, 0, 0}
	return Visitor.Visit([]PermissionLevelWeight{{*permission, 1}}) > 0
}

func (ac *AuthorityChecker) SatisfiedAcd(authority *SharedAuthority, cachedPermissions *PermissionCacheType, depth uint16) bool {
	var metaPermission []interface{}
	metaPermission = append(metaPermission, authority.Waits)
	metaPermission = append(metaPermission, authority.Keys)
	metaPermission = append(metaPermission, authority.Accounts)
	visitor := WeightTallyVisitor{ac, cachedPermissions, depth, 0}
	for _, permission := range metaPermission {
		if visitor.Visit(permission) >= authority.Threshold {
			ac.KeyReverterCancel()
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

func (ac *AuthorityChecker) GetUsedKeys() treeset.Set {
	f := treeset.NewWith(ecc.ComparePubKey)
	for i, usedKey := range ac.UsedKeys {
		if usedKey == true {
			f.AddItem(&ac.ProvidedKeys[i])
		}
	}
	return *f
}

func (ac *AuthorityChecker) GetUnusedKeys() treeset.Set {
	f := treeset.NewWith(ecc.ComparePubKey)
	for i, usedKey := range ac.UsedKeys {
		if usedKey == false {
			f.AddItem(&ac.ProvidedKeys[i])
		}
	}
	return *f
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
	ac.ProvidedPermissions = *treeset.NewWith(ecc.ComparePubKey)
	itr := ac.ProvidedPermissions.Iterator()
	for itr.Next() {
		val := itr.Value()
		map[PermissionLevel]PermissionCacheStatus(*cachedPermission)[*(val.(*PermissionLevel))] = PermissionSatisfied
	}
	return cachedPermission
}

func (ac *AuthorityChecker) KeyReverterCancel() {
	for i := range ac.UsedKeys {
		ac.UsedKeys[i] = true
	}
}

type WeightTallyVisitor struct {
	Checker           *AuthorityChecker
	CachedPermissions *PermissionCacheType
	RecursionDepth    uint16
	TotalWeight       uint32
}

func (wtv *WeightTallyVisitor) Visit(permission interface{}) uint32 {
	switch v := permission.(type) {
	case []WaitWeight:
		for _, p := range v {
			wtv.VisitWaitWeight(p)
		}
		return wtv.TotalWeight
	case []KeyWeight:
		for _, p := range v {
			wtv.VisitKeyWeight(p)
		}
		return wtv.TotalWeight
	case []PermissionLevelWeight:
		for _, p := range v {
			wtv.VisitPermissionLevelWeight(p)
		}
		return wtv.TotalWeight
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
		if key.Compare(permission.Key) {
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
	if status == 0 {
		if wtv.RecursionDepth < wtv.Checker.RecursionDepthLimit {
			r := false
			propagateError := false
			isNotFound := false
			Try(func() {
				auth := wtv.Checker.permissionToAuthority(&permission.Permission)
				if auth.Threshold == 0 {
					return
				}
				propagateError = true
				map[PermissionLevel]PermissionCacheStatus(*wtv.CachedPermissions)[permission.Permission] = BeingEvaluated
				r = wtv.Checker.SatisfiedAcd(&auth, wtv.CachedPermissions, wtv.RecursionDepth+1)
			}).Catch(func(e *PermissionQueryException) {
				isNotFound = true
				if propagateError {
					EosThrow(e, "authority_check::VisitPermissionLevelWeight is error: %v", GetDetailMessage(e))
				} else {
					return
				}
			}).End()
			if isNotFound {
				return wtv.TotalWeight
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
	providedKeys *treeset.Set, //[]*ecc.PublicKey,
	providedPermission *treeset.Set, //[]*PermissionLevel,
	providedDelay common.Microseconds,
	checkTime *func()) AuthorityChecker {
	//noopChecktime := func() {}
	providedKeysArray := make([]ecc.PublicKey, 0)
	usedKeysArray := make([]bool, providedKeys.Size())
	itr := providedKeys.Iterator()
	for itr.Next() {
		providedKeysArray = append(providedKeysArray, itr.Value().(ecc.PublicKey))
	}
	/*for i := 0; i < len(providedKeys.Values()); i++ {
		data:=providedKeys.Values()
		element := data[i].(ecc.PublicKey)
		providedKeysArray[i] = element
	}*/
	return AuthorityChecker{permissionToAuthority: pta, RecursionDepthLimit: recursionDepthLimit,
		ProvidedKeys: providedKeysArray, ProvidedPermissions: *providedPermission,
		ProvidedDelay: providedDelay, UsedKeys: usedKeysArray, CheckTime: checkTime,
	}
}
