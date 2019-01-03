package chain

import (
	abi "github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/common"
)

func CommonTypeDefs() []abi.TypeDef {
	ts := make([]abi.TypeDef, 0)
	ts = append(ts, abi.TypeDef{"account_name", "name"})
	ts = append(ts, abi.TypeDef{"permission_name", "name"})
	ts = append(ts, abi.TypeDef{"action_name", "name"})
	ts = append(ts, abi.TypeDef{"table_name", "name"})
	ts = append(ts, abi.TypeDef{"transaction_id_type", "checksum256"})
	ts = append(ts, abi.TypeDef{"block_id_type", "checksum256"})
	ts = append(ts, abi.TypeDef{"weight_type", "uint16"})
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
			{"actor", "account_name"},
			{"permission", "permission_name"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "action",
		Base: "",
		Fields: []abi.FieldDef{
			{"account", "account_name"},
			{"name", "action_name"},
			{"authorization", "permission_level[]"},
			{"data", "bytes"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "extension",
		Base: "",
		Fields: []abi.FieldDef{
			{"type", "uint16"},
			{"data", "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "transaction_header",
		Base: "",
		Fields: []abi.FieldDef{
			{"expiration", "time_point_sec"},
			{"ref_block_num", "uint16"},
			{"ref_block_prefix", "uint32"},
			{"max_net_usage_words", "varuint32"},
			{"max_cpu_usage_ms", "uint8"},
			{"delay_sec", "varuint32"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "transaction",
		Base: "transaction_header",
		Fields: []abi.FieldDef{
			{"context_free_actions", "action[]"},
			{"actions", "action[]"},
			{"transaction_extensions", "extension[]"},
		},
	})

	//block_header
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "producer_key",
		Base: "",
		Fields: []abi.FieldDef{
			{"producer_name", "account_name"},
			{"block_signing_key", "public_key"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "producer_schedule",
		Base: "",
		Fields: []abi.FieldDef{
			{"version", "uint32"},
			{"producers", "producer_key[]"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "block_header",
		Base: "",
		Fields: []abi.FieldDef{
			{"timestamp", "uint32"},
			{"producer", "account_name"},
			{"confirmed", "uint16"},
			{"previous", "block_id_type"},
			{"transaction_mroot", "checksum256"},
			{"action_mroot", "checksum256"},
			{"schedule_version", "uint32"},
			{"new_producers", "producer_schedule?"},
			{"header_extensions", "extension[]"},
		},
	})

	//authority
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "key_weight",
		Base: "",
		Fields: []abi.FieldDef{
			{"key", "public_key"},
			{"weight", "weight_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "permission_level_weight",
		Base: "",
		Fields: []abi.FieldDef{
			{"permission", "permission_level"},
			{"weight", "weight_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "wait_weight",
		Base: "",
		Fields: []abi.FieldDef{
			{"wait_sec", "uint32"},
			{"weight", "weight_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "authority",
		Base: "",
		Fields: []abi.FieldDef{
			{"threshold", "uint32"},
			{"keys", "key_weight[]"},
			{"accounts", "permission_level_weight[]"},
			{"waits", "weight_type[]"},
		},
	})

	//action payloads
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "newaccount",
		Base: "",
		Fields: []abi.FieldDef{
			{"creator", "account_name"},
			{"name", "account_name"},
			{"owner", "authority"},
			{"active", "authority"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "setcode",
		Base: "",
		Fields: []abi.FieldDef{
			{"account", "account_name"},
			{"vmtype", "uint8"},
			{"vmversion", "uint8"},
			{"code", "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "setabi",
		Base: "",
		Fields: []abi.FieldDef{
			{"account", "account_name"},
			{"abi", "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "updateauth",
		Base: "",
		Fields: []abi.FieldDef{
			{"account", "account_name"},
			{"permission", "permission_name"},
			{"parent", "permission_name"},
			{"auth", "authority"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "deleteauth",
		Base: "",
		Fields: []abi.FieldDef{
			{"account", "account_name"},
			{"permission", "permission_name"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "linkauth",
		Base: "",
		Fields: []abi.FieldDef{
			{"account", "account_name"},
			{"code", "account_name"},
			{"type", "account_name"},
			{"requirement", "permission_name"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{Name: "unlinkauth", Base: "",
		Fields: []abi.FieldDef{
			{"account", "account_name"},
			{"code", "permission_name"},
			{"type", "permission_name"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "canceldelay",
		Base: "",
		Fields: []abi.FieldDef{
			{"canceling_auth", "permission_level"},
			{"trx_id", "transaction_id_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "onerror",
		Base: "",
		Fields: []abi.FieldDef{
			{"sender_id", "uint128"},
			{"send_trx", "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{
		Name: "onblock",
		Base: "",
		Fields: []abi.FieldDef{
			{"header", "block_header"},
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
