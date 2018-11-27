package abi_serializer

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"io"
	"strconv"
	"strings"
	"time"
)

type ABIEncoder struct {
	abiReader  io.Reader
	eosEncoder *rlp.Encoder
	abi        *AbiDef
	pos        int
}

func (a *AbiDef) EncodeAction(actionName common.ActionName, json []byte) ([]byte, error) {

	action := a.ActionForName(actionName)
	if action == nil {
		return nil, fmt.Errorf("encode action: action %s not found in abi", actionName)
	}

	var buffer bytes.Buffer
	encoder := rlp.NewEncoder(&buffer)

	err := a.encode(encoder, action.Type, json)
	if err != nil {
		return nil, fmt.Errorf("encode action: %s", err)
	}
	return buffer.Bytes(), nil
}

func (a *AbiDef) encode(binaryEncoder *rlp.Encoder, structureName string, json []byte) error {
	abiLog.Debug("abi encode struct %s", structureName)

	structure := a.StructForName(structureName)
	if structure == nil {
		return fmt.Errorf("encode struct [%s] not found in abi", structureName)
	}

	if structure.Base != "" {
		abiLog.Debug("struct has base struct %s : %s", structureName, structure.Base)
		err := a.encode(binaryEncoder, structure.Base, json)
		if err != nil {
			return fmt.Errorf("encode base [%s]: %s", structureName, err)
		}
	}
	err := a.encodeFields(binaryEncoder, structure.Fields, json)
	return err
}

func (a *AbiDef) encodeFields(binaryEncoder *rlp.Encoder, fields []FieldDef, json []byte) error {
	for _, field := range fields {

		abiLog.Debug("encode field: name: %s, type: %s", field.Name, field.Type)

		fieldType, isOptional, isArray := analyzeFieldType(field.Type)
		typeName := a.TypeNameForNewTypeName(fieldType)
		fieldName := field.Name
		if typeName != field.Type {
			abiLog.Debug("type is an alias, from  %s to %s", field.Type, typeName)
		}

		err := a.encodeField(binaryEncoder, fieldName, typeName, isOptional, isArray, json)
		if err != nil {
			return fmt.Errorf("encoding fields: %s", err)
		}
	}
	return nil
}

func (a *AbiDef) encodeField(binaryEncoder *rlp.Encoder, fieldName string, fieldType string, isOptional bool, isArray bool, json []byte) (err error) {

	abiLog.Debug("encode field json : %#v", json)

	value := gjson.GetBytes(json, fieldName)
	if isOptional {
		if value.Exists() {
			abiLog.Debug("field is optional and present, name : %s, type: %s", fieldName, fieldType)
			if e := binaryEncoder.WriteByte(1); e != nil {
				return e
			}
		} else {
			abiLog.Debug("field is optional and *not* present", zap.String("name", fieldName), zap.String("type", fieldType))
			return binaryEncoder.WriteByte(0)
		}

	} else if !value.Exists() {
		return fmt.Errorf("encode field: none optional field [%s] as a nil value", fieldName)
	}

	if isArray {

		abiLog.Debug("field is an array, name is %s,type is %s", fieldName, fieldType)
		if !value.IsArray() {
			return fmt.Errorf("encode field: expected array for field [%s] got [%s]", fieldName, value.Type.String())
		}

		results := value.Array()
		binaryEncoder.WriteUVarInt(len(results))

		for _, r := range results {
			a.writeField(binaryEncoder, fieldName, fieldType, r)
		}

		return nil
	}

	return a.writeField(binaryEncoder, fieldName, fieldType, value)
}

