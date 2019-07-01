package common

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/eosspark/eos-go/crypto"
)

type SizeT = int

type ChainIdType = crypto.Sha256
type NodeIdType = crypto.Sha256
type BlockIdType = crypto.Sha256
type TransactionIdType = crypto.Sha256
type CheckSum256Type = crypto.Sha256
type DigestType = crypto.Sha256
type IdType = int64
type KeyType = uint64

func BlockIdNil() BlockIdType {
	return crypto.NewSha256Nil()
}

func TransactionIdNil() TransactionIdType {
	return crypto.NewSha256Nil()
}

func DecodeIdTypeString(str string) (id [4]uint64, err error) {
	b, err := hex.DecodeString(str)
	if err != nil {
		return
	}

	for i := range id {
		id[i] = binary.LittleEndian.Uint64(b[i*8 : (i+1)*8])
	}

	return
}

func DecodeIdTypeByte(b []byte) (id [4]uint64, err error) {
	for i := range id {
		id[i] = binary.LittleEndian.Uint64(b[i*8 : (i+1)*8])
	}

	return id, nil
}

// CurrencyName

type CurrencyName string

// HexBytes
type HexBytes []byte

func (t HexBytes) Size() int {
	return len([]byte(t))
}

func (t HexBytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(t))
}

func (t *HexBytes) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	*t, err = hex.DecodeString(s)
	return
}

type Compare interface {
	String() string
}

func CompareString(a, b Compare) int {
	return strings.Compare(a.String(), b.String())
}

func Min(x, y uint64) uint64 {
	if x < y {
		return x
	} else {
		return y
	}
}

func Max(x, y uint64) uint64 {
	if x > y {
		return x
	} else {
		return y
	}
}

type Vuint32 uint32

func (v Vuint32) Pack() ([]byte, error) {
	return WriteUVarInt(int(v)), nil
}
func (v *Vuint32) Unpack(in []byte) (l int, err error) {
	re, l, err := ReadUvarint64(in)
	if err != nil {
		return 0, nil
	}
	*v = Vuint32(re)
	return l, nil
}

type Vint32 int32

func (v Vint32) Pack() ([]byte, error) {
	return WriteVarInt(int(v)), nil
}
func (v *Vint32) Unpack(in []byte) (l int, err error) {
	re, l, err := ReadVarint64(in)
	if err != nil {
		return 0, nil
	}
	*v = Vint32(re)
	return l, nil
}

type PermissionLevel struct {
	Actor      AccountName    `json:"actor"`
	Permission PermissionName `json:"permission"`
}

func ComparePermissionLevel(first interface{}, second interface{}) int {
	if first.(PermissionLevel).Actor > second.(PermissionLevel).Actor {
		return 1
	} else if first.(PermissionLevel).Actor < second.(PermissionLevel).Actor {
		return -1
	}
	if first.(PermissionLevel).Permission > second.(PermissionLevel).Permission {
		return 1
	} else if first.(PermissionLevel).Permission < second.(PermissionLevel).Permission {
		return -1
	} else {
		return 0
	}
}

func (level PermissionLevel) String() string {
	return "{ actor: " + level.Actor.String() + ", " + "permission: " + level.Permission.String() + "}"
}

type AccountDelta struct {
	Account AccountName
	Delta   int64
}

func NewAccountDelta(name AccountName, d int64) *AccountDelta {
	return &AccountDelta{name, d}
}

func CompareAccountDelta(first interface{}, second interface{}) int {
	if first.(AccountDelta).Account == second.(AccountDelta).Account {
		return 0
	}
	if first.(AccountDelta).Account < second.(AccountDelta).Account {
		return -1
	}
	return 1
}

type NamePair struct {
	First  AccountName
	Second ActionName
}

func CompareNamePair(a, b interface{}) int {
	aPair, bPair := a.(NamePair), b.(NamePair)
	aFirst, bFirst := aPair.First, bPair.First
	aSecond, bSecond := aPair.Second, bPair.Second

	switch {
	case aFirst > bFirst:
		return 1
	case aFirst < bFirst:
		return -1
	}

	switch {
	case aSecond > bSecond:
		return 1
	case aSecond < bSecond:
		return -1
	}

	return 0
}
