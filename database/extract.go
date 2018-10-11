
package database
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	tagPrefix	 	= 	"multiIndex"
	tagID        	= 	"id"
	tagIdx       	= 	"orderedNonUnique"
	tagUniqueIdx 	= 	"orderedUnique"
	tagIncrement 	= 	"increment"
	tagLess		 	= 	"less"
	tagGreater	 	= 	"greater"
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

func (f *fieldInfo,)addFieldInfo(fieldName string,fieldValue *reflect.Value){
	f.fieldName = append(f.fieldName,fieldName)
	f.fieldValue = append(f.fieldValue,fieldValue)
}

func extractSort(tag *string,f *fieldInfo)error{

	tmp := *tag
	if strings.Contains(tmp,":"){

		parts := strings.Split(tmp,":")
		if len(parts) != 2 && parts[1] != tagLess && parts[1] != tagGreater {
			return ErrUnknownTag
		}

		tagSort := parts[1]
		*tag = parts[0]
		if tagSort == tagGreater{
			f.greater = true
		}
	}
	return nil
}

func extractF(value *reflect.Value, field *reflect.StructField, m *structInfo)error{

	tag := field.Tag.Get(tagPrefix)
	if tag == ""{
		return nil
	}

	fieldName := field.Name
	fieldValue := value

	tags := strings.Split(tag,",")
	if len(tags) == 1{
		// 单一 index or unique
		tag = tags[0]

		f  := fieldInfo{}
		f.typeName = m.Name
		err := extractSort(&tag,&f)
		if err != nil{
			return err
		}


		if tag != tagIdx && tag!= tagUniqueIdx{
			return ErrUnknownTag
		}
		if tag == tagUniqueIdx{
			f.unique = true
		}
		if v,ok := m.Fields[fieldName];ok{
			v.addFieldInfo(fieldName,fieldValue)
		}else{

			f.addFieldInfo(fieldName,fieldValue)
			m.Fields[fieldName] = &f
		}
		return nil
	}

	for _,tag := range tags{
		if tag == tagIdx || tag == tagUniqueIdx{
			continue
		}

		if tag == tagID{
			f  := fieldInfo{}
			f.typeName = m.Name
			f.unique = true
			f.addFieldInfo(fieldName,fieldValue)
			m.Fields[tag] = &f
			m.Id = value
			m.IncrementStart = 1
			continue
		}

		if strings.HasPrefix(tag, tagIncrement) {
			parts := strings.Split(tag, "=")
			if parts[0] != tagIncrement {
				return ErrUnknownTag
			}
			if _,ok := m.Fields[tagID];!ok{
				return ErrUnknownTag
			}
			if len(parts) > 1 {
				incrementStart, err := strconv.ParseInt(parts[1], 0, 64)
				if err != nil {
					return err
				}
				m.IncrementStart = incrementStart
			}
		}


		f  := fieldInfo{}
		f.typeName = m.Name
		err := extractSort(&tag,&f)
		if err != nil{
			return err
		}

		if v,ok := m.Fields[tag];ok{
			v.addFieldInfo(fieldName,fieldValue)
		}else{
			f.addFieldInfo(fieldName,fieldValue)
			m.Fields[tag] = &f
		}
	}

	return  nil
}