func (a *AbiDef) writeField(binaryEncoder *rlp.Encoder, fieldName string, fieldType string, value gjson.Result) error {

	abiLog.Debug("write field, name is %s, type is %s,json is %s", fieldName, fieldType, value.Raw)

	structure := a.StructForName(fieldType)
	if structure != nil {
		abiLog.Debug("field is a struct, name is %s", fieldName)

		err := a.encodeFields(binaryEncoder, structure.Fields, []byte(value.Raw))
		if err != nil {
			return err
		}
		return nil
	}

	var object interface{}
	switch fieldType {
	case "int8":
		i, err := valueToInt(fieldName, value, 8)
		if err != nil {
			return err
		}
		object = int8(i)
	case "uint8":
		i, err := valueToUint(fieldName, value, 8)
		if err != nil {
			return err
		}
		object = uint8(i)
	case "int16":
		i, err := valueToInt(fieldName, value, 16)
		if err != nil {
			return err
		}
		object = int16(i)
	case "uint16":
		i, err := valueToUint(fieldName, value, 16)
		if err != nil {
			return err
		}
		object = uint16(i)
	case "int32", "varint32":
		i, err := valueToInt(fieldName, value, 32)
		if err != nil {
			return err
		}
		object = int32(i)
	case "uint32", "varuint32":
		i, err := valueToUint(fieldName, value, 32)
		if err != nil {
			return err
		}
		object = uint32(i)
	case "int64":
		var in Int64
		if err := json.Unmarshal([]byte(value.Raw), &in); err != nil {
			return fmt.Errorf("encoding int64: %s", err)
		}
		object = in
	case "uint64":
		var in Uint64
		if err := json.Unmarshal([]byte(value.Raw), &in); err != nil {
			return fmt.Errorf("encoding uint64: %s", err)
		}
		object = in
	case "int128":
		var in Int128
		if err := json.Unmarshal([]byte(value.Raw), &in); err != nil {
			return err
		}
		object = in
	case "uint128":
		var in Uint128
		if err := json.Unmarshal([]byte(value.Raw), &in); err != nil {
			return err
		}
		object = in
	case "float32":
		f, err := valueToFloat(fieldName, value, 32)
		if err != nil {
			return err
		}
		object = float32(f)
	case "float64":
		f, err := valueToFloat(fieldName, value, 64)
		if err != nil {
			return err
		}
		object = f
	case "float128":
		var in Float128
		if err := json.Unmarshal([]byte(value.Raw), &in); err != nil {
			return err
		}
		object = in
	case "bool":
		object = value.Bool()
	case "time_point_sec":
		t, err := time.Parse("2006-01-02T15:04:05", value.Str)
		if err != nil {
			return fmt.Errorf("writing field: time_point_sec: %s", err)
		}
		object = common.TimePointSec(t.UTC().Second())
	case "time_point":
		t, err := time.Parse("2006-01-02T15:04:05.999", value.Str)
		if err != nil {
			return fmt.Errorf("writing field: time_point: %s", err)
		}
		object = common.TimePoint(t.UTC().Nanosecond() / int(time.Millisecond))
	case "block_timestamp_type":
		t, err := time.Parse("2006-01-02T15:04:05.999999-07:00", value.Str)
		if err != nil {
			return fmt.Errorf("writing field: block_timestamp_type: %s", err)
		}
		slot := uint32(t.Unix() - 946684800)
		object = types.BlockTimeStamp(slot)
	case "name":
		if len(value.Str) > 12 {
			return fmt.Errorf("writing field: name: %s is to long. expected length of max 12 characters", value.Str)
		}
		object = common.N(value.Str)
	case "bytes":
		data, err := hex.DecodeString(value.String())
		if err != nil {
			return fmt.Errorf("writing field: bytes: %s", err)
		}
		object = data

	case "string":
		object = value.String()
	case "checksum160":
		if len(value.Str) != 40 {
			return fmt.Errorf("writing field: checksum160: expected length of 40 got %d for value %s", len(value.Str), value.String())
		}
		data, err := hex.DecodeString(value.Str)
		if err != nil {
			return fmt.Errorf("writing field: checksum160: %s", err)
		}
		object = crypto.NewRipemd160Byte(data)
	case "checksum256":
		if len(value.Str) != 64 {
			return fmt.Errorf("writing field: checksum256: expected length of 64 got %d for value %s", len(value.Str), value.String())
		}
		data, err := hex.DecodeString(value.Str)
		if err != nil {
			return fmt.Errorf("writing field: checksum256: %s", err)
		}
		object = crypto.NewSha256Byte(data)
	case "checksum512":
		if len(value.Str) != 128 {
			return fmt.Errorf("writing field: checksum512: expected length of 128 got %d for value %s", len(value.Str), value.String())
		}
		data, err := hex.DecodeString(value.String())
		if err != nil {
			return fmt.Errorf("writing field: checksum512: %s", err)
		}
		object = crypto.NewSha512Byte(data)
	case "public_key":
		pk, err := ecc.NewPublicKey(value.String())
		if err != nil {
			return fmt.Errorf("writing field: public_key: %s", err)
		}
		object = pk
	case "signature":
		signature, err := ecc.NewSignature(value.String())
		if err != nil {
			return fmt.Errorf("writing field: public_key: %s", err)
		}
		object = signature
	case "symbol":
		parts := strings.Split(value.Str, ",")
		if len(parts) != 2 {
			return fmt.Errorf("writing field: symbol: symbol should be of format '4,EOS'")
		}

		i, err := strconv.ParseUint(parts[0], 10, 8)
		if err != nil {
			return fmt.Errorf("writing field: symbol: %s", err)
		}
		object = common.Symbol{ //TODO
			Precision: uint8(i),
			Symbol:    parts[1],
		}
	case "symbol_code":
		//object = SymbolCode(value.Uint())
		object = uint64(value.Uint())
	case "asset":
		asset, err := common.NewAsset(value.String())
		if err != nil {
			return fmt.Errorf("writing field: asset: %s", err)
		}
		object = asset
	case "extended_asset":
		var extendedAsset common.ExtendedAsset
		err := json.Unmarshal([]byte(value.Raw), &extendedAsset)
		if err != nil {
			return fmt.Errorf("writing field: extended_asset: %s", err)
		}
		object = extendedAsset
	default:
		return fmt.Errorf("writing field of type [%s]: unknown type", fieldType)
	}

	abiLog.Debug("write object %#v", object)

	return binaryEncoder.Encode(object)
}

func valueToInt(fieldName string, value gjson.Result, bitSize int) (int64, error) {
	i, err := strconv.ParseInt(value.Raw, 10, bitSize)
	if err != nil {
		return i, fmt.Errorf("writing field: [%s] type int%d : %s", fieldName, bitSize, err)
	}
	return i, nil
}

func valueToUint(fieldName string, value gjson.Result, bitSize int) (uint64, error) {
	i, err := strconv.ParseUint(value.Raw, 10, bitSize)
	if err != nil {
		return i, fmt.Errorf("writing field: [%s] type uint%d : %s", fieldName, bitSize, err)
	}
	return i, nil
}

func valueToFloat(fieldName string, value gjson.Result, bitSize int) (float64, error) {
	f, err := strconv.ParseFloat(value.Raw, bitSize)
	if err != nil {
		return f, fmt.Errorf("writing field: [%s] type float%d : %s", fieldName, bitSize, err)
	}
	return f, nil
}
