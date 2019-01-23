package abi_serializer

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/eosspark/eos-go/exception"
	. "github.com/eosspark/eos-go/exception/try"
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
	abi           *AbiDef
	typeDefs      map[typeName]typeName
	structs       map[typeName]StructDef
	actions       map[common.Name]typeName
	tables        map[common.Name]typeName
	errorMessages map[uint64]string
	variants      map[typeName]VariantDef
	builtInTypes  map[string]int
}

func (a *AbiSerializer) ConfigureBuiltInTypes() {
	a.builtInTypes = make(map[string]int, 31)
	a.builtInTypes["bool"] = 1
	a.builtInTypes["int8"] = 1
	a.builtInTypes["uint8"] = 1
	a.builtInTypes["int16"] = 1
	a.builtInTypes["uint16"] = 1
	a.builtInTypes["int32"] = 1
	a.builtInTypes["uint32"] = 1
	a.builtInTypes["int64"] = 1
	a.builtInTypes["uint64"] = 1
	a.builtInTypes["int128"] = 1
	a.builtInTypes["uint128"] = 1
	a.builtInTypes["varint32"] = 1
	a.builtInTypes["varuint32"] = 1

	// TODO: Add proper support for floating point types. For now this is good enough.
	a.builtInTypes["float32"] = 1
	a.builtInTypes["float64"] = 1
	a.builtInTypes["float128"] = 1

	a.builtInTypes["time_point"] = 1
	a.builtInTypes["time_point_sec"] = 1
	a.builtInTypes["block_timestamp_type"] = 1

	a.builtInTypes["name"] = 1

	a.builtInTypes["bytes"] = 1
	a.builtInTypes["string"] = 1

	a.builtInTypes["checksum160"] = 1
	a.builtInTypes["checksum256"] = 1
	a.builtInTypes["checksum512"] = 1

	a.builtInTypes["public_key"] = 1
	a.builtInTypes["signature"] = 1

	a.builtInTypes["symbol"] = 1
	a.builtInTypes["symbol_code"] = 1
	a.builtInTypes["asset"] = 1
	a.builtInTypes["extended_asset"] = 1

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

	for _, td := range abi.Types { //TODO valide types
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
	for _, v := range abi.Variants { //may not exit
		a.variants[v.Name] = v
	}

	EosAssert(len(a.typeDefs) == len(abi.Types), &exception.DuplicateAbiTypeDefException{}, "duplicate type definition detected")
	EosAssert(len(a.structs) == len(abi.Structs), &exception.DuplicateAbiStructDefException{}, "duplicate struct definition detected")
	EosAssert(len(a.actions) == len(abi.Actions), &exception.DuplicateAbiActionDefException{}, "duplicate action definition detected")
	EosAssert(len(a.tables) == len(abi.Tables), &exception.DuplicateAbiTableDefException{}, "duplicate table definition detected")
	EosAssert(len(a.errorMessages) == len(abi.ErrorMessages), &exception.DuplicateAbiErrMsgDefException{}, "duplicate error message definition detected")
	EosAssert(len(a.variants) == len(abi.Variants), &exception.DuplicateAbiVariantDefException{}, "duplicate variant definition detected")
	a.validate() //TODO always return true
	a.abi = abi
}

func (a AbiSerializer) IsBuiltinType(stype typeName) bool { //TODO
	for p := range a.builtInTypes {
		if p == stype {
			return true
		}
	}
	return false
}

func (a AbiSerializer) IsInteger(stype typeName) bool {
	return strings.HasPrefix(stype, "int") || strings.HasPrefix(stype, "uint")
}

func (a AbiSerializer) GetIntegerSize(stype typeName) int {
	EosAssert(a.IsInteger(stype), &exception.InvalidTypeInsideAbi{}, "%v is not an integer type", stype)
	var num int
	if strings.HasPrefix(stype, "uint") {
		num, _ = strconv.Atoi(string([]byte(stype)[4:]))
		return num
	} else {
		num, _ = strconv.Atoi(string([]byte(stype)[3:]))
		return num
	}
}

func (a AbiSerializer) IsStruct(stype typeName) bool {
	name := a.ResolveType(stype)
	for _, p := range a.structs {
		if p.Name == name {
			return true
		}
	}
	return false
}

func (a AbiSerializer) IsArray(stype typeName) bool {
	return strings.HasSuffix(stype, "[]")
}

func (a AbiSerializer) IsOptional(stype typeName) bool {
	return strings.HasSuffix(stype, "?")
}

//bool abi_serializer::is_type(const type_name& type, const fc::microseconds& max_serialization_time)const {
//impl::abi_traverse_context ctx(max_serialization_time);
//return _is_type(type, ctx);
//}
//func (a AbiSerializer)IsType(rtype *string,maxSerializationTime common.Microseconds) bool{
//	return false //TODO
//}

func (a AbiSerializer) FundamentalType(stype typeName) string {
	btype := []byte(stype)
	if a.IsArray(stype) {
		return string(string(btype[0 : len(btype)-2]))
	} else if a.IsOptional(stype) {
		return string(string(btype[0 : len(btype)-1]))
	} else {
		return stype
	}
}

func (a AbiSerializer) IsType(rtype typeName, recursionDepth common.SizeT, deadline *common.TimePoint, maxSerializationTime common.Microseconds) bool {
	EosAssert(common.Now() < *deadline, &exception.AbiSerializationDeadlineException{}, "serialization time limit %vus exceeded", maxSerializationTime)
	recursionDepth++
	if recursionDepth > maxRecursionDepth {
		return false
	}
	ftype := a.FundamentalType(rtype)
	if a.IsBuiltinType(ftype) {
		return true
	}

	if a.IsStruct(ftype) {
		return true
	}
	return false
}

func (a AbiSerializer) GetStruct(stype typeName) *StructDef {
	itr, ok := a.structs[a.ResolveType(stype)]
	EosAssert(ok, &exception.InvalidTypeInsideAbi{}, "Unknown struct %s", stype)
	return &itr
}

func (a AbiSerializer) ResolveType(stype typeName) typeName {
	var ok bool
	var itr, t string
	itr, ok = a.typeDefs[stype]
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
	if itr, ok := a.actions[action]; ok {
		return itr
	}
	return ""
}

func (a AbiSerializer) GetTableType(action common.Name) typeName {
	if itr, ok := a.tables[action]; ok {
		return itr
	}
	return ""
}

func (a AbiSerializer) GetErrorMessage(errorCode uint64) string {
	if itr, ok := a.errorMessages[errorCode]; ok {
		return itr
	}
	return ""
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
		fmt.Println(err)
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
	return a.variantToBinary(name, body, true, 0, common.Now().AddUs(maxSerializationTime), maxSerializationTime)
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
func (a *AbiSerializer) variantToBinary(name typeName, data *common.Variants, allowExtensions bool, recursion_depth common.SizeT,
	deadline common.TimePoint, maxSerializationTime common.Microseconds) (re []byte) {
	Try(func() {
		buf, err := json.Marshal(data)
		if err != nil {
			abiLog.Error("Marshal action is error: %s", err.Error())
			Throw(fmt.Sprintf("Marshal action is error: %s", err.Error()))
		}

		re, err = a.abi.EncodeStruct(name, buf)
		if err != nil {
			abiLog.Error("encode action is error:%s", err)
			Throw(fmt.Errorf("encode actoin is error:%s", err))
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
	Try(func() {
		var bytes []byte
		var err error
		if a.abi.StructForName(rtype) != nil {
			bytes, err = a.abi.DecodeStruct(rtype, binary)
		}
		//else if a.abi.TableForName(rtype) != nil {
		//	bytes, err = a.abi.DecodeTableRow(rtype, binary)
		//}

		if err != nil {
			fmt.Println(err.Error())
			Throw(fmt.Sprintf("binary_to_variant is error: %s", err.Error()))
		}
		err = json.Unmarshal(bytes, &re)
		if err != nil {
			fmt.Println(err.Error())
			Throw(fmt.Sprintf("unmarshal variants is error: %s", err.Error()))
		}
	}).EosRethrowExceptions(&exception.UnpackException{}, "Unable to unpack %s from bytes", string(binary)).End()
	return re
}

func (a AbiSerializer) BinaryToVariant2(rtype typeName, binary []byte, maxSerializationTime common.Microseconds, shortPath bool) []byte {
	var re []byte
	Try(func() {
		var bytes []byte
		var err error
		if a.abi.StructForName(rtype) != nil {
			bytes, err = a.abi.DecodeStruct(rtype, binary)
		}
		//else if a.abi.TableForName(rtype) != nil {
		//	bytes, err = a.abi.DecodeTableRow(rtype, binary)
		//}

		if err != nil {
			fmt.Println(err.Error())
			Throw(fmt.Sprintf("binary_to_variant is error: %s", err.Error()))
		}
		re = bytes

	}).EosRethrowExceptions(&exception.UnpackException{}, "Unable to unpack %s from bytes", string(binary)).End()
	return re
}

func (a AbiSerializer) BinaryToVariantPrint(rtype typeName, binary []byte, maxSerializationTime common.Microseconds, shortPath bool) interface{} {
	var bytes []byte
	Try(func() {
		var err error
		if a.abi.StructForName(rtype) != nil {
			bytes, err = a.abi.DecodeStruct(rtype, binary)
		}
		if err != nil {
			Throw(fmt.Sprintf("binary_to_variant is error: %s", err.Error()))
		}
	}).EosRethrowExceptions(&exception.UnpackException{}, "Unable to unpack %s from bytes", string(binary)).End()
	return bytes
}

//template<typename T, typename Resolver>
//void abi_serializer::to_variant( const T& o, variant& vo, Resolver resolver, const fc::microseconds& max_serialization_time ) try {
//mutable_variant_object mvo;
//impl::abi_traverse_context ctx(max_serialization_time);
//impl::abi_to_variant::add(mvo, "_", o, resolver, ctx);
//vo = std::move(mvo["_"]);
//} FC_RETHROW_EXCEPTIONS(error, "Failed to serialize type", ("object",o))
//
//template<typename T, typename Resolver>
//void abi_serializer::from_variant( const variant& v, T& o, Resolver resolver, const fc::microseconds& max_serialization_time ) try {
//impl::abi_traverse_context ctx(max_serialization_time);
//impl::abi_from_variant::extract(v, o, resolver, ctx);
//} FC_RETHROW_EXCEPTIONS(error, "Failed to deserialize variant", ("variant",v))
//
//
//} } // eosio::chain

func ToVariant() {

}

//func FromVariant(v *common.Variants, o *types.SignedTransaction, resolver *AbiSerializer, maxSerialization common.Microseconds) {
//	data, err := json.Marshal(v)
//	fmt.Println(data, err)
//	fmt.Printf("%s\n", string(data))
//	err = json.Unmarshal(data, o)
//	fmt.Println(o, err)
//
//	type Actions struct {
//		Actions []types.Action `json:'actions'`
//	}
//	var actions Actions
//	err = json.Unmarshal(data, &actions)
//	fmt.Println(err, actions)
//}

//template<typename M, typename Resolver, not_require_abi_t<M> = 1>
//static void extract( const variant& v, M& o, Resolver, abi_traverse_context& ctx )
//{
//   auto h = ctx.enter_scope();
//   from_variant(v, o);
//}
func FromVariant(v *common.Variants, o *types.SignedTransaction, resolver func(account common.AccountName) *AbiSerializer, maxSerialization common.Microseconds) {
	data, err := json.Marshal(v)
	fmt.Println(data, err)
	fmt.Printf("%s\n", string(data))
	err = json.Unmarshal(data, o)
	fmt.Println(o, err)

	type Actions struct {
		Actions []types.Action `json:'actions'`
	}
	var actions Actions
	err = json.Unmarshal(data, &actions)
	fmt.Println(err, actions)
}

func FromVariantToActionData(actionParams *types.ContractTypesInterface, v *common.Variants, resolver func(account common.AccountName) *AbiSerializer, maxSerialization common.Microseconds) {
	data, err := json.Marshal(v)
	fmt.Println(data, err)
	fmt.Printf("%s\n", string(data))
	err = json.Unmarshal(data, actionParams)
	fmt.Println(actionParams, err)
}

func ToVariantFromActionData(actionParams *types.ContractTypesInterface, v *common.Variants, resolver func(account common.AccountName) *AbiSerializer, maxSerialization common.Microseconds) {
	bytes, err := json.Marshal(actionParams)
	if err != nil {
		fmt.Printf("marshal actionParams is error :%s\n", err.Error())
	}
	err = json.Unmarshal(bytes, v)
	if err != nil {
		fmt.Printf("Unmarshal variants is error: %s\n", err.Error())
	}

}
