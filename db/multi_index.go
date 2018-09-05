package eosiodb

type multi_index struct {
	Code      uint64
	Scope     uint64
	TableName uint64
}

func (index *multi_index) Get_code() uint64 {
	return index.Code
}

func (index *multi_index) Get_scope() uint64 {
	return index.Scope
}

func (index *multi_index) Get_table() uint64 {
	return index.TableName
}
