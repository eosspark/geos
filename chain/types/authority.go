package types

import (
	"github.com/eosspark/eos-go/common"
)

type PermissionLevelWeight struct {
	Permission common.PermissionLevel `json:"permission"`
	Weight     common.WeightType      `json:"weight"` // weight_type
}

type KeyWeight struct {
	Key    common.PublicKeyType `json:"key"`
	Weight common.WeightType	 `json:"weight"`
}

type WaitWeight struct {
	WaitSec uint32         `json:"wait_sec"`
	Weight  common.WeightType `json:"weight"`
}

type Authority struct { 
	Threshold uint32             	  `json:"threshold"`
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

func (auth Authority) Equals(author Authority) bool {
	return true
}

func (sharedAuth SharedAuthority) Equals(sharedAuthor SharedAuthority) bool {
	return true
}

func (sharedAuth SharedAuthority) GetBillableSize() uint64 {           //返回值类型不确定
	accountSize := uint64(len(sharedAuth.Accounts)) * common.BillableSizeV("permission_level_weight")
	waitsSize := uint64(len(sharedAuth.Waits)) * common.BillableSizeV("wait_weight")
	keysSize := uint64(0)
	for _, key := range sharedAuth.Keys {
		keysSize += common.BillableSizeV("key_weight")
		keysSize += uint64(key.Weight) * 0                             //待修改
	}
	return accountSize + waitsSize + keysSize
}

func Validate(auth Authority) bool {
	var totalWeight uint32 = 0
	if len(auth.Accounts)+len(auth.Keys)+len(auth.Waits) > 1 << 16 {
		return false
	}
	if auth.Threshold == 0 {
		return false
	}
	var prevKey KeyWeight
	prevKey.Weight = 0
	for _, k := range auth.Keys{
		if k.Weight == 0 {
			return false
		}
		totalWeight += uint32(k.Weight)
		prevKey = k
	}
	var prevAcc PermissionLevelWeight
	prevAcc.Weight = 0
	for _, a := range auth.Accounts{
		if a.Weight == 0 {
			return false
		}
		totalWeight += uint32(a.Weight)
		prevAcc = a
	}
	var prevWts WaitWeight
	prevWts.Weight = 0
	for _, w := range auth.Waits{
		if w.Weight == 0 {
			return false
		}
		totalWeight += uint32(w.Weight)
		prevWts = w
	}

	return totalWeight >= auth.Threshold
}