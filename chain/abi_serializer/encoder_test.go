package abi_serializer

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"strings"
	"testing"
)

var abiString = `
{
	"version": "eosio::abi/1.0",
	"types": [{
		"new_type_name": "new_type_name_1",
		"type": "name"
	}],
	"structs": [
	{
		"name": "struct_name_1",
		"base": "struct_name_2",
		"fields": [
			{"name":"struct_1_field_1", "type":"new_type_name_1"},
			{"name":"struct_1_field_2", "type":"struct_name_3"},
			{"name":"struct_1_field_3", "type":"string?"},
			{"name":"struct_1_field_4", "type":"string?"},
			{"name":"struct_1_field_5", "type":"struct_name_4[]"}
		]
   },{
		"name": "struct_name_2",
		"base": "",
		"fields": [
			{"name":"struct_2_field_1", "type":"string"}
		]
   },{
		"name": "struct_name_3",
		"base": "",
		"fields": [
			{"name":"struct_3_field_1", "type":"string"}
		]
   },{
		"name": "struct_name_4",
		"base": "",
		"fields": [
			{"name":"struct_4_field_1", "type":"string"}
		]
   }
	],
  "actions": [{
		"name": "action_name_1",
		"type": "struct_name_1",
		"ricardian_contract": ""
  }],
  "tables": [{
      "name": "table_name_1",
      "index_type": "i64",
      "key_names": [
        "key_name_1"
      ],
      "key_types": [
        "string"
      ],
      "type": "struct_name_1"
    }
  ]
}
`

var abiData = []byte(`{
	"struct_2_field_1": "struct_2_field_1_value",
	"struct_1_field_1": Name("eoscanadacom"),
	"struct_1_field_2": M{
		"struct_3_field_1": "struct_3_field_1_value",
	},
	"struct_1_field_3": "struct_1_field_3_value",
	//"struct_1_field_4": "struct_1_field_4_value",
	"struct_1_field_5": ["struct_1_field_5_value_1","struct_1_field_5_value_2"],
}`)

func TestABIEncoder_Encode(t *testing.T) {

	testCases := []map[string]interface{}{
		{"caseName": "sunny path", "actionName": "action_name_1", "expectedError": nil, "abi": abiString},
		{"caseName": "missing action", "actionName": "badactionname", "expectedError": fmt.Errorf("encode action: action badactionname not found in abi"), "abi": abiString},
	}

	for _, c := range testCases {
		caseName := c["caseName"].(string)
		t.Run(caseName, func(t *testing.T) {

			abi, err := NewABI(strings.NewReader(c["abi"].(string)))
			assert.NoError(t, err)
			_, err = abi.EncodeAction(common.AccountName(common.N(c["actionName"].(string))), abiData)
			assert.Equal(t, c["expectedError"], err)

			if c["expectedError"] != nil {
				return
			}

			//decoder := NewABIDecoder(buf.Bytes(), strings.NewReader(abiString))
			//result := make(M)
			//err = decoder.Decode(result, ActionName(c["actionName"].(string)))
			//assert.NoError(t, err)

			//assert.Equal(t, abiData, result)
			//fmt.Println(result)
		})
	}
}

func TestABIEncoder_encodeMissingActionStruct(t *testing.T) {

	abiString := `
{
	"version": "eosio::abi/1.0",
	"types": [{
		"new_type_name": "new.type.name.1",
		"type": "name"
	}],
	"structs": [
	],
  "actions": [{
		"name": "action.name.1",
		"type": "struct.name.1",
		"ricardian_contract": ""
  }]
}
`
	abi, err := NewABI(strings.NewReader(abiString))
	assert.NoError(t, err)
	_, err = abi.EncodeAction(common.ActionName(common.N("action.name.1")), abiData)
	assert.Equal(t, fmt.Errorf("encode action: encode struct [struct.name.1] not found in abi"), err)
}

func TestABIEncoder_encodeErrorInBase(t *testing.T) {

	abiString := `
{
	"version": "eosio::abi/1.0",
	"types": [{
		"new_type_name": "new.type.name.1",
		"type": "name"
	}],
	"structs": [
	{
		"name": "struct.name.1",
		"base": "struct.name.2",
		"fields": [
			{"name":"struct.1.field.1", "type":"new.type.name.1"}
		]
   }
	],
  "actions": [{
		"name": "action.name.1",
		"type": "struct.name.1",
		"ricardian_contract": ""
  }]
}
`
	abi, err := NewABI(strings.NewReader(abiString))
	assert.NoError(t, err)
	_, err = abi.EncodeAction(common.ActionName(common.N("action.name.1")), abiData)
	assert.Equal(t, fmt.Errorf("encode action: encode base [struct.name.1]: encode struct [struct.name.2] not found in abi"), err)
}

