package database

import (
	"reflect"
)

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


func splicingString(k,v[]byte) []byte {
	key := cloneByte(k)
	key = append(key, '_')
	key = append(key, '_')
	key = append(key, v...)
	return key
}

func keyEnd(key []byte) []byte { /* non unique fields --> regexp*/
	end := cloneByte(key)
	end[len(end)-1] = end[len(end)-1] + 1
	return end
}


func fieldValueToByte(info *fieldInfo,zero...bool) ([]byte,error) { /* fieldValue[0]__fieldValue[1]... */
	cloneKey := []byte{}

	skipZero := false
	if len(zero) > 0{
		skipZero = zero[0]
	}

	for _, v := range info.fieldValue { // typeName__tag__fieldValue...
		cloneKey = append(cloneKey, '_')
		cloneKey = append(cloneKey, '_')
		if skipZero{ //FIXME field value is zero ?
			if v.Kind() != reflect.Bool && isZero(v) {
				continue
			}
		}
		value, err := EncodeToBytes(v.Interface())
		if err != nil {
			return nil,err
		}
		cloneKey = append(cloneKey, value...)
	}
	return cloneKey,nil
}
