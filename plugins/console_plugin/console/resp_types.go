package console

import (
	"encoding/json"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/abi_serializer"
)

type GetInfoResult struct {
	ServerVersion            string             `json:"server_version"` // "2cc40a4e"
	ChainID                  common.ChainIdType `json:"chain_id"`
	HeadBlockNum             uint32             `json:"head_block_num"`              // 2465669,
	LastIrreversibleBlockNum uint32             `json:"last_irreversible_block_num"` // 2465655
	LastIrreversibleBlockID  common.BlockIdType `json:"last_irreversible_block_id"`  // "00000008f98f0580d7efe7abc60abaaf8a865c9428a4267df30ff7d1937a1084"
	HeadBlockID              common.BlockIdType `json:"head_block_id"`               // "00259f856bfa142d1d60aff77e70f0c4f3eab30789e9539d2684f9f8758f1b88",
	HeadBlockTime            common.TimePoint   `json:"head_block_time"`             //  "2018-02-02T04:19:32"
	HeadBlockProducer        common.AccountName `json:"head_block_producer"`         // "inita"

	VirtualBlockCPULimit uint64 `json:"virtual_block_cpu_limit"`
	VirtualBlockNetLimit uint64 `json:"virtual_block_net_limit"`
	BlockCPULimit        uint64 `json:"block_cpu_limit"`
	BlockNetLimit        uint64 `json:"block_net_limit"`
	ServerVersionString  string `json:"server_version_string"`
}

type GetBlockResult struct {
	SignedBlock    types.SignedBlock  `json:"signed_block"`
	ID             common.BlockIdType `json:"id"`
	BlockNum       uint32             `json:"block_num"`
	RefBlockPrefix uint32             `json:"ref_block_prefix"`
}

type AccountResp struct {
	AccountName            common.AccountName        `json:"account_name"`
	HeadBlockNum           uint32                    `json:"head_block_num"`
	HeadBlockTime          types.BlockTimeStamp      `json:"head_block_time"`
	Privileged             bool                      `json:"privileged"`
	LastCodeUpdate         types.BlockTimeStamp      `json:"last_code_update"`
	Created                types.BlockTimeStamp      `json:"created"`
	CoreLiquidBalance      common.Asset              `json:"core_liquid_balance"`
	RAMQuota               int64                     `json:"ram_quota"`
	RAMUsage               int64                     `json:"ram_usage"`
	NetWeight              int64                     `json:"net_weight"`
	CPUWeight              int64                     `json:"cpu_weight"`
	Permissions            []types.Permission        `json:"permissions"`
	TotalResources         common.TotalResources     `json:"total_resources"`
	SelfDelegatedBandwidth common.DelegatedBandwidth `json:"self_delegated_bandwidth"`
	RefundRequest          common.Variant            `json:"refund_request"`
	VoterInfo              common.VoterInfo          `json:"voter_info"`
}

type GetAbiResult struct {
	AccountName common.Name           `json:"account_name"`
	Abi         abi_serializer.AbiDef `json:"abi"`
}

type GetCodeResult struct {
	AccountName common.Name           `json:"account_name"`
	Wast        string                `json:"wast"`
	Wasm        string                `json:"wasm"`
	CodeHash    crypto.Sha256         `json:"code_hash"`
	Abi         abi_serializer.AbiDef `json:"abi"`
}

type GetTableRowsResp struct {
	More bool            `json:"more"`
	Rows json.RawMessage `json:"rows"` // defer loading, as it depends on `JSON` being true/false.
}
