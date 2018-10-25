
package database

import (
	"fmt"
	"github.com/eosspark/eos-go/crypto/rlp"
	"reflect"
	"strings"
)

const (
	tagPrefix	 	= 	"multiIndex"
	tagID        	= 	"id"
	tagNoUniqueIdx  = 	"orderedNonUnique"
	tagUniqueIdx 	= 	"orderedUnique"
	tagIncrement 	= 	"increment"
	tagLess		 	= 	"less"
	tagGreater	 	= 	"greater"
	tagInline	 	= 	"inline"
)

/*

tag
	unique
	greater
	typeName
	fieldName	fieldName	fieldName
	fieldValue	fieldValue	fieldValue
*/


type fieldInfo struct{
	unique 		bool
	greater		bool
	typeName 	string
	fieldName 	[]string
	fieldValue 	[]*reflect.Value
}

/*

tag
	name 			--> TypeName
	IncrementStart 	--> for id
	Id
	fields			--> tag-fieldInfo tag-fieldInfo tag-fieldInfo
*/
//// TODO A separate module for external use in the future
type structInfo struct{
	Name 			string
	IncrementStart 	int64
	Id				*reflect.Value
	Fields 			map[string]*fieldInfo
}

func isZero(v *reflect.Value) bool {
	zero := reflect.Zero(v.Type()).Interface()
	current := v.Interface()
	return reflect.DeepEqual(current, zero)
}


func (s*structInfo)showStructInfo(){
	fmt.Println("name : ",s.Name)
	fmt.Println("IncrementStart : ",s.IncrementStart)

	for k,v := range s.Fields{
		fmt.Println("fields key is 				: ",k)
		fmt.Println("unique			 			: ",v.unique)
		fmt.Println("greater			 			: ",v.greater)
		fmt.Println("fields fieldName 	is 		: ",v.fieldName)

		for _,va := range v.fieldValue{
			fmt.Println("fields fieldValue 	is 		: ",va.Interface())
		}
	}
}

func cloneInterface(data interface{}) interface{} {

	src := reflect.ValueOf(data)
	dst := reflect.New(reflect.Indirect(src).Type())
	if src.Kind() == reflect.Ptr{
		src = src.Elem()
	}
	dstElem := dst.Elem()
	NumField := src.NumField()
	for i := 0; i < NumField; i++ {
		sf := src.Field(i)
		df := dstElem.Field(i)
		df.Set(sf)
	}
	return dst.Interface()
}

func cloneByte(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

func parseObjectToCfg(in interface{})( *structInfo,error){
	ref := reflect.ValueOf(in)
	if !ref.IsValid() || ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
		return nil,ErrStructPtrNeeded
	}

	cfg, err := extractStruct(&ref)
	if err != nil {
		return nil,err
	}

	if _, ok := cfg.Fields[tagID]; !ok {
		return nil,ErrNoID
	}
	return cfg,nil
}

func extractStruct(s *reflect.Value,mi ...*structInfo) (*structInfo,error) {

	if s.Kind() == reflect.Ptr {
		e := s.Elem()
		s = &e
	}
	if s.Kind() != reflect.Struct {
		return nil, ErrBadType
	}

	typ := s.Type()
	var m *structInfo
	if len(mi) > 0 {
		m = mi[0]
	} else {
		m = &structInfo{}
		m.Fields = make(map[string]*fieldInfo)
	}

	if m.Name == "" {
		m.Name = typ.Name()
	}

	numFields := s.NumField()
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		value := s.Field(i)

		if field.PkgPath != "" {
			continue
		}

		err := extractF(&value, &field, m)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}

func extractF(value *reflect.Value, field *reflect.StructField, m *structInfo)error{

	tag := field.Tag.Get(tagPrefix)
	if tag == ""{
		return nil
	}

	tags := strings.Split(tag,":")

	//fmt.Println(tags)

	for _,tag := range tags{
		err := splitSubTag(field.Name,value,tag,m)
		if err != nil{
			return err
		}
	}
	return  nil
}

func splitSubTag(fieldName string,fieldValue *reflect.Value,tag string,m *structInfo)error{

	tags := strings.Split(tag,",")
	//fmt.Println(tags)
	tagPre := tags[0]
	if tagPre == tagID{

		return doIdTag(tags,fieldValue,m)
	}else if tagPre == tagUniqueIdx || tagPre == tagNoUniqueIdx{

		return doUniqueOrNoUniqueTag(tagPre,fieldName,tags,fieldValue,m)
	}else if tagPre == tagInline{
		_,err := extractStruct(fieldValue,m)
		if err != nil{
			return err
		}
		return nil
	}
	return doOtherTag(tagPre,fieldName,tags,fieldValue,m)
}

func doIdTag(tags []string,fieldValue *reflect.Value,m *structInfo)error{
	for _,subTag := range tags{
		//fmt.Println(subTag)
		if subTag == tagGreater || subTag ==  tagLess {
			return ErrIdNoSort
		}
		if subTag == tagIncrement{
			continue
		}
		f  := fieldInfo{}
		f.unique = true
		m.Id = fieldValue
		m.IncrementStart = 1
		addFieldInfo(subTag,tagID,fieldValue,&f,m)

	}
	return nil
}