func TestABIEncoder_encodeField(t *testing.T) {

	testCases := []map[string]interface{}{
		{"caseName": "sunny path", "fieldName": "field_name", "fieldType": "string", "expectedValue": "0f6669656c642e312e76616c75652e31", "json": "{\"field_name\": \"field.1.value.1\"}", "isOptional": false, "isArray": false, "expectedError": nil, "writer": new(bytes.Buffer)},
		{"caseName": "optional present", "fieldName": "field_name", "fieldType": "string", "expectedValue": "010f6669656c642e312e76616c75652e31", "json": "{\"field_name\": \"field.1.value.1\"}", "isOptional": true, "isArray": false, "expectedError": nil, "writer": new(bytes.Buffer)},
		{"caseName": "optional not present", "fieldName": "field_name", "fieldType": "string", "expectedValue": "00", "json": "{\"field_name_other\": \"field.1.value.2\"}", "isOptional": true, "isArray": false, "expectedError": nil, "writer": new(bytes.Buffer)},
		{"caseName": "optional present write flag err", "fieldName": "field_name", "fieldType": "string", "expectedValue": "010f6669656c642e312e76616c75652e31", "json": "{\"field_name\": \"field.1.value.1\"}", "isOptional": true, "isArray": false, "expectedError": fmt.Errorf("error.1"), "writer": mockWriter{err: fmt.Errorf("error.1")}},
		{"caseName": "not optional not present", "fieldName": "field_name", "fieldType": "string", "expectedValue": "00", "json": "{\"field_name_other\": \"field.1.value.2\"}", "isOptional": false, "isArray": false, "expectedError": fmt.Errorf("encode field: none optional field [field_name] as a nil value"), "writer": new(bytes.Buffer)},
		{"caseName": "array", "fieldName": "field_name", "fieldType": "string", "expectedValue": "020f6669656c642e312e76616c75652e310f6669656c642e312e76616c75652e32", "json": "{\"field_name\": [\"field.1.value.1\",\"field.1.value.2\"]}", "isOptional": false, "isArray": true, "expectedError": nil, "writer": new(bytes.Buffer)},
		{"caseName": "expected array got string", "fieldName": "field_name", "fieldType": "string", "expectedValue": "", "json": "{\"field_name\": \"field.1.value.1\"}", "isOptional": false, "isArray": true, "expectedError": fmt.Errorf("encode field: expected array for field [field_name] got [String]"), "writer": new(bytes.Buffer)},
	}

	for _, c := range testCases {
		caseName := c["caseName"].(string)
		t.Run(caseName, func(t *testing.T) {
			buf := c["writer"].(mockWriterable)
			encoder := rlp.NewEncoder(buf)

			abi := AbiDef{}

			json := c["json"].(string)
			fieldName := c["fieldName"].(string)
			fieldType := c["fieldType"].(string)
			isOptional := c["isOptional"].(bool)
			isArray := c["isArray"].(bool)
			expectedError := c["expectedError"]

			err := abi.encodeField(encoder, fieldName, fieldType, isOptional, isArray, []byte(json))
			assert.Equal(t, expectedError, err, caseName)

			if c["expectedError"] == nil {
				assert.Equal(t, c["expectedValue"], hex.EncodeToString(buf.Bytes()), c["caseName"])
			}

		})

	}
}

