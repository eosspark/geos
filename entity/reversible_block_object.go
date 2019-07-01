package entity

import (
	"fmt"

	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
)

type ReversibleBlockObject struct {
	ID          common.IdType   `multiIndex:"id,increment"`
	BlockNum    uint32          `multiIndex:"byNum,orderedUnique"`
	PackedBlock common.HexBytes //TODO c++ shared_string
}

func (rbo *ReversibleBlockObject) SetBlock(b *types.SignedBlock) {
	bo, err := rlp.EncodeToBytes(b)
	if err != nil {
		fmt.Println("ReversibleBlockObject SetBlock is error:", err)
	}
	rbo.PackedBlock = bo
}

func (rbo *ReversibleBlockObject) GetBlock() *types.SignedBlock {
	result := types.SignedBlock{}
	rlp.DecodeBytes(rbo.PackedBlock, result)
	return &result
}

func (rbo ReversibleBlockObject) IsEmpty() bool {
	return rbo.ID == 0 && rbo.BlockNum == 0 && rbo.PackedBlock.Size() == 0
}
