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


type MetaPermission []interface{}

func (m MetaPermission) Len() int { return len(m) }

func (m MetaPermission) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m MetaPermission) Less(i, j int) bool {
	iType, jType := 0, 0
	iWeight, jWeight := uint16(0), uint16(0)
	switch v := m[i].(type) {
	case WaitWeight:
		iWeight = uint16(v.Weight)
		iType = 1
	case KeyWeight:
		iWeight = uint16(v.Weight)
		iType = 2
	case PermissionLevelWeight:
		iWeight = uint16(v.Weight)
		iType = 3
	}
	switch v := m[j].(type) {
	case WaitWeight:
		jWeight = uint16(v.Weight)
		iType = 1
	case KeyWeight:
		jWeight = uint16(v.Weight)
		iType = 2
	case PermissionLevelWeight:
		jWeight = uint16(v.Weight)
		iType = 3
	}
	if iWeight < jWeight {
		return true
	} else if iWeight > jWeight{
		return false
	} else {
		if iType < jType {
			return true
		} else {
			return false
		}
	}
}

func (m MetaPermission) Sort() {
	for i := 0; i < m.Len() - 1; i++ {
		for j := 0; j < m.Len() - 1 - i; j++ {
			if m.Less(j, j + 1) {
				m.Swap(j, j + 1)
			}
		}
	}
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
	return Visitor.Visit(PermissionLevelWeight{*permission, 1}) > 0
}

func (ac *AuthorityChecker) SatisfiedAcd(authority *SharedAuthority, cachedPerms *PermissionCacheType, depth uint16) bool {
	cachedPermissions := make(PermissionCacheType)
	if cachedPerms == nil {
		cachedPerms = ac.initializePermissionCache(&cachedPermissions)
	}
	var metaPermission MetaPermission
	satisfied := false
	oldUsedKeys := make([]bool, len(ac.UsedKeys))
	copy(oldUsedKeys, ac.UsedKeys)
	defer func() {
		if !satisfied {
			ac.UsedKeys = oldUsedKeys
		}
	}()
	for _, wait := range authority.Waits {
		metaPermission = append(metaPermission, wait)
	}
	for _, key := range authority.Keys {
		metaPermission = append(metaPermission, key)
	}
	for _, account := range authority.Accounts {
		metaPermission = append(metaPermission, account)
	}
	metaPermission.Sort()
	visitor := WeightTallyVisitor{ac, cachedPerms, depth, 0}
	for _, permission := range metaPermission {
		if visitor.Visit(permission) >= authority.Threshold {
			satisfied = true
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
	f := treeset.NewWith(ecc.TypePubKey, ecc.ComparePubKey)
	for i, usedKey := range ac.UsedKeys {
		if usedKey == true {
			f.AddItem(ac.ProvidedKeys[i])
		}
	}
	return *f
}

func (ac *AuthorityChecker) GetUnusedKeys() treeset.Set {
	f := treeset.NewWith(ecc.TypePubKey, ecc.ComparePubKey)
	for i, usedKey := range ac.UsedKeys {
		if usedKey == false {
			f.AddItem(ac.ProvidedKeys[i])
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

type PermissionCacheType = map[PermissionLevel]PermissionCacheStatus

func (ac *AuthorityChecker) PermissionStatusInCache(permissions PermissionCacheType, level *PermissionLevel) PermissionCacheStatus {
	itr, ok := permissions[*level]
	if ok {
		return itr
	}
	itr2, ok := permissions[PermissionLevel{level.Actor, common.PermissionName(common.N(""))}]
	if ok {
		return itr2
	}
	return 0
}

func (ac *AuthorityChecker) initializePermissionCache(cachedPermission *PermissionCacheType) *PermissionCacheType {
	itr := ac.ProvidedPermissions.Iterator()
	for itr.Next() {
		val := itr.Value()
		(*cachedPermission)[(val.(PermissionLevel))] = PermissionSatisfied
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
		wtv.VisitWaitWeight(v)
		return wtv.TotalWeight
	case KeyWeight:
		wtv.VisitKeyWeight(v)
		return wtv.TotalWeight
	case PermissionLevelWeight:
		wtv.VisitPermissionLevelWeight(v)
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
				(*wtv.CachedPermissions)[permission.Permission] = BeingEvaluated
				r = wtv.Checker.SatisfiedAcd(&auth, wtv.CachedPermissions, wtv.RecursionDepth+1)
			}).Catch(func(e *PermissionQueryException) {
				isNotFound = true
				if propagateError {
					EosThrow(e, "authority_check::VisitPermissionLevelWeight is error: %v", e.DetailMessage())
				} else {
					return
				}
			}).End()
			if isNotFound {
				return wtv.TotalWeight
			}
			if r {
				wtv.TotalWeight += uint32(permission.Weight)
				(*wtv.CachedPermissions)[permission.Permission] = PermissionSatisfied
			} else {
				(*wtv.CachedPermissions)[permission.Permission] = PermissionUnsatisfied
			}
		}
	} else if status == PermissionSatisfied {
		wtv.TotalWeight += uint32(permission.Weight)
	}
	return wtv.TotalWeight
}

func MakeAuthChecker(pta PermissionToAuthorityFunc,
	recursionDepthLimit uint16,
	providedKeys *treeset.Set,
	providedPermission *treeset.Set,
	providedDelay common.Microseconds,
	checkTime *func()) AuthorityChecker {
	providedKeysArray := make([]ecc.PublicKey, 0)
	usedKeysArray := make([]bool, providedKeys.Size())
	itr := providedKeys.Iterator()
	for itr.Next() {
		providedKeysArray = append(providedKeysArray, itr.Value().(ecc.PublicKey))
	}
	return AuthorityChecker{permissionToAuthority: pta, RecursionDepthLimit: recursionDepthLimit,
		ProvidedKeys: providedKeysArray, ProvidedPermissions: *providedPermission,
		ProvidedDelay: providedDelay, UsedKeys: usedKeysArray, CheckTime: checkTime,
	}
}
