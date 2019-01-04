package chain_plugin

import (
	"github.com/eosspark/eos-go/chain/abi_serializer"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"unsafe"
)

type GetInfoResult struct {
	ServerVersion            string             `json:"server_version"`
	ChainID                  common.ChainIdType `json:"chain_id"`
	HeadBlockNum             uint32             `json:"head_block_num"`
	LastIrreversibleBlockNum uint32             `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID  common.BlockIdType `json:"last_irreversible_block_id"`
	HeadBlockID              common.BlockIdType `json:"head_block_id"`
	HeadBlockTime            common.TimePoint   `json:"head_block_time"`
	HeadBlockProducer        common.AccountName `json:"head_block_producer"`
	VirtualBlockCPULimit     uint64             `json:"virtual_block_cpu_limit"`
	VirtualBlockNetLimit     uint64             `json:"virtual_block_net_limit"`
	BlockCPULimit            uint64             `json:"block_cpu_limit"`
	BlockNetLimit            uint64             `json:"block_net_limit"`
	ServerVersionString      string             `json:"server_version_string"`
}

type GetBlockParams struct {
	BlockNumOrID string `json:"block_num_or_id"`
}
type GetBlockResult struct {
	SignedBlock    *types.SignedBlock `json:"signed_block"`
	ID             common.BlockIdType `json:"id"`
	BlockNum       uint32             `json:"block_num"`
	RefBlockPrefix uint32             `json:"ref_block_prefix"`
}

type GetBlockHeaderStateParams struct {
	BlockNumOrID string `json:"block_num_or_id"`
}
type GetBlockHeaderStateResult = types.BlockHeaderState

type Permission struct {
	PermName     common.Name
	Parent       common.Name
	RequiredAuth types.Authority
}

type GetAccountParams struct {
	AccountName        common.AccountName `json:"account_name"`
	ExpectedCoreSymbol *common.Symbol     `json:"expected_core_symbol"`
}
type GetAccountResult struct {
	AccountName            common.AccountName         `json:"account_name"`
	HeadBlockNum           uint32                     `json:"head_block_num"`
	HeadBlockTime          common.TimePoint           `json:"head_block_time"`
	Privileged             bool                       `json:"privileged"`
	LastCodeUpdate         common.TimePoint           `json:"last_code_update"`
	Created                common.TimePoint           `json:"created"`
	CoreLiquidBalance      common.Asset               `json:"core_liquid_balance"`
	RAMQuota               int64                      `json:"ram_quota"`
	NetWeight              int64                      `json:"net_weight"`
	CPUWeight              int64                      `json:"cpu_weight"`
	NetLimit               types.AccountResourceLimit `json:"net_limit"`
	CpuLimit               types.AccountResourceLimit `json:"cpu_limit"`
	RAMUsage               int64                      `json:"ram_usage"`
	Permissions            []Permission               `json:"permissions"`
	TotalResources         common.Variant             `json:"total_resources"`
	SelfDelegatedBandwidth common.Variant             `json:"self_delegated_bandwidth"`
	RefundRequest          common.Variant             `json:"refund_request"`
	VoterInfo              common.Variant             `json:"voter_info"`
}

type GetAbiParams struct {
	AccountName common.Name `json:"account_name"`
}
type GetAbiResult struct {
	AccountName common.Name           `json:"account_name"`
	Abi         abi_serializer.AbiDef `json:"abi"`
}

type GetCodeParams struct {
	AccountName common.Name `json:"account_name"`
	CodeAsWasm  bool        `json:"code_as_wasm"`
}
type GetCodeResult struct {
	AccountName common.Name           `json:"account_name"`
	Wast        string                `json:"wast"`
	Wasm        string                `json:"wasm"`
	CodeHash    crypto.Sha256         `json:"code_hash"`
	Abi         abi_serializer.AbiDef `json:"abi"`
}

type GetCodeHashParams struct {
	AccountName common.AccountName `json:"account_name"`
}
type GetCodeHashResults struct {
	AccountName common.Name   `json:"account_name"`
	CodeHash    crypto.Sha256 `json:"code_hash"`
}

type GetRawCodeAndAbiParams struct {
	AccountName common.Name `json:"account_name"`
}
type GetRawCodeAndAbiResults struct {
	AccountName common.Name `json:"account_name"`
	Wasm        Blob        `json:"wasm"` //chain::blob
	Abi         Blob        `json:"abi"`  //chain::blob
}

type Blob struct {
	Data []byte `json:"data"`
}

type GetRawAbiParams struct {
	AccountName common.Name   `json:"account_name"`
	AbiHash     crypto.Sha256 `json:"abi_hash"` //optional
}
type GetRawAbiResult struct {
	AccountName common.Name   `json:"account_name"`
	CodeHash    crypto.Sha256 `json:"code_hash"`
	AbiHash     crypto.Sha256 `json:"abi_hash"`
	Abi         *Blob         `json:"abi"` //TODO C++ optional<chain::blob>
}

type GetRequiredKeysParams struct {
	Transaction   common.Variant  `json:"transaction"`
	AvailableKeys []ecc.PublicKey `json:"available_keys"`
}
type GetRequiredKeysResult struct {
	RequiredKeys []ecc.PublicKey `json:"required_keys"`
}

type GetCurrencyBalanceParams struct {
	Code    common.Name `json:"code"`
	Account common.Name `json:"account"`
	Symbol  string      `json:"symbol"`
}
type GetCurrencyBalanceResult = []common.Asset

type GetTableRowsParams struct {
	JSON          bool        `json:"json"`
	Code          common.Name `json:"code"`
	Scope         string      `json:"scope"`
	Table         common.Name `json:"table"`
	TableKey      string      `json:"table_key"`
	LowerBound    string      `json:"lower_bound"`
	UpperBound    string      `json:"upper_bound"`
	Limit         uint32      `json:"limit,omitempty"` // defaults to 10 => chain_plugin.hpp:struct get_table_rows_params
	KeyType       string      `json:"key_type"`        // type of key specified by index_position
	IndexPosition string      `json:"index_position"`  // 1 - primary (first), 2 - secondary index (in order defined by multi_index), 3 - third index, etc
	EncodeType    string      `json:"encode_type"`     //dec, hex , default=dec
}
type GetTableRowsResult struct {
	Rows []common.Variants `json:"rows"` // true if last element in data is not the end and sizeof data() < limit
	More bool              `json:"more"` // one row per item, either encoded as hex String or JSON object
}

type GetTableByScopeParams struct {
	Code       common.Name `json:"code"`        // mandatory
	Table      common.Name `json:"table"`       // optional, act as filter =0
	LowerBound string      `json:"lower_bound"` // lower bound of scope, optional
	UpperBound string      `json:"upper_bound"` // upper bound of scope, optional
	Limit      uint32      `json:"limit"`       //=10
}
type GetTableByScopeResultRow struct {
	Code  common.Name `json:"code"`
	Scope common.Name `json:"scope"`
	Table common.Name `json:"table"`
	Payer common.Name `json:"payer"`
	Count uint32      `json:"count"`
}

type GetCurrencyStatsParams struct {
	Code   common.Name `json:"code"`
	Symbol string      `json:"symbol"`
}
type GetCurrencyStatsResult struct {
	Supply    common.Asset       `json:"supply"`
	MaxSupply common.Asset       `json:"max_supply"`
	Issuer    common.AccountName `json:"issuer"`
}

const SizeofGetCurrencyStatsResult = int(unsafe.Sizeof(GetCurrencyStatsResult{}))

type GetProducerScheduleParams struct {
}
type GetProducerScheduleResult struct {
	Active   common.Variants `json:"active"`
	Pending  common.Variants `json:"pending"`
	Proposed common.Variants `json:"proposed"`
}

type GetProducersParams struct {
	Json       bool   `json:"json"` //defaults false
	LowerBound string `json:"lower_bound"`
	Limit      uint32 `json:"limit"` //defaults 50
}
type GetProducersResult struct {
	Rows                    []common.Variant `json:"rows"`                       //one row per item, either encoded as hex string or JSON object
	TotalProducerVoteWeight float64          `json:"total_producer_vote_weight"` //TODO C++ double
	More                    string           `json:"more"`                       //fill lower_bound with this value to fetch more rows
}

type GetScheduledTransactionParams struct {
	Json       bool   `json:"json"`
	LowerBound string `json:"lower_bound"` //timestamp OR transaction ID
	Limit      uint32 `json:"limit"`       //default 50
}
type GetScheduleTransactionReuslt struct {
	Transaction common.Variants `json:"transaction"`
	More        string          `json:"more"` //fill lower_bound with this to fetch next set of transactions
}

type AbiJsonToBinParams struct {
	Code   common.Name     `json:"code"`
	Action common.Name     `json:"action"`
	Args   common.Variants `json:"args"`
}
type AbiJsonToBinReuslt struct {
	Binargs []byte `json:"binargs"`
}

type AbiBinToJsonParams struct {
	Code    common.Name `json:"code"`
	Action  common.Name `json:"action"`
	Binargs []byte      `json:"binargs"`
}
type AbiBinToJsonResult struct {
	Args common.Variant `json:"args"`
}

// func (resp *GetTableRowsResp) JSONToStructs(v interface{}) error {
// 	return json.Unmarshal(resp.Rows, v)
// }

// func (resp *GetTableRowsResp) BinaryToStructs(v interface{}) error {
// 	var rows []string

// 	err := json.Unmarshal(resp.Rows, &rows)
// 	if err != nil {
// 		return err
// 	}

// 	outSlice := reflect.ValueOf(v).Elem()
// 	structType := reflect.TypeOf(v).Elem().Elem()

// 	for _, row := range rows {
// 		bin, err := hex.DecodeString(row)
// 		if err != nil {
// 			return err
// 		}

// 		// access the type of the `Slice`, create a bunch of them..
// 		newStruct := reflect.New(structType)

// 		decoder := NewDecoder(bin)
// 		if err := decoder.Decode(newStruct.Interface()); err != nil {
// 			return err
// 		}

// 		outSlice = reflect.Append(outSlice, reflect.Indirect(newStruct))
// 	}

// 	reflect.ValueOf(v).Elem().Set(outSlice)

// 	return nil
// }

//type Currency struct {
//	Precision uint8
//	Name      common.CurrencyName
//}
//
//type GetRequiredKeysResp struct {
//	RequiredKeys treeset.Set
//	//RequiredKeys []ecc.PublicKey `json:"required_keys"`
//}
//
//// PushTransactionFullResp unwraps the responses from a successful `push_transaction`.
//// FIXME: REVIEW the actual output, things have moved here.
//type PushTransactionFullResp struct {
//	StatusCode    string
//	TransactionID string               `json:"transaction_id"`
//	Processed     TransactionProcessed `json:"processed"` // WARN: is an `fc::variant` in server..
//}
//
//type TransactionProcessed struct {
//	Status               string                   `json:"status"`
//	ID                   common.TransactionIdType `json:"id"`
//	ActionTraces         []Trace                  `json:"action_traces"`
//	DeferredTransactions []string                 `json:"deferred_transactions"` // that's not right... dig to find what's there..
//}
//
//type Trace struct {
//	Receiver common.AccountName `json:"receiver"`
//	// Action     Action       `json:"act"` // FIXME: how do we unpack that ? what's on the other side anyway?
//	Console    string       `json:"console"`
//	DataAccess []DataAccess `json:"data_access"`
//}
//
//type DataAccess struct {
//	Type     string             `json:"type"` // "write", "read"?
//	Code     common.AccountName `json:"code"`
//	Scope    common.AccountName `json:"scope"`
//	Sequence int                `json:"sequence"`
//}
//
//type PushTransactionShortResp struct {
//	TransactionID string `json:"transaction_id"`
//	Processed     bool   `json:"processed"` // WARN: is an `fc::variant` in server..
//}

//// //
//
//type WalletSignTransactionResp struct {
//	// Ignore the rest of the transaction, so the wallet server
//	// doesn't forge some transactions on your behalf, and you send it
//	// to the network..  ... although.. it's better if you can trust
//	// your wallet !
//
//	Signatures []ecc.Signature `json:"signatures"`
//}

//type MyStruct struct {
//	Currency
//	Balance uint64
//}

//// NetConnectionResp
//type NetConnectionsResp struct {
//	Peer          string                      `json:"peer"`
//	Connecting    bool                        `json:"connecting"`
//	Syncing       bool                        `json:"syncing"`
//	LastHandshake net_plugin.HandshakeMessage `json:"last_handshake"`
//}

//type NetStatusResp struct {
//}
//
//type NetConnectResp string
//
//type NetDisconnectResp string

// type Global struct {
// 	MaxBlockNetUsage               int              `json:"max_block_net_usage"`
// 	TargetBlockNetUsagePct         int              `json:"target_block_net_usage_pct"`
// 	MaxTransactionNetUsage         int              `json:"max_transaction_net_usage"`
// 	BasePerTransactionNetUsage     int              `json:"base_per_transaction_net_usage"`
// 	NetUsageLeeway                 int              `json:"net_usage_leeway"`
// 	ContextFreeDiscountNetUsageNum int              `json:"context_free_discount_net_usage_num"`
// 	ContextFreeDiscountNetUsageDen int              `json:"context_free_discount_net_usage_den"`
// 	MaxBlockCPUUsage               int              `json:"max_block_cpu_usage"`
// 	TargetBlockCPUUsagePct         int              `json:"target_block_cpu_usage_pct"`
// 	MaxTransactionCPUUsage         int              `json:"max_transaction_cpu_usage"`
// 	MinTransactionCPUUsage         int              `json:"min_transaction_cpu_usage"`
// 	MaxTransactionLifetime         int              `json:"max_transaction_lifetime"`
// 	DeferredTrxExpirationWindow    int              `json:"deferred_trx_expiration_window"`
// 	MaxTransactionDelay            int              `json:"max_transaction_delay"`
// 	MaxInlineActionSize            int              `json:"max_inline_action_size"`
// 	MaxInlineActionDepth           int              `json:"max_inline_action_depth"`
// 	MaxAuthorityDepth              int              `json:"max_authority_depth"`
// 	MaxRAMSize                     string           `json:"max_ram_size"`
// 	TotalRAMBytesReserved          common.JSONInt64 `json:"total_ram_bytes_reserved"`
// 	TotalRAMStake                  common.JSONInt64 `json:"total_ram_stake"`
// 	LastProducerScheduleUpdate     string           `json:"last_producer_schedule_update"`
// 	LastPervoteBucketFill          int64            `json:"last_pervote_bucket_fill,string"`
// 	PervoteBucket                  int              `json:"pervote_bucket"`
// 	PerblockBucket                 int              `json:"perblock_bucket"`
// 	TotalUnpaidBlocks              int              `json:"total_unpaid_blocks"`
// 	TotalActivatedStake            float64          `json:"total_activated_stake,string"`
// 	ThreshActivatedStakeTime       int64            `json:"thresh_activated_stake_time,string"`
// 	LastProducerScheduleSize       int              `json:"last_producer_schedule_size"`
// 	TotalProducerVoteWeight        float64          `json:"total_producer_vote_weight,string"`
// 	LastNameClose                  string           `json:"last_name_close"`
// }

//type Producer struct {
//	Owner         string             `json:"owner"`
//	TotalVotes    float64            `json:"total_votes,string"`
//	ProducerKey   string             `json:"producer_key"`
//	IsActive      int                `json:"is_active"`
//	URL           string             `json:"url"`
//	UnpaidBlocks  int                `json:"unpaid_blocks"`
//	LastClaimTime common.JSONFloat64 `json:"last_claim_time"`
//	Location      int                `json:"location"`
//}
//type ProducersResp struct {
//	Producers []Producer `json:"producers"`
//}
