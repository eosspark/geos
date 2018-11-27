package abi_serializer

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"strconv"
	"strings"
)

var maxRecursionDepth = 32

func Encode_Decode() common.Pair {
	decode := func() {
	}
	encode := func() {

	}
	return common.MakePair(decode, encode)
}

type AbiSerializer struct {
	typeDefs      map[string]string
	structs       map[string]StructDef
	actions       map[common.Name]string
	tables        map[common.Name]string
	errorMessages map[uint64]string
	variants      map[string]VariantDef
	builtInTypes  map[string]common.Pair
}

//func (a AbiSerializer) ConfigureBuiltInTypes() {
//	a.builtInTypes["bool"]
//}

func (a AbiSerializer) SetAbi(abi *AbiDef, maxSerializationTime *common.Microseconds) {
	deadline := common.Now() + common.TimePoint(*maxSerializationTime)
	a.typeDefs = make(map[string]string)
	a.structs = make(map[string]StructDef)
	a.actions = make(map[common.Name]string)
	a.tables = make(map[common.Name]string)
	a.errorMessages = make(map[uint64]string)

	for _, st := range abi.Structs {
		a.structs[st.Name] = st
	}

	for _, td := range abi.Types {
		try.EosAssert(a.IsType(&td.Type, 0, &deadline, *maxSerializationTime), &exception.InvalidTypeInsideAbi{}, "invalid type : %v", td.Type)
		try.EosAssert(!a.IsType(&td.NewTypeName, 0, &deadline, *maxSerializationTime), &exception.DuplicateAbiTypeDefException{}, "type already exists : %v", td.Type)
		a.typeDefs[td.NewTypeName] = td.Type
	}

	for _, ac := range abi.Actions {
		a.actions[ac.Name] = ac.Type
	}

	for _, t := range abi.Tables {
		a.tables[t.Name] = t.Type
	}

	for _, e := range abi.ErrorMessages {
		a.errorMessages[e.Code] = e.Message
	}

	try.EosAssert(len(a.typeDefs) == len(abi.Types), &exception.DuplicateAbiTypeDefException{}, "duplicate type definition detected")
	try.EosAssert(len(a.structs) == len(abi.Structs), &exception.DuplicateAbiStructDefException{}, "duplicate struct definition detected")
	try.EosAssert(len(a.actions) == len(abi.Actions), &exception.DuplicateAbiActionDefException{}, "duplicate action definition detected")
	try.EosAssert(len(a.tables) == len(abi.Tables), &exception.DuplicateAbiTableDefException{}, "duplicate table definition detected")
	try.EosAssert(len(a.errorMessages) == len(abi.ErrorMessages), &exception.DuplicateAbiErrMsgDefException{}, "duplicate error message definition detected")

	a.validate() //TODO always return true
}

func (a AbiSerializer) IsBuiltinType(rtype *string) bool {
	for p := range a.builtInTypes {
		if p == *rtype {
			return true
		}
	}
	return false
}

func (a AbiSerializer) IsInteger(rtype *string) bool {
	stype := string(*rtype)
	return strings.HasPrefix(stype, "int") || strings.HasPrefix(stype, "uint")
}

func (a AbiSerializer) GetIntegerSize(rtype *string) int {
	stype := string(*rtype)
	try.EosAssert(a.IsInteger(rtype), &exception.InvalidTypeInsideAbi{}, "%v is not an integer type", stype)
	var num int
	if strings.HasPrefix(stype, "uint") {
		num, _ = strconv.Atoi(string([]byte(stype)[4:]))
		return num
	} else {
		num, _ = strconv.Atoi(string([]byte(stype)[3:]))
		return num
	}
}

func (a AbiSerializer) IsStruct(rtype *string) bool {
	for p := range a.structs {
		if p == *rtype {
			return true
		}
	}
	return false
}

func (a AbiSerializer) IsArray(rtype *string) bool {
	//TODO: [] in go is prefix.
	return strings.HasSuffix(string(*rtype), "[]")
}

func (a AbiSerializer) IsOptional(rtype *string) bool {
	return strings.HasPrefix(string(*rtype), "?")
}

func (a AbiSerializer) FundamentalType(rtype *string) string {
	stype := string(*rtype)
	btype := []byte(stype)
	if a.IsArray(rtype) {
		return string(string(btype[0 : len(btype)-2]))
	} else if a.IsOptional(rtype) {
		return string(string(btype[0 : len(btype)-1]))
	} else {
		return *rtype
	}
}

func (a AbiSerializer) IsType(rtype *string, recursionDepth common.SizeT, deadline *common.TimePoint, maxSerializationTime common.Microseconds) bool {
	try.EosAssert(common.Now() < *deadline, &exception.AbiSerializationDeadlineException{}, "serialization time limit %vus exceeded", maxSerializationTime)
	recursionDepth++
	if recursionDepth > maxRecursionDepth {
		return false
	}
	ftype := a.FundamentalType(rtype)
	if a.IsBuiltinType(&ftype) {
		return true
	}

	if a.IsStruct(&ftype) {
		return true
	}
	return true
}

func (a AbiSerializer) validate() bool {
	return true
}

func isEmptyABI(abiVec common.HexBytes) bool {
	return len(abiVec) <= 4
}

func ToABI(abiVec common.HexBytes, abi *AbiDef) bool {
	if isEmptyABI(abiVec) { // 4 == packsize of empty Abi
		return false
	}
	err := rlp.DecodeBytes(abiVec, abi)
	if err != nil {
		return false
	}
	return true
}
