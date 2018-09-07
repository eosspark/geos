package types

import (
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/ecc"
)

type BlockHeader struct {
	Timestamp        common.BlockTimeStamp       `json:"timestamp"`
	Producer         common.AccountName          `json:"producer"`
	Confirmed        uint16                      `json:"confirmed"`
	Previous         common.BlockIDType          `json:"previous"`
	TransactionMRoot common.CheckSum256Type      `json:"transaction_mroot"`
	ActionMRoot      common.CheckSum256Type      `json:"action_mroot"`
	ScheduleVersion  uint32                      `json:"schedule_version"`
	NewProducers     SharedProducerScheduleType `json:"new_producers" eos:"optional"`
	HeaderExtensions []*Extension                `json:"header_extensions"`
}

func (b *BlockHeader) Digest() [4]uint64 {
	return common.Hash(*b)
}

func (b *BlockHeader) BlockNumber() uint32 {
	return common.EndianReverseU32(uint32(b.Previous[0])) + 1
}

func NumFromID(id common.BlockIDType) uint32 {
	return common.EndianReverseU32(uint32(id[0]))
}

func (b *BlockHeader) BlockID() common.BlockIDType {
	// Do not include signed_block_header attributes in id, specifically exclude producer_signature.
	result := b.Digest()
	result[0] &= 0xffffffff00000000
	result[0] += uint64(common.EndianReverseU32(b.BlockNumber())) // store the block num in the ID, 160 bits is plenty for the hash
	return result
}

type SignedBlockHeader struct {
	BlockHeader
	ProducerSignature ecc.Signature `json:"producer_signature"`
}

type HeaderConfirmation struct {
	BlockId           common.BlockIDType `json:"block_id"`
	Producer          common.AccountName `json:"producer"`
	ProducerSignature ecc.Signature      `json:"producers_signature"`
}
