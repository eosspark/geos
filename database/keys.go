package database

import (
	"github.com/eosspark/eos-go/crypto/rlp"
	"reflect"
)

func idKey(id, typeName []byte) []byte { /* key --> typeName__id */
	key := append(typeName, '_')
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
		//return nil, ErrStructNeeded
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

func nonUniqueValue(info *fieldInfo) []byte { /* non unique fields --> find function */
	for _, v := range info.fieldValue {
		if isZero(v) && v.Kind() != reflect.Bool {
			return nil
		}
	}
	val, _ := getNonUniqueFieldValue(info)
	return val
}

// TODO The function is unchanged, need to modify the implementation
func getNonUniqueFieldValue(info *fieldInfo) ([]byte, []byte) { /* non unique fields --> get function */
	values := []byte{}
	prefix := []byte{}
	//regexp := []byte{40,46,42,41}
	for _, v := range info.fieldValue {
		values = append(values, '_')
		values = append(values, '_')
		if v.Kind() != reflect.Bool && isZero(v) {
			//values = append(values,regexp...)
			return prefix, prefix
		}
		re, err := rlp.EncodeToBytes(v.Interface())
		if err != nil {
			return nil, nil
		}

		prefix = append(prefix, re...)
		values = append(values, re...)
	}

	return values, prefix /*  NON unique values , unique use prefix*/
}

func typeNameFieldName(typeName, tagName []byte) []byte { /* typeName__fieldName*/
	key := []byte(typeName) // TODO copy ?
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
		value, err := rlp.EncodeToBytes(v.Interface())
		if err != nil {
			return nil
		}
		cloneKey = append(cloneKey, value...)
	}
	return cloneKey
}
