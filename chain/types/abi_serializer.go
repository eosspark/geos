package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/exception"
	"strings"
)

//func Encode_Decode() {
//
//}

type AbiSerializer struct {
	typeDefs      map[TypeName]TypeName
	structs       map[TypeName]StructDef
	actions       map[common.Name]TypeName
	tables        map[common.Name]TypeName
	errorMessages map[uint64]string
	builtInTypes  map[TypeName]common.Pair
}

func (a AbiSerializer) ConfigureBuiltInTypes() {

}

func (a AbiSerializer) SetAbi(abi *AbiDef, maxSerializationTime *common.Microseconds){
	deadline := common.Now() + common.TimePoint(*maxSerializationTime)
	a.typeDefs = make(map[TypeName]TypeName)
	a.structs = make(map[TypeName]StructDef)
	a.actions = make(map[common.Name]TypeName)
	a.tables = make(map[common.Name]TypeName)
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

	a.validate()
}

func (a AbiSerializer) IsBuiltinType(rtype *TypeName) bool {
	for p := range a.builtInTypes {
		if p == *rtype {
			return true
		}
	}
	return false
}

func (a AbiSerializer) IsInteger(rtype *TypeName) bool {
	stype := string(*rtype)
	return strings.HasPrefix(stype,"int") || strings.HasPrefix(stype,"uint")
}

func (a AbiSerializer) IsType(rtype *TypeName, recursionDepth common.SizeT, deadline *common.TimePoint, maxSerializationTime common.Microseconds) bool{
	return true
}

func (a AbiSerializer) validate() bool {
	return true
}