func TestABI_Write(t *testing.T) {
	testCases := []map[string]interface{}{
		{"caseName": "string", "typeName": "string", "expectedValue": "0e746869732e69732e612e74657374", "json": "{\"testField\":\"this.is.a.test\""},
		{"caseName": "min int8", "typeName": "int8", "expectedValue": "80", "json": "{\"testField\":-128}"},
		{"caseName": "max int8", "typeName": "int8", "expectedValue": "7f", "json": "{\"testField\":127}", "expectedError": nil},
		{"caseName": "out of range int8", "typeName": "int8", "expectedValue": "", "json": "{\"testField\":128}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int8 : strconv.ParseInt: parsing \"128\": value out of range")},
		{"caseName": "out of range int8", "typeName": "int8", "expectedValue": "", "json": "{\"testField\":-129}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int8 : strconv.ParseInt: parsing \"-129\": value out of range")},
		{"caseName": "min uint8", "typeName": "uint8", "expectedValue": "00", "json": "{\"testField\":0}", "expectedError": nil},
		{"caseName": "max uint8", "typeName": "uint8", "expectedValue": "ff", "json": "{\"testField\":255}", "expectedError": nil},
		{"caseName": "out of range uint8", "typeName": "uint8", "expectedValue": "", "json": "{\"testField\":-1}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint8 : strconv.ParseUint: parsing \"-1\": invalid syntax")},
		{"caseName": "out of range uint8", "typeName": "uint8", "expectedValue": "", "json": "{\"testField\":256}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint8 : strconv.ParseUint: parsing \"256\": value out of range")},
		{"caseName": "min int16", "typeName": "int16", "expectedValue": "0080", "json": "{\"testField\":-32768}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max int16", "typeName": "int16", "expectedValue": "ff7f", "json": "{\"testField\":32767}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "out of range int16", "typeName": "int16", "expectedValue": "", "json": "{\"testField\":-32769}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int16 : strconv.ParseInt: parsing \"-32769\": value out of range")},
		{"caseName": "out of range int16", "typeName": "int16", "expectedValue": "", "json": "{\"testField\":32768}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int16 : strconv.ParseInt: parsing \"32768\": value out of range")},
		{"caseName": "min uint16", "typeName": "uint16", "expectedValue": "0000", "json": "{\"testField\":0}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max uint16", "typeName": "uint16", "expectedValue": "ffff", "json": "{\"testField\":65535}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "out of range uint16", "typeName": "uint16", "expectedValue": "", "json": "{\"testField\":-1}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint16 : strconv.ParseUint: parsing \"-1\": invalid syntax")},
		{"caseName": "out of range uint16", "typeName": "uint16", "expectedValue": "", "json": "{\"testField\":65536}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint16 : strconv.ParseUint: parsing \"65536\": value out of range")},
		{"caseName": "min int32", "typeName": "int32", "expectedValue": "00000080", "json": "{\"testField\":-2147483648}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max int32", "typeName": "int32", "expectedValue": "ffffff7f", "json": "{\"testField\":2147483647}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "out of range int32", "typeName": "int32", "expectedValue": "", "json": "{\"testField\":-2147483649}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int32 : strconv.ParseInt: parsing \"-2147483649\": value out of range")},
		{"caseName": "out of range int32", "typeName": "int32", "expectedValue": "", "json": "{\"testField\":2147483648}", "expectedError": fmt.Errorf("writing field: [test_field_name] type int32 : strconv.ParseInt: parsing \"2147483648\": value out of range")},
		{"caseName": "min uint32", "typeName": "uint32", "expectedValue": "00000000", "json": "{\"testField\":0}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max uint32", "typeName": "uint32", "expectedValue": "ffffffff", "json": "{\"testField\":4294967295}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "out of range uint32", "typeName": "uint32", "expectedValue": "", "json": "{\"testField\":-1}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint32 : strconv.ParseUint: parsing \"-1\": invalid syntax")},
		{"caseName": "out of range uint32", "typeName": "uint32", "expectedValue": "", "json": "{\"testField\":4294967296}", "expectedError": fmt.Errorf("writing field: [test_field_name] type uint32 : strconv.ParseUint: parsing \"4294967296\": value out of range")},
		{"caseName": "min int64", "typeName": "int64", "expectedValue": "0000000000000080", "json": "{\"testField\":\"-9223372036854775808\"}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "mid int64", "typeName": "int64", "expectedValue": "00f0ffffffffffff", "json": "{\"testField\":-4096}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max int64", "typeName": "int64", "expectedValue": "ffffffffffffff7f", "json": "{\"testField\":\"9223372036854775807\"}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "out of range int64 lower", "typeName": "int64", "expectedValue": "", "json": "{\"testField\":-9223372036854775809}", "expectedError": fmt.Errorf("encoding int64: json: cannot unmarshal number -9223372036854775809 into Go value of type int64")},
		{"caseName": "out of range int64 upper", "typeName": "int64", "expectedValue": "", "json": "{\"testField\":9223372036854775808}", "expectedError": fmt.Errorf("encoding int64: json: cannot unmarshal number 9223372036854775808 into Go value of type int64")},
		{"caseName": "min uint64", "typeName": "uint64", "expectedValue": "0000000000000000", "json": "{\"testField\":0}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "mid uint64", "typeName": "uint64", "expectedValue": "c06ddb095f285813", "json": "{\"testField\":\"1393908473323548096\"}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "max uint64", "typeName": "uint64", "expectedValue": "ffffffffffffffff", "json": "{\"testField\":\"18446744073709551615\"}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "out of range uint64 lower", "typeName": "uint64", "expectedValue": "", "json": "{\"testField\":-1}", "expectedError": fmt.Errorf("encoding uint64: json: cannot unmarshal number -1 into Go value of type uint64")},
		{"caseName": "out of range uint64 upper", "typeName": "uint64", "expectedValue": "", "json": "{\"testField\":18446744073709551616}", "expectedError": fmt.Errorf("encoding uint64: json: cannot unmarshal number 18446744073709551616 into Go value of type uint64")},
		{"caseName": "int128", "typeName": "int128", "expectedValue": "01020000000000000200000000000000", "json": "{\"testField\":\"0x01020000000000000200000000000000\"}"},
		{"caseName": "uint128", "typeName": "uint128", "expectedValue": "01000000000000000200000000000000", "json": "{\"testField\":\"0x01000000000000000200000000000000\"}"},
		{"caseName": "varint32", "typeName": "varint32", "expectedValue": "00000080", "json": "{\"testField\":-2147483648}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "varuint32", "typeName": "varuint32", "expectedValue": "ffffffff", "json": "{\"testField\":4294967295}", "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"}, //{"caseName": "min varuint32", "typeName": "varuint32", "expectedValue": "0", "json": Varuint32(0), "expectedError": nil, "isOptional": false, "isArray": false, "fieldName": "testedField"},
		{"caseName": "min float32", "typeName": "float32", "expectedValue": "01000000", "json": "{\"testField\":0.000000000000000000000000000000000000000000001401298464324817}", "expectedError": nil},
		{"caseName": "max float32", "typeName": "float32", "expectedValue": "ffff7f7f", "json": "{\"testField\":340282346638528860000000000000000000000}", "expectedError": nil},
		{"caseName": "err float32", "typeName": "float32", "expectedValue": "ffff7f7f", "json": "{\"testField\":440282346638528860000000000000000000000}", "expectedError": fmt.Errorf("writing field: [test_field_name] type float32 : strconv.ParseFloat: parsing \"440282346638528860000000000000000000000\": value out of range")},
		{"caseName": "min float64", "typeName": "float64", "expectedValue": "0100000000000000", "json": "{\"testField\":0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005}", "expectedError": nil},
		{"caseName": "max float64", "typeName": "float64", "expectedValue": "ffffffffffffef7f", "json": "{\"testField\":179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000}", "expectedError": nil},
		{"caseName": "err float64", "typeName": "float64", "expectedValue": "ffffffffffffef7f", "json": "{\"testField\":279769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000}", "expectedError": fmt.Errorf("writing field: [test_field_name] type float64 : strconv.ParseFloat: parsing \"279769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\": value out of range")},
		{"caseName": "float128", "typeName": "float128", "expectedValue": "ffffffffffffef7fffffffffffffef7f", "json": "{\"testField\":\"0xffffffffffffef7fffffffffffffef7f\"}"},
		{"caseName": "bool true", "typeName": "bool", "expectedValue": "01", "json": "{\"testField\":true}", "expectedError": nil},
		{"caseName": "bool false", "typeName": "bool", "expectedValue": "00", "json": "{\"testField\":false}", "expectedError": nil},
		{"caseName": "time_point", "typeName": "time_point", "expectedValue": "0100000000000000", "json": "{\"testField\":\"1970-01-01T00:00:00.001\"", "expectedError": nil},
		{"caseName": "time_point err", "typeName": "time_point", "expectedValue": "0100000000000000", "json": "{\"testField\":\"bad.date\"", "expectedError": fmt.Errorf("writing field: time_point: parsing time \"bad.date\" as \"2006-01-02T15:04:05.999\": cannot parse \"bad.date\" as \"2006\"")},
		{"caseName": "time_point_sec", "typeName": "time_point_sec", "expectedValue": "3be10b5e", "json": "{\"testField\":\"2020-01-01T00:00:59\"", "expectedError": nil},
		{"caseName": "time_point_sec err", "typeName": "time_point_sec", "expectedValue": "01000000", "json": "{\"testField\":\"bad date\"", "expectedError": fmt.Errorf("writing field: time_point_sec: parsing time \"bad date\" as \"2006-01-02T15:04:05\": cannot parse \"bad date\" as \"2006\"")},
		{"caseName": "block_timestamp_type", "typeName": "block_timestamp_type", "expectedValue": "368d2223", "json": "{\"testField\":\"2018-09-05T12:48:54.000\"}", "expectedError": nil},
		{"caseName": "block_timestamp_type err", "typeName": "block_timestamp_type", "expectedValue": "76c52223", "json": "{\"testField\":\"this is not a date\"}", "expectedError": fmt.Errorf("writing field: block_timestamp_type: parsing time \"this is not a date\" as \"2006-01-02T15:04:05.000\": cannot parse \"this is not a date\" as \"2006\"")},
		{"caseName": "Name", "typeName": "name", "expectedValue": "0000000000ea3055", "json": "{\"testField\":\"eosio\"}", "expectedError": nil},
		{"caseName": "Name", "typeName": "name", "expectedValue": "", "json": "{\"testField\":\"waytolongnametomakethetestcrash\"}", "expectedError": fmt.Errorf("writing field: name: waytolongnametomakethetestcrash is to long. expected length of max 13 characters")}, //12
		{"caseName": "bytes", "typeName": "bytes", "expectedValue": "0e746869732e69732e612e74657374", "json": "{\"testField\":\"746869732e69732e612e74657374\"}", "expectedError": nil},
		{"caseName": "bytes err", "typeName": "bytes", "expectedValue": "0e746869732e69732e612e74657374", "json": "{\"testField\":\"those are not bytes\"}", "expectedError": fmt.Errorf("writing field: bytes: encoding/hex: invalid byte: U+0074 't'")},
		{"caseName": "checksum160", "typeName": "checksum160", "expectedValue": "0000000000000000000000000000000000000000", "json": "{\"testField\":\"0000000000000000000000000000000000000000\"}", "expectedError": nil},
		{"caseName": "checksum256", "typeName": "checksum256", "expectedValue": "0000000000000000000000000000000000000000000000000000000000000000", "json": "{\"testField\":\"0000000000000000000000000000000000000000000000000000000000000000\"}", "expectedError": nil},
		{"caseName": "checksum512", "typeName": "checksum512", "expectedValue": "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "json": "{\"testField\":\"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\"}", "expectedError": nil},
		{"caseName": "checksum160 to long", "typeName": "checksum160", "expectedValue": "", "json": "{\"testField\":\"10000000000000000000000000000000000000000\"}", "expectedError": fmt.Errorf("writing field: checksum160: expected length of 40 got 41 for value 10000000000000000000000000000000000000000")},
		{"caseName": "checksum256 to long", "typeName": "checksum256", "expectedValue": "", "json": "{\"testField\":\"10000000000000000000000000000000000000000000000000000000000000000\"}", "expectedError": fmt.Errorf("writing field: checksum256: expected length of 64 got 65 for value 10000000000000000000000000000000000000000000000000000000000000000")},
		{"caseName": "checksum512 to long", "typeName": "checksum512", "expectedValue": "", "json": "{\"testField\":\"100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\"}", "expectedError": fmt.Errorf("writing field: checksum512: expected length of 128 got 129 for value 100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")},
		{"caseName": "checksum160 hex err", "typeName": "checksum160", "expectedValue": "", "json": "{\"testField\":\"BADX000000000000000000000000000000000000\"}", "expectedError": fmt.Errorf("writing field: checksum160: encoding/hex: invalid byte: U+0058 'X'")},
		{"caseName": "checksum256 hex err", "typeName": "checksum256", "expectedValue": "", "json": "{\"testField\":\"BADX000000000000000000000000000000000000000000000000000000000000\"}", "expectedError": fmt.Errorf("writing field: checksum256: encoding/hex: invalid byte: U+0058 'X'")},
		{"caseName": "checksum512 hex err", "typeName": "checksum512", "expectedValue": "", "json": "{\"testField\":\"BADX0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\"}", "expectedError": fmt.Errorf("writing field: checksum512: encoding/hex: invalid byte: U+0058 'X'")},
		{"caseName": "public_key", "typeName": "public_key", "expectedValue": "00000000000000000000000000000000000000000000000000000000000000000000", "json": "{\"testField\":\"EOS1111111111111111111111111111111114T1Anm\"}", "expectedError": nil},
		{"caseName": "public_key err", "typeName": "public_key", "expectedValue": "", "json": "{\"testField\":\"EOS1111111111111111111111114T1Anm\"}", "expectedError": fmt.Errorf("writing field: public_key: checkDecode: invalid checksum")},
		{"caseName": "signature", "typeName": "signature", "expectedValue": "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "json": "{\"testField\":\"SIG_K1_111111111111111111111111111111111111111111111111111111111111111116uk5ne\"}", "expectedError": nil},
		{"caseName": "signature err", "typeName": "signature", "expectedValue": "", "json": "{\"testField\":\"SIG_K1_BADX11111111111111111111111111111111111111111111111111111111111116uk5ne\"}", "expectedError": fmt.Errorf("writing field: public_key: signature checksum failed, found 3aea1e96 expected e72f76ff")},
		{"caseName": "symbol", "typeName": "symbol", "expectedValue": "0403454f53", "json": "{\"testField\":\"4,EOS\"}", "expectedError": nil},
		{"caseName": "symbol format error", "typeName": "symbol", "expectedValue": "", "json": "{\"testField\":\"4EOS\"}", "expectedError": fmt.Errorf("writing field: symbol: symbol should be of format '4,EOS'")},
		{"caseName": "symbol format error", "typeName": "symbol", "expectedValue": "", "json": "{\"testField\":\"abc,EOS\"}", "expectedError": fmt.Errorf("writing field: symbol: strconv.ParseUint: parsing \"abc\": invalid syntax")},
		{"caseName": "symbol_code", "typeName": "symbol_code", "expectedValue": "ffffffffffffffff", "json": "{\"testField\":18446744073709551615}", "expectedError": nil},
		{"caseName": "asset", "typeName": "asset", "expectedValue": "a08601000000000004454f5300000000", "json": "{\"testField\":\"10.0000 EOS\"}", "expectedError": nil},
		{"caseName": "asset err", "typeName": "asset", "expectedValue": "", "json": "{\"testField\":\"AA.0000 EOS\"}", "expectedError": fmt.Errorf("writing field: asset: strconv.ParseInt: parsing \"AA0000\": invalid syntax")},
		{"caseName": "extended_asset", "typeName": "extended_asset", "expectedValue": "0a0000000000000004454f5300000000202932c94c833055", "json": "{\"testField\":{\"asset\":\"0.0010 EOS\",\"Contract\":\"eoscanadacom\"}}", "expectedError": nil},
		{"caseName": "extended_asset err", "typeName": "extended_asset", "expectedValue": "", "json": "{\"testField\":{\"asset\":\"abc.0010 EOS\",\"Contract\":\"eoscanadacom\"}}", "expectedError": fmt.Errorf("writing field: extended_asset: strconv.ParseInt: parsing \"abc0010\": invalid syntax")},
		{"caseName": "bad type", "typeName": "bad.type.1", "expectedValue": nil, "json": "{\"testField\":0}", "expectedError": fmt.Errorf("writing field of type [bad.type.1]: unknown type")},
		{"caseName": "optional present", "typeName": "string", "expectedValue": "0776616c75652e31", "json": "{\"testField\":\"value.1\"}", "expectedError": nil},
		{"caseName": "struct", "typeName": "struct_name_1", "expectedValue": "0e746869732e69732e612e74657374", "json": "{\"testField\": {\"field_name_1\":\"this.is.a.test\"}}", "expectedError": nil},
		{"caseName": "struct err", "typeName": "struct_name_1", "expectedValue": "0e746869732e69732e612e74657374", "json": "{\"testField\": {}", "expectedError": fmt.Errorf("encoding fields: encode field: none optional field [field_name_1] as a nil value")},
	}

	for _, c := range testCases {

		t.Run(c["caseName"].(string), func(t *testing.T) {
			var buffer bytes.Buffer
			encoder := rlp.NewEncoder(&buffer)

			abi := AbiDef{
				Structs: []StructDef{
					{
						Name:   "struct_name_1",
						Base:   "",
						Fields: []FieldDef{{Name: "field_name_1", Type: "string"}},
					},
				},
			}
			fieldName := "test_field_name"
			result := gjson.Get(c["json"].(string), "testField")
			err := abi.writeField(encoder, fieldName, c["typeName"].(string), result)
			if err != nil {
				fmt.Println(err.Error())
			}
			assert.Equal(t, c["expectedError"], err, c["caseName"])

			if c["expectedError"] == nil {
				assert.Equal(t, c["expectedValue"], hex.EncodeToString(buffer.Bytes()), c["caseName"])
			}
		})
	}
}

type mockWriterable interface {
	Write(p []byte) (n int, err error)
	Bytes() []byte
}
type mockWriter struct {
	length int
	err    error
}

func (w mockWriter) Write(p []byte) (n int, err error) {
	return w.length, w.err
}

func (w mockWriter) Bytes() []byte {
	return []byte{}
}

func TestMyabi(t *testing.T) {
	var abidef AbiDef
	ToABI([]byte(abiString), &abidef)
	fmt.Println("abiDef:  ", abidef)
}

var voteData = common.Variants{
	"voter":     common.N("eosio.token"),
	"proxy":     common.N("walker"),
	"producers": []common.AccountName{common.N("aaa"), common.N("hello")},
}

func TestVoteProducers(t *testing.T) {
	abi, err := NewABI(strings.NewReader(systemABIString))
	assert.NoError(t, err)
	re, err := abi.EncodeAction(common.AccountName(common.N("voteproducer")), abiData)
	fmt.Println(err, re)

}

var systemABIString = `
{
   "version": "eosio::abi/1.0",
   "types": [{
      "new_type_name": "account_name",
      "type": "name"
   },{
      "new_type_name": "permission_name",
      "type": "name"
   },{
      "new_type_name": "action_name",
      "type": "name"
   },{
      "new_type_name": "transaction_id_type",
      "type": "checksum256"
   },{
      "new_type_name": "weight_type",
      "type": "uint16"
   }],
   "____comment": "eosio.bios structs: set_account_limits, setpriv, set_global_limits, producer_key, set_producers, require_auth are provided so abi available for deserialization in future.",
   "structs": [{
      "name": "permission_level",
      "base": "",
      "fields": [
        {"name":"actor",      "type":"account_name"},
        {"name":"permission", "type":"permission_name"}
      ]
    },{
      "name": "key_weight",
      "base": "",
      "fields": [
        {"name":"key",    "type":"public_key"},
        {"name":"weight", "type":"weight_type"}
      ]
    },{
      "name": "bidname",
      "base": "",
      "fields": [
        {"name":"bidder",  "type":"account_name"},
        {"name":"newname", "type":"account_name"},
        {"name":"bid", "type":"asset"}
      ]
    },{
      "name": "permission_level_weight",
      "base": "",
      "fields": [
        {"name":"permission", "type":"permission_level"},
        {"name":"weight",     "type":"weight_type"}
      ]
    },{
      "name": "wait_weight",
      "base": "",
      "fields": [
        {"name":"wait_sec", "type":"uint32"},
        {"name":"weight",   "type":"weight_type"}
      ]
    },{
      "name": "authority",
      "base": "",
      "fields": [
        {"name":"threshold", "type":"uint32"},
        {"name":"keys",      "type":"key_weight[]"},
        {"name":"accounts",  "type":"permission_level_weight[]"},
        {"name":"waits",     "type":"wait_weight[]"}
      ]
    },{
      "name": "newaccount",
      "base": "",
      "fields": [
        {"name":"creator", "type":"account_name"},
        {"name":"name",    "type":"account_name"},
        {"name":"owner",   "type":"authority"},
        {"name":"active",  "type":"authority"}
      ]
    },{
      "name": "setcode",
      "base": "",
      "fields": [
        {"name":"account",   "type":"account_name"},
        {"name":"vmtype",    "type":"uint8"},
        {"name":"vmversion", "type":"uint8"},
        {"name":"code",      "type":"bytes"}
      ]
    },{
      "name": "setabi",
      "base": "",
      "fields": [
        {"name":"account", "type":"account_name"},
        {"name":"abi",     "type":"bytes"}
      ]
    },{
      "name": "updateauth",
      "base": "",
      "fields": [
        {"name":"account",    "type":"account_name"},
        {"name":"permission", "type":"permission_name"},
        {"name":"parent",     "type":"permission_name"},
        {"name":"auth",       "type":"authority"}
      ]
    },{
      "name": "deleteauth",
      "base": "",
      "fields": [
        {"name":"account",    "type":"account_name"},
        {"name":"permission", "type":"permission_name"}
      ]
    },{
      "name": "linkauth",
      "base": "",
      "fields": [
        {"name":"account",     "type":"account_name"},
        {"name":"code",        "type":"account_name"},
        {"name":"type",        "type":"action_name"},
        {"name":"requirement", "type":"permission_name"}
      ]
    },{
      "name": "unlinkauth",
      "base": "",
      "fields": [
        {"name":"account",     "type":"account_name"},
        {"name":"code",        "type":"account_name"},
        {"name":"type",        "type":"action_name"}
      ]
    },{
      "name": "canceldelay",
      "base": "",
      "fields": [
        {"name":"canceling_auth", "type":"permission_level"},
        {"name":"trx_id",         "type":"transaction_id_type"}
      ]
    },{
      "name": "onerror",
      "base": "",
      "fields": [
        {"name":"sender_id", "type":"uint128"},
        {"name":"sent_trx",  "type":"bytes"}
      ]
    },{
      "name": "buyrambytes",
      "base": "",
      "fields": [
         {"name":"payer", "type":"account_name"},
         {"name":"receiver", "type":"account_name"},
         {"name":"bytes", "type":"uint32"}
      ]
    },{
      "name": "sellram",
      "base": "",
      "fields": [
         {"name":"account", "type":"account_name"},
         {"name":"bytes", "type":"uint64"}
      ]
    },{
      "name": "buyram",
      "base": "",
      "fields": [
         {"name":"payer", "type":"account_name"},
         {"name":"receiver", "type":"account_name"},
         {"name":"quant", "type":"asset"}
      ]
    },{
      "name": "delegatebw",
      "base": "",
      "fields": [
         {"name":"from", "type":"account_name"},
         {"name":"receiver", "type":"account_name"},
         {"name":"stake_net_quantity", "type":"asset"},
         {"name":"stake_cpu_quantity", "type":"asset"},
         {"name":"transfer", "type":"bool"}
      ]
    },{
      "name": "undelegatebw",
      "base": "",
      "fields": [
         {"name":"from", "type":"account_name"},
         {"name":"receiver", "type":"account_name"},
         {"name":"unstake_net_quantity", "type":"asset"},
         {"name":"unstake_cpu_quantity", "type":"asset"}
      ]
    },{
      "name": "refund",
      "base": "",
      "fields": [
         {"name":"owner", "type":"account_name"}
      ]
    },{
      "name": "delegated_bandwidth",
      "base": "",
      "fields": [
         {"name":"from", "type":"account_name"},
         {"name":"to", "type":"account_name"},
         {"name":"net_weight", "type":"asset"},
         {"name":"cpu_weight", "type":"asset"}
      ]
    },{
      "name": "user_resources",
      "base": "",
      "fields": [
         {"name":"owner", "type":"account_name"},
         {"name":"net_weight", "type":"asset"},
         {"name":"cpu_weight", "type":"asset"},
         {"name":"ram_bytes", "type":"uint64"}
      ]
    },{
      "name": "total_resources",
      "base": "",
      "fields": [
         {"name":"owner", "type":"account_name"},
         {"name":"net_weight", "type":"asset"},
         {"name":"cpu_weight", "type":"asset"},
         {"name":"ram_bytes", "type":"uint64"}
      ]
    },{
      "name": "refund_request",
      "base": "",
      "fields": [
         {"name":"owner", "type":"account_name"},
         {"name":"request_time", "type":"time_point_sec"},
         {"name":"net_amount", "type":"asset"},
         {"name":"cpu_amount", "type":"asset"}
      ]
    },{
      "name": "blockchain_parameters",
      "base": "",
      "fields": [

         {"name":"max_block_net_usage",                 "type":"uint64"},
         {"name":"target_block_net_usage_pct",          "type":"uint32"},
         {"name":"max_transaction_net_usage",           "type":"uint32"},
         {"name":"base_per_transaction_net_usage",      "type":"uint32"},
         {"name":"net_usage_leeway",                    "type":"uint32"},
         {"name":"context_free_discount_net_usage_num", "type":"uint32"},
         {"name":"context_free_discount_net_usage_den", "type":"uint32"},
         {"name":"max_block_cpu_usage",                 "type":"uint32"},
         {"name":"target_block_cpu_usage_pct",          "type":"uint32"},
         {"name":"max_transaction_cpu_usage",           "type":"uint32"},
         {"name":"min_transaction_cpu_usage",           "type":"uint32"},
         {"name":"max_transaction_lifetime",            "type":"uint32"},
         {"name":"deferred_trx_expiration_window",      "type":"uint32"},
         {"name":"max_transaction_delay",               "type":"uint32"},
         {"name":"max_inline_action_size",              "type":"uint32"},
         {"name":"max_inline_action_depth",             "type":"uint16"},
         {"name":"max_authority_depth",                 "type":"uint16"}

      ]
    },{
      "name": "eosio_global_state",
      "base": "blockchain_parameters",
      "fields": [
         {"name":"max_ram_size",                  "type":"uint64"},
         {"name":"total_ram_bytes_reserved",      "type":"uint64"},
         {"name":"total_ram_stake",               "type":"int64"},
         {"name":"last_producer_schedule_update", "type":"block_timestamp_type"},
         {"name":"last_pervote_bucket_fill",      "type":"uint64"},
         {"name":"pervote_bucket",                "type":"int64"},
         {"name":"perblock_bucket",               "type":"int64"},
         {"name":"total_unpaid_blocks",           "type":"uint32"},
         {"name":"total_activated_stake",         "type":"int64"},
         {"name":"thresh_activated_stake_time",   "type":"uint64"},
         {"name":"last_producer_schedule_size",   "type":"uint16"},
         {"name":"total_producer_vote_weight",    "type":"float64"},
         {"name":"last_name_close",               "type":"block_timestamp_type"}
      ]
    },{
      "name": "producer_info",
      "base": "",
      "fields": [
         {"name":"owner",           "type":"account_name"},
         {"name":"total_votes",     "type":"float64"},
         {"name":"producer_key",    "type":"public_key"},
         {"name":"is_active",       "type":"bool"},
         {"name":"url",             "type":"string"},
         {"name":"unpaid_blocks",   "type":"uint32"},
         {"name":"last_claim_time", "type":"uint64"},
         {"name":"location",        "type":"uint16"}
      ]
    },{
      "name": "regproducer",
      "base": "",
      "fields": [
        {"name":"producer",     "type":"account_name"},
        {"name":"producer_key", "type":"public_key"},
        {"name":"url",          "type":"string"},
        {"name":"location",     "type":"uint16"}
      ]
    },{
      "name": "unregprod",
      "base": "",
      "fields": [
        {"name":"producer",     "type":"account_name"}
      ]
    },{
      "name": "setram",
      "base": "",
      "fields": [
        {"name":"max_ram_size",     "type":"uint64"}
      ]
    },{
      "name": "regproxy",
      "base": "",
      "fields": [
        {"name":"proxy",     "type":"account_name"},
        {"name":"isproxy",   "type":"bool"}
      ]
    },{
      "name": "voteproducer",
      "base": "",
      "fields": [
        {"name":"voter",     "type":"account_name"},
        {"name":"proxy",     "type":"account_name"},
        {"name":"producers", "type":"account_name[]"}
      ]
    },{
      "name": "voter_info",
      "base": "",
      "fields": [
        {"name":"owner",                "type":"account_name"},
        {"name":"proxy",                "type":"account_name"},
        {"name":"producers",            "type":"account_name[]"},
        {"name":"staked",               "type":"int64"},
        {"name":"last_vote_weight",     "type":"float64"},
        {"name":"proxied_vote_weight",  "type":"float64"},
        {"name":"is_proxy",             "type":"bool"}
      ]
    },{
      "name": "claimrewards",
      "base": "",
      "fields": [
        {"name":"owner",   "type":"account_name"}
      ]
    },{
      "name": "setpriv",
      "base": "",
      "fields": [
        {"name":"account",    "type":"account_name"},
        {"name":"is_priv",    "type":"int8"}
      ]
    },{
      "name": "rmvproducer",
      "base": "",
      "fields": [
        {"name":"producer", "type":"account_name"}
      ]
    },{
      "name": "set_account_limits",
      "base": "",
      "fields": [
        {"name":"account",    "type":"account_name"},
        {"name":"ram_bytes",  "type":"int64"},
        {"name":"net_weight", "type":"int64"},
        {"name":"cpu_weight", "type":"int64"}
      ]
    },{
      "name": "set_global_limits",
      "base": "",
      "fields": [
        {"name":"cpu_usec_per_period",    "type":"int64"}
      ]
    },{
      "name": "producer_key",
      "base": "",
      "fields": [
        {"name":"producer_name",      "type":"account_name"},
        {"name":"block_signing_key",  "type":"public_key"}
      ]
    },{
      "name": "set_producers",
      "base": "",
      "fields": [
        {"name":"schedule",   "type":"producer_key[]"}
      ]
    },{
      "name": "require_auth",
      "base": "",
      "fields": [
        {"name":"from", "type":"account_name"}
      ]
    },{
      "name": "setparams",
      "base": "",
      "fields": [
        {"name":"params", "type":"blockchain_parameters"}
      ]
    },{
      "name": "connector",
      "base": "",
      "fields": [
        {"name":"balance", "type":"asset"},
        {"name":"weight", "type":"float64"}
      ]
    },{
      "name": "exchange_state",
      "base": "",
      "fields": [
        {"name":"supply", "type":"asset"},
        {"name":"base", "type":"connector"},
        {"name":"quote", "type":"connector"}
      ]
    }, {
       "name": "namebid_info",
       "base": "",
       "fields": [
          {"name":"newname", "type":"account_name"},
          {"name":"high_bidder", "type":"account_name"},
          {"name":"high_bid", "type":"int64"},
          {"name":"last_bid_time", "type":"uint64"}
       ]
    }
   ],
   "actions": [{
     "name": "newaccount",
     "type": "newaccount",
     "ricardian_contract": ""
   },{
     "name": "setcode",
     "type": "setcode",
     "ricardian_contract": ""
   },{
     "name": "setabi",
     "type": "setabi",
     "ricardian_contract": ""
   },{
     "name": "updateauth",
     "type": "updateauth",
     "ricardian_contract": ""
   },{
     "name": "deleteauth",
     "type": "deleteauth",
     "ricardian_contract": ""
   },{
     "name": "linkauth",
     "type": "linkauth",
     "ricardian_contract": ""
   },{
     "name": "unlinkauth",
     "type": "unlinkauth",
     "ricardian_contract": ""
   },{
     "name": "canceldelay",
     "type": "canceldelay",
     "ricardian_contract": ""
   },{
     "name": "onerror",
     "type": "onerror",
     "ricardian_contract": ""
   },{
      "name": "buyrambytes",
      "type": "buyrambytes",
      "ricardian_contract": ""
   },{
      "name": "buyram",
      "type": "buyram",
      "ricardian_contract": ""
   },{
      "name": "sellram",
      "type": "sellram",
      "ricardian_contract": ""
   },{
      "name": "delegatebw",
      "type": "delegatebw",
      "ricardian_contract": ""
   },{
      "name": "undelegatebw",
      "type": "undelegatebw",
      "ricardian_contract": ""
   },{
      "name": "refund",
      "type": "refund",
      "ricardian_contract": ""
   },{
      "name": "regproducer",
      "type": "regproducer",
      "ricardian_contract": ""
   },{
      "name": "setram",
      "type": "setram",
      "ricardian_contract": ""
   },{
      "name": "bidname",
      "type": "bidname",
      "ricardian_contract": ""
   },{
      "name": "unregprod",
      "type": "unregprod",
      "ricardian_contract": ""
   },{
      "name": "regproxy",
      "type": "regproxy",
      "ricardian_contract": ""
   },{
      "name": "voteproducer",
      "type": "voteproducer",
      "ricardian_contract": ""
   },{
      "name": "claimrewards",
      "type": "claimrewards",
      "ricardian_contract": ""
   },{
      "name": "setpriv",
      "type": "setpriv",
      "ricardian_contract": ""
   },{
      "name": "rmvproducer",
      "type": "rmvproducer",
      "ricardian_contract": ""
   },{
      "name": "setalimits",
      "type": "set_account_limits",
      "ricardian_contract": ""
    },{
      "name": "setglimits",
      "type": "set_global_limits",
      "ricardian_contract": ""
    },{
      "name": "setprods",
      "type": "set_producers",
      "ricardian_contract": ""
    },{
      "name": "reqauth",
      "type": "require_auth",
      "ricardian_contract": ""
    },{
      "name": "setparams",
      "type": "setparams",
      "ricardian_contract": ""
    }],
   "tables": [{
      "name": "producers",
      "type": "producer_info",
      "index_type": "i64",
      "key_names" : ["owner"],
      "key_types" : ["uint64"]
    },{
      "name": "global",
      "type": "eosio_global_state",
      "index_type": "i64",
      "key_names" : [],
      "key_types" : []
    },{
      "name": "voters",
      "type": "voter_info",
      "index_type": "i64",
      "key_names" : ["owner"],
      "key_types" : ["account_name"]
    },{
      "name": "userres",
      "type": "user_resources",
      "index_type": "i64",
      "key_names" : ["owner"],
      "key_types" : ["uint64"]
    },{
      "name": "delband",
      "type": "delegated_bandwidth",
      "index_type": "i64",
      "key_names" : ["to"],
      "key_types" : ["uint64"]
    },{
      "name": "rammarket",
      "type": "exchange_state",
      "index_type": "i64",
      "key_names" : ["supply"],
      "key_types" : ["uint64"]
    },{
      "name": "refunds",
      "type": "refund_request",
      "index_type": "i64",
      "key_names" : ["owner"],
      "key_types" : ["uint64"]
    },{
       "name": "namebids",
       "type": "namebid_info",
       "index_type": "i64",
       "key_names" : ["newname"],
       "key_types" : ["account_name"]
    }
   ],
   "ricardian_clauses": [],
   "abi_extensions": []
}
`
