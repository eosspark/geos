package common

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/eosspark/eos-go/crypto"
	"strings"
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
	return *crypto.NewSha256Nil()
}

func TransactionIdNil() TransactionIdType {
	return *crypto.NewSha256Nil()
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

type Varuint32 struct {
	V uint32 `eos:"vuint32"`
}
type Varint32 struct {
	V int32 `eos:"vint32"`
}
