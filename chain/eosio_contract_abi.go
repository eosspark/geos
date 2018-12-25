package chain

import (
	abi "github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/common"
)

func CommonTypeDefs() []abi.TypeDef {
	ts := []abi.TypeDef{}
	ts = append(ts, abi.TypeDef{"AccountName", "name"})
	ts = append(ts, abi.TypeDef{"PermissionName", "name"})
	ts = append(ts, abi.TypeDef{"ActionName", "name"})
	ts = append(ts, abi.TypeDef{"TableName", "name"})
	ts = append(ts, abi.TypeDef{"TransactionIdType", "checksum256"})
	ts = append(ts, abi.TypeDef{"BlockIdType", "checksum256"})
	ts = append(ts, abi.TypeDef{"WeightType", "uint16"})
	return ts
}

func EosioContractAbi(eosioSystemAbi abi.AbiDef) abi.AbiDef {
	eosAbi := eosioSystemAbi
	if len(eosAbi.Version) == 0 {
		eosAbi.Version = "eosio::abi/1.0"
	}
	eosAbi.Types = CommonTypeDefs()
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"permission_level", "",
		[]abi.FieldDef{
			abi.FieldDef{"actor", "account_name"},
			abi.FieldDef{"permission", "permission_name"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"action", "",
		[]abi.FieldDef{
			abi.FieldDef{"account", "account_name"},
			abi.FieldDef{"name", "action_name"},
			abi.FieldDef{"authorization", "permission_level[]"},
			abi.FieldDef{"data", "bytes"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"extension", "",
		[]abi.FieldDef{
			abi.FieldDef{"type", "uint16"},
			abi.FieldDef{"data", "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"transaction_header", "",
		[]abi.FieldDef{
			abi.FieldDef{"expiration", "time_point_sec"},
			abi.FieldDef{"ref_block_num", "uint16"},
			abi.FieldDef{"ref_block_prefix", "uint32"},
			abi.FieldDef{"max_net_usage_words", "varuint32"},
			abi.FieldDef{"max_cpu_usage_ms", "uint8"},
			abi.FieldDef{"delay_sec", "varuint32"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"transaction", "transaction_header",
		[]abi.FieldDef{
			abi.FieldDef{"context_free_actions", "action[]"},
			abi.FieldDef{"actions", "action[]"},
			abi.FieldDef{"transaction_extensions", "extension[]"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"producer_key", "",
		[]abi.FieldDef{
			abi.FieldDef{"producer_name", "account_name"},
			abi.FieldDef{"block_signing_key", "public_key"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"producer_schedule", "",
		[]abi.FieldDef{
			abi.FieldDef{"version", "uint32"},
			abi.FieldDef{"producers", "producer_key[]"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"block_header", "",
		[]abi.FieldDef{
			abi.FieldDef{"timestamp", "uint32"},
			abi.FieldDef{"producer", "account_name"},
			abi.FieldDef{"confirmed", "uint16"},
			abi.FieldDef{"previous", "block_id_type"},
			abi.FieldDef{"transaction_mroot", "checksum256"},
			abi.FieldDef{"action_mroot", "checksum256"},
			abi.FieldDef{"schedule_version", "uint32"},
			abi.FieldDef{"new_producers", "producer_schedule?"}, //TODO c++ producer_schedule?
			abi.FieldDef{"header_extensions", "extension[]"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"key_weight", "",
		[]abi.FieldDef{
			abi.FieldDef{"key", "public_key"},
			abi.FieldDef{"weight", "weight_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"permission_level_weight", "",
		[]abi.FieldDef{
			abi.FieldDef{"permission", "permission_level"},
			abi.FieldDef{"weight", "weight_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"wait_weight", "",
		[]abi.FieldDef{
			abi.FieldDef{"wait_sec", "uint32"},
			abi.FieldDef{"weight", "weight_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"authority", "",
		[]abi.FieldDef{
			abi.FieldDef{"threshold", "uint32"},
			abi.FieldDef{"keys", "weight_type[]"},
			abi.FieldDef{"accounts", "permission_level_weight[]"},
			abi.FieldDef{"waits", "weight_type[]"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"newaccount", "",
		[]abi.FieldDef{
			abi.FieldDef{"creator", "account_name"},
			abi.FieldDef{"name", "account_name"},
			abi.FieldDef{"owner", "authority"},
			abi.FieldDef{"active", "authority"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"setcode", "",
		[]abi.FieldDef{
			abi.FieldDef{"account", "account_name"},
			abi.FieldDef{"vmtype", "uint8"},
			abi.FieldDef{"vmversion", "uint8"},
			abi.FieldDef{"code", "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"setabi", "",
		[]abi.FieldDef{
			abi.FieldDef{"account", "account_name"},
			abi.FieldDef{"abi", "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"updateauth", "",
		[]abi.FieldDef{
			abi.FieldDef{"account", "account_name"},
			abi.FieldDef{"permission", "permission_name"},
			abi.FieldDef{"parent", "permission_name"},
			abi.FieldDef{"auth", "authority"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"deleteauth", "",
		[]abi.FieldDef{
			abi.FieldDef{"account", "account_name"},
			abi.FieldDef{"permission", "permission_name"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"linkauth", "",
		[]abi.FieldDef{
			abi.FieldDef{"account", "account_name"},
			abi.FieldDef{"code", "account_name"},
			abi.FieldDef{"type", "account_name"},
			abi.FieldDef{"requirement", "permission_name"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"unlinkauth", "",
		[]abi.FieldDef{
			abi.FieldDef{"account", "account_name"},
			abi.FieldDef{"code", "permission_name"},
			abi.FieldDef{"type", "permission_name"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"canceldelay", "",
		[]abi.FieldDef{
			abi.FieldDef{"canceling_auth", "permission_level"},
			abi.FieldDef{"trx_id", "transaction_id_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"onerror", "",
		[]abi.FieldDef{
			abi.FieldDef{"sender_id", "uint128"},
			abi.FieldDef{"send_trx", "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, abi.StructDef{"onblock", "",
		[]abi.FieldDef{
			abi.FieldDef{"header", "block_header"},
		},
	})

	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{common.ActionName(common.N("newaccount")), "newaccount", ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{common.ActionName(common.N("setcode")), "setcode", ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{common.ActionName(common.N("setabi")), "setabi", ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{common.ActionName(common.N("updateauth")), "updateauth", ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{common.ActionName(common.N("deleteauth")), "deleteauth", ""})

	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{common.ActionName(common.N("linkauth")), "linkauth", ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{common.ActionName(common.N("unlinkauth")), "unlinkauth", ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{common.ActionName(common.N("canceldelay")), "canceldelay", ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{common.ActionName(common.N("onerror")), "onerror", ""})
	eosAbi.Actions = append(eosAbi.Actions, abi.ActionDef{common.ActionName(common.N("onblock")), "onblock", ""})
	return eosAbi
}
