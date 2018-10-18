
package database

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	fieldName	fieldName	fieldName
	fieldValue	fieldValue	fieldValue

tag
	fieldName	fieldName	fieldName
	fieldValue	fieldValue	fieldValue

tag
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


// from
func numberfromb(raw []byte) (int64, error) {
	r := bytes.NewReader(raw)
	var to int64
	err := binary.Read(r, binary.BigEndian, &to)
	if err != nil {
		return 0, err
	}
	return to, nil
}
// to
func numbertob(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
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

func splitSubTag(fieldName string,fieldValue *reflect.Value,tag string,m *structInfo)error{

	tags := strings.Split(tag,",")
	//fmt.Println(tags)
	tagLen := len(tags)
	tagPre := tags[0]
	if tagPre == tagID{
		for _,subTag := range tags{
			//fmt.Println(subTag)
			if subTag == tagGreater || subTag ==  tagLess {
				return ErrIdNoSort
			}

			f  := fieldInfo{}
			f.unique = true
			m.Id = fieldValue
			m.IncrementStart = 1
			addFieldInfo(subTag,tagID,fieldValue,&f,m)

		}
		return nil
	}else if tagPre == tagUniqueIdx || tagPre == tagNoUniqueIdx{
		f  := fieldInfo{}
		if tagPre == tagUniqueIdx{
			f.unique = true
		}
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

