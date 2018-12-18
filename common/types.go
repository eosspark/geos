package common

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/crypto"
				"strconv"
	"strings"
)

type SizeT = int

// For reference:
// https://github.com/mithrilcoin-io/EosCommander/blob/master/app/src/main/java/io/mithrilcoin/eoscommander/data/remote/model/types/EosByteWriter.java
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

type DelegatedBandwidth struct {
	From      AccountName `json:"from"`
	To        AccountName `json:"to"`
	NetWeight Asset       `json:"net_weight"`
	CPUWeight Asset       `json:"cpu_weight"`
}

type TotalResources struct {
	Owner     AccountName `json:"owner"`
	NetWeight Asset       `json:"net_weight"`
	CPUWeight Asset       `json:"cpu_weight"`
	RAMBytes  JSONInt64   `json:"ram_bytes"`
}

type VoterInfo struct {
	Owner             AccountName   `json:"owner"`
	Proxy             AccountName   `json:"proxy"`
	Producers         []AccountName `json:"producers"`
	Staked            JSONInt64     `json:"staked"`
	LastVoteWeight    JSONFloat64   `json:"last_vote_weight"`
	ProxiedVoteWeight JSONFloat64   `json:"proxied_vote_weight"`
	IsProxy           byte          `json:"is_proxy"`
}

// CurrencyName

type CurrencyName string

type Bool bool

func (b *Bool) UnmarshalJSON(data []byte) error {
	var num int
	err := json.Unmarshal(data, &num)
	if err == nil {
		*b = Bool(num != 0)
		return nil
	}

	var boolVal bool
	if err := json.Unmarshal(data, &boolVal); err != nil {
		return fmt.Errorf("couldn't unmarshal bool as int or true/false: %s", err)
	}

	*b = Bool(boolVal)
	return nil
}


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

type JSONFloat64 float64

func (f *JSONFloat64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}

		*f = JSONFloat64(val)

		return nil
	}

	var fl float64
	if err := json.Unmarshal(data, &fl); err != nil {
		return err
	}

	*f = JSONFloat64(fl)

	return nil
}

type JSONInt64 int64

func (i *JSONInt64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return errors.New("empty value")
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}

		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}

		*i = JSONInt64(val)

		return nil
	}

	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*i = JSONInt64(v)

	return nil
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
