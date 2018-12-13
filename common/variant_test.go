package common_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/eosspark/eos-go/common"
		"encoding/json"
)

func TestToVariant(t *testing.T) {
	example := struct {
		A int
		B string
		C []int
	}{1, "a", []int{2, 3}}

	var variant common.Variant //:= make(map[string]interface{}, 3)

	err := common.ToVariant(&example, &variant)
	assert.NoError(t, err)

	variant.(common.Variants)["B"] = "b"

	err = common.FromVariant(&variant, &example)

	assert.NoError(t, err)
	assert.Equal(t, "b", example.B)
	assert.Equal(t, []int{2, 3}, example.C)
}

func TestToVariant_Simple(t *testing.T) {
	example := 100
	var variant common.Variant

	err := common.ToVariant(&example, &variant)
	assert.NoError(t, err)

	vs, ok := common.VariantToVariants(variant)
	assert.Empty(t, vs)
	assert.Equal(t, false, ok)

	var a int
	err = common.FromVariant(&variant, &a)
	assert.NoError(t, err)
	assert.Equal(t, 100, a)
}

func TestVariant(t *testing.T) {
	example := struct {
		A int
		B string
		C []int
	}{1, "a", []int{2, 3}}

	var variant common.Variant
	err := common.ToVariant(&example, &variant)
	assert.NoError(t, err)

	data, err := json.Marshal(variant)
	assert.NoError(t, err)

	assert.Equal(t, "{\"A\":1,\"B\":\"a\",\"C\":[2,3]}", string(data))

	vs,_ := common.VariantToVariants(variant)

	(vs)["B"] = "b"

	err = common.FromVariant(variant, &example)

	assert.Equal(t, "b", example.B)
}

func Example_StaticVariant() {
	var example common.StaticVariant = int(1)

	if vt, ok := example.(int); ok {
		vt ++
	} else if _, ok := example.(string); ok {
		//else operation
	}
}
