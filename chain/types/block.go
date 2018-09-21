package types

import (
	"encoding/json"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
)

type TransactionStatus uint8

const (
	TransactionStatusExecuted TransactionStatus = iota ///< succeed, no error handler executed
	TransactionStatusSoftFail                          ///< objectively failed (not executed), error handler executed
	TransactionStatusHardFail                          ///< objectively failed and error handler objectively failed thus no state change
	TransactionStatusDelayed                           ///< transaction delayed
	TransactionStatusUnknown  = TransactionStatus(255)
)

type BlockStatus uint8

const (
	Irreversible BlockStatus = iota ///< this block has already been applied before by this node and is considered irreversible
	Validated                       ///< this is a complete block signed by a valid producer and has been previously applied by this node and therefore validated but it is not yet irreversible
	Complete                        ///< this is a complete block signed by a valid producer but is not yet irreversible nor has it yet been applied by this node
	Incomplete                      ///< this is an incomplete block (either being produced by a producer or speculatively produced by a node)
)

func (s *TransactionStatus) UnmarshalJSON(data []byte) error {
	var decoded string
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	switch decoded {
	case "executed":
		*s = TransactionStatusExecuted
	case "soft_fail":
		*s = TransactionStatusSoftFail

	case "hard_fail":
		*s = TransactionStatusHardFail
	case "delayed":
		*s = TransactionStatusDelayed
	default:
		*s = TransactionStatusUnknown
	}
	return nil
}

func (s TransactionStatus) MarshalJSON() (data []byte, err error) {
	out := "unknown"
	switch s {
	case TransactionStatusExecuted:
		out = "executed"
	case TransactionStatusSoftFail:
		out = "soft_fail"
	case TransactionStatusHardFail:
		out = "hard_fail"
	case TransactionStatusDelayed:
		out = "delayed"
	}
	return json.Marshal(out)
}
func (s TransactionStatus) String() string {

	switch s {
	case TransactionStatusExecuted:
		return "executed"
	case TransactionStatusSoftFail:
		return "soft fail"
	case TransactionStatusHardFail:
		return "hard fail"
	case TransactionStatusDelayed:
		return "delayed"
	default:
		return "unknown"
	}

}

//type TransactionID SHA256Bytes

// type ShardLock struct {
// 	AccountName common.AccountName `json:"account_name"`
// 	ScopeName   common.ScopeName   `json:"scope_name"`
// }

// type ShardSummary struct {
// 	ReadLocks    []ShardLock          `json:"read_locks"`
// 	WriteLocks   []ShardLock          `json:"write_locks"`
// 	Transactions []TransactionReceipt `json:"transactions"`
// }

// type Cycles []ShardSummary
// type RegionSummary struct {
// 	Region        uint16   `json:"region"`
// 	CyclesSummary []Cycles `json:"cycles_summary"`
// }

type ReversibleBlockObject struct {
	ID          uint64
	BlockNum    uint32
	PackedBlock string
}

type ReversibleBlockIndex struct {
	rbObject ReversibleBlockObject
	byId     uint64
	byNum    uint32
}

type TransactionReceiptHeader struct {
	Status               TransactionStatus `json:"status"`
	CPUUsageMicroSeconds uint32            `json:"cpu_usage_us"`
	NetUsageWords        uint32            `json:"net_usage_words" eos:"vuint32"`
}

type TransactionReceipt struct {
	TransactionReceiptHeader
	Transaction TransactionWithID `json:"trx"`
}

type SignedBlock struct {
	SignedBlockHeader
	Transactions    []TransactionReceipt `json:"transactions"`
	BlockExtensions []*Extension         `json:"block_extensions"`
}

func (m *SignedBlock) String() string {
	return "SignedBlock"
}

type ProducerConfirmation struct {
	BlockID     common.BlockIDType
	BlockDigest [4]uint64
	Producer    common.AccountName
	Sig         ecc.Signature
}

type Optional struct {
	Valid bool
	Pair  map[common.ChainIDType][]ecc.PublicKey
}

type TransactionWithID struct {
	// ID     common.TransactionIDType
	Tag    uint8              `json:"-"`
	Packed *PackedTransaction `json:"packed_transaction"`
}

func (t TransactionWithID) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{
		t.Packed,
	})
}
func (t *TransactionWithID) UnmarshalJSON(data []byte) error {
	var packed PackedTransaction
	if data[0] == '{' {
		if err := json.Unmarshal(data, &packed); err != nil {
			return err
		}
		*t = TransactionWithID{
			// ID:     packed.ID(),
			Packed: &packed,
		}
		// 	else if data[0] == '"' {
		// 	var id string
		// 	err := json.Unmarshal(data, &id)
		// 	if err != nil {
		// 		return err
		// 	}

		// 	shaID, err := hex.DecodeString(id)
		// 	if err != nil {
		// 		return fmt.Errorf("decoding id in trx: %s", err)
		// 	}

		// 	*t = TransactionWithID{
		// 		ID: SHA256Bytes(shaID),
		// 	}

		// 	return nil
		// }

		return nil
	}
	return nil
}

// func (t *TransactionWithID) UnmarshalJSON(data []byte) error {
// 	var packed PackedTransaction
// 	if data[0] == '{' {
// 		if err := json.Unmarshal(data, &packed); err != nil {
// 			return err
// 		}
// 		*t = TransactionWithID{
// 			ID:     packed.ID(),
// 			Packed: &packed,
// 		}

// 		return nil
// 	} else if data[0] == '"' {
// 		var id string
// 		err := json.Unmarshal(data, &id)
// 		if err != nil {
// 			return err
// 		}

// 		shaID, err := hex.DecodeString(id)
// 		if err != nil {
// 			return fmt.Errorf("decoding id in trx: %s", err)
// 		}

// 		*t = TransactionWithID{
// 			ID: SHA256Bytes(shaID),
// 		}

// 		return nil
// 	}

// 	var in []json.RawMessage
// 	err := json.Unmarshal(data, &in)
// 	if err != nil {
// 		return err
// 	}

// 	if len(in) != 2 {
// 		return fmt.Errorf("expected two params for TransactionWithID, got %d", len(in))
// 	}

// 	typ := string(in[0])
// 	switch typ {
// 	case "0":
// 		var s string
// 		if err := json.Unmarshal(in[1], &s); err != nil {
// 			return err
// 		}

// 		*t = TransactionWithID{}
// 		if err := json.Unmarshal(in[1], &t.ID); err != nil {
// 			return err
// 		}
// 	case "1":

// 		// ignore the ID field right now..
// 		err = json.Unmarshal(in[1], &packed)
// 		if err != nil {
// 			return err
// 		}

// 		*t = TransactionWithID{
// 			ID:     packed.ID(),
// 			Packed: &packed,
// 		}
// 	default:
// 		return fmt.Errorf("unsupported multi-variant trx serialization type from C++ code into Go: %q", typ)
// 	}
// 	return nil
// }
