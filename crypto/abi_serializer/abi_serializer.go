package abi_serializer

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"strconv"
	"strings"
)

var maxRecursionDepth = 32

type typeName = string

func Encode_Decode() common.Pair {
	decode := func() {
	}
	encode := func() {

	}
	return common.MakePair(decode, encode)
}

type AbiSerializer struct {
	abi           *AbiDef
	typeDefs      map[string]string
	structs       map[string]StructDef
	actions       map[common.Name]string
	tables        map[common.Name]string
	errorMessages map[uint64]string
	variants      map[string]VariantDef
	builtInTypes  map[string]common.Pair
}

func (a AbiSerializer) ConfigureBuiltInTypes() {
	//
}

func NewAbiSerializer(abi *AbiDef, maxSerializationTime common.Microseconds) *AbiSerializer {
	abiSer := AbiSerializer{}
	abiSer.ConfigureBuiltInTypes()
	abiSer.SetAbi(abi, maxSerializationTime)
	abiSer.abi = abi
	return &abiSer
}

func (a *AbiSerializer) SetAbi(abi *AbiDef, maxSerializationTime common.Microseconds) {
	//deadline := common.Now() + common.TimePoint(maxSerializationTime)
	a.typeDefs = make(map[string]string)
	a.structs = make(map[string]StructDef)
	a.actions = make(map[common.Name]string)
	a.tables = make(map[common.Name]string)
	a.errorMessages = make(map[uint64]string)

	for _, st := range abi.Structs {
		a.structs[st.Name] = st
	}

	for _, td := range abi.Types {
		//try.EosAssert(a.IsType(&td.Type, 0, &deadline, maxSerializationTime), &exception.InvalidTypeInsideAbi{}, "invalid type : %v", td.Type)
		//try.EosAssert(!a.IsType(&td.NewTypeName, 0, &deadline, maxSerializationTime), &exception.DuplicateAbiTypeDefException{}, "type already exists : %v", td.Type)
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

//bool abi_serializer::is_type(const type_name& type, const fc::microseconds& max_serialization_time)const {
//impl::abi_traverse_context ctx(max_serialization_time);
//return _is_type(type, ctx);
//}
//func (a AbiSerializer)IsType(rtype *string,maxSerializationTime common.Microseconds) bool{
//	return false //TODO
//}

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
	return false
}

func (a AbiSerializer) GetStruct(rtype *typeName) *StructDef {

	itr, ok := a.structs[a.ResolveType(rtype)]
	try.EosAssert(ok, &exception.InvalidTypeInsideAbi{}, "Unknown struct %s", rtype)
	return &itr
}

func (a AbiSerializer) ResolveType(rtype *typeName) typeName {
	var ok bool
	var itr, t string
	itr, ok = a.typeDefs[*rtype]
	if ok {
		for i := len(a.typeDefs); i > 0; i-- { // avoid infinite recursion
			t = itr
			itr, ok = a.typeDefs[t]
			if !ok {
				return t
			}
		}
	}
	return ""
}
func (a AbiSerializer) validate() bool {
	return true //TODO need to check Abi
}

//void abi_serializer::validate( impl::abi_traverse_context& ctx )const {
//   for( const auto& t : typedefs ) { try {
//      vector<type_name> types_seen{t.first, t.second};
//      auto itr = typedefs.find(t.second);
//      while( itr != typedefs.end() ) {
//         ctx.check_deadline();
//         EOS_ASSERT( find(types_seen.begin(), types_seen.end(), itr->second) == types_seen.end(), abi_circular_def_exception, "Circular reference in type ${type}", ("type",t.first) );
//         types_seen.emplace_back(itr->second);
//         itr = typedefs.find(itr->second);
//      }
//   } FC_CAPTURE_AND_RETHROW( (t) ) }
//   for( const auto& t : typedefs ) { try {
//      EOS_ASSERT(_is_type(t.second, ctx), invalid_type_inside_abi, "${type}", ("type",t.second) );
//   } FC_CAPTURE_AND_RETHROW( (t) ) }
//   for( const auto& s : structs ) { try {
//      if( s.second.base != type_name() ) {
//         struct_def current = s.second;
//         vector<type_name> types_seen{current.name};
//         while( current.base != type_name() ) {
//            ctx.check_deadline();
//            const auto& base = get_struct(current.base); //<-- force struct to inherit from another struct
//            EOS_ASSERT( find(types_seen.begin(), types_seen.end(), base.name) == types_seen.end(), abi_circular_def_exception, "Circular reference in struct ${type}", ("type",s.second.name) );
//            types_seen.emplace_back(base.name);
//            current = base;
//         }
//      }
//      for( const auto& field : s.second.fields ) { try {
//         ctx.check_deadline();
//         EOS_ASSERT(_is_type(_remove_bin_extension(field.type), ctx), invalid_type_inside_abi, "${type}", ("type",field.type) );
//      } FC_CAPTURE_AND_RETHROW( (field) ) }
//   } FC_CAPTURE_AND_RETHROW( (s) ) }
//   for( const auto& s : variants ) { try {
//      for( const auto& type : s.second.types ) { try {
//         ctx.check_deadline();
//         EOS_ASSERT(_is_type(type, ctx), invalid_type_inside_abi, "${type}", ("type",type) );
//      } FC_CAPTURE_AND_RETHROW( (type) ) }
//   } FC_CAPTURE_AND_RETHROW( (s) ) }
//   for( const auto& a : actions ) { try {
//     ctx.check_deadline();
//     EOS_ASSERT(_is_type(a.second, ctx), invalid_type_inside_abi, "${type}", ("type",a.second) );
//   } FC_CAPTURE_AND_RETHROW( (a)  ) }
//
//   for( const auto& t : tables ) { try {
//     ctx.check_deadline();
//     EOS_ASSERT(_is_type(t.second, ctx), invalid_type_inside_abi, "${type}", ("type",t.second) );
//   } FC_CAPTURE_AND_RETHROW( (t)  ) }
//}
//

func (a AbiSerializer) GetActionType(action common.Name) typeName {
	itr, ok := a.actions[action]
	if !ok {
		return ""
	}
	return itr
}

func (a AbiSerializer) GetTableType(action common.Name) typeName {
	itr, ok := a.tables[action]
	if !ok {
		return ""
	}
	return itr
}

func (a AbiSerializer) GetErrorMessage(errorCode uint64) string {
	itr, ok := a.errorMessages[errorCode]
	if !ok {
		return ""
	}
	return itr
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

//   bool abi_serializer::_is_type(const type_name& rtype, impl::abi_traverse_context& ctx )const {
//      auto h = ctx.enter_scope();
//      auto type = fundamental_type(rtype);
//      if( built_in_types.find(type) != built_in_types.end() ) return true;
//      if( typedefs.find(type) != typedefs.end() ) return _is_type(typedefs.find(type)->second, ctx);
//      if( structs.find(type) != structs.end() ) return true;
//      if( variants.find(type) != variants.end() ) return true;
//      return false;
//   }
//
//func (a AbiSerializer)isType(rtype typeName,ctx)

func (a *AbiSerializer) VariantToBinary(name typeName, body *common.Variants, maxSerializationTime common.Microseconds) []byte {
	return a._VariantToBinary(name, body, true, 0, common.Now().AddUs(maxSerializationTime), maxSerializationTime)
}

//bytes abi_serializer::_variant_to_binary( const type_name& type, const fc::variant& var, bool allow_extensions,
//                                          size_t recursion_depth, const fc::time_point& deadline, const fc::microseconds& max_serialization_time )const
//{ try {
//   EOS_ASSERT( ++recursion_depth < max_recursion_depth, abi_recursion_depth_exception, "recursive definition, max_recursion_depth ${r} ", ("r", max_recursion_depth) );
//   EOS_ASSERT( fc::time_point::now() < deadline, abi_serialization_deadline_exception, "serialization time limit ${t}us exceeded", ("t", max_serialization_time) );
//   if( !_is_type(type, recursion_depth, deadline, max_serialization_time) ) {
//      return var.as<bytes>();
//   }
//
//   bytes temp( 1024*1024 );
//   fc::datastream<char*> ds(temp.data(), temp.size() );
//   _variant_to_binary(type, var, ds, allow_extensions, recursion_depth, deadline, max_serialization_time);
//   temp.resize(ds.tellp());
//   return temp;
//} FC_CAPTURE_AND_RETHROW( (type)(var) ) }
func (a *AbiSerializer) _VariantToBinary(name typeName, data *common.Variants, allowExtensions bool, recursion_depth common.SizeT,
	deadline common.TimePoint, maxSerializationTime common.Microseconds) (re []byte) {
	try.Try(func() {
		buf, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("tester GetAction Marshal is error:%s", err)
			try.Throw(fmt.Errorf("tester GetAction Marshal is error:%s", err))
		}

		re, err = a.abi.EncodeAction(common.N(name), buf)
		if err != nil {
			fmt.Printf("encode actoin is error:%s", err)
			try.Throw(fmt.Errorf("encode actoin is error:%s", err))
		}
	}).FcCaptureAndRethrow(name, data).End()
	return re

}

//fc::variant abi_serializer::binary_to_variant( const type_name& type, const bytes& binary, const fc::microseconds& max_serialization_time, bool short_path )const {
//impl::binary_to_variant_context ctx(*this, max_serialization_time, type);
//ctx.short_path = short_path;
//return _binary_to_variant(type, binary, ctx);
//}
//todo shortPath default is false
func (a AbiSerializer) BinaryToVariant(rtype typeName, binary []byte, maxSerializationTime common.Microseconds, shortPath bool) common.Variants {
	var re common.Variants
	try.Try(func() {
		bytes, err := a.abi.DecodeAction(rtype, binary)
		if err != nil {
			try.Throw(fmt.Sprintf("binary_to_variant is error: %s", err.Error()))
		}
		err = json.Unmarshal(bytes, &re)
		if err != nil {
			try.Throw(fmt.Sprintf("unmarshal variants is error: %s", err.Error()))
		}
	}).EosRethrowExceptions(&exception.UnpackException{}, "Unable to unpack %s from bytes", string(binary)).End()
	return re
}
