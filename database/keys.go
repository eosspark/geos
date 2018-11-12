package database

import (
	"log"
	"reflect"
)

func idKey(id, typeName []byte) []byte { /* key --> typeName__id */
	key := cloneByte(typeName)
	key = append(key, '_')
	key = append(key, '_')
	key = append(key, id...)

	return key
}

func getFieldInfo(fieldName string, value interface{}) (*fieldInfo, error) {
	ref := reflect.ValueOf(value)
	if !ref.IsValid() || reflect.Indirect(ref).Kind() != reflect.Struct {
		return nil, ErrBadType
	}
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}
	cfg, err := extractObjectTagInfo(&ref)
	if err != nil {
		return nil, err
	}

	fields, ok := cfg.Fields[fieldName]
	if !ok {
		return nil, ErrNotFound
	}
	return fields, nil
}

// TODO The function is unchanged, need to modify the implementation
func getFieldValue(info *fieldInfo) []byte { /* non unique fields --> get function */
	values := []byte{}
	if len(info.fieldValue) == 0{
		log.Println("object is empty")
		return nil
	}
	for _, v := range info.fieldValue {

		if v.Kind() != reflect.Bool && isZero(v) {
			return values
		}
		values = append(values, '_')
		values = append(values, '_')
		re, err :=    EncodeToBytes(v.Interface())
		if err != nil {
			return nil
		}

		values = append(values, re...)
	}

	return values
}

func typeNameFieldName(typeName, tagName []byte) []byte { /* typeName__fieldName*/
	key := cloneByte(typeName)
	key = append(key, '_')
	key = append(key, '_')
	key = append(key, tagName...)
	return key
}

func getNonUniqueEnd(key []byte) []byte { /* non unique fields --> regexp*/
	end := make([]byte, len(key))
	copy(end, key)
	end[len(end)-1] = end[len(end)-1] + 1
	return end
}

func fieldValueToByte(key []byte, info *fieldInfo) []byte { /* fieldValue[0]__fieldValue[1]... */
	cloneKey := cloneByte(key)
	for _, v := range info.fieldValue { // typeName__tag__fieldValue...
		cloneKey = append(cloneKey, '_')
		cloneKey = append(cloneKey, '_')
		//fmt.Println(v.Interface())
		value, err := EncodeToBytes(v.Interface())
		if err != nil {
			return nil
		}
		//fmt.Println(value)
		cloneKey = append(cloneKey, value...)
		//fmt.Println(cloneKey)
	}
	return cloneKey
}
