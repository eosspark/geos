package types

import (
	"github.com/eosspark/eos-go/common"
	"go/types"
	"time"
)

type PermissionToAuthorityFunc uint64
type AuthorityChecker struct {
	permissionToAuthority PermissionToAuthorityFunc
	CheckTime             types.Func
	ProvidedKeys          []common.PublicKeyType
	ProvidePermissions    []PermissionLevel
	UsedKeys              []bool
	ProvideDelay          time.Duration
	RecursionDepthLimit   uint16
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
	providedDelay time.Duration,
	checkTime types.Func) AuthorityChecker {
	//noopChecktime := func() {}
	return AuthorityChecker{permissionToAuthority: pta, RecursionDepthLimit: recursionDepthLimit,
		ProvidedKeys: providedKeys, ProvidePermissions: providedPermission,
		ProvideDelay: providedDelay, CheckTime: checkTime,
	}
}
