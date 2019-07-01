package common_test

import (
	"encoding/json"
	"testing"

	"github.com/eosspark/eos-go/common"
	"github.com/stretchr/testify/assert"
)

func TestToVariant(t *testing.T) {
	example := struct {
		A int
		B string
		C []int
	}{1, "a", []int{2, 3}}

	var variant common.Variant //:= make(map[string]interface{}, 3)

	common.ToVariant(&example, &variant)

	variant.(common.Variants)["B"] = "b"

	common.FromVariant(&variant, &example)

	assert.Equal(t, "b", example.B)
	assert.Equal(t, []int{2, 3}, example.C)
}

func TestToVariant_Simple(t *testing.T) {
	example := 100
	var variant common.Variant

	common.ToVariant(&example, &variant)

	var a int
	common.FromVariant(&variant, &a)
	assert.Equal(t, 100, a)
}

func TestVariant(t *testing.T) {
	example := struct {
		A int
		B string
		C []int
	}{1, "a", []int{2, 3}}

	var variant common.Variant
	common.ToVariant(&example, &variant)

	data, err := json.Marshal(variant)
	assert.NoError(t, err)

	assert.Equal(t, "{\"A\":1,\"B\":\"a\",\"C\":[2,3]}", string(data))

	(variant.(common.Variants))["B"] = "b"

	common.FromVariant(variant, &example)

	assert.Equal(t, "b", example.B)
}

func Example_Variant() {
	var example common.StaticVariant = int(1)

	if vt, ok := example.(int); ok {
		vt++
	} else if _, ok := example.(string); ok {
		//else operation
	}
}
