package database

import (
	"github.com/eosspark/eos-go/crypto/rlp"
	"reflect"
)

func idKey(id ,typeName []byte) []byte {	/* key --> typeName__id */
	key := 	append(typeName,'_')
	key = 	append(key,'_')
	key = 	append(key,id...)

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

func nonUniqueValue(info *fieldInfo)[]byte{ 					/* non unique fields --> find function */
	for _, v := range info.fieldValue  {
		if isZero(v) && v.Kind() != reflect.Bool {
			return nil
		}
	}
	reg,_ :=  getNonUniqueFieldValue(info)
	return reg
}

// TODO The function is unchanged, need to modify the implementation
func getNonUniqueFieldValue(info *fieldInfo)([]byte,[]byte){ 	/* non unique fields --> get function */
	values := []byte{}
	prefix := []byte{}
	//regexp := []byte{40,46,42,41}
	count := 0
	for _, v := range info.fieldValue {
		values = append(values,'_')
		values = append(values,'_')
		if v.Kind() != reflect.Bool &&  isZero(v) {
			//values = append(values,regexp...)
			count++
			return prefix,prefix
			continue
		}
		re, err := rlp.EncodeToBytes(v.Interface())
		if err != nil {
			return nil,nil
		}

		if count == 0{
			prefix = append(prefix,re...)
		}
		values = append(values,re...)
	}
	if count == len(info.fieldValue){
		return nil,nil
	}
	if len(info.fieldValue) - count == 1{
		return values,prefix
	}
	return values,nil
}

func typeNameFieldName(typeName,tagName []byte)[]byte{ 		/* typeName__fieldName*/
	key := []byte(typeName)// TODO copy ?
	key = append(key, '_')
	key = append(key, '_')
	key = append(key, tagName...)
	return key
}

func getNonUniqueEnd(key []byte)[]byte{ 					/* non unique fields --> regexp*/
	end := make([]byte, len(key))
	copy(end, key)
	end[len(end)-1] = end[len(end)-1] + 1
	return end
}

func getFieldValue(key []byte, info *fieldInfo) []byte { 		/* fieldValue[0]__fieldValue[1]... */
	cloneKey :=  cloneByte(key)
	for _, v := range info.fieldValue { 						// typeName__tag__fieldValue...
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

func modifyField(cfg, oldCfg *structInfo, callBack func(newKey, oldKey []byte) error) error { 	/* modify function*/
	id,err := rlp.EncodeToBytes(cfg.Id.Interface())
	if err != nil{
		return err
	}

	typeName := []byte(cfg.Name)

	for tag, fieldCfg := range cfg.Fields {
		key := append(typeName, '_') 	// typeName__
		key = append(key, '_')
		key = append(key, tag...) 			// typeName__tag__

		newKey := getFieldValue(key, fieldCfg)
		oldKey := getFieldValue(key, oldCfg.Fields[tag])
		if !fieldCfg.unique && len(fieldCfg.fieldValue) == 1{
			oldKey = append(oldKey, id...)
			newKey = append(newKey, id...)
		}

		err := callBack(newKey, oldKey)
		if err != nil {
			return err
		}
	}
	return nil
}
