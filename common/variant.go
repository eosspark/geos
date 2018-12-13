package common

import "encoding/json"

type Variant = interface{}

type StaticVariant = Variant // use type-assert to get static_variant

type Variants = map[string]interface{}

func ToVariant(T interface{}, variant Variant) error {
	data, err := json.Marshal(T)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, variant); err != nil {
		return err
	}

	return nil
}

func FromVariant(variant Variant, T interface{}) error {
	return ToVariant(variant, T)
}

func VariantToVariants(variant Variant) (Variants, bool) {
	vs, ok := variant.(Variants)
	return vs, ok
}
