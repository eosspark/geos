package types

import (
	"reflect"

	"github.com/eosspark/eos-go/common"
	. "github.com/eosspark/eos-go/common/eos_math"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
)

func IsIntegral(T interface{}) bool {
	switch T.(type) {
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case uint:
		return true
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	case bool:
		return true
	default:
		return false
	}
}

func ConvertToWordT(T interface{}) WordT {
	switch v := T.(type) {
	case int:
		if v >= 0 {
			return Uint128{Low: uint64(v), High: 0}
		} else {
			return Uint128{Low: 0, High: 0}
		}
	case int8:
		if v >= 0 {
			return Uint128{Low: uint64(v), High: 0}
		} else {
			return Uint128{Low: 0, High: 0}
		}
	case int16:
		if v >= 0 {
			return Uint128{Low: uint64(v), High: 0}
		} else {
			return Uint128{Low: 0, High: 0}
		}
	case int32:
		if v >= 0 {
			return Uint128{Low: uint64(v), High: 0}
		} else {
			return Uint128{Low: 0, High: 0}
		}
	case int64:
		if v >= 0 {
			return Uint128{Low: uint64(v), High: 0}
		} else {
			return Uint128{Low: 0, High: 0}
		}
	case uint:
		return Uint128{Low: uint64(v), High: 0}
	case uint8:
		return Uint128{Low: uint64(v), High: 0}
	case uint16:
		return Uint128{Low: uint64(v), High: 0}
	case uint32:
		return Uint128{Low: uint64(v), High: 0}
	case uint64:
		return Uint128{Low: uint64(v), High: 0}
	case bool:
		if v == true {
			return Uint128{Low: 1, High: 0}
		} else {
			return Uint128{Low: 0, High: 0}
		}

	default:
		return Uint128{Low: 0, High: 0}
	}
}

func SizeOf(T interface{}) int {
	return int(reflect.TypeOf(T).Size())
}

func IsSame(T interface{}, compType interface{}) bool {
	return reflect.TypeOf(T) == reflect.TypeOf(compType)
}

type WordT = Uint128

var wordTSize = SizeOf(WordT{})

var SizeOfWordT = SizeOf(WordT{})

type FixedKey struct {
	Size common.SizeT
	data []WordT
}

func (f FixedKey) numWords() common.SizeT {
	return (f.Size + SizeOfWordT - 1) / SizeOfWordT
}

func (f FixedKey) EnableFirstWord(first interface{}) bool {
	var boolType bool
	return IsIntegral(first) && !IsSame(first, boolType) && SizeOf(first) <= wordTSize
}

func (f *FixedKey) SetFromWordSequence(arr []interface{}) {

	wordSize := SizeOf(reflect.TypeOf(arr[0]))
	tempWord := CreateUint128(int(0))
	subWordShift := 8 * wordSize
	numSubWords := wordTSize / wordSize
	subWordsLeft := numSubWords
	f.data = make([]WordT, f.numWords())
	i := 0
	for _, w := range arr {
		if subWordsLeft > 1 {
			tempWord = tempWord.Or(ConvertToWordT(w))
			tempWord.LeftShifts(subWordShift)
			subWordsLeft--
			continue
		}

		EosAssert(subWordsLeft == 1, &FixedKeyTypeException{}, "unexpected error in fixed_key constructor")
		tempWord = tempWord.Or(ConvertToWordT(w))
		subWordsLeft = numSubWords

		f.data[i] = tempWord
		tempWord = Uint128{Low: 0, High: 0}
		i++
	}
	if subWordsLeft != numSubWords {
		if subWordsLeft > 1 {
			tempWord.LeftShifts(8 * int(subWordsLeft-1))
		}
		f.data[i] = tempWord
	}
}

func (f *FixedKey) MakeFromWordSequence(first interface{}, rest ...interface{}) *FixedKey {
	Try(func() {
		FcAssert(f.EnableFirstWord(first), "The first word is invalid")
		FcAssert(wordTSize == wordTSize/SizeOf(first)*SizeOf(first),
			"size of the backing word size is not divisible by the size of the words supplied as arguments")
		FcAssert(SizeOf(first)*(1+len(rest)) <= f.Size, "too many words supplied to make_from_word_sequence")
		var arr []interface{}
		arr = append(arr, first)
		arr = append(arr, rest...)
		f.SetFromWordSequence(arr)
	}).Catch(func(e Exception) {
		log.Error(e.DetailMessage())
	}).End()
	return f
}

func (f *FixedKey) GetArray() []WordT {
	return f.data
}
