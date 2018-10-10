package entity

import (
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto/rlp"
	"fmt"
)

type ReversibleBlockObject struct {
	ID          uint64          `storm:"id,increment"`
	BlockNum    uint32          `storm:"unique,blockNum"`
	PackedBlock common.HexBytes //TODO c++ shared_string
}

func (rbo *ReversibleBlockObject) SetBlock(b *types.SignedBlock) {
	bo,err:= rlp.EncodeToBytes(b)
	if err!=nil{
		fmt.Println("ReversibleBlockObject SetBlock is error:",err)
	}
	rbo.PackedBlock = bo
}

func (rbo *ReversibleBlockObject) GetBlock() *types.SignedBlock{
	result := types.SignedBlock{}
	rlp.DecodeBytes(rbo.PackedBlock,result)
	return &result
}
