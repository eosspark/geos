package chain

import (
	abi "github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/common"
)

func CommonTypeDefs() []abi.TypeDef {
	ts := make([]abi.TypeDef, 0)
	ts = append(ts, abi.TypeDef{NewTypeName: "account_name", Type: "name"})
	ts = append(ts, abi.TypeDef{NewTypeName: "permission_name", Type: "name"})
	ts = append(ts, abi.TypeDef{NewTypeName: "action_name", Type: "name"})
	ts = append(ts, abi.TypeDef{NewTypeName: "table_name", Type: "name"})
	ts = append(ts, abi.TypeDef{NewTypeName: "transaction_id_type", Type: "checksum256"})
	ts = append(ts, abi.TypeDef{NewTypeName: "block_id_type", Type: "checksum256"})
	ts = append(ts, abi.TypeDef{NewTypeName: "weight_type", Type: "uint16"})
	return ts
}

func EosioContractAbi(eosioSystemAbi abi.AbiDef) *abi.AbiDef {
	eosAbi := eosioSystemAbi
	if len(eosAbi.Version) == 0 {
		eosAbi.Version = "eosio::abi/1.0"
	}

	eosAbi.Types = append(eosAbi.Types, CommonTypeDefs()...)

	//transaction
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "permission_level",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "actor", Type: "account_name"},
			{Name: "permission", Type: "permission_name"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "action",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "account", Type: "account_name"},
			{Name: "name", Type: "action_name"},
			{Name: "authorization", Type: "permission_level[]"},
			{Name: "data", Type: "bytes"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "extension",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "type", Type: "uint16"},
			{Name: "data", Type: "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "transaction_header",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "expiration", Type: "time_point_sec"},
			{Name: "ref_block_num", Type: "uint16"},
			{Name: "ref_block_prefix", Type: "uint32"},
			{Name: "max_net_usage_words", Type: "varuint32"},
			{Name: "max_cpu_usage_ms", Type: "uint8"},
			{Name: "delay_sec", Type: "varuint32"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "transaction",
		Base: "transaction_header",
		Fields: []abi.FieldDef{
			{Name: "context_free_actions", Type: "action[]"},
			{Name: "actions", Type: "action[]"},
			{Name: "transaction_extensions", Type: "extension[]"},
		},
	})

	//block_header
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "producer_key",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "producer_name", Type: "account_name"},
			{Name: "block_signing_key", Type: "public_key"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "producer_schedule",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "version", Type: "uint32"},
			{Name: "producers", Type: "producer_key[]"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "block_header",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "timestamp", Type: "uint32"},
			{Name: "producer", Type: "account_name"},
			{Name: "confirmed", Type: "uint16"},
			{Name: "previous", Type: "block_id_type"},
			{Name: "transaction_mroot", Type: "checksum256"},
			{Name: "action_mroot", Type: "checksum256"},
			{Name: "schedule_version", Type: "uint32"},
			{Name: "new_producers", Type: "producer_schedule?"},
			{Name: "header_extensions", Type: "extension[]"},
		},
	})

	//authority
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "key_weight",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "key", Type: "public_key"},
			{Name: "weight", Type: "weight_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "permission_level_weight",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "permission", Type: "permission_level"},
			{Name: "weight", Type: "weight_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "wait_weight",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "wait_sec", Type: "uint32"},
			{Name: "weight", Type: "weight_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "authority",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "threshold", Type: "uint32"},
			{Name: "keys", Type: "key_weight[]"},
			{Name: "accounts", Type: "permission_level_weight[]"},
			{Name: "waits", Type: "weight_type[]"},
		},
	})

	//action payloads
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "newaccount",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "creator", Type: "account_name"},
			{Name: "name", Type: "account_name"},
			{Name: "owner", Type: "authority"},
			{Name: "active", Type: "authority"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "setcode",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "account", Type: "account_name"},
			{Name: "vmtype", Type: "uint8"},
			{Name: "vmversion", Type: "uint8"},
			{Name: "code", Type: "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "setabi",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "account", Type: "account_name"},
			{Name: "abi", Type: "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "updateauth",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "account", Type: "account_name"},
			{Name: "permission", Type: "permission_name"},
			{Name: "parent", Type: "permission_name"},
			{Name: "auth", Type: "authority"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "deleteauth",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "account", Type: "account_name"},
			{Name: "permission", Type: "permission_name"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "linkauth",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "account", Type: "account_name"},
			{Name: "code", Type: "account_name"},
			{Name: "type", Type: "account_name"},
			{Name: "requirement", Type: "permission_name"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{Name: "unlinkauth", Base: "",
		Fields: []abi.FieldDef{
			{Name: "account", Type: "account_name"},
			{Name: "code", Type: "permission_name"},
			{Name: "type", Type: "permission_name"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "canceldelay",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "canceling_auth", Type: "permission_level"},
			{Name: "trx_id", Type: "transaction_id_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "onerror",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "sender_id", Type: "uint128"},
			{Name: "send_trx", Type: "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "onblock",
		Base: "",
		Fields: []abi.FieldDef{
			{Name: "header", Type: "block_header"},
		},
	})

	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{Name: common.N("newaccount"), Type: "newaccount", RicardianContract: ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{Name: common.N("setcode"), Type: "setcode", RicardianContract: ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{Name: common.N("setabi"), Type: "setabi", RicardianContract: ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{Name: common.N("updateauth"), Type: "updateauth", RicardianContract: ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{Name: common.N("deleteauth"), Type: "deleteauth", RicardianContract: ""})

	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{Name: common.N("linkauth"), Type: "linkauth", RicardianContract: ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{Name: common.N("unlinkauth"), Type: "unlinkauth", RicardianContract: ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{Name: common.N("canceldelay"), Type: "canceldelay", RicardianContract: ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{Name: common.N("onerror"), Type: "onerror", RicardianContract: ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{Name: common.N("onblock"), Type: "onblock", RicardianContract: ""})

	return &eosAbi
}
