package types

import "github.com/eosspark/eos-go/common"

type TypeName string
type FieldName string

type TypeDef struct {
	NewTypeName TypeName
	Type        TypeName
}

type FieldDef struct {
	Name FieldName `json:"name"`
	Type TypeName `json:"type"`
}

type StructDef struct {
	Name   TypeName     `json:"name"`
	Base   TypeName     `json:"base"`
	Fields []FieldDef   `json:"fields,omitempty"`
}

type ActionDef struct {
	Name              common.ActionName `json:"name"`
	Type              TypeName          `json:"type"`
	RicardianContract string            `json:"ricardian_contract"`
}

type TableDef struct {
	Name      common.TableName `json:"name"`
	IndexType TypeName         `json:"index_type"`
	KeyNames  []FieldName      `json:"key_names,omitempty"`
	KeyTypes  []TypeName       `json:"key_types,omitempty"`
	Type      TypeName           `json:"type"`
}

type ClausePair struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

type ErrorMessage struct {
	Code    uint64 `json:"error_code"`
	Message string `json:"error_msg"`
}

type AbiDef struct {
	Version          string            `json:"version"`
	Types            []TypeDef         `json:"types,omitempty"`
	Structs          []StructDef       `json:"structs,omitempty"`
	Actions          []ActionDef       `json:"actions,omitempty"`
	Tables           []TableDef        `json:"tables,omitempty"`
	RicardianClauses []ClausePair      `json:"ricardian_clauses,omitempty"`
	ErrorMessages    []ErrorMessage `json:"error_messages,omitempty"`
	Extensions       []*Extension      `json:"abi_extensions,omitempty"`
}

func CommonTypeDefs() []TypeDef {
	types := make([]TypeDef, 7)
	types[0] = TypeDef{"account_name", "name"}
	types[1] = TypeDef{"permission_name", "name"}
	types[2] = TypeDef{"action_name", "name"}
	types[3] = TypeDef{"table_name", "name"}
	types[4] = TypeDef{"transaction_id_type","checksum256"}
	types[5] = TypeDef{"block_id_type", "checksum256"}
	types[6] = TypeDef{"weight_type", "uint16"}
	return types
}