package common

import (
	"strconv"
	"strings"
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/exception"
	"math"
)

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

type SymbolCode uint64

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

func (sym *Symbol) ToSymbolCode() SymbolCode {
	return SymbolCode(N(sym.Symbol)) >> 8
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