func doUniqueOrNoUniqueTag(tagPre,fieldName string,tags []string,fieldValue *reflect.Value,m *structInfo)error{

	f  := fieldInfo{}
	if tagPre == tagUniqueIdx{
		f.unique = true
	}
	tagLen := len(tags)
	subTag := fieldName
	if tagLen > 1{
		sor := tags[1]
		if sor != tagGreater && sor != tagLess{
			return ErrTagInvalid
		}
		if sor == tagGreater{
			f.greater = true
		}
	}
	addFieldInfo(subTag,fieldName,fieldValue,&f,m)
	return nil
}

func doOtherTag(tagPre,fieldName string,tags []string,fieldValue *reflect.Value,m *structInfo)error{

	tagLen := len(tags)
	if tagLen < 2{
		return ErrTagInvalid
	}

	f  := fieldInfo{}
	f.typeName = m.Name

	if tagLen > 2{
		sor := tags[2]
		if sor != tagGreater && sor != tagLess{
			return ErrTagInvalid
		}
		if sor == tagGreater{
			f.greater = true
		}
	}
	tagIdx := tags[1]
	if tagIdx != tagUniqueIdx && tagIdx != tagNoUniqueIdx{
		return ErrTagInvalid
	}

	if tagIdx == tagUniqueIdx{
		f.unique = true
	}
	addFieldInfo(tagPre,fieldName,fieldValue,&f,m)
	return nil
}

func addFieldInfo(tag ,fieldName string,fieldValue *reflect.Value,f *fieldInfo,m * structInfo){
	if v,ok := m.Fields[tag];ok{
		v.fieldName = append(v.fieldName,fieldName)
		v.fieldValue = append(v.fieldValue,fieldValue)
	}else{
		f.typeName = m.Name
		f.fieldName = append(f.fieldName,fieldName)
		f.fieldValue = append(f.fieldValue,fieldValue)
		m.Fields[tag] = f
	}

}



/*
								all key
increment			-->	typeName
id field  			--> typeName__tagName__fieldValue
unique fields 		--> typeName__tagName__fieldValue
non unique field 	--> typeName__fieldName__idFieldValue__fieldValue
non unique fields 	--> typeName__tagName__fieldValue[0]__fieldValue[1]...

								all value

increment			-->	val
id field  			--> objectValue
unique fields 		--> idFieldValue
non unique field 	--> idFieldValue
non unique fields 	--> idFieldValue

*/

type kv struct {
	key		[]byte
	value	[]byte
}

type incrementKV struct {
	key			[]byte
	oldValue	[]byte
	newValue	[]byte
	delete 		bool

}

type dbKeyValue struct{
	id 			kv
	index 		[]kv
	typeName 	[]byte
	increment 	incrementKV
	first 		bool
}

func (kv *kv)showKV(){
	space := " : "
	fmt.Println(kv.key,space,kv.value)
}

func (increment *incrementKV)showIncrement(){
	space := " : "
	fmt.Println(increment.key,space,increment.oldValue,space,increment.newValue)
}

func (dbKV *dbKeyValue)showDbKV(){
	fmt.Println("--------------------- show db kv begin ---------------------")
	dbKV.id.showKV()
	for _,v := range dbKV.index{
		v.showKV()
	}
	fmt.Println(dbKV.typeName)
	dbKV.increment.showIncrement()
	fmt.Println("--------------------- show db kv end  ---------------------")
}

func structKV(in interface{} ,dbKV *dbKeyValue,cfg *structInfo)error{

	objValue, err := rlp.EncodeToBytes(in)
	if err != nil {
		return err
	}
	objId, err := rlp.EncodeToBytes(cfg.Id.Interface())
	if err != nil {
		return err
	}

	idk := idKey(objId,[]byte(cfg.Name))
	kv_ := kv{}
	kv_.key 	= idk
	kv_.value 	= objValue
	dbKV.id 	= kv_
	dbKV.typeName = []byte(cfg.Name)

	cfgToKV(objId,cfg,dbKV)

	return nil
}

func cfgToKV (objId []byte,cfg *structInfo,dbKV *dbKeyValue) {

	typeName := []byte(cfg.Name)

	for tag, fieldCfg := range cfg.Fields {
		prefix 	:= append(typeName, '_') 					/* 			typeName__ 				*/
		prefix 	=  append(prefix, '_')
		prefix 	=  append(prefix, tag...) 						/* 			typeName__tagName__ 	*/
		key 	:= getFieldValue(prefix, fieldCfg)
		if !fieldCfg.unique && len(fieldCfg.fieldValue) == 1{ 	/* 			non unique 				*/
			 key = append(key, objId...)
		}

		kv_ := kv{}
		kv_.key 	= key
		kv_.value	= objId
		dbKV.index 	= append(dbKV.index,kv_)
	}
}


