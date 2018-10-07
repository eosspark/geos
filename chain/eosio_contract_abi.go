package chain

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
)

func CommonTypeDefs() []types.TypeDef {
	ts := []types.TypeDef{}
	ts = append(ts, types.TypeDef{"AccountName", "name"})
	ts = append(ts, types.TypeDef{"PermissionName", "name"})
	ts = append(ts, types.TypeDef{"ActionName", "name"})
	ts = append(ts, types.TypeDef{"TableName", "name"})
	ts = append(ts, types.TypeDef{"TransactionIdType", "checksum256"})
	ts = append(ts, types.TypeDef{"BlockIdType", "checksum256"})
	ts = append(ts, types.TypeDef{"WeightType", "uint16"})
	return ts
}

func EosioContractAbi(eosioSystemAbi types.AbiDef) types.AbiDef {
	eosAbi := eosioSystemAbi
	if len(eosAbi.Version) == 0 {
		eosAbi.Version = "eosio::abi/1.0"
	}
	eosAbi.Types = CommonTypeDefs()
	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"permission_level", "",
		[]types.FieldDef{
			types.FieldDef{"actor", "account_name"},
			types.FieldDef{"permission", "permission_name"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"action", "",
		[]types.FieldDef{
			types.FieldDef{"account", "account_name"},
			types.FieldDef{"name", "action_name"},
			types.FieldDef{"authorization", "permission_level[]"},
			types.FieldDef{"data", "bytes"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"extension", "",
		[]types.FieldDef{
			types.FieldDef{"type", "uint16"},
			types.FieldDef{"data", "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"transaction_header", "",
		[]types.FieldDef{
			types.FieldDef{"expiration", "time_point_sec"},
			types.FieldDef{"ref_block_num", "uint16"},
			types.FieldDef{"ref_block_prefix", "uint32"},
			types.FieldDef{"max_net_usage_words", "varuint32"},
			types.FieldDef{"max_cpu_usage_ms", "uint8"},
			types.FieldDef{"delay_sec", "varuint32"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"transaction", "transaction_header",
		[]types.FieldDef{
			types.FieldDef{"context_free_actions", "action[]"},
			types.FieldDef{"actions", "action[]"},
			types.FieldDef{"transaction_extensions", "extension[]"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"producer_key", "",
		[]types.FieldDef{
			types.FieldDef{"producer_name", "account_name"},
			types.FieldDef{"block_signing_key", "public_key"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"producer_schedule", "",
		[]types.FieldDef{
			types.FieldDef{"version", "uint32"},
			types.FieldDef{"producers", "producer_key[]"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"block_header", "",
		[]types.FieldDef{
			types.FieldDef{"timestamp", "uint32"},
			types.FieldDef{"producer", "account_name"},
			types.FieldDef{"confirmed", "uint16"},
			types.FieldDef{"previous", "block_id_type"},
			types.FieldDef{"transaction_mroot", "checksum256"},
			types.FieldDef{"action_mroot", "checksum256"},
			types.FieldDef{"schedule_version", "uint32"},
			types.FieldDef{"new_producers", "producer_schedule?"}, //TODO c++ producer_schedule?
			types.FieldDef{"header_extensions", "extension[]"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"key_weight", "",
		[]types.FieldDef{
			types.FieldDef{"key", "public_key"},
			types.FieldDef{"weight", "weight_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"permission_level_weight", "",
		[]types.FieldDef{
			types.FieldDef{"permission", "permission_level"},
			types.FieldDef{"weight", "weight_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"wait_weight", "",
		[]types.FieldDef{
			types.FieldDef{"wait_sec", "uint32"},
			types.FieldDef{"weight", "weight_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"authority", "",
		[]types.FieldDef{
			types.FieldDef{"threshold", "uint32"},
			types.FieldDef{"keys", "weight_type[]"},
			types.FieldDef{"accounts", "permission_level_weight[]"},
			types.FieldDef{"waits", "weight_type[]"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"newaccount", "",
		[]types.FieldDef{
			types.FieldDef{"creator", "account_name"},
			types.FieldDef{"name", "account_name"},
			types.FieldDef{"owner", "authority"},
			types.FieldDef{"active", "authority"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"setcode", "",
		[]types.FieldDef{
			types.FieldDef{"account", "account_name"},
			types.FieldDef{"vmtype", "uint8"},
			types.FieldDef{"vmversion", "uint8"},
			types.FieldDef{"code", "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"setabi", "",
		[]types.FieldDef{
			types.FieldDef{"account", "account_name"},
			types.FieldDef{"abi", "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"updateauth", "",
		[]types.FieldDef{
			types.FieldDef{"account", "account_name"},
			types.FieldDef{"permission", "permission_name"},
			types.FieldDef{"parent", "permission_name"},
			types.FieldDef{"auth", "authority"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"deleteauth", "",
		[]types.FieldDef{
			types.FieldDef{"account", "account_name"},
			types.FieldDef{"permission", "permission_name"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"linkauth", "",
		[]types.FieldDef{
			types.FieldDef{"account", "account_name"},
			types.FieldDef{"code", "account_name"},
			types.FieldDef{"type", "account_name"},
			types.FieldDef{"requirement", "permission_name"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"unlinkauth", "",
		[]types.FieldDef{
			types.FieldDef{"account", "account_name"},
			types.FieldDef{"code", "permission_name"},
			types.FieldDef{"type", "permission_name"},
		},
	})
	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"canceldelay", "",
		[]types.FieldDef{
			types.FieldDef{"canceling_auth", "permission_level"},
			types.FieldDef{"trx_id", "transaction_id_type"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"onerror", "",
		[]types.FieldDef{
			types.FieldDef{"sender_id", "uint128"},
			types.FieldDef{"send_trx", "bytes"},
		},
	})

	eosAbi.Structs = append(eosAbi.Structs, types.StructDef{"onblock", "",
		[]types.FieldDef{
			types.FieldDef{"header", "block_header"},
		},
	})

	eosAbi.Actions = append(eosAbi.Actions, types.ActionDef{common.ActionName(common.S("newaccount")), "newaccount", ""})
	eosAbi.Actions = append(eosAbi.Actions, types.ActionDef{common.ActionName(common.S("setcode")), "setcode", ""})
	eosAbi.Actions = append(eosAbi.Actions, types.ActionDef{common.ActionName(common.S("setabi")), "setabi", ""})
	eosAbi.Actions = append(eosAbi.Actions, types.ActionDef{common.ActionName(common.S("updateauth")), "updateauth", ""})
	eosAbi.Actions = append(eosAbi.Actions, types.ActionDef{common.ActionName(common.S("deleteauth")), "deleteauth", ""})

	eosAbi.Actions = append(eosAbi.Actions, types.ActionDef{common.ActionName(common.S("linkauth")), "linkauth", ""})
	eosAbi.Actions = append(eosAbi.Actions, types.ActionDef{common.ActionName(common.S("unlinkauth")), "unlinkauth", ""})
	eosAbi.Actions = append(eosAbi.Actions, types.ActionDef{common.ActionName(common.S("canceldelay")), "canceldelay", ""})
	eosAbi.Actions = append(eosAbi.Actions, types.ActionDef{common.ActionName(common.S("onerror")), "onerror", ""})
	eosAbi.Actions = append(eosAbi.Actions, types.ActionDef{common.ActionName(common.S("onblock")), "onblock", ""})
	return eosAbi
}
