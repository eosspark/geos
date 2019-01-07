package common

import (
	"encoding/json"
	. "github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
)

type Variant = interface{}

type StaticVariant = Variant // use type-assert to get static_variant

type Variants = map[string]interface{}

func ToVariant(T interface{}, variant Variant) {
	data, err := json.Marshal(T)
	if err != nil {
		EosThrow(&ParseErrorException{}, err.Error())
	}

	err = json.Unmarshal(data, variant)
	if err != nil {
		EosThrow(&ParseErrorException{}, err.Error())
	}
}

func FromVariant(variant Variant, T interface{}) {
	ToVariant(variant, T)
}

func VariantsFromData(data []byte) Variants {
	result := Variants{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		EosThrow(&ParseErrorException{}, err.Error())
	}

	return result
}