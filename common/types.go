package common

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"math"
	"strconv"
	"strings"
	"time"
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

type AccountName = Name
type PermissionName = Name
type ActionName = Name
type TableName = Name
type ScopeName = Name

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

type RefundRequest struct {
	Owner       AccountName `json:"owner"`
	RequestTime JSONTime    `json:"request_time"` //         {"name":"request_time", "type":"time_point_sec"},
	NetAmount   Asset       `json:"net_amount"`
	CPUAmount   Asset       `json:"cpu_amount"`
}

type CompressionType uint8

const (
	CompressionNone = CompressionType(iota)
	CompressionZlib
)

func (c CompressionType) String() string {
	switch c {
	case CompressionNone:
		return "none"
	case CompressionZlib:
		return "zlib"
	default:
		return ""
	}
}

func (c CompressionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *CompressionType) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	switch s {
	case "zlib":
		*c = CompressionZlib
	default:
		*c = CompressionNone
	}
	return nil
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

// Asset

// NOTE: there's also ExtendedAsset which is a quantity with the attached contract (AccountName)
type Asset struct {
	Amount int64 `eos:"asset"`
	Symbol
}

func (a Asset) Add(other Asset) Asset {
	if a.Symbol != other.Symbol {
		panic("Add applies only to assets with the same symbol")
	}
	return Asset{Amount: a.Amount + other.Amount, Symbol: a.Symbol}
}

func (a Asset) Sub(other Asset) Asset {
	if a.Symbol != other.Symbol {
		panic("Sub applies only to assets with the same symbol")
	}
	return Asset{Amount: a.Amount - other.Amount, Symbol: a.Symbol}
}

func (a Asset) String() string {
	strInt := fmt.Sprintf("%d", a.Amount)
	if len(strInt) < int(a.Symbol.Precision+1) {
		// prepend `0` for the difference:
		strInt = strings.Repeat("0", int(a.Symbol.Precision+uint8(1))-len(strInt)) + strInt
	}

	var result string
	if a.Symbol.Precision == 0 {
		result = strInt
	} else {
		result = strInt[:len(strInt)-int(a.Symbol.Precision)] + "." + strInt[len(strInt)-int(a.Symbol.Precision):]
	}

	return fmt.Sprintf("%s %s", result, a.Symbol.Symbol)
}

func (a Asset) FromString(from *string) Asset {
	spacePos := strings.Index(*from, " ")
	try.EosAssert(spacePos != -1, &exception.AssetTypeException{}, "Asset's amount and symbol should be separated with space")
	symbolStr := string([]byte(*from)[spacePos+1:])
	amountStr := string([]byte(*from)[:spacePos])

	dotPos := strings.Index(amountStr, ".")
	if dotPos != -1 {
		try.EosAssert(dotPos != len(amountStr)-1, &exception.AssetTypeException{}, "Missing decimal fraction after decimal point")
	}

	var precisionDigitStr string
	if dotPos != -1 {
		precisionDigitStr = strconv.Itoa(len(amountStr) - dotPos - 1)
	} else {
		precisionDigitStr = "0"
	}

	symbolPart := precisionDigitStr + "," + symbolStr
	sym := Symbol{}.FromString(&symbolPart)

	var intPart, fractPart int64
	if dotPos != -1 {
		intPart, _ = strconv.ParseInt(string([]byte(amountStr)[:dotPos]), 10, 64)
		fractPart, _ = strconv.ParseInt(string([]byte(amountStr)[dotPos+1:]), 10, 64)
		if amountStr[0] == '-' {
			fractPart *= -1
		}
	} else {
		intPart, _ = strconv.ParseInt(amountStr, 10, 64)
	}
	amount := intPart
	amount += fractPart
	return Asset{Amount: amount, Symbol: sym}
}

type ExtendedAsset struct {
	Asset    Asset `json:"asset"`
	Contract AccountName
}

type SymbolCode = uint64

// NOTE: there's also a new ExtendedSymbol (which includes the contract (as AccountName) on which it is)
type Symbol struct {
	Precision uint8
	Symbol    string
}

var MaxPrecision = uint8(18)

func (sym Symbol) FromString(from *string) Symbol {
	//TODO: unComplete
	try.EosAssert(!Empty(*from), &exception.SymbolTypeException{}, "creating symbol from empty string")
	commaPos := strings.Index(*from, ",")
	try.EosAssert(commaPos != -1, &exception.SymbolTypeException{}, "missing comma in symbol")
	precPart := string([]byte(*from)[:commaPos])
	p, _ := strconv.ParseInt(precPart, 10, 64)
	namePart := string([]byte(*from)[commaPos+1:])
	try.EosAssert(uint8(p) <= MaxPrecision, &exception.SymbolTypeException{}, "precision %v should be <= 18", p)
	return Symbol{Precision: uint8(p), Symbol: namePart}
}

// EOSSymbol represents the standard EOS symbol on the chain.  It's
// here just to speed up things.
var EOSSymbol = Symbol{Precision: 4, Symbol: "EOS"}

func NewEOSAssetFromString(amount string) (out Asset, err error) {
	if len(amount) == 0 {
		return out, fmt.Errorf("cannot be an empty string")
	}

	if strings.Contains(amount, " EOS") {
		amount = strings.Replace(amount, " EOS", "", 1)
	}
	if !strings.Contains(amount, ".") {
		val, err := strconv.ParseInt(amount, 10, 64)
		if err != nil {
			return out, err
		}
		return NewEOSAsset(val * 10000), nil
	}

	parts := strings.Split(amount, ".")
	if len(parts) != 2 {
		return out, fmt.Errorf("cannot have two . in amount")
	}

	if len(parts[1]) > 4 {
		return out, fmt.Errorf("EOS has only 4 decimals")
	}

	val, err := strconv.ParseInt(strings.Replace(amount, ".", "", 1), 10, 64)
	if err != nil {
		return out, err
	}
	return NewEOSAsset(val * int64(math.Pow10(4-len(parts[1])))), nil
}

func NewEOSAsset(amount int64) Asset {
	return Asset{Amount: amount, Symbol: EOSSymbol}
}

// NewAsset parses a string like `1000.0000 EOS` into a properly setup Asset
func NewAsset(in string) (out Asset, err error) {
	sec := strings.SplitN(in, " ", 2)
	if len(sec) != 2 {
		return out, fmt.Errorf("invalid format %q, expected an amount and a currency symbol", in)
	}

	if len(sec[1]) > 7 {
		return out, fmt.Errorf("currency symbol %q too long", sec[1])
	}

	out.Symbol.Symbol = sec[1]
	amount := sec[0]
	amountSec := strings.SplitN(amount, ".", 2)

	if len(amountSec) == 2 {
		out.Symbol.Precision = uint8(len(amountSec[1]))
	}

	val, err := strconv.ParseInt(strings.Replace(amount, ".", "", 1), 10, 64)
	if err != nil {
		return out, err
	}

	out.Amount = val

	return
}

func (a *Asset) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	asset, err := NewAsset(s)
	if err != nil {
		return err
	}

	*a = asset

	return nil
}

func (a Asset) MarshalJSON() (data []byte, err error) {
	return json.Marshal(a.String())
}

// JSONTime

type JSONTime struct {
	time.Time
}

const JSONTimeFormat = "2006-01-02T15:04:05"

func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", t.Format(JSONTimeFormat))), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return nil
	}

	t.Time, err = time.Parse(`"`+JSONTimeFormat+`"`, string(data))
	return err
}

// HexBytes

type HexBytes []byte

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
