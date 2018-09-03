package types

import (
	"crypto/sha256"
	"encoding/binary"
	// "encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
	"github.com/eosspark/eos-go/log"
	"github.com/eosspark/eos-go/rlp"
	"sort"
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

type SharedProducerScheduleType struct {
	Version   uint32
	Producers []ProducerKey
}

type BlockHeader struct {
	Timestamp        common.BlockTimeStamp       `json:"timestamp"`
	Producer         common.AccountName          `json:"producer"`
	Confirmed        uint16                      `json:"confirmed"`
	Previous         common.BlockIDType          `json:"previous"`
	TransactionMRoot common.CheckSum256Type      `json:"transaction_mroot"`
	ActionMRoot      common.CheckSum256Type      `json:"action_mroot"`
	ScheduleVersion  uint32                      `json:"schedule_version"`
	NewProducers     *SharedProducerScheduleType `json:"new_producers" eos:"optional"`
	HeaderExtensions []*Extension                `json:"header_extensions"`
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
	return id, nil
}

/*type OptionalProducerSchedule struct {
	ProducerScheduleType
}*/

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

type HeaderConfirmation struct {
	BlockId           common.BlockIDType `json:"block_id"`
	Producer          common.AccountName `json:"producer"`
	ProducerSignature ecc.PublicKey      `json:"producers_signature"`
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
	ProducerToLastProduced           map[common.AccountName]uint32
	ProducerToLastImpliedIrb         map[common.AccountName]uint32
	BlockSigningKey                  ecc.PublicKey
	ConfirmCount                     []uint8              `json:"confirm_count"`
	Confirmations                    []HeaderConfirmation `json:"confirmations"`
}

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

func (bs *BlockHeaderState) GetScheduledProducer(t common.BlockTimeStamp) ProducerKey {
	index := uint32(t) % uint32(len(bs.ActiveSchedule.Producers)*12)
	index /= 12
	return bs.ActiveSchedule.Producers[index]
}

func (bs *BlockHeaderState) CalcDposLastIrreversible() uint32 {
	blockNums := make([]int, 0, len(bs.ProducerToLastImpliedIrb))
	for _, value := range bs.ProducerToLastImpliedIrb {
		blockNums = append(blockNums, int(value))
	}
	/// 2/3 must be greater, so if I go 1/3 into the list sorted from low to high, then 2/3 are greater

	if len(blockNums) == 0 {
		return 0
	}
	/// TODO: update to nth_element
	sort.Ints(blockNums)
	return uint32(blockNums[(len(blockNums)-1)/3])
}

func (bs *BlockHeaderState) GenerateNext(when *common.BlockTimeStamp) *BlockHeaderState {
	result := new(BlockHeaderState)

	if when != nil {
		if *when <= bs.Header.Timestamp {
			panic("next block must be in the future") //block_validate_exception
		}
	} else {
		*when = bs.Header.Timestamp
		*when++
	}

	result.Header.Timestamp = *when
	result.Header.Previous = bs.ID
	result.Header.ScheduleVersion = bs.ActiveSchedule.Version

	proKey := bs.GetScheduledProducer(*when)
	result.BlockSigningKey = proKey.BlockSigningKey
	result.Header.Producer = proKey.AccountName

	result.PendingScheduleLibNum = bs.PendingScheduleLibNum
	result.PendingScheduleHash = bs.PendingScheduleHash
	result.BlockNum = bs.BlockNum + 1
	result.ProducerToLastProduced = bs.ProducerToLastProduced
	result.ProducerToLastImpliedIrb = bs.ProducerToLastImpliedIrb
	result.ProducerToLastProduced = make(map[common.AccountName]uint32)
	result.ProducerToLastProduced[proKey.AccountName] = result.BlockNum
	result.BlockrootMerkle = bs.BlockrootMerkle
	result.BlockrootMerkle.Append(bs.ID)

	result.ActiveSchedule = bs.ActiveSchedule
	result.PendingSchedule = bs.PendingSchedule
	result.DposProposedIrreversibleBlocknum = bs.DposProposedIrreversibleBlocknum
	result.BftIrreversibleBlocknum = bs.BftIrreversibleBlocknum

	result.ProducerToLastImpliedIrb = make(map[common.AccountName]uint32)
	result.ProducerToLastImpliedIrb[proKey.AccountName] = result.DposProposedIrreversibleBlocknum
	result.DposIrreversibleBlocknum = result.CalcDposLastIrreversible()

	/// grow the confirmed count
	if common.DefaultConfig.MaxProducers*2/3+1 > 0xff {
		panic("8bit confirmations may not be able to hold all of the needed confirmations")
	}

	// This uses the previous block active_schedule because thats the "schedule" that signs and therefore confirms _this_ block
	numActiveProducers := len(bs.ActiveSchedule.Producers)
	requiredConfs := uint32(numActiveProducers*2/3) + 1

	if len(bs.ConfirmCount) < common.DefaultConfig.MaxTrackedDposConfirmations {
		result.ConfirmCount = make([]uint8, len(bs.ConfirmCount)+1)
		copy(result.ConfirmCount, bs.ConfirmCount)
		result.ConfirmCount[len(result.ConfirmCount)-1] = uint8(requiredConfs)
	} else {
		result.ConfirmCount = make([]uint8, len(bs.ConfirmCount))
		copy(result.ConfirmCount, bs.ConfirmCount[1:])
		result.ConfirmCount[len(result.ConfirmCount)-1] = uint8(requiredConfs)
	}

	return result
}

func (bs *BlockHeaderState) MaybePromotePending() bool {
	if len(bs.PendingSchedule.Producers) > 0 && bs.DposIrreversibleBlocknum >= bs.PendingScheduleLibNum {
		bs.ActiveSchedule = bs.PendingSchedule

		var newProducerToLastProduced map[common.AccountName]uint32
		var newProducerToLastImpliedIrb map[common.AccountName]uint32
		for _, pro := range bs.ActiveSchedule.Producers {
			existing, hasExisting := bs.ProducerToLastProduced[pro.AccountName]
			if hasExisting {
				newProducerToLastProduced[pro.AccountName] = existing
			} else {
				newProducerToLastProduced[pro.AccountName] = bs.DposIrreversibleBlocknum
			}

			existingIrb, hasExistingIrb := bs.ProducerToLastImpliedIrb[pro.AccountName]
			if hasExistingIrb {
				newProducerToLastImpliedIrb[pro.AccountName] = existingIrb
			} else {
				newProducerToLastImpliedIrb[pro.AccountName] = bs.DposIrreversibleBlocknum
			}
		}

		bs.ProducerToLastProduced = newProducerToLastProduced
		bs.ProducerToLastImpliedIrb = newProducerToLastImpliedIrb
		bs.ProducerToLastProduced[bs.Header.Producer] = bs.BlockNum

		return true
	}
	return false
}

func (bs *BlockHeaderState) SetConfirmed(numPrevBlocks uint16) {
	bs.Header.Confirmed = numPrevBlocks

	i := len(bs.ConfirmCount) - 1
	blocksToConfirm := numPrevBlocks + 1 /// confirm the head block too
	for i >= 0 && blocksToConfirm > 0 {
		bs.ConfirmCount[i]--
		if bs.ConfirmCount[i] == 0 {
			blockNumFori := bs.BlockNum - uint32(len(bs.ConfirmCount)-1-i)
			bs.DposProposedIrreversibleBlocknum = blockNumFori

			if i == len(bs.ConfirmCount)-1 {
				bs.ConfirmCount = make([]uint8, 0)
			} else {
				bs.ConfirmCount = bs.ConfirmCount[i+1:]
			}

			return
		}
		i--
		blocksToConfirm--
	}

}

func (bs *BlockHeaderState) SigDigest() []byte {
	result := make([]byte, 32)
	headerBmroot := common.Hash([2]interface{}{bs.Header, bs.BlockrootMerkle.GetRoot()})
	digest := common.Hash([2]interface{}{headerBmroot, bs.PendingScheduleHash})

	binary.LittleEndian.PutUint64(result[0:8], digest[0])
	binary.LittleEndian.PutUint64(result[8:16], digest[1])
	binary.LittleEndian.PutUint64(result[16:24], digest[2])
	binary.LittleEndian.PutUint64(result[24:32], digest[3])

	return result
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

func (spst *SharedProducerScheduleType) clear() {
	spst.Version = 0
	spst.Producers = []ProducerKey{}
}

func (spst *SharedProducerScheduleType) SharedroducerScheduleType(a ProducerScheduleType) *ProducerScheduleType {
	var result ProducerScheduleType = ProducerScheduleType{}
	spst.Version = a.Version
	spst.Producers = nil
	//spst.Producers = a.Producers
	for i := 0; i < len(a.Producers); i++ {
		spst.Producers[i] = a.Producers[i]
	}
	return &result
}

func (spst *SharedProducerScheduleType) producerScheduleType() *ProducerScheduleType {
	var result ProducerScheduleType = ProducerScheduleType{}
	result.Version = spst.Version
	if len(result.Producers) == 0 {
		result.Producers = spst.Producers
	} else {
		var step int = len(result.Producers)
		for _, p := range spst.Producers {
			result.Producers[step] = p
			step++
		}
	}
	return &result
}

func (bhs *BlockState) SetNewProducers(pending SharedProducerScheduleType) {
	if pending.Version == bhs.ActiveSchedule.Version+1 {
		log.Error("wrong producer schedule version specified")
		return
	}
	/*	bhs.Header.NewProducers = pending.Producers
		bhs.PendingScheduleHash = bhs.Header.NewProducers
		bhs.PendingSchedule = bhs.Header.NewProducers*/
	bhs.PendingScheduleLibNum = bhs.BlockNum

}
func (bs *BlockHeaderState) AddConfirmation(conf HeaderConfirmation) {
	//TODO
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
