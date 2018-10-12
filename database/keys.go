package database

import (
	"github.com/eosspark/eos-go/crypto/rlp"
	"reflect"
)

// key --> typeName__id
func idKey(id ,typeName []byte) []byte {// FIXME
	key := 	append(typeName,'_')
	key = 	append(key,'_')
	key = 	append(key,id...)

	return key
}

// 	fieldValue[0]__fieldValue[1]...
func getFieldValue(key []byte, info *fieldInfo) []byte {

	for _, v := range info.fieldValue {
		// typeName__tag__fieldValue...
		key = append(key, '_')
		key = append(key, '_')
		value, err := rlp.EncodeToBytes(v.Interface())
		if err != nil {
			return nil
		}
		key = append(key, value...)
	}

	//fmt.Println("func fieldKey value is : ",string(key))
	//fmt.Println("func fieldKey value is : ",key)
	return key
}

func getFieldInfo(fieldName string,value interface{})(*fieldInfo,error){
	ref := reflect.ValueOf(value)
	if !ref.IsValid() || reflect.Indirect(ref).Kind() != reflect.Struct {
		return nil,ErrBadType
	}
	if ref.Kind() == reflect.Ptr {
		return nil, ErrStructNeeded
	}
	cfg, err := extractStruct(&ref)
	if err != nil {
		return nil, err
	}

	fields, ok := cfg.Fields[fieldName]
	if !ok {
		return nil, ErrNotFound
	}
	return fields,nil
}

// non unique fields --> find function
func nonUniqueValue(info *fieldInfo)[]byte{
	for _, v := range info.fieldValue {
		if isZero(v) {
			return nil
		}
	}
	return getNonUniqueFieldValue(info)
}

// non unique fields --> get function
func getNonUniqueFieldValue(info *fieldInfo)[]byte{
	values := []byte{}
	regexp := []byte{40,46,42,41}
	for _, v := range info.fieldValue {
		values = append(values,'_')
		values = append(values,'_')
		if isZero(v) {
			values = append(values,regexp...)
			continue
		}
		re, err := rlp.EncodeToBytes(v.Interface())
		if err != nil {
			return nil
		}
		values = append(values,re...)
	}
	return values
}

// typeName__fieldName
func typeNameFieldName(typeName,fieldName string)[]byte{
	key := []byte(typeName)
	key = append(key, '_')
	key = append(key, '_')
	key = append(key, []byte(fieldName)...)
	return key
}

// non unique fields --> regexp
func indexEnd(key []byte)[]byte{
	end := make([]byte, len(key))
	copy(end, key)
	end[len(end)-1] = end[len(end)-1] + 1
	return end
}

// remove and insert function
func doCallBack(id, typeName []byte, cfg *structInfo, callBack func(key, value []byte) error) error {
	for tag, fieldCfg := range cfg.Fields {
		// typeName__
		key := append(typeName, '_')
		key = append(key, '_')
		// typeName__tag__
		key = append(key, tag...)
		key =getFieldValue(key, fieldCfg)
		if !fieldCfg.unique {
			key = append(key, id...)
		}
		//fmt.Println("func fieldIndex value is : ",string(key))
		//fmt.Println("func fieldIndex value is : ",key)
		err := callBack(key, id)
		if err != nil {
			return err
		}
	}
	return nil
}

// modify function
func modifyField(cfg, oldCfg *structInfo, callBack func(newKey, oldKey []byte) error) error {

	id, err := numbertob(cfg.Id.Interface())
	if err != nil {
		return err
	}
	typeName := []byte(cfg.Name)

	for tag, fieldCfg := range cfg.Fields {
		// typeName__
		key := append(typeName, '_')
		key = append(key, '_')
		// typeName__tag__
		key = append(key, tag...)

		oldKey := getFieldValue(key, fieldCfg)
		newKey := getFieldValue(key, oldCfg.Fields[tag])
		if !fieldCfg.unique {
			newKey = append(newKey, id...)
			oldKey = append(oldKey, id...)
		}

		err := callBack(newKey, oldKey)
		if err != nil {
			return err
		}
	}
	return nil
}
