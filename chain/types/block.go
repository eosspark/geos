package types

import (
	"crypto/sha256"
	"encoding/binary"
	// "encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"github.com/eosspark/eos-go/rlp"
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

type ProducerKey struct {
	AccountName     common.AccountName `json:"account_name"`
	BlockSigningKey ecc.PublicKey      `json:"block_signing_key"`
}

type ProducerScheduleType struct {
	Version   uint32        `json:"version"`
	Producers []ProducerKey `json:"producers"`
}

type BlockHeader struct {
	Timestamp        common.BlockTimeStamp     `json:"timestamp"`
	Producer         common.AccountName        `json:"producer"`
	Confirmed        uint16                    `json:"confirmed"`
	Previous         common.BlockIDType        `json:"previous"`
	TransactionMRoot common.CheckSum256Type    `json:"transaction_mroot"`
	ActionMRoot      common.CheckSum256Type    `json:"action_mroot"`
	ScheduleVersion  uint32                    `json:"schedule_version"`
	NewProducers     *OptionalProducerSchedule `json:"new_producers" eos:"optional"`
	HeaderExtensions []*Extension              `json:"header_extensions"`
}

func (b *BlockHeader) BlockNumber() uint32 {
	return common.EndianReverseU32(uint32(b.Previous[0])) + 1
}

func (b *BlockHeader) BlockID() (id common.BlockIDType, err error) {

	cereal, err := rlp.EncodeToBytes(b)
	if err != nil {
		return id, err
	}

	h := sha256.New()
	_, _ = h.Write(cereal)
	hashed := h.Sum(nil)

	binary.BigEndian.PutUint32(hashed, b.BlockNumber())
	fmt.Println(hashed)

	id[0] = binary.LittleEndian.Uint64(hashed[:8])
	id[1] = binary.LittleEndian.Uint64(hashed[8:16])
	id[2] = binary.LittleEndian.Uint64(hashed[16:24])
	id[3] = binary.LittleEndian.Uint64(hashed[24:32])
	return
}

type OptionalProducerSchedule struct {
	ProducerScheduleType
}

type SignedBlockHeader struct {
	BlockHeader
	ProducerSignature ecc.Signature `json:"producer_signature"`
}

type SignedBlock struct {
	SignedBlockHeader
	Transactions    []TransactionReceipt `json:"transactions"`
	BlockExtensions []*Extension         `json:"block_extensions"`
}

func (m *SignedBlock) String() string {
	return "SignedBlock"
}

type IncrementalMerkle struct {
	NodeCount   uint64    `json:"node_count"`
	ActiveNodes [4]uint64 `json:"active_nodes"`
}

type FlatMap struct {
	AccountName common.AccountName `json:"account_name"`
	ProducerKey uint32             `json:"producer_key"`
}

type HeaderConfirmation struct {
	BlockId           common.BlockIDType
	Producer          common.AccountName
	ProducerSignature ecc.PublicKey
}
type BlockHeaderState struct {
	ID                               common.BlockIDType `storm:"id,unique"`
	BlockNum                         uint32             `storm:"block_num,unique"`
	Header                           SignedBlockHeader
	DposProposedIrreversibleBlocknum uint32    `json:"dpos_proposed_irreversible_blocknum"`
	DposIrreversibleBlocknum         uint32    `json:"dpos_irreversible_blocknum"`
	BftIrreversibleBlocknum          uint32    `json:"bft_irreversible_blocknum"`
	PendingScheduleLibNum            uint32    `json:"pending_schedule_lib_num"`
	PendingScheduleHash              [4]uint64 `json:"pending_schedule_hash"`
	PendingSchedule                  ProducerScheduleType
	ActiveSchedule                   ProducerScheduleType
	BlockrootMerkle                  IncrementalMerkle
	ProducerToLastProduced           FlatMap
	ProducerToLastImpliedIrb         FlatMap
	BlockSigningKey                  ecc.PublicKey
	ConfirmCount                     []uint8              `json:"confirm_count"`
	Confirmations                    []HeaderConfirmation `json:"confirmations"`
}

type BlockState struct {
	BlockHeaderState
	SignedBlock    SignedBlock
	Validated      bool `json:validated`
	InCurrentChain bool `json:"in_current_chain"`
	Trxs           []TransactionMetadata
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

type Optional struct {
	Valid bool
	Pair  map[common.ChainIDType][]ecc.PublicKey
}
type TransactionMetadata struct {
	ID          common.TransactionIDType
	SignedID    common.TransactionIDType
	Trx         SignedTransaction
	PackedTrx   PackedTransaction
	SigningKeys Optional
	Accepted    bool
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
// 			// ID:     packed.ID(),
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
