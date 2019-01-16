package types

import (
	"encoding/binary"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"reflect"
	"strconv"
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

func (level *PermissionLevel) GetKey() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(level.Actor))
	return b
}

//for treeset
var PermissionLevelType = reflect.TypeOf(PermissionLevel{})

func ComparePermissionLevel(first interface{}, second interface{}) int {
	if first.(PermissionLevel).Permission > second.(PermissionLevel).Permission {
		return 1
	} else if first.(PermissionLevel).Permission < second.(PermissionLevel).Permission {
		return -1
	}
	if first.(PermissionLevel).Actor > second.(PermissionLevel).Actor {
		return 1
	} else if first.(PermissionLevel).Actor < second.(PermissionLevel).Actor {
		return -1
	} else {
		return 0
	}
}

type PermissionLevelWeight struct {
	Permission PermissionLevel `json:"permission"`
	Weight     WeightType      `json:"weight"`
}

type KeyWeight struct {
	Key    ecc.PublicKey `json:"key"`
	Weight WeightType    `json:"weight"`
}

func (key KeyWeight) Compare(kw KeyWeight) bool {
	if !key.Key.Compare(kw.Key) {
		return false
	}
	if key.Weight != kw.Weight {
		return false
	}
	return true
}

type WaitWeight struct {
	WaitSec uint32     `json:"wait_sec"`
	Weight  WeightType `json:"weight"`
}

type Authority struct {
	Threshold uint32                  `json:"threshold"`
	Keys      []KeyWeight             `json:"keys"`
	Accounts  []PermissionLevelWeight `json:"accounts"`
	Waits     []WaitWeight            `json:"waits"`
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
	a.Keys = append(a.Keys, KeyWeight{k, 1})
	if delaySec > 0 {
		a.Threshold = 2
		a.Waits = append(a.Waits, WaitWeight{delaySec, 1})
	}
	a.Accounts = make([]PermissionLevelWeight, 0)
	return a
}

func (auth *Authority) ToSharedAuthority() SharedAuthority {
	return SharedAuthority{auth.Threshold, auth.Keys, auth.Accounts, auth.Waits}
}

func (sharedAuth *SharedAuthority) ToAuthority() Authority {
	return Authority{sharedAuth.Threshold, sharedAuth.Keys, sharedAuth.Accounts, sharedAuth.Waits}
}

func (weight WeightType) String() string {
	return strconv.FormatInt(int64(weight), 10)
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

func (wait WaitWeight) String() string {
	return "{ weightSec: " + strconv.FormatInt(int64(wait.WaitSec), 10) + "weight" + wait.Weight.String() + "}"
}

func (auth Authority) String() string {
	ThresholdStr := "threshold: " + strconv.FormatInt(int64(auth.Threshold), 10)
	KeysStr := "keys: ["
	for _, key := range auth.Keys {
		KeysStr += "key: " + key.String()
		if !key.Compare(auth.Keys[len(auth.Keys)-1]) {
			KeysStr += ", "
		}
	}
	KeysStr += "]"
	AccountsStr := "accounts: ["
	for _, account := range auth.Accounts {
		AccountsStr += "account: " + account.String()
		if account != auth.Accounts[len(auth.Accounts)-1] {
			AccountsStr += ", "
		}
	}
	AccountsStr += "]"
	WaitsStr := "waits: ["
	for _, wait := range auth.Waits {
		WaitsStr += "account: " + wait.String()
		if wait != auth.Waits[len(auth.Waits)-1] {
			WaitsStr += ", "
		}
	}
	WaitsStr += "]"
	return "{ " + ThresholdStr + ", " + KeysStr + ", " + AccountsStr + ", " + WaitsStr + "}"
}

func (auth Authority) Equals(author Authority) bool {
	return auth.ToSharedAuthority().Equals(author.ToSharedAuthority())
}

func (sharedAuth SharedAuthority) Equals(sharedAuthor SharedAuthority) bool {
	if sharedAuth.Threshold != sharedAuthor.Threshold {
		return false
	}
	if len(sharedAuth.Keys) != len(sharedAuthor.Keys) || len(sharedAuth.Accounts) != len(sharedAuthor.Accounts) || len(sharedAuth.Waits) != len(sharedAuthor.Waits) {
		return false
	}
	for i := range sharedAuth.Keys {
		if sharedAuth.Keys[i] != sharedAuthor.Keys[i] {
			return false
		}
	}
	for j := range sharedAuth.Accounts {
		if sharedAuth.Accounts[j] != sharedAuthor.Accounts[j] {
			return false
		}
	}
	for k := range sharedAuth.Waits {
		if sharedAuth.Waits[k] != sharedAuthor.Waits[k] {
			return false
		}
	}
	return true
}

func (sharedAuth SharedAuthority) GetBillableSize() uint64 {
	accountSize := uint64(len(sharedAuth.Accounts)) * common.BillableSizeV("permission_level_weight")
	waitsSize := uint64(len(sharedAuth.Waits)) * common.BillableSizeV("wait_weight")
	keysSize := uint64(0)
	keySize := 0
	for _, key := range sharedAuth.Keys {
		keysSize += common.BillableSizeV("key_weight")
		keySize, _ = rlp.EncodeSize(key.Key)
		keysSize += uint64(keySize)
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

	for i, k := range auth.Keys {
		if i > 0 && ! (ecc.ComparePubKey(auth.Keys[i-1].Key, k.Key) == -1){
			return false
		}
		totalWeight += uint32(k.Weight)
	}

	for i, a := range auth.Accounts {
		if i > 0 && ! (ComparePermissionLevel(auth.Accounts[i-1].Permission, a.Permission) == -1){
			return false
		}
		totalWeight += uint32(a.Weight)
	}

	for i, w := range auth.Waits {
		if i > 0 && auth.Waits[i-1].WaitSec >= w.WaitSec {
			return false
		}
		totalWeight += uint32(w.Weight)
	}

	return totalWeight >= auth.Threshold
}
