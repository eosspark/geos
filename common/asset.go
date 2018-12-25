package common

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"math"
	"strconv"
	"strings"
)

var maxAmount int64 = int64(1)<<62 - 1

type Asset struct {
	Amount int64 `eos:"asset"`
	Symbol
}

//func (a Asset) Pack(p *fcbuffer.PackStream) error {
//	p.WriteInt64(a.Amount)
//	return a.Symbol.Pack(p)
//
//}
//
//func (a *Asset) Unpack(u *fcbuffer.UnPackStream) error {
//	a.Amount, _ = u.ReadInt64()
//	a.Symbol.Unpack(u)
//	return nil
//}

func (a *Asset) assert() {
	try.EosAssert(a.isAmountWithinRange(), &exception.AssetTypeException{}, "magnitude of asset amount must be less than 2^62")
	try.EosAssert(a.Symbol.Valid(), &exception.AssetTypeException{}, "invalid symbol")
}

func (a *Asset) isAmountWithinRange() bool {
	return -maxAmount <= a.Amount && a.Amount <= maxAmount
}

func (a *Asset) isValid() bool {
	return a.isAmountWithinRange() && a.Symbol.Valid()
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
	sign := ""
	abs := a.Amount
	if a.Amount < 0 {
		sign = "-"
		abs = -1 * a.Amount
	}
	strInt := fmt.Sprintf("%d", abs)
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

	return fmt.Sprintf("%s %s", sign + result, a.Symbol.Symbol)
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
	for i := uint8(0); i < sym.Precision; i++ {
		amount *= 10
	}
	amount += fractPart
	asset := Asset{Amount: amount, Symbol: sym}
	asset.assert()
	return asset
}

type ExtendedAsset struct {
	Asset    Asset `json:"asset"`
	Contract AccountName
}

type SymbolCode uint64

// NOTE: there's also a new ExtendedSymbol (which includes the contract (as AccountName) on which it is)
type Symbol struct {
	Precision uint8
	Symbol    string
}

var MaxPrecision = uint8(18)

//func (s Symbol) Pack(p *fcbuffer.PackStream) error {
//	p.WriteUint8(s.Precision)
//	symbol := make([]byte, 7, 7)
//	copy(symbol[:], []byte(s.Symbol))
//	p.ToWriter(symbol)
//	return nil
//}
//
//func (s *Symbol) Unpack(u *fcbuffer.UnPackStream) error {
//	s.Precision, _ = u.ReadUint8()
//
//	if u.Remaining() < 7 {
//		u.Log.Error("asset symbol required [%d] bytes, remaining [%d]", 7, u.Remaining())
//		return nil
//	}
//	data := u.Data[u.Pos : u.Pos+7]
//	u.Pos += 7
//	s.Symbol = strings.TrimRight(string(data), "\x00")
//	return nil
//}

func (sym Symbol) FromString(from *string) Symbol {
	//TODO: unComplete
	try.EosAssert(!Empty(*from), &exception.SymbolTypeException{}, "creating symbol from empty string")
	commaPos := strings.Index(*from, ",")
	try.EosAssert(commaPos != -1, &exception.SymbolTypeException{}, "missing comma in symbol")
	precPart := string([]byte(*from)[:commaPos])
	p, _ := strconv.ParseInt(precPart, 10, 64)
	namePart := string([]byte(*from)[commaPos+1:])
	try.EosAssert(sym.ValidName(namePart), &exception.SymbolTypeException{}, "invalid symbol: %s", namePart)
	try.EosAssert(uint8(p) <= MaxPrecision, &exception.SymbolTypeException{}, "precision %v should be <= 18", p)
	return Symbol{Precision: uint8(p), Symbol: namePart}
}

func (sym Symbol) String() string {
	try.EosAssert(sym.Valid(), &exception.SymbolTypeException{}, "symbol is not valid")
	v := sym.Precision
	ret := strconv.Itoa(int(v))
	ret += ","+sym.Symbol
	return ret
}

func (sym *Symbol) SymbolValue() uint64 {
	result := uint64(0)
	for i := len(sym.Symbol) - 1; i >= 0; i-- {
		if sym.Symbol[i] < 'A' || sym.Symbol[i] > 'Z' {
			log.Error("symbol cannot exceed A~Z")
		} else {
			result |= uint64(sym.Symbol[i])
		}
		result = result << 8
	}
	result |= uint64(sym.Precision)
	return result
}

func (sym *Symbol) ToSymbolCode() SymbolCode {
	return SymbolCode(sym.SymbolValue()) >> 8
}

func (sym *Symbol) Decimals() uint8 {
	return sym.Precision
}

func (sym *Symbol) Name() string {
	return sym.Symbol
}

func (sym *Symbol) Valid() bool {
	return sym.Decimals() <= MaxPrecision && sym.ValidName(sym.Symbol)
}

func (sym *Symbol) ValidName(name string) bool {
	return -1 == strings.IndexFunc(name, func(r rune) bool {
		return !(r >= 'A' && r <= 'Z')
	})
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
	asset := Asset{Amount: amount, Symbol: EOSSymbol}
	asset.assert()
	return asset
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
	out.assert()
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
