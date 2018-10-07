package types

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"strings"
)

type WeightType uint16

type Permission struct {
	PermName     string    `json:"perm_name"`
	Parent       string    `json:"parent"`
	RequiredAuth Authority `json:"required_auth"`
}

type PermissionLevel struct {
	Actor      common.AccountName    `json:"actor"`
	Permission common.PermissionName `json:"permission"`
}

type PermissionLevelWeight struct {
	Permission PermissionLevel `json:"permission"`
	Weight     WeightType      `json:"weight"` // weight_type
}

type KeyWeight struct {
	// Key    common.PublicKeyType `json:"key"`
	Key    ecc.PublicKey `json:"key"`
	Weight WeightType    `json:"weight"`
}

type WaitWeight struct {
	WaitSec uint32     `json:"wait_sec"`
	Weight  WeightType `json:"weight"`
}

type Authority struct {
	Threshold uint32                  `json:"threshold"`
	Keys      []KeyWeight             `json:"keys,omitempty"`
	Accounts  []PermissionLevelWeight `json:"accounts,omitempty"`
	Waits     []WaitWeight            `json:"waits,omitempty"`
}

type SharedAuthority struct {
	Threshold uint32
	Keys      []KeyWeight             `json:"keys"`
	Accounts  []PermissionLevelWeight `json:"accounts"`
	Waits     []WaitWeight            `json:"waits"`
}

// NewPermissionLevel parses strings like `account@active`,
// `otheraccount@owner` and builds a PermissionLevel struct. It
// validates that there is a single optional @ (where permission
// defaults to 'active'), and validates length of account and
// permission names.
func NewPermissionLevel(in string) (out PermissionLevel, err error) {
	parts := strings.Split(in, "@")
	if len(parts) > 2 {
		return out, fmt.Errorf("permission %q invalid, use account[@permission]", in)
	}

	if len(parts[0]) > 12 {
		return out, fmt.Errorf("account name %q too long", parts[0])
	}

	out.Actor = common.AccountName(common.StringToName(parts[0]))
	out.Permission = common.PermissionName(common.StringToName("active"))

	if len(parts) == 2 {
		if len(parts[1]) > 12 {
			return out, fmt.Errorf("permission %q name too long", parts[1])
		}

		out.Permission = common.PermissionName(common.StringToName("active"))
	}

	return
}

func (auth Authority) Equals(author Authority) bool {
	return true
}

func (sharedAuth SharedAuthority) Equals(sharedAuthor SharedAuthority) bool {
	return true
}

func (sharedAuth SharedAuthority) GetBillableSize() uint64 { //返回值类型不确定
	accountSize := uint64(len(sharedAuth.Accounts)) * common.BillableSizeV("permission_level_weight")
	waitsSize := uint64(len(sharedAuth.Waits)) * common.BillableSizeV("wait_weight")
	keysSize := uint64(0)
	for _, key := range sharedAuth.Keys {
		keysSize += common.BillableSizeV("key_weight")
		keysSize += uint64(key.Weight) * 0 //待修改
	}
	return accountSize + waitsSize + keysSize
}

func Validate(auth Authority) bool {
	var totalWeight uint32 = 0
	if len(auth.Accounts)+len(auth.Keys)+len(auth.Waits) > 1<<16 {
		return false
	}
	if auth.Threshold == 0 {
		return false
	}
	var prevKey KeyWeight
	prevKey.Weight = 0
	for _, k := range auth.Keys {
		if k.Weight == 0 {
			return false
		}
		totalWeight += uint32(k.Weight)
		prevKey = k
	}
	var prevAcc PermissionLevelWeight
	prevAcc.Weight = 0
	for _, a := range auth.Accounts {
		if a.Weight == 0 {
			return false
		}
		totalWeight += uint32(a.Weight)
		prevAcc = a
	}
	var prevWts WaitWeight
	prevWts.Weight = 0
	for _, w := range auth.Waits {
		if w.Weight == 0 {
			return false
		}
		totalWeight += uint32(w.Weight)
		prevWts = w
	}

	return totalWeight >= auth.Threshold
}
