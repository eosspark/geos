package database

// key 	 --> typeName__fieldName__fieldValue__id
func indexKey (fieldName ,typeName ,fieldValue []byte) ([]byte,error)  {// FIXME
	key := 	append(typeName,'_')
	key = 	append(key,'_')
	key = 	append(key,fieldName...)
	key = 	append(key,'_')
	key = 	append(key,'_')
	key = 	append(key,fieldValue...)
	key = 	append(key,'_')
	key = 	append(key,'_')

	return key,nil
}

// key 	 --> typeName__fieldName__fieldValue
func uniqueKey(fieldName ,typeName,fieldValue []byte) ([]byte,error) {// FIXME
	key := 	append(typeName,'_')
	key = 	append(key,'_')
	key = 	append(key,fieldName...)
	key = 	append(key,'_')
	key = 	append(key,'_')
	key = 	append(key,fieldValue...)

	return key,nil
}

// key --> typeName__id
func idKey(id ,typeName []byte) []byte {// FIXME
	key := 	append(typeName,'_')
	key = 	append(key,'_')
	key = 	append(key,id...)

	return key
}