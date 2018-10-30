package types

import (
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"strings"
	"strconv"
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
	Weight     WeightType      `json:"weight"`
}

type KeyWeight struct {
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

	out.Actor = common.AccountName(common.N(parts[0]))
	out.Permission = common.PermissionName(common.N("active"))

	if len(parts) == 2 {
		if len(parts[1]) > 12 {
			return out, fmt.Errorf("permission %q name too long", parts[1])
		}

		out.Permission = common.PermissionName(common.N("active"))
	}

	return
}

func NewAuthority(k ecc.PublicKey, delaySec uint32) (a Authority) {
	a.Threshold = 1
	a.Keys[0] = KeyWeight{k,1}
	if delaySec > 0 {
		a.Threshold = 2
		a.Waits[0] = WaitWeight{delaySec, 1}
	}
	return a
}

func (auth *Authority) ToSharedAuthority() SharedAuthority {
	return SharedAuthority{auth.Threshold, auth.Keys, auth.Accounts, auth.Waits}
}

func (sharedAuth *SharedAuthority) ToAuthority() Authority {
	return Authority{sharedAuth.Threshold,sharedAuth.Keys,sharedAuth.Accounts,sharedAuth.Waits}
}

func (weight WeightType) String() string {
	return strconv.FormatInt(int64(weight),10)
}

func (level PermissionLevel) String() string {
	return "{ actor: " + level.Actor.String() + ", " + "permission: " + level.Permission.String() + "}"
}

func (key KeyWeight) String() string {
	return "{ key: " + key.Key.String() + ", " + " weight: " + key.Weight.String() + "} "
}

func (permLevel PermissionLevelWeight) String() string {
	return "{ permission: " + permLevel.Permission.String() + ", " + "weight: " + permLevel.Weight.String() + "}"
}

func (wait WaitWeight) String() string{
	return "{ weightSec: " + strconv.FormatInt(int64(wait.WaitSec), 10) + "weight" + wait.Weight.String() + "}"
}

func (auth Authority) String() string {
	ThresholdStr := "threshold: " + strconv.FormatInt(int64(auth.Threshold),10)
	KeysStr := "keys: ["
	for _, key := range auth.Keys {
		KeysStr += "key: " + key.String()
		if key != auth.Keys[len(auth.Keys)-1]{
			KeysStr += ", "
		}
	}
	KeysStr += "]"
	AccountsStr := "accounts: ["
	for _, account := range auth.Accounts {
		fmt.Println(account)
		fmt.Println(auth.Accounts[len(auth.Accounts)-1])
		AccountsStr += "account: " + account.String()
		if account != auth.Accounts[len(auth.Accounts)-1]{
			AccountsStr += ", "
		}
	}
	AccountsStr += "]"
	WaitsStr := "waits: ["
	for _, wait := range auth.Waits {
		WaitsStr += "account: " + wait.String()
		if wait != auth.Waits[len(auth.Waits)-1]{
			WaitsStr += ", "
		}
	}
	WaitsStr += "]"
	return "{ "+ ThresholdStr + ", " + KeysStr + ", " + AccountsStr + ", " + WaitsStr + "}"
}

func (auth Authority) Equals(author Authority) bool {
	return true
}

func (sharedAuth SharedAuthority) Equals(sharedAuthor SharedAuthority) bool {
	return true
}

func (sharedAuth SharedAuthority) GetBillableSize() uint64 { //TODO
	accountSize := uint64(len(sharedAuth.Accounts)) * common.BillableSizeV("permission_level_weight")
	waitsSize := uint64(len(sharedAuth.Waits)) * common.BillableSizeV("wait_weight")
	keysSize := uint64(0)
	for _, key := range sharedAuth.Keys {
		keysSize += common.BillableSizeV("key_weight")
		keysSize += uint64(key.Weight) * 0 //TODO
	}
	return accountSize + waitsSize + keysSize
}

func Validate(auth Authority) bool {      //TODO: sort.
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